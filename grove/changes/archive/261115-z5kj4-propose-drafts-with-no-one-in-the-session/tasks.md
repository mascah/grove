## T-01 — Propose instructions for a session with no one in it

Result: `grove instructions propose --unattended` prints instructions
that write one draft from the request, asking no one, as the propose and
schemas change specs' requirements say, and without the flag prints what
it printed before.

Files:
- `internal/instruct/requires/propose-unattended.md`
- `internal/instruct/requires/propose.md`
- `internal/instruct/instruct.go`
- `internal/instruct/skills.yaml`
- `internal/cli/cli.go`
- `internal/cli/instructions.go`
- `internal/cli/walkthrough_test.go`

Interfaces: `grove instructions propose --unattended [--schema <name>]
"<request>"`, the flag also read as the first word of the quoted
description; `instruct.Propose` takes whether no one is in the session.

Constraints: the mode's text is Grove's own, under `requires/`, so that a
project's propose guide cannot replace it (rule 10).

Acceptance:
- `just check` passes, with at most 7,954 instruction words.
- `grove instructions propose --unattended "x"` and `grove instructions
  propose "--unattended x"` print the same instructions, and the request
  they print is `x`.
- Those instructions say to ask and wait for no one; to write the round,
  and whatever the check of the change would ask, as open questions; to
  record as decisions only what the request says; to write nothing for a
  request that names no outcome; to write one change, naming a split, a
  change to land first or `bugfix` as open questions; and to offer no
  `/grove-implement` on the branch.
- `grove instructions propose --unattended --schema bugfix "x"` prints
  them following `bugfix`.
- `grove instructions propose "x"` prints what it printed at efdb8af.
- `grove instructions propose --unattended` with a change's ID, or with
  no request, prints one sentence and exits non-zero.
- The walkthrough test checks each of these.

## T-02 — Docs, reference and walkthrough show the mode

Result: the user documentation says how to have a session nobody attends
draft a change from a request and what the owner then reads, and the
walkthrough runs it on the fixture.

Files:
- `docs/running.md`
- `docs/reference.md`
- `dev/walkthrough.md`
- `grove/specs/propose/spec.md`

Acceptance:
- `docs/running.md` has a section with the line that starts such a
  session, the `rules.skills.propose` lines that commit and push the
  draft, and what a request with no outcome leaves.
- `docs/reference.md` lists the flag beside `grove instructions propose`.
- A walkthrough step under "Running unattended" drafts "Editing an
  unknown ID prints an error and exits 2" with no one in the session and
  says what the proposal's Decisions and open questions hold.
- The Purpose of `grove/specs/propose/spec.md` names the mode.

## T-03 — A skill trial for a draft nobody attends

Result: a sixth trial case counts how each harness drafts a change from
a request with no one to answer it.

Files:
- `dev/trials.md`
- `dev/trials/f.sh`
- `dev/trials/run.sh`

Acceptance:
- Case (f) runs `/grove-propose --unattended` with a request that settles
  the exit status and leaves where the check goes open. A good run ends
  its first turn without waiting, on a branch other than `main`, with one
  change that `grove check` passes, whose Decisions quote only the
  request and whose open questions each give a recommended answer and
  name a task.
- `just trials` runs every case and its table goes in the pull request.
