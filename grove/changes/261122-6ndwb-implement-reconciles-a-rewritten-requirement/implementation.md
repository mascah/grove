# Implementation: Implement reconciles a requirement rewritten on the target

<!-- Progress on change 261122-6ndwb, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The implement instructions name a rewritten requirement

Status: done

`stale` in `internal/cli/check.go` now calls a per-change helper,
`rewritten`, which returns each requirement a change modifies or removes
that the maintained spec in the working tree rewrote or removed since its
`base`; `grove check` prints the same line as before from it.
`instructionsChange` in `internal/cli/instructions.go` runs `rewritten` on
the working tree's copy of the change, for a change in no cycle, a draft
and a cycle member, with or without `--runner`, and when it finds any
appends one line to what Grove requires, such as:

    - Since base 1a2b3c4, `record-editing`'s "Keep a category" was rewritten; before building any task, revise the change with Grove's `instructions propose <id>` so that each holds what the maintained spec now says together with what this change alters, and write a `Check:` line under a task that builds it, naming it and what was kept from each text.

Checks: `TestCoordinating` in `internal/cli/walkthrough_test.go` asserts
the line for a rewritten requirement, a removed one and a cycle member
built with `--runner` after an earlier member was archived on the branch,
and its absence before anything rewrote the requirement and when the
change depends on the rewriting one. With `instructions.go` reverted, three
of those assertions fail. `just check` passes (4574 Go, 2395 test lines;
8227 instruction words).

Check: the guidance is printed in the line itself, only when a requirement is named, and `internal/instruct/requires/implement.md` is unchanged, because the task's constraint keeps the instructions for a change with no such requirement as they are now; this mirrors the "Changed on main since" line, which is also built in `instructions.go`.
Check: the requirement is found by running check's own logic on the working tree's copy of the change, so once the agent revises and re-bases the change on the branch, the line goes away, as check's report does.

### T-02 — Docs and walkthrough

Status: done

`docs/coordinating.md`, "When the plan is stale", after check's report and
how to revise the change, now shows the line `/grove-implement` prints for
the same requirement, says the agent revises the change on its branch
before building any task and writes a `Check:` line, and that the merge
approves the revision with the code; it links the implement spec's new
requirement. `dev/walkthrough.md`, step 4 of "Coordinating changes", where
`grove check` on main reports show-saved's requirement as rewritten, has a
fresh session's `/grove-implement <show-saved>` print instructions naming
it. `TestCoordinating` runs that step.

Check: the walkthrough step was not run by hand in Claude Code or Codex; `TestCoordinating` runs the same step on a disposable fixture repository and asserts the line.

## Scenarios

- implement "Rewritten by a change that landed first": internal/cli/walkthrough_test.go TestCoordinating, the "implement of a stale change" assertion, passed; the revise-and-archive half is the agent following the line, unverified by a test since it is the agent's behavior
- implement "Removed by a change that landed first": internal/cli/walkthrough_test.go TestCoordinating, the "requirement was removed" assertion, passed
- implement "Nothing rewritten": internal/cli/walkthrough_test.go TestCoordinating, the "before anything rewrote it" assertion, passed
- implement "Rewritten by a change it depends on": internal/cli/walkthrough_test.go TestCoordinating, the "rewritten by a change it depends on" assertion, passed
- implement "Two members of one cycle": internal/cli/walkthrough_test.go TestCoordinating, the "implement --runner of a cycle member" assertion, passed

## Reconciliation
