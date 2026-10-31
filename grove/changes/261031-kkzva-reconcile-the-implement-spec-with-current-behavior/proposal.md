# Reconcile the implement spec with current behavior

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

The implement spec has grown by one requirement per change, in the order
the changes arrived, until a person cannot take it in at a glance.

- At a5b5674, `grove/specs/implement/spec.md` holds 24 requirements, 59
  scenarios and 4,992 words, the largest of Grove's 12 specs (propose,
  next, has 16 requirements and 3,575 words). Its own instructions,
  `internal/instruct/requires/implement*.md`, are about 1,150 words.
- Related requirements sit apart because each arrived with a different
  change: "Honor a task's delegation" and "Ignore a task's delegation
  unattended"; "Suggest a rules line" and "Name a rules line unattended".
  "Build unattended" restates the rules of five others for the unattended
  case.
- Two read as history to anyone who never used the apply skill: "Report
  rules for apply" and "No apply skill". The second has no code of its
  own: `grove init` removes any skill Grove no longer lists
  (`grove/specs/init/spec.md`, "A skill Grove no longer has"), and
  `grove instructions apply` fails like any unknown skill.
- The spec and the instructions have drifted apart the other way too:
  five rules the implement instructions print are in no requirement (see
  "Rules the instructions already print").

## Intended outcome

The implement spec says what implement does today in 8 requirements and
under 3,000 words. Every requirement it held is folded into one of them,
or retired with its reason, below. Nothing implement does changes, and
the reasons for each rule stay one link away in the spec's History and
the archived changes it links.

## Scope

### In scope

- One change spec for implement: 7 requirements ADDED, "Verify before
  finishing" MODIFIED, and the other 23 REMOVED, each naming where it
  went.
- implement's Purpose and Boundaries, as given below, and the clause in
  init's Boundaries that points to "No apply skill" (T-01).

### Out of scope

- Any code. The two lines that report a `rules.skills.apply` key
  (`internal/cli/check.go:413`) stay as behavior no spec agrees to; a
  later change may remove them.
- A guide to reconciling a maintained spec. It is written from this
  change once it has been built, as its own pull request, as the guide to
  splitting a spec was.
- Every other maintained spec, and a spec of its own for `grove review`.

## Decisions

- **Rewrite to current behavior, keeping history in Git.** The owner:
  "some higher level review of the maintained specs is necessary" and "We
  don't want to destroy history", betting the implement spec "is able to
  be rewritten to represent current functionality in much fewer
  lines/words and still describe the current state correctly". History
  keeps its links; archive adds this change's.
- **Spec text only.** The owner leaned "towards A (spec text only)": a
  reconcile change records what the code does and changes none of it.
- **A guide, not a schema, for now.** The owner: "lets do the guide for
  now", deciding on a built-in schema after using it a few times.

Accepted in the round before writing (the owner: "yes"):

1. The five rules the instructions already print become agreed, each as
   one clause (see "Rules the instructions already print").
2. "No apply skill" is retired.
3. "Report rules for apply" is retired; its code stays.
4. The two review requirements stay in implement, merged into "Review
   names what the branch approves and builds".
5. The scenario "A large draft" is retired: it states only that a size
   warning does not exist.
6. One task, T-01, for the text around the requirements in implement and
   init.
7. "Verify before finishing" keeps its name and is MODIFIED; every other
   requirement is REMOVED and its rewrite ADDED.
8. "Settle what the change does not" stays one requirement, though the
   densest, since merging each with-the-owner and unattended pair is most
   of the saving.
9. The change is built on its own branch, so one merge approves it with
   the rewritten spec.

## Coverage

Where each requirement at a5b5674 goes. "Fold" means its rule is a
clause or scenario of the named requirement.

