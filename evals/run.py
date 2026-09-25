#!/usr/bin/env python3
"""Behavioral evaluations of Grove's headless shaping workflow (G-108).

    python3 evals/run.py run --runs N --budget USD --model MODEL \\
        --permission-mode MODE --config-dir DIR [--case NAME]... [--out DIR]
    python3 evals/run.py run --harness codex --runs N --model MODEL --effort EFFORT \\
        --permission-mode MODE --config-dir DIR --max-seconds S --max-plan-percent P [--case NAME]... [--out DIR]
    python3 evals/run.py selftest

`run` spends money: every spend parameter is required and has no default.
`selftest` spends nothing: it drives the same runner with a fake `claude` and a
fake `codex` and asserts that the checks pass and fail where they should. evals/README.md says
what a run retains, what each check means, and how to score the rubric.
"""
import argparse, datetime, glob, json, os, pwd, re, shlex, shutil, signal, subprocess, sys, tempfile, time, tomllib

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
FIXTURE = os.path.join(ROOT, "evals", "fixture")
# Variables through which Git takes a repository from its caller; a hook exports GIT_DIR (G-089).
GIT_LOCATION = ("GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_NAMESPACE")
IDENTITY = ["-c", "user.name=Grove Eval", "-c", "user.email=eval@example.invalid", "-c", "commit.gpgsign=false"]
TIMEOUT = 1800  # seconds per run; failure detection only, the budget is the real bound
PAIR = ("presumes choice", "planted question", "brief constraint", "handoff")
KNOWN = ("constraint applied", "brief constraint", "handoff")
# records: (key, type, title, fields) added to the case's copy of the fixture, body from evals/fixture/records/KEY.md,
# a list field naming other keys; holding is the record whose constraint the proposal must apply (G-154).
CASES = {
    "missing-choice": {"topic": "hide finished tasks from tasks list by default", "question": True, "rubric": PAIR},
    "companion": {"topic": "let tasks list filter by tag", "question": False, "rubric": PAIR},
    "listed-constraint": {"topic": "make due dates for tasks ready to assign", "question": False, "rubric": KNOWN,
                          "records": (("export", "work", "Add tasks export", {"status": "done"}),
                                      ("csv", "work", "Export tasks as CSV", {}),
                                      ("remind", "work", "Remind the owner of tasks due today", {}),
                                      ("due", "work", "Give tasks a due date", {"depends_on": ["export"], "relates_to": ["csv", "remind"]})),
                          "holding": "export", "distractors": ("csv", "remind")},
    "code-constraint": {"topic": "add a tasks tag command that adds or removes tags on an existing task", "question": False, "rubric": KNOWN,
                        "records": (("notes", "decision", "Keep the owner's notes in task files", {"status": "accepted"}),),
                        "holding": "notes", "distractors": ()},
}
# The G-108 pair is the regression rerun for a guide edit; the G-154 cases run only when --case names them.
DEFAULT = ("missing-choice", "companion")
# Files some step of the shaping guide needs for these topics; any other read is listed as unneeded.
READERS = ("cat", "head", "tail", "sed", "nl", "less", "awk")
NEEDED = {"AGENTS.md", "CLAUDE.md", "grove.yaml", "grove/brief.md", "tasks.py", ".agents/skills/grove-shape/SKILL.md"}
# ponytail: the largest five-hour plan use one G-135 run took (gpt-6-astra, G-143); every run is assumed to take at least this, measure per model if it binds
PLAN_POINTS_FLOOR = 13
CUSTOMIZATION = ("CLAUDE.md", "agents", "commands", "output-styles", "hooks", "settings.local.json")
# A login writes settings.json; these keys shape the terminal and memory, not what the agent reads or may do.
SETTINGS = {"tui", "theme", "autoMemoryEnabled"}
GROVE = "grove version (its revision can lag in a linked worktree; the digest pins the guides)"
UNTOUCHED = {"proposed", "open", "current", None}
# Claude's own auth variables pass through; every other CLAUDE* variable is the caller's session leaking in.
AUTH = ("CLAUDE_CODE_OAUTH_TOKEN",)
# What in a CODEX_HOME shapes a session beyond Codex's defaults; config.toml and skills/ are checked apart.
CODEX_CUSTOMIZATION = ("AGENTS.md", "AGENTS.override.md", "rules", "prompts", "hooks.json", "hooks")
# approve-for-me is the workspace-write sandbox with Codex's automatic reviewer on approvals, the nearest to Claude's auto.
CODEX_MODES = ("read-only", "workspace-write", "danger-full-access", "approve-for-me")


def env():
    return {k: v for k, v in os.environ.items() if (k in AUTH or not k.startswith("CLAUDE")) and k not in GIT_LOCATION}


class Failed(Exception):
    """A command the runner needed failed."""


def sh(*args, cwd=None, check=True, environ=None):
    r = subprocess.run(args, cwd=cwd, capture_output=True, text=True, env=environ or env())
    if check and r.returncode:
        raise Failed(f"{' '.join(args)}: exit {r.returncode}\n{r.stderr}")
    return r


def git(cwd, *args, check=True):
    return sh("git", *IDENTITY, "-C", cwd, *args, check=check).stdout.strip()


def frontmatter(text):
    fields = {}
    if text.startswith("---\n"):
        key = None
        for line in text[4:].partition("\n---")[0].splitlines():
            if key and re.match(r"\s+- ", line):  # a block-style list item of the key above
                fields[key] = (fields[key] or []) + [line.split("- ", 1)[1].strip(" '\"")]
                continue
            key, _, value = line.partition(":")
            value = value.strip()
            try:
                key = key.strip()
                fields[key] = json.loads(value) if value[:1] in '["' else value
            except ValueError:  # YAML Grove accepts but does not write, such as [G-002]
                fields[key] = [v.strip(" '\"") for v in value[1:-1].split(",") if v.strip()] if value[:1] == "[" else value
    return fields


def records(clone, ref):
    """Every record's frontmatter and source on ref, by path."""
    out = {}
    for path in git(clone, "ls-tree", "-r", "--name-only", ref, "grove/").splitlines():
        if path.endswith(".md") and path != "grove/brief.md":
            text = git(clone, "show", f"{ref}:{path}")
            out[path] = {"fields": frontmatter(text), "source": text}
    return out


def codex_home(home):
    """Why a CODEX_HOME is not clean; empty when it is. Codex writes a trust table per directory it runs in,
    a login writes tui state, and skills/.system holds Codex's bundled skills: those are recorded, not refused."""
    found = [c for c in CODEX_CUSTOMIZATION if os.path.exists(os.path.join(home, c))]
    memories = os.path.join(home, "memories")
    found += ["memories/"] if os.path.isdir(memories) and os.listdir(memories) else []
    skills = os.path.join(home, "skills")
    found += [f"skills/{e}" for e in sorted(os.listdir(skills)) if e != ".system"] if os.path.isdir(skills) else []
    try:
        config = tomllib.load(open(os.path.join(home, "config.toml"), "rb")) if os.path.exists(os.path.join(home, "config.toml")) else {}
    except (tomllib.TOMLDecodeError, OSError) as err:
        return found + [f"config.toml unreadable: {err}"]
    for key, value in config.items():
        if key == "projects" and isinstance(value, dict):
            found += [f"config.toml projects.{p}.{k}" for p, t in value.items() for k in (t if isinstance(t, dict) else {"(not a table)": 0}) if k != "trust_level"]
        elif key != "tui":
            found.append(f"config.toml key {key}")
    return found


def codex_env(home, grove):
    """Codex runs each command in the user's login shell, whose profile can put an installed grove before the built
    one (G-135 review). A ZDOTDIR of the runner's own keeps a zsh user's startup files out, and its .zprofile restores
    the PATH that /etc/zprofile reorders, so the session has the Claude row's PATH; run_on verifies the result."""
    e = {k: v for k, v in env().items() if not k.startswith("CODEX_") or k == "CODEX_API_KEY"}
    e["PATH"] = os.path.dirname(grove) + os.pathsep + e.get("PATH", "")
    zdotdir = os.path.join(os.path.dirname(os.path.dirname(grove)), "zdotdir")
    os.makedirs(zdotdir, exist_ok=True)
    open(os.path.join(zdotdir, ".zprofile"), "w").write(f"export PATH={shlex.quote(e['PATH'])}\n")
    return dict(e, CODEX_HOME=home, ZDOTDIR=zdotdir)


