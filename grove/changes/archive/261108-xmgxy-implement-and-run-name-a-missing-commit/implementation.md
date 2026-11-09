# Implementation: Implement and claim name a missing base or planned commit

<!-- Progress on change 261108-xmgxy, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement and claim name a missing base or planned commit

Status: done

- `internal/cli/check.go`: `absent` returns check's line, `<root>/<file>:
  <field>: <commit> ...`, as an error, from the existing `missing`
  sentence (delegated).
- `internal/cli/instructions.go`: `instructionsChange` returns `absent`
  for the change's `base` when `moved` fails, and for the cycle's
  `planned` when `sincePlan` fails; `instructionsCycle` runs `sincePlan`
  for a cycle not done with a `planned` and refuses the same way before
  printing.
- `internal/cli/claim.go`: `runnable` returns `absent` for the cycle's
  `planned` when `changedSincePlan` fails, before anything is recorded;
  `eligible` still skips such a cycle silently.
- `internal/cli/walkthrough_test.go` `TestShallowClone`: adds change B in
  no cycle and marks the cycle `unattended`; in the shallow clone,
  `instructions implement` B, A, A `--unattended`, the cycle, and
  `claim <cycle-id>` each exit 1 with only `grove: ` and check's line;
  `claim` given no cycle says there is nothing to claim, and no `grove/`
  branch or record exists in the clone or its origin; in the full
  repository, B's `base: 0000000` committed gives the "is not a commit in
  this repository" line. Failed before the fix, with Git's `rev-parse`
  and `diff` errors; passes with it.
- `dev/walkthrough.md` step 6 and `docs/guides/run-grove-in-ci.md`
  "What to expect" gain the implement and claim lines.
- On the fixture (the proposal's Reproduction, A in cycle K and B in
  none, Grove built from this branch): in a `--depth 1` clone, each
  command above printed only `grove: ` and check's shallow-clone line and
  exited 1, and `grove claim` recorded nothing and created no branch.
  With B's `base` and K's `planned` set to `0000000` and committed in the
  full repository, each printed the "is not a commit in this repository"
  line.
- `just check`: passes; 4386 Go, 2044 test, 410 page lines, 8613
  instruction words.

## Scenarios

None: the change has no change spec, since the owner kept this observed
behavior (decision 3). Its acceptance is shown by
`internal/cli/walkthrough_test.go` `TestShallowClone`, passed.

## Reconciliation

No change spec; no maintained spec changes.
