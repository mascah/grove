---
id: "G-069"
type: page
title: "Identity and path migration map"
created: "2026-09-21T21:11:16Z"
updated: "2026-09-21T21:11:16Z"
---

W-029 reconciled every record, legacy plan and legacy review into neutral IDs
directly under `grove/` on 2026-09-21. This page is the one durable mapping
from an old ID or path to its current counterpart. To read an old commit, use
the old ID or path there with that commit's own CLI (`go run ./cmd/grove`);
to find what it became, look it up here or search `formerly:`. Numbers follow
document date (a record's `created`, a legacy document's first commit), ties
broken on the old ID or path; the order conveys no authority or priority.

| Old ID | Old path | ID | Path | Document date |
| --- | --- | --- | --- | --- |
| D-001 | `grove/decisions/D-001-starter-defaults.md` | G-001 | `grove/G-001-starter-defaults.md` | 2026-09-19T14:08:40Z |
| Q-001 | `grove/questions/Q-001-branch-versions.md` | G-002 | `grove/G-002-branch-versions.md` | 2026-09-19T14:08:40Z |
| W-001 | `grove/work/W-001-inspect-records.md` | G-003 | `grove/G-003-inspect-records.md` | 2026-09-19T14:08:40Z |
| D-002 | `grove/decisions/D-002-sequential-ids.md` | G-004 | `grove/G-004-sequential-ids.md` | 2026-09-19T14:32:32Z |
| - | `docs/plans/W-001-inspection.md` | G-005 | `grove/G-005-inspection-plan.md` | 2026-09-19T15:19:41Z |
| D-003 | `grove/decisions/D-003-allocator-mechanism.md` | G-006 | `grove/G-006-allocator-mechanism.md` | 2026-09-19T15:25:51Z |
| W-002 | `grove/work/W-002-create-records.md` | G-007 | `grove/G-007-create-records.md` | 2026-09-19T15:25:51Z |
| - | `docs/plans/W-002-create.md` | G-008 | `grove/G-008-create-plan.md` | 2026-09-19T15:31:29Z |
| W-003 | `grove/work/W-003-update-records.md` | G-009 | `grove/G-009-update-records.md` | 2026-09-19T15:36:19Z |
| W-004 | `grove/work/W-004-record-versions.md` | G-010 | `grove/G-010-record-versions.md` | 2026-09-19T17:49:58Z |
| W-005 | `grove/work/W-005-record-workspace.md` | G-011 | `grove/G-011-record-workspace.md` | 2026-09-19T17:50:01Z |
| - | `docs/plans/W-003-update.md` | G-012 | `grove/G-012-update-plan.md` | 2026-09-19T18:42:44Z |
| - | `docs/plans/W-004-W-005-coordination.md` | G-013 | `grove/G-013-coordination-plan.md` | 2026-09-19T19:28:26Z |
| W-006 | `grove/work/W-006-workspace-provenance.md` | G-014 | `grove/G-014-workspace-provenance.md` | 2026-09-19T20:13:44Z |
| W-007 | `grove/work/W-007-preserve-updates.md` | G-015 | `grove/G-015-preserve-updates.md` | 2026-09-19T20:13:54Z |
| W-008 | `grove/work/W-008-git-paths.md` | G-016 | `grove/G-016-git-paths.md` | 2026-09-19T20:13:57Z |
| W-009 | `grove/work/W-009-terminal-picker.md` | G-017 | `grove/G-017-terminal-picker.md` | 2026-09-19T20:23:01Z |
| - | `docs/plans/W-006-workspace-provenance.md` | G-018 | `grove/G-018-workspace-provenance-plan.md` | 2026-09-19T20:32:24Z |
| - | `docs/plans/W-007-preserve-updates.md` | G-019 | `grove/G-019-preserve-updates-plan.md` | 2026-09-19T20:32:24Z |
| - | `docs/plans/W-008-git-paths.md` | G-020 | `grove/G-020-git-paths-plan.md` | 2026-09-19T20:32:24Z |
| - | `docs/plans/W-009-terminal-picker.md` | G-021 | `grove/G-021-terminal-picker-plan.md` | 2026-09-19T20:32:24Z |
| - | `docs/reviews/2026-09-19-integrated-cli.md` | G-022 | `grove/G-022-integrated-cli-review.md` | 2026-09-19T20:32:24Z |
| W-010 | `grove/work/W-010-work-handoffs.md` | G-023 | `grove/G-023-work-handoffs.md` | 2026-09-19T20:44:28Z |
| - | `docs/reviews/2026-09-19-predecessor-work.md` | G-024 | `grove/G-024-predecessor-work-review.md` | 2026-09-19T20:59:57Z |
| W-011 | `grove/work/W-011-shaping-entrypoint.md` | G-025 | `grove/G-025-shaping-entrypoint.md` | 2026-09-19T21:20:41Z |
| - | `docs/reviews/2026-09-19-shaping-and-runner-evidence.md` | G-026 | `grove/G-026-shaping-and-runner-evidence-review.md` | 2026-09-19T21:24:18Z |
| - | `docs/plans/W-010-W-011-agent-handoffs.md` | G-027 | `grove/G-027-agent-handoffs-plan.md` | 2026-09-19T21:39:23Z |
| - | `docs/reviews/2026-09-19-repairs-W-006-W-008.md` | G-028 | `grove/G-028-repairs-review.md` | 2026-09-19T22:07:59Z |
| - | `docs/reviews/2026-09-19-board-W-009.md` | G-029 | `grove/G-029-board-review.md` | 2026-09-19T23:28:01Z |
| W-012 | `grove/work/W-012-card-lineage.md` | G-030 | `grove/G-030-card-lineage.md` | 2026-09-20T04:36:57Z |
| W-013 | `grove/work/W-013-load-scaling.md` | G-031 | `grove/G-031-load-scaling.md` | 2026-09-20T04:36:57Z |
| - | `docs/reviews/2026-09-19-W-010-dogfood.md` | G-032 | `grove/G-032-dogfood-review.md` | 2026-09-20T05:50:31Z |
| - | `docs/plans/W-012-card-lineage.md` | G-033 | `grove/G-033-card-lineage-plan.md` | 2026-09-20T14:55:55Z |
| - | `docs/reviews/2026-09-20-card-lineage-W-012.md` | G-034 | `grove/G-034-card-lineage-review.md` | 2026-09-20T16:11:14Z |
| D-004 | `grove/decisions/D-004-interactive-adoption.md` | G-035 | `grove/G-035-interactive-adoption.md` | 2026-09-21T00:54:13Z |
| W-018 | `grove/work/W-018-interactive-adoption.md` | G-036 | `grove/G-036-interactive-adoption.md` | 2026-09-21T00:54:13Z |
| W-019 | `grove/work/W-019-knowledge-artifacts.md` | G-037 | `grove/G-037-knowledge-artifacts.md` | 2026-09-21T00:54:14Z |
| W-020 | `grove/work/W-020-review-lifecycle.md` | G-038 | `grove/G-038-review-lifecycle.md` | 2026-09-21T00:54:14Z |
| W-021 | `grove/work/W-021-interactive-loop.md` | G-039 | `grove/G-039-interactive-loop.md` | 2026-09-21T00:54:14Z |
| W-022 | `grove/work/W-022-portable-bootstrap.md` | G-040 | `grove/G-040-portable-bootstrap.md` | 2026-09-21T00:54:15Z |
| W-023 | `grove/work/W-023-nullsec-pilot.md` | G-041 | `grove/G-041-nullsec-pilot.md` | 2026-09-21T00:54:15Z |
| W-024 | `grove/work/W-024-current-view.md` | G-042 | `grove/G-042-current-view.md` | 2026-09-21T00:54:15Z |
| W-025 | `grove/work/W-025-board-detail.md` | G-043 | `grove/G-043-board-detail.md` | 2026-09-21T00:54:15Z |
| W-026 | `grove/work/W-026-review-integration.md` | G-044 | `grove/G-044-review-integration.md` | 2026-09-21T00:54:16Z |
| W-027 | `grove/work/W-027-durable-attempt.md` | G-045 | `grove/G-045-durable-attempt.md` | 2026-09-21T00:54:16Z |
| W-028 | `grove/work/W-028-managed-runs.md` | G-046 | `grove/G-046-managed-runs.md` | 2026-09-21T00:54:16Z |
| - | `docs/plans/W-018-adoption-roadmap.md` | G-047 | `grove/G-047-adoption-roadmap-plan.md` | 2026-09-21T01:08:42Z |
| - | `docs/reviews/2026-09-20-direction-evaluation.md` | G-048 | `grove/G-048-direction-evaluation-review.md` | 2026-09-21T01:08:42Z |
| - | `docs/plans/W-011-shaping-entrypoint.md` | G-049 | `grove/G-049-shaping-entrypoint-plan.md` | 2026-09-21T02:53:09Z |
| - | `docs/reviews/2026-09-20-W-011-shaping.md` | G-050 | `grove/G-050-shaping-review.md` | 2026-09-21T03:02:56Z |
| D-005 | `grove/decisions/D-005-typed-knowledge-records.md` | G-051 | `grove/G-051-typed-knowledge-records.md` | 2026-09-21T04:37:58Z |
| W-029 | `grove/work/W-029-migrate-knowledge.md` | G-052 | `grove/G-052-migrate-knowledge.md` | 2026-09-21T04:37:58Z |
| - | `docs/plans/W-019-knowledge-artifacts.md` | G-053 | `grove/G-053-knowledge-artifacts-plan.md` | 2026-09-21T04:55:41Z |
| T-001 | `grove/terms/T-001-work.md` | G-054 | `grove/G-054-work.md` | 2026-09-21T05:01:55Z |
| T-002 | `grove/terms/T-002-preparation.md` | G-055 | `grove/G-055-preparation.md` | 2026-09-21T05:01:56Z |
| T-003 | `grove/terms/T-003-attempt.md` | G-056 | `grove/G-056-attempt.md` | 2026-09-21T05:01:56Z |
| T-004 | `grove/terms/T-004-candidate.md` | G-057 | `grove/G-057-candidate.md` | 2026-09-21T05:01:56Z |
| T-005 | `grove/terms/T-005-review.md` | G-058 | `grove/G-058-review.md` | 2026-09-21T05:01:56Z |
| T-006 | `grove/terms/T-006-approval.md` | G-059 | `grove/G-059-approval.md` | 2026-09-21T05:01:56Z |
| T-007 | `grove/terms/T-007-integration.md` | G-060 | `grove/G-060-integration.md` | 2026-09-21T05:01:57Z |
| T-008 | `grove/terms/T-008-source.md` | G-061 | `grove/G-061-source.md` | 2026-09-21T05:01:57Z |
| T-009 | `grove/terms/T-009-revision.md` | G-062 | `grove/G-062-revision.md` | 2026-09-21T05:01:57Z |
| R-001 | `grove/reviews/R-001-w-019-knowledge-records.md` | G-063 | `grove/G-063-knowledge-records-review.md` | 2026-09-21T05:59:58Z |
| D-006 | `grove/decisions/D-006-stable-knowledge.md` | G-064 | `grove/G-064-stable-knowledge.md` | 2026-09-21T15:40:59Z |
| W-030 | `grove/work/W-030-flexible-records.md` | G-065 | `grove/G-065-flexible-records.md` | 2026-09-21T15:40:59Z |
| P-001 | `grove/plans/P-001-w-030-flexible-records.md` | G-066 | `grove/G-066-flexible-records-plan.md` | 2026-09-21T16:02:50Z |
| R-002 | `grove/reviews/R-002-w-030-flexible-records.md` | G-067 | `grove/G-067-flexible-records-review.md` | 2026-09-21T16:32:01Z |
| P-002 | `grove/plans/P-002-w-029-reconciliation.md` | G-068 | `grove/G-068-reconciliation-plan.md` | 2026-09-21T21:00:03Z |
| - | `docs/restart-brief.md` | - | `grove/brief.md` | the brief, not a record |

