# Record model

This is Grove's current record contract: configuration, record types, fields,
statuses, validation and the lifecycle rules software enforces. How each rule
was chosen is in the record linked beside it and in Git, not here; the
[brief](../grove/brief.md) owns product direction, and the
[README](../README.md) describes the commands.

## Identity and placement apart from classification

`schema_version: 3` is the one schema Grove reads
([G-064](../grove/G-064-stable-knowledge.md) selected it,
[G-065](../grove/G-065-flexible-records.md) implemented it). Grove keeps no
backward compatibility before its first release: schemas 1 and 2, their type
folders, typed `W-`/`Q-`/`D-`/`T-`/`P-`/`R-` IDs and per-type counters were
deleted by [G-052](../grove/G-052-migrate-knowledge.md), and
[G-069](../grove/G-069-migration-map.md) maps this repository's old IDs and
paths. A commit from before the conversion is inspected with the CLI in that
commit (`go run ./cmd/grove` there); the current CLI reports such a branch in
`versions` and the board as a source it cannot inspect.

- **Discovery.** Every `.md` file beneath the record root, at any depth, is a
  record, except the configured brief. No folder names a type: the root and
  any nested folder are equally valid. A `.md` file there that is
  not a valid record is a diagnostic, never skipped, dot folders included; only
  the lowercase `.md` extension counts. The brief may be any clean
  project-relative `.md` path.
  Symlinks are refused.
- **Identity.** An ID is `G-`, and a canonical number of at least
  three digits (`G-001`, `G-1000`). It carries no type: `new` issues `G-NNN`
  for every type, and a record keeps its ID whatever happens to its `type`.
  IDs are unique per checkout and matched exactly. Any other spelling,
  a typed `W-001` included, is an invalid ID.
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
- **Reclassification.** `update ID --set type=TYPE` changes
  classification in place. The result must satisfy the new type's whole
  contract in that one update, for example `--set type=work --set
  status=proposed` from a page, or `--set type=page --unset status` toward
  one, and the project must still validate, so work that a plan or question
  names cannot stop being work. Reclassifying grants nothing: an `accepted`
  decision is a claim in a file, as it always was.
- **`formerly`.** An optional string that only `convert` writes and `update`
  refuses: the typed ID or document path a record replaced. Two records with
  one `formerly` are an error. A merge from a branch that predates the
  conversion cannot quietly restore a second owner: a restored typed-ID record
  fails `check` on its ID.

**Allocation.** Numbers come from the one counter file `grove/neutral-ids` in
the Git common directory, under `grove/lock`, in the form `G 12`, with the
floor scan (record files beneath the record root on every branch,
remote-tracking ref and tag, and in every worktree, nested included) and recovery notices described under
[Identity and dates](#identity-and-dates). A `grove/next-ids` file left by the
deleted typed counters is never read or written.

**Conversion (`grove convert`).** Turns a Markdown document outside the record
root into a record, one source per run. Its other form, which gave a typed-ID
record a neutral ID, did G-052's conversion and was deleted with typed IDs:

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
  since a case-insensitive filesystem opens `docs/Plan.md` as `docs/plan.md`)
  and a missing source are refused before any
  ID is reserved, so a rerun neither duplicates a record nor remaps an
  identity. A refusal after the reservation, such as an existing target file,
  consumes the number and says so; gaps are acceptable, as for `new`.
- Not rewritten: body prose, Markdown links (the moved file's own relative
  links included), `examined`, and anything outside the record root. The
  caller repairs links from the mapping; `check` does not verify them.
- Everything is validated before the write, and an existing target file
  refuses the run untouched. The original document is never modified or
  removed.

Not provided: lookup of a record by its former ID, batch conversion, link
rewriting, and per-type folders or prefixes.

## Knowledge records and the brief

Terms, plans and reviews ([G-051](../grove/G-051-typed-knowledge-records.md),
implemented by [G-037](../grove/G-037-knowledge-artifacts.md)):

