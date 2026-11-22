# Implementation: Explore recommends and asks less

<!-- Progress on change 261122-m78w6, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Explore recommends, investigates first when asked, and stops asking

Status: done

Rewrote `internal/instruct/schemas/default/skills/explore.md` (291 words at
the base, 274 now) into Start, Recommend and Finish: investigate first and
reply with findings and a recommendation when asked to research,
investigate or recommend; a recommendation beside every set of options; one
question a turn, only one that would change the recommendation, then
finish; an earlier conflicting decision in one line. Dropped "the owner
sees the options ... before anything is recommended" and the file-and-commit
citation of every finding. In `requires/explore.md`, "your own inference
stays open, and you say so" is gone; a question is settled when the owner
answers it. `TestExplore` checks the new wording and that the old
recommendation ban is gone. `just check` passed; instruction words 8227 at
the base, 8206 after this task, 8208 at the end of the change.

Check: open question 2 (explore's limit) built as recommended: one question a turn, at most two in a first reply after investigating.
Check: open question 7 (a new `explore` capability) built as recommended; the maintained spec is written by archive.
Check: "Where documents or code disagree, ask which is intended" is dropped: it made a question whatever its answer would change, against "Ask only what would change the recommendation"; whether behavior is intended stays the owner's to say in `requires/explore.md`.

### T-02 — Every skill asks the owner plainly

Status: done

Added a paragraph to the shared opening (`meanings` in
`internal/instruct/instruct.go`), after "Grove itself commits nothing":
question first after at most a few lines of context, standing on its own;
options by what the owner gains and gives up, with a recommendation and its
reason; about 150 words unless the owner asks for more; a file or commit
cited only where the owner asks or a document needs it; no change ID, task
ID, rule number or decision number without saying in words what it is.
`TestExplore` checks that `grove instructions explore` and `propose` print
it before "## What Grove requires of". Ran `bin/grove instructions` for
explore, propose, plan and prototype "q?": each prints it on line 24, before
the skill's own section on line 32. `just check` passed.

Check: open question 1 (where the asking rules live) built as recommended: in the shared opening, so every skill, implement and archive included, prints it; it applies only when the agent asks the owner something.
Check: open question 6 (length) built as recommended: "about 150 words unless the owner asks for more", guidance Grove does not count.

### T-03 — Propose asks at most five decisions, in one message

Status: done

`internal/instruct/requires/propose.md`: the separate bullets that asked
what must land first, offered a split and asked a round of every decision
("Leave none out") are one bullet: before writing, one message, the round,
that states the scope, offers a split where the idea is too big, names an
open change touching the same capabilities or code and asks whether it must
land first, and asks at most five decisions, the most consequential,
numbered, each with a recommended answer and its reason worded so "yes"
accepts it, behavior/scope/compatibility/data first, then architecture or
interfaces costly to reverse; nothing to ask means saying so and writing
once the scope is confirmed. A new bullet writes every other decision,
routine choices on one open question of their own, as an open question with
its recommended answer, its reason and the task it decides.
`schemas/default/skills/propose.md`: steps 2 and 3 are one step, "Confirm
the scope and ask one round, in one message"; routine choices are not asked;
more than five that matter means offering the split; steps renumbered.
`TestPropose` checks the new wording and that neither the default nor the
`bugfix` output says "Leave none out"; `TestCoordinating` checks the new
land-first wording. `go test ./internal/cli` passes except `TestServe`,
which fails the same way at the base: a `grove serve` the owner runs
outside this work holds port 7070.

Check: open questions 3, 4 and 5 built as recommended: five, the most consequential, in the existing order; scope, split, land-first and the round in one message; other choices as open questions with recommended answer, reason and task, routine ones on one open question of their own.
Check: `internal/instruct/schemas/bugfix/skills/propose.md` is unchanged: its only mention of the round ("Ask, in the round before writing, whether the corrected behavior should become agreed") still holds, and the round's rules come from `requires/propose.md`, which both schemas print.
Check: "Unattended propose's instructions are unchanged" read as `requires/propose-unattended.md` unchanged (it is): unattended output prints the shared `requires/propose.md` too, so its text shows the new round, which the unattended part still replaces with open questions.

### T-04 — The docs and walkthrough describe the new explore and round

Status: done

`docs/getting-started.md` §2: explore recommends beside every set of
options, asks one question at a time and only one that would change its
recommendation, and investigates first when asked to recommend; §3: one
message with the scope and at most five decisions, every other choice an
open question with a recommended answer, the split offered in the same
message. `docs/concepts.md` (the loop's Explore and Propose),
`docs/reference.md` (the explore row, and "propose asking at most five
decisions in one round") and `docs/guides/spec-existing-code.md` (at most
five questions, the rest open questions) no longer say explore interviews
or that propose asks every decision. `dev/walkthrough.md` main line §2
(explore recommends keeping an older record's category, I say "yes", it
does not ask about the README/CHANGELOG wording, finishes), §3 (one message,
at most five, routine choices as open questions), Coordinating step 1 (no
"says so" when no open change touches the same code) and "The round in
order" (the test's name and place not asked, an open question naming T-01).
`grep -i "one question at a time|every decision|interview|line of defaults"`
over `docs`, `README.md` and `dev/walkthrough.md` finds nothing.

Check: the walkthrough's explore and propose steps were not run in Claude Code or Codex in this session; they describe what the new instructions say a run does, and `dev/walkthrough.md` is run by hand (CONTRIBUTING.md).

### T-05 — A skill trial runs explore

Status: done

`dev/trials/g.sh`: case (g) runs `/grove-explore Should ledger keep its
records in SQLite instead of records.json? Look into it and recommend.` on
branch `explore`; to a message holding a recommendation the owner answers
"Yes, I'll go with your recommendation."; `good` checks that `explore`
changes no file against main, that the scripted answer was given or the
last message recommends, and that the run took at most three of the
owner's replies. `dev/trials.md` describes (g) and says there are seven
cases; `dev/trials/run.sh` runs (g) by default. `RUNS=1 JUDGINGS=1 just
trials g` printed:

    | (g) | Claude Code | 1/1 | 0 | 0 | 0 | 0 | 0 | 236 |
    | (g) | Codex | 1/1 | 0 | 0 | 0 | 0 | 0 | 234 |

The first attempt stopped both runs: a run that writes no change leaves no
`grove/changes`, and `find` failing there ended the trial under `set -e`;
`run.sh` now tolerates it. That attempt's Claude Code run recommended
"Keep `records.json`" under "Short answer:", which the first answer
pattern ("recommend") missed, so the pattern also matches "short answer"
and "I'd keep/go/…".

Check: "a good run recommends" accepts either the scripted answer to a recommending message or a last message that recommends, because Codex recommended and finished in one turn, asking nothing, which the explore spec allows (at most two questions); the proposal's open question 8 named only the first.
Check: (g) was run once per harness, not the default three runs and three judgings, to keep the cost down; the full `just trials` table was not run.
Check: case (b) was also run once per harness on the new round: Codex 1/1; Claude Code 0/1 because it asked "Is it OK that deleting the highest ID lets the next `add` reuse that number?", which `b.sh`'s answer pattern does not match, so the scripted answer was not given. `b.sh` is left as it is: no task names it.
Also changed: `CLAUDE.md` — its layout table said the skill trials are six cases.

## Scenarios

- explore "Two ways to keep an older record editable": unverified; `TestExplore` checks the instructions say "say which you recommend and why" and that the owner's answer settles it, but the walkthrough's §2 run on the fixture was not done in this session.
- explore "A question that changes nothing": unverified for the README case; `TestExplore` checks "only one whose answer would change your recommendation"; in trial (g) both harnesses asked at most one question, judged 0 unneeded questions.
- explore "Look into it and recommend": `dev/trials/g.sh` case (g), passed in Claude Code and Codex (first reply findings and a recommendation after reading `ledger/records.py`, at most one question).
- explore "A description that asks nothing to be looked into": unverified; the instructions keep "one question or a short account of what you will look up", not run on the fixture.
- explore "Records move out of records.json": unverified; `TestExplore` checks "say so in one line", but no fixture holds such an archived decision and it was not run.
- instructions "Every skill is told": `internal/cli/walkthrough_test.go` `TestExplore` (explore and propose), passed; `bin/grove instructions plan` and `prototype "q?"` print it on line 24, before line 32's skill section.
- instructions "An archived change named in a question": unverified; how an agent words a question is for a person to judge, and it was not run on the fixture.
- propose "The owner accepts the round": unverified; `TestPropose` checks the round wording, and trial (b) in Codex recorded its answers as decisions, but no run checked the quoted "yes".
- propose "The round is ordered, and routine choices are not asked": partly by trial (b): both harnesses sent one message with the scope and at most five decisions, and Claude Code wrote routine choices as an open question; the unknown-ID case itself was not run.
- propose "More than five decisions": unverified; `TestPropose` checks "at most five decisions", not run with seven.
- propose "The owner declines one": unverified; the instructions' deferred-question rule is unchanged, not run.
- propose "Nothing to assume": unverified; the instructions say "If there is nothing to ask, it says so", not run.

## Reconciliation

- explore: created `grove/specs/explore/spec.md` from the change spec's ADDED requirements.
- instructions: applied to `grove/specs/instructions/spec.md` (ADDED "Ask the owner plainly").
- propose: applied to `grove/specs/propose/spec.md` (MODIFIED "Ask a round before writing"); its Purpose drops "asking one round of every decision it would otherwise assume", as the proposal's scope says.
