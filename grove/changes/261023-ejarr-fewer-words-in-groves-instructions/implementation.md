# Implementation: Fewer words in Grove's instructions

<!-- Progress on change 261023-ejarr, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Say each rule once

Status: done

- Built on the files as cycle B and 261023-d6r79 leave them (this branch
  is stacked on 261023-d6r79's). At the start, `internal/instruct` held
  7,890 words of Markdown and YAML; after, 6,979.
- The opening text (`meanings` in `instruct.go`) is shorter and still says
  what a maintained spec, a change, observed behavior and a draft are,
  that merging is the only approval, and that Grove commits nothing; it
  now holds the one "commit only when the owner or the project's rules
  ask", taken out of every requires file.
- Each line in both a requires file and its default schema guide, or the
  bugfix propose guide, is kept in the requires file only; the pull
  request lists each and where its rule lives. Being over 7,000 words
  after that, the same files were reworded more tightly, as the proposal's
  decision says, dropping no rule.
- `internal/cli/walkthrough_test.go`: each assertion on a reworded line
  checks the same rule in its new words; the shared implement list also
  checks the commit sentence in the opening text. `go test ./...`: passed.
  `grove check` on Grove's repository: `No problems in grove/.`
- Words `grove instructions` printed on a ledger fixture with an approved
  change and a cycle, before and after, all shorter: explore 759 → 674,
  prototype 717 → 626, propose 3,201 → 2,719, plan 1,087 → 921, implement
  a change 3,108 → 3,003 and unattended 2,867 → 2,770, a cycle 863 → 772
  and unattended 623 → 540, archive 1,613 → 1,561.
- The trials before are 261023-kh3gg's baseline; after: see T-01's line in
  the pull request.

### T-02 — `just check` counts the instruction words

Status: done

- `justfile`: `check` counts the words of `internal/instruct`'s Markdown
  and YAML files, prints them beside the line counts, and fails above
  7,000 saying "over the word budget"; its comment names both budgets.
  `CONTRIBUTING.md`'s `just check` line names the line and word budgets,
  with no number.
- `just check`: `lines: 3877 Go, 1579 test; words: 6979 instruction`,
  passed. With a 100-word file added under `internal/instruct/requires/`:
  `words: 7079 instruction`, `over the word budget`, failed. With it
  removed: passed.

## Scenarios

None: the change has no change spec.

## Reconciliation

None: the change has no change spec.
