---
id: "G-197"
type: term
title: "Policy"
status: proposed
created: "2026-09-26T03:29:15Z"
updated: "2026-09-26T03:29:23Z"
relates_to: ["G-057", "G-058", "G-059", "G-060", "G-180", "G-182", "G-192"]
---

## Meaning

The owner's standing instruction, written as `policy:` in `grove.yaml`, for
what may happen to a [candidate](G-057-candidate.md) in review with no
per-candidate human act: start one [resolution](G-192-resolution.md)
attempt for a conflict, and approve and integrate a candidate that meets
the written conditions once its merged result passed the policy's
verification. Its absence means nothing is automatic. `grove sweep` applies
it ([G-180](G-180-policy-driven-integration.md)), and each act it takes is
**delegated**: attributed to the policy's `grove.yaml` revision and the
evidence it relied on, told apart from the owner's own verdict, and
reversible by an ordinary revert. The delegation itself is decision
[G-182](G-182-standing-policy-delegation.md).

## Relationships

A delegated [approval](G-059-approval.md) is approval by someone the owner
delegated to, which that term already allows. A [review](G-058-review.md)
stays evidence: the policy reads its closing line, `Open findings: none`,
and the review gives no verdict. [Integration](G-060-integration.md) under
a policy is the ordinary merge and `done`, preceded by the verification.
