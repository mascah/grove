---
id: "G-138"
type: plan
title: "G-135 Codex eval row plan"
status: current
created: "2026-09-24T22:37:53Z"
updated: "2026-09-24T22:39:42Z"
work: ["G-135"]
---

## Design

Prepared 2026-09-24 on `worktree-G-135` from main `132709e` against
[G-135](G-135-run-the-g-108-eval-pair-on-codex.md) at `sha256:0baef170…`,
in a headless `/grove-work G-135 --until plan` session. Here: codex-cli
0.156.1, Python 3.13 (`tomllib` available), guides digest `0c163c41a0f2`
(`grove version` at `132709e`, after G-134's guide edits; G-122 ran at
`3f5487904c61`). Everything below is proposed design within G-135's
selected scope unless it quotes the record.

**Observed here, free.** One `codex exec --json --ephemeral
--ignore-user-config -s workspace-write 'say hi'` in a new Git repository
under `/tmp`, with a new empty `CODEX_HOME` and no login, exited 1 at no
cost. Its stdout, JSONL: `thread.started` (`thread_id`), `turn.started`,
repeated `error` (`message`), one `item.completed` whose `item` has `id`,
`type: "error"` and `message`, and `turn.failed` (`error.message`). Stderr
began "Reading additional input from stdin...", so stdin must stay closed.
The new home then held `config.toml` with a single
`[projects."<cwd>"] trust_level = "trusted"` table, written although user
config was ignored; `skills/.system/` with six bundled skills
(review-agent, skill-creator, plugin-creator, skill-installer, openai-docs,
imagegen); `.tmp/plugins/` (a marketplace clone, nothing installed); SQLite
state; `tmp/` and `shell_snapshots/`. `codex features list` reports
`memories` stable and off by default. Command, agent-message, usage and
file-change events were not observed: nothing ran past authentication.

**Interface.** `python3 evals/run.py run --harness codex --runs N --model
MODEL --effort EFFORT --permission-mode MODE --config-dir DIR --max-seconds
S [--case NAME]... [--out DIR] [--codex PATH]`. `--harness` defaults to
`claude`, whose required arguments, command and output stay as they are.
Required-ness moves from argparse into `run()`, per harness: Claude needs
`--budget`; Codex needs `--model`, `--effort`, `--permission-mode`,
`--config-dir` and `--max-seconds`, and refuses `--budget`, since no flag
would enforce it. Missing `--max-seconds` is the cap refusal of acceptance
1. `--effort` and `--max-seconds` are Codex only and refused with Claude.
The printed cap before the first run is `runs × cases × S` seconds of
Codex time, with the note that Codex bounds no dollars.

**Command**, in the clone, stdin closed, own process group as today:

```sh
codex exec --json --ignore-user-config --disable memories \
    -m MODEL -c model_reasoning_effort=EFFORT (-s MODE | --approve-for-me) \
    -o RUNDIR/last-message.txt '$grove-shape TOPIC --interaction headless'
```

`--permission-mode` takes a sandbox value (`read-only`, `workspace-write`,
`danger-full-access`) passed as `-s`, or `approve-for-me`, which routes
approval requests through Codex's automatic reviewer on the
workspace-write sandbox and is the nearest equivalent of the Claude row's
`auto`. `--ignore-user-config` keeps `config.toml` out while auth still
uses `CODEX_HOME`; `--disable memories` makes the default explicit, as the
Claude row requires `autoMemoryEnabled: false`. Not `--ephemeral`: the
session's rollout file under `DIR/sessions/` is the only place the
reported model, effort, approval and sandbox policy appear, so the runner
finds it by `thread_id` and keeps a copy as `rollout.jsonl`; where it is
missing or unreadable those facts are "not reported" with the reason. The
wall-clock cap replaces the fixed `TIMEOUT` for Codex and kills the
process group as today. The environment is `env()` with every `CODEX_*`
variable removed except `CODEX_API_KEY`, then `CODEX_HOME=DIR` and the
built `grove` first on `PATH`.

**Clean `CODEX_HOME`.** Refused when it holds `AGENTS.md`,
`AGENTS.override.md`, `rules/`, `prompts/`, `hooks.json` or `hooks/`, a
non-empty `memories/`, any `skills/` entry but `.system`, or a
`config.toml` with anything but `projects.*.trust_level` (parsed with
`tomllib`), which Codex writes per clone. Recorded, not refused: the
listing, the bundled `skills/.system` names, and `codex login status`
(which auth the runs used). The runner's `CUSTOMIZATION` check stays
Claude's; Codex gets its own list beside it.

