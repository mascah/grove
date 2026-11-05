# Implementation: A cycle builds each change in its own context

<!-- Progress on change 261105-pr77s, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — A cycle builds each change in its own context

Status: done

- `internal/instruct/requires/implement-cycle.md`: the cycle's instructions,
  in both modes, build each member in a fresh context ("such as a
  subagent", no harness named) by `instructions implement <change-id>
  --unattended`, which neither archives nor commits; start at once members
  whose `after` changes are archived and whose `Files:` share no file with
  a member being built; only the session archives and commits, one at a
  time; where two members built together changed the same file, run the
  tests and, if they fail, archive the earlier-started one and rebuild the
  other on its own; with the owner, put returned checkpoints to the owner
  at once and rebuild with their words, unattended leave them for the
  report; with no fresh context, build one after another by the command
  the list names.
- `internal/instruct/instruct.go`: `Cycle` ends a member's line
  `after <ids>`, naming the listed members it depends on; a member that
  depends on none has no `after`. The guide line names the command for
  building members yourself.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks the new
  phrases in both modes of the cycle's instructions; `TestCoordinating`
  checks the cycle of two lists the second `after` the first and the first
  with no `after`.
- `docs/coordinating.md` ("Building a cycle unattended") and
  `dev/walkthrough.md` (Coordinating, step 8) explain and run it.
- `just check`: passed; 4,057 Go lines, 1,758 test lines, 8,614 words of
  instruction text (+188 over 8,426).

Also changed: `docs/reference.md` — its row for `grove instructions implement <cycle-id>` said each change is built one at a time, which T-01 makes false

## Scenarios

- implement "Two members that share nothing": `internal/cli/walkthrough_test.go` TestImplementInstructions checks the instructions say to start members at once in fresh contexts and archive and commit one at a time; TestCoordinating checks the `after` marks; passed. That an agent does it is unverified: it needs a session with subagents (`dev/trials.md`).
- implement "Members that share a file": TestImplementInstructions checks the start rule on `Files:`, passed; the agent's behaviour unverified, as above.
- implement "A checkpoint, with the owner": TestImplementInstructions checks the owner-mode text on returned checkpoints, passed; the agent's behaviour unverified, as above.
- implement "Two members changed the same file": TestImplementInstructions checks the test-and-rebuild text, passed; the agent's behaviour unverified, as above.
- implement "A harness with no subagents": TestImplementInstructions checks the fallback text and the mode's command, passed; the agent's behaviour unverified, as above.
- implement "An approved change": TestImplementInstructions, passed (unchanged).
- implement "A draft": TestImplementInstructions (the draft case), passed (unchanged).
- implement "A change that waits": TestCoordinating, passed (unchanged).
- implement "A cycle": TestImplementInstructions and TestCoordinating check the `after` mark and each member's place in the plan, passed.
- implement "A member left out": TestCoordinating ("leaves out a change that waits"), passed.

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
