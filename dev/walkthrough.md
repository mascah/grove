# Walkthrough

The working test script. It runs every step of
[Getting started](../docs/getting-started.md) on the ledger fixture
(`dev/fixture.md`), with what to type and what to expect, in Claude Code and
again in Codex. The documentation is the specification (CONSTITUTION.md rule
8); each step here names the page section it tests. Grove's output is quoted
only where it carries meaning.

Setup: in the checkout or worktree under test, `eval "$(just dev-path)"`
in the shell the agent is started from. It builds `bin/grove` and puts that
`bin/` first on PATH, so `grove` there is this build and the agent inherits
it; other shells keep the everyday `grove` from `just install`. Start Codex
as `codex -c allow_login_shell=false`: otherwise it runs commands in a login
shell, whose profile can put another `grove` first. `grove version`
(reference, "Terminal commands") shows which one runs: a build in a checkout
prints a pseudo-version naming its commit, such as
`grove 0.4.1-0.<date>-<commit>`, ending `+dirty` when there are
uncommitted changes; a build in a worktree prints `grove dev`, since Go
reads the commit only where `.git` is a folder. After
changing a guide, template or a schema, `just build` is enough; after
changing skill text, run `grove init` again in the fixture. `just fixture ~/tmp/ledger`, `cd ~/tmp/ledger`. For
Codex, use a second fixture and type `$grove-…` where this says `/grove-…`.

## Main line: allow editing older records

1. **Set up** (getting started §1). `grove init`
   > Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.

   The fixture has no `.claude/`, `CLAUDE.md`, `AGENTS.md` or other sign of
   an agent, so `grove.yaml` lists `claude: skills` and `codex: skills`
   under `agents`. `.claude/skills/` and `.agents/skills/` each hold
   `grove-explore`, `grove-prototype`, `grove-propose`, `grove-plan`,
   `grove-implement` and `grove-archive`, and
   every `SKILL.md` has the same body after its name. Under `context` in
   `grove.yaml`, write "Run python3 -m unittest before saying a task is
   done." Then `git add -A && git commit -m "Set up Grove"`. Every skill's
   instructions, such as `grove instructions explore`, then hold a
   `## Project context` section with that line before
   `## What Grove requires of explore`.
   Running `grove init` again rewrites the skills and leaves `grove.yaml` as
   edited.

2. **Explore** (§2). `git switch -c older-edits`, start `claude`, and type
   `/grove-explore Editing an old record fails because its category was retired`.

   The session reads the code, reproduces
   `error: category 'misc' is not allowed` for `ledger edit 3 --note x`, and
   asks one question at a time. It notices that the README ("every record")
   and the CHANGELOG ("new records") disagree, and asks which is intended. I
   answer that older records keep their category when it is unchanged, and
   leave open what happens if the category is changed. It suggests no
   prototype: no open question is one an experiment would answer.
   `git status` shows nothing new: exploring wrote no file.

3. **Propose** (§3). `/grove-propose` in the same conversation. The session
   states the scope and I confirm it. Before writing, it lists, numbered,
   each decision it would otherwise assume, such as which file the fix goes
   in and whether the README changes, each with a recommended answer; I
   say "yes". It then creates `grove/changes/261015-k3m9q-allow-editing-older-records/` with
   `.grove.yaml` (its title and created time), `proposal.md` (my decision
   and my "yes" to the round quoted, no assumptions section, the category
   question open and blocking no task, since T-01 leaves the category
   unchanged),
   `specs/record-editing/spec.md` (an ADDED requirement as an EARS sentence
   with scenarios), and `tasks.md` (T-01 with its acceptance, a `Files:`
   line naming `ledger/cli.py`, the file under fix, and `tests/test_cli.py`,
   which the session said in the conversation it would touch, and an
   `Interfaces:` line, with no code or steps). Before presenting it, the
   session checks that T-01's acceptance covers each scenario and says so,
   or says what it added. It runs `grove check` until:
   > No problems in grove/.

   Its instructions, `grove instructions propose`, give
   `## What Grove requires of propose` before `## How to propose`.

4. **Review and approve** (§4).
   `git add -A && git commit -m "Propose allowing edits to older records"`,
   then `grove review`:
   > Merging older-edits into main approves change 261015-k3m9q as future work; it does not implement it.

   Then each document, labelled change or maintained spec, and the diff.
   `git switch main && git merge --no-ff older-edits`, then `grove changes`:
   > 261015-k3m9q  Allow editing older records  approved

5. **Implement** (§5). `git switch -c older-edits-impl`, a fresh `claude`,
   `/grove-implement 261015-k3m9q`. Its instructions,
   `grove instructions implement 261015-k3m9q`, give
   `## What Grove requires of implement` before `## How to implement`. The
   open question blocks no task, so the session asks nothing, works T-01
   with a test, editing no file but the two T-01 names without asking me,
   tells me how to try it, and writes `implementation.md` with
   T-01's status and evidence. Before finishing it runs the tests, rechecks
   T-01's evidence, and writes `## Scenarios` there, naming each scenario's
   test and "passed", or `unverified` and why; its report names each
   unverified one, and so does the archive before it folds the spec.
   `python3 -m unittest` passes. With every task
   done, the same session runs `grove instructions archive 261015-k3m9q`,
   adds the requirement to a new `grove/specs/record-editing/spec.md`, with
   `grove/specs/record-editing/.grove.yaml` beside it, runs `grove check`,
   then `grove archive 261015-k3m9q`:
   > Archived 261015-k3m9q to grove/changes/archive/261015-k3m9q-allow-editing-older-records/; nothing is committed.