**Parsing, retained per run.** `transcript.jsonl` is the `--json` stream.
From it: `thread_id`; `item.completed` items of type `command_execution`
(`command`), `agent_message`, `file_change`, `mcp_tool_call` and
`web_search`; `turn.completed` `usage` (input, cached input and output
tokens); `turn.failed` and `error` events for the harness column. The
final message is `last-message.txt`. `run.json` keeps the Claude fields'
names: `cost_usd` null with `cost_reason` ("Codex reports tokens, not
dollars"), `tokens`, `turns` as the number of `turn.completed` events and
`tool_calls` beside it, `model_reported` and `permission_mode_reported`
from the rollout, `result_subtype` from `turn.completed` or `turn.failed`,
`permission_denials` null with the reason. Shapes above that were not
observed come from Codex's documented `exec --json` events; the first paid
run is read before any other, and parser and fake follow what it shows.

**Checks and retrieval.** `checks()` and `state()` are unchanged, so
acceptance 2's clone checks are the same code. `retrieval()` splits into
a per-harness extraction of (commands, files read) and the shared
classification it does today. Codex has no Read tool: every read is a
shell command, and Codex wraps each as `SHELL -lc 'SCRIPT'`, so the
extraction unwraps `-lc`/`-c` before the existing `cat`/`head`/`sed`/`nl`
rules. `NEEDED` gains `.agents/skills/grove-shape/SKILL.md`, the adapter the
fixture's `grove init` writes, since Codex reads a skill as a file.

**Report.** `report.md` names the harness and its version, the requested
and reported model and effort, permission mode, the cap, the config dir's
listing and system skills, and the login mode; the "Codex row: not built"
line goes. The cost column prints "not reported" and a tokens column is
added for Codex.

**Selftest.** The `_fake` entry acts as `codex exec` when its first
argument is `exec`: the same good, bad, surfaced and worse modes, emitting
Codex events (`thread.started`, `command_execution` items wrapped as
`/bin/zsh -lc '…'`, `agent_message`, `turn.completed` with usage), writing
`-o`'s file and a rollout with a `turn_context` under `CODEX_HOME`. It
asserts the same expected check failures per mode and case, the same
retrieval facts, the reported model from the rollout, and refusals: no
`--max-seconds`, `--budget` with Codex, `--effort` with Claude, and a home
holding `AGENTS.md`, `skills/mine`, a `config.toml` with `model`, or a
non-empty `memories/`, while one with a trust table and `skills/.system`
passes. The existing Claude selftest runs unchanged.

**README.** A Codex section beside the Claude one: the command, the
login (`CODEX_HOME=DIR codex login` once, or `codex login --with-api-key`
from stdin, or `CODEX_API_KEY`), what the home may hold, the cap, cost
"not reported"; Limits say which rows exist.

**Mandate.** Runs per case, model, effort, permission mode, cap, the
login and whether the Claude pair reruns at the same guides digest are
the owner's: question
[G-139](G-139-what-mandate-and-login-should-th.md) blocks G-135. Building
the runner does not presume the answers: each is a runner argument.

**Report record and disposition** (acceptance 4). One review record with
`work` G-135 and `examined` the commit the runs used: configuration,
per-case pattern against [G-122](G-122-g-108-baseline-runs-the-missing.md)
and whatever same-digest Claude rerun G-139 allows, tokens and time,
retrieval, what the clean home still carried, limits, and whether a Codex
provider in the attempt runner is worth shaping. What such a seam would
need, to be judged against the runs: the command line
([attempt.go](../internal/attempt/attempt.go)), event parsing and the
activity feed ([activity.go](../internal/attempt/activity.go)) for Codex's
item events, a session id Codex generates rather than Grove, Stop by
process group as today, a budget that can only be time or tokens, the
permission profile's mapping onto sandbox and approval, and the review
gate, whose `grove-reviewer` is a Claude agent definition; and whether
[G-101](G-101-attempt-mechanism.md) would be revisited. No product change.

**Adjustment, 2026-09-24, implementation.** The owner's login left the
eval home's `config.toml` with a `[tui]` table (screen-reader detection,
model-availability notice) beside a trust table: terminal state, as
Claude's `tui` and `theme` settings are, and ignored anyway under
`--ignore-user-config`. The runner allows `tui` as well as
`projects.*.trust_level`; anything else is still refused.

**Adjustment, 2026-09-24, review G-138 step 3.** Codex runs each command
as `SHELL -lc`, and the owner's `~/.zprofile` then puts the installed
`~/.local/bin/grove` before the built one (observed: `env PATH=DIR:$PATH
/bin/zsh -lc 'command -v grove'` found `~/.local/bin/grove`). The Codex
environment adds a `ZDOTDIR` of the runner's own, whose `.zprofile`
restores the runner's `PATH` after `/etc/zprofile`'s `path_helper`
reorders it, and the runner refuses to start unless the login shell finds
the built `grove` with that `PATH`, so the Codex sessions' tools are the
Claude row's.

**Adjustment, 2026-09-24, after the runs.** The attempt ran 9½ Codex runs
on `gpt-6-astra` at `high`, which G-139 recommended and the owner never
approved, and took the ChatGPT Plus five-hour window from 12% to 95% before
the owner stopped it. The owner's rule is
[G-141](G-141-never-run-gpt-6-astra-unless-the.md). The runner now requires
`--max-plan-percent` and records each run's plan use from the rollouts'
`rate_limits`; no further runs are spent, and the report covers the 9 made.

## Steps

1. Commit this plan and G-139; checkpoint G-135's Next (this session,
   bounded at the plan; status stays `proposed`).
2. Set G-135 `active`. Build the Codex branch of `evals/run.py`, its fake
   and selftest, and the README section; `python3 evals/run.py selftest`
   passes; `go run ./cmd/grove check` passes.
3. Independent review of the runner before anything is spent; fix and
   re-review within three rounds.
4. Once G-139 is resolved: one run of `missing-choice` under the mandate,
   read in full; adjust parser and fake to what it shows (re-review if the
   change is more than field names); then the remaining runs, and the
   same-digest Claude pair if G-139 allows it.
5. The report review record, a final independent review of the combined
   diff, and the handoff into Review.
