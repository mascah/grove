# Implementation: grove claim and the project value read the target's schemas

<!-- Progress on change 261122-pm1d7, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Read the target's schemas for claim and the project value

Status: done

- `internal/cli/project.go`: `project.at(commit)`, which `claim` uses, now
  loads the commit's documents with the schemas in that commit's
  `<root>/schemas` (`fs.Sub` of its tree). `project.agreed()` keeps the
  working tree's schemas, for `instructions`, `check`, `review` and the
  plain `changes` list (decision 2).
- `internal/cli/serve.go`: `server.value`, behind both `/api/project` and
  `grove changes --json`, takes its schemas from the target tree it already
  reads.
- `internal/cli/walkthrough_test.go`: `TestTargetSchemas`. A clone made
  before schema `lite` lands on the target, standing on another branch,
  with A claimed and B sharing A's `record-editing` capability through
  `lite`. Its `grove changes --json` gives B `claimable` false with a
  verdict naming `grove/specs/record-editing/spec.md` and A's attempt, and
  `grove claim B` is refused with the same message and records nothing.
  Without the fix both checks fail (B "would claim it", and the claim
  printed B's ID with status 0).
- Acceptance 2 is held by `TestSchemas` (`instructions propose --schema
  short` for a schema only in the working tree), and acceptance 3 by
  `TestServe` (`changes --json` equals `/api/project` but for `read`). Both
  pass unchanged.
- `just check` passes: 4709 Go, 2595 test, 411 page lines, 8227
  instruction words.

Also changed: `internal/instruct/instruct.go` — the comment on `Project.Schemas` said they were always the working tree's, which is no longer true.

Check: `internal/cli/claim.go` needed no edit. `claim` already read the target through `project.at`, so fixing `at` fixed it.
Check: `claim.go`, `serve.go` and `walkthrough_test.go` changed on main since base baea991, in a4ad6e2 (claim a change by name). I built on them as they are now. a4ad6e2 added `busy` and the loose-change verdicts in `server.value`, which are the code that missed the overlap, so its diff is the code this fix corrects. Nothing in it conflicts with the fix. Before archiving, `grove check` reported these files changed since base. I did not re-base: the diff changes nothing the change says.
Check: the plain `grove changes` list, without `--json`, still reads the working tree's schemas through `agreed`. It shows no overlap verdicts, and the change scopes only `--json` and `/api/project`.

## Scenarios

The change has no change spec. The scenarios of the requirement it restores:

- runs "One capability, different requirements", where the change's schema is only on the target: internal/cli/walkthrough_test.go TestTargetSchemas, passed
- serve "A change that overlaps one in flight", where the change's schema is only on the target: internal/cli/walkthrough_test.go TestTargetSchemas, passed

## Reconciliation

None: the change has no change spec, as what is agreed does not change.
