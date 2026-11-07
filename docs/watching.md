# Watching the project

`grove serve` shows where a project stands on one page in your browser:
what `grove run` would do with each cycle, which branches hold work the
target does not have yet, the runs recorded in this checkout, how changes
and cycles depend on each other, what the target's history did to the
record, and the maintained specs. Click any change, cycle, spec or file
to read it beside the list.

The page reads the target branch, this checkout's branches, the working
tree and `.git/grove/`. It writes nothing, fetches nothing from `origin`
and starts no agent, so what it shows is what this checkout has: `git
fetch` first to see your teammates' latest branches. The
[serve spec](https://github.com/mascah/grove/blob/main/grove/specs/serve/spec.md)
has the exact rules.

## Start it

In a checkout of the project:

```sh
grove serve
```

Grove prints the address and keeps serving until you press Ctrl-C:
> http://127.0.0.1:7070/

Open it in your browser; Grove does not open one for you. The page reads
the project again every 10 seconds; uncheck "every 10 s" to stop it, or
press Refresh. To serve on another port, or a second project at once,
give another loopback address:

```sh
grove serve --addr 127.0.0.1:8080
```

The page is only for this machine: Grove refuses an address that is not
loopback, and if the port is taken it says so and exits.

## Readiness

One card for each cycle that is not done, with the answer `grove run`
would give it today, the first that applies:

- a member needs a person, as the cycle's `needs_human` says;
- the cycle waits on another cycle that is not done;
- a member waits on a change that is neither in the cycle nor archived as
  finished;
- a member changed on the target since the plan, so the plan is stale and
  needs `/grove-plan <cycle-id>` again;
- otherwise `grove run` would build it, or `grove run <cycle-id>` would,
  when it is not marked `unattended`.

Each card lists the members in the order they would be built, with their
status and where each is being built, then the members already archived.
Below the cards are the approved changes in no cycle, each ready or
waiting, the cycles done, and the problems `grove check` would report on
the target's documents. The tab counts the cycles not done.

## In progress

One card for each branch, and for the working tree, that holds a change's
files the target never held: a draft the target does not hold yet, or a
change it approves that is being built. Each change shows how many of its
tasks are done, each task's status from that branch's `implementation.md`,
and the files that are new there. Click a file to read that branch's
version.

Grove looks at the branches not merged into the target and counts a file
only when the target's history never held that version of it. A branch
squash-merged into the target holds nothing the target never held, so it
is left out, however much the target has moved since.
A branch you pushed shows once, saying origin has it too; a teammate's
branch you fetched but never checked out shows as `origin/<name>`, only
on origin. Newest first.

Below the branches are the runs recorded in this checkout, newest first:
each attempt's cycle, when it started and ended, its outcome and its
branch, as [Running unattended](running.md) describes them.

## Graph

Each cycle not done is a band holding its members, and the approved
changes in no cycle are below them. An arrow runs from a change to each
change that depends on it, and a dashed arrow from a cycle to each cycle
that follows it. A green box is ready, an orange one waits, and a dashed
blue one is being built on a branch. Check "show archived changes and done
cycles" to see the past as well.

## History

A timeline with one row for each cycle not done and each change open, and
below it each
commit on the target's first-parent history that did something to the
record, newest first, with what it did:

- **approved**: added a change;
- **revised**: changed one of an approved change's files, such as its
  proposal or `implementation.md`;
- **archived**: moved a change to `changes/archive/`;
- **planned** and **re-planned**: added or changed a cycle's file;
- **spec**: edited a maintained spec.

A commit that touched nothing under the document root is not listed.

## Specs

The maintained specs on the target, with their descriptions.

## The pane

Clicking a change, cycle or spec opens it in the pane on the right. A
change shows its proposal and a button for each of its files on the
target, and for its files new on each branch. Markdown is rendered and a
cycle's YAML shown as text. A link from one document to another under the
root opens in the pane. A refresh keeps the document you are reading, and
where you scrolled to, in both columns.
