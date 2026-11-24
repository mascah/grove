# Implementation: Propose writes every part of an agreed split

<!-- Progress on change 261124-mhdgj, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose writes every part of a split, against its dependency's text

Status: done

- `internal/instruct/requires/propose.md`: the split bullet now says to write
  every agreed part now, in this session, each later part against the
  requirements the earlier parts leave, and to defer no part because an
  earlier one rewrites the same requirements.
- `internal/instruct/schemas/default/documents/change-spec/guide.md`: a new
  bullet puts a requirement that a change in `depends_on` adds or modifies
  under MODIFIED, from that change's text, naming it, even with no maintained
  spec at `base`. `bugfix` has no change-spec guide of its own, so it prints
  the same text.
- `internal/cli/walkthrough_test.go` `TestPropose` checks the three split
  phrases and the guide bullet in `grove instructions propose`, and the guide
  bullet in `grove instructions propose --schema bugfix`. It passes, and fails
  on all of them with the two instruction files reverted.
- Disposable repository (`grove init`, the worktree's `bin/grove`): change A
  ADDS "Show a record" to `record-view`, which has no maintained spec; change
  B, `depends_on: [A]`, MODIFIES it from A's text, naming A. `grove check`
  prints "No problems in grove/."
- `just check` passes: 4,842 Go lines, 2,716 test lines, 8,458 instruction
  words.

Check: the guide's "no maintained spec" bullet said "Do not use MODIFIED or
REMOVED", which the new bullet contradicts; I reworded it to "Do not use
MODIFIED or REMOVED, but as the next bullet says". The wording and the phrases
the test checks were the routine choices the change's open question left to
the implementer.

## Scenarios

Whether an agent follows the instructions is a skill trial (`dev/trials.md`),
not run here; each line below shows the instruction text, and the agent's
part is unverified for that reason.

- propose "A requirement only the dependency adds": `internal/cli/walkthrough_test.go` TestPropose checks the guide bullet, passed; `grove check` on a disposable repository reports neither change, passed; the agent writing it unverified
- propose "A requirement the dependency modifies": TestPropose checks the guide bullet, passed; the agent writing it unverified
- propose "Two parts, one needing the other": TestPropose checks the offer to split and `depends_on`, passed; the agent writing it unverified
- propose "Parts that rewrite the same requirement": TestPropose checks "write every agreed part now", "against the requirements the earlier parts leave" and "Defer no part", passed; the agent writing it unverified
- propose "The owner keeps it whole": the split bullet ends "Otherwise write one", not checked by a test; unverified, agent behavior

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
