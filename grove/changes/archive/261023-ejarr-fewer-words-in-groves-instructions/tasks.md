## T-01 — Say each rule once

Result: the opening text every skill prints is shorter, no line is in
both a skill's requires file and its default schema guide unless removing
it would need the copy that stays reworded, and every rule the
instructions held before is still printed for its skill.

Files: `internal/instruct/instruct.go`,
`internal/instruct/requires/propose.md`,
`internal/instruct/requires/plan.md`,
`internal/instruct/requires/archive.md`,
`internal/instruct/requires/explore.md`,
`internal/instruct/requires/prototype.md`,
`internal/instruct/requires/implement.md`,
`internal/instruct/requires/implement-cycle.md`,
`internal/instruct/requires/implement-owner.md`,
`internal/instruct/requires/implement-unattended.md`,
`internal/instruct/schemas/default/skills/propose.md`,
`internal/instruct/schemas/default/skills/plan.md`,
`internal/instruct/schemas/default/skills/archive.md`,
`internal/instruct/schemas/default/skills/explore.md`,
`internal/instruct/schemas/default/skills/prototype.md`,
`internal/instruct/schemas/default/skills/implement.md`,
`internal/instruct/schemas/bugfix/skills/propose.md`,
`internal/cli/walkthrough_test.go`, `dev/trials/run.sh`, `dev/trials.md`,
`dev/trials/b.sh`

Interfaces: none; it is instruction text.

Constraints: written against the files as cycle B leaves them;
`requires/implement-owner.md` and `requires/implement-unattended.md` are
261023-g0a3h's, which also deletes `requires/apply.md`,
`requires/apply-cycle.md` and `skills/apply.md`. The `Status:` lines
stay exactly as the templates have them (rule 3). A test assertion on a
line that moves or is reworded still checks that rule.

Acceptance:
- The implementing pull request lists each removed or shortened line and
  the file where its rule still lives; for example, "archiving brings the
  maintained specs up to date" stays in `requires/implement.md` only.
- No rule a file in `requires/` states is stated again in the schema
  guide of the same skill in `schemas/default/skills/`, or in
  `schemas/bugfix/skills/propose.md`, except where removing it would need
  the copy that stays reworded; the pull request lists each such repeat.
- The opening text, printed before "What Grove requires of", is shorter
  and still says what a maintained spec, a change, observed behavior and
  a draft are, that Grove commits nothing and that merging is the only
  approval.
- The pull request's description gives, on the ledger fixture, the words
  printed by `grove instructions` for propose, plan, prototype, explore,
  archive, and implement for a change and a cycle, with and without
  `--unattended`, before and after; none is longer after.
- The trials of 261023-kh3gg, in `dev/trials.md`, are run on the fixture
  before and after, and both results are in the pull request's
  description.
- `internal/cli/walkthrough_test.go` passes, and `grove check` on Grove's
  own repository prints "No problems in grove/.".

## T-02 — `just check` counts the instruction words

Result: `just check` counts the words of the Markdown and YAML files in
`internal/instruct` and fails above 8,000, as constitution rule 12 says.

Files: `justfile`, `CONTRIBUTING.md`

Interfaces: none.

Constraints: the number 8,000 is written in rule 12 and the `justfile`
only; `CONTRIBUTING.md` says `just check` checks the budgets and names no
number. If `internal/instruct` holds more than 8,000 words once T-01 is
done, make more cuts first; failing that, stop and ask the owner. Only a
constitution amendment raises the budget (rule 12).

Acceptance:
- `just check` prints the word count beside the line counts and passes.
- With a file added to `internal/instruct` that takes the count above
  8,000, `just check` fails saying the budget is exceeded; with it
  removed, it passes.
- `CONTRIBUTING.md`'s `just check` line, and the `justfile`'s comment on
  `check`, name the word budget with the line budget.
