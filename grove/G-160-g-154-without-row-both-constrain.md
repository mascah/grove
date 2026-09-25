---
id: "G-160"
type: review
title: "G-154 without row: both constraints found in every run"
status: current
created: "2026-09-25T20:27:37Z"
updated: "2026-09-25T20:28:38Z"
work: ["G-154"]
examined: "eeec725"
---

## Examined

The `without` row of [G-154](G-154-listed-constraint-eval.md) acceptance 2:
the two cases of plan [G-155](G-155-g-154-listed-and-code-constraint.md),
five headless shaping runs each on Claude, under the mandate
[G-158](G-158-what-mandate-should-the-g-154-wi.md) set ("Use the G-122
settings", which are [G-122](G-122-g-108-baseline-runs-the-missing.md)'s
command). Run 2026-09-25 20:15 to 20:26 UTC from `worktree-G-154` at
`eeec725`, the `examined` commit, clean, by a headless `/grove-work G-154`
session, as two foreground pieces, one per case, so that each fits the
harness's ten-minute command limit:

```sh
python3 evals/run.py run --case listed-constraint --runs 5 --budget 5 \
    --model claude-opus-5-5 --permission-mode auto \
    --config-dir ~/.cache/grove-evals/claude \
    --out ~/.cache/grove-evals/runs/2026-09-25-G-154-without/listed-constraint
python3 evals/run.py run --case code-constraint --runs 5 --budget 5 \
    --model claude-opus-5-5 --permission-mode auto \
    --config-dir ~/.cache/grove-evals/claude \
    --out ~/.cache/grove-evals/runs/2026-09-25-G-154-without/code-constraint
```

Configuration, from the two reports: Claude Code 2.1.282, model reported
`claude-opus-5-5`, permission mode `auto`, guides digest `41324c3655a1`
(before G-153), fixture commits `85a718a` (listed-constraint) and `2fa3c4f`
(code-constraint), cap $25 per piece printed first. The config directory
held what G-122 finding 5 describes: the login, `settings.json` with
`autoMemoryEnabled: false`, `theme` and `tui`, seven synced Anthropic skills
and three synced plugins. The output directories hold every transcript,
clone, remote, `state.json`, `run.json` and each piece's `report.md`,
outside every checkout and not committed. Each score below was given by the
implementing session from `run.json`, the transcript's commands and the
records on each proposal branch, so its scorer is `judge`; the owner column
is open.

In both fixtures the holding record is G-002: the `done` "Add tasks export"
in listed-constraint, whose distractors are G-003 "Export tasks as CSV" and
G-004 "Remind the owner of tasks due today"; the `accepted` decision "Keep
the owner's notes in task files" in code-constraint, which has no
distractor.

## Findings

**1. listed-constraint: found and applied in 5 of 5 runs.** Every run
refined G-005 "Give tasks a due date" in place, `proposed`, with no
question, and every acceptance has an item that `tasks export` keeps
exactly its seven keys and no `due`, citing G-002 or the widget in the
item or in Scope. Every run also found in `tasks.py` that `export` passes
every frontmatter key through, and, unplanted, that `done` and `drop`
rewrite through `write` and would drop `due`; each added an acceptance item
keeping it. Every check passed in every run.

