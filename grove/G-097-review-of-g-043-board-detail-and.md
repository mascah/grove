---
id: "G-097"
type: review
title: "Review of G-043 board, detail and search"
status: current
created: "2026-09-22T23:57:39Z"
updated: "2026-09-22T23:57:46Z"
work: ["G-043"]
examined: "ed11f51"
---

## Examined

Branch `worktree-G-043`, base `main` dc3b9b6. An independent reviewer
(a subagent given the work record, plan [G-096](G-096-g-043-board-and-detail-design-vi.md),
[G-017](G-017-terminal-picker.md)'s contract and the diff, with no write
access) examined `git diff dc3b9b6..b020601` in round 1 and ed11f51, the
fix commit, in round 2, which is the `examined` field. Its probes were Go
test files kept outside the worktree and added through `go test -overlay`.
The candidate differs from ed11f51 by this record, [G-043](G-043-board-detail.md)'s
evidence and a README wording change (disposition of finding 6).

## Findings

1. **Blocking: HTML character references reached the terminal decoded.**
   `escapeLines` ran `safe` before glamour, but Markdown decodes `&#x1b;`
   afterwards, so a body could emit OSC 52 (clipboard), OSC 2 (title), DCS,
   SGR (`[8m` conceals text) and U+202E, in text, code spans, headings, HTML
   blocks and table cells; bubbletea's renderer writes every non-SGR
   sequence to the terminal. The hostile-input test had no reference case.
2. **Medium-low: a shrinking resize left focus on a Done card the smaller
   page hides**, drawn without `▶`, with Enter doing nothing.
3. **Low: a record under the open one vanishing on refresh** left the stack
   holding it: the board drew but keys went to the detail until another Esc.
4. **Low: below 100 columns Tab could not reach the sidebar** of a record
   with no linked records and no commits, so Sources and a history error
   were unreachable.
5. **Low: the G-043 record was stale** (Next said implementation had not
   started; lipgloss v2.0.6 named where go.mod has v2.0.4).
6. **Low, cosmetic (round 2):** README said a reference "shows as typed",
   which holds in prose and code spans, while a fenced or indented code
   block or a link target shows the inserted `&amp;`.

No finding: titles, IDs, commit subjects and the search query stay escaped;
the OSC 8 regexp covers both terminators; no Git process per branch; history
is read only with a detail open and any key cancels it; the versions screen
and selector resolution are unchanged; the Done page and its footer, the
Abandoned toggle, the Esc stack and the examined-marker prefix match hold;
nine sizes from 40x10 rendered every screen without a panic; the largest
record (31 KB) renders in 8.3 ms; `TestTerminal` keeps its `-short` skip.

## Disposition

1. Fixed in ed11f51 with two layers: every reference is made literal before
   glamour, and each rendered row keeps only glamour's SGR styles
   afterwards. Round 2 reran every round-1 probe and a bypass set (double
   encoding, inline HTML, link titles, alt text, no semicolon, `&#0;`,
   surrogates, uppercase hex, shortcodes) through the full detail view: no
   control byte or planted style survived. Every record in `grove/` rendered
   at three widths shows no new escape.
2. Fixed: a resize settles the focus; tested at three sizes.
3. Fixed: every vanished entry leaves the stack; the detail closes only when
   it empties; tested for a lower and a top record.
4. Fixed: narrow Tab shows the sidebar with nothing to select; tested at 80x24.
5. Written in the handoff below the examined commit, as planned.
6. README reworded in the handoff commit; the pre-escape stays uniform since
   the second layer alone would still let a decoded SGR through in code spans.
   No record in this repository carries a reference in a code block.

Verdict: clear with notes (round 2). This record holds evidence; it is not
the owner's judgment, which acceptance 1 and 3 still require in a terminal.
