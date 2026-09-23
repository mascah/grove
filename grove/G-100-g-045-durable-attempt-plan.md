---
id: "G-100"
type: plan
title: "G-045 durable attempt plan"
status: current
created: "2026-09-23T02:23:05Z"
updated: "2026-09-23T02:23:09Z"
work: ["G-045"]
---

## Design

Prepared 2026-09-22 at `5ae87d5` for [G-045](G-045-durable-attempt.md)
(record revision `c976f07d…`), from the runner research in
[G-026](G-026-shaping-and-runner-evidence-review.md), the evaluation's
runtime boundary in [G-048](G-048-direction-evaluation-review.md), the
predecessor adapter findings in [G-024](G-024-predecessor-work-review.md),
and the installed Claude Code 2.1.280 (`claude --help`, `claude agents
--json`, the official headless, sessions and agent-view documentation, read
2026-09-22). The plan owns implementation steps; G-045 owns outcome,
bounds and acceptance.

### Mechanism comparison (acceptance 1)

| Contract item | Native background session (`claude --bg`, `agents`, `attach`, `logs`, `stop`, `rm`) | Grove-owned `claude -p` process | tmux host |
| --- | --- | --- | --- |
| Invocation | Interactive session under Claude's supervisor daemon; "rejects `--bg` combined with `-p`" | Print mode with `--output-format stream-json --verbose` | Hosts either; adds a socket namespace and a dependency (tmux 3.6b present) |
| Identity | Supervisor id plus session id; state under `~/.claude/jobs/<id>/`, documented as "not a stable interface"; `claude agents --json` is the supported read | Grove attempt id and a `--session-id` Grove generates before launch, so the session is known even without a result event | Session name; no attempt bookkeeping |
| Survives the terminal | Yes, through the supervisor; machine restart "stops running background sessions", restarted on attach within 48 h | Yes: an owner process in its own session (`setsid`) with files, no terminal; no reboot recovery (not selected) | Yes while the tmux server lives |
| Events and result | `logs` prints recent terminal text; no structured result | Newline-delimited JSON events; `system/init` records model, permission mode, version, capabilities; `result` carries `subtype`, `is_error`, `session_id`, `total_cost_usd`, `num_turns`, `permission_denials` | Terminal text only |
| Budget | No `--max-budget-usd` ("only works with `--print`") | `--max-budget-usd` | none |
| Permission profile | Starts "the way a new `claude` session in that directory would"; prompts wait under "Needs input" for an attach | `--permission-mode MODE --permission-prompts none`: nobody answers, denials are recorded | inherits |
| Workspace | Claude moves the session into its own worktree under `.claude/worktrees/` before editing (`worktree.bgIsolation`) | Grove creates or reuses the worktree with explicit destination and provenance | caller's |
| Duplicate start, Stop, owner loss | Not per work item; `stop` keeps the conversation | Grove: per-work refusal while an attempt runs, `stop` through the owner, owner loss detected and reconciled | none |
| Daemon | Requires Claude's supervisor (`~/.claude/daemon`, currently `daemon-auth-status: auth_required` on this machine) | None: one on-demand owner per attempt | tmux server |

Selection: the Grove-owned `claude -p` process under an on-demand owner.
It is the only option that gives structured events, a spend bound, a
recorded permission profile and Grove-owned workspace and identity, and it
keeps ordinary CLI access daemon-free. Native sessions are the right
tool for a person's interactive parallel work, not for a bounded assignment
whose result Grove must reconcile; tmux would host the same owner and add
nothing Grove needs. Capabilities are pinned to Claude Code 2.1.280: the
`result` fields above, `--permission-prompts` (v2.1.259+), `--session-id`,
`--max-budget-usd`, and the `capabilities` array in `system/init`
(v2.1.205+), which the result records verbatim. Claude Code rejects
`--bg` with `-p` and exits 143 on SIGTERM, leaving the turn unfinished;
SIGINT ends the turn.

### Attempt contract

- **Identity.** `WORK.YYYYMMDDTHHMMSSZ`, for example `G-045.20260922T210000Z`,
  serialized per work by the launch lock so two starts in one second cannot
  share it. Its directory is `<git common dir>/grove/attempts/<attempt>/`,
  beside the ID counter: shared by every worktree of the repository, never
  committed, readable without any allocator or owner state.
- **Inputs and provenance** (`attempt.json`, written before the owner
  starts): work ID, record path and revision, the launching checkout and
  its HEAD (the base), branch, worktree, the exact command argv, the
  requested model, budget, permission mode, the generated session id, the
  `claude` executable (`GROVE_CLAUDE` overrides `claude` for fakes) and its
  `--version`, Grove's version line, and the start time. The owner adds its
  pid, session id, and the child's pid and process group.
- **Working directory.** `grove run ID` creates `worktree-ID` at
  `<root>/.claude/worktrees/worktree-ID` from the launching checkout's HEAD
  (`--branch` and `--worktree` override both), adds the directory to
  `info/exclude` when Git does not already ignore it, and reuses an
  existing registered worktree of that branch so a second attempt continues
  from preserved partial work. A directory that exists but is not that
  branch's registered worktree, a branch checked out elsewhere, or an
  uncommitted change to the record in the launching checkout (the attempt
  would not see it) refuses the launch before anything is written.
