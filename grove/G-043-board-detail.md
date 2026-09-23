---
id: "G-043"
type: work
title: "Make the board and item detail clear and visually polished"
status: done
created: "2026-09-21T00:54:15Z"
updated: "2026-09-23T00:06:54Z"
kind: feature
size: medium
priority: 3
depends_on: ["G-042"]
relates_to: ["G-035", "G-017", "G-030", "G-041", "G-044", "G-064", "G-065"]
formerly: "W-025"
candidate: "178b305"
---

## Outcome

A visually polished terminal board and detail view make current work, its
content, blockers, artifacts and next action immediately understandable.

## Scope and bounds

Use G-042's project-wide projection. Board columns follow the target lifecycle;
show compact cards and clear borders/focus, hide Abandoned by default, and bound
the recent Done column while retaining searchable older work. Detail leads with
rendered Markdown; metadata, artifacts and a navigable timeline support it.
Keep alternate sources available as secondary information. Include a useful
list/search path; dependency graph visualization is later.

Bound Done through presentation, never by relocating files. Resolve type and
artifact roles from supported metadata, not ID prefix or folder. List/search
also makes G-065's general knowledge discoverable without turning it into work
cards or loading every page into an agent's context.

Design board, detail and the future review view together before implementation,
but implement review actions in G-044. Use Bubble Tea v2 with compatible Charm
layout/components/Markdown rendering as appropriate. Preserve terminal text
escaping and separate trusted renderer styling from untrusted content.

## Acceptance

1. The owner reviews a concrete visual proposal and the implemented terminal
   experience, including representative narrow/wide sizes and many Done items.
2. Current content is the first useful detail; artifacts and timeline entries
   can be reached without interpreting a list of checkout observations.
3. Focus, scrolling, search, empty/loading/incomplete states and navigation
   remain understandable; meaning does not depend on color alone.
4. Lifecycle, alternate-screen restoration, resize and cancellation checks pass;
   verify connected navigation using real work and attached artifacts.
5. On-demand history and measured loading stay responsive. Bare invocation and
   explicit noninteractive command contracts remain intact.

## Evidence

Implemented on 2026-09-22 on branch `worktree-G-043` (worktree
`.claude/worktrees/G-043`) from `main` dc3b9b6, which holds G-042. It
started from this record at 8bf691b and plan
[G-096](G-096-g-043-board-and-detail-design-vi.md), whose visual proposal
the owner answered on 2026-09-22: implement as drawn, links as plain text.
The plan's Adjustments section records what changed while implementing.
The candidate is the commit that adds this evidence; `git log
dc3b9b6..178b305` lists it all.

Behavior, against the acceptance:

1. The owner reviewed the proposal (G-096, the 120- and 80-column mockups
   with real records and a bounded Done) and answered it. Judging the
   implemented terminal experience is still theirs: run `go run ./cmd/grove`
   in this worktree at a wide and a narrow size; this repository has 27 Done
   records, so the bound and its `+ N older · / to search` footer show.
2. Enter on a card opens the detail: a boxed header, the body rendered from
   Markdown, and a sidebar of linked records by role (a review with what it
   examined against the candidate), the timeline with Enter showing the
   record as of a commit, and the Sources last. Alternate sources are behind
   `v`, which opens the unchanged versions screen.
3. The focused card has a heavy border and `▶`; each column an accent colour
   nothing depends on; Done is bounded to its page newest first and counted;
   Abandoned is hidden until `a`; `/` searches every record of every type by
   ID, type, status and title; empty, reading, error and vanished states are
   worded; below 100 columns one column or one pane shows at a time; below
   16 rows the header shrinks and cards drop to four rows.
4. The pty scenarios (lifecycle, alternate screen, resize, cancellation,
   blocked Git) pass through the detail. Connected navigation runs on real
   Git in `internal/cli`'s board workflow and, on this repository's own
   records, in a scripted pty: board, detail of G-043 with its plan, needs
   and related links, Tab to the timeline, as-of 3ae5df7, `v`, and `/G-06`.
5. History is read only with a detail open and any key cancels it; the board
   load starts no Git process it did not before; the render of the largest
   record (31 KB) takes 8 ms and is cached per content and width. Bare
   invocation and every subcommand keep their contracts; `terminal.py`'s
   noninteractive checks are unchanged.

Decisions: glamour v2.0.1 with lipgloss v2.0.4, the pair that keeps
bubbletea v2.0.9's pins; an ASCII base style with bold headings and ANSI
colours; record bodies escaped per line before glamour and, after the review,
filtered to glamour's own SGR styles per row, with every HTML character
reference made literal; OSC 8 hyperlinks stripped by the owner's choice;
relative link targets shown root-relative; the Done page bound by the
column's rows; Esc unwinds as-of, then the previous record, then the board,
and returns from the versions to the detail.

Verification at ed11f51 (the review's examined commit; the candidate differs
by records and one README paragraph): `gofmt -l .` empty; `go vet ./...`
clean; `go run ./cmd/grove check` OK, 93 records; `go test -count=1
-timeout 120s ./...` all ok (cli 4.2 s, tui 6.2 s with the pty scenarios,
0.4 s under `-short`, versions 7.8 s as before); `go test -race -timeout
120s ./internal/tui` ok in 13 s.

Review: [G-097](G-097-review-of-g-043-board-detail-and.md), two rounds.
Round 1 found a blocking trust-boundary break (character references
decoded after the escaping) and three small state bugs; ed11f51 fixed all
four with tests, and round 2 confirmed each and found one cosmetic README
wording, corrected in the candidate. Verdict: clear with notes.

Limits: the owner's judgment of the visual result in a real terminal
(acceptance 1 and 3) is not supplied by these checks; the plan's nullsec
pass was not run; a reference inside a code block or link target shows an
extra `&amp;`; the sidebar clips long titles; a timeline commit's content
is held only while its history is; dependency graph visualization and the
review view's actions remain G-044's.

## Next

The integrator judges the candidate in a terminal and, on acceptance,
records the verdict and integrates. From the repository root:

```sh
git -C .claude/worktrees/G-043 diff --stat 178b305 HEAD   # only this record
(cd .claude/worktrees/G-043 && go run ./cmd/grove)          # the board; Enter, Tab, v, /, a
git merge worktree-G-043
# edit grove/G-043-board-detail.md: add "Verdict: <the owner's words>" under Evidence
go run ./cmd/grove update G-043 --set status=done --commit
```

Verdict: Significant improvement over previous TUI experience.