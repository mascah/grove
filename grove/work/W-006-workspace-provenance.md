---
id: "W-006"
type: work
title: "Bind workspace routing to the actual project and live checkout"
status: done
created: "2026-09-19T20:13:44Z"
updated: "2026-09-19T22:07:49Z"
kind: fix
priority: 1
size: medium
relates_to: ["W-004", "W-005", "Q-001"]
---

## Outcome

A selected version resolves only to the actual project in its registered
checkout, and an observed disappearing/foreign checkout cannot yield a success
path. This repairs W-004/W-005's existing contracts (implemented; see
Evidence); it is not a new workspace-opening feature. [Review R1/R2](../../docs/reviews/2026-09-19-integrated-cli.md)
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

## Evidence

Implemented 2026-09-19 on branch `worktree-W-006-W-008` (base `2d6de36`), not
yet integrated into main. Code: `d5666dd` (committed projects two or more
levels deep), `5316dbe` (ownership through the prefix, second-inventory
re-entry), `892a842` (final check in `Resolve`), `9e8430c` (review fixes), `7f02b71`
(combined-review fix).
The [plan](../../docs/plans/W-006-workspace-provenance.md#implementation-notes-2026-09-19)
records the bounded adjustments. Selector grammar, JSON, schema, and CLI
surface are unchanged; both commands still write nothing.

Reproduced before repair, on a pristine export of `2d6de36` holding only the
new tests (plus `d5666dd`, without which a deep prefix cannot be read at all):
`TestInspectProjectLocation` reported `Complete:true` and a valid feature
source for an external, middle, internal, and dangling symlink and for nested
repositories at the prefix and at the record folder (R1);
`TestInspectPrunableDuringRead` and the CLI Git-wrapper test
`TestWorkspaceCheckoutDeletedDuringInspection` returned exit 0 and a JSON
workspace inside the deleted checkout (R2). All pass on the candidate.

Acceptance: (1) `TestInspectProjectLocation` covers nine locations; a missing
directory or parent is absent, the rest invalid with an attributable
diagnostic, main's records stay visible, an earlier selection refuses, and
the enclosing directory hashes unchanged; `TestForeignProjectPrefixCLI`
proves exit 1, `complete: false`, no attributed version, and no stdout from
`workspace`. (2) `TestInspectPrunableDuringRead`,
`TestInspectForeignDuringRead`, `TestInspectConfigurationRemovedDuringRead`,
`TestInspectWorktreeReplacedByPlainDirectory`, and the wrapper test, with no
sleeps. (3) `TestResolveFinalCheck`: seventeen mutations between inspection
and return (record bytes, path, deletion, configuration, invalid source,
removed or foreign or symlinked project, deleted checkout, moved or removed
registration, HEAD, branch, detaching, second checkout of a committed
route's branch) each refuse without writing; `TestResolveFinalCheckAdmits`
keeps unchanged live and committed routes, an unrelated invalid source, a
dirty unrelated file, and an explicit live selection beside a new duplicate.
(4) The existing W-004/W-005 fixtures pass unchanged.

Independent review (separate reviewer agent, range `2d6de36..892a842`, own
export): no P1/P2; reproducer confirmed real; five adversarial probes of its
own (foreign repository at an intermediate component, worktree root or its
parent swapped for a symlink, symlinked `grove.yaml`, a swap inside the final
check's hook) were all refused. P3 findings fixed in `9e8430c`: configuration
vanishing during the read, an over-broad parent-entry filter, the cost note,
and a CLI-level proof for acceptance 1.

The combined review then found a P1 introduced by a W-008 review fix: a
foreign repository registered below `<common>/worktrees` was admitted once
the common-directory comparison had been removed. `7f02b71` restores it with
`TestInspectForeignRepositoryRegisteredAsWorktree`, which fails at `40e882f`.

Limits: changes after the final check remain possible, as the contract
says; a crafted `.git` file pointing a registered path at another worktree's
Git directory is not detected; a live checkout now costs about eight
`rev-parse` processes per inspection, twelve under a nested prefix, where it
cost one (for main plus three worktrees, 10 Git processes became 39, or 55
nested, measured at the final revision); macOS only,
no Windows or case-insensitive collision testing. Suite results are in the
[combined repair evidence](../../docs/reviews/2026-09-19-repairs-W-006-W-008.md).

## Next

Integrate branch `worktree-W-006-W-008` into main as a separate, explicit
step; this record being done asserts completion on that branch only. W-009's
board can then depend on `versions` and `workspace` for routing.