def system_skills(home):
    d = os.path.join(home, "skills", ".system")
    return ", ".join(sorted(e for e in os.listdir(d) if not e.startswith("."))) if os.path.isdir(d) else "none yet"


def fill(project, path, key):
    """Give a record grove new wrote the body evals/fixture/records/KEY.md, keeping its frontmatter."""
    with open(os.path.join(project, path), "r+") as f:
        head = f.read().split("\n---\n", 1)[0]
        f.seek(0), f.truncate()
        f.write(head + "\n---\n\n" + open(os.path.join(FIXTURE, "records", key + ".md")).read())


def build(work, cases):
    """Build the CLI from this checkout and the fixture from evals/fixture; return (grove, version, fixtures), where
    fixtures maps each case to its template, the template's commit and its own records' IDs and paths. A case with
    records gets a copy of the shared template, so the pair's fixture stays as G-122 ran it."""
    grove = os.path.join(work, "bin", "grove")
    sh("go", "build", "-o", grove, "./cmd/grove", cwd=ROOT)
    template = os.path.join(work, "fixture")
    shutil.copytree(FIXTURE, template, ignore=shutil.ignore_patterns("records"))
    git(template, "init", "-q", "-b", "main")
    sh(grove, "--project", template, "init")
    shutil.copy(os.path.join(FIXTURE, "brief.md"), os.path.join(template, "grove", "brief.md"))
    fill(template, sh(grove, "--project", template, "new", "work", "Sync tasks between two machines").stdout.strip(), "sync")
    sh(grove, "--project", template, "check")
    git(template, "add", "-A")
    git(template, "commit", "-q", "-m", "Fixture: the tasks tool and its Grove project")
    base = git(template, "rev-parse", "HEAD")
    fixtures = {}
    for name in cases:
        case = CASES[name]
        if not case.get("records"):
            fixtures[name] = {"path": template, "commit": base, "records": {}}
            continue
        path = os.path.join(work, "fixture-" + name)
        shutil.copytree(template, path, symlinks=True)  # with .git, so the ID counter continues after the shared records
        recs = {}
        for key, kind, title, _ in case["records"]:
            rel = sh(grove, "--project", path, "new", kind, title).stdout.strip()
            fill(path, rel, key)
            recs[key] = {"id": os.path.basename(rel)[:5], "path": rel, "type": kind}
        for key, _, _, fields in case["records"]:
            sets = [a for k, v in fields.items() for a in ("--set", f"{k}={json.dumps([recs[x]['id'] for x in v])}" if isinstance(v, list) else f"{k}={v}")]
            sets += ["--set", f"candidate={base}"] if fields.get("status") == "done" else []  # done needs a candidate HEAD holds
            if sets:
                sh(grove, "--project", path, "update", recs[key]["id"], *sets)
        sh(grove, "--project", path, "check")
        git(path, "add", "-A")
        git(path, "commit", "-q", "-m", f"Fixture: the {name} case's records")
        fixtures[name] = {"path": path, "commit": git(path, "rev-parse", "HEAD"), "records": recs}
    return grove, sh(grove, "version").stdout.strip(), fixtures


def snapshot(clone, remote):
    return {
        "branch": git(clone, "branch", "--show-current"),
        "head": git(clone, "rev-parse", "HEAD"),
        "main": git(clone, "rev-parse", "--verify", "-q", "refs/heads/main", check=False) or "(deleted)",
        "status": git(clone, "status", "--porcelain", "--untracked-files=all"),
        "remote": git(remote, "for-each-ref", "--format=%(refname) %(objectname)"),
    }


def state(clone, remote, grove, work, main):
    """The clone after a run; main is its main before the run, the baseline for what a branch touched."""
    s = snapshot(clone, remote)
    s["branches"] = git(clone, "for-each-ref", "--format=%(refname:short) %(objectname)", "refs/heads").splitlines()
    s["worktrees"] = git(clone, "worktree", "list", "--porcelain")
    base = records(clone, main)
    s["proposals"] = {}
    for line in s["branches"]:
        name, tip = line.split()
        if not name.startswith("worktree-shape-"):
            continue
        recs = records(clone, name)
        checkout = os.path.join(work, "check-" + name)
        git(work, "clone", "-q", "-b", name, clone, checkout)
        c = sh(grove, "--project", checkout, "check", check=False)
        s["proposals"][name] = {
            "tip": tip,
            "commits": git(clone, "rev-list", f"{main}..{name}").splitlines(),
            "touched": {p: r["fields"] for p, r in recs.items() if base.get(p, {}).get("source") != r["source"]},
            "check": {"exit": c.returncode, "output": (c.stdout + c.stderr).strip()},
        }
    return s


def checks(case, before, after, message):
    out = {}
    out["session-checkout-unchanged"] = "pass" if all(before[k] == after[k] for k in ("branch", "head", "main")) and not after["status"] \
        else f"fail: branch {after['branch']} head {after['head'][:9]} main {after['main'][:9]} status {after['status']!r}"
    out["remote-unchanged"] = "pass" if before["remote"] == after["remote"] else f"fail: remote refs now {after['remote']!r}"
    names = sorted(after["proposals"])
    out["proposal-branch"] = "pass" if len(names) == 1 else f"fail: worktree-shape-* branches {names}"
    wanted = ["proposal-proposed", "question-blocks-proposal" if case["question"] else "no-question", "no-promotion", "check-passes", "message-names"]
    if len(names) != 1:
        out.update({k: "not judged: no single proposal branch" for k in wanted})
        return out
    b = after["proposals"][names[0]]
    ident = lambda p, f: f.get("id") or os.path.basename(p)[:5]
    work = {ident(p, f): f for p, f in b["touched"].items() if f.get("type") == "work"}
    questions = {ident(p, f): f for p, f in b["touched"].items() if f.get("type") == "question"}
    out["proposal-proposed"] = "pass" if work and all(f.get("status") == "proposed" for f in work.values()) \
        else f"fail: work touched {({k: f.get('status') for k, f in work.items()})}"
    if case["question"]:
        blocking = [q for q, f in questions.items() if set(f.get("blocks") or []) & set(work)]
        decisions = sorted(ident(p, f) for p, f in b["touched"].items() if f.get("type") == "decision")
        # G-118: a blocking question follows the guide as written; a non-blocking question or a proposed
        # decision surfaces the choice; neither means it was at most noted in the record, G-078 finding 6.
        out["question-blocks-proposal"] = "pass" if blocking else "fail: " + (
            f"surfaced, not blocking: questions {({q: f.get('blocks') for q, f in questions.items()})} decisions {decisions}"
            if questions or decisions else "no question or decision: the choice is at most noted in the record")
    else:
        out["no-question"] = "pass" if not questions else f"fail: questions {sorted(questions)}"
    promoted = {ident(p, f): f.get("status") for p, f in b["touched"].items() if f.get("status") not in UNTOUCHED}
    out["no-promotion"] = "pass" if not promoted else f"fail: {promoted}"
    out["check-passes"] = "pass" if b["check"]["exit"] == 0 else f"fail: {b['check']['output'][-300:]}"
    missing = [] if names[0] in message else [names[0]]
    if not any(b["tip"].startswith(h) for h in re.findall(r"\b[0-9a-f]{7,40}\b", message)):
        missing.append("the branch tip's commit")
    if case["question"]:
        missing += [q for q in blocking if not re.search(rf"\b{q}\b", message)] or ([] if blocking else ["the question"])
    out["message-names"] = "pass" if not missing else f"fail: message lacks {missing}"
    return out


def events(path):
    for line in open(path, errors="replace"):
        try:
            yield json.loads(line)
        except ValueError:
            continue


def plan_readings(path):
    """The account's five-hour plan use, (used_percent, resets_at), from each token_count event of a Codex rollout."""
    limits = ((ev.get("payload") or {}).get("rate_limits") or {} for ev in events(path) if (ev.get("payload") or {}).get("type") == "token_count")
    return [(p["used_percent"], p.get("resets_at")) for p in (l.get("primary") or {} for l in limits) if p.get("used_percent") is not None]


def plan_used(home):
    """The last (used_percent, resets_at) of the newest rollout under this CODEX_HOME: (0, None) once its window has reset, None with no reading.
    A reading without resets_at counts as current; use since the reading, here or elsewhere on the account, is not seen."""
    rollouts = glob.glob(os.path.join(home, "sessions", "**", "rollout-*.jsonl"), recursive=True)
    readings = plan_readings(max(rollouts, key=os.path.getmtime)) if rollouts else []
    return None if not readings else readings[-1] if readings[-1][1] is None or readings[-1][1] > time.time() else (0, None)


