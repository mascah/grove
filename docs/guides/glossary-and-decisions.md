# A glossary and decision records with `rules`

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

Words drift as much as behavior does. One person's "old record" is
another's "existing record", and an agent will happily invent a third name.
This recipe keeps a glossary of your project's own terms, agreed the way
specs are. A change carries the terms it adds or sharpens. Your merge
approves them, and archiving folds them into one glossary. The recipe also
keeps a short record of the few decisions a later reader would otherwise
trip over.

It needs nothing beyond `rules` lines in `grove.yaml` (see
[Reference](../reference.md), "grove.yaml").

## What to copy

Put these in `grove.yaml` in place of the `rules: {}` that `grove init`
writes, or merge them into the `rules` you have. They assume the default
root, `grove`; change the paths if yours differs.

```yaml
rules:
  skills:
    explore:
      - "Read grove/glossary.md if it exists. When I use a term that conflicts with it, or a vague word for something it names, challenge it and suggest the precise term. Ask about edge cases at a term's boundary."
    propose:
      - "If this change introduces or sharpens a domain term, write glossary.md in the change folder in the glossary format below, holding only the terms this change adds or changes."
      - "Glossary format: '# Glossary', one or two sentences naming the domain, then '## Terms': each term in bold, a definition of one or two sentences, and an '_Avoid_:' line of rejected synonyms. Only words specific to this project; no implementation details."
      - "In the proposal's decisions, mark a decision '(ADR)' only if all three hold: it is hard to reverse, it would surprise a reader without context, and it came from a real trade-off. Most decisions are not ADRs."
    implement:
      - "If building a task shows that a term in grove/glossary.md or the change's glossary.md is wrong or missing, say so under '## Glossary' in implementation.md."
    archive:
      - "Fold the change folder's glossary.md, if any, into grove/glossary.md, creating it if needed; a changed term replaces the old one. Leave the change's own copy in place."
      - "For each decision marked '(ADR)' in the proposal, write grove/decisions/<change-id>-<few-words>.md - a title, one paragraph saying what was decided and why, the options set aside if the proposal names them, and the change ID."
```

Run `grove check` to see that `grove.yaml` still reads cleanly, and commit.

## What to expect

- **Propose** writes `glossary.md` in the change folder with the terms the
  conversation settled, and only those.
- **Archive** folds them into `grove/glossary.md`, creating it on the first
  change, and writes one file in `grove/decisions/` for each decision marked
  `(ADR)`.
- **Decision records are rare.** In the trial, propose marked no decision
  `(ADR)` on its own. If a decision deserves a record, say so while
  proposing.
- **Later changes use the agreed terms.** In the trial, the second change
  asked what its new term meant, and did not ask again about a term the
  glossary already defined.

Grove itself knows nothing of these files. `grove check` does not read
them. `grove review` lists `grove/glossary.md` and `grove/decisions/` as
`maintained spec`: that is its label for any file under the root outside
`changes/`, `schemas/` and `cycles/`.
