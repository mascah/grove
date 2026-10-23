# Plain docs, exact specs

## Problem

The docs open with the path most readers do not take first. At 0c6f419,
the second paragraph of `README.md` brings in `grove run`, cycles and
unattended runs before the reader has taken one change through. In
`docs/index.md`, the first section of the page says the `grove` command
"can run apply unattended".

The docs are hard to read. Many sentences carry several conditions. One
example from `docs/planning.md`: "A change that waits on one outside the
cycle keeps `grove run` from the cycle until that one is archived as
finished." Terms such as cycle, attempt, base and stale have their
meaning given where they are used, in several places, and the glossary in
`docs/concepts.md` defines few of them. One sentence is now false:
`docs/concepts.md`, "Two halves", says "Grove starts no agent and makes
no commits", and `grove run` starts one.

The same rules are written twice. The amendment to rule 8 (PR #47) makes
the maintained specs in `grove/specs/` the specification and has `docs/`
explain them and link them. At 0c6f419 the docs restate rules the specs
hold: `docs/running.md` gives the conditions for selecting a cycle that
`grove/specs/runs/spec.md`, "Start one attempt", holds. Other exact rules
are in the docs and in no spec. Examples: what `grove check` reports on
`depends_on` and on a stale requirement (`docs/coordinating.md`), the
JSON that `grove changes --json` prints, when `grove archive` refuses, and
how `grove init` chooses harnesses and writes their files
(`docs/reference.md`). The `config` spec's Boundaries says the rest of
`grove.yaml` "is specified by `docs/reference.md`".

## Intended outcome

`README.md` and `docs/index.md` show one change taken through first:
propose, review and merge, implement, merge. Cycles and `grove run` come
after, under "Doing more".

Every page says one idea per sentence. A term is defined once, in the
glossary, and other pages use it.

Each exact rule lives in a maintained spec. A page explains it plainly
and links the spec. The rules found so far that no spec holds become
agreed here: two `grove check` reports in `propose`, and two new
capabilities, `changes` and `init`. The code already does what they say.
Rules found later while auditing are added to this change, with the
owner's agreement, before they leave the docs.

The rules are in [specs/propose/spec.md](specs/propose/spec.md),
[specs/changes/spec.md](specs/changes/spec.md) and
[specs/init/spec.md](specs/init/spec.md).

## Scope

### In scope

- `README.md` and every page in `docs/`, guides included.
- Moving each exact rule from `docs/` into a maintained spec, adding a
  requirement where none holds it.
- The `config` spec's Boundaries line that points to `docs/reference.md`.
- `dev/walkthrough.md`, so it still runs the docs' steps (rule 8), and
  `mkdocs.yml` where a page title changes.

### Out of scope

- Go code. If a rule the docs state is not what the code does, the owner
  decides; this change does not fix code.
- Grove's instruction text in `internal/instruct`: change 261023-ejarr.
- How maintained specs are laid out: change 261023-d6r79.
- What changes in cycle B do to the docs. This change starts from the
  docs as they leave them.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **The README and the home page lead with the attended path:** propose,
  review and merge, implement, merge. Cycles and `grove run` come after,
  under a "Doing more" heading. Reason: most readers start with one
  change, and today the second paragraph of `README.md` is about
  unattended runs.
- **A plain-language pass over `docs/`:** one idea per sentence, each term
  defined once in the glossary (`docs/concepts.md`, "Glossary"), long
  conditional sentences broken up. The sentence in "Two halves" that
  `grove run` made false is fixed. Reason: the docs explain; precision
  belongs to the specs.
- **Each exact rule taken out of `docs/` is kept in `grove/specs/`,** with
  a requirement there if no maintained requirement holds it. The docs
  link the spec instead of restating its rules. Reason: amended rule 8.
- **Scoped as three tasks:** audit and move the rules, restructure the
  README and home page, and the plain-language pass. The requirements
  found now are ADDED here. Others found while auditing go in a revision
  of this change agreed with the owner while implementing. Reason: no
  one can list every rule before the audit.
- **It comes last in cycle C.** It depends on 261023-kh3gg, 261023-ejarr
  and 261023-d6r79, and follows cycle B, because all of them edit the
  docs.

## Open questions

None must be answered before review.

- **How a docs page links a spec.** The docs site is built from `docs/`
  alone, so a relative link to `grove/specs/` does not resolve there.
  Recommended: a full GitHub link to the spec on `main`. It can wait; it
  decides the link form in T-01.
- **Whether the lead names one merge for a small change.** Change
  261023-6yafg lets a change drafted on its implementing branch be
  approved with its code in one merge. Recommended: one sentence after the
  four steps, saying a small change can be one merge. It can wait; it
  decides a sentence of T-02.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md),
  [specs/changes/spec.md](specs/changes/spec.md),
  [specs/init/spec.md](specs/init/spec.md)
- Tasks: [tasks.md](tasks.md)
