#!/usr/bin/env python3
"""Pseudo-terminal checks for Grove's board: the terminal's modes are restored
however a session ends, the interface and the result use separate streams, a
blocked Git read is killed on the way out, and nothing on disk changes.

    python3 internal/tui/testdata/terminal.py /path/to/grove [scenario ...]

`go test ./internal/tui` builds the binary and runs this; without scenario
names every scenario runs at once in its own process. Unix only (stdlib
pty, termios, fcntl); it says nothing about Windows consoles.
"""
import fcntl, hashlib, json, os, pty, select, shutil, signal, struct, subprocess, sys, tempfile, termios, time

GROVE = os.path.abspath(sys.argv[1])
ONLY = sys.argv[2:]
TIMEOUT = 20  # failure detection only; nothing waits on a guessed delay
GIT = shutil.which("git")
# Variables through which Git takes a repository from its caller; a hook exports GIT_DIR (G-089).
GIT_LOCATION = ("GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_NAMESPACE")
ENTER, DOWN, ESC, CTRL_C = b"\r", b"\x1b[B", b"\x1b", b"\x03"
ALT_ON, ALT_OFF = b"\x1b[?1049h", b"\x1b[?1049l"

SESSIONS = []  # every grove started, killed on the way out so a failed scenario leaves none running

WORK = "---\nid: G-001\ntype: work\ntitle: {title}\nstatus: {status}\n---\nAn outcome.\n"


def git(cwd, *args):
    subprocess.run([GIT, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "-C", cwd, *args],
                   check=True, capture_output=True, env=clean_env())


def fixture(base):
    """main has G-001 proposed; a linked worktree on feature has it active."""
    root, wt = os.path.join(base, "main"), os.path.join(base, "feature-wt")
    os.makedirs(os.path.join(root, "grove", "work"))
    with open(os.path.join(root, "grove.yaml"), "w") as f:
        f.write("schema_version: 3\nrecords: grove\n")
    with open(os.path.join(root, "grove", "work", "G-001-first.md"), "w") as f:
        f.write(WORK.format(title="First on main", status="proposed"))
    git(root, "init", "-q", "-b", "main")
    git(root, "add", "-A")
    git(root, "commit", "-q", "-m", "main")
    git(root, "worktree", "add", "-q", "-b", "feature", wt)
    with open(os.path.join(wt, "grove", "work", "G-001-first.md"), "w") as f:
        f.write(WORK.format(title="First on feature", status="active"))
    git(wt, "add", "-A")
    git(wt, "commit", "-q", "-m", "feature")
    return os.path.realpath(root), os.path.realpath(wt)


def short(cwd, rev):
    """The seven-character commit ID a history row shows."""
    out = subprocess.run([GIT, "-C", cwd, "rev-parse", rev], check=True, capture_output=True, env=clean_env()).stdout
    return out.decode()[:7]


def tree(base):
    """Every file below base, Git's included, by content."""
    found = {}
    for folder, _, names in os.walk(base):
        for name in names:
            path = os.path.join(folder, name)
            if not os.path.islink(path) and os.path.isfile(path):
                with open(path, "rb") as f:
                    found[path] = hashlib.sha256(f.read()).hexdigest()
    return found


class Session:
    """grove with stdin and stderr on a pty (or as given) and stdout on a pipe."""

    def __init__(self, cwd, args=(), env=None, stdin="pty", stderr="pty", size=(30, 120)):
        self.master, self.slave = pty.openpty()
        fcntl.ioctl(self.slave, termios.TIOCSWINSZ, struct.pack("HHHH", size[0], size[1], 0, 0))
        self.before = termios.tcgetattr(self.slave)
        self.master2 = self.slave2 = None
        if stderr == "pty2":  # a second terminal whose far end the test can close
            self.master2, self.slave2 = pty.openpty()
            fcntl.ioctl(self.slave2, termios.TIOCSWINSZ, struct.pack("HHHH", size[0], size[1], 0, 0))
        streams = {"pty": self.slave, "pty2": self.slave2, "null": subprocess.DEVNULL, "pipe": subprocess.PIPE}
        self.screen = b""
        self.proc = subprocess.Popen([GROVE, *args], cwd=cwd, env=env or clean_env(), stdin=streams[stdin],
                                     stdout=subprocess.PIPE, stderr=streams[stderr])
        SESSIONS.append(self)

    def pump(self, wait=0.05):
        fds = [fd for fd in (self.master, self.master2) if fd is not None]
        ready, _, _ = select.select(fds, [], [], wait)
        for fd in ready:
            try:
                self.screen += os.read(fd, 65536)
            except OSError:
                pass

    def expect(self, text, since=0):
        """Wait until text has been drawn after offset since; return the new offset."""
        deadline = time.monotonic() + TIMEOUT
        while text.encode() not in self.screen[since:]:
            if time.monotonic() > deadline or self.proc.poll() is not None:
                raise AssertionError(f"never drew {text!r}; exit={self.proc.poll()} screen tail={self.screen[-600:]!r}")
            self.pump()
        return len(self.screen)

    def send(self, data):
        os.write(self.master, data)

    def resize(self, rows, cols):
        fcntl.ioctl(self.slave, termios.TIOCSWINSZ, struct.pack("HHHH", rows, cols, 0, 0))
        self.proc.send_signal(signal.SIGWINCH)

    def finish(self):
        """Wait for exit, draining the screen; return (code, stdout)."""
        deadline = time.monotonic() + TIMEOUT
        while self.proc.poll() is None:
            if time.monotonic() > deadline:
                self.proc.kill()
                raise AssertionError(f"did not exit; screen tail={self.screen[-600:]!r}")
            self.pump()
        for _ in range(3):
            self.pump(0.005)
        out = self.proc.stdout.read()
        self.after = termios.tcgetattr(self.slave)
        for fd in (self.master, self.slave, self.master2, self.slave2):
            if fd is not None:
                try:
                    os.close(fd)
                except OSError:
                    pass
        return self.proc.returncode, out

    def restored(self):
        """Modes equal to before the run, and the alternate screen left last."""
        assert same_modes(self.after, self.before), f"terminal modes changed:\n before {self.before}\n after  {self.after}"
        on, off = self.screen.rfind(ALT_ON), self.screen.rfind(ALT_OFF)
        assert on >= 0 and off > on, f"alternate screen entered at {on}, left at {off}"
        assert self.screen.rfind(b"\x1b[?25h") > self.screen.rfind(b"\x1b[?25l"), "the cursor stayed hidden"


