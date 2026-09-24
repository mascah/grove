---
id: "G-134"
type: work
title: "Bound an attempt at its plan and choose each phase's model and reviewer"
status: active
created: "2026-09-24T21:48:22Z"
updated: "2026-09-24T22:07:22Z"
relates_to: ["G-045", "G-046", "G-055", "G-101", "G-108", "G-109", "G-114", "G-135"]
---

## Outcome

The owner can launch an attempt that prepares and stops at a committed plan,
read that plan, and relaunch implementation from it as a fresh attempt with
its own model and effort; every independent review runs on a Grove-owned
reviewer definition with one standard brief; and each attempt records the
model, effort and reviewer definition it actually used.

Owner intent, shaping conversation 2026-09-24: the owner wants the handoff
points inside execution handled more distinctly, with a human look at the
plan before implementation when they want one, as superpowers gates spec and
plan; and wants to be able to consider different models per phase, with
pre-defined agent definitions carrying tailored instructions, as bench's
roles did. The owner asked to understand what that takes before committing
to it, read the attempt analysis below, and agreed to this scope: the
per-launch bound and the per-phase model and reviewer, with role profiles and
a second harness left until [G-135](G-135-run-the-g-108-eval-pair-on-codex.md)
reports. The routine choices below (flag spelling, the reviewer definition's
delivery, the experiment's proposed models and effort) were made in that
session with the owner present. The owner also expects, 2026-09-24, that preparation and review
should run at a higher reasoning effort than implementation, and chose the
experiment's levels (acceptance 5); that is the routing the experiment
tests, not a default this record sets. The flag
spelling, the reviewer definition's delivery and the experiment's target are
proposed design where labelled.

## Constraints

Observed at main `282d282`, from the sixteen Grove-owned attempts under
`.git/grove/attempts/` (event logs and results read in full) and the code:

- **Phase order is uniform and no human reads the plan before code.** Every
  implementing attempt ran context, code reading, a plan in one to three
  minutes, `status=active`, implementation, one fresh reviewer subagent, one
  to three fix rounds by messaging that reviewer, then the review record,
  evidence and `status=review`. Four attempts wrote a plan record (G-103,
  G-115, G-116, G-131) and implemented against it within minutes. The plan
  is read by a person only when a question stops the attempt.
- **The mid-run gate exists and works.** Three attempts stopped on a
  question written during preparation: `G-109.20260923T194926Z` (G-117, a
  plan and a question in 4.4 minutes for $1.62, record left `proposed`),
  `G-108.20260923T194708Z` (G-118) and `G-108.20260924T004935Z` (G-121).
  The owner answered in 17 minutes, 5 hours and 7 minutes; each resumed
  attempt reused the worktree and implemented. Re-entry cost 2 to 4
  minutes and $1 to $2.50 per resumed attempt.
- **The phase after the reviewer is the longest and runs at the highest
  context.** Fix rounds, the review record, evidence and the handoff took
  24 minutes on G-046 (278k to 352k tokens of context), 18 on G-108's first
  attempt, 14 on G-109's second (210k to 254k) and 10 on G-125 (208k to
  249k). No attempt checkpointed between `active` and the handoff.
