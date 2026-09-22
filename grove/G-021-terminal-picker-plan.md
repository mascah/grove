---
id: "G-021"
type: plan
title: "G-017 terminal Kanban and version selection plan"
status: current
formerly: "docs/plans/W-009-terminal-picker.md"
work: ["G-017"]
created: "2026-09-19T20:32:24Z"
updated: "2026-09-21T21:11:16Z"
---

# G-017 terminal Kanban and version selection plan

> For Fable: execute serially in an isolated worktree using repository
> instructions and an inline execution workflow. Do not merge automatically.

**Goal:** launch Grove's TUI when called without a subcommand, initially showing
a checkout-scoped Kanban board with grouped version inspection and explicit
existing-workspace selection. Do not add a `grove board` subcommand.

**Spec:** [G-017](G-017-terminal-picker.md), including the
selected experience, proposed screen, keyboard contract, output and cancellation
semantics, dependency rationale, and acceptance.

**Architecture:** a small `internal/tui` model renders `versions.Result` and
issues cancellable read operations. Existing explicit CLI and version contracts retain
their noninteractive behavior. UI rendering uses stderr, with a workspace result
on stdout only after successful selection and terminal restoration.

**Execution base:** verify G-014/G-015/G-016 are integrated with their reviews.
Do not use their status fields as integration proof. This plan targets their
repaired ownership, final-check, and exact-path helpers. The current review base
is `9b7f730`; no TUI code or dependency was installed during shaping.

## Global constraints and review focus

Keep G-002, schema 1, Go 1.26, current local Git scope, and read-only inspection.
No worktree creation, record edits, editor/shell/agent launches, watchers,
clipboard, file logs, or network service. Pin Bubble Tea v2.0.9 and only needed
terminal/width helpers as specified in the work record; use the v2 API.
Review: group rows accidentally authorizing a source, stale asynchronous replies,
raw terminal-control text in records, incomplete sources disappearing from the
screen, and failure/cancellation leaving terminal mode or Git children behind.
Each task below includes the corresponding acceptance checks.

## Task 1: Make source reads cancellable without changing old CLI behavior

Files: `internal/repo/repo.go`, G-016's shared inventory helper,
`internal/versions/versions.go`, `workspace.go`, `tree.go`, new context tests.

Add these entrypoints, retaining existing functions as Background wrappers:

```go
func GitContext(ctx context.Context, dir string, args ...string) (string, error)
func InspectContext(ctx context.Context, root, id string) (*Result, error)
func ResolveContext(ctx context.Context, root, selector string) (*Workspace, error)
```

Keep internal path/inventory helpers context-aware too. `ResolveContext` must
use the context for inspection and G-014's final source check; an uncancelled
helper or direct cat-file invocation would defeat cancellation.

- [ ] Use a test Git helper that reports it started and waits for cancellation;
  cancel through a context after that handshake. Assert timely exit and process
  collection, with no state writes. Do not depend on sleep scheduling.
- [ ] Confirm new tests fail before implementation. Use `exec.CommandContext`
  at every Git execution, including cat-file; return cancellation distinctly
  using `ctx.Err()` when applicable. Check cancellation between source loads and
  live filesystem reads; do not promise interruption of arbitrary blocked OS
  filesystem calls. Preserve all existing diagnostic text for ordinary errors.
- [ ] Run all repo/versions and relevant CLI tests, then commit
  `feat(versions): support cancellation of source inspection`.

## Task 2: Model board context, explicit selection, and read outcomes

Files: new `internal/tui/model.go`, `model_test.go`, `backend.go`.
Install the pinned UI dependency in this task's implementation, not before the
repair prerequisite check. Keep UI state and backend effects separate:

```go
type Backend struct {
    Inspect func(context.Context, string, string) (*versions.Result, error)
    Resolve func(context.Context, string, string) (*versions.Workspace, error)
}
```

Use `tea.Model` (`Init`, `Update`, `View`) and v2 `tea.KeyPressMsg`.
Model fields include root, result, board source identity (locator/path/ref),
column/card/shelf focus, current card ID, focused version identity, checkout
chooser/version/detail/source view, dimensions/scroll offsets, one active operation/cancel func,
request generation, error text, and optional final Workspace. Row identity is
record ID plus selector, never only an array position. Group rows have no selector.
Messages carry request generation plus result/error; resolve replies also bind
the exact requested selector. Context belongs to the UI session and each request.

- [ ] Create fake committed/live rows with conflicting statuses/titles and
  identical-content/different-selector variants. Check event sequences directly:

