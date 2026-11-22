# grove claim and the project value read the target's schemas

<!-- A bug fix describes future work. It becomes approved future work when
it is merged into the target branch, which says nothing about
implementation. -->

## Observed

`grove claim <change-id>` records an attempt for a change that shares a
capability with a change in flight, when both use a project schema the
claiming checkout's working tree does not hold. Reproduced on the ledger
fixture with Grove at a4ad6e2 (see "Reproduction"): `grove claim` of the
second change printed its ID and exited 0, and `origin`'s `grove/runs` then
held two records, where a refusal naming
`grove/specs/record-editing/spec.md` was due. Codex found it first,
reviewing pull request #185 after it merged.

The cause, at a4ad6e2: `project.instruct` (`internal/cli/project.go`)
takes the schemas from the working tree,
`os.DirFS(filepath.Join(p.Top, p.Root, "schemas"))`, and `project.at`
loads the target commit's changes with that layout. A change whose schema
the working tree lacks gets no capabilities (`Change.Touches`, in
`internal/doc/doc.go`), so `Change.Planned` leaves out its maintained
specs, and `busy` (`internal/cli/claim.go`) finds no overlap. The project
value is built the same way (`server.value` in `internal/cli/serve.go`,
`changes` in `internal/cli/changes.go` through `project.agreed`), so its
verdicts miss the same overlap. A runner's checkout is usually on its last
attempt's branch, whose working tree can predate a schema the target has
since gained.

## Expected

`grove claim B` refuses B, saying it shares
`grove/specs/record-editing/spec.md` with A, naming A's attempt; and
before that claim, B's verdict in `grove changes --json` says the same.

## Reproduction

1. `just fixture ~/tmp/ledger-schemas`, `grove init`, commit. Create a
   bare repository with `git init --bare -b main`, add it as `origin`,
   push `main`, and clone it twice, as `h1` and `h2`.
2. In the owner's checkout, copy `internal/instruct/schemas/default` to
   `grove/schemas/lite`. Add changes A and B with `schema: lite`, each
   with a change spec for `record-editing`, A's one task naming
   `ledger/cli.py` and B's `ledger/report.py`. `grove check` reports no
   problems. Commit and push.
3. In `h1`, `grove claim <A>` prints A's ID. In `h2`, `grove claim <B>`
   prints B's ID and exits 0, and `git -C <origin> log --oneline
   grove/runs` lists two records.

## The requirement it breaks

`grove/specs/runs/spec.md`, "Refuse what is in flight or overlaps it",
its scenario "One capability, different requirements", with the
Purpose's rule that planned files and overlap are "read from the fetched
target's documents and from `grove/runs`, never from an attempt's
branch". `grove/specs/serve/spec.md`, "Readiness of each change in no
cycle" and "Readiness of each cycle", which read the same "from the
target the checkout holds". What is agreed does not change, so this
change has no change spec.

## Decisions

The owner answered "yes" on 2026-11-22 to each recommended answer:

1. `grove claim` and the project value (`grove changes --json` and
   `/api/project`) read the schemas from the target commit whose changes
   they read, so the page's verdicts agree with what `claim` does.
2. `grove instructions`, `check`, `archive` and `review` keep reading the
   working tree's schemas, since a schema drafted on a branch must work
   for them; they are out of scope.
3. It needs no other change to land first.
4. Routine: one task, whose first acceptance is a test that fails before
   the fix; private helpers and test names are the implementer's; no
   user-facing message changes.

## Open questions

None.

## Related documents

- Tasks: [tasks.md](tasks.md)
- Requirements: [runs](../../specs/runs/spec.md),
  [serve](../../specs/serve/spec.md); no change spec
