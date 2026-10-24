# Implementation: Revising a change re-plans its cycle

<!-- Progress on change 261023-bgtcb, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose re-plans the cycle of a revised member

Status: in progress

- `internal/cli/instructions.go`: `instructionsPropose`, revising a
  change, finds the cycle not done in the working tree that lists it and
  passes it, with its file, to `instruct.Revise`.
- `internal/instruct/instruct.go`: `Revise` takes that cycle or none;
  given one, the task line says it re-plans the cycle, the cycle document
  is written from `metadata/cycle.yaml` with `planned` the instructions'
  target commit and `created` the cycle's, and the file is printed under
  "## The cycle", as `Replan` does.
- `internal/instruct/requires/propose.md`: revising a member, re-plan its
  cycle on this branch: ask whether the change needs a person as planning
  does, set `planned` from the template, change nothing else, never write
  or remove `unattended`, and tell the owner the revision and re-plan must
  reach the target as one commit.
- `internal/instruct/requires/plan.md`: the finish for a member that needs
  a person names revising it with `/grove-propose`, which re-plans the
  cycle on the same branch, and one merge.
- `internal/cli/walkthrough_test.go`: `TestCycles` checks that `instructions
  propose <member>` of a marked cycle prints the re-planning task line, the
  cycle document with `planned: <main's commit>`, the rule never to write
  or remove `unattended` and the cycle's file; that a change in no cycle
  gets no cycle; and that `instructions plan` names the shorter way.

### T-02 — The stale check counts a member changed with its plan

Status: done

- `internal/cli/check.go`: `changedSincePlan`, which `stalePlans` and
  `runnable` in `internal/cli/run.go` share, reads the target's
  first-parent history since `planned` once per cycle with `git log
  --first-parent --diff-merges=first-parent --name-only`, so a merge commit
  counts as one change against its first parent, as a squash does; a
  member is stale when some commit there changed its folder or a named
  file without also editing the cycle file. Its messages are unchanged.
- `internal/cli/walkthrough_test.go`: `TestReplan`, on the fixture, with a
  marked cycle of one member: the revision and a re-plan setting `planned`
  to main's commit, squash-merged, leave `grove check` clean and
  `eligible`, which `grove run` and `grove run --retry` select by, picks
  the cycle; as two commits merged with `--no-ff`, the same; as two
  commits fast-forwarded onto main, `grove check` reports the cycle stale
  with its old message and `eligible` picks nothing; a member revised in a
  merge and the cycle marked in a later commit is reported stale. The
  existing stale-plan tests in `TestCycles`, `TestMovedFiles` and
  `TestRun` pass unchanged.
- The pause after an attempt is not replayed with a stub session here;
  `--retry` starts any attempt `eligible` or `runnable` allows, which the
  test checks after the squash merge.
