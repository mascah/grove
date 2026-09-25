---
id: "G-165"
type: plan
title: "G-161 dependency view: layouts, shared preview model and CLI contract"
status: current
created: "2026-09-25T20:45:47Z"
updated: "2026-09-25T20:47:55Z"
work: ["G-161"]
---

## Inputs

Prepared on 2026-09-25 for [G-161](G-161-dependency-view.md) at revision
`sha256:705f8dda…`, headless, on branch `worktree-G-161` from `main`
`05892a2`. Read in full: G-161, [G-163](G-163-selected-work-review-boundary.md)
(which does not block this read-only work), the settled terms Work, Candidate,
Integration, Source and Revision, the G-116/G-117 layout-question precedent,
`internal/project/graph.go`, `internal/handoff/context.go`,
`internal/versions/current.go`, `internal/versions/changes.go`,
`internal/tui/detail.go` (`linked`), `internal/tui/model.go`
(`currentCards`), `internal/cli/versions.go`, `docs/board.md`,
`docs/commands.md`, the dependency rules of `docs/record-model.md` and step 5
of `docs/work-shaping.md`.

No board layout is approved. [G-166](G-166-g-161-dependency-layout.md) asks
the owner to choose one; the board steps below wait for its answer. The
shared model, the noninteractive command and the shaping-guide change do not
depend on it and are implemented first.

## Observed at `05892a2`

- 50 work records, 24 `depends_on` edges. The largest connected group has 13
  items, all done: G-025, G-037, G-038, G-039, G-040, G-041, G-042, G-043,
  G-044, G-045, G-046, G-052, G-065. Its longest chain is seven layers
  (G-037 → G-065 → G-038 → G-039 → G-040 → G-045 → G-046). The others are
  G-014/15/16 → G-017, G-010 → G-011, and G-150/G-151 → G-110 (G-110 is the
  only unfinished work with prerequisites; both are done).
- `context` orders a selection through transitive prerequisites, including
  paths through unselected work (`selection` in `internal/handoff`), and
  lists unselected prerequisites and blocking questions without adding them.
  That function is the ordering both the CLI and the board must share.
- The board's detail sidebar lists direct `needs`/`needed by` and
  `blocked by`/`blocks` both ways; nothing shows more than one hop.
- Candidate ancestry is read on demand with `git merge-base --is-ancestor`
  (`update`, `integrate`, `versions.ChangesContext`); the board never reads
  it during a load.

## Design (shared, layout-independent)

A new package `internal/deps` holds the one interpretation that the command
and the board both render.

- **Order.** `handoff.selection` moves there as `deps.Order` unchanged, and
  `context` calls it, so `context`, `deps` and the board cannot disagree.
  Only `depends_on` orders; `members`, `priority` and `blocks` never do.
- **Layer and group.** Among the work shown (the unfinished work, or the
  selection), a layer is one more than the deepest shown item it needs,
  directly or through any other work (0 when none), and shown items either
  of which needs the other share a group. A separate group is unrelated
  work; a done prerequisite two items share does not join them. Equal layers
  mean no declared order, never permission to run together.
- **Needs and unlocks.** Each item lists every direct prerequisite and the
  unfinished work it directly unlocks. Prerequisites that are not shown
  (done or abandoned ones in the overview) are listed after the rows with
  their delivery, collapsed on the board, never dropped.
- **Selection preview** for explicit IDs: the IDs as given, their order,
  every transitive prerequisite outside the selection ("not added"), the
  questions that block the selection or those prerequisites with their
  status, and notes. Nothing is added to the selection.
