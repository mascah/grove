## T-01 — Explore recommends, investigates first when asked, and stops asking

Result: the explore instructions meet every requirement in the `explore`
change spec, and are shorter than at the change's base.

Files:
- `internal/instruct/schemas/default/skills/explore.md`
- `internal/instruct/requires/explore.md`
- `internal/cli/walkthrough_test.go`

Acceptance:
- `grove instructions explore` says to recommend beside every set of
  options, to ask only a question whose answer would change the
  recommendation and stop once none would, to investigate first and reply
  with findings and a recommendation when the description asks for
  research or a recommendation, and to give an earlier conflicting
  decision one line; and no longer says that nothing is recommended before
  the owner sees the options, nor that every finding is cited by file and
  commit in the conversation.
- `TestExplore` checks the new wording in place of "Interview before you
  conclude anything".
- `just check` passes, and the instruction word count is lower than at the
  change's base.

## T-02 — Every skill asks the owner plainly

Result: the opening every skill's instructions share meets "Ask the owner
plainly" in the `instructions` change spec.

Files:
- `internal/instruct/instruct.go`
- `internal/cli/walkthrough_test.go`

Acceptance:
- The output of `grove instructions explore`, `propose`, `plan` and
  `prototype "<question>"` holds the asking rules before the skill's own
  instructions.
- A walkthrough test checks that the explore and propose output hold them.
- `just check` passes.

## T-03 — Propose asks at most five decisions, in one message

Result: both schemas' propose instructions meet "Ask a round before
writing" in the `propose` change spec.

Files:
- `internal/instruct/requires/propose.md`
- `internal/instruct/schemas/default/skills/propose.md`
- `internal/instruct/schemas/bugfix/skills/propose.md`
- `internal/cli/walkthrough_test.go`

Acceptance:
- `grove instructions propose` and `grove instructions propose --schema
  bugfix` have the agent state the scope, offer a split, ask what must
  land first and ask at most five decisions in one message, and write
  every other decision, routine choices among them, as an open question
  with its recommended answer, its reason and the task it decides; neither
  says "Leave none out".
- Unattended propose's instructions are unchanged.
- `TestPropose` checks the new round wording.
- `just check` passes.

## T-04 — The docs and walkthrough describe the new explore and round

Result: no documentation step says explore interviews one question at a
time or that propose asks every decision it would otherwise assume.

Files:
- `docs/getting-started.md`
- `docs/concepts.md`
- `docs/reference.md`
- `docs/guides/spec-existing-code.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/getting-started.md` says explore recommends and asks only what
  would change its recommendation, and that propose asks at most five
  decisions with the scope and writes the rest as open questions with
  recommended answers.
- `docs/concepts.md`, `docs/reference.md` and
  `docs/guides/spec-existing-code.md` say nothing contrary to that.
- The walkthrough's explore and propose steps (§2, §3) describe what a run
  of this change's instructions does on the fixture.

## T-05 — A skill trial runs explore

Result: the skill trials have a case (g) that runs `/grove-explore` on a
description asking for a recommendation and counts what it costs the
owner.

Files:
- `dev/trials.md`
- `dev/trials/g.sh`
- `dev/trials/run.sh`

Acceptance:
- `dev/trials.md` describes case (g): its setup, its run, the owner's
  scripted answers and what a good run is: it writes nothing, its
  scripted answer to a message holding a recommendation was given, and it
  took at most three of the owner's replies.
- `just trials g` runs the case in both harnesses and prints its row,
  with the words the owner had to read.
- `dev/trials/run.sh` runs (g) among its default cases, and
  `dev/trials.md` says there are seven.
