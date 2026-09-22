---
id: "G-042"
type: work
title: "Derive a project-wide current view of work"
status: active
created: "2026-09-21T00:54:15Z"
updated: "2026-09-22T22:40:21Z"
kind: feature
size: medium
priority: 3
depends_on: ["G-037", "G-038"]
relates_to: ["G-035", "G-002", "G-010", "G-011", "G-030", "G-031", "G-064", "G-065"]
formerly: "W-024"
candidate: "046150e"
---

## Outcome

Opening Grove from any linked checkout presents the same useful project-wide
current work view, with source observations and genuine divergence accessible
without making every user choose a branch first.

## Selected direction

Collapse identical observations; move demonstrably superseded states into
history; expose unintegrated progress and label live uncommitted changes.
Keep genuine competing changes visible. Never use the largest status, latest
timestamp or newest branch tip as authority. A revert is a real change, not
automatically an older state. Exact source routing and stale-selection checks
remain intact; view selection does not merge content or grant write authority.

Use G-065's stable identities and recursive discovery. Neither type prefixes
nor a completed/history folder establish which record is current or integrated.
"History" here is a view of evidence, not a filesystem move. Include supported
old/new schema sources and legacy/new IDs in the projection fixtures.

## Acceptance

1. Specify and fixture-test a deterministic projection for old branches,
   integrated work, unmerged progress, dirty/deleted records, branch-only work,
   genuine divergence, reverts and missing/invalid sources.
2. The policy distinguishes record history from unrelated branch-tip commits
   and committed observations from live overlays; ambiguity is explicit.
3. Equivalent observations from different invoking checkouts produce the same
   default view. Users can inspect the chosen evidence and other sources.
4. For the observed G-030 history, an old Proposed copy on G-023 does not
   obscure the Done/integrated result on main. Do not hard-code these branches.
5. Preserve read-only behavior, incomplete-result diagnostics, safe action
   targeting and G-031's batched Git reads. Avoid per-branch Git processes;
   any added history work must retain acceptable measured board load behavior.
6. Update G-002's implementation note and current documentation. Full visual
   redesign is G-043; projection should be independently testable and explainable.

## Evidence

Implementation is on branch `worktree-G-042` from base `939d090` (main). It
started from this record and [plan G-093](G-093-current-view-plan.md) as
committed at `2b2831a`. The candidate is the commit that adds this section;
the commit after it only sets `status=review`.

**Decisions.** The owner decided the following on 2026-09-22; the plan records
them with their words:

- Uncommitted edits count, labelled.
- A divergence is one card, placed in the earliest status among its current
  states.
- There is no integration target, so nothing is called unintegrated.

Routine technical choices, made in the implementation:

- Observations are ordered by the record's bytes at the two commits' merge
  base: the older one is the one whose bytes the base holds.
- Merge bases come from Git's paint-down-to-common walk, read through the
  existing `cat-file` process.
- A cycle, which reverts carried across merges can produce, is decided by its
  components, with a note.
- Schema 3 is the only schema, so a pre-G-052 branch is an invalid source
  rather than a legacy observation.

**Acceptance.**

1. `internal/versions/current.go` holds the projection. `TestCurrentView` covers:
   - a stale branch (the G-030/G-023 shape) and a merged branch whose record
     main moved on;
   - unmerged progress, and work that exists only on a branch;
   - committed and uncommitted deletions, and an uncommitted edit and addition;
   - divergence, and a revert to earlier bytes (current, while a later-looking
     branch copy is older);
   - a detached checkout, and an old-schema branch as an invalid source.

   `TestCurrentViewUnorderedPair` covers an unreadable base.
   `TestCurrentViewCycle` covers cycles. `TestMergeBases` compares the walk
   with `git merge-base --all` on a criss-cross history.
2. Only the record's bytes at merge bases matter, so commits touching other
   files change nothing (the stale branch). An uncommitted edit is its own
   observation on top of HEAD, labelled `uncommitted`. A pair that cannot be
   ordered stays current with a note, shown in `versions` (stderr and
   `notes`) and in the card's "Could not order" lines.
