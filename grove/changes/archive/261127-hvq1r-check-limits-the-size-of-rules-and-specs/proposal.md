# Check limits the size of rules and specs

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

Maintained specs only grow, and nothing fails when they do. A requirement's
"rule" below is its text before its first scenario; words are counted with
`wc -w`.

- **Spec text grew by two thirds in five weeks.** The files
  `grove/specs/*/spec.md` hold 18,065 words at ea7fbaf (2026-10-25) and
  30,463 at 3cfd2ef.
- **The one rewrite lasted 22 days.** `grove/specs/implement/spec.md` held
  4,992 words before pull request #112 rewrote it, 2,713 after it
  (2026-10-31), and 5,102 on 2026-11-22, as it still does at 3cfd2ef.
- **Rules are no longer a sentence a reader can hold.** At 3cfd2ef the 107
  rules average 72 words; 57 are over 60 and 22 over 100. The longest,
  "Print the project as data" in `grove/specs/changes/spec.md`, is 189.
  The guide asks for "one EARS sentence"
  (`internal/instruct/schemas/default/documents/change-spec/guide.md`).
- **A modification nearly always lengthens.** Of the 106 MODIFIED
  requirements in the change specs under `grove/changes/archive/` at
  3cfd2ef whose earlier text is there too, taken in order of change ID, 77
  are longer than the text they replace, rule and scenarios together.
  Those add 6,474 words; the 22 shorter ones cut 989.
- **Nothing pushes back.** `grove check` asks a requirement only for
  "shall" and a scenario (requirements spec, "Each requirement says shall
  and has a scenario"). Archive refuses a fold that alters a requirement
  the change does not name (`CONSTITUTION.md` rule 3), so a change cannot
  tidy in passing.
- **Text with a failing limit did not grow.** The instruction text under
  `internal/instruct`, which `just check` holds to 9,000 words
  (`justfile`; `CONSTITUTION.md` rule 12), is 7,931 words at ea7fbaf and
  8,458 at 3cfd2ef: the same repository, agents and weeks.
- **It is not only Grove's own specs.** Two other projects of the owner's
  that use Grove, private and measured on the owner's machine on
  2026-11-27, show the same: one holds 71,805 words in 37 specs five weeks
  after its first commit, its rules averaging 121 words; the other 22,543
  words in 9 specs after two weeks.
- **The project Grove borrows its documents from reports the same.**
  OpenSpec's lead maintainer writes that "each new behavior lands as
  another clause on the nearest requirement"
  (<https://github.com/Fission-AI/OpenSpec/issues/2139>). It has since
  shipped warnings for a requirement over 500 characters and for a
  MODIFIED one that grows past that
  (<https://github.com/Fission-AI/OpenSpec/blob/main/CHANGELOG.md>,
  1.14.1 and 2.1.1).

## Intended outcome

`grove check` fails a change that adds a rule over 60 words, lengthens a
rule that is over 60 words, or leaves a maintained spec over 2,500 words
larger than it found it. Long text already in the specs is left alone
until a change touches it. From then on every change leaves the specs it
touches no worse, and the pressure to shorten a rule or split a spec
arrives while proposing, where the owner reviews the result.

## Scope

### In scope

- Three reports in `grove check`, each about an open change: a new rule
  over the limit, a modified rule that grew, and a spec the change grows
  past the limit.
- One line in the guide to a change spec stating both limits, so an agent
  writes to them instead of finding them through `check`.
- The reference and the walkthrough, and the first sentence of the
  requirements spec's Purpose.

### Out of scope

- Reporting text already over a limit that no open change touches.
- A setting in `grove.yaml` for the limits.
- A limit on the number or length of scenarios, other than through the
  spec's total.
- `grove archive` or `grove claim` refusing anything for size: `check`
  reports and blocks nothing (`CONSTITUTION.md` rule 10).
- Rewriting the five specs that are over 2,500 words at 3cfd2ef (propose,
  implement, cycles, runs and serve). Each is its own change, by the
  guides to splitting and reconciling a maintained spec.
- A budget for Grove's own specs in `just check`.
- Trimming the archive, which is pull request #205.
- Linking a scenario to the test that shows it.

## Decisions

- **Limit growth with a failing check.** The owner, on the review that
  recommended it first of four options: "yes i agree with this. draft 1
  and 2". Reason: the budgeted instruction text held its size while the
  unbudgeted specs grew, in the same repository over the same weeks.

Accepted in the round before writing (the owner: "yes to all five"):

1. **The limits are 60 words for a rule and 2,500 for a spec.** At 3cfd2ef
   57 of 107 rules and 5 of 14 specs are over, so those can stay or
   shrink and cannot grow.
2. **Only growth is reported.** Text already over a limit is left alone
   until a change touches it; otherwise `grove check` would fail in every
   project on the day it shipped.
3. **A report fails `grove check` like every other.** Grove has no warning
   level, and the instruction budget holds because it fails.
4. **The numbers are fixed in Grove.** No `grove.yaml` setting until a
   project needs other numbers (`CONSTITUTION.md` rule 8).
5. **Scenarios are not counted separately.** They are limited only
   through the spec's total.

## The text around the requirements

Archive replaces requirements, not the text around them, so T-01 writes
it. The first sentence of the requirements spec's Purpose, which is also
the spec's description in the project index, becomes:

> `grove check` wants "shall" and a scenario in each requirement and
> reports a change that grows a rule or a spec past its limit, and `grove
> archive` refuses a fold that does not hold the change's requirements
> word for word.

## Open questions

None must be answered before review. These can wait:

- **A MODIFIED requirement that `base` does not hold, because a change in
  `depends_on` adds it: longer than what?** Recommended: than the text
  that dependency gives it, which the stale check already reads (`gives`
  in `internal/cli/check.go`), so that a change shortening a dependency's
  long rule is not reported. The same reading serves the spec's total.
  Decides T-01.
- **Routine choices inside the agreed approach:** the wording of the three
  messages, each one sentence; the name of the word-count method; the
  tests' names and places; which fixture requirement the walkthrough
  lengthens. Recommended: the implementer chooses. Decides T-01 and T-02.
- **Whether the limits suit the owner's other projects,** where rules
  average 121 words. Using it will tell; a `grove.yaml` setting is a
  later change if they do not. Decides no task.
- **Whether and in what order Grove's five long specs are rewritten.**
  Decides no task.

## Related documents

- Requirements: [specs/requirements/spec.md](specs/requirements/spec.md)
- Tasks: [tasks.md](tasks.md)
