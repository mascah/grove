---
id: "G-130"
type: review
title: "Review of G-124 board re-reads on focus and moved tips"
status: current
created: "2026-09-24T05:08:09Z"
updated: "2026-09-24T05:08:23Z"
work: ["G-124"]
examined: "70509fba23b50e57f3c66441423706879b5debea"
---

## Examined

Commit `a6807a4` (round 1) and the fixes in `70509fb` (round 2) on
`worktree-G-124`, against G-124's Outcome, Constraints and Acceptance, by an
independent reviewer subagent that edited nothing. It read the diff from
`2491eda`, and ran `go test -short ./internal/tui`, `TestTerminal` without
`-short`, `go vet` and `gofmt -l`. Round 2 also reviewed the combined diff
`2491eda..70509fb` once. The only later commit before the handoff,
`eba3197`, rewords docs/board.md for the round-2 nit and changes no code.

## Findings

Round 1, `a6807a4`:

1. Should-fix. An automatic re-read closed a detail showing a timeline
   commit or a diff, since `inspectMsg` clears `asOf` and `diff`. With an
   attempt running, each of its commits would yank the owner back from its
   own changes.
2. Nit. Focus clears the versions cursor and a refusal, as `r` does.
3. Nit. docs/board.md said "unless a read is already under way" where a
   pending action also counts, and "starts no process by itself" beside a
   re-read on focus.
4. Nit. `FocusMsg` did not check `m.done`; harmless, since Run cancels the
   read.

The reviewer checked and found sound: no Model access off the Update
goroutine; branch tips and committed sources come from the same
`listBranches`, so current tips never re-read in a loop; a moved tip never
restarts an inspection under way; Git runs through `repo.GitContext`; the
tests fail with the logic removed.

Round 2, `70509fb`: findings 1, 3 and 4 confirmed fixed with no regression.
One new nit: `pinned()` ignores the screen, so a detail left at a diff
still holds automatic re-reads while the attempts screen is open above it,
where the docs spoke only of a detail showing one.

## Disposition

- Finding 1 fixed in `70509fb`: `pinned()` holds the focus and moved-tip
  re-reads (not the one when an attempt ends); both tests cover it.
- Finding 2 kept: the record says focus does what `r` does.
- Findings 3 and 4 fixed in `70509fb`.
- The round-2 nit is kept as behaviour, so Esc still returns to the diff
  left open, and docs/board.md says so in `eba3197`, self-checked, not
  re-reviewed.
