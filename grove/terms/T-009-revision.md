---
id: "T-009"
type: term
title: "Revision"
status: settled
created: "2026-09-21T05:01:57Z"
updated: "2026-09-21T14:23:35Z"
relates_to: ["T-004", "T-008"]
---

## Meaning

The exact content identity of one file: `sha256:` plus the hash of every byte,
as `show --json` and `context` print it. `update --expect` refuses to write
unless the file still has the revision the caller read, so a stale edit cannot
silently overwrite someone else's change. Dates in frontmatter are an
authoring convention and are never used for this.

Not a Git commit: a revision identifies a file's bytes wherever they are, a
commit identifies a whole tree in history. A review's `examined` field is a
commit, because a [candidate](T-004-candidate.md) is one.

## Relationships

Each version of a record in a [source](T-008-source.md) has a revision; equal
revisions mean identical files.
