---
id: "G-104"
type: review
title: "G-046 managed runs review"
status: current
created: "2026-09-23T04:04:15Z"
updated: "2026-09-23T04:14:44Z"
work: ["G-046"]
examined: "8385f51"
---

# G-046 managed runs review

Evidence for [G-046](G-046-managed-runs.md) against plan
[G-103](G-103-g-046-managed-runs-plan.md). Independent review by a
read-only reviewer subagent in this harness (Claude Code, the `reviewer`
agent), which ran builds, tests and scratch tests in a throwaway copy
and edited nothing. A review is evidence, not approval.

## Independent review, round one (examined `6518342`)

Ten findings; dispositions in `3bde002`:

1. Consequential: `R` could run an attempt on an unrelated branch and in
   its checkout. The shown version is the first current one in Inspect's
   source order, which sorts committed branches by name, so identical
   bytes on `feature` and `main` made `feature` the launch branch. The
   model tests missed it because their fixture kept `main` first. Fixed:
   a branch is passed only when the target does not hold the current
   state (`OnTarget`), or without a target when this checkout's revision
   differs; `TestLaunchPlace` uses Inspect's order.
2. Consequential: this checkout was taken as the live source with locator
   `.`, which is the main worktree, so Grove opened in a linked worktree
   sent the main worktree's revision as `Expect` and Start refused every
   launch as changed. Fixed: the source whose Git directory is the one
   the board was read from; `TestLaunchPlace` opens from a linked
   worktree.
3. Consequential for acceptance 3: `candidate ready` read the result's
   record, which comes from the worktree's files, so a provider that wrote
   `review` and then failed before committing read as ready. Fixed:
   ready needs the branch's tip, as the board read it, to hold that
   record revision (falling back to a clean worktree for a branch the
   board did not read); otherwise the outcome adds "its record says
   review, uncommitted". A first fix that required a clean worktree was
   refuted by the pseudo-terminal scenario: an earlier stopped attempt's
   untracked file blocked a committed candidate.
4. Minor: every past attempt read as `waiting on question` while a
   question blocked the work. Fixed: stopped and failed come first, and a
   wait belongs only to the work's latest attempt.
5. Minor: the open attempt's poll scanned all of `events.jsonl` every
   2 s (610 ms on 118 MB). Fixed: the board's read skips that scan
   (`ShowContext(..., false)`); the activity is the bounded tail.
6. Minor: an inspection already running when an attempt ended was kept,
   and could predate the attempt's last commit. Fixed: restarted.
7. Note: a failed pseudo-terminal scenario could leave the owner and the
   fake provider running. Fixed: the harness kills every attempt's owner
   and process group on the way out.
8. Note: `o` from an attempt did not return to it. Fixed.
9. Note: the plan's tag glyph and polling description differed from the
   code. Fixed in G-103.
10. Note: the hints offered `R` on done work. Fixed.

Checked and found sound in round one: the separate attempts read (one in
flight, stale replies ignored, collected at quit, never touching the Git
read slot), the tick's bounds, the prompt state machine, the result
screen's Esc targets, escaping of activity and the final report,
`Request.Expect`, `go owner.Wait()`, `List` without events, `ShowContext`
cancellation, `ReadActivity` bounds and UTF-8 clipping, and the CLI's
`attempt` output byte for byte.

## Independent review, round two (examined `3bde002`)

Eight of the ten round-one findings confirmed closed, fix 3 found to
cause a new problem, and three findings; dispositions in `8385f51`:

1. Consequential: readiness compared the result's record with the
   branch's current tip, which `approve` and `feedback` then move, so the
   attempt that produced a candidate read "ended without a handoff … with
   no candidate; … uncommitted" exactly while the owner judged it, and
   after feedback the earlier attempt was misreported. Fixed: the owner
   records at exit whether the record's file was committed
   (`result.json` `record_uncommitted`, from `git status --porcelain` on
   the record's path), readiness uses only that, and the wording names
   the candidate; `attempt` prints ", uncommitted" on the record line when
   it applies. `TestRunToResult` covers a clean record and one edited
   before a failure.
2. Minor: the terminal checks' cleanup could SIGKILL reused pids of
   finished attempts. Fixed: only attempts whose owner holds the lock.
3. Minor: a fresh launch from inside `worktree-ID`, or with that branch
   checked out elsewhere, was refused with advice to pass a CLI flag the
   board lacks. Fixed: a live checkout of `worktree-ID` is passed as the
   worktree; `TestLaunchPlace` covers one at another path.

## Independent review, round three (examined `8385f51`)

All three round-two findings confirmed fixed; one consequential finding
about the evidence and three notes. This was the third and last round the
cap allows; the dispositions after it were checked by the implementing
session, not by the reviewer:

1. Consequential for the evidence: `attempt_lifecycle` failed 1 run in 8
   (2 in 15) from a clean copy of `8385f51`: Esc pressed on a launch's
   result screen while the board was being re-read is ignored by design,
   and the scenario then waited forever. The implementing session had
   seen the same failure under the full suite's load and fixed it in
   `bca6397` (wait for "The board has been re-read." before each Esc
   after a launch or stop). Self-checked: 16 of 16 runs of the whole
   terminal suite passed at `79095a6`, four at a time, and the full
   suite passed.
2. Note: that ignored Esc gave no message. Fixed in `79095a6`: the banner
   says Esc waits for the re-read.
3. Note: a failing `git status` at exit left `record_uncommitted` false,
   so an unverified record could read as ready. Fixed in `79095a6`: an
   error counts as uncommitted and is logged.
4. Note: the terminal harness's cleanup could raise on an attempt
   directory without its lock yet. Fixed in `79095a6`.

The reviewer's verdict at round three: acceptance items 1 to 4 met by
code and tests, item 5 met in part (terminal checks cover launch,
observe, exit, reconnect, wait, stop, continuation and high-volume
output; review only as reaching the review detail; the owner's own
review in a terminal is outstanding). Not verified by the reviewer:
Linux, a real `claude` provider, `-race` in rounds two and three, and the
changes after `8385f51`, which are the self-checked ones above.
