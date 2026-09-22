---
id: "G-058"
type: term
title: "Review"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:33Z"
relates_to: ["G-056", "G-057", "G-059"]
formerly: "T-005"
---

## Meaning

An examination of a [candidate](G-057-candidate.md) against the work's
acceptance and evidence, recorded with what it looked at. An independent agent
review looks for defects the implementer missed; a human review leads with
changed behavior, decisions, open issues, and checks before any diff.

A review is evidence, never a verdict by itself: it is not
[approval](G-059-approval.md), and a self-check is not an independent review.
The word also names the work status between active and done, in which a
candidate awaits human judgment; say "review record" or "Review status" when
the difference matters.

## Relationships

A review record names its work and the commit it examined; work in Review
status names its `candidate`, so the two can be compared. Feedback that
starts another [attempt](G-056-attempt.md) keeps earlier reviews.
