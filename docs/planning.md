# Planning cycles

A cycle is approved changes you mean to build together, such as a release.
A cycle is also how work reaches `grove claim`: it claims a cycle you name or
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
reads each one you agree to. It titles the cycle itself. It tells you
when a change waits on one outside the cycle that is not archived as
finished. Whether a change needed a person is settled when you review
its build, not in the plan.

Before it writes the cycle, the agent runs `grove check`. A change you
agreed to may be stale on its base: a requirement it modifies was
rewritten on main, or a file its tasks name changed there, often by the
merge that archived a change it waits on. The agent tells you this edits
approved work, and brings the change up to date on the same branch, as
`/grove-propose <change-id>` would. It asks only where that needs a
decision. So one merge brings the revision and the plan, with no second
session.

Then it writes `grove/cycles/<cycle-id>-<slug>.md`. It edits no other
change:

```markdown
---
title: "0.4 release"
created: 2026-10-20T09:12:40Z
goal: Delete and budgets in the README's planned list
changes: [261019-p4x8d, 261019-r2w7k]
planned: 8f2d6c1a4b9e0f7d3c5a2b1e6f4d8c0a9b7e5f3d
---

## Found on main

- 261019-r2w7k was brought up to date: its T-01 now names
  `ledger/report.py` as main has it after 261019-p4x8d's dependency.

## Left out

- 261019-t5c2v waits on a login the runner does not have.
```

The fields sit at the top, as YAML between two lines of `---`. `created` is
when the plan started, and `planned` the commit your main branch and
origin's both held when it read them. Grove's template writes both, so you
never type them.

Below the fields are the plan's notes, for whoever builds the cycle: what
the plan found on main at `planned`, what the builder should know above any
one change, and why these changes, in this order, and what was left out.
The agent writes them only where there is something to say, so a cycle may
have none. Only plan writes them, and re-planning brings them up to date.
Grove reads none of them. They change no change: where the notes and a
change's documents disagree, the change's documents hold, and
`/grove-implement` says so when it prints the cycle file.

The agent ends with a verdict. It names `/grove-implement <cycle-id>`,
which builds the cycle with you. It says whether `grove claim <cycle-id>`
would claim it. It says that `unattended: true`, merged on its own, hands
the cycle to a `grove claim` given no cycle. And for each obstacle, it says
what removes it. For a cycle with nothing in the way:

> Build this cycle with `/grove-implement 261020-m4k8t`, or have `grove claim 261020-m4k8t` claim it for an agent you start. Adding `unattended: true` to its file in a merge of its own hands it to a `grove claim` given no cycle, such as a timer's.

Say a change in the cycle waits on a change outside it. Then `grove claim`
does not claim the cycle until that outside change is archived as
finished. The verdict also names what the runner's checkout must have,
such as a login one of its changes calls for.

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

`grove claim` claims nothing of the cycle while a cycle it follows is not
done. `grove changes` says so: `in cycle 261021-q2w9e, waits on cycle
261020-m4k8t`. A cycle it follows, directly or through others, lands in
merges that archive its changes. Those merges may edit files that this
cycle's tasks name, but they leave this plan fresh. The agent's verdict
says both. `/grove-implement <cycle-id>` ignores `depends_on`: work you
build with the agent, you order yourself.

## Merge the plan

**Terminal:** commit, then `grove review` labels the file `cycle`:

> adds     cycle            grove/cycles/261020-m4k8t-0-4-release.md

Merge it, and the cycle is planned. `grove changes` names each change's
cycle:

> 261019-p4x8d  Editing an unknown ID prints an error  approved  in cycle 261020-m4k8t
>
> 261019-r2w7k  Add ledger delete ID  approved  in cycle 261020-m4k8t

A planned cycle is built by you, unless you mark it. **Chat:**
`/grove-implement 261020-m4k8t` builds each of its changes in turn with
you, as `/grove-implement <change-id>` does.

## Mark it unattended

Handing a cycle to `grove claim` is a merge of its own. On a branch, add one
line to the cycle file's fields, above the second `---`:

```yaml
unattended: true
```

and merge it. `grove changes` then says `in cycle 261020-m4k8t, unattended`.
No skill writes or removes the line, and re-planning keeps it.

`grove check` reports a marked cycle that `grove claim` would skip because
its plan is stale (see "Re-plan a cycle").

## Re-plan a cycle

A plan goes stale when a commit on your main branch, after `planned`,
changes one of two things: the folder of a change in the cycle that is
not archived, or a file one of its tasks names on its `Files:` line. A
commit that also edits the cycle file does not count. Nor does one that
archives a change of this cycle, or of a cycle this one follows. So the
cycle's changes can reach main one pull request at a time, each archiving
its own change, and the plan stays fresh. Grove reads main's
first-parent history, so a squash or a merge commit counts as one commit.
`grove check` reports a stale cycle, marked or not:

> grove/cycles/261020-m4k8t-0-4-release.md: planned: 261019-p4x8d changed on main since this plan; re-plan it with /grove-plan 261020-m4k8t

For a marked cycle, the line says `since this plan, so grove claim skips
this cycle`.

```text
/grove-plan 261020-m4k8t
```

revises the cycle as you ask. The agent reads its changes as your main
branch holds them, as for a new plan. It sets
`planned` to the commit your main branch and origin's both hold, and
brings the notes up to date with it: what holds there, and nothing main
made false. Merge
it. Say a file a change's task names has moved, and you are content to
build the change on the new text. Then re-planning is all it takes; the
change itself stays as it is.

Revising a change in the cycle re-plans it too.
`/grove-propose 261019-r2w7k` prints the cycle's file beside the change.
On the same branch, the agent sets `planned` to the commit your main branch and
origin's both hold. It never writes or removes `unattended`. Merge the
branch as one commit, a squash or a merge commit, and the plan holds. A
rebase or fast-forward brings the revision and the re-plan to main as two
commits. Then the revision is a change the plan did not check, and the
cycle is stale.

A cycle is done when every change in it is archived, finished or dropped.
Its file stays. A change is in at most one cycle that is not done.
`grove check` reports a cycle file that is malformed:

- it does not start with its fields between two lines of `---`;
- it has no `title`, `created`, `changes` or `planned`;
- it has a field Grove does not know;
- it lists a change that is no change; or
- it lists a change that another cycle not done lists:

> grove/cycles/261021-q2w9e-later.md:4: changes: 261019-p4x8d is also in cycle 261020-m4k8t, which is not done

For a cycle not done, it also reports a `depends_on` that names no cycle,
and cycles that follow each other in a loop:

> grove/cycles/261021-q2w9e-later.md:6: depends_on: 261020-m4k8t waits on this cycle in turn, so neither can run
