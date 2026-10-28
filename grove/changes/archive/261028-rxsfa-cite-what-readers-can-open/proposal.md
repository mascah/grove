# Cite what the project's readers can open

## Problem

Grove documents are read by everyone who reads the project: teammates with a
fresh clone, agents in a fresh checkout and, for a public project, anyone.
Some of what they cite exists only on the author's machine.

- Grove's own instructions cause one case. At 451bb6f,
  `internal/instruct/requires/propose.md` (the prototype line) and the
  `prototype` spec's "Propose uses a prototype's findings" have a proposal
  cite the prototype branch and its commit. The prototype instructions
  (`internal/instruct/requires/prototype.md`) say to push nothing, so that
  citation resolves only where the prototype ran.
- On a private project using Grove, three archived proposals cite a
  prototype branch at a commit that exists only on the author's machine:
  the project's remote has none of its `prototype/*` branches. One of them
  cites a prototype that was never committed, which resolves nowhere.
- Agents cite other machine-local things out of habit, with no instruction
  asking them to. For example, Grove's own archived change `261021-2aj5w`
  links a note in another private repository on the author's machine.

Nothing in the opening every skill prints (`meanings`,
`internal/instruct/instruct.go` at 451bb6f) says what a document may cite.

## Intended outcome

- Every skill's instructions tell the agent to write Grove documents for
  whoever reads the project: cite only what those readers can open, and quote
  or summarise what only this machine holds. What readers can open depends on
  the project: an internal company project may cite internal documents, and a
  public one public sources.
- A proposal written after a prototype stands without the prototype branch.
  It states what was tried, what was observed and the verdict, and names the
  branch and its commit only as where the findings came from.

## Scope

### In scope

- One principle in the opening every skill prints, and a new `instructions`
  maintained spec that holds its requirement.
- The prototype line in `internal/instruct/requires/propose.md` and the
  MODIFIED requirement "Propose uses a prototype's findings".
- The same prototype wording where the docs and the walkthrough repeat it:
  `docs/concepts.md`, `docs/getting-started.md` and `dev/walkthrough.md`.

### Out of scope

- Rewriting archived change records. They are history, and the rule applies
  to documents written from now on.
- `internal/instruct/requires/prototype.md`: committing the experiment so a
  proposal can cite its commit still holds, as the note of where the findings
  came from.
- A `grove check` that looks for such citations: no documented step needs it
  (`CONSTITUTION.md` rule 8), and which sources a project's readers can open
  is not something Grove can know.
- New documentation of the principle beyond the prototype lines.
- The change-spec guide's sentence on what a capability is; that is its own
  change.

## Decisions

- **Who the records serve.** Teammates cloning the repository, unattended
  agents in a fresh checkout, and public readers. The owner, choosing all
  three: "I think 1 2 and 3". A record that makes sense only on one machine
  fails all three.
- **The line is the project's audience, not the repository.** The owner: "I
  don't think it's "everything outside the repository", it's what's private
  versus what's public to the audience of the project. An internal company
  repo can reference other internal company docs, and this public repo can
  reference other public things, but not my other private projects on my
  machine that we use to exercise grove, that makes so sense to anyone
  else."
- **One principle in Grove's instructions.** The owner chose to put one
  principle in Grove's instructions: cite only what the project's readers
  can open, and quote or summarise what only this machine holds.
- **Its wording**, accepted by the owner in the round ("yes"): "Write Grove
  documents for whoever reads this project: cite only what they can open,
  and quote or summarise what only this machine holds, such as an unpushed
  branch, a file outside the checkout or a private project."
- **It lives in the opening every skill prints**, not in the guides,
  accepted ("yes"). It reaches every skill in every schema, and a project
  that overrides a guide cannot drop it. It costs Go lines, not instruction
  words.
- **A proposal carries a prototype's findings itself.** The owner chose the
  option in which the branch and commit are only a note of where the
  findings came from. The MODIFIED requirement's wording, with a second
  scenario for an experiment the owner chose not to commit, was accepted
  ("yes").
- **Archived records stay as they are**, accepted ("yes").
- **A new maintained spec, `instructions`, holds the principle's
  requirement**, accepted ("yes"). The opening has no spec at 451bb6f, and
  the rule spans propose, implement and archive, so no one skill's spec
  fits.
- **`requires/prototype.md` stays unchanged**, accepted ("yes").
- **Routine choices**, accepted with one "yes":
  - the title;
  - two tasks;
  - no new documentation of the principle;
  - no test that pins the sentence;
  - no `Delegated:` lines;
  - this branch.

## Open questions

None.

## Related documents

- Requirements: [specs/instructions/spec.md](specs/instructions/spec.md),
  [specs/prototype/spec.md](specs/prototype/spec.md)
- Tasks: [tasks.md](tasks.md)