def same_modes(a, b):
    """Compare termios attributes, ignoring PENDIN: the kernel's own note that
    typed input awaits re-reading after a raw-to-canonical switch. A program
    does not set it and the next read clears it."""
    pendin = getattr(termios, "PENDIN", 0)
    return a[:3] == b[:3] and a[3] & ~pendin == b[3] & ~pendin and a[4:] == b[4:]


def clean_env(**extra):
    env = {k: v for k, v in os.environ.items() if not k.startswith(("TEA_", "UV_")) and k not in GIT_LOCATION}
    env.update(TERM="xterm-256color", **extra)
    return env


def check(condition, message):
    if not condition:
        raise AssertionError(message)


def main_board(s):
    """Leave the current view for main's own board, which the lifecycle checks drive."""
    s.expect("Board: current view")
    s.send(b"b" + DOWN + ENTER)
    s.expect("First on main")


# --- scenarios ---------------------------------------------------------------

def select_and_show(root, wt, base):
    """current view -> card -> version -> resolve -> show reads the selected bytes."""
    for flags in ([], ["--json"]):
        s = Session(root, flags)
        s.expect("Board: current view")
        s.expect("First on feature")  # feature changed G-001 after main: its state is current
        s.send(b"l" + ENTER)  # the card is Active
        # The detail opens on the history of its current state, read from Git.
        s.expect("Timeline on branch feature")
        s.expect(f"active     {short(root, 'feature')}  feature")
        check(s.proc.poll() is None, "opening a card must not resolve or exit")
        s.send(b"v")
        mark = s.expect("versions differ")
        s.expect("History on branch feature")
        # Rows fold by content, current first: feature's (branch, checkout),
        # then main's older one. Enter on a fold only lists its places.
        s.send(DOWN + ENTER)
        s.expect("checkout feature-wt (feature)", mark)
        check(s.proc.poll() is None, "opening a fold must not resolve or exit")
        s.send(DOWN * 2)
        s.expect("First on feature", mark)
        s.expect("Selector: live:feature-wt:", mark)
        s.send(ENTER)
        code, out = s.finish()
        s.restored()
        check(code == 0, f"exit {code}")
        check(b"\x1b" not in out, f"interface bytes reached stdout: {out!r}")
        project = json.loads(out)["project"] if flags else out.decode()[:-1]
        check(project == wt, f"stdout names {project!r}, want {wt!r}; raw {out!r}")
        check(flags or out == (wt + "\n").encode(), f"plain result is exactly the path and a newline: {out!r}")
        tail = s.screen[s.screen.rfind(ALT_OFF):]
        check(b"Checkout: " + wt.encode() in tail and b"refs/heads/feature" in tail, f"context belongs on stderr after the screen: {tail!r}")
        shown = subprocess.run([GROVE, "--project", project, "show", "G-001"], capture_output=True, cwd=base, env=clean_env())
        with open(os.path.join(wt, "grove", "work", "G-001-first.md"), "rb") as f:
            check(shown.returncode == 0 and shown.stdout == f.read(), "show did not read the selected bytes")


def leave_without_selecting(root, wt, base):
    """q and Esc exit 0 with no output; Ctrl-C exits 1 with no output."""
    for keys, want, where in ((b"q", 0, "board"), (ESC, 0, "board"), (CTRL_C, 1, "board"), (b"q", 0, "versions"), (CTRL_C, 1, "versions")):
        s = Session(root)
        main_board(s)
        if where == "versions":
            s.send(ENTER + b"v")
            s.expect("versions differ")
            s.send(DOWN)
        s.send(keys)
        code, out = s.finish()
        s.restored()
        check(code == want and out == b"", f"{keys!r} on {where}: exit {code}, stdout {out!r}")
        check((b"interrupted" in s.screen) == (want == 1), f"{keys!r}: interruption message mismatch")


