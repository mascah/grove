# Starter record model

Status: core model accepted, 2026-09-18, on the owner's response "yup this looks
good" to the draft. This covers the representation, fields, relationships,
lifecycles, and initial validation boundary below. The owner accepted the
[on-disk defaults below](#on-disk-contract) on 2026-09-19, then replaced the
random-ID/long-filename trial with shared sequential IDs and short filenames.
The first Go CLI now reads and validates this model; G-003 records its evidence.
[The restart brief](../grove/brief.md) owns product direction;
`grove/` owns operational work, questions, and decision receipts.

## Schema 3: identity and placement apart from classification

[G-064](../grove/G-064-stable-knowledge.md) selected stable identity
and placement, general knowledge pages, flat creation and recursive discovery
independent of type folders. [G-065](../grove/G-065-flexible-records.md)
implements it as `schema_version: 3`. Everything in the schema-1 and schema-2
sections below still holds under schema 3 except where this section says
otherwise, and a schema-1 or schema-2 project keeps exactly its old rules.
[G-052](../grove/G-052-migrate-knowledge.md) moved this repository to schema 3
and converted every record; [G-069](../grove/G-069-migration-map.md) is the
old-to-new mapping.

- **Discovery.** Every `.md` file beneath the record root, at any depth, is a
  record, except the configured brief. No folder names a type: the root, a
  former type folder and any other nested folder are equally valid, so records
  written under schema 2 stay valid where they are. A `.md` file there that is
  not a valid record is a diagnostic, never skipped, dot folders included; only
  the lowercase `.md` extension counts. The brief may be any clean
  project-relative `.md` path, including one inside a former type folder.
  Symlinks are refused as before.
- **Identity.** An ID is one of the letters `W Q D T P R G`, a hyphen, and a
  canonical number (`G-001`, `G-1000`). It carries no type: `new` issues
  `G-NNN` for every type, and an existing `W-019` stays `W-019` whatever
  happens to its `type`. IDs are still unique per checkout and matched exactly.
- **Placement.** `new` writes `ROOT/G-NNN-slug.md`, flat. No command moves or
  renames a record when its title, type or status changes. `convert` below is
  the only command that moves a file.
- **Pages.** `type: page` is general knowledge with no lifecycle. Its envelope
  is `id`, `type` and `title`, with optional `relates_to`, `created`, `updated`
  and `formerly`; `status` and every work, question, plan or review field are
  unknown fields on a page. This is the whole boundary: a record is a page only
  because its `type` says so. A missing or unknown `type` is an error, so a
  damaged operational record never degrades into a valid page, and the six
  known types keep all their rules wherever they sit. A page is never a work
  card, cannot be selected by `context`, cannot be a `depends_on`, `members`,
  `blocks` or `work` target, and gains nothing from its folder or its prose.
  `context` lists a related page (its status shown as `-`) and reads it only
  through `--include PATH`; `show G-NNN` prints it.
- **Reclassification.** `update ID --expect REVISION --set type=TYPE` changes
  classification in place. The result must satisfy the new type's whole
  contract in that one update, for example `--set type=work --set
  status=proposed` from a page, or `--set type=page --unset status` toward
  one, and the project must still validate, so work that a plan or question
  names cannot stop being work. Reclassifying grants nothing: an `accepted`
  decision is a claim in a file, as it always was. Schemas 1 and 2 still
  refuse `type`.
- **`formerly`.** An optional string that only `convert` writes and `update`
  refuses: the typed ID or document path a record replaced. Two records with
  one `formerly`, or a `formerly` naming an ID that still exists in the
  checkout, are errors, which is how a merge from a branch that predates a
  conversion is caught instead of quietly restoring a second owner.

**Allocation.** Neutral numbers come from `grove/neutral-ids` in the Git common
directory, beside `grove/next-ids` and under the same `grove/lock`, with the
same format (`G 12`), floor scan (every `.md` beneath the record root in every
local ref and worktree, nested included) and recovery notices. The two files
are separate on purpose. A CLI that predates schema 3 refuses a `next-ids`
holding a prefix it does not know, so `G` never goes there: an older checkout
keeps allocating `W-`/`Q-`/… IDs for its schema-1 or schema-2 worktree, never
reads or rewrites the neutral counter, and cannot issue an ID a schema-3
worktree could also issue, because the namespaces are disjoint. The current CLI
likewise still issues typed IDs in type folders for a schema-1 or schema-2
worktree of the same repository. Exercised with a binary built from `bd6debe`
beside this one; see G-065's evidence.

**Moving to schema 3** is the deliberate one-line edit of `schema_version`, as
schema 2 was. No command rewrites it and every schema-2 record stays valid and
in place. It cannot be undone by editing the number back once a page, a neutral
ID or a record outside a type folder exists: schema 2 refuses those. An older
CLI refuses the checkout with "unsupported version 3; expected 1 or 2" on every
command, reads included, until that checkout has newer code. `versions`,
`workspace` and the board judge each branch and checkout by its own
`grove.yaml`, so schema-2 and schema-3 sources are read side by side. Where
sources disagree about a record's type, the board follows the record in the
checkout it is showing, and shelves work that checkout does not hold. One
limit: a group deleted from every source carries no record, so only a typed
`W-` ID can still say it was work; deleted neutral-ID work is not shelved.

**Conversion (`grove convert`, schema 3 only).** The one deliberate identity
change, bounded to one source per run; it is what
[G-052](../grove/G-052-migrate-knowledge.md) uses, never hand-numbering:

- `convert ID [--slug SLUG]` takes a record with a typed ID. It reserves the
  next neutral ID, changes only `id`, appends `formerly: "OLD-ID"`, and moves
  the file to `ROOT/G-NNN-slug.md` (slug from the old filename unless given).
  `created`, `updated`, `status`, `examined` and every body byte are kept.
  Relationship lists in other records that name the old ID are rewritten to
  the new one (a rewritten list is written in flow style), again without
  touching `updated`.
- `convert PATH --type TYPE --title TITLE [--slug SLUG]` takes a Markdown
  document outside the record root, such as a plan written before plan
  records. It creates a record whose body is the document's bytes (less a
  leading byte-order mark), with the type's first status, `formerly: "PATH"`
  and no invented dates. The brief is refused: it is not a record, and moving
  it is a file move plus an edit of `brief:` in `grove.yaml`. The original
  is left for the caller to remove; set `work`, `status` or `examined`
  afterwards with `update`.
