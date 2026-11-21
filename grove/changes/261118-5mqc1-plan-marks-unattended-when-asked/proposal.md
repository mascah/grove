# Plan marks a cycle unattended when the owner asks

## Problem

The plan instructions say "Never write or remove `unattended`. It hands
the cycle to a `grove claim` given no cycle and is the owner's to set, in
a later merge of its own." (`internal/instruct/requires/plan.md` at
dc8d335). The cycles spec requires it too ("Mark a cycle unattended",
"Say whether the cycle can run unattended", `grove/specs/cycles/spec.md`
at dc8d335), and `docs/planning.md` says "No skill writes or removes the
line".

So an owner who asks plan for an unattended cycle gets a cycle file
without the mark and needs a second pull request to add it. Cycle
261118-5b98a went that way: the owner ran `/grove-plan Please plan the
next changes into an unattended cycle`, and the cycle took #172 to plan
and #173 to mark. The owner, afterwards: "you could have just written it
when creating the cycle".

The rule came from change 261022-1grb0, where the owner said "A cycle
should only be elligible for execution after it's explicitly been marked
as ready." (`grove/changes/archive/261022-1grb0-cycles-for-unattended-runs/proposal.md`).
Asking plan for it in the session is as explicit, and the owner still
reads and merges the file. `CONSTITUTION.md` rule 4 says only that
`unattended: true` in the target's cycle file hands it to a `claim` given
no cycle; it asks for no separate merge.

## Intended outcome

When the owner asks for an unattended cycle, in the description or in the
conversation, plan writes `unattended: true` in the cycle file, and
merging that one file plans the cycle and hands it to a `grove claim`
given no cycle. Revising a cycle, plan writes or removes the mark when the
owner asks and keeps it otherwise. Plan never offers the mark. Where the
owner asks nothing, everything stays as it is, and the finish still says
how to add the mark in a merge of its own.

## Scope

### In scope

- The plan instructions: when to write or remove `unattended`, and what
  the finish says when the file holds it.
- The cycles spec's "Mark a cycle unattended" and "Say whether the cycle
  can run unattended".
- `docs/planning.md`, `dev/walkthrough.md` and the walkthrough test.

### Out of scope

- `grove claim`, `grove check`, `grove changes` and `grove serve`: they
  read `unattended` as they do now.
- Propose re-planning the cycle of a member it revised: it still never
  writes or removes `unattended` (`grove/specs/propose/spec.md`, "Re-plan
  the cycle of a revised member", unchanged).
- Any amendment to `CONSTITUTION.md`: rule 4 already allows the mark in
  the merge that plans the cycle.

## Decisions

The owner asked for this change after cycle 261118-5b98a needed two pull
requests, and answered the round before writing with "yes to all" to each
recommended answer.

1. Plan writes `unattended: true` only when the owner's description or the
   conversation asks for an unattended cycle. It never offers or suggests
   the mark, and its finish still says how to add it later. Reason: no
   cycle runs unattended without the owner saying so.
2. Revising a cycle with `/grove-plan <cycle-id>`, plan writes or removes
   `unattended` when the owner asks and keeps it as it is otherwise.
   Reason: one rule for plan, not two.
3. Propose's re-plan still never writes or removes `unattended`. Reason: a
   re-plan made with a revision is not the owner asking for the mark.
4. Where the owner asked and a member waits on a change neither in the
   cycle nor archived as finished, or the cycle follows one not done,
   plan still writes the mark and names the obstacle. Reason: `grove
   claim` already skips such a cycle, and Grove's gates are only those of
   `CONSTITUTION.md` rule 10.
5. When the cycle file holds `unattended: true`, the finish says merging
   it hands the cycle to a `grove claim` given no cycle; otherwise it says
   what it says now.
6. This change names 261118-40s1a in `depends_on`, and so follows
   261118-emse7. Reason: all three edit
   `internal/instruct/requires/plan.md`, `dev/walkthrough.md` and
   `internal/cli/walkthrough_test.go`, and both are in cycle 261118-5b98a,
   marked unattended.
7. The commit is `feat`, not breaking: no project that works now needs an
   edit.
8. Routine, accepted together: two tasks, the instructions, spec and test
   in one and the docs and walkthrough in the other; at most 60 words of
   instruction text added; test layout and private helpers are the
   implementer's to choose.

## Open questions

None.

## Related documents

- Requirements: [cycles](specs/cycles/spec.md)
- Tasks: [tasks.md](tasks.md)
