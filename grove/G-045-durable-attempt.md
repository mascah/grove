---
id: "G-045"
type: work
title: "Run one bounded implementation independently of the viewing terminal"
status: review
created: "2026-09-21T00:54:16Z"
updated: "2026-09-23T02:54:54Z"
kind: feature
size: large
priority: 4
depends_on: ["G-038", "G-040"]
relates_to: ["G-035", "G-044", "G-046", "G-101"]
formerly: "W-027"
candidate: "d8ed3159447bdc7f4930499aaf6a57e182e68c04"
---

## Outcome

One explicitly assigned implementation runs independently of the viewing
terminal, leaves a durable result, and can be inspected, reconnected or stopped
without accidentally starting it again.

## Scope and bounds

Reuse the interactive assignment/preparation/review contract. Begin with one
work item and one attempt, then bounded independent review where the workflow
requires it. No schedules, dream loop, automatic scope expansion, arbitrary
multi-item fan-out or automatic merge. Keep ordinary CLI access daemon-free.

Define attempt identity, source/input revisions, working directory, process
owner, budgets, permission profile, events/logs, checkpoint, result, wait,
cancellation and interruption/recovery. Persist raw output separately from
durable review evidence. Record actual model/role configuration. Bound subagents
and retries explicitly. Process exit and streamed claims are not acceptance.

## Acceptance

1. Compare a Grove-owned `claude -p` process with available native Claude
   background sessions; select the smallest mechanism that meets this contract.
   Evaluate tmux only as an option. Pin capabilities to the installed version.
2. Demonstrate competing-start refusal, TUI/terminal exit survival, reconnect
   without replayed work, Stop after reconnect and explicit owner-loss handling.
3. Bound and handle partial/unknown/oversized events, provider errors, budget
   exhaustion, stale inputs and missing results. Preserve partial code/evidence.
4. A missing human decision persists a question and wait; an unchanged wait
   cannot trigger an uncontrolled retry. A result identifies exact revisions.
5. Exercise fake-process lifecycle failures first and one explicitly bounded
   real-provider trial. Distinguish observed recovery from untested shutdown
   behavior; machine-reboot recovery is not a selected requirement.

## Evidence

Branch `worktree-G-045` from `main` at `5ae87d5`, in
`.claude/worktrees/worktree-G-045`. Plan [G-100](G-100-g-045-durable-attempt-plan.md)
(`a041d22`), started from this record at revision `c976f07d…`. Commits:
`3af1c93` active, `83092cb` the implementation, `cab470c` documentation and
decision [G-101](G-101-attempt-mechanism.md), `69b9d9e` the review refusal
the trial exposed, `9ae71ac` review round one, `b4a93a7` the test budget, `30b9c20` review
round two;
the candidate is the evidence commit named in `candidate`. Review record:
[G-102](G-102-g-045-durable-attempt-review.md), with the independent
findings, their dispositions, and the trial evidence.

Changed behaviour, against each acceptance item:

1. **Mechanism.** The comparison in G-100 and the decision in G-101: native
   background sessions (`--bg`, `agents`, `attach`, `logs`, `stop`, `rm`)
   are interactive sessions under Claude's supervisor daemon, reject `-p`,
   print terminal text rather than events, have no budget flag and move
   into their own worktree; tmux would host the same process and add a
   dependency. Selected: a Grove-owned `claude -p "/grove-work ID
   --interaction headless" --output-format stream-json --verbose
   --session-id UUID --max-budget-usd N --permission-mode M
   --permission-prompts none` under an on-demand owner. Pinned to Claude
   Code 2.1.280; each attempt records `claude --version`, the `system/init`
   model, permission mode, version and `capabilities`, and the `result`
   fields, from the run itself.
2. **Lifecycle.** `grove run` refuses a second start while an attempt of
   the work runs or is orphaned; the owner runs in its own session with the
   lock inherited from the launcher, so the launching shell exiting changes
   nothing (the real trial launched from `sh -c` that exited at once;
   the owner reparented to launchd and finished 62 s later); `attempts`
   and `attempt` are reads of files, and the fake counts one provider
   start across any number of reconnects; `stop` after a reconnect sends
   SIGINT through the owner (the real provider ended within a second with
   an `error_during_execution` result, `stopped: true`); an owner killed
   with SIGKILL reads as `orphaned` while the provider's process group
   lives and `interrupted` once nothing is left, `run` refuses over an
   orphan, and `stop` ends and reconciles it.