- Stdout is one JSON line, `{from, from_path, id, path}`: the mapping entry.
  Collecting these is the caller's durable old-to-new mapping.
- A source that some record's `formerly` already names (compared without case,
  since a case-insensitive filesystem opens `docs/Plan.md` as `docs/plan.md`),
  an ID that is already neutral, and a missing source are refused before any
  ID is reserved, so a rerun neither duplicates a record nor remaps an
  identity. A refusal after the reservation, such as an existing target file,
  consumes the number and says so; gaps are acceptable, as for `new`.
- Not rewritten: body prose, Markdown links (the moved file's own relative
  links included), `examined`, and anything outside the record root. The
  caller repairs links from the mapping; `check` does not verify them.
- Everything is validated before the first write. The writes themselves are
  not atomic across files: the new file is created first (an existing target
  refuses the run untouched), then the old file is removed, then referrers
  are replaced. A failure after the new file exists exits 1 but still prints
  the mapping line, which is already true of the files. An interruption leaves
  a checkout that fails `check`, naming the half-converted record; recover with
  Git and run it again. Convert from a
  clean, committed record root.

Not provided: lookup of a record by its former ID, batch conversion, link
rewriting, per-type folders or prefixes, and a command that edits
`schema_version`.

## Schema 2: knowledge records and the brief

On 2026-09-20 [G-035](../grove/G-035-interactive-adoption.md) selected
terms and linked work artifacts, and
[G-051](../grove/G-051-typed-knowledge-records.md) selected their
representation. [G-037](../grove/G-037-knowledge-artifacts.md) implements
it as `schema_version: 2`:

| Type | ID | Folder | Statuses (first is what `new` writes) | Extra fields |
| --- | --- | --- | --- | --- |
| `term` | `T-001` | `terms/` | `proposed`, `settled` | none |
| `plan` | `P-001` | `plans/` | `current`, `superseded` | `work` |
| `review` | `R-001` | `reviews/` | `current`, `superseded` | `work`, `examined` |

