---
id: "G-046"
type: work
title: "Launch and inspect managed attempts from the TUI"
status: review
created: "2026-09-21T00:54:16Z"
updated: "2026-09-23T04:14:44Z"
kind: feature
size: medium
priority: 4
depends_on: ["G-044", "G-045"]
relates_to: ["G-035", "G-043"]
formerly: "W-028"
candidate: "61edd53390c4459462e79ededf5a25eb0644a1f7"
---

## Outcome

Launch one eligible work item from the board, inspect its durable attempt,
return after closing the TUI, and review the resulting candidate in Grove.

## Scope and bounds

Implement is an explicit assignment action, not a side effect of changing
status. Show a runs list and attempt detail or an integrated detail section,
using G-045's ownership/events rather than making the TUI own the process.
Present progress, waiting/failure, budget and Stop with linked logs as needed.
Retain the final report and reviews when ephemeral activity scrolls away.

## Acceptance

1. Eligible selection creates exactly one bounded assignment; duplicate starts,
   missing readiness and changed inputs produce useful actionable outcomes.
2. Closing and reopening the TUI reconnects to the same attempt. Stop is a
   separate explicit action and preserves partial work/evidence.
3. Running, waiting, failed/interrupted and candidate-ready outcomes are honest;
   a successful process exit alone cannot fabricate Review readiness.
4. Review feedback can launch a new bounded attempt on the correct candidate
   branch, preserve prior evidence and return work to Active.
5. Owner-reviewed connected workflow and terminal lifecycle checks cover launch,
   observe, exit, reconnect, wait, stop, review and continuation. Test high-volume
   output without flooding the main interface or blocking navigation.

## Evidence

Branch `worktree-G-046` in `.claude/worktrees/worktree-G-046`, from `main`
at `bd2418f`, started from this record at revision `sha256:4dc799d6…`.
Plan [G-103](G-103-g-046-managed-runs-plan.md) (`e1d3693`, which also set
this record active). Commits: `85a36c5` the attempt package, `dbab7c6` the
board, `6518342` the pseudo-terminal scenario and docs, `3bde002`,
`8385f51`, `bca6397` and `79095a6` the review rounds' fixes; the candidate
is the evidence commit named in `candidate`. Review record
[G-104](G-104-g-046-managed-runs-review.md): three rounds of independent
review, the findings and their dispositions.

