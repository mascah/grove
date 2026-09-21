---
id: "W-022"
type: work
title: "Bootstrap projects with portable Grove workflows"
status: proposed
created: "2026-09-21T00:54:15Z"
updated: "2026-09-21T15:45:41Z"
kind: feature
size: medium
priority: 2
depends_on: ["W-021"]
relates_to: ["D-004", "W-018", "W-011", "W-019", "W-020", "D-006", "W-030"]
---

## Outcome

An existing repository can adopt a versioned Grove CLI and shared workflows,
then shape and execute work through thin Claude/Codex entrypoints without a
dependency on this repository's development-only instructions.

## Scope and bounds

Provide minimal setup with valid `grove.yaml`, discoverable brief/knowledge
locations and harness entrypoints. Create directories as needed and preserve
existing AGENTS/CLAUDE instructions. Setup can leave a clearly incomplete brief
for an interactive shaping session; it must not invent project intent.

The target default uses W-030's one root, flat creation, neutral new IDs and
general knowledge pages, as selected in
[D-006](../decisions/D-006-stable-knowledge.md). Do not pre-create type folders,
require knowledge classification, add per-type routing/prefix settings, or move
existing records when setup is rerun. Keep shared workflow
instructions versioned with Grove; adapters load one owner rather than divergent
editable copies. Choose packaging/update behavior during preparation. A TUI
setup wizard is optional later, not a prerequisite for this outcome.

## Acceptance

1. A disposable existing repository can initialize, validate, shape and hand
   off selected work through the new CLI without predecessor initialization.
2. Re-running setup preserves user content and makes managed updates explicit;
   conflicting configuration or instruction blocks receive actionable outcomes.
3. The actual executable and workflow version are unambiguous. Document staged
   coexistence with the predecessor's installed `grove` without replacing it
   globally as a side effect of setup or testing.
4. Exercise fresh-session discovery in Claude and Codex where available; report
   real behavior separately from file existence and simulated harness calls.
5. Package boundaries separate portable workflow from this repo's Go tests,
   worktree naming, `go run` invocation and development history.

## Preparation and next

No implementation plan exists yet. Create a linked plan through the supported
CLI from W-030/W-020 interfaces and W-021 trial feedback. Select a distribution and
upgrade mechanism using installed harness capabilities. Test generated projects
in disposable clones with explicit absolute project paths. Live sibling
installation and migration belong to W-023, not this assignment.
