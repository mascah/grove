---
id: "G-154"
type: work
title: "Evaluate whether agents find a constraint held in a listed record or in the code they touch"
status: active
created: "2026-09-25T19:15:12Z"
updated: "2026-09-25T19:29:49Z"
relates_to: ["G-108", "G-122", "G-135", "G-153"]
kind: investigation
---

## Outcome

The owner has repeatable evidence of whether an agent shaping in an adopting
project finds and applies a constraint that the selected record does not
hold: one kept in a listed prerequisite or plan among plausible distractors,
and one kept in a record that describes the file the proposal must touch;
and of whether [G-153](G-153-search-and-code-links.md)'s search changes
that, at what cost in reading.

Owner intent, shaping conversation 2026-09-25: evaluate body search and
code-to-record links against Grove's existing retrieval before paying for
an index or a code graph. The first case is the one
[G-122](G-122-g-108-baseline-runs-the-missing.md) named as the only case
that would show the `context` listing, and G-108's Next lists it first
among follow-on cases.

Owner choices, 2026-09-25, after the records were shaped: this work runs
before [G-153](G-153-search-and-code-links.md) starts, so its `without`
row is the baseline G-153's preparation reads; and it runs on Claude only.
The Codex row and `gpt-6-astra` are outside this work, whatever a later
mandate says about model, runs, budget and permission mode.

## Constraints

Observed at main `001b271`, 2026-09-25:

- The runner ([evals/run.py](../evals/run.py),
  [evals/README.md](../evals/README.md)) runs headless shaping on a fixture
  created by `grove init`, checks the clone, and reports retrieval facts
  from the trace, never scored: whether the brief, the guide and `context`
  or `show` were used, every file read, and the reads no step needed. The
  fixture's only knowledge is its brief, so today no run can show the
  listing being used or ignored (G-122, Disposition).
- The G-108 pair costs about $3 for ten runs on Claude Opus 5.5 and is the
  regression check for a guide edit ([G-114](G-114-capture-and-reuse-terms-question.md)
  Evidence). The Codex row exists ([G-135](G-135-run-the-g-108-eval-pair-on-codex.md))
  with plan-window caps, and is not used here: the runner's `--harness`
  defaults to `claude` and runs Codex only when a command names it, and
  `gpt-6-astra` never runs unless the owner names it. Model, repeat count,
  budget and permission mode are the owner's explicit answer, never a default
  ([G-118](G-118-what-mandate-should-the-g-108-pa.md),
  [G-141](G-141-never-run-gpt-6-astra-unless-the.md)); no paid run is ever
  part of `go test` or CI.
- The runner pins the guides digest `grove version` prints, and builds the
  CLI from the checkout it runs in. A with-and-without comparison is
  therefore two checkouts: one before G-153's guide text and command land,
  one after. The "without" row needs nothing from G-153.
- The work row through `grove run` is unbuilt (G-108 Next). Both cases fit
  shaping: a topic whose proposal must respect a constraint, and a topic
  whose implementation would change a file a fixture record describes.
- Mex's published evaluation ([evaluate/RESULTS.md](https://github.com/mex-memory/mex/blob/57f565cd595e803816b413451f9f87906e218747/evaluate/RESULTS.md))
  compared ordinary file tools with a forced-first retrieval call on 12
  tasks and reports 54.5% fewer new tokens and 7 of 12 correct against 6 of
  12, on one model, and says those limits plainly. Borrowed: the paired
  design and the honesty; not the forced-first call, which the guides do not
  impose.

**Proposed design.** Two cases on the existing fixture, each with a
`without` row at the digest before G-153 and a `with` row after, the same
run count, on Claude only:

- `listed-constraint`: a `done` prerequisite work record, or a `current`
  plan, holds a constraint the topic must respect, and two records with
  plausible titles hold nothing relevant. Checks on the clone: the proposal
  stays `proposed` and respects the constraint (rubric); retrieval facts
  say whether the holding record was read and which distractors were.
- `code-constraint`: a `decision` record names the file the topic's
  implementation must change, in a code span, and states a constraint the
  brief does not. The same checks; retrieval facts say whether that decision
  was read, and, on the `with` row, whether `grove search` ran.

A linked review record compares rows per case: found and applied, holding
record read, unneeded reads, turns and cost, and says whether the difference
justifies the next lever or none.

Out of scope: any product change; the work row; a code-graph case, which
has no product to test; scoring the retrieval facts.

