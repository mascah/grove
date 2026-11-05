# A cycle builds each change in its own context

## Problem

`internal/instruct/requires/implement-cycle.md` (bc9c29b) has one session
build every change of a cycle "in the order given, one at a time", so that
session carries the whole history of every change, and nothing is built
in parallel even where `depends_on` allows it. The owner: "I've had some
sessions that have run over days and compacted many times over a single
cycle. It wasn't necessarily an issue other than it probably took more
nudges from me than it really needed", and "If a cycle has 10 changes in
it, we shouldn't have to wait for every task of every change to be done
in order unless that's necessary. Its also unclear if we should be
bloating the context with history of every change."

Change 261105-va9qn already keeps each change's state in its
`implementation.md`, so a fresh context can pick a change up from
`grove instructions implement <change-id>` without the conversation.

## Intended outcome

A session building a cycle, with the owner or unattended, builds each
change in a fresh context of its own where the harness offers one (in
Claude Code, a subagent), starts at once every change whose dependencies
are archived and whose tasks' `Files:` share no file with a change being
built, and itself only holds the cycle's list and each change's report:
it archives and commits each change one at a time and is the only one
that talks to the owner. Where the harness has no subagents, the session
builds each change itself, one after another, as today.

## Scope

### In scope

- The cycle's implement instructions, in both modes, and the `after`
  marks in the cycle's list.
- The docs and walkthrough steps for building a cycle.

### Out of scope

- Launchers such as GitHub agentic workflows (gh-aw), one job per change:
  prototyped next, and any guide or `grove changes --json` field in a
  change of its own.
- Worktrees or any workspace (`CONSTITUTION.md` rules 2 and 11); a
  subagent per task; `grove run` starting more than one session.

## Decisions

1. **A fresh context per change by default.** Asked whether by default or
   only on a project's opt-in, the owner chose "By default, where the
   harness can": fresh context per change, ready changes at once when
   their `Files:` do not overlap, otherwise one after another.
2. **The constitution allows it.** The owner chose to amend rules 1, 2
   and 9 so the user's agent may start its own subagents where Grove's
   instructions ask (pull request #133, no code). This change needs it in
   the target first.
3. **The main session archives and commits**, one change at a time, after
   each subagent returns (the owner chose it over subagents archiving),
   so no two writers touch the maintained specs or Git's index at once.
   Subagents neither archive nor commit.
4. **A file beyond a task's `Files:`.** The subagent edits it and writes
   `Also changed:` as today; the main session checks for collisions on
   return, before archiving. The owner: "I think blocking will get too
   slow. Lets go with 2 for now".
5. **Questions.** A subagent records a checkpoint and goes on; the main
   session puts the questions to the owner when a change returns and
   reruns that change's subagent with the answers (the owner chose this
   over returning at the first question).
6. **Both modes.** The same with the owner and with `--unattended`;
   unattended, the main session lists every checkpoint in its report
   (the owner chose it over attended only).
7. **Launchers wait for a prototype** of one gh-aw job per change (the
   owner chose it over a guide now).
8. The owner said "yes" to each of these in the round before writing:
   - Subagents run each member by `grove instructions implement
     <change-id> --unattended`, in both modes; with the owner, the main
     session puts every returned checkpoint to them at once and reruns
     the member with their words, as "Take a checkpoint's answer" allows,
     and a revision their answer needs goes through them and propose in
     the main session.
   - A dependent member starts only once its dependencies are archived;
     Grove already refuses it before.
   - Where two members built together changed the same file, the main
     session runs the project's tests once both returned, and if they
     fail rebuilds the later-started one on its own before archiving it.
   - The instructions say "a fresh context of your own, such as a
     subagent", naming no harness, and fall back to building one after
     another.
   - Grove caps nothing; the harness's own limits apply.
   - Each line of the cycle's list names the members it depends on,
     `after <ids>`; the agent reads `Files:` from each `tasks.md`, and
     Grove prints no `Files:`.
   - In implement's spec, "Build a change, a draft or a cycle" is
     modified and "Build a cycle's changes in their own contexts" added;
     the runs spec is unchanged.
   - One task; test layout and private helper names delegated; at most
     about 250 more words of instruction text (8,426 of rule 12's 9,000
     at bc9c29b).

### What the gh-aw prototype found

A throwaway prototype (branch `prototype/gh-aw-handoff`, 9481bf4, never
merged) asked whether a Claude Code GitHub agentic workflow could install
grove, build a cycle unattended into a pull request, then take the
owner's answer from a PR comment and push to the same branch. Verdict:
yes, for a cycle of one change built by one job. It tested nothing about
context or one job per change, which is why launchers wait for a
prototype of their own.

## Open questions

- Can wait: whether one gh-aw job per change needs each task's `Files:`
  or the cycle a change is in from `grove changes --json` (the next
  prototype).

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