| Type | Statuses (first is what `new` writes) | Extra fields |
| --- | --- | --- |
| `term` | `proposed`, `settled` | none |
| `plan` | `current`, `superseded` | `work` |
| `review` | `current`, `superseded` | `work`, `examined` |

- A term's title is the term; its body gives meaning and relationships, not
  execution instructions. Two terms whose titles match, ignoring case and
  surrounding space, are an error.
- `work` is an optional list of work IDs the plan or review belongs to, checked
  like `depends_on` targets. One plan can name several work items. Work does
  not name its plans or reviews back: that side is derived, and
  `context G-NNN` lists them without reading them. `new` takes no fields, so
  set it with `update ID --set 'work=["G-001"]'`.
- `examined` is an optional quoted Git commit, 7 to 40 lowercase hex digits:
  what the review looked at. Whether the reviewed content has changed since is
  a comparison a reader makes, not stored state, against the work's
  `candidate` ([Work lifecycle](#work-lifecycle)). A review record is
  evidence, never approval: approval is the work record's `approved` field
  with the verdict appended to its body. There is no `report` type, since
  the work record's Evidence is the report.
- Optional `brief: PATH` in `grove.yaml` names the one project brief: a clean
  project-relative `.md` path without `..`, anywhere in the project, including
  directly under the record root (`grove/brief.md`). It is not a record and
  has no ID or frontmatter rules; that one path is exempt from discovery
  (compared without case). Every live command requires a
  regular, non-symlink file there and names `grove.yaml: brief` when it is
  missing. `grove brief [--json]` prints it like `show`; `context` never adds
  it, and `--include PATH` still can.

Committed sources in `versions` and the board check only the form of `brief`,
never that the file exists: they read `grove.yaml` and the record root, and
only a live checkout is required to hold the brief.

## Minimum representation

Use one Markdown file per record, with YAML frontmatter for facts that the
CLI interprets and a freeform Markdown body for explanation.

| Field | Purpose |
| --- | --- |
| `id` | Stable identity; renaming the title or file does not change it |
| `type` | `work`, `question`, `decision`, `term`, `plan`, `review`, or `page` |
| `title` | A readable label for lists, search, and the board |
| `status` | Explicit lifecycle state for that record type; a `page` has none |

Relationship fields:

- Work `depends_on`: prerequisite work IDs. Missing means no declared prerequisites.
- Question `blocks`: work IDs whose outcome needs the answer. An unresolved
  question may be relevant without blocking work.
- Any record `relates_to`: related record IDs, with no implied ordering or gate.
- Plan and review `work`: the work they belong to.

## Work planning metadata

Optional work fields for filtering, grouping and presentation; no automatic
execution policy is implied. Their allowed values are under
[Planning values](#planning-values).

| Field | Meaning and use |
| --- | --- |
| `kind` | What sort of work this is, for filtering and distinct presentation |
| `priority` | Importance for selection, independent of dependencies |
| `size` | Coarse scope/effort estimate, without automatic time estimates or execution policy |
| `members` | Child work included in this outcome; expandable groups and member completion counts |
| `depends_on` | Prerequisite work that must be delivered before this work can proceed |
| `created`, `updated` | Creation and modification timestamps for chronology, allowed on every type |

Keep these optional so quick capture remains useful. An absent size or priority
means unspecified; do not silently turn it into an estimate or urgency decision.

These distinctions hold:

- Membership describes decomposition. Member order can express presentation or
  preferred sequence, but only dependencies impose prerequisite ordering.
- Allow grouped work to retain its own outcome and acceptance. Completed children
  do not automatically establish parent completion or integration.
- Incomplete members can prevent declaring a parent done without blocking work
  on that parent. Do not conflate an unfinished group with an execution blocker.
- Derive member counts and blocker explanations from relationships and the
  selected branch context. Do not store parallel progress percentages or an
  independently editable `blocked` flag.
- A spike/investigation is a kind, not a size. A size sets no preparation
  depth or execution rule.
- Member targets and membership cycles are validated as well as dependency
  cycles. Nesting and shared membership follow the rules below; relationships
  resolve within one checkout.

## On-disk contract

[G-001](../grove/G-001-starter-defaults.md) and
[G-004](../grove/G-004-sequential-ids.md) record how these defaults were
accepted. [G-003](../grove/G-003-inspect-records.md) implements discovery,
validation and inspection, [G-007](../grove/G-007-create-records.md)
creation and allocation per [G-006](../grove/G-006-allocator-mechanism.md),
and [G-009](../grove/G-009-update-records.md) field updates with content
revisions and a shared write lock. No command renames, moves (other than
`convert`), deletes or edits the body of a record.

### Configuration and discovery

Keep one `grove.yaml` at the project root:

```yaml
schema_version: 3
records: grove
brief: grove/brief.md
```

`grove init` writes exactly this file when none exists, creates the record
root and a placeholder brief, and keeps an existing configuration that
validates, following its own `records` and `brief`.

Optional `target: BRANCH` names the integration target: the local branch that
work is merged into, `main` in this repository. `init` does not write it.
[G-042](../grove/G-042-current-view.md) uses it only to label the current view
(on the target, or not), never to decide which state is current. It is
compared with branch names and never passed to Git; a value with surrounding
spaces or a `refs/` prefix is refused as a likely mistake. Every valid source
whose `grove.yaml` names a target must agree; a source that names none has no
say. Conflicting names, or a branch that is missing or unreadable, leave no
target and give a note. A target branch with no project yet, as while Grove
is adopted on another branch, lacks every record.

`schema_version` versions the configuration and record schema together. Require
both keys, with `brief` and `target` optional; accept exactly 3 and refuse a missing or
other version without guessing, migrating, or rewriting files. The number is
this CLI's, unrelated to the sibling skills CLI's schema numbering or
`grove.toml` configuration.

Resolve `records` relative to the directory containing `grove.yaml`. Permit a
different relative folder, such as `docs/grove`; require it to remain inside the
project and reject absolute paths, `..` components, and symlink traversal. The
record root must be a dedicated subdirectory, not the project root itself.

Without an explicit project, search from the current directory upward for the
nearest `grove.yaml`, stopping at the current Git checkout root when inside one.
Outside Git, search up to the filesystem root. An explicit `--project <dir>`
selects the directory containing the configuration without upward searching.
Git is not required to read a configured project. `list`, `show` and `check`
never consult another worktree, the Git common directory, or another tool's
configuration.

### Folders and files

```text
grove.yaml
grove/
  brief.md
  G-001-starter-defaults.md
  G-003-inspect-records.md
```

Read `.md` files recursively beneath the record root. Nested folders may
organize records but confer no type or lifecycle meaning; `new` writes flat.
Closed records stay discoverable in the same tree. Reject symlinks in the
record tree; report a `.md` file that is not a valid record rather than
silently dropping it (the configured brief excepted). Other file extensions
are not records. A missing record root is an error.

Generate short filenames as `<id>-<slug>.md`.

Keep dates in frontmatter. The number supplies allocation order; it
does not prove creation time, priority, or execution order.
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

Use neutral sequential IDs, `G-001` for every type, from one counter. The
full ID is the canonical identity, not an alias for a hidden random value.
Start at 1, pad to a minimum of three digits, and expand beyond 999
(`G-1000`) without wrapping or renumbering older records. Require canonical
padding. Store IDs as strings and match references exactly; `show G-001`
needs no abbreviated-ID lookup. Numeric ordering must not rely on
lexicographic sorting once the counter expands.

Creation coordinates through the directory that
`git rev-parse --path-format=absolute --git-common-dir` returns, never a
worktree's own `.git` path, since linked worktrees have private metadata as
well as a shared common directory
([Git's worktree documentation](https://git-scm.com/docs/git-worktree#_details)).
Allocation in one local repository:

1. Take `grove/lock` in the common directory, one lock for every worktree and
   every record type.
2. Reserve the next number and durably save the advanced counter while holding
   the lock. If locking or persistence fails, no ID is issued and no record is
   created.
3. Release the lock, then create the record in its checkout without
   overwriting an existing file. A failed or abandoned creation consumes the
   reservation; gaps are acceptable and numbers are not recycled.

The counter, `grove/neutral-ids` in the form `G 12` (the next number), belongs
to the local repository, not to a branch, worktree, or record-root path. It is
local coordination state, not a tracked record, and needs no daemon;
[G-006](../grove/G-006-allocator-mechanism.md) owns its encoding, lock
primitive and recovery protocol. `new` and `update` serialize publication
through `grove/write.lock` beside it. The read-only commands never create any
of these files, and none is ever unlinked.

Separate clones do not share reservations. Directly authored IDs also bypass
allocation. Imported records and independently allocated clone histories require
collision checks and explicit reconciliation; matching IDs alone cannot prove
two independently created records are the same item. Two files with one ID in
one checkout remain an error even when their contents match. Genuine branch
copies of one record retain their identity.

Counter state never silently restarts at 1. Each allocation floors the
counter by the highest ID in use: an `id:` line in any text file beneath the
record root on every branch, remote-tracking ref and tag, and every live
record in every worktree, nested folders included. A missing counter is
initialized from that floor with a notice on stderr that reservations for
records never written or since deleted cannot be recovered; a counter below
the floor continues above it, with a notice. A corrupt counter refuses
allocation until it is fixed or removed. Concurrent imports and manual edits
are outside the allocator's exclusivity guarantee. Reading needs no Git or
allocator state; allocation outside Git is not provided.

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
and `update --expect REVISION` refuses any other current content. The flag is
optional: an agent session passes it because its read may be old, while a
person at a shell, reading and writing within seconds under the same write
lock, omits it and lets the update apply to the file as it is.

### Planning values

| Field | Allowed values | When absent |
| --- | --- | --- |
| `kind` | `feature`, `fix`, `refactor`, `investigation`, `tooling`, `release` | Unspecified |
| `priority` | Integer 1 (highest) through 5 (lowest) | Unspecified |
| `size` | `small`, `medium`, `large` | Unspecified |
| `members`, `depends_on` | Ordered lists of work IDs | No declared relationships |

The planning fields above belong only on work. `blocks` belongs only on
questions; `relates_to`, `created`, and `updated` can appear on every type.
Do not synthesize a default kind, priority, or size. An investigation uses
`kind: investigation` and an independent size; there is no separate `spike`
value.

Allow nested and shared membership; grouping does not assign exclusive
ownership. Reject duplicate IDs within a relationship list and self-links.
Check dependency and membership cycles separately; do not combine the two edge
types into one precedence graph. A group may depend on delivery of its own
members without making membership an execution prerequisite for each child.
Relationships resolve within one checkout.

### Reading and writing records

`list`, `show <id>`, and `check` read the selected checkout's live files,
including uncommitted records. Present the selected project path so the source
is clear. Those commands change no records, dates, configuration, or Git state.
`new <type> <title> [--slug SLUG]` allocates the next ID as specified above
(refusing an unknown type before reserving anything),
writes `<id>-<slug>.md` with a body skeleton and equal `created`/`updated`
timestamps, prints the root-relative path, and fails without deleting the file
if the project no longer validates. It requires Git and never overwrites.
`show <id> --json` prints one object with `id`, `path`, `revision`, and
`source`. `update <id> [--expect REVISION] [--commit]` with `--set FIELD=VALUE` and
`--unset FIELD` changes `title`, `status`, `relates_to`, work planning fields,
`candidate` and `approved`, question `blocks`, plan and review `work`, or review `examined` by editing only those frontmatter entries plus `updated`;
[G-009](../grove/G-009-update-records.md) owns its request, preservation,
locking, and failure-reporting contract, and prints `{id, path, revision, changed}`.
[G-079](../grove/G-079-update-a-record-by-hand-without.md) made `--expect`
optional and added `--commit`, which after a change runs `git add` and
`git commit` for the record's file alone with a generated message, adds
`commit` to the result (`null` when nothing changed), and reports a commit Git
refused as an applied, uncommitted update.
`versions [ID] [--json]` reads the same project location on every local
branch tip and in every registered worktree, validating each source alone by
these rules, and prints one row or JSON object per observed version with a
selector; [G-010](../grove/G-010-record-versions.md) owns its source,
output, incomplete-result, and selector contract, and
[G-042](../grove/G-042-current-view.md) marks each version current or older
by Git ancestry, which the board's default view shows. `workspace --source
SELECTOR [--json]` re-inspects that selection and prints the project
directory of the existing checkout that still holds exactly that version;
[G-011](../grove/G-011-record-workspace.md) owns its resolution and
refusal contract. Both require Git and read only.

- `list`: show ID, type, status, and title, ordered by `created` ascending with
  undated records last, then the ID's number as the tie-breaker.
  Do not infer urgency from that order. `--status VALUE`, repeatable and only
  on `list`, keeps the records whose status equals any given value; a value
  outside the union of the type table's status vocabularies, or an empty one, is a
  usage error (exit 2), and a status no record holds prints the header alone.
  Without it, every record is printed.
- `show <id>`: show the file path and complete Markdown source, including
  frontmatter and relationships. The original bytes go to stdout; project and
  file context go to stderr. Missing or ambiguous identity is an error.
- `check`: report all discovered configuration, file, metadata, and relationship
  problems with paths and field names where available. Successful validation
  establishes structural consistency only.

Require nonempty string values for required fields, valid lifecycle values, and the
declared types for optional fields. Reject duplicate YAML keys and unknown
frontmatter/configuration keys so misspellings cannot
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
or missing targets within this checkout. The inspection commands use the
same validation boundary and return nonzero for an invalid project; `list` and
`show` never present a partial valid subset as the whole project.
Concurrent direct edits can invalidate a read; the reader does not promise a
transactional snapshot or ownership of the files.

## Lifecycle and validation boundary

- Work: `proposed`, `active`, `review`, `done`, `abandoned`, in that order;
  the rules are under [Work lifecycle](#work-lifecycle).
- Question: `open`, `resolved`; retain the answer in its body or link to the
  durable decision instead of deleting the question's identity.
- Decision: `proposed`, `accepted`, `rejected`, `superseded`. `superseded`
  is an accepted decision that a later one replaced: `relates_to` names the
  replacement and the body says why, with no dedicated field
  ([G-041](../grove/G-041-nullsec-pilot.md)). Prior versions remain available
  in Git.

Reopening changes status explicitly. Directory movement does not determine
completion. A resolved question stops blocking named work; an abandoned
prerequisite does not count as delivered.

### Work lifecycle

Proposed → Active → Review → Done, with Abandoned only through an explicit
human decision ([G-035](../grove/G-035-interactive-adoption.md) selected it,
[G-038](../grove/G-038-review-lifecycle.md) implemented it, and
[G-044](../grove/G-044-review-integration.md) added `approved`, `approve`,
`feedback` and `integrate`). Preparation,
implementation, independent review and waiting are activities inside
`active`, recorded in the body, never statuses. The settled terms
[candidate](../grove/G-057-candidate.md), [review](../grove/G-058-review.md),
[approval](../grove/G-059-approval.md) and
[integration](../grove/G-060-integration.md) name the facts.

- **`candidate`** is an optional work field: a quoted Git commit, the same
  form as a review's `examined`, naming the commit offered for judgment
  together with the evidence gathered at it. It is required while the status
  is `review` and allowed on every other status. It names the last
  implementation commit; the commit that sets `review` changes only the
  record, so `git diff --stat CANDIDATE TIP` shows one file. A changed
  candidate is a new value set through `update`, so the reviewed and the
  approved commit can be compared to each review's `examined`; the earlier
  value stays in Git history.
- **Review** means a candidate awaits human judgment. The record's Evidence
  and Next carry the handoff the work guide describes, so a new session can
  judge it without the originating chat. A failed or interrupted attempt does
  not enter Review: it stays `active` with a checkpoint.
- **Approval** is the owner's verdict on one candidate. `approve ID VERDICT`,
  in a clean checkout of the branch that holds the record in review, sets
  the optional work field **`approved`**, a quoted commit that must equal
  `candidate` and is valid only while the status is `review` or `done`, and
  appends `Verdict on candidate X, DATE: VERDICT` as the body's last
  paragraph, committed alone. A changed candidate cannot inherit it: `check`
  rejects an `approved` that differs from `candidate`, and `approve` refuses
  a tip that changed any file but the record after the candidate. Feedback
  that asks for more implementation, `feedback ID TEXT` in the same checkout,
  sets `active`, unsets `approved`, keeps `candidate` so earlier reviews
  still compare to it, and appends `Feedback on candidate X, DATE: TEXT`;
  nothing earlier is removed.
- **Done** means the candidate was accepted and merged into the target, for
  research and design deliverables too, since those are files. `update` writes
  `done`, or changes a done record's candidate, only when the candidate is an
  ancestor of the checkout's `HEAD` (`git merge-base --is-ancestor`), so a
  checkout that lacks the code cannot close the work; a record that stays
  done cannot lose its candidate. With a configured `target`, `update` also
  refuses `done` in a checkout on any other branch. `integrate ID` does the
  integration: in a clean checkout of the target it finds the one branch
  holding an approved candidate of ID, merges it with a plain `git merge`
  (a conflict is aborted and refused before anything changes), writes done
  there committed alone, and with `--cleanup` removes the branch's worktree
  and the branch only where Git agrees and the worktree holds no ignored
  files. The merge is of the commit the checks read, so a branch that moves
  meanwhile is not merged. A squash or rebase that lands a
  different commit is a manual merge that names that commit as the candidate
  in the same `update`. The check needs Git, as `update` already does;
  `check` verifies the form only.
- **Historical Done.** A `done` work record without `candidate` was completed
  before this rule and asserts only that its outcome was achieved in that
  record's own branch context, as its Evidence says; it is not proof of a
  merge. Nothing rewrites it, `check` accepts it, and its other fields stay
  editable; `update` never writes a new one. Delivery of such a prerequisite
  is established by Git ancestry or observed behavior, as before.
- Not enforced by software: the order of transitions; that Abandoned needs a
  human decision; that a review record exists before Review; that done is
  written on the target where `grove.yaml` names none; and that a reopened
  record's candidate is moved to its new commits before it is closed again,
  which the guide's `git diff --stat CANDIDATE TIP` check and `approve`'s
  refusal catch. These are guide rules, since software cannot verify a
  person, and the owner edits by hand. The merge is local, and every action
  is a person's command or key.

Body organization is for readers. The CLI does not infer readiness or
completion from exact headings, populated prose, or checked boxes. Validate
required metadata, supported values, ID uniqueness within the current checkout,
relationship targets/types, and work dependency cycles. Copies of one ID on
different branches are versions to reconcile, not automatically ID collisions.

## Not records

Attempts are not records: `grove run`
([G-045](../grove/G-045-durable-attempt.md)) keeps each attempt's inputs, raw
events and result as files under the Git common directory, shared by every
worktree and never committed, and the work record's own status on the
attempt's branch is the only handoff. The board
([G-046](../grove/G-046-managed-runs.md)) reads the same files and derives an
outcome from them for display; it writes nothing about an attempt. There is
no report type, and no assignee field; ordinary Markdown links and prose carry
other supporting material.