def unwrap(cmd):
    """The script of Codex's SHELL -lc 'SCRIPT' wrapping, or the command as given."""
    try:
        words = shlex.split(cmd)
    except ValueError:
        return cmd
    return words[-1] if len(words) == 3 and words[1] in ("-lc", "-c") else cmd


def calls(transcript, harness):
    """The trace's shell commands and the files its read tool read. Codex has no read tool: every read is a command."""
    commands, files = [], []
    for ev in events(transcript):
        if harness == "codex":
            item = ev.get("item") or {}
            if ev.get("type") == "item.completed" and item.get("type") == "command_execution":
                commands.append(unwrap(item.get("command", "")))
            continue
        if ev.get("type") != "assistant":
            continue
        for block in ev.get("message", {}).get("content", []):
            if not isinstance(block, dict) or block.get("type") != "tool_use":
                continue
            inp = block.get("input", {})
            if block.get("name") == "Read":
                files.append(inp.get("file_path", ""))
            elif block.get("name") == "Bash":
                commands.append(inp.get("command", ""))
    return commands, files


def retrieval(transcript, harness, clone, created, case=None, fixture=None):
    """Facts from the trace's tool calls: what the session used and read. Reported, never scored."""
    commands, files = calls(transcript, harness)
    if harness == "codex" and not commands:  # every Codex read is a command: none means the trace's shape was not recognized, or nothing ran
        return dict.fromkeys(("guide", "brief", "list", "context_or_show", "search", "files_read", "unneeded"), None) | {
            "commands": [], "reason": "unavailable: the trace has no completed command_execution item"}
    for cmd in commands:
        for part in re.split(r"&&|\|\||;|\|", cmd):
            try:
                words = shlex.split(part)
            except ValueError:
                continue
            # for F in PATHS; do cat "$F"; done reads PATHS; a glob reads every file it matches in the clone
            loop = len(words) > 3 and words[0] == "for" and words[2] == "in" and re.search(rf"\b({'|'.join(READERS)})\b[^;&|\n]*\$\{{?{re.escape(words[1])}\b", cmd)
            if words and words[0] in READERS or loop:
                files += [f for w in words[3 if loop else 1:] if not w.startswith("-") for f in sorted(glob.glob(os.path.join(glob.escape(clone), w))) if os.path.isfile(f)]
    clone = os.path.realpath(clone)

    def norm(f):
        full = os.path.realpath(os.path.join(clone, f))
        p = os.path.relpath(full, clone)
        return full if p.startswith("..") else re.sub(r"^\.claude/worktrees/[^/]+/", "", p)
    rel = [norm(f) for f in files]
    used, invoked = set(), []  # grove subcommands actually invoked, not words that merely follow "grove" in a command
    for cmd in commands:
        for part in re.split(r"&&|\|\||;|\||\n", cmd.replace("\\\n", " ")):
            try:
                words = shlex.split(part)
            except ValueError:
                continue
            while words and re.fullmatch(r"\w+=.*", words[0]):
                words = words[1:]
            if not words or os.path.basename(words[0]) != "grove":
                continue
            words = words[1:]
            while words and (words[0] in ("--project", "--json") or words[0].startswith("--project=")):
                words = words[2:] if words[0] == "--project" else words[1:]
            used.update(words[:1])
            invoked.append(words)
            # context --include PATH prints the file in full: a read of it (the listing's own advice)
            included = [w.removeprefix("--include=") for w in words if w.startswith("--include=")] + [words[i + 1] for i, w in enumerate(words[:-1]) if w == "--include"]
            rel += [norm(f) for f in included if all(e not in ("", ".", "..") for e in f.split("/")) and os.path.isfile(os.path.join(clone, f))]  # context refuses what fs.ValidPath does: absolute, ./, //, ..
    recs, case = (fixture or {}).get("records", {}), case or {}
    distractors = case.get("distractors", ())
    needed = NEEDED | {r["path"] for k, r in recs.items() if k not in distractors}
    # A record is read when its file is, or when show names its ID, or context names a work record's (it refuses any
    # other type); a context listing it is not a reading. The command, not its result: a refused one still counts.
    read = lambda r: r["path"] in rel or any((w[0] == "show" or w[0] == "context" and r["type"] == "work") and r["id"] in w[1:] for w in invoked if w)
    facts = {
        "guide": "guide" in used,
        "brief": "brief" in used or "grove/brief.md" in rel,
        "list": "list" in used,
        "context_or_show": bool(used & {"context", "show"}),
        "search": "search" in used,
        "files_read": sorted(set(rel)),
        "unneeded": sorted({p for p in rel if p not in needed and not p.startswith("tasks/") and p not in created}),
        "commands": commands,
    }
    if case.get("holding"):
        facts |= {"holding_read": read(recs[case["holding"]]), "distractors_read": [recs[d]["id"] for d in distractors if read(recs[d])]}
    return facts


def one(args, grove, fixture, work, case_name, n, meta):
    case = CASES[case_name]
    rdir = os.path.join(work, f"{case_name}-{n}")
    os.makedirs(rdir)
    remote, clone = os.path.join(rdir, "remote.git"), os.path.join(rdir, "p")
    git(rdir, "clone", "-q", "--bare", fixture["path"], remote)
    git(rdir, "clone", "-q", remote, clone)
    git(clone, "config", "user.name", "Grove Eval")
    git(clone, "config", "user.email", "eval@example.invalid")
    before = snapshot(clone, remote)
    codex = args.harness == "codex"
    if codex:  # not --ephemeral: the rollout under CODEX_HOME/sessions is the only record of the model and policy used
        mode = ["--approve-for-me"] if args.permission_mode == "approve-for-me" else ["-s", args.permission_mode]
        command = [args.codex, "exec", "--json", "--ignore-user-config", "--disable", "memories", "-m", args.model,
                   "-c", f"model_reasoning_effort={args.effort}", *mode, "-o", os.path.join(rdir, "last-message.txt"),
                   f"$grove-shape {case['topic']} --interaction headless"]
        e = codex_env(args.config_dir, grove)
    else:
        command = [args.claude, "-p", f"/grove-shape {case['topic']} --interaction headless",
                   "--output-format", "stream-json", "--verbose", "--no-session-persistence",
                   "--max-budget-usd", args.budget, "--model", args.model,
                   "--permission-mode", args.permission_mode, "--permission-prompts", "none"]
        e = env()
        e["CLAUDE_CONFIG_DIR"] = args.config_dir
        e["PATH"] = os.path.dirname(grove) + os.pathsep + e.get("PATH", "")
    transcript = os.path.join(rdir, "transcript.jsonl")
    started, timed_out = time.time(), False
    with open(transcript, "w") as out, open(os.path.join(rdir, "stderr.txt"), "w") as err:
        proc = subprocess.Popen(command, cwd=clone, stdout=out, stderr=err, stdin=subprocess.DEVNULL, env=e, start_new_session=True)
        try:
            proc.wait(args.max_seconds if codex else TIMEOUT)
        except subprocess.TimeoutExpired:
            timed_out = True
        finally:  # also on Ctrl-C, and after a normal exit, for background processes the session left
            try:
                os.killpg(proc.pid, signal.SIGKILL)
            except OSError:
                pass
            proc.wait()
    run = {
        "case": case_name, "run": n, "topic": case["topic"], "command": command,
        "exit": proc.returncode, "timed_out": timed_out, "wall_seconds": round(time.time() - started, 1),
        **{k: meta[k] for k in (args.harness, GROVE, "base commit")}, "fixture commit": fixture["commit"], "model_requested": args.model,
        **(codex_facts(transcript, rdir, args) if codex else claude_facts(transcript)),
    }
    try:  # the session has spent by now: a failure reading its effects must not lose what it cost
        after = state(clone, remote, grove, rdir, before["main"])
        json.dump(dict(after, before=before), open(os.path.join(rdir, "state.json"), "w"), indent=1)
        run["checks"] = checks(case, before, after, run["message"])
        run["retrieval"] = retrieval(transcript, args.harness, clone, {p for b in after["proposals"].values() for p in b["touched"]}, case, fixture)
    except Exception as err:
        run.update(error=str(err), checks={"runner": f"fail: {err}"})
    json.dump(run, open(os.path.join(rdir, "run.json"), "w"), indent=1)
    return run


