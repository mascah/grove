# A stale cycle is reported whether or not it is marked

<!-- A bug fix describes future work. It becomes approved future work when
it is merged into the target branch, which says nothing about
implementation. -->

## Observed

A change in a cycle that is not done and not marked `unattended` gets no
report when a file one of its tasks names changes on the target after the
cycle's `planned` commit. `grove check` is clean, and
`grove instructions implement <change-id>` names no moved file. Observed
on the ledger fixture with main's build at 2d07433 (see Reproduction):

- `grove check` prints `No problems in grove/.` and exits 0.
- `grove instructions implement 261024-aaaaa`, with or without
  `--unattended`, prints no `Changed on main` line.

Two other states give a report:

- Marked `unattended: true`, the same cycle is reported:
  `grove/cycles/261024-kkkkk-keep.yaml: planned: 261024-aaaaa changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan 261024-kkkkk`.
- With the cycle gone, the change is reported against its `base`, in
  check and in the implement instructions.

Source, `internal/cli/check.go` at 2d07433: `stalePlans` skips every
cycle not marked unattended (`if !cy.Meta.Unattended || cy.Done || ...`).
`movedSinceBase` skips every member of a cycle that is not done, leaving it
to the plan's report. `internal/cli/instructions.go` at 2d07433,
`instructionsChange`, names moved files only when the change is in no
cycle that is not done (`cy == nil && !inOpenCycle(tree, id)`). So a
member of an unmarked cycle gets neither report.

The same check without the unattended guard, run on Grove's own main at
2d07433, reports 7 problems where main's build prints `No problems in
grove/.`. They cover every member of cycles 261023-9kfxn and
261023-eyaqa. Both cycles are stale by the cycles spec's rule: for one,
3e40d51 (#65) edited `CLAUDE.md`, which 261023-kh3gg's T-01 names.

## Expected

`grove check` reports the cycle as stale. The scenario "A member of a
cycle" under "Report a named file that moved" in
[grove/specs/propose/spec.md](../../specs/propose/spec.md) says so,
under the same name under "Tell the owner of a moved file" in
[grove/specs/implement/spec.md](../../specs/implement/spec.md). Neither
sets a condition on `unattended`. For an unmarked cycle, the report does
not say that `grove run` skips the cycle.

`grove instructions implement <change-id>` names a member's task files
that changed on the target since the cycle's plan, as it names them
against `base` for a change in no cycle. The owner agreed to this as a new
requirement ([the change spec](specs/implement/spec.md)).

Constitution rule 3 lets `check` and `instructions implement` read tasks
"to tell a task whose named file changed on the target since the change's
`base`, or since the `planned` commit of the cycle it is in".

## Reproduction

From a clean checkout of main at 2d07433, with `grove` being
`go build -o bin/grove ./cmd/grove`:

1. `just fixture <dir>`, then in `<dir>`: `grove init`, which prints
   `Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.`
   Commit everything.
2. Write change `grove/changes/261024-aaaaa-keep-records/`:
   - `.grove.yaml` with `title`, `created`, `schema: default` and
     `base:` the commit from step 1
   - a `proposal.md`
   - a `tasks.md` whose `## T-01 — Keep records` has
     ``Files: `ledger/records.py` ``

   Commit it: this approves it.
3. Write `grove/cycles/261024-kkkkk-keep.yaml` with `title`, `created`,
   `changes: [261024-aaaaa]` and `planned:` the commit from step 2, and no
   `unattended`. Commit it. `grove check` prints `No problems in grove/.`
4. Append a comment to `ledger/records.py` and commit. Then:
   - `grove check` prints `No problems in grove/.` and exits 0.
   - `grove instructions implement 261024-aaaaa | grep "Changed on"`
     prints nothing, and the same with `--unattended`.
5. To compare, add `unattended: true` to the cycle file without
   committing. `grove check` prints
   `grove/cycles/261024-kkkkk-keep.yaml: planned: 261024-aaaaa changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan 261024-kkkkk`
   and `grove: 1 problem found`, and exits 1. Discard the edit.
6. To compare, remove the cycle file and commit. Then:
   - `grove check` prints
     `grove/changes/261024-aaaaa-keep-records/tasks.md: T-01's ledger/records.py changed on main since base f7dc8c0; revise the change with propose`.
   - `grove instructions implement 261024-aaaaa` prints
     ``- Changed on main since base f7dc8c0: T-01's `ledger/records.py`.``,
     and so does `--unattended`.

   The short commits differ from run to run.

