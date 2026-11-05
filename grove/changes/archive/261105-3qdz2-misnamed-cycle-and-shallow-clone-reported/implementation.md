# Implementation: A misnamed cycle and a shallow clone are reported plainly

<!-- Progress on change 261105-3qdz2, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Report a cycle file named without an ID

Status: done

- Test first: `TestMisnamedCycle` in `internal/cli/walkthrough_test.go`
  plans `261105-proto-unknown-id-cycle.yaml` (IDs leave out `o`) and a
  well-named `261105-pr0t0-cycle.yaml`. Before the fix it failed as the
  proposal observed: `check` printed "No problems in grove/." with status
  0, and `changes` printed "in cycle " with no ID.
- `internal/doc/doc.go`: loading a cycle whose file name starts with no
  ID reports "the file name must start with a cycle ID, such as
  261015-k3m9q", as a misnamed change folder is reported.
- `internal/cli/changes.go`: names no cycle without an ID, matching
  `--json`. The well-named cycle is still shown (`in cycle 261105-pr0t0`).
- `dev/walkthrough.md` "Planning cycles" step 6: rename the 0.4 cycle's
  file to `261105-proto-0-4.yaml`, run `grove check` and `grove changes`,
  rename it back.
- `just check`: passed.

### T-02 — Say a missing commit is missing from a shallow clone

Status: done

- Test first: `TestShallowClone` in `internal/cli/walkthrough_test.go`
  clones with `git clone --depth 1` a repository whose change's `base` and
  cycle's `planned` are older than its one commit. Before the fix, check
  said both were "not a commit in this repository".
- `internal/gitx/gitx.go`: `Shallow` reports whether a repository is a
  shallow clone (`git rev-parse --is-shallow-repository`).
- `internal/cli/check.go`: `missing` gives "<field>: <sha> is not in this
  shallow clone; fetch its full history, such as with git fetch
  --unshallow" in a shallow clone, and the old message otherwise, for both
  `base` and `planned`. The same test shows a `base` of `0000000` in the
  full repository is reported as before.
- `docs/guides/run-grove-in-ci.md` "What to expect": shows the message
  and says the checkout needs `fetch-depth: 0`.
- `dev/walkthrough.md` "Planning cycles" step 6: a `--depth 1` clone of
  the fixture and its `grove check`.
- `just check`: passed (4050 Go lines, 1752 test lines, 8426 words).

## Scenarios

- cycles "A letter IDs leave out": internal/cli/walkthrough_test.go TestMisnamedCycle, passed

## Reconciliation

- cycles: applied to grove/specs/cycles/spec.md
