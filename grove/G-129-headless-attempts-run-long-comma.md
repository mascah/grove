---
id: "G-129"
type: work
title: "Headless attempts run long commands in the foreground, never as a background continuation"
status: active
created: "2026-09-24T04:49:40Z"
updated: "2026-09-24T15:36:30Z"
relates_to: ["G-114", "G-108", "G-045"]
candidate: "347e3abdebc4b80adfc5a3f9cdeede6d14bbc88c"
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

## Evidence

Branch `worktree-G-129`, based on main `c38d914` (record revision
`sha256:d53d511c`). This was run as headless `/grove-work G-129`. The
candidate is the commit that holds this text; the status change that
follows it names that commit. The reviewed content is `4e8cd92`: the rule
`88312db`, review fixes `9ad05c5` and `4e8cd92`, and the review record
[G-132](G-132-g-129-independent-review-of-the.md). No plan was written,
because the In scope item is one guide rule. The owner's choice about the
attempt owner is left open, as the Constraints say, and the rule does not
depend on it.

Acceptance:

1. Met. [Step 5](../docs/work-execution.md#5-implement-through-evidence)
   of the work guide now ends with the rule, stated once. It opens
   "Headless, no command outlives the session." What it says:
   - Run a long command in the foreground with an explicit timeout, up to
     the harness's maximum, in bounded pieces where the command allows.
   - If the command cannot finish inside that maximum, or it times out,
     the work stays active. The session returns a committed checkpoint
     that names the command, the path of any partial output, why it must
     run, what its result decides, and who can run it. A headless rerun
     with nothing changed returns the same checkpoint.
   - Never end the turn on a background job as an implied continuation.

   The Inputs sentence about mode now points to this rule. Step 8's return
   list now includes a wait on a command. No other text restates the rule.
   `go build -o /tmp/grove-g129 ./cmd/grove`, then `/tmp/grove-g129 guide
   work | diff - docs/work-execution.md`, shows no difference.
   `/tmp/grove-g129 version` prints guides digest `3c9996e33e42`. Before
   this change, main `c38d914` printed `5a224350feae`.
2. Met, before integration, following the owner's feedback on `347e3ab`.
   Headless attempt `G-129.20260924T153802Z` ran two fixture attempts in a
   disposable clone, `/tmp/g129-fixture`. The clone's branch was
   `worktree-G-129` at `f0ffd50`. Its guide matches `0597319` byte for
   byte, since everything after the rule changed only this record. Both
   fixtures were launched as `/tmp/grove-g129 run ID --budget 3
   --permission-mode auto`, with that binary built from `f0ffd50` (guides
   `3c9996e33e42`). Each fixture work record in the clone had a single
   mandate: run the command as written and record its output, without
   shortening, splitting or simulating it. The clone's counter numbered the
   fixtures G-134 and G-135; those IDs exist only in the clone. The attempt
   files, plus each fixture branch's commits as patches, are copied to
   `~/.cache/grove/G-129-acceptance-2/ATTEMPT/`. They are not committed,
   because attempt files hold provider text.
   - `G-134.20260924T153900Z`, for `sleep 900 && echo done`, which is longer
     than the ten-minute foreground limit. It finished in 45 s, exit 0,
     $0.32. It set the fixture active and committed a checkpoint (`0b57f96`)
     naming the command. The checkpoint says why the command was not
     started (15 minutes against a 600000 ms maximum), that there was no
     partial output and no owned command, what the result decides, and that
     an interactive session or a person can run it. It also says a headless
     rerun returns the same checkpoint. None of its tool calls ran in the
     background. Its last message begins "G-134 is waiting on one command
     that I didn't run: `sleep 900 && echo done`". This is the guide's "do
     not start it" branch, not the owner's expected "run, time out,
     checkpoint". The outcome is the same checkpoint without the wasted
     ten minutes.
   - `G-135.20260924T154028Z`, for `sleep 420 && echo done`, a seven-minute
     command. That is the shape of the G-114 failure, a wait that fits the
     limit. It ran the command in the foreground with `timeout: 600000`,
     and the tool call printed `done` and `exit=0` at 15:47:53Z. It
     recorded the output in Evidence (`52c480f`) and handed the fixture into
     review with that candidate (`30b8b5e`). It finished at 15:48:18Z,
     exit 0, $0.35.
3. Met. At `4e8cd92`, `go run ./cmd/grove check` printed `OK: 126 records`
   before G-132 was written. `gofmt -l .` printed nothing, and `go vet
   ./...` was clean. At `9ad05c5`, `go test -count=1 -timeout 120s ./...`
   passed for every package. Several packages took more than five seconds,
   as they already did on main, and this change touches no Go code.
   `4e8cd92` changes only wording in the guide.

Review: G-132 ran two rounds with an independent reviewer subagent. Round 1
raised 2 consequential findings, 4 minor ones and 1 nit; round 2 raised 2
nits. All were fixed except one step 8 nit, left as it is on purpose. The
knowledge check found nothing: no new domain concept, and no conflict with
G-056 or G-101.

## Next

In review; `candidate` names the commit. All three acceptance items are
met. Commit `4e56eea` added the owner's recommendation, to check before
integrating in a throwaway clone. That check ran as Evidence acceptance 2
describes, and so it is removed from this section. The guide has not
changed since G-132 examined `4e8cd92`. Every later commit changes only
this record. To integrate, run the first command in a clean checkout of
`worktree-G-129` and the second in a clean checkout of `main`:

```sh
go run ./cmd/grove approve G-129 "VERDICT"
go run ./cmd/grove integrate G-129 --cleanup
```

`/tmp/g129-fixture` is disposable and can be deleted
(`rm -rf /tmp/g129-fixture`); the evidence copy stays in `~/.cache/grove/`.

Still open, for shaping and not for this record: should the attempt
owner's result reconciliation flag an attempt whose last message promises
a continuation?

Feedback on candidate 347e3ab, 2026-09-24: poke
