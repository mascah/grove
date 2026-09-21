---
id: "G-011"
type: work
title: "Locate the workspace for a selected record version"
status: done
kind: feature
priority: 2
size: medium
depends_on: ["G-010"]
relates_to: ["G-002", "G-009"]
created: "2026-09-19T17:50:01Z"
updated: "2026-09-19T20:05:45Z"
formerly: "W-005"
---

## Outcome

Given an explicitly selected record version, locate and validate its existing
checkout so a person or future UI can enter the right editing context without
switching main. This implemented CLI foundation is one part of Open workspace;
automatic checkout creation and interactive opening remain follow-on design.

## Why now

Version visibility becomes actionable when the user can reach the selected
checkout. Priority 2 follows the selected coordination experience. Medium
reflects stale selection and ambiguous/missing-worktree handling. G-010 is the
prerequisite because this operation consumes its source identity contract.

## Constraints

- Keep current branch, working directory, refs, records, and staged/dirty files
  unchanged. Do not acquire an agent claim or infer that a checkout is idle.
- G-002 requires explicit version selection before opening a workspace. Use
  that version's source identity. Never guess from record ID, newest timestamp,
  a branch-name-derived directory, or content equality alone.
- No worktree creation/deletion, forced checkout, editor configuration, child
  shell, agent launch, or cross-branch mutation in this unit.
- G-002 owns the selection policy. G-009 remains responsible for safe updates
  after a caller enters a checkout; a returned path is not future write authority.

## Implemented technical contract

CLI surface: `grove workspace --source SELECTOR [--json]`, consuming
G-010's version selector. Success prints the absolute project directory within
the target checkout; JSON includes checkout path, project path, record path,
branch/HEAD, and the validated current content revision. Human diagnostics go
to stderr. Documentation shows how a caller uses the result with `--project`.

For a live source, match its registered worktree identity and verify repository,
project configuration, branch/detached state, HEAD, record identity/path, and
content revision against the selection. For a committed branch source, locate
an existing checkout of that full branch ref and require its observed HEAD and
live record bytes to match the selected committed version. If live content
differs, direct the caller to refresh and select the live observation; do not
silently route a committed selection to different content.

Dirty unrelated files do not prevent location and remain untouched. A missing
checkout, removed/moved worktree, branch/HEAD change, changed configuration,
changed/missing record, invalid project, or ambiguous mapping returns an
attributable error and no success path. Do not pick the first match when forced
multiple checkouts of a branch exist. A user can disambiguate by selecting one
live source in G-010. Preserve arbitrary paths through JSON without shell eval.

Revalidate immediately before returning. Another process can still change the
workspace afterward; document this boundary. A future integrated mutation must
perform its own checks at execution time. A Git worktree lock protects worktree
administration and must not be described as exclusive editing ownership.

## Acceptance

- From main, select a feature's live version and obtain its exact project path,
  including a project below repository root, without changing main or Git state.
- A matching committed selection resolves to its existing checkout; differing
  live content refuses with a refresh/reselect diagnostic.
- Branch/HEAD changes after selection, detached state changes, moved/removed
  worktrees, changed configuration, and changed/missing record bytes refuse.
- Missing and ambiguous checkouts return distinct errors and create nothing;
  an explicitly selected live checkout disambiguates matching branch checkouts.
- Unrelated staged and dirty files survive unchanged. Paths with spaces and
  control characters round-trip through JSON; no path becomes executable text.
- A joint fixture lists versions, selects one, resolves the workspace, and uses
  existing `show --project` to read exactly the selected record bytes there.

## Preparation and execution boundary

Prepare against G-010's completed source/selector interface. Expected ownership
is its Git-source package and `internal/cli`; reuse discovery and validation.
Recommend a single Fable agent after G-010, with full suite, race suite, vet,
and independent review of stale selections and unintended Git writes.

## Evidence

Done 2026-09-19 on branch `worktree-W-004-W-005` at `d232aa2`, after G-010
closed at `a465d88` in the same branch. The
[coordination plan](G-013-coordination-plan.md) records
the consumed selector, this command's resolution and refusal contract, the
commits, the fixture list per acceptance item, suite, race, vet, gofmt, and
`check` results, real use in this repository, the joint verification, and
the combined review. `workspace --source SELECTOR [--json]` re-runs the
G-010 inspection for the selected record and prints the absolute project
directory of the existing checkout that still holds exactly that version;
JSON adds checkout, record, branch or detached HEAD, and current revision.

Every acceptance bullet has fixtures: a live feature selection from main with
the project below the repository root, returning the exact path with both
checkouts and Git state hashed unchanged; a matching committed selection
resolving to its checkout while differing live content, path, or
configuration refuses with a refresh-and-reselect diagnostic; branch and
HEAD moves, detaching and re-attaching, moved and removed worktrees, changed
configuration, and changed, renamed, or deleted records each refusing with
their own reason; missing and ambiguous checkouts as distinct errors that
create nothing, with an explicit live selection disambiguating; staged and
unstaged unrelated files surviving; a worktree path with a newline preserved
in JSON and escaped on stderr; and the joint fixture listing versions,
selecting one, resolving its workspace, and reading the selected bytes
there with `show --project`.

Review disposition: no blocking findings; three should-fix items (absent
project wording, prunable duplicates, committed-route configuration check)
were fixed and covered. This is the CLI foundation for Open workspace: it
locates, it does not create worktrees, switch branches, launch editors,
agents, or shells, claim ownership, or edit records, and the printed path is
a location rather than write authority. A later `update` performs its own
revision check because the workspace can change after resolution.

## Next

Integrated into main at `5041ae1` on 2026-09-19. Repair the routing/freshness
and Git-path defects in [G-014](G-014-workspace-provenance.md) and
[G-016](G-016-git-paths.md) before building actions on this command. The
[integrated review](G-022-integrated-cli-review.md) supplies
reproducers and distinguishes tested behavior from remaining limits. Then shape
the interactive existing-workspace experience; checkout creation needs its own
write/failure design. Path resolution is not that entire experience.
