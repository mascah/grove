---
id: "W-003"
type: work
title: "Update record status and fields from the CLI"
status: proposed
kind: feature
priority: 3
size: medium
members: []
depends_on: []
relates_to: ["W-002", "D-003"]
created: "2026-09-19T15:36:19Z"
updated: "2026-09-19T15:36:19Z"
---

## Outcome

Change a record's status and planning fields from the CLI so agents and the
future workspace can move work without hand-editing frontmatter. Grove's own
records are the first data.

## Why now

Created by `grove new` as W-002's first real record. Priority 3: creation and
inspection are enough to dogfood; this becomes urgent when a UI or run needs
to change state. Medium reflects revision checks and frontmatter rewriting
that preserve the human-authored body byte for byte.

## Constraints

- Preserve `created`; set `updated` only when content changes. Compare actual
  content, not timestamps, before writing; refuse when the file changed
  underneath the command.
- Rewrite only the changed frontmatter fields; leave the body and unrelated
  fields untouched. Validate the whole project after the write.
- Respect the lifecycle values and planning vocabulary in the record model.
  Defer renames, moves, deletes, and cross-branch edits.

## Acceptance

- Set status and each optional work field, including clearing one, with the
  file otherwise identical.
- A concurrent edit between read and write is detected and refused.
- Fixture tests cover valid transitions, invalid values, unchanged no-ops that
  leave `updated` alone, and unchanged file hashes on failure.

## Next

Shape after W-002 closes: decide whether the revision check compares a content
hash supplied by the caller or re-reads and diffs, and whether status changes
need a reason recorded in the body.