- **Process owner.** The launcher starts the Grove binary itself
  (`os.Executable()`) with `GROVE_ATTEMPT_OWNER=<dir>` in its own session
  (`Setsid`), stdin from `/dev/null`, stdout and stderr to `owner.log`,
  holding `owner.lock` (`repo.Lock`, an flock) from before the spawn: the
  launcher opens and locks the file, passes it as an inherited descriptor,
  and closes its copy, so the lock is held continuously from launch until
  the owner exits and never inherited by `claude` (close-on-exec). The
  owner runs `claude -p "/grove-work ID --interaction headless"
  --output-format stream-json --verbose --session-id UUID --max-budget-usd
  N --permission-mode MODE --permission-prompts none [--model M]` in the
  worktree, in its own process group, with stdin `/dev/null`, stdout to
  `events.jsonl` and stderr to `stderr.log` opened by the owner: raw capture
  is the kernel writing a file, so a closed or slow viewer cannot block the
  stream and no buffer grows. The child's environment drops
  `GROVE_ATTEMPT_OWNER`, `CLAUDECODE`, `CLAUDE_PID`, `CLAUDE_CODE_*` and the
  Git location variables, keeping `CLAUDE_CONFIG_DIR`. `go run` deletes its
  binary after the launcher exits; the owner keeps running on the mapped
  image, which is fine for one attempt and is why the owner never re-executes.
- **Budgets and bounds.** `--budget USD` and `--permission-mode` are
  required: Grove sets no default spend and no default permission profile,
  since either would make the choice in practice. Grove starts exactly one
  process per attempt and never retries; subagents Claude starts share the
  attempt's dollar budget and are not bounded separately (recorded as the
  bound). The result is whatever the one process left.
- **Events.** `events.jsonl` is raw. Readers parse it bounded: lines over
  1 MiB are counted as oversized and skipped, malformed JSON and unknown
  `type` values are counted, only `system/init` and `result` fields are
  extracted, and a final line without a newline is a partial line. No
  provider text is printed by the CLI; `attempt` prints counts, fields and
  paths, and the TUI (G-046) renders activity behind Grove's escaping.
- **Result** (`result.json`, written once by the owner after the child
  exits, temp file then rename): exit code or signal, whether Stop was
  requested, the extracted `system/init` and `result` fields, the event
  counts, the worktree's HEAD and whether it is dirty, and the work record
  as the worktree holds it (status, candidate, revision). A missing
  `result` event is `null`, a nonzero exit or `is_error` is recorded as it
  is: process exit and streamed claims are facts, never acceptance, and the
  record's own status in the worktree is the only handoff.
- **Liveness and owner loss.** running = `owner.lock` held. Owner lost with
  the child's process group still alive (`kill(-pgid, 0)`) = `orphaned`;
  owner lost with nothing alive and no result = `interrupted`. No pid is
  trusted alone. Machine reboot reads as `interrupted`, not selected.
- **Stop.** `grove stop ATTEMPT` sends SIGTERM to a running owner; the
  owner sends SIGINT to the child (ends the turn), waits 15 s, then SIGKILL
  to the child's process group, and writes the result with `stopped: true`.
  For an orphaned attempt, `stop` kills the group itself and writes the
  result with `reconciled_by: stop`. Stop never touches the worktree or the
  work record: partial code and evidence stay.
- **Reconnect and wait.** `grove attempts [ID]` lists attempts with status,
  start, branch, exit and cost; `grove attempt ATTEMPT [--json]` reads the
  files and prints the facts, plus `inputs changed` when the record on the
  target branch no longer hashes to the launch revision. Reading never
  starts a process and never resumes the conversation: resuming with
  `claude -p --resume SESSION` is a different, explicit action, out of scope.
- **Missing decision and unchanged wait.** The headless guide persists a
  question with `blocks`. `run` refuses while any open question blocks the
  work, naming it, and while an attempt for the work is running or
  orphaned; it also refuses a record whose status is `review`, `done` or
  `abandoned`. Rerunning with nothing changed returns the same refusal:
  an unchanged wait cannot spend another process.

### Placement

- `internal/attempt`: `Launch`, `List`, `Show`, `Stop`, `Own`, and the
  bounded event reader; `cmd/grove/main.go` runs `Own` when
  `GROVE_ATTEMPT_OWNER` is set, and the package's `TestMain` does the same
  so the test binary is its own owner.
- `internal/cli`: `run`, `attempts`, `attempt`, `stop` and their usage text.
- A decision record for the mechanism selection with the alternatives above;
  README command list; `docs/record-model.md` (attempts are files under the
  Git common directory, not records); `docs/work-execution.md` invocation
  table, whose headless row `grove run` now exercises.

## Steps

1. Set G-045 active; commit this plan.
2. `internal/attempt`: contract types, directory layout, `Launch` with the
   worktree rules and refusals, `Own`, bounded reader, `Show`, `List`,
   `Stop`. Fake-process tests first (acceptance 5), each behind `-short`
   because they spawn processes: competing start; owner outside the
   launcher's session; reconnect without a second start; Stop after
   reconnect; owner killed with the child alive (orphaned, then Stop) and
   with the child gone (interrupted); partial, unknown, oversized and
   malformed events; result with `is_error` and nonzero exit; no result
   event; record changed on the target after launch; a blocking question
   refusing the next run; partial uncommitted work surviving Stop and
   reused by the next run.
3. CLI commands and usage; `go run ./cmd/grove check`, vet, gofmt, tests.
4. Documentation and the decision record.
5. One bounded real-provider trial (acceptance 5) in a disposable
   repository with a clone-built binary, a small budget, a tiny work item
   with acceptance, and `--permission-mode` chosen at the trial; demonstrate
   launch, the launching shell exiting, `attempts`/`attempt` from another
   process, and the reconciled result; a second short real run for Stop
   after reconnect if the budget allows, otherwise Stop stays fake-only and
   is reported so.
6. Independent review of the combined diff; fix and re-review within the cap.
7. Evidence and handoff in G-045; `status=review` with the candidate.
