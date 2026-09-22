---
id: "G-040"
type: work
title: "Bootstrap projects with portable Grove workflows"
status: active
created: "2026-09-21T00:54:15Z"
updated: "2026-09-22T16:06:20Z"
kind: feature
size: medium
priority: 2
depends_on: ["G-039"]
relates_to: ["G-035", "G-036", "G-025", "G-037", "G-038", "G-064", "G-065"]
formerly: "W-022"
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

The target default uses G-065's one root, flat creation, neutral new IDs and
general knowledge pages, as selected in
[G-064](G-064-stable-knowledge.md). Do not pre-create type folders,
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

Assigned alone (`/grove-work G-040`, interactive) on 2026-09-22; branch
`worktree-G-040` in `.claude/worktrees/G-040`, base main `ccdc92d`, this
record at `sha256:5d171634…`. The plan is
[G-080](G-080-portable-bootstrap-plan.md): the binary embeds the two guides
and prints them (`grove guide`), `grove version` names the executable and so
the workflow, and `grove init` writes `grove.yaml`, the record root, a
placeholder brief and six marked adapters, keeping user files and reporting
managed updates. Distribution is a build from a named commit; the predecessor
stays installed. Test generated projects in disposable clones with explicit
absolute project paths. Live sibling installation and migration belong to
G-041, not this assignment.
