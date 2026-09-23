---
id: "G-106"
type: review
title: "G-105 green CI review"
status: current
created: "2026-09-23T04:59:33Z"
updated: "2026-09-23T04:59:37Z"
work: ["G-105"]
examined: "64f0a0c"
---

# G-105 green CI review

Evidence for [G-105](G-105-restore-green-ci-linux-build-go.md). Independent
review by a read-only reviewer subagent in this harness (Claude Code, the
`reviewer` agent), which ran commands and edited nothing. A review is
evidence, not approval.

## Examined

`git diff 8310af1 64f0a0c`: the `unix.Getsid` fix, the Go 1.26.8 and module
upgrades, and three Linux-only test fixes. The reviewer ran on macOS
`GOOS=linux go vet ./...`, `GOOS=linux GOARCH=amd64 go vet ./...`,
`go mod tidy -diff`, the full suite and govulncheck@v1.8.0 ("No
vulnerabilities found"); in Docker `golang:1.26` (go1.26.8 linux/arm64)
`go vet ./...` and the full suite, with python3 present so `TestTerminal`
really ran; and `TestStopAfterReconnect -count=60` in Docker, idle and under
16 CPU hogs. All passed. It diffed goldmark 1.7.13 against 1.8.6 in the
module cache: `ast.Walk` is byte-identical and the `util` functions
`internal/handoff/sources.go` calls behave the same. Its evidence for
lipgloss, x/ansi and ultraviolet is the tui and cli suites on both systems,
not a source diff.

## Findings

Nothing blocking; three notes.

1. `internal/cli/tui_test.go`: the approve/feedback line was still checked
   on the wrapped screen, and a temp path of 50+ extra characters wraps it
   after "feature in". The new unwrapped integrate check is not vacuous: the
   header shows the root alone and truncated. Side observation, before this
   diff: with such a path the review card drops its third wrapped row, so
   the integrate target is not shown at all.
2. `internal/attempt/attempt_test.go`: the relaxed exit check still means
   something, since only the owner's `Kill(-pgid, SIGINT)` sends SIGINT, and
   it still fails on exit 0, another code, or the SIGKILL escalation. The
   first stop in the same test has the same trap gap in principle but did
   not fail in 120 Linux runs.
3. `internal/tui/testdata/terminal.py`: after the fix the first frame alone
   can satisfy "step 19999", so the comment "polled while it runs" claimed
   more than the check proves; polling is covered by
   `internal/tui/attempts_test.go`.

## Disposition

1. Fixed in `1be9030`: both targets are read from the unwrapped card. With
   `TMPDIR` 60 characters longer the approve check now passes and the
   integrate check fails on the product limit the reviewer observed: the
   detail's review card has no room for a third wrapped row. That limit is
   G-043/G-044's card, outside G-105, and is left open. Default temp paths on
   macOS and Linux pass.
2. No change; recorded here to watch if the test flakes.
3. Fixed in `1be9030`: the comment says "shown while it runs".
