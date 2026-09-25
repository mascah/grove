---
id: "G-140"
type: work
title: "Default an attempt's budget, mode, model and effort from grove.yaml and launch from one line"
status: active
created: "2026-09-24T23:52:07Z"
updated: "2026-09-25T04:12:06Z"
relates_to: ["G-045", "G-046", "G-134", "G-135"]
---

## Outcome

The owner launches an attempt from the board with one keypress after `R`,
or with `grove run ID` and no flags, when `grove.yaml` holds the project's
launch defaults; a flag or a typed override still changes any value for one
launch; and every attempt records the values it ran with, as it does now.

Owner intent, shaping conversation 2026-09-24: after
[G-134](G-134-bound-an-attempt-at-its-plan-and.md) there are too many
questions when launching an attempt, and the defaults need a home such as
`grove.yaml`. The owner read the investigation below and chose this design
over the same defaults with five prefilled prompts, over provider-side
settings alone, and over only skipping the optional prompts. The key
spelling and the launch line are proposed design where labelled.

## Constraints

Observed at main `132709e`, 2026-09-24:

- **Five prompts on the board, two required flags.** `R` asks budget,
  permission mode, bound, model and effort in sequence, one line each
  (`launchKey` in [attempts.go](../internal/tui/attempts.go), the prompt
  texts in [review.go](../internal/tui/review.go)); Enter alone on the last
  three asks for nothing. `run` requires `--budget` and `--permission-mode`
  and refuses without them (`Start` in
  [attempt.go](../internal/attempt/attempt.go)); `--until`, `--model` and
  `--effort` are optional. G-134 added the three optional prompts; before it
  there were two.
- **What was typed.** The nineteen attempts under `.git/grove/attempts/`,
  from every `attempt.json`: budget `50` on seventeen, `30` once (G-135's
  bounded run under its experiment cap) and `50` once; permission mode
  `auto` on all nineteen; model empty on seventeen and `opus` on the two
  G-135 attempts; effort set once (`high`); bound set once (`plan`). Two
  prompts retype constants and three are Enter.
- **Why they are required today.** [G-045](G-045-durable-attempt.md) made
  both flags required "since a default for either would make the choice in
  practice" ([G-100](G-100-g-045-durable-attempt-plan.md), Budgets and
  bounds); [G-046](G-046-managed-runs.md) kept them typed with no prefill
  for that reason; G-134 put "role profiles or model defaults in
  `grove.yaml`" out of scope, "as G-045 chose". This record reverses that
  one sentence of G-045, on the owner's request: the choice has been made
  in practice by repetition, and a default in a committed file is the choice
  made once, attributably, in Git history, where a retyped value is not.
