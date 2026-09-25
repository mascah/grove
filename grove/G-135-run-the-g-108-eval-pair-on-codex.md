---
id: "G-135"
type: work
title: "Run the G-108 eval pair on Codex"
status: review
created: "2026-09-24T21:48:22Z"
updated: "2026-09-25T04:02:07Z"
relates_to: ["G-050", "G-101", "G-108", "G-122", "G-134"]
candidate: "fec6e94d382aea3ab8b5b50c681145f91ad7c35b"
---

## Outcome

The owner has repeatable evidence of how the shaping guide behaves under
Codex, on the same fixture, cases, checks and rubric as the Claude baseline,
and a report that says whether a Codex provider in the attempt runner is
worth shaping and what it would take.

Owner intent, shaping conversation 2026-09-24: the owner wants to be able to
consider different harnesses for different execution phases, and agreed that
this row of the [G-108](G-108-workflow-evals.md) suite should decide whether
a provider seam is built before any is. G-108's Next listed the Codex row as
a follow-on shaped only if the skeleton's report said so; the owner's
harness question is the reason it is shaped now.

## Constraints

Observed at main `282d282`:

- [evals/run.py](../evals/run.py) drives one executable (`--claude`, default
  `claude`), composes `claude -p … --output-format stream-json` in a
  disposable clone with a bare remote, reads the init and result events for
  model, cost and turns, takes the retrieval facts from the transcript, and
  runs the clone checks; its selftest uses a fake `claude`.
  [evals/README.md](../evals/README.md) says the suite is headless shaping
  on Claude only. The output directory, `~/.cache/grove-evals/runs/`, is
  outside every checkout.
- [G-122](G-122-g-108-baseline-runs-the-missing.md): the Claude row at
  guides digest `3f5487904c61` on Claude Opus 5.5, 5 of 5 missing-choice
  runs blocking on the planted question and 5 of 5 companion runs asking
  nothing, $0.27 to $0.32 and 45 to 62 seconds a run; a login syncs account
  plugins and skills into the clean configuration. G-114's rerun at
  `5a224350feae` repeated both patterns. The current guides digest is
  `3c9996e33e42` (G-129); pin whichever the runs use.
- Codex 0.156.1 (`codex exec --help`): `--json` prints events as JSONL,
  `-m MODEL`, `-c model_reasoning_effort=…`, `-p PROFILE`, `-s read-only |
  workspace-write | danger-full-access`, `--approve-for-me`,
  `--dangerously-bypass-approvals-and-sandbox`, `-C DIR`, `--ephemeral`,
  `--ignore-user-config`, `--output-last-message FILE`; no budget flag.
  `CODEX_HOME` is the configuration directory, the analogue of
  `CLAUDE_CONFIG_DIR`. The event shapes are unobserved here. The owner's
  `~/.codex/config.toml` sets a model, high reasoning effort and plugins.
- [G-050](G-050-shaping-review.md) ran headless shaping on Codex 0.155.1 by
  hand with `codex exec --sandbox workspace-write`: it found the adapter
  (`.agents/skills/grove-shape/SKILL.md`), read the guide, refined the
  overlapping proposal, wrote the question with `blocks`, and committed on
  a proposal branch; the sandbox refused the default `GOCACHE`, which it
  moved under `/private/tmp`. That is the only Codex evidence, once, on an
  older guide.
- Grove-owned attempts run under the owner's configuration by design: each
  of the sixteen attempts' first events carry the owner's SessionStart hook
  output. The eval's clean configuration does not. A comparison between a
  harness row and the attempts must say so.
- [G-101](G-101-attempt-mechanism.md) pins attempts to a Claude process;
  a Codex provider would revisit it. The runner's Claude-specific parts are
  the command line, the init and result parsing and the activity feed
  ([attempt.go](../internal/attempt/attempt.go),
  [activity.go](../internal/attempt/activity.go)); the rest is provider
  neutral.

In scope, proposed design: a harness branch in `evals/run.py`
(`--harness codex`, spelling proposed) that composes `codex exec --json` in
the clone with a clean `CODEX_HOME`, a sandbox and approval setting
equivalent to the Claude row's permission mode, an explicit cap (turns or
time) since Codex has no budget flag, the same cases, clone checks and
rubric, retrieval facts and cost from Codex's events where they exist and
"not reported" where they do not, a selftest with a fake `codex`, and a
README update. One linked review record reports the pattern per case against
G-122, configuration, cost, limits and the disposition: whether a Codex
provider in the attempt runner is worth shaping, what the seam would need
(command, events, stop, budget, review), and whether the guide text needs a
Codex-specific change. Paid runs need the owner's mandate at assignment.

