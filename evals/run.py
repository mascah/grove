#!/usr/bin/env python3
"""Behavioral evaluations of Grove's headless shaping workflow (G-108).

    python3 evals/run.py run --runs N --budget USD --model MODEL \\
        --permission-mode MODE --config-dir DIR [--case NAME]... [--out DIR]
    python3 evals/run.py selftest

`run` spends money: every spend parameter is required and has no default.
`selftest` spends nothing: it drives the same runner with a fake `claude` and
asserts that the checks pass and fail where they should. evals/README.md says
what a run retains, what each check means, and how to score the rubric.
"""
import argparse, datetime, glob, json, os, re, shlex, shutil, signal, subprocess, sys, tempfile, time

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
FIXTURE = os.path.join(ROOT, "evals", "fixture")
# Variables through which Git takes a repository from its caller; a hook exports GIT_DIR (G-089).
GIT_LOCATION = ("GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_NAMESPACE")
IDENTITY = ["-c", "user.name=Grove Eval", "-c", "user.email=eval@example.invalid", "-c", "commit.gpgsign=false"]
TIMEOUT = 1800  # seconds per run; failure detection only, the budget is the real bound
CASES = {
    "missing-choice": {"topic": "hide finished tasks from tasks list by default", "question": True},
    "companion": {"topic": "let tasks list filter by tag", "question": False},
}
# Files some step of the shaping guide needs for these topics; any other read is listed as unneeded.
NEEDED = {"AGENTS.md", "CLAUDE.md", "grove.yaml", "grove/brief.md", "tasks.py"}
CUSTOMIZATION = ("CLAUDE.md", "agents", "commands", "output-styles", "hooks", "settings.local.json")
# A login writes settings.json; these keys shape the terminal and memory, not what the agent reads or may do.
SETTINGS = {"tui", "theme", "autoMemoryEnabled"}
GROVE = "grove version (its revision can lag in a linked worktree; the digest pins the guides)"
UNTOUCHED = {"proposed", "open", "current", None}
# Claude's own auth variables pass through; every other CLAUDE* variable is the caller's session leaking in.
AUTH = ("CLAUDE_CODE_OAUTH_TOKEN",)


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


def build(work):
    """Build the CLI from this checkout and the fixture from evals/fixture; return (grove, template, version)."""
    grove = os.path.join(work, "bin", "grove")
    sh("go", "build", "-o", grove, "./cmd/grove", cwd=ROOT)
    template = os.path.join(work, "fixture")
    shutil.copytree(FIXTURE, template, ignore=shutil.ignore_patterns("records"))
    git(template, "init", "-q", "-b", "main")
    sh(grove, "--project", template, "init")
    shutil.copy(os.path.join(FIXTURE, "brief.md"), os.path.join(template, "grove", "brief.md"))
    path = sh(grove, "--project", template, "new", "work", "Sync tasks between two machines").stdout.strip()
    with open(os.path.join(template, path), "r+") as f:
        head = f.read().split("\n---\n", 1)[0]
        f.seek(0), f.truncate()
        f.write(head + "\n---\n\n" + open(os.path.join(FIXTURE, "records", "sync.md")).read())
    sh(grove, "--project", template, "check")
    git(template, "add", "-A")
    git(template, "commit", "-q", "-m", "Fixture: the tasks tool and its Grove project")
    return grove, template, sh(grove, "version").stdout.strip()


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


