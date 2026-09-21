---
id: "W-019"
type: work
title: "Represent domain terms and linked work artifacts"
status: proposed
created: "2026-09-21T00:54:14Z"
updated: "2026-09-21T01:03:37Z"
kind: feature
size: medium
priority: 1
relates_to: ["D-004", "W-018", "W-011"]
---

## Outcome

Give domain vocabulary and work artifacts durable, discoverable homes so agents
and the TUI can retrieve the relevant knowledge or evidence without loading the
whole project. Selected direction is in [D-004](../decisions/D-004-interactive-adoption.md).

## Scope and bounds

Extend the record/configuration contract for brief discovery/location, terms,
and linked plans, reports and reviews. Keep one authoritative brief, separately
identifiable from typed records, with support for the selected default beneath
the configured Grove root. Preserve an existing brief's ownership and links
through an explicit migration rather than creating a competing copy.
Terms describe domain meaning and relationships, not execution
instructions. Start with Grove's work, preparation, attempt, review, approval,
integration, source and revision vocabulary. Define terms during real modeling,
not to satisfy a fixed count.

Artifacts remain ordinary inspectable files and can relate to several work
items. Include enough identity, kind, provenance and relevant revision linkage
for navigation and later review; W-020 owns review approval semantics. Small
preparation can live in a work body. Do not require a standalone plan ticket,
adopt every Bench type, or build semantic search as part of this slice.

## Acceptance

1. CLI/configuration support a discoverable brief location, the selected term
   format and artifact references, with stable identity and clear missing-target
   errors. A brief under the configured Grove root validates; it is retrieved
   deliberately when relevant, not automatically added to every work context.
2. An agent can discover an attached artifact and deliberately include its
   contents through staged context; default context stays bounded and read-only.
3. Shared artifacts, changed reviewed content and historical references remain
   distinguishable; no duplicate editable copy is required for each work item.
4. Configuration/schema compatibility is explicit. Existing schema-1 projects
   remain readable or get a deliberate migration path; unknown values are not
   silently accepted. Document new file locations before using them.
5. Existing links and evidence survive any chosen migration. Update the record
   model, commands and guides to match actual supported behavior.

## Preparation and next

No implementation plan exists yet. Inspect the record reader, allocator, update
contract and staged context. Choose the smallest term and artifact representation,
including schema versioning and multi-work ownership; write a concise plan at
`docs/plans/W-019-knowledge-artifacts.md`. Existing `grove/brief.md` or
`grove/artifacts/*.md` would be rejected today: do not move files before support.
Human judgment is needed only if the design changes the selected information
ownership or meaning. Assign after W-011 in the preferred investment sequence.