Decisions (the roadmap's run-experience questions, G-047): attempts get
two screens, a list (`A`) and one attempt (Enter), plus a row in work's
detail naming the latest attempt and its outcome and a `● running` tag on
the card; activity is summarized one short line per event from the last
1 MiB of `events.jsonl`, newest first, below the final report, which is
rendered like a record body behind the same escaping; relaunch after
feedback is `R` on the record, which runs on the candidate's branch.
Budget and permission mode are typed at every launch with no prefill,
since G-045 made a default of either a product choice.

Changed behaviour, against each acceptance item:

1. **One bounded assignment, useful refusals.** `R` on proposed or active
   work asks for a budget (checked) and a permission mode, then runs
   G-045's `attempt.Start` once: exactly one attempt, as `run` makes it.
   Refused up front with what to do: work in review ("judge its
   candidate"), done or abandoned, not work, blocked by an open question,
   or with an attempt running or orphaned ("A shows it, x stops it");
   Start's own refusals show on the result screen as `NOT DONE: …`. A
   changed input is refused by the new `Request.Expect`: the revision of
   the record in the checkout Grove was opened in, as the board read it,
   must still be the one Start reads. The launch runs from this checkout
   (found by its Git directory, so a linked worktree works) and continues
   on the branch holding the current state when the target (or, without
   one, this checkout) does not hold it, in that branch's live checkout;
   otherwise it starts afresh, reusing a live checkout of `worktree-ID`
   wherever it is.
2. **Reconnect and Stop.** The board only reads attempt files, so
   quitting leaves the attempt running and the next session shows the
   same one; `attempt.Start` now waits on the owner in a goroutine, so the
   long-lived board reaps it (G-102's note). `x` asks `y/n` and stops
   through `attempt.Stop`; partial work stays. The pseudo-terminal
   scenario quits while an attempt runs, finds it running in a new session
   with one provider start, stops it there, and finds `partial.txt` kept.
3. **Honest outcomes.** Derived for display, never written: `running`,
   `orphaned`, `interrupted`, and for a finished attempt `candidate
   ready` only when the owner found the record committed in review with a
   candidate at exit (the new `record_uncommitted` in `result.json`, which
   an unverifiable status counts as uncommitted), then `stopped`,
   `failed` (no result event, an error result or a nonzero exit),
   `waiting on question` (the work's latest attempt while a question
   blocks it) and `ended without a handoff`. A clean exit alone never
   reads as ready (`TestAttemptOutcomes`, fifteen cases). Running
   attempts are re-read every 2 s, and an attempt's end re-reads the
   board.
4. **Feedback continuation.** After `f`, the result names `R` on the
   record; `R` then runs on the candidate's branch in its checkout
   (`TestFeedbackContinuesOnTheCandidateBranch`, `TestLaunchPlace`, and
   the scenario's third attempt, which continues on `worktree-G-001` to a
   committed candidate after its question is answered). Earlier evidence,
   reviews and attempts stay; the attempt changes nothing but what the
   headless guide writes, and earlier attempts keep their outcomes after
   later commits.
5. **Connected and terminal checks.** `attempt_lifecycle` in
   `internal/tui/testdata/terminal.py` drives the built binary through a
   pseudo-terminal with a fake provider (`GROVE_CLAUDE`) that writes
   20,000 events: launch, observe the newest activity, page while it
   floods, quit (modes restored, nothing on stdout), reconnect, stop,
   a second attempt that persists a blocking question (`waiting on
   question`, then `R` refused), answer it, and a third attempt that ends
   in review (`candidate ready`, the detail becomes the review view).
   Model tests cover the same paths, escaping of provider text, and every
   row's width at 40, 80 and 160 columns with 200 activity lines. The
   owner's judgment in a terminal is Next.

Verification at `79095a6`, uncached: `gofmt -l .` clean; `go vet ./...`
ok; `go run ./cmd/grove check` OK (100 records, this record's review
included); `go test -count=1 -timeout 120s ./...` ok in every package;
`-race` alone on `internal/tui` and `internal/attempt` ok; `-short`:
tui 0.6 s, attempt 0.8 s, cli 2.5 s. The whole terminal suite (ten
scenarios) passed 16 of 16 runs, four at a time. Links in this record,
G-103, G-104, the README and both docs resolve.

Limits: no real-provider run from the board (G-045's two real trials ran
the same `Start`, owner and `Stop`; here the provider is a fake); Linux
not run; the activity summary knows Claude Code 2.1.280's event shapes
and shows nothing for others; Stop of an orphan holds the keys for up to
its grace period; the board's `candidate ready` for an attempt whose
`result.json` predates `record_uncommitted` trusts the record; the final
report is shown only while the result event is within the last 1 MiB.

## Next

In Review. The candidate is the evidence commit named in `candidate`; the
branch tip adds only this status change. To judge it, from this worktree:

```sh
go run ./cmd/grove context G-046 --include grove/G-103-g-046-managed-runs-plan.md
go run ./cmd/grove show G-104
git diff --stat bd2418f..HEAD
go test -count=1 -timeout 120s ./...
# the board with a fake provider, in a disposable repository, with a binary of this branch:
go build -o /tmp/grove-G-046 ./cmd/grove
printf '#!/bin/sh\n[ "$1" = --version ] && { echo fake; exit 0; }\necho "{\\"type\\":\\"system\\",\\"subtype\\":\\"init\\",\\"model\\":\\"fake\\"}"\nfor i in $(seq 1 60); do echo "{\\"type\\":\\"assistant\\",\\"message\\":{\\"content\\":[{\\"type\\":\\"text\\",\\"text\\":\\"step $i\\"}]}}"; sleep 1; done\necho "{\\"type\\":\\"result\\",\\"subtype\\":\\"success\\",\\"result\\":\\"Done.\\"}"\n' > /tmp/fake-claude && chmod +x /tmp/fake-claude
mkdir -p /tmp/grove-demo && cd /tmp/grove-demo && git init -q -b main && /tmp/grove-G-046 init >/dev/null && /tmp/grove-G-046 new work "Demo" && git add -A && git commit -qm demo
GROVE_CLAUDE=/tmp/fake-claude /tmp/grove-G-046   # Enter on the card, R, 1, auto; A, Enter; q; reopen; x, y
# a real bounded run: the same without GROVE_CLAUDE, in a repository whose AGENTS.md says grove is on PATH
```

Approve: `go run ./cmd/grove approve G-046 "VERDICT"` in this worktree,
then `go run ./cmd/grove integrate G-046` in `main`'s checkout. Feedback:
`go run ./cmd/grove feedback G-046 "TEXT"` here.
