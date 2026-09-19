# Starter record model

Status: core model accepted, 2026-09-18, on the owner's response "yup this looks
good" to the draft. This covers the representation, fields, relationships,
lifecycles, and initial validation boundary below. The owner accepted the
[on-disk defaults below](#on-disk-contract) on 2026-09-19, then replaced the
random-ID/long-filename trial with shared sequential IDs and short filenames.
The first Go CLI now reads and validates this model; W-001 records its evidence.
[The restart brief](restart-brief.md) owns product direction;
`grove/` owns operational work, questions, and decision receipts.

## Minimum representation

Use one Markdown file per record, with YAML frontmatter for facts that the
CLI interprets and a freeform Markdown body for explanation.

| Field | Purpose |
| --- | --- |
| `id` | Stable identity; renaming the title or file does not change it |
| `type` | `work`, `question`, or `decision` |
| `title` | A readable label for lists, search, and the eventual board |
| `status` | Explicit lifecycle state for that record type |

The accepted relationship fields cover the starter records:

- Work `depends_on`: prerequisite work IDs. Missing means no declared prerequisites.
- Question `blocks`: work IDs whose outcome needs the answer. An unresolved
  question may be relevant without blocking work.
- Any record `relates_to`: related record IDs, with no implied ordering or gate.

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
These revised defaults govern the operational records. W-001 implements the
discovery, validation, graph, and inspection behavior below; W-002 implements
creation and allocation per [D-003](../grove/decisions/D-003-allocator-mechanism.md).
[W-003](../grove/work/W-003-update-records.md) implements field updates with
content revisions and a shared write lock; renames, moves, deletes, and body
edits remain unimplemented. This section owns the schema; the brief owns direction.

### Configuration and discovery

Keep one `grove.yaml` at the project root:

```yaml
schema_version: 1
records: grove
```

`schema_version` versions the configuration and record schema together. Require
both keys; refuse missing or unsupported versions without guessing, migrating,
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

Read `.md` files recursively within those three type folders. Require frontmatter
`type` to match its folder (`work`, `question`, or `decision`); nested folders
may organize records but do not confer lifecycle meaning. Closed records stay
discoverable in the same tree. Reject symlinks in the record tree; report `.md`
files outside the three type folders as misplaced rather than silently dropping
them. Other file extensions are not records. A missing record root is an error;
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

Use type-prefixed sequential IDs: `W-001` for work, `Q-001` for questions, and
`D-001` for decisions. Each type has its own counter; the full prefixed ID is the
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
the Git common directory is the intended home. [D-003](../grove/decisions/D-003-allocator-mechanism.md)
owns the accepted state encoding, lock primitive, and recovery protocol for
[W-002](../grove/work/W-002-create-records.md) to test with the creation command. `grove new` initializes and maintains `grove/next-ids` under `grove/lock` in
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
is retained in [D-002](../grove/decisions/D-002-sequential-ids.md); this is not a
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
`new <type> <title> [--slug SLUG]` allocates the next ID as specified above,
writes `<id>-<slug>.md` with a body skeleton and equal `created`/`updated`
timestamps, prints the root-relative path, and fails without deleting the file
if the project no longer validates. It requires Git and never overwrites.
`show <id> --json` prints one object with `id`, `path`, `revision`, and
`source`. `update <id> --expect REVISION` with `--set FIELD=VALUE` and
`--unset FIELD` changes `title`, `status`, `relates_to`, work planning fields,
or question `blocks` by editing only those frontmatter entries plus `updated`;
[W-003](../grove/work/W-003-update-records.md) owns its request, preservation,
locking, and failure-reporting contract, and prints `{id, path, revision, changed}`.
`versions [ID] [--json]` reads the same project location on every local
branch tip and in every registered worktree, validating each source alone by
these rules, and prints one row or JSON object per observed version with a
selector; [W-004](../grove/work/W-004-record-versions.md) owns its source,
output, incomplete-result, and selector contract. It requires Git and reads only.

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

- [Work: inspect project records](../grove/work/W-001-inspect-records.md)
  owns the first CLI's completed outcome and verification evidence.
- [Question: branch versions](../grove/questions/Q-001-branch-versions.md)
  retains the accepted grouping and explicit-version-selection policy; source
  and routing implementation contracts remain in W-004 and W-005.
- [Decision: starter defaults](../grove/decisions/D-001-starter-defaults.md)
  records acceptance and points here for the schema rather than copying it.
- [Work: create records](../grove/work/W-002-create-records.md) owns the
  proposed creation command and shared allocation.
- [Decision: sequential IDs](../grove/decisions/D-002-sequential-ids.md)
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
while refusing stale writes, and shows each record's versions across local
branches and worktrees. Workspace navigation builds on that foundation.

Structured attachments, review/report records, artifact ingestion, and agent
attempts are deferred. Ordinary Markdown links and prose can carry supporting
material in the meantime. The work planning metadata above is accepted for
the starting schema; attachment deferral does not require
deferring useful planning fields. Assignees and richer record types remain
future additions when the first workflow needs them.