def refuses_without_terminal(root, wt, base):
    """No terminal: prompt refusal, exit 1, no stdout, no terminal modes touched."""
    for stdin, stderr in (("null", "pty"), ("pipe", "pty"), ("pty", "pipe"), ("null", "pipe")):
        for flags in ([], ["--json"]):
            s = Session(root, flags, stdin=stdin, stderr=stderr)
            if stdin == "pipe":
                s.proc.stdin.close()
            piped = s.proc.stderr.read() if stderr == "pipe" else b""
            code, out = s.finish()
            text = piped + s.screen
            check(code == 1 and out == b"", f"{stdin}/{stderr}: exit {code}, stdout {out!r}")
            check(b"needs a terminal" in text and b"grove list" in text and b"--help" in text, f"no guidance: {text!r}")
            check(b"\x1b" not in text and same_modes(s.after, s.before), "the refusal touched the terminal")
    # Help and explicit commands stay noninteractive, with or without a terminal.
    for args, cwd in ((["--help"], base), (["help"], base), (["list"], root), (["versions", "G-001"], root)):
        for stdin, stderr in (("null", "pipe"), ("pty", "pty")):
            s = Session(cwd, args, stdin=stdin, stderr=stderr)
            code, out = s.finish()
            check(code == 0 and out and b"\x1b" not in out + s.screen and same_modes(s.after, s.before), f"{args} {stdin}/{stderr}: exit {code} {s.screen!r}")
    s = Session(base, ["board"])
    code, out = s.finish()
    check(code == 2 and out == b"" and b"unknown command board" in s.screen and ALT_ON not in s.screen, "board must not be a command")
    s = Session(base)  # a terminal, but no project: said plainly, outside the interface
    code, out = s.finish()
    check(code == 1 and out == b"" and b"grove.yaml" in s.screen and ALT_ON not in s.screen, f"no project: {code} {s.screen!r}")


def blocked_git(root, wt, base):
    """A Git read that never returns is killed and collected when leaving."""
    tools = os.path.join(base, "tools")
    os.makedirs(tools, exist_ok=True)
    flag, fifo = os.path.join(tools, "block"), os.path.join(tools, "started")
    with open(os.path.join(tools, "git"), "w") as f:
        f.write(f'#!/bin/sh\nif [ -e "{flag}" ]; then echo $$ > "{fifo}"; exec sleep 600; fi\nexec "{GIT}" "$@"\n')
    os.chmod(os.path.join(tools, "git"), 0o755)
    env = clean_env(PATH=tools + os.pathsep + os.environ["PATH"])
    stages = (("refresh", b"q", 0), ("refresh", CTRL_C, 1), ("resolve", b"q", 0), ("resolve", CTRL_C, 1), ("start", b"q", 0),
              ("history", b"q", 0), ("history", CTRL_C, 1), ("history", ESC, None))
    for stage, keys, want in stages:
        if os.path.exists(fifo):
            os.remove(fifo)
        os.mkfifo(fifo)
        if stage == "start":
            open(flag, "w").close()
        s = Session(root, env=env)
        if stage != "start":
            main_board(s)
            if stage == "resolve":
                # Each history is on screen before the next key, so that the
                # only Git left to block is the resolution's.
                s.send(ENTER)
                s.expect(f"{short(root, 'main')}  main")
                s.send(b"v" + DOWN)
                s.expect(f"{short(root, 'feature')}  feature")
                s.send(DOWN + ENTER + DOWN * 2)
                s.expect("Selector: live:.:")
            open(flag, "w").close()
            s.send(b"r" if stage == "refresh" else ENTER)  # for history, Enter opens the card
        # The handshake: the fake git says it started before anything is cancelled.
        reader = os.open(fifo, os.O_RDONLY | os.O_NONBLOCK)
        deadline, data = time.monotonic() + TIMEOUT, b""
        while not data.endswith(b"\n"):
            check(time.monotonic() < deadline, f"{stage}: the blocked git never started")
            s.pump()
            try:
                data += os.read(reader, 64)
            except BlockingIOError:
                pass
        os.close(reader)
        pid = int(data)
        s.expect({"resolve": "Resolving the selected workspace", "history": "reading…"}.get(stage, "Reading branches and checkouts"))
        s.send(keys)
        if want is None:  # Esc leaves the card, not the session: the read dies and keys still work
            deadline = time.monotonic() + TIMEOUT
            while True:
                try:
                    os.kill(pid, 0)
                except ProcessLookupError:
                    break
                check(time.monotonic() < deadline, f"{stage}: Esc left git child {pid} running")
                s.pump()
            os.remove(flag)
            s.send(b"q")
            want = 0
        code, out = s.finish()
        if os.path.exists(flag):
            os.remove(flag)
        s.restored()
        check(code == want and out == b"", f"{stage} {keys!r}: exit {code}, stdout {out!r}")
        try:
            os.kill(pid, 0)
            raise AssertionError(f"{stage} {keys!r}: git child {pid} outlived the session")
        except ProcessLookupError:
            pass


