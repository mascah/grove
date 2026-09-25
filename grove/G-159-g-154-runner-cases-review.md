---
id: "G-159"
type: review
title: "G-154 runner cases review"
status: current
created: "2026-09-25T20:07:54Z"
updated: "2026-09-25T20:08:15Z"
work: ["G-154"]
examined: "47223ed"
---

## Examined

The runner built for [G-154](G-154-listed-constraint-eval.md) acceptance 1,
at the consequential boundary of plan
[G-155](G-155-g-154-listed-and-code-constraint.md) step 3: the two cases,
their fixture records, the retrieval facts `search`, `holding read` and
`distractors read`, per-case rubric columns, and the README. Three rounds by
fresh `grove-reviewer` agents on `worktree-G-154`, 2026-09-25, each told
it may run the free selftest and never a paid harness: round 1 examined
`b684951..b2c7ef0 -- evals/`, round 2 `b2c7ef0..946459a` and the combined
change, round 3 `946459a..47223ed` and the combined change, the `examined`
commit. `c895666` came after the cap and is self-checked only (below).
This is evidence for the runner, not the eval's result: acceptance 2 and 3
wait on the mandate in
[G-158](G-158-what-mandate-should-the-g-154-wi.md).

## Findings

Round 1 (at `b2c7ef0`):

1. Medium. A holding record read through `grove context ID --include PATH`,
   the listing's own advice, counted as unread. Fixed in `946459a`: an
   `--include` path on a `grove` command is a file read.
2. Low to medium. The brief lists four frontmatter keys, so a careful run
   might ask whether a `due` key is allowed and fail `no-question` for a
   reason the case does not test. Fixed: the due-date record's owner note
   says "a `due` frontmatter key".
3. Low. "The decision's title shares no word with the topic" was false
   ("task"). Wording corrected in the README and G-155.
4. Low. Selftest gaps: `context` by ID, a listing not being a reading, a
   distractor read through `show`, the default case set, the decision's
   status. Each now exercised by the fake and asserted; the default-set
   and include assertions were mutation-checked.
5. Info. `distractors read` counts `show`, `unneeded` counts files only.
   Documented: compare the two together.
6. Low. "Claude only" read as a runner rule. Reworded as a rule of G-154
   and G-141.
7. Low. The listed-constraint rubric's top anchor was ambiguous for a
   proposal citing only the brief. Scored 1, with `holding read` noted.

Round 2 (at `946459a`):

- A. Medium, introduced by the round 1 fix's selftest. `context` on a
  decision is refused by the CLI but counted as reading it. Fixed in
  `47223ed`: `context` by ID counts only for a `work` record; the fake's bad
  code-constraint run tries it and is asserted unread. Mutation-checked.
- B. Low. A missing `--include` file counted. Filtered to files in the
  clone; that a refused command still counts is documented as a limit.
- C. Low. The spaced `--include PATH` form was unexercised. The Claude fake
  uses it, the Codex fake `--include=`. Mutation-checked.
- D. Info. G-158 had duplicate empty headings. Removed.

Round 3 (at `47223ed`): every round 2 disposition verified, two by the
reviewer's own mutations; the pair's fixture, default and columns
unchanged; G-158 accurate. Two new points:

- Low. An absolute `--include` path counted as a read, though `context`
  refuses one. Fixed after the cap in `c895666` and self-checked only: a
  scratch trace with the relative path counts the file, the absolute one
  does not; `selftest: ok`.
- Info, open. No fake passes an `--include` of a missing file, so deleting
  the existence filter would still pass the selftest.

Knowledge: no term needed ("holding record", "distractor" and "row" are
the eval's own words); nothing contradicts G-141, G-064 or the rule that
fixture records are made only in disposable directories.

## Disposition

Acceptance 1 and 4 are met by the runner at `c895666`. Every consequential
finding is fixed; one informational coverage gap stays open. The runner is
ready for the `without` row once G-158 is answered. Limits: no real
harness has run the new cases, so whether a real agent's reads land in the
facts as the fakes' do is unobserved; the facts come from commands, not
their results, and `grep` output is never a read, so any `holding read`
the review of the runs rests on should be checked against the transcript.
