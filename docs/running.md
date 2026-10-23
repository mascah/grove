# Running apply unattended

`grove run` builds a cycle you have marked `unattended` without anyone
typing `/grove-apply`. In a checkout you set up for one project, it starts
one Claude Code session on the apply skill for that cycle, waits for it
within a time limit, and then pauses the project until the cycle is merged
or dropped. It builds nothing else: an approved change in no cycle, or in a
cycle you have not marked, is yours to build with `/grove-implement`. To
send it one change, plan a cycle of one (see
[Planning cycles](planning.md)). Grove does not
schedule itself: run it by hand, or from a timer as in
[Run apply on a timer](guides/run-apply-on-a-timer.md).

## Set up a checkout

Grove builds in the checkout you run it in and creates none. Give each
project a clone of its own that nobody works in:

```sh
git clone git@github.com:you/ledger.git ~/grove-runs/ledger
cd ~/grove-runs/ledger
```

What you prepare, since Grove does not:

- **The remote.** `origin` holds the target branch. Grove fetches it and
  pushes nothing.
- **The harness.** `grove.yaml` on the target lists `claude` under
  `agents`. `grove run` types the apply skill as that form has it,
  `/grove-apply` for skills or `/grove:apply` for commands, followed by the
  cycle's ID.
- **Claude Code's permissions.** The session asks you nothing, so it can do
  only what its settings allow, and Grove passes it no permission flags.
  Put them in `.claude/settings.local.json` in the checkout, such as

  ```json
  {
    "permissions": {
      "defaultMode": "acceptEdits",
      "allow": ["Bash(grove:*)", "Bash(git:*)", "Bash(gh:*)", "Bash(python3:*)"]
    }
  }
  ```

  with your project's test command in place of `python3`, and keep the file
  out of `git status`, since `grove run` refuses a checkout with anything
  uncommitted:

  ```sh
  echo .claude/settings.local.json >> .git/info/exclude
  ```

- **GitHub.** `gh auth login` as the account whose pull requests these
  are.
- **`PATH`.** `grove`, `claude`, `git` and `gh` are on the `PATH` of the
  shell that runs `grove run`; the session inherits it.

## Tell the agent to commit and open a pull request

Grove's skills commit nothing unless a rule asks. For an unattended run,
add these under `rules` in `grove.yaml` and merge them into the target:

```yaml
rules:
  skills:
    apply:
      - Commit after each task, naming its task ID.
      - When you are done, push the branch to origin and open a pull request with gh pr create, listing the changes archived and each blocked task with its question.
```

`grove run` reads `grove.yaml` from the target it fetched, not from the
checkout's branch, so the rules count once they are merged.

## An ordinary run

```sh
grove run
```

Grove refuses while another `grove run` is going in this checkout, or while
it has uncommitted changes. Then it:

1. Fetches the target from `origin`. If the local target branch has
   commits `origin`'s does not, it refuses and leaves them alone.
2. Selects one cycle on the fetched target: the one with the earliest
   `created` that is marked `unattended`, is not done, has no change under
   `needs_human`, has no change waiting on one neither in the cycle nor
   archived as finished, and whose changes are as its plan read them. With
   none, it says
   > Nothing to apply on main at 243a63d.

   `grove check` says why a marked cycle is skipped.
3. Creates the branch `grove/run-<attempt>` from the fetched target,
   `<attempt>` being a new ID such as `261021-6jbj1`, fast-forwards the
   local target branch to it, and starts
   `claude -p "/grove-apply 261020-m4k8t"` on the new branch.

The agent builds the cycle's changes in turn, as "Applying a cycle" in
[Coordinating changes](coordinating.md) says, archives each one it finishes,
commits, pushes and opens the pull request, as the rules say. When the
session ends, Grove says
> Attempt 261021-6jbj1 ended: exit status 0; its output is /Users/you/grove-runs/ledger/.git/grove/261021-6jbj1.log.

Exit status 0 means the session ended normally, not that every task is
done; `implementation.md` and the pull request say that. Grove keeps two
files per attempt in the checkout's `.git/grove/`, which `git status` does
not show and nobody commits:

- `<attempt>.log`, everything the session printed.
- `<attempt>.json`, the record: when it started and `ended`, the `cycle`
  it selected and its `changes` in build order, the target `commit`, the
  `checkout`, the `branch`, the `output` file and the `outcome`, one of `exit status <n>`, `timed out`,
  `cancelled` or `failed to start`.

## The pause

After an attempt, `grove run` starts nothing while a change of the cycle
that attempt selected is not archived on the target, as finished or
dropped:

> grove: the project is paused after attempt 261021-6jbj1, whose output is /Users/you/grove-runs/ledger/.git/grove/261021-6jbj1.log; once you have dealt with it, run grove run --retry

Merging the pull request archives the changes and ends the pause; the next
`grove run` selects the next eligible cycle. Other merges, even of another
marked cycle, do not end it. A change left blocked by a question stays
open, so the project stays paused: answer the question with
`/grove-propose <change-id>`, merge the revision, re-plan the cycle with
`/grove-plan <cycle-id>` and merge that, since the revision made its plan
stale, then run `grove run --retry`; or drop the change with
`/grove-archive <change-id>` and merge that.

## When a session is stopped

A session has 2 hours unless you give another limit, such as
`grove run --timeout 30m`. When the limit passes, or you press Ctrl-C, or
`grove run` is sent `SIGTERM` (as launchd does to stop a job), Grove asks
every process in the session's process group to stop, kills those still
running 5 seconds later, and ends with
> grove: attempt 261021-6jbj1 ended: timed out; its output is /Users/you/grove-runs/ledger/.git/grove/261021-6jbj1.log

or `cancelled`. The branch, its uncommitted files and the output stay as
they were. A process that left the group, such as a server that started a
session of its own, keeps running. If `grove run` itself is killed
outright, it stops nothing: check with `ps` that no `claude -p` from it is
still running before you start another.

`claude -p` prints only its last message, so a stopped session's output is
usually empty; the branch and its files show what it did.

Grove does not recover an attempt. Look at the branch and the output, then
commit what is worth keeping and open its pull request yourself, or
discard it:

```sh
git checkout -- . && git clean -fd
```

Until the checkout is clean, `grove run` refuses:

> grove: the checkout has uncommitted changes; commit, stash or discard them first

## Retry

```sh
grove run --retry
```

starts a new attempt as an unpaused `grove run` does, whatever the latest
attempt selected. It starts from the fetched target, not from the last
attempt's branch: to keep tasks that branch finished, merge its pull request
first; otherwise the agent builds them again.

If Claude Code cannot be started, such as when `claude` is not on `PATH`,
the attempt ends `failed to start`, with the reason in its output, and the
project is paused like after any attempt. Fix it and retry.

Different projects, each in its own checkout, run independently.
