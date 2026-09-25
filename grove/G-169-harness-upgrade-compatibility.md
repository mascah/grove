---
id: "G-169"
type: work
title: "Keep installed harness entrypoints compatible with Grove upgrades"
status: proposed
created: "2026-09-25T21:04:46Z"
updated: "2026-09-25T21:04:46Z"
---

## Outcome

An adopting project can upgrade Grove and continue interactive shaping and work
through Claude and Codex without silently combining incompatible installed
instructions and executable behavior. Workflow changes travel with the binary;
project-owned policy stays with the project.

Owner intent, review and `$grove-shape` conversation 2026-09-25: remove reliance
on Grove's development checkout and the owner's environment, and make the skills
that `init` installs understandable to maintain and distribute. The owner asked
to shape the review's recommendations; the design below remains proposed.

## Observed evidence

At main `05892a2`, [init.go](../internal/cli/init.go) generates skills which
load `grove guide work|shape`, but also encode allowed arguments. G-040's
candidate `4edd800` rejects input other than IDs and interaction mode; today's
template accepts `--until plan`. An older installed adapter can therefore
contradict a newer guide. The full [reviewer](../.claude/agents/grove-reviewer.md)
is copied, rather than loaded at use time.

The review ran a binary built from `git archive HEAD` in a disposable external
Git repository with a stripped environment. After simulated edits to managed
skill and reviewer files, `check` still printed `OK: 0 records`; rerunning
`init` restored both and preserved an unmarked custom skill. This proves current
refresh and ownership behavior, not a real two-release or harness upgrade trial.
[G-150](G-150-launch-attempts-only-where-the-w.md) checks skill presence before
launch; it does not establish compatibility. Existing worktrees retain their
own files after another checkout runs `init`.

## Proposed design and constraints

- Keep one binary-owned workflow. Move evolving assignment grammar into the
  guide. Make the reviewer load binary-owned review instructions, for example
  through `grove guide review`; retain necessary harness metadata and the
  reviewer's read-only authority in its adapter. Preserve argument-as-data
  handling and explicit invocation.
- Give the small entrypoint interface an identifiable compatibility revision,
  separate from package releases and record schema. Specify how legacy files
  without metadata are diagnosed. Content inequality alone is not evidence of
  incompatibility; compatible older entrypoints should remain usable.
- Provide a read-only installation diagnostic, for example `init --check`,
  that distinguishes current, compatible older, incompatible, missing and
  project-owned files. Make supported interactive entrypoints and the runner
  surface incompatibility before work or paid execution begins. Preparation
  must explain how legacy entrypoints reach that diagnosis and recovery.
- Retain explicit `init` refresh and the existing ownership marker. It may
  replace marked files; it must preserve unmarked custom files, configuration,
  brief, records and AGENTS/CLAUDE policy. Never silently repair another
  worktree, commit changes, or claim a custom file is compatible without evidence.
- Explain commit requirements, existing-worktree repair, and restarting or
  reloading sessions which already loaded older instructions. Keep guide access
  usable without a project. Keep this repository's development adapters working.

Alternatives: separately released harness plugins add an installation/version
lifecycle; full copied guides add synchronization work. Neither is proposed.
No new execution provider, daemon, self-update, schema migration or broad
compatibility promise is authorized. [G-101](G-101-attempt-mechanism.md) remains
the runner decision; [G-152](G-152-shipped-document.md) names the shipped-document
boundary. G-040's marker-only approach would be extended, with existing
project ownership preserved. Release policy belongs to
[G-172](G-172-first-release-policy.md), not an inferred schema change here.

## Acceptance

1. In a disposable adopting project, changing only the installed binary changes
   the work, shaping and review instructions used by a fresh session without
   rewriting compatible entrypoints. Evolving argument rules have one owner.
2. Read-only diagnosis covers current, compatible older, incompatible, legacy,
   missing and custom entrypoints without writes. Incompatible supported
   invocations stop actionably; runner refusal occurs before paid execution.
3. Explicit refresh is idempotent, preserves project-owned content, and gives a
   usable path for a committed older checkout and an existing worktree. Tests
   exercise an actual old/new fixture pair, not only repeated fresh `init`.
4. Evidence distinguishes deterministic checks from fresh Claude/Codex session
   discovery and interactive invocation. Live trials need an explicit resource
   mandate; unavailable trials are reported as limits, not simulated successes.
5. Command documentation and shipped instructions explain the contract and
   repair path; repository verification and shipped-document checks pass.

## Next

Proposed, unassigned. Prepare this work independently of release packaging,
covering the compatibility interface, legacy transition, reviewer loading and
interactive/runner checks. The owner can assign `$grove-work G-169`; live agent
trials require separately stated resource bounds. G-110 needs this outcome for
its package upgrade rehearsal.
