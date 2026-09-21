---
id: "G-061"
type: term
title: "Source"
status: settled
created: "2026-09-21T05:01:57Z"
updated: "2026-09-21T14:23:34Z"
relates_to: ["G-062"]
formerly: "T-008"
---

## Meaning

One place where a version of the project's records was observed: either the
committed tip of a local branch, or the live files of a registered worktree,
uncommitted edits included. `versions` prints one line per source and a
selector for each record version in it. A source is judged alone, by the same
rules as any checkout, and an invalid one is reported rather than hidden.

Not source code, and not the `source` field of `show --json`, which is a
record file's exact text.

## Relationships

A record has one version per source that holds it. `workspace` resolves a
selected version to the existing checkout of its source. Each version has a
[revision](G-062-revision.md).
