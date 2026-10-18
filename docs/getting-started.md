# Getting started

This page takes one change from an idea to a spec that describes the new
behavior. You need Git, a repository, Claude Code or Codex, and the `grove` command on
your PATH.

Commands marked **terminal** go in your shell. Commands marked **chat** go in
your agent's conversation. In Codex, skills start with `$` instead of `/`.

## 1. Set up

**Terminal**, at the root of your repository:

```sh
grove init
```

> Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.

Open `grove.yaml` and describe your project under `context`: the language,
how tests run, and anything an agent should always know. Keep it short; it is
given to every skill. Then commit:

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

## 5. Apply

**Terminal:** `git switch -c rename-projects-impl`. **Chat:**

```text
/grove-apply 261015-k3m9q
```

The agent works through the tasks with tests and records each task's status
and evidence in the change's `implementation.md`. If an open question decides
how a task is built, the agent does not guess: it marks that task blocked,
writes the question there, and builds the rest. To unblock it, answer the
question with `/grove-propose 261015-k3m9q`, merge the revision, and apply
again.

## 6. Archive

**Chat**, when every task is done:

```text
/grove-archive 261015-k3m9q
```

The agent merges the change's specs into the maintained specs in
`grove/specs/`, then runs `grove archive`, which moves the change to
`grove/changes/archive/`. **Terminal:** commit, and merge the branch. The
code, the updated spec and the archived change land together.

```sh
grove specs
```

> projects  Projects  grove/specs/projects/spec.md

`grove changes` no longer lists the change. Run `grove` with no command to
read the specs and changes.

## Next

A new `/grove-explore` now starts from the updated spec. See
[Concepts](concepts.md) for what each document means, and
[Reference](reference.md) for every command and setting.
