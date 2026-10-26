# A title with a colon passes check

## Observed

Grove's templates for a change's `.grove.yaml` and a cycle file print the
title unquoted: `title: <title>` in
`internal/instruct/metadata/change.grove.yaml` and
`internal/instruct/metadata/cycle.yaml` at ea7fbaf. An agent filling the
template with a title holding `: ` writes YAML that does not parse.

In memex, two of five proposals had a colon in the title and failed `grove
check` with "invalid YAML: mapping values are not allowed in this context"
(`dev/memex-review-2026-10-26.md`, finding 2).

Reproduced on the ledger fixture with grove 0.5.0: a change whose
`.grove.yaml` holds `title: Records: delete by ID` makes `grove check`
print

    grove/changes/261026-aaaaa-x/.grove.yaml: invalid YAML: mapping values are not allowed in this context
    grove: 1 problem found

and exit 1.

## Expected

A title written into either template as an agent writes it, colon or not,
passes `grove check`. The owner decided, in the review of 2026-10-26:
"Quote the title in the templates", and accepted, in the round before
writing, quoting it in double quotes, `title: "<title>"`, in both the
change and the cycle template ("yes"). A title holding a double quote
then needs escaping; titles hold apostrophes far more often, which single
quotes would break.

## Reproduction

1. `just fixture <dir>`, `cd <dir>`, `grove init`, commit.
2. Create `grove/changes/261026-aaaaa-x/` holding `proposal.md` with
   `# Records: delete by ID`, and `.grove.yaml` with `title: Records:
   delete by ID`, `created`, `schema: default` and `base` set to `HEAD`.
3. Run `grove check`.

## The requirement it breaks

No maintained requirement at ea7fbaf covers what the templates print. The
owner agreed, in the round, that the fix stays observed behavior with no
change spec ("yes"); the test that reproduces the bug guards it.

## Open questions

None.

## Related documents

- Tasks: [tasks.md](tasks.md)