- **`grove.yaml` today.** `parseConfig` in
  [project.go](../internal/project/project.go) accepts exactly
  `schema_version`, `records`, `brief` and `target` and reports any other key
  as unknown; values are yaml.v3 nodes read as scalars, with no helper for a
  nested mapping; `check` validates the file; `init` writes the first three
  keys (`defaultConfig` in [init.go](../internal/cli/init.go)). The versions
  loader parses every branch's committed `grove.yaml` through the same
  function, but a launch loads only the launching checkout's project, so
  branches need not agree on defaults. The
  [record model](../docs/record-model.md#configuration-and-discovery) owns
  the key list.
- **The provider.** Claude Code 2.1.282 (`claude --help`): `--effort` takes
  low, medium, high, xhigh or max; `--permission-mode` takes acceptEdits,
  auto, bypassPermissions, manual, dontAsk or plan; `--max-budget-usd` works
  only with `--print`. Its documentation, read on 2026-09-24 through the
  claude-code-guide agent: settings.json can set `model` and `effortLevel`;
  nothing but the flag sets a budget; `permissions.defaultMode` is honoured
  from user or managed settings only, never from a project's settings file.
  The owner's user settings already set `defaultMode: auto`. A provider-side
  default reaches an attempt's facts only as the init event's actual model
  and mode, never as a requested effort.
- **Code shape.** `internal/cli` imports `internal/tui`, and both import
  `internal/attempt`, so a flag table shared by `run` and the board belongs
  in the attempt package; today the table is in
  [cli.go](../internal/cli/cli.go), and `Start` validates the budget before
  it loads the project. The board's resolver carries the live checkout's
  `target` (`versions.Result`) and can carry its launch defaults the same
  way. `TestLaunchFromTheDetail`, `TestRefusals`, `TestAttemptCommandsUsage`
  and three scenarios in
  [terminal.py](../internal/tui/testdata/terminal.py) drive the prompts and
  refusals.

In scope, proposed design:

1. **`run:` in `grove.yaml`.** An optional mapping with optional `budget`,
   `permission_mode`, `model` and `effort`, named as `run`'s flags:

   ```yaml
   run:
     budget: 50
     permission_mode: auto
   ```

   `check` refuses a budget `ValidBudget` refuses, a value with whitespace,
   and any other key under `run`. No default for the bound: a bound is one
   launch's mandate (G-134), and the brief selects no universal plan gate.
   `init` does not write `run:`, so a fresh project keeps G-045's rule until
   its owner adds a default. The alternative spelling, flat top-level keys,
   needs no nested-mapping helper at the cost of a `model:` beside
   `records:`; the implementer chooses and the record model says which.
2. **`run` and `Start` default from it.** `Start` loads the project first,
   fills each empty value from `run:`, then validates as now; the refusal
   for a missing budget or mode stays and names both ways to supply one.
   `attempt.json` records the resolved values, and nothing distinguishes a
   default from a flag. `run`'s usage and [commands.md](../docs/commands.md)
   say the two flags are required unless `grove.yaml` defaults them.
3. **One launch line on the board.** `R` replaces the five prompts with one
   line: the resolved launch (budget, mode, bound, model, effort) and where
   it runs, where Enter launches; text typed in `run`'s own flag spelling,
   such as `--until plan --effort xhigh`, overrides for that launch and is
   parsed by the table `run` uses, moved into the attempt package. Without a
   default for the budget or the mode, the same line says what must be
   typed and refuses to launch until it is. Esc cancels. The plan-ready
   outcome's "read it, then R" is unchanged.
4. **This repository's defaults.** Add `run:` with `budget: 50` and
   `permission_mode: auto`, the values the nineteen attempts used, to this
   repository's `grove.yaml`, as the first use.
5. Reconcile the [record model](../docs/record-model.md) (configuration),
   [commands.md](../docs/commands.md) (run and init) and
   [board.md](../docs/board.md) (R), each saying its part once.

Out of scope: per-phase or role profiles, such as one effort for `--until
plan` and another for the implementation, which G-134 defers until
[G-135](G-135-run-the-g-108-eval-pair-on-codex.md) reports; a default
bound; remembering the last launch; user-level Grove configuration; writing
defaults in `init`; and any change to what an attempt records.

## Acceptance

1. `grove.yaml` with `run: {budget: 50, permission_mode: auto}` passes
   `check`; `run: {budget: x}`, `run: {effort: "a b"}` and `run: {until:
   plan}` are refused by `check` with the key named; a file without `run:`
   means what it means today.
2. With those defaults, `grove run ID` with no flags starts an attempt whose
   command carries `--max-budget-usd 50 --permission-mode auto` and whose
   `attempt.json` shows them; `grove run ID --budget 2` carries 2 and
   `auto`; without `run:` and without the flags, `run` is refused as today
   by a message that also names `run:` in `grove.yaml`. Shown with the fake
   provider, and once with the real one in a disposable clone reached by an
   absolute `--project` path.
3. `R` on the board shows one line naming the resolved launch and where it
   runs; Enter starts the attempt with those values; typing `--until plan
   --effort xhigh` on that line starts one with the bound and effort and the
   defaults for the rest; a line that leaves the budget or the mode
   unsupplied is refused without launching; a word the table does not know
   is refused with the message `run` gives. `TestLaunchFromTheDetail` and
   the three pseudo-terminal scenarios cover it, and the five prompts they
   expect today are gone.
4. The usage text, the record model, `commands.md` and `board.md` describe
   the defaults and the launch line once each, every link resolves, and no
   document still states G-045's sentence as current.
5. `gofmt`, `go vet`, `check`, the full uncached suite and `terminal.py`
   pass, with no package over five seconds.

## Next

Proposed, unassigned. To carry it out: `/grove-work G-140`, or `R` on the
board; the design above stands in for a plan unless the implementer finds
the nested `run:` mapping costs more than a small helper in `parseConfig`,
in which case the plan records the flat spelling and why. The owner's
verdict on the candidate is the check that one keypress after `R` is what
they wanted.
