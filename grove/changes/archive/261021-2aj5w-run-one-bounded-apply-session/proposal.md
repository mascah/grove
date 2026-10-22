# Run one bounded apply session

## Problem

Approved changes are built only when someone opens their agent and types
`/grove-apply`. The owner wants approved work started without anyone
watching, the same way in every project, and with execution in Grove rather
than in a script each project copies (handoff note
`memex/projects/grove/2026-10-20-grove-managed-execution.md`).

At dca5543 Grove cannot do that. Rules 1, 2 and 11 of `CONSTITUTION.md` say
it launches no agent, `grove instructions apply` prints and exits, and the
harness table in `internal/cli/init.go` holds only where each harness reads
its skills. A timer that runs the agent itself has to handle overlapping
starts, runaway sessions, stray child processes and repeating unfinished
work on its own. The agent-authored draft `dev/managed-execution-proposal.md`
(at 696c8b4) sets out the simplest version. This change specifies it.

## Intended outcome

In a checkout the operator set up for one project, `grove run` updates the
target from `origin` and selects what no-ID apply would build. It creates a
branch, launches one Claude Code session on the apply skill, and waits for
the session within a time limit. The agent builds, archives, commits and
opens the pull request, as `rules` in `grove.yaml` tell it. Grove keeps a
record of the attempt and its output under `.git/grove/` and pauses the
project. The next `grove run` launches nothing until every selected change
is archived on the target, or until the operator runs `grove run --retry`.
On timeout or cancellation, Grove stops the session's process group, keeps
the branch, files and output, and waits for the operator.

Different projects, each in its own checkout, run independently. The rules
are in [specs/runs/spec.md](specs/runs/spec.md).

## Scope

### In scope

- `grove run` and `grove run --retry`, with a `--timeout` flag, for Claude
  Code on macOS.
- One run at a time per checkout; the time limit; stopping the session's
  whole process group on timeout or cancellation.
- The attempt record and output under `.git/grove/`, and the pause.
- Documentation, walkthrough steps and a guide for running `grove run` on a
  launchd timer.

### Out of scope

- Recovering or resuming an interrupted attempt, and repairing a checkout.
  The operator does this.
- Creating workspaces or worktrees, disposable runners such as GitHub
  Actions, and coordinating runs across hosts (amended rules 2 and 11).
- Grove pushing branches or opening pull requests. The agent does this,
  steered by `rules`.
- Grove scheduling itself. launchd or cron calls `grove run`.
- Harnesses other than Claude Code, and passing permission flags to the
  harness. The operator configures permissions in the harness's own
  settings.
- Running several changes of one project in parallel, or watching
  pull-request state.

## Decisions

- **Execution lives in Grove's CLI, and the constitution is amended for it.**
  The owner, 2026-10-20 (handoff note): "If the runner is gong to be
  inherently coupled to grove, I'd be willing to consider changing grove's
  constitution to support those responsibilities. Grove's go CLI is a good
  footing for that type of orchestration." In this session: "The
  constitution changes and such have already been answered", and the
  amendment was drafted as its own pull request, #17, at the owner's
  request: "draft it as its own pr". The code waits for #17 to be merged.
- **Simplicity over unattended recovery.** The owner, 2026-10-21 (handoff
  note): "In essense our implementation must be simple, even if it doesn't
  handle every edge case right now." Reason: an earlier build leaned too far
  toward orchestration. Interrupted attempts stop and wait for the operator.
- **No ephemeral execution yet.** The owner, quoted in
  `dev/managed-execution-proposal.md`: "I think we need to defer ephemeral
  execution until the simplest possible solution works".
- **macOS and Claude Code first.** The owner, in this session, answering
  "macOS and Claude Code only, with GitHub Actions later?": "yes".
- **Commands and storage.** The owner, in this session, on `grove run` /
  `grove run --retry` with records under `.git/grove/`: "fine for now".
- **Agent behavior goes through existing customization.** The owner,
  2026-10-21 (handoff note): "yes, we can use grove's override
  functionality like rules or schemas to pass addiitional instructions to
  the agent about behavior." Committing, opening the pull request and
  reporting blocked work are `rules`, not new configuration.
- **Not tied to Hermes.** The owner, 2026-10-21 (handoff note): "This
  functionality should work with or without hermes."

### Assumptions to confirm in review

These are the agent's recommendations, not the owner's decisions. Merging
this change agrees to them.

- Grove fetches the target from the remote `origin` and fast-forwards the
  local target to it. It refuses when it cannot fast-forward, and never
  moves the local target off commits of its own.
- Grove reads `grove.yaml` and selects changes from the fetched target, so
  the attempt's record and the agent's instructions describe the same work.
- The default time limit is 2 hours.
- The branch is `grove/run-<attempt>`. `<attempt>` is a new ID in the form
  Grove gives changes (`doc.NewID`), so quick retries never collide.
- A dropped change releases the pause just as a finished one does.
- An interrupted attempt is paused by the same rule as any other. The
  files it leaves make the checkout unclean, so `grove run` refuses to start
  until the operator deals with them.
- `--retry` starts a new attempt from the fetched target, not from the
  previous branch. To keep a partial branch's finished work, the operator
  merges its pull request first; otherwise the agent builds those tasks
  again.
- Grove launches `claude -p` with the apply skill as the project's `agents`
  form types it (`/grove-apply` for skills, `/grove:apply` for commands) and
  passes no permission flags. An omitted `agents` lists Claude Code, as
  today.
- A session that cannot be started is recorded as failed to start and
  pauses the project like any other attempt. A timer does not retry it.
- Stopping a session stops its process group. A process that leaves the
  group, such as a server that starts its own session, is not stopped.
- The lock is an operating-system file lock, so a killed `grove run` leaves
  no stale lock. Its agent may survive and a clean checkout would let
  `--retry` start a second one beside it. The operator checks for a
  surviving agent before retrying; Grove does not.
- The agent pushes and opens the pull request only because the project's
  `rules` tell it to, using the operator's `gh` login. The documentation
  shows those rules.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/runs/spec.md](specs/runs/spec.md)
- Tasks: [tasks.md](tasks.md)
- Constitution amendment: pull request #17
