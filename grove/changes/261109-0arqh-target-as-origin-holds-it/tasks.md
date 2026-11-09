## T-01 — The target's commit is the merge base with origin's

Result: every command reads the target at the newest commit the local
branch and `refs/remotes/origin/<target>` both hold, or at the one that
exists, fetching nothing; `run`'s fast-forward refusal still compares the
local tip with the fetched origin; a checkout with no local target branch
works.

Files:
- `internal/gitx/gitx.go`
- `internal/cli/run.go`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers; the name of the local-tip read

Interfaces: `gitx.ResolveBranch(dir, name)` returns the merge base of
`refs/heads/<name>` and `refs/remotes/origin/<name>`, or the one that
exists; a read of the local tip for `run`.

Constraints: Git runs only in `internal/gitx` (`CLAUDE.md`); nothing is
fetched (rule 7); only the remote named `origin`; the Go and test line
budgets of rule 12 hold; the pull request is `feat(cli)!` with a
`BREAKING CHANGE:` footer telling owners to push local merges for Grove
to count them; the constitution amendment of decision 4 is merged first.

Acceptance:
- Every scenario of the config change spec holds on the fixture.
- A test in `internal/cli/walkthrough_test.go` builds a bare origin with
  two clones and a checkout with no local `main`, and shows `base` as the
  shared commit, the second clone's clean `check`, and `grove changes`
  and `grove check` working with no local `main`.
- `grove run`'s refusal when the local target has commits origin's does
  not still holds.
- `just check` passes.

## T-02 — Documentation and the walkthrough

Result: the documentation says the target is read as this checkout and
origin both hold it, that `base` and `planned` are that commit, that a
local merge counts once pushed, and that a pull-request job's checkout
needs no local `main`; the walkthrough runs the two-clone and
no-local-main steps on the fixture.

Files:
- `docs/concepts.md`
- `docs/coordinating.md`
- `docs/planning.md`
- `docs/reference.md`
- `docs/guides/run-grove-in-ci.md`
- `dev/walkthrough.md`

Acceptance:
- Where concepts.md, coordinating.md, planning.md and reference.md say
  "your main branch's commit" for `base`, `planned` or `target`, they say
  the commit your main branch and origin's both hold, and concepts.md
  says a local merge is approved once pushed.
- reference.md's `target` row says how the target is read.
- The CI guide says a pull-request job's checkout, with `origin/main` and
  no local `main`, works.
- Each documented step is run on the fixture and recorded in
  `dev/walkthrough.md`.
