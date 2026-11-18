# The project index names requirements two changes both change

## Problem

Two approved changes may each rewrite the same requirement. Once the
first is archived, the second is stale on that requirement and must be
revised before it is built (`grove/specs/propose/spec.md` at f7e0140,
"Report a stale requirement"). In a cycle that holds both, that happens
in the middle of the build.

Nothing says so beforehand. `grove check` reads each open change's
MODIFIED and REMOVED requirements, by capability, to compare them with the
maintained spec (`internal/cli/check.go` at f7e0140, `stale` and `gives`),
and compares no change with another. The project index every skill reads
lists, for an approved change, its capabilities and what it waits on
(`internal/instruct/instruct.go` at f7e0140, `index`), so two changes
sharing a capability look alike whether or not they share a requirement.

In another project of the owner's, private, a plan of three cycles over
ten changes found four such pairs by reading every change spec, and cut
the cycles so that the second of each pair sits in a later cycle, "revised
once against the target when its turn comes instead of stalling
mid-cycle". The pairs were recorded only in that plan's pull request
description.

## Intended outcome

Where two approved changes, neither archived, each hold a change-spec
block for the same requirement of the same capability, and neither
depends on the other, the project index names, under each, the other
change, the capability and the requirement. Every skill's instructions
show it, since they share the index.

Plan, agreeing a cycle's members, tells the owner of each pair among them
and what it means, and the owner chooses whether both stay in the cycle
or one goes to a later one. Nothing is refused or reordered.

## Scope

### In scope

- The project index's line for an approved change that shares a
  requirement with another.
- The plan instructions telling the owner of such a pair among the
  members it proposes.
- The maintained cycles spec, `docs/`, `dev/walkthrough.md` and the tests.

### Out of scope

- A report in `grove check`, `grove changes` or `grove serve`: two
  approved changes that change one requirement are not a problem to fix,
  and `grove check` already reports the second once it is stale.
- Any gate or ordering: `grove claim` claims a cycle that holds such a
  pair, and plan lists what the owner agrees to (`CONSTITUTION.md` rule
  10).
- `/grove-implement <cycle-id>` starting two such members at once: it
  holds a member back only for a file its tasks share with one being
  built (`grove/specs/implement/spec.md`, "Build a cycle's changes in
  their own contexts"), and that stays as it is.
- Drafts: the index lists approved changes, and only those are compared.
- Two changes whose tasks name the same file: `grove check` and the
  implement instructions already deal with those.

## Decisions

The owner decided 1 while exploring, and answered the round before
writing, 2 to 8, with "yes" to each recommended answer.

1. Grove reports this itself, in plan's index, as a change of its own.
   Asked whether any of what a plan found should be acted on by Grove,
   the owner said "mostly a, a bit of c", where (c) was "some of it should
   be Grove's own check, not a note", and then "yes to rest" to it being a
   separate change reporting in the project index, not in `grove check`.
   Reason: no one has to find the pairs by reading or write them down,
   and two approved changes changing one requirement is not an error.
2. A pair is two approved changes, neither archived, whose change specs
   each hold an ADDED, MODIFIED or REMOVED block for the same requirement
   of the same capability. A pair where one depends on the other, directly
   or through other changes, is not listed. Reason: that pair is already
   ordered, and the later is written on the earlier's text.
3. It shows in the project index of every skill's instructions, as a line
   under each of the two naming the other change, the capability and the
   requirement. `grove check`, `grove changes` and `grove serve` are
   unchanged. Reason: the index is built once for all skills, and propose
   reads it to ask what a new change should wait on.
4. Plan tells the owner each pair among the members it proposes, and that
   the second built is stale on that requirement once the first is
   archived. The owner chooses the same cycle or a later one. Nothing is
   gated.
5. `/grove-implement <cycle-id>` starting two such members at once stays
   as it is.
6. This change names 261118-emse7 in `depends_on`. Reason: both edit
   `internal/instruct/requires/plan.md`, `internal/instruct/instruct.go`,
   the walkthrough and its test.
7. The pairs are found in `internal/doc` when it loads the changes, and
   the index in `internal/instruct` prints them.
8. Routine, accepted together: tests follow new walkthrough steps in
   `internal/cli/walkthrough_test.go`; this change adds at most 80 words
   of instruction text; test layout and private helpers are the
   implementer's to choose.

## Open questions

This one can wait until review.

- The words of the index's line. The owner agreed what it names, not how
  it reads. It decides how T-01 is built; the tasks suggest `also changed
  by <change-id>: <capability> "<requirement>"`, which `/grove-implement`
  builds and records for the reviewer unless the owner gives another.

## Related documents

- Requirements: [cycles](specs/cycles/spec.md)
- Tasks: [tasks.md](tasks.md)
