# Propose asks what matters

## Problem

Propose's round gives every question the same weight. At 57af2e3,
`internal/instruct/requires/propose.md` has the agent ask "every decision
you would otherwise assume, numbered", and the propose guide
(`internal/instruct/schemas/default/skills/propose.md`, step 3) lists
"which option, which requirements and scenarios, where each task goes in
the code, what is out of scope, which files outside the code change". A
question that changes what users see sits in the same list as the name of
a test file, in no particular order. A long flat round is easy to answer
"yes" without reading. In the 2026-10-23 review, Spec Kit caps its
clarification questions, and AI-DLC's lesson was that "tired people still
approve anything".

The round settles too much for implement, which then asks again. At
57af2e3, `internal/instruct/requires/implement.md` has the agent change
only the files the tasks name and ask the owner before touching any other,
even a test helper or a private module the owner does not care about.
There is no way for the owner to say once, while proposing, which small
choices the implementer may make alone.

Implement has no explicit stop when the approach turns out wrong.
`requires/implement.md` asks the owner when "a task cannot be built as
the change says", and revises the change only with the owner's agreement.
A task that can still be built as written, but that the agent finds rests
on something false, is not named: nothing says to stop building it.

Propose never offers to split. The propose guide's step 2 confirms the
scope of one change. An idea with parts the owner could agree separately
becomes one large change, reviewed and merged as one.

## Intended outcome

Propose's round asks the questions that change behavior, scope,
compatibility or lose data first, then those that fix architecture or
interfaces costly to reverse, and puts routine choices inside the agreed
approach last, as recommended defaults on one line that one "yes"
accepts. Nothing is assumed silently: every decision is still asked.

While confirming scope, propose says when an idea is too big for one
change or has parts the owner could agree separately. If the owner
agrees, it writes several changes on the same branch, each from its own
`grove instructions propose`, linked with `depends_on`.

A task may carry a `Delegated:` line, written only when the owner agrees
in the round, naming choices the implementer may make without asking and
paths it may change beyond `Files:`. Implement with the owner honors it.
`/grove-implement --unattended`, which builds without the owner, ignores
it and keeps to `Files:`. Grove reads none of it.

While implementing, a discovery that makes the agreed approach wrong stops
that task, and the agent tells the owner what it found and the revision
the change needs.

The rules are in [specs/propose/spec.md](specs/propose/spec.md) and
[specs/implement/spec.md](specs/implement/spec.md).

This change builds on change 261023-g0a3h, "Implement builds unattended,
and apply goes", and depends on it: that change folds apply into
`/grove-implement --unattended`, splits implement's requires text into a
shared part, the owner's part and the unattended part, and modifies
"Ask a round before writing", which this change modifies as it gives it.

## Scope

### In scope

- The order and batching of propose's round, and offering to split, in
  propose's requires text and the default schema's propose guide.
- The `Delegated:` line in the default schema's tasks guide and template.
- Implement honoring a task's delegation and stopping on a discovery, in
  the owner's part of its requires text, and ignoring the delegation, in
  the unattended part, as change 261023-g0a3h names those parts.
- `docs/concepts.md`, `docs/getting-started.md`, `docs/reference.md`,
  `dev/walkthrough.md` and `internal/cli/walkthrough_test.go`.

### Out of scope

- Grove reading the `Delegated:` line (`CONSTITUTION.md` rule 3 lists
  what Grove reads). So `grove check` reports a moved file only for paths
  a `Files:` line names, never for a delegated path.
- A cap on the number of questions in the round: every decision is still
  asked; the round is ordered and batched, not shortened.
- A stop on a discovery in an unattended build: change 261023-g0a3h's
  "Build unattended" already marks blocked, with the reason, a task that
  cannot be built as the change says.
- Folding apply into `/grove-implement --unattended`: change 261023-g0a3h.
- Propose checking its documents before presenting them: change
  261023-m2she, which also edits propose's requires text and guide. The
  two add separate requirements and neither needs the other first.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes").

- **The round is ordered by consequence.** Questions that change
  behavior, scope, compatibility or lose data come first; then
  architecture or interfaces costly to reverse; routine choices inside an
  agreed approach come last, as recommended defaults on one line that a
  single "yes" accepts. Nothing is assumed silently: the existing "every
  decision you would otherwise assume" holds, and this orders and batches
  it. Reason: Spec Kit caps clarification questions, and AI-DLC's lesson
  is "tired people still approve anything"; the owner's attention goes to
  what is costly to get wrong.
- **An optional `Delegated:` line on a task.** It names choices the
  implementer may make without asking, such as private helpers or test
  layout, and paths it may change beyond `Files:`. Propose writes it only
  when the owner agrees in the round. Implement with the owner honors it;
  an unattended build ignores it and keeps to `Files:`. The owner's brief
  named today's apply and said change 261023-g0a3h carries it into
  implement's unattended mode; building on that change, this one writes
  it there directly.
  Reason: the owner settles once, while proposing, what is not worth a
  question later, and nobody widens what an unattended build may touch.
- **A discovery that invalidates the agreed approach stops the task.**
  The agent explains the revision the change needs. Today's implement
  already asks before revising; this makes the stop explicit. Reason: a
  task built on a premise found false is wasted work, and the revision is
  the owner's to agree.
- **Propose offers to split.** When one idea is too big for one change,
  or has parts the owner could agree separately, propose says so while
  confirming scope; if the owner agrees, it writes several changes on the
  same branch, each from its own `grove instructions propose`, linked
  with `depends_on`. Reason: smaller changes are reviewed and approved on
  their own, and `depends_on` already orders them.
- **Grove parses nothing new.** `Delegated:` is read by people and
  agents, not by Grove. Reason: rule 3 lists what Grove reads.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md),
  [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
- Builds on: change 261023-g0a3h
