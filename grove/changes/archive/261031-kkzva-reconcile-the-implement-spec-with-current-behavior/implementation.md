# Implementation: Reconcile the implement spec with current behavior

<!-- Progress on change 261031-kkzva, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The text around implement's requirements says what it covers today

Status: done

implement's Purpose and Boundaries are the proposal's, word for word;
init's Boundaries ends at config, naming no "No apply skill". `just check`
passed (gofmt, vet, every test; 4,024 Go lines, 1,681 test lines, 8,023
instruction words), and `grove check` reports no problem.

## Scenarios

"Instruction text" scenarios describe what an agent does; the test named
checks the instructions print the rule, and the agent obeying it needs an
agent session, as the skill trials in `dev/trials.md` run. Tests are in
`internal/cli/walkthrough_test.go`; all passed.

- implement "An approved change": unverified; it needs an agent session (TestImplementInstructions checks the approved change's instructions in both modes, and `requires/implement.md` says to work on a branch other than the target and archive when every task is done)
- implement "A draft": TestImplementInstructions (draft instructions; unattended prints "main holds no approved change") and TestOneMerge, passed
- implement "A change that waits": TestCoordinating (approved, both modes) and TestImplementInstructions (draft), passed
- implement "A cycle": TestImplementInstructions (cycle instructions in both modes, a member's place in its plan, `needs_human` ignored) and TestCoordinating ("change 2 of 2 in cycle"), passed
- implement "A member left out": TestCoordinating ("Leave out … it waits on 260101-zzzzz, which is neither being built nor archived as finished."), passed
- implement "A large change": TestImplementInstructions (documents named by path, none of their text, the implementation template and no other), passed
- implement "A document that differs": TestImplementInstructions (`tasks.md` named as differing, with its `git show`), passed; which text the agent builds is unverified, needing an agent session
- implement "An open question, with the owner": unverified; it needs an agent session with the owner (TestImplementInstructions checks "ask the owner and wait")
- implement "Feedback changes the change": unverified; it needs an agent session with the owner (TestImplementInstructions checks the edit goes through `instructions propose`; telling the owner how to try a task and setting a changed task back to in progress are in `requires/implement-owner.md`, which no test asserts)
- implement "The approach rests on something false": unverified; it needs an agent session (TestImplementInstructions checks "stop that task and build nothing more of it")
- implement "Unattended" (Settle what the change does not): unverified; it needs an unattended session (TestImplementInstructions checks "ask no one and wait for no one's answer"; editing none of the change's files but `implementation.md` is in `requires/implement-unattended.md`, which no test asserts)
- implement "A delegated path": unverified; it needs an agent session (TestImplementInstructions checks the `Delegated:` line is honored with the owner)
- implement "A file nothing names": unverified; it needs an agent session (TestImplementInstructions checks the `Also changed:` line with the owner's agreement)
- implement "Unattended" (Change only the files the tasks name): unverified; it needs an unattended session (TestImplementInstructions checks `Delegated:` ignored and no `Also changed:` line)
- implement "Changed since base": TestMovedFiles (named in both modes), passed; asking or blocking is unverified, needing an agent session (TestImplementInstructions checks "mark that task `Status: blocked`, naming the file")
- implement "Its dependency landed": unverified; it needs an agent session with the owner (TestImplementInstructions checks "recommend propose, which re-bases the change")
- implement "A member of a cycle": TestMovedFiles (named since the plan, not since base) and TestFollows (nothing named once the cycle depended on lands, or once re-planned), passed
- implement "A rule for both modes": TestImplementInstructions, passed
- implement "The owner agrees": unverified; it needs an agent session with the owner (TestImplementInstructions checks "from one correction, quoting it")
- implement "Unattended" (Rules for implement): unverified; it needs an unattended session (TestImplementInstructions checks "Do not add it to `grove.yaml`.")
- implement "A draft, before and after building": TestApprove (as future work) and TestOneMerge (with its code), passed
- implement "Planned, agreed and named nowhere": TestPlanned, passed
- implement "Every scenario has a test": unverified; it needs an agent session (TestImplementInstructions checks the verify step's text)
- implement "A done task later broken": unverified; it needs an agent session (TestImplementInstructions checks "Before you finish, run the project's tests, re-read each done task's")
- implement "A scenario only a person can judge": unverified; it needs an agent session (TestImplementInstructions checks the report names unverified scenarios, and TestArchive that archive names them without stopping)

## Reconciliation

- implement: applied to specs/implement/spec.md
