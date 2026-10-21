# Coordinating changes

[Getting started](getting-started.md) takes one change from an idea to an
updated spec. When several changes are approved at once, some must land
before others. This page shows how a change says what it waits on, how one
session applies every approved change in order, what happens when two
changes rewrite the same requirement, and how to get the whole picture as
data.

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

A change waits on each change in its `depends_on` until it is archived as
finished. `grove changes` reads your main branch; applying reads your
working tree, which on a branch holds everything main held when the branch
started, plus what you archived on it. A dropped change never finishes.

Applying a change that waits is refused. **Chat:**

```text
/grove-apply 261019-r2w7k
```

> grove: 261019-r2w7k waits on 261019-p4x8d, which this working tree has not archived as finished

The agent tells you and builds nothing. There are three ways past: apply the
change it waits on, which archives it when every task is done; apply every
approved change at once, as the next section shows; or, if it need not wait
after all, remove the dependency from `depends_on` on a branch and merge
that. Any field of an approved change's `.grove.yaml` but `created` and
`base` can change this way; `base` changes only when the change is revised
(see "When the plan is stale").

`grove check` reports a `depends_on` that can never be met: one naming no
change, open or archived; one naming a dropped change; and two or more
changes that wait on each other in a loop.

> grove/changes/261019-r2w7k-add-ledger-delete/.grove.yaml:4: depends_on: 260101-zzzzz is no change, open or archived

The order comes from `depends_on`: a change can be applied once every
change it waits on is archived as finished, and changes that do not wait on
each other can be built in parallel, on separate branches.

## Applying every approved change

One session can build everything that is approved and not yet archived, on
one branch, and you merge once. **Terminal:** `git switch -c apply-all`.
**Chat:**

```text
/grove-apply
```

With no ID, Grove lists the approved changes your working tree has not
archived, in build order: each after the changes it depends on, otherwise
oldest first. The agent takes them in turn, each with its own apply
instructions, and when a change's tasks are all done it archives the change
on the branch, as `/grove-apply <change-id>` does. The next change in the
list then finds its dependency archived as finished.

A change that waits on a change that is neither in the list nor archived as
finished, because it was dropped or does not exist, is left out and named;
`grove check` says why. A change whose task is blocked by an open question
stays open, with the question in its `implementation.md`; a change that
waits on it is refused and stays as it is. The agent's report says which
changes were archived, which are blocked and on what, and which were left
out.

The branch holds the code, the updated specs and each archived change.
Review the pull request and merge once. A blocked change is finished by
revising it with `/grove-propose <change-id>`, merging, and applying again: a
task its `implementation.md` already marks done is not built twice.

Nothing in the session needs you. It creates a branch if it is on your main
branch, settles no question itself, and commits only when a rule in
`grove.yaml` asks, such as `Commit after each task, naming its task ID.`
under `rules.skills.apply` (see [Reference](reference.md), "grove.yaml"). A
pipeline that runs your agent on each push to main, with `/grove-apply` as
the prompt, builds whatever that push approved and opens one pull request.

## When the plan is stale

A change's spec MODIFIES or REMOVES a maintained requirement as it read at
the change's `base`. If another change rewrites or removes that requirement
first, the change was planned against text that is no longer there. `grove
check` reports each such requirement:

> grove/changes/261019-s5v3m-show-what-was-saved/specs/record-editing/spec.md:5: requirement "Keep a category" was rewritten since base 4f1c2a9; revise the change with propose

The rule: a requirement is stale when the maintained spec in your working
tree gives it a text that is neither what it was at `base`, nor what this
change gives it, nor what a change it depends on gives it. So a change whose
MODIFIED text has been reconciled into the spec is not reported, and a
change that names the rewriting change in `depends_on` is not reported
either: its author planned against that text. Only a rewrite the change did
not plan for is.

Because `grove check` reads your working tree, the branch that archives the
first change is where the second is reported, before anything is merged,
whether one session applies them both or not. `/grove-archive` runs `grove
check` first, and stops if it reports the change being archived.

To bring the change up to date, revise it on a branch. **Chat:**

```text
/grove-propose 261019-s5v3m
```

The agent rewrites that MODIFIED or REMOVED requirement from your main
branch's text, keeping what the change itself alters, and sets `base` to
your main branch's commit. `grove check` is clean again. Merge it, and the
revised change is approved. Do not edit `base` by hand: that hides the
report without bringing the change up to date.

In a change spec, each REMOVED requirement keeps its heading,
`### Requirement: <name>`, followed by why it goes, so `grove check` can
find it.

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
  usually reaches your main branch with the archive merge, so an approved
  change's tasks mostly have none.
- A field with no value is left out.

To keep a tracker's issues in step with this data, see the guide
[Push to a tracker](guides/push-to-a-tracker.md).
