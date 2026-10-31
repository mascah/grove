# Implementation: Skill trials that measure what they claim

<!-- Progress on change 261028-3ntsa, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Outside edits counted by grove review

Status: done

`dev/trials/run.sh`: `allowed` and `named` (the awk copy of Grove's rule)
and the per-run `snapshot` are gone. After the case's last run the script
commits what it left (`save`), and `outside` lists, on each branch the
case used, each file outside `grove/` the branch changes against main that
`grove review` marks neither `planned` nor `agreed`, as `branch: path` in
`outside.txt`. A branch whose review lists nothing, as one holding no
change, counts every such file. The column is "Edits outside Files/Also
changed"; `dev/trials.md`'s counts table says how it is taken.

Checks, on `SETUP_ONLY=1 HARNESSES=claude dev/trials/run.sh a c` with the
end state written by hand and `outside` run in each fixture:
- (a): on `fix` from main, a new change whose T-01 names only
  `ledger/rules.py`, and `ledger/cli.py` changed: `fix: ledger/cli.py`,
  1. The awk of a5b5674 on the same state lists `ledger/cli.py`,
  `ledger/rules.py`, `tests/test_cli.py` as allowed, so 0.
- (c): on `unknown-id-impl`, `ledger/records.py` and `ledger/cli.py`
  changed: `unknown-id-impl: ledger/records.py`, 1. With `Also changed:
  \`ledger/records.py\`` under T-01 in `implementation.md`: none, 0
  (`grove review` marks it `agreed 261015-c0000:T-01`).

### T-02 — Pinned models and an isolated Codex

Status: done

`dev/trials/run.sh`: `CLAUDE_MODEL` (default `claude-opus-5-5`) and
`CODEX_MODEL` (default `gpt-6.1-sol`) set the agents under test, both run
with medium effort (`--effort medium`, `-c model_reasoning_effort="medium"`);
the judge is fixed at `claude-opus-5-5` with `--effort medium` and the
waiting check at `claude-sonnet-5-5`. The waiting check gets no effort
setting, since the proposal fixes effort only for the agents and the judge,
and the first line says "at its default effort". Each Codex trial runs with
`HOME=<trial folder>` and `CODEX_HOME=<trial folder>/codex`, holding only a
symlink to the user's `auth.json` (`${CODEX_HOME:-$HOME/.codex}`). The
first line of the table names the Grove commit, every model with its
effort, and `claude --version` and `codex --version`. `dev/trials.md`'s
"Running" says so.

Checks so far, on this machine (2026-10-31):
- The user's `~/.codex/config.toml` enables hooks (`[features] hooks =
  true`) and trusts those of `~/.codex/hooks.json`, whose SessionStart
  hook appends to Crew's `journal.jsonl`; `~/.agents/skills` holds
  `hindsight-coding-agent`.
- With `CODEX_HOME` a folder holding only the `auth.json` link and the real
  `HOME`, `codex exec -m gpt-6.1-sol` asked to list its skills lists
  `hindsight-coding-agent` among them. With `HOME` and `CODEX_HOME` set as
  the script sets them, it lists only Codex's bundled `imagegen`,
  `openai-docs`, `skill-creator`, `skill-installer`, signed in, and
  `journal.jsonl`'s size and mtime are unchanged.

End to end, on 2026-10-31 at 4edf29a: `RUNS=2 HARNESSES=claude just trials c`
and `RUNS=1 JUDGINGS=1 HARNESSES=codex just trials e`, both exit 0, with
empty `stderr.txt`. Their output, folders in this machine's scratch space:

    Trials at 4edf29a (…/e2e-c): Claude Code claude-opus-5-5 and Codex gpt-6.1-sol at medium effort, judge claude-opus-5-5 at medium effort, waiting check claude-sonnet-5-5 at its default effort; 2.1.316 (Claude Code), codex-cli 0.166.0

    | Case | Harness | Unneeded questions | Missed decisions | Edits outside Files/Also changed | Done claims without evidence | Listed documents not read | Words read |
    | --- | --- | --- | --- | --- | --- | --- | --- |
    | (c) | Claude Code | 0 (0–0) | 0 (0–0) | 0 (0–0) | 0 (0–0) | 0 (0–0) | 734.5 (714–755) |

    | (e) | Codex | 0 | 0 | 0 | 0 | 0 | 134 |

- Every Claude Code session under test reports `claude-opus-5-5` in its
  `init` event; the Codex trial's sessions, in its own `codex/sessions`,
  `gpt-6.1-sol` with `reasoning_effort` `medium`. The judge and the waiting
  check use `judge_model` and `waiting_model`, set apart from
  `CLAUDE_MODEL` (read from the script; not run with another
  `CLAUDE_MODEL`).
- The Codex trial ran to its end signed in. Crew's `journal.jsonl` holds
  no `codex` entry whose `cwd` is under the trial folder: the user's
  SessionStart hook did not run.

### T-03 — A judge that sees results and the repository

Status: done

`dev/trials/run.sh`: the jq function `result` writes, after each tool
line, `RESULT:` or `RESULT (error):` and the result's first 20 lines,
indented, with `… N more lines` when cut, or every line when the command
runs `grove`. Claude Code's results come from the `user` events'
`tool_result` blocks, matched to their call by `tool_use_id`; Codex's from
`aggregated_output` and `exit_code`, and a `file_change`'s `status`. After
setup the script keeps `setup_rev`; before judging it extracts that commit
with `git archive` into `<trial>/setup/` and runs the judge there with
`--tools Read,Grep,Glob`. Its prompt names the RESULT lines, says the
folder holds the fixture as the setup left it, before the run, and lets it
check what the repository already answered; "Count only from the
transcript" is gone. The files under `grove/changes/` at the run's end are
listed from the run's fixture before judging. `dev/trials.md`'s "The
counts" says what the judge reads.

Not in the judge's copy: case (e)'s setup leaves T-01's work and
`implementation.md` uncommitted, so its copy, at the setup's commit, does
not hold them; `dev/trials.md` says so.

Checks so far:
- The Claude Code program on a `claude -p` session running a blocked `ls`
  and `echo hi`: `RESULT (error):` with the block message, then `RESULT:`
  `hi`. The Codex program on a `codex exec` session running `ls
  /nonexistent` (exit 1) and `echo hi`: the same. On 25-line outputs:
  `seq 25` shows lines 1 to 20 and `… 5 more lines`; `grove instructions
  x` all 25.
- End to end (T-02's runs): (c) run 1's transcript has 29 TOOL lines and
  29 RESULT lines, (e)'s Codex transcript 12 and 12. Each run's `setup/`
  holds 261015-c0000's or 261015-e0000's `.grove.yaml`, `proposal.md`,
  `specs/`, `tasks.md`, and no `implementation.md` the run wrote.
- (b)'s setup commit (`SETUP_ONLY=1`), the judge's copy: `git archive
  HEAD` lists no file under `grove/changes/`; the run's change is
  committed later, on `delete`.

### T-04 — Repeated runs and judgings

Status: done

`dev/trials/run.sh`: `RUNS` (default 3) and `JUDGINGS` (default 3). Each
case and harness runs its k runs one after another, in folders
`<case>-<harness>-<run>`, the cases and harnesses at once; each run is
judged j times from its saved `transcript.txt`, into `judge-1.json` to
`judge-<j>.json`, and leaves `outside.txt`, `words.txt`, and `failed.txt`
when it failed. A Python summary prints one row per case and harness: a
run's judged count is the median of its judgings, and a row's count the
median over its runs, with `(min–max)` when k > 1. Below the table, a line
for each failed run and each run whose judgings disagree on a count.
`dev/trials.md`'s "Running" and "Results" show both settings, the summary
row and the lines.

Checks:
- `RUNS=2 HARNESSES=claude just trials c` (T-02's run): `c-claude-1` and
  `c-claude-2`; run 1's `words.txt` written at 1793480003, run 2's
  `setup-check.txt` at 1793480004; one row for (c), Claude Code, each count
  as `median (min–max)`; each folder holds `judge-1.json` to
  `judge-3.json` (the default `JUDGINGS`), the judgings agreed and no
  line followed.
- `RUNS=1 JUDGINGS=1 HARNESSES=codex just trials e`: one number per count,
  as at a5b5674, and only `judge-1.json`.
- The summary on hand-written folders, k = 2, j = 3: judgings 1, 2, 1 give
  the line `(c) Claude Code run 1: Unneeded questions judged 1, 2, 1` and
  run 1 the count 1; with run 2's 4, 4, the row shows `2.5 (1–4)`; an
  unreadable `judge-3.json` is left out; a run with only `failed.txt`
  gives `(c) Codex run 2 failed: usage limit`.
- Defaults of 3 runs and 3 judgings: `runs_each=${RUNS:-3}`,
  `judgings=${JUDGINGS:-3}`; a full default `just trials` was not run.

## Scenarios

The change has no change specs, so no scenarios: it alters no behavior of
Grove (its proposal's "Related documents").

`just check`: passed (4024 Go lines, 1681 test lines, 8023 instruction
words, all unchanged).

## Reconciliation

None: the change has no change specs.
