## T-01 — Outside edits counted by grove review

Result: the trials count edits outside `Files:` and `Also changed:` from
`grove review`'s `named nowhere` lines, as the proposal's decisions say.

Files:
- `dev/trials/run.sh`
- `dev/trials.md`

Interfaces: none.

Constraints: rule 17; no Go, so rule 12's budget is unchanged. A branch
that holds no change of its own still counts every file it changes outside
`grove/`, which `grove review` does not list.

Acceptance:
- On case (a)'s setup, a branch from main holding a new change whose task
  names only `ledger/rules.py`, and changing `ledger/cli.py`, counts 1
  edit outside; the script at a5b5674 counts 0.
- On case (c)'s setup, a branch changing `ledger/records.py` and
  `ledger/cli.py` counts 1; with an `Also changed:` line naming
  `ledger/records.py` under T-01, it counts 0.
- `dev/trials.md` names the column "Edits outside Files/Also changed" and
  says how it is taken.

## T-02 — Pinned models and an isolated Codex

Result: every model the trials use is pinned, the agents under test and
the judge at medium effort, and printed with the table, and Codex runs
without the user's Codex setup, as the proposal's decisions say.

Files:
- `dev/trials/run.sh`
- `dev/trials.md`

Interfaces: the `CLAUDE_MODEL` and `CODEX_MODEL` settings of `just
trials`.

Constraints: the user's `auth.json` is linked, never copied.

Acceptance:
- The table's first line names the Grove commit, the Claude Code and Codex
  models under test, the judge's and the waiting check's models, each with
  its effort, and the `claude --version` and `codex --version` output.
- Without `CLAUDE_MODEL` or `CODEX_MODEL`, the Claude Code sessions under
  test run `claude-opus-5-5` at medium effort, the Codex sessions
  `gpt-6.1-sol` at medium effort, and the judge `claude-opus-5-5` at
  medium effort.
- `CLAUDE_MODEL` changes the model of the Claude Code sessions under test
  and not the judge's or the waiting check's.
- With a hook in the user's `~/.codex/hooks.json` that writes a file, and
  a skill in `~/.agents/skills`, a Codex trial runs to its end signed in,
  the hook writes nothing and the skill is not offered.
- `dev/trials.md`'s "Running" names both settings, the medium effort,
  that `gpt-6-astra` is not used for the trials, and how Codex is
  isolated.

## T-03 — A judge that sees results and the repository

Result: the judge reads each tool's result and may read the fixture, as
the proposal's decisions say.

Files:
- `dev/trials/run.sh`
- `dev/trials.md`

Interfaces: none.

Acceptance:
- In both harnesses, `transcript.txt` has a result line after each tool
  line: an error mark when the tool failed, and its first 20 lines, or
  every line for a `grove` command.
- The judge runs with only Read, Grep and Glob, in a copy of the fixture
  at the commit the case's setup left, and its prompt says which state it
  reads and no longer says to count only from the transcript.
- In case (b), the judge's copy holds no change the run proposed.
- `dev/trials.md`'s "How it is taken" says what the judge reads.

## T-04 — Repeated runs and judgings

Result: `RUNS=k` runs each case and harness k times, `JUDGINGS=j` judges
each run j times, and the table prints one row per case and harness, as
the proposal's decisions say.

Files:
- `dev/trials/run.sh`
- `dev/trials.md`

Interfaces: the `RUNS` and `JUDGINGS` settings of `just trials`.

Acceptance:
- `RUNS=2 HARNESSES=claude just trials c` leaves `c-claude-1` and
  `c-claude-2`, the second started after the first ended, and prints one
  row for (c) and Claude Code, each count its median and min–max over the
  two runs.
- Without `RUNS` or `JUDGINGS`, each case and harness runs 3 times and
  each run is judged 3 times.
- `RUNS=1 JUDGINGS=1` prints each count as one number, as at a5b5674.
- Each run's folder holds `judge-1.json` to `judge-<j>.json`, judged from
  its saved `transcript.txt` with no agent session started after its last
  turn; a run's judged count is their median.
- A run whose judgings give a count different values has a line below the
  table naming its case, harness, run, count and each value; with judgings
  that agree there is none.
- `dev/trials.md`'s "Running" and "Results" show both settings, the
  summary row and the judgings' lines.