def retrieval(transcript, clone, created):
    """Facts from the trace's tool calls: what the session used and read. Reported, never scored."""
    commands, files = [], []
    for ev in events(transcript):
        if ev.get("type") != "assistant":
            continue
        for block in ev.get("message", {}).get("content", []):
            if not isinstance(block, dict) or block.get("type") != "tool_use":
                continue
            inp = block.get("input", {})
            if block.get("name") == "Read":
                files.append(inp.get("file_path", ""))
            elif block.get("name") == "Bash":
                cmd = inp.get("command", "")
                commands.append(cmd)
                for part in re.split(r"&&|\|\||;|\|", cmd):
                    try:
                        words = shlex.split(part)
                    except ValueError:
                        continue
                    if words and words[0] in ("cat", "head", "tail", "sed", "nl", "less", "awk"):
                        files += [w for w in words[1:] if not w.startswith("-") and os.path.isfile(os.path.join(clone, w))]
    rel, clone = [], os.path.realpath(clone)
    for f in files:
        full = os.path.realpath(os.path.join(clone, f))
        p = os.path.relpath(full, clone)
        rel.append(full if p.startswith("..") else re.sub(r"^\.claude/worktrees/[^/]+/", "", p))
    used = set()  # grove subcommands actually invoked, not words that merely follow "grove" in a command
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
    return {
        "guide": "guide" in used,
        "brief": "brief" in used or "grove/brief.md" in rel,
        "list": "list" in used,
        "context_or_show": bool(used & {"context", "show"}),
        "files_read": sorted(set(rel)),
        "unneeded": sorted({p for p in rel if p not in NEEDED and not p.startswith("tasks/") and p not in created}),
        "commands": commands,
    }


def one(args, grove, template, work, case_name, n, meta):
    case = CASES[case_name]
    rdir = os.path.join(work, f"{case_name}-{n}")
    os.makedirs(rdir)
    remote, clone = os.path.join(rdir, "remote.git"), os.path.join(rdir, "p")
    git(rdir, "clone", "-q", "--bare", template, remote)
    git(rdir, "clone", "-q", remote, clone)
    git(clone, "config", "user.name", "Grove Eval")
    git(clone, "config", "user.email", "eval@example.invalid")
    before = snapshot(clone, remote)
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
            proc.wait(TIMEOUT)
        except subprocess.TimeoutExpired:
            timed_out = True
        finally:  # also on Ctrl-C, and after a normal exit, for background processes the session left
            try:
                os.killpg(proc.pid, signal.SIGKILL)
            except OSError:
                pass
            proc.wait()
    init = next((ev for ev in events(transcript) if ev.get("type") == "system" and ev.get("subtype") == "init"), {})
    result = next((ev for ev in events(transcript) if ev.get("type") == "result"), {})
    run = {
        "case": case_name, "run": n, "topic": case["topic"], "command": command,
        "exit": proc.returncode, "timed_out": timed_out, "wall_seconds": round(time.time() - started, 1),
        **{k: meta[k] for k in ("claude", GROVE, "base commit", "fixture commit")}, "model_requested": args.model, "model_reported": init.get("model"),
        "permission_mode_reported": init.get("permissionMode"),
        "cost_usd": result.get("total_cost_usd"), "turns": result.get("num_turns"), "duration_ms": result.get("duration_ms"),
        "result_subtype": result.get("subtype"), "is_error": result.get("is_error"),
        "permission_denials": len(result.get("permission_denials") or []),
        "message": result.get("result", ""),
    }
    try:  # the session has spent by now: a failure reading its effects must not lose what it cost
        after = state(clone, remote, grove, rdir, before["main"])
        json.dump(dict(after, before=before), open(os.path.join(rdir, "state.json"), "w"), indent=1)
        run["checks"] = checks(case, before, after, run["message"])
        run["retrieval"] = retrieval(transcript, clone, {p for b in after["proposals"].values() for p in b["touched"]})
    except Exception as err:
        run.update(error=str(err), checks={"runner": f"fail: {err}"})
    json.dump(run, open(os.path.join(rdir, "run.json"), "w"), indent=1)
    return run


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
    lines += ["- Codex row: not built (G-108 follow-on)", f"- Cases not run: {', '.join(unrun) or 'none'}", ""]
    for name in CASES:
        rs = [r for r in runs if r["case"] == name]
        if not rs:
            continue
        keys = list(dict.fromkeys(k for r in rs for k in r["checks"]))
        lines += [f"## {name}: {CASES[name]['topic']}", "", "| run | harness | " + " | ".join(keys) + " | cost | turns | seconds | guide | brief | list | context/show | unneeded reads |",
                  "|" + " --- |" * (len(keys) + 10)]
        for r in rs:
            f = r.get("retrieval") or dict.fromkeys(("guide", "brief", "list", "context_or_show"), "-") | {"unneeded": []}
            lines.append(f"| {r['run']} | {harness(r)} | " + " | ".join(mark(r["checks"].get(k)) for k in keys)
                         + f" | {r.get('cost_usd')} | {r.get('turns')} | {(r.get('duration_ms') or 0) / 1000:.0f} | {f['guide']} | {f['brief']} | {f['list']} | {f['context_or_show']} | {', '.join(f['unneeded']) or '-'} |")
        lines += ["", "Failures:", ""]
        lines += [f"- run {r['run']} {k}: {v}" for r in rs for k, v in r["checks"].items() if v != "pass"] or ["- none"]
        lines += ["", "Rubric (evals/README.md), scorer `owner` or `judge`:", "", "| run | scorer | presumes choice | planted question | brief constraint | handoff | notes |", "| --- | --- | --- | --- | --- | --- | --- |"]
        lines += [f"| {r['run']} |  |  |  |  |  |  |" for r in rs]
        lines.append("")
    open(path, "w").write("\n".join(lines))