```text
initial board card + Enter -> versions header, zero Resolve calls
version header + Enter -> zero Resolve calls
version cursor + Enter -> exactly one Resolve with that row's selector
b + valid live source -> columns change, zero Git mutations
branch-only ID -> Other sources shelf, no invented status
pending read + Enter/r -> no second operation
resolve error -> remain open, show reason, no final Workspace
refresh -> clear version focus, increment generation, request Inspect
old-generation reply -> ignore, keep current state
q -> cancel request, no result; Ctrl-C -> cancel request, interrupted result
```

Derive four columns from Version.Record where Type == work and Source equals
the explicitly scoped live source, retaining version ordering from Inspect.
Derive the shelf from work groups with no live record in that source. Test
divergent titles/statuses, question/decision exclusion, empty columns, detached
context, deleted rows, incomplete valid subsets, no work records, invalid
context versus empty board, inventory failure, and a disappeared card after
refresh. On context identity changes require a new checkout choice.
- [ ] Implement asynchronous `tea.Cmd` effects with owned cancellable contexts.
  Wire q/Esc and interrupt to cancellation; do not assume `tea.Quit` cancels
  backend work. After Run returns, cancel and collect in-flight read cleanup.
- [ ] Run model tests under the race detector with out-of-order fake replies.
  Commit `feat(tui): model Kanban context and explicit version selection`.

## Task 3: Render a usable terminal screen and restore its state

Files: new `internal/tui/view.go`, `view_test.go`, `run.go`, `run_test.go`.
Expose:

```go
// nil workspace and nil error means ordinary user cancellation.
func Run(ctx context.Context, root string, input, screen *os.File) (*versions.Workspace, error)
```

`context.Canceled` represents interruption to the CLI; use `tea.WithInput`,
`tea.WithOutput`, `tea.WithContext`, and v2 `tea.View.AltScreen`.
Pin direct terminal/width helpers at versions from the spec. Pass a framework
environment with TEA_DEBUG removed and never enable file logging.

- [ ] Add rendering tests at 120x30, 80x24, 40x10, and below minimum size:
  checkout-labelled board, four columns at wide sizes, one column with tabs
  below 100 columns, context-title cards, an unplaced Other sources shelf,
  ID-only version headers, status/source per version, selected-version details,
  always-visible incomplete/loading/error notice and key hints. Cards, shelf,
  versions and source diagnostics must remain reachable by scrolling. Tab switches
  board/shelf or version/detail focus; s shows source diagnostics. Esc returns
  from overlays, and q quits everywhere. At narrow sizes show the focused pane.
- [ ] Test an OSC clipboard sequence, ANSI cursor erase sequence, C1 controls,
  combining/wide Unicode and newline/tab paths in every displayed input type.
  Escape controls before cell-aware clipping. Show source line breaks as lines;
  display CR/other controls visibly. Do not modify the source buffer or hash it
  after display escaping. Color is optional; semantics must survive no color.
- [ ] Render only sanitized record content; use framework-generated control
  sequences solely for the interface. Add no Markdown renderer or clickable links.
- [ ] Pseudo-terminal tests prove alternate-screen/raw-mode restoration after
  selection, q, Ctrl-C, output failure, and cancelled blocked Git. Use a Unix
  PTY integration harness (Python stdlib `pty` plus `termios` is sufficient on
  the currently supported development platforms); retain the harness under
  `internal/tui/testdata/terminal.py` and document its invocation. Do not claim
  Windows PTY coverage. Unit tests alone do not prove terminal lifecycle.
- [ ] Run view/model/PTY tests and commit `feat(tui): render the terminal Kanban board and version view`.

## Task 4: Make the TUI the default invocation and prove the connected workflow

Files: `internal/cli/cli.go`, `cli_test.go`, new `internal/cli/tui.go`, `tui_test.go`,
`internal/cli/workspace.go` (shared result formatting), README and work/plan bodies.
Keep `cli.Run` signature unchanged: the TUI requires `os.Stdin` and errOut to be
terminal `*os.File` handles; ordinary tests with buffers exercise refusal/help.
Run real interactive behavior through the PTY harness. Extract existing workspace
result formatting into a private CLI helper shared by workspace and the TUI.

- [ ] Add parser tests for no arguments, project-only/JSON-only invocation and
  both flag orders in `grove [--project DIR] [--json]`. Default mode accepts no
  positionals; reject `board`, bare work IDs, duplicate/unknown options, and
  subcommand-only flags (`--source`, `--slug`, `--expect`, `--set`, `--unset`).
  Do not let the new default bypass normal option validation. Retain explicit
  subcommand parsing. `--help`, `-h`, and `help` need no project/TTY and never
  enter the UI. Update the existing nil-arguments usage-error expectation in
  `TestUsageAndMissingID`: no command now selects the TUI, not a usage error.
  Reject nonterminal stdin/stderr before entering raw mode with exit 1 and
  no stdout, with guidance to list/versions/workspace and --help. stdout may be
  redirected; --json selects only the final result encoding and still needs a TTY.
