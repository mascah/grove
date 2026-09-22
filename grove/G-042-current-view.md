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
committed at `2b2831a`. The first candidate, `046150e`, went to the owner
for review. The owner reopened the work to add an integration target (see
Next); the plan's revised decision was committed at `070671d`. The current
candidate is the commit that writes this revision of the section. The commit
after it only sets `status=review`.

**Decisions.** The owner decided the following on 2026-09-22; the plan records
them with their words:

- Uncommitted edits count, labelled.
- A divergence is one card, placed in the earliest status among its current
  states.
- ~~There is no integration target, so nothing is called unintegrated.~~
  Revised the same day, while reviewing `046150e`. `grove.yaml` names the
  target (`target: main` here), and the target labels states without
  deciding anything. Having `update` refuse `done` off the target is to be
  separate work.

Routine technical choices, made in the implementation:

- Observations are ordered by the record's bytes at the two commits' merge
  base: the older one is the one whose bytes the base holds.
- Merge bases come from Git's paint-down-to-common walk, read through the
  existing `cat-file` process.
- A cycle, which reverts carried across merges can produce, is decided by its
  components, with a note.
- Schema 3 is the only schema, so a pre-G-052 branch is an invalid source
  rather than a legacy observation.
- The target is the branch that every valid `grove.yaml` naming one agrees
  on. This keeps it independent of the invoking checkout, and it works on a
  branch before that branch merges. A target branch with no project yet lacks
  every record.

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
   `TestCurrentViewTarget` and `TestCurrentViewTargetBeforeAdoption` cover
   the target: named before merging, conflicting, missing, deletions, and
   adoption.
2. Only the record's bytes at merge bases matter, so commits touching other
   files change nothing (the stale branch). An uncommitted edit is its own
   observation on top of HEAD, labelled `uncommitted`. A pair that cannot be
   ordered stays current with a note, shown in `versions` (stderr and
   `notes`) and in the card's "Could not order" lines. With a target, the
   view also exposes unintegrated progress:
   - a card none of whose committed current states is on the target is
     tagged `not on main`;
   - each state in its details says whether it is on main;
   - `versions` shows a `TARGET` column and JSON `on_target`.

   This is tested in `TestCurrentViewTarget` (TUI) and
   `TestVersionsTargetCLI`.
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
6. Updated:
   - G-002's disposition and Next;
   - README (board, `versions`, target);
   - AGENTS.md;
   - the brief's current-view paragraph (direction only: ancestry means merge
     bases, and the target is optional configuration that labels);
   - `docs/record-model.md` (the `target` key);
   - the usage text;
   - this repository's `grove.yaml`, which now has `target: main`.

   The visual redesign is left to G-043.

**Measured load**, `versions` wall time on this Mac, with 3 runs each for the
synthetic repositories:

| Repository | Before (`939d090`) | After |
| --- | --- | --- |
| This repository (4 sources) | 0.07 s | 0.06 s |
| 300 branches, no record edits | 0.22 s | 0.29 s |
| 300 branches forked along 200 edits of one record on main | 0.50 s | 0.79 s |
| 1,000 branches each editing one record differently | 1.87 s | 3.0 s |

Re-timed after the target change, at `00ddefe`: 0.29 s and 0.78 to 0.81 s
for the two 300-branch cases, and 3.1 s for the 1,000-branch case, which is
unchanged. Resolving the target costs one pass over the sources.

**Verification.** At `3851182`, whose code the candidate shares (only records
follow), these all passed:

- `go vet ./...`
- `gofmt -l .` (no output)
- `go run ./cmd/grove check` (OK: 90 records; 91 with G-095 added after)
- `go test -count=1 -timeout 120s ./...`, including the pseudo-terminal
  lifecycle test `TestTerminal`

The versions package takes 4.1 to 4.4 s alone under `-short`. Under
whole-suite load it takes 8.1 s, against 7.3 s for main's.

**Review.**

- [G-094](G-094-current-view-review.md), on the first candidate: three
  rounds, with five round-1 defects and one round-2 remainder, all fixed with
  regression tests. Round 3 found none. It examined `96900e0`; `1c328cf` (a
  test made parallel) and `84115c4` (usage text) followed.
- [G-095](G-095-integration-target-review.md), on the target: two rounds.
  Round 1 found one defect, adoption with a target that has no project yet,
  which is fixed with its test. Round 2 found none. It examined `3851182`.

**Limits.**

- A record that diverges n ways costs n² merge-base walks (marked
  `ponytail:`). A cycle costs the full relation.
- Which newer place a reason names can depend on commit dates. It is the
  same from every checkout.
- In a shallow clone, a pair whose history crosses the boundary is noted as
  unordered.
- `independent()` walks full history under criss-cross bases.
- On the target means the target's tip holds the same bytes. A stale target
  can therefore "hold" a state that a revert restored elsewhere.
- Every build older than this branch refuses the `target` key. After the
  merge, `~/.local/bin/grove` rejects this repository's `grove.yaml`, and
  treats each branch carrying the key as an invalid source, until it is
  rebuilt (see Next).
- `init` does not write `target`; add it by hand.

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

**What the target adds.** A work branch moves G-050 to `review`, and main
still says `active`. The card shows Review, tagged `not on main`, until the
merge. In the two-sessions example above, neither active state is on main:
`⑂ 2 states not on main`. In the typo example, main holds one side, so there
is no tag; the details say which state is on main.

## Next

**Owner feedback on candidate `046150e`, 2026-09-22 (reviewing its
Decisions):** "I think I want to revisit the idea of an integration target.
In my case that's main. During initial scoping I was asked if this should be a
grove.yaml setting and I think it probably should be." Asked where it lands,
the owner chose to reopen G-042 rather than merge first. They deferred to
separate work having `update` refuse `done` off the target. The target
itself is now implemented, as recorded above.

Judge the candidate (the owner). Demo from the branch:

```sh
cd /Users/mascah/GitHub/mascah/grove/.claude/worktrees/G-042
go run ./cmd/grove            # "Board: current view, target main"; this record is tagged not on main
go run ./cmd/grove versions   # the CURRENT and TARGET columns; stderr names the target
```

On approval, in the main checkout, after quoting the verdict under Evidence
in this record:

```sh
cd /Users/mascah/GitHub/mascah/grove
git merge worktree-G-042
go run ./cmd/grove update G-042 --set status=done --commit
just install                  # older builds refuse the new target key
```

On feedback, write it here and set `status=active` on the branch. The
done-on-target check in `update` is open for shaping as its own work.