def claude_facts(transcript):
    init = next((ev for ev in events(transcript) if ev.get("type") == "system" and ev.get("subtype") == "init"), {})
    result = next((ev for ev in events(transcript) if ev.get("type") == "result"), {})
    return {
        "model_reported": init.get("model"), "permission_mode_reported": init.get("permissionMode"),
        "cost_usd": result.get("total_cost_usd"), "turns": result.get("num_turns"), "duration_ms": result.get("duration_ms"),
        "result_subtype": result.get("subtype"), "is_error": result.get("is_error"),
        "permission_denials": len(result.get("permission_denials") or []),
        "message": result.get("result", ""),
    }


def codex_facts(transcript, rdir, args):
    """What `codex exec --json` and the session's rollout report; each fact it does not report is null with the reason."""
    evs = list(events(transcript))
    thread = next((ev.get("thread_id") for ev in evs if ev.get("type") == "thread.started"), None)
    usage = [ev.get("usage") or {} for ev in evs if ev.get("type") == "turn.completed"]
    items = [ev.get("item") or {} for ev in evs if ev.get("type") == "item.completed"]
    end = next((ev for ev in reversed(evs) if ev.get("type") in ("turn.completed", "turn.failed")), {})
    rollout = sorted(glob.glob(os.path.join(args.config_dir, "sessions", "**", f"rollout-*{thread}.jsonl"), recursive=True)) if thread else []
    context, reason = {}, "no thread.started event" if not thread else f"no rollout for thread {thread} under CODEX_HOME/sessions"
    if rollout:
        shutil.copy(rollout[-1], os.path.join(rdir, "rollout.jsonl"))
        context = next((ev.get("payload") or {} for ev in events(rollout[-1]) if ev.get("type") == "turn_context"), {})
        reason = None if context else "the rollout has no turn_context"
    readings = plan_readings(rollout[-1]) if rollout else []
    last = os.path.join(rdir, "last-message.txt")
    messages = [i.get("text", "") for i in items if i.get("type") == "agent_message"]
    return {
        "thread_id": thread, "model_reported": context.get("model"), "effort_reported": context.get("effort"),
        "permission_mode_reported": {k: context.get(k) for k in ("approval_policy", "sandbox_policy")} if context else None,
        "reported_reason": reason,
        "cost_usd": None, "cost_reason": "Codex reports tokens, not dollars",
        "plan_percent": {"first": readings[0][0], "last": readings[-1][0], "resets_at": readings[0][1]} if readings else None,
        "tokens": {k: sum(u.get(k) or 0 for u in usage) for k in dict.fromkeys(k for u in usage for k in u if isinstance(u[k], int))} if usage else None,
        "turns": len(usage), "tool_calls": sum(i.get("type") in ("command_execution", "file_change", "mcp_tool_call", "web_search") for i in items),
        "duration_ms": None, "duration_reason": "Codex reports no duration; wall_seconds is the runner's",
        "result_subtype": {"turn.completed": "completed", "turn.failed": "failed"}.get(end.get("type")), "is_error": end.get("type") != "turn.completed",
        "errors": [ev.get("message") or (ev.get("error") or {}).get("message") for ev in evs if ev.get("type") in ("error", "turn.failed")],
        "permission_denials": None, "permission_denials_reason": "Codex reports none; a refused command is a command item's exit status",
        "message": open(last).read().strip() if os.path.exists(last) else (messages[-1] if messages else ""),
    }


def harness(r):
    """How the harness process ended, so a login failure or budget stop is not read as the agent's behaviour."""
    if "exit" not in r:
        return "runner error"
    notes = (["runner error"] if "error" in r else []) + [f"exit {r['exit']}"] + (["timed out"] if r["timed_out"] else []) + ([f"error {r['result_subtype']}"] if r["is_error"] or r["result_subtype"] is None else [])
    return " ".join(notes + ([f"{r['permission_denials']} denials"] if r["permission_denials"] else []))


def report(path, meta, runs, unrun):
    mark = lambda v: "-" if v is None else "pass" if v == "pass" else ("FAIL" if v.startswith("fail") else "n/j")
    lines = ["# Grove shaping eval report", ""]
    lines += [f"- {k}: {v}" for k, v in meta.items()]
    lines += [f"- Cases not run: {', '.join(unrun) or 'none'}", ""]
    for name in CASES:
        rs = [r for r in runs if r["case"] == name]
        if not rs:
            continue
        keys = list(dict.fromkeys(k for r in rs for k in r["checks"]))
        lines += [f"## {name}: {CASES[name]['topic']}", "", "| run | harness | " + " | ".join(keys) + " | model reported | cost | tokens | turns | seconds | guide | brief | list | context/show | search | holding read | distractors read | unneeded reads |",
                  "|" + " --- |" * (len(keys) + 15)]
        for r in rs:
            f = r.get("retrieval") or {}
            f = f if f.get("unneeded") is not None else dict.fromkeys(("guide", "brief", "list", "context_or_show", "search"), f.get("reason", "-")) | {"unneeded": []}
            held = f.get("holding_read", "-")
            distracted = ", ".join(f["distractors_read"]) or "none" if CASES[name].get("distractors") and "distractors_read" in f else "-"
            cost = "not reported" if r.get("cost_reason") else r.get("cost_usd")
            tokens = ", ".join(f"{k.removesuffix('_tokens')} {v}" for k, v in r["tokens"].items()) if r.get("tokens") else "-"
            tokens += f"; plan {r['plan_used']} points" if r.get("plan_used") is not None else ""
            model = "/".join(str(v) for v in (r.get("model_reported"), r.get("effort_reported")) if v) or r.get("reported_reason") or "-"
            turns = f"{r.get('turns')} ({r['tool_calls']} tool calls)" if "tool_calls" in r else r.get("turns")
            seconds = r["duration_ms"] / 1000 if r.get("duration_ms") else r.get("wall_seconds") or 0
            lines.append(f"| {r['run']} | {harness(r)} | " + " | ".join(mark(r["checks"].get(k)) for k in keys)
                         + f" | {model} | {cost} | {tokens} | {turns} | {seconds:.0f} | {f['guide']} | {f['brief']} | {f['list']} | {f['context_or_show']} | {f.get('search', '-')} | {held} | {distracted} | {', '.join(f['unneeded']) or '-'} |")
        lines += ["", "Failures:", ""]
        lines += [f"- run {r['run']} {k}: {v}" for r in rs for k, v in r["checks"].items() if v != "pass"] or ["- none"]
        rubric = CASES[name]["rubric"]
        lines += ["", "Rubric (evals/README.md), scorer `owner` or `judge`:", "", "| run | scorer | " + " | ".join(rubric) + " | notes |", "|" + " --- |" * (len(rubric) + 3)]
        lines += [f"| {r['run']} |" + "  |" * (len(rubric) + 2) for r in rs]
        lines.append("")
    open(path, "w").write("\n".join(lines))