3. The same sources give the same view: `TestCurrentView` inspects from two
   checkouts, `TestCurrentViewBoard` from two invoking Git directories.
   Evidence stays inspectable:
   - `versions` has a `CURRENT` column, and its JSON adds `current`, `older`,
     and `notes`;
   - a card lists current rows first, each older row with its reason;
   - `b` still opens one checkout's own board;
   - selection is unchanged.
4. In the G-030/G-023 shape, the stale branch's Proposed copy is older and
   main's Done is current. The code names no branch.
5. The code writes nothing and still reports incomplete sources. A committed
   deletion row has no selector and refuses. `TestCommittedReadIsScopedAndShared`
   asserts one `cat-file` process and no `merge-base` or `rev-list` process.
6. Updated: G-002's disposition and Next, README (board, `versions`),
   AGENTS.md, the brief's current-view paragraph (direction only: ancestry
   means merge bases, and no target yet), `docs/record-model.md`, and the
   usage text. The visual redesign is left to G-043.

**Measured load**, `versions` wall time on this Mac, with 3 runs each for the
synthetic repositories:

| Repository | Before (`939d090`) | After |
| --- | --- | --- |
| This repository (4 sources) | 0.07 s | 0.06 s |
| 300 branches, no record edits | 0.22 s | 0.29 s |
| 300 branches forked along 200 edits of one record on main | 0.50 s | 0.79 s |
| 1,000 branches each editing one record differently | 1.87 s | 3.0 s |

**Verification.** At `84115c4`, whose code the candidate shares (only records
follow), these all passed:

- `go vet ./...`
- `gofmt -l .` (no output)
- `go run ./cmd/grove check` (OK: 90 records)
- `go test -count=1 -timeout 120s ./...`, including the pseudo-terminal
  lifecycle test `TestTerminal`

The versions package takes 4.0 to 4.3 s alone under `-short`; under
whole-suite load it takes 7.2 s, as main's does (7.3 s).

**Review.** [G-094](G-094-current-view-review.md): three rounds, with five
round-1 defects and one round-2 remainder, all fixed with regression tests.
Round 3 found none. It examined `96900e0`; `1c328cf` (a test made parallel) and
`84115c4` (usage text) followed.

**Limits.**

- A record that diverges n ways costs n² merge-base walks (marked
  `ponytail:`). A cycle costs the full relation.
- Which newer place a reason names can depend on commit dates. It is the
  same from every checkout.
- In a shallow clone, a pair whose history crosses the boundary is noted as
  unordered.
- `independent()` walks full history under criss-cross bases.
- Without a target the view cannot say "not yet on main". Ask for that as
  its own work when wanted.

**When a card shows `⑂ N states`.** Divergence means two places each changed
the record since they last shared a commit, into different bytes:

- **Two sessions on one record.** Worktree A sets G-050 to `active` with its
  Next, and worktree B, started from the same main, does too with a different
  Next. Two active states: the card stays in Active with `⑂ 2 states`.
- **Main edited while a branch worked.** A work branch moves G-050 to
  `review`. Meanwhile someone fixes a typo in G-050 on main. Main's `active`
  and the branch's `review` both changed since the split, so the card sits in
  Active, the earlier status, marked `⑂ 2 states`, until the branch merges
  main or main merges the branch.
- **Not divergence.** In each of these, one state is current and the others
  are older rows under the card: a branch that changed the record while main
  did not, a stale branch that never touched it, or a checkout's uncommitted
  edit on top of its own HEAD.

## Next

**Owner feedback on candidate `046150e`, 2026-09-22 (reviewing its
Decisions):** "I think I want to revisit the idea of an integration target.
In my case that's main. During initial scoping I was asked if this should be a
grove.yaml setting and I think it probably should be." Asked where it lands,
the owner chose to reopen G-042 rather than merge first. They deferred to
separate work having `update` refuse `done` off the target.

Implement an optional `target` key in `grove.yaml`, set to `main` here. The
current view labels each current state on or not on the target, following
the plan's revised decision. The candidate above and review G-094 remain as
evidence of the earlier attempt.