- **Outcomes.** All ten items reached done; nine were accepted on the first
  candidate; the two feedback rounds were about intent (G-107 "missed the
  intent", G-129 "poke"), which no code reviewer catches and which a plan
  read before implementation is the earliest point that could have caught,
  once in ten. Total spend about $84.
- **Every review ran on an unrecorded, user-level definition.** Each attempt
  dispatched the `reviewer` subagent from `~/.claude/agents/reviewer.md`,
  the owner's startup-team persona with `model: opus`, outside this
  repository and unversioned; the attempt facts record only that an agent
  by that name existed (the init event). Its team instructions did not
  appear in any of 28 reviewer messages checked; the brief each attempt
  wrote dominated, and those sixteen briefs differ in what the reviewer may
  run and whether it gets a commit, a range or a branch. A preview user has
  no such agent. `G-108.20260924T010336Z`, the one attempt whose main model
  differed, shows the split: $7.78 on Fable for the session, $1.61 on Opus
  for the reviewer, from the result event's per-model usage.
- **Runner and harness.** [attempt.go](../internal/attempt/attempt.go)
  composes `claude -p "/grove-work ID --interaction headless"` with budget,
  permission mode, session and an optional `--model`, records the requested
  model in `attempt.json` and the actual one from the init event; there is
  no `--effort`. Claude Code 2.1.281 accepts `--effort`, `--agents
  <json-or-file>` in print mode, and project agents in `.claude/agents/`.
  The adapters ([Claude](../.claude/skills/grove-work/SKILL.md),
  `.agents/skills/grove-work/`) report any argument other than IDs and
  `--interaction MODE` as an error. The board's attempt outcomes are
  "candidate ready", "waiting on question" and "question answered: R
  again", grouped Needs you, Running, Settled
  ([attempts.go](../internal/tui/attempts.go), G-109, G-125).
- **Guide and terms.** The [work guide](../docs/work-execution.md) step 4
  ends at a committed plan and step 5 sets `active` when implementation
  starts; its Lifecycle paragraph says preparation is a fact inside
  `active`, while G-109's first attempt left the record `proposed` with a
  plan, which is the behaviour this record keeps. [G-055](G-055-preparation.md)
  says a technical plan does not wait for a human sign-off unless it needs
  a product choice; a caller bounding an attempt at the plan is the
  mandate's end, not a gate the plan needs, and G-055 may need that
  sentence. The brief selects "no universal plan gate" (kept: the bound is
  per launch) and "route model strength by uncertainty and consequence and
  retain actual configuration per attempt" (this record implements the
  retention and the per-phase choice). [G-101](G-101-attempt-mechanism.md)
  pins attempts to a Claude process; unchanged here.
- **What transfers from elsewhere.** Superpowers (plugin cache 6.4.1)
  front-loads human gates at spec, plan and execution method, then executes
  without asking; its execution guidance says a fully specified plan runs
  well on a mid-tier model and the strongest model earns its cost at the
  final review. Bench (`../bench/meta/team/*.md`) separated a role's
  mandate (`writes: false` for review) from the adapter naming a provider,
  with `model: strong|fast` in the role. The separation transfers; the
  catalogs do not. The owner's own `eng` on Sonnet and `pm` and `reviewer`
  on Opus already encode that routing.
- **The coupling.** Routing implementation to a faster model pays only if
  preparation writes a plan an implementer with less judgment can execute.
  Today's plans are written for the session that implements them. The
  experiment below is the test, and it is the reason the model choice and
  the plan bound are one record.

In scope, proposed design:

1. **A per-launch bound at the plan.** `/grove-work G-030 --until plan`,
   `$grove-work G-030 --until plan` and `grove run G-030 --until plan`,
   with the board's `R` asking. The spelling is chosen: it names the
   artifact the bound ends at and leaves room for a later bound at another
   step. The work guide's step 4
   says once: under that bound, commit the plan and any question,
   checkpoint the runnable continuation in the record's Next naming the
   plan's revision, leave the status as found, and return. The adapters
   accept the argument. The runner carries it in the prompt and in
   `attempt.json`. A rerun with nothing changed returns the same checkpoint
   under the guide's unchanged-wait rule. The relaunch (`R`, or `grove run`
   without the bound) is the owner's approval of the plan, attributable
   through the plan revision the next attempt reads: no plan status, field
   or approval record.
2. **A plan-ready outcome on the board**, beside the existing ones, under
   Needs you, with the continuation as its next action.
3. **Per-phase model and effort.** `grove run` gains `--effort`, passed
   through and recorded with the requested model, so a bounded preparation
   attempt and the implementation attempt that follows it can differ in
   both; the attempt facts add the per-model cost the result event carries,
   so the two attempts show their split.
4. **A Grove-owned reviewer.** A `grove-reviewer` agent definition carried
   in the binary like the guides and written by `init` beside the skills
   (`.claude/agents/grove-reviewer.md`, managed marker, created, kept or
   updated as the skills are). Written by `init` rather than passed with
   `--agents`, because `--agents` reaches only sessions the runner launches,
   while an interactive `/grove-work` session reviews through the same
   definition only if the checkout holds it. Its body is the standard brief: read-only; what it
   receives (checkout, exact range or commit, the record's acceptance and
   constraints, what it may run); what it checks, including the knowledge
   check G-114 added; what it returns (findings with evidence, dispositions
   after fixes, limits). `model: inherit` by default, with its reasoning effort set in
   the definition's frontmatter, where Claude Code reads an agent's model
   and effort, `high` as the owner chose, so the review's effort is
   independent of the implementing session's. Step 6 names it and says what the session
   passes it. The attempt facts record the
   definition's digest or its absence. The distinct name leaves the owner's
   `reviewer` untouched. Codex has no agent-definition equivalent found in
   `codex exec --help`; step 6's "if the harness cannot supply an
   independent reviewer, say so" stands and the Codex adapter is unchanged.
