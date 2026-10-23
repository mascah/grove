# Implementation: Prototype before proposing

<!-- Progress on change 261022-ka8bx, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The prototype skill

Status: done

What changed:

- `internal/instruct/skills.yaml`: the `prototype` skill, between explore
  and propose, passing `"<args>"`; so `grove init` writes it and
  `rules.skills.prototype` is accepted.
- `internal/instruct/requires/prototype.md` and
  `schemas/default/skills/prototype.md` (in `schema.yaml`): what Grove
  requires of prototype and how it goes.
- `internal/instruct/instruct.go` `Prototype`, `internal/cli` `instructions
  prototype "<question>"`, the usage text and `init`'s skills comment.
- `requires/explore.md`: suggest a prototype when its answer could change
  the proposal, build nothing; explore's guide lists it as a result.
  `requires/propose.md`: write a prototype's findings in, citing branch and
  commit, linking no file outside the change and the target.
- Docs: getting started §2 "Try it first"; reference (command, six skills,
  `rules`); concepts "The loop"; walkthrough step 1, step 2, "Prototype
  before proposing", "Commands", "Agents from the repository".
- Tests: `TestPrototype`; `TestInit`, `TestHarnesses` and `TestExplore`
  updated.

Checks:

- `just check`: passes; 2855 Go lines, 891 test lines.
- On the ledger fixture, Claude Code headless, resumed per turn, grove at
  this branch: `/grove-explore Monthly budgets in report, from the README's
  planned list`; when I said I did not know how a budget should read, it
  named the question and asked whether to prototype it, with nothing
  written. On yes it built `prototypes/budget-report-format/show.py` on a
  new `prototype/budget-report-format` branch from `budgets`, ran under, at
  and over budget cases, and reported the command, observations, a
  verdict and its limits; on my yes it committed `1cf7433` there only.
  Asked to propose, it switched back to `budgets`; the proposal states the
  finding and cites `prototype/budget-report-format` at `1cf7433`, linking
  no prototype file. `git diff main prototype/budget-report-format --
  grove` is empty and main is unchanged.
- Main line step 2's explore, on a fresh branch: no prototype suggested,
  nothing written.
- `grove init` on the fixture writes `grove-prototype` for claude and codex;
  `rules: skills: prototype: [...]` passes `grove check` and is printed with
  the prototype instructions.

## Reconciliation

- prototype: applied to specs/prototype/spec.md
