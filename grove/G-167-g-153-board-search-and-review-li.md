---
id: "G-167"
type: review
title: "G-153 board search and review listing review"
status: current
created: "2026-09-25T20:54:27Z"
updated: "2026-09-25T20:54:44Z"
work: ["G-153"]
examined: "702b576"
---

## Examined

Round 1: an independent `grove-reviewer` agent, dispatched by the headless
`/grove-work G-153` session, on `worktree-G-153` from main `3f2b923` to
`702b576` (narrowing, plan [G-164](G-164-g-153-board-body-search-and-revi.md),
implementation), against [G-153](G-153-search-and-code-links.md)'s
acceptance and constraints. It read the diff, ran `go vet`, `gofmt -l`,
`grove check`, the tui and handoff packages uncached, and experiments in a
disposable clone it removed; it wrote nothing in the checkout.

## Findings

1. **Consequential, fixed.** The review listing stripped `Prefix + "/"`,
   but `Result.Prefix` comes from `git rev-parse --show-prefix` and already
   ends in a slash, so in a project below the repository root every link
   match was lost and only code spans matched. The test used `"proj"`
   without a slash and passed. Fixed in `7f00b39`: strip `Prefix` as Git
   gives it, a file outside the project matches by code span only, and the
   test uses `"proj/"`, a rename and a file outside the prefix; reverting
   the fix fails it.
2. **Consequential, fixed.** Acceptance 5's evidence was not in the record.
   Written into G-153's Evidence; the reviewer reproduced the link pairs of
   the five merges independently (3, 1, 9, 0, 5 over 4, 1, 9, 0, 6 code
   files). Plan G-164's scope sentence, which still called the nullsec
   observation out, is corrected.
3. **Minor, fixed.** No test for a rename, a real prefix, or a checkout's
   own board: added in `7f00b39`.
4. **Minor, fixed.** PgUp/PgDn paged by rows while each hit takes two:
   search now pages by visible hits (`7f00b39`).
5. **Informational, noted.** The listing is recomputed per frame, about
   21 ms for 60 files over 153 records; a `ponytail:` comment names the
   ceiling and the cache to add.
6. **Informational, no change.** The narrowing follows the owner's rule
   in G-153's former Next faithfully; G-160 is judge-scored and not yet
   integrated, and the owner can reverse the narrowing when judging.

Knowledge: no term or decision is contradicted; "describes a file" is
board vocabulary defined in board.md and needs no term.
