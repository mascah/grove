# Implementation: Run one bounded apply session

<!-- Progress on change 261021-2aj5w, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Run one bounded apply session

Status: done

Constitution amendment #17 is merged (dcab551), so code followed it.

What changed:

- `internal/cli/run.go`: `grove run [--retry] [--timeout <duration>]`.
  Standard library only (`syscall.Flock` for the lock, `Setpgid` and
  `syscall.Kill(-pgid, …)` to stop the group: SIGTERM, then SIGKILL after
  5 s). Records and output under `.git/grove/<attempt>.json` and `.log`.
- `internal/cli/project.go`: `parseProject` and `project.at(commit)`, so
  `run` reads `grove.yaml` and the documents from the fetched target.
- `internal/cli/instructions.go`: no-ID apply's selection moved into
  `buildable`, which `run` reuses, so both select the same changes.
- `internal/cli/cli.go`: usage and the `run` flags.
- `docs/running.md` (new page, in `mkdocs.yml`), `docs/reference.md` rows
  for `grove run`, `--retry`, `--timeout`, and `dev/walkthrough.md`
  "Running apply unattended".
- No new package, so `CLAUDE.md`'s layout table is unchanged.

Checks:

- `TestRun` in `internal/cli/walkthrough_test.go`, stand-in `claude` on
  `PATH`, operator clone of a bare remote: nothing to apply; agents without
  `claude`; local target with its own commit; unclean checkout; a run in
  progress; an ordinary run whose session leaves a child (group gone, record
  names A then B and the target commit, `git status` clean); paused, still
  paused after an unrelated approval; batch archived (A finished, B
  dropped) releases the pause and selects only C; failed to start, then
  paused; cancelled by SIGINT after an edit (group gone), retry refuses
  while the edit is uncommitted; timed out with a child that ignores
  SIGTERM (group gone). Disabling the group stop makes it fail three ways.
  The test does not run a second project's checkout beside the first; the
  lock is a file in each checkout's own `.git/grove/`.
- `just check`: passes; 2,811 Go lines, 863 test lines.
- Real Claude Code 2.1.306 on the ledger fixture, local bare remote,
  `.claude/settings.local.json` as in `docs/running.md`, the two
  `rules.skills.apply` lines:
  - `grove run` with A (unknown ID) and B (delete, depends on A): one
    session, about 80 s; a second `grove run` meanwhile said a run is in
    progress. It built A then B, archived both, made four commits naming
    task IDs, pushed `grove/run-261021-6jbj1`, and ended exit status 0.
    `gh pr create` failed, as the remote is not on GitHub; the agent said
    so in its last message.
  - `grove run` then said paused; after merging the branch into `origin`'s
    main, `grove run` said nothing to apply.
  - With change C approved, `grove run`, SIGINT once `ledger/report.py`
    was modified: cancelled, no process left in the session's group, edit
    and branch kept; `grove run --retry` refused and left the edit.
  - After discarding it, `grove run --retry --timeout 1m`: timed out at
    61 s, no `claude -p` left, its uncommitted archive left in the checkout.
  - Found: `claude -p` prints only its last message, so a stopped
    session's `.log` is empty. Documented in `docs/running.md`.

- Against GitHub, with the owner's go-ahead: a private, disposable
  repository `mascah/grove-run-trial` holding the fixture at A and B
  approved, cloned as the operator's checkout. `grove run`: attempt
  `261021-f8b73` ended exit status 0; the agent built A then B, archived
  both, made four commits and opened pull request #1, "Apply 261021-a0001
  and 261021-b0002", listing both changes and "Blocked tasks: None".
  `grove run` before the merge said paused; after `gh pr merge 1`, it said
  "Nothing to apply on main at 536d352."

### T-02 — Run it on a launchd timer

Status: done

`docs/guides/run-apply-on-a-timer.md` (in `mkdocs.yml`): the job
definition, the Claude Code settings and login an unattended session
needs, `PATH`, and where each attempt's output is.

Trial on macOS 26.6.2, the guide's plist with `StartInterval` 60, loaded
with `launchctl bootstrap gui/501` and unloaded afterwards, in a fresh clone
of the trial's bare remote with change C approved, real Claude Code:

- Loaded 14:58:02; launchd did not fire at load (runs 0), so the guide
  says the first firing is one interval later.
- 14:59:02 firing: one attempt, `261021-64m36`, which built and archived C
  in two commits and ended exit status 0 at 15:01:04. The 15:00 firing
  came while it ran and launchd started nothing (runs stayed 1). Claude
  Code's Terminal login worked from the launchd job.
- 15:02:04 firing: "the project is paused after attempt 261021-64m36 …",
  last exit code 1, still one attempt in `.git/grove/`.
- Stopping: with a stand-in `claude` that sleeps, and `--retry`, the job
  fired, then `launchctl bootout` during the session: the log says
  cancelled, and no process was left in the session's group.

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->

- runs: applied to specs/runs/spec.md, created with its eight ADDED requirements
