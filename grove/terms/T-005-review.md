---
id: "T-005"
type: term
title: "Review"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:33Z"
relates_to: ["T-003", "T-004", "T-006"]
---

## Meaning

An examination of a [candidate](T-004-candidate.md) against the work's
acceptance and evidence, recorded with what it looked at. An independent agent
review looks for defects the implementer missed; a human review leads with
changed behavior, decisions, open issues, and checks before any diff.

A review is evidence, never a verdict by itself: it is not
[approval](T-006-approval.md), and a self-check is not an independent review.
The word also names a selected future work status between active and done,
which is not implemented; say "review record" or "Review status" when the
difference matters.

## Relationships

A review record (`R-NNN`) names its work and the commit it examined. Feedback
that starts another [attempt](T-003-attempt.md) keeps earlier reviews.
