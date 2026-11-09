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
title: "Add ledger delete ID"
created: 2026-10-19T14:02:11Z
schema: default
depends_on: [261019-p4x8d]
base: 4f1c2a9e0b7d6c5e4f3a2b1c0d9e8f7a6b5c4d3e
```

`base` is the commit your main branch and origin's both held when the
change was proposed. It stays in origin's history whatever you later
squash or rebase on your machine. Grove's template writes it, so you never
type it.

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

A line also ends with `open questions` when the change's proposal has
any, and with the tasks its `implementation.md` marks `Status: blocked`:

> 261019-r2w7k  Add ledger delete ID  approved  waits on 261019-p4x8d  open questions  T-02, T-04 blocked

The "Open questions" section counts unless, without its HTML comments, it
is empty or says "None.". The line reports what your main branch holds: a
task blocked on a branch not merged yet shows on that branch's card in
`grove serve` (see [Watching the project](watching.md)), not here.

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
A change that depends on others in the list ends `after` their IDs:

```text
1. 261019-p4x8d "Editing an unknown ID prints an error"
2. 261019-r2w7k "Add ledger delete ID" after 261019-p4x8d
```

The agent builds each change in a fresh context of its own, such as a
Claude Code subagent, so no one context carries the whole cycle. Each
follows `grove instructions implement <change-id> --unattended`, which
prints the change's instructions with the cycle's list, naming the
change's place in it, but neither archives nor commits. The agent starts at once every change
whose `after` changes are archived and whose tasks' `Files:` share no file
with a change being built; the others wait their turn. As each returns, the
agent itself archives and commits it, one at a time. If two changes built
together changed the same file anyway, through an `Also changed:` line, it
runs the project's tests; if they fail, it archives the one started first
and builds the other again on its own. Where the harness has no subagents,
the agent builds each change itself, one after another. A change in a
cycle that is not done always gets the cycle's list, with you or without.
So `/grove-implement 261019-p4x8d` shows you its place too.

`--unattended` builds a cycle whether or not it is marked `unattended`,
and whatever its `needs_human` says. Only `grove claim` reads those. Given a
change's ID, `/grove-implement 261019-p4x8d --unattended` builds that one
change the same way, with no cycle to plan first.

With you or unattended, the agent changes the files a change's tasks
name on their `Files:` lines and the paths their `Delegated:` lines name.
It may find another file the task needs changed, or makes false. It edits
that file without asking and writes an `Also changed:` line under the task
in `implementation.md`, with why:

```markdown
Also changed: `README.md` — T-01 makes its usage example false
```

`grove review` then shows that file as
`also changed <change-id>:T-01: T-01 makes its usage example false`, so
you read the reason beside it.

A task may name a file that changed on your main branch since the
change's `base`, or, for a change in a cycle, since the cycle's plan.
Unattended, that task is marked blocked, naming the file, and is not
built. With you, the agent asks what to do.

Without `--unattended`, `/grove-implement 261020-m4k8t` builds the same
cycle with you. Each change is still built in its own context, asking no
one; when it returns, the agent puts its checkpoints to you at once and
builds it again with your answers. An answer that changes what the change
says goes through you and `/grove-propose`, as for one change.

Some changes cannot be built in the session. A change may wait on one
that is neither in the cycle nor archived as finished, because it was
dropped or does not exist. That change is left out and named, and
`grove check` says why. A change with a checkpoint (see
[Getting started](getting-started.md)) stays open, with its `Check:` line
in `implementation.md`. A change that waits on it is refused and stays as
it is. The agent's report lists the checkpoints first, then which changes
were archived, which are blocked and on what, and which were left out. The [implement spec](https://github.com/mascah/grove/blob/main/grove/specs/implement/spec.md) has the exact rules for building a cycle.

The branch holds the code, the updated specs and each archived change.
Review the pull request and merge once. A task its `implementation.md`
already marks done is not built twice.

A later session on the same branch picks up the checkpoints. Run
`/grove-implement 261019-p4x8d` there: the agent puts each `Check:` line
to you, writes your answer under its task as an `Answer:` line, builds on
it and archives the change once every task is done. Or start an
unattended session with your answer beside the command, as a workflow
could pass on a pull request comment: `/grove-implement 261019-p4x8d
--unattended` with "use six-digit IDs". An answer that changes what the
change says still needs `/grove-propose <change-id>`: unattended, the task
stays blocked, naming the revision. Revise the change, merge, and build
again.

Nothing in the session needs you. It creates a branch if it is on your
main branch. It settles no question itself. It commits only when a rule in
`grove.yaml` asks, such as `Commit after each task, naming its task ID.`
under `rules.skills.implement`. That rule counts with you and without (see
[Reference](reference.md), "grove.yaml"). `grove claim` claims a cycle you name
or mark `unattended` for such a session that you start (see
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
the commit your main branch and origin's both hold. `grove check` is clean
again. Merge it, and the revised change is approved. Do not edit `base`
by hand: that hides the report without bringing the change up to date.

The code can move too. Each task names the files it touches on its
`Files:` line. `grove check` reports a named file that changed on your
main branch since the change's `base`:

> grove/changes/261019-s5v3m-show-what-was-saved/tasks.md: T-01's ledger/cli.py changed on main since base 4f1c2a9; revise the change with propose, which also re-bases it after a dependency lands

Only main counts. A branch building the change edits that file, and is
not reported for it. `/grove-implement 261019-s5v3m` names the task and
file. The agent tells you what changed and asks before building that task.

A change that depends on another is usually proposed before that one
lands, so the files the dependency touches are reported once it merges.
That is expected: run `/grove-propose 261019-s5v3m`. The agent brings the
tasks up to date with the files as your main branch has them and sets
`base` to its commit, and `grove check` is clean again. With you,
`/grove-implement` recommends this when the file came with a dependency;
`--unattended` marks the task blocked.

A change in a cycle that is not done is not reported against its `base`.
A moved file makes its cycle stale instead (see
[Planning cycles](planning.md), "Re-plan a cycle"). The change's
instructions name the task and file as changed since the cycle's plan,
and the agent tells you and asks, as for a file changed since `base`:

> - Changed on main since cycle 261020-m4k8t's plan 8f2d6c1: T-01's `ledger/cli.py`.

Re-planning the cycle is then enough to build on the new text. Revising a
change in such a cycle with `/grove-propose <change-id>` re-plans the
cycle on the same branch (see "Re-plan a cycle" there).

In a change spec, each REMOVED requirement keeps its heading,
`### Requirement: <name>`, followed by why it goes. That way `grove check`
can find it.

