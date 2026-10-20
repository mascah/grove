# Coordinating changes

[Getting started](getting-started.md) takes one change from an idea to an
updated spec. When several changes are approved at once, some must land
before others. This page shows how a change says what it waits on.

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