def hangup(root, wt, base):
    """SIGHUP (the window closed) still kills a blocked Git child and restores modes."""
    tools = os.path.join(base, "tools-hup")
    os.makedirs(tools, exist_ok=True)
    flag, fifo = os.path.join(tools, "block"), os.path.join(tools, "started")
    with open(os.path.join(tools, "git"), "w") as f:
        f.write(f'#!/bin/sh\nif [ -e "{flag}" ]; then echo $$ > "{fifo}"; exec sleep 600; fi\nexec "{GIT}" "$@"\n')
    os.chmod(os.path.join(tools, "git"), 0o755)
    os.mkfifo(fifo)
    s = Session(root, env=clean_env(PATH=tools + os.pathsep + os.environ["PATH"]))
    main_board(s)
    open(flag, "w").close()
    s.send(b"r")
    reader = os.open(fifo, os.O_RDONLY | os.O_NONBLOCK)
    deadline, data = time.monotonic() + TIMEOUT, b""
    while not data.endswith(b"\n"):
        check(time.monotonic() < deadline, "the blocked git never started")
        s.pump()
        try:
            data += os.read(reader, 64)
        except BlockingIOError:
            pass
    os.close(reader)
    s.proc.send_signal(signal.SIGHUP)
    code, out = s.finish()
    os.remove(flag)
    s.restored()
    check(code == 1 and out == b"", f"exit {code}, stdout {out!r}")
    try:
        os.kill(int(data), 0)
        raise AssertionError("the git child outlived a hangup")
    except ProcessLookupError:
        pass


def output_failure(root, wt, base):
    """The screen goes away mid-session: exit 1, no result, input modes restored."""
    s = Session(root, stderr="pty2")
    main_board(s)
    os.close(s.master2)
    s.master2 = None
    s.send(ENTER)  # forces a redraw onto the dead screen
    s.send(DOWN)
    code, out = s.finish()
    check(code == 1 and out == b"", f"exit {code}, stdout {out!r}")
    check(same_modes(s.after, s.before), "input terminal modes were not restored after the screen failed")


def resize(root, wt, base):
    s = Session(root)
    mark = s.expect("Active 1")
    s.resize(24, 80)
    mark = s.expect("[Proposed 0] Active 1", mark)
    s.resize(8, 30)
    mark = s.expect("Grove needs 40x10", mark)
    s.resize(30, 120)
    s.expect("Active 1", mark)
    s.send(b"q")
    code, out = s.finish()
    s.restored()
    check(code == 0 and out == b"", f"exit {code}")


def focus_rereads(root, wt, base):
    """The board asks for focus reports and re-reads when focus returns (G-124), and turns them off on the way out."""
    s = Session(root)
    mark = s.expect("read 2 branches")
    check(b"\x1b[?1004h" in s.screen, "the board did not ask for focus reports")
    git(root, "branch", "focus-probe")
    s.send(b"\x1b[I")  # focus in
    mark = s.expect("Reading branches and checkouts", mark)
    s.send(b"s")  # only changed cells are redrawn; the sources screen lists what the re-read found
    s.expect("focus-probe", mark)
    s.send(b"q")
    code, out = s.finish()
    s.restored()
    check(code == 0 and out == b"", f"exit {code}, stdout {out!r}")
    check(s.screen.rfind(b"\x1b[?1004l") > s.screen.rfind(b"\x1b[?1004h"), "focus reports were left on")


focus_rereads.mutates = True  # the probe's branch changes the repository on purpose


def writes_no_logs(root, wt, base):
    """The framework's log switches are in the environment and must do nothing."""
    logs = os.path.join(base, "logs")
    os.makedirs(logs, exist_ok=True)
    env = clean_env(TEA_DEBUG="true", TEA_TRACE=os.path.join(logs, "trace.log"), UV_DEBUG=os.path.join(logs, "uv.log"))
    s = Session(root, env=env)
    main_board(s)
    s.send(ENTER + DOWN + b"q")
    code, _ = s.finish()
    s.restored()
    check(code == 0 and os.listdir(logs) == [], f"log files appeared: {os.listdir(logs)}")


