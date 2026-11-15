# Propose drafts a change with no one in the session

## Problem

Since 0.10.0 `/grove-implement` asks no one and ends at a branch ready
for review (`grove/specs/implement/spec.md` at efdb8af). Proposing still
needs the owner in the session: `internal/instruct/requires/propose.md`
at efdb8af has the agent confirm the scope, ask one round and "Wait for
the owner", and `grove/specs/propose/spec.md` requires it ("Ask a round
before writing").

The owner wants agents to draft proposals too, from a tracker issue, a
review finding or a problem `grove check` reports, and to read each draft
as a pull request: "I need something that is also creating proposals."
Those instructions have a session nobody attends wait at its first
question, before it writes any document.

A prototype on 2026-11-14 showed the building half already runs this way.
One script, keeping no state of its own, took a cycle of one on the
ledger fixture through `grove claim`, a Claude Code session running
`/grove-implement` in a terminal pane, a Codex review of the head commit,
one repair and a second review, in three runs of the same line, with no
change to Grove. The defect it repaired was planted by hand. It came from the branch `prototype/herdr-loop`, commit
4c5804d, which is never merged. Nothing like it can start earlier,
because nothing writes the proposal.

## Intended outcome

`/grove-propose --unattended <request>` writes one new change from the
request, on a branch, asking no one and waiting for no one, and ends with
the branch ready for the owner to read. The round it would have asked is
in the proposal as open questions, each with the answer it recommends;
Decisions holds only what the request itself says. A request that names
no outcome gets no change. The draft is approved, like any other, by the
merge that brings it into the target, and nothing builds it before then.
An open question the owner merges unanswered is built as
`/grove-implement` builds any other: its best answer with a `Check:`
line, or a blocked task where a wrong answer would be costly.

`/grove-propose` without the flag is as it was.

This rests on `CONSTITUTION.md` rule 4 as
[#163](https://github.com/mascah/grove/pull/163) amends it, and is merged
only after that pull request.

## Scope

### In scope

- `--unattended` on `grove instructions propose` for a new change, as its
  own argument or as the first word of the quoted description, which is
  how a skill or command passes what the user typed.
- The instructions that mode prints, in `internal/instruct/requires/`.
- A section of `docs/running.md`, the reference's row, a walkthrough step
  and one skill trial.

### Out of scope

- Revising a drafted change with no one in the session, such as from the
  owner's comment on its pull request: a later change.
- Starting the session, reading a tracker, pushing the branch or opening
  the pull request: Grove launches no agent and writes to no tracker
  (rules 1 and 11). Whoever starts the session passes the request in,
  and a line under `rules.skills.propose` has the agent commit and push.
- `/grove-explore`, `/grove-prototype` and `/grove-plan`, which stay
  conversations with the owner.

## Decisions

The owner answered the round on 2026-11-14 with "yes", accepting each
recommended answer:

1. The mode is asked for as `/grove-propose --unattended "<request>"`,
   and attended propose is unchanged, so nothing breaks. A prompt, script
   or timer that types the flag is the only thing that reaches it.
2. New changes only. "Revising a draft from your PR comment is a
   separate, later change."
3. "Every decision of the round becomes an open question, in the round's
   order, with a recommended answer and the task it decides. Decisions
   holds only what the request itself says, quoted." An agent records
   nothing it inferred as agreed.
4. "A request that does not say what outcome is wanted gets no change.
   The agent writes nothing and says what is missing." A vague issue
   costs the owner no review.
5. One request, one change: "it never splits, adds `depends_on`, or
   switches schema on its own. It names each of those as an open question
   instead." A `depends_on` or a `--schema` the request or its caller
   names is theirs to give.
6. It ends "at a branch ready for review, with no offer to run
   `/grove-implement` on that branch": rule 4 keeps the one merge that
   approves a draft with its code for a change drafted with the owner in
   the session.
7. It works on a branch other than the target; it commits, pushes and
   opens a pull request only as `rules.skills.propose` say; it is
   documented as a section of `docs/running.md`, with a walkthrough step
   and one new skill trial; and its instruction text stays under 250
   words (rule 12, at 7,704 of 9,000 at efdb8af).

## Open questions

Neither was in the round. Each can wait for the review of this draft,
and is written here as the change itself would have it.

1. T-01: what `--unattended` does with a change's ID or with no request.
   Recommended: one sentence and a non-zero exit, as the change spec's
   "Refuse an unattended propose it cannot write" says. An ID is decision
   2; a request is what rule 4 says the draft's decisions come from.
2. T-01: where the mode's text lives. Recommended: one file,
   `internal/instruct/requires/propose-unattended.md`, printed with the
   propose instructions when the flag is given and saying which of them
   it replaces. A project's schema can replace a propose guide but not
   what Grove requires (rule 10), so the mode holds under any schema.

## Related documents

- Requirements: [propose](specs/propose/spec.md),
  [schemas](specs/schemas/spec.md)
- Tasks: [tasks.md](tasks.md)
