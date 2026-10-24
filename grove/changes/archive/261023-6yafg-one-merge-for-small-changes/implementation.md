# Implementation: One merge for small changes

<!-- Progress on change 261023-6yafg, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement builds a draft with the owner

Status: done

- `internal/cli/instructions.go`: `instructionsChange` reads a change the
  target does not hold from the working tree, with the owner only; with
  `--unattended` it is refused as an unknown change is ("main holds no
  approved change <id>"). A draft that waits is refused as any change is.
  A change the target holds is read from the target, as before.
- `internal/instruct/instruct.go`: `Implement` calls a change the target
  does not hold a drafted change, "as this working tree holds it, with the
  owner", in its task line.
- `internal/instruct/requires/implement.md`: a rule for a change the line
  above calls a draft: build it as the working tree holds it; merging this
  branch approves it with its code, so its proposal, change specs and tasks
  are edited only with the owner's agreement, as when revising a change on
  its branch. The instructions for an approved change gain this one
  conditional bullet; what they require of an approved change is unchanged,
  and their task line and every other rule are as before.
- `internal/instruct/skills.yaml`: implement's description says "an
  approved or drafted change"; `grove init` rewrote
  `.claude/skills/grove-implement/SKILL.md` and
  `.agents/skills/grove-implement/SKILL.md`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks a
  change proposed and committed on a branch, not merged: implement
  instructions call it a draft, say merging approves it with its code and
  that its documents are edited only with the owner's agreement; with
  `--unattended` it is refused with "main holds no approved change <id>";
  a draft whose `depends_on` names the unfinished first draft is refused as
  waiting. The approved change's checks are as before. `TestInit` checks the
  description in both skill files.

### T-02 — Propose offers one merge or two

Status: done

- `internal/instruct/requires/propose.md`: the finish says, for a change
  the project index does not list as approved, to offer both next steps and
  choose neither for the owner: merge it to approve it, or run
  `/grove-implement <change-id>` on this branch so one merge approves it
  with its code; for a change it lists as approved, the next step is
  merging the revision, since `/grove-implement` builds it as the target
  holds it. The requires text is the same for every propose session, so
  the revise instructions for an approved change hold the offer only under
  its condition; the agent tells the two apart from the project index,
  which lists the change as approved.
- `internal/cli/walkthrough_test.go`: `TestPropose` checks the two next
  steps and the approved-change case in `grove instructions propose
  "<description>"`.

### T-03 — Review names a change approved with its code

Status: done

- `internal/cli/review.go`: a change the branch holds archived as finished,
  which the target holds neither open nor archived, is reported as
  "Merging <branch> into <target> approves change <id> with its code."
- `internal/cli/walkthrough_test.go`: `TestOneMerge` proposes, builds and
  archives a change on `rename` and checks review's first line, then merges
  once: `grove changes` prints nothing and `grove specs` lists
  `record-editing`. `TestApprove` still sees a draft not yet built as
  future work, and `TestArchive` an approved change archived on its
  branch as "approves no new change".

### T-04 — Document one merge for small changes

Status: done

- `docs/concepts.md`: "Approval is a merge" says a small change proposed
  and built on one branch with the owner is approved with its code by one
  merge, Grove setting no size, and that `--unattended` and `grove run`
  build only changes the target holds; the loop's step 4 offers skipping
  the first merge.
- `docs/getting-started.md`: after step 3, the two next steps the agent
  offers, building on `rename-projects` with `/grove-implement`, and
  review's line "approves change 261015-k3m9q with its code".
- `docs/reference.md`: the `grove instructions implement <change-id>` row
  takes a draft, read from the working tree; the `--unattended` row only a
  change the target holds; the `grove review` row names a change approved
  with its code; the `grove instructions propose` row offers both next
  steps.
- `dev/walkthrough.md`: "Implement refuses a draft" is now "One merge for a
  small change": proposed, built and archived on `rename`, review's line,
  one merge, `grove changes` and `grove specs` after it, and
  `--unattended` refusing the draft. `TestOneMerge` and
  `TestImplementInstructions` follow it; it was not walked by hand in
  Claude Code in this run.
- Recheck before finishing: T-01 to T-04's evidence holds against the code
  as it now is. `just check` passes (3721 Go lines, 1320 test lines);
  `grove check` prints "No problems in grove/.".

## Scenarios

- implement "Drafted and built on one branch": `walkthrough_test.go` `TestOneMerge` (instructions, archive on the branch, review) and `TestImplementInstructions` (the draft's instructions), passed; an agent building it with the owner is unverified, no owner in this run.
- implement "Unattended refuses a draft": `TestImplementInstructions`, passed.
- implement "Building shows the draft is wrong": unverified, needs an owner trying the build; the draft rule pointing to revising with the owner's agreement is checked by `TestImplementInstructions`.
- implement "A large draft": unverified by test; the code prints nothing about size, as `instructionsChange` and `Implement` show.
- implement "A draft that waits": `TestImplementInstructions`, passed.
- implement "Drafted, built and archived on one branch": `TestOneMerge`, passed.
- implement "A draft not yet built": `TestApprove`, passed.
- implement "An approved change archived on its branch": `TestArchive`, passed.
- implement "A change with two tasks": unverified, needs an owner in the session; the instructions are `TestImplementInstructions`'s, passed.
- implement "A change that waits": `TestCoordinating`, passed.
- propose "A new change on its branch": `TestPropose`, passed, for the instructions; an agent's finish is unverified.
- propose "Revising an approved change": unverified; the requires text gives the offer only for a change the index does not list as approved, checked by `TestPropose`, and an agent's finish was not tried.

## Reconciliation

- implement: applied to specs/implement/spec.md; Purpose and Boundaries name drafts and review's line.
- propose: applied to specs/propose/spec.md
