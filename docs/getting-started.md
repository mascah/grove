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
or `AGENTS.md`. When it finds none, it chooses Claude Code and Codex.
`agents` in `grove.yaml` lists them. To add an agent, change it and run
`grove init` again. The same way, you can have Claude Code's skills as
commands, such as `/grove:explore`. The
[init spec](https://github.com/mascah/grove/blob/main/grove/specs/init/spec.md) has the exact rules.

Open `grove.yaml` and write, under `context`, what an agent working through
Grove should always do, such as how to show a task is done. It is given to
every skill, and to no other session. Project facts every session needs
stay in `CLAUDE.md` or `AGENTS.md` ([config spec](https://github.com/mascah/grove/blob/main/grove/specs/config/spec.md)). Keep it short:

```yaml
context: Run python3 -m unittest before saying a task is done.
```

Then commit:

```sh
git add -A && git commit -m "Set up Grove"
```

A new project has no specs, and that is fine: you do not need to describe
your whole system before you change it. To write down what some of it
already does, see [Specs of the code you already have](guides/spec-existing-code.md).

## 2. Explore

**Terminal:** start a branch for the change.

```sh
git switch -c rename-projects
```

**Chat:** open your agent in the repository and describe what you want.

```text
/grove-explore Users want to rename a project after creating it
```

The agent reads the code and the project's specs and changes. Then it
interviews you, one question at a time. It keeps apart what is agreed, what
is approved and what the code merely does. Exploring writes nothing. When
you want something kept, ask, or go on to propose.

### Try it first

Some questions are best answered by trying something, such as how a new
screen feels or whether a library can do what you need. When the answer
could change what you propose, the agent names the question and asks
whether to prototype it. Say yes, or ask yourself:

```text
/grove-prototype Should renaming happen inline in the project list?
```

First the agent agrees the question and its limits with you. Then it
builds the simplest throwaway experiment on a branch such as
`prototype/inline-rename`, and runs it. It reports how to run it, what it
observed, a verdict, and what it cannot show. It writes no Grove document.
When it asks, let it commit the experiment there, so the proposal can name
where its findings came from. That branch is never merged; push it if you
want to keep it.

When you go on to propose, the agent switches back to `rename-projects` and
writes the findings into the proposal: what was tried, what was observed and
the verdict, so the proposal stands without the branch. It names the
prototype's branch and commit as their source. A prototype that kills the
idea leaves only its branch. The
[prototype spec](https://github.com/mascah/grove/blob/main/grove/specs/prototype/spec.md) has the exact rules.

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
- `tasks.md`: the work, each task naming the files it touches and saying
  how to check it is done.

Before writing, the agent asks you one round of questions. It lists every
decision it would otherwise assume, each with a recommended answer, and
waits for yours. The questions come in order of consequence:

- first, what changes behavior, scope or compatibility, or could lose data;
- then, choices that are costly to reverse;
- last, routine choices, on one line of defaults that one "yes" accepts.

The idea may be too big for one change, or have parts you could agree
separately. Then the agent offers to write several changes on this branch.
Each later one names the earlier in `depends_on`.

The agent may also recommend letting the implementer make some small
choices alone, such as test layout or a new test file. If you agree, that
task gets a `Delegated:` line naming them.

Anything you leave undecided stays an open question. The agent says which
task, if any, that question blocks.

Before showing you the change, the agent checks it. Every requirement needs
an outcome you can observe, and a task whose acceptance covers it. The
scope, specs and tasks must agree. The agent fixes what it can, tells you
what it changed, and asks you about the rest. It runs `grove check` until
the folder has no problems.

You can also propose without exploring: `/grove-propose <description>`. To
revise a drafted change later, run `/grove-propose <change-id>`.

When it is done, the agent offers two next steps and picks neither:

- merge the change now to approve it, as §4 does;
- or build it on this branch first, so one merge approves it with its code.

For the second, stay on `rename-projects` and run
`/grove-implement 261015-k3m9q` there, as §5 describes. The agent treats the
change as a draft, and edits it only with your agreement. After the
archive, `grove review` says

> Merging rename-projects into main approves change 261015-k3m9q with its code.

and one merge lands everything. The rest of this page takes the two-merge
path. The [propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md) has the exact rules for proposing.

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

The agent works through the tasks with tests, and tells you how to try each
one. It records each task's status and evidence in the change's
`implementation.md`. If an open question decides how a task is built, the
agent does not guess: it asks you. If you say to leave it, it marks that
task blocked, writes the question there, and builds the rest.

Before it finishes, the agent checks its work:

- It runs the tests.
- It rechecks each done task's evidence against the code as it now is. A
  task whose evidence no longer holds goes back to in progress.
- It lists each scenario of the change's specs under `## Scenarios` in
  `implementation.md`, with the test that shows it, or `unverified` and
  why.

Its report, and archiving, name each unverified scenario.

A task's `Delegated:` line names choices and paths the agent may take
without asking you. Anything else beyond the task's `Files:` it asks about.
Once you agree, it edits the file and writes an `Also changed:` line under
the task in `implementation.md`, naming the file and the reason;
`tasks.md` stays as approved. `/grove-implement --unattended` ignores
`Delegated:`, keeps to `Files:` and writes no `Also changed:` line.

The agent may find that the agreed approach is wrong. Then it stops that
task, even if the task could still be built as written. It tells you what
it found and the revision the change needs.

When your answer or your feedback changes what the change says, the agent
tells you which of its proposal, specs and tasks it would edit and how. Once
you agree, it edits them on this branch, through
`grove instructions propose <change-id>`, and builds to them; it never edits
a `Files:` line by hand. The implement instructions give only the template
of `implementation.md`, and leave the others to propose. Merging the
branch approves the revision with the code, and `grove review` shows it.

`/grove-implement 261015-k3m9q --unattended` builds the change without
you. It asks no one. It marks blocked a task that an open question decides,
and it never edits the change, or `grove.yaml`. Given a cycle's ID, it builds each change of
the cycle the same way (see [Planning cycles](planning.md)). `grove run`
uses it too (see [Running unattended](running.md)).

If you corrected the agent in a way that would apply to future work, such
as how to name tests, it suggests one line for `context` or `rules` in
`grove.yaml` when it finishes, quoting you. It adds the line on this branch
only if you agree. `--unattended` may name such a line in its report, for
you to add yourself.

When every task is done, the agent archives the change. It merges the
change's specs into the maintained specs in `grove/specs/`. Then it runs
`grove archive`, which moves the change to `grove/changes/archive/`.
`/grove-archive 261015-k3m9q` does the same by hand. It is also how you
drop a change you will not build. The exact rules are in the
[implement spec](https://github.com/mascah/grove/blob/main/grove/specs/implement/spec.md) and, for archiving, the
[requirements spec](https://github.com/mascah/grove/blob/main/grove/specs/requirements/spec.md).

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
