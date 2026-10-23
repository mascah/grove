# Check requirements and the archive fold

## Problem

Nothing checks the requirements themselves. At 57af2e3, `grove check`
reports a document's metadata, headings, dependencies, cycles and stale
plans (`internal/cli/check.go`), but not a requirement with no "shall" or
no scenario, though the spec and change-spec guides ask for both
(`internal/instruct/schemas/default/documents/spec/guide.md`).

Nothing checks the fold either. `grove archive` refuses a change only
when a task is not `Status: done` (`internal/cli/instructions.go`,
`archive`). Folding the change's specs into the maintained specs is the
agent's, by the archive guide
(`internal/instruct/schemas/default/skills/archive.md`), and nothing
compares the result with the change. A requirement reworded, left out or
left in after a REMOVED, or an unrelated requirement edited on the way, is
merged as agreed behavior. In the 2026-10-23 review against OpenSpec,
Spec Kit and AI-DLC, drift between maintained specs and what was agreed
was the most common complaint about such tools, and OpenSpec's own
CHANGELOG lists several silent fold failures.

## Intended outcome

`grove check` reports a requirement that has no "shall" or no scenario,
in a maintained spec or an open change, and a change heading left in a
maintained spec. `grove archive` refuses a finished change when a
maintained spec it touches does not hold what it should: the requirements
where this branch left the target, with what each change archived on this
branch since then gives, and then what this change gives. Each problem is
printed as `path:line: message`, and nothing moves until the fold is
fixed. The archive instructions tell the agent to copy requirements word
for word. The rules are in
[specs/requirements/spec.md](specs/requirements/spec.md).

## Scope

### In scope

- The requirement shape check and the change-heading check in `grove
  check`.
- The fold check in `grove archive`, with `archive` in a file of its own.
- The archive instructions and guide; `docs/reference.md`,
  `docs/concepts.md` and `dev/walkthrough.md`.

### Out of scope

- Whether the code meets the scenarios: change 2 of cycle A, verify
  before done.
- Grove folding the specs itself: the agent still does the fold.
- Sections of a maintained spec other than its requirements, such as
  Purpose and Known discrepancies.
- Archived changes, which `grove check` does not read.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to both").

- **A new capability, `requirements`,** for how a requirement is written
  and how archive folds it. Reason: Grove has no spec for archive, and
  these rules cover maintained and change specs alike.
- **Every requirement not under REMOVED has "shall" and a `#### Scenario:`,**
  in maintained specs and open changes; archived changes are not checked.
  Reason: the guides already ask for both. It is a new `check` problem on
  documents that were clean, so the change is breaking (`CONTRIBUTING.md`).
- **`check` reports an ADDED, MODIFIED or REMOVED heading in a maintained
  spec.** Reason: the archive guide says they never appear there.
- **The fold is compared with where this branch left the target, plus each
  change archived on this branch since then, plus this change.** Not the
  target now, which may have moved, and not the change's `base`, which may
  be older. Reason: it is what this branch's specs should hold whatever
  happened on the target meanwhile.
- **Whitespace is ignored, wording is not; sections other than
  requirements are not compared.** Reason: reflowed lines change nothing
  agreed.
- **Archive prints each problem as `path:line: message` and moves nothing;
  a dropped change is not compared,** since it skips the fold.
- **The archive instructions say to copy each requirement word for word,**
  and after a refusal to fix the spec and run `grove archive` again.
- **`archive` moves into `internal/cli/archive.go`.** Reason: one file per
  command (`CLAUDE.md`).

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/requirements/spec.md](specs/requirements/spec.md)
- Tasks: [tasks.md](tasks.md)
