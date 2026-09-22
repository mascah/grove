---
id: "G-075"
type: plan
title: "Trial plan for the interactive loop on a list status filter"
status: current
created: "2026-09-22T14:56:44Z"
updated: "2026-09-22T14:57:26Z"
work: ["G-039"]
---

## Design

Prepared 2026-09-22 on `worktree-G-039` from main `5f07c94` against
[G-039](G-039-interactive-loop.md). The owner selected the real change in that
session on 2026-09-22: a status filter for `grove list`. Observed at `5f07c94`:
`grove list --status proposed` is refused with `unknown option --status`, and
`list` takes no filter (`internal/cli/cli.go:156`); finding open work today
means piping `list` through `grep`. Harness for the trial: Claude Code 2.1.278
with model Fable 5.1, go 1.26.2, git 2.55.0, macOS Darwin 25.6.0.

The trial is the sequence below, run in fresh sessions that the owner starts.
This plan only sequences them and says what to record; it is not the change's
plan, which the shaped work's own session prepares. Each observation is
labelled one of: real trial (a fresh session on this repository), simulation
(a disposable clone or fixture), or unexercised.

## Steps

1. **Shape** (owner, fresh interactive session in the main checkout):
   `/grove-shape add a status filter to grove list`. Expected: one proposed
   work record with outcome, scope, acceptance and Next, related to nothing
   invented; no implementation, no status promotion. Record: the invocation,
   whether the session found the brief, `list`, and the code without being
   told, the records and knowledge it created or changed with revisions,
   where they were committed, and friction.
2. **Implement** (owner, fresh session): `/grove-work G-NNN` for the shaped
   record. Expected: `worktree-G-NNN` from main, preparation proportional to
   the size (a plan record or "no plan needed" in Next), `active` set through
   the CLI, only genuine human decisions escalated, an independent review
   record with its `examined` commit, and a handoff into `review` with a
   candidate. Record: escalations, checkpoints, the candidate, the review IDs.
3. **Judge and integrate** (owner, from files and the CLI only): in
   `.claude/worktrees/G-NNN`, `go run ./cmd/grove context G-NNN`, confirm
   `git diff --stat CANDIDATE HEAD` touches only the record, then either write
   feedback into Next and set `active`, or accept: on main
   `git merge --ff-only worktree-G-NNN` and `update ... --set status=done`,
   committed there. Record: whether chat was needed, the verdict verbatim.
4. **Continue** (owner): resume `/grove-work G-039` in a fresh session. It
   must find this plan and G-039's checkpoint and carry on without the chat
   that wrote them. If the feedback case (step 3) and the missing-human case
   did not arise, that session exercises them in a disposable clone reached by
   an explicit `--project` path, for example a headless `/grove-shape` on a
   topic that needs a product choice, and labels the result a simulation.
5. **Record** (that session): a review record (`new review`, `work` G-039,
   `examined` the trial's candidate) with the evidence per acceptance item;
   lessons written into `docs/work-execution.md`, `docs/work-shaping.md`, or
   G-036's Next as their owner; then G-039 handed into `review`.
