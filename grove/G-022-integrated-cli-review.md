---
id: "G-022"
type: review
title: "Integrated CLI review, 2026-09-19"
status: current
formerly: "docs/reviews/2026-09-19-integrated-cli.md"
work: ["G-009", "G-010", "G-011", "G-014", "G-015", "G-016"]
updated: "2026-09-21T21:11:16Z"
---

# Integrated CLI review, 2026-09-19

Review base: `9b7f73000ecb52f151fb69d078f363df678e77dd` on main. Scope:
G-009 updates, G-010 versions, G-011 workspace resolution, their shared
reader/coordination code, integration evidence, and next implementation handoffs.
No product implementation or merge was performed during this review.

## Integration and evidence

- Main and the retained `.claude/worktrees/W-004-W-005` checkout were clean
  at entry. The latter is on `worktree-W-004-W-005` at `c6b6e9d`.
  `git rev-list --left-right --count main...worktree-W-004-W-005` returned
  `2 0`: no branch-only commits await integration. Leave that worktree intact.
- G-009's `14ed115` review fix and `6a5e297` closure are ancestors of main,
  integrated by `b20d2b0`. G-010's `ca420f5` and G-011's `d232aa2`
  review fixes and both closures are ancestors, integrated by `5041ae1`.
- The [update plan](G-012-update-plan.md) and
  [coordination plan](G-013-coordination-plan.md) retain reviewed
  ranges, test mappings, fixes, and reviewer summaries. Their tests and fixes
  exist. The claimed historical fuzz run and original reviewer sessions were
  not independently replayed; these are retained author reports, not new evidence.
- Done means completion asserted in a branch context. Integration here is
  established by Git ancestry, independently of the work records' statuses.

## Fresh verification

All commands below ran against the unchanged implementation at the review base:

| Check | Actual result |
| --- | --- |
| `go test -count=1 ./...` | PASS, all six internal packages; CLI entrypoint has no tests |
| `go test -race -count=1 ./...` | PASS, all six internal packages |
| `go vet ./...` | PASS |
| `gofmt -l cmd internal` | No output |
| `go run ./cmd/grove check` | PASS, 9 records before shaping |
| `go test -count=1 -timeout=60s -v ./internal/versions` | PASS, all inspection and workspace tests |
| Real `versions G-011 --json` → all four selectors → `workspace --json` → `show --project` | All resolved and read the exact selected bytes |
| Hashes before/after that real read workflow | All 395 Git metadata/configuration/record files unchanged, including both checkouts' records |

Additional probes used a temporary copy of the Go sources and disposable Git
repositories, with no source changes in the working repository. Six new
expectation tests failed on the current implementation (as described below):
`TestReviewPrunableDuringRead`, `TestReviewKeyComment`,
`TestReviewFlowAdjacentUnset`, `TestReviewConfigOnlyChange`,
`TestReviewCreateConfigChange`, and `TestReviewOddWorktreeFloor`.
The temporary copy used a writable `/tmp` Go build cache after the default
cache was denied. These failures do not contradict the passing existing suite;
they expose missing coverage. A CLI binary built from the same base reproduced
foreign-prefix routing, removed-worktree success, explicit-key refusal, and
newline-path problems.

## Actionable findings

Severity: P1 = fix before workspace actions depend on it; P2 = concrete
correctness/contract defect requiring a bounded repair; P3 = documentation drift.
Line references below are at the review base.

### R1 — P1: A live project prefix can route into another repository

**Demonstrated bug.** `internal/versions/versions.go:129-144` joins the selected
prefix to each worktree and loads it without checking the prefix's path
components or Git ownership. Only the worktree root was checked earlier.
`internal/versions/workspace.go:194` then returns that textual project path.

Reproduction: main holds a project at `sub/`; a feature worktree removes that
project and has an untracked `sub` symlink to a separate Git repository with
valid Grove records. From main's `sub`, `versions W-001 --json` returned
`complete: true` and attributed the external record to live `refs/heads/feature`.
Resolving its selector returned exit 0, the feature ref/HEAD, and
`<feature>/sub`, which actually opens the external project. No content or Git
mutation was needed to demonstrate this mismatch.

