# Delivery rules a runner can switch off

## Problem

Grove gives a project only one place to say how implemented work leaves
the working tree: `rules.skills.implement`. Its documentation tells
projects to put "push the branch to origin and open a pull request" there
for unattended runs (`docs/running.md` at f18de5b, "Tell the agent to
commit and open a pull request"). Its runs spec says the agent pushes and
opens pull requests "only as the project's `rules` tell it"
(`grove/specs/runs/spec.md` at f18de5b, Boundaries).

Those lines reach every implement session. Grove prints them in the
cycle's instructions and in each member's, which a subagent reads
(`internal/instruct/instruct.go:634` at f18de5b). A caller that delivers
the work itself, such as a runner that reviews the branch before opening
the pull request, has no way to keep them out.

One such runner is tend, which claims a cycle, has it built, has the head
reviewed and then publishes it. On 2026-11-16 it started
`/grove-implement` on a claimed memex cycle of three changes. memex's
`rules.skills.implement` said to work in a worktree per change and open a
stacked pull request for each. The agent followed them: it created nested
worktrees, pushed three branches and opened three pull requests itself,
and left the attempt's branch empty, so the runner's review never ran. The
full account is in tend's
`docs/work/261116-runner-contract-and-pipeline/design.md`, §1 and §2.

## Intended outcome

A project keeps how it wants work delivered apart from how the work is
done, under `rules.delivery`. A session the owner starts prints both, as
today. A caller that delivers the work itself runs implement with
`--runner`. Grove then leaves the delivery lines out and says, in its own
words, that this session delivers nothing.

A cycle's subagents always get that paragraph, because the session that
started them delivers the work. Projects that change nothing see the same
instructions as before, except that the cycle's member commands now carry
`--runner`.

## Scope

### In scope

- A `rules.delivery` key in `grove.yaml`, keyed by skill like
  `rules.skills`. Its lines print in that skill's instructions, after the
  skill's rules, under their own heading.
- `--runner` on `grove instructions implement`, for a change or a cycle.
  It leaves out `rules.delivery.implement` and prints Grove's fixed
  paragraph instead.
- The cycle's instructions naming `instructions implement --runner
  <change-id>` as the command for each member.
- `grove check` accepting `delivery` under `rules`, with the same skill
  keys as `skills`.
- Documentation:
  - the usage line;
  - `docs/reference.md`;
  - `docs/running.md`'s example, moved to `rules.delivery.implement`;
  - a walkthrough step;
  - the runs spec's Boundaries line, reconciled at archive.

### Out of scope

- Moving any project's existing lines; memex moves its own.
- tend's side: passing `--runner` and checking the result. That is tend's
  own change.
- `--runner` for skills other than implement; `propose --unattended` is
  unchanged.
- Who merges, and anything Grove writes to GitHub (`CONSTITUTION.md`
  rule 11).
- A warning for delivery-like lines left under `rules.skills`: Grove
  cannot classify free text.

## Decisions

All were accepted by the owner with "yes" to the round asked before
writing, on 2026-11-16.

1. **`--runner` means "whoever started this session delivers the work".**
   Grove leaves out `rules.delivery.implement` and prints a fixed paragraph
   saying to deliver nothing: no push, no pull request, no branch or
   worktree, and the work left on this branch in this working tree.
   - **Committing.** The paragraph says nothing about committing. That
     stays with Grove's existing instructions, the project's rules and the
     caller's prompt.
   - **Reason.** Saying "commit here" would contradict the cycle
     instructions, which tell a member's subagent to neither archive nor
     commit.
2. **A cycle's members are always listed with `--runner`**, whether or not
   the cycle's own session runs under one.
   - **Reason.** A subagent never delivers; the session that started it
     does.
   - **Effect.** An owner-run cycle's subagents stop seeing the delivery
     lines, which is the intent.
3. **`rules.delivery` is keyed by every skill, like `rules.skills`.** Only
   `instructions implement` accepts `--runner` for now.
   - **Reason.** It mirrors an existing key and its check. Only implement
     is started by a caller that delivers.
4. **Grove's documented recipe keeps the agent delivering.**
   `K=$(grove claim) && claude -p "/grove-implement $K"` stays without
   `--runner`, and its push-and-PR line moves to
   `rules.delivery.implement`.
   - **Reason.** With no runner, the agent is the deliverer. Only a caller
     that delivers passes `--runner`.
5. **Lines left under `rules.skills.implement` print as today**, with no
   warning.
   - **Reason.** Existing projects keep working. Grove cannot tell a
     delivery line from a work line.
6. **The requirements go in the implement spec.**
   - A new requirement sits beside "Rules for implement", and "Build a
     cycle's changes in their own contexts" is modified to name the flag.
   - The runs spec's Boundaries line becomes: "the agent pushes and opens
     pull requests only as `rules.delivery` tells it, and not at all under
     `--runner`".
7. **Routine choices:**
   - one task;
   - the paragraph lives under `internal/instruct/requires/`, so that no
     project file replaces it (rule 10);
   - one walkthrough test;
   - test layout and private helpers are delegated.

The constitution needs no amendment:
- rule 1 still has Grove launch nothing;
- rule 2 leaves starting, bounding and stopping the agent to the caller,
  which this extends to delivery;
- rule 11 still keeps Grove from pushing or writing to GitHub.

The documented step rule 8 asks for is the reference row, the running
guide and the walkthrough step.

## Open questions

None.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