def run(args):
    for name in args.case:
        if name not in CASES:
            raise SystemExit(f"unknown case {name}; cases: {', '.join(CASES)}")
    if not re.fullmatch(r"[0-9]+(\.[0-9]+)?", args.budget) or float(args.budget) <= 0:
        raise SystemExit(f"--budget must be a positive dollar amount, not {args.budget!r}")
    if args.runs < 1:
        raise SystemExit("--runs must be at least 1")
    args.config_dir = os.path.abspath(args.config_dir)
    os.makedirs(args.config_dir, exist_ok=True)
    found = [c for c in CUSTOMIZATION if os.path.exists(os.path.join(args.config_dir, c))]
    plugins = os.path.join(args.config_dir, "plugins", "installed_plugins.json")
    if os.path.exists(plugins) and json.load(open(plugins)).get("plugins"):
        found.append("installed plugins")
    settings = os.path.join(args.config_dir, "settings.json")
    settings = json.load(open(settings)) if os.path.exists(settings) else {}
    settings = settings if isinstance(settings, dict) else {"not an object": settings}
    found += [f"settings.json key {k}" for k in sorted(settings.keys() - SETTINGS)]
    if settings.get("autoMemoryEnabled"):  # memory is written under the config directory and would carry across runs
        found.append("settings.json autoMemoryEnabled true")
    # A login syncs the account's Anthropic skills and plugins under skills/synced/ID and plugins/synced/ID, listed
    # in each ID's manifest.json; they cannot be kept out and a preview user has them too, so they are recorded,
    # not refused. Anything else under skills/, or in a synced ID directory that its manifest does not name, is authored.
    found += [f"skills/{e}" for e in sorted(os.listdir(os.path.join(args.config_dir, "skills"))) if e != "synced"] if os.path.isdir(os.path.join(args.config_dir, "skills")) else []
    synced = {}
    for kind in ("skills", "plugins"):
        for d in sorted(glob.glob(os.path.join(args.config_dir, kind, "synced", "[!.]*"))):
            m = os.path.join(d, "manifest.json")
            names = {e.get("name", "?") for e in json.load(open(m)).get(kind, [])} if os.path.exists(m) else set()
            synced[kind] = sorted(synced.get(kind, []) + sorted(names))
            found += [f"{kind}/synced/{e}" for e in sorted(os.listdir(d)) if os.path.isdir(os.path.join(d, e)) and not e.startswith(".") and e not in names]
    if found:
        raise SystemExit(f"--config-dir {args.config_dir} is not clean: {', '.join(found)}")
    cases = args.case or list(CASES)
    work = os.path.abspath(args.out) if args.out else tempfile.mkdtemp(prefix="grove-evals-")
    os.makedirs(work, exist_ok=True)
    meta = {"started": datetime.datetime.now(datetime.timezone.utc).isoformat(timespec="seconds"), "output": work,
            "runs per case": args.runs, "budget per run (USD)": args.budget,
            "cap (USD)": f"{float(args.budget) * args.runs * len(cases):.2f}", "model": args.model,
            "permission mode": args.permission_mode, "config dir": args.config_dir,
            "config dir holds": ", ".join(sorted(os.listdir(args.config_dir))) or "nothing",
            "config dir settings": json.dumps(settings, sort_keys=True), "config dir synced skills": ", ".join(synced.get("skills", [])) or "none",
            "config dir synced plugins": ", ".join(synced.get("plugins", [])) or "none"}
    exe = shutil.which(args.claude)
    if not exe:
        meta["claude"] = f"unavailable: {args.claude} not found; no case ran"
        report(os.path.join(work, "report.md"), meta, [], list(CASES))
        raise SystemExit(meta["claude"] + f"\nreport: {work}/report.md")
    args.claude = exe
    version = sh(exe, "--version", environ=dict(env(), CLAUDE_CONFIG_DIR=args.config_dir)).stdout.strip()
    grove, template, grove_version = build(work)
    meta.update({"claude": version, GROVE: grove_version, "base commit": git(ROOT, "rev-parse", "HEAD") + (" with uncommitted changes" if git(ROOT, "status", "--porcelain") else ""),
                 "fixture commit": git(template, "rev-parse", "HEAD")})
    print(f"output {work}; spending at most ${meta['cap (USD)']}", file=sys.stderr)
    runs = []
    for name in cases:
        for n in range(1, args.runs + 1):
            print(f"{name} {n}/{args.runs}", file=sys.stderr)
            try:
                runs.append(one(args, grove, template, work, name, n, meta))
            except Exception as err:  # a runner failure on one run is reported, and the rest still run
                r = {"case": name, "run": n, "error": str(err), "checks": {"runner": f"fail: {err}"}}
                os.makedirs(os.path.join(work, f"{name}-{n}"), exist_ok=True)
                json.dump(r, open(os.path.join(work, f"{name}-{n}", "run.json"), "w"), indent=1)
                runs.append(r)
            report(os.path.join(work, "report.md"), meta, runs, [c for c in CASES if c not in cases])
    print(os.path.join(work, "report.md"))
    return runs


