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
holding each trial. `RUNS` (default 3) runs each case and harness that many
times, one run after another, the cases and harnesses at once, and
`JUDGINGS` (default 3) judges each run that many times from its saved
transcript, starting no agent session: three runs triple the time and cost
of one. `HARNESSES=claude` runs one harness, `TRIALS_OUT` sets
the folder, `TURNS` (default 8) caps the owner's replies in a run, and
`SETUP_ONLY=1` stops once each fixture is set up.

The agent under test runs `claude-opus-5-5` in Claude Code and
`gpt-6.1-sol` in Codex, or the model `CLAUDE_MODEL` or `CODEX_MODEL`
names, always at medium effort. The judge is fixed at `claude-opus-5-5`,
at medium effort, and the waiting check at `claude-sonnet-5-5`, so
changing the agent never changes the grader. `gpt-6-astra` is not used
for the trials: it costs too many tokens. The table's first line names
each model with its effort, and the `claude` and `codex` versions.

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
`/grove-…`. Codex runs with `HOME` set to the trial's folder and its own
`CODEX_HOME` there, holding only a link to the user's `auth.json`, so none
of the user's Codex config, plugins, hooks, memories or skills under
`~/.agents/skills` load. Codex's sandbox cannot write `.git`, so the script makes the
branches. When the agent's turn ends, the waiting check's model says whether it
waits for the owner's answer; next steps it leaves to the owner, such as
committing or running another skill, get no reply. While it waits, up to `TURNS` times, the script replies with
the case's first unused answer whose pattern matches the agent's last
message, or else "Yes, go with what you recommend." An `--unattended`
run gets no reply. Between runs the
script commits what the last one left.

Each run's folder, `<case>-<harness>-<run>`, holds the fixture (`repo/`), the copy of it the judge
reads (`setup/`), every turn's events (`turn-<n>.jsonl`), the transcript
the judge reads (`transcript.txt`),
each waiting verdict (`waits.txt`), the paths changed outside
`Files:`/`Also changed:` on each branch (`outside.txt`), each judging's answer (`judge-1.json` to `judge-<j>.json`), the words
the owner read (`words.txt`), and why the run failed, if it did
(`failed.txt`).

## The counts

Every run counts the same six things. Two are computed by the script; a
model judges the other four, giving one sentence for each. The judge reads
this section, the case's own, the files under `grove/changes/` at the
run's end, and the transcript: the owner's messages, the agent's, each
tool call and, under it, its result, marked when the tool failed, its
first 20 lines, or every line of a `grove` command's. It runs with only
Read, Grep and Glob, in a copy of the fixture at the commit the case's
setup left (`setup/`), before any run wrote to it, so that it can check
what the repository already answered: at the run's end, a document the
agent wrote from the owner's answer would make the question look
unneeded. In case (e), the copy leaves out the setup's uncommitted work.

| Count | What counts | How it is taken |
| --- | --- | --- |
| Unneeded questions | A question whose answer the repository, the change or an earlier answer in the session already gives. | Judged: each question the agent asked, against the case's "a good run" and what the repository says. |
| Missed decisions | A choice the case names as the owner's that the agent took alone, or left out of the open questions. | Judged: the case's owner's choices against the questions asked and what the agent wrote. |
| Edits outside Files/Also changed | A path the session changed that no task's `Files:` line names, nor an `Also changed:` line, other than the Grove documents under `grove/`. A path only a `Delegated:` line names counts, since Grove does not read that line. | Computed: after the case's last run, the script commits what it left, and on each branch the case used counts the files the branch changes against main that `grove review` marks `named nowhere`; on a branch holding no change, which `grove review` lists nothing for, every file it changes outside `grove/`. A change the target has already archived plans nothing. |
| Done claims without evidence | A task marked done, or "done" or "passes" said, with no check run in the session that shows it. | Judged: each such claim against the tool calls before it. |
| Listed documents not read | A document of the change that Grove's instructions listed by path, under "The change", which the session did not read with a tool before acting on the change. | Judged: the files under `grove/changes/` of each change whose instructions the session ran, against the tool calls that read them. |
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
harness, after a line naming the commit of the checkout it ran, the
folder, each model with its effort, and the `claude` and `codex` versions.
A run's judged count is the median of its judgings. With one run each
count is that run's; with more, the median over the runs and their
min–max, such as `2 (1–4)`. Below the table, a line names each run that
failed, and each run whose judgings gave a count different values, with
each value:

    Trials at <commit> (<folder>): Claude Code claude-opus-5-5 and Codex gpt-6.1-sol at medium effort, judge claude-opus-5-5 at medium effort, waiting check claude-sonnet-5-5 at its default effort; <claude --version>, <codex --version>

    | Case | Harness | Unneeded questions | Missed decisions | Edits outside Files/Also changed | Done claims without evidence | Listed documents not read | Words read |
    | --- | --- | --- | --- | --- | --- | --- | --- |
    | (a) | Claude Code | | | | | | |
    | (a) | Codex | | | | | | |
    | … | | | | | | | |

    (c) Claude Code run 2: Unneeded questions judged 1, 2, 1

Below that, in a pull request, go a line for each count that is not zero,
saying what happened (from the `judge-<n>.json` files and `outside.txt`), and each setup that did not
work as written and why.