def review_and_integrate(root, wt, base):
    """A candidate in review, judged from the board: standing, changes and a diff; a approves on feature; i merges into main and writes done."""
    for key, value in (("user.name", "t"), ("user.email", "t@t"), ("commit.gpgsign", "false"), ("maintenance.auto", "false")):
        git(root, "config", key, value)
    with open(os.path.join(root, "grove.yaml"), "w") as f:
        f.write("schema_version: 3\nrecords: grove\ntarget: main\n")
    git(root, "commit", "-qam", "target")
    with open(os.path.join(wt, "code.txt"), "w") as f:
        f.write("hello\n")
    git(wt, "add", "-A")
    git(wt, "commit", "-qm", "feat: code")
    candidate = subprocess.run([GIT, "-C", wt, "rev-parse", "HEAD"], check=True, capture_output=True, env=clean_env()).stdout.decode().strip()
    subprocess.run([GROVE, "--project", wt, "update", "G-001", "--set", "status=review", "--set", f"candidate={candidate}", "--commit"],
                   check=True, capture_output=True, cwd=base, env=clean_env())
    s = Session(root)
    s.expect("Board: current view")
    s.send(b"ll" + ENTER)  # the Review column's card
    mark = s.expect(f"Review: candidate {candidate[:7]} · not yet approved")  # the renderer redraws lines from their first changed cell, so expectations stay short
    s.expect("code.txt  +1 −0")
    s.send(b"\t" + ENTER)  # no linked records, so Tab lands on the first changed file
    s.expect("Diff of code.txt", mark)  # siblings of one frame are searched from the same offset
    mark = s.expect("+hello", mark)
    s.send(ESC)  # alone: an Esc followed at once by a letter reads as Alt
    mark = s.expect("An outcome.", mark)
    s.send(b"a")
    mark = s.expect("Approve G-001 on branch feature", mark)
    s.send(b"Ship it" + ENTER)
    s.expect("Approved G-001", mark)
    mark = s.expect("The board has been re-read.", mark)
    with open(os.path.join(wt, "grove", "work", "G-001-first.md")) as f:
        record = f.read()
    check(f'approved: "{candidate}"' in record and record.endswith(": Ship it\n"), f"feature's record after approval: {record!r}")
    s.send(ESC)
    # The renderer scrolls and redraws lines from their first changed cell, so
    # the approval is checked in the file above; the re-read changes are new.
    mark = s.expect("only the record changed since it", mark)
    s.send(b"i")
    mark = s.expect("mark G-001 done? y/n", mark)
    s.send(b"y")
    mark = s.expect("remove its worktree? y/n", mark)
    s.send(b"n")
    s.expect("Integration of G-001", mark)
    s.expect("merge: merge commit", mark)  # main gained the target commit after feature branched
    s.expect("done: G-001 done at commit", mark)
    mark = s.expect("The board has been re-read.", mark)  # Esc waits for the re-read
    s.send(ESC)
    s.expect("· done", mark)
    s.send(b"q")
    code, out = s.finish()
    s.restored()
    check(code == 0 and out == b"", f"exit {code}, stdout {out!r}")
    with open(os.path.join(root, "grove", "work", "G-001-first.md")) as f:
        record = f.read()
    check("status: done" in record and f'approved: "{candidate}"' in record, f"main's record after integration: {record!r}")
    check(os.path.isdir(wt), "n kept the worktree")


review_and_integrate.mutates = True  # approval and the merge change the repository on purpose

FAKE_CLAUDE = r"""#!/bin/sh
[ "$1" = --version ] && { echo 'fake 0.1'; exit 0; }
echo start >> "@STARTS@"
trap 'echo "{\"type\":\"result\",\"subtype\":\"error_during_execution\",\"is_error\":true}"; exit 130' INT
echo '{"type":"system","subtype":"init","model":"fake-model"}'
if [ -e "@QUESTION@" ] && [ ! -e grove/G-002-colour.md ]; then
  printf -- '---\nid: G-002\ntype: question\ntitle: Which colour?\nstatus: open\nblocks: ["G-001"]\n---\nRed or blue?\n' > grove/G-002-colour.md
  git add grove/G-002-colour.md
  git -c user.name=t -c user.email=t@t -c commit.gpgsign=false -c maintenance.auto=false commit -qm question
  echo '{"type":"result","subtype":"success","is_error":false,"result":"Waiting on G-002."}'
  exit 0
fi
if [ -e "@FINISH@" ]; then
  head=$(git rev-parse HEAD)
  printf -- '---\nid: G-001\ntype: work\ntitle: First on main\nstatus: review\ncandidate: "%s"\n---\nAn outcome.\n' "$head" > grove/work/G-001-first.md
  git -c user.name=t -c user.email=t@t -c commit.gpgsign=false -c maintenance.auto=false commit -qam review
  echo '{"type":"result","subtype":"success","is_error":false,"total_cost_usd":0.25,"result":"## Handoff\n\nReady for review."}'
  exit 0
fi
awk 'BEGIN { for (i = 0; i < 20000; i++) printf "{\"type\":\"assistant\",\"message\":{\"content\":[{\"type\":\"text\",\"text\":\"step %d\"}]}}\n", i }'
echo partial > partial.txt
n=0
while [ $n -lt 600 ]; do sleep 0.1; n=$((n+1)); done
"""


def attempts_of(root):
    """Each attempt directory of the repository, by name, with whether an owner holds its lock."""
    found = {}
    top = os.path.join(root, ".git", "grove", "attempts")
    for name in sorted(os.listdir(top)) if os.path.isdir(top) else []:
        if name.endswith(".tmp") or not os.path.isdir(os.path.join(top, name)):
            continue
        try:
            fd = os.open(os.path.join(top, name, "owner.lock"), os.O_RDWR)
        except OSError:  # a launch between renaming its directory and creating the lock
            continue
        try:
            fcntl.flock(fd, fcntl.LOCK_SH | fcntl.LOCK_NB)
            running = False
        except BlockingIOError:
            running = True
        finally:
            os.close(fd)
        found[name] = (running, os.path.exists(os.path.join(top, name, "result.json")))
    return found