def fake(argv):
    """A stand-in for `claude -p` acting out a scripted outcome, chosen by GROVE_EVAL_FAKE:
    good follows the guide; bad is G-078's divergence (no question) plus a write to the
    session checkout; surfaced asks the question without `blocks`; worse pushes, promotes,
    breaks `check`, names a stale commit and, on the companion, leaves two proposal branches."""
    if "--version" in argv:
        return print("0.0.0 (fake claude)")
    topic = argv[argv.index("-p") + 1].removeprefix("/grove-shape ").removesuffix(" --interaction headless")
    missing, mode = topic == CASES["missing-choice"]["topic"], os.environ["GROVE_EVAL_FAKE"]
    emit = lambda ev: print(json.dumps(ev), flush=True)
    tool = lambda name, **inp: emit({"type": "assistant", "message": {"content": [{"type": "tool_use", "name": name, "input": inp}]}})
    emit({"type": "system", "subtype": "init", "model": "fake", "permissionMode": argv[argv.index("--permission-mode") + 1]})
    if mode == "worse":  # commands that mention grove subcommands without running them
        for cmd in ('grove new work "Let tasks list filter by tag"', 'grove new question "Should tasks list show dropped tasks?"',
                    "cd /tmp/grove-evals-x/p && git worktree list", "cat > /tmp/x.md <<EOF\nthe brief and the guide\nEOF"):
            tool("Bash", command=cmd)
    else:
        tool("Bash", command="grove guide shape")
        tool("Read", file_path=os.path.abspath("grove/brief.md"))
        tool("Bash", command="grove list && cat tasks.py")
        tool("Read", file_path=os.path.expanduser("~/.claude/CLAUDE.md"))
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
    emit({"type": "result", "subtype": "success", "is_error": False, "total_cost_usd": 0, "num_turns": 5, "duration_ms": 1000, "result": text})


