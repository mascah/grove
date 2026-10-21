# Persistent implementation runs

Agent-authored proposal for review, 2026-10-21. The behavior and amendment
wording below are recommendations. [CONSTITUTION.md](../CONSTITUTION.md)
remains authoritative; this draft is separate from the current user
specification in `docs/`.

The owner's instruction for this first version:

> I think we need to defer ephemeral execution until the simplest possible solution works

## Proposed first version

An operator starts one bounded implementation attempt in a dedicated,
persistent checkout. Grove runs the existing apply process and retains a
record and output for inspection. Once this works manually, an external
timer can invoke the same operation. Different projects run independently.

The operator supplies one designated checkout per project, credentials and
the execution environment. macOS and Claude Code are the recommended first
environment and harness; those are assumptions to confirm.

## Ordinary run

1. Grove prevents overlapping invocations for the designated project and
   checks its previous attempt before starting anything new.
2. In an idle, clean checkout, Grove fetches the remote target and updates
   the local target reference its instruction code reads. It selects the
   approved batch using existing dependency ordering. No eligible work
   means no attempt.
3. Grove creates a branch and records the attempt before launching one
   agent session. The selected target snapshot stays fixed during the run.
4. The agent implements the batch, records task progress and archives
   finished changes. Existing project rules cover verification, commits
   and PR creation. Several changes can share one branch and PR.
5. Grove captures output and records the process outcome. Process success
   does not establish implementation success; task progress and the PR
   describe the work.

After a normal exit, the project stays paused until every selected change
is archived as finished on the refreshed target, or the operator
explicitly resolves or retries the attempt. Unrelated commits and newly
approved changes do not release the pause. This needs no PR-state polling.

## Interruption and partial work

Each attempt has one session and a wall-clock limit. On cancellation or
timeout, Grove terminates the agent's process group before releasing
execution exclusion. Its record, branch, uncommitted files and output stay
available. An interrupted attempt always needs operator attention before
another start. After an abrupt supervisor failure, the operator also
checks for surviving processes.

The operator inspects and prepares the checkout before explicitly
retrying. Retry uses the existing apply process, which skips tasks marked
done and changes archived there. Automatic restoration and checkout repair
are deferred under the proposed rule 11 boundary below.

For example, A can finish while B is blocked by a question. The agent can
open a draft PR containing A and explain B's blocker. Automation pauses;
the owner can review useful work, revise the change and explicitly permit
another attempt. A publication failure similarly leaves work available
for operator attention.

## Local record

Outside the implementation checkout, retain the selected changes, starting
target commit, checkout, branch, start/end times, process outcome and output
location. Save updates atomically and retain resolved records for inspection.
`implementation.md` continues to own task progress. Operational records
establish neither task completion nor approval.

## Constitutional changes to propose

The constitution requires a separate amendment PR, read by the owner
before code follows. The proposed changes are:

- **Rules 1 and 2:** allow opt-in bounded execution in an operator-prepared
  checkout; retain print-and-exit behavior for `instructions apply`.
- **Rules 6 and 7:** permit local operational records and output outside
  the implementation checkout; retain Markdown in Git as authoritative
  for specs, changes and task progress.
- **Rule 11:** permit bounded launching; defer session resumption,
  workspace creation, cross-host coordination, ephemeral execution,
  Grove-owned publication and scheduling.

Existing archive status and customization suffice under rules 3 and 10.
The proposal retains the code and dependency budgets in rules 12 and 16.

## Demonstrating it

Use the ledger fixture in real disposable Git repositories (rule 17):

1. Apply two dependent changes, inspect the output and implementation PR,
   and invoke Grove before merging: no second agent starts. Merge the
   finished batch; a subsequent invocation can consider new work.
2. Interrupt after a visible edit. Confirm process termination, retained
   files and output, and no automatic restart. Prepare the checkout and
   explicitly retry. Also exercise the time limit.

Focused lifecycle checks cover simultaneous starts and the unresolved
record left by an abruptly killed supervisor. Controllable worker
processes can verify supervision; a real agent run must demonstrate the
complete workflow.

## Remaining detail and cost

Before implementation, confirm the first harness and document concrete
start, inspection and retry commands, host settings and the default time
limit in the user specification and walkthrough (rule 8).

Expected code areas: one execution package, CLI wiring and small Git
helpers, reusing the document and instruction code. The preliminary sizing
judgment is 600–1,200 production Go lines plus focused tests, assuming one
harness, an existing checkout and agent-owned publication. Process-group
termination needs a fixture trial before tightening that estimate.
