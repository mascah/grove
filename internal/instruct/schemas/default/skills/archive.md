**For a finished change, archive first brings the maintained specs up to
date,** so the code, the specs and the archived change land in the target
in one merge.

### Reconcile the specs

Apply each change spec to its maintained spec, `specs/<capability>/spec.md`
under the document root: copy each ADDED requirement in, and each MODIFIED
one over the requirement of the same name, word for word, scenarios and
all; delete each REMOVED one; change no other requirement. The
ADDED/MODIFIED/REMOVED headings never appear in a maintained spec. Remove
the known discrepancies the work resolved. Where there is no maintained
spec, create it from the template under "Documents" below. List each
reconciled spec under "Reconciliation" in `implementation.md`, then run
Grove's `check` and fix what it reports.
