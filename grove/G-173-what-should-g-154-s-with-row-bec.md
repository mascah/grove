---
id: "G-173"
type: question
title: "What should G-154's with row become now that G-153 shipped no agent-facing search?"
status: resolved
created: "2026-09-25T21:19:37Z"
updated: "2026-09-25T21:24:49Z"
blocks: ["G-154"]
relates_to: ["G-153", "G-158", "G-160"]
---

## Question

[G-154](G-154-listed-constraint-eval.md) acceptance 2 runs each case "on
its `with` row once G-153's candidate is integrated", and acceptance 3
compares the rows. G-153 is integrated on main (`1db5a75`, done at
`47852e3`), but its preparation narrowed it to the board's body search
and the review listing ([G-153](G-153-search-and-code-links.md), Scope at
preparation): no `grove search` command, no guide or reviewer sentence.
A headless shaping run cannot use the board, so a `with` row at main
would measure the guide edits of G-150 and G-151 that landed on main
meanwhile, not search. G-153's own record says so ("G-154's `with` row
compares nothing new"), and the owner's verdict on its candidate was
"fine with deferring the grove search command for now". G-158's mandate
("Use the G-122 settings") was given for the `without` row; its item 2
names the `with` row only for its run count.

This session did not spend, since which of these G-154 does is a scope
change to its acceptance, not a routine choice. Options:

1. **Drop the `with` row from G-154.** Amend acceptance 2 and 3 to the
   `without` row; the review completes on that row with the recomputed
   facts, and G-154 is handed into Review. A later search command or
   guide sentence, reshaped as its own work, carries its own `with` row
   on these cases, whose runner and `without` row stay on main. Cost:
   nothing further.
2. **Run the `with` row at main anyway**, under G-158's settings (about
   $3.35), as a check that main's guide edits since digest
   `41324c3655a1` neither break these cases nor add cost, then complete
   the review across both rows and hand off. It says nothing about
   search. The G-108 pair is already the regression rerun for a guide
   edit (G-114 Evidence).
3. **Keep G-154 active**, waiting for a future work that ships `grove
   search` or the guide sentences, and run the `with` row then.

Recommendation: 1. G-160 found no miss for search to recover at these
fixture sizes, so neither 2 nor 3 can show what G-154's outcome asks;
evidence for search needs the larger-fixture case G-160's Disposition
names, which is its own owner choice.

Who can answer: the owner. Answer below, set `status=resolved`, then
assign `/grove-work G-154` again on `worktree-G-154`.

## Next

Waits for the owner.

## Answer
Drop the with row
