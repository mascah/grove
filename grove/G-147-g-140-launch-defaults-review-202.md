---
id: "G-147"
type: review
title: "G-140 launch defaults review, 2026-09-25"
status: current
created: "2026-09-25T04:27:49Z"
updated: "2026-09-25T04:28:44Z"
work: ["G-140"]
examined: "ad52710ff3d7ead867d0b2a68df80d82b0e051eb"
---

## Examined

Independent review by a fresh `grove-reviewer` agent, 2026-09-25, of
`worktree-G-140` from base `670ca9c` to `ad52710` (the implementation
`d9ca2ef` and this repository's `run:` defaults `ad52710`), against
[G-140](G-140-default-an-attempt-s-budget-mode.md)'s acceptance and
constraints, with no plan beyond the record's design.

The reviewer ran these at `ad52710`: `gofmt -l .` (clean), `go vet ./...`,
`go run ./cmd/grove check` (OK, 141 records), `go test -count=1 -timeout
120s ./...` and `go test -short -count=1 ./...` (all ok), and `terminal.py`
(11 of 11 ok). It also probed `check` and `run` against `run:` values in a
temporary fixture, and built base `670ca9c` to compare timings and to see
how an older binary reads the new `grove.yaml`. It did not run the real
provider or `-race`.

## Findings

Acceptance 1 to 5 met. `check` also refuses `1e3`, `.5`, `"50 "`, `~`,
`permission_mode: 5`, a bare `run:` and duplicate keys. No regression was
found in CLI option parsing, the `parseMapping` refactor or `versions`
(`Source.Run` is never serialized).

1. **Knowledge, medium.** G-140 does not mention accepted decision
   [G-141](G-141-never-run-gpt-6-astra-unless-the.md): "what a paid run
   spends … is the owner's explicit answer". The reviewer called the reading
   that a committed, owner-authored `run:` is that answer defensible but
   unrecorded. G-101's required budget and mode still hold, now from a flag
   or `grove.yaml`.
2. **Scope, low to medium.** The defaults come from the checkout that
   launches, so `grove run` inside a work branch's worktree takes that
   branch's `run:`, which an attempt could have edited. The CLI prints the
   values only after the start; the board shows them before Enter.
3. **Compatibility, low.** A binary older than G-140 refuses the merged
   `grove.yaml` (`run: unknown configuration key`), including the installed
   `~/.local/bin/grove` at `670ca9c` and older branches reading `main`.
4. **Documentation, low.** The work guide's invocation table still wrote
   both flags as mandatory.
5. **Record, low.** G-140's Next still said "Proposed, unassigned".
6. **Notes, not defects.** `run` in an invalid or missing project now exits
   1 before the usage check; control characters in `model` and similar
   values pass `check` but are escaped wherever they are shown.

## Disposition

1. Recorded in G-140's Evidence as the author's reading, with G-141 added
   to `relates_to`, and put first among the owner's judgments. Not settled
   by the author.
2. Recorded as a limit in G-140. The record chose the launching checkout,
   and a worktree's other files (its `AGENTS.md`, for one) already shape
   what launches from it. Owner's judgment.
3. Recorded in G-140's Next: rebuild the installed binary when integrating.
   No compatibility is kept before the first release (G-065).
4. Fixed in `b917fc2`: the guide's row now says either flag is optional
   where `grove.yaml` sets it under `run:`.
5. Fixed with the handoff.
6. Accepted as is. The first follows from the order the record asks for.
