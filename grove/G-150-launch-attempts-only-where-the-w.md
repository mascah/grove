---
id: "G-150"
type: work
title: "Launch attempts only where the worktree holds the entrypoints init wrote"
status: active
created: "2026-09-25T19:06:01Z"
updated: "2026-09-25T19:22:56Z"
kind: fix
priority: 1
size: small
relates_to: ["G-040", "G-101", "G-110", "G-134", "G-152"]
---

## Outcome

`grove run` and the board's `R` start a provider only in a worktree that
holds the `grove-work` skill their prompt names, tell the launcher before
any spend when the reviewer definition is absent, and `init` and the
adoption steps say to commit what `init` wrote. Owner intent, review
conversation 2026-09-25: Grove must work in other projects as an installed
CLI without this repository present, and the owner called this class of
dependency a significant blocker to that.

## Constraints

Observed at `main` `001b271` (this checkout), 2026-09-25:

- `Start` in [attempt.go](../internal/attempt/attempt.go) creates
  `worktree-ID` with `git worktree add -b BRANCH DIR HEAD` from the
  launching checkout's HEAD, then composes `claude -p "/grove-work ID
  --interaction headless" …`. It reads `.claude/agents/grove-reviewer.md`
  from the worktree only to record its digest or `none` in `attempt.json`;
  it never looks for `.claude/skills/grove-work/SKILL.md`, and neither
  absence stops the launch. The board's `R` builds the same request and
  launches through the same function.
- In a disposable Git project with one commit, `grove init` from a binary
  built at `001b271` wrote its files, and a worktree added from HEAD held no
  `.claude` directory at all. `init`'s closing note says `Next: grove check`
  and what `AGENTS.md` may add; the README's adoption block ends at `check`;
  the command reference's Init section says "Upgrading is building again and
  rerunning `init`". None says to commit.
- Claude Code's skills documentation, read 2026-09-25: in a linked worktree
  the search for skills stops at the worktree root, and from 2.1.277 a
  worktree with no `.claude/skills/` loads the main checkout's project
  skills. So on this machine (2.1.282) the slash command resolves through
  the launching checkout's uncommitted copy; on an earlier version the
  provider receives the assignment as plain text with no workflow and spends
  the budget on it. The documentation says nothing about `.claude/agents/`,
  so whether the reviewer definition is found the same way is unverified.
  The work guide's step 6 dispatches the reviewer only "where the checkout
  holds" the definition, so a worktree without it runs with no independent
  review and, where the record requires one, cannot hand off.
- No trial exercised this path: [G-082](G-082-portable-bootstrap-review.md)'s
  disposable repositories were committed before their headless runs, the
  eval fixture (`evals/run.py`, `build`) commits right after `init`, and
  [G-134](G-134-bound-an-attempt-at-its-plan-and.md)'s real trial ran in a
  clone of this repository, where the files are committed. G-134 chose to
  record the reviewer's digest or absence as a fact; it did not consider a
  launch whose absence comes from an uncommitted `init`.

Proposed design, labelled proposed:

1. `Start` checks the prepared worktree after `prepareWorktree` and before
   the attempt directory is written, so the refusal comes before any spend
   and, as for the other late refusals, a branch that had no worktree keeps
   the one `run` made. Missing `grove-work` skill: refuse, naming the path
   and saying to commit the files `init` wrote. Missing reviewer definition:
   a warning through `report`, and `none` recorded as today; the guide
   already defines that degraded path, and a project that deliberately
   removed the file keeps running. The board's `R` shows both as it shows
   the other launch messages, since it uses the same function.
2. `init`'s closing note, the README's adoption block and the command
   reference's Init and Attempts sections say to commit what `init` wrote
   before the first attempt; `grove --help`'s `run` text lists the refusal.
3. The runner's fake-provider tests cover the refusal and the warning. No
   file is copied into the worktree: passing the reviewer with `--agents`
   was rejected in G-134 because an interactive `/grove-work` session would
   not get it, and copying uncommitted files onto the work branch would let
   the attempt commit them.

Out of scope: gating the check on the provider's version (the file is what
the guide reads, whatever the harness does), a Codex runner
([G-143](G-143-g-135-codex-eval-row-pattern-pla.md)), and the predecessor.

## Acceptance

1. In a disposable project made by `grove init` with nothing committed after
   it, `grove run ID --budget … --permission-mode …` and the board's `R`
   refuse before starting the provider, name the missing skill path and say
   to commit `init`'s files; after `git add -A && git commit`, the same
   launch starts.
2. With the skill committed and the reviewer definition absent, the launch
   proceeds, the launcher sees a warning naming the path before the provider
   starts, and `attempt.json` and `grove attempt` still show `no reviewer
   definition`.
3. `init`'s closing note, the README's adoption block and the command
   reference's Init and Attempts sections say to commit the written files,
   and `grove --help` lists the refusal; a fresh reader following the README
   alone reaches a working first `run`.
