---
id: "W-006"
type: work
title: "Bind workspace routing to the actual project and live checkout"
status: proposed
created: "2026-09-19T20:13:44Z"
updated: "2026-09-19T20:21:16Z"
kind: fix
priority: 1
size: medium
relates_to: ["W-004", "W-005", "Q-001"]
---

## Outcome

A selected version resolves only to the actual project in its registered
checkout, and an observed disappearing/foreign checkout cannot yield a success
path. This is a proposed repair of W-004/W-005's existing contracts, not a new
workspace-opening feature. [Review R1/R2](../../docs/reviews/2026-09-19-integrated-cli.md)
reproduced wrong-repository attribution and success for a deleted checkout at
`9b7f730`. [Implementation plan](../../docs/plans/W-006-workspace-provenance.md).

## Constraints

Keep Q-001's explicit selection and per-source status policy. Keep commands
read-only: no refs, index, worktree administration, records, configuration,
allocator state, editor, shell, or agent writes. Keep plain-directory inspection
unchanged. No selector grammar or schema change, caching, or new dependencies.
A successful location is still not future write authority.

## Design

Validate the path from a registered checkout to the selected project, including
every prefix component, before loading it. A genuinely missing component means
absent; symlinks (even back into this checkout), a file in place of a directory,
unreadable paths, and another nested repository mean invalid with attributable
diagnostics. Confirm the project directory has the expected common directory,
worktree identity, and repository-relative prefix. Do not infer ownership from
the worktree root alone or silently follow an external project.

At the second inventory, invalidate entries that became prunable, inaccessible,
or foreign as well as removed/moved/ref/HEAD changes. Before returning from
Resolve, freshly validate the selected target's registration, ownership, project
configuration and complete source, record path/revision, and attached/detached
identity. Committed routing must recheck branch tip and uniqueness of enterable
checkouts. Keep unrelated invalid sources from preventing an otherwise valid
explicit live selection. Refuse changes with refresh/reselect guidance; never
substitute a new selection. No atomic guarantee against changes after the last
check is promised.

## Acceptance

1. Main's `sub` project versus a feature `sub` symlink to an external Grove
   repository reports the feature source invalid/incomplete with no selectable
   records; resolving an earlier selection refuses with no success stdout.
   Cover internal and dangling symlinks, a regular-file prefix, a nested foreign
   Git repository, and a genuinely absent project. Healthy sources still print.
2. Removing a worktree directory between the first and second inventories leaves
   a prunable registration but invalidates its source. A CLI-level workspace
   test proves it cannot return the removed directory; no timing sleeps.
3. Final-check mutations of target bytes/configuration, path, registration,
   branch/HEAD, ownership, or committed-route multiplicity refuse. Unchanged
   live and committed routes still resolve; dirty unrelated files survive.
4. Existing partial-result, source-local validation, detached/duplicate-checkout,
   and exact-selector fixtures remain green. Hash actual Git state and all
   checkouts across success and refusal paths to prove reads have no writes.
5. Full suite, race suite, vet, formatting, this checkout's `check`, and an
   independent correctness review pass; record actual revisions and residual limits.

## Dependencies and handoff

W-004/W-005 and review fixes are integrated in main. No unfinished product
prerequisite. Recommend one Fable agent in an isolated worktree, serial with
W-007/W-008 because repo/loader/CLI ownership overlaps. The plan supplies the
reproducers and implementation boundaries. Leave the retained implementation
worktree alone. Completion here is separate from integration into main.

## Next

Fable: start from main plus these review artifacts in a fresh worktree; reproduce
R1/R2 before changing code, follow the linked plan, and return focused commits,
verification results, and independent review evidence. Do not add interactive
opening, create worktrees for users, or merge as part of this unit.
