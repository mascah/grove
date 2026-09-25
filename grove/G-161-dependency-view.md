---
id: "G-161"
type: work
title: "See work dependencies and preview a selected assignment"
status: review
created: "2026-09-25T20:35:15Z"
updated: "2026-09-25T22:54:10Z"
kind: feature
relates_to: ["G-035", "G-042", "G-043", "G-023", "G-047", "G-054", "G-060", "G-162"]
candidate: "83e7f3823f3444832b5c9866300e201f06691674"
approved: "83e7f3823f3444832b5c9866300e201f06691674"
---

## Outcome

After several shaping sessions, the owner can return to Grove, see the
prerequisite structure across unfinished work, and choose what to do next or
assign together without asking an agent to reconstruct the intended order.

Owner intent, conversation 2026-09-25: uncertainty about dependencies among
roughly ten proposed items is preventing adoption in other hobby projects.
The owner suggested a visual DAG and recalled selecting several items for the
predecessor's work agent, while wanting more human control over that selection.
The owner then requested shaping this exploration. The design below remains
proposed; no layout or execution policy was selected.

## Constraints

### Observed evidence

At main `3f2b923`, the [record model](../docs/record-model.md) already defines
`depends_on`, `members`, question `blocks`, and priority separately.
[Validation](../internal/project/graph.go) rejects dependency and membership
cycles independently. [Context](../internal/handoff/context.go) orders selected
work by transitive prerequisites, including paths through unselected work,
without adding prerequisites to the assignment. The [detail
sidebar](../internal/tui/detail.go) exposes direct relationships in both
directions, but neither the board nor the CLI provides a dependency overview
or an interactive selection preview.

The main checkout contains 50 work records and 24 dependency edges; its largest
connected component contains 13 items. Only four records there are unfinished,
so this is layout evidence, not a representative ten-item proposed backlog.
`versions --json` also shows active G-154 on its work branch while main still
holds its proposed copy: the current view matters here. The [adoption
roadmap](G-047-adoption-roadmap-plan.md) names dependency visualization as a
later candidate; no existing proposal owns it.

### Proposed design and scope

- Add a focused dependency view to the board. Start with unfinished work;
  make a selected item's upstream prerequisites and downstream dependants
  easy to follow, with completed prerequisite history collapsed and
  expandable. Keep dependencies outside the visible filter discoverable and
  counted rather than making a truncated graph appear complete.
- Explore compact terminal lanes or layers with a readable list fallback.
  Arrows run from prerequisite to dependent work, with a visible legend.
  The owner judges concrete layouts before one is settled. A browser view
  or export remains an alternative if actual terminal readability warrants
  it, not a required second interface.
- Let the person select explicit work IDs and preview their dependency order,
  prerequisites outside the selection, and blocking questions. Expose the
  same interpretation through noninteractive inspection; choose the command
  and output contract during preparation. Browsing or selecting starts no
  attempt, adds no prerequisite to the assignment, and changes no status.
- Explain recorded facts relative to a named checkout/base and, when
  configured, the integration target: awaiting implementation, awaiting
  review, candidate present in the base, or candidate integrated. Status
  alone does not establish delivery. Old done records without candidates
  and unavailable evidence retain their uncertainty. The preview does not
  invent a requirement to merge between selected items; that execution
  choice belongs to [G-163](G-163-selected-work-review-boundary.md).
- Keep membership, preferred sequence, and actual prerequisites distinct.
  Equal graph depth means no declared ordering between those items, not
  permission or evidence for parallel writes. An absent dependency field
  means no declared prerequisites, not proven implementation readiness.
- Extend the [shaping workflow](../docs/work-shaping.md) at the relevant
  handoff: capture real prerequisite edges and explain their reasons in
  their owning work records; preserve preferred order as such. Do not
  manufacture an edge to encode priority or duplicate relationships in a
  separately editable graph file. Reconcile the owning command/board docs
  when behavior ships.