def stop_attempts(root):
    """Kill the owner and provider group of every attempt still running, so a failed scenario leaves none; a finished one's pids may belong to others by now."""
    top = os.path.join(root, ".git", "grove", "attempts")
    for name, (running, _) in attempts_of(root).items():
        if not running:
            continue
        for file, key, sign in (("attempt.json", "owner_pid", 1), ("child.json", "pgid", -1)):
            try:
                with open(os.path.join(top, name, file)) as f:
                    pid = json.load(f).get(key, 0)
                if pid > 0:
                    os.kill(sign * pid, signal.SIGKILL)
            except (OSError, ValueError):
                pass


def attempt_lifecycle(root, wt, base):
    """R launches a bounded attempt of a fake provider that floods its events; the board quits while it runs, a new session reconnects to the same attempt and stops it; the next waits on a question that then refuses a launch; e on that attempt suspends the board for the owner's editor, resumes it, and resolves and commits the answer on the branch; R then continues on the same branch to a candidate in review."""
    for key, value in (("user.name", "t"), ("user.email", "t@t"), ("commit.gpgsign", "false"), ("maintenance.auto", "false")):
        git(root, "config", key, value)
    git(root, "worktree", "remove", "--force", wt)
    git(root, "branch", "-qD", "feature")
    with open(os.path.join(root, "grove.yaml"), "w") as f:
        f.write("schema_version: 3\nrecords: grove\ntarget: main\n")
    git(root, "commit", "-qam", "target")
    # The skill as init leaves it: written, not committed (G-150).
    skill = os.path.join(root, ".claude", "skills", "grove-work", "SKILL.md")
    os.makedirs(os.path.dirname(skill))
    with open(skill, "w") as f:
        f.write("---\nname: grove-work\n---\n")
    tools = os.path.join(base, "tools-attempt")
    os.makedirs(tools)
    starts, finish, question = os.path.join(tools, "starts"), os.path.join(tools, "finish"), os.path.join(tools, "question")
    fake = os.path.join(tools, "claude")
    with open(fake, "w") as f:
        f.write(FAKE_CLAUDE.replace("@STARTS@", starts).replace("@FINISH@", finish).replace("@QUESTION@", question))
    os.chmod(fake, 0o755)
    # The owner's editor: it must get the terminal back in canonical mode, and
    # it appends the answer to the file it is given.
    editor, edits = os.path.join(tools, "editor"), os.path.join(tools, "edits")
    with open(editor, "w") as f:
        f.write(f"#!/bin/sh\n[ -t 0 ] && [ -t 1 ] || exit 3\necho \"$1 $(stty -a | grep -o -- '-*icanon')\" >> '{edits}'\n"
                "echo EDITOR-RAN\nprintf 'Blue.\\n' >> \"$1\"\n")
    os.chmod(editor, 0o755)
    env = clean_env(GROVE_CLAUDE=fake, VISUAL=editor)
    count = lambda: open(starts).read().count("start") if os.path.exists(starts) else 0

    s = Session(root, env=env)
    s.expect("Board: current view")
    s.send(ENTER)
    mark = s.expect("Attempts: none")
    # Without run: in grove.yaml the one launch line needs the budget and the
    # mode typed, in run's own flags, and refuses what run refuses.
    s.send(b"R")
    mark = s.expect("Launch G-001 ▏ · no budget, no mode, to the handoff, model default, effort default", mark)
    s.send(ENTER)
    mark = s.expect("type --budget USD and --permission-mode MODE", mark)
    s.send(b"--frob" + ENTER)
    mark = s.expect("unknown option --frob", mark)
    s.send(b"\x7f" * len("--frob") + b"--budget 1 --permission-mode auto")
    mark = s.expect("$1, mode auto, to the handoff", mark)
    # Uncommitted, the skill refuses the launch before the provider starts;
    # committed, the same launch starts and warns that no reviewer is there.
    s.send(ENTER)
    mark = s.expect("NOT DONE: .claude/skills/grove-work/SKILL.md is not committed at HEAD", mark)
    mark = s.expect("The board has been re-read.", mark)
    check(count() == 0 and not os.path.exists(os.path.join(root, ".claude", "worktrees")), "a refused launch started nothing")
    git(root, "add", skill)
    git(root, "commit", "-qm", "commit what init wrote")
    s.send(ESC)
    mark = s.expect("R launches one", mark)
    s.send(b"R")
    mark = s.expect("Launch G-001", mark)
    s.send(b"--budget 1 --permission-mode auto" + ENTER)
    s.expect("Launch of an attempt of G-001", mark)
    s.expect("warning: .claude/agents/grove-reviewer.md is not in", mark)
    s.expect("started; owner pid", mark)
    mark = s.expect("The board has been re-read.", mark)  # Esc waits for the re-read; both may be one frame
    s.send(ESC)
    mark = s.expect("A lists them", mark)
    s.send(b"A")
    mark = s.expect("Attempts of G-001", mark)
    s.send(ENTER)
    s.expect("Running.", mark)
    mark = s.expect("step 19999", mark)  # the newest of 20,000 events, shown while it runs
    for _ in range(20):  # keys stay immediate however much the provider writes
        s.send(b"\x1b[6~")
    s.send(b"q")
    code, out = s.finish()
    s.restored()
    check(code == 0 and out == b"", f"quit while running: exit {code}, stdout {out!r}")
    found = attempts_of(root)
    check(len(found) == 1 and list(found.values())[0] == (True, False) and count() == 1, f"the attempt outlives the board: {found}, {count()} starts")
    first = list(found)[0]

    # From here this checkout's run: defaults launch with Enter alone.
    with open(os.path.join(root, "grove.yaml"), "a") as f:
        f.write("run:\n  budget: 1\n  permission_mode: auto\n")
    git(root, "commit", "-qam", "defaults")

    s = Session(root, env=env)  # reconnect: the same attempt, never a second start
    s.expect("Board: current view")
    mark = s.expect("running")  # the card's tag
    # A status committed on the running attempt's branch moves its card with no key pressed (G-124).
    wt1 = os.path.join(root, ".claude", "worktrees", "worktree-G-001")
    record = os.path.join(wt1, "grove", "work", "G-001-first.md")
    with open(record) as f:
        text = f.read()
    with open(record, "w") as f:
        f.write(text.replace("status: proposed", "status: active"))
    git(wt1, "commit", "-qam", "active")
    s.expect("Active 1", mark)
    s.send(ENTER)
    mark = s.expect("A lists them")
    s.send(b"A")
    mark = s.expect("Attempts of G-001", mark)
    s.send(ENTER)
    mark = s.expect("Running.", mark)
    check(count() == 1, "reconnecting started nothing")
    s.send(b"x")
    mark = s.expect("Stop attempt", mark)
    s.send(b"y")
    s.expect(f"Stop of attempt {first}", mark)
    mark = s.expect("The board has been re-read.", mark)
    s.send(ESC)
    mark = s.expect("Stopped (exit 130)", mark)  # only the changed cells are redrawn
    check(os.path.exists(os.path.join(root, ".claude", "worktrees", "worktree-G-001", "partial.txt")), "Stop kept the partial work")
    s.send(ESC)
    mark = s.expect("0 need you · 0 running · 1 settled", mark)  # the list, redrawn where it differs
    s.send(ESC)
    mark = s.expect("R launches one", mark)
    # The next attempt persists a question and ends: a wait, which a launch then refuses.
    open(question, "w").close()
    s.send(b"R")
    mark = s.expect("Launch G-001 ▏ · $1, mode auto, to the handoff, model default, effort default", mark)
    s.send(ENTER)
    s.expect("started; owner pid", mark)
    mark = s.expect("The board has been re-read.", mark)  # Esc waits for the re-read; both may be one frame
    s.send(ESC)
    mark = s.expect("answer question G-002", mark)
    s.send(b"R")
    mark = s.expect("blocked by open question G-002", mark)
    check(count() == 2, "the wait started nothing more")
    # The owner answers from the attempt (G-125): e suspends the board for the
    # editor on the branch's copy, then resolves and commits it there.
    s.send(b"A")
    mark = s.expect("answer question G-002", mark)
    s.send(ENTER)
    mark = s.expect("e answers G-002 in your editor", mark)
    s.send(b"e")
    mark = s.expect("EDITOR-RAN", mark)
    after = s.expect("Resolve G-002 and commit it with your answer on branch worktree-G-001? y/n", mark)
    suspended = s.screen[mark:after]
    check(suspended.find(ALT_ON) >= 0, f"the board did not resume on the alternate screen: {suspended[:300]!r}")
    path = os.path.join(wt1, "grove", "G-002-colour.md")
    with open(edits) as f:
        check(f.read() == path + " icanon\n", f"the editor ran on {open(edits).read()!r}, want {path} in canonical mode")
    s.send(b"y")
    s.expect("Answer to G-002", mark)
    s.expect("next: R on G-001 launches its next attempt", mark)
    mark = s.expect("The board has been re-read.", mark)
    with open(path) as f:
        text = f.read()
    check("status: resolved" in text and text.endswith("## Answer\n\nBlue.\n"), f"the answer and the status: {text!r}")
    log = subprocess.run([GIT, "-C", wt1, "log", "-1", "--name-only", "--format=%s"], check=True, capture_output=True, env=clean_env()).stdout.decode()
    check(log == "docs(G-002): set status=resolved\n\ngrove/G-002-colour.md\n", f"one commit of the question alone on the branch: {log!r}")
    s.send(ESC)  # to the question
    mark = s.expect("G-002 · question · resolved", mark)
    s.send(ESC)  # to the attempt
    mark = s.expect("question answered: R again", mark)
    s.send(b"o")  # G-001's detail, launchable again
    mark = s.expect("R launches one", mark)
    open(finish, "w").close()
    s.send(b"R")
    mark = s.expect("$1, mode auto", mark)
    s.send(b"--effort xhigh")  # a typed flag overrides for this launch only
    mark = s.expect("$1, mode auto, to the handoff, model default, effort xhigh", mark)
    s.send(ENTER)
    s.expect("worktree: reusing", mark)
    s.expect("started; owner pid", mark)
    mark = s.expect("The board has been re-read.", mark)  # Esc waits for the re-read; both may be one frame
    s.send(ESC)
    s.expect("judge candidate", mark)  # the poll saw it end
    s.expect("Review: candidate", mark)  # and the board was re-read: the detail is a review
    s.send(b"q")
    code, out = s.finish()
    s.restored()
    check(code == 0 and out == b"", f"exit {code}, stdout {out!r}")
    found = attempts_of(root)
    check(len(found) == 3 and all(v == (False, True) for v in found.values()) and count() == 3, f"three attempts, all ended: {found}, {count()} starts")
    launched = [json.load(open(os.path.join(root, ".git", "grove", "attempts", name, "attempt.json"))) for name in sorted(found)]
    asked = [(l["budget_usd"], l["permission_mode"], l.get("effort", "")) for l in launched]
    check(asked == [("1", "auto", ""), ("1", "auto", ""), ("1", "auto", "xhigh")], f"typed, defaulted, then overridden: {asked}")


