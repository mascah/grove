## T-01 — Implement builds a draft with the owner

Result: `grove instructions implement <change-id>` builds a change only
this working tree holds as
[specs/implement/spec.md](specs/implement/spec.md) "Implement a draft
with the owner" says, and "Implement an approved change" no longer
refuses it.

Files: `internal/cli/instructions.go`, `internal/instruct/instruct.go`,
`internal/instruct/requires/implement.md`,
`internal/instruct/skills.yaml`, `.claude/skills/grove-implement/SKILL.md`,
`.agents/skills/grove-implement/SKILL.md`,
`internal/cli/walkthrough_test.go`

Interfaces: `grove instructions implement <change-id>` accepts a change
the working tree holds open and the target does not, read from the
working tree; `instruct.Implement` tells such a draft from an approved
change in the task line and requirements it prints. With
`--unattended`, as change 261023-g0a3h adds it, a draft is refused with
the message an unknown change gets today.

Constraints: a change the target holds is still read from the target,
as today.

Acceptance:
- In `walkthrough_test.go`, on the fixture: a change proposed and
  committed on a branch, not merged, gets implement instructions saying
  it is a draft, that merging the branch approves it with its code, and
  that its documents are edited only with the owner's agreement; the same
  change with `--unattended` is refused with "main holds no approved
  change <id>"; a draft whose `depends_on` names an unfinished change is
  refused as waiting.
- The implement instructions for an approved change are unchanged.
- The implement skill's description says it takes a drafted change as
  well as an approved one, and `grove init` writes it to both skill
  files.

## T-02 — Propose offers one merge or two

Result: propose's finish offers both next steps as
[specs/propose/spec.md](specs/propose/spec.md) "Offer one merge or two"
says.

Files: `internal/instruct/requires/propose.md`,
`internal/cli/walkthrough_test.go`

Interfaces: none.

Acceptance:
- `grove instructions propose "<description>"` prints both next steps:
  merge the proposal to approve it, or `/grove-implement <change-id>` on
  this branch so one merge approves it with its code.
- The revise instructions for a change the target holds do not offer to
  implement it on the branch.

## T-03 — Review names a change approved with its code

Result: `grove review` reports a change drafted and archived on the
branch as [specs/implement/spec.md](specs/implement/spec.md) "Review
names a change approved with its code" says.

Files: `internal/cli/review.go`, `internal/cli/walkthrough_test.go`

Interfaces: `grove review`'s line "Merging <branch> into <target>
approves change <id> with its code."

Acceptance:
- In `walkthrough_test.go`, on the fixture: a change proposed, built and
  archived on one branch is reviewed as approved with its code; a draft
  not yet archived is still reviewed as future work; an approved change
  archived on its implementation branch still reviews as "approves no
  new change".

## T-04 — Document one merge for small changes

Result: the documentation describes building a draft on its own branch,
review's new line and propose's two next steps (rule 8).

Files: `docs/concepts.md`, `docs/getting-started.md`, `docs/reference.md`,
`dev/walkthrough.md`

Interfaces: none.

Acceptance:
- `docs/concepts.md` says, under approval and in the loop, that a change
  drafted and built on one branch with the owner is approved with its
  code by one merge, and that `--unattended` and `grove run` build only
  changes the target holds.
- `docs/getting-started.md` offers, after step 3, implementing the change
  on the same branch instead of merging it first, and shows review's
  line for it.
- `docs/reference.md`'s rows for `grove instructions implement
  <change-id>` and `grove review` say what changes.
- `dev/walkthrough.md`'s "Implement refuses a draft" step becomes a step
  on the ledger fixture where a change is proposed, built and archived on
  one branch, `grove review` shows "approves change <id> with its code",
  one merge lands it, `grove changes` does not list it and `grove specs`
  lists its capability; and `--unattended` on a draft is refused.
