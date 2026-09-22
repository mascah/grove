---
id: "G-052"
type: work
title: "Reconcile all Grove content into neutral IDs and one flat layout"
status: done
created: "2026-09-21T04:37:58Z"
updated: "2026-09-22T00:43:59Z"
kind: refactor
size: large
depends_on: ["G-065"]
relates_to: ["G-051", "G-036", "G-041", "G-064"]
formerly: "W-029"
---

## Outcome

Reconcile all existing Grove records, legacy plans/reviews and the brief into
one flat layout under the configured Grove root using G-065's record contract.
Every record uses the neutral ID/filename convention, with one editable owner
per document. On 2026-09-21 the owner explicitly selected "Reconcile both IDs
and file locations" because multiple conventions were already causing friction;
[G-064](G-064-stable-knowledge.md) records that authority and revises
[G-051](G-051-typed-knowledge-records.md). The brief remains a
separate configured document, not a record.

## Constraints

In:

- Convert every existing record under `grove/`, including done work, decisions,
  questions, terms, plans and reviews, to a neutral ID and flat filename. Use
  G-065's supported conversion/allocation path; never hand-number replacements.
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
on completion, sibling writes (nullsec is G-041), and relocating repository
entrypoints or product/workflow documentation merely because they are Markdown.
The README, AGENTS.md, adapters, `docs/work-execution.md`, `docs/work-shaping.md`
and `docs/record-model.md` keep their functional homes with references updated.
Inventory other documents, including any `docs/prompts/` sources, and account
for their role rather than silently excluding project knowledge/evidence.

Observed on 2026-09-21 in the shaping checkout based on main `76da081`: 43 Grove
records, 13 legacy plans, nine legacy reviews and one brief. A scoped search
found 407 directory-reference occurrences in 62 Markdown/Go files, including
14 Go files; this is a search inventory, not an exact migration-edit count.
Refresh the inventory after G-065: its own plan/review and new records also
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

## Evidence

Branch `worktree-W-029` from main `70de539`. Plan
[G-068](G-068-reconciliation-plan.md); mapping and historical literals
[G-069](G-069-migration-map.md); independent review
[G-070](G-070-reconciliation-review.md), two rounds, all consequential and
minor findings fixed. Not merged or pushed.

1. 46 records, 13 legacy plans and 9 legacy reviews became `G-001` to `G-068`
   flat under `grove/`, in document-date order, through `grove convert`; G-069
   and G-070 were created there with `new`. The type folders, `docs/plans/` and
   `docs/reviews/` are gone; `grove.yaml` names `grove/brief.md`.
   `docs/prompts/*.txt` stays as spent evidence, accounted for in G-069.
2. G-069 has one row per source and each record's `formerly` matches its row.
   Round 1 diffed all 69 pairs: only mechanical ID, path and link edits, with
   `created`, `updated`, `status` and `examined` byte-identical. Converted
   legacy documents carry no `created`; their document dates are in G-069.
3. `grove check` OK at 70 records. Every relative link and anchor resolves
   (the reviewers' checkers and the migration script's audit). Remaining typed
   IDs and old folder names are the literals G-069 accounts for: fixtures,
   trial clones, the predecessor's records, verbatim quotations, branch names,
   paths, `formerly`.
4. `context G-030 G-031`, `context G-010` (shared plan G-013, review G-022),
   `show`, `brief` and `versions G-052` exercised on migrated proposed and done
   work. The board, driven in a pseudo-terminal on this checkout, shows
   Proposed 10, Active 1, Done 15 under neutral IDs and opens a card with its
   history; `git log --follow` on `grove/G-030-card-lineage.md` reaches its
   pre-migration commits. Limits: `show W-029` finds nothing, so an old ID is
   resolved through G-069 or `formerly`; a branch that predates the conversion
   is a source the current CLI cannot inspect (`unsupported version 2;
   expected 3`, exit 1, valid sources still print), and the board says
   INCOMPLETE while such local branches exist. Read them with their own CLI.
5. README, AGENTS.md, the record model, both guides, the adapters and the
   brief describe one schema. Rehearsed in disposable clones, by the
   implementer and again by the round-2 reviewer: the migration reproduced an
   identical mapping after reset, clean and deleting `neutral-ids`; a rerun was
   refused reserving nothing; a merge from an old-layout branch carried an edit
   to the record's new path, conflicted on a record added in the old layout,
   and a restored typed-ID file fails `check`. `new` issues the next `G-` ID
   flat, and title, status and type updates keep ID and path.
6. The owner's judgment is outstanding.

Schemas 1 and 2 are deleted (`dd3a6f5`): `schema_version` must be 3, an ID is
`G-NNN` only, one `neutral-ids` counter, no type folders, and `convert` takes
only a document path. Verified uncached on the final revision: `gofmt -l .`
and `go vet ./...` clean, `go test -count=1 ./...` ok (four whole-suite runs
in a row at `1a85fc7`), `grove check` ok. `go test -race -count=1 -p 1 ./...`
at `1a85fc7`: seven packages ok, `internal/versions` failed
`TestResolveFinalCheck` once. Examined on resume: that test alone passed under
race, and 20 race runs of the package (10 with Homebrew Git 2.55.0, 10 with
Apple Git 2.50.1) failed 3 times, each in a different test and each as a
fixture or `repo.GitContext` child dying before it ran (`git worktree` or
`git rev-parse`: `signal: segmentation fault`) or, once, as a forked child
spinning until the 10-minute test timeout. A `sample` of that child,
symbolized against a race-built test binary, was `__tsan::TraceSwitchPartImpl`
under `syscall.forkAndExecInChild`: the race detector's runtime in the forked
child before `exec`, with no Grove frame. That child also held copies of other
tests' `cat-file --batch` stdin pipes, so three unrelated subtests blocked in
`Wait` until it was killed. 10 non-race runs of the package in a row were
ok. Recorded as an environment flake of Go 1.26.2's race detector on macOS,
not a regression of `dd3a6f5`; no code changed. Rerun a failed race package,
and kill any orphaned `versions.test` process a timeout leaves behind. Also
seen and rerun green, likewise outside this change: `internal/tui`
`TestTerminal` ("files in the repository changed") in whole-suite runs.

## Next

The race failure noted at the 2026-09-21 checkpoint is examined and recorded
in Evidence as an environment flake; nothing else changed.

Implementation complete and independently reviewed; awaiting the owner's
judgment of the flat tree and ordinary CLI and board browsing (acceptance 6):
`ls grove/` and `go run ./cmd/grove` in `.claude/worktrees/W-029`. Then mark
this done and merge `worktree-W-029` (fast-forward from `70de539`). Open for
the owner while browsing: the 22 converted legacy plans and reviews carry no
`created`, because `convert` sets none and `update` keeps the field fixed, so
`list` sorts them last and their document dates live only in G-069; say
whether a follow-up should set them from Git history. The merged
local branches `worktree-W-004-W-005`, `worktree-W-006-W-008`, `worktree-W-010`,
`worktree-W-030` and `worktree-direction-reconciliation` still hold schema 1 or
2, so `versions` exits 1 and the board says INCOMPLETE until the owner deletes
them; they are all ancestors of main. G-041 reuses the script's approach from
`8bf691b` for nullsec but owns its own inventory.
