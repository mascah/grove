# Coordinating changes

[Getting started](getting-started.md) takes one change from an idea to an
updated spec. When several changes are approved at once, some must land
before others. This page shows:

- how a change says what it waits on;
- how one session builds a cycle of changes in order without you;
- what happens when two changes rewrite the same requirement; and
- how to get the whole picture as data.

## A change that waits

When you propose a change, the agent asks whether it needs another open
change to land first, approved or still a draft, and names any that touch
the same capabilities. Each one you name goes in the change's `.grove.yaml`
as `depends_on`:

```yaml
title: Add ledger delete ID
created: 2026-10-19T14:02:11Z
schema: default
depends_on: [261019-p4x8d]
base: 4f1c2a9e0b7d6c5e4f3a2b1c0d9e8f7a6b5c4d3e
```

`base` is the commit of your main branch that the change was proposed from.
Grove's template writes it, so you never type it.

**Terminal**, once both changes are merged:

```sh
grove changes
```

> 261019-p4x8d  Editing an unknown ID prints an error  approved
>
> 261019-r2w7k  Add ledger delete ID  approved  waits on 261019-p4x8d

A change waits on each change in its `depends_on` until that one is
archived as finished. A dropped change never finishes. `grove changes`
reads your main branch. Implementing reads your working tree. On a
branch, that holds everything main held when the branch started, plus
what you archived on it.

Building a change that waits is refused. **Chat:**

```text
/grove-implement 261019-r2w7k
```

> grove: 261019-r2w7k waits on 261019-p4x8d, which this working tree has not archived as finished

The agent tells you and builds nothing. There are three ways past:

- Build the change it waits on. That archives it when every task is done.
- Build both in one cycle, as the next section shows.
- If it need not wait after all, remove the dependency from `depends_on`
  on a branch, and merge that.

Any field of an approved change's `.grove.yaml` can change this way, but
`created` and `base`. `base` changes only when the change is revised (see
"When the plan is stale").

`grove check` reports a `depends_on` that can never be met:

- one naming no change, open or archived;
- one naming a dropped change; and
- two or more changes that wait on each other in a loop.

> grove/changes/261019-r2w7k-add-ledger-delete/.grove.yaml:4: depends_on: 260101-zzzzz is no change, open or archived

The [propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md), "Report a dependency that can never be met", has
the exact rule.

The order comes from `depends_on`. A change can be built once every change
it waits on is archived as finished. Changes that do not wait on each
other can be built in parallel, on separate branches.

## Building a cycle unattended

One session can build every change of a cycle (see
[Planning cycles](planning.md)), on one branch, without asking you, and
you merge once. **Terminal:** `git switch -c build-0.4`. **Chat:**

```text
/grove-implement 261020-m4k8t --unattended
```

Grove lists the cycle's changes your working tree has not archived, in
build order: each after the changes it depends on, otherwise oldest first.
The agent takes them one at a time. For each, it runs
`grove instructions implement <change-id> --unattended`. That prints the
change's instructions with the cycle's list, naming the change's place in
it. When the change's tasks are all done, the agent archives it on the
branch. The next change in the list then finds its dependency archived as
finished. A change in a cycle that is not done always gets the cycle's
list, with you or without. So `/grove-implement 261019-p4x8d` shows you
its place too.

`--unattended` builds a cycle whether or not it is marked `unattended`,
and whatever its `needs_human` says. Only `grove run` reads those. Given a
change's ID, `/grove-implement 261019-p4x8d --unattended` builds that one
change the same way, with no cycle to plan first.

Unattended, the agent changes only the files a change's tasks name on
their `Files:` lines. It may find another file that should change too.
It leaves that file as it is and names it in the task's evidence in
`implementation.md`. With you, the agent asks to add it to the task
instead. A change in no cycle may have a task naming a file that changed
on your main branch since the change's `base`. Unattended, that task is
marked blocked, naming the file, and is not built. With you, the agent
asks what to do.

Without `--unattended`, `/grove-implement 261020-m4k8t` builds the same
cycle with you, each change as `/grove-implement <change-id>` does.

