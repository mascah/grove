## T-01 — Write the skill trials

Result: `dev/trials.md` holds the five cases and the five counts the
proposal's "Intended outcome" lists, the `CLAUDE.md` layout table names
it, and `CONTRIBUTING.md` says when to run it.

Files: `dev/trials.md`, `CLAUDE.md`, `CONTRIBUTING.md`

Interfaces: none; it is a test script read by people.

Constraints: each fact is written once (`CLAUDE.md`, "How agents work
here"): the `CLAUDE.md` row says what the file is, and when to run the
trials and where their results go is only in `CONTRIBUTING.md`, "Tests".
Every setup names only files, commands and records the ledger fixture
has (`dev/fixture.md`, `testdata/ledger/`) and steps
`dev/walkthrough.md` has.

Acceptance:
- Each case in `dev/trials.md` gives its setup from `just fixture`, what
  the tester types in Claude Code and in Codex (`/grove-…` and
  `$grove-…`), and what a good run does, as the proposal's table says.
- `dev/trials.md` defines the five counts and says how each is taken in
  each harness, including where the agent's messages are read to count
  their words.
- `dev/trials.md` gives the table a pull request's description carries:
  one row per case and harness, one column per count, and the commit
  `grove version` names.
- `CLAUDE.md`'s layout table lists `dev/trials.md`; `CONTRIBUTING.md`,
  "Tests", says to run the trials with `just trials` in Claude Code and
  in Codex when a pull request changes instruction text in
  `internal/instruct`, and to put the counts in its description.

## T-02 — Record the baseline

Result: the five trials have been run once in Claude Code and once in
Codex on Grove as cycle B leaves it, by `just trials`, and their counts
are in the implementing pull request's description.

Files: `dev/trials.md`, `dev/trials`

Interfaces: none.

Constraints: `dev/trials.md` and `dev/trials/` change here only where
the run shows a setup does not work on the fixture as written; the counts are not
written into the repository.

Acceptance:
- The implementing pull request's description holds the table T-01
  defines, filled in for all five cases in both harnesses, with the
  commit the runs used.
- A setup that did not work as written is corrected in `dev/trials.md`,
  and the description says which and why.

## T-03 — A script that runs the trials

Result: `just trials` runs each case of `dev/trials.md`, or the ones
named, in Claude Code and in Codex, and prints the results table T-01
defines, filled in.

Files: `dev/trials`, `justfile`, `dev/trials.md`

Interfaces: `just trials [case ...]`, a recipe running the script in
`dev/trials/`.

Constraints: rule 17 as amended in PR #72. No Go, so rule 12's line
budget is unchanged. The script writes each setup directly on a fresh
`just fixture`, with no agent session; answers each question from the
case's fixed list, with one fallback reply; and runs the agents with the
build of the checkout it is in.

Acceptance:
- Each case's setup, written by the script, holds what `dev/trials.md`
  says it does, and `grove check` is clean on it.
- One run of a case in each harness leaves its fixture, its transcript
  and the judge's answer, with a reason for each count, in the output
  folder, and prints its rows of the table.
- Edits outside `Files:`/`Delegated:` and words read are computed, not
  judged; the other three counts come from the judge.
