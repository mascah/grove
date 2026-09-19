# Executing prepared Grove work

This is the restart's reusable implementation guidance for an explicitly assigned
work ID or set of IDs. It is ordinary repository documentation, not an installed
`/work` skill, a prompt-generation command, or an agent runner. The active user
assignment supplies authorization and scope. Reading this file alone does not
start work or authorize a launch, merge, or push.

This guide is a partial baseline, not a proven replacement for the predecessor's
controller, recovery, and closure workflow. The
[predecessor review](reviews/2026-09-19-predecessor-work.md) maps those
responsibilities for W-010; its deferred mechanisms are not instructions to
invoke unsupported commands here.

## Authority and scope

Read `docs/restart-brief.md` first, then `AGENTS.md`, `docs/record-model.md`, the
assigned records, their linked implementation plans, and relevant review evidence.
Repository/user instructions take precedence over predecessor workflows. Records
own outcomes, constraints, and acceptance; plans describe implementation steps;
the brief owns product direction. Preserve accepted/proposed/observed distinctions.

Use `go run ./cmd/grove` in the selected restart checkout, never the installed
predecessor executable. Use `show ID` to retrieve assigned records, `new` for new
records, and `update ID --expect REVISION` for field changes, getting the revision
from `show ID --json`. Edit explanatory bodies and linked plans as ordinary text.
Do not invent IDs, statuses, fields, schema, `grove.toml`, `docs/grove/`, or the
predecessor's capabilities/history trees. Do not invoke predecessor work/close
workflows. Leave sibling repositories unchanged.

## Establish the execution base

Inspect status, HEAD, branches, registered worktrees, and actual prerequisite
commits before editing. A done record does not prove integration. If assigned
work already has implementation, inspect its diff/evidence and resume or verify
it; do not create a duplicate implementation based only on stale Next prose.

Use an isolated worktree from the assignment's required base. Reuse an existing
one only when it is clearly this assignment's workspace and doing so preserves
all concurrent work. Never reset, clean, remove, or repurpose another session's
checkout. State the selected path, branch, base revision, and assigned IDs.
Follow explicit sequencing; independent IDs do not imply independent interfaces.

## Implement through evidence

An assignment to implement prepared work authorizes routine technical decisions
inside its documented outcome. Do not ask again for blanket implementation
permission. Ask only for a consequential product choice, incompatible scope or
contract change, or a genuine external blocker; continue independent useful work.
Update plans when evidence requires a bounded technical adjustment, retaining
why it changed. Do not quietly widen the assignment.

Set each work record active when its implementation starts. Reproduce specified
bugs with deterministic fixtures before repairing them, then implement against
the record's acceptance and the existing shared interfaces. Preserve unrelated
bytes/state and error semantics. Keep commits focused and Conventional. One
implementer should own overlapping interfaces; use independent reviewers at the
specified boundaries without letting them concurrently edit those interfaces.

Run the relevant targeted checks, `go test ./...`, `go test -race ./...`,
`go vet ./...`, formatting, and `go run ./cmd/grove check` as required by the
record/plan and repository. Prefer uncached test runs for final evidence. For
TUI work, run terminal lifecycle/connected-workflow checks too. Documentation
changes need link/consistency checks. Record actual commands, results, and tested
code revisions; distinguish new evidence from inherited reports.

## Review, reconcile, and return

Obtain the requested independent review and address consequential findings with
regressions. For a multi-record assignment, retain each unit's evidence and
review the final combined diff for regressions across shared helpers. If an
independent reviewer is unavailable, state that limit rather than labelling a
self-review independent. Do not conceal outstanding required verification.

Reconcile the assigned records' evidence and concrete Next actions, linked plans,
README/model when contracts need clarification, and the brief's next action.
Mark done through the CLI only when that record's acceptance is met. Do not claim
owner usability approval from automated checks or a screenshot; if required
judgment is outstanding, record it explicitly and leave the appropriate work
active. Commit evidence alongside the code. No new structured review/run schema
is needed: use prose and links.

Unless the assignment explicitly says otherwise, do not merge, push, deploy, or
remove worktrees. Branch completion and integration are separate. Return:

- Assigned IDs, branch/worktree, base and final code/evidence commits.
- Concrete behavior changes and acceptance evidence per record.
- Verification results, review findings/disposition, and unverified limits.
- Exact next action for integration or required owner judgment, including a
  runnable demo command for interactive work.

Clean up only this assignment's disposable probes/processes. Never leave an
untracked background agent running as an implied continuation of the handoff.