Some changes cannot be built in the session. A change may wait on one
that is neither in the cycle nor archived as finished, because it was
dropped or does not exist. That change is left out and named, and
`grove check` says why. A change whose task is blocked by an open question
stays open, with the question in its `implementation.md`. A change that
waits on it is refused and stays as it is. The agent's report says which
changes were archived, which are blocked and on what, and which were left
out. The [implement spec](https://github.com/mascah/grove/blob/main/grove/specs/implement/spec.md) has the exact rules for building a cycle.

The branch holds the code, the updated specs and each archived change.
Review the pull request and merge once. To finish a blocked change, revise
it with `/grove-propose <change-id>`, merge, and build again. A task its
`implementation.md` already marks done is not built twice.

Nothing in the session needs you. It creates a branch if it is on your
main branch. It settles no question itself. It commits only when a rule in
`grove.yaml` asks, such as `Commit after each task, naming its task ID.`
under `rules.skills.implement`. That rule counts with you and without (see
[Reference](reference.md), "grove.yaml"). `grove run` starts such a
session for a cycle you name or mark `unattended` (see
[Running unattended](running.md)).

## When the plan is stale

A change's spec MODIFIES or REMOVES a maintained requirement as it read at
the change's `base`. Another change may rewrite or remove that requirement
first. Then the change was planned against text that is no longer there.
`grove check` reports each such requirement:

> grove/changes/261019-s5v3m-show-what-was-saved/specs/record-editing/spec.md:5: requirement "Keep a category" was rewritten since base 4f1c2a9; revise the change with propose

Grove compares the requirement's text in your working tree's maintained
spec with three texts: its text at `base`, the text this change gives it,
and the text a change it depends on gives it. If it matches none of them,
the requirement is stale. So a change whose MODIFIED text has been
reconciled into the spec is not reported. Nor is a change that names the
rewriting change in `depends_on`: its author planned against that text.
Only a rewrite the change did not plan for is reported. The [propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md), "Report a stale requirement",
has the exact rule.

`grove check` reads your working tree. So the second change is reported on
the branch that archives the first, before anything is merged. That holds
whether one session builds them both or not. `/grove-archive` runs `grove
check` first, and stops if it reports the change being archived.

To bring the change up to date, revise it on a branch. **Chat:**

```text
/grove-propose 261019-s5v3m
```

The agent rewrites that MODIFIED or REMOVED requirement from your main
branch's text, keeping what the change itself alters. It sets `base` to
your main branch's commit. `grove check` is clean again. Merge it, and the
revised change is approved. Do not edit `base` by hand: that hides the
report without bringing the change up to date.

The code can move too. Each task names the files it touches on its
`Files:` line. `grove check` reports a named file that changed on your
main branch since the change's `base`:

> grove/changes/261019-s5v3m-show-what-was-saved/tasks.md: T-01's ledger/cli.py changed on main since base 4f1c2a9; revise the change with propose

Only main counts. A branch building the change edits that file, and is
not reported for it. `/grove-implement 261019-s5v3m` names the task and
file. The agent tells you what changed and asks before building that task.

A change in a cycle that is not done is not reported against its `base`.
A moved file makes its cycle stale instead (see
[Planning cycles](planning.md), "Mark it unattended"). Re-planning the
cycle is then enough to build on the new text. Revising a change in such
a cycle with `/grove-propose <change-id>` re-plans the cycle on the same
branch (see "Re-plan a cycle" there).

In a change spec, each REMOVED requirement keeps its heading,
`### Requirement: <name>`, followed by why it goes. That way `grove check`
can find it.

## The project as data

**Terminal:**

```sh
grove changes --json
```

prints your main branch's changes as JSON, archived changes included, for your agent or a script to read:

```json
{
  "root": "grove",
  "target": "main",
  "commit": "8f2d6c1a4b9e0f7d3c5a2b1e6f4d8c0a9b7e5f3d",
  "changes": [
    {"id": "261019-r2w7k", "path": "changes/261019-r2w7k-add-ledger-delete/proposal.md",
     "title": "Add ledger delete ID", "created": "2026-10-19T14:02:11Z",
     "schema": "default", "depends_on": ["261019-p4x8d"],
     "base": "4f1c2a9e0b7d6c5e4f3a2b1c0d9e8f7a6b5c4d3e", "state": "approved",
     "tasks": [{"id": "T-01", "title": "Add ledger delete ID"}]},
    {"id": "261019-p4x8d", "path": "changes/archive/261019-p4x8d-editing-an-unknown-id/proposal.md",
     "title": "Editing an unknown ID prints an error", "created": "2026-10-19T13:40:02Z",
     "schema": "default", "base": "4f1c2a9e0b7d6c5e4f3a2b1c0d9e8f7a6b5c4d3e",
     "state": "archived", "tasks": [{"id": "T-01", "title": "Report an unknown ID", "status": "done"}]}
  ]
}
```

- `path` is a change's main document, under `root`.
- `state` is `approved`, `archived`, or `dropped` for a change archived
  without being built.
- `waits_on` lists what `depends_on` names that is not archived as finished;
  it is left out when the change waits on nothing.
- A task's `status` is its `Status:` line in `implementation.md`. That file
  usually reaches your main branch with the archive merge. So an approved
  change's tasks mostly have none.
- A field with no value is left out.

The [changes spec](https://github.com/mascah/grove/blob/main/grove/specs/changes/spec.md) has the exact rules.

To keep a tracker's issues in step with this data, see the guide
[Push to a tracker](guides/push-to-a-tracker.md).
