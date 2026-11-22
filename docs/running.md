# Running unattended

`grove claim` hands a piece of work to an agent that nobody watches: a
cycle, or an approved change in no cycle. It picks the piece you name,
or else a cycle you have marked `unattended`. It records the attempt on
the branch `grove/runs`, creates the attempt's branch in the checkout and
prints the piece's ID. Then you start the agent, in any harness, bounded
as you like:

```sh
K=$(grove claim) && timeout 2h claude -p "/grove-implement $K"
```

Grove starts no process but Git. It does not start, bound, stop or log
the agent; your shell, `timeout` and the harness do. The attempt's
changes are then in flight until they are archived on the target, as
finished or dropped. While any change is in flight, a `grove claim` given
nothing claims nothing, from every checkout (see "The pause"); a claim
that names a piece is refused only for what is in flight or overlaps it
(see "Claim a change you name").

Grove does not choose what is claimed next or how many attempts run at
once; whoever runs `grove claim` does. Grove does not schedule itself: run the line by hand, or from a timer as
in [Run on a timer](guides/run-on-a-timer.md). The
[runs spec](https://github.com/mascah/grove/blob/main/grove/specs/runs/spec.md)
has the exact rules.

## Set up a checkout

Grove builds in the checkout you run it in and creates none. Give each
project a clone of its own that nobody works in:

```sh
git clone git@github.com:you/ledger.git ~/grove-runs/ledger
cd ~/grove-runs/ledger
```

What you prepare, since Grove does not:

- **The remote.** `origin` holds the target branch. Grove fetches it, and
  the only branch it pushes is `grove/runs`. A checkout with no `origin`
  works too: see "Without a remote".
- **The harness.** The examples here use Claude Code, as
  `/grove-implement`; with commands in place of skills, it is
  `/grove:implement`. Another harness's skill files from `grove init`
  work the same way, given that harness's command for one headless
  prompt.
- **Claude Code's permissions.** The session asks you nothing, so it can do
  only what its settings allow. Put them in `.claude/settings.local.json`
  in the checkout, such as

  ```json
  {
    "permissions": {
      "defaultMode": "acceptEdits",
      "allow": ["Bash(grove:*)", "Bash(git:*)", "Bash(gh:*)", "Bash(python3:*)"]
    }
  }
  ```

  with your project's test command in place of `python3`. Keep the file
  out of `git status`, since `grove claim` refuses a checkout with anything
  uncommitted:

  ```sh
  echo .claude/settings.local.json >> .git/info/exclude
  ```

- **GitHub.** `gh auth login` as the account whose pull requests these
  are. It must be allowed to push to `origin`, since `grove claim` pushes
  `grove/runs` there.
- **`PATH`.** `grove`, `claude`, `git`, `gh` and `timeout` are on the
  `PATH` of the shell that runs the line; the session inherits it. On a
  Mac, `timeout` comes with Homebrew's `coreutils`.

## Tell the agent to commit and open a pull request

Grove's skills commit nothing unless a rule asks. For an unattended run,
add these under `rules` in `grove.yaml`, and merge them into the target:

```yaml
rules:
  skills:
    implement:
      - Commit after each task, naming its task ID.
  delivery:
    implement:
      - When you are done, push the branch to origin and open a pull request with gh pr create, listing the changes archived and each blocked task with its question.
```

`grove claim` reads `grove.yaml` from the target it fetched, not from the
checkout's branch. So the rules count once they are merged. Lines under
`implement` are also given when you build with the agent yourself, so
word them to hold for both.

The line under `delivery` is how the work leaves the checkout. When what
starts the agent delivers the work itself, such as a pipeline that has the
branch reviewed before it opens the pull request, have it run
`/grove-implement --runner $K` instead: the `delivery` lines are left out,
and the session pushes nothing, opens no pull request and leaves its
commits on the attempt's branch.

## Claim a cycle

```sh
grove claim
```

Grove refuses while the checkout has uncommitted changes. Then it:

1. Fetches the target and `grove/runs` from `origin`, if the checkout has
   one. If the local target branch has commits `origin`'s does not, it
   refuses and leaves them alone.
2. Refuses while the project is paused (see "The pause").
3. Selects one cycle on the fetched target: of the cycles that pass
   every test below, the one with the earliest `created`. A cycle passes
   when:
   - it is marked `unattended`;
   - it is not done;
   - it follows no cycle that is not done (see
     [Planning cycles](planning.md), "Follow another cycle");
   - none of its changes waits on a change that is neither in the cycle
     nor archived as finished;
   - its plan is not stale (see [Planning cycles](planning.md), "Mark it
     unattended").

   With none, it prints nothing on standard output, says
   > grove: nothing to claim on main at 243a63d

   and exits non-zero. `grove check` says why a marked cycle is skipped.
4. Records the attempt as one commit on `grove/runs` and pushes it to
   `origin` (see "The record").
5. Creates the branch `grove/<attempt>` from the fetched target,
   `<attempt>` being a new ID such as `261021-6jbj1`, and switches to it.
   It fast-forwards the local target branch to the fetched target too.
6. Prints the cycle's ID, alone, on standard output:
   > 261020-m4k8t

In the line at the top, `$(...)` keeps that ID, and `&&` starts the
agent only when `grove claim` exits 0. So when there is nothing to claim,
the project is paused or a refusal names something to fix, nothing
starts and the line exits non-zero. Otherwise the agent builds the
cycle's changes in turn on the attempt's branch, as "Building a cycle"
in [Coordinating changes](coordinating.md) says. It archives
each one it finishes and, as the rules say, it commits, pushes and opens
the pull request. The session's output is the harness's own; redirect it
to a file to keep it.

## The record

Each attempt is one commit on `grove/runs` that holds no file. Its
message is the record:

```sh
git log -1 --format=%B origin/grove/runs
```

```json
{
  "id": "261021-6jbj1",
  "started": "2026-10-21T09:12:40-07:00",
  "cycle": "261020-m4k8t",
  "changes": ["261019-p4x8d", "261019-r2w7k"],
  "commit": "243a63d9c1e5b7f04a8e2d6c3b9f1a7e5d4c2b80",
  "host": "runner.local",
  "checkout": "/Users/you/grove-runs/ledger",
  "branch": "grove/261021-6jbj1"
}
```

It names the cycle and its changes in build order, the target commit, the
host and checkout that claimed it, and the attempt's branch. A record for
a change in no cycle lists that change alone and has no `cycle`. It says an
attempt started, not how the session ended nor whether any task is done:
`implementation.md` and the pull request say that. `grove/runs` never
shares a commit with your work, the attempt's branch holds only the
target's commits until the agent adds its own, and `git status` shows
nothing of the record. [Watching the project](watching.md) shows the
records with the latest commit on each attempt's branch.

## Two checkouts at once

Each claim needs a checkout of its own: the attempt's branch is created
in the checkout that claimed it, and the agent works there. To run two
attempts at once, clone the project twice.

The push to `origin` succeeds only while `origin`'s `grove/runs` is where
this checkout just read it. If two checkouts, on one host or two, claim
at the same moment, `origin` takes one record. The other checkout's push
is refused, so it reads `grove/runs` again and decides again as a fresh
`grove claim` would. Given nothing, it finds the project paused after the
first one's attempt. Given a piece, it records its own attempt when the
two share nothing, and otherwise says the piece is in flight or shares a
path with the first, and creates no branch.

A push refused for another reason, such as an account that may not push
to `origin`, records nothing and creates no branch. Grove says so,
quoting Git.

## Claim a cycle you name

To build a planned cycle once, without marking it, name it:

```sh
K=$(grove claim 261020-m4k8t) && timeout 2h claude -p "/grove-implement $K"
```

Grove claims that cycle whether or not it is marked `unattended`. The
pause does not stop it: it is refused only when the cycle is in flight or
overlaps a change in flight, as "Claim a change you name" says. A
`grove claim` given nothing, such as a timer's, still claims only marked
cycles, and is paused while this attempt is in flight.

If the cycle cannot be built, Grove records nothing, creates no branch,
prints nothing on standard output and exits non-zero. It says why, which
is one of these:

- there is no such cycle on the target;
- it is done;
- it follows a cycle that is not done;
- a change in it waits on one that is neither in the cycle nor archived
  as finished;
- its plan is stale.

Such as:

> grove: cycle 261021-q2w9e waits on cycle 261020-m4k8t, which is not done

## Claim a change you name

An approved change that no cycle lists is claimed by its ID:

```sh
C=$(grove claim 261019-t6n3b) && timeout 2h claude -p "/grove-implement $C"
```

Grove records an attempt naming that change alone, creates its branch and
prints the change's ID. A change stale against its `base`, or with an open
question, is claimed all the same; the agent reconciles or answers it as
it builds. If the change cannot be built, Grove records nothing, creates
no branch, prints nothing on standard output, exits non-zero and says
why, which is one of these:

- there is no such cycle or change on the target:
  > grove: main at 243a63d holds no cycle or change 261019-zzzzz
- it is archived:
  > grove: change 261019-p4x8d is archived
- a cycle not done lists it, which you claim instead:
  > grove: 261019-r2w7k is in cycle 261020-m4k8t; claim that with grove claim 261020-m4k8t
- it waits on a change not archived as finished on the target; finished
  on an attempt's branch does not count:
  > grove: 261019-t6n3b waits on 261019-p4x8d, which is not archived as finished

A claim you name, of a cycle or a change, is not stopped by the pause. It
is refused only for what is in flight. A change is in flight while a
record on `grove/runs` names it and the target does not hold it archived;
a cycle is in flight while one of its changes is. Grove refuses a piece
that is itself in flight:

> grove: 261019-t6n3b is in flight in attempt 261021-6jbj1, which runner.local started at 2026-10-21T09:12:40-07:00 on branch grove/261021-6jbj1

and a piece that overlaps a change in flight:

> grove: 261019-t6n3b shares ledger/cli.py with 261019-r2w7k, in flight in attempt 261021-6jbj1

A change's planned files are the paths its tasks' `Files:` lines name,
and the maintained spec of each capability it has a change spec for,
such as `grove/specs/record-editing/spec.md`; a cycle's are those of its
changes not archived. Two pieces overlap when a planned file of one is
a planned file of the other, or lies under a folder the other names. Two
changes on one capability overlap whatever requirements they touch,
since each fold adds to that spec's History. Grove reads the planned
files from the target, never from an attempt's branch, so a file an
agent adds to its tasks on its branch counts only once merged. The check
only makes collisions rarer: an agent may touch a file its tasks do not
name, and merging the pull requests is what settles it.

Pieces that share nothing can be claimed at once, each in a checkout of
its own (see "Two checkouts at once"). `--retry` claims despite what is in
flight or overlaps it (see "Retry").

## Without a remote

In a checkout with no remote named `origin`, `grove claim` fetches
nothing. The local target branch stands in for the fetched target, and
the record goes on the local `grove/runs`, which Grove updates only while
it is where Grove just read it. Nobody opens a pull request. So in place
of the push and `gh pr create` rule, have the agent commit what it
leaves, such as the archive:

```yaml
      - When you are done, commit everything left uncommitted.
```

and merge the attempt's branch into the target yourself:

```sh
git switch main
git merge grove/261021-6jbj1
```

That merge archives the changes and ends the pause, as merging the pull
request does with `origin`.

## The pause

The pause is for a `grove claim` given nothing. While any change is in
flight, that is, named by a record on `grove/runs` and not yet archived
on the fetched target, as finished or dropped, it claims nothing, from
this checkout or any other clone of the project, and names the newest
record in flight:

> grove: the project is paused after attempt 261021-6jbj1, which runner.local started at 2026-10-21T09:12:40-07:00 on branch grove/261021-6jbj1; once you have dealt with it, run grove claim --retry

Merging the pull request archives the changes. Once no change is in
flight, the pause ends, and the next `grove claim` selects the next cycle
it can build. Other merges, even of another marked cycle, do not end the
pause, and an older attempt still in flight keeps it though a newer one's
change is archived. Nothing expires it with time: Git cannot tell that a
session somewhere has died. A claim you name is not paused.

A change left blocked by a question stays open, so the project stays
paused. There are two ways on:

- Answer the question with `/grove-propose <change-id>`, which re-plans
  the cycle on the same branch. Merge that, then `grove claim --retry`.
  The revision and the re-plan must reach main as one commit, a squash or
  a merge commit. Rebased or fast-forwarded as two commits, the revision
  makes the plan stale (see [Planning cycles](planning.md), "Re-plan a
  cycle").
- Drop the change with `/grove-archive <change-id>`, and merge that.

## When a session is stopped

Grove is not running while the agent is, so it does not know when or how
the session ends. `timeout 2h` stops it after two hours, and Ctrl-C
stops it sooner. The branch and its uncommitted files stay as they were,
and the project stays paused.

Grove does not recover an attempt. Look at the branch and the session's
output. Then commit what is worth keeping and open its pull request
yourself, or discard it:

```sh
git checkout -- . && git clean -fd
```

Until the checkout is clean, `grove claim` refuses, with `--retry` too:

> grove: the checkout has uncommitted changes; commit, stash or discard them first

## Retry

```sh
K=$(grove claim --retry) && timeout 2h claude -p "/grove-implement $K"
```

records a new attempt even while the project is paused, as an unpaused
`grove claim` would. With an ID, before or after `--retry`, it claims
that piece despite its being in flight or overlapping a change in
flight; the newest record naming a change is its attempt, so the earlier
one is no longer in flight. It still refuses a piece that cannot be
built, such as a change that waits. It starts from the fetched target, not from the last
attempt's branch. To keep the tasks that branch finished, merge its pull
request first; otherwise the agent builds them again. The earlier records
and their branches stay as they are.

Different projects, each with its own `grove/runs`, are claimed
independently.

## Draft a change with no one in the session

An agent nobody watches can write a proposal too, from a request such as
a tracker issue, a review finding or a problem `grove check` reports.
Grove reads no tracker: whoever starts the session passes the request in.

```sh
timeout 30m claude -p '/grove-propose --unattended "Editing an unknown ID prints an error and exits 2"'
```

The agent creates a branch, writes one change from the request and asks
nothing. What it would have asked you in the round is in the proposal's
open questions, each with the answer it recommends and the task it
decides; its Decisions quote only the request. A split, an open change
that may need to land first or the `bugfix` schema is an open question
too: it never decides one alone. It offers no `/grove-implement` on the
branch, so the draft is approved only when you merge it, never with its
code. An open question you merge unanswered is built as
`/grove-implement` builds any choice left open: its best answer, with a
`Check:` line for the reviewer, or a blocked task where a wrong answer
would be costly.

To have it commit and push the draft, add under `rules` in `grove.yaml`:

```yaml
rules:
  skills:
    propose:
      - With no one in the session, commit the change folder once grove check passes, push the branch to origin and open a pull request with gh pr create.
```

Lines under `propose` are also given when you propose with the agent
yourself, so this one names the session it is for.

A request that says no outcome, such as "Look into the report command",
gets no change: the agent writes nothing and its last message says what
the request lacks. The session cannot revise a draft from your review;
run `/grove-propose <change-id>` yourself for that. The
[propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md)
has the exact rules.
