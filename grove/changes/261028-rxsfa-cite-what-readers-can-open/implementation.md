# Implementation: Cite what the project's readers can open

<!-- Progress on change 261028-rxsfa, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Every skill is told to cite what readers can open

Status: done

- `internal/instruct/instruct.go`: the `meanings` constant gains, after the
  paragraph on drafts, the sentence agreed in the proposal: "Write Grove
  documents for whoever reads this project: cite only what they can open,
  and quote or summarise what only this machine holds, such as an unpushed
  branch, a file outside the checkout or a private project."
- `bin/grove instructions explore`, `bin/grove instructions propose` and
  `bin/grove instructions implement 261028-rxsfa` each print it at line 17,
  before `## What Grove requires of <skill>` at line 24.
- `just check` passed: 4024 Go lines, 1681 test lines, 7997 instruction words.

### T-02 — A proposal carries a prototype's findings

Status: done

- `internal/instruct/requires/propose.md`: the prototype line now has the
  agent write the findings into the proposal so it stands without the
  prototype branch (what was tried, what was observed and the verdict),
  citing the branch and commit only as where they came from, or saying the
  experiment was not committed. `bin/grove instructions propose` prints it.
- `docs/concepts.md`, `docs/getting-started.md`, `dev/walkthrough.md`: the
  prototype steps say the proposal states the findings and names the branch
  and commit as their source.
- `internal/cli/walkthrough_test.go` `TestPrototype` checks the propose
  instructions for the phrase "citing the prototype branch and the commit".
  That file is on no task's `Files:` line, so it stays as it is; the
  propose line keeps that phrase, on one line, which the requirement's
  "naming the branch and its commit only as where they came from" allows.
  The test would otherwise have needed its phrase updated.
- `just check` passed: 4024 Go lines, 1681 test lines, 8023 instruction words.

## Scenarios

- instructions "Every skill is told": `bin/grove instructions explore`,
  `propose` and `implement 261028-rxsfa` each print the sentence at line 17,
  before the skill's own instructions at line 24, passed. No test pins the
  sentence, as the proposal decided.
- instructions "A note in a private project": unverified; it is how an agent
  writes a proposal, judged by a person reading it or by a skill trial.
- instructions "A source the readers share": unverified; as above, an
  agent's choice judged by a person.
- prototype "Proposing after a prototype": unverified; the propose
  instructions say it (`internal/cli/walkthrough_test.go` `TestPrototype`
  checks their citing phrase, passed), but what the proposal states is
  judged by a person running the walkthrough's "Prototype before proposing".
- prototype "An experiment not committed": unverified; the propose
  instructions say to say the experiment was not committed, but the
  proposal is judged by a person.

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
