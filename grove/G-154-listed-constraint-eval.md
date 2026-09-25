---
id: "G-154"
type: work
title: "Evaluate whether agents find a constraint held in a listed record or in the code they touch"
status: review
created: "2026-09-25T19:15:12Z"
updated: "2026-09-25T21:42:28Z"
relates_to: ["G-108", "G-122", "G-135", "G-153"]
kind: investigation
candidate: "357c5079248dfec989cf155a1bf7833f20e5cc5a"
---

## Outcome

The owner has repeatable evidence of whether an agent shaping in an adopting
project finds and applies a constraint that the selected record does not
hold: one kept in a listed prerequisite or plan among plausible distractors,
and one kept in a record that describes the file the proposal must touch;
at what cost in reading.

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

Owner answer to [G-173](G-173-what-should-g-154-s-with-row-bec.md),
2026-09-25: "Drop the with row". G-153 shipped body search on the board
only, with no `grove search` and no guide text, so a headless `with` row
had nothing to measure. This work ends on its `without` row; a later
agent-facing search, shaped as its own work, carries its own `with` row
on these cases. The outcome no longer asks whether G-153's search changes
the result, and acceptance 2 and 3 were amended to match.

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

The `with` row was dropped (G-173); the runner and the README keep the
comparison for the work that ships a search.

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
   its `without` row, retaining per run what the README lists. (The `with`
   row was dropped by the owner's answer to G-173.)
3. A review record linked by `work` reports the pattern per case on the
   `without` row:
   whether the constraint was found and applied, whether the holding record
   was read, unneeded reads and cost; names the lever it points at or
   concludes that no change is justified; and says whether a symbol-level
   or index question deserves a case. It records its limits as G-122 does.
4. Nothing is spent without the mandate, nothing runs in `go test` or CI,
   and fixture records are created only in disposable directories.

## Evidence

Branch `worktree-G-154` in `.claude/worktrees/worktree-G-154`, base main
`b684951`, which merges into main `47852e3` cleanly. The final session
started from this record at `sha256:8fc9ba32…` and plan
[G-155](G-155-g-154-listed-and-code-constraint.md) at `sha256:64da8623…`;
the candidate is the commit that adds this section.

Against each acceptance item:

1. `listed-constraint` and `code-constraint` exist in
   [evals/run.py](../evals/run.py), each on its own copy of the fixture
   with records from `evals/fixture/records/`, per-case rubric columns and
   the retrieval facts `search`, `holding read` and `distractors read`;
   [evals/README.md](../evals/README.md) documents them and the two-digest
   comparison, now saying G-153 shipped no search. `selftest` covers every
   check and fact, including reads through a glob or a `for` loop
   (`40e1263`), a glob character in the clone path and an `--include` path
   `context` refuses (`7026cbe`); each of those assertions fails with its
   fix reverted. With no `--case` the runner still runs only the G-108
   pair, on its unchanged fixture. Commits `b2c7ef0`, `946459a`,
   `47223ed`, `c895666`, `40e1263`, `7026cbe`.
2. The `without` row ran under [G-158](G-158-what-mandate-should-the-g-154-wi.md)'s
   answer, "Use the G-122 settings" (`claude-opus-5-5`, 5 runs a case, $5
   a run, `auto`, `~/.cache/grove-evals/claude`), at `eeec725`, guides
   digest `41324c3655a1`, never the Codex harness: $3.35, every retained
   file the README lists under
   `~/.cache/grove-evals/runs/2026-09-25-G-154-without/`. The `with` row was
   dropped (G-173).
3. [G-160](G-160-g-154-without-row-both-constrain.md), examined `eeec725`:
   both constraints found and applied in 10 of 10 runs, the holding record
   read in 10 of 10 (recomputed with the fixed runner), distractors read in
   every listed-constraint run, G-001 in 8 of 10, $0.29 to $0.38 a run;
   lever none, because at five and two records every run reads nearly
   every record; a symbol-level or index case is not worth it without a
   fixture large enough that reading everything fails; limits as G-122
   records them.
4. No paid run without the mandate; no run in `go test` or CI; fixture
   records only under each run's output directory or the selftest's
   temporary one.

Decisions taken: the Outcome and acceptance 2 and 3 were amended to the
`without` row by the owner's answer to G-173; no decision record, since
reversing it costs one more row on the unchanged runner. The recomputed
facts are a throwaway script's, reproduced by both final reviewers; the
reports under `--out` keep the `eeec725` values.

Verification at `144329d` (the candidate adds only records):
`python3 evals/run.py selftest` prints `selftest: ok`; `go run
./cmd/grove check` prints `OK: 155 records`; `go test -count=1 -timeout
120s ./...` passes (no Go changed on the branch); every relative link in
the changed records and the README resolves.

Reviews: [G-159](G-159-g-154-runner-cases-review.md), the runner at
`47223ed`, three rounds, every consequential finding fixed;
[G-160](G-160-g-154-without-row-both-constrain.md), the row;
[G-181](G-181-g-154-final-review-runner-fixes.md), the final gate, two
rounds, examined `c45af99`, every finding minor and fixed but two
informational limits; `144329d` after it rewords G-160 as its round 2
asks, self-checked.

Open limits: the `for`-loop rule looks for a reader anywhere in the
command (G-181 round 1 finding 2); the include rule does not copy
`context`'s `.git` and symlink refusals (round 2 finding 2); the rubric's
owner column in G-160 is unscored.

## Next

In review: the owner judges the candidate. From this checkout:

```sh
go run ./cmd/grove approve G-154 "VERDICT"
```

then, in main's checkout:

```sh
go run ./cmd/grove integrate G-154
```

or `go run ./cmd/grove feedback G-154 "TEXT"` here for changes.

Pending judgments, the owner's: the Outcome's amendment after G-173; the
rubric's owner column in G-160; whether a larger-fixture case, which could
also carry a later search's `with` row, deserves its own work (G-160
Disposition).
