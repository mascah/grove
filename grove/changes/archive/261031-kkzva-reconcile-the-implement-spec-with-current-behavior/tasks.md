## T-01 — The text around implement's requirements says what it covers today

Result: implement's Purpose and Boundaries read as the proposal gives
them, and init's Boundaries no longer points to a requirement this change
removes.

Files:
- `grove/specs/implement/spec.md`
- `grove/specs/init/spec.md`

Acceptance:
- implement's Purpose and Boundaries match the proposal's "The text
  around the requirements" word for word.
- init's Boundaries names no "No apply skill".
- `grove check` reports no problem, and archiving the change folds its
  requirements without a refusal.
- `implementation.md`'s `## Scenarios` names, for each scenario of the
  change spec, the test or instruction text that shows it holds today.
