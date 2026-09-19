---
id: "W-009"
type: work
title: "Browse a terminal Kanban board with explicit record versions"
status: proposed
created: "2026-09-19T20:23:01Z"
updated: "2026-09-19T20:31:23Z"
kind: feature
priority: 2
size: medium
depends_on: ["W-006", "W-007", "W-008"]
relates_to: ["Q-001", "W-004", "W-005"]
---

## Outcome

See work in a terminal Kanban board, inspect differing committed/live versions
inside a card, and explicitly select an existing workspace without copying an
opaque selector. On 2026-09-19 the owner preferred the board over a standalone
version picker as the first terminal experience after reliability fixes.
Board-first is selected direction; the detailed layout, context policy and
framework choice below remain proposed design, not shipped behavior.
[Implementation plan](../../docs/plans/W-009-terminal-picker.md).

## What the TUI means here

The terminal becomes a persistent keyboard-driven board while the command runs.
Left/Right move between status columns; Up/Down move between cards. Enter opens a
card's versions and source details. Leaving restores the ordinary shell screen.
Grove still reads local files without a server, and the board stores no duplicate
state. The first workspace action resolves and returns a location; it does not
change the parent shell's directory or launch an editor.

Version selection serves a specific purpose inside a card: main may still say
proposed while Fable's checkout says done. Showing both lets the person inspect
the difference and choose the intended workspace. It is supporting navigation;
the board is the starting screen and everyday overview.

Proposed board (illustrative data, not current project status):

```text
Grove  /project    Board: live main    4 sources inspected
Proposed           Active             Done               Abandoned
> W-006            W-003              W-001
  Workspace safety Record updates     Inspect records
  4 versions       4 versions         4 versions
  W-007                               W-002
  Preserve updates                    Create records

Other sources: W-010 [2 versions]   (absent from this checkout)
Left/Right columns   Up/Down cards   Enter versions   b checkout   r refresh
```

The label identifies one live checkout as the board context, initially the
invocation's checkout. Only that context's live work status determines columns.
Its title labels its card. This is an explicitly scoped view, not an aggregate
status across branches. Opening the card shows all observed versions grouped by
ID, each with its own title, status and source. A committed `done` version on a
feature branch can coexist with main's `active` card. Nothing infers integration.

A separate Other sources shelf lists grouped work IDs absent from the selected
live checkout but present in another valid source. It has no lifecycle status;
Enter opens the same version view. This keeps newly created branch-only work
visible from main without assigning it a made-up main status. Questions and
decisions remain available through the CLI; they do not gain work-status columns.

## Command and interaction proposal

`grove board [WORK_ID] [--json]` uses existing project discovery and Git scope.
Show only work groups; an optional work ID restricts the cards/shelf but still
requires a concrete version choice. Reject a non-work ID as usage error.

- Left/Right (h/l) switch columns; Up/Down (j/k) move cards. At widths below
  100 columns show one status column at a time with all four status tabs and
  counts; below 40 columns or 10 rows show a resize message and retain quit.
  Tab switches focus between the status columns and Other sources shelf.
- Enter on any card opens its grouped versions; it never resolves directly from
  the board, even when one version exists. The version view initially highlights
  the ID header, not a default source. Up/Down moves across source rows; Enter on
  a concrete selectable version means Select workspace and calls fresh Resolve
  with exactly that selector. Identical bytes remain separate explicit choices.
- The version detail pane shows title/status, committed/live source, branch or
  detached HEAD, checkout when live, change classification, paths, revision,
  and scrollable source. Tab toggles row/detail focus; PageUp/PageDown scrolls.
  Narrow terminals show the focused pane. Esc returns from versions to the board.
- `b` opens a checkout chooser using live sources from the same inspection.
  Choosing a valid source rebuilds columns from its records; it does not switch
  any Git branch or shell directory. Invalid/absent sources show their reason
  and cannot become a healthy board context. If the invocation checkout is
  invalid, show that state explicitly and offer b; never call it an empty board.
- Successful workspace resolution exits, restores the terminal, prints the
  absolute project path plus newline on stdout (or workspace's existing JSON
  shape under `--json`), and prints escaped target/source context on stderr.
  The UI uses stderr; stdin and stderr must be terminals. Nonterminal use refuses
  before raw mode and points to versions/workspace. Help needs no project/TTY.
  JSON is the lossless path interface for paths with control characters.
- A refused resolution stays in the version view with its reason and r refresh
  guidance. Missing checkout never creates one. Deleted rows remain visible but
  cannot be selected; invalid sources have diagnostics and no admitted versions.
- Incomplete results retain a persistent banner and a scrollable source/diagnostic
  view (`s`). Healthy explicit versions remain selectable. Total inventory failure
  offers retry/quit with no selectable rows, rather than a successful empty board.
- `r` refreshes explicitly. Disable workspace selection while a read is pending;
  clear version selection and return focus to its card/header. Keep board context
  only if its locator/path/attached identity still matches; otherwise require b
  to choose again. Ignore older asynchronous results. Do not poll, watch, cache
  persistently, filter/search, or add mouse interaction in this slice.
