# Behavioral evaluations

These evaluations measure how Grove's headless shaping workflow guides an
agent, which the Go tests cannot: those check what `context` prints, not what
an agent decides after reading it. [G-108](../grove/G-108-workflow-evals.md)
owns the outcome and acceptance, and [G-115](../grove/G-115-g-108-eval-skeleton-plan.md)
the design. The first increment asks one question: does the workflow tell a
missing owner decision, which must become a question, from a choice the
project's brief already answers? [G-135](../grove/G-135-run-the-g-108-eval-pair-on-codex.md)
adds the same cases on Codex.

## Run

```sh
python3 evals/run.py selftest     # free: a fake claude and codex, asserts every check and retrieval fact
python3 evals/run.py run --runs 5 --budget 3 --model MODEL \
    --permission-mode MODE --config-dir ~/.cache/grove-evals/claude
python3 evals/run.py run --harness codex --runs 5 --model MODEL --effort EFFORT \
    --permission-mode MODE --max-seconds 600 --max-plan-percent 50 --config-dir ~/.cache/grove-evals/codex
```

`run` spends up to runs × cases × budget dollars and prints that cap first.
It has no default for any spend parameter, is never part of `go test` or CI,
and runs only under a mandate that names the model, repeat count, budget and
permission mode. `--case NAME` runs one case; `--out DIR` keeps the output
somewhere other than a new temporary directory, and must not sit under a
directory holding a `CLAUDE.md`, which the session would load.

`--config-dir` becomes `CLAUDE_CONFIG_DIR` for every run, so the owner's
settings, plugins, skills and `CLAUDE.md` stay out. The runner refuses a
directory holding any of them. A new directory has no login: log in once with
`CLAUDE_CONFIG_DIR=DIR claude`, or export `ANTHROPIC_API_KEY` or
`CLAUDE_CODE_OAUTH_TOKEN`, which pass through; every other `CLAUDE*`
variable is removed. `settings.json` must set `autoMemoryEnabled` to
false, since memory is on by default and is written under the directory;
it may also hold `tui` and `theme`, and nothing else. A login syncs
the account's Anthropic skills and plugins under `skills/synced` and
`plugins/synced`, which a preview user has too and which load into every
session; the report records both, and anything there that the sync's
manifest does not name is refused as authored. Managed (policy) settings
still apply and are outside the runner's control.

### Codex

`--harness codex` runs `codex exec --json --ignore-user-config --disable
memories -m MODEL -c model_reasoning_effort=EFFORT` with `$grove-shape TOPIC
--interaction headless` as the prompt. Codex has no budget flag, so no dollar
cap exists: `--max-seconds` is required, each run is killed at it, and the
runner prints the cap, runs × cases × seconds, before the first run.
A ChatGPT login spends the plan's five-hour window instead, and seconds do
not bound that: G-135's `gpt-6-astra` runs took 6 to 13 points of a Plus
window each ([G-143](../grove/G-143-g-135-codex-eval-row-pattern-pla.md),
[G-141](../grove/G-141-never-run-gpt-6-astra-unless-the.md)).
`--max-plan-percent P` is required too: before each run the runner takes
the last reading of the newest rollout under `CODEX_HOME` (0 once its window
has reset), and starts no run once that reading plus the largest run's use
so far, never less than 13 points, would reach P (a new home with no
reading counts as 0), or once a run left no reading; the report says where it stopped. P is not a ceiling: a run
already started is not stopped, and use elsewhere on the account since the
last reading is not seen. Model, effort, runs and both
caps are the owner's explicit answer, never a default (G-141).
`--budget` is refused with Codex, and `--effort`, `--max-seconds` and
`--max-plan-percent` with Claude. `--permission-mode` is a sandbox (`read-only`, `workspace-write`,
`danger-full-access`), passed as `-s`, or `approve-for-me`, the
workspace-write sandbox with Codex's automatic reviewer on approval
requests, the nearest to Claude's `auto`.

`--config-dir` becomes `CODEX_HOME`; every other `CODEX_*` variable but
`CODEX_API_KEY` is removed. Log in once with `CODEX_HOME=DIR codex login`, or
`codex login --with-api-key` reading the key from stdin, or export
`CODEX_API_KEY`; the report records `codex login status` and which of
`CODEX_API_KEY` and `OPENAI_API_KEY` are set. Codex runs each command in the
user's login shell, whose profile can put an installed `grove` before the
built one: the runner gives it a `ZDOTDIR` of its own, which keeps a zsh
user's startup files out and whose `.zprofile` restores the `PATH` that
`/etc/zprofile` reorders, and refuses to start unless `SHELL -lc` then finds
the built `grove` with the runner's `PATH`, the one the Claude row's
sessions get. Another login shell, such as bash, is refused. The runner refuses
a directory holding `AGENTS.md`, `AGENTS.override.md`, `rules`, `prompts`,
`hooks.json`, `hooks`, a non-empty `memories`, any `skills` entry but
`.system`, or a `config.toml` with anything but the `tui` state a login
writes and the per-directory `trust_level` Codex writes for each clone.
Codex installs its bundled skills under `skills/.system` on the first run;
the report lists them. Skills Codex discovers outside `CODEX_HOME`, such as
under `~/.agents`, are not checked. Cost is "not reported": Codex reports
tokens, which the report shows per kind with the plan points; tokens spent
by `approve-for-me`'s reviewer are not among them, though its plan use is
in the account-wide readings.

## What a run does

1. Builds the CLI from this checkout and records `grove version`, whose guide
   digest pins the guides under evaluation.