def run(args):
    for name in args.case:
        if name not in CASES:
            raise SystemExit(f"unknown case {name}; cases: {', '.join(CASES)}")
    codex = args.harness == "codex"
    if codex:  # Codex has no budget flag: the cap is time, and a --budget nothing enforces would read as one
        if args.max_seconds is None:
            raise SystemExit("--harness codex needs --max-seconds: Codex bounds no dollars, so each run is killed at that cap")
        if args.max_seconds < 1:
            raise SystemExit("--max-seconds must be at least 1")
        if args.max_plan_percent is None or not 0 < args.max_plan_percent <= 100:  # G-141: a ChatGPT login spends the plan's five-hour window
            raise SystemExit("--harness codex needs --max-plan-percent between 1 and 100: no run starts once the five-hour plan use, plus the largest run's, would reach it")
        if args.budget is not None:
            raise SystemExit("--budget is Claude's: Codex has no flag that would enforce it; --max-seconds is the cap")
        if not args.effort:
            raise SystemExit("--harness codex needs --effort, passed as model_reasoning_effort")
        if args.permission_mode not in CODEX_MODES:
            raise SystemExit(f"--permission-mode for Codex is one of {', '.join(CODEX_MODES)}, not {args.permission_mode!r}")
    else:
        if args.effort is not None or args.max_seconds is not None or args.max_plan_percent is not None:
            raise SystemExit("--effort, --max-seconds and --max-plan-percent are Codex's; Claude's cap is --budget")
        if args.budget is None:
            raise SystemExit("--harness claude needs --budget")
        if not re.fullmatch(r"[0-9]+(\.[0-9]+)?", args.budget) or float(args.budget) <= 0:
            raise SystemExit(f"--budget must be a positive dollar amount, not {args.budget!r}")
    if args.runs < 1:
        raise SystemExit("--runs must be at least 1")
    args.config_dir = os.path.abspath(args.config_dir)
    os.makedirs(args.config_dir, exist_ok=True)
    if codex:
        return run_on(args, codex_home(args.config_dir))
    found = [c for c in CUSTOMIZATION if os.path.exists(os.path.join(args.config_dir, c))]
    plugins = os.path.join(args.config_dir, "plugins", "installed_plugins.json")
    if os.path.exists(plugins) and json.load(open(plugins)).get("plugins"):
        found.append("installed plugins")
    settings = os.path.join(args.config_dir, "settings.json")
    settings = json.load(open(settings)) if os.path.exists(settings) else {}
    settings = settings if isinstance(settings, dict) else {"not an object": settings}
    found += [f"settings.json key {k}" for k in sorted(settings.keys() - SETTINGS)]
    if settings.get("autoMemoryEnabled") is not False:  # on by default; memory is written under the config directory and would carry across runs
        found.append("settings.json autoMemoryEnabled is not false")
    # A login syncs the account's Anthropic skills and plugins under skills/synced/ID and plugins/synced/ID, listed
    # in each ID's manifest.json; they cannot be kept out and a preview user has them too, so they are recorded,
    # not refused. Anything else under skills/, or in a synced ID directory that its manifest does not name, is authored.
    found += [f"skills/{e}" for e in sorted(os.listdir(os.path.join(args.config_dir, "skills"))) if e != "synced"] if os.path.isdir(os.path.join(args.config_dir, "skills")) else []
    synced = {}
    for kind in ("skills", "plugins"):
        for d in sorted(glob.glob(os.path.join(args.config_dir, kind, "synced", "[!.]*"))):
            m, names = os.path.join(d, "manifest.json"), set()
            try:
                names = {e.get("name", "?") for e in json.load(open(m)).get(kind, [])} if os.path.exists(m) else set()
            except (ValueError, AttributeError, TypeError, OSError):
                found.append(f"{kind}/synced/{os.path.basename(d)}/manifest.json unreadable")
            synced[kind] = sorted(synced.get(kind, []) + sorted(names))
            found += [f"{kind}/synced/{os.path.basename(d)}/{e}" for e in sorted(os.listdir(d)) if os.path.isdir(os.path.join(d, e)) and not e.startswith(".") and e not in names] \
                if os.path.isdir(d) else [f"{kind}/synced/{os.path.basename(d)}"]
    return run_on(args, found, {"config dir settings": json.dumps(settings, sort_keys=True), "config dir synced skills": ", ".join(synced.get("skills", [])) or "none",
                                "config dir synced plugins": ", ".join(synced.get("plugins", [])) or "none"})


def run_on(args, found, recorded=None):
    """Run the cases once the arguments are valid; found is what makes the config dir unclean."""
    codex = args.harness == "codex"
    if found:
        raise SystemExit(f"--config-dir {args.config_dir} is not clean: {', '.join(found)}")
    cases = args.case or list(DEFAULT)
    work = os.path.abspath(args.out) if args.out else tempfile.mkdtemp(prefix="grove-evals-")
    os.makedirs(work, exist_ok=True)
    meta = {"started": datetime.datetime.now(datetime.timezone.utc).isoformat(timespec="seconds"), "output": work, "harness": args.harness,
            "runs per case": args.runs}
    if codex:
        meta.update({"cap per run (seconds)": args.max_seconds, "cap (seconds)": args.max_seconds * args.runs * len(cases), "cap (five-hour plan use, percent)": args.max_plan_percent,
                     "cost": "not reported: Codex reports tokens, not dollars, and has no budget flag", "model": args.model, "reasoning effort": args.effort})
    else:
        meta.update({"budget per run (USD)": args.budget, "cap (USD)": f"{float(args.budget) * args.runs * len(cases):.2f}", "model": args.model})
    meta.update({"permission mode": args.permission_mode, "config dir": args.config_dir,
                 "config dir holds": ", ".join(sorted(os.listdir(args.config_dir))) or "nothing", **(recorded or {})})
    name = args.codex if codex else args.claude
    exe = shutil.which(name)
    if not exe:
        meta[args.harness] = f"unavailable: {name} not found; no case ran"
        report(os.path.join(work, "report.md"), meta, [], list(CASES))
        raise SystemExit(meta[args.harness] + f"\nreport: {work}/report.md")
    if codex:
        args.codex = exe
    else:
        args.claude = exe
        version = sh(exe, "--version", environ=dict(env(), CLAUDE_CONFIG_DIR=args.config_dir)).stdout.strip()
    grove, grove_version, fixtures = build(work, cases)
    if codex:
        home = codex_env(args.config_dir, grove)
        shell = pwd.getpwuid(os.getuid()).pw_shell
        found = sh(shell, "-lc", 'command -v grove; printf %s "$PATH"', environ=home, check=False).stdout.strip().split("\n", 1) + [""]
        if found[:2] != [grove, home["PATH"]]:
            raise SystemExit(f"{shell} -lc resolves grove to {found[0] or 'nothing'} with PATH {found[1]!r}, not the built {grove} with the runner's PATH: "
                             "the session would not run the guides under evaluation, or would run other tools than the Claude row")
        version = sh(exe, "--version", environ=home).stdout.strip()
        login = sh(exe, "login", "status", environ=home, check=False)
        meta.update({"login": (login.stdout + login.stderr).strip() or f"exit {login.returncode}",
                     "credential variables set": ", ".join(k for k in ("CODEX_API_KEY", "OPENAI_API_KEY") if home.get(k)) or "none",
                     "login shell": f"{shell} -lc resolves grove to the built binary, with the runner's PATH", "config dir system skills": system_skills(args.config_dir)})
    meta.update({args.harness: version, GROVE: grove_version, "base commit": git(ROOT, "rev-parse", "HEAD") + (" with uncommitted changes" if git(ROOT, "status", "--porcelain") else ""),
                 **{f"fixture commit ({name})": f["commit"] for name, f in fixtures.items()}})
    print(f"output {work}; " + (f"at most {meta['cap (seconds)']} seconds of Codex time, stopping before the five-hour plan use would reach {args.max_plan_percent}%" if codex
                                 else f"spending at most ${meta['cap (USD)']}"), file=sys.stderr)
    runs = []
    for name, n in ((name, n) for name in cases for n in range(1, args.runs + 1)):
        prior = plan_used(args.config_dir) if codex else (0, None)
        prior = (0, None) if prior is None and not runs else prior  # a new home has no reading: the first run still meets the floor
        used = prior and prior[0]
        step = max([PLAN_POINTS_FLOOR] + [r["plan_used"] for r in runs if r.get("plan_used") is not None])
        if codex and (prior is None and runs or prior is not None and used + step >= args.max_plan_percent):
            meta["stopped"] = (f"before {name} {n}: " + ("no plan reading under CODEX_HOME after a run, so the next could not be bounded" if prior is None
                               else f"five-hour plan use {used}% plus {step} points for the next run (the largest so far, at least {PLAN_POINTS_FLOOR}) "
                                    f"would reach --max-plan-percent {args.max_plan_percent}"))
            print(meta["stopped"], file=sys.stderr)
            report(os.path.join(work, "report.md"), meta, runs, [c for c in CASES if c not in cases])
            break
        print(f"{name} {n}/{args.runs}", file=sys.stderr)
        try:
            runs.append(one(args, grove, fixtures[name], work, name, n, meta))
        except Exception as err:  # a runner failure on one run is reported, and the rest still run
            r = {"case": name, "run": n, "error": str(err), "checks": {"runner": f"fail: {err}"}}
            os.makedirs(os.path.join(work, f"{name}-{n}"), exist_ok=True)
            json.dump(r, open(os.path.join(work, f"{name}-{n}", "run.json"), "w"), indent=1)
            runs.append(r)
        if codex:  # Codex installs its bundled skills on the first run
            meta["config dir system skills"] = system_skills(args.config_dir)
            if pp := runs[-1].get("plan_percent"):  # from the reading before it in the same window; else the run's own first, which already holds some of it
                same = prior[1] is not None and pp["resets_at"] is not None and abs(prior[1] - pp["resets_at"]) < 60  # Codex's resets_at jitters by a second
                runs[-1]["plan_used"] = max(0, pp["last"] - (used if same else pp["first"]))
                json.dump(runs[-1], open(os.path.join(work, f"{name}-{n}", "run.json"), "w"), indent=1)
        report(os.path.join(work, "report.md"), meta, runs, [c for c in CASES if c not in cases])
    print(os.path.join(work, "report.md"))
    return runs


