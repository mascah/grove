---
id: "W-030"
type: work
title: "Decouple record identity and storage from knowledge classification"
status: done
created: "2026-09-21T15:40:59Z"
updated: "2026-09-21T16:45:05Z"
kind: feature
size: medium
priority: 1
depends_on: ["W-019"]
relates_to: ["D-006", "W-018", "W-029", "W-020"]
---

## Outcome

People and agents can keep ordinary project knowledge alongside work and
evidence without choosing a schema category first, while Grove keeps identity,
placement and operational validation dependable. Owner-selected direction,
2026-09-21: [D-006](../decisions/D-006-stable-knowledge.md). This is the one
foundation change before W-029's migration, not another product restart.

## Constraints

Start from W-019's delivered schema 2. Preserve existing records, IDs, paths
and relationships during ordinary operations; W-029 owns the owner's explicit
one-time reconciliation of this repo's existing IDs and paths. Preserve brief
discovery, staged context and safe writes. New records share a neutral
sequential namespace; suggested spelling is `G-NNN` with the
existing minimum-three-digit convention. The exact schema and command surface
are preparation decisions within this outcome.

In: minimal general knowledge pages; flat creation; recursive, path-independent
discovery; type-independent identity; generic-page retrieval and deliberate
context inclusion; compatibility across old/new checkouts and shared allocator
state; documentation and both workflow guides. General pages have no mandatory
operational lifecycle. Retain known work/question/decision/term/plan/review
contracts and relationships where software depends on them. Reclassification
does not grant authority and must satisfy any newly applicable contract.

Out: migrating this repo's content (W-029), review lifecycle changes (W-020),
custom field/schema plugins, arbitrary lifecycle engines, per-type routing or
prefix settings, root relocation, automatic renames/moves, archive commands,
semantic search, wikilink syntax, a new TUI design or sibling migrations.

Provide a supported conversion path for W-029 to allocate neutral replacement
IDs through Grove, retain original metadata/provenance and rebuild references.
Preparation may choose a bounded conversion interface rather than a general
migration engine. Never require hand-numbered IDs or an ordinary field update
that silently changes identity. W-030 ships the foundation; W-029 performs the
repository-wide reconciliation and owns its mapping and verification.

Observed at `76da081`: the loader recursively walks the record root but filters
by first folder; `ParseRecord` checks folder/type and ID/type agreement;
creation chooses `TypeInfo.Folder` and allocates per type prefix. Prefix parsing,
CLI argument validation, shared counter recovery and source selectors need
inspection too: this is not just removing the folder check. Existing filenames
are not identities. Content revisions and exact source/path targeting must stay
honest. W-013's batched Git reads and W-012's on-demand history remain constraints.

## Acceptance

1. The CLI can create, show, list, update and check a general knowledge page
   without selecting a predefined semantic type or assigning work status.
   The documented minimal envelope preserves a stable ID and readable title.
   Explicit context inclusion works without preloading all knowledge bodies.
2. In the new schema, a valid record at the root or beneath arbitrary safe
   subfolders is discovered by identity and metadata, not location. The one
   brief remains separately identifiable. Invalid Markdown/metadata produces
   clear diagnostics; malformed operational records are not silently treated
   as valid generic pages. Define the generic-page boundary explicitly.
3. All newly created record types share the neutral namespace and allocator;
   concurrent local worktrees cannot allocate duplicate IDs. Recovery scans
   include nested records and local refs. Existing W/Q/D/T/P/R identities stay
   readable and resolvable before deliberate conversion and in historical
   sources, with no automatic renumbering or rewrite during ordinary operations.
4. Title, classification and lifecycle updates leave ID and path unchanged.
   Relationships, linked work artifacts and duplicate-ID checks remain correct;
   known type rules remain enforced. Generic knowledge cannot acquire work or
   approval behavior merely from folder placement or prose.
5. Document and exercise explicit schema migration, old schema readers and
   writers, and shared allocator state with mixed-version worktrees. An older
   checkout must not corrupt counters or issue a conflicting ID; unsupported
   mutation fails clearly. No command silently upgrades a project.
6. CLI/context, versions, workspace resolution and the existing board load
   both compatible old records and the new representation. Generic pages do
   not become work cards. Preserve freshness checks, escaping, cancellation,
   on-demand history and batched Git reads; test connected terminal behavior
   where affected. Routine reads remain read-only and require no allocator.
7. README, record model, guides and agent instructions describe actual shipped
   behavior and give W-029 an unambiguous conversion path for existing records
   as well as legacy documents. Exercise conversion on a disposable set with
   cross-record and shared-artifact references, preserving original metadata
   and examined commits. A rerun must not duplicate already converted records
   or silently remap identities. Relevant Go
   tests, uncached full/race suites, vet, formatting and `grove check` pass.

## Evidence

Branch `worktree-W-030` from main `bd6debe`. Plan
[P-001](../plans/P-001-w-030-flexible-records.md); independent review
[R-002](../reviews/R-002-w-030-flexible-records.md), two rounds, all findings
fixed. Verified on the final revision, uncached: `gofmt -l .` and `go vet
./...` clean, `go test -count=1 ./...` ok, `go test -race -count=1 -p 1 ./...`
ok, `grove check` ok. Known flakes that predate this work:
`TestContextLeavesEverythingUnchanged` against Git's maintenance lock, and
`internal/tui TestTerminal` in a parallel whole-suite race run (R-002
reproduced it at the base).

1. `internal/cli/flexible_test.go` creates, shows, lists, updates and checks a
   page; a related page is listed by `context` and read only by `--include`.
2. `internal/project/project_test.go`: root, former type folder and nested
   discovery, the brief beside them, and each page-boundary refusal by file
   and field; schema 2 still refuses every schema-3 form in its old words.
3. `internal/create/create_test.go`: separate `neutral-ids`, floor from a
   nested record in a ref and in another worktree, 12 concurrent creations
   across two worktrees. With a binary from `bd6debe` in a disposable clone,
   the old CLI kept issuing `W-` IDs for a schema-2 worktree and never saw the
   neutral counter; R-002 repeated this concurrently.
4. `internal/update/convert_test.go` `TestReclassifyKeepsIdentityAndPath`.
5. The record model's schema-3 section; the old CLI refuses a schema-3
   checkout with "unsupported version 3; expected 1 or 2". No command edits
   `schema_version`, and this repository is still schema 2.
6. The CLI test reads a schema-1 branch beside schema 3 through `versions`
   and `workspace`; `internal/tui` `TestBoardFollowsTypeNotIDOrPlacement`.
   Limit: neutral-ID work deleted from every source is not shelved.
7. `convert` tests cover cross-record and shared-plan references, `examined`,
   CRLF/BOM, reruns (by another case too) and refusals that reserve nothing.
   Rehearsed in a disposable clone on W-019, R-001, W-030, P-001, D-006 and
   the legacy W-019 plan: mechanical diffs only, `check` ok, reruns refused.

## Next

Implementation complete and independently reviewed; not merged or pushed. The
owner's judgment of the page envelope, `G-NNN`, and `convert` as W-029's
interface is theirs to give. Then merge `worktree-W-030` and assign W-029,
which sets `schema_version: 3` here, converts every record and legacy document
with `convert`, collects the mapping lines, moves the brief by hand and repairs
body links, none of which `convert` does.
