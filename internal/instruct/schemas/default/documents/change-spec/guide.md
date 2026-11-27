Write one change spec for each affected capability, inside the change folder.

- When a maintained spec exists for the capability, use its ID and link it. New rules go under ADDED. A changed
  rule goes under MODIFIED, written as the full intended requirement with all
  of its scenarios. A dropped rule goes under REMOVED, as `### Requirement: <name>` with its
  name in the maintained spec, followed by the reason.
- When no maintained spec exists, name the capability by choosing its spec ID,
  put every proposed requirement under ADDED, and describe what the code does
  today under "Observed current behavior", each observation with its source
  path and revision. Do not use MODIFIED or REMOVED, but as the next bullet
  says: nothing agreed exists to modify. Do not present observed behavior as
  agreed behavior.
- When a change in `depends_on` adds or modifies a requirement of the
  capability, put that requirement under MODIFIED, written from the text
  that change gives it and naming that change, even where no maintained
  spec exists at `base`.
- Name what the proposal starts from by commit, "no maintained spec at <base
  commit>", never as a claim about the target branch that a later merge could
  make false.
- Each requirement is one EARS sentence about behavior a user or caller can
  observe, followed by at least one scenario at the boundary that matters.
- Keep a rule, the text before its first scenario, to 60 words, and leave
  no maintained spec over 2,500 words larger than you found it: `grove
  check` reports a change that does either.
- Do not edit the maintained specs to describe this change. Future behavior
  stays in the change until it is implemented and archived.
