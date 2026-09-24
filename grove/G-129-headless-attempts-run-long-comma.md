---
id: "G-129"
type: work
title: "Headless attempts run long commands in the foreground, never as a background continuation"
status: proposed
created: "2026-09-24T04:49:40Z"
updated: "2026-09-24T04:50:02Z"
relates_to: ["G-114", "G-108", "G-045"]
---

## Outcome

A headless work attempt that must run a command longer than one tool call,
such as the G-108 eval pair, either runs it in the foreground with a timeout
or returns a checkpoint that names the wait, so that no attempt ends its turn
on a background job it will never be re-invoked for.

## Constraints

Observed 2026-09-24 on two Grove-owned attempts of
[G-114](G-114-capture-and-reuse-terms-question.md),
`G-114.20260924T042824Z` and `G-114.20260924T043857Z`, under
`.git/grove/attempts/`:

- Both started `python3 evals/run.py run …` as a background shell job (the
  second wrapped in an `until grep … sleep 15` loop with `run_in_background`),
  reported "I'll be notified when it finishes", and ended the turn. Under
  `claude -p` there is no next turn: the provider exited with code 0, the
  owner recorded the attempt as finished and successful, and the runner died
  with the session's process group. Each rerun cost about $2 to re-read the
  assignment and about $0.30 of eval before the kill, and left a half-empty
  run directory under `~/.cache/grove-evals/runs/`.
- The G-108 baseline (`~/.cache/grove-evals/runs/2026-09-23-G-108`) was
  launched detached from an interactive session, which is why it has a
  `runner.pid` and why the [work guide](../docs/work-execution.md) never had
  to say how a headless attempt runs a nine-minute command.
- The guide's "Never leave an untracked background agent running as an
  implied continuation" (step 8) is about agents; the model did not read a
  background shell job as one. The headless missing-decision path says how
  to return a wait on a human, not a wait on a process.
- The Claude Code Bash tool's foreground timeout tops out at ten minutes,
  which is about the eval pair's runtime, so the background was the path of
  least resistance.
- The attempt owner ([attempt.go](../internal/attempt/attempt.go)) records
  exit 0 and `success` for such a run; nothing in `grove attempts` or the
  board distinguishes an attempt that finished from one that abandoned a
  running job.

In scope: one rule in the work guide's headless path, that a long command
runs in the foreground with a timeout, or its wait is returned as a
checkpoint, never backgrounded as an implied continuation; whether the
result reconciliation should flag an attempt whose last message promises a
continuation is a question for shaping, not assumed.

Out of scope: a Grove-owned way to run and await the eval pair, and changing
`evals/run.py`.

## Acceptance

1. The work guide's headless path states the rule once, `grove guide work`
   prints it, and the Evidence records the new guides digest.
2. One headless attempt of a record whose mandate includes a command longer
   than a tool call, run after the change, either finishes the command or
   returns a checkpoint naming it; its attempt files are cited.
3. `grove check` passes.

## Next

Proposed 2026-09-24 from the diagnosis of G-114's two abandoned attempts.
Shape before assigning: the guide sentence is small, but whether the attempt
owner should surface an abandoned job is an open choice.
