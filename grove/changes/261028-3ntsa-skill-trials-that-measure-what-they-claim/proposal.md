# Skill trials that measure what they claim

## Problem

The skill trials (`dev/trials.md`, `dev/trials/`, from change 261023-kh3gg)
print counts that cannot tell a better skill from a worse one. A review of
them on 2026-10-28, at a5b5674, found:

- **Edits outside `Files:` miss case (a)'s failure.** `allowed` in
  `dev/trials/run.sh` (line 116) reads every `tasks.md` under
  `grove/changes/`, the archive included. Case (a)'s setup archives change
  261015-k3m9q, whose T-01 names `ledger/cli.py` and `tests/test_cli.py`,
  so on a fresh (a) setup both paths are allowed before any agent runs
  (checked with `SETUP_ONLY=1`). A fix that touches files its own tasks do
  not name counts 0. `grove review` already marks each changed file
  `planned`, `agreed` or `named nowhere` (`internal/cli/review.go`, `mark`)
  and leaves out changes the target has already archived; the script's awk
  is a second copy of that rule.
- **Nothing is pinned or recorded but the Grove commit.** `turn` and
  `judge` pass no `--model`, so the agent and the judge run on the
  account's default model (`claude-opus-5-5` on 2026-10-28, from a session
  with the script's flags), and `waiting` uses the `sonnet` alias. `codex
  exec` reads the user's `~/.codex/config.toml` (on this machine
  `gpt-6-astra`, high effort, and the enabled plugins `mascah-skills`,
  `browser`, `documents`, `spreadsheets`, `presentations`), its
  `hooks.json` (Hindsight memory and Crew hooks) and `~/.agents/skills`.
  So the Codex column measures Grove plus one user's setup, and a run on
  another day or machine is not comparable. Claude Code is isolated
  already: `--setting-sources project,local` loaded none of the user's
  plugins.
- **The judge cannot see what it judges.** For Claude Code, the
  transcript holds tool calls without their results (`turn`, the
  `assistant` events only), so "done claims without evidence" cannot tell
  a passing test run from a failing one; Codex lines carry the exit code.
  "Unneeded questions" depends on what the repository says, but the judge
  runs with `--tools ""` and is told to "Count only from the transcript".
  "Listed documents not read" depends on what `grove instructions` listed,
  which is tool output and so not in the transcript.
- **One run per case and harness.** PR #73 (8705c42) found "the judge
  counted the same question differently between runs" and, over three
  runs, unneeded questions of 34 before and 34 after; the owner then
  called the trials "a record here, not a gate". The script has no way to
  repeat a run other than running it again by hand, nor to judge one saved
  transcript twice, so a count that moves cannot be put down to the agent
  or to the judge. `CONTRIBUTING.md` asks for the table in every pull
  request that changes instruction text, and that table is one run.

## Intended outcome

`just trials` prints the same table, with these differences:

- Edits outside are counted from `grove review`'s `named nowhere` lines on
  each branch a case used, so case (a) counts an edit its new change does
  not name. The column is "Edits outside Files/Also changed".
- The agent's model in each harness is pinned at medium effort,
  overridable by `CLAUDE_MODEL` and `CODEX_MODEL`; the judge and the
  waiting check are pinned in the script. The table's first line names
  every model with its effort and the `claude` and `codex` versions.
- Each Codex trial runs with its own `CODEX_HOME` and `HOME`, so none of
  the user's Codex config, plugins, hooks, memories or user skills load.
- The transcript the judge reads shows each tool's result in both
  harnesses, and the judge can read the fixture as the case's setup left
  it.
- `RUNS=k` (default 3) runs each case and harness k times, and
  `JUDGINGS=j` (default 3) judges each run j times from its saved
  transcript, starting no agent session. The table keeps one row per case
  and harness, each count its median and min–max over the runs when k > 1.
  Below it, each run whose judgings disagree on a count gets a line.

## Scope

### In scope

- `dev/trials/run.sh` and `dev/trials.md`.

### Out of scope

- Whether a run met its case's "a good run": change 261029-wtbc0 (the
  next change on this branch), which depends on this one.
- Any change to what Grove does, so no change spec and no `docs/` task
  (rule 8): `dev/` is not the specification.
- Replacing the scripted answers with a model playing the owner: rule 17
  says "answering them from a script".
- New cases, and `CONTRIBUTING.md`'s rule to run the trials when a pull
  request changes instruction text.

## Decisions

The owner accepted the scope, its split into two changes and no
`depends_on` for this one ("yes", 2026-10-28), then every item of the round
and its defaults ("yes", 2026-10-28). After a review comment on PR #106,
the owner chose three runs by default ("yes default to 3", 2026-10-30) and
accepted judgings, the summary row and the judge reading the setup's state.
On 2026-10-31 the owner said "we absolutely cannot use gpt 6 astra at high
effort for the trials", having asked for "gpt 6.1 sol medium or opus 5.5
medium for the trials to converse tokens", and accepted the revised models
and efforts below ("yes", 2026-10-31):

- **Outside edits are counted by `grove review`.** After a case's last run
  the script commits what it left and counts the `named nowhere` lines of
  `grove review` on each branch the case used. Reason: it fixes case (a)'s
  leak and drops a second copy of Grove's rule. A path on a `Delegated:`
  line now counts as outside, since Grove does not read that line, and the
  column is "Edits outside Files/Also changed".
- **Codex runs isolated:** a fresh `CODEX_HOME` per trial holding only a
  link to the user's `auth.json`, and `HOME` set to the trial folder.
  Reason: `--ignore-user-config` skips `config.toml` but not, as far as
  the review found, `hooks.json` or user skills.
- **Pinned models, printed with the table:** the agent under test uses
  `CLAUDE_MODEL` (default `claude-opus-5-5`) and `CODEX_MODEL` (default
  `gpt-6.1-sol`); the judge is fixed at `claude-opus-5-5` and the waiting
  check at `claude-sonnet-5-5`. Reason: changing the agent never changes
  the grader. The table's first line names each model with its effort and
  the `claude` and `codex` versions.
- **Medium effort, always:** both agents under test and the judge run at
  medium effort, fixed in the script with no setting to change it, on
  every run rather than only once trials take longer than 30 minutes.
  Reason: the owner's token budget, and a switch by run time would make
  two runs of one case incomparable. The judge runs j times per run, so it
  is most of the cost.
- **`gpt-6-astra` is not used for the trials:** `dev/trials.md` says so,
  and the script does not refuse it as a `CODEX_MODEL`. Reason: a guard
  would be code for a choice only the owner makes (rule 13).
- **The judge sees what it judges:** both transcripts show each tool's
  result after its call, an error mark and its first 20 lines, and a
  `grove` command's output in full; the judge gets Read, Grep and Glob on
  a copy of the fixture at the commit the case's setup left, and its prompt
  says so and lets it check what the repository already answers. Reason:
  in the end state, a document the agent wrote from the owner's answer
  makes the question that got it look unneeded.
- **`RUNS=k`** (default 3) runs each case and harness k times, one after
  another, with cases and harnesses in parallel as now. Reason: at one run
  a count that moves by 1 cannot be told from noise, and PR #73's 34 and
  34 over three runs said nothing either way. It triples the time and cost
  of `just trials`.
- **`JUDGINGS=j`** judges each run j times from its saved `transcript.txt`,
  with no new agent session; a run's judged count is the median of its
  judgings. Below the table, each run whose judgings give a count different
  values gets a line naming its case, harness, run, count and each value.
  Reason: pinning the judge's model does not make it repeat itself, and
  reruns alone mix its noise with the agent's.
- **One row per case and harness.** At k = 1 each count is as now; at
  k > 1 each is its median and min–max over the runs, such as `2 (1–4)`.
  Each run's detail stays in its trial folder.
- **`CONTRIBUTING.md` stays as it is.** With `RUNS` at 3, the table it
  asks for is already three runs. Still a record, not a gate.
- Defaults: trial folders are `<case>-<harness>-<run>`; `JUDGINGS` is 3;
  judgings are saved as `judge-<n>.json`; no task has a `Delegated:` line.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: none; this change alters no behavior of Grove.
- Tasks: [tasks.md](tasks.md)