Out of scope: changing `attempt.go` or the guides, the work row on Codex,
the owner-configuration comparison, and any Codex-specific adapter change
unless the report shows the adapter failing.

## Acceptance

1. `python3 evals/run.py run --harness codex …` runs both cases a chosen
   number of times with a clean `CODEX_HOME`, refuses without an explicit
   cap, and retains per run the transcript, the clone's final state, Codex
   version, model, guides digest, and cost or "not reported" with the
   reason. `--harness claude` (the default) behaves as before.
2. The clone checks of G-108 acceptance 2 are reported per run unchanged;
   retrieval facts come from the Codex trace or are reported unavailable
   with the reason.
3. `python3 evals/run.py selftest` covers the Codex branch with a fake
   `codex`, including the cap refusal.
4. A linked review record reports the pattern per case against G-122, with
   configuration, cost and limits, and states whether a Codex provider in
   the attempt runner is worth shaping and what it would take, or that it
   is not. No product change happens in this work.

## Next

Owner decision, 2026-09-24: assign this after
[G-134](G-134-bound-an-attempt-at-its-plan-and.md) lands, as the target of
G-134's routing experiment: a bounded preparation attempt at `xhigh`, then
implementation at `medium`, both on Opus 5.5. A Codex provider, if the
report favours it, is separate work that revisits G-101.

Handoff, 2026-09-25, interactive `/grove-work G-135` on `worktree-G-135`,
base main `132709e`, from G-135 at `sha256:0baef170…` and plan
[G-138](G-138-g-135-codex-eval-row-plan.md) at `sha256:4f455a34…`.

- **Incident.** The headless attempt built the runner (`c57fd37`,
  `e00a92e`, `1e1a345`, reviewed at plan step 3), then ran Codex on
  `gpt-6-astra`/`high`, which G-139 recommended and the owner never
  approved: G-139 let a blank item take the recommendation. Nine and a half
  runs took the ChatGPT Plus five-hour window from 12% to 95%; the owner
  stopped the attempt. The owner's rule is
  [G-141](G-141-never-run-gpt-6-astra-unless-the.md); the guide fix is
  captured, unassigned, as [G-142](G-142-keep-a-blank-mandate-answer-from.md).
  The owner rejected putting the rule in `AGENTS.md`.
- **Acceptance 1.** `--harness codex` as built, plus a required
  `--max-plan-percent` (`48358b7`, `f37be4e`, `4513c69`): no run starts
  once the newest five-hour reading plus the largest run's use, at least
  13 points, would reach it; each run's plan use is in `run.json` and the
  report. `--harness claude` unchanged.
- **Acceptance 2.** Clone checks unchanged, reported per run; retrieval
  from the Codex commands.
- **Acceptance 3.** `python3 evals/run.py selftest` covers the Codex branch,
  the time and plan cap refusals, the stop, the new-home floor, a reset
  window and a missing `resets_at`: "selftest: ok" at `4513c69`.
- **Acceptance 4.** [G-143](G-143-g-135-codex-eval-row-pattern-pla.md), on
  the nine completed runs only, by the owner's direction: the pattern
  matched Claude's (5 of 5 blocking, 4 of 4 no question), Codex read
  outside the clone in every run and read the owner's Grove checkout in
  five, about 9 plan points a run; a Codex provider is not worth shaping
  now. Its general guide finding is captured, unassigned, as
  [G-144](G-144-give-adopting-projects-the-recor.md). Companion 5 was not
  completed and nothing more was spent.
- **Verification** at `4513c69`: `python3 evals/run.py selftest` ok;
  `go run ./cmd/grove check` OK. No Go file changed since `132709e`, so the
  Go suite was not rerun.
- **Review.** [G-145](G-145-g-135-plan-cap-guard-and-report.md): three
  independent rounds, every consequential finding fixed; four low points
  open at the cap.
- **Next action.** Owner judgment: in this checkout, `go run ./cmd/grove
  approve G-135 "VERDICT"` or `go run ./cmd/grove feedback G-135 "TEXT"`;
  after approval, in main's checkout, `go run ./cmd/grove integrate G-135`.
  G-141, G-142, G-144 merge with it; G-142 and G-144 then await shaping.
