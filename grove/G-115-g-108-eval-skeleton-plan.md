---
id: "G-115"
type: plan
title: "G-108 eval skeleton plan"
status: current
created: "2026-09-23T19:50:02Z"
updated: "2026-09-23T19:50:05Z"
work: ["G-108"]
---

## Design

Prepared 2026-09-23 on `worktree-G-108` from main `6208e82` against
[G-108](G-108-workflow-evals.md) at `sha256:ff01c9f0…`, in a headless
`/grove-work G-108` session. Harness here: Claude Code 2.1.281, go 1.26,
git 2.55.0, Python 3.13, macOS Darwin 25.6.0. Everything below is proposed
design within G-108's selected skeleton unless it quotes the record.

**Where it lives.** A top-level `evals/` directory: `evals/run.py` (stdlib
Python, like `internal/tui/testdata/terminal.py`), `evals/fixture/` (the
synthetic project's files, editable for review) and `evals/README.md` (the
documented command, what a run retains, the rubric, limits). No Go code, so
nothing joins `go test` or CI.

**Fixture.** A to-do tool, `tasks`: one person's tasks as Markdown files with
frontmatter under `tasks/`, a 60-line `tasks.py`, a brief with purpose,
constraints (a task file never moves or is renamed; no index or cache; `list`
output is for people and `export` is the stable format) and conventions
(repeated values of one filter flag match any of them, different flags all
must match; a command that hides tasks by default says how many and offers
`--all`). Task statuses are `open`, `done` and `dropped`. One unrelated
proposed record, "Sync tasks between two machines", gives `grove list` a
neighbour. `AGENTS.md` says only how proposal branches are named
(`worktree-shape-SLUG` under `.claude/worktrees/`), which acceptance 2 checks,
and `CLAUDE.md` imports it. The runner builds it with `grove init` in a
disposable directory, so the adapters are the binary's.

**Cases.** Missing choice, the G-078 finding 6 seed transposed: "hide finished
tasks from tasks list by default". The planted choice is which statuses are
"finished": `done` only, or `done` and `dropped`. The brief defines both
statuses and settles the flag, the count and the default output's audience,
so this is the one open product choice, and acceptance cannot be written
without it. Companion: "let tasks list filter by tag". Its apparent choice,
whether repeated `--tag` values combine as AND or OR and how they combine
with `--status`, is answered by the brief's filter convention, and tags are
already lowercased on `add`. Constraints a proposal must respect: no moving
done tasks to an archive (hide), no tag index (tag).

**Runner.** `python3 evals/run.py run --runs N --budget USD --model MODEL
--permission-mode MODE --config-dir DIR [--case NAME] [--out DIR]`. Every
spend parameter is required: there are no defaults, so it cannot run by
accident, and it prints the total cap before starting. It builds the CLI from
this checkout, builds the fixture once, and per case and run makes a bare
remote from the fixture and a clone of it, then runs `claude -p "/grove-shape
TOPIC --interaction headless" --output-format stream-json --verbose
--max-budget-usd USD --model MODEL --permission-mode MODE
--permission-prompts none --no-session-persistence` in the clone, with
`CLAUDE_CONFIG_DIR=DIR`, the built `grove` first on `PATH`, and Git's
location variables and other `CLAUDE*` variables removed. It refuses a config
directory holding `CLAUDE.md`, skills, agents, commands, plugins, output
styles or hooks, and records its listing. A clean directory needs one login,
or `ANTHROPIC_API_KEY`, which passes through.

**Retained per run:** `transcript.jsonl`, `stderr.txt`, `state.json` (the
clone's branches, HEAD, worktrees, status, the remote's refs, each proposal
branch's records), `run.json` (CLI version and guide digest, harness version,
model requested and reported, cost, turns, duration, exit status, checks,
retrieval facts). A `report.md` across runs reports every check per run,
unavailable harnesses, unrun cases and the Codex row as not built.

**Checks** (acceptance 2), on the clone after the process exits, each pass,
fail or the reason it could not be judged: `proposal-branch` (one
`worktree-shape-*` branch); `question-blocks-proposal` or `no-question`;
`proposal-proposed`; `message-names` (question ID, branch and the branch tip's
abbreviated commit in the final result text); `session-checkout-unchanged`
(branch, HEAD, `main` and a clean status); `remote-unchanged`
(`git ls-remote` equal before and after); `check-passes` (`grove check` in a
checkout of the proposal branch); `no-promotion` (no record the run touched
has a status other than `proposed`, `open` or `current`).

**Retrieval facts** (acceptance 3), from `tool_use` blocks: whether the
guide was printed (`grove guide shape`), the brief read (`grove brief` or
the file), `list`, `context` or `show` run; every file read, with those no
step needed (outside the case's expected set) listed separately. Reported,
not scored.

**Rubric** (acceptance 4) in `evals/README.md`: the four questions, a 0 to 2
scale with anchors, and a scoring table whose scorer column says `owner` or
`judge`. No automated judge in the skeleton; a judge is whoever the owner
names, labelled as such.

**Offline self-check.** `python3 evals/run.py selftest` runs the whole runner
against a fake `claude` (the runner's `--claude PATH`) that performs a
scripted good outcome for both cases and a G-078-style bad one (no question,
acceptance presuming the choice, a write to the session checkout), and asserts
that every check passes for the good runs and that exactly the expected
checks fail for the bad one. It spends nothing; it builds, so it is not a Go
test.

**Owner inputs.** The model, the repeat count, the per-run budget, the
permission mode and the fixture as built are the owner's to agree (G-108
Next). A headless session cannot, so they go to one question that blocks
G-108; building the runner and fixture does not presume the answers, because
every one is a runner argument or an editable fixture file, and nothing paid
runs without them.

## Steps

1. Commit this plan; set G-108 `active`.
2. Build `evals/fixture/`, `evals/run.py` and `evals/README.md`; the selftest
   passes. Add the README row to "Where each subject lives".
3. Persist the mandate question with `blocks` G-108; checkpoint G-108's Next.
4. Independent review of the runner, fixture and checks; fix and re-review
   within three rounds.
Status, 2026-09-23: steps 1 to 4 done (`d9dbf23`, `4b5b345`, question
[G-118](G-118-what-mandate-should-the-g-108-pa.md) at `275bd41`, review
[G-119](G-119-g-108-eval-skeleton-review.md) examined `1f03a12` with no open
findings). G-118 was resolved at `1d3ad82` and the eval login
[G-121](G-121-how-do-the-g-108-eval-runs-log-i.md) at `e57a4c6`.

5. After the owner answers: the paid runs, then a review record reporting
   the pattern against G-078 finding 6 (acceptance 5), then handoff.
   Done 2026-09-24: runs at `d565fcf` (which also let the runner accept a
   login's settings and synced skills), review
   [G-122](G-122-g-108-baseline-runs-the-missing.md), handoff in G-108's
   Next.