- [ ] Dispatch the default TUI before synchronous `project.Load` in `cli.Run`.
  Pass an empty ID filter to InspectContext; the first screen includes all work.
  Route the TUI through the same project discovery/invalid-current-checkout
  allowance as versions. Begin reads in the UI so startup/loading stays responsive.
  Do not synchronously pre-load the full project before launching the TUI: extract
  a discovery-only helper in `internal/project` shared with `Load` if necessary;
  keep all original inspection validation and boundaries unchanged.
- [ ] After Run restores the terminal, emit exactly workspace's path or JSON
  shape on stdout and escaped target context on stderr. q/Esc emits nothing;
  Ctrl-C/error returns 1, usage 2. Rendering must never contaminate result stdout.
- [ ] Use an actual main/feature fixture: inspect scoped status columns, choose another
  checkout with b, open a card, inspect differing version details,
  select a concrete version, resolve, and `show --project` reads exactly those
  bytes. Then change the target before selection and verify refusal stays visible
  until refresh and explicit reselection. Repeat with one invalid unrelated source.
  Hash all Git/config/records/unrelated files across every path; no state changes.
- [ ] Run targeted tests, `go test -count=1 ./...`, `go test -race -count=1 ./...`,
  `go vet ./...`, formatting, `go run ./cmd/grove check`, and the PTY harness.
  Commit `feat(cli): launch the TUI by default` with documentation. Show
  `go run ./cmd/grove` as this checkout's demo command; document explicit help
  and noninteractive subcommands without claiming the installed CLI was replaced.
- [ ] Obtain independent review against G-017 acceptance and provide a runnable
  owner demo. Record human feedback separately from automated evidence. Reconcile
  the work record and brief via current CLI/body edits; no branch merge or
  worktree provisioning is part of this unit.

## Adjustments made while implementing, 2026-09-19

Implemented on branch `worktree-W-009` from `acfc905`, after confirming by Git
ancestry that the G-014 to G-016 repair and review-fix commits are in main.
The [evidence](G-029-board-review.md) holds commits and results.
Each adjustment below stays inside G-017's outcome; the reason is what was
found when the proposal met the code and the pinned framework.

- **Debug switches are process-wide.** Bubble Tea v2.0.9 reads `TEA_DEBUG`
  (panic log file) and `TEA_TRACE` (trace file) with `os.Getenv`, not from the
  environment passed to a program; its renderer reads `UV_DEBUG` from the
  passed environment. A filtered `tea.WithEnvironment` would therefore not
  have stopped log files. `tui.Run` unsets all three in the process instead,
  and the terminal harness runs with all three set and asserts no file appears.
- **The framework discards screen write errors** (`_ = p.flush()` in its
  render loop), so a dead screen left a live session in which Enter could
  still select a workspace unseen. `tui.Run` wraps the screen in a writer that
  records the first failed write and cancels the session: exit 1, no result.
- **The invocation's checkout is identified by Git, not by path.**
  `versions.Result.GitDir` (one more `rev-parse --git-dir`) names root's own
  worktree; the board starts on the live source with that Git directory. It is
  not part of `versions --json`, which is built field by field.
- **Tasks 2 and 3 are one commit.** The model calls the view for scroll
  limits, so a model-only commit would not compile. `Backend` lives in
  `model.go`; there is no `backend.go`.
- **Discovery is shared, loading is not.** `project.Discover` is `Load`'s own
  first step, exported; the board is dispatched before `project.Load`.
- **The connected-workflow fixture is a Go test** (`internal/cli/tui_test.go`)
  that drives the real model against real Git and asserts on rendered screens,
  because a cell-diffing renderer's byte stream is not a dependable screen to
  assert against. The pseudo-terminal harness covers what only a terminal can
  show: mode restoration, stream separation, select-to-`show` through the
  built binary, resize, refusal without a terminal, and killing a blocked Git
  child after a started-handshake.
- **Exact values are hard-wrapped.** Word wrapping split selectors and paths
  at hyphens and spaces on real data; field values and source lines in the
  sources view break only at the pane edge. Prose wraps between words.
- **Essential state comes before long paths** in the header (board label
  first) and the checkout chooser (valid/unavailable first): real temporary
  paths clipped both off screen in the first version.
- **macOS sets `PENDIN`** in a terminal's local flags when raw mode ends with
  input queued. The harness compares every termios field but masks that one
  kernel-owned bit.
- Direct dependencies as pinned by the record: `charm.land/bubbletea/v2
  v2.0.9`, `github.com/charmbracelet/x/ansi v0.11.7` (display width, clipping,
  wrapping), `github.com/charmbracelet/x/term v0.2.2` (terminal check). All
  built and tested under Go 1.26.2; no Bubbles or Lip Gloss.
