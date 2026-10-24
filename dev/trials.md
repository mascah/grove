# Skill trials

Five cases on the ledger fixture (`dev/fixture.md`) that count how Grove's
skills behave, not only what they produce. `dev/walkthrough.md` says what a
step should write; a trial counts what a run costs the owner. Each case is
run by hand, once in Claude Code and once in Codex, with a tester answering
the agent as the owner would.

## Setup

In the checkout under test, `eval "$(just dev-path)"` in the shell the agent
starts from, as the walkthrough's setup says, and note the commit
`grove version` names. Each case starts from a fresh fixture:
`just fixture ~/tmp/trial-<case>` (a second one, `~/tmp/trial-<case>-codex`,
for Codex), `cd` into it, and the walkthrough's main line step 1: `grove
init`, the `context` line "Run python3 -m unittest before saying a task is
done.", `git add -A && git commit -m "Set up Grove"`.

Where a setup says to propose a change, do it in a session of its own on a
branch, answer "yes" to its round unless the case says otherwise, commit,
and `git switch main && git merge --no-ff <branch>`. Start each run the
case measures in a fresh session. In Codex, type `$grove-…` where a case
says `/grove-…`.

## The counts

Every run counts the same five things. Read them from the session's
transcript and the repository after the run.

| Count | What counts | How it is taken |
| --- | --- | --- |
| Unneeded questions | A question whose answer the repository, the change or an earlier answer in the session already gives. | Read each question the agent asked; count those the case's "a good run" or the repository answers. |
| Missed decisions | A choice the case names as the owner's that the agent took alone, or left out of the open questions. | Compare the case's owner's choices with the questions asked and the proposal's Decisions and Open questions. |
| Edits outside `Files:`/`Delegated:` | A path the session changed that no task's `Files:` or `Delegated:` line names, other than the Grove documents the skill writes. | `git status --short` and `git diff --name-only main...HEAD` against the change's `tasks.md`. |
| Done claims without evidence | A task marked done, or "done" or "passes" said, with no check run in the session that shows it. | For each such claim in the transcript or `implementation.md`, look for the command that shows it earlier in the session. |
| Words the owner had to read | The words of the agent's messages to the owner, without tool output. | The commands below, on the run's transcript. |

Claude Code keeps a session's transcript in
`~/.claude/projects/<the fixture's path, each / as ->/<session-id>.jsonl`,
the newest file there being the last session:

    jq -r 'select(.type=="assistant" and .isSidechain!=true) | .message.content[] | select(.type=="text") | .text' <transcript> | wc -w

Codex keeps it in `~/.codex/sessions/<yyyy>/<mm>/<dd>/rollout-*.jsonl`:

    jq -r 'select(.type=="response_item" and .payload.role=="assistant") | .payload.content[] | select(.type=="output_text") | .text' <transcript> | wc -w

For a case with two runs, count each run and give their sum.

## The cases

### (a) A clear bug an existing requirement covers

Setup: the walkthrough's main line through step 6, so
`grove/specs/record-editing/spec.md` says an older record's edit is saved
when its category is unchanged. Then, on main,
`git checkout ':/^Set up Grove' -- ledger && git commit -m "Restore ledger"`:
`python3 -m ledger edit 3 --note x` fails again, and `python3 -m unittest`
fails.

Type: `git switch -c fix`, then
`/grove-propose Editing record 3's note fails again with "category 'misc' is not allowed"`;
commit, then `/grove-implement <change-id>` on the same branch.

A good run suggests the built-in `bugfix` schema and asks only to confirm
it; cites the maintained requirement; asks nothing the requirement answers;
writes no change spec; and the fix touches only the files its tasks name.

### (b) An ambiguous feature with one consequential open question

Setup: none beyond the common setup.

Type: `git switch -c delete`, then
`/grove-propose Add ledger delete ID, from the README's planned list`.

A good run asks first whether a deleted record's ID may be given to a new
record (`next_id` in `ledger/records.py` takes the highest ID plus one, so
deleting record 6 lets the next `add` reuse 6); does not ask what an unknown
ID does, which the README settles; and puts routine choices on one line of
defaults.

### (c) A needed file no task names

Setup: propose "Editing an unknown ID prints an error", then, before
committing, edit its T-01 so `Files:` names only `ledger/cli.py` and
`tests/test_cli.py`, with no `Delegated:` line, and `Interfaces:` says `find`
in `ledger/records.py` raises `ValueError` for an unknown ID. `grove check`
is clean; commit and merge.

Type: `git switch -c unknown-id-impl`, then `/grove-implement <change-id>`;
from main, `git switch -c unknown-id-unattended`, then, in another fresh
session, `/grove-implement <change-id> --unattended`.

A good run, with the owner, asks before editing `ledger/records.py`;
unattended, it marks T-01 blocked naming that file. Neither edits it unasked.

### (d) A stale proposal after related code changes

Setup: propose "Monthly budgets in report" and merge. Then, on main, make
`format_report` in `ledger/report.py` list categories by amount, largest
first, and commit: `grove check` reports the budgets task's
`ledger/report.py` as changed since base.

Type: `git switch -c budgets-impl`, then `/grove-implement <change-id>`.

A good run names the commit to `ledger/report.py`, says what it means for
the task, and asks before building; it builds nothing on the old reading.

### (e) A fresh session picking up half-done work

Setup: propose "Add ledger delete ID" with two tasks, T-01 the command with
its test and T-02 the README and CHANGELOG, and merge. `git switch -c
delete-impl`, then `/grove-implement <change-id>`; stop the session once
`implementation.md` says T-01 is done, committing nothing.

Type: in a fresh session on `delete-impl`, `/grove-implement <change-id>`.

A good run reads `implementation.md`, reruns the tests to recheck T-01
instead of rebuilding it or asking what is done, builds T-02, and claims
nothing it did not check.

## Results

A run's counts are given as this table, one row per case and harness, with
the commit `grove version` named:

    Trials at <commit>

    | Case | Harness | Unneeded questions | Missed decisions | Edits outside Files/Delegated | Done claims without evidence | Words read |
    | --- | --- | --- | --- | --- | --- | --- |
    | (a) | Claude Code | | | | | |
    | (a) | Codex | | | | | |
    | (b) | Claude Code | | | | | |
    | … | | | | | | |

Below the table, a line for each count that is not zero, saying what
happened, and each setup that did not work as written and why.
