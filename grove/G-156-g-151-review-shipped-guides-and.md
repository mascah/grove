---
id: "G-156"
type: review
title: "G-151 review: shipped guides and model free of Grove pointers"
status: current
created: "2026-09-25T19:31:29Z"
updated: "2026-09-25T19:31:41Z"
work: ["G-151"]
examined: "02a0e08"
---

## Examined

Three independent `grove-reviewer` rounds on `worktree-G-151`, each a fresh
agent, read-only, against
[G-151](G-151-strip-grove-repository-pointers.md)'s acceptance (no plan: the
record is its own):

- Round 1: `28aaf95..114db6b`.
- Round 2: `28aaf95..be4efcb`, the fix commit read line by line.
- Round 3: `28aaf95..02a0e08`, the fix commit read line by line.

Each ran `go test -short ./internal/cli`, `go vet ./...`, `gofmt -l .` and
`grove check`, all passing, rendered `grove guide work|shape|model` and
mutation-checked the test in a `/tmp` copy; none ran the full suite, which
the author ran (G-151 Evidence).

## Findings

Round 1 (4 notes, none blocking):

1. The phrase check was case-sensitive and narrow: "Command reference",
   "The Predecessor tool" and "Grove's repository" passed.
2. [G-146](G-146-how-should-an-adopting-project-r.md)'s Answer says the
   model cites the command reference by name; G-151 drops that, and a later
   `https://` link to it ([G-110](G-110-external-preview.md)) would also
   need the test's denylist changed.
3. "predecessor" is a common English word a future sentence could use.
4. AGENTS.md's "never names a record" read as forbidding the example IDs.

Round 2 (1 should-fix, 3 notes): a phrase split across a line break passed;
a typographic apostrophe passed; the reworded AGENTS.md sentence was
ungrammatical; Evidence still to be written.

Round 3 (no blocking or should-fix): all round 2 items but Evidence
resolved and mutation-checked (split lines, tabs, U+2019, U+00A0, upper
case, "Predecessors"); ID lists confirmed per document; no false positive
in the current documents. Two notes: `guides.go`'s comment says the model
"names no record" though it carries example IDs; Markdown between the words
("Grove's **repository**") still passes.

## Disposition

- Fixed in `be4efcb`: round 1's 1 and 4. Fixed in `02a0e08`: round 2's
  should-fix, apostrophe and wording.
- Round 1's 2: G-146 stays as answered; G-151's Evidence records the
  reversal and its Next the G-110 consequence.
- Not changed: round 1's 3 (a failure names its cause); round 3's comment
  note (example IDs are not records, so it is true as read) and Markdown
  note (a phrase denylist's accepted limit, like G-149 finding 7's example
  ceiling).
