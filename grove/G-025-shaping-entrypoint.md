---
id: "G-025"
type: work
title: "Shape project work through reusable agent instructions"
status: done
created: "2026-09-19T21:20:41Z"
updated: "2026-09-21T04:09:41Z"
kind: feature
priority: 1
size: small
relates_to: ["G-023", "G-017", "G-036", "G-037", "G-035"]
formerly: "W-011"
---

## Outcome

Discuss an idea or selected outcome in an interactive Claude or Codex session
and have Grove-owned instructions maintain useful proposed work and related
knowledge, without bespoke prompts or predecessor skills. This is the first
assignment in [G-036](G-036-interactive-adoption.md)'s adoption sequence.

## Selected scope

Deliver thin repository-local `grove-shape` adapters loading one shared
`docs/work-shaping.md` guide. Start open exploration with the brief and relevant
knowledge; when work is selected, use staged context and inspect actual code,
existing proposals, branches and worktrees before creating duplicates.

Separate intent, observed evidence, proposed design and decisions with actual
authority. Capture outcomes, constraints, meaningful acceptance and Next. Persist
real unresolved human questions; investigate routine technical unknowns rather
than escalating them. Creating a proposal does not assign implementation.

Use the current work/question/decision schema and supported CLI operations.
Terms and structured artifacts follow in G-037: link existing ordinary documents
when useful, but do not invent unsupported record types or metadata. Context is
facts, not permission. Shared guidance owns how and when to use the CLI;
AGENTS.md owns this repository's invocation and development policy.

Document the same guide's bounded headless adaptation: explicit mandate and
human availability, durable questions/waits, isolated reviewable proposal branch,
supporting knowledge identified for selective integration, no self-authorized
implementation or merge. Do not add a launcher, timer or process supervision.

## Acceptance

1. Short Claude/Codex entrypoints load the same guide; verify discovery in fresh
   available harness sessions and state unexercised behavior explicitly.
2. A real requirements conversation produces or refines useful work and any
   justified question/decision, with attributable authority and a concrete Next.
   Do not manufacture records to meet a checklist.
3. Exercise overlapping proposals, work on another branch, a stale revision,
   unanswered human choice and headless wait in disposable fixtures where needed.
   Distinguish observed agent behavior from source checks and simulations.
4. Records validate and links resolve. Report the actual checkout/publication
   location using today's board behavior; do not promise G-042's current view.
5. No implementation, sibling migration, unsupported schema, status promotion,
   agent launch or automatic merge is a side effect of shaping.

## Evidence and next

Checkpoint 2026-09-20. Assigned alone (`/grove-work G-025`, interactive);
branch `worktree-W-011` in `.claude/worktrees/W-011`, base `main` `70f19c5`.
The current plan is
[G-049-shaping-entrypoint-plan.md](G-049-shaping-entrypoint-plan.md); the
G-025 portion of the old
[shared handoff plan](G-027-agent-handoffs-plan.md) is
superseded. The [roadmap](G-047-adoption-roadmap-plan.md) holds
coordination only.

Delivered: [the shaping guide](../docs/work-shaping.md), the `grove-shape`
adapters for Claude and Codex, and the AGENTS.md and README sections. On the
owner's instruction of 2026-09-20 the brief no longer tracks progress or a next
action, so this record's Next is the only one. [Evidence](G-050-shaping-review.md)
separates source checks, observed `claude -p` and `codex exec` trials in
disposable clones, CLI simulation, and what was not exercised.

| Acceptance | State |
| --- | --- |
| 1 | Met for explicit headless invocation in both harnesses; interactive typing of the command is unexercised and stated. |
| 2 | Not exercised through `grove-shape`; closed on the owner's verdict below, with that limit. |
| 3 | Overlap, other-branch work, unanswered choice, headless wait and unchanged rerun observed; stale revision simulated with the CLI only. |
| 4 | `check` OK and links resolve; the guide's return reports checkout, branch and commit with today's checkout-scoped board. |
| 5 | Observed in both trials: status stayed `proposed`, no code, launch or merge. |

**Owner's verdict, 2026-09-20 (interactive, in the implementing session):** merge
and close G-025. The owner's recent requirements conversation already produced
the work now on `main` (G-036 to G-046) without this guide, and there is no
fresh idea to shape until that work advances. The owner will run
`/grove-shape G-037` from the main checkout, inspect the result, and adjust the
guide through further dogfooding if needed. Done here therefore means delivered,
independently reviewed and accepted by the owner on the headless evidence; a
real interactive shaping conversation remains unobserved.

**Next:** none for G-025. Findings from the owner's `/grove-shape G-037` run
belong to new or related work, not to reopening this record by default.
