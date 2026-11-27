# Trim the archive

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

`grove archive` moves a finished change to `changes/archive/` with every
document it had. Within weeks the archive holds several times the words of
the maintained specs, and much of it is requirement text that later
changes rewrote. An agent searching the project finds those old versions
beside the live ones, under the same headings.

Grove reads two things from an archived change: its `.grove.yaml`, and, if
it was dropped, the `Status: dropped` line of its `implementation.md`.
This recipe deletes every other archived document from the working tree.
Git keeps them all, and each spec's History still names the changes that
shaped it.

## Before you start

- **Trim only what the target holds.** While a change archived on a branch
  waits for its merge, its documents are what the reviewer reads.
- **Leave a change that an open change depends on.** For an open change
  that lists an archived one in `depends_on`, `grove check` reads the
  archived change's specs to tell which rewritten requirements that
  dependency explains. Trim it once the open change is archived too.
  `grove changes` lists what is open.

## Steps

**Terminal:** on a branch from the target, delete all but each change's
`.grove.yaml`, and a dropped change's `implementation.md`. The paths assume
the default root, `grove`.

```sh
git switch -c trim-the-archive
cd grove/changes/archive
grep -L '^Status: dropped' */implementation.md | xargs rm
find . -type f ! -name .grove.yaml ! -name implementation.md -delete
find . -mindepth 2 -type d -empty -delete
cd -
```

**Terminal:** check, commit and merge.

```sh
grove check
git add -A grove/changes/archive
git commit -m "Trim the archive"
```

To read a trimmed document later, ask Git for the commit that removed it
and show the file as it was before:

```sh
f=grove/changes/archive/<folder>/proposal.md
git show "$(git log -1 --format=%h -- "$f")^:$f"
```

## What to expect

It was tried on Grove's own archive of 62 changes, with no change open:
348 files and 202,293 words became 62 files and 2,504 words.

- **Nothing Grove reports changes.** With the trim merged, `grove check`
  reports no problems, `grove changes` lists every change as archived and
  every cycle as done, and `grove serve` answers as before.
- **A dropped change stays dropped.** On a copy with one change marked
  dropped, that change kept its `implementation.md` and Grove still
  listed it as dropped.
- **History links still resolve.** Each links a folder that now holds only
  `.grove.yaml`; the command above reads the rest.
- **Most of what goes was already stale.** Of the 253 requirement texts
  in the archive's change specs, 145 no longer matched the maintained
  spec: 97 had been rewritten since and 48 named a requirement that no
  longer exists.
- **Not tried:** trimming a change that an open change depends on. The
  rule above comes from reading how `grove check` works, not from running
  it.
