---
id: "G-164"
type: plan
title: "G-153 board body search and review listing"
status: current
created: "2026-09-25T20:40:14Z"
updated: "2026-09-25T20:40:35Z"
work: ["G-153"]
---

## Scope

[G-153](G-153-search-and-code-links.md) as narrowed at preparation: its
Next made the owner's rule that if agents already find and apply the
listed constraint without search, preparation narrows it to the board
search and the review listing. [G-154](G-154-listed-constraint-eval.md)'s
`without` row, reviewed in G-160 on `worktree-G-154` at `5ff6eb3`, found
both constraints applied in 10 of 10 runs without search, and named no
lever. So this plan covers G-153 acceptance 1 (board `/` over bodies), 3
(the review listing), the first observation of 5, and 6. `grove search`,
the guide and reviewer sentences, and the nullsec observation through
`grove search` are out, recorded in G-153.

## Design

- **Where the matcher lives.** `internal/handoff` already owns the body
  split, the goldmark link walk and `resolve`. A new exported
  `handoff.Mentions(r)` walks a record's body once and returns each real
  link that resolves to a project path, and each code span's text, with
  the body line it sits on. Everything else, the tiers and their order,
  is one function in `internal/tui/search.go` used by both surfaces:
  `match(query, record, mentions) (tier, snippet)`, tiers `title`, `link`,
  `code span`, `text` in that order, first that applies wins, no score.
  - `title`: the existing ID, type, status, title substring test,
    case-insensitive.
  - `link`: a mention's path equals the cleaned query or lies under it as
    a directory.
  - `code span`: the cleaned span equals the query, or equals a trailing
    run of the query's `/` components (`ws.rs` names
    `crates/server/src/ws.rs`).
  - `text`: a body line contains the query, case-insensitive; that line
    is the snippet.
  Snippets and reasons reach the screen only through `line`, which
  escapes.
- **Cache.** Mentions are parsed once per record revision into a map the
  model drops with the rendered-Markdown cache on each new result; nothing
  is written anywhere.
- **Board `/`.** In the current view each current state of a record
  (`currentStates`) is matched on its own, so a divergent record shows
  every current content, each hit naming its source (`label`); in one
  checkout's board the shown record alone. Hits sort by tier, then by the
  inspection's order. With a query each hit takes two rows: the existing
  row, then `tier · source: snippet`; with an empty query one row, as
  today. Enter opens the record's detail, as today.
- **Review listing.** In the detail's Changes section, under each changed
  file, one plain (not selectable) row: `  described by G-140 link,
  G-121 code span` or `  no record names it`. The file's path is made
  project-relative by stripping the result's prefix; renames match both
  sides. Only `link` and `code span` count, over every record's current
  states except the open record, each ID once. It uses the loaded records
  and the changes read the view already makes: no Git process, no write.

## Steps

1. `handoff.Mentions` with a unit test (links resolved, spans, lines,
   fences and images ignored).
2. Matcher, cache and board `/` with tests: body word, tier order, path
   query with link and span, suffix span, divergent record, hostile
   snippet escaped. Update [board.md](../docs/board.md) Search.
3. Review listing with a test (link, span, none, prefix, self excluded,
   no extra backend call). Update board.md's review section.
4. Evidence for acceptance 5: the review listing's matcher over the five
   merges in G-153 Constraints, by a throwaway harness outside the commit.
5. AGENTS.md checks, terminal lifecycle script, independent review, hand
   off into Review for the owner's layout judgment.
