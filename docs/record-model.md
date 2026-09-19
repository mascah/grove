# Starter record model

Status: core model accepted, 2026-09-18, on the owner's response "yup this looks
good" to the draft. This covers the representation, fields, relationships,
lifecycles, and initial validation boundary below. Identity allocation and
storage layout/versioning remain open. No implementation exists yet.
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

Three optional relationship fields are sufficient for the initial examples:

- Work `depends_on`: prerequisite work IDs. Missing means no declared prerequisites.
- Question `blocks`: work IDs whose outcome needs the answer. An unresolved
  question may be relevant without blocking work.
- Any record `relates_to`: related record IDs, with no implied ordering or gate.

The illustrative IDs below are readable placeholders, not a selected allocation
scheme. Choose an allocation method that handles independent branch creation
before implementing record creation. Links resolve by stable ID, not filename.
The first implementation also needs a schema-version marker; its location and
project configuration remain to be settled with storage layout.

Proposed defaults for those remaining choices:

- Generate type-prefixed random IDs independently in each checkout, preserving
  them through renames and branch creation. Avoid a shared sequential counter.
  The encoding/length and abbreviated lookup behavior still need specification.
- Use root `grove.yaml` for the schema version and a configurable record folder,
  defaulting to `docs/grove/records/`. Keep all three types in that folder; type
  is metadata rather than inferred from the directory name.
- Use `<id>-<readable-slug>.md` filenames for navigation. Relationships use the
  stored ID, so a filename change does not invalidate them.

These defaults are recommendations, not part of the earlier core-model approval.

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
material in the meantime. Assignees, dates, priorities, and richer record types
are future additions when the first workflow needs them; their absence here
does not remove them from the broader product direction.
