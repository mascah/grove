# Plan's verdict, runs without a remote, streamed output

## Problem

Three observations from the owner's use of Grove on itself and on the
ledger fixture, 2026-10-22 and 2026-10-23.

- **Plan stops at the questions.** `internal/instruct/requires/plan.md` at
  5383ad8 has the agent ask, per member, whether the change needs a person
  and tell the owner of members that wait on a change outside the cycle.
  Nothing requires it to say whether `grove run` would build the cycle once
  marked, or what would remove an obstacle. The owner: "When planning a
  cycle, agent should help recommend whether it can be `unattended` or not
  and how to make it so if not."
- **`grove run` needs a remote.** `internal/cli/run.go:72` at 5383ad8
  fetches the target from `origin` before anything else, and the runs
  spec's Purpose defines "the fetched target" from `origin`. In a checkout
  with no remote the owner saw: "grove: git fetch: fatal: 'origin' does not
  appear to be a git repository fatal: Could not read from remote
  repository. ...: exit status 128", Git's words rather than one sentence
  (`CONSTITUTION.md` rule 14). The owner: "This seems like an arbitrary
  restriction but please confirm." It is specified, not accidental, and
  nothing in the constitution needs it.
- **Nothing shows until the session ends.** `session` in
  `internal/cli/run.go` at 5383ad8 runs `claude -p <skill>` in its default
  text output, which prints only the session's last message;
  `docs/running.md` says a stopped session's output is usually empty for
  that reason. The owner: "the terminal just hangs for however long the
  session takes and then prints the final message."

## Intended outcome

- A plan ends with a verdict: whether `grove run` would build the cycle
  once it is marked `unattended`, and for each obstacle what removes it.
- `grove run` works in a checkout with no remote named `origin`, using the
  local target branch; the owner merges the attempt's branch themselves.
- While the session runs, `grove run` prints one line per assistant
  message, tool call and result, and the attempt's log holds every event
  the session emitted, so a stopped session's log shows what it did.

## Scope

### In scope

- The plan instructions' finish and the cycles spec requirement for it.
- `grove run` in a checkout with no `origin`.
- The session's output format, what `grove run` prints while it runs, and
  what the attempt's log holds.
- The documentation and walkthrough steps for each.

### Out of scope

- Plan revising a change in the same session. The owner chose the verdict
  alone; a change that needs a person is revised with `/grove-propose` and
  merged, and the cycle re-planned.
- Any new dependency of the `grove` binary (`CONSTITUTION.md` rule 16).
- Rendering token-by-token output, or a TUI for runs.
- Recovering an interrupted attempt (`CONSTITUTION.md` rule 11).

## Decisions

- **Plan ends with a verdict and what removes each obstacle.** The owner,
  2026-10-23, offered the verdict alone, the verdict plus a hand-off to
  propose in the same session, or leaving plan as it is: "6 - 1." Reason:
  the owner's "help recommend whether it can be `unattended` or not and how
  to make it so if not"; plan reads members from the target, so a revision
  has to be merged before plan can read it.
- **`grove run` works without a remote.** The owner, 2026-10-23, offered a
  no-remote path, a one-sentence refusal, or leaving it: "7 - 1." Reason:
  the restriction served the pull request flow, and a local flow where the
  owner merges the branch is as sound.
- **Stream the session; keep the full JSON Lines log.** The owner,
  2026-10-22: "We might want to use the streaming JSON output to parse a
  session of `grove run`", and 2026-10-23, offered a JSONL log with
  rendered lines on stdout, a log of rendered lines, or plain verbose text:
  "8 - 1." Reason: the raw events stay parseable and complete, and the
  terminal shows progress.
- **One change, three tasks.** The owner, 2026-10-23, offered one change or
  three: "1". Reason: each part is one requirement with its docs and test.

### Assumptions to confirm in review

- The lines go to standard output whether or not it is a terminal, so a
  launchd job's log receives one line per event rather than one per firing.
- The rendering uses the standard library only; rule 16 is unchanged.
- The session is started with Claude Code's `--output-format stream-json`
  and `--verbose`; `--include-partial-messages` is not used.
- A line is printed for each assistant text message, each tool call (its
  name and the first line of its input, such as a command or a file path)
  and each tool result (its first line); system and other events go only
  to the log.
- Without `origin`, nothing is fetched, the fast-forward is a no-op, the
  attempt's record holds the local target's commit, and the pause ends
  when the owner merges the attempt's branch into the local target. With
  `origin`, nothing changes, including the refusal when the local target
  has commits `origin` does not.
- When `origin` exists but the fetch fails, the error Git returns is
  reported as today; handling a returned error is not a guard (rule 13).
- The verdict is instruction text in `requires/plan.md`, given with the
  members in build order at plan's finish; `grove check` keeps reporting a
  marked cycle that `grove run` would skip, and gains nothing.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/cycles/spec.md](specs/cycles/spec.md) (maintained
  spec: `grove/specs/cycles/spec.md`), [specs/runs/spec.md](specs/runs/spec.md)
  (maintained spec: `grove/specs/runs/spec.md`)
- Tasks: [tasks.md](tasks.md)
