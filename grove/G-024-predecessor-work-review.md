---
id: "G-024"
type: review
title: "Predecessor `/work` review for G-023"
status: current
formerly: "docs/reviews/2026-09-19-predecessor-work.md"
work: ["G-023"]
updated: "2026-09-21T21:11:16Z"
---

# Predecessor `/work` review for G-023

Inspected 2026-09-19 at sibling skills revision
`ec87bb2eb19a350634c4976cc5d215a57deb58de`. This is observed source behavior
and a proposed adaptation, not adoption of the predecessor workflow. No skill,
runner, or agent was launched, and the sibling repository was left unchanged.

## Sources and evidence limits

The primary source is the actual
[work skill](../../skills/skills/work/SKILL.md), especially preparation
(lines 10–15), controller execution (17–46), and handoff/recovery/closure (48–54).
Also inspected its [discipline](../../skills/references/discipline.md),
[planning](../../skills/references/planning.md),
[implementer](../../skills/skills/work/implementer.md),
[reviewer](../../skills/skills/work/reviewer.md), and
[close](../../skills/skills/close/SKILL.md) instructions.

After reading the sibling's AGENTS.md, retrieved its project knowledge through
its installed CLI from that repository: `grove status`, `grove find controller`,
`grove find export`, and `grove context --work W-016 --phase shape --budget 12000`
and the equivalent W-002 context with budget 10000. Those work IDs belong to
the predecessor. Neither context response reported required omissions.

Inspected [launch.py](../../skills/cli/grove/launch.py), its
[tests](../../skills/cli/tests/test_launch.py), and the actual
[claude-p adapter](../../skills/scripts/adapters/claude-p.sh).
Source inspection establishes their authored behavior; their tests and the
headless runtime were not executed in this review. Prior trial reports retrieved
through context are inherited evidence. In particular, W-016 still has a pending
user trial; persistence of its `/goal` instruction across compaction is unproven.
Sibling links refer to evolving files; the revision above identifies this review.

## What `/work` actually did

1. **Prepare, including missing plans.** Inspect ownership and dependencies,
   fetch each work unit's context, resolve required omissions, investigate code,
   and prepare or repair a plan. Choose delegation separately from durable
   execution. Structural readiness does not prove correctness or authorization.
   Ask only for human preferences or mandate changes; continue independent work.
2. **Own execution.** Use an isolated worktree and acquire claims. Small work
   runs in-session. Bounded/large work uses a controller that dispatches
   implementers and reviewers rather than writing code or running tests itself.
   Its Sonnet/Opus assignments, size vocabulary, and Agent/SendMessage APIs are
   explicit predecessor policies, not restart defaults.
3. **Bound review and retries.** Tie briefs, reports, diffs, and reviews to an
   attempt and tested revision; consume each result once. Fix rounds 1–3 reuse
   the implementer, 4–5 escalate to a fresh implementer; then park unresolved
   findings. Reconcile or stop the previous writer before starting another.
   Workers own checks through completion, or report the command owner and wake
   condition. Duplicate notifications must not cause duplicate work.
4. **Verify the connected result.** A final whole-branch review maps affected
   sequences to evidence, missing checks, defects, or reserved human judgment.
   Passing component checks and uninspected screenshots do not prove an end-user
   sequence. Run additional verification for identified evidence gaps.
5. **Recover and reconcile.** Keep a compact checkpoint in the ledger plus work
   Next, or Next/Evidence alone for small work. Resume surviving workers; record
   lost attempts as interrupted and preserve partial files before replacement.
   Keep per-unit acceptance and joint checks for batches. Close completed units
   in dependency order with evidence, project-knowledge reconciliation, delivery
   checks, and an explicit integration handoff. Human acceptance stays human.

## Three different handoffs already existed

**An instruction generator:** `grove launch <ids>` calls batch selection and
prints a single `/goal …` message. It launches no process. It refuses external
blockers but permits preparation gaps because `/work` prepares missing plans.
The message defines both finished and parked terminal conditions, controller
discipline, and checkpoint recovery. This directly addresses the repeated
prompt-writing problem behind G-023; do not confuse its name with a runner.

**A frozen execution contract:** `/work` describes `export` for prepared work
and `reconcile` for its result. This carries source/checkout identity, acceptance,
bounds, and tested revision; changed exported inputs prevent reconciliation.
Unlike the instruction generator, export refuses unprepared work. These are
different stages, not contradictory readiness rules. Acceptance remains judgment.

**An optional process adapter:** `scripts/adapters/claude-p.sh` consumes a
contract and attempt directory, creates or reuses a worktree, invokes `claude -p`
with the plugin, and asks for a result document. Thus a concrete predecessor
adapter already exists; future runner investigation must start from it as well
as the skill. It is not evidence that the restart has a working runner.

Do not copy this adapter unchanged: it uses `--dangerously-skip-permissions`,
reuses an existing worktree directory without validating its identity in the
script, and can fill in a missing tested HEAD from the final branch HEAD without
proving tests ran there. Its optional test command ignores failure status and
appends only the last output line as evidence. These are visible script behaviors,
not end-to-end failure demonstrations; surrounding driver guarantees were not
audited here. They warrant specific investigation before any launcher reuse.

## Proposed responsibility mapping for the restart

| Responsibility | G-023 treatment |
| --- | --- |
| Current source context, dependency order, missing-input diagnostics | Deterministic CLI facts; explicit checkout and record revisions. |
| Prepare/repair a plan within an authorized outcome | Retain in judgment instructions; a missing plan can be a preparation task, not an automatic refusal of every work invocation. |
| Outcome, constraints, terminal conditions, technical autonomy | Retain in the shared guide; return a concrete completion or blocked/awaiting-judgment handoff. “Parked” need not become a record status. |
| Isolation, duplicate-work avoidance, integration awareness | Adapt to actual Git/worktree evidence; no invented restart claim commands or proof from `done`. |
| Implementer/reviewer coordination | Choose and document supported harness behavior during G-023 preparation. Do not silently imply controller parity from a prompt template. |
| Evidence tied to revisions, connected checks, human judgment | Retain; use current record bodies and linked artifacts. |
| Bounded retries and interruption recovery | Specify for the selected manual execution path; checkpoint Next/Evidence. Durable multi-agent attempt machinery is separate scope unless explicitly selected. |
| Knowledge reconciliation and integration handoff | Retain responsibility using restart records/docs/check; do not invoke predecessor close/archive/delivery commands. |
| Frozen contracts, result reconciliation, supervised `claude -p` | Defer implementation to a separately specified runner; retain provenance requirements as design input. |
| Fixed models, size-triggered controller, `/goal`, claims, release/history schema | Do not import automatically. Verify harness capabilities and deliberately select any equivalent. |

The current [execution guide](../docs/work-execution.md) covers only part of this
behavior. It is a useful baseline, not a replacement proven equivalent to the
old `/work`. [G-023](G-023-work-handoffs.md) must explain what
the first supported path retains and what it defers, then dogfood that path.
An eventual Kanban Implement action should consume the same assignment rather
than maintaining another prompt, after a separate launch/lifecycle contract.
