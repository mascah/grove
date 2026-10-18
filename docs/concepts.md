# Concepts

## Three things Grove keeps apart

- **A maintained spec** says how part of your system should behave. It lives
  at `grove/specs/<capability>/spec.md` on your main branch, and it is
  agreed: someone merged it.
- **A change** is a proposal for future work. It lives in
  `grove/changes/<id>-<slug>/`. Once it is merged into your main branch it is
  approved, which means agreed, not built.
- **What the code does** is evidence. Where it differs from a spec, the spec
  says what was agreed and the code shows what happens; neither silently
  wins.

Grove's instructions tell your agent which is which, so a conversation never
treats a draft as agreed or an approved change as done.

## Approval is a merge

Grove has no approve button. A change is approved when its proposal is in
your main branch, which you choose in `grove.yaml` as `target`. Committing,
saving or ending a conversation approves nothing. Merge commits and squash
merges both count. If you commit a proposal straight onto your main branch,
that commit approves it.

`grove review` shows what merging the branch you are on would approve, before
you merge it.

## A change

A change folder holds:

| File | Written by | Says |
| --- | --- | --- |
| `.grove.yaml` | propose | Title, when it was created, and an optional one-line description. |
| `proposal.md` | propose | Why, what changes, decisions with who made them, open questions. |
| `specs/<capability>/spec.md` | propose | Requirements the change adds, changes or removes, for one capability each. |
| `tasks.md` | propose | The work as numbered tasks, each with how to check it is done. |
| `implementation.md` | implement | Each task's status and evidence, and which specs were updated. |

`tasks.md` is part of what was approved, so applying it never edits it;
progress goes in `implementation.md`.

## The loop

1. **Explore.** A conversation with your agent about what you want and what
   is true today. It writes nothing unless you ask.
2. **Propose.** The agent writes a change from the conversation or from a
   description.
3. **Review and merge.** You read the change and merge it. Now it is
   approved.
4. **Apply.** The agent builds the change on a branch and records
   progress.
5. **Archive.** The agent folds the change's specs into the maintained specs
   and moves the change to `grove/changes/archive/`. You merge code, specs and
   archive together.

A change you decide not to build is archived too, without touching the
specs; its `implementation.md` says `Status: dropped` and why.

## Two halves

The `grove` command runs in your terminal: it sets up, lists, checks, reviews
and reads. Skills run in your agent's chat: each one runs
`grove instructions <skill>` and follows what it prints. Grove starts no
agent and makes no commits, so it works with the agent you already use, in
the checkout and branch you already have.

## IDs

A change's ID is the date and five characters, such as `261015-k3m9q`. Its
folder adds a short name: `261015-k3m9q-let-users-rename-projects`. A
capability's ID is a short lowercase name, such as `projects`.

## Glossary

- **Approved:** merged into the target branch.
- **Archive:** fold a change's specs into the maintained specs and move it to
  `changes/archive/`.
- **Capability:** one area of behavior, with one maintained spec.
- **Change:** a folder proposing future work.
- **Context:** project facts in `grove.yaml` given to every skill.
- **Harness:** the agent you run, such as Claude Code or Codex.
- **Maintained spec:** agreed behavior of one capability.
- **Skill:** a chat command Grove installs for your harness.
- **Target:** the branch whose contents count as agreed, usually `main`.