6. **Merge** (§6). Commit, `git switch main && git merge --no-ff older-edits-impl`. `grove specs`:
   > record-editing  Record editing  grove/specs/record-editing/spec.md

   The spec says "When a user edits an existing record and leaves its
   category unchanged, the system shall save the edit, even if that category
   is no longer allowed for new records." `grove changes` lists nothing. `grove`
   opens the reader on the spec.

7. **Next** (§Next). A fresh `/grove-explore Editing a missing record crashes`
   cites that requirement as agreed and lists no approved changes. Its
   instructions are under 12 KiB.

## Coordinating changes

It starts where the main line's step 6 ends: one archived change and the
`record-editing` spec. Each step names the section of
[Coordinating changes](../docs/coordinating.md) it tests.

1. **A change that waits** ("A change that waits"). From main,
   `git switch -c unknown-id`, and in a fresh session
   `/grove-propose Editing an unknown ID prints an error`. It finds no other
   open change, says so, and leaves `depends_on` out. Commit, then
   `git switch main && git merge --no-ff unknown-id`.
   `git switch -c delete`, and in a fresh session
   `/grove-propose Add ledger delete ID, from the README's planned list`. The
   session asks whether it needs another open change first and names the
   unknown-ID change, since `delete` would use `find`. I say yes.
   `.grove.yaml` has `depends_on: [<unknown-id>]` and `base:` with main's
   full commit. Commit and merge. `grove changes`:
   > \<unknown-id>  Editing an unknown ID prints an error  approved

   > \<delete-id>  Add ledger delete ID  approved  waits on \<unknown-id>

2. **Implement waits** ("A change that waits"). `git switch -c delete-impl`,
   a fresh session, `/grove-implement <delete-id>`. Grove prints
   > grove: \<delete-id> waits on \<unknown-id>, which this working tree has not archived as finished

   The session tells me and builds nothing. `git switch main && git branch -D
   delete-impl && git switch -c unknown-id-impl`, a fresh session,
   `/grove-implement <unknown-id>`, which builds and archives it. Commit and
   merge. `grove changes` lists the delete change without `waits on`, and
   `/grove-implement <delete-id>`, on a new branch `delete-impl`, builds it.
   Leave that branch unmerged for now.

3. **A change that waits on nothing** ("A change that waits"). From main,
   `git switch -c budgets`, a fresh session, `/grove-propose Monthly
   budgets in report`. It needs no other change first, and its
   `.grove.yaml` has no `depends_on`. Commit and merge: two changes are
   approved, delete and budgets, the unknown-ID change being archived.
   `grove changes` lists both without `waits on`.

4. **When the plan is stale** ("When the plan is stale"). From main,
   `git switch -c recategorize`, a fresh session, `/grove-propose An older
   record's category can be changed, only to a current category`. When it
   offers a new requirement, I ask it to fold the rule into the main line's
   requirement instead, so its change spec MODIFIES that requirement.
   Commit and merge. Then
   `git switch main && git switch -c show-saved`, a fresh session,
   `/grove-propose Editing a record prints what was saved`, which MODIFIES
   the same requirement from the same text. Commit and merge. Implement
   `recategorize` on `recategorize-impl`; the `grove check` its archive runs
   after reconciling the spec reports show-saved, not recategorize, and
   exits non-zero:
   > grove/changes/\<show-saved>/specs/record-editing/spec.md:\<line>: requirement "\<name>" was rewritten since base \<sha7>; revise the change with propose

   The session archives recategorize anyway, since the report is not about
   it. Commit and merge; `grove check` on main reports the same line. Set
   show-saved's `depends_on` to `[<recategorize-id>]` in the working tree:
   `grove check` is clean, since show-saved then builds on the change that
   rewrote the requirement. Discard that edit.
   `git switch -c show-saved-revise`, a fresh session,
   `/grove-propose <show-saved>`, asked to bring it up to date with main.
   The session adds recategorize's scenario to show-saved's MODIFIED text
   and sets `base` to main's commit. `grove check` is clean. Commit and
   merge.

5. **A file that moved** ("When the plan is stale"). The budgets change's
   T-01 names `ledger/report.py` on its `Files:` line. On a branch, add a
   comment to `ledger/report.py`: `grove check` there is clean, since only
   the branch changed it. Commit and merge; on main, `grove check` reports
   > grove/changes/\<budgets>/tasks.md: T-01's ledger/report.py changed on main since base \<sha7>; revise the change with propose

   `grove instructions implement <budgets-id>` names T-01 and
   `ledger/report.py` as changed since base, and a fresh session's
   `/grove-implement <budgets-id>`, reaching T-01, tells me what changed
   and asks before building; I stop it there. Write `Files: ledger/report.py`
   without backticks in a task: `grove check` reports the line and exits
   non-zero. Discard that edit.

6. **The project as data** ("The project as data"). `grove changes --json`
   shows the unknown-ID and recategorize changes as `archived` with `T-01` `done`, and the delete change with no
   `waits_on`, since the unknown-ID change is archived.

## Planning cycles

It starts where coordinating changes ends: the delete, budgets and
show-saved changes approved. Each step names the section of
[Planning cycles](../docs/planning.md) it tests.

1. **Plan a cycle** ("Plan a cycle"). From main, `git switch -c plan-0.4`,
   a fresh session, `/grove-plan 0.4: delete and budgets`. Its
   instructions, `grove instructions plan "0.4: delete and budgets"`, give
   `## What Grove requires of plan` before `## How to plan`. The session
   lists the approved changes and I agree to delete and budgets. It asks
   about budgets, whose task has the report read clearly beside each
   category's total, whether it needs a person; I say yes, to judge the
   wording. It asks nothing about delete, nor about other cycles, there
   being none. It writes
   `grove/cycles/<cycle-id>-<slug>.yaml` with both under `changes`,
   budgets under `needs_human` with my reason, `created:` with the time,
   `planned:` with main's full commit and no `unattended`. `git status` shows only that file, and
   `grove check`:
   > No problems in grove/.

   It ends by naming `/grove-implement <cycle-id>` first, then saying
   `grove run` would refuse the cycle, given or marked, because budgets
   needs a person, and that revising budgets with
   `/grove-propose <budgets-id>` so its acceptance can be checked without
   a person, which re-plans the cycle on the same branch, and one merge
   removes that.

