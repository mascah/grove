---
id: "G-150"
type: work
title: "Launch attempts only where the worktree holds the entrypoints init wrote"
status: proposed
created: "2026-09-25T19:06:01Z"
updated: "2026-09-25T19:08:15Z"
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

## Next

Proposed 2026-09-25 from the owner's self-containment review of `main`
`001b271`. [G-110](G-110-external-preview.md) depends on it: its acceptance
4 executes a bounded assignment in a disposable repository from the preview
documentation alone. Assign with `/grove-work G-150`; a plan is optional,
since the change is one guard in `Start`, three documents and tests. If the
owner prefers a missing reviewer to refuse rather than warn, say so at
assignment.
