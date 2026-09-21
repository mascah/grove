---
id: "P-001"
type: plan
title: "W-030 flexible records: schema 3, neutral IDs, pages, conversion"
status: current
created: "2026-09-21T16:02:50Z"
updated: "2026-09-21T16:02:54Z"
work: ["W-030"]
---

## Design

Plan for [W-030](../work/W-030-flexible-records.md) within
[D-006](../decisions/D-006-stable-knowledge.md). Base `bd6debe`, branch
`worktree-W-030`. Settled here as preparation decisions:

**Schema 3, opt-in.** `schema_version: 3` is the deliberate one-line edit that
schema 2 was; no command rewrites it. Schemas 1 and 2 keep their rules and
diagnostics byte-for-byte. This repository stays at schema 2: W-029 migrates it.

**Discovery (schema 3).** Every `.md` beneath the record root, at any depth, is
a record, except the configured brief. Folders mean nothing: no type-folder
check, no "type folder must be a directory", and the brief may sit anywhere,
including the record root. Symlink and regular-file rules are unchanged. A
`.md` that is not a valid record is still a diagnostic, never skipped.

**Identity (schema 3).** An ID is `[WQDTPRG]-NNN` with canonical padding and no
tie to `type`. `new` issues `G-NNN` for every type and writes
`ROOT/G-NNN-slug.md`. Existing typed IDs and nested paths stay valid and are
never rewritten by ordinary commands.

**Generic page boundary.** A page is exactly a record with `type: page`.
Envelope: `id`, `type`, `title`; optional `relates_to`, `created`, `updated`,
`formerly`. `status` and every operational field are refused on a page. A
missing or unknown `type` is an error, so a damaged operational record cannot
fall through to a page. Known types keep every current rule. `grove new page
TITLE`; list, versions and context print `-` for a page's status. Pages are
never work cards (the board selects `type == work`; R-002 found it read the
first source's type, now the board source's own), cannot be selected
by `context`, and are listed when related and read only by `show` or
`--include`.

**Reclassification (schema 3).** `update --set type=T` is accepted; the
candidate must satisfy T's whole contract in the same request (for example
`--set type=work --set status=proposed`, or `--unset status` toward a page).
ID and path do not change. Schemas 1 and 2 still refuse `type`.

**Allocator.** Neutral IDs use their own counter file, `grove/neutral-ids`
(same one-line format, `G N`), under the same `grove/lock`. `next-ids` is never
given a `G` line, so an older checkout keeps allocating typed IDs for its
schema-1/2 project and cannot see, corrupt or collide with the neutral
counter; the two namespaces cannot produce the same ID. The recovery scan is
unchanged: already recursive over the record root in every local ref and
worktree. An older CLI refuses a schema-3 checkout with "unsupported version 3".

**Conversion interface for W-029 (schema 3 only), `grove convert`:**

- `convert ID [--slug SLUG]`: a record with a typed ID gets the next neutral
  ID. Only its `id` changes and `formerly: "OLD-ID"` is added; timestamps,
  status, `examined` and body bytes are kept and `updated` is not bumped. The
  file moves to `ROOT/G-NNN-slug.md` (slug from the old filename). Every other
  record's relationship lists naming the old ID are rewritten, again without
  touching `updated`.
- `convert PATH --type TYPE --title TITLE [--slug SLUG]`: a project-relative
  legacy Markdown document outside the record root becomes a new record whose
  body is the document's bytes, with `formerly: "PATH"`. The original is left
  for the caller to remove.
- Stdout is one JSON mapping line `{from, from_path, id, path}`. A source
  already named by some record's `formerly`, or an ID that is already neutral,
  is refused before anything is reserved, so a rerun neither duplicates nor
  remaps. `formerly` must be unique across records.
- Body prose and Markdown links are not rewritten: W-029 owns them and the
  durable mapping. The whole candidate set is validated in memory before the
  first write. Adjusted during implementation: the new file is created first
  (`O_EXCL`, the likeliest refusal, which then leaves everything untouched),
  then the old file is removed, then referrers are replaced by temp+rename. It
  is not atomic across files: an interruption leaves a project that fails
  `check`, recovered with Git.

Not built: alias lookup of old IDs, link rewriting, batch conversion, per-type
settings, a schema-upgrade command.

## Steps

1. `internal/project`: `page` type, schema 3 loader and `ParseRecord` rules,
   `formerly`, uniqueness; tests for root/nested discovery, the page boundary,
   brief placement, unchanged schema-1/2 diagnostics.
2. `internal/create`: neutral allocation and flat paths in schema 3; tests for
   the separate counter, nested/ref recovery, concurrent worktrees, mixed
   schema-2 and schema-3 worktrees.
3. `internal/update`: `type` in schema 3; `convert.go`; tests with
   cross-record and shared-plan references, `examined`, rerun refusal.
4. `internal/cli`: `new page`, `convert`, `-` status, usage; context, versions,
   workspace and board tests over a schema-3 tree with a page.
5. Docs: record model, README, both guides, `AGENTS.md`; brief only where
   direction wording says support is unshipped.
6. Full verification, independent review (one gate, final revision), reconcile.
