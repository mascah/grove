## T-01 — A cycle builds each change in its own context

Result: the cycle's implement instructions, in both modes, say what
"Build a cycle's changes in their own contexts" and the modified "Build a
change, a draft or a cycle" in [specs/implement/spec.md](specs/implement/spec.md)
say, the cycle's list marks each member `after` the members it depends
on, and the docs and walkthrough explain and run it.

Files:
- `internal/instruct/requires/implement-cycle.md`
- `internal/instruct/instruct.go`
- `internal/cli/walkthrough_test.go`
- `docs/coordinating.md`
- `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: each line of the list `grove instructions implement
<cycle-id>` prints ends `after <ids>` when the member depends on other
members of the cycle.

Constraints: rule 12's instruction-text budget, adding at most about 250
words; the instructions name no harness ("a fresh context of your own,
such as a subagent"); `CONSTITUTION.md` rules 1, 2 and 9 as amended by
#133.

Acceptance:
- In the walkthrough's cycle of two, where the second depends on the
  first, the cycle's instructions list the second `after` the first, and
  a member that depends on nothing has no `after`.
- The cycle's instructions, with the owner and with `--unattended`, say
  to build each member in a fresh context by `grove instructions
  implement <change-id> --unattended` that neither archives nor commits,
  to start at once members whose dependencies are archived and whose
  `Files:` share no file with a member being built, to archive and commit
  each returned member one at a time, to run the tests and rebuild the
  later-started member where two changed the same file, how checkpoints
  reach the owner or the report, and to build one after another where the
  harness has no subagents.
- `docs/coordinating.md` explains building a cycle this way, and
  `dev/walkthrough.md` has the step the test follows.
- `just check` passes within the word budget.
