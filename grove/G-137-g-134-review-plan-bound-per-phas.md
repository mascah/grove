---
id: "G-137"
type: review
title: "G-134 review: plan bound, per-phase effort and the grove-reviewer definition"
status: current
created: "2026-09-24T22:25:45Z"
updated: "2026-09-24T22:25:59Z"
work: ["G-134"]
examined: "ccbb73ce87338c3b8e397eae0bbdb03ad23ed237"
---

## Examined

An independent reviewer subagent of this headless session (Claude, a
general-purpose agent on Opus, no edits) examined
[G-134](G-134-bound-an-attempt-at-its-plan-and.md) against its acceptance and
[plan G-136](G-136-g-134-plan-plan-bound-per-phase.md). Its instructions were
the body of the new [`grove-reviewer`](../.claude/agents/grove-reviewer.md)
definition, read from the file: the definition itself was not dispatchable
in this session, which started before the file existed, so the review ran at
the session's effort, not the definition's `high`. Round 1 was at `a84d01d`
(`git diff 98ce628 a84d01d`) and round 2 at `ccbb73c`
(`git diff a84d01d ccbb73c`). It ran `go test -count=1 -timeout 120s` on
`./internal/attempt`, `./internal/cli` and `./internal/tui` (the
pseudo-terminal script included), `go vet ./...`, `gofmt -l .`,
`grove check`, `grove guide work`, and `grove context G-134 --until plan`.

## Findings

Round 1 (`a84d01d`), no blockers:

1. Medium: the board showed `plan ready` for any clean, bounded end, even
   with the record uncommitted, the worktree dirty or the record unreadable,
   and offered `R`, which is the approval of a plan that might not be
   committed.
2. Low to medium: the adapters and init's `assignmentData` said to pass
   the bound to commands, and `grove context` refuses `--until` (verified
   by running).
3. Low: the Lifecycle paragraph said preparation comes before `active`,
   which is untrue for work already active; step 4 called the relaunch "the
   caller's reading" where the record and board doc say approval.
4. Low: the reviewer is read-only by instruction for Bash; only the
   editing tools are disallowed.

## Disposition

Round 1: (1) fixed in `ccbb73c`: plan ready also needs a readable,
committed record and a clean worktree, and anything else ends without a
handoff, with three regression cases; one limit remains, reported in
G-134: an attempt that ignores the bound and commits the record as `active`
still reads as plan ready, since the launch does not record the status it
found. (2) Fixed: "pass IDs and mode to commands ...; the bound is for the
guide, not an argument to any command." (3) Both reworded. (4) Kept as a
limit: the brief forbids writes.

Round 2 (`ccbb73c`): every round-1 disposition confirmed; no new
consequential finding. One choice was noted and kept: `!Dirty` counts
untracked files, so a stray scratch file makes a bounded end read as without
a handoff, which is stricter than the candidate rule and on the safe side
given the guide's "no change outside records".

Knowledge: no contradiction with the settled terms or accepted decisions.
The amended [G-055](G-055-preparation.md) and the guide's Lifecycle and step
4 agree that the bound ends the caller's mandate and is not a gate. The
Attempt term's "Produces a candidate" is looser than a bounded attempt, as
it already was for an attempt that stops on a question. This review is
evidence, not approval.
