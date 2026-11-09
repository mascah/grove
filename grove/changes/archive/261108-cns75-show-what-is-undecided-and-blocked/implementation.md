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

Status: done

- `docs/coordinating.md`, "A change that waits": a `grove changes` line
  ending `open questions  T-02, T-04 blocked`, when the "Open questions"
  section counts, and that the line reports what the main branch holds,
  a task blocked on an unmerged branch showing on that branch's card in
  `grove serve` instead. "The project as data": an `open_questions`
  bullet, and `branches` carrying it as the branch holds it.
- `docs/watching.md`: Readiness's member rows show `open questions` and
  the tasks the target marks blocked; an In progress card shows `open
  questions`, and a task blocked only on a branch shows on its card, not
  in Readiness.
- `docs/reference.md`: the `grove changes` row names both markers.
- `dev/walkthrough.md`, "Watching the project" steps 2 and 3: A gets an
  open question and T-02, T-04 blocked on `main`, B a question, D's
  revision the template's comment and "None.", G on `draft-notes` a
  question; expected Readiness rows, the `build-b` and `draft-notes`
  cards, the `grove changes` lines and the `open_questions` values. These
  are the steps `TestServe` follows (rule 17).
- Run on a fresh fixture (`just fixture`, `grove init`, a bare `origin`,
  this worktree's `bin/grove`), with A, B and D approved as plain commits
  and cycles K (A, B, unattended) and M (D, unattended), `build-b` and
  `draft-notes` as the steps say; the steps of those sections this change
  does not touch (C, E, F, L, N, `teammate`, the claim from a second
  clone) were left out. `grove changes` printed A's line ending `in cycle
  261108-k0000, unattended  open questions  T-02, T-04 blocked`, B's
  `waits on 261108-a0000  in cycle 261108-k0000, unattended  open
  questions`, D's `in cycle 261108-m0000, unattended`. `grove changes
  --json` gave A and B their text, D no key, `draft-notes`'s G "Should
  notes be dated?", and equalled `curl /api/project` but for `read`. The
  page in headless Chrome (`--dump-dom`): A's row `ready  open questions
  T-02, T-04 blocked`, B's `waits on 261108-a0000  open questions` with no
  blocked task, `build-b`'s card T-02 blocked with `open questions`, G's
  card `open questions` beside its tasks, no script run.
- `just check` passed: 4371 Go, 2023 test, 410 page lines, 8613
  instruction words. `grove check`: no problems.

## Scenarios

- changes "Questions and a blocked task": `internal/cli/walkthrough_test.go`
  `TestServe` (A's `grove changes` line ends `open questions  T-02, T-04
  blocked`, its `open_questions` "Should IDs be six digits? Blocks
  T-02."), passed; also on the fixture.
- changes "None": `internal/doc/doc_test.go` `TestOpenQuestions` and
  `TestServe` (D, the template's comment and "None.": no marker, no key),
  passed.
- changes "Blocked only on a branch": `TestServe` (B, T-02 blocked only on
  `build-b`: its line names no blocked task), passed.
- changes "Waiting and in a cycle": `TestServe` (B's line ends `waits on
  <A>  in cycle <K>, unattended  open questions`; K is marked unattended,
  so the line carries `, unattended`, which the scenario's K has not),
  passed.
- serve "A member with questions and a blocked task": `TestServe` checks
  A's `open_questions` and tasks in `/api/project`, passed; the row's
  chips seen in headless Chrome on the fixture (`--dump-dom` and a
  screenshot), as the serve spec's other page scenarios are checked by
  hand.
- serve "A draft with a question": `TestServe` (`draft-notes`'s row gives
  G `open_questions` "Should notes be dated?"), passed; the card's chip
  seen in headless Chrome on the fixture.
- serve "Blocked only on a branch": `TestServe` (`build-b`'s row gives T-02
  blocked; B on the target has no blocked task), passed; in headless
  Chrome, `build-b`'s card shows T-02 blocked and B's row under K none.

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

- changes: applied to specs/changes/spec.md
- serve: applied to specs/serve/spec.md
