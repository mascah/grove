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
meaning given where they are used, and the glossary in `docs/concepts.md`
defines few of them, so no page says which meaning holds. One sentence is now false:
`docs/concepts.md`, "Two halves", says "Grove starts no agent and makes
no commits", and `grove run` starts one.

The docs and the specs can disagree with nothing saying which one holds.
Rule 8, amended in #78, makes the maintained specs in `grove/specs/` the
specification, has `docs/` explain them in plain words and link them,
lets a page restate a rule where a reader needs it to act, and says the
spec holds where the two differ. At 0c6f419 the docs restate rules the
specs hold without linking them: `docs/running.md` gives the conditions
for selecting a cycle that `grove/specs/runs/spec.md`, "Start one
attempt", holds. Other exact rules are in the docs and in no spec, so no
spec holds when the page is wrong. Examples: what `grove check` reports on
`depends_on` and on a stale requirement (`docs/coordinating.md`), the
JSON that `grove changes --json` prints, when `grove archive` refuses, and
how `grove init` chooses harnesses and writes their files
(`docs/reference.md`). The `config` spec's Boundaries says the rest of
`grove.yaml` "is specified by `docs/reference.md`".

## Intended outcome

`README.md` and `docs/index.md` show one change taken through first:
propose, review and merge, implement, merge. Cycles and `grove run` come
after, under "Doing more".

Every page says one idea per sentence. Each term is defined in the
glossary, which holds where a page's wording differs.

Each exact rule lives in a maintained spec. A page explains it plainly
and links the spec; it may restate the rule where a reader needs it to
act, and where the two differ, the spec holds and the page is fixed. The rules found so far that no spec holds become
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
- Holding each exact rule `docs/` states in a maintained spec, adding a
  requirement where none holds it, and fixing a page that differs from
  its spec.
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
  defined in the glossary (`docs/concepts.md`, "Glossary"), long
  conditional sentences broken up. The sentence in "Two halves" that
  `grove run` made false is fixed. Reason: the docs explain; precision
  belongs to the specs. Revised on 2026-10-24 (below): the glossary is
  the authority for a term, not its only definition.
- **Each exact rule a page states is held in `grove/specs/`,** with a
  requirement there if no maintained requirement holds it, and the page
  links the spec. Reason: amended rule 8. Revised on 2026-10-24 (below):
  a page may keep a rule it restates.
- **The spec is the authority, not the only copy (2026-10-24).** The
  owner said "Agreed" to this principle: "repeats cause bugs when copies
  disagree, not when they exist. So the rule should say which copy is the
  authority, not forbid copies." Rule 8, amended in #78, says so for the
  docs. So a page may restate a rule where a reader needs it to act, and
  links the spec. Where page and spec differ, the page is fixed to match
  the spec. The implementer proposes which restatements stay, and the
  owner agrees in the session. Likewise a page may say what a term means
  where a reader needs it, agreeing with the glossary. This replaces
  saying each thing once, and no page restating a requirement's
  conditions. The owner agreed to these revisions in this session.
- **Scoped as three tasks:** audit and move the rules, restructure the
  README and home page, and the plain-language pass. The requirements
  found now are ADDED here. Others found while auditing go in a revision
  of this change agreed with the owner while implementing. Reason: no
  one can list every rule before the audit.
- **It comes last in cycle C.** It depends on 261023-kh3gg, 261023-ejarr
  and 261023-d6r79, and follows cycle B, because all of them edit the
  docs.
- **Docs link a spec with a full GitHub link to it on `main`
  (2026-10-24).** The owner said "yes". Reason: the docs site is built
  from `docs/` alone, so a relative link to `grove/specs/` does not
  resolve there. It decides the links in T-01.
- **The README names one merge for a small change (2026-10-24).** The
  owner said "yes": after the README's four steps, one sentence says a
  small change can be proposed and built in one merge. Reason: change
  261023-6yafg lets a change drafted on its implementing branch be
  approved with its code in one merge. It decides that sentence in T-02.

## Open questions

None.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md),
  [specs/changes/spec.md](specs/changes/spec.md),
  [specs/init/spec.md](specs/init/spec.md)
- Tasks: [tasks.md](tasks.md)
