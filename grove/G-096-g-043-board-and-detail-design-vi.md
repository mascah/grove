---
id: "G-096"
type: plan
title: "G-043 board and detail design: visual proposal and implementation plan"
status: current
created: "2026-09-22T22:59:51Z"
updated: "2026-09-22T23:04:02Z"
work: ["G-043"]
---

## Inputs

Prepared on 2026-09-22 for [G-043](G-043-board-detail.md) at revision
`7c64e190`, on branch `worktree-G-043` from `main` `dc3b9b6`, which holds
G-042 (candidate `b56646c`, done). Read in full: G-043, G-042, G-017, G-030,
G-044, G-036, the brief's "Current view and TUI" section, the record model's
type table, and `internal/tui` (`model.go`, `view.go`, `run.go`, the tests).

The owner's recorded feedback on the board, which this plan answers:

- G-017, 2026-09-19: "an acceptable basic starting point; clearer boundaries
  and colour can come later"; the version view "did not help"; "source" and
  `b` were unclear.
- G-030, 2026-09-20: lineage is "a step in the right direction"; the owner
  wants to iterate on "the information architecture of the TUI".
- G-041 recorded no board feedback: the owner's judgment of nullsec's
  converted board (its acceptance 6) was left open when G-041 closed.

## Design

Principles, from G-043 and the brief: the board is the project-wide current
view (G-042) in the target's lifecycle order; the first useful detail is the
record's own content, rendered; sources and versions are secondary; meaning
never depends on colour; every displayed byte from a record, path, or Git goes
through the existing escaping (`safe`) before any renderer or style touches
it; Done is bounded by presentation, never by moving files; type and artifact
roles come from `type`, `work`, `depends_on`, `blocks`, `members` and
`relates_to`, never from an ID or a folder.

### Screens and keys

| Screen | Reached by | Keys |
| --- | --- | --- |
| Board | bare `grove` | ←→↑↓/hjkl move; Enter open; Tab Deleted/Elsewhere; `/` search; `a` show/hide Abandoned; `b` view or checkout; `s` sources; `r` refresh; Esc/q quit |
| Detail | Enter on a card, a search hit, or a linked record | ↑↓/PgUp/PgDn scroll the focused pane; Tab cycles Content → Linked → Timeline; Enter opens a linked record or shows the record as of a timeline commit; `v` versions and places; Esc back (as-of view first, then the previous record, then the board) |
| Search | `/` on the board | type to filter; ↑↓ move; Enter open; Esc close |
| Versions and places | `v` in a detail | unchanged from G-017/G-042: rows, folds, Enter selects a workspace, Tab details |
| Chooser, Sources | `b`, `s` | unchanged |

Enter on a card opens the detail, not the version list: the list the owner
found unhelpful becomes the secondary screen `v`. Workspace selection lives
there unchanged, so G-002, G-011 and G-017's selection contract hold.

### The board

Columns are Proposed, Active, Review, Done; Abandoned is hidden until `a`
shows it, and the header counts it so nothing is silently missing. Cards are
bordered boxes: ID and tag on the first line, the title wrapped to two lines,
and one metadata line (kind · size · priority, or for Done the `updated` date
and candidate). The focused card has a heavy border and a `▶` marker, so focus
reads without colour; columns get one accent colour each from the ANSI 16
palette, which follows the terminal's own theme without a background query.

Done shows the most recent cards by `updated` (then `created`, then ID) that
fit the column, newest first, and a footer counts the rest: the bound is the
page, so it adapts to the terminal, and older done work stays reachable by
search. Focus never lands on a card the bound hides.

Wide (120×36) with this repository's real records, once G-043 is active on
this branch:

