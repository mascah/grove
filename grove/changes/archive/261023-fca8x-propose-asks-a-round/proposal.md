# Propose asks a round

## Problem

The proposal guide (`internal/instruct/schemas/default/documents/proposal/guide.md`
at b143991) says what the author infers "is not a decision ... Put it
under open questions, or in a list headed 'Assumptions to confirm in
review'". Nothing says what confirms an assumption: `grove review` shows
the diff and the merge approves the change, so an assumption is confirmed
by nobody objecting. On the ledger fixture (session `d50a0a86`,
2026-10-23), the explore session closed with "these are details for
whoever implements it, not questions for you", propose wrote one of them
as an assumption, and the implement session built on it without asking,
citing "the proposal's 'can wait' assumption". The owner: "This leaves
ambiguity for the implement session."

## Intended outcome

After the owner confirms the scope and before anything is written, propose
lists every decision it would otherwise assume, numbered, each with a
recommended answer worded so that "yes" accepts it, and waits. Each answer
becomes a decision in the proposal, quoted. A decision the owner declines
to make becomes an open question, naming the task it blocks. The proposal
guide and template have no assumptions section.

## Scope

### In scope

- The round in Grove's propose requirements, for every schema.
- The proposal guide and template of the default schema, and the bugfix
  example in `docs/reference.md`.
- The documentation and walkthrough steps.

### Out of scope

- Explore, which keeps one question at a time: its frontier is not known
  until it has looked.
- Plan and prototype, which ask as they do today.
- What a task names: change 261023-6j5bw.

## Decisions

- **Propose asks a round; explore does not.** The owner, 2026-10-23,
  offered a round in propose, rounds in explore too, one assumption at a
  time, or keeping the section with implement asking: "1", after "I think
  I lean towards 2, which is sort of how mattpocock and superpowers skills
  already work I think" on the exploration's earlier framing. Reason: once
  the scope is confirmed, the open decisions are known and one message
  settles them, as mattpocock's `grilling` skill does; what is not settled
  is then an open question, which blocks the task it decides under apply
  and is asked under implement, instead of an assumption nobody confirmed.
- **No dependency on 261023-6j5bw.** The owner accepted, 2026-10-23, that
  both changes add to a new `propose` spec and whichever archives second
  adds to the spec the first created.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md) (no
  maintained spec at b143991)
- Tasks: [tasks.md](tasks.md)
