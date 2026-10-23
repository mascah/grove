## T-01 — grove run builds a cycle it is given

Result: `grove run <cycle-id>` behaves as [specs/runs/spec.md](specs/runs/spec.md)
"Start one attempt" and "Refuse a given cycle that cannot be built" say,
and the documentation specifies it.

Files: `internal/cli/cli.go`, `internal/cli/run.go`, `internal/cli/walkthrough_test.go`, `docs/running.md`, `docs/reference.md`, `dev/walkthrough.md`

Interfaces: `grove run [--retry] [--timeout <duration>] [<cycle-id>]`; `runAttempt` takes the given cycle ID, empty for none.

Constraints: plain `grove run` is unchanged; each refusal is one sentence
(rule 14); within the rule 12 budget.

Acceptance:
- `docs/running.md` has a section "Run a cycle you name" showing
  `grove run <cycle-id>` on a cycle not marked, and a refusal;
  `docs/reference.md` has a `grove run <cycle-id>` row and the usage line
  names the argument; `dev/walkthrough.md` runs the section (rule 8).
- `TestRun` in `internal/cli/walkthrough_test.go` gives an unmarked
  buildable cycle and sees one session on it while a marked cycle exists,
  and gives a cycle with a member under `needs_human` and sees a refusal
  naming it and no branch.
- `just check` passes.

## T-02 — Plan's finish names /grove-implement first

Result: `/grove-plan` finishes as [specs/cycles/spec.md](specs/cycles/spec.md)
"Say whether the cycle can run unattended" says, and the documentation
specifies it.

Files: `internal/instruct/requires/plan.md`, `docs/planning.md`, `dev/walkthrough.md`, `internal/cli/walkthrough_test.go`

Constraints: the plan instructions stay data under `internal/instruct`
(rule 10).

Acceptance:
- `docs/planning.md` shows the finish naming `/grove-implement <cycle-id>`,
  `grove run <cycle-id>` and the mark, in that order, for a cycle with
  nothing in the way and for one with a member that needs a person;
  `dev/walkthrough.md`'s planning steps match.
- `grove instructions plan` prints the finish in that order, checked in
  `internal/cli/walkthrough_test.go`; `just check` passes.
