## T-01 — List the change's documents by path

Result: the implement, archive and propose revise instructions behave as
[specs/implement/spec.md](specs/implement/spec.md) "List the change's
documents" and [specs/propose/spec.md](specs/propose/spec.md) "List a
revised change's documents" say.

Files: `internal/instruct/instruct.go`, `internal/cli/instructions.go`,
`internal/instruct/requires/implement.md`,
`internal/instruct/requires/archive.md`, `docs/reference.md`,
`dev/walkthrough.md`, `internal/cli/walkthrough_test.go`

Interfaces: `Instructions.Change` holds each document's path and, for an
approved change implement builds, whether its working-tree copy differs
from the target's, and no content; `Markdown` prints "The change" as a
list of paths, a differing one with `git show <commit>:<path>`.

Acceptance:
- On the fixture, `grove instructions implement`, `archive` and `propose
  <change-id>` print each document's path and none of its text; their
  length does not grow with the documents' text.
- After an edit of `tasks.md` on the branch, `grove instructions implement`
  names it as differing, with the `git show` command; without one, none is
  named.
- The cycle file is still printed in full for a member of a cycle.
- `dev/walkthrough.md`'s step saying the agent "is given that change's
  files in full" and `docs/reference.md` say the documents are listed
  (rule 8).
- `walkthrough_test.go` checks the listing, the differing document, and
  that no document's text is printed.
- `just check` passes, the instruction text included.

## T-02 — Implement prints only the implementation template

Result: implement behaves as [specs/implement/spec.md](specs/implement/spec.md)
"Implement leaves revising to propose" says.

Files: `internal/instruct/schemas/default/schema.yaml`,
`internal/instruct/requires/implement-owner.md`,
`docs/getting-started.md`, `dev/walkthrough.md`,
`internal/cli/walkthrough_test.go`

Interfaces: the default schema's implement `writes` is
`[implementation]`.

Acceptance:
- On the fixture, `grove instructions implement <change-id>` prints the
  implementation template and no proposal, change spec, tasks or
  `.grove.yaml` template; `bugfix`, which takes default's skills, does
  the same.
- With the owner, an agreed revision is made after running `grove
  instructions propose <change-id>`.
- The docs and the walkthrough say so (rule 8).
- `walkthrough_test.go` checks the templates printed.
- `just check` passes, the instruction text included.

## T-03 — Trials count listed documents not read

Result: the trials report, per run, each document the instructions listed
that the session never read.

Files: `dev/trials/run.sh`, `dev/trials.md`

Acceptance:
- `dev/trials.md` defines "Listed documents not read" as a judged count,
  and the results table has its column.
- The judge prompt in `dev/trials/run.sh` asks for it, and a run of one
  case reports it.
