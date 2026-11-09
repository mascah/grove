# Show what is undecided and blocked

## Problem

- Every proposal has an "Open questions" section, and the proposal guide
  asks it to say which questions must be answered before review and which
  task each blocks (`internal/instruct/schemas/default/documents/proposal/template.md`
  at ccaee4d). Nothing Grove prints reads it: `grove changes` gives each
  approved change's ID, title, what it waits on and its cycle
  (changes spec at ccaee4d, "List approved changes"), and `grove changes
  --json` gives its tasks but not its questions ("Print the project as
  data").
- `/grove-implement` marks a task `Status: blocked` in `implementation.md`
  (implement spec at ccaee4d, "Settle what the change does not" and "Stop
  at a checkpoint"). `grove changes --json` carries each task's status,
  but the terminal list does not say that a change is stuck.
  `grove serve` shows blocked task chips only on the cards of branches
  holding work the target lacks (serve spec at ccaee4d, "Work in progress
  by branch"; `internal/cli/page.html` at ccaee4d).
- Issue #107 reports that finding what is undecided or stuck across a
  project means opening every `proposal.md` and `implementation.md`, and
  that the reporter keeps a hand-written status file with "Open
  decisions" and "Blocked on" sections for that gap, which drifts. The
  owner narrowed the issue's scope in its comment of 2026-11-08 after
  #139 added `grove serve`.

## Intended outcome

- `grove changes` adds `open questions` to the line of a change the target
  holds whose proposal has any, and `T-02 blocked` (or `T-02, T-04
  blocked`) when the target's `implementation.md` marks tasks blocked.
- `grove changes --json` gives each change `open_questions`: the text of
  its main document's "Open questions" section, left out when the section
  is empty or says "None". Because `grove serve`'s project value gives
  each change as `grove changes --json` does, the page's value carries it
  too; its branch rows carry it as well, read at that branch's tip or in
  the working tree.
- The page shows the same markers in the status column of each cycle's
  member rows, and an `open questions` marker on a branch card's change
  beside the task chips it already shows.
- What the target holds and what a branch holds stay apart (rule 5): a
  task blocked on an unmerged build branch shows on that branch's card,
  not in `grove changes`.

## Scope

### In scope

- Reading the "Open questions" section of each change's main document,
  on the target and on each branch the page lists.
- The markers in `grove changes`, the `open_questions` field in `grove
  changes --json` and in the page's branch rows, and the markers on the
  page.
- The documentation and walkthrough steps that show them (rule 8).

### Out of scope

- A blocked task's reason. It stays in `implementation.md`, which the
  page's pane opens; rule 3 is unchanged.
- `grove changes` scanning branches; the page already lists them.
- The terminal reader that `grove` opens.
- A count of open questions, and any change to a template or document
  convention.
- Tasks with `Status: to check`.
- Markers on the page's dependency graph.

## Decisions

The owner confirmed the scope after an adversarial review of it, and
accepted the round in this conversation on 2026-11-08 with "yes".

1. **The blocked reason stays out.** Accepted as recommended: "Rule 3 is
   unchanged. The page's pane opens `implementation.md` for the detail,
   and a question that blocks a task already reaches the JSON through
   `open_questions`." Reason: the reason has no fixed form today (a
   `Check:` line, or free text "with the reason"), and reading it would
   need an amendment to rule 3 and a new implement rule.
2. **The "None" rule.** Accepted: take the "Open questions" section of the
   change's main document and remove HTML comments; if what is left is
   empty, or is "None" with an optional period in any case, the change
   has no open questions; anything else counts. Reason: it changes no
   document convention, so it is not breaking, and Grove's own proposals
   already write "None." when there are none.
3. **Terminal markers.** Accepted: `open questions` with no count, and
   blocked tasks as `T-02, T-04 blocked`, after the existing `waits on` /
   `in cycle` text, separated by two spaces like today, reflecting only
   what the target holds. Reason: the section is prose, so there is no
   count to trust, and the target is what `grove changes` reports.
4. **JSON field.** Accepted: `open_questions` is a string holding the
   section's text, left out when there are none; the same rule applies to
   every change, archived ones included; branch rows
   (`branches[].changes[]`) carry it too, read at that branch tip or in
   the working tree. Reason: drafts are where most open questions live,
   and a branch row today copies only the fields it names
   (`branchChange`, `internal/cli/serve.go` at ccaee4d).
5. **Page.** Accepted: markers in the status column of cycle member rows,
   plus an `open questions` marker on branch cards, which already show
   blocked tasks; the dependency graph's nodes get nothing. Reason: the
   page's 500-line budget (rule 12) holds 408 lines at ccaee4d.
6. **No `depends_on`.** Accepted: it waits on neither 261108-7grt3 nor
   261108-jpg8n. Reason: the page's change rows already embed the change
   as `grove changes --json` gives it (`pageChange`,
   `internal/cli/serve.go` at ccaee4d), so the field reaches both without
   7grt3.
7. **ADDED requirements only.** Accepted: one in the changes spec and one
   in the serve spec, editing no requirement that 261108-7grt3 or
   261108-jpg8n modifies. Reason: neither change goes stale on a
   requirement when the other lands.
8. **Two tasks**, the code with its tests and the documentation with the
   walkthrough. Accepted as listed in the round.
9. **Routine defaults**, accepted with the round: the section is read in
   `internal/doc` beside the task parsing; test layout, private helpers
   and a small direct test of the "None" rule (rule 17) are delegated to
   the implementer.

## Open questions

None.

## Related documents

- Requirements: [changes](specs/changes/spec.md), [serve](specs/serve/spec.md)
- Tasks: [tasks.md](tasks.md)
