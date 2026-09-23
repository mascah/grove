---
id: "G-102"
type: review
title: "G-045 durable attempt review"
status: current
created: "2026-09-23T02:37:36Z"
updated: "2026-09-23T02:53:46Z"
work: ["G-045"]
examined: "30b9c20"
---

# G-045 durable attempt review

Evidence for [G-045](G-045-durable-attempt.md) against plan
[G-100](G-100-g-045-durable-attempt-plan.md). Two rounds of independent
review by a read-only reviewer subagent in this harness (Claude Code
2.1.280), then two bounded real-provider trials by the implementing session.
A review is evidence, not approval.

## Independent review, round one (examined `69b9d9e`)

Fifteen findings; dispositions in `9ae71ac` and `b4a93a7`:

1. Consequential: a project below the checkout's top skipped the branch
   guards, read the wrong record state and never resolved the inputs
   check. Fixed: the prefix is recorded and used for the guards, the
   provider's directory, the record state and `git show`; covered by
   `TestInputsChanged`, which now runs under `sub/`.
2. Consequential: a Stop arriving before the owner's signal handler killed
   the owner without a result while `stop` reported success. Fixed: the
   handler is installed first, the owner tells the launcher it is ready
   through a pipe, and a signal already pending before the provider starts
   writes a stopped result without starting it; tested with a SIGTERM
   immediately after `Start`.
3. Minor: one refusal (an existing branch without a worktree) comes after
   the worktree is added. Documented rather than changed; the checkout is
   the branch's own and is kept.
4. Minor: the liveness probe took an exclusive lock, so readers could
   misread each other. Fixed with a shared probe and a blocking exclusive
   launcher lock.
5. Minor: `--budget` accepted NaN, Inf and hex. Fixed in `Start`.
6. Minor: the branch's record status and candidate were printed
   unescaped. Fixed with `visible()`; record diagnostics are printed too.
7. Minor: the environment-scrub test set nothing. Fixed; it sets the
   session and Git variables and asserts only `CLAUDE_CONFIG_DIR` passes.
8. Minor: untracked partial work read as clean. Fixed.
9. Minor: an attempt directory without `attempt.json` broke `attempts`.
   Fixed: built under a temporary name and renamed; skipped if incomplete.
10. Minor: the Git variable list was copied from `repo`. Fixed: exported.
11. Minor: the owner's lock check could not fail. Fixed: `flock` on fd 3
    must succeed.
12. Minor: a launcher killed before recording the owner pid made Stop
    refuse forever. Fixed: falls back to the owner's own record.
13. Minor: Stop of an orphan could wait unbounded after SIGKILL. Fixed.
14. Notes: test loops without deadlines (fixed), the package at the 5 s
    budget (three tests merged into siblings), `--version` without `-short`
    (kept, trivial), Git through `repo.Command` throughout (confirmed).
15. Notes for G-046: zombie owners under a long-lived launcher, no timeout
    on `claude --version` (added, 30 s), grandchildren after a normal exit,
    usage text refusals (added).

Verdict at round one: items 1, 2 and 5 met by code and tests, 3 mostly,
4 only for a top-level project; the consequential findings above closed
that gap.

## Independent review, round two (examined `b4a93a7`)

Twelve of the fifteen round-one findings closed, finding 3 accepted as
documented, and six new findings; dispositions in `30b9c20`:

1. Consequential regression: the round-one prefix fix compared Git's real
   path with the caller's path, so a project reached through a symlink
   (any `/tmp` or `/var` path on macOS) was refused. Fixed: the prefix
   comes from `git rev-parse --show-prefix`; the test reruns through a
   symlinked root.
2. Minor: the "owner exited while starting" check could not fire, because
   the launcher still held its own descriptor of the lock. Fixed: the
   launcher closes it right after the spawn; tested with an owner that
   exits at once (`Start` reports it, the attempt lists as interrupted).
3. Minor: the fake's own commit could start Git's detached maintenance
   (the recorded `maintenance.lock` flake). Fixed with `-c
   maintenance.auto=false`.
4. Note: a bad `--budget` had become exit 1. Fixed: `attempt.ValidBudget`
   is checked at parse time too, exit 2.
5. Note: a test ignored `Start`'s error. Fixed.
6. Note: the launcher's wait for the owner's readiness was unbounded.
   Fixed: 10 s, then the lock decides.

Verdict at round two: items 1, 2, 3 and 5 met by code and tests, 4 met
except for symlinked paths, which the regression above closed.

## Independent review, round three (examined `30b9c20`)

All six round-two findings closed, no regression in the fix commit, and
one minor note: the dying-owner case in `TestRefusals` starts an owner
process without a `-short` skip. Applied in the evidence commit as a
`testing.Short()` return with the reason, after the third round, so it is
self-checked rather than re-reviewed (the review cap). Verdict: every
acceptance item met by code and tests as far as they show it, with budget
exhaustion covered only by a fake result subtype; the reviewer recommends
the candidate on code and tests.

## Real-provider trials (binary from `cab470c`, Claude Code 2.1.280)

In a disposable repository under the session scratchpad, initialized with
`grove init`, whose `AGENTS.md` says `grove` is on PATH and the checkout
is the work branch, with a one-line script and README:

- **G-001 "Greet the user by name"**, `grove run G-001 --budget 4
  --permission-mode auto` from `sh -c` that exited immediately. The owner
  ran under launchd in its own session and process group; `attempts` from
  a new process showed `running`; `attempt` showed the init event (model
  `claude-opus-5-5[1m]`, permission mode `auto`, 26 tools, five
  capabilities) while it ran. Finished in 62 s, exit 0, result `success`,
  9 turns, $0.4971, 0 permission denials, 53 events (20 system, 18
  assistant, 11 user, 3 rate limit, 1 result), no oversized or malformed
  line. On the branch: three commits, the record in `review` with candidate
  `a083a9f`, the tip differing from it by the record alone; `./hello.sh
  Ada` prints `hello Ada`; the headless session had run an independent
  review subagent of its own and written the handoff. A second `run` of
  G-001 was then refused, since the branch's record is in review.
- **G-002 "Add a farewell script"**, `--budget 1`, reconnected with
  `attempt` 20 s after the init event, then `stop`: the owner logged
  SIGTERM, sent SIGINT, the provider exited with code 0 within a second
  and an `error_during_execution` result (`is_error: true`, 6 turns,
  $0.2259); `stopped: true`; the worktree kept an uncommitted `bye.sh`,
  README and record edit, HEAD at the base.

Not observed on the real provider: budget exhaustion, SIGKILL after the
15 s grace, owner loss, machine restart.
