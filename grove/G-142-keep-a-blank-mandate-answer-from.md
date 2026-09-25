---
id: "G-142"
type: work
title: "Keep a blank mandate answer from approving spend"
status: proposed
created: "2026-09-25T03:47:15Z"
updated: "2026-09-25T03:47:30Z"
relates_to: ["G-135", "G-139", "G-141"]
---

## Outcome

A work session never spends on an answer the owner did not give: a question
about what a paid step spends (model, effort, runs, budget or cap) takes
only explicit values, and a blank or missing item keeps it open and
blocking.

## Constraints

Observed at `worktree-G-135` `48358b7`:

- [G-139](G-139-what-mandate-and-login-should-th.md), written by a
  headless `/grove-work G-135 --until plan` session, said "any item left
  blank takes the recommendation"; the owner's answer covered the login
  only, and the continuation ran 9½ paid Codex runs on the recommended
  `gpt-6-astra` at `high`, 12% to 95% of the ChatGPT Plus five-hour window.
  The owner's rule is [G-141](G-141-never-run-gpt-6-astra-unless-the.md).
- [docs/work-shaping.md](../docs/work-shaping.md) already says "Your
  recommendation is not a decision. Silence is not agreement."
  [docs/work-execution.md](../docs/work-execution.md) has no equivalent: its
  step 5 and "When a human decision is missing" say how to ask, not that a
  resolved question's blank item is still unanswered.
- G-135 excluded changing the guides, so this is separate work.

Proposed design: the work guide says a question must not offer a default
for a blank item that spends, and that before acting on a resolved question
the session checks every item the step depends on has the owner's value,
treating a missing one as a missing human decision.

## Acceptance

1. `docs/work-execution.md` states the rule where it asks and where it
   resumes on a resolved question, consistent with the shaping guide's.
2. The guides digest changes and `go run ./cmd/grove check` passes.

## Next

Captured 2026-09-24 from the G-135 incident; not assigned.
