# A built-in bugfix schema

## Problem

Fixing a bug is the most common small change, and Grove makes every
project build the schema for it. At 57af2e3 Grove's only built-in change
schema is `default` (`internal/instruct/instruct.go`, `SchemaNames`), the
software change, whose proposal asks for a problem, an outcome, a scope and
decisions (`internal/instruct/schemas/default/documents/proposal/template.md`).
`docs/reference.md`, "Example: a bug schema", gives a `bugfix` schema as
files a project copies into `grove/schemas/bugfix/`, and
`dev/walkthrough.md`, "A bug", has the tester write them before proposing a
fix. A project that has not copied them proposes a bug fix as a software
change, with no section for what was observed, how to reproduce it or the
requirement it breaks, and nothing asking for a failing test first.

Explore and propose list the project's schemas only when it has more than
`default` (`internal/instruct/instruct.go`, `choices`), and switch schema
only when the owner says which kind of work it is
(`internal/instruct/requires/propose.md`,
`internal/instruct/requires/explore.md`). In the 2026-10-23 review against
OpenSpec, Spec Kit and AI-DLC, the owner agreed Grove should ask for less
ceremony on small changes; this is change 3 of cycle B.

## Intended outcome

Every project has a `bugfix` schema beside `default` without writing a
file. Its proposal has Observed, Expected, Reproduction, The requirement it
breaks and Open questions; its first task is a test that fails before the
fix; it writes a change spec only when the owner agrees behavior that no
maintained requirement gives. A project still replaces any of its files, in
`grove/schemas/bugfix/`, and a file of its own for `default` still reaches
`bugfix` unless `bugfix` has its own. When the owner describes behavior
that is wrong, explore and propose suggest `bugfix` and use it once the
owner confirms. The rules are in
[specs/schemas/spec.md](specs/schemas/spec.md).

## Scope

### In scope

- The built-in `bugfix` schema, embedded beside `default`, from the
  reference's example: its `schema.yaml`, propose guide, and proposal
  template and guide.
- Where each of a schema's files comes from, with the built-in `bugfix`
  between the project's `bugfix` and its `default`.
- The explore and propose instructions suggesting `bugfix`.
- `docs/reference.md`, `docs/concepts.md`, `dev/walkthrough.md` and the
  layout table in `CLAUDE.md`.

### Out of scope

- Other built-in schemas, such as one for docs or chores.
- Grove deciding a change's schema itself: the owner confirms it.
- Running the failing test: building and verifying the change is
  implement's, and verifying before done is change 261023-m2she.
- Changing `default`, or what Grove requires of any skill in every schema
  (`CONSTITUTION.md` rule 10).

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes").

- **Grove ships `bugfix` beside `default`, made from the example in
  `docs/reference.md` ("Example: a bug schema"), trimmed.** Its proposal's
  sections are Observed, Expected, Reproduction, The requirement it breaks
  and Open questions; the first task is a test that fails before the fix;
  a change spec is written only when the owner agrees new agreed behavior,
  and restoring an existing requirement writes none. Reason: the example is
  already trialed in the walkthrough, and a bug fix that restores what is
  agreed changes nothing agreed. Rule 10, as amended at 57af2e3, has a
  project override "a built-in" schema, which leaves room for this one.
- **Each file of `bugfix` comes from the project's `grove/schemas/bugfix/`,
  then Grove's built-in `bugfix`, then the project's
  `grove/schemas/default/`, then Grove's `default`.** Reason: a project's
  own `default` file, such as its tasks template, still reaches `bugfix`
  unless `bugfix` has its own, as it does for a project's schemas today.
- **Explore and propose suggest `bugfix` when the owner describes wrong
  behavior; the owner confirms.** Propose already lists the project's
  schemas when there are others than `default`; now there always are.
  Reason: the owner should not have to know the schema's name to get it,
  and choosing it stays the owner's.
- **`docs/reference.md`'s "Example: a bug schema" becomes a description of
  the built-in, and "A schema of your own" keeps a short example. It is not
  breaking:** a project's existing `grove/schemas/bugfix/` still overrides
  the built-in file by file. Reason: the docs are the specification
  (rule 8), and a project that copied the example keeps what it has.

## Open questions

None must be answered before review. None blocks a task.

- Can wait: a project with files in `grove/schemas/bugfix/` but no
  `schema.yaml` there, which Grove ignored at 57af2e3, would now have them
  read as `bugfix` files, and `grove check` would report one `bugfix` does
  not use. Whether that edge makes the change breaking under
  `CONTRIBUTING.md`, "Breaking changes", is the owner's to say at review.

## Related documents

- Requirements: [specs/schemas/spec.md](specs/schemas/spec.md)
- Tasks: [tasks.md](tasks.md)
