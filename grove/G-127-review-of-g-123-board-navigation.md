---
id: "G-127"
type: review
title: "Review of G-123 board navigation and cards"
status: current
created: "2026-09-24T04:35:31Z"
updated: "2026-09-24T04:35:49Z"
work: ["G-123"]
examined: "6e78245f4c0ad02e32ab1fae01299dbc7c2b9446"
---

## Examined

Commit `46f1322` (round 1) and the fixes in `6e78245` (round 2) on
`worktree-G-123`, against G-123's Outcome, Constraints and Acceptance, by an
independent reviewer subagent that edited nothing. It read the diff from
`b262d00`, ran `go test -short ./internal/tui`, `go vet`, `gofmt -l .` and
`go run ./cmd/grove check`, and proved findings with throwaway tests it
deleted afterwards.

## Findings

Round 1, `46f1322`:

1. Should-fix. When a refresh dropped a record from the detail path,
   `settleFocus` did not move `workDepth`, the depth of the record `o`
   opened. The breadcrumb could then name the wrong layer, and Esc from a
   record opened later could return to the attempts screen. The mechanism
   predates G-123, but the new breadcrumb puts the depth on screen.
2. Nit. No test covered the reset when a cut drops the layer `o` opened.
3. Nit. An orphaned attempt gets the same border as a running one.
4. Nit. `runningAccent` duplicated the existing `cyan` style.
5. Nit. The `not on main` exception in docs/board.md read as if it applied
   to the header.

Round 2, `6e78245`: each fix confirmed correct, with no new findings. The
reviewer traced every write of `workDepth` and found no way for it to
exceed the stack length.

## Disposition

- Findings 1, 2, 4 and 5 are fixed in `6e78245`. The test for finding 2
  also covers finding 1, and it fails with the `settleFocus` fix removed.
- Finding 3 is kept as designed, since `live()` counts an orphan as live.
  docs/board.md says so. It is left to the owner's visual judgment under
  Acceptance 3.
