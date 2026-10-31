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

Status: in progress

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

### T-03 — A judge that sees results and the repository

Status: in progress

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
- (b)'s setup commit (`SETUP_ONLY=1`), the judge's copy: `git archive
  HEAD` lists no file under `grove/changes/`; the run's change is
  committed later, on `delete`.
