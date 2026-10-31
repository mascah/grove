# A cycle's own member merges keep its plan fresh

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

A cycle's members can be built one at a time, each with
`/grove-implement <change-id>` (the implement spec's "naming the member's
place in build order"). When each one reaches the target in its own merge,
every such merge makes the plan stale for the members still open, as soon
as it edits a file one of their tasks names.

At a5b5674, `sincePlan` (`internal/cli/check.go`) skips a commit on the
target's first-parent history only when it edits the cycle file or
archives a change of a cycle this one depends on. A commit that archives
a member of the same cycle counts like any other. `grove check`,
`grove run` and the implement instructions' "Changed on main since cycle
<id>'s plan" line all share that count. The cycles spec's Purpose and its
requirement "Count a member changed with its plan" say the same, and so
does `CONSTITUTION.md` rule 10, which names only "the merges that archived
the changes of the cycles it depends on".

Observed in a private project of the owner's, on Grove 0.6.0. A cycle of
seven changes was delivered as a stack of six squash-merged pull requests,
one per member, each archiving its own change, plus one acceptance-fix
pull request that archived nothing. From the first member's merge onward,
`grove check` reported the last open member as `changed on main since this
plan`. Its tasks named shared files, such as the app's API module and its
Tauri `main.rs`. Every edit those six member merges made to the open
member's named files was to a file the archiving member's own tasks also
named, so the plan had seen each overlap. The acceptance-fix commit edited
three of the open member's named files outside any task, and that edit is
the drift the check exists for.

The documented path does not hit this. `docs/coordinating.md` ("Building
a cycle unattended") builds the cycle on one branch and merges once, and
`sincePlan` reads only the target, so members archived on that branch
never count.

## Intended outcome

A commit on the target's first-parent history that archives a change the
cycle itself lists is one the plan expected, as is one that archives a
change of a cycle it depends on. The whole commit is skipped, whatever
files it edits. This holds wherever the count is used: `grove check`'s
stale report, `grove run`'s refusal and selection, and the line in a
member's implement instructions naming files changed since the plan. A
commit that archives nothing, or archives a change of an unrelated cycle,
still makes the plan stale.

So a cycle can reach the target one pull request per member, squashed or
merged, and stay as fresh as the one-branch build leaves it. The rules
are in [specs/cycles/spec.md](specs/cycles/spec.md).

## Scope

### In scope

- The count of commits since a cycle's plan (`sincePlan`), shared by
  `grove check`, `grove run` and `grove instructions implement`.
- The cycles spec's Purpose definition of stale and its requirement
  "Count a member changed with its plan".
- `docs/planning.md` ("Re-plan a cycle"), `dev/walkthrough.md`
  ("Building a cycle unattended" step) and the test.

### Out of scope

- Counting only the edits to files the archived member's own tasks name.
  The owner chose to skip the whole commit (decision 2).
- The stale report's wording, and naming the commits or files in it.
  `grove instructions implement <change-id>` already names each task and
  file.
- Which commit a re-plan records as `planned`, and the plan and implement
  instructions. Neither changes.
- The amendment to `CONSTITUTION.md` rule 10. It is its own pull request
  with no code (the constitution's preamble), and this change needs it
  merged first (decision 3).
- The runs and implement specs. They refer to the cycles spec's rule and
  keep their text.

## Decisions

The owner asked for this direction ("ok lets /grove-propose this", after
the investigation in the same session, which recommended skipping the
whole commit) and accepted every answer below in one round ("yes",
2026-10-30).

1. **The `default` schema, not `bugfix`.** Reason: today's behavior is
   what the cycles spec and rule 10 agree to, so this changes agreed
   behavior rather than fixing wrong behavior.
2. **A commit that archives a change the cycle itself lists is skipped
   whole, whatever files it edits.** Reason: it matches how the one-branch
   build behaves, where siblings archived on the branch never count, and
   how an archive of a cycle it depends on is treated. In the observed
   project, counting only the archived member's own named files would have
   reported nothing more, since every overlapping edit was to a file the
   archived member named.
3. **`CONSTITUTION.md` rule 10 is amended first, in its own pull request
   with no code,** to read "other than by the merges that archived the
   changes of the cycle itself or of the cycles it depends on". This
   change needs that amendment merged before it is built, as change
   261023-2ej3w needed PR #61. Reason: rule 10 names the only merges
   `run`'s stale gate may skip, and this change adds one.
4. **No change to the stale report's wording, the plan instructions or the
   implement instructions.** Reason: the implement line already names each
   task and file, and with this change it names only what the plan did not
   expect.
5. **Only the cycles spec changes.** "Count a member changed with its
   plan" is MODIFIED to add the cycle's own archived changes, with two new
   scenarios: a member landing in its own squash merge, and a later commit
   archiving nothing. The Purpose's definition of stale changes to match,
   stated in the change spec's preamble. Reason: the runs and implement
   specs refer to this rule rather than restating it.
6. **Docs and walkthrough.** In `docs/planning.md`, "Re-plan a cycle"'s
   sentence on archives covers this cycle's own changes too, and a new
   sentence says the changes can then reach main one pull request at a
   time. In `dev/walkthrough.md`, the "Building a cycle unattended" step
   gains a paragraph: delete is built alone, its branch also edits
   `ledger/report.py`, which budgets' T-01 names, and it is squash-merged;
   `grove check` on main is clean; main is then reset and the step goes on
   as written. Reason: the documentation is the specification, and a
   behavior exists only where a documentation step needs it (rule 8).
7. **Routine defaults:** one task. It touches `internal/cli/check.go`
   (with its comment), `internal/cli/walkthrough_test.go`,
   `docs/planning.md` and `dev/walkthrough.md`. The test extends
   `TestFollows`, giving cycle L's member d a task naming `ledger/cli.py`,
   so that landing c alone must leave L clean and selected. No
   `depends_on`. A `Delegated:` line covers test layout and private helper
   names.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/cycles/spec.md](specs/cycles/spec.md)
- Tasks: [tasks.md](tasks.md)