attempt_lifecycle.mutates = True  # attempts, a worktree and the fake's commit change the repository on purpose

def dependencies(root, wt, base):
    """board -> g -> focus and select -> preview -> back -> a record's detail -> back -> refresh -> resize -> exit, reading nothing into stdout (G-161)."""
    with open(os.path.join(root, "grove", "work", "G-002-second.md"), "w") as f:
        f.write('---\nid: G-002\ntype: work\ntitle: Second needs first\nstatus: proposed\ndepends_on: ["G-001"]\n---\nAn outcome.\n')
    git(root, "add", "-A")
    git(root, "commit", "-qm", "second")
    s = Session(root)
    mark = s.expect("Board: current view")
    s.send(b"g")  # each frame's lines are searched from where the frame began
    s.expect("Connected · 2 work", mark)
    s.expect("← Needs", mark)
    mark = s.expect("G-002 proposed · layer 1", mark)  # the board's focused card
    s.send(b"\x1b[A")  # up
    mark = s.expect("active · layer 0", mark)  # G-001 in the current view: feature's state; redraws start at the first changed cell
    s.send(DOWN)
    mark = s.expect("proposed · layer 1", mark)
    s.send(b" p")
    s.expect("Selection preview", mark)
    s.expect("Order:    G-002", mark)
    s.expect("Outside the selection, not added", mark)
    mark = s.expect("awaiting implementation", mark)  # G-001 is proposed in this checkout, whatever feature holds
    s.send(ESC)
    time.sleep(0.2)  # alone: an Esc followed at once by a letter reads as Alt
    mark = s.expect("p preview (1 selected)", mark)
    s.send(ENTER)
    mark = s.expect("dependencies › G-002", mark)
    s.send(ESC)
    time.sleep(0.2)
    mark = s.expect("Connected · 2 work", mark)
    s.send(b"r")
    mark = s.expect("Reading branches and checkouts", mark)
    s.resize(24, 80)
    mark = s.expect("Tab tree", mark)
    s.send(b"\t")
    s.expect("G-002 proposed · layer 1", mark)
    s.send(b"q")
    code, out = s.finish()
    s.restored()
    check(code == 0 and out == b"", f"exit {code}, stdout {out!r}")


