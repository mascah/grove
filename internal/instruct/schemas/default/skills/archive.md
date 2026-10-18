**For a finished change, archive first brings the maintained specs up to
date,** so the code, the specs and the archived change land in the target
in one merge.

### Reconcile the specs

Apply each change spec to its maintained spec, `specs/<capability>/spec.md`
under the document root. ADDED requirements become requirements, MODIFIED
ones replace the requirement of the same name, REMOVED ones are deleted; the
ADDED/MODIFIED/REMOVED headings never appear in a maintained spec. Remove the
known discrepancies the work resolved. Where there is no maintained spec,
create it, and its `.grove.yaml`, from the templates under "Documents" below.
List each reconciled spec under "Reconciliation" in `implementation.md`. Then
run Grove's `check` and fix what it reports.
