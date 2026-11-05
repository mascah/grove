# Implement stops at checkpoints and hands off in the branch

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

With the owner, implement stops after every task: "When a task is built,
tell the owner how to try it and take their feedback before the next"
(`internal/instruct/requires/implement-owner.md` at 8a8d064). A cycle
builds each change in one session (`implement-cycle.md`), so a large
cycle is many waits in one long conversation. Unattended, implement
never stops for a look. It blocks only on an open question or a task it
cannot build, and such a task "is settled by revising the change, not
here" (`implement-unattended.md`), which takes a `/grove-propose`
revision and a merge.

The owner, in a session on 2026-11-04, said they have barely used
`grove run` since it was added because "most of the sessions I felt like
I needed to inspect the running app between tasks". Building a cycle
with the owner "can take a long while and sometimes its up to the user
to digest things to be able to evaluate between tasks without letting
the claude cache go stale so there feels like a weird time crunch". They
cannot always say ahead of time where to look: "the agent has to tell me
why and I try and make a decision".

Neither mode lets the agent say "this needs you, because ...", build
everything else, and leave a handoff on the branch that a later
session, local or remote, interactive or not, can pick up.

## Intended outcome

In both modes, the agent finds the checkpoints as it builds. A
checkpoint is a task whose acceptance no check it can run shows, or one
resting on a choice the change does not settle. For each checkpoint it
either builds its best guess, marking the task `Status: to check`, or,
where a wrong guess would be costly, builds nothing and marks it
`Status: blocked`. A `Check:` line under the task says what to try or
decide and why.

The agent builds every task that waits on no checkpoint. Then, with the
owner, it stops once and puts the checkpoints to them. Unattended, it
finishes. Its report lists the checkpoints first. A task `to check` is
not done, so archive refuses the change until the owner answers.

A later session on the branch picks up the checkpoints from
`implementation.md`: a short one with the owner, or an unattended one
started with the owner's words, such as a pull request comment an
agentic workflow passes on. The agent writes each answer under its task
as an `Answer:` line and builds on it. An answer that changes what the
change says still goes through `/grove-propose`. The rules are in
[specs/implement/spec.md](specs/implement/spec.md).

## Scope

### In scope

- The implement instructions, with the owner, unattended and for a
  cycle, and their wording test.
- The `implementation.md` guide and template.
- The implement spec's "Settle what the change does not", and the new
  "Stop at a checkpoint" and "Take a checkpoint's answer".
- `docs/concepts.md`, `docs/getting-started.md`, `docs/coordinating.md`
  and a walkthrough step.

### Out of scope

- `grove run`: an attempt, and `--retry`, still start from the target
  branch, and the project stays paused until the cycle is archived.
  Whether `run` stays is for later ("Checkpoints first, run later").
- `needs_human` and the plan instructions.
- Go code other than the test, and `CONSTITUTION.md`.
- The owner writing `implementation.md` by hand: rule 3 has the
  implementing agent write it.
- A guide for GitHub agentic workflows; it follows a prototype.
- Managing a session's context across a cycle, and building in
  parallel; they are their own question.
- The skill trials.

## Decisions

The owner settled the direction while exploring on 2026-11-04: the agent
decides where to stop and says why ("the agent has to tell me why and I
try and make a decision"); while a checkpoint waits, it should "build the
rest, then end"; it builds a guess or nothing, "Agent says which, per
case"; and the loop should be "a fluid human in the loop process
regardless of whether the sessions are running locally or remotely",
with "Checkpoints first, run later". The owner accepted every answer
below in one round ("Yes to all", 2026-11-05).

1. **With the owner, the stop after every task goes.** The agent builds
   every task that waits on no checkpoint, then stops once, listing each
   checkpoint with what to try and why. It still stops at once when it
   finds the change's approach is wrong. Reason: fewer, explained stops,
   and none while the owner digests.
2. **What a checkpoint is, and guess or nothing.** A task whose
   acceptance no check the agent can run shows, such as how it looks or
   a manual step, or that rests on a choice the change does not settle
   that the owner would notice. The agent builds its guess and marks it
   `Status: to check` with what it chose and why, or builds nothing and
   marks it `Status: blocked` with the question when a wrong guess is
   costly: data, an interface other tasks build on, something hard to
   undo. Reason: the owner chose "Agent says which, per case".
3. **`Status: to check` is new and is not done**, so archive refuses the
   change until the owner confirms, with no code change. A `Check:` line
   under the task says what to try and why. A task that builds on a
   checkpoint waits and names it. Reason: a guess must not be archived
   unseen.
4. **Answers.** The agent writes each answer under its task as
   `Answer:` with the owner's words quoted, then builds on it.
   Unattended, it follows an answer already under a task or given in the
   words that started the session, as an agentic workflow could pass a
   pull request comment, only when the answer stays within the change;
   one that changes the change leaves the task blocked, to revise with
   `/grove-propose` as today. The owner writing `implementation.md`
   themselves stays out (rule 3). Reason: the handoff lives in the
   branch, so any session anywhere can pick it up.
5. **The report lists the checkpoints first**, so a launcher can post
   them. The project's rules still decide commits, pushes and pull
   requests. Reason: launchers differ; Grove says what, not where.
6. **Spec:** implement's "Settle what the change does not" is MODIFIED;
   "Stop at a checkpoint" and "Take a checkpoint's answer" are ADDED.
7. **Routine defaults:** the `default` schema; one task; `depends_on:
   [261030-t2pt7, 261101-mg8hd, 261105-gmpa4]`, since all three edit
   these files and landing first would make cycle 261104-z05ey's plan
   stale; a `Delegated:` line for test layout and private helper names;
   within rule 12's 9,000-word budget.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
