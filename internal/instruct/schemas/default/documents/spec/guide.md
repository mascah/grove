A maintained spec states the behavior agreed for one capability. It becomes
agreed only when the owner merges it into the target branch.

- The capability's ID is lowercase kebab-case, such as `record-editing`; use
  the change spec's when there is one.
- The first heading is the spec's title, and the first sentence of
  Purpose its description: make that sentence read as one line saying
  what the spec covers.
- State intended behavior, not whatever the code does. Go through each rule
  with the owner before writing it. Where the code differs, keep the rule and
  record the difference under "Known discrepancies" with its source.
- Each requirement is one EARS sentence about behavior a user or caller can
  observe, followed by at least one scenario at the boundary that matters.
- Cover only the area at hand.
- End with `## History`: one line per change folded in, in the order they
  were archived, linking its archived folder, where its reasons are.
