# Implementation: Fewer words in Grove's instructions

<!-- Progress on change 261023-ejarr, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Say each rule once

Status: in progress

Left: the trials before and after, rerun with the fixed script (below).

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
- `dev/trials/run.sh`, added to T-01's `Files:` with the owner's
  agreement: case (a)'s setup kept a spec `.grove.yaml` that this branch's
  Grove reports, because its check read `grove check`'s pipe status under
  `pipefail`; it now reads the output. `SETUP_ONLY=1` case (a): clean.
- `dev/trials/run.sh` and `dev/trials.md`, both added to T-01's `Files:`
  with the owner's agreement. In the first "after" run, case (a) in Codex,
  the scripted owner answered propose's finish ("either merge … or run
  `/grove-implement`") with the fallback reply, so the agent built and
  archived the change in the propose session. The waiting check now uses
  Sonnet and gives no reply to next steps left to the owner (Haiku's
  answers varied between runs); a run needing an earlier run's change ID
  fails the trial if there is none.
- Trials: the baseline (261023-kh3gg) and nine "after" trials used the old
  waiting check, which also took two commit offers for the owner. The
  owner chose to finish in a new session; the comparison is to be rerun
  with the fixed script, before (on #70's Grove) and after (this branch).

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