Reject symlink traversal through the project prefix and reject a prefix that
belongs to a different nested repository. Check the selected project's actual
worktree/common directory and prefix, including at final resolution. Keep absent
projects distinct from invalid/foreign ones. Repair: [G-014](G-014-workspace-provenance.md).

### R2 — P2: A worktree deleted during inspection remains selectable

**Demonstrated bug.** `internal/versions/versions.go:169-174` compares path,
HEAD, and branch but ignores a newly `prunable` entry in the second inventory.
Removing the checkout directory leaves its registration, HEAD, and branch
unchanged, so the source stays valid. `workspace` trusts that result.

`inspect(..., between)` with `os.RemoveAll(feature)` returned `Complete == true`
and admitted the vanished source. A CLI probe wrapped Git, deleted the disposable
feature checkout immediately before the second `worktree list`, and delegated
to real Git. `workspace --source <previous-live-selector> --json` returned exit 0
with a project path that no longer existed. The deletion was visible in the
final Git inventory; this is not merely the acknowledged race after return.

Invalidate newly prunable/inaccessible/foreign identities; perform a final
selected-source check before returning a workspace. Preserve the honest
post-check race boundary. Repair: [G-014](G-014-workspace-provenance.md).

### R3 — P2: An update silently removes a comment outside the value

**Demonstrated data-preservation bug.** `internal/update/edit.go:242-246`
replaces everything from the key's colon through a later-line value.
The reader accepts this fragment:

```yaml
title: # retain this comment
  First
```

Updating `title=New` succeeded but produced `title: "New"`, deleting the comment.
The comment precedes the value's syntax span and is protected by G-009's
preservation contract. Replace the value's bytes while retaining the key-line
comment and surrounding standalone comments; exact-byte fixtures must cover
LF, CRLF, quoted values, lists, and explicit keys.
Repair: [G-015](G-015-preserve-updates.md).

### R4 — P2: Valid multi-field and explicit-key updates are refused

**Demonstrated functional gaps.** `internal/update/edit.go:248-269` plans each
flow removal's comma independently; removing the final two entries claims the
same separator twice and the overlap guard rejects the request. Given:

```yaml
{id: W-001, type: work, title: T, status: proposed, kind: fix, size: small}
```

`update W-001 --expect <current> --unset kind --unset size` fails with
`frontmatter: overlapping edits`. Both fields are optional and the combined
request is valid. Separately, the reader accepts `? status` followed by
`: proposed`, but `--set status=active` fails at `edit.go:298-299` with
`cannot locate the key's colon`. Both refusals leave bytes unchanged.

Explicit keys were acknowledged in the previous review but waived on the basis
that no accepted fixture used them. The actual reader accepts them; G-009 says
accepted-form refusals are implementation gaps. Do not narrow the reader to
silently ratify this limitation. Plan separators for the whole edit request
and locate explicit-key colon tokens safely.
Repair: [G-015](G-015-preserve-updates.md).

### R5 — P2: Configuration-only changes evade update's snapshot comparison

**Demonstrated contract gap.** `internal/update/update.go:336-338` compares the
parsed record directory and record sources but never `Project.Config`.
A `compare` fault hook changing `grove.yaml` from the usual two lines to the
same configuration with a comment let the update publish successfully.
The configuration bytes changed before the comparison, so G-009 requires refusal.
`internal/create/create.go:92` likewise compares only the parsed record root
although the creation contract requires refusing a changed configuration after
allocation; a second probe loaded the creation input, changed only a config comment,
and `New` still created a record instead of refusing. Compare exact
configuration bytes in both operations, retaining
consumed reservations on creation failure.
Repair: [G-015](G-015-preserve-updates.md).

### R6 — P2: Newlines in Git paths corrupt discovery and misplace write locks

**Demonstrated bug and unintended write.** `internal/repo/repo.go:41-48`
combines two `rev-parse` outputs and splits on newline, although paths can
contain that character. `internal/versions/versions.go:114-115` repeats the
same assumption for Git/common directories. In a repository named
`new\nline`, `versions` truncated the repository to `.../new`, set prefix to
`line/.git`, and marked its valid live source invalid. `update` succeeded but
created `.../new/grove/write.lock` outside the repository; the correct
`.../new\nline/.git/grove/write.lock` did not exist.