4. Runner and TUI tests cover 1 and 2 with the fake provider; `go vet`,
   `gofmt`, `grove check` and the full uncached suite pass; `internal/attempt`
   stays under five seconds.

## Evidence

Implemented 2026-09-25 by a headless `/grove-work G-150` attempt on
`worktree-G-150` from `main` `28aaf95`, record revision
`sha256:085b9a43…` at start; no plan, as Next allowed. Implementation
`cd3a1c5`, review fixes `b322825`.

- **Guard** (`internal/attempt/attempt.go`, `SkillPath`). `Start` refuses
  when the attempt's worktree would lack `.claude/skills/grove-work/SKILL.md`
  under the project's prefix: for a branch that does not yet exist, by
  `git cat-file -e HEAD_SHA:PATH` before the branch or worktree is made;
  for a reused or existing branch, on disk after `prepareWorktree` and
  after the branch's own review and blocking-question refusals. Both come
  before the attempt directory and the owner. The message names the path
  and says to commit the files `grove init` wrote. An absent reviewer
  definition is a `warning:` fact through `report` before the attempt
  directory, and `attempt.json` still records `none`.
- **Deviation from the proposed design, bounded and technical.** The
  design checked only after `prepareWorktree`. That keeps a fresh
  `worktree-ID` at a HEAD without the skill, which the next launch would
  reuse after the commit and refuse again, contrary to acceptance 1. So a
  new branch is checked in HEAD first. An existing branch without a
  worktree still keeps the one `run` made, as documented.
- **Acceptance 1.** `TestRefusals`: a skill that is on disk but not in
  HEAD is refused with no branch made. An existing branch without the
  skill is refused in its checkout. After the commit, the launch passes
  every check and reaches the owner. By hand, with a binary built from
  `cd3a1c5`: a disposable project made by `grove init`, with only a work
  record committed, was refused (`…SKILL.md is not committed at HEAD…`,
  exit 1, no branch). After the commit it started (exit 0). The board's
  `R` is covered by `terminal.py` `attempt_lifecycle`: the launch is
  refused with `NOT DONE: .claude/skills/grove-work/SKILL.md is not
  committed at HEAD`, with no start and no worktree. After the commit,
  the same launch starts.
- **Acceptance 2.** `TestRunToResult` asserts the exact warning fact and
  `Reviewer == "none"`. By hand, `grove run` printed the warning before
  `attempt: … started`, `requested:` said `no reviewer definition`, and
  `attempt.json` held `"reviewer": "none"`. `terminal.py` asserts that the
  board draws the warning. Open (G-157 finding 2): the board draws it only
  once `Start` returns, after the provider has started.
- **Acceptance 3.** `init`'s closing note (asserted in
  `internal/cli/init_test.go`), the README's adoption block (it now
  commits the paths `init` writes and names `grove run`), the command
  reference's Attempts and Init sections, and `grove --help`'s `run` text
  all say to commit, and `--help` lists the refusal.
- **Acceptance 4.** Run at `b322825`: `go vet ./...` was clean, `gofmt -l .`
  was empty, `go run ./cmd/grove check` reported `OK: 150 records`, and
  `go test -count=1 -timeout 120s ./...` passed every package. `python3
  internal/tui/testdata/terminal.py` passed 11 of 11. Limit (G-157 finding
  3): `go test -count=1 ./internal/attempt` alone took 5.06 to 5.22 s at
  `b322825`, and the same at base `28aaf95` (a `git archive` copy, three
  runs each, load average about 4). So the package was already over five
  seconds before this change.
- **Review.** [G-157](G-157-g-150-skill-guard-review.md), two rounds.
  Findings 1, 4 and 5 were fixed. Finding 2 is left for the owner, and
  finding 3 is inherited.

## Next

Proposed 2026-09-25 from the owner's self-containment review of `main`
`001b271`. [G-110](G-110-external-preview.md) depends on it: its acceptance
4 executes a bounded assignment in a disposable repository from the preview
documentation alone.

In review 2026-09-25 on `worktree-G-150` (base `28aaf95`). The candidate
is the evidence commit this record names in `candidate`. For the owner:

- Judge the guard and its order: `git diff 28aaf95 b322825 --
  internal/attempt/attempt.go`.
- Decide on G-157 finding 2. The board draws the reviewer warning after
  the provider has started. Either accept that as the design's "as it
  shows the other launch messages", or capture new work to stream launch
  facts into the board, or make a missing reviewer refuse. The record left
  that choice to you at assignment, and none was given, so it warns.
- Demo: in a scratch `git init` repository, run `grove init`, then
  `grove new work "Try"`. Commit only `grove.yaml` and `grove`. Then
  `GROVE_CLAUDE=/path/to/fake grove run G-001 --budget 1 --permission-mode
  auto` is refused, and after committing `.claude .agents` it starts.

Integration, as given:

```sh
go run ./cmd/grove approve G-150 "VERDICT"   # in this worktree
go run ./cmd/grove integrate G-150           # in the main checkout
```
