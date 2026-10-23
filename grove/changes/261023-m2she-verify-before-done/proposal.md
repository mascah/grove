# Verify before done

## Problem

A change can be finished without anyone showing it does what it agreed.
At 57af2e3, implement and apply record each task's status and evidence
(`internal/instruct/requires/implement.md`, `requires/apply.md`), and
archive moves a change once every task says `Status: done`. Nothing ties
a scenario in the change specs to a test, so a task list that missed a
requirement still finishes, and the maintained spec then says agreed what
was never shown. Implement also trusts a task already marked done
(`requires/implement.md`: "A task it marks done when you start is done"),
though later work may have broken it.

Before review, nothing checks that a proposal agrees with itself. Propose
writes the proposal, change specs and tasks in one pass
(`internal/instruct/schemas/default/skills/propose.md`), and the owner is
the first to look for a requirement no task covers or a scope that
contradicts a task.

In the 2026-10-23 review, OpenSpec's `verify` and Spec Kit's `analyze`
and `converge` were the ideas most worth taking; both reviews (Claude's
and Codex's) ranked this first.

## Intended outcome

Before presenting a change, propose checks that each requirement has an
observable outcome, each ADDED or MODIFIED requirement is covered by a
task's acceptance, the scope, specs and tasks do not contradict, and no
task is out of scope. It fixes what it can, tells the owner what it
changed, and asks about what needs a decision.

Before finishing, implement and apply run the project's tests, re-read
each done task's evidence against the code as it is now, and write a
`## Scenarios` section in `implementation.md`: one line per scenario in
the change specs, naming the check that shows it and its result, or
`unverified` and why. A done task whose evidence no longer holds is in
progress again. The finishing report, and archive, name each unverified
scenario. Grove reads none of it (rule 3).

The rules are in [specs/propose/spec.md](specs/propose/spec.md),
[specs/implement/spec.md](specs/implement/spec.md) and
[specs/apply/spec.md](specs/apply/spec.md).

## Scope

### In scope

- Propose's check before presenting: requires text and the default
  schema's propose guide.
- The scenario list and the recheck of done tasks: implement's and apply's
  requires text and guides, the implementation template and guide, and
  archive's requires text.
- `docs/concepts.md`, `docs/getting-started.md`, `docs/reference.md` and
  `dev/walkthrough.md`.

### Out of scope

- Grove reading the scenario list, or archive refusing on it.
- Reading test reports (JUnit and the like): parked in the vault's
  `test-evidence.md`.
- A new verify command or document kind.
- Checking the fold: change 261023-bfwp5.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to both").

- **A requirement in each of `propose`, `implement` and `apply`.** Apply's
  moves into implement when cycle B folds the two skills together.
- **Propose's check covers four things:** observable outcomes, every
  ADDED or MODIFIED requirement covered by a task's acceptance, no
  contradiction between scope, specs and tasks, and no task out of scope.
  The agent fixes what it can, says what it changed, and asks about what
  needs a decision. Reason: no new command or document.
- **`implementation.md` gets a `## Scenarios` section,** one line per
  scenario: the check and its result, or `unverified` and why. Grove does
  not read it (rule 3).
- **Before finishing, run the project's tests and re-read every done
  task's evidence;** a task whose evidence no longer holds is in progress
  again. With the owner or without.
- **Archive names unverified scenarios to the owner before the fold, and
  does not refuse.**

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md),
  [specs/implement/spec.md](specs/implement/spec.md),
  [specs/apply/spec.md](specs/apply/spec.md)
- Tasks: [tasks.md](tasks.md)