```text
Grove  current view · target main  ·  read 2 branches, 2 checkouts                     ~/GitHub/mascah/grove
 Proposed 3                  Active 1                    Review 0                     Done 27 · 4 recent
 ╭──────────────────────────╮ ┏━━━━━━━━━━━━━━━━━━━━━━━━━━┓                             ╭──────────────────────────╮
 │ G-044                    │ ┃▶G-043       not on main  ┃   (none)                    │ G-042                    │
 │ Review candidates and    │ ┃ Make the board and item  ┃                             │ Derive a project-wide    │
 │ integrate approved work… │ ┃ detail clear and visual… ┃                             │ current view of work     │
 │ feature · large · P3     │ ┃ feature · medium · P3    ┃                             │ done 2026-09-22 · b56646c│
 ╰──────────────────────────╯ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━┛                             ╰──────────────────────────╯
 ╭──────────────────────────╮                                                          ╭──────────────────────────╮
 │ G-045                    │                                                          │ G-036                    │
 │ Run one bounded          │                                                          │ Complete the interactive │
 │ implementation independ… │                                                          │ Grove adoption milestone…│
 │ feature · large · P4     │                                                          │ done 2026-09-22 · 977d98b│
 ╰──────────────────────────╯                                                          ╰──────────────────────────╯
 ╭──────────────────────────╮                                                          ╭──────────────────────────╮
 │ G-046                    │                                                          │ G-041                    │
 │ Launch and inspect       │                                                          │ Cut nullsec over to this │
 │ managed attempts from…   │                                                          │ Grove and uninstall the… │
 │ feature · medium · P4    │                                                          │ done 2026-09-22 · b6db474│
 ╰──────────────────────────╯                                                          ╰──────────────────────────╯
                                                                                       ╭──────────────────────────╮
                                                                                       │ G-081                    │
                                                                                       │ Run secure CI and        │
                                                                                       │ Dependabot on GitHub     │
                                                                                       │ done 2026-09-22 · 81cd0c0│
                                                                                       ╰──────────────────────────╯
                                                                                         + 23 older · / to search
 Deleted: none                                                          Abandoned 0 hidden · a shows
 ←→↑↓ move   Enter open   / search   a abandoned   Tab deleted   b view or checkout   s sources   r refresh   q quit
```

Narrow (80×24): one column at a time with tabs, as today, and the same cards
at full width:

```text
Grove  current view · target main · 2 branches, 2 checkouts    ~/GitHub/mascah/grove
 Proposed 3   [Active 1]   Review 0   Done 27   (Abandoned 0 hidden)
 ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
 ┃▶G-043                                                        not on main    ┃
 ┃ Make the board and item detail clear and visually polished                 ┃
 ┃ feature · medium · P3                                                      ┃
 ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
 Deleted: none
 ←→ columns  ↑↓ cards  Enter open  / search  a  Tab  b  s  r  q quit
```

Below 40×10 the resize message stays. A card's tags are the G-042 ones
(`⑂ 2 states`, `uncommitted`, `not on main`) unchanged.

### The detail

The header is a boxed card: ID and status, the title, then one line of
metadata (kind, size, priority, candidate, whether the current state is on the
target, where it is held, `updated`). Below it, two panes at 100 columns and
more, one at a time below that (Tab cycles them):

- **Content**: the record body rendered from Markdown (headings, emphasis,
  lists, code, tables) at the pane width, scrollable. This is the first useful
  detail (G-043 acceptance 2).
- **Linked**: every record connected to this one, each with its role, ID,
  title and status: `plan` and `review` (records whose `work` names it),
  `needs` (`depends_on`), `needed by` (work whose `depends_on` names it),
  `blocked by` (questions whose `blocks` names it), `member` and `part of`
  (`members`), `related` (`relates_to`). Enter opens the linked record's own
  detail, of any type, and Esc returns.
- **Timeline**: G-030's history, newest first, with an `uncommitted` first
  row where the checkout's files differ. Enter on a commit shows the record as
  it was at that commit in the Content pane, labelled, since the history read
  already fetches those bytes; Esc returns to now.
- **Sources**: one line per current state and a count of older places, as
  the version list's header describes them today, ending with `v` for the
  full list. Not a pane: it cannot be focused, only read.

Wide (120×36), G-042 in this repository:

