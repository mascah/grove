---
id: "G-118"
type: question
title: "What mandate should the G-108 paid eval runs have?"
status: resolved
created: "2026-09-23T19:54:34Z"
updated: "2026-09-23T19:54:49Z"
blocks: ["G-108"]
---

## Question

[G-108](G-108-workflow-evals.md) says paid trials need "a separately
authorized bounded execution mandate", and its Next asks the owner to agree,
at assignment, the model, the repeat count, the per-run budget, and the
fixture's brief and two topics. The headless `/grove-work G-108` session of
2026-09-23 was not given them, so it built everything that does not spend
and stops here. Four answers, and one approval, are needed before
`python3 evals/run.py run` can be started:

1. **Model** (`--model`). Recommendation: pin the model `claude -p` uses by
   default today, by full ID, since that is what a preview user gets; G-078
   finding 6 ran on that default (then Claude Opus 5). Alternative: the model
   the owner works with interactively.
2. **Repeat count** (`--runs`, per case). Recommendation: 5, as G-108
   proposes; 3 if the budget should be smaller. Few runs show patterns, not
   rates.
3. **Budget per run** (`--budget`, passed as `--max-budget-usd`).
   Recommendation: $3. G-078's headless shaping run took 92 seconds; the cap
   for 2 cases × 5 runs is then $30, printed before the runner starts.
4. **Permission mode** (`--permission-mode`). Headless shaping must create a
   worktree, run `grove new`, edit files and commit, and a denied prompt
   fails the run. Options: `bypassPermissions` (comparable with G-078's
   "write tools allowed"; the session can also write outside the disposable
   clone, which the checks cannot see but the trace shows) or `acceptEdits`
   with allowed Bash tools (safer, but the allow-list itself shapes behaviour
   and would need a runner option). Recommendation: `bypassPermissions`, in
   the disposable clone, with the transcripts read for writes outside it.
5. **Fixture approval.** The fixture and topics as built at `1f03a12` on
   `worktree-G-108`: [`evals/fixture/`](../evals/fixture),
   with the case table in [`evals/README.md`](../evals/README.md#what-a-run-does)
   and the design in [G-115](G-115-g-108-eval-skeleton-plan.md). The planted
   choice is whether `dropped` tasks count as "finished"; the companion's
   AND/OR choice is answered by the brief's filter convention. Approve, or
   name the change.

Also say where the clean `CLAUDE_CONFIG_DIR` should live and whether it
logs in with the owner's account or an `ANTHROPIC_API_KEY` (recommendation:
`~/.cache/grove-evals/claude`, logged in once with the owner's account).

Who can answer: the owner. The answer is the mandate: a successor session
runs `python3 evals/run.py run` with exactly those values and writes the
review G-108 acceptance 5 asks for.

## Next

Waits for the owner. Answer in this record, set `status=resolved`, then
assign `/grove-work G-108` again on `worktree-G-108`.

1. Use opus 5.5 for now
2. 5
3. $5
4. I think auto mode would be better than bypassPermissions
5. Approved as built. The planted choice is reversible, not fundamental:
   a run that blocks follows the guide as written, a run that ships the
   proposal with a non-blocking question or a proposed decision is what I
   want, and a run that only notes the choice in Next is the G-078
   failure. Report each run against those three. The lever, if the
   pattern shows it, is the headless bound's "must block" sentence.