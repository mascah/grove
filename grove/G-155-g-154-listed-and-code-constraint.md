---
id: "G-155"
type: plan
title: "G-154 listed and code constraint cases plan"
status: current
created: "2026-09-25T19:29:24Z"
updated: "2026-09-25T19:29:45Z"
work: ["G-154"]
---

## Design

Prepared 2026-09-25 on `worktree-G-154` from main `b684951`, against
[G-154](G-154-listed-constraint-eval.md) at `sha256:d569262b…`, in a
headless `/grove-work G-154` session. Everything below is proposed design
inside G-154's selected scope unless it quotes the record.

**Fixture.** The G-108 pair keeps its fixture byte for byte, so its
reruns stay comparable. Each new case gets a variant of the built
template: a copy of it (counter included), its records created with the
built CLI and an explicit absolute `--project` inside the run's output
directory, bodies from `evals/fixture/records/`, fields through `update`,
one commit. A `done` record gets the template's first commit as its
`candidate`, which `update` requires to be in HEAD. Each variant's commit
is reported per run beside the base fixture's.

**`listed-constraint`**, topic "make due dates for tasks ready to assign".
A proposed work record, "Give tasks a due date", carries the owner's
intent for the feature (`tasks add --due`, the date shown in `list`, no
sorting or reminders) and needs acceptance; it `depends_on` a `done`
record, "Add tasks export", and relates to two proposed records,
"Export tasks as CSV" and "Remind the owner of tasks due today". Only the
`done` record holds the constraint: the owner's phone widget reads
`export` and rejects any key but the seven it has, so a new frontmatter
key stays out of `export` until the widget changes, which is its own work.
`load` in `tasks.py` passes every frontmatter key through, so a `due` key
reaches `export` unless the proposal keeps it out. The brief says only
that `export`'s JSON is stable. The refined record is the natural target:
`context` on it lists the prerequisite and both distractors.

**`code-constraint`**, topic "add a tasks tag command that adds or removes
tags on an existing task". An `accepted` decision, "Keep the owner's
notes in task files", names `write` in `tasks.py` in code spans, never a
link, and states that a command rewriting a task file keeps everything
below the title, and that `write` today drops it. The brief says nothing
of notes. Its title shares no word with the topic but "task".

Both cases expect a proposal and no question, like the companion: each
constraint is answered by a record, so a question about it is an over-ask,
and the `no-question` check reports it. The other clone checks apply
unchanged.

**Retrieval facts added**, from the trace as today: `holding_read`,
whether the holding record was read (a `grove show` or `grove context`
naming its ID, or a read of its file, `context --include` among them);
`distractors_read`, the distractor
IDs read the same way; and `search`, whether `grove search` ran, on every
case, since it does not exist before G-153. A distractor's file counts as
unneeded; the refined record and the holding record count as needed.
`grep` stays uncounted, as the README says.

**Rubric.** Per-case columns: the two new cases score "constraint
applied", "brief constraint" and "handoff"; the pair's columns stay as
they are.

**Two-digest comparison.** The `without` row runs from this branch as it
stands (guides digest before G-153); the `with` row from this branch
after G-153's candidate is on main and merged in, same cases, model, runs,
budget and permission mode. The README says so.

**Spend.** No mandate came with this assignment. Per G-141, a missing
item is not a default: the session persists a question naming model,
runs, budget, permission mode and config directory, and stops before any
paid run.

## Steps

1. Fixture variants, the two cases, retrieval facts and per-case rubric
   columns in `evals/run.py`; the fake covers them; `selftest` ok.
2. `evals/README.md`: cases, facts, rubric rows, the two-digest
   comparison.
3. Independent review of steps 1 and 2 (consequential boundary: the
   runner every later row depends on).
4. Mandate question, blocking G-154; checkpoint.
5. Under the mandate: the `without` row, one run set per case, and a
   review record of its pattern (G-154 acceptance 3 for that row).
6. After G-153 integrates: merge main, the `with` row, the review
   completed across both rows; handoff into Review.

Steps 1 to 4 done, 2026-09-25 (G-154's Next has the evidence); step 5
waits on G-158.