```text
Grove  current view · target main  ·  read 2 branches, 2 checkouts                     ~/GitHub/mascah/grove
 ┏━ G-042 · done ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
 ┃ Derive a project-wide current view of work                                                                     ┃
 ┃ feature · medium · P3 · candidate b56646c · on main · same on 2 branches, 2 checkouts · updated 2026-09-22     ┃
 ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
 ▶ Content  1-26 of 214                                        │ Linked
                                                               │   plan       G-093  Project-wide current view: per-…   current
 ## Outcome                                                    │   review     G-094  Review of G-042 current view        current
                                                               │   review     G-095  Review of G-042 integration target current
 Opening Grove from any linked checkout presents the same      │   needs      G-037  Represent domain terms and linke…  done
 useful project-wide current work view, with source            │   needs      G-038  Hand implementation candidates i…  done
 observations and genuine divergence accessible without        │   needed by  G-043  Make the board and item detail c…  active
 making every user choose a branch first.                      │   related    G-035  Adopt the interactive adoption l…  accepted
                                                               │   related    G-002  How should the board present dif…  resolved
 ## Selected direction                                         │   … 6 more
                                                               │ Timeline  on branch main
 Collapse identical observations; move demonstrably            │   2026-09-22 22:55  done      dc3b9b6  docs(G-042): set status=done
 superseded states into history; expose unintegrated           │   2026-09-22 22:52  review    4b288e7  docs(G-042): set status=review…
 progress and label live uncommitted changes. Keep genuine     │   2026-09-22 22:41  active    b56646c  docs(G-042): clarify what the r…
 competing changes visible. Never use the largest status,      │   2026-09-22 22:30  active    d1dd581  docs(G-042): record the integra…
 latest timestamp or newest branch tip as authority. A revert  │   … 9 more
 is a real change, not automatically an older state. Exact     │ Sources
 source routing and stale-selection checks remain intact;      │   current: done on branch main, checkout . · on main
 view selection does not merge content or grant write          │   older: 2 places (branch worktree-G-042, checkout G-042)
 authority.                                                    │   v lists every version and place, and selects a workspace
 ↑↓ PgUp/PgDn scroll   Tab linked/timeline   Enter open   v versions and places   Esc board   q quit
```

The same record at 80×24 shows one pane, Content first, with the header
shortened to ID, status, title and the metadata line wrapped; Tab reaches
Linked and Timeline. A page, term, decision, question, plan or review opens in
the same screen: no status columns apply, the metadata line shows what the
type has (`work` for a plan or review, `examined` for a review, `blocks` for a
question), and Linked shows `relates_to` both ways plus the work a plan or
review belongs to.

### Search

