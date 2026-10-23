# Fewer words in Grove's instructions

## Problem

Grove's instructions say some rules twice, and the agent reads every word
each time a skill runs.

- Every skill opens with the same text. `Markdown` in
  `internal/instruct/instruct.go` at 0c6f419 prints the `meanings` const
  before anything else: about 100 words on the three meanings, drafts and
  commits. Its last sentence, "Commit only when the owner or the project's
  rules ask", is also in the finish rule of every requires file under
  `internal/instruct/requires/`. Its "merging into the target branch is
  the only approval" repeats its own bullet on approval.
- A skill's requires file and its schema guide say the same thing. At
  0c6f419, `requires/implement.md` ends with "archiving brings the
  maintained specs up to date on this branch", and
  `schemas/default/skills/implement.md` ends with "Archiving the change,
  once every task is done, brings the maintained specs up to date on this
  branch". `requires/apply.md` and `skills/apply.md` repeat the same pair.
  `requires/propose.md` asks one round, numbered, with recommended answers
  worded so that "yes" accepts them; step 3 of `skills/propose.md` says it
  again. `requires/explore.md` and `skills/explore.md` both say to write
  nothing until the owner asks.
- At 0c6f419, `internal/instruct` holds 6,414 words of Markdown and YAML.
  On the ledger fixture, `grove instructions propose` prints 2,404 words,
  `plan` 971, `prototype` 670 and `explore` 653. Cycle B adds text to
  these files before this change is built.
- `just check` counts Go lines and nothing else (`justfile` at 0c6f419).
  Constitution rule 12, as PR #47 amends it, adds a budget of 7,000 words
  of instruction text in `internal/instruct`, counted by `just check`.

## Intended outcome

Each rule in Grove's instructions is said once. The opening text printed
for every skill is shorter. Where a line is in both a skill's requires
file and its schema guide, only the requires file keeps it. No rule is
dropped: the agent is asked to do the same things in fewer words.

`just check` counts the words of `internal/instruct`'s Markdown and YAML
files and fails above 7,000.

This change has no change spec. No requirement a user or caller observes
changes: the instructions keep every rule, and the budget is a rule about
Grove's own source (constitution rule 12), not behavior. No maintained
spec at 0c6f419 fits it: apply, config, cycles, implement, propose,
prototype and runs each describe a skill or a command. The work is
instruction text and the `justfile`.

## Scope

### In scope

- The `meanings` const and the commit sentence in
  `internal/instruct/instruct.go`.
- Lines in both a requires file in `internal/instruct/requires/` and its
  schema guide in `internal/instruct/schemas/default/skills/`, as cycle B
  leaves them, and in the built-in `bugfix` schema's propose guide.
- The word count and budget in `justfile`, and the `just check` line of
  `CONTRIBUTING.md`.
- Measuring each skill's instructions on the fixture, and running the
  trials of 261023-kh3gg, before and after, in the implementing pull
  request's description.

### Out of scope

- Dropping or changing any rule. A rule that should go is a change of its
  own, with a change spec.
- Repetition between a document kind's guide and its template, and in
  `skills.yaml` or `metadata/`: not in the decisions below.
- `docs/` and `dev/walkthrough.md`: no step changes, so no page does (rule
  8).
- Raising the budget: a change to the constitution (rule 12).

## Decisions

The owner agreed to the direction on 2026-10-23 ("yes to all") and to
every answer below in one round ("yes to C & D").

- **Remove repetition.** Shorten the opening text printed for every skill,
  the `meanings` const and the commit sentence in `instruct.go`; where a
  line is in both a skill's requires file and its schema guide, keep it
  only in the requires file. The implementing pull request says where
  every removed rule still lives. No rule is dropped. Reason: the agent
  reads each rule once, and a reviewer can check nothing was lost.
- **`just check` counts the words of `internal/instruct` (Markdown and
  YAML) and fails above 7,000,** as rule 12 says once PR #47 amends it.
  Reason: the instruction text has a budget the way the Go code does, and
  the check keeps it.
- **Measured before and after.** The implementing pull request's
  description gives each skill's word count on the fixture before and
  after, and the results of the trials of 261023-kh3gg (`dev/trials.md`)
  run before and after. Reason: fewer words must not change what the
  agent does.

The owner answered the proposal's two open questions with "yes"
(2026-10-23):

- **The built-in `bugfix` schema's propose guide follows the same rule.**
  `internal/instruct/schemas/bugfix/skills/propose.md` (261023-vqsrt)
  repeats no line already in a requires file. Reason: it replaces the
  default propose guide, so the same repetition would come back there.
- **Over 7,000 words, cut further; else stop and ask.** If
  `internal/instruct` still holds more than 7,000 words once cycle B has
  landed and T-01 is done, the implementer makes more cuts first; failing
  that, it stops and asks the owner. Only a constitution amendment raises
  the budget (rule 12). Reason: the budget is the constitution's, not the
  implementer's, to change.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: none; this change has no change spec (see Intended
  outcome).
- Tasks: [tasks.md](tasks.md)