## The requirement it breaks

- `grove check`: the scenarios "A member of a cycle" in
  `grove/specs/propose/spec.md` ("Report a named file that moved") and in
  `grove/specs/implement/spec.md` ("Tell the owner of a moved file"). They
  are agreed, so no change spec is written for this part.
  `grove/specs/cycles/spec.md`'s "Report a cycle that cannot run
  unattended" covers only the marked case, and still holds.
- `grove instructions implement`: no maintained requirement covers naming
  a cycle member's files that moved since the plan. The owner agreed that
  it should become agreed behavior. The change spec ADDs "Tell the owner
  of a file moved since the plan" to `implement`.

## Decisions

The owner answered the round "yes" (2026-10-24), accepting each decision
below as recommended.

1. For a cycle not marked `unattended`, `grove check` reports
   `<cycle file>: planned: <change-id> changed on main since this plan; re-plan it with /grove-plan <cycle-id>`.
   A marked cycle keeps today's message, which adds `, so grove run skips
   this cycle`. The report is a problem like any other: check exits 1.
   Reason: the specs call the cycle "reported as stale", and an unmarked
   cycle is not one `grove run` given no cycle would pick.
2. `planned: <x> is not a commit in this repository` is reported for
   unmarked cycles too. Reason: `grove run <cycle-id>` can still be given
   such a cycle and reads `planned`.
3. Implementing this change is breaking under CONTRIBUTING.md "Breaking
   changes": `check` reports a new problem on documents that were clean
   before. Its commit is `fix(cli)!:` with a `BREAKING CHANGE:` footer
   saying to re-plan each stale cycle with `/grove-plan <cycle-id>`.
4. This change does not re-plan Grove's own cycles 261023-9kfxn and
   261023-eyaqa. They are re-planned with `/grove-plan` on their own,
   eyaqa after this change is implemented. Its T-02 edits files that
   261023-mfgew names, which makes eyaqa stale again. CI's `just check`
   does not run `grove check`, so nothing fails meanwhile.
5. `grove instructions implement <change-id>` names a member's moved files
   both with the owner and with `--unattended`. Unattended, the task is
   then blocked, as for a change in no cycle.
6. "Changed since the plan" counts commits exactly as the cycles spec's
   "Count a member changed with its plan" does. A commit that also edits
   the cycle file, or that archives a change of a cycle this one depends
   on, does not count. Reason: a plain diff from `planned` would name a
   file that the merge archiving a cycle it depends on changed. The plan
   expected that merge, and an unattended build would block the task for
   it.
7. `grove instructions implement <cycle-id>` prints nothing new. Each
   member's own `grove instructions implement <change-id>` names its files.
8. The line in a member's instructions reads
   ``- Changed on main since cycle <cycle-id>'s plan <planned sha7>: T-01's `<file>`.``
9. The implement requires text says "changed on the target", dropping
   "since the change's base", so one rule covers both lines with fewer
   words (rule 12).
10. The change spec ADDs one requirement to `implement` with four
    scenarios: with the owner, unattended, after the merge archiving a
    cycle it depends on, and after a re-plan. The maintained requirements,
    including both "A member of a cycle" scenarios, stay as they are.
11. No `depends_on`. No open change at 2d07433 names
    `internal/cli/check.go` or `internal/cli/instructions.go`, or changes
    cycle staleness. Text overlaps without a dependency:
    - 261023-mfgew edits the implement requires files and
      `internal/cli/walkthrough_test.go`, and ADDs to the implement spec.
    - 261023-338jz rewrites `docs/coordinating.md` and
      `docs/planning.md`.
12. `stalePlans` reports every cycle not done that has a `planned`,
    choosing the message by the mark. `movedSinceBase` stays as it is, so
    a member is still not reported against its `base`.
13. Counting the commits since a plan becomes one function, which returns
    the paths those commits changed. Check and `grove run` derive the
    stale members from it, as today. `grove instructions implement`
    filters a member's task files by those paths, against the cycle the
    target holds. A member of a cycle only drafted on the current branch
    is not named, as today.

The routine defaults, accepted with the rest:

- Two tasks: T-01 for check, T-02 for implement.
- The tests extend `TestMovedFiles` and `TestFollows` in
  `internal/cli/walkthrough_test.go`.
- The docs and walkthrough lines named in the tasks change (rule 8).
- Each task carries a `Delegated:` line for test layout and private
  helper names.

## Open questions

None.

## Related documents

- Tasks: [tasks.md](tasks.md)
- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