- q quits from anywhere; Esc on the board cancels, Esc on an overlay returns.
  Ordinary cancellation exits 0 without a workspace result; Ctrl-C exits 1.
  Terminal setup/output errors exit 1, malformed CLI usage 2. Cancel in-flight
  Git reads and restore terminal state; an old response cannot select a workspace
  after cancellation. No drag/drop or keyboard status mutation yet.

Escape terminal controls (including ESC and C1 characters) in all displayed
record/body/path text. Keep source line structure and original model bytes;
do not emit file-provided ANSI/OSC sequences, open URLs, or touch the clipboard.
Use display-cell-aware clipping for Unicode; meaning must remain clear without
color. Board columns and counts are derived, never written back to records.

## Technical recommendation and alternatives

Recommend Bubble Tea v2 for a small explicit state machine: UI events and read
results update state, and rendering derives from it. Its documented Model,
Update, View and command separation fits deterministic stale-result tests;
this is a design inference, not a Grove prototype result.
[Official documentation](https://pkg.go.dev/charm.land/bubbletea/v2).

Pin `charm.land/bubbletea/v2 v2.0.9` for the proposed implementation baseline;
its tagged module declares Go 1.25, below Grove's Go 1.26 floor.
[Release](https://github.com/charmbracelet/bubbletea/releases/tag/v2.0.9),
[tagged module](https://github.com/charmbracelet/bubbletea/blob/v2.0.9/go.mod).
Reuse its versioned `github.com/charmbracelet/x/term v0.2.2` for terminal checks
and `github.com/charmbracelet/x/ansi v0.11.7` for display width/clipping when
needed, promoting direct imports explicitly. No Bubbles/Lip Gloss dependency
is necessary for the first board and version inspector. Dependency compatibility has been
inspected in documentation, not built in this review.

`tview` is the credible alternative: it supplies tree, table and text widgets,
which could shorten screen assembly. Its widget-first API is less aligned with
the explicitly testable selection/read-result state chosen here; that is an
engineering preference, not a measured performance claim.
[Official tview documentation](https://github.com/rivo/tview).
Avoid a custom raw-terminal framework; it would add lifecycle/input work to the
slice. No benchmark or hands-on framework comparison was conducted.

Keep UI code in `internal/tui`; call existing version contracts in process.
Add context-aware read wrappers to versions/repo while retaining existing
non-context entrypoints for the CLI. Propagate cancellation to all Git commands,
including cat-file, and check it between filesystem/source reads. Allow only one
read operation in flight. Bubble Tea commands do not themselves cancel arbitrary
work, so the UI must own and cancel the read context explicitly.
[Runtime source](https://github.com/charmbracelet/bubbletea/blob/v2.0.9/tea.go).

## Constraints

W-006/W-007/W-008 repairs are prerequisites under the selected after-fixes
sequence; verify their integration, not only done statuses. Q-001 stays accepted.
No record or configuration changes, index/ref changes, worktree provisioning,
claims, agent launches, child shell/editor, body editing, or new record schema.
Terminal mode/output is the only intended interactive side effect. Do not create
debug/panic log files in the repository; strip the framework's debug-log switch
from the UI environment and do not enable file logging. Keep existing CLI output
and ordinary local-file use intact.

## Acceptance

1. Main/feature fixtures with differing work statuses and titles produce columns
   from exactly the chosen live source. b switches displayed context without Git
   writes. One work card exists per ID in that context; branch-only work remains
   discoverable in the unplaced Other sources shelf. Questions/decisions are not
   coerced into work columns. Invalid context is distinct from an empty board.
2. Card Enter opens grouped versions without resolving. A version Enter passes
   exactly its displayed selector to Resolve; no source is selected automatically.
   Identical bytes remain separate choices, and each version shows its own status.
   Deleted/invalid observations cannot resolve; Q-001 stays intact.
3. Fresh resolution exits with workspace's output contract; stale/moved/foreign/
   disappearing/ambiguous/missing targets keep the version view with a reason.
   Incomplete warnings stay visible and attributable. Refresh clears old selection
   and never routes by a previously occupied row position.
4. Board/overlay navigation, responsive columns, scrolling, refresh, context
   changes, and cancellation have deterministic model tests, including late read
   replies and no second read while one is pending. Escape OSC/ANSI/C1 controls;
   test wide Unicode, narrow terminals and large bodies without changing bytes.
5. A pseudo-terminal fixture proves terminal restoration, separate UI/result
   streams, no-output cancellation and cancellation of a blocked Git read. An
   actual board → card → version → resolve → show workflow reads exactly the
   selected bytes. Hash Git/configuration/records across success and refusal.
6. Targeted/full/race suites, vet, formatting, Grove check, dependency review,
   and independent behavior review pass. Owner usability feedback on a real
   Kanban demo remains separate evidence from automated tests or screenshots.

## Next

After the three repair candidates are reviewed and integrated, Fable should
verify this design against those interfaces, implement the linked plan in an
isolated worktree, and return a runnable demo plus test/review evidence. The
owner has selected the terminal experience and requested earlier Kanban value.
The board-first recommendation and its checkout-scoped status policy should be
kept explicit in the handoff; do not silently invent a cross-branch card status. Worktree creation and agent execution
remain separate later work. No implementation is authorized in this review turn.