Reuse [G-042's current view](G-042-current-view.md), including visible
divergence, and the existing checkout selector. Conflicting branch versions
must not be silently combined into a supposedly authoritative DAG or order;
expose the ambiguity and bind a preview to a coherent source before acting.
Preserve exact source targeting, freshness, text escaping, batched Git reads,
and on-demand history. This work neither launches batches nor edits graph
edges by gesture. New merge-gate fields and an editable scheduling system
are outside this initial proposal.

[G-162](G-162-bounded-work-selection.md) owns managed execution of several
selected items. Graph-first is the assistant's suggested investment order,
not a technical prerequisite between the two proposals.

## Acceptance

1. From a representative ten-item backlog, the owner can identify a chain,
   a shared prerequisite, a convergence point, unrelated work, and the work
   a selected item unlocks without opening every record or asking an agent.
2. A selection preview preserves the explicit selection, orders its members
   correctly through transitive dependencies, and identifies outside
   prerequisites and relevant open questions without silently adding work.
   It explains facts and uncertainty rather than claiming authorization,
   safe parallelism, or a guaranteed completion time.
3. Review candidates not integrated into the intended base, historical done
   records without candidates, abandoned prerequisites, hidden dependencies,
   uncommitted records, divergent edges, and incomplete source reads each
   produce an honest, actionable explanation. A changed source invalidates
   any stale selection handoff.
4. Membership and priority cannot alter dependency order. A same-level pair
   is not labelled safe to run concurrently; a blocking question remains
   distinct from a work-to-work edge.
5. The owner judges 80- and 120-column layouts using the real 13-item
   component and a synthetic unfinished backlog with long titles, branches,
   convergence, and disconnected items. A larger case demonstrates bounded
   navigation rather than shrinking the whole project into illegibility.
6. Connected terminal checks cover board to graph, focus and selection,
   preview, record detail, return, refresh, resize, and exit. CLI and TUI
   agree on the same source's relationships and ordering. Loading retains
   the repository's Git-read and responsiveness constraints.
7. A fresh shaping-session handoff leaves real dependencies and their reasons
   discoverable in records, with preferred sequence identified separately;
   the resulting view is useful before multi-item execution exists.

## Next

**Handoff into review, 2026-09-25, headless `/grove-work G-161`.** The
candidate is the commit holding this evidence, and `candidate` names it; the
next commit sets only `status=review`. Awaiting the owner's judgment, including acceptance item 5.

- **Branch and inputs.** Branch `worktree-G-161` in
  `.claude/worktrees/worktree-G-161`, based on `main` `05892a2`. This
  session resumed at `7e1acdb`, after [G-166](G-166-g-161-dependency-layout.md)
  was resolved. It started from G-161 `sha256:ac898c03…` and plan
  [G-165](G-165-g-161-dependency-view-plan.md) `sha256:f1062b6d…`. The plan
  is reconciled with what was built: the Binding and Freshness bullets, step
  4, and a "B as built" drawing.
- **What was built.**
  - The first session (`c19e945`..`c92d16f`, review
    [G-168](G-168-g-161-deps-review.md)): `internal/deps`, `grove deps`, and
    the shaping guide's step 5.
  - This session (`476b5ac`, `ffaca01`, `867a9d8`): `g` on the board opens
    the dependency view in layout B (`internal/tui/deps.go`, documented in
    `docs/board.md#dependencies`). The list shows connected groups, then
    unconnected work, with every layer indented. Each row carries its
    status, `? question` and the card's tags (`⑂`, `uncommitted`,
    `not on main`, running). The focused row's `← Needs` and `→ Unlocks`
    trees mark done `✓` and abandoned `✗`, write a repeat as
    `(shown above)`, keep blocking questions apart from work, and list each
    diverging state's own edges. From 100 columns the trees stand beside the
    list; Tab gives them focus to scroll, and below 100 columns shows one
    pane.
  - Keys: Space selects and `c` clears. `p` previews, bound to one checkout
    (the board's, else the one Grove started in), with delivery read from
    Git on demand in that checkout. Enter opens a record, and Esc returns
    here. `h` shows every work, and `b` chooses a checkout and returns here.
- **Decisions.**
  - Layout B and no printed handoff are the owner's G-166 answer. It is
    linked from the plan and needs no decision record, following G-117's
    precedent.
  - Divergent overview rows use the state the board places the card by, and
    the trees show each state's edges. This replaced the plan's per-item
    `Compare` after review.
  - Each preview is recomputed on every re-read and names the records that
    changed. With no handoff, nothing acts on a stale preview.
- **Acceptance.**
  1. Tests show a chain, a shared prerequisite, a convergence, unrelated
     work and unlocks in `TestDepsListAndTree`, on the synthetic ten-item
     backlog. The owner's reading is item 5's judgment.
  2. Met. `TestDepsPreview` covers the marked order, the transitive order,
     outside prerequisites "not added", and the questions and notes.
  3. Met. The cases and where each shows:
     - A review candidate off main: `not on main` and "candidate … not in
       HEAD".
     - Done without a candidate: its delivery is unrecorded.
     - An abandoned prerequisite: `✗` and a note.
     - Hidden prerequisites: counted in the heading.
     - Uncommitted records: the tag, and a preview note.
     - Divergent edges: `⑂` with each state's edges.
     - An incomplete read: the banner and a preview note.
     - A changed source: recomputed and named.
  4. Met. Only `depends_on` orders. The legend says equal layers have no
     declared order, and questions stay apart from work.
  5. **Pending: the owner's judgment.** The plan's "B as built" drawing
     shows 120 columns on the synthetic backlog. Run `go run ./cmd/grove` in
     this checkout, press `g`, and resize to 80 and to 120. Press `h` for
     the real 13-item group (seven layers) and `Tab` for the trees. Paging
     the list and scrolling the trees show a larger case.
  6. Met.
     - The `dependencies` scenario in `terminal.py` drives board, graph,
       focus, selection, preview, record detail, return, refresh, resize and
       exit.
     - `TestDepsPreview` checks that the board's order equals
       `deps.Preview`'s.
     - The list asks Git nothing, and the preview reads through
       `Model.read`.
  7. The shaping half was met in the first session. The view needs no
     multi-item execution.
- **Verification** at `867a9d8`:
  - `gofmt -l .` printed nothing.
  - `go vet ./...` was clean.
  - `grove check` printed "OK: 158 records".
  - `go test -count=1 -timeout 120s ./...` passed every package. `tui` took
    15.4s, since `TestTerminal` runs all 12 `terminal.py` scenarios,
    `dependencies` included.
- **Reviews.**
  - [G-168](G-168-g-161-deps-review.md) is the first gate.
  - [G-175](G-175-g-161-deps-second-review-gate-on.md) is the board gate
    and final combined candidate, examined at `ffaca01`. Its round 1 had one
    high, four medium and three low findings plus a docs gap, fixed with
    regressions in `ffaca01`; one low was documented. Its round 2 found
    nothing consequential and two low findings. One was fixed in `867a9d8`,
    self-checked with a regression. The other is the existing chooser `s`
    trap.
- **Limits.**
  - `A`, then `o`, from a record opened here returns to the board, not to
    the dependency view (one return slot).
  - The chooser's `s` then Esc trap predates this work.
  - The trees' scroll is clamped on the next key after a resize.
  - `-race` was not run, since this is not a concurrency change beyond
    one more read kind.
- **Integrator's next action**, after the owner's judgment:
  - In this checkout: `grove approve G-161 "VERDICT"`, or
    `grove feedback G-161 "TEXT"`.
  - In the target's checkout: `grove integrate G-161`.

Verdict on candidate 83e7f38, 2026-09-25: approved
