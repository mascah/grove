# One merge for small changes

## Problem

Every change takes two merges: one to approve the proposal, one to land
the code. At 57af2e3, `grove instructions implement <change-id>` reads
the change only from the target, and for a change the working tree holds
but the target does not it prints "main holds no approved change <id>"
(`internal/cli/instructions.go`, `instructionsChange`; the maintained
`grove/specs/implement/spec.md`, scenario "A change that is only a
draft"). Propose's finish names one next step, merging the proposal
(`internal/instruct/requires/propose.md`). For a one-task fix the owner
reviews the same small change twice, in two pull requests, a day apart.

Rule 4 of `CONSTITUTION.md`, amended in PR #40, now allows one merge: "A
change drafted on the branch that implements it, with the owner in the
session, is approved with its code by the merge that brings both into
the target; `run` builds only changes the target already holds." Nothing
in Grove does it yet.

`grove review` would also misreport such a branch. It names only changes
the branch holds open (`internal/cli/review.go`), so a change drafted and
archived on one branch shows as "approves no new change", though the
merge approves it.

## Intended outcome

The owner proposes a change on a branch and, in the same session or a
later one, runs `/grove-implement <change-id>` on that branch. Grove
prints the implement instructions for the change as the working tree
holds it, saying it is a draft, that merging the branch approves it with
its code, and that its documents are edited only with the owner's
agreement. The agent builds it and archives it there. `grove review`
says "Merging <branch> into main approves change <id> with its code."
One merge approves the change and lands the code, the updated specs and
the archived change together.

`/grove-implement --unattended`, and so `grove run`, still build only
changes the target holds. Propose's finish offers both paths. Which
changes are small enough is the owner's call; Grove sets no limit.

The rules are in [specs/implement/spec.md](specs/implement/spec.md) and
[specs/propose/spec.md](specs/propose/spec.md).

## Scope

### In scope

- `grove instructions implement <change-id>` for a change only the
  working tree holds, and the implement instructions for it.
- `--unattended`, as change 261023-g0a3h adds it, refusing a draft.
- Propose's finish offering both paths.
- `grove review`'s line for a change drafted and archived on the branch.
- The implement skill's description, `docs/concepts.md`,
  `docs/getting-started.md`, `docs/reference.md` and
  `dev/walkthrough.md`.

### Out of scope

- Drafts in cycles: a cycle lists only approved changes, and `run`
  builds only what the target holds (rule 4).
- A size rule or a check on how big a change is (decision 5).
- `--unattended` itself and apply's removal: change 261023-g0a3h.
- A review line for a draft dropped on its branch.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes").

- **`/grove-implement <change-id>` accepts a change that exists only in
  this working tree, a draft, with the owner in the session;
  `--unattended` refuses a draft.** Reason: rule 4 approves a draft with
  its code only when the owner was in the session, and `run` builds only
  changes the target already holds.
- **The implement instructions for a draft say it is a draft, that the
  owner's merge of this branch approves it with the code, and that the
  change is edited only with the owner's agreement, as for revisions
  while implementing.** Reason: the agent must not treat a draft as
  agreed (rule 5), and one rule for editing a change while building it
  is simpler than two.
- **Propose's finish offers both paths: merge the proposal now to
  approve it, or implement it on this branch so one merge approves and
  lands it.** Reason: the owner chooses per change; neither path is the
  default.
- **`grove review` says when a merge approves and archives a change at
  once: "approves change C with its code" instead of "as future work".**
  Reason: review is what merging would agree to (rule 6), and "approves
  no new change" would be false.
- **No size rule: the owner decides what is small.** Reason: any rule
  would be a gate Grove does not need (rule 10), and the owner reviews
  the merge anyway.

## Open questions

None must be answered before review. None blocks a task.

Change 261023-g0a3h, which this change waits on, may also modify
"Implement an approved change". If it does, `grove check` will report
this change's MODIFIED requirement as rewritten since base once that
change is archived, and this change is revised with propose to keep its
text and drop only the draft scenario. This can wait until then.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md),
  [specs/propose/spec.md](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
