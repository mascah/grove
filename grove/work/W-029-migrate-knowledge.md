---
id: "W-029"
type: work
title: "Reconcile all Grove content into neutral IDs and one flat layout"
status: active
created: "2026-09-21T04:37:58Z"
updated: "2026-09-21T21:00:42Z"
kind: refactor
size: large
depends_on: ["W-030"]
relates_to: ["D-005", "W-018", "W-023", "D-006"]
---

## Outcome

Reconcile all existing Grove records, legacy plans/reviews and the brief into
one flat layout under the configured Grove root using W-030's record contract.
Every record uses the neutral ID/filename convention, with one editable owner
per document. On 2026-09-21 the owner explicitly selected "Reconcile both IDs
and file locations" because multiple conventions were already causing friction;
[D-006](../decisions/D-006-stable-knowledge.md) records that authority and revises
[D-005](../decisions/D-005-typed-knowledge-records.md). The brief remains a
separate configured document, not a record.

## Constraints

In:

- Convert every existing record under `grove/`, including done work, decisions,
  questions, terms, plans and reviews, to a neutral ID and flat filename. Use
  W-030's supported conversion/allocation path; never hand-number replacements.
- Migrate all legacy documents in `docs/plans/` and `docs/reviews/` to the same
  convention. Preserve plan/review roles, shared ownership and examined commits.
- Move `docs/restart-brief.md` to `grove/brief.md` and update `grove.yaml`.
- Rewrite current ID relationships, Markdown links and operational references
  throughout records, README, guides, adapters, instructions and code where
  they target migrated content. Retire superseded copies and empty type folders.
- Retain one durable old-ID/path to new-ID/path mapping, including legacy files
  that had no record ID. It must let a reader of an old commit locate the current
  counterpart. No permanent duplicate records or compatibility symlink tree.
- After the conversion, delete schema 1 and 2 support and the compatibility
  kept for them: folder/type and prefix/type rules, per-type counters,
  schema-gated wording, and old/new-checkout tests. The owner directed on
  2026-09-21 that Grove keeps no backward compatibility before its first
  release; one current schema remains. An old commit stays inspectable with
  the CLI in that commit (`go run ./cmd/grove` there), not the current one.

Preserve original timestamps, statuses, authority, historical conclusions and
evidence. Only mechanical schema/identity/path transformations are in scope.
Changing a reference does not make an old review examine the migration commit.
Old identifiers or paths may remain in clearly marked historical quotations,
evidence and the migration mapping; current instructions must use the new ones.

Out: rewriting Git history, changing completed work's meaning, automatic filing
on completion, sibling writes (nullsec is W-023), and relocating repository
entrypoints or product/workflow documentation merely because they are Markdown.
The README, AGENTS.md, adapters, `docs/work-execution.md`, `docs/work-shaping.md`
and `docs/record-model.md` keep their functional homes with references updated.
Inventory other documents, including any `docs/prompts/` sources, and account
for their role rather than silently excluding project knowledge/evidence.

Observed on 2026-09-21 in the shaping checkout based on main `76da081`: 43 Grove
records, 13 legacy plans, nine legacy reviews and one brief. A scoped search
found 407 directory-reference occurrences in 62 Markdown/Go files, including
14 Go files; this is a search inventory, not an exact migration-edit count.
Refresh the inventory after W-030: its own plan/review and new records also
belong in the reconciliation. Unmerged branches retain old files and IDs.

Preparation must produce the complete mapping before publication, verify a
disposable rehearsal and define recovery/rerun behavior. Proposed ID order is
document date with a stable tie-breaker; it conveys no authority or priority.
Use one record for an artifact shared across work items. Address reintegration
from old branches explicitly so a later merge cannot quietly restore duplicate
old-layout records. Do not rewrite other sessions' worktrees.

## Acceptance

1. An inventory accounts for every prior record, plan, review and the brief.
   All resulting records live directly under `grove/` with neutral IDs and
   matching generated filenames. The former type folders and legacy plan/review
   locations contain no remaining content; there is one configured brief.
2. Every original document has exactly one mapped counterpart, including shared
   artifacts. Bodies differ only by documented mechanical ID/path/schema edits;
   metadata and historical evidence retain their meaning. The mapping resolves
   old identities and paths without retaining a second editable authority.
3. `grove check` and a repository-wide link/reference audit pass. No current
   operational reference targets a removed path or obsolete ID. Historical
   literals are explicitly accounted for. Check bare IDs and inline-code paths
   as well as Markdown links; do not blindly replace historical evidence.
4. Context, attachments, dependencies, explicit lookup, board/detail and history
   are exercised on migrated proposed and done work. Relationships use the new
   IDs; shared plans/reviews remain discoverable without preloading their bodies.
   Historical commits remain inspectable with their own CLI, and lineage
   limits are documented.
5. README, instructions, guides and adapters use the reconciled convention;
   no new authoring path produces the old layout. A disposable rehearsal checks
   restart/rollback and reintegration from a branch containing old IDs/paths.
6. The owner judges the flat file tree and normal CLI/board browsing coherent.
   Subsequent title/type/status changes keep the new ID/path stable. No ongoing
   archive or completion-driven move is introduced.

## Next

W-019 is delivered. Wait for W-030's foundation, then assign this complete
one-time reconciliation. Prepare a linked plan covering full inventory and
mapping, allocation, references/evidence, rehearsal, recovery and unmerged
branches. Full ID and location reconciliation is already selected; do not ask
again whether existing content should be included. W-023 can reuse the lessons
but still owns its separate nullsec inventory and live-cutover authority.
