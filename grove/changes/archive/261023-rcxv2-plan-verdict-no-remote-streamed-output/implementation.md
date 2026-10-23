# Implementation: Plan's verdict, runs without a remote, streamed output

## Tasks

### T-01 — Plan ends with a verdict

Status: done

- `internal/instruct/requires/plan.md`: the finish has the agent say
  whether `grove run` would build the cycle once marked, and what removes
  each obstacle (a member needing a person, one waiting outside the cycle,
  something the runner's checkout must have), or how to add the mark.
- `docs/planning.md` shows the verdict for a cycle with a member that
  needs a person and for one with nothing in the way; `dev/walkthrough.md`
  planning steps 1 and 3 run them.
- `TestCycles` checks `grove instructions plan` prints the requirement.
  `just check` passes.
- The ledger fixture run of `/grove-plan` in Claude Code was not done: the
  owner, 2026-10-23, "I'll skip the fixture check, it's getting too time
  consuming to setup".

### T-02 — Run without a remote

Status: done

- `internal/cli/run.go`: `grove run` fetches and checks the fast-forward
  only when `git remote` lists `origin`; otherwise the local target's
  commit is the fetched target.
- `docs/running.md` has "Without a remote", saying the owner merges the
  attempt's branch; "The remote" and the ordinary run's step 1 say
  `origin` is optional. `docs/reference.md`'s `grove run` row no longer
  requires `origin`. `dev/walkthrough.md` running step 6 runs it.
- `TestRun` ends with the checkout's `origin` removed: a local commit of
  its own is no refusal, the record holds the local target's commit, the
  project is paused, and merging the attempt's branch locally ends the
  pause. The `origin` cases above it are unchanged. Without the code
  change it fails with Git's "'origin' does not appear to be a git
  repository". `just check` passes.

### T-03 — Show the session as it runs

Status: done

- `internal/cli/run.go`: the session is started with `--output-format
  stream-json --verbose`, its standard output through a pipe; `relay`
  writes every line to the attempt's log as it comes and prints the first
  line of each assistant text, each tool call (name and the first of its
  `command`, `file_path`, `pattern`, `url` or `description`) and each tool
  result, indented. Standard error still goes to the log. Once the group
  is stopped, reading waits at most 5 seconds for a process that left the
  group and holds the pipe.
- `docs/running.md` shows the lines and the JSON Lines log and no longer
  says a stopped session's output is usually empty;
  `docs/guides/run-apply-on-a-timer.md` says the launchd log receives the
  same lines; `dev/walkthrough.md` running steps 2 and 4 run them.
- `TestRun`: the stand-in `claude` emits stream-json events (text, tool
  calls, string and array results, system and result events); the test
  checks the log holds them before the session ends, the lines on
  standard output, the log's content after, and the new flags.
- On the ledger fixture, 2026-10-23, with no remote and Claude Code: `grove
  run` on a marked cycle of one change printed a line per message, tool
  call and result, ended `exit status 0`, and its 45-line log parsed as
  JSON on every line. The same run showed T-02: the record held local
  `main`'s commit, the project was paused, and merging the branch into
  `main` gave "Nothing to apply". The agent left the archive uncommitted,
  its rules asking only for a commit per task, so "Without a remote" in
  `docs/running.md` now suggests a rule to commit what is left.
- `just check` passes.

## Reconciliation

- cycles: applied to specs/cycles/spec.md
- runs: applied to specs/runs/spec.md
