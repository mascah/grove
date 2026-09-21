---
id: "G-043"
type: work
title: "Make the board and item detail clear and visually polished"
status: proposed
created: "2026-09-21T00:54:15Z"
updated: "2026-09-21T15:45:42Z"
kind: feature
size: medium
priority: 3
depends_on: ["G-042"]
relates_to: ["G-035", "G-017", "G-030", "G-041", "G-044", "G-064", "G-065"]
formerly: "W-025"
---

## Outcome

A visually polished terminal board and detail view make current work, its
content, blockers, artifacts and next action immediately understandable.

## Scope and bounds

Use G-042's project-wide projection. Board columns follow the target lifecycle;
show compact cards and clear borders/focus, hide Abandoned by default, and bound
the recent Done column while retaining searchable older work. Detail leads with
rendered Markdown; metadata, artifacts and a navigable timeline support it.
Keep alternate sources available as secondary information. Include a useful
list/search path; dependency graph visualization is later.

Bound Done through presentation, never by relocating files. Resolve type and
artifact roles from supported metadata, not ID prefix or folder. List/search
also makes G-065's general knowledge discoverable without turning it into work
cards or loading every page into an agent's context.

Design board, detail and the future review view together before implementation,
but implement review actions in G-044. Use Bubble Tea v2 with compatible Charm
layout/components/Markdown rendering as appropriate. Preserve terminal text
escaping and separate trusted renderer styling from untrusted content.

## Acceptance

1. The owner reviews a concrete visual proposal and the implemented terminal
   experience, including representative narrow/wide sizes and many Done items.
2. Current content is the first useful detail; artifacts and timeline entries
   can be reached without interpreting a list of checkout observations.
3. Focus, scrolling, search, empty/loading/incomplete states and navigation
   remain understandable; meaning does not depend on color alone.
4. Lifecycle, alternate-screen restoration, resize and cancellation checks pass;
   verify connected navigation using real work and attached artifacts.
5. On-demand history and measured loading stay responsive. Bare invocation and
   explicit noninteractive command contracts remain intact.

## Preparation and next

No implementation plan exists yet. Create a linked plan through the supported
CLI with visual examples using real records and G-041 feedback. Check dependency
compatibility rather than assuming all Charm modules share the same major
version. Design can begin during the interactive loop, without making the
complete visual redesign a prerequisite for the hobby-project pilot.