- A term's title is the term; its body gives meaning and relationships, not
  execution instructions. Two terms whose titles match, ignoring case and
  surrounding space, are an error.
- `work` is an optional list of work IDs the plan or review belongs to, checked
  like `depends_on` targets. One plan can name several work items. Work does
  not name its plans or reviews back: that side is derived, and
  `context W-NNN` lists them without reading them. `new` takes no fields, so
  set it with `update ID --expect REVISION --set 'work=["W-001"]'`.
- `examined` is an optional quoted Git commit, 7 to 40 lowercase hex digits:
  what the review looked at. Whether the reviewed content has changed since is
  a comparison a reader makes, not stored state. Approval, candidates, and
  dispositions belong to G-038, which may extend the review type. A `report`
  type is not defined yet.
- Optional `brief: PATH` in `grove.yaml` names the one project brief: a clean
  project-relative `.md` path without `..`, anywhere in the project, including
  directly under the record root (`grove/brief.md`), but never inside a type
  folder (compared without case). It is not a record and has no ID or frontmatter rules. That one path
  is exempt from the misplaced-Markdown rule. Every live command requires a
  regular, non-symlink file there and names `grove.yaml: brief` when it is
  missing. `grove brief [--json]` prints it like `show`; `context` never adds
  it, and `--include PATH` still can.

