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

## Scenarios

## Reconciliation