dependencies.mutates = True  # the second record's commit changes the repository on purpose


SCENARIOS = [select_and_show, leave_without_selecting, refuses_without_terminal, blocked_git, hangup, output_failure, resize, focus_rereads, writes_no_logs, review_and_integrate,
             attempt_lifecycle, dependencies]


def main():
    if not ONLY:  # every scenario at once, each in its own process with its own repository
        runs = [subprocess.Popen([sys.executable, __file__, GROVE, s.__name__], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                for s in SCENARIOS]
        for proc in runs:
            sys.stdout.write(proc.communicate()[0].decode(errors="replace"))
        sys.exit(1 if any(proc.returncode for proc in runs) else 0)
    base = os.path.realpath(tempfile.mkdtemp(prefix="grove-terminal-"))
    failed = 0
    try:
        root, wt = fixture(os.path.join(base, "repo"))
        before = tree(os.path.join(base, "repo"))
        for scenario in SCENARIOS:
            if ONLY and scenario.__name__ not in ONLY:
                continue
            try:
                scenario(root, wt, base)
                check(getattr(scenario, "mutates", False) or tree(os.path.join(base, "repo")) == before, "files in the repository changed")
                print(f"ok    {scenario.__name__}")
            except AssertionError as e:
                failed += 1
                print(f"FAIL  {scenario.__name__}: {e}")
    finally:
        for s in SESSIONS:
            if s.proc.poll() is None:
                s.proc.kill()
                s.proc.wait()
        stop_attempts(os.path.join(base, "repo", "main"))
        shutil.rmtree(base, ignore_errors=True)
    sys.exit(1 if failed else 0)


if __name__ == "__main__":
    main()
