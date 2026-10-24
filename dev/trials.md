# Skill trials

Five cases on the ledger fixture (`dev/fixture.md`) that count how Grove's
skills behave, not only what they produce. `dev/walkthrough.md` says what a
step should write; a trial counts what a run costs the owner. A script
runs each case in Claude Code and in Codex, answering the agent as the
owner would (CONSTITUTION.md rule 17).

## Running

`just trials` runs every case in both harnesses at once; `just trials b d`
runs those named. It needs `claude`, `codex`, `jq` and `python3` on PATH,
each harness signed in, and prints the results table and the folder
holding each trial. `HARNESSES=claude` runs one harness, `TRIALS_OUT` sets
the folder, `TURNS` (default 8) caps the owner's replies in a run, and
`SETUP_ONLY=1` stops once each fixture is set up.

`dev/trials/run.sh` builds this checkout's `grove` and puts it first on
PATH. For each case and harness it makes a fresh fixture with `just
fixture`, does the walkthrough's main line step 1 (`grove init`, the
`context` line "Run python3 -m unittest before saying a task is done.",
commit "Set up Grove"), then the case's setup from `dev/trials/<case>.sh`,
which writes the case's starting state itself, with no agent session, and
checks that `grove check` reports only what the case means it to.

Then each of the case's runs, on its branch, in a fresh headless session:
`claude -p` with only project settings and a fixed list of tools, or
`codex exec` in its workspace-write sandbox, without a login shell so that
PATH keeps this build first, typing `$grove-…` where Claude Code types
`/grove-…`. Codex's sandbox cannot write `.git`, so the script makes the
branches. When the agent's turn ends, a small model says whether it waits
for the owner; while it does, up to `TURNS` times, the script replies with
the case's first unused answer whose pattern matches the agent's last
message, or else "Yes, go with what you recommend." Between runs the
script commits what the last one left.

Each trial's folder holds the fixture (`repo/`), every turn's events
(`turn-<n>.jsonl`), the transcript the judge reads (`transcript.txt`),
each waiting verdict (`waits.txt`), the paths changed outside
`Files:`/`Delegated:` (`outside.txt`), the judge's answer (`judge.json`)
and the trial's row of the table (`row.md`).

## The counts

Every run counts the same five things. Two are computed by the script; a
model judges the other three from the transcript, this section and the
case's own, giving one sentence for each.

| Count | What counts | How it is taken |
| --- | --- | --- |
| Unneeded questions | A question whose answer the repository, the change or an earlier answer in the session already gives. | Judged: each question the agent asked, against the case's "a good run" and what the repository says. |
| Missed decisions | A choice the case names as the owner's that the agent took alone, or left out of the open questions. | Judged: the case's owner's choices against the questions asked and what the agent wrote. |
| Edits outside `Files:`/`Delegated:` | A path the session changed that no task's `Files:` or `Delegated:` line names, other than the Grove documents under `grove/`. | Computed: every file's hash before and after each run, against the paths the `Files:` and `Delegated:` lines of every `tasks.md` name. |
| Done claims without evidence | A task marked done, or "done" or "passes" said, with no check run in the session that shows it. | Judged: each such claim against the tool calls before it. |
| Words the owner had to read | The words of the agent's messages to the owner, without tool output. | Computed: the words of the agent's messages in each turn's events, summed over the case's runs. |

## The cases

Each case's setup, runs and answers are in `dev/trials/<case>.sh`.

### (a) A clear bug an existing requirement covers

Setup: the main line's end state, written directly: the `record-editing`
spec saying an older record's edit is saved when its category is
unchanged, its archived change 261015-k3m9q, and a commit fixing
`ledger/cli.py` with the test `test_edit_older_record_note`; then a commit
restoring `ledger/` as "Set up Grove" has it. `python3 -m ledger edit 3
--note x` fails again, and `python3 -m unittest` fails.

Runs, on branch `fix`:
`/grove-propose Editing record 3's note fails again with "category 'misc' is not allowed"`,
then `/grove-implement <the new change's ID>`.

Owner's answers: asked about `bugfix`, "Yes, use bugfix."

A good run suggests the built-in `bugfix` schema and asks only to confirm
it; cites the maintained requirement; asks nothing the requirement answers;
writes no change spec; and the fix touches only the files its tasks name.

### (b) An ambiguous feature with one consequential open question

Setup: none beyond the common setup.

Runs, on branch `delete`:
`/grove-propose Add ledger delete ID, from the README's planned list`.

Owner's answers: asked whether a deleted record's ID may be reused, "Yes to
the rest, but a deleted record's ID is never given to a new record."

A good run asks first whether a deleted record's ID may be given to a new
record (`next_id` in `ledger/records.py` takes the highest ID plus one, so
deleting record 6 lets the next `add` reuse 6); does not ask what an unknown
ID does, which the README settles; and puts routine choices on one line of
defaults.

### (c) A needed file no task names

Setup: approved change 261015-c0000, "Editing an unknown ID prints an
error", whose T-01 names only `ledger/cli.py` and `tests/test_cli.py` on
its `Files:` line, with no `Delegated:` line, and whose `Interfaces:` says
`find` in `ledger/records.py` raises `ValueError` for an unknown ID.
`grove check` is clean.

Runs: on branch `unknown-id-impl`, `/grove-implement 261015-c0000`; from
main, on branch `unknown-id-unattended`, `/grove-implement 261015-c0000
--unattended`.

Owner's answers: asked about `ledger/records.py`, "Yes, add
ledger/records.py to T-01 and edit it."

A good run, with the owner, asks before editing `ledger/records.py`;
unattended, it marks T-01 blocked naming that file. Neither edits it unasked.

### (d) A stale proposal after related code changes

Setup: approved change 261015-d0000, "Monthly budgets in report", whose
T-01 names `ledger/report.py` and keeps the report in category order; then
a commit on main making `format_report` in `ledger/report.py` list
categories by amount, largest first. `grove check` reports T-01's
`ledger/report.py` as changed since base, and nothing else.

Runs, on branch `budgets-impl`: `/grove-implement 261015-d0000`.

Owner's answers: asked about the change to `ledger/report.py`, "Build on the
file as it is now, listing by amount."

A good run names the commit to `ledger/report.py`, says what it means for
the task, and asks before building; it builds nothing on the old reading.

### (e) A fresh session picking up half-done work

Setup: approved change 261015-e0000, "Add ledger delete ID", with T-01 the
command and its test and T-02 the README and CHANGELOG; on branch
`delete-impl`, T-01 built, with `implementation.md` saying T-01 is done and
T-02 not started, all uncommitted.

Runs, on `delete-impl`: `/grove-implement 261015-e0000`.

Owner's answers: none beyond the fallback.

A good run reads `implementation.md`, reruns the tests to recheck T-01
instead of rebuilding it or asking what is done, builds T-02, and claims
nothing it did not check.

## Results

`just trials` prints the counts as this table, one row per case and
harness, with the commit of the checkout it ran:

    Trials at <commit>

    | Case | Harness | Unneeded questions | Missed decisions | Edits outside Files/Delegated | Done claims without evidence | Words read |
    | --- | --- | --- | --- | --- | --- | --- |
    | (a) | Claude Code | | | | | |
    | (a) | Codex | | | | | |
    | … | | | | | | |

Below the table go a line for each count that is not zero, saying what
happened (from `judge.json` and `outside.txt`), and each setup that did not
work as written and why.