**Compatibility.** The CLI reads schemas 1, 2 and 3
([schema 3](#schema-3-identity-and-placement-apart-from-classification) has its
own compatibility rules). A schema-1 project
keeps exactly the rules below: the three original types, no `brief` key, and
`new term` refused before any ID is reserved. Moving to schema 2 is the
deliberate one-line edit of `schema_version`; no command rewrites it, and every
schema-1 record stays valid. An older CLI refuses a schema-2 project with
"unsupported version 2", and refuses `new` once a newer CLI has written a `T`,
`P`, or `R` line to the repository's shared counter file, until that checkout
has the newer code. Committed sources in `versions` and the board check only
the form of `brief`, never that the file exists: they read `grove.yaml` and the
record root, and only a live checkout is required to hold the brief.

Plans and reviews written before this support were ordinary files in
`docs/plans/` and `docs/reviews/` until
[G-052](../grove/G-052-migrate-knowledge.md) converted them and moved the
brief. Not implemented: a Review work status and accepted-and-integrated
completion ([G-038](../grove/G-038-review-lifecycle.md)). Do not write
`status: review` yet. Historical Done records retain their original
branch-local acceptance meaning, not proof of merge.

## Minimum representation

Use one Markdown file per record, with YAML frontmatter for facts that the
CLI interprets and a freeform Markdown body for explanation.

| Field | Purpose |
| --- | --- |
| `id` | Stable identity; renaming the title or file does not change it |
| `type` | `work`, `question`, or `decision`; with schema 2 also `term`, `plan`, or `review`; with schema 3 also `page` |
| `title` | A readable label for lists, search, and the eventual board |
| `status` | Explicit lifecycle state for that record type; a schema-3 `page` has none |

The accepted relationship fields cover the starter records:

- Work `depends_on`: prerequisite work IDs. Missing means no declared prerequisites.
- Question `blocks`: work IDs whose outcome needs the answer. An unresolved
  question may be relevant without blocking work.
- Any record `relates_to`: related record IDs, with no implied ordering or gate.
- Plan and review `work` (schema 2): the work they belong to.

## Work planning metadata

On 2026-09-18, the owner proposed richer work metadata after comparing nullsec's
W-032 release, specifically members, dependencies, priority, size, and a work
sub-kind, to support useful UI behavior. The owner then accepted adding these
optional fields and the distinctions below ("yeah this is great"). Concrete
values were accepted with the on-disk defaults on 2026-09-19; no automatic
execution policy is implied.

| Field | Meaning and use |
| --- | --- |
| `kind` | What sort of work this is, for filtering and distinct presentation |
| `priority` | Importance for selection, independent of dependencies |
| `size` | Coarse scope/effort estimate, without automatic time estimates or execution policy |
| `members` | Child work included in this outcome; expandable groups and member completion counts |
| `depends_on` | Existing accepted field: prerequisite work that must be delivered before this work can proceed |
| `created`, `updated` | Creation and modification timestamps for chronology; potentially common to all record types |

Keep these optional so quick capture remains useful. An absent size or priority
means unspecified; do not silently turn it into an estimate or urgency decision.
Dates should be written by CLI mutations when those exist. The on-disk contract
below specifies the timestamp format and a proposed direct-editor policy.

Preserve these distinctions when implementing the extension:

- Membership describes decomposition. Member order can express presentation or
  preferred sequence, but only dependencies impose prerequisite ordering.
- Allow grouped work to retain its own outcome and acceptance. Completed children
  do not automatically establish parent completion or integration.
- Incomplete members can prevent declaring a parent done without blocking work
  on that parent. Do not conflate an unfinished group with an execution blocker.
- Derive member counts and blocker explanations from relationships and the
  selected branch context. Do not store parallel progress percentages or an
  independently editable `blocked` flag.
- A spike/investigation is a kind, not a size. The predecessor's size vocabulary
  partly chooses preparation depth; do not inherit those execution rules merely
  by accepting a size field.
- Validate member targets and membership cycles as well as dependency cycles.
  Nesting and shared membership follow the rules below; cross-branch
  resolution remains open.

Evidence: `grove context --work W-032 --phase shape` in nullsec returned
`kind: release`, `size: large`, `priority: 2`, `depends_on: []`, eight ordered
members, creation/update dates, and parent-level acceptance. That context call
also reported unrelated supporting sources omitted by its token budget; the
owning work record itself was present and inspected. Its `scope` points to
capability records, a type outside the current three-type starter model.

## On-disk contract

Defaults accepted 2026-09-19 on the owner's response "YeaI accept those defaults":
`grove.yaml` with schema version and configurable record root, the three type
folders, timestamps, and optional planning values. Later that day, the owner
accepted sequential IDs allocated across local worktrees and short filenames
after finding the original random-ID/timestamp filenames difficult to browse.
These revised defaults govern the operational records. G-003 implements the
discovery, validation, graph, and inspection behavior below; G-007 implements
creation and allocation per [G-006](../grove/G-006-allocator-mechanism.md).
[G-009](../grove/G-009-update-records.md) implements field updates with
content revisions and a shared write lock; renames, moves, deletes, and body
edits remain unimplemented. This section owns the schema; the brief owns direction.

### Configuration and discovery

Keep one `grove.yaml` at the project root:

```yaml
schema_version: 1
records: grove
```

`schema_version` versions the configuration and record schema together. Require
both keys; accept 1 and 2 ([what 2 adds](#schema-2-knowledge-records-and-the-brief));
refuse missing or unsupported versions without guessing, migrating,
or rewriting files. This is the new CLI's schema 1, unrelated to the sibling
skills CLI's schema numbering or `grove.toml` configuration.

Resolve `records` relative to the directory containing `grove.yaml`. Permit a
different relative folder, such as `docs/grove`; require it to remain inside the
project and reject absolute paths, `..` components, and symlink traversal. The
record root must be a dedicated subdirectory, not the project root itself.

Without an explicit project, search from the current directory upward for the
nearest `grove.yaml`, stopping at the current Git checkout root when inside one.
Outside Git, search up to the filesystem root. An explicit `--project <dir>`
selects the directory containing the configuration without upward searching.
Git is not required to read a configured project. Never consult another
worktree, the Git common directory, or the predecessor's configuration for
these initial commands.

### Folders and files

```text
grove.yaml
grove/
  work/
  questions/
  decisions/
```

Schema 2 adds `terms/`, `plans/`, and `reviews/` beside them; schema 3
[drops type folders](#schema-3-identity-and-placement-apart-from-classification). Read `.md` files
recursively within the type folders of the project's schema. Require frontmatter
`type` to match its folder; nested folders
may organize records but do not confer lifecycle meaning. Closed records stay
discoverable in the same tree. Reject symlinks in the record tree; report `.md`
files outside the type folders as misplaced rather than silently dropping
them (the configured brief excepted). Other file extensions are not records. A missing record root is an error;
missing type folders simply contain no records.

Generate short filenames as `<id>-<slug>.md`, for example:

```text
work/W-001-inspect-records.md
questions/Q-001-branch-versions.md
decisions/D-001-starter-defaults.md
```

Keep dates in frontmatter. The number supplies allocation order within each
record type; it does not prove creation time, priority, or execution order.
Use a short descriptive slug, with the full title in frontmatter. `grove new`
accepts an explicit `--slug`, or derives lowercase ASCII
letters/digits separated by hyphens from the title, trim to at most 32 characters
and strip trailing hyphens, falling back to `record` when empty. Generated names
stay short even when titles are long. Title edits do not extend the filename.

The generated filename is a convention, not a validation requirement. Read the
ID, type, and dates from frontmatter; accept manually named or renamed Markdown
files. Title edits do not automatically rename files. Relationships use IDs and
survive renaming; ordinary Markdown path links still need updating when moved.

### Identity and dates

Before schema 3, which
[issues neutral `G-` IDs](#schema-3-identity-and-placement-apart-from-classification)
from a counter of their own, use type-prefixed sequential IDs: `W-001` for work, `Q-001` for questions, and
`D-001` for decisions; schema 2 adds `T-`, `P-`, and `R-`. Each type has its own counter; the full prefixed ID is the
canonical identity, not an alias for a hidden random value. Start at 1, pad to a
minimum of three digits, and expand beyond 999 (`W-1000`) without wrapping or
renumbering older records. Require canonical padding and a prefix matching
`type`. Store IDs as strings and match references exactly; `show W-001` needs
no abbreviated-ID lookup. Numeric ordering must not rely on lexicographic
sorting once the counter expands.

For Git projects, creation commands coordinate through the shared
directory returned by `git rev-parse --path-format=absolute --git-common-dir`.
Do not derive it from a worktree's `.git` path: linked worktrees have private
metadata as well as a shared common directory. This is supported by
[Git's worktree documentation](https://git-scm.com/docs/git-worktree#_details).

Allocation requirements for cooperating Grove commands in one local repository:

1. Acquire an exclusive lock shared by every worktree and all record types.
2. Reserve the next number for the requested type and durably save the advanced
   counter while holding the lock. If locking or persistence fails, do not issue
   an ID or create a record.
3. Release the lock, then create the record in its selected checkout without
   overwriting an existing file. A failed or abandoned creation consumes the
   reservation; gaps are acceptable and numbers are not deliberately recycled.

The counter belongs to the local repository's Grove project, not to a branch,
worktree, or configured record-root path. It is local coordination state, not a
tracked project record; no daemon is required. A Grove-owned directory under
the Git common directory is the intended home. [G-006](../grove/G-006-allocator-mechanism.md)
owns the accepted state encoding, lock primitive, and recovery protocol for
[G-007](../grove/G-007-create-records.md) to test with the creation command. `grove new` initializes and maintains `grove/next-ids` under `grove/lock` in
that common directory, and `new` and `update` serialize publication through
`grove/write.lock` beside them; the read-only inspection commands never create
any of these files, and none is ever unlinked.

Separate clones do not share reservations. Directly authored IDs also bypass
allocation. Imported records and independently allocated clone histories require
collision checks and explicit reconciliation; matching IDs alone cannot prove
two independently created records are the same item. Two files with one ID in
one checkout remain an error even when their contents match. Genuine branch
copies of one record retain their identity.

Missing, corrupt, or restored-old counter state must not silently restart at 1.
Before enabling creation, define explicit initialization/recovery that considers
committed records across relevant refs and live records in linked worktrees.
Scanning records can find used numbers but cannot recover reservations for
deleted or never-written records; the recovery policy must expose that limit.
Concurrent imports and manual edits are outside the allocator's exclusivity
guarantee. Plain-directory reading remains supported without Git or allocator
state; automatic allocation outside Git is deferred.

The starter records were renumbered once before any new CLI existed. The mapping
is retained in [G-004](../grove/G-004-sequential-ids.md); this is not a
general ID-renaming feature. Schema 1 is still the unshipped starter contract.

Keep `created` and `updated` optional on every type. When present, require quoted
UTC timestamps in `YYYY-MM-DDTHH:MM:SSZ` form, and require `updated >= created`
when both exist. `new` writes both with the same current time; `update`
preserves `created` and sets `updated` when record content changes.
Do not invent a missing creation date for an existing file or refresh dates on
a read or no-op mutation.

Direct editors should update `updated` when changing content, but that is an
authoring convention rather than a provable freshness guarantee. Readers never
repair dates or infer them from filenames, filesystem modification time, or Git.
Revision checks for safe writes compare actual content, not these timestamps:
`show --json` reports `sha256:` plus the hex digest of the exact file bytes,
and `update --expect` refuses any other current content.

### Initial planning values

Use these accepted values when a field is supplied:

| Field | Initial allowed values | When absent |
| --- | --- | --- |
| `kind` | `feature`, `fix`, `refactor`, `investigation`, `tooling`, `release` | Unspecified |
| `priority` | Integer 1 (highest) through 5 (lowest) | Unspecified |
| `size` | `small`, `medium`, `large` | Unspecified |
| `members`, `depends_on` | Ordered lists of work IDs | No declared relationships |

The planning fields above belong only on work. `blocks` belongs only on
questions; `relates_to`, `created`, and `updated` can appear on every type.
Do not synthesize a default kind, priority, or size. An investigation uses
`kind: investigation` and an independent size; there is no separate `spike`
value in this starting vocabulary.

Allow nested and shared membership; grouping does not assign exclusive
ownership. Reject duplicate IDs within a relationship list and self-links.
Check dependency and membership cycles separately; do not combine the two edge
types into one precedence graph. A group may depend on delivery of its own
members without making membership an execution prerequisite for each child.
Cross-branch relationships remain outside this reader's scope.

### First commands

`list`, `show <id>`, and `check` read the selected checkout's live files,
including uncommitted records. Present the selected project path so the source
is clear. Those commands change no records, dates, configuration, or Git state.
`new <type> <title> [--slug SLUG]` allocates the next ID as specified above
(refusing a type the project's schema does not have before reserving anything),
writes `<id>-<slug>.md` with a body skeleton and equal `created`/`updated`
timestamps, prints the root-relative path, and fails without deleting the file
if the project no longer validates. It requires Git and never overwrites.
`show <id> --json` prints one object with `id`, `path`, `revision`, and
`source`. `update <id> --expect REVISION` with `--set FIELD=VALUE` and
`--unset FIELD` changes `title`, `status`, `relates_to`, work planning fields,
question `blocks`, plan and review `work`, or review `examined` by editing only those frontmatter entries plus `updated`;
[G-009](../grove/G-009-update-records.md) owns its request, preservation,
locking, and failure-reporting contract, and prints `{id, path, revision, changed}`.
`versions [ID] [--json]` reads the same project location on every local
branch tip and in every registered worktree, validating each source alone by
these rules, and prints one row or JSON object per observed version with a
selector; [G-010](../grove/G-010-record-versions.md) owns its source,
output, incomplete-result, and selector contract. `workspace --source
SELECTOR [--json]` re-inspects that selection and prints the project
directory of the existing checkout that still holds exactly that version;
[G-011](../grove/G-011-record-workspace.md) owns its resolution and
refusal contract. Both require Git and read only.

- `list`: show ID, type, status, and title, ordered by `created` ascending with
  undated records last, then ID prefix and numeric suffix as the tie-breaker.
  Do not infer urgency from that order.
- `show <id>`: show the file path and complete Markdown source, including
  frontmatter and relationships. The original bytes go to stdout; project and
  file context go to stderr. Missing or ambiguous identity is an error.
- `check`: report all discovered configuration, file, metadata, and relationship
  problems with paths and field names where available. Successful validation
  establishes structural consistency only.

Require nonempty string values for required fields, valid lifecycle values, and the
declared types for optional fields. Reject duplicate YAML keys and unknown
frontmatter/configuration keys in this initial schema so misspellings cannot
silently change behavior. Keep arbitrary supporting material in the Markdown
body; adding structured fields requires an explicit schema choice.

Frontmatter requires opening and closing `---` lines. UTF-8 BOM and CRLF files
are accepted and retained by `show`; invalid UTF-8 is rejected. Configuration and
frontmatter each hold one YAML mapping. Aliases, merge keys, custom tags, extra
YAML documents, null optional values, and scalar coercion are rejected. Normal
block/flow lists and string forms remain available. One-line list output escapes
control characters so a multiline title cannot alter the table structure.

`--project DIR` can appear before or after the command. `--help` works without a
project. Exit codes are 0 for success/help, 1 for inspection or output errors, and
2 for invalid usage. `check` reports the record count when validation succeeds.

Load the complete record set before resolving relationships, and report invalid
or missing targets within this checkout. Initial inspection commands use the
same validation boundary and return nonzero for an invalid project; `list` and
`show` must not silently present a partial valid subset as the whole project.
Concurrent direct edits can invalidate a read; this first reader does not
promise a transactional snapshot or ownership of the files.

### Why these defaults

Type folders support manual browsing; a configurable root accommodates projects
that keep knowledge under `docs/`. A flat directory remains possible future
design if three folders become friction. Shared sequential IDs make filenames
and references readable for the initial local-worktree audience. The earlier
random-ID/timestamp filename trial was rejected after actual browsing; a hidden
random identity with numbered aliases would add a second identity without being
needed for allocation within one repository. Reconsider allocation if independent
clones become a normal collaboration path. One project schema marker avoids repeating
versions in every file; mixed record versions and automatic migrations are
deferred until a real schema change needs them.

## Operational records

The earlier illustrative examples have been replaced by real records:

- [Work: inspect project records](../grove/G-003-inspect-records.md)
  owns the first CLI's completed outcome and verification evidence.
- [Question: branch versions](../grove/G-002-branch-versions.md)
  retains the accepted grouping and explicit-version-selection policy; source
  and routing implementation contracts remain in G-010 and G-011.
- [Decision: starter defaults](../grove/G-001-starter-defaults.md)
  records acceptance and points here for the schema rather than copying it.
- [Work: create records](../grove/G-007-create-records.md) owns the
  implemented creation command and shared allocation.
- [Decision: sequential IDs](../grove/G-004-sequential-ids.md)
  records the naming revision and starter-record ID migration.

Branch-context direction and its evidence remain in the restart brief.

## Lifecycle and validation boundary

- Work: `proposed`, `active`, `done`, `abandoned`.
- Question: `open`, `resolved`; retain the answer in its body or link to the
  durable decision instead of deleting the question's identity.
- Decision: `proposed`, `accepted`, `rejected`. Supersession can be added when
  an actual replacement needs it; prior versions remain available in Git.

A work item marked done asserts its intended outcome was achieved in that
record's branch context. It does not establish integration into main or the
truth of its evidence. Reopening changes status explicitly. Directory movement
does not determine completion. A resolved question stops blocking named work;
an abandoned prerequisite does not count as delivered.

Body organization is for readers. The first CLI should not infer readiness or
completion from exact headings, populated prose, or checked boxes. Validate
required metadata, supported values, ID uniqueness within the current checkout,
relationship targets/types, and work dependency cycles. Copies of one ID on
different branches are versions to reconcile, not automatically ID collisions.

## First dogfooding boundary

The CLI now lists, shows, and validates the records tracking its own development
in one checkout, creates records with shared IDs, updates their fields
while refusing stale writes, shows each record's versions across local
branches and worktrees, and locates the existing checkout holding a selected
version. The [integrated CLI review](../grove/G-022-integrated-cli-review.md) found
contract defects that G-014/G-015/G-016 repair; their
[evidence](../grove/G-028-repairs-review.md) lists the remaining limits.
[G-017](../grove/G-017-terminal-picker.md) adds a read-only terminal board over
these operations; it changes no schema and writes no record. Editing from an
interactive view and automatic checkout creation are future investments.

Plans and reviews are records from schema 2 on. Report records, artifact
ingestion, and agent attempts are deferred. Ordinary Markdown links and prose
can carry other supporting material in the meantime. The work planning metadata above is accepted for
the starting schema; attachment deferral does not require
deferring useful planning fields. Assignees and richer record types remain
future additions when the first workflow needs them.
