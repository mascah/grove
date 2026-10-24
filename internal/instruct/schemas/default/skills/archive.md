**For a finished change, archive first brings the maintained specs up to
date,** so the code, the specs and the archived change land in the target
in one merge.

### Reconcile the specs

Apply each change spec to its maintained spec, `specs/<capability>/spec.md`
under the document root. Copy each ADDED requirement in, and each MODIFIED
one over the requirement of the same name, word for word, scenarios and
all; delete each REMOVED one; change no other requirement. The
ADDED/MODIFIED/REMOVED headings never appear in a maintained spec. Remove the
known discrepancies the work resolved. Where there is no maintained spec,
create it from the template under "Documents" below.
At the end of the `## History` section of each maintained spec the change
has a change spec for, adding the section last if it has none, add the line
`- [<change-id>](../../changes/archive/<folder>/) <title>`, with the
change's folder name and title. History is not a requirement; `archive`
does not compare it.
List each reconciled spec under "Reconciliation" in `implementation.md`. Then
run Grove's `check` and fix what it reports.
