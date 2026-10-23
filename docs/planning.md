# Planning cycles

A cycle is approved changes you mean to build together, such as a release.
Planning one finds out, with you, which of its changes need a person. A
cycle is also how work reaches `grove run`: it builds only a cycle you have
marked `unattended`, never a change on its own (see
[Running apply unattended](running.md)).

## Plan a cycle

Start from your main branch with the changes approved. **Terminal:**
`git switch -c plan-0.4`. **Chat:**

```text
/grove-plan 0.4 release
```

The agent lists the approved changes, asks which belong in the cycle, and
reads each one you agree to. Where a change calls for something an
unattended session may not be able to do, such as a login the runner
lacks, a manual step, or judging how a screen looks, it says what and asks
whether the change needs a person:

> 261019-r2w7k's T-01 asks that the new list "reads clearly". Does this change need a person to build it?

Only a yes is written down, with your reason. It also tells you of a change
that waits on one outside the cycle that is not archived as finished. Then
it writes one file, `grove/cycles/<cycle-id>-<slug>.yaml`, and edits no
change:

```yaml
title: 0.4 release
created: 2026-10-20T09:12:40Z
goal: Delete and budgets in the README's planned list
changes: [261019-p4x8d, 261019-r2w7k]
needs_human:
  261019-r2w7k: judging whether the new list reads clearly
planned: 8f2d6c1a4b9e0f7d3c5a2b1e6f4d8c0a9b7e5f3d
```

`created` is when the plan started, and `planned` the commit of your main
branch it read. Grove's template writes both, so you never type them.

To send one change on its own, plan a cycle of one from its ID:

```text
/grove-plan 261019-p4x8d
```

The cycle takes the change's title unless you give another.

## Merge the plan

**Terminal:** commit, then `grove review` labels the file `cycle`:

> adds     cycle            grove/cycles/261020-m4k8t-0-4-release.yaml

Merge it, and the cycle is planned. `grove changes` names each change's
cycle:

> 261019-p4x8d  Editing an unknown ID prints an error  approved  in cycle 261020-m4k8t
>
> 261019-r2w7k  Add ledger delete ID  approved  in cycle 261020-m4k8t

A planned cycle is built by you, unless you mark it. **Chat:**
`/grove-implement 261020-m4k8t` builds each of its changes in turn with
you, as `/grove-implement <change-id>` does.

## Mark it unattended

Handing a cycle to `grove run` is a merge of its own. On a branch, add one
line to the cycle file:

```yaml
unattended: true
```

and merge it. `grove changes` then says `in cycle 261020-m4k8t, unattended`.
No skill writes or removes the line, and re-planning keeps it.

`grove check` reports a marked cycle that `grove run` would skip. One whose
`needs_human` names a change:

> grove/cycles/261020-m4k8t-0-4-release.yaml:8: unattended: 261019-r2w7k needs a person, as needs_human says, so grove run skips this cycle; remove unattended, or build the cycle with /grove-implement 261020-m4k8t

And one whose plan is stale, because the folder of a change in it that is
not archived differs on your main branch from what it was at `planned`:

> grove/cycles/261020-m4k8t-0-4-release.yaml: planned: 261019-p4x8d changed on main since this plan, so grove run skips this cycle; re-plan it with /grove-plan 261020-m4k8t

## Re-plan a cycle

```text
/grove-plan 261020-m4k8t
```

revises the cycle as you ask: the agent reads its changes as your main
branch holds them, asks about them as for a new plan, and sets `planned`
to your main branch's commit. Merge it.

A cycle is done when every change in it is archived, finished or dropped;
its file stays. A change is in at most one cycle that is not done.
`grove check` reports a cycle file that has no `title`, `created`,
`changes` or `planned`, has a field Grove does not know, lists a change that is no
change, or lists one that another cycle not done lists:

> grove/cycles/261021-q2w9e-later.yaml:3: changes: 261019-p4x8d is also in cycle 261020-m4k8t, which is not done
