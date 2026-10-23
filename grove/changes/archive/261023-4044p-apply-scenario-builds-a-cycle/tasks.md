## T-01 — Rewrite the apply scenario

Result: the maintained implement spec's scenario under "Apply asks no one"
builds C as a member of a cycle, as
[specs/implement/spec.md](specs/implement/spec.md) says, and the spec
records no known discrepancy.

Acceptance:
- After archiving, `grove/specs/implement/spec.md` holds the requirement
  as the change spec states it, and its "Known discrepancies" section
  says none is known.
- `grove check` reports no problems; no code and no documentation page
  changes.