3. **Bounds.** `events.jsonl` is raw; the reader counts oversized (over
   1 MiB), malformed, unknown and partial lines and keeps only the init and
   result fields (tested with a crafted file). A result with `is_error`, a
   nonzero exit, a missing result event, and a provider that cannot start
   are recorded as they are. Budget exhaustion is a result subtype Claude
   emits; it was not reached in the trials and is covered only by a fake.
   The record on the target changing after launch is reported as `inputs
   changed`. Partial work is untouched by Stop and counted as dirty
   (untracked included); the next attempt reuses the branch's worktree.
4. **Waits.** `run` refuses while an open question blocks the work, in the
   launching checkout or on the branch, and while the branch's record is
   already in review; rerunning with nothing changed refuses identically
   and writes no attempt (tested). The result names the base, the
   worktree's HEAD and the record's revision.
5. **Order of evidence.** Nine fake-process tests came first
   (`internal/attempt`, behind `-short`), then two bounded real trials in a
   disposable repository with a clone-built binary (G-102): G-001 ran to a
   handoff in review with a candidate for $0.50 in 9 turns, and G-002 was
   stopped mid-run after a reconnect for $0.23. Observed: launcher exit
   survival, reconnect, Stop, owner loss (fake only), result
   reconciliation. Untested: machine reboot (not selected), budget
   exhaustion on the real provider, SIGKILL of the real provider.

Decisions inside the outcome: budget and permission mode are required
flags, since a default for either would make the choice in practice;
subagents are bounded by the attempt's budget and not separately; the
worktree default is `.claude/worktrees/worktree-ID`, the path Claude's own
worktrees use, added to `info/exclude` when not ignored; the owner is the
Grove binary itself re-executed with `GROVE_ATTEMPT_OWNER`, so the test
binary can be its own owner; a branch that had no worktree keeps the one
`run` made when its records then refuse the run. `GROVE_CLAUDE` names a
fake provider. Left for [G-046](G-046-managed-runs.md): a long-lived
launcher such as the TUI must reap or ignore the owner it spawns (the CLI
releases it), and background grandchildren of a provider that exited
normally are not killed by Grove.

Verification at `30b9c20` and again on the candidate's tree (which adds
one `-short` guard in a test and these records), uncached: `gofmt -l .`
clean, `go vet ./...` ok, `go run ./cmd/grove check` → `OK: 98 records`,
`go test -count=1 -timeout 120s ./...` ok in every package
(`internal/attempt` 4.6 to 4.9 s alone, near the budget; 9 s inside the
parallel suite; 0.7 s with `-short`). Every link written
here, in G-100, G-101 and G-102 resolves in this checkout.

## Next

In Review. The candidate is the evidence commit named in `candidate`; the
branch tip adds only this status change. To judge it, from a checkout of
`worktree-G-045`:

```sh
go run ./cmd/grove context G-045 --include grove/G-100-g-045-durable-attempt-plan.md
go run ./cmd/grove show G-102
git diff --stat 5ae87d5..HEAD
go test -count=1 -timeout 120s ./...
# demo with a fake provider, in a disposable repository, with a binary built from this branch:
go build -o /tmp/grove-bin/grove ./cmd/grove
printf '#!/bin/sh\n[ "$1" = --version ] && { echo fake; exit 0; }\necho "{\"type\":\"system\",\"subtype\":\"init\",\"model\":\"fake\"}"\nsleep 30\necho "{\"type\":\"result\",\"subtype\":\"success\",\"session_id\":\"x\"}"\n' > /tmp/fake-claude && chmod +x /tmp/fake-claude
mkdir -p /tmp/grove-demo && cd /tmp/grove-demo && git init -q -b main && PATH=/tmp/grove-bin:$PATH grove init >/dev/null && PATH=/tmp/grove-bin:$PATH grove new work "Demo" && git add -A && git commit -qm demo
GROVE_CLAUDE=/tmp/fake-claude PATH=/tmp/grove-bin:$PATH grove run G-001 --budget 1 --permission-mode auto   # then exit this shell
PATH=/tmp/grove-bin:$PATH grove attempts && PATH=/tmp/grove-bin:$PATH grove attempt G-001.TIMESTAMP && PATH=/tmp/grove-bin:$PATH grove stop G-001.TIMESTAMP
# a real bounded run: the same with GROVE_CLAUDE unset, --budget 2, in a repository whose AGENTS.md says grove is on PATH
```

Approve: `go run ./cmd/grove approve G-045 "VERDICT"` in this checkout,
then `go run ./cmd/grove integrate G-045` in `main`'s checkout. Feedback:
`go run ./cmd/grove feedback G-045 "TEXT"` here. Then
[G-046](G-046-managed-runs.md) can be prepared on this package.
