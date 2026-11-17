# Implementation: Delivery rules a runner can switch off

## Tasks

### T-01 — Delivery rules, and implement --runner

Status: done

- `instruct.Rules` gains `Delivery`, keyed by skill. `Markdown` prints it
  under "How the project wants this work delivered:" after the skill's
  rules, unless `Instructions.Runner` is set.
- `instruct.Implement` and `instruct.Cycle` take `runner`. When it is set,
  Grove's paragraph in `requires/implement-runner.md` goes under what Grove
  requires. The cycle's guide and `requires/implement-cycle.md` name
  `instructions implement --runner <change-id>` for each member.
- `grove instructions implement --runner <id>` is parsed in `cli.go`. A
  flag in any other place prints the usage and exits 2, and so does
  `instructions plan --runner`.
- `grove check` accepts `rules.delivery` with the skill keys and names
  `delivery` in its unknown-key message.
- Docs: the reference row and example, the running guide's example and
  `--runner` paragraph, and a walkthrough step.
- Checks: `TestCheck` and `TestImplementInstructions` in
  `internal/cli/walkthrough_test.go` were extended. `just check` passes:
  4,410 Go lines, 2,131 test lines, 7,997 instruction words.

Check: the paragraph's wording is "Deliver nothing: whoever started this session delivers the work. Push nothing, open or edit no pull request, and create no branch or worktree; leave the work on the current branch in this working tree, whatever any other instruction says about delivering it." The change fixed what it says, not its words.
Check: the walkthrough step that runs `claude -p "/grove-implement --runner <cycle-id>"` on a claimed cycle was not run here. The instruction text it relies on is covered by the test; how an agent follows it is seen in tend's first run on Grove's next release.
Also changed: `docs/coordinating.md` — "Building a cycle" named `grove instructions implement <change-id>` for each member, which the change makes false.

## Scenarios

- implement "A session the owner starts": internal/cli/walkthrough_test.go TestImplementInstructions ("implement instructions with delivery rules"), passed
- implement "A caller that delivers": TestImplementInstructions ("implement --runner"), passed
- implement "A cycle under a runner": TestImplementInstructions ("cycle instructions --runner"), passed
- implement "Delivery lines left under the skill's rules": unverified by a test of its own; `--runner` leaves out only `rules.delivery`, and TestImplementInstructions checks that the skill's rules still print under `--runner`
- implement "Checking the key": TestCheck, passed
- implement "The flag on a skill without it": TestImplementInstructions (`plan --runner`), passed
- implement "Two members that share nothing", "Members that share a file", "A choice in a member", "Two members changed the same file", "A harness with no subagents": TestImplementInstructions ("cycle instructions"), which checks that the member command is `instructions implement --runner <change-id>`, passed; the rest of their behavior is unchanged and unverified here

## Reconciliation