`/` opens a list over every record of the project in its current state, of
every type, including hidden Abandoned and older Done work, pages and terms:
typing filters by case-insensitive substring of ID, type and title; Enter
opens the detail. This is the list/search path G-043 asks for and how G-065's
general knowledge is discoverable without loading it (acceptance 2 of G-043,
and the brief's "searchable older work").

```text
 / current view                                                     23 of 97 records
 > review
   G-022  review   current   Integrated CLI review, 2026-09-19
   G-024  review   current   Predecessor `/work` review for G-023
   G-038  work     done      Hand implementation candidates into revision-bound human review
 ▶ G-044  work     proposed  Review candidates and integrate approved work locally
   G-058  term     settled   Review
   …
 type to filter   ↑↓ move   Enter open   Esc close
```

### The review view, designed here and built in G-044

For a work record in `review`, the detail header gains a Review block between
the metadata line and the panes: the candidate, whether the branch tip is the
candidate (`git diff --stat CANDIDATE TIP` touching only the record), whether
the candidate is on the target, and each review with its `examined` commit
marked `= candidate` or `before candidate`. The Content pane starts at the
record's Evidence section when one exists. G-044 adds its actions as keys on
this screen (approve, feedback, integrate) and their CLI equivalents; G-043
shows the facts only.

```text
 ┏━ G-043 · review ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
 ┃ Make the board and item detail clear and visually polished                                                     ┃
 ┃ feature · medium · P3 · candidate 1a2b3c4 · not on main · on branch worktree-G-043 · updated 2026-09-23         ┃
 ┃ Review: candidate 1a2b3c4 is the branch tip · review G-097 examined 1a2b3c4 (= candidate)                       ┃
 ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
```

### Rendering and trust

Charm modules, checked by building them together with the pinned
`charm.land/bubbletea/v2 v2.0.9` on 2026-09-22 (a throwaway module outside
every checkout, `go build` clean):

| Module | Version | Needs | Effect on `go.mod` |
| --- | --- | --- | --- |
| `charm.land/glamour/v2` | v2.0.1 | lipgloss/v2 v2.0.4, x/ansi v0.11.7, goldmark v1.7.8 | adds goldmark and goldmark-emoji |
| `charm.land/lipgloss/v2` | v2.0.6 | ultraviolet 20260811, x/ansi v0.11.8 | lifts ultraviolet from 20260703 (bubbletea's) and x/ansi from v0.11.7 |
| `charm.land/bubbles/v2` | v2.2.1 | bubbletea v2.0.8, lipgloss v2.0.5 | not taken: a text input and a viewport are a few lines here |

The lifted `ultraviolet` is the one compatibility risk; step 1 proves it with
the existing suite, including the pseudo-terminal test, before anything is
built on it. Glamour's own module path is `charm.land/glamour/v2`, not the
`github.com/charmbracelet/glamour` v1 line.

Observed in that build: glamour passes a raw `ESC` in its input through to
its output unchanged, and turns every Markdown link into an OSC 8 terminal
hyperlink. So the content pipeline is: record bytes → `safe` (the existing
escaping of controls, format characters and invalid bytes) → glamour, which
adds only its own styles → strip OSC 8 → one row per line, clipped by display
cells (`ansi.Truncate`). A test feeds a body holding `ESC`, an OSC, C1 bytes,
a bidirectional override, a wide title and a Markdown link, and asserts that
the rendered rows hold no control byte but glamour's SGR sequences and no
OSC. Rendering is cached per (revision, width) so a scroll or a key never
re-renders; the largest record here (G-042, about 1,800 words) is timed in
step 2 and must render well under a frame.

Hyperlinks are the owner's call, asked alongside this proposal (see Open
choices): the default is to strip them, because G-017 says the board emits no
file-provided OSC, and because record links are relative paths that a
terminal would resolve nowhere.

Styles: glamour runs with one fixed style (a copy of its `notty` style with
bold headings and ANSI 16 accents), never the terminal-querying `auto` style,
so tests render the same bytes as a terminal and no background query is
needed. Lip Gloss draws borders and column accents from the ANSI 16 palette.
Every heading, focus mark and tag keeps a textual form.

### What does not change

The current view and its tags (G-042), `b`, `s`, `r`, the versions screen and
workspace selection (G-017, G-002, G-011), on-demand history and its
cancellation (G-030), one read at a time, no history during the board load,
identity-based focus across refreshes, the resize floor, the nonterminal
refusal and exit codes, and every explicit subcommand. Records, refs, index
and worktrees are never written.

## Steps

1. **Dependencies.** `go get charm.land/glamour/v2@v2.0.1
   charm.land/lipgloss/v2@v2.0.6`, tidy, and run `go test -count=1 ./...`
   including `TestTerminal`: proves bubbletea v2.0.9 on the lifted
   ultraviolet before any code depends on it. Commit alone.
2. **Markdown pipeline** (`internal/tui/markdown.go`): `body` (the source
   after its frontmatter), `render(source, width)` as above with its cache,
   the hostile-input test, and a timing of G-042's body. Commit.
3. **Board.** Bordered cards, heavy focus border and marker, column accents,
   the header, Abandoned hidden with `a`, the Done bound and footer, the
   Deleted/Elsewhere row. Focus, `window`, and the shelf keep working; a card
   the bound hides is never focused. Tests: layout at every size with many
   Done cards, the bound's order and footer, `a`, focus after refresh.
4. **Detail.** New screen and state (a stack of open IDs, the focused pane,
   sidebar cursor, content scroll, as-of commit): header, Content, Linked
   (derived from the result's groups; roles as above), Timeline reusing
   `History` with `Source []byte` added to `versions.Commit`, Sources line,
   `v` to the versions screen and Esc back. Enter on a card opens it. Tests:
   the first rows are the rendered body; every role from a fixture; opening
   a linked record and returning; as-of content and return; narrow one-pane
   cycling; cancellation and stale replies unchanged.
5. **Search.** Overlay, filter, open. Tests: hidden Abandoned and older Done
   are found; a page and a term open in the detail; hostile titles stay inert.
6. **Docs.** README's board section (the screens, keys, search, bound,
   Abandoned), AGENTS.md's TUI sentence (G-043 owns the board and detail),
   and this plan's adjustments. The record model changes nothing.
7. **Verification.** Per package while iterating; then `go vet ./...`,
   `gofmt -l .`, `go run ./cmd/grove check`, `go test -count=1 -timeout 120s
   ./...`; `-race` for `internal/tui` only, `-timeout 120s`; a real-terminal
   pass of board → detail → linked record → timeline → `v` → workspace
   selection, resize at 80 and 120 columns, Ctrl-C and `q`, on this
   repository and on nullsec with the installed CLI's checkout (read-only).
8. **Review and handoff.** An independent review of the final diff, recorded
   as a review record with `work` and `examined`; then the record's Evidence
   and Next, and `status=review` with the candidate.

## Open choices

Asked of the owner with this proposal on 2026-09-22:

1. Implement the proposal as drawn, or adjust it first.
2. Markdown links: plain text (default, G-017's wording) or OSC 8 terminal
   hyperlinks that a terminal can open on click.

Everything else above is a routine technical choice inside G-043's outcome.

## Limits

Dependency graph visualization stays later, as G-043 says. Review actions are
G-044. Search matches ID, type and title, not bodies. Colour is an accent on
top of a design that reads without it; the owner judges both in a terminal.
