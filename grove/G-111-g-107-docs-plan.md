---
id: "G-111"
type: plan
title: "G-107 documentation inventory and ownership plan"
status: current
created: "2026-09-23T16:48:53Z"
updated: "2026-09-23T16:49:34Z"
work: ["G-107"]
---

## Design

Inventory taken at `768efab` on `worktree-G-107` (main at the same commit),
before any document was edited. Word counts are `wc -w` at that commit.

### Ownership

| Fact | One editable owner | Others |
| --- | --- | --- |
| Repository development policy, constraints, CLI invocation here, verification | `AGENTS.md` | README and guides link to it |
| Purpose, constraints, selected direction, phase proposals | `grove/brief.md` | Work records own progress and Next |
| What Grove is, what this build does, command usage | `README.md` | Links to the record model for schema rules |
| Configuration, record schema, validation, lifecycle rules | `docs/record-model.md` | README states usage and links here |
| Work and shaping workflow | `docs/work-execution.md`, `docs/work-shaping.md`, embedded by `guides.go` | Adapters load them; AGENTS.md holds only this repository's policy for them |
| How a fact came to be | The record that decided or delivered it, and Git | Current documents link, never retell |

### Inventory and dispositions

| Surface | Words | Finding | Disposition |
| --- | --- | --- | --- |
| `AGENTS.md` | 1557 | Restart/read-only opening; "The runner contract remains open" though [G-101](G-101-attempt-mechanism.md) is accepted and `run` merged; implementation history per work ID beside policy; "The hook also scrubs them; keep both layers" though `ff43559` removed the pre-push hook | **Reconciled**: policy only, grouped as orientation, workflow routing, constraints, verification; each constraint below kept as a rule with the record that owns its detail |
| `grove/brief.md` | 1592 | Adoption milestone as current; Done "will mean"; branch-local Done still permitted; installed `grove` as the predecessor; shipped features as future; observed state at `c9904ea`; sequence routes to the closed G-036 roadmap | **Reconciled**: purpose with the owner's preview audience; G-036 closure with its nullsec substitution; shipped state as fact; sequence names G-108, G-109, G-110 as proposed, unordered. Headings `Current view and TUI` and `Suggested sequence` kept, since records link them |
| `README.md` | 6015 | Adoption as next milestone; "Grove launches no agent"; "Agent execution remains future work"; pre-push "runs `go test`" (removed in `ff43559`); "Resume this conversation" prompt; lifecycle, allocation and schema rules repeated from the record model | **Reconciled**: new opening (what, what this build does, where next); command reference retained; repeated schema rules cut to links; stale claims corrected; resume prompt **retired** |
| `docs/record-model.md` | 5454 | Titled "Starter record model"; acceptance quotes and dates, nullsec W-032 evidence, "Why these defaults", "Operational records", "First dogfooding boundary" narrate history beside the contract | **Reconciled**: current contract only; acceptance history left to the records it cites; anchors `identity-and-dates`, `on-disk-contract`, `folders-and-files`, `identity-and-placement-apart-from-classification` kept, since records link them; "Attempts are not records" **relocated** into a short section |
| `docs/work-execution.md` | 3821 | Current | **Retained** |
| `docs/work-shaping.md` | 2140 | "Grove starts no agent and schedules nothing" contradicts `grove run` | **Reconciled**: that sentence only, stating that `grove run` starts only work attempts; no workflow step changes |
| `guides.go` | – | Embeds the two guides | **Retained** |
| `.claude/skills/*/SKILL.md`, `.agents/skills/*/SKILL.md`, `agents/openai.yaml` | – | Thin: load AGENTS.md and the guide | **Retained** unchanged |
| `docs/prompts/W-006-W-008-implementation.txt`, `W-009-…`, `W-010-W-011-…` | 935 | Spent handoffs naming `docs/restart-brief.md`, `grove/work/…`, typed IDs; two prescribe race suites policy forbids | **Retired** (deleted). Git keeps them at `768efab`; G-023's two links become historical, as G-107 accepts for records |
| `justfile`, `lefthook.yml`, `.github/` | – | Accurate; `lefthook.yml` comment still names G-089 | **Retained** |
| `grove/` records | – | History | Not edited, per G-107 |

Constraints AGENTS.md keeps as rules: Git subprocesses only through
`repo.Command`; IDs only from `new` with the shared counter and write lock,
fixtures in a disposable clone by absolute `--project`; one `git cat-file`
process for all branches, history read only while a card is open; source
freshness and exact targeting on the board and `workspace`; terminal escaping
before glamour and the style filter after it; explicit subcommands
noninteractive; lifecycle authority (`review` with a candidate, `done` only on
the target after the merge, no backfilled candidates); staged retrieval and
`context` as facts; schema 3 only, no pre-release compatibility; G-069 for old
references; verification and race-detector policy; installed `grove` lagging
the checkout; siblings and the archive as evidence only.

Found outside the documentation scope, flagged for follow-up, not changed:
`grove --help` still says the board "Reads only", though it writes through
the prompted review actions and starts and stops attempts.

## Steps

1. Done: plan `0c42435`, G-107 active `ef911b1`.
2. Done: `d6cc1df`. The README also dropped the "Local restart" history and
   the unselected web UI remark. The archive's location moved to the brief.
3. Done: see G-107 Evidence.
4. Done: see G-107 Evidence.
5. Done: review [G-112](G-112-g-107-review.md), fixes `1398458` and
   `6a9f146`; handed into Review.
