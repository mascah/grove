---
id: "W-029"
type: work
title: "Migrate the brief, plans and reviews into typed Grove records"
status: proposed
created: "2026-09-21T04:37:58Z"
updated: "2026-09-21T04:39:27Z"
kind: refactor
size: medium
depends_on: ["W-019"]
relates_to: ["D-005", "W-018", "W-023"]
---

## Outcome

This repository keeps one layout: the brief, every plan and every review live
as typed records under the Grove root, and nothing editable remains in the old
`docs/` locations. Owner's intent, 2026-09-20, recorded in
[D-005](../decisions/D-005-typed-knowledge-records.md): "it will be weird to
have [two] different layouts of things."

## Constraints

In: `docs/restart-brief.md` to the brief location W-019 supports under the
record root; the plans in `docs/plans/` and reviews in `docs/reviews/` to `P-`
and `R-` records created with `grove new`, each with its `work` list; every
link to them rewritten. Out: changing what any document says, the guides
(`docs/work-execution.md`, `docs/work-shaping.md`) and `docs/record-model.md`
unless W-019 gives them a home, sibling repositories (nullsec's cutover is
W-023), and `docs/prompts/`.

Observed at main `42c077d`: 12 plans and 9 reviews; 66 Markdown links in 42
files point into `docs/plans/` or `docs/reviews/`; 10 Markdown files link the
brief, including `AGENTS.md`; 6 Go source lines mention these paths. Most done
records' bodies, and so their content revisions, will change. Unmerged
branches (`versions` lists them) still hold the old paths.

Proposed, not decided: record the old-path to new-ID mapping in one place, as
[D-002](../decisions/D-002-sequential-ids.md) did for the renumbering, and
order IDs by each document's date. A review file covering several work items
becomes one record listing them all.

## Acceptance

1. `grove check` passes; a link check finds no Markdown link or Go path that
   still targets the old locations, and the old files no longer exist.
2. Every migrated document's body is byte-identical apart from rewritten
   links; the mapping lets a reader of an old commit or branch find the record.
3. `grove context` for a done work item lists its migrated plan and review
   through the mechanism W-019 ships, without including them by default.
4. `AGENTS.md`, the README, guides and adapters name the new locations, and
   there is one editable brief.
5. The owner judges the result browsable in the file tree and on the board.

## Next

Blocked on W-019 (`depends_on`). Then assign; it needs a short plan covering
the mapping, ID order, and how unmerged branches are handled. Treat it as the
rehearsal for W-023's nullsec migration.