2. **Merge the plan** ("Merge the plan"). Commit, then `grove review`:
   > adds     cycle            grove/cycles/\<cycle-id>-\<slug>.yaml

   Merge it. `grove changes` ends the delete and budgets lines with
   `in cycle <cycle-id>`, and the show-saved line with nothing.

3. **A cycle of one** ("Plan a cycle"). `git switch -c plan-show-saved`, a
   fresh session, `/grove-plan <show-saved>`. It asks whether to keep the
   change's title for the cycle and, finding nothing that needs a person,
   writes a cycle listing only show-saved. It ends by naming
   `/grove-implement <cycle-id>`, then saying `grove run <cycle-id>` would
   build the cycle, then that `unattended: true`, added to its file in a
   merge of its own, hands it to a `grove run` given no cycle. Commit and
   merge.

4. **Mark it unattended** ("Mark it unattended"). `git switch -c mark`, add
   `unattended: true` to the show-saved cycle, `grove check` is clean;
   commit and merge. `grove changes` ends the show-saved line with
   `in cycle <cycle-id>, unattended`. Add the same line to the 0.4 cycle:
   `grove check` reports it and exits non-zero:
   > grove/cycles/\<0.4-cycle>.yaml:\<line>: unattended: \<budgets-id> needs a person, as needs_human says, so grove run skips this cycle; remove unattended, or build the cycle with /grove-implement \<0.4-cycle-id>

   Discard the edit.

5. **A stale plan** ("Mark it unattended", "Re-plan a cycle"). On a branch,
   add a sentence to show-saved's proposal, commit and merge. On main,
   `grove check`:
   > grove/cycles/\<show-saved-cycle>.yaml: planned: \<show-saved> changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan \<show-saved-cycle-id>

   `git switch -c replan`, a fresh session, `/grove-plan <show-saved-cycle-id>`.
   It reads show-saved as main holds it, asks nothing new, keeps
   `unattended`, and sets `planned` to main's commit. `grove check` is
   clean. Commit and merge.

   Then add a comment to a file show-saved's T-01 names on its `Files:`
   line, commit and merge: `grove check` on main reports the same
   `planned:` line, and nothing against show-saved's `base`. Re-plan it
   again as above: `grove check` is clean, and show-saved's files are
   unchanged.

   Then, on a branch, add a sentence to show-saved's proposal and commit,
   and set the cycle's `planned` to main's commit and commit: two commits.
   Fast-forward main onto them (`git merge --ff-only`): `grove check`
   reports the same `planned:` line, since the revision reached main
   without the cycle file. Re-plan it as above.

6. **One change in two cycles** ("Re-plan a cycle"). Copy the 0.4 cycle to
   `grove/cycles/261099-zzzzz-later.yaml`. `grove check` reports
   > grove/cycles/261099-zzzzz-later.yaml:\<line>: changes: \<delete-id> is also in cycle \<0.4-cycle-id>, which is not done

   and the same for budgets. Delete the copy.

