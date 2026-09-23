---
id: "G-107"
type: work
title: "Reconcile current documentation and give each fact one owner"
status: active
created: "2026-09-23T15:51:37Z"
updated: "2026-09-23T17:24:43Z"
relates_to: ["G-036", "G-108", "G-109", "G-110"]
candidate: "71a650ec59c7e0b881a295741ede86bda2c84d09"
---

## Outcome

Give fresh agent sessions and readers of this repository an accurate,
navigable account of Grove as it works today, with each fact owned by one
editable document, and without carrying the completed adoption push as current
instructions.

Owner intent, shaping conversation 2026-09-23: Grove is a "minimum functional
product," useful for its own development, and nullsec has cut over. The owner
requested documentation reconciliation to avoid stale conversation context and
selected "Prepare for a small external preview" as the next-phase audience.
"Preview" was the assistant's shorthand for initial use outside the owner's
projects, not a selected release channel or compatibility promise; the owner
asked for that clarification in the same conversation. On 2026-09-23 the owner
narrowed this work, in a review of the first draft, to reconciliation and
ownership: newcomer onboarding and any preview-readiness judgment belong to
[G-110](G-110-external-preview.md), and agent behavioral trials to
[G-108](G-108-workflow-evals.md).

## Scope and constraints

Observed at main `4b2a01c6d301557b82ad989aabb80d826ef38a55`:

- [AGENTS.md](../AGENTS.md) opens with restart/read-only framing, says "The
  runner contract remains open" although [G-101](G-101-attempt-mechanism.md)
  is accepted and `run` is merged, and carries accumulated implementation
  history beside its policy.
- [The brief](brief.md) describes shipped features as future work, the installed
  CLI as the predecessor, and branch-local Done as still permitted. Its
  "Observed state" is dated at `c9904ea`, and its "Suggested sequence" points at
  the adoption roadmap that [G-036](G-036-interactive-adoption.md) closed; G-036
  records the owner's closure, including the explicit substitution of Grove
  dogfooding for a real nullsec change. Do not erase that evidence distinction.
- [README.md](../README.md) calls adoption the next milestone and ends with a
  "Resume this conversation" prompt. Its roughly 6,000 words mix onboarding,
  command reference, and history. [The record model](../docs/record-model.md)
  mixes the implemented contract with starter-era design narrative.
- `docs/prompts/` retains three spent assignments naming removed paths and typed
  IDs; two prescribe full race suites, which repository policy now forbids.
- [The work guide](../docs/work-execution.md) and
  [shaping guide](../docs/work-shaping.md) are the shared workflow sources,
  embedded through [guides.go](../guides.go). Repository adapters load those
  sources; generated adapters load the binary's guides.

Proposed approach: inventory current instruction and documentation entrypoints,
then give each fact one editable owner. AGENTS.md owns durable repository policy;
the brief owns compact purpose, constraints and selected direction; README owns
the product introduction and routes to reference; the record model owns the
current schema; shared guides own workflow. Keep adapters thin. The inventory
is a plan record linked by `work`, written before any document is edited.

In scope: documents outside `grove/`, and the brief. Records, including done
work and dated reviews, are history and are not edited: their identity, paths,
acceptance limits and verdicts stay as written. Retire spent prompts from the
live instruction surface; Git keeps their provenance.

Retain consequential constraints, including repository-safe Git subprocesses,
shared allocation/write locking, source freshness, terminal escaping, lifecycle
authority, staged retrieval, and verification policy. Removing the narrative of
how a constraint arose must not remove the constraint. Do not infer a need to
load every implementation record into every session.

This is documentation reconciliation. Behavior changes to the CLI, skills or
context assembly, an eval runner, Attempts redesign, release automation, Pages
publication, GitHub settings, and sibling-repository changes are separate work.
Flag newly found behavioral contradictions for follow-up rather than silently
changing the workflow contract. No release compatibility promise is selected.

## Acceptance

1. A plan record inventories current documentation, agent entrypoints, and
   spent prompts, giving each a retained, relocated, reconciled, or retired
   disposition. Current entrypoints no longer present completed adoption as
   pending, the predecessor as installed, or the runner contract as open.
   Claims match the inspected revision and commands.
2. The brief states product purpose and the owner's selected preview audience,
   and its sequence section names [G-108](G-108-workflow-evals.md),
   [G-109](G-109-attempts-usability.md) and
   [G-110](G-110-external-preview.md) as the proposals of that phase, their
   order labelled proposed until the owner selects one. Next actions remain in
   work records; no second editable roadmap or status account is introduced.
   Historical verdicts remain attributable.
3. A fresh session given only AGENTS.md can say how to shape work, how to
   execute assigned work, and how to retrieve context in stages, without being
   routed through the old adoption roadmap. AGENTS.md states each rule in a
   line with the record that owns its reasoning, and routes to the documents
   that own everything else; it does not explain. This is a reading check,
   not a behavioral eval: one bounded interactive session, its harness and
   revision recorded, launched only with the assignment's mandate.
4. The README is an introduction and a router, not a manual (owner feedback,
   2026-09-23). It says what Grove is, how to run it, and where each subject
   is owned, and links there, without requiring old work IDs. It gives a
   command overview, not a paragraph-per-command retelling of `grove --help`.
   The board's key-by-key description and other reference-grade text move to
   a document that owns them, or are cut where another owner already holds
   them. There is no word cap: the test is that a reader new to the
   repository finds the place that owns each subject without reading the
   rest. Whether that serves an external newcomer is G-110's judgment, not
   this work's.
