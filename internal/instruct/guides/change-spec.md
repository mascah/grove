Write one change spec for each affected capability, inside the change folder.

- When a maintained spec exists for the capability, or a baseline is drafted
  on this branch, use its ID and link it. New rules go under ADDED. A changed
  rule goes under MODIFIED, written as the full intended requirement with all
  of its scenarios. A dropped rule goes under REMOVED with the reason.
- When no maintained spec exists, name the capability by choosing its spec ID,
  put every proposed requirement under ADDED, and describe what the code does
  today under "Observed current behavior", each observation with its source
  path and revision. Do not use MODIFIED or REMOVED: nothing agreed exists to
  modify. Do not present the observed baseline as agreed behavior.
- Describe the baseline relative to what the proposal is reviewed against,
  never as a claim about the target branch that a merge could make false:
  "the baseline proposed alongside this change in specs/<id>/spec.md", or
  "no maintained spec at <base commit>".
- Each requirement is one EARS sentence about behavior a user or caller can
  observe, followed by at least one scenario at the boundary that matters.
- Do not edit the maintained specs to describe this change. Future behavior
  stays in the change until it is implemented and archived.