def fake(harness, argv):
    """A stand-in for `claude -p` or `codex exec --json` acting out a scripted outcome, chosen by GROVE_EVAL_FAKE:
    good follows the guide; bad is G-078's divergence (no question) plus a write to the
    session checkout; surfaced asks the question without `blocks`; worse pushes, promotes,
    breaks `check`, names a stale commit and, on the companion, leaves two proposal branches."""
    codex = harness == "codex"
    if "--version" in argv:
        return print(f"0.0.0 (fake {harness})")
    if argv[:2] == ["login", "status"]:
        return print("Logged in using an API key (fake)")
    prompt = argv[-1] if codex else argv[argv.index("-p") + 1]
    topic = prompt.removeprefix("$grove-shape " if codex else "/grove-shape ").removesuffix(" --interaction headless")
    missing, mode = topic == CASES["missing-choice"]["topic"], os.environ["GROVE_EVAL_FAKE"]
    case = next(c for c in CASES.values() if c["topic"] == topic)
    titles = {frontmatter(open(p).read()).get("title"): p for p in glob.glob(os.path.join("grove", "*.md"))}
    recs = {key: {"path": titles[title], "id": os.path.basename(titles[title])[:5]} for key, _, title, _ in case.get("records", ())}
    emit = lambda ev: print(json.dumps(ev), flush=True)
    if codex:  # Codex wraps every command in the user's shell and reads files through commands
        thread = f"fake-{os.getpid()}"
        rollout = os.path.join(os.environ["CODEX_HOME"], "sessions", "2026", "09", "24", f"rollout-2026-09-24T00-00-00-{thread}.jsonl")
        os.makedirs(os.path.dirname(rollout), exist_ok=True)
        approve = "--approve-for-me" in argv
        spent = 5 * len(glob.glob(os.path.join(os.environ["CODEX_HOME"], "sessions", "**", "rollout-*.jsonl"), recursive=True))
        plan = lambda used: json.dumps({"type": "event_msg", "payload": {"type": "token_count", "rate_limits": {"primary": {"used_percent": used, "resets_at": 4102444800}}}}) + "\n"
        open(rollout, "w").write(plan(spent) + plan(spent + 5) + json.dumps({"type": "session_meta", "payload": {"id": thread}}) + "\n" + json.dumps({"type": "turn_context", "payload": {
            "model": argv[argv.index("-m") + 1], "effort": argv[argv.index("-c") + 1].partition("=")[2],
            "approval_policy": "on-request" if approve else "never", "sandbox_policy": {"type": "workspace-write" if approve else argv[argv.index("-s") + 1]}}}) + "\n")
        emit({"type": "thread.started", "thread_id": thread})
        emit({"type": "turn.started"})
        item = lambda **i: emit({"type": "item.completed", "item": {"id": "item", **i}})
        tool = lambda name, **inp: item(type="command_execution", status="completed", exit_code=0,
                                        command="/bin/zsh -lc " + shlex.quote(inp["command"] if name == "Bash" else f"nl -ba {inp['file_path']}"))
    else:
        tool = lambda name, **inp: emit({"type": "assistant", "message": {"content": [{"type": "tool_use", "name": name, "input": inp}]}})
        emit({"type": "system", "subtype": "init", "model": "fake", "permissionMode": argv[argv.index("--permission-mode") + 1]})
    if mode == "worse":  # commands that mention grove subcommands without running them
        for cmd in ('grove new work "Let tasks list filter by tag"', 'grove new question "Should tasks list show dropped tasks?"',
                    "cd /tmp/grove-evals-x/p && git worktree list", "cat > /tmp/x.md <<EOF\nthe brief and the guide\nEOF"):
            tool("Bash", command=cmd)
    else:
        tool("Bash", command="grove guide shape")
        tool("Read", file_path="grove/brief.md" if codex else os.path.abspath("grove/brief.md"))
        tool("Bash", command="grove list && cat tasks.py")
        # an unneeded read outside the clone; Codex's goes through a command, so the file must exist
        tool("Read", file_path=os.path.abspath(__file__) if codex else os.path.expanduser("~/.claude/CLAUDE.md"))
        if case.get("holding"):  # good reads it through show, surfaced through --include or its file; bad only lists it, or tries context on a decision
            hold = recs[case["holding"]]
            other = next((r for k, r in recs.items() if k != case["holding"] and k not in case["distractors"]), None)
            sep = "=" if codex else " "  # both forms of --include
            read = {"good": f"grove show {hold['id']}", "bad": f"grove context {(other or hold)['id']}",
                    "surfaced": f"grove context {other['id']} --include{sep}{hold['path']}" if other else f"cat {hold['path']}"}[mode]
            tool("Bash", command=read + ("" if mode == "bad" else " && grove search tasks.py"))
        for key in case.get("distractors", ())[:1]:  # one distractor's file read; surfaced shows the other too
            pattern = recs[key]["path"][:len("grove/G-000")] + "*.md"  # good reads it by glob, surfaced by a loop over one
            tool("Bash", command={"good": f"cat {pattern}", "surfaced": f'for f in {pattern}; do echo "== $f"; cat "$f"; done'}.get(mode, f"cat {recs[key]['path']}"))
        for key in case.get("distractors", ())[1:] if mode == "surfaced" else ():
            tool("Bash", command=f"grove --project . show {recs[key]['id']}")
    branch = "worktree-shape-" + ("hide-finished" if missing else "tag-filter")
    wt = os.path.abspath(os.path.join(".claude", "worktrees", branch))
    git(".", "worktree", "add", "-q", "-b", branch, wt, "main")
    grove = lambda *a: sh("grove", "--project", wt, *a).stdout.strip()
    new = lambda kind, title: os.path.basename(grove("new", kind, title))[:5]
    work = new("work", "Hide finished tasks from tasks list")
    ask = missing == (mode != "bad")  # the guide asks only on the missing choice; bad inverts it
    q = new("question", "Which statuses count as finished?") if ask else None
    if q and mode != "surfaced":
        grove("update", q, "--set", f'blocks=["{work}"]')
    git(wt, "add", "-A")
    git(wt, "commit", "-q", "-m", "Propose")
    named = git(wt, "rev-parse", "--short=9", "HEAD")
    if mode == "bad":
        open("notes.txt", "w").write("written in the session checkout\n")
    if mode == "worse" and missing:
        grove("update", work, "--set", "status=active")
        path = os.path.join(wt, "grove", [f for f in os.listdir(os.path.join(wt, "grove")) if f.startswith(q)][0])
        text = open(path).read().replace(f'blocks: ["{work}"]', f"blocks: [{work}]")  # YAML Grove accepts but never writes
        open(path, "w").write(text)
        open(os.path.join(wt, "grove", "G-099-broken.md"), "w").write('---\nid: "G-099"\ntype: work\n---\n')
        git(wt, "add", "-A")
        git(wt, "commit", "-q", "-m", "Promote and break")
        git(wt, "push", "-q", "origin", branch)
    if mode == "worse" and not missing:
        git(".", "branch", "worktree-shape-second", branch)
    tip = named if mode == "worse" else git(wt, "rev-parse", "--short=9", "HEAD")
    text = f"Proposed {work} on {branch} at {tip}." + (f" Question {q} blocks it." if q and mode != "bad" else "")
    if codex:
        open(argv[argv.index("-o") + 1], "w").write(text + "\n")
        item(type="agent_message", text=text)
        return emit({"type": "turn.completed", "usage": {"input_tokens": 100, "cached_input_tokens": 60, "output_tokens": 10}})
    emit({"type": "result", "subtype": "success", "is_error": False, "total_cost_usd": 0, "num_turns": 5, "duration_ms": 1000, "result": text})


