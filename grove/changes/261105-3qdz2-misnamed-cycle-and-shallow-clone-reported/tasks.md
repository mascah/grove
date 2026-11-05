## T-01 — Report a cycle file named without an ID

Result: `grove check` reports a cycle file whose name starts with no cycle
ID and exits non-zero, and `grove changes` names no cycle for its
members, as "Report a cycle file named without an ID" in
[specs/cycles/spec.md](specs/cycles/spec.md) says.

Files:
- `internal/doc/doc.go`
- `internal/cli/changes.go`
- `internal/cli/walkthrough_test.go`
- `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: `grove check` reports "the file name must start with a cycle
ID, such as 261015-k3m9q" for the file.

Acceptance:
- A test with a cycle file named `261105-proto-unknown-id.yaml` fails at
  8a8d064 and passes after the fix: `grove check` reports the file and
  exits non-zero, and `grove changes` shows its member with no
  "in cycle".
- A correctly named cycle is reported and shown as before.
- A walkthrough step renames a cycle file so, runs `grove check` and
  `grove changes`, and discards the rename.
- `just check` passes.

## T-02 — Say a missing commit is missing from a shallow clone

Result: in a shallow clone, `grove check` says a change's `base` or a
cycle's `planned` commit is not in this shallow clone and how to fetch
its full history, as decision 3 of the proposal gives it; in a full
clone the message is unchanged.

Files:
- `internal/gitx/gitx.go`
- `internal/cli/check.go`
- `internal/cli/walkthrough_test.go`
- `docs/guides/run-grove-in-ci.md`
- `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: `internal/gitx` reports whether a repository is shallow.

Acceptance:
- A test in a `git clone --depth 1` of a repository whose change's `base`
  and cycle's `planned` are older than the clone's one commit fails at
  8a8d064 and passes after the fix: each message says the commit is not
  in this shallow clone and to fetch its full history.
- A `base` that is no commit in a full clone is reported as before.
- `docs/guides/run-grove-in-ci.md` "What to expect" shows the message and
  says the checkout needs `fetch-depth: 0`.
- A walkthrough step clones the fixture with `--depth 1` and runs
  `grove check`.
- `just check` passes.
