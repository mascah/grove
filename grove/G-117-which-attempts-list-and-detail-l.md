---
id: "G-117"
type: question
title: "Which Attempts list and detail layout should G-109 implement?"
status: open
created: "2026-09-23T19:51:07Z"
updated: "2026-09-23T19:53:32Z"
blocks: ["G-109"]
---

## Question

[G-109](G-109-attempts-usability.md) requires comparing concrete Attempts
layouts with the owner before implementing one. Plan
[G-116](G-116-g-109-attempts-layouts-observed.md) draws both at 80×24 from
this repository's real attempts and sanitized examples of the other states.
Answer three parts:

1. **Layout.** A: the list grouped into `Needs you`, `Running` and
   `Settled`, and one attempt showing work title, State, Next, then the
   report or latest activity, with the full facts folded behind `d`.
   B: one newest-first list with `!` on rows that need you, and the same
   attempt screen with the facts always shown below the report or activity.
2. **What needs you.** Proposed: only the latest attempt of work that is
   still proposed, active or review, when it is a candidate to judge, a
   blocking question, a failure, an orphan, an interruption or an end
   without a handoff. An earlier failure followed by another attempt, a
   stopped attempt, and any attempt of work now done or abandoned are
   settled, each saying why.
3. **Activity noise.** Proposed: count provider `system:` notices other than
   `init` (for example Claude Code's `thinking_tokens`) instead of listing
   them one per line; the raw log keeps them all.

## Evidence

G-116's Observed section: the current list shows the timestamped attempt ID
but not the work title, and says `candidate ready` for work integrated since.
At 80×24 the report starts on page two, after 18 lines of provenance. Most
of a running attempt's activity is `system: thinking_tokens`.

## Recommendation

A, with parts 2 and 3 as proposed. The list keeps growing, since every
attempt is kept, so grouping keeps what needs you at the top however many
settled attempts accumulate. With the details folded, the raw paths are one
key away however long the report is. B is less code and meets the acceptance
while attempts are few.

## Next

Open; blocks G-109. Answer here (for example "A as drawn", or which parts
change), set `status=resolved`, and relaunch G-109; implementation follows
G-116's Steps for the chosen option.
