---
id: "W-019"
type: work
title: "Represent domain terms and linked work artifacts"
status: done
created: "2026-09-21T00:54:14Z"
updated: "2026-09-21T14:23:43Z"
kind: feature
size: medium
priority: 1
relates_to: ["D-004", "D-005", "W-018", "W-011", "W-029"]
---

## Outcome

Give domain vocabulary and work artifacts durable, discoverable homes so agents
and the TUI can retrieve the relevant knowledge or evidence without loading the
whole project. Selected direction is in [D-004](../decisions/D-004-interactive-adoption.md).

## Scope and bounds

Selected by the owner on 2026-09-20,
[D-005](../decisions/D-005-typed-knowledge-records.md): terms, plans and reviews
are typed records with sequential IDs in their own folders under the record
root (`terms/` `T-`, `plans/` `P-`, `reviews/` `R-`); a plan or review names its
work items in a `work` list, and work's attachments are derived from that. No
`artifacts` folder, `artifact` type or `kind` field. A `report` type waits for
W-020, its first consumer.

In: those three types through the existing reader, `new`, `update`, `show`,
`versions` and `check`; a discoverable brief location, valid both beneath the
record root and where a project's brief is today, separately identifiable from
typed records, one per project; listing a work item's plans and reviews in
staged `context`; the starting Grove terms. Terms describe domain meaning and
relationships, not execution instructions. Start with work, preparation,
attempt, review, approval, integration, source and revision; define terms
during real modeling, not to satisfy a count.

Out: moving this repository's brief, plans or reviews
([W-029](W-029-migrate-knowledge.md) owns that; do not move files before
support ships); review approval, candidate and disposition semantics (W-020
extends the review type); standalone plan tickets with their own lifecycle;
other Bench types; semantic search; wikilink parsing.

Observed at main `42c077d`: the loader refuses any `schema_version` but 1, any
unknown configuration key, and any `.md` outside the three type folders
(`internal/project/project.go:45-52,99`). A type is one table row of prefix,
folder, initial status and skeleton (`internal/create/create.go:26-31`), a
folder map entry (`project.go:67`) and the ID pattern `^[WQD]-`
(`internal/project/metadata.go:162-165`). `context` already lists every body
link unopened and `--include PATH` adds a file with its revision
(`internal/handoff/sources.go:91-138`), so shared plans work today through
plain links; what is missing is telling a plan from any other link, the reverse
lookup, and the revision a review examined. `check` does not verify link
targets. nullsec, the pilot target, holds six slug-ID terms and a brief under
`docs/grove/`; their ID mapping belongs to W-023.

Proposed, for preparation to settle: statuses such as term
`proposed`/`settled` and plan `current`/`superseded`; an optional revision
field on reviews naming what was examined; a `brief:` key in `grove.yaml`;
whether new types and the key need `schema_version: 2`.

## Acceptance

1. `new`, `show`, `update`, `versions`, `check` and the board's loader accept
   term, plan and review records under the existing ID, folder and revision
   rules. Unknown types, fields and values are still refused, and a `work`
   entry naming a missing or non-work record is an error naming the file.
2. `context W-NNN` lists the plans and reviews whose `work` names it, including
   one shared by several work items, without including their bodies; the caller
   can deliberately include one. Default context stays bounded and read-only.
3. A review can record the revision it examined, so a later change to the
   reviewed content is distinguishable from what was reviewed; no per-work
   duplicate of a shared plan is needed.
4. The configured brief is discoverable from the CLI, validates beneath the
   record root and at this repository's current `docs/restart-brief.md`, is
   never a typed record or added to context by default, and a missing target
   is a clear error.
5. Compatibility is explicit: an existing schema-1 project stays readable or
   gets a deliberate, documented migration; no file is silently rewritten.
6. The starting Grove terms exist as records, written from real modeling of
   this repository's vocabulary, and the owner judges them accurate.
7. The record model, CLI help, both guides and `AGENTS.md` describe the
   supported behavior, including the new locations, before anything uses them.

## Evidence

Branch `worktree-W-019` from main `d9fc2a5`; implementation through `fc9bef1`.
[Plan](../../docs/plans/W-019-knowledge-artifacts.md);
[R-001](../reviews/R-001-w-019-knowledge-records.md) is the independent
review, with its findings and their fixes. Verified at `fc9bef1`, uncached:
`gofmt -l .` and `go vet ./...` clean, `go test -count=1 ./...` ok, and
`go test -race -count=1 -p 1` ok (the reviewer's run; in the implementer's
run `internal/versions` flaked as R-001 describes and passed alone).
`grove check` passes.

1. `internal/cli/knowledge_test.go` creates, updates, lists and validates the
   three types; `internal/project/project_test.go` covers schema-1 refusal and
   each bad field with file and field; `versions` and `workspace` read a
   committed schema-2 tree.
2. `context W-019` in this checkout lists `R-001  review  current  listed
   review of W-019`; a plan shared by two work items is listed for both.
3. `examined` on R-001 names `fc9bef1`.
4. `grove brief` prints `docs/restart-brief.md`; `TestBrief` covers both
   locations, a missing file, and refused paths.
5. Schema 1 is unchanged against a binary built from the base (R-001).
6. Nine terms, `T-001` to `T-009`: the owner read them on 2026-09-21, agreed
   with them for now, and they are `settled`.
7. Record model, README, `AGENTS.md` and both guides reconciled in `4bfb5ff`.

## Preparation and next

Done on the owner's verdict of 2026-09-21: the terms (Candidate included) and
the new commands were accepted, and the owner asked for the merge into `main`.
Nothing remains here. W-020 depends on this; W-029 follows it.
