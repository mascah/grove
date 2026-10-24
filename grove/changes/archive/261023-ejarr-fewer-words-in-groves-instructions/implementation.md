# Implementation: Fewer words in Grove's instructions

<!-- Progress on change 261023-ejarr, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Say each rule once

Status: done

- Acceptance, checked on 4bafd7f's code and #73's description: the
  description lists each removed line and where its rule still prints; the
  repeats left between a requires file and its guide are each one that
  can't go without rewording the copy that stays, and the description
  lists them ("What 801eda2 cut that stays"); the opening text is shorter
  and says what a maintained spec, a change, observed behavior and a draft
  are, that Grove commits nothing and that approval is the proposal being
  in the target branch; the per-skill table shows none longer after;
  trials before and after are in the description; `go test ./...` passes;
  `grove check` prints "No problems in grove/.".

The owner's decisions, 2026-10-24:

- Some rules a requires file states are still stated again in the same
  skill's guide, kept because removing them needed the surviving line
  reworded: in `skills/propose.md`, §3's numbered round answered by "yes"
  and its order by consequence, the `Delegated:` line and §5's check, all
  in `requires/propose.md`; in `skills/plan.md`, "It puts approved changes
  in a cycle", "Read each one" and "Ask about each change that may need a
  person"; in `skills/prototype.md`, "answers one question"; in
  `skills/archive.md`, landing in one merge. Asked whether to reword them
  away or revise T-01's acceptance to accept them, the owner said "Revise
  T-01's acceptance (recommended) preferred". T-01's result and second
  acceptance line, and the proposal's intended outcome and decisions, now
  accept a repeat that can't go without rewording the copy that stays.
- The budget is 8,000 words: asked on 2026-10-24 whether it becomes 8,000,
  the owner said "yes", and constitution rule 12 says so as PR #75 amends
  it. `internal/instruct` holds 7,705 words. The trials were rerun as a
  record in the pull request, not a gate.
- Two removals go beyond the proposal's listed cases: "It describes future
  work." in `skills/propose.md` and `bugfix/skills/propose.md`, and
  `requires/propose.md`'s first bullet ("A change describes future work. It
  is approved only when the owner merges it into the target branch; nothing
  you do approves it."), both covered by the opening text (`meanings`).
  Asked whether they stay out, the owner said "keep them out" (2026-10-24).

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
  8,000 saying "over the word budget"; its comment names both budgets.
  It failed above 7,000 until the owner agreed to 8,000 (2026-10-24).
  `CONTRIBUTING.md`'s `just check` line names the line and word budgets,
  with no number.
- `just check`: `lines: 3877 Go, 1579 test; words: 6979 instruction`,
  passed. With a 100-word file added under `internal/instruct/requires/`:
  `words: 7079 instruction`, `over the word budget`, failed. With it
  removed: passed.
- At 7,000, with T-01 redone without rewording, `just check` printed
  `words: 7705 instruction` and failed, saying `over the word budget`.
- At 8,000: `just check`: `lines: 3881 Go, 1579 test; words: 7705
  instruction`, passed. With a 300-word file added under
  `internal/instruct/requires/`: `words: 8005 instruction`, `over the word
  budget`, failed. With it removed: passed.

## Scenarios

None: the change has no change spec.

## Reconciliation

None: the change has no change spec.