def selftest():
    tmp = tempfile.mkdtemp(prefix="grove-evals-selftest-")
    shim = os.path.join(tmp, "claude")
    open(shim, "w").write(f"#!/bin/sh\nexec {shlex.quote(sys.executable)} {shlex.quote(os.path.abspath(__file__))} _fake \"$@\"\n")
    os.chmod(shim, 0o755)
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
    reasons = {"bad": "fail: no question or decision", "surfaced": "fail: surfaced, not blocking"}
    for mode in ("good", "bad", "surfaced", "worse"):
        os.environ["GROVE_EVAL_FAKE"] = mode
        args = argparse.Namespace(runs=1, budget="0.01", model="fake", permission_mode="fake", config_dir=config,
                                  case=[], out=os.path.join(tmp, mode), claude=shim)
        for r in run(args):
            failed = {k for k, v in r["checks"].items() if v != "pass"}
            assert failed == expected[(mode, r["case"])], (mode, r["case"], r["checks"])
            if mode in reasons and r["case"] == "missing-choice":
                assert r["checks"]["question-blocks-proposal"].startswith(reasons[mode]), r["checks"]
            f = r["retrieval"]
            if mode == "worse":
                assert not (f["guide"] or f["brief"] or f["list"] or f["context_or_show"]), f
                continue
            assert f["guide"] and f["brief"] and f["list"] and not f["context_or_show"], f
            assert "tasks.py" in f["files_read"] and len(f["unneeded"]) == 1 and f["unneeded"][0].endswith(".claude/CLAUDE.md"), f
    assert open(os.path.join(tmp, "good", "report.md")).read().count("- config dir synced skills: pdf") == 1
    for name, content in (("settings.json", '{"hooks": {}}'), ("settings.json", '{"autoMemoryEnabled": true}'), ("settings.json", "[]"),
                          ("skills/mine/SKILL.md", ""), ("skills/synced/x/mine/SKILL.md", ""), ("plugins/synced/y/mine/plugin.json", "")):
        os.makedirs(os.path.dirname(os.path.join(config, name)), exist_ok=True)
        saved = open(os.path.join(config, name)).read() if os.path.exists(os.path.join(config, name)) else None
        open(os.path.join(config, name), "w").write(content)
        try:
            run(args)
            raise AssertionError(f"{name} accepted")
        except SystemExit as err:
            assert "not clean" in str(err), (name, err)
        open(os.path.join(config, name), "w").write(saved) if saved is not None else shutil.rmtree(os.path.dirname(os.path.join(config, name)))
    shutil.rmtree(tmp)
    print("selftest: ok")


def main():
    for sig in (signal.SIGTERM, signal.SIGHUP):  # unwind like Ctrl-C, so a running harness is killed
        signal.signal(sig, lambda n, _: sys.exit(128 + n))
    if sys.argv[1:2] == ["_fake"]:
        return fake(sys.argv[2:])
    parser = argparse.ArgumentParser(prog="evals/run.py", description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = parser.add_subparsers(dest="command", required=True)
    p = sub.add_parser("run", help="run the cases on Claude; spends up to runs x cases x budget")
    p.add_argument("--runs", type=int, required=True)
    p.add_argument("--budget", required=True, help="USD per run, passed as --max-budget-usd")
    p.add_argument("--model", required=True)
    p.add_argument("--permission-mode", required=True)
    p.add_argument("--config-dir", required=True, help="CLAUDE_CONFIG_DIR for every run; must hold no customization")
    p.add_argument("--case", action="append", default=[], help=f"one of {', '.join(CASES)}; default all")
    p.add_argument("--out", help="output directory; default a new temporary one")
    p.add_argument("--claude", default="claude", help="the harness executable")
    sub.add_parser("selftest", help="check the runner against a fake claude; spends nothing")
    args = parser.parse_args()
    try:
        run(args) if args.command == "run" else selftest()
    except Failed as err:
        raise SystemExit(str(err))


if __name__ == "__main__":
    main()