| Requirement at a5b5674 (scenarios) | Goes to |
| --- | --- |
| Implement an approved change (2) | Build a change, a draft or a cycle |
| Ask the owner instead of blocking (2) | Settle what the change does not |
| Revise the change on its branch (3) | Settle what the change does not |
| Implement a cycle (1) | Build a change, a draft or a cycle |
| Tell the owner of a moved file (3) | A task's file changed on the target; its file no task names to Change only the files the tasks name |
| Tell the owner of a file moved since the plan (4) | A task's file changed on the target |
| Verify before finishing (4) | Verify before finishing (MODIFIED) |
| Archive names unverified scenarios (1) | Verify before finishing |
| Build unattended (6) | The unattended case of each requirement above |
| A member gets its cycle's plan (2) | Build a change, a draft or a cycle |
| Implement's rules in both modes (1) | Rules for implement |
| Report rules for apply (1) | Retired |
| No apply skill (2) | Retired |
| Implement a draft with the owner (5) | Build a change, a draft or a cycle; Settle what the change does not; "A large draft" retired |
| Review names a change approved with its code (3) | Review names what the branch approves and builds |
| Honor a task's delegation (2) | Change only the files the tasks name |
| Stop on a discovery that breaks the approach (1) | Settle what the change does not |
| Ignore a task's delegation unattended (1) | Change only the files the tasks name |
| Suggest a rules line (4) | Rules for implement |
| Name a rules line unattended (1) | Rules for implement |
| Review shows planned against actual (3) | Review names what the branch approves and builds |
| List the change's documents (4) | Name the change's documents |
| Implement leaves revising to propose (2) | Name the change's documents; Settle what the change does not |
| Recommend re-basing after a dependency lands (1) | A task's file changed on the target |

## Rules the instructions already print

Each is printed by the implement instructions at a5b5674 and stated by no
requirement; this change agrees to them as they are.

- "When a task is built, tell the owner how to try it and take their
  feedback before the next." (`internal/instruct/requires/implement-owner.md`)
- If `check` reports a requirement of the change as rewritten or removed
  since its base, stop: the change must be revised with propose and
  merged again first. (`implement-owner.md`)
- A done task whose result a revision changes is in progress again until
  rebuilt. (`implement-owner.md`)
- A task already done is not built again unless the check before
  finishing finds its evidence no longer holds. (`implement.md`)
- A cycle's instructions leave out a member that waits on a change neither
  being built nor archived as finished, saying why
  (`implement-cycle.md`; `TestCoordinating` in
  `internal/cli/walkthrough_test.go`).

## The text around the requirements

Archive replaces requirements, not the text around them, so T-01 writes
it.

implement's Purpose:

> `/grove-implement` builds an approved change, a draft, or each change
> of a cycle on a branch, asking the owner where the change does not
> settle a task; with `--unattended` it asks no one and marks the task
> blocked instead. "The owner" is the person in whose session the skill
> runs.

implement's Boundaries:

> This covers the implement instructions in both modes, `grove.yaml`'s
> rules for implement, and how `grove review` names what an
> implementation branch approves and builds. Starting a session with no
> one watching belongs to `grove run` ([runs](../runs/spec.md)); when a
> cycle member counts as changed since its plan, to
> [cycles](../cycles/spec.md); revising a change off its implementation
> branch, to [propose](../propose/spec.md). What makes a revision
> approved is the merge into the target (`CONSTITUTION.md` rule 4).

init's Boundaries loses its last clause, "removing the apply skill to
[implement](../implement/spec.md), "No apply skill"", and the sentence
before it ends at config.

## Before proposing

A throwaway rewrite, made on a prototype branch and not kept, answered
whether this could be done: 24 requirements and 4,992 words became 7
requirements and 2,222 words, `grove check` passed on a copy of the
repository holding it, and every scenario it kept named a test in
`internal/cli/walkthrough_test.go` or a line of the implement
instructions that shows it. Only one rule had no test asserting its
instruction text: unattended, the agent edits none of the change's files
but `implementation.md`. This change starts from that rewrite, with the
round's decisions.

## Open questions

None must be answered before review. These can wait:

- Whether a later change removes the `rules.skills.apply` report in
  `internal/cli/check.go`.
- Whether `grove review` gets a spec of its own.
- Whether the guide to reconciling a spec ships schema files to copy and
  an archive `rules` line saying when to reconcile.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