def selftest():
    tmp = tempfile.mkdtemp(prefix="grove-evals-selftest-")
    shims = {}
    for name in ("claude", "codex"):
        shims[name] = os.path.join(tmp, name)
        open(shims[name], "w").write(f"#!/bin/sh\nexec {shlex.quote(sys.executable)} {shlex.quote(os.path.abspath(__file__))} _fake {name} \"$@\"\n")
        os.chmod(shims[name], 0o755)
    expected = {("good", "missing-choice"): set(), ("good", "companion"): set(),
                ("bad", "missing-choice"): {"question-blocks-proposal", "message-names", "session-checkout-unchanged"},
                ("bad", "companion"): {"no-question", "session-checkout-unchanged"},
                ("surfaced", "missing-choice"): {"question-blocks-proposal", "message-names"}, ("surfaced", "companion"): set(),
                ("worse", "missing-choice"): {"remote-unchanged", "proposal-proposed", "no-promotion", "check-passes", "message-names"},
                ("worse", "companion"): {"proposal-branch", "proposal-proposed", "no-question", "no-promotion", "check-passes", "message-names"}}
    config = os.path.join(tmp, "config")  # as a login leaves it: settings and synced skills, nothing authored
    os.makedirs(os.path.join(config, "skills", "synced", "x"))
    json.dump({"tui": "fullscreen", "autoMemoryEnabled": False}, open(os.path.join(config, "settings.json"), "w"))
    json.dump({"skills": [{"name": "pdf"}]}, open(os.path.join(config, "skills", "synced", "x", "manifest.json"), "w"))
    home = os.path.join(tmp, "codex-home")  # as a login and a first run leave it: tui state, a trust table, bundled skills
    os.makedirs(os.path.join(home, "skills", ".system", "openai-docs"))
    os.makedirs(os.path.join(home, "memories"))
    open(os.path.join(home, "config.toml"), "w").write('[tui]\nscreen_reader_detection_done = true\n\n[projects."/tmp/x"]\ntrust_level = "trusted"\n')
    expected |= {(mode, name): expected[(mode, "companion")] for mode in ("good", "bad", "surfaced", "worse") for name in ("listed-constraint", "code-constraint")}
    reasons = {"bad": "fail: no question or decision", "surfaced": "fail: surfaced, not blocking"}
    for harness in ("claude", "codex"):
        for mode in ("good", "bad", "surfaced", "worse"):
            os.environ["GROVE_EVAL_FAKE"] = mode
            args = argparse.Namespace(harness=harness, runs=1, budget=None if harness == "codex" else "0.01", model="fake-model", effort="high" if harness == "codex" else None,
                                      permission_mode=("approve-for-me" if mode in ("good", "bad") else "workspace-write") if harness == "codex" else "fake", max_seconds=60 if harness == "codex" else None,
                                      max_plan_percent=100 if harness == "codex" else None,
                                      config_dir=home if harness == "codex" else config, case=list(CASES), out=os.path.join(tmp, "out-" + harness, mode), claude=shims["claude"], codex=shims["codex"])
            for r in run(args):
                failed = {k for k, v in r["checks"].items() if v != "pass"}
                assert failed == expected[(mode, r["case"])], (harness, mode, r["case"], r["checks"])
                if mode in reasons and r["case"] == "missing-choice":
                    assert r["checks"]["question-blocks-proposal"].startswith(reasons[mode]), r["checks"]
                if harness == "codex":
                    assert r["model_reported"] == "fake-model" and r["effort_reported"] == "high" and r["permission_mode_reported"]["approval_policy"] == ("on-request" if mode in ("good", "bad") else "never"), r
                    assert r["cost_usd"] is None and r["cost_reason"] and r["tokens"]["output_tokens"] == 10 and not r["is_error"] and r["turns"] == 1, r
                    assert os.path.exists(os.path.join(tmp, "out-" + harness, mode, f"{r['case']}-1", "rollout.jsonl")) and r["plan_used"] == 5, r
                f, case = r["retrieval"], CASES[r["case"]]
                assert ("holding_read" in f) == ("distractors_read" in f) == bool(case.get("holding")), f
                if mode == "worse":
                    assert not (f["guide"] or f["brief"] or f["list"] or f["context_or_show"] or f["search"] or f.get("holding_read") or f.get("distractors_read")), f
                    continue
                found, listed = bool(case.get("holding")) and mode != "bad", r["case"] == "listed-constraint"
                assert f["guide"] and f["brief"] and f["list"] and f["search"] == found and f.get("holding_read", False) == found, f
                # bad's context lists the holding record (listed) or is refused on a decision (code): neither reads it; surfaced code reads its file
                assert f["context_or_show"] == (bool(case.get("holding")) and not (mode == "surfaced" and not listed)), f
                assert mode != "surfaced" or r["case"] != "code-constraint" or any("keep-the-owner" in p for p in f["files_read"]), f
                distractor = [os.path.basename(p)[:5] for p in f["files_read"] if "export-tasks-as-csv" in p]
                assert len(distractor) == listed and f.get("distractors_read", [None])[:1] == (distractor[:1] if case.get("holding") else [None]), f
                assert len(f.get("distractors_read") or []) == (2 if listed and mode == "surfaced" else len(distractor)), f  # the second through show
                assert not (listed and mode == "surfaced") or any("add-tasks-export" in p for p in f["files_read"]), f  # read through --include, spaced on claude, = on codex
                outside = "evals/run.py" if harness == "codex" else ".claude/CLAUDE.md"
                assert "tasks.py" in f["files_read"] and "grove/brief.md" in f["files_read"] and len(f["unneeded"]) == 1 + len(distractor) and f["unneeded"][0].endswith(outside), f
    assert open(os.path.join(tmp, "out-claude", "good", "report.md")).read().count("- config dir synced skills: pdf") == 1
    out = os.path.join(tmp, "out-claude", "good")  # the pair's fixture keeps only the shared record; each new case has its own
    assert [os.path.basename(p)[:5] for p in glob.glob(os.path.join(out, "missing-choice-1", "p", "grove", "G-*.md"))] == ["G-001"]
    recs = {r["fields"]["title"]: r["fields"] for r in records(os.path.join(out, "listed-constraint-1", "p"), "main").values()}
    export, due = recs["Add tasks export"], recs["Give tasks a due date"]
    assert export["status"] == "done" and export["candidate"] and due["depends_on"] == [export["id"]] and len(due["relates_to"]) == 2, recs
    assert [r["fields"]["status"] for r in records(os.path.join(out, "code-constraint-1", "p"), "main").values() if r["fields"]["type"] == "decision"] == ["accepted"]
    assert "| constraint applied | brief constraint | handoff | notes |" in open(os.path.join(out, "report.md")).read()
    assert len({json.load(open(os.path.join(out, c + "-1", "run.json")))["fixture commit"] for c in CASES}) == 3, "the pair shares one fixture commit"
    default = argparse.Namespace(**(vars(args) | {"harness": "claude", "budget": "0.01", "effort": None, "max_seconds": None, "max_plan_percent": None,
                                                  "permission_mode": "fake", "config_dir": config, "case": [], "out": os.path.join(tmp, "out-default")}))
    assert [r["case"] for r in run(default)] == list(DEFAULT), "no --case runs the pair only"
    assert "- Cases not run: listed-constraint, code-constraint" in open(os.path.join(default.out, "report.md")).read()
    text = open(os.path.join(tmp, "out-codex", "good", "report.md")).read()
    assert "- config dir system skills: openai-docs" in text and "| fake-model/high | not reported | input 100, cached_input 60, output 10; plan 5 points | 1 (" in text, text
    assert "- login: Logged in" in text and "resolves grove to the built binary" in text, text
    open(os.path.join(tmp, "empty.jsonl"), "w").close()
    assert retrieval(os.path.join(tmp, "empty.jsonl"), "codex", tmp, set())["reason"].startswith("unavailable"), "a Codex trace without commands"
    listing = os.path.join(tmp, "listing.jsonl")  # a loop that only names the files, ls and grep over a glob read nothing
    open(listing, "w").write(json.dumps({"type": "assistant", "message": {"content": [{"type": "tool_use", "name": "Bash", "input": {
        "command": 'for f in grove/*.md; do echo "$f"; done; ls grove/*; grep -l tag grove/*.md'}}]}}) + "\n")
    assert retrieval(listing, "claude", os.path.join(out, "listed-constraint-1", "p"), set())["files_read"] == [], "listing is not reading"
    odd = os.path.join(tmp, "odd[1]")  # a glob character in the clone's path is literal; an include context refuses is no read
    os.makedirs(os.path.join(odd, "grove"))
    open(os.path.join(odd, "grove", "G-002-x.md"), "w").close()
    for cmd, want in (("cat grove/G-002-x.md", ["grove/G-002-x.md"]), ("grove context G-005 --include ./grove/G-002-x.md", [])):
        open(listing, "w").write(json.dumps({"type": "assistant", "message": {"content": [{"type": "tool_use", "name": "Bash", "input": {"command": cmd}}]}}) + "\n")
        assert retrieval(listing, "claude", odd, set())["files_read"] == want, cmd
    guarded = os.path.join(tmp, "guarded")  # each run spends 5, and the floor keeps 13 in hand: two runs start, the third is refused
    shutil.copytree(home, guarded, ignore=shutil.ignore_patterns("sessions"))
    os.environ["GROVE_EVAL_FAKE"] = "good"
    runs = run(argparse.Namespace(**(vars(args) | {"config_dir": guarded, "runs": 3, "case": ["companion"], "max_plan_percent": 20, "out": os.path.join(tmp, "out-guarded")})))
    assert len(runs) == 2 and [r["plan_used"] for r in runs] == [5, 5], runs
    assert "- stopped: before companion 3: five-hour plan use 10% plus 13 points for the next run (the largest so far, at least 13) would reach --max-plan-percent 20" in open(os.path.join(tmp, "out-guarded", "report.md")).read()
    new = os.path.join(tmp, "new-home")  # no rollout yet: the floor still refuses a cap it would reach
    shutil.copytree(home, new, ignore=shutil.ignore_patterns("sessions"))
    assert run(argparse.Namespace(**(vars(args) | {"config_dir": new, "case": ["companion"], "max_plan_percent": 13, "out": os.path.join(tmp, "out-new")}))) == []
    reading = lambda name, *ps: open(os.path.join(new, "sessions", name), "w").write("".join(json.dumps({"payload": {"type": "token_count", "rate_limits": {"primary": p}}}) + "\n" for p in ps))
    os.makedirs(os.path.join(new, "sessions"))
    reading("rollout-a.jsonl", {"used_percent": 90, "resets_at": 1})
    assert plan_used(new) == (0, None), "a reset window reads as 0"
    reading("rollout-a.jsonl", {"used_percent": 90})
    assert plan_used(new) == (90, None), "a reading without resets_at is current"
    shutil.rmtree(new)
    for harness, change, reason in (("codex", {"max_seconds": None}, "needs --max-seconds"), ("codex", {"budget": "1"}, "--budget is Claude's"),
                                    ("codex", {"max_plan_percent": None}, "needs --max-plan-percent"), ("codex", {"max_plan_percent": 101}, "needs --max-plan-percent"),
                                    ("codex", {"effort": None}, "needs --effort"), ("codex", {"permission_mode": "auto"}, "one of"),
                                    ("claude", {"effort": "high"}, "are Codex's"), ("claude", {"max_seconds": 60}, "are Codex's"), ("claude", {"max_plan_percent": 50}, "are Codex's"),
                                    ("claude", {"budget": None}, "needs --budget")):
        valid = dict(harness=harness, runs=1, budget=None, model="fake", effort="high", permission_mode="workspace-write", max_seconds=60, max_plan_percent=100,
                     config_dir=home, case=[], out=os.path.join(tmp, "refused"), claude=shims["claude"], codex=shims["codex"])
        if harness == "claude":
            valid.update(budget="0.01", effort=None, max_seconds=None, max_plan_percent=None, permission_mode="fake", config_dir=config)
        try:
            run(argparse.Namespace(**(valid | change)))
            raise AssertionError(f"{harness} {change} accepted")
        except SystemExit as err:
            assert reason in str(err), (harness, change, err)
    args = argparse.Namespace(**(vars(args) | {"out": os.path.join(tmp, "refused")}))
    for name, content, reason in (("AGENTS.md", "", "AGENTS.md"), ("skills/mine/SKILL.md", "", "skills/mine"), ("memories/note.md", "x", "memories/"),
                                  ("config.toml", 'model = "x"\n', "config.toml key model"), ("config.toml", '[projects."/x"]\nsandbox_mode = "x"\n', "projects./x.sandbox_mode"),
                                  ("config.toml", "[", "config.toml unreadable")):
        path = os.path.join(home, name)
        os.makedirs(os.path.dirname(path), exist_ok=True)
        saved = open(path).read() if os.path.exists(path) else None
        open(path, "w").write(content)
        try:
            run(args)
            raise AssertionError(f"{name} accepted")
        except SystemExit as err:
            assert "not clean" in str(err) and reason in str(err), (name, reason, err)
        if saved is not None:
            open(path, "w").write(saved)
        else:
            os.remove(path)
            if name.startswith("skills/"):
                os.rmdir(os.path.dirname(path))
    args = argparse.Namespace(runs=1, budget="0.01", model="fake", permission_mode="fake", config_dir=config, case=[], claude=shims["claude"], codex=shims["codex"],
                              harness="claude", effort=None, max_seconds=None, max_plan_percent=None)
    for name, content, reason in (("settings.json", '{"autoMemoryEnabled": false, "hooks": {}}', "key hooks"),
                                  ("settings.json", '{"autoMemoryEnabled": true}', "autoMemoryEnabled"), ("settings.json", '{"tui": "fullscreen"}', "autoMemoryEnabled"),
                                  ("settings.json", "[]", "key not an object"), ("skills/mine/SKILL.md", "", "skills/mine"),
                                  ("skills/synced/x/mine/SKILL.md", "", "skills/synced/x/mine"), ("skills/synced/x/manifest.json", "{", "manifest.json unreadable"),
                                  ("skills/synced/stray.json", "", "skills/synced/stray.json"), ("plugins/synced/y/mine/plugin.json", "", "plugins/synced/y/mine")):
        os.makedirs(os.path.dirname(os.path.join(config, name)), exist_ok=True)
        saved = open(os.path.join(config, name)).read() if os.path.exists(os.path.join(config, name)) else None
        open(os.path.join(config, name), "w").write(content)
        args.out = os.path.join(tmp, "refused-" + reason.replace("/", "-"))  # fresh, so an accepted shape runs and fails the assertion below
        try:
            run(args)
            raise AssertionError(f"{name} accepted")
        except SystemExit as err:
            assert "not clean" in str(err) and reason in str(err), (name, reason, err)
        open(os.path.join(config, name), "w").write(saved) if saved is not None else shutil.rmtree(os.path.dirname(os.path.join(config, name)))
    shutil.rmtree(tmp)
    print("selftest: ok")