## Acceptance

1. Both cases exist in the runner with their fixture records, checks and
   rubric rows; `python3 evals/run.py selftest` covers each check and each
   retrieval fact they add; [evals/README.md](../evals/README.md) documents
   them and the two-digest comparison.
2. Under an owner mandate naming a Claude model, runs, budget and
   permission mode, and never through the Codex harness, each case runs on
   its `without` row, and on its `with` row once G-153's
   candidate is integrated, retaining per run what the README lists.
3. A review record linked by `work` reports the pattern per case and row:
   whether the constraint was found and applied, whether the holding record
   was read, unneeded reads and cost; names the lever it points at or
   concludes that no change is justified; and says whether a symbol-level
   or index question deserves a case. It records its limits as G-122 does.
4. Nothing is spent without the mandate, nothing runs in `go test` or CI,
   and fixture records are created only in disposable directories.

## Next

Assign the `without` row now: build both cases, then one mandated run each
on Claude. The `with` row waits for G-153 to be integrated. The owner
supplies the mandate before any paid run; it names a Claude model, and no
attempt passes `--harness codex` or names `gpt-6-astra`. G-153 starts after
the `without` row's review is read (owner decision, 2026-09-25).

**Checkpoint, 2026-09-25, second headless `/grove-work G-154`.** The
`without` row is run and reviewed; the `with` row waits for G-153 to be
integrated, and G-153 waits for the owner to read
[G-160](G-160-g-154-without-row-both-constrain.md). G-154 alone, on
`worktree-G-154` in `.claude/worktrees/worktree-G-154`, base main
`b684951`; plan [G-155](G-155-g-154-listed-and-code-constraint.md) at
`sha256:fc27deb1…`, this record at `sha256:8efad9e4…` when the session
started. Status stays `active`: acceptance 2 and 3 are half done.

- Done, plan steps 1 and 2 (acceptance 1 and 4): `listed-constraint` and
  `code-constraint` in `evals/run.py`, each on its own copy of the fixture
  with records from `evals/fixture/records/`; the retrieval facts
  `search`, `holding read` and `distractors read`; per-case rubric
  columns; README sections for the cases, facts, rubric rows and the
  two-digest comparison. With no `--case` the runner still runs only the
  G-108 pair, whose fixture is unchanged. Commits `b2c7ef0`, `946459a`,
  `47223ed`, `c895666`.
- Done, step 3: independent review
  [G-159](G-159-g-154-runner-cases-review.md), three rounds, examined
  `47223ed`; every consequential finding fixed, one post-cap fix
  (`c895666`) self-checked only, one informational coverage gap open.
- Done, step 5 (acceptance 2 and 3 for the `without` row): under
  [G-158](G-158-what-mandate-should-the-g-154-wi.md)'s answer, "Use the
  G-122 settings" (`claude-opus-5-5`, 5 runs, $5 a run, `auto`,
  `~/.cache/grove-evals/claude`), both cases ran at `eeec725`, guides
  digest `41324c3655a1`, as two foreground pieces, one per case, into
  `~/.cache/grove-evals/runs/2026-09-25-G-154-without/{listed-constraint,code-constraint}`;
  $3.35 spent, never the Codex harness. Review
  [G-160](G-160-g-154-without-row-both-constrain.md), examined `eeec725`:
  both constraints found and applied in 10 of 10 runs, because at five and
  two records every run reads every record; code-constraint runs 4 and 5
  also blocked on an unplanted syntax question; the runner's `holding
  read` and `distractors read` miss reads through a glob (5 of 10
  false negatives). Lever: none from this row.
- Verification at the checkpoint commit: `python3 evals/run.py selftest`
  prints `selftest: ok`; `go run ./cmd/grove check` passes; links in
  G-160 resolve. No code changed.
- No command is running.

Pending judgments, the owner's: read G-160, then start G-153 or not, since
the row shows these fixtures cannot show search finding what the listing
misses; and whether a larger-fixture case deserves its own work.

Continuation, once G-153's candidate is on main: assign `/grove-work
G-154` again on this branch. It runs `git merge main`, first makes the
runner's `holding read` and `distractors read` count a glob read, with a
selftest case, and recomputes the `without` row's facts from its
transcripts (G-160 Disposition); then runs the `with` row with the same
two commands under G-158's mandate and a new `--out`
(`2026-MM-DD-G-154-with/…`), completes the review across both rows, and
hands off into Review (plan step 6).