**2. code-constraint: found and applied in 5 of 5 runs; 2 of 5 also asked
an unplanted question.** Every run read G-002, and every proposal requires
that a tagged task keeps what is below its title, either through a separate
prerequisite that fixes `write` (runs 2, 4, 5: a new "Keep a task's notes
when a command rewrites it" the tag work depends on) or inside the tag
record itself (runs 1, 3), citing G-002. No question asked what the
decision answers. Runs 4 and 5 fail `no-question`: each opened G-005 on the
command's syntax (`--add`/`--remove`, `+tag -tag`, or `tag`/`untag`),
blocking the tag record, with option 1 recommended; runs 1 to 3 wrote the
same syntax as a non-binding suggestion. The brief's conventions do not
settle the syntax, so this is the topic leaving a user-facing choice open,
not the constraint being missed; it is noise in this case's check, not a
retrieval result.

**3. At this size every run reads nearly every record, so the listing is
not what finds the constraint.** Reading the transcripts, not the facts
alone: all ten runs read the holding record, all five listed-constraint
runs read both distractors, and 8 of 10 read the unrelated G-001 (not
listed-constraint runs 3 and 4). Runs 2 and 4 of listed-constraint and 1,
4, 5 of code-constraint used `grove show` on the records they read; the
others `cat` the
record files by glob (`cat grove/G-00*.md`, a `for` over `grove/G-*.md`,
or `cat grove/G-005*.md grove/G-004*.md …`). listed-constraint run 3 ran
`grove context G-005` after it had already read G-002 to G-005. With five or two
records, `ls grove` and one `cat` cost less than choosing, so neither
`context`'s listing nor the distractors' titles decide anything, and the
case cannot show G-153's search helping to find the constraint: the
`without` row has no misses to recover.

**4. The runner's `holding read` and `distractors read` miss glob reads.**
They report `holding read` false for listed-constraint runs 1, 3, 5 and
code-constraint runs 2, 3, and `distractors read` `none` for
listed-constraint runs 1, 3, 5, while each of those runs read the records
through a glob path (finding 3). The README's warning to read the
transcript before resting on one fact covers this, and the facts are never
scored, but a comparison across rows on those two columns alone would be
wrong: 5 of 10 `holding read` values are false negatives. Unneeded reads
count only paths too: the one reported is `.gitignore` (code-constraint
run 2); the distractors read by glob appear nowhere.

Recomputed after the fix. Commit `40e1263` makes both facts, and the
unneeded reads, count a read through a glob or a `for` loop over paths,
with a selftest case. `retrieval()` at that commit, rerun on the ten
retained transcripts and post-run clones (a throwaway script importing
`evals/run.py`, not committed; rerun at `7026cbe`, the final runner, with
the same result) gives: `holding read` true in 10 of 10;
listed-constraint `distractors read` G-003 and G-004 in 5 of 5; unneeded
reads G-003 and G-004 in listed-constraint runs 1, 3, 5, the shared
fixture record G-001 "Sync tasks between two machines" in listed-constraint
runs 1 and 5 and code-constraint runs 2 and 3, and `.gitignore` in
code-constraint run 2. The unneeded reads count files only, so G-001 read
through `grove show` appears in no fact: counting those (listed-constraint
run 2, code-constraint runs 1, 4, 5), G-001 was read in 8 of 10 runs, and
a comparison on unneeded reads must add them from the transcripts, or
`show` looks cheaper than `cat`. The reports under `--out` keep the values
from `eeec725`; these, with that addition, are the facts a later row
compares against. G-154's final independent review reproduced them.

**5. Cost and shape.** listed-constraint $0.29 to $0.33 a run, 6 to 9
turns, 50 to 64 seconds, $1.55 for five; code-constraint $0.34 to $0.38,
9 turns, 56 to 72 seconds, $1.80 for five; $3.35 for the row against the
$50 cap. Exit 0, no permission denials, no timeouts, no error results.
Tools across the ten transcripts: 71 `Bash` calls and one `Write`, inside
the clone; no skill, agent or MCP tool, and no tool input names a path
outside the clone. This is the G-108 pair's cost: reading the records
added nothing measurable.

Rubric ([`evals/README.md`](../evals/README.md)), scorer `judge` (the
implementing session), owner column open:

| case | run | scorer | constraint applied | brief constraint | handoff |
| --- | --- | --- | --- | --- | --- |
| listed-constraint | 1 to 5 | judge | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |
| code-constraint | 1 to 5 | judge | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |

Notes. Brief constraint: every due-date proposal keeps `due` as one
frontmatter key in a file that stays in place, and every tag proposal
keeps tags lowercase and the file in place. Handoff: every final message
names the branch, the tip commit, the records with revisions, and the next
action or, in code-constraint runs 4 and 5, the question it waits on.

## Disposition

**Lever.** No change is justified by this row. Both constraints were found
and applied in every run on Claude Opus 5.5 at guides digest
`41324c3655a1`, at the G-108 pair's cost, without search, because in
fixtures of five and two records the agent reads nearly all of them. G-153 then
shipped no agent-facing search, so a `with` row at main would have
measured only the guide edits landed since this digest, and even with a
search it could not show search finding what the listing missed, since
this fixture has nothing missed; the owner dropped it (G-173). Evidence
that would justify G-153's search, or its code-to-record links, needs a
fixture large enough that reading every record costs more than choosing,
tens of records with plausible titles, which neither case has. That is a
new case and an owner choice, not a change to these two, whose `without`
row the `with` row must match.

The runner's `holding read` and `distractors read` now count a read
through a glob (finding 4, `40e1263`), and this row's facts are
recomputed on that definition, so a later `with` row compares against
them.

Finding 2's syntax question is fixture noise. Removing it means naming the
syntax in the topic, which would change the case after its `without` row;
the owner may instead read `no-question` for code-constraint as passing
when the only question is the syntax.

**Symbol-level or index question.** Not worth a case now. Every run read
`tasks.py` whole, 74 lines, and found `write` and `load` from it,
without a code span or link pointing there; an index or a symbol graph
has nothing to shorten at this size. The question deserves a case only
together with a fixture large enough that reading everything fails, and
the same case would test search first.

**Limits.** Claude Opus 5.5 only, five runs per case, one fixture of five
or two records, a constraint each held by one record, and a judge that is
the implementing session. Five runs show a pattern, not a rate. The row
was run as two pieces with separate reports rather than G-158's single
command with one `--out`; the arguments are otherwise identical, and the
cap per piece is half. The checks saw the clone; the trace showed no write
outside it. The `with` row was dropped by the owner's answer to
[G-173](G-173-what-should-g-154-s-with-row-bec.md), since G-153 shipped no
agent-facing search, so this record makes no comparison across rows; a
later search command or guide text carries its own `with` row on these
cases.
