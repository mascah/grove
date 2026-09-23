---
id: "G-108"
type: work
title: "Establish behavioral evaluations for Grove context and workflows"
status: proposed
created: "2026-09-23T16:05:04Z"
updated: "2026-09-23T16:07:57Z"
relates_to: ["G-078", "G-040", "G-107", "G-110"]
---

## Outcome

Give the owner repeatable evidence of how well Grove's context retrieval and
shared shaping/work workflows guide an agent, and concrete opportunities to
improve them before a small external preview.

Owner intent, shaping conversation 2026-09-23: create an eval suite for context
injection and the corresponding skills, explore improvements, and prepare for
a small external preview. The suite design below is proposed, not selected.

## Constraints

Observed at main `f27444e`: [handoff tests](../internal/handoff/context_test.go)
and [CLI tests](../internal/cli/context_test.go) cover source selection,
revisions, ordering, refusals and read-only behavior. They do not measure an
agent's decisions after reading that context.
[G-078](G-078-g-039-trial-evidence-for-the-int.md) records a headless shaping
run that made product choices instead of persisting a question; the guide was
subsequently tightened, so this is a seed case, not a current reproduction.
[G-040](G-040-portable-bootstrap.md) records discovery trials and interference
from a then-installed global predecessor plugin. Current interference is an
unknown to measure, not an assumed defect.

Proposed first scope: a small local scenario runner, disposable repository
fixtures, observable checks, and a baseline report. Exercise the actual Grove
entrypoints and context commands. Keep software contract tests distinct from
real-model behavioral trials and from human assessment of judgment quality.
No required service, scheduled agents, automatic product merges, or broad eval
platform. Fixture record creation uses explicit absolute project paths in
disposable clones and never consumes this repository's shared counter.

Candidate scenarios: bounded shaping and execution; a missing human product
choice; an unchanged wait; checkpoint continuation after inputs change;
stale candidate/review evidence; and irrelevant or instruction-like linked
material that must not expand authority. Include both successful retrieval and
unnecessary loading. A shortest prompt is not inherently a better prompt.

Record CLI/guide revision, fixture revision, harness/model configuration,
surrounding instructions, trial budget, traces, artifacts, elapsed time, and
reported usage/cost when available. Compare a controlled harness configuration
with the owner's normal configuration to expose instruction interactions.
Use repeated trials and per-scenario outcomes, not a single aggregate pass mark.
Claude and Codex workflow coverage must be explicit; managed attempts currently
use Claude, so do not imply provider parity or build a new runner provider here.

Paid trials require a separately authorized bounded execution mandate. Their
model selection, repeat count, budgets and supported mode matrix are preparation
inputs to agree with the owner. Offline checks must run without paid calls.
[G-107](G-107-current-documentation.md) can proceed independently; pin the
documentation revision under evaluation so subsequent cleanup is comparable.

## Acceptance

1. A documented local command runs the fixtures and checks in isolation, retains
   attributable results, and reports unavailable harnesses and unrun cases as
   such. It cannot silently run paid trials as an ordinary Go test or CI check.
2. Cases above specify expected observable behavior and forbidden side effects.
   Missing-choice and unchanged-wait cases check repository artifacts and
   actions, not merely whether the final answer says the right thing.
3. An authorized bounded baseline exercises real agents with repeated cases,
   retaining successes and failures with their configuration. Simulated or
   scripted responses are clearly separate. Uncovered modes remain explicit.
4. Measures distinguish outcome correctness, authority/scope violations,
   relevant retrieval, unnecessary context, and usable handoff. Human-scored
   criteria have a written rubric; process exit and record validation alone
   cannot pass a behavioral case.
5. A linked report identifies reproducible failure patterns and ranks proposed
   guide/context improvements with evidence and limits. A comparison can be
   rerun against another guide revision; improvement is not claimed without
   comparison evidence. Implementing those product changes is separate work.

## Next

Review and assign this proposal after committing it. Prepare the minimal
scenario/measurement plan and agree the live-trial matrix and resource bounds
with the owner before starting agents. Begin with the historical missing-choice
case and fresh-session discovery, then expand within this bounded scope.
No framework, paid-run budget, or quality threshold has been selected yet.