7. **A cycle that follows another** ("Follow another cycle", "Re-plan a
   cycle"). Approve `/grove-propose Import records from a CSV` (import), its
   task naming a file show-saved's T-01 names. `git switch -c plan-0.5`, a
   fresh session, `/grove-plan 0.5: import`. It asks whether the new cycle
   follows the 0.4 or the show-saved cycle; I name show-saved. It writes
   `depends_on: [<show-saved-cycle-id>]`, and ends saying `grove run` would
   not build it until the show-saved cycle is done, and that the merges
   archiving its changes leave the plan fresh. Commit and merge, then add
   `unattended: true` in a merge of its own. `grove changes` ends import's line with
   `in cycle <0.5-cycle-id>, unattended, waits on cycle <show-saved-cycle-id>`.

   Set the 0.5 cycle's `depends_on` to `[261099-zzzzz]`: `grove check` reports
   > grove/cycles/\<0.5-cycle>.yaml:\<line>: depends_on: 261099-zzzzz is no cycle

   Put it back, and add `depends_on: [<0.5-cycle-id>]` to the show-saved
   cycle: `grove check` reports each one's `depends_on` line, as
   > grove/cycles/\<0.5-cycle>.yaml:\<line>: depends_on: \<show-saved-cycle-id> waits on this cycle in turn, so neither can run

   Discard both edits. On a branch, a fresh session,
   `/grove-implement <show-saved-cycle-id>` builds and archives show-saved,
   editing the file import names; squash-merge it. On main, `grove check`
   is clean, and `grove changes` ends import's line with
   `in cycle <0.5-cycle-id>, unattended`. Then add a comment to that file,
   commit and merge: `grove check` reports
   > grove/cycles/\<0.5-cycle>.yaml: planned: \<import> changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan \<0.5-cycle-id>

   Re-plan it as step 5 does.

8. **Building a cycle unattended** (coordinating, "Building a cycle
   unattended"). From main, `git switch -c build-0.4`, a fresh session,
   `/grove-implement <0.4-cycle-id> --unattended`. Its instructions list
   delete and budgets in build order, each to be built with
   `grove instructions implement <change-id> --unattended`. The session
   runs `grove instructions implement <delete-id> --unattended`, whose
   `## The cycle` lists both and marks delete `this change`, builds and
   archives delete, then does the same for budgets, whose list marks delete
   `archived in this working tree`; it builds budgets though `needs_human`
   names it, and asks nothing. `grove instructions implement
   <0.4-cycle-id> --unattended` on the branch then lists nothing to build.
   One commit, one merge;
   `grove changes` lists neither. `/grove-implement <0.4-cycle-id>` on
   another branch instead builds the two with me, asking me where each
   change's instructions say. `git branch -D delete-impl`.

## Running unattended

Each step names the section of
[Running unattended](../docs/running.md) it tests. Claude Code only.

1. **Set up a checkout** ("Set up a checkout", "Tell the agent to commit and
   open a pull request"). `just fixture ~/tmp/ledger`, `grove init`, add
   the two `rules.skills.implement` lines from the page to `grove.yaml`, commit.
   Create a private, disposable GitHub repository, add it as `origin` and
   push `main`. `git clone` it to `~/tmp/ledger-run`, the operator's
   checkout, and in it write `.claude/settings.local.json` and the
   `.git/info/exclude` line from the page. In `~/tmp/ledger-run`,
   `grove run`:
   > Nothing to apply on main at \<sha7>.

2. **An ordinary run** ("An ordinary run"). In `~/tmp/ledger`, approve two
   changes as coordinating step 1 does, the delete change waiting on the
   unknown-ID change, and push `main`. In `~/tmp/ledger-run`, `grove run`
   says there is nothing to apply: neither change is in a cycle. In
   `~/tmp/ledger`, plan the two as a cycle as planning step 1 does, with
   nothing under `needs_human`, merge and push: still nothing to apply.
   Mark it `unattended`, merge and push. In `~/tmp/ledger-run`, `grove run`.
   While it runs, `grove run` in a second shell there says
   > grove: a run is in progress in this checkout

   and `grove run` in another project's checkout proceeds. While the first
   runs, it prints one row per event as it comes, each with a time and a
   mark: a green `▶` start naming the model, `•` messages, cyan `›` tool
   calls such as `› Bash grove instructions implement <change-id> --unattended`,
   no row for a tool result unless it is a red `✗` error, faint `·`
   thinking whose `(×N)` count goes up in place, and a `■` result with its
   turns and cost. Narrow the terminal: new rows are cut to its width. It
   ends
   > Attempt \<attempt> ended: exit status 0; its output is \<checkout>/.git/grove/\<attempt>.log.

   `head -1` of that log is a JSON object, as is every line after it.

   The checkout is on `grove/run-<attempt>`, whose commits build and archive
   both changes, the unknown-ID one first, each task with its ID. The
   record's prompt was `/grove-implement <cycle-id> --unattended`, and the
   output shows the session ran `grove instructions implement <change-id>
   --unattended` once for each. The pull request lists both.
   `.git/grove/<attempt>.json` names the cycle, both changes and `main`'s
   commit, and `git status` shows nothing.

3. **The pause** ("The pause"). `grove run` says
   > grove: the project is paused after attempt \<attempt>, whose output is \<checkout>/.git/grove/\<attempt>.log; once you have dealt with it, run grove run --retry

   In `~/tmp/ledger`, approve the budgets change, plan it as a cycle of one
   with `/grove-plan <budgets-id>`, merge it, mark it `unattended`, merge
   and push: still paused. Merge the pull request:
   `grove run > ~/tmp/run.txt` starts an attempt that selects only the
   budgets cycle. Let it finish and close its pull request unmerged.
   `~/tmp/run.txt` holds the same rows as a terminal shows, with no color
   codes, none cut short, and each thinking count on one line.

4. **When a session is stopped** ("When a session is stopped").
   `grove run --retry`, and once it prints an `Edit:` or `Write:` line,
   press Ctrl-C:
   > grove: attempt \<attempt> ended: cancelled; its output is \<checkout>/.git/grove/\<attempt>.log

   `ps -o pid,pgid,command -g <pid of that claude -p>` lists nothing, and
   the branch, the edit and the output are there; the log holds each event
   up to the stop, one JSON object per line, the edit's tool call among them. `grove run --retry` says
   > grove: the checkout has uncommitted changes; commit, stash or discard them first

   and the edit is untouched. Discard it as the page shows, then
   `grove run --retry --timeout 1m`: after a minute it ends `timed out`,
   leaving no process of the session.

5. **Refusals and failures** ("An ordinary run", "Retry"). Commit a file
   on the operator's `main` and `grove run` from it:
   > grove: the local main has commits origin's does not, so it cannot be fast-forwarded

   Reset it. With `agents` holding only `codex` on `origin`'s `main`,
   `grove run --retry` says
   > grove: run needs Claude Code listed under agents in grove.yaml, as skills or commands

   Put `claude` back. With `PATH` missing `claude`, `grove run --retry`
   ends `failed to start`, the reason in its output, and the next
   `grove run` says the project is paused.

6. **Run a cycle you name** ("Run a cycle you name"). In `~/tmp/ledger`,
   approve `/grove-propose Editing a record prints what was saved`
   (show-saved), plan it as a cycle of one, merge and push, without
   marking it. Approve `/grove-propose Export the ledger as CSV` (export),
   plan it as a cycle of one, telling the session export needs a person to
   judge the file, merge and push. In
   `~/tmp/ledger-run`, `grove run --retry <export-cycle-id>`:
   > grove: \<export-id> needs a person, as cycle \<export-cycle-id>'s needs_human says; build the cycle with /grove-implement \<export-cycle-id>

   and `git branch` lists no new branch. In `~/tmp/ledger`, on a branch, a
   fresh session, `/grove-propose <export-id>`, asking for export's
   acceptance to be checked by a test that reads the CSV back. The session
   prints the export cycle's file with the change, revises export, asks
   whether it still needs a person, takes it out of `needs_human` when I
   say no, sets `planned` to main's commit and adds no `unattended`.
   `grove check` is clean. Squash-merge the branch and push. In
   `~/tmp/ledger-run`, `grove run --retry <export-cycle-id>` starts a
   session; stop it with Ctrl-C and discard what it left.
   `grove run --retry 261099-zzzzz`
   says `main` holds no such cycle. In `~/tmp/ledger`, approve
   `/grove-propose Import records from a CSV` (import), plan it as a cycle
   of one that follows the export cycle, merge and push. In
   `~/tmp/ledger-run`, `grove run --retry <import-cycle-id>`:
   > grove: cycle \<import-cycle-id> waits on cycle \<export-cycle-id>, which is not done

   and `git branch` lists no new branch. `grove run --retry <show-saved-cycle-id>`
   starts a session on `/grove-implement <show-saved-cycle-id> --unattended` though the
   budgets cycle is marked, and its record names the show-saved cycle.

7. **Without a remote** ("Without a remote"). `just fixture
   ~/tmp/ledger-local`, `grove init`, commit, approve the unknown-ID
   change, plan it as a cycle of one, mark it `unattended` and merge, with
   `Commit after each task, naming its task ID.` and the page's
   commit-what-is-left line under `rules.skills.implement`; write `.claude/settings.local.json` and the
   `.git/info/exclude` line from the page. `git remote` lists nothing.
   `grove run` fetches nothing, ends `exit status 0`, and its record holds
   `main`'s commit. `grove run` says the project is paused. Merge
   `grove/run-<attempt>` into `main`: `grove run` says there is nothing to
   apply.

## Other paths

- **Propose from a description** (§3). On a new branch, with no exploration,
  `/grove-propose Editing a missing record should print an error` writes a
  change.
- **Prototype before proposing** (§2, "Try it first"; concepts, "The
  loop"). After the main line, from main, `git switch -c budgets`, start
  `claude`, and type `/grove-explore Monthly budgets in report, from the
  README's planned list`. When I say I do not know how a budget should read
  beside each category's total, the session names that question and asks
  whether to prototype it, building nothing. I say yes; it runs
  `grove instructions prototype "<question>"`, creates
  `prototype/<slug>` from `budgets`, builds a throwaway script that prints
  the report with budgets for cases under, at and over budget, runs it, and
  reports the command, what it showed, a verdict and its limits. It asks
  whether to commit; I say yes, and it commits on `prototype/<slug>` only.
  Asked to propose, it switches back to `budgets` and the proposal states
  the finding and cites `prototype/<slug>` at that commit, linking no
  prototype file. `git diff main prototype/<slug> -- grove` is empty, and
  main is unchanged.
- **Propose checks the change** (§3). On a new branch, `/grove-propose Add
  ledger delete ID, and print the deleted record`, answering "yes" to the
  round. Before presenting the change, the session checks that a task's
  acceptance covers each requirement, deleting and printing alike; where
  one did not, it says what it added.
- **The round in order** (§3). From main, `git switch -c unknown-id`, a
  fresh session, `/grove-propose Editing an unknown ID prints an error`,
  keeping the `default` schema. Its round asks the exit status before
  whether `find` or `ledger/cli.py` holds the check, and puts the test's
  name and place on a last line of defaults that one "yes" accepts. Say
  "Yes, and let the implementer lay out tests and add
  `tests/test_records.py` without asking": T-01 gets a `Delegated:` line
  naming them, and its `Files:` line does not name `tests/test_records.py`.
  Answering just "Yes" to a round that offers no delegation leaves no task
  with a `Delegated:` line. `grove check` is clean.
- **Split a change** (§3). From main, `git switch -c delete`, a fresh
  session, `/grove-propose Add ledger delete ID, and an error for an
  unknown ID`. While confirming scope, the session says these are two
  changes I could agree separately. Say to split it: it runs `grove
  instructions propose` once for each and writes both on `delete`, the
  delete change naming the unknown-ID change in `depends_on`. `grove check`
  is clean.
- **Delegated with me, ignored unattended** (§5; coordinating, "Building a
  cycle unattended"). Approve a change whose T-01 names `ledger/records.py`
  and `tests/test_cli.py` on its `Files:` line and `tests/test_records.py`
  on a `Delegated:` line right after it. `/grove-implement <change-id>`
  creates `tests/test_records.py` for `find`'s own tests without asking.
  On a fresh fixture with the same change in a cycle of one,
  `grove run <cycle-id>` leaves `tests/test_records.py` uncreated, puts
  `find`'s tests in `tests/test_cli.py`, and T-01's evidence names
  `tests/test_records.py` and why it was left.
- **Revise a draft** (§3). `/grove-propose <change-id>` in a new conversation
  is given that change's files in full and edits them.
- **A question that blocks a task** (§5). On a new branch,
  `/grove-propose Editing a missing record should print an error`; when
  asked whether it should exit 2 like other errors, defer it. The proposal
  names T-01 as blocked by that question. Merge it, plan it as a cycle of
  one with `/grove-plan <change-id>` and merge that, then, on another
  branch, build it headless: `claude -p "/grove-implement <cycle-id> --unattended"`. It builds
  nothing, `implementation.md` says `Status: blocked` under T-01 with the
  question, and its last message names the question. `grove archive
  <change-id>` refuses.
- **Implement asks, and revises with the owner** (§5). Continuing, on
  another branch, a fresh `claude`, `/grove-implement <change-id>`. Before
  building T-01 the session asks whether a missing record should exit 2. I
  answer yes; it says it would edit the change spec's requirement and its
  scenario, and the open question in the proposal, to say exit 2, and
  waits. I agree: it edits them, runs `grove check` until
  > No problems in grove/.

  builds T-01 to exit 2 with a test, and archives the change. Commit, then
  `grove review` shows the edited change spec beside the code. Had I
  declined, the change's files would stay as they are. Had I said to leave
  it, T-01 would be `Status: blocked` with the question.
- **One merge for a small change** (§3, §4; concepts, "Approval is a
  merge"). From main, `git switch -c rename`, a fresh session,
  `/grove-propose Editing a missing record should print an error`. Its
  finish offers both next steps and picks neither. Commit, then
  `/grove-implement <change-id>` on `rename`: the instructions call the
  change a draft that merging `rename` approves with its code, and the
  session builds T-01 with a test and archives the change. Commit.
  `grove review` says
  > Merging rename into main approves change \<change-id> with its code.

  Merge `rename` once: `grove changes` does not list the change, and
  `grove specs` lists `record-editing`. Before that merge,
  `grove instructions implement <change-id> --unattended` on `rename` said
  > grove: main holds no approved change \<change-id>
- **A title that drifted** (reference, "Files"). Change `title` in a
  change's `.grove.yaml` and run `grove check`:
  > grove/changes/<change-id>-<slug>/proposal.md:1: the first heading should be "# <new title>", the title in .grove.yaml
- **A change without its proposal** (reference, "Files"). After step 3,
  delete the change's `proposal.md` and run `grove check`. It exits non-zero:
  > grove/changes/<change-id>-<slug>/proposal.md: missing; it is the change's main document in schema default

  `grove instructions propose <change-id>` still finds the change, so a
  session can write the proposal again.
- **Your own tasks template** (reference, "Change schemas"). On a new
  branch, copy the tasks template from `grove instructions propose "x"` (the
  `<file>` whose path ends in `tasks.md`) to
  `grove/schemas/default/documents/tasks/template.md` with `{{change}}` where it
  shows the new change's ID, and add a line
  `Rollback: [how to undo this task]` to its task. `grove instructions
  propose "x"` now gives that template. `/grove-propose Editing a missing
  record should print an error` writes `tasks.md` with a `Rollback:` line in
  each task, and `grove check` is clean.
- **A bug** (reference, "Grove's bugfix"). From a fixture with no
  `grove/schemas/`, on a new branch, in a fresh session,
  `/grove-propose Editing a missing record crashes`. The session sees
  `default` and `bugfix` under "Change schemas", suggests `bugfix` itself,
  and waits; I confirm. It runs `grove instructions propose --schema bugfix
  "..."`, and follows that: it reproduces the `StopIteration` traceback of
  `ledger edit 99 --note x`, and finds no maintained requirement for a
  missing record, so it asks whether the corrected behavior should be
  agreed. I say yes. It writes `proposal.md` with Observed, Expected,
  Reproduction and The requirement it breaks, a change spec adding "editing
  a missing record prints an error" to `record-editing`, `tasks.md` whose
  T-01 starts with a test that reproduces the crash, and `.grove.yaml` with
  `schema: bugfix`. Commit; `grove check` is clean. Merge it, then on another branch
  `/grove-implement <change-id>` builds it (its instructions name schema
  `bugfix`) and `/grove-archive <change-id>` adds the requirement to
  `grove/specs/record-editing/spec.md` and archives the change.
- **A mistake in a schema** (reference, "Change schemas", "Grove's
  bugfix"). After "A bug", write a `grove/schemas/bugfix/` of my own: a
  `schema.yaml` of `description: A fix.` and `instrution: x`, and a
  proposal template at `documents/propsal/template.md`, misspelled. Commit;
  `grove review` labels the schema's files `change schema`, not
  `maintained spec`.
  `grove check` reports both and exits non-zero:
  > grove/schemas/bugfix/schema.yaml:2: unknown field "instrution"; schema.yaml takes description, documents and skills

  > grove/schemas/bugfix/documents/propsal/template.md: schema bugfix does not use this file

  `grove changes` still works.
- **A mistake in grove.yaml** (reference, "grove.yaml"). Add `derp: [x]`
  under `rules`, `propsal: [x]` under `rules: documents:`, and a top-level
  `rule: {}` in `grove.yaml`. `grove check` reports all three and exits
  non-zero:
  > grove.yaml:<line>: unknown rules key "derp"; rules takes documents and skills

  > grove.yaml:<line>: unknown rules.documents key "propsal"; kinds are proposal, change-spec, tasks, implementation, spec

  > grove.yaml:<line>: unknown field "rule"; grove.yaml takes root, target, agents, schema, context and rules

  `grove changes` still works.
- **Agents from the repository** (reference, "Harnesses"). On a fresh
  fixture, `mkdir -p .claude .opencode .agents/skills/notes`, then
  `grove init`:
  > Wrote grove.yaml and the Grove skills for claude and opencode; nothing is committed.

  `grove.yaml` lists `claude: skills` and `opencode: skills`;
  `.claude/skills/` and `.agents/skills/` hold the six skills, the same
  file in each.
- **Commands** (reference, "Harnesses"). Continuing, add an `AGENTS.md`,
  set `claude: commands` in `grove.yaml`, and run `grove init`. Later runs
  read only `grove.yaml`, so codex stays out:
  > Kept grove.yaml and wrote the Grove skills for opencode and commands for claude, and removed grove-archive, grove-explore, grove-implement, grove-plan, grove-propose and grove-prototype from .claude/skills; nothing is committed.

  `.claude/commands/grove/` holds `explore.md`, `prototype.md`,
  `propose.md`, `plan.md`, `implement.md` and `archive.md`. In Claude Code, `/grove:propose Editing a missing record
  should print an error` runs `grove instructions propose "Editing a
  missing record should print an error"`.
- **Drop a harness** (reference, "Harnesses"). Continuing, delete the
  `opencode` line and run `grove init`:
  > Kept grove.yaml and wrote the Grove commands for claude, and removed grove-archive, grove-explore, grove-implement, grove-plan, grove-propose and grove-prototype from .agents/skills; nothing is committed.

  A skill of my own, `.agents/skills/notes/`, is still there.
- **A harness Grove does not have** (reference, "grove.yaml"). Set
  `agents` to `codex: commands` and `zed: skills`. `grove check` reports
  both and exits non-zero:
  > grove.yaml:\<line>: harness "codex" has no form "commands"; it has skills

  > grove.yaml:\<line>: unknown harness "zed"; agents takes claude, codex, cursor, hermes, opencode and pi

  `grove init` stops at the first, writing nothing.
- **A dependency that can never be met** (coordinating, "A change that
  waits"). After coordinating step 2, set the delete change's `depends_on`
  to `[260101-zzzzz]`, and `grove check` reports it and exits non-zero:
  > grove/changes/\<delete-id>-\<slug>/.grove.yaml:\<line>: depends_on: 260101-zzzzz is no change, open or archived

  Then propose another change on a branch with `depends_on: [<delete-id>]`,
  and set the delete change's to `[<that change>]`. `grove check` reports
  each with its `path:line`:
  > grove/changes/\<delete-id>-\<slug>/.grove.yaml:\<line>: depends_on: \<other-id> waits on this change in turn, so neither can be applied

  `grove changes` still works. Discard the edits.
- **A change left out of a cycle** (coordinating, "Building a cycle
  unattended"). With an approved change whose `depends_on` names
  `260101-zzzzz` in a cycle, `grove instructions implement <cycle-id>
  --unattended` ends its list with
  > Leave out \<id>, "\<title>": it waits on 260101-zzzzz, which is neither being built nor archived as finished. Grove's `check` says why.
- **Implement a cycle unattended, headless** (coordinating, "Building a
  cycle unattended"). From main with a cycle of two changes planned,
  `claude -p "/grove-implement <cycle-id> --unattended"` with
  `implement: [Commit after each task, naming its task ID.]` under
  `rules: skills:` in `grove.yaml`. It creates a branch, builds and
  archives both, commits, asks nothing, and its last message says which
  changes were archived.
- **Commit as you go** (reference, `rules`). With
  `implement: [Commit after each task, naming its task ID.]` under
  `rules: skills:` in `grove.yaml`, `/grove-implement <cycle-id>
  --unattended` commits T-01 with its ID in the message, and so does
  `/grove-implement <change-id>` with me.
- **One change, unattended** (coordinating, "Building a cycle
  unattended"). Approve a change in no cycle whose proposal has an open
  question naming T-02, and from main, `claude -p "/grove-implement
  <change-id> --unattended"`. It creates a branch, builds T-01, marks T-02
  `Status: blocked` with the question, asks nothing, leaves the proposal,
  specs and tasks as they were, and does not archive.
- **Upgrading from apply** (reference, "grove.yaml", "Harnesses"). In a
  project set up by a Grove that still had apply, whose `grove.yaml` keeps
  its rules under `rules.skills.apply`, `grove check` reports
  > grove.yaml:\<line>: unknown rules.skills key "apply"; its lines go under rules.skills.implement

  With a `.claude/skills/grove-apply/` an earlier Grove wrote, `grove init`
  removes it and names `grove-apply` among what it removed. Moved under
  `rules.skills.implement`, the lines are given with the implement
  instructions, and `grove check` says there are no problems.
- **Archive checks the fold** (reference, `grove archive`). After the main
  line, propose and merge a change that MODIFIES the requirement the main
  line added, and on `fold-impl` mark its task done and paste the
  requirement into `grove/specs/record-editing/spec.md` with one word
  changed. `grove archive <change-id>` says
  > grove/specs/record-editing/spec.md:\<line>: requirement "\<name>" is not as \<change-id> gives it

  and leaves the folder where it is. Copy the requirement word for word,
  wrapped at other places, and `grove archive <change-id>` archives it.
  Paste `## ADDED Requirements` above it, or delete its scenario, and
  `grove check` reports that line.
- **Drop a change** (concepts, "The loop"). `/grove-archive <change-id>` on
  an approved change nobody will build asks why, writes `Status: dropped`
  and the reason in `implementation.md`, and moves it without touching the
  specs. Before that, `grove archive <change-id>` says
  > T-01 is not done in implementation.md; archive a change when every task is done or the file says Status: dropped

## Decided

Owner, asked whether to add each:

- A review session or `instructions review <id>`, 2026-10-14: "No".
- A top-level `implement <change-id>`, 2026-10-14: "No".
- Whether a session can confirm the Grove files it wrote are valid,
  2026-10-14: "yes"; how: "by calling the CLI?".
- Bubble Tea and Glamour, 2026-10-14: "We need bubble tea and glamour".
- Codex reads repository skills from `.agents/skills/`, typed
  `$grove-<skill>`, 2026-10-15: "I know this is how it works. We don't
  need to verify this."
- The rest of this page's shape is decisions D-33 to D-43 of 2026-10-15.

## Open questions

- Does the reader showing `implementation.md` count as reading it under rule
  3? Owner: "No, but im not certain." Until settled, the reader shows it.
- The wording of phase 3's `check` messages and the JSON's `root`-relative
  paths. Owner's lean, from Remote Control, 2026-10-20: "the wording is
  fine for now".

## Assumptions to confirm

Agent inference, not the owner's decisions (`CLAUDE.md`).

- `grove.yaml` fields and defaults: `root: grove`, `target: main`,
  `agents` (left out or empty, `claude` and `codex` as skills), `schema: default`, `context`, and `rules`, with `documents` keyed by
  `proposal`, `change-spec`, `tasks`, `implementation` and `spec`, and
  `skills` keyed by a skill.
- `init` never overwrites `grove.yaml` and always rewrites the skills.
  Its second run says "Kept grove.yaml and wrote the Grove skills for claude
  and codex; nothing is committed." A field left out of `grove.yaml` takes
  its default. A harness Grove has no row for, or a form it does not
  have, stops `init` before it writes anything, with `check`'s words.
- `agents` takes only the map of harness to form. The list form,
  `agents: [claude, codex]`, which no release wrote, is a YAML error, not a
  `check` problem; v0.1.0 is the first release, so the change is not marked
  breaking.
- `init` removes Grove's files it no longer writes before it writes the
  rest, so a skill folder linked to another, such as `.claude/skills` to
  `.agents/skills`, keeps its skills.
- `init`'s message groups harnesses by form, skills first ("the Grove
  skills for opencode and commands for claude"), and names what it removed
  by folder ("grove-archive, … from .claude/skills"). It leaves a folder it
  emptied, such as `.claude/commands/grove/`, in place; Git ignores it.
- Detection looks only at the top of the repository, and a `CLAUDE.md` or
  `AGENTS.md` there counts even when it is a link to the other.
- A skill runs `grove` from the PATH of the agent's shell. Every
  `SKILL.md` has one body: it names the user's words `<args>` and says
  that is what the user typed after the skill's name, or nothing (Claude
  Code appends them as `ARGUMENTS:`). A command passes them as
  `$ARGUMENTS`, and its file has only a `description` above its body.
- OpenCode's skills are not typed: Grove writes no `opencode/slash`
  metadata, since the skill files are shared with other harnesses; typing
  is what OpenCode's commands are for.
- Pi's one-time step is trusting the project when Pi asks, or `/trust`.
- `review` takes no argument and compares the current branch's last commit
  with the target, so uncommitted files are not shown. A change counts as
  approved by the merge when the target holds no change with its ID; with
  none, the first line is "Merging <branch> into main approves no new
  change."
- `instructions propose` tells a change ID (`YYMMDD-xxxxx`) from a
  description by its form.
- `archive` moves the folder, refusing a fold that does not match; the
  agent does the reconciliation, from `instructions archive`. Both read the change from the working tree, where
  `implementation.md` is, and neither checks that the change is approved.
- Every skill commits only when the owner or its `rules` ask; otherwise
  the owner commits after archiving (getting started §6).
- A blocked task's question is answered by revising the change with
  propose and merging it, not by the owner telling the unattended session.
- Implementing unattended leaves files no task names, such as the README or a changelog,
  as they are; propose makes such updates tasks when the owner wants them.
- `changes` lists approved changes that are not archived; drafts on the
  current branch are shown by `review`, not by `changes`.
- IDs use the UTC date.
- Explore's instructions offer no document to fill in.
- A change's ID is the start of its folder name. `check` reports a change
  folder whose name does not start with an ID.
- A new change's folder name is its ID and a slug the agent chooses from the
  title; Grove allocates the ID and `created` timestamp.
- The fixture's "Planned for 0.4" leaves out the direction note's
  `imports/2024.csv` and its import item; no step needs them.
- The project index lists each approved change's capabilities (the
  `<capability>` folders of its files) and what it waits on, rather than
  dependencies Grove works out as candidates. Propose names the likely ones
  from it.
- `grove changes` writes `waits on <id>, <id>` after `approved`, and
  nothing for a change whose dependencies are all archived as finished.
- A change with no dependency leaves `depends_on` out; `depends_on: []` is
  reported like any other empty field.
- `base` comes from `{{commit}}`, a placeholder only Grove's `.grove.yaml`
  template fills; `{{base}}` stays the short commit in schema templates.
- The docs page for coordinating changes is `docs/coordinating.md`.
- `check` and `instructions implement` look for what `depends_on` names among
  the working tree's changes, open and archived, and not on the target; a
  branch holds the target's archive as of where it started. `changes` reads
  the target.
- The stale check compares a requirement's text below its heading, trimmed,
  by name: the working tree's maintained spec against what it was at
  `base`, what the change itself gives it, and what each change it depends
  on, directly or through others, gives it. One missing from the working
  tree was removed. A requirement the maintained spec did not have at
  `base` is not compared. A `base` that is not a commit is reported.
- `instructions implement <cycle-id>` orders the target cycle's changes
  that the working tree has not archived: each after its dependencies,
  otherwise by `created`. It prints the list, the loop and the cycle file;
  each change's own instructions come from `instructions implement
  <change-id>`, with `--unattended` when the cycle's had it, run by the
  session, with the cycle's members in build order and that change's place
  among them. A member of a cycle the target holds that is not done gets
  that list whichever way it is run.
- Implement ends by running the archive instructions; `/grove-archive` stays
  for dropping a change and for archiving by hand.
- `grove changes --json` gives paths under the root, with `root` at the
  top, rather than from the repository's top; `state` is worked out when
  the project is read; a field with no value is left out; a change's
  capabilities and problems are not in it.
- `grove version` prints one line, `grove <version>`, without a leading
  `v`, as `grove 0.1.0`; a release build sets it, `go install` of a tag
  takes the tag, and a build with no version recorded says `dev`.
- `install.sh` takes a tag as the release names it (`v0.1.0`), not a bare
  version, and replaces a `grove` already in the folder without asking. It
  says "Installed grove in <dir>." and, when that folder is not on PATH,
  "<dir> is not on your PATH; add it to run grove."; each failure is one
  line starting "grove: ".
