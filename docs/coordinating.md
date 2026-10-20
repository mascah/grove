# Coordinating changes

[Getting started](getting-started.md) takes one change from an idea to an
updated spec. When several changes are approved at once, some must land
before others, and some belong to one piece of work, such as a release.
This page shows how a change says what it waits on, how to plan a cycle,
and what happens when two changes rewrite the same requirement.

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

A change waits on each change in its `depends_on` until your main branch
holds that change archived as finished. A dropped change never finishes.

Applying a change that waits is refused. **Chat:**

```text
/grove-apply 261019-r2w7k
```

> grove: 261019-r2w7k waits on 261019-p4x8d, which main has not archived as finished

The agent tells you and builds nothing. There are two ways past: apply and
archive the change it waits on, and merge that; or, if it need not wait
after all, remove the dependency from `depends_on` on a branch and merge
that. Any field of an approved change's `.grove.yaml` but `created` can
change this way.

`grove check` reports a `depends_on` that can never be met: one naming no
change, open or archived; one naming a dropped change; and two or more
changes that wait on each other in a loop.

> grove/changes/261019-r2w7k-add-ledger-delete/.grove.yaml:4: depends_on: 260101-zzzzz is no change, open or archived

## Plan a cycle

A cycle is approved changes you mean to finish together. It is a file,
`grove/cycles/<id>-<slug>.yaml`, holding its `title` and, if you want
them, a one-line `goal` and `start` and `end` dates as `YYYY-MM-DD`:

```yaml
title: "0.4"
goal: Everything the README plans for 0.4.
start: 2026-10-20
end: 2026-11-06
```

A change joins a cycle by naming it in its `.grove.yaml` as
`cycle: <cycle-id>`, so it is in at most one. The cycle file does not list
its changes.

**Terminal:** `git switch -c plan-04`. **Chat:**

```text
/grove-plan 0.4: the README's planned list
```

The agent reads the approved changes, says what each waits on and which
can be built in parallel, and asks whether changes that touch the same
capability should wait on one another. It asks for the cycle's title and
goal, then writes the cycle file and `cycle:` in each change's
`.grove.yaml`, telling you first that this edits approved changes. It runs
`grove check` until there are no problems. To revise a cycle later, run
`/grove-plan <cycle-id>`.

Commit, read `grove review`, which labels the file `cycle`, and merge. The
merge accepts the plan, as it approves a change. **Terminal:**

```sh
grove changes
```

> 261019-j83h9  0.4  2026-10-20 to 2026-11-06  0 of 2 archived
>
> &nbsp;&nbsp;261019-c4t8n  Monthly budgets in report  approved
>
> &nbsp;&nbsp;261019-r2w7k  Add ledger delete ID  approved

`grove changes` lists each cycle that is not done, with its dates if it has
them and how many of its changes are archived, and under it its approved
changes. Changes in no cycle follow.

Once a project has an open cycle, `/grove-propose` asks which cycle, if
any, a new change joins. `grove check` reports a `cycle` that names no file
in `grove/cycles/`.

The order comes from `depends_on`: a change can be applied once every
change it waits on is archived as finished, and changes that do not wait on
each other can be built in parallel, on separate branches. A cycle is done
when every change in it is archived. It stays in `grove/cycles/`.

## When the plan is stale

A change's spec MODIFIES or REMOVES a maintained requirement as it read at
the change's `base`. If another change rewrites or removes that requirement
on your main branch first, the change was planned against text that is no
longer there. `grove check` reports each such requirement:

> grove/changes/261019-s5v3m-show-what-was-saved/specs/record-editing/spec.md:5: requirement "Keep an unchanged category on edit" was rewritten on main since base 4f1c2a9; revise the change with propose

It compares `base` with your main branch, not with your working tree, so
the branch that archives the first change does not report the second.
`/grove-archive` runs `grove check` first, and stops if it reports the
change being archived.

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
