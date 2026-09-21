---
id: "G-003"
type: work
title: Inspect Grove project records from the CLI
status: done
kind: feature
priority: 2
size: medium
members: []
depends_on: []
relates_to: ["G-001", "G-004", "G-002"]
created: "2026-09-19T14:08:40Z"
updated: "2026-09-19T15:12:15Z"
formerly: "W-001"
---

## Outcome

Use a Go CLI to list work, questions, and decisions in the current project,
inspect a record by ID, and diagnose invalid records without modifying files.
Grove's own development records are the first dogfooding data.

## Why now

The file defaults are accepted and real records exist. Reading and validating
them makes the format useful and exposes problems before adding mutations,
cross-branch aggregation, or agent execution. Priority 2 reflects this being
the next investment; medium reflects configuration, parsing, graph validation,
and command integration, without prescribing an execution policy.

## Constraints

- Follow the [record model](../docs/record-model.md), using its accepted
  defaults and resolving the supporting reader proposals during preparation.
- Read ordinary local files without a service, index database, or required Git
  repository. Inspect one selected checkout's working files.
- Keep IDs independent of filenames and chronology. Do not alter records,
  configuration, timestamps, branches, or the index while inspecting them.
- Read sequential `W-`, `Q-`, and `D-` IDs without needing allocator state in
  Git metadata. Shared ID allocation belongs to future creation commands.
- Keep the installed predecessor CLI intact. Run/build this implementation from
  its own source; do not install a replacement into the user's command path.
- Defer creation/update commands, TUI, cross-branch selection, claims, runners,
  structured attachments, and migration of sibling projects.

## Acceptance

- `list` identifies each valid record by full ID, type, title, and status.
- `show <id>` identifies the file and shows its complete Markdown source,
  including relationships. Missing or duplicate identity produces a diagnostic.
- `check` identifies malformed metadata, duplicate IDs, invalid relationship
  targets/types, and dependency/membership cycles with useful file diagnostics.
- The project path makes the source context clear. A configured nondefault
  record root works as well as `grove/`.
- Inspection leaves all project file contents unchanged, including for invalid
  projects. Structural validation makes no claim that acceptance is proven.
- The CLI successfully inspects Grove's operational records; focused
  fixture tests demonstrate failure behavior and preserve direct-editor use.

## Preparation

Execution: single agent; configuration, loading, validation, then commands.
Reason: the command layer depends on a shared reader and validation contract.
Delegation: none.
Runtime: interactive session; no runner is required.
Reassess: split work only if preparation reveals independently verifiable scope.

The [implementation plan](G-005-inspection-plan.md) adopts the reader
proposals for this work, specifies the remaining parsing/output choices, and
tracks implementation and verification. The owner authorized implementation
with "lets go" after the sequential-ID revision.

Verify discovery in a plain directory and a linked-worktree fixture, configured
roots, renamed records, missing optional fields, malformed frontmatter, duplicate
IDs, unresolved/wrong-type references, and distinct graph-cycle checks. Include
an integration check that hashes fixture contents before and after commands.
Test the reader as software once it exists; these records alone do not justify
creating an application test suite now.

## Evidence

Closed 2026-09-19. The [implementation plan](G-005-inspection-plan.md)
records the verification: passing unit, race, vet, and format checks; `list`,
`show`, and `check` against these records with unchanged file hashes; and a
linked-worktree fixture covering live edits, an explicit project override, a
malformed neighbor, and nested-directory discovery. Structural validation
passed; this does not prove acceptance of any record's content.

## Next

Shape safe record creation and shared sequential-ID allocation, per the
[restart brief](brief.md#suggested-sequence-and-current-next-action).
