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

## Scenarios

## Reconciliation
