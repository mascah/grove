---
id: "G-136"
type: plan
title: "G-134 plan: plan bound, per-phase effort and the grove-reviewer definition"
status: current
created: "2026-09-24T22:06:52Z"
updated: "2026-09-24T22:07:12Z"
work: ["G-134"]
---

## Design

Written at `98ce628` from the code G-134's Constraints name, read in full:
`internal/attempt/{attempt,facts}.go`, `internal/cli/{cli,init}.go`,
`guides.go`, `internal/tui/{attempts,review}.go`, both `grove-work`
adapters, and the provider's own help (`claude --help`, 2.1.282) and a
real result event (`G-108.20260924T010336Z`).

- **The bound is data in the assignment.** The assignment grammar becomes
  IDs, then optionally `--until plan`, then optionally `--interaction MODE`.
  Step 4 of the work guide states what the bound does, once; Inputs points
  at it. The Lifecycle paragraph is reconciled: preparation precedes
  `active`, so a bounded attempt leaves the status as found. Both of this
  repository's adapters and init's managed `assignmentData` accept it, and
  the Claude `argument-hint` names it.
- **Runner.** `Request` and `Launch` gain `Until` and `Effort`; the prompt
  is `/grove-work ID --until plan --interaction headless` when bounded, and
  `--effort E` follows `--model M` when given. `--until` accepts only
  `plan`; `--effort` is one token passed through, since the provider owns
  its vocabulary (`low` … `max` today) and refuses others before spending.
  `Launch.Reviewer` records `sha256:…` of the worktree's
  `.claude/agents/grove-reviewer.md` at launch, or `none`; empty means an
  attempt from before this field. `Final.ModelCostUSD` keeps
  `modelUsage[model].costUSD` from the result event. `Facts` prints a
  `Requested:` line (bound, model, effort, reviewer) and a per-model cost
  line.
- **Reviewer definition.** The file `.claude/agents/grove-reviewer.md` is
  the one owner: embedded by package `grove` beside the guides (a named
  dot-path embeds; checked), written verbatim by `init` as a managed file,
  and committed here as init writes it, so `init` in this checkout reports
  it `unchanged`. Frontmatter: `model: inherit`, `effort: high`,
  `disallowedTools` for the editing tools (the provider reads model,
  effort and tools from agent frontmatter). Its body is the standard brief.
  Step 6 names it and says what the session passes it. The Codex adapter
  is unchanged.
- **Board.** `R` asks, after budget and mode, three optional prompts where
  Enter skips: bound (`plan` or nothing), model, effort. A finished,
  successful, latest attempt launched with `until: plan`, no open question
  and the record not in review is a new outcome `plan ready`, under Needs
  you, next action "o opens WORK and its plan; R launches the
  implementation". The attempt screen shows the requested bound, model and
  effort beside the actual model, and the per-model cost.
- **Docs.** `docs/commands.md` (run flags, attempt facts, init output),
  `docs/board.md` (the prompts and the outcome), and one sentence in
  G-055.

Not done here: acceptance 5 is the experiment on the first real assignment
after this lands (G-135), outside this work by its own terms.

## Steps

1. Guide, adapters, `assignmentData`, reviewer definition, embed, init.
2. Runner fields, command, reviewer digest, per-model cost, facts; CLI
   flags and usage.
3. Board prompts, plan-ready outcome, attempt screen.
4. Docs and G-055.
5. Tests: runner command/facts with the fake provider, init verdicts for
   the new file, board standings and prompts; full verification; the
   pseudo-terminal script.
6. A real bounded attempt in a disposable clone (acceptance 3), then an
   independent review through the new definition's brief on the combined
   diff.