- **Delivery facts**, relative to a named checkout (its HEAD is the base) and
  the configured target, read only on demand (the command, or opening the
  board's preview), never during a board load, and only for prerequisites
  with a candidate:
  - proposed or active: awaiting implementation;
  - review: awaiting review, candidate `C` in the base or not;
  - done with a candidate: candidate `C` in the base or not, on the target or
    not;
  - done without a candidate: delivery unrecorded; establish it by ancestry
    or behaviour;
  - abandoned: will not be delivered; the dependency needs a decision;
  - a candidate Git cannot read here: unavailable, said as such.
- **Current-view notes** from `versions` (G-042), binding everything above
  to the one checkout's records: an involved record whose version here is
  older than the current one elsewhere (with that source and status),
  current states that diverge, current states elsewhere whose `depends_on`
  differs from this checkout's, uncommitted changes to it here, and an
  incomplete inspection. The preview never merges edges across sources.
- **Binding.** `Compare` takes the source a View was built from, never
  assuming the root checkout. The command's View is its checkout's. The
  board's overview may be built from the current view, as its columns are,
  with divergent records flagged, but a preview binds to one checkout (the
  one `b` chose, else this one) before it orders anything, so edges from
  different sources are never combined (G-161's Constraints). Review finding
  2 of the first gate asked for this.
- **Freshness.** Every involved record carries its revision (G-062). The
  board's preview is recomputed on every re-read and says when a selected
  record changed; any handoff it offers rereads before acting.

## Command contract (chosen here, layout-independent)

`grove [--project DIR] deps [WORK_ID...] [--json]`, read-only.

- Without IDs: the overview of unfinished work in the checkout, one row each,
  `GROUP LAYER ID STATUS NEEDS UNLOCKS TITLE`, ordered by group, layer and
  ID, then the notes (open blocking questions, abandoned or undelivered
  prerequisites, current-view notes).
- With IDs: the selection preview: `Selected:` as given, `Order:`, a table
  `ID STATUS ROLE DELIVERY` for the selection and every outside prerequisite,
  then questions and notes, ending with the reminder that the preview adds no
  work and authorizes nothing.
- The checkout (root, ref, HEAD) and target head the output. `--json` prints
  the same model with revisions. An incomplete inspection prints everything,
  says so on stderr and exits 1, as `versions` does; an unknown or non-work
  ID fails as `context` fails.

## Board layouts (owner's choice, G-166)

Both drawings use the synthetic unfinished backlog below; the real 13-item
group appears only when history is expanded, since all of it is done. Titles
are clipped with `…`; nothing else is. `S-` IDs are synthetic.

| ID | Status | Needs | Title |
| --- | --- | --- | --- |
| S-01 | done (candidate on main) | | Scaffold the ledger project and its test harness |
| S-02 | proposed | S-01 | Store accounts, balances and currencies in one shared ledger file |
| S-03 | active | S-02 | Import bank statements from CSV and OFX exports into the ledger |
| S-04 | proposed | S-02 | Categorize transactions with editable, ordered matching rules |
| S-05 | proposed | S-03 S-04 | Monthly budget report with category totals, carry-over and warnings |
| S-06 | proposed | S-05 | Export the monthly report as CSV and printable HTML |
| S-07 | review (candidate not on main) | S-03 | Reconcile imported balances against statement closing balances |
| S-08 | proposed | S-04 S-10 | Detect recurring transactions and predict next month's bills |
| S-09 | proposed | | Dark theme for the report viewer |
| S-10 | abandoned | | Legacy OFX 1.x parser |
| S-11 | active, blocked by open S-12 | | Fix rounding of foreign-currency amounts on import |
| S-13 | proposed | S-05 | Send budget alerts when a category passes its limit |

### A: lanes

Layers are columns, prerequisites to the left; each card names what it needs
(`←`) and unlocks (`→`). 120×30:

```text
Dependencies · current view · target main        10 unfinished · 1 done, 1 abandoned collapsed (h)
 Layer 0                     Layer 1                     Layer 2                     Layer 3
┏━━━━━━━━━━━━━━━━━━━━━━━━━┓ ┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
┃▶ S-02 proposed     grp 1┃ │ S-03 active       grp 1 │ │ S-05 proposed     grp 1 │ │ S-06 proposed     grp 1 │
┃Store accounts, balances…┃ │Import bank statements f…│ │Monthly budget report wi…│ │Export the monthly repor…│
┃← ✓S-01   → S-03 S-04    ┃ │← S-02   → S-05 S-07     │ │← S-03 S-04  → S-06 S-13 │ │← S-05                   │
┗━━━━━━━━━━━━━━━━━━━━━━━━━┛ └─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘
┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐ ┌─────────────────────────┐
│ S-09 proposed     grp 2 │ │ S-04 proposed     grp 1 │ │ S-07 review       grp 1 │ │ S-13 proposed     grp 1 │
│Dark theme for the repor…│ │Categorize transactions …│ │Reconcile imported balan…│ │Send budget alerts when …│
│                         │ │← S-02   → S-05 S-08     │ │← S-03  not on main      │ │← S-05                   │
└─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘ └─────────────────────────┘
┌─────────────────────────┐                             ┌─────────────────────────┐
│ S-11 active       grp 3 │                             │ S-08 proposed     grp 1 │
│Fix rounding of foreign-…│                             │Detect recurring transac…│
│? S-12 open              │                             │← S-04 ✗S-10 abandoned   │
└─────────────────────────┘                             └─────────────────────────┘
 ← prerequisite · → unlocks · ✓ done · ✗ abandoned · ? blocking question
 Space select · p preview (0 selected) · Enter open · h history · ←→ layer · Esc back
```

At 80 columns two lanes show at a time and ←/→ scroll by layer
(`Layers 0–1 of 4 ▸`); titles clip at 34 characters. A seven-layer group (the
real 13 items expanded) needs four screens across. Lanes show breadth and
convergence at a glance, but a card's prerequisites may be off screen and the
lines between cards are IDs, not drawn arrows.

### B: layered list with a focus tree (recommended)

One row per unfinished item, grouped and indented by layer; the focused
item's upstream and downstream trees fill the right pane from 100 columns, or
replace the list with Tab below that. 120×30:

```text
Dependencies · current view · target main        10 unfinished · 1 done, 1 abandoned collapsed (h)
 Group 1 · 8 unfinished                             S-05 proposed · layer 2 · group 1
   S-02 proposed  Store accounts, balances and c…   Monthly budget report with category totals,
     S-03 active    Import bank statements from…    carry-over and warnings
     S-04 proposed  Categorize transactions wit…
▶      S-05 proposed  Monthly budget report wit…    Needs
       S-07 review    Reconcile imported balanc…    S-05
       S-08 proposed  Detect recurring transac…     ├─ S-03 active     Import bank statements f…
         S-06 proposed  Export the monthly rep…     │  └─ S-02 proposed Store accounts, balanc…
         S-13 proposed  Send budget alerts whe…     │     └─ ✓ S-01 done
 Group 2 · 1 unfinished                             └─ S-04 proposed   Categorize transactions…
   S-09 proposed  Dark theme for the report view…      └─ S-02 (shown above)
 Group 3 · 1 unfinished
   S-11 active    Fix rounding of foreign-curre…    Unlocks
                                                    S-05
                                                    ├─ S-06 proposed   Export the monthly repo…
                                                    └─ S-13 proposed   Send budget alerts when…

                                                    Same layer, no declared order: S-07 S-08
 ← needs · → unlocks · ✓ done · ✗ abandoned · ? blocking question · indent = layer
 Space select · p preview (0 selected) · Enter open · h history · Tab tree · Esc back
```

At 80×24 the list alone shows, titles clipped at 44 characters, and Tab
swaps to the focused tree. Any size of project reads the same way: the list
scrolls, and each tree shows only the focused item's ancestry, with a
repeated node written `(shown above)`.

### Preview (both layouts)

`p` with S-05 and S-03 selected (Space marks `●` in the list):

```text
Selection preview · bound to the current view, main 05892a2 · target main
Selected: S-05 S-03 (as marked)
Order:    S-03 → S-05   (S-05 needs S-03)

 In order   ID    Status    Needs              Delivery
 1          S-03  active    S-02               awaiting implementation
 2          S-05  proposed  S-03 S-04          awaiting implementation
 Outside the selection, not added
            S-02  proposed  S-01               awaiting implementation · needed by S-03 S-05
            S-04  proposed  S-02               awaiting implementation · needed by S-05
            S-01  done                         candidate 9c0ffee in main 05892a2, on main
 Questions: none block the selection or its prerequisites.

 · S-02 and S-04 are unfinished and not selected: the selection cannot finish until they are delivered.
 · Order comes from depends_on only. It does not say that unordered work can run in parallel.
 · A preview adds no work, starts nothing, and authorizes nothing.
 Enter print `grove context S-03 S-05` and exit · Esc back
```

## Steps

1. `internal/deps`: move `selection` as `Order`; add `Build` (layers,
   groups, needs, unlocks, preview), `Deliver` (on-demand ancestry, through
   `repo`) and `Compare` (current-view notes from a `versions.Result`), with
   table tests: chain, shared prerequisite, convergence, unrelated work,
   priority and members not changing order, abandoned and historical done
   prerequisites, blocking questions, divergent edges, uncommitted and older
   records, incomplete inspection. `context` calls `deps.Order`.
2. `grove deps`: usage, text and `--json`, CLI tests over a fixture
   repository including a candidate in review off the target.
3. Docs: `docs/commands.md` gains `Deps`; `docs/work-shaping.md` step 5 says
   to record real prerequisites with their reasons in the dependent work and
   keep preferred order as order, never as an edge.
4. **Waits for G-166.** The board view in the chosen layout from the board
   (`g`), with Space selection, the preview, Enter into a record's detail and
   back, history collapse, re-read and resize; terminal checks in
   `internal/tui/testdata/terminal.py`; `docs/board.md`.
5. Owner judgment at 80 and 120 columns on the real 13-item group (expanded)
   and the synthetic backlog, then review and handoff.
