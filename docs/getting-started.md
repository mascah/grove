# Getting started

This page takes one change from an idea to a spec that describes the new
behavior. You need Git, a repository, a coding agent Grove has skills for
(Claude Code, Codex, Cursor, OpenCode, Pi or Hermes), and Grove
([Install](install.md)).

Commands marked **terminal** go in your shell. Commands marked **chat** go in
your agent's conversation. This page types skills as Claude Code does; the
[Reference](reference.md), under "Harnesses", shows how your agent types
them, such as `$grove-explore` in Codex.

## 1. Set up

**Terminal**, at the root of your repository:

```sh
grove init
```

> Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.

It chooses the agents from what your repository holds, such as `.claude/`
or `AGENTS.md`, and Claude Code and Codex when it finds none. `agents` in
`grove.yaml` lists them; change it and run `grove init` again to add an
agent, or to have Claude Code's as commands, `/grove:explore`.

Grove adds nothing about your project to what your agent already reads, such
as `CLAUDE.md` or `AGENTS.md`; keep project facts there. Then commit:

```sh
git add -A && git commit -m "Set up Grove"
```

A new project has no specs, and that is fine: you do not need to describe
your whole system before you change it.

## 2. Explore

**Terminal:** start a branch for the change.

```sh
git switch -c rename-projects
```

**Chat:** open your agent in the repository and describe what you want.

```text
/grove-explore Users want to rename a project after creating it
```

The agent reads the code and the project's specs and changes, then interviews
you, one question at a time. It keeps apart what is agreed, what is approved
and what the code merely does. Exploring writes nothing; when you want
something kept, ask, or go on to propose.

### Try it first

When a question's answer could change what you propose, such as how a new
screen feels or whether a library can do what you need, the agent names the
question and asks whether to prototype it. Say yes, or ask yourself:

```text
/grove-prototype Should renaming happen inline in the project list?
```

The agent agrees the question and its limits with you, builds the simplest
throwaway experiment on a branch such as `prototype/inline-rename`, runs
it, and reports how to run it, what it observed, a verdict and what it
cannot show. It writes no Grove document. When it asks, let it commit the
experiment there, so the proposal can cite it. That branch is never merged;
push it if you want to keep it.

When you go on to propose, the agent switches back to `rename-projects` and
writes the findings into the proposal, citing the prototype's branch and
commit. A prototype that kills the idea leaves only its branch.

## 3. Propose

**Chat**, in the same conversation:

```text
/grove-propose
```

The agent writes a change folder, such as
`grove/changes/261015-k3m9q-let-users-rename-projects/`, with:

- `.grove.yaml`: the change's title and when it was created.
- `proposal.md`: why, what changes, and your decisions in your words.
- `specs/<capability>/spec.md`: the requirements this adds or changes.
- `tasks.md`: the work, each task with how to check it is done.

It asks rather than assumes, keeps anything you leave undecided as an open
question, and says which task, if any, that question blocks. It runs
`grove check` until the folder has no problems. You can
also propose without exploring: `/grove-propose <description>`. To revise a
drafted change later, run `/grove-propose <change-id>`.

## 4. Review and approve

**Terminal:**

```sh
git add -A && git commit -m "Propose renaming projects"
grove review
```

> Merging rename-projects into main approves change 261015-k3m9q as future work; it does not implement it.

Then each document the merge adds or changes, and the diff. Read the specs
first: they are what you are agreeing to. When you are satisfied, merge:

```sh
git switch main && git merge --no-ff rename-projects
grove changes
```

> 261015-k3m9q  Let users rename projects  approved

The merge is the approval. Approved means agreed future work, not built.

## 5. Implement

**Terminal:** `git switch -c rename-projects-impl`. **Chat:**

```text
/grove-implement 261015-k3m9q
```

The agent works through the tasks with tests, tells you how to try each one,
and records each task's status and evidence in the change's
`implementation.md`. If an open question decides how a task is built, the
agent does not guess: it asks you. If you say to leave it, it marks that
task blocked, writes the question there, and builds the rest.

When your answer or your feedback changes what the change says, the agent
tells you which of its proposal, specs and tasks it would edit and how. Once
you agree, it edits them on this branch and builds to them. Merging the
branch approves the revision with the code, and `grove review` shows it.

`/grove-apply 261015-k3m9q` builds the change without you: it asks no one,
marks a task an open question decides blocked, and never edits the change.
`grove run` uses it (see [Running apply unattended](running.md)).

When every task is done, the agent archives the change: it merges the
change's specs into the maintained specs in `grove/specs/`, then runs
`grove archive`, which moves the change to `grove/changes/archive/`.
`/grove-archive 261015-k3m9q` does the same by hand, and is how a change
you will not build is dropped.

## 6. Merge

**Terminal:** commit, and merge the branch. The code, the updated spec and
the archived change land together.

```sh
grove specs
```

> projects  Projects  grove/specs/projects/spec.md

`grove changes` no longer lists the change. Run `grove` with no command to
read the specs and changes.

## Next

A new `/grove-explore` now starts from the updated spec. See
[Coordinating changes](coordinating.md) when one change waits on another,
[Concepts](concepts.md) for what each document means, and
[Reference](reference.md) for every command and setting.
