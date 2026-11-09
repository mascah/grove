# Implementation: Show what is undecided and blocked

<!-- Progress on change 261108-cns75, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Mark open questions and blocked tasks

Status: done

- `internal/doc/doc.go`: `Change.OpenQuestions` (JSON `open_questions`,
  left out when empty), read in `loadChange` from the main document by
  `openQuestions`: HTML comments removed, the text under `## Open
  questions` up to the next `#` or `##` heading, "" when that is empty
  or "None" with an optional period in any case.
- `internal/cli/changes.go`: after the `waits on` / `in cycle` text, two
  spaces and `open questions`, then two spaces and the IDs of the tasks
  the target's `implementation.md` gives `Status: blocked`, with
  `blocked`.
- `internal/cli/serve.go`: `branchChange.OpenQuestions`, read at the
  branch tip or in the working tree. The page's change rows embed
  `doc.Change`, so they carry the key as `grove changes --json` does. No
  existing key changed shape.
- `internal/cli/page.html`: in a member row's status column, `open
  questions` and `T-02, T-04 blocked` chips; on a branch card, `open
  questions` after the task chips (shown even with no tasks).
- Delegated: a direct test of the "None" rule, `internal/doc/doc_test.go`
  `TestOpenQuestions` (the template's comment with "None.", a comment
  only, lowercase "none", no section, a section with a sub-heading,
  "None, but…").
- `internal/cli/walkthrough_test.go` `TestServe`: change A (in cycle K)
  has an open question and T-02, T-04 blocked on main; B (waits on A, in
  K) has an open question and T-02 blocked only on `build-b`; D's
  proposal has the template's comment and "None."; draft G on
  `draft-notes` has a question. It checks the `grove changes` lines, each
  change's `open_questions` in `/api/project` (A's and B's text, none for
  D and the rest) and each branch row's `open_questions`. Mutating the
  `blocked` marker in `changes.go` fails it.
- On a fixture (`just fixture`, `grove init`; change C with a question
  and T-02, T-04 blocked, B with "None." and T-02 blocked only on
  `build-b`, draft D with a question on `draft-notes`, cycle K of C and
  B): `grove changes` printed `261108-c0000  Notes on records  approved
  in cycle 261108-k0000  open questions  T-02, T-04 blocked` and B's line
  with no marker; `grove changes --json` gave C and the draft D
  `open_questions` and B no key. The page, rendered in headless Chrome
  (`--dump-dom` and a screenshot), showed `ready`, `open questions` and
  `T-02, T-04 blocked` in C's row under K, nothing more in B's row,
  `T-02 blocked` on `build-b`'s card, and `open questions` beside D's
  task chips on `draft-notes`'s card.
- The guide "Run Grove in CI"'s `jq '.changes[] | [.id, .state, .path,
  .title]'` reads only keys that did not change.
- `just check` passed: 4371 Go, 2023 test, 410 page lines, 8613
  instruction words.

### T-02 — Documentation and the walkthrough

Status: not started

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