Obtain each path separately and remove exactly one output terminator; never
split path-bearing output on newline. Also replace the allocator's non-NUL
`worktree list --porcelain` parser (`internal/create/create.go:264`) with a
shared NUL-delimited inventory reader. An additional allocation probe put an untracked generated-format W-090 in
a newline-named linked checkout, with no counter state: allocation from main
returned 2 instead of the required 91, silently omitting that live source.
The old path tests cover a linked checkout's unusual name under a normal common
directory, leaving this case untested.
Repair: [G-016](G-016-git-paths.md).

### R7 — P3: Current documentation mixes historical proposals and integrated facts

**Demonstrated documentation drift, reconciled by this review's docs changes.**
G-009's Next still asked for its already-completed integration; G-002 said no
combined view was implemented; G-010/G-011 opened with proposed-only wording
despite later implementation evidence. Preserve historical review evidence and
G-002's accepted answer, but make current Next actions and implementation state
explicit. The record model's G-007 operational link also called creation proposed.

## Design concerns and limits, not additional demonstrated bugs

- `Resolve` re-inspects all branches/worktrees, and tree loading reads every
  Markdown blob under the project prefix. This multiplies refresh cost for an
  interactive view. No representative large-repository benchmark was run;
  measure before adding caching or asynchronous refresh. A future selected-source
  revalidation helper can improve correctness without promising a performance win.
- An update revision binds file content, not a previously selected branch.
  The current `workspace` contract explicitly returns a location, not write
  authority. A future cross-workspace edit or agent action must bind and recheck
  repository/worktree/configuration/record identity at execution time; do not
  claim today's `--expect` alone makes such an action safe.
- Existing tests cover cooperating same-revision writers, reciprocal dependency
  attempts, new/update serialization, linked-worktree locks, killed lock holders,
  and publication failures. They passed afresh. The existing same-revision and reciprocal-dependency
  tests use goroutines; additional pairs of actual CLI processes each produced
  exactly one success and one refusal, with the final project valid. Arbitrary
  editors or Git writes
  after the final check remain outside the serialization guarantee.
- Unknown HEAD comparisons, partial results, absent sources, deleted rows without
  selectors, dirty unrelated files, and committed/live distinction are exercised
  by passing tests. The prefix and prunable cases above qualify the overall claim.
- Windows locking, separate-clone coordination, filesystem crash durability, and
  all possible concurrent Git administration schedules remain unverified here.
  No TUI, worktree provisioning, editor/agent launch, claim, or run lifecycle exists.

## Recommended next investment

First make selecting an existing workspace trustworthy and make existing field
updates meet their preservation contract. The three proposed repair records
have linked plans and need no new product preference. Use one Fable worktree
serially in order G-014 → G-015 → G-016; the order is execution coordination,
not a semantic dependency between these independently testable fixes. Review
each unit and the combined diff before separately deciding integration.

After these repairs, the owner selected a terminal experience, then emphasized
early Kanban value. The resulting recommendation is a board for a clearly
labelled live checkout, with grouped cross-branch version inspection inside each
card and explicit existing-workspace selection. Other sources expose branch-only
work without an invented aggregate status. Incomplete results and stale-selection
refusals remain visible. [G-017](G-017-terminal-picker.md)
now owns the proposed detailed interaction and technical handoff. The terminal experience is selected and early Kanban value requested;
the board-first interaction, framework recommendation and detailed design are
proposals, not shipped behavior. Defer worktree creation to a separate bounded design:
it introduces filesystem/Git writes and destination/failure policies. Agent
launching, editing through the picker, claims, and run recovery follow later.
This keeps G-002 intact and avoids building actions on the defects found here.

## Shaping artifact verification

After creating the follow-up records through this checkout's `new` command and
setting their fields through `update`, `check` passed with 13 records. Local
Markdown path/anchor checks passed across the review, plans, operational records,
README, and direction/model documents. Product source files and go.mod/go.sum
remain unchanged; the retained G-010/G-011 checkout remains clean. The only
intentional coordination writes were the new-record reservations and CLI
mutation locks. No predecessor commands or sibling-project writes were used.
