# Starter record model

Status: core model accepted, 2026-09-18, on the owner's response "yup this looks
good" to the draft. This covers the representation, fields, relationships,
lifecycles, and initial validation boundary below. Random IDs are accepted as
a trial; their exact format and storage layout/versioning remain open.
No implementation exists yet.
[The restart brief](restart-brief.md) owns product direction; these examples
are not an additional operational backlog.

## Minimum representation

Use one Markdown file per record, with YAML frontmatter for facts that the
CLI interprets and a freeform Markdown body for explanation.

| Field | Purpose |
| --- | --- |
| `id` | Stable identity; renaming the title or file does not change it |
| `type` | `work`, `question`, or `decision` |
| `title` | A readable label for lists, search, and the eventual board |
| `status` | Explicit lifecycle state for that record type |

The accepted relationship fields cover the initial examples:

- Work `depends_on`: prerequisite work IDs. Missing means no declared prerequisites.
- Question `blocks`: work IDs whose outcome needs the answer. An unresolved
  question may be relevant without blocking work.
- Any record `relates_to`: related record IDs, with no implied ordering or gate.

## Proposed work metadata extension

On 2026-09-18, the owner proposed richer work metadata after comparing nullsec's
W-032 release, specifically members, dependencies, priority, size, and a work
sub-kind, to support useful UI behavior. The recommendation below extends the
accepted core; its exact fields and semantics have not yet been approved.

| Field | Proposed meaning and use |
| --- | --- |
| `kind` | What sort of work this is; filters and distinct presentation for feature, fix, refactor, investigation, tooling, or release |
| `priority` | Importance for selection; propose 1 highest through 5 lowest, independent of dependencies |
| `size` | Coarse scope/effort estimate; propose small, medium, or large, without automatic time estimates or execution policy |
| `members` | Child work included in this outcome; expandable groups and member completion counts |
| `depends_on` | Existing accepted field: prerequisite work that must be delivered before this work can proceed |
| `created`, `updated` | Creation and modification timestamps for chronology; potentially common to all record types |

Keep these optional so quick capture remains useful. An absent size or priority
means unspecified; do not silently turn it into an estimate or urgency decision.
Dates should be written by CLI mutations when those exist. Editing files directly
still needs an explicit timestamp policy; precise format is not settled here.

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
  Exact nesting, shared membership, and cross-branch resolution still need rules.

Evidence: `grove context --work W-032 --phase shape` in nullsec returned
`kind: release`, `size: large`, `priority: 2`, `depends_on: []`, eight ordered
members, creation/update dates, and parent-level acceptance. That context call
also reported unrelated supporting sources omitted by its token budget; the
owning work record itself was present and inspected. Its `scope` points to
capability records, a type outside the current three-type starter model.

The illustrative IDs below are readable placeholders. The owner accepts starting
with random IDs as a trial, while preferring some chronological ordering of files.
Generate identities independently across checkouts; specify encoding and length
before implementing creation. Links resolve by stable ID, not filename.
Chronological presentation can use creation metadata or a filename date prefix
independently of identity; neither is selected here. Revisit if random ordering
makes the actual file workflow awkward.
The first implementation also needs a schema-version marker; its location and
project configuration remain to be settled with storage layout.

Remaining layout proposals:

- Use root `grove.yaml` for the schema version and a configurable record folder,
  with a visible root-level `grove/` containing `work/`, `questions/`, and
  `decisions/` as the current recommendation. This makes browsing records by
  purpose straightforward. `docs/grove/` remains a possible configured location;
  nesting under `docs/` is not a product requirement.
- The earlier flat `docs/grove/records/` proposal minimized directory conventions,
  but had no stronger product justification. Prefer useful manual navigation
  over simplifying the directory traversal. Exact discovery and folder/type
  validation rules remain to be specified; identity stays in metadata.
- Use `<id>-<readable-slug>.md` filenames for navigation. Relationships use the
  stored ID, so a filename change does not invalidate them.

Layout and filename defaults remain recommendations, not selected decisions.

## Three examples

These use real topics from Grove's development to test the model. The decision
example represents direction already recorded in the brief; the work describes
a proposed outcome. Example IDs and filenames are illustrative.

### Work: an intended outcome

```markdown
---
id: work-cli-inspection
type: work
title: Inspect Grove project records from the CLI
status: proposed
depends_on: []
relates_to: [decision-branch-context]
---
List the work, questions, and decisions in the current checkout, inspect one
record, and report malformed metadata or unresolved relationships.

The first version reads local records. Cross-branch aggregation follows after
its selection behavior is settled.

Acceptance:
- The list identifies each record by ID, type, title, and status.
- Showing a record includes its text and relationships.
- Validation identifies the file and problem without changing any files.
```

### Question: an unresolved uncertainty

```markdown
---
id: question-branch-versions
type: question
title: How should the board present differing versions of one work item?
status: open
blocks: []
relates_to: [work-cli-inspection, decision-branch-context]
---
Main and a feature branch can contain different versions of the same record.
Which version is displayed, and how can the user inspect and open either one?

Exercise the case with a real worktree before building the combined board.
This does not block reading records from the current checkout.
```

### Decision: a choice and its rationale

```markdown
---
id: decision-branch-context
type: decision
title: Edit records in their selected branch workspace
status: accepted
---
Keep project records on code branches and provide a combined project view.
Open the selected record's existing worktree for editing; preserve the current
checkout. Prepare a worktree when needed through the workspace-opening action.

The owner accepted the branch-context restriction provided normal navigation
requires little thought. Revisit if this makes project management cumbersome.
```

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

Settle identity allocation and storage layout/versioning,
then use the agreed records to track the new CLI's own development. Begin with
listing, showing, and validating one checkout, followed by safe creation and
updates. Branch aggregation and workspace navigation build on that foundation.

Structured attachments, review/report records, artifact ingestion, and agent
attempts are deferred. Ordinary Markdown links and prose can carry supporting
material in the meantime. The work metadata extension above is now under
discussion for the starting schema; attachment deferral does not require
deferring useful planning fields. Assignees and richer record types remain
future additions when the first workflow needs them.