2. Builds the fixture once: [`evals/fixture/`](fixture) copied into a new Git
   repository, `grove init`, the fixture's brief, one unrelated proposed
   record, one commit. The fixture is `tasks`, a small to-do tool whose brief
   states constraints and conventions the cases depend on.
3. Per case and run: a bare remote cloned from the fixture, a clone of it,
   then `claude -p "/grove-shape TOPIC --interaction headless"`, or `codex
   exec … "$grove-shape TOPIC --interaction headless"`, in the clone with the
   built `grove` first on `PATH`.
4. Checks the clone after the process exits, and reads the trace.

| Case | Topic | Expected |
| --- | --- | --- |
| `missing-choice` | hide finished tasks from tasks list by default | A question blocking the proposal: whether `dropped` tasks count as finished, which the brief leaves open |
| `companion` | let tasks list filter by tag | A proposal and no question: the brief's filter convention answers how repeated `--tag` values combine |

The first transposes [G-078](../grove/G-078-g-039-trial-evidence-for-the-int.md)
finding 6 onto the fixture. The second catches a guide change that makes
agents ask about what the brief already answers.

## What it retains

In the output directory, per run `CASE-N/`: `transcript.jsonl` (the
stream-json or `--json` trace), `stderr.txt`, for Codex `last-message.txt`
and `rollout.jsonl` (the session's rollout, the only record of the model,
effort and policies it used), `state.json` (the clone's branches, HEAD,
status and worktrees, the remote's refs, and each proposal branch's commits,
touched records and `grove check` output, with the state before the run) and
`run.json` (command, harness version, `grove version` with the guide digest,
base and fixture commits, requested and reported model, cost, turns,
duration, exit, checks, retrieval facts, final message; for Codex also the
thread, reported effort, approval and sandbox policy, tokens and tool calls,
the first and last five-hour plan readings and the points the run used,
with a reason beside each fact Codex does not report). The clone
and remote stay too. `report.md` summarizes every run, lists cases not run,
and says when the harness was unavailable. Its harness column (exit status,
timeout, an error result, permission denials) separates a run the harness
never carried out, such as a missing login or a budget stop, from an agent's
behaviour. Ctrl-C, SIGTERM or SIGHUP stops the runner and kills the running session;
a runner failure on one run is reported, with what the run cost, and the
next run still starts.

## Checks

On the clone, never the final message except where named. Each is `pass`,
`fail` with the reason, or not judged when no single proposal branch exists.

| Check | Passes when |
| --- | --- |
| `session-checkout-unchanged` | The clone's branch, HEAD and `main` are as before, and its status is clean |
| `remote-unchanged` | The bare remote's refs are as before: nothing pushed |
| `proposal-branch` | Exactly one `worktree-shape-*` branch exists |
| `proposal-proposed` | The branch adds or changes at least one work record, and every one is `proposed` |
| `question-blocks-proposal` | (missing-choice) A question on the branch has `blocks` naming that work. The failure reason tells a choice surfaced without blocking (a non-blocking question or a decision record) from one at most noted in the record, the G-078 finding 6 outcome |
| `no-question` | (companion) The branch adds or changes no question |
| `no-promotion` | No record the branch touches is in any status but `proposed`, `open` or `current` |
| `check-passes` | `grove check` passes in a checkout of the branch |
| `message-names` | The final message names the branch, its tip commit (7 or more hex digits) and, for missing-choice, the blocking question's ID |

## Retrieval facts

From the trace's tool calls, reported and never scored: whether the session
printed the guide (`grove guide shape`), read the brief (`grove brief` or the
file), and ran `list`, `context` or `show`, counting only a command whose
program is `grove` and whose subcommand is that word; every file it read; and the reads
no step needed, meaning anything but `AGENTS.md`, `CLAUDE.md`, `grove.yaml`,
the brief, `tasks.py`, `tasks/`, the Codex adapter
`.agents/skills/grove-shape/SKILL.md` and the records it wrote. Reads through
`cat`, `head`, `tail`, `sed`, `nl`, `less` or `awk` count; `grep` and other
tools do not. Codex has no read tool: every read is a command, which it
wraps as `SHELL -lc 'SCRIPT'` and the runner unwraps; a Codex trace
with no command is reported unavailable, not as reading nothing. A `grove` run through a wrapper such as `timeout` or
inside `$(…)` is missed, and a heredoc line starting with `grove` is counted:
read the transcript before resting a conclusion on one fact.

## Rubric

A person (`owner`) or a model they name (`judge`) scores each run from its
`state.json`, the proposal branch and the final message, and writes the
scorer's label beside every score in the report's table. 2 is met, 1 partly
met, 0 not met.

| Question | 2 | 1 | 0 |
| --- | --- | --- | --- |
| Does acceptance presume the missing choice? (missing-choice) | No item depends on which statuses are finished | One item leans on an answer but says it is open | Acceptance hides `dropped`, or keeps it, as settled |
| Is the question the planted choice? (missing-choice) | It asks whether `dropped` tasks are finished, with options and a recommendation | It asks that among unrelated choices, or without options | It asks something else, or there is none |
| Does the proposal respect the brief's constraints? | Hide: no file moves or renames, count and `--all` per convention. Tag: no index, OR within `--tag`, AND with `--status` | One convention missed | A constraint broken: an archive directory, an index, AND within `--tag` |
| Is the handoff usable? | Branch, commit, records and the wait or next action, without the transcript | One of them missing | The owner needs the transcript |

The companion scores only the last two.

## Limits

- Headless shaping on Claude and Codex only. The work row through
  `grove run`, and every other case in G-108's Next, are not built.
- The checks see the clone. A session could write outside it, for example to
  the owner's home; the trace shows such writes, the checks do not.
- Few runs show patterns, not rates, and no result gates anything.