def main():
    for sig in (signal.SIGTERM, signal.SIGHUP):  # unwind like Ctrl-C, so a running harness is killed
        signal.signal(sig, lambda n, _: sys.exit(128 + n))
    if sys.argv[1:2] == ["_fake"]:
        return fake(sys.argv[2], sys.argv[3:])
    parser = argparse.ArgumentParser(prog="evals/run.py", description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = parser.add_subparsers(dest="command", required=True)
    p = sub.add_parser("run", help="run the cases on Claude (up to runs x cases x budget) or Codex (up to runs x cases x max-seconds)")
    p.add_argument("--harness", choices=("claude", "codex"), default="claude")
    p.add_argument("--runs", type=int, required=True)
    p.add_argument("--budget", help="Claude only, required: USD per run, passed as --max-budget-usd")
    p.add_argument("--model", required=True)
    p.add_argument("--effort", help="Codex only, required: passed as -c model_reasoning_effort=EFFORT")
    p.add_argument("--permission-mode", required=True, help=f"Claude's --permission-mode; for Codex one of {', '.join(CODEX_MODES)}")
    p.add_argument("--max-seconds", type=int, help="Codex only, required: each run is killed at this cap")
    p.add_argument("--max-plan-percent", type=int, help="Codex only, required: no run starts once the five-hour plan use, plus the largest run's, would reach it")
    p.add_argument("--config-dir", required=True, help="CLAUDE_CONFIG_DIR or CODEX_HOME for every run; must hold no customization")
    p.add_argument("--case", action="append", default=[], help=f"one of {', '.join(CASES)}; default {', '.join(DEFAULT)}")
    p.add_argument("--out", help="output directory; default a new temporary one")
    p.add_argument("--claude", default="claude", help="the Claude executable")
    p.add_argument("--codex", default="codex", help="the Codex executable")
    sub.add_parser("selftest", help="check the runner against a fake claude and codex; spends nothing")
    args = parser.parse_args()
    try:
        run(args) if args.command == "run" else selftest()
    except Failed as err:
        raise SystemExit(str(err))


if __name__ == "__main__":
    main()