5. Local links and command examples resolve or are explicitly labelled
   historical. `go run ./cmd/grove check` passes. Shared guide ownership and thin
   adapter behavior remain consistent, and no safeguarded behavior or authority
   boundary changes merely to shorten text. Report context reduction as evidence,
   not as a substitute for correctness.

## Evidence

Executed headless via `/grove-work G-107 --interaction headless` in Claude
Code (Opus 5.5), 2026-09-23. Branch `worktree-G-107` in
`.claude/worktrees/worktree-G-107`, based on main `768efab`, which held this
record at `sha256:c411c9a8…`. The candidate is the commit that adds this
Evidence (the `candidate` field). Plan: [G-111](G-111-g-107-docs-plan.md).
Commits: plan `0c42435`, active `ef911b1`, reconciliation `d6cc1df`, review
fixes `1398458` and `6a9f146`.

Against the acceptance:

1. G-111 inventories every documentation file, agent entrypoint and spent
   prompt at `768efab`, each with a disposition. The spent prompts in
   `docs/prompts/` are deleted and remain in Git at `768efab`. As a result,
   G-023's two links to them no longer resolve, which leaves them as
   history. Current entrypoints no longer say that adoption is pending, that
   the predecessor is installed, that the runner contract is open, that
   "Grove launches no agent", or that agent execution is future work. The
   same holds for the removed pre-push hook: `ff43559` removed it, and the
   entrypoints no longer describe it.
2. The brief now states:
   - the purpose;
   - G-036's closure, with the nullsec substitution attributed to the owner;
   - Keyborg, still selected and without a record;
   - the preview audience, marked as not a release channel or promise.

   Its "Suggested sequence" names G-108, G-109 and G-110, says the owner has
   selected no order, and sends G-047 and G-048 to history. It keeps no
   progress account. The headings that records link to are kept.
3. Reading check: a fresh `general-purpose` subagent in this session read
   only AGENTS.md at `d6cc1df` and answered how to shape, execute and
   retrieve context in stages. It went through no roadmap and never mentioned
   G-036 or G-047. It asked what "the headless form" meant, and `1398458`
   answered that.
   Limit: its injected system context still held the session-start copy of
   AGENTS.md (`768efab`). It reports answering from the file on disk. This
   was a subagent, not a separate interactive session.
4. The README now opens with what Grove is, a list of what this build does,
   and where to go next, with no work IDs. Whether that serves a newcomer is
   G-110's judgment.
5. Local links and anchors were checked with a script over AGENTS.md, the
   README, `docs/*.md`, the brief and G-111: no problems. Every anchor that a
   record links into the record model or the brief resolves.

   The non-writing README examples all exit 0 with a built binary at
   `d6cc1df`: `list`, `list --status active --status review`, `show`,
   `show --json`, `brief`, `versions`, `versions --json`, `guide work`,
   `guide shape`, `version`, `attempts`, `context` and `--help`. Writing
   commands and `workspace` were not run; their README text is unchanged.

   The adapters are unchanged. Of the shared guides, only one sentence in
   `docs/work-shaping.md` changed: Grove starts no *shaping* agent, and
   `grove run` starts only work attempts. No workflow step changed.

Verification at `6a9f146` content:

- `go run ./cmd/grove check`: OK.
- `go vet ./...`: clean.
- `gofmt -l .`: empty.
- `go test -count=1 -timeout 120s ./...`: all packages pass.
  `attempt`, `cli`, `tui` and `versions` take 6–13 s. That is inherited;
  documentation cannot change it.

Context size (`wc -w`), from `768efab` to candidate:

| File | Before | After |
| --- | --- | --- |
| AGENTS.md | 1557 | 1323 |
| README | 6015 | 5485 |
| Record model | 5454 | 4682 |
| Brief | 1592 | 1566 |
| Prompts | 935 | 0 |
| **Total**, with the shaping guide | 17693 | 15206 |

This is evidence of reduction, not of correctness.

Review: [G-112](G-112-g-107-review.md), two rounds, with an independent
agent reviewer. Its seven findings were fixed and confirmed. One round-2
wording fix was self-checked. No findings are open.

Flagged for follow-up, outside this documentation scope: `grove --help`
says the board "Reads only", though it approves, gives feedback, integrates,
and starts and stops attempts behind prompts.

## Next

Returned to active by the owner's feedback below. The next pass starts from
this branch, keeps the ownership and reconciliation work already on it, and
restructures against the revised acceptance 3 and 4: the README becomes an
introduction and router, the board manual and other reference text move to
an owning document, and AGENTS.md keeps its rules as one line each. G-111's
dispositions are revised where they change. `grove --help` improvements and
a board help overlay are separate work; propose them, do not do them here.
Then hand off into Review again with a new candidate.

Feedback on candidate 71a650e, 2026-09-23: Missed the intent. The record asked for reconciliation and got it, but never said that the README and AGENTS.md are routers, not books, so they still are: the README is 34.6k characters, a prose retelling of grove --help plus a full board manual. There is no word cap; the test is that the right content is in the right place and someone else can understand it. The README keeps a command overview but does not retell --help, which needs its own work and does not replace the README. The board's key-by-key manual leaves the README. AGENTS.md matters less. Acceptance revised in the next commit.