## Kept in place

- `docs/prompts/*.txt`: three spent handoff prompts, historical evidence that
  the handoff work links. Their old IDs and paths are literals of their time.
- `README.md`, `AGENTS.md`, the skill adapters, `docs/record-model.md`,
  `docs/work-execution.md` and `docs/work-shaping.md`: functional homes, with
  references updated.

## Historical literals

Old identifiers remain in `formerly` fields, in this table, in Git branch
names such as `worktree-W-012` and in commit messages, which name what existed
then. The typed IDs named on these lines were left as written: they belong to
test fixtures, disposable trial clones or the predecessor project rather than
to this repository's records, or they sit in a verbatim quotation (the
owner's words, a commit subject, the 2026-09-19 renumbering table). Every
other typed ID that had a counterpart was rewritten, command transcripts
included:

- `grove/G-008-create-plan.md:71: W-003, W-005`
- `grove/G-008-create-plan.md:73: W-004`
- `grove/G-008-create-plan.md:93: W-030`
- `grove/G-008-create-plan.md:95: W-007`
- `grove/G-034-card-lineage-review.md:51: W-001`
- `grove/G-022-integrated-cli-review.md:77: W-001`
- `grove/G-022-integrated-cli-review.md:131: W-001`
- `grove/G-022-integrated-cli-review.md:134: W-001`
- `grove/G-030-card-lineage.md:29: W-006`
- `grove/G-026-shaping-and-runner-evidence-review.md:24: W-003`
- `grove/G-026-shaping-and-runner-evidence-review.md:38: W-003`
- `grove/G-032-dogfood-review.md:54: W-001, W-002, W-003`
- `grove/G-027-agent-handoffs-plan.md:329: W-001, W-002, W-003`
- `grove/G-027-agent-handoffs-plan.md:330: W-001, W-002, W-003`
- `grove/G-027-agent-handoffs-plan.md:402: W-001`
- `grove/G-027-agent-handoffs-plan.md:406: W-001`
- `grove/G-027-agent-handoffs-plan.md:409: W-001`
- `grove/G-006-allocator-mechanism.md:14: D-003`
- `grove/G-029-board-review.md:52: W-001`
- `grove/G-029-board-review.md:53: W-002`
- `grove/G-029-board-review.md:66: W-001`
- `grove/G-004-sequential-ids.md:40: W-001`
- `grove/G-004-sequential-ids.md:41: Q-001`
- `grove/G-004-sequential-ids.md:42: D-001`
- `grove/G-020-git-paths-plan.md:88: W-001`
- `grove/G-005-inspection-plan.md:117: W-001`
- `grove/G-024-predecessor-work-review.md:32: W-002`
- `grove/G-013-coordination-plan.md:264: W-001`
- `grove/G-013-coordination-plan.md:293: W-001`
- `grove/G-063-knowledge-records-review.md:31: T-001`
- `grove/G-068-reconciliation-plan.md:62: W-001`
- `grove/G-018-workspace-provenance-plan.md:63: W-001`
- `grove/G-018-workspace-provenance-plan.md:106: W-001`
- `grove/G-019-preserve-updates-plan.md:56: W-001`
- `grove/G-019-preserve-updates-plan.md:85: W-001`
- `grove/G-019-preserve-updates-plan.md:86: W-001`
- `grove/G-019-preserve-updates-plan.md:119: W-001`
- `grove/G-002-branch-versions.md:43: W-003`
- `grove/G-012-update-plan.md:144: W-002`
- `grove/G-050-shaping-review.md`: every `W-029` and `W-030`, which the trial wrote in its own clones before these numbers existed here
- `grove/G-015-preserve-updates.md:106: W-003`
- `grove/G-065-flexible-records.md:130-131: W-019, R-001, W-030, P-001, D-006` (the inputs of a `convert` rehearsal, which takes typed IDs)
- `grove/G-064-stable-knowledge.md:80: W-019` (evidence about schema 2's folder rule)

A typed ID inside a filename or path (`docs/prompts/W-009-implementation.txt`,
a fixture's `grove/work/W-001-first.md`, a hypothetical `work/W-019/plan.md`)
was never rewritten either: it names a file, not a record. A typed ID with no row in the table (`W-014` to `W-017`, `W-090`, `Q-002` and
the like) never named a record here: it is a fixture's or the predecessor's.
Records and evidence dated before this migration also name the folders of
their time (`grove/work/`, `docs/plans/`, `docs/reviews/`) where they describe
what was then true; nothing lives there now.