## The project as data

**Terminal:**

```sh
grove changes --json
```

prints the project as one JSON object, for your agent or a script to read:
your main branch's changes, archived ones included, and around them what
`grove serve` shows (see [Watching the project](watching.md)). The page
reads this same object, so the two agree but for the time of the read.

```json
{
  "root": "grove",
  "target": "main",
  "commit": "8f2d6c1a4b9e0f7d3c5a2b1e6f4d8c0a9b7e5f3d",
  "read": "2026-10-20T09:14:03Z",
  "specs": [{"id": "record-editing", "path": "specs/record-editing/spec.md", "title": "Record editing"}],
  "changes": [
    {"id": "261019-r2w7k", "path": "changes/261019-r2w7k-add-ledger-delete/proposal.md",
     "title": "Add ledger delete ID", "created": "2026-10-19T14:02:11Z",
     "schema": "default", "depends_on": ["261019-p4x8d"],
     "base": "4f1c2a9e0b7d6c5e4f3a2b1c0d9e8f7a6b5c4d3e", "state": "approved",
     "tasks": [{"id": "T-01", "title": "Add ledger delete ID"}],
     "cycle": "261020-m4k8t", "ready": true,
     "files": ["changes/261019-r2w7k-add-ledger-delete/.grove.yaml", "changes/261019-r2w7k-add-ledger-delete/proposal.md",
               "changes/261019-r2w7k-add-ledger-delete/tasks.md"]},
    {"id": "261019-p4x8d", "path": "changes/archive/261019-p4x8d-editing-an-unknown-id/proposal.md",
     "title": "Editing an unknown ID prints an error", "created": "2026-10-19T13:40:02Z",
     "schema": "default", "base": "4f1c2a9e0b7d6c5e4f3a2b1c0d9e8f7a6b5c4d3e",
     "state": "archived", "tasks": [{"id": "T-01", "title": "Report an unknown ID", "status": "done"}],
     "ready": false, "files": ["changes/archive/261019-p4x8d-editing-an-unknown-id/proposal.md"]}
  ],
  "cycles": [...], "branches": [...], "runs": [...], "history": [...], "problems": [...]
}
```

- `read` is when Grove read the project, and `specs` the maintained specs.
- `path` is a change's main document, under `root`.
- `state` is `approved`, `archived`, or `dropped` for a change archived
  without being built.
- `waits_on` lists what `depends_on` names that is not archived as finished;
  it is left out when the change waits on nothing.
- `open_questions` is the text of the "Open questions" section of the
  change's main document, without its HTML comments; it is left out when
  that is empty or says "None".
- A task's `status` is its `Status:` line in `implementation.md`. That file
  usually reaches your main branch with the archive merge. So an approved
  change's tasks mostly have none.
- `cycle` is the cycle not done that lists the change, and `ready` says it
  is approved and waits on nothing; `files` are its documents, under
  `root`.
- A field with no value is left out.
- `cycles`: each cycle, its members in build order, and what `grove
  claim` would do with it, or why not.
- `branches`: each branch, and the working tree, holding work the main
  branch lacks, with each change's tasks and their status there, and its
  `open_questions` as that branch holds it.
- `runs`: the attempts recorded on `grove/runs`, newest first, with
  cycle, host, start time and branch.
- `history`: the main branch's commits that did something to the record,
  newest first, with what each did.
- `problems`: what `grove check` would report on the main branch's
  documents.
- A list with nothing in it is `null`, but for `changes`, which is `[]`.

The [changes spec](https://github.com/mascah/grove/blob/main/grove/specs/changes/spec.md) has the exact rules.

To keep a tracker's issues in step with this data, see the guide
[Push to a tracker](guides/push-to-a-tracker.md).
