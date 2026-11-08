# grove claim replaces grove run, with attempt records on a branch

## Problem

- `grove run` does two jobs. It selects and records the cycle to build,
  which only Grove can do, and it launches and supervises one Claude Code
  session: `internal/cli/run.go` at 8174594 starts
  `claude -p ... --output-format stream-json`, relays its events to the
  terminal, bounds it with `--timeout` and stops its process group. The
  runs spec at 8174594 puts "harnesses other than Claude Code" outside
  `run`, so the launcher serves one of the six harnesses `grove init`
  writes skills for, and it is the one exception `CONSTITUTION.md` rule 1
  makes to "Grove launches no agent".
- The claim is per checkout. `run` records each attempt as
  `.git/grove/<attempt>.json` and reads the latest one for the pause
  (runs spec at 8174594, "Pause after an attempt"; `run.go`, `latest`).
  Two checkouts of one project, on two machines or on one, both read "not
  paused" and both build the oldest marked cycle. Rule 11 at 8174594 names
  the gap: "coordinating runs across hosts" and "runs on disposable
  runners" are out of scope.
- A second opinion on an earlier draft of this change (a Codex review in
  this conversation, 2026-11-08) found two faults in that draft. A timeout
  cannot prove the earlier runner is dead, since a killed `grove run`
  leaves its session running, as the runs spec itself says under "One run
  per checkout". And a record kept on the agent's own work branch is
  rewritten by the rebases and stacked pull requests the project's `rules`
  ask for, while an end-of-run commit would publish whatever an interrupted
  agent left staged.
- An experiment in this conversation on 2026-11-08, two clones of one bare
  repository, showed that `git push` of a new branch is accepted for the
  first clone and rejected for the second, that `--force-with-lease`
  against a stale tip is rejected with "stale info", and that it succeeds
  against the true tip. The experiment was not committed.

## Intended outcome

`grove run` is gone. `grove claim [cycle-id] [--retry]` does what only
Grove can: it fetches the target, refuses while the project is paused,
selects the cycle it is given or the oldest marked `unattended` that can
be built, records the attempt as one commit on the branch `grove/runs`,
pushed to `origin` with a compare-and-swap so that two hosts cannot record
the same claim, creates the attempt's branch and prints the cycle's ID.
Whoever called it starts the agent, in any harness, bounded as they see
fit:

```sh
K=$(grove claim) && timeout 2h claude -p "/grove-implement $K --unattended"
```

With nothing to claim or the project paused, `claim` prints nothing on
standard output and exits non-zero, so nothing starts. The record on
`grove/runs` is the one authority on whether the project is paused, from
any host. `grove serve` reads it, and shows who claimed what, when, and
how recently the attempt's branch moved. Nothing is written under
`.git/grove/`.

## Scope

### In scope

- `grove claim`, with its optional cycle argument and `--retry`: the
  fetch, the pause, the selection, the record, the compare-and-swap push,
  the attempt's branch and the printed ID.
- The branch `grove/runs` and the record's shape.
- `grove serve` reading the records and naming `grove claim` in its
  verdicts.
- Removing `grove run`, `--timeout`, the session relay, the process-group
  stop, the lock and `.git/grove/`.
- Instruction text, documentation and the walkthrough naming `grove
  claim`.

### Out of scope

- Starting, bounding, stopping or logging an agent: the caller's, with
  `timeout`, a workflow's time limit or the harness's own output.
- An end record. Grove no longer knows when a session ends; `serve` shows
  the attempt branch's latest commit time instead.
- Expiring a claim by time. A stale claim is released by archiving its
  cycle on the target or by `--retry`.
- Building two cycles at once. The pause stays project-wide; parallel
  cycles would be a change of their own.
- `serve` fetching from `origin`; the owner fetches, as today.
- Grove writing to GitHub, GitLab or any tracker, and new guides for a
  scheduled runner or a tracker projection, until a recipe has been
  trialed (rule 8).
- Amending `CONSTITUTION.md`: its own pull request with no code, which
  lands before this change is built. Rules 1, 2, 3, 4, 7, 10 and 11 name
  `run`, launching, `.git/grove/` and cross-host coordination.

## Decisions

The owner accepted every item of the round in this conversation on
2026-11-08 with "yes". Each item below is one of them, as worded there.

1. **Remove the launcher.** "`grove run` goes, `grove claim [cycle-id]
   [--retry]` replaces it, Grove starts no process. Breaking, so
   `feat(cli)!`." Reason: launching is the harness's job and the caller's;
   it tied Grove to Claude Code's command line and output format, and was
   the one exception to rule 1.
2. **Records live on branch `grove/runs`**, "one commit per attempt, the
   record JSON as the commit message on an empty tree, so no new Git
   plumbing. Pushed with `--force-with-lease` against the tip just read;
   `update-ref` with the old value does the same locally without origin."
   Reason: a Git ref update is atomic, so the push is the compare-and-swap,
   with no service; and the record never shares a branch with the agent's
   work.
3. **Keep the global pause**, "read from `grove/runs`: while the latest
   attempt's cycle has a member not archived on the fetched target, `claim`
   refuses, naming the host and start time. `--retry` overrides. No
   automatic expiry, since Git cannot prove a runner is dead." Reason: the
   pause is also the circuit breaker that keeps a broken harness or an
   unmerged pull request from piling up attempts.
4. **No end record and no `.git/grove/`.** "The lock, logs and JSON files
   go. `serve`'s runs table reads `grove/runs`, preferring
   `origin/grove/runs` when it exists, and shows the branch's last commit
   time." Reason: nothing of Grove's runs while the session does, so
   nothing can write an end; the branch's commits are the evidence of
   activity.
5. **`claim` still prepares the checkout** "as `run` does: refuses an
   unclean tree, fast-forwards the local target, creates `grove/<attempt>`
   from the fetched target, prints only the cycle ID on stdout. It no
   longer requires Claude under `agents`." Reason: the implement
   instructions read the local target, and a bare ID on standard output
   composes in a shell with any harness.
6. **Two changes**, this one and [the project as
   JSON](../261108-7grt3-the-project-as-json/proposal.md), which depends on
   it. Reason: each can be agreed and built on its own.
7. **Documentation.** "Rewrite running.md and the timer guide around
   `claim` plus the invoker's command with `timeout`. No Actions or tracker
   guide yet." Reason: rule 8, guides are trialed recipes.
8. **Routine defaults.** "Record fields are today's minus `ended`,
   `outcome`, `output`, plus `host` from the hostname; `TestRun` becomes
   `TestClaim` and adds two clones racing for one claim, per rule 17; the
   cycles, config and serve spec Purpose sentences lose their mentions of
   `grove run`; test layout and private helpers delegated to the
   implementer." The Purpose sentences change when this change is
   archived.

## Open questions

None must be answered before review. One can wait, and decides nothing
about how a task is built: whether the record should also carry the
committer's name from Git's `user.name`, beside the host, for projects
where several people's checkouts share one machine.

## Related documents

- Requirements: [runs](specs/runs/spec.md), [cycles](specs/cycles/spec.md),
  [serve](specs/serve/spec.md), [config](specs/config/spec.md),
  [implement](specs/implement/spec.md), [propose](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
