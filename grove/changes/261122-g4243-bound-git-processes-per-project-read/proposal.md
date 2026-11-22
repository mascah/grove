# Reading the project value starts a bounded number of git processes

<!-- A bug fix describes future work. It becomes approved future work when
it is merged into the target branch, which says nothing about
implementation. -->

## Observed

[GitHub issue 196](https://github.com/mascah/grove/issues/196) reports that one
`grove changes --json` on a project with 9 changes and no cycles or runs
starts 44 `git` processes (21 `diff`, 6 `rev-parse`, 5 `ls-tree`, 5
`cat-file`, 2 `merge-base`, and one each of `status`, `rev-list`, `log`,
`for-each-ref`). That is about 25 ms a call. In tend's test suite,
about 925 such calls started about 85% of its 44,000 `git` processes.

Measured at 522f060 with a shim on `PATH` that logs each `git` call and
runs the real Git, on a local clone of this repository with its remote
removed:

| State of the clone | `git` processes per `grove changes --json` |
| --- | --- |
| no unmerged branch, no open change | 11 |
| 5 unmerged branches, each one empty commit | 16 |
| the 3 newest archived changes moved back into `changes/` | 36 |
| the 6 newest archived changes moved back into `changes/` | 47 |

In Grove's own checkout, with 51 unmerged refs and no approved change, one call
started 33 `git diff`s.

The project value that `grove changes --json` prints, and that `grove serve`
answers `/api/project` with, is built by `server.value` in
`internal/cli/serve.go` at 522f060. It starts processes like this:

- For each change, `moved` in `internal/cli/check.go` runs one `git diff`
  between the change's base and the target. `rewritten` reads the whole tree
  at the change's base with `gitx.TreeFS` (one `ls-tree`, one `cat-file`),
  only to read the maintained specs the change modifies or removes.
- For each branch the target has not merged, `scan` in `serve.go` runs one
  `git diff --raw target...tip`.
- For each record on `grove/runs`, `value` runs one `git for-each-ref`.
- For each cycle that is not done, `changedSincePlan` (`rev-parse`, `log`,
  `cat-file`) runs once for the cycle's stale mark, again in `runnable` and
  again in `stalePlans`, which also resolves and reads the target again.
  `movedSinceBase` resolves the target again too.

`grove check` shares `diagnose`, so it pays the per-change cost as well.

## Expected

What the owner accepted ("yes") in the conversation that proposed this: the
number of `git` processes per call does not grow with the number of changes or run
records. It may grow by 2 for each branch that holds Grove documents the
target's history lacks, and by 3 for each cycle that is not done.

What each command prints stays as it is, with one exception. Where a file
was renamed since the change's base and the change's tasks name both its old
and its new path, through a folder or one by one, `grove check`'s one
moved-file problem for the change names both paths, and `grove instructions
implement` lists each path under the task that names it. Today both name only
the new path, and only when the user's Git config detects renames.

## Reproduction

1. Build `grove` at 522f060.
2. Put first on `PATH` a `git` script that appends its arguments to a log
   file and then runs the real `git` with them.
3. Clone this repository into a scratch folder and remove the clone's
   `origin` remote.
4. In the clone, run `grove changes --json` once and count the logged calls.
5. Move the three newest folders in `grove/changes/archive/` back into
   `grove/changes/` and commit. Run `grove changes --json` again and count
   again: the count rises with each reopened change.

## The requirement it breaks

No maintained requirement covers how many `git` processes a command
starts. The owner decided that none is agreed: the bound stays observed
behavior, guarded by the test T-01 adds.

## Decisions

- Schema: "bugfix is fine."
- Accepting 2 more processes for each branch with Grove documents the
  target's history lacks: "yes". A strictly fixed count would mean parsing
  Git's binary tree objects, which costs much of the Go line budget left.
- No change spec, and the reproducing test puts a logging shim in front of
  the real `git` on `PATH`: "yes". The shim runs real Git, so it is not the
  fake Git the constitution's rule on tests (rule 17) bars.
- The fix, as reviewed in three rounds by Codex and accepted by the owner
  ("ok"):
  - Unmerged branches are diffed in one `git diff-tree --stdin --always -r
    --raw --no-abbrev --no-renames -- <root>`, fed the branch tip, then the
    target, for each branch. diff-tree takes
    the first ID as the commit and the second as its parent, so this diffs
    from the target to the tip. Today's diff starts from the merge-base
    instead. The two give the same files once the existing filter drops
    every blob the target's history already holds.
  - The maintained specs at each change's base are read in one
    `git cat-file --batch` of `<base>^{commit}` and
    `<base>:<root>/specs/<capability>/spec.md`. The `^{commit}` reply both
    tells a base that is no commit, for the existing diagnostic, and gives
    the base's full ID.
  - Moved files are found in one `git diff-tree --stdin --always -r
    --name-only --no-renames`, fed the target, then the full base, for each
    distinct base, with no path filter. The existing `naming` matches the
    result to tasks. A base that is no commit is skipped, as today.
    `--no-renames` matches every other diff in this code and causes the one
    output change named under Expected. `moved`, which `grove instructions
    implement` still runs for its one change, takes `--no-renames` too, so
    the two commands agree.
  - `--always` makes each request print a header even when its diff is
    empty, so output blocks are matched to requests in input order.
  - Run records' branch dates come from one `git for-each-ref` over
    `refs/heads` and `refs/remotes/origin`.
  - `changedSincePlan` runs once per cycle in one project value.
    `server.value` passes the target commit and documents it already holds
    to `diagnose`. `grove check` passes none, and it keeps resolving the
    target only when a cycle or task needs it, so a checkout without the
    target still checks.
- Left as is: a base missing from a shallow clone still costs one `git
  rev-parse --is-shallow-repository` per diagnostic. That only happens on
  an error.

## Open questions

- Routine choices inside the agreed approach, which can wait for review
  and block nothing: where the batching code lives between
  `internal/gitx` and `internal/cli`, and the name of any `gitx` function
  that runs Git with standard input. Recommended: one exported `gitx`
  function that runs Git with given standard input, since Git runs only
  in `gitx`, and the parsing next to its callers in `internal/cli`.
  Decides how T-01 is built.

## Related documents

- Tasks: [tasks.md](tasks.md)
