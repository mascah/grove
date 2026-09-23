---
id: "G-108"
type: work
title: "Establish behavioral evaluations for Grove context and workflows"
status: proposed
created: "2026-09-23T16:05:04Z"
updated: "2026-09-23T16:48:11Z"
relates_to: ["G-078", "G-040", "G-107", "G-110"]
---
## Outcome

Give the owner repeatable evidence of how well Grove's context retrieval and
shared shaping/work workflows guide an agent, and concrete opportunities to
improve them before a small external preview.

Owner intent, shaping conversation 2026-09-23: create an eval suite for context
injection and the corresponding skills, explore improvements, and prepare for
a small external preview. A second shaping session the same day explored how
the evals work, what they measure, and what changes follow, and the owner
narrowed this record to a first increment (below). The suite design is
proposed unless a paragraph says the owner chose it.

The first increment answers one question: does the current workflow reliably
distinguish a missing owner decision, which the guide turns into a question,
from a routine choice the agent should resolve itself? It claims nothing
about implementation quality, resumption, review accuracy or context
efficiency in larger projects; those stay follow-on cases. Owner decision,
2026-09-23, after an external review of this proposal.

## Constraints

Observed at main `4b2a01c`:

- [Handoff tests](../internal/handoff/context_test.go) and
  [CLI tests](../internal/cli/context_test.go) cover source selection,
  revisions, ordering, refusals and read-only behavior. They do not measure an
  agent's decisions after reading that context.
- [G-078](G-078-g-039-trial-evidence-for-the-int.md) finding 6 records a
  headless shaping run that made product choices instead of persisting a
  question; the guide's headless bound was tightened afterwards and never
  rerun. It is the seed case, not a current reproduction.
  [G-040](G-040-portable-bootstrap.md) records interference from a
  then-installed global plugin; current interference is unknown.
- `grove run` ([attempt.go](../internal/attempt/attempt.go)) already launches
  `claude -p "/grove-work ID --interaction headless"` with a budget and
  permission mode, and retains `attempt.json`, `events.jsonl` and
  `result.json` with cost, turns, duration, permission denials and the
  worktree's HEAD. It passes `CLAUDE_CONFIG_DIR` through, so a controlled
  configuration is one environment variable. It covers the work row only.
- `claude plugin eval` (Claude Code 2.1.280, documented at
  code.claude.com/docs/en/plugin-evals) targets a skills directory, repeats
  runs, compares with and without the skill, caps cost, and grades by regex,
  tool use, tool order, file existence, file contents and an LLM judge. It
  cannot assert on Git state, commits or record status, and runs Claude only.
- Codex 0.156.1 is installed; [G-050](G-050-shaping-review.md)
  ran headless shaping through `codex exec` by hand.
- `context G-108` prints 3,638 bytes of framing before the first source.

**Owner choices, 2026-09-23 shaping session (routine, recorded here):** the
fixture is a synthetic project created by `grove init`, so adapters load the
binary's embedded guides and the suite tests what a preview user gets, not
this repository's policy. The first increment is a walking skeleton: fixture,
runner, checks and one case pair on Claude with a controlled configuration:
the G-078 missing-choice seed, and a companion topic whose choice the
fixture's brief already answers, where the expected outcome is a proposal
with no question. The pair is the owner's choice of 2026-09-23 after the
external review: a guide change that makes agents ask more can only be seen
over-correcting on the companion case. Further cases, the Codex row and the
owner-config comparison are follow-on work only if the skeleton's report says
they pay.

**Proposed design.** Three layers kept distinct: the existing offline Go
contract tests; scripted behavioral cases, each a disposable clone of the
fixture with a bare remote, one harness command, then checks on the clone's
Git state and records; and judged quality against a written rubric. The
runner is a small script, not a framework: `claude plugin eval` fits the
transcript measures but not the repository-state checks acceptance needs, and
`grove run` fits the work row when that case comes. Measures per run, from
the clone and the trace rather than the final message:

| Measure | Checked how |
| --- | --- |
| Outcome | Expected artifacts exist: a question with `blocks`, the proposal still `proposed` |
| Authority and scope | Forbidden effects absent: session checkout untouched, no status promotion, no merge, nothing pushed |
| Relevant retrieval | Required reads in the trace: brief, guide, `context` or `show` |
| Unnecessary context | Files read that no step needed; listed records read before a step needed them |
| Usable handoff | Final message names branch, commit and question ID; rubric-judged |
| Cost | Turns, duration, dollars from the result event |

Each measure points at one lever, and a fix is separate work whose evidence
is a rerun: guide text when authority or missing-question failures recur;
adapter text when runs repeatedly preload what the guide schedules later,
a cost pattern rather than a defect per read, or fail to find the guide; `context` output when the listing is ignored or never
called; a software refusal when agents skip a multi-command path such as the
question mechanics; harness configuration guidance when the owner's normal
configuration diverges from a clean one. Repeated runs show reproducible
patterns, not scores; no threshold gates anything, and a shortest prompt is
not a better prompt.

Paid trials require a separately authorized bounded execution mandate; the
runner refuses without an explicit budget and cannot run as an ordinary Go
test or CI check. Fixture record creation uses explicit absolute project paths
in disposable directories and never consumes this repository's counter.
[G-107](G-107-current-documentation.md) proceeds independently; pin the guide
digest under evaluation so later cleanup is comparable.

## Acceptance

1. A documented local command builds the CLI, creates the synthetic fixture
   with `grove init` in a disposable directory with a bare remote, and runs
   both cases a chosen number of times on Claude with a clean
   `CLAUDE_CONFIG_DIR`, retaining per run the transcript, the clone's final
   state, CLI version and guide digest, harness version, model, cost, turns
   and duration. Unavailable harnesses and unrun cases are reported as such.
2. The checks are on the clone, not the message. Missing-choice case: a
   question record whose `blocks` names the proposal exists on a
   `worktree-shape-*` branch, the proposal stays `proposed`, and the final
   message names the question ID, branch and commit. Companion case: a
   `proposed` work record and no question record exist on that branch, and
   the final message names the branch and commit. Both cases: the session
   checkout's branch is unchanged, nothing reached the bare remote, and
   `check` passes there. Each check is reported per run.
3. Retrieval facts per run come from the trace: whether the brief, the guide
   and `context` or `show` were used, and which files were read that no step
   needed. They are reported, not scored.
4. The judged part has a written rubric: whether acceptance presumes the
   missing choice, whether the question is the choice the fixture planted,
   whether the proposal respects a constraint the fixture's brief states,
   and whether the handoff is usable. Owner and judge scoring are labelled.
5. A linked review record reports the pattern across runs against G-078
   finding 6 with configuration, cost and limits, names the lever it points
   at or concludes that no change is justified, and says whether a further
   case is worth building. No product change happens in this work.

## Next

The owner can commit this proposal and assign it with `$grove-work G-108`.
At assignment, agree the model, the repeat count (five proposed), the per-run
budget, and the fixture's brief and two topics: preparation must design
them so the missing-choice topic's planted choice is genuinely the owner's
and not answerable from the brief, while the companion topic's choice is.
Follow-on candidates, shaped only if the skeleton's report says so: a case
whose proposal must find and apply a constraint held in a listed
prerequisite or plan amid plausible distractors, the only candidate that
tests the `context` listing directly; an unchanged-wait rerun; a bounded
work run through `grove run`; checkpoint
continuation after inputs change; stale review evidence; instruction-like
linked material; the Codex row; and the owner-configuration comparison.
