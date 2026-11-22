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
the base, 8199 after this task.

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

Also changed: none beyond the task's files.

Check: open questions 3, 4 and 5 built as recommended: five, the most consequential, in the existing order; scope, split, land-first and the round in one message; other choices as open questions with recommended answer, reason and task, routine ones on one open question of their own.
Check: `internal/instruct/schemas/bugfix/skills/propose.md` is unchanged: its only mention of the round ("Ask, in the round before writing, whether the corrected behavior should become agreed") still holds, and the round's rules come from `requires/propose.md`, which both schemas print.
Check: "Unattended propose's instructions are unchanged" read as `requires/propose-unattended.md` unchanged (it is): unattended output prints the shared `requires/propose.md` too, so its text shows the new round, which the unattended part still replaces with open questions.

## Scenarios

## Reconciliation
