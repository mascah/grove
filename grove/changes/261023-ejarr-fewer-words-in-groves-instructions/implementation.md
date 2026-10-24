# Implementation: Fewer words in Grove's instructions

<!-- Progress on change 261023-ejarr, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Say each rule once

Status: in progress

Left: waits on the rule 12 amendment. `internal/instruct` holds 7,705
words, over the 7,000 budget, and the owner chose to raise the budget in a
constitution PR of its own rather than reword (2026-10-23). Then the trials
are rerun as a record in the pull request, not a gate.

- Redone as repetition removal only, on the files as 261023-d6r79 leaves
  them (801eda2 had also reworded them; the owner separated the two). A
  line goes only where the same rule already prints for the same skill,
  and the surviving copy keeps d6r79's words; `git diff --word-diff
  implement/261023-d6r79 -- internal/instruct` adds no word, only the
  punctuation where a sentence was cut. Where a removal needed the survivor
  reworded, both stay. The pull request lists each removed line and where
  its rule still prints.
- The commit sentence stays only in the opening text (`meanings`), which
  every skill prints: removed from `requires/explore.md`, `propose.md`,
  `plan.md`, `archive.md` (with its "commits nothing"), `implement.md` and
  `implement-cycle.md`. The opening text drops its repeat of its own
  approval bullet ("merging into the target branch is the only approval").
- Words `internal/instruct` holds: 7,890 at d6r79, 7,705 after (`just
  check`).
- Words `grove instructions` prints, d6r79's build → this branch's, all
  fewer: explore 743 → 708, prototype 703 → 694, propose 3,176 → 3,073,
  plan 1,065 → 1,046, implement a change 2,876 → 2,838 and unattended
  2,635 → 2,597, a cycle 844 → 825 and unattended 604 → 585, archive
  1,384 → 1,363. Fixture: `just fixture`, `grove init`, committed; change
  "Allow editing older records" proposed on a branch from `instructions
  propose` with the walkthrough test's files and merged with `--no-ff`;
  cycle `261022-m0000` listing it, `planned` at that merge, committed on
  `main`; then on a new branch, `grove check` clean, each of `instructions
  explore`, `prototype "Can an older record keep its category?"`,
  `propose`, `plan`, `implement <change>` and `<cycle>`, each with and
  without `--unattended`, and `archive <change>`, piped to `wc -w`.
- Each surviving rule checked in the printed instructions of its skill on
  that fixture, and each removed line checked gone.
- `internal/cli/walkthrough_test.go`: d6r79's assertions, with those on a
  removed line or a rewrapped one checking the surviving text; the shared
  implement list checks the commit sentence in the opening text. `go test
  ./...`: passed. `just check`: fails only on the word budget, as expected.
- `dev/trials/run.sh` and `dev/trials.md`, added to `Files:` with the
  owner's agreement (636c03c), and `dev/trials/b.sh`, added the same way
  for the fix below. Script fixes: case (a)'s setup check reads `grove
  check`'s output; the waiting check uses Sonnet and gives no reply to a
  next step left to the owner; a run needing an earlier run's change ID
  fails the trial if there is none; `reply` matches the case's answers
  against all the agent's messages of the turn, not only the last, since
  Codex can split a turn; case (b)'s pattern no longer matches "`delete`
  would reuse that lookup", only a question on a deleted record's ID.
- Trials so far measured 801eda2's reworded text, not this one: the
  three-run comparison found unneeded questions even (34 → 34) and missed
  decisions 2 → 4, all in case (b), partly from the two script bugs fixed
  above. Details are in the pull request.

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
- With T-01 redone without rewording, `just check` prints `words: 7705
  instruction` and fails, saying `over the word budget`, until the rule 12
  amendment raises the budget.

## Scenarios

None: the change has no change spec.

## Reconciliation

None: the change has no change spec.
