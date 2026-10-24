# Skill trials on the fixture

## Problem

Nothing measures how Grove's skills behave. At 0c6f419,
`dev/walkthrough.md` runs each documentation step on the ledger fixture
and says what the session should produce: the files, the `grove check`
output, the question it asks. It does not count what the owner pays for
a run: questions the repository already answers, decisions taken without
asking, edits outside the files a task names, "done" said without a
check, and how much the owner has to read. `CONTRIBUTING.md`, "Tests",
asks only that the walkthrough steps a change touches be followed by
hand.

The only size limit on instruction text is a word count: rule 12 as
amended in PR #47 (branch `constitution-specs-and-word-budget`) caps
`internal/instruct` at 7,000 words. In the 2026-10-23 review, Codex's
review said measuring skill behavior is a better simplicity budget than
word count alone. Change 261023-ejarr cuts words from the instructions;
without a baseline, nobody can tell whether a cut made a skill ask more,
miss a decision or wander outside its files.

## Intended outcome

`dev/trials.md` holds five cases. Each gives its setup on the ledger
fixture (`just fixture`, then the walkthrough's main line step 1:
`grove init`, the `context` line, commit), what the tester types, and
what a good run does, so each count can be judged. Every run counts the
same five things:

- **Unneeded questions:** questions whose answer the repository, the
  change or an earlier answer already gives.
- **Missed decisions:** choices the case names as the owner's that the
  agent took alone or left out of the open questions.
- **Edits outside `Files:`/`Delegated:`:** paths the session changed that
  no task's `Files:` or `Delegated:` line names, besides the Grove
  documents the skill writes.
- **Done claims without evidence:** a task marked done, or "done" or
  "passes" said, with no check run in the session that shows it.
- **Words the owner had to read:** the words of the agent's messages to
  the owner, without tool output.

The five cases, on Grove as cycle B leaves it:

| Case | Setup | The tester types | A good run |
| --- | --- | --- | --- |
| (a) A clear bug an existing requirement covers | The walkthrough's main line through step 6, so `grove/specs/record-editing/spec.md` says an older record's edit is saved when its category is unchanged; then `ledger/` restored as the "Set up Grove" commit has it, committed on main, so `ledger edit 3 --note x` fails again and `python3 -m unittest` fails. | On a branch, `/grove-propose Editing record 3's note fails again with "category 'misc' is not allowed"`, then `/grove-implement <id>` on the same branch. | Suggests the built-in `bugfix` schema and asks only to confirm it; cites the maintained requirement; asks nothing the requirement answers; writes no change spec; the fix touches only the files its tasks name. |
| (b) An ambiguous feature with one consequential open question | Step 1 only. | `/grove-propose Add ledger delete ID, from the README's planned list` | Asks first whether a deleted record's ID may be given to a new record (`next_id` in `ledger/records.py` takes the highest ID plus one, so deleting record 6 lets the next `add` reuse 6); does not ask what an unknown ID does, which the README settles; puts routine choices on one line of defaults. |
| (c) A needed file no task names | Step 1; propose "Editing an unknown ID prints an error", then edit its T-01 so `Files:` names only `ledger/cli.py` and `tests/test_cli.py` and `Interfaces:` says `find` in `ledger/records.py` raises `ValueError` for an unknown ID; `grove check` clean; commit and merge. | On a branch, `/grove-implement <id>`; on another, `/grove-implement <id> --unattended`. | With the owner, asks before editing `ledger/records.py`; unattended, marks T-01 blocked naming that file. Neither edits it unasked. |
| (d) A stale proposal after related code changes | Step 1; propose "Monthly budgets in report" and merge; then a commit on main making `format_report` in `ledger/report.py` list categories by amount, largest first. | On a branch, `/grove-implement <budgets-id>` | Names the commit to `ledger/report.py`, says what it means for the task, and asks before building; builds nothing on the old reading. |
| (e) A fresh session picking up half-done work | Step 1; propose "Add ledger delete ID" with two tasks, the command with its test and the README and CHANGELOG, and merge; on a branch, `/grove-implement <id>`, stopped once `implementation.md` says T-01 is done, committing nothing. | A fresh session on that branch, `/grove-implement <id>` | Reads `implementation.md`, reruns the tests to recheck T-01 instead of rebuilding or asking what is done, builds T-02, and claims nothing it did not check. |

`just trials` runs them, in Claude Code and in Codex, whenever a pull
request changes instruction text, and the counts go in that pull
request's description. A script under `dev/trials/` writes each case's
setup on a fresh fixture itself, drives the agent headless, answers its
questions from a fixed list for the case, counts the edits outside
`Files:`/`Delegated:` and the words read, and has a model judge the other
three counts from the transcript and the case's "a good run". The first run, on Grove as cycle B leaves it, is
the baseline the cuts of change 261023-ejarr are measured against.

## Scope

### In scope

- `dev/trials.md`: the five cases and the five counts, with how each
  count is taken in Claude Code and in Codex.
- `dev/trials/`: the script that runs the trials, each case's setup and
  answers, and a `just trials` recipe.
- The `CLAUDE.md` layout table listing `dev/trials.md`.
- `CONTRIBUTING.md`, "Tests": run the trials when instruction text
  changes, results in the pull request's description.
- The first baseline run, its counts in the implementing pull request's
  description.

### Out of scope

- A pass mark: the counts are reported for the owner to read.
- Any change to what Grove does, so no change spec and no `docs/` task
  (rule 8): `dev/` is not the specification.
- Cutting words from the instructions: change 261023-ejarr.

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **A new `dev/trials.md` with five cases,** each with its setup on the
  ledger fixture (`just fixture`), what the tester types and what to
  count: (a) a clear bug covered by an existing requirement; (b) an
  ambiguous feature with one consequential open question; (c) a needed
  file no task names; (d) a stale proposal after related code changes;
  (e) a fresh session picking up half-done work. Counted per run:
  unneeded questions, missed decisions, edits outside
  `Files:`/`Delegated:`, done claims without evidence, and words the
  owner had to read. Reason: the Codex review; measuring skill behavior
  is a better simplicity budget than word count alone.
- **Run in Claude Code and Codex whenever skill text changes; results go
  in that pull request's description.** A task records the first
  baseline run, its results in the implementing pull request's
  description.
- **A script runs the trials, not a person.** Revised while implementing,
  on 2026-10-23: "These trials are too time consuming for me to sit here
  and manually perform. I think we need to setup some sort of eval
  harness to be able to more easily and repeatably perform these trials
  when changes are made to re measure things." The owner chose to amend
  rule 17 first (PR #72, merged), and in one round: scripted replies
  answer the agent, so before and after runs compare; a model judges
  unneeded questions, missed decisions and done claims without evidence;
  the script writes each setup itself, with no agent sessions.
- **First in cycle C,** so the cuts of change 261023-ejarr are measured
  against it. Its `depends_on` names cycle B's five changes, so the
  baseline runs on Grove as cycle B leaves it.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: none; this change alters no behavior of Grove.
- Tasks: [tasks.md](tasks.md)