5. **One routing experiment on real work** (acceptance 5).
6. Reconcile [docs/commands.md](../docs/commands.md) (run flags, attempt
   facts, init output) and [docs/board.md](../docs/board.md), and add one
   sentence to [G-055](G-055-preparation.md): a caller may bound an attempt
   at its plan, and that is the mandate's end, not a gate the plan needs.
   The owner chose to change the settled term on 2026-09-24.

Out of scope: role profiles or model defaults in `grove.yaml` (budget and
permission mode stay required per launch, as G-045 chose); a second harness
in the runner (G-101; G-135's report decides whether to shape it); per-task
implementer subagents by default (a plan may call for them); splitting the
guide into phase files (the fixed reading is about 35k tokens against peaks
of 65k to 352k, and every phase needs the authority and staging rules); a
universal plan gate; and bounds at other steps.

## Acceptance

1. The work guide states the bound once, in step 4, and step 6 names the
   reviewer definition and what the session passes it; both adapters accept
   `--until plan` and still reject other text; `grove guide work` prints
   the change; the Evidence records the new guides digest against
   `3c9996e33e42`.
2. `grove run ID --until plan --model M --effort E`, and `R` on the board,
   start an attempt whose command carries the bound, the model and the
   effort; `attempt.json`, `grove attempt` and the attempt screen show the
   three, the init's actual model, and the result's per-model cost; `run`
   without them behaves as before.
3. With the runner's fake provider and in one disposable clone with the real
   one: a bounded attempt ends with a committed current plan record whose
   `work` names the ID, the work record unchanged in status with a
   checkpoint naming the plan revision and the continuation, and no change
   outside records; the board lists it under Needs you with the plan-ready
   outcome; `R` there launches the implementation attempt on the reused
   worktree; a second bounded run with nothing changed returns the same
   checkpoint.
4. `init` writes `.claude/agents/grove-reviewer.md` with the managed marker,
   reports it as `created`, `kept`, `unchanged` or `updated` like the skills,
   and the existing init tests cover it; an attempt's facts show the
   definition's digest, or that none was present.
5. The routing experiment, on the first real assignment after this lands
   whose plan is more than a "no plan needed" note, under a budget the
   owner confirms at assignment (proposed cap $30 for the two attempts,
   against the sixteen-attempt mean of about $5 per attempt): a bounded
   attempt on Claude Opus 5.5 at `xhigh` effort, then an implementation
   attempt on Opus 5.5 at `medium`, reviewed through `grove-reviewer` on
   Opus 5.5 at `high`. The owner chose these on 2026-09-24: one model
   throughout, so the experiment measures effort alone, with the plan
   sufficiency question then bearing on effort rather than on model
   strength. The assignment's record reports, in Evidence, each
   attempt's model, effort, cost and duration, the review's
   findings and rounds, whether the plan sufficed for the implementer or
   what it lacked, and the owner's verdict, read against the sixteen-attempt
   figures above. The experiment's result changes nothing inside this work.
6. `grove check` passes; `gofmt`, `go vet`, the full uncached suite and the
   pseudo-terminal script pass; documentation links resolve.

## Next

The owner can commit this proposal and assign it with `$grove-work G-134`.
Everything is chosen above except two things that cannot be known before
assignment: the experiment's target item, and confirmation of its budget,
which is a paid mandate. The owner chose the target on 2026-09-24:
[G-135](G-135-run-the-g-108-eval-pair-on-codex.md), held until this lands.
Fallbacks if that changes: the follow-on G-129 left open (an attempt owner
that flags a last message promising a continuation), Go work with tests of
the kind the sixteen attempts mostly did, needing a short shaping first; or
the prerequisite-constraint eval case in G-108's Next. [G-110](G-110-external-preview.md) is not a good first
target: it is large and would stop on owner choices during preparation,
which muddies the plan-sufficiency reading. The independent review should read the plan bound
against the amended G-055 and the guide's Lifecycle paragraph. Follow-on, shaped only after this lands: the bounded run as the
first work-row evaluation case for the G-108 suite, about $1.60 a run, with
clone checks on the plan record, the unchanged status and the checkpoint;
role profiles and a second harness wait on G-135's report.
