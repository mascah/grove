---
id: "G-158"
type: question
title: "What mandate should the G-154 without-row runs have?"
status: open
created: "2026-09-25T19:36:20Z"
updated: "2026-09-25T19:48:20Z"
blocks: ["G-154"]
relates_to: ["G-118", "G-122", "G-141"]
---

## Question

## Next

## Question

[G-154](G-154-listed-constraint-eval.md) acceptance 2 runs its two cases
under "an owner mandate naming a Claude model, runs, budget and permission
mode", and its Next says the owner supplies it before any paid run. The
headless `/grove-work G-154` session of 2026-09-25 was given none, so it
built the cases (`b2c7ef0` on `worktree-G-154`, `selftest: ok`) and stops
here. Under [G-141](G-141-never-run-gpt-6-astra-unless-the.md), an item left
blank does not take the recommendation: it leaves this question open and the
runs unstarted. Answer each:

1. **Model** (`--model`, a Claude model by full ID). Recommendation:
   `claude-opus-5-5`, as [G-122](G-122-g-108-baseline-runs-the-missing.md)
   ran, so the rows compare with the pair's baseline.
2. **Runs per case** (`--runs`). Recommendation: 5, as G-122; the `with`
   row after [G-153](G-153-search-and-code-links.md) uses the same number.
3. **Budget per run** (`--budget`, USD). Recommendation: $5, as G-118 set
   for G-122, whose runs cost $0.27 to $0.32. These runs read more records,
   so expect more, not above $1. The cap printed first is runs × 2 × budget:
   $50 at the recommendation.
4. **Permission mode** (`--permission-mode`). Recommendation: `auto`, as
   G-118 chose and G-122 ran with no denials.
5. **Config directory** (`--config-dir`). Recommendation:
   `~/.cache/grove-evals/claude`, the G-122 login, which the runner still
   has to accept as clean at run time.

The command the answer authorizes, from `worktree-G-154` as it stands (the
`without` row, guides digest before G-153):

```sh
python3 evals/run.py run --case listed-constraint --case code-constraint \
    --runs RUNS --budget USD --model MODEL --permission-mode MODE \
    --config-dir DIR --out ~/.cache/grove-evals/runs/2026-09-25-G-154-without
```

It never passes `--harness codex` and never names `gpt-6-astra` (G-154).

Who can answer: the owner. Answer in this record, set `status=resolved`,
then assign `/grove-work G-154` again on `worktree-G-154`.

## Next

Waits for the owner.
