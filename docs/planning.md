# Planning cycles

A cycle is approved changes you mean to build together, such as a release.
Planning one finds out, with you, which of its changes need a person. A
cycle is also how work reaches `grove run`: it builds a cycle you name or
have marked `unattended`, never a change on its own (see
[Running unattended](running.md)). The [cycles spec](https://github.com/mascah/grove/blob/main/grove/specs/cycles/spec.md) has the exact
rules for cycles.

## Plan a cycle

Start from your main branch with the changes approved. **Terminal:**
`git switch -c plan-0.4`. **Chat:**

```text
/grove-plan 0.4 release
```

The agent lists the approved changes, asks which belong in the cycle, and
reads each one you agree to. Some changes call for something an
unattended session may not be able to do: a login the runner lacks, a
manual step, or judging how a screen looks. For each, the agent says what
it found and asks whether the change needs a person:

> 261019-r2w7k's T-01 asks that the new list "reads clearly". Does this change need a person to build it?

Only a yes is written down, with your reason. The agent also tells you
when a change waits on one outside the cycle that is not archived as
finished. Then it writes one file, `grove/cycles/<cycle-id>-<slug>.yaml`.
It edits no change:

```yaml
title: "0.4 release"
created: 2026-10-20T09:12:40Z
goal: Delete and budgets in the README's planned list
changes: [261019-p4x8d, 261019-r2w7k]
needs_human:
  261019-r2w7k: judging whether the new list reads clearly
planned: 8f2d6c1a4b9e0f7d3c5a2b1e6f4d8c0a9b7e5f3d
```

`created` is when the plan started, and `planned` the commit of your main
branch it read. Grove's template writes both, so you never type them.

The agent ends with a verdict. It names `/grove-implement <cycle-id>`,
which builds the cycle with you. It says whether `grove run <cycle-id>`
would build it. It says that `unattended: true`, merged on its own, hands
the cycle to a `grove run` given no cycle. And for each obstacle, it says
what removes it. For this one:

> Build this cycle with `/grove-implement 261020-m4k8t`. grove run would refuse it, given as `grove run 261020-m4k8t` or marked: 261019-r2w7k needs a person. To remove that, revise it with `/grove-propose 261019-r2w7k` so its acceptance can be checked without a person; that re-plans the cycle on the same branch, and one merge removes it.

For a cycle with nothing in the way:

> Build this cycle with `/grove-implement 261020-m4k8t`, or have `grove run 261020-m4k8t` build it. Adding `unattended: true` to its file in a merge of its own hands it to a `grove run` given no cycle, such as a timer's.

Say a change in the cycle waits on a change outside it. Then `grove run`
does not build the cycle until that outside change is archived as
finished. The verdict also names what the runner's checkout must have,
such as a login you said the runner has.

To send one change on its own, plan a cycle of one from its ID:

```text
/grove-plan 261019-p4x8d
```

The cycle takes the change's title unless you give another.

## Follow another cycle

When other cycles are not done, the agent asks whether the new one
follows any of them, and writes only those you name:

```yaml
depends_on: [261020-m4k8t]
```

`grove run` builds nothing of the cycle while a cycle it follows is not
done. `grove changes` says so: `in cycle 261021-q2w9e, waits on cycle
261020-m4k8t`. A cycle it follows, directly or through others, lands in
merges that archive its changes. Those merges may edit files that this
cycle's tasks name, but they leave this plan fresh. The agent's verdict
says both. `/grove-implement <cycle-id>` ignores `depends_on`: work you
build with the agent, you order yourself.

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

And one whose plan is stale (see "Re-plan a cycle").

## Re-plan a cycle

A plan goes stale when a commit on your main branch, after `planned`,
changes one of two things: the folder of a change in the cycle that is
not archived, or a file one of its tasks names on its `Files:` line. A
commit that also edits the cycle file does not count. Nor does one that
archives a change of a cycle this one follows. Grove reads main's
first-parent history, so a squash or a merge commit counts as one commit.
`grove check` reports a stale cycle, marked or not:

> grove/cycles/261020-m4k8t-0-4-release.yaml: planned: 261019-p4x8d changed on main since this plan; re-plan it with /grove-plan 261020-m4k8t

For a marked cycle, the line says `since this plan, so grove run skips
this cycle`.

```text
/grove-plan 261020-m4k8t
```

revises the cycle as you ask. The agent reads its changes as your main
branch holds them and asks about them as for a new plan. It sets
`planned` to your main branch's commit. Merge it. Say a file a change's
task names has moved, and you are content to build the change on the new
text. Then re-planning is all it takes; the change itself stays as it
is.

Revising a change in the cycle re-plans it too.
`/grove-propose 261019-r2w7k` prints the cycle's file beside the change.
On the same branch, the agent asks whether the change needs a person, as
planning does. It sets `planned` to your main branch's commit. It never
writes or removes `unattended`. Merge the branch as one commit, a squash
or a merge commit, and the plan holds. A rebase or fast-forward brings the
revision and the re-plan to main as two commits. Then the revision is a
change the plan did not check, and the cycle is stale.

A cycle is done when every change in it is archived, finished or dropped.
Its file stays. A change is in at most one cycle that is not done.
`grove check` reports a cycle file that is malformed:

- it has no `title`, `created`, `changes` or `planned`;
- it has a field Grove does not know;
- it lists a change that is no change; or
- it lists a change that another cycle not done lists:

> grove/cycles/261021-q2w9e-later.yaml:3: changes: 261019-p4x8d is also in cycle 261020-m4k8t, which is not done

For a cycle not done, it also reports a `depends_on` that names no cycle,
and cycles that follow each other in a loop:

> grove/cycles/261021-q2w9e-later.yaml:5: depends_on: 261020-m4k8t waits on this cycle in turn, so neither can run
