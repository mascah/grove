# A misnamed cycle and a shallow clone are reported plainly

<!-- A bug fix describes future work. It becomes approved future work when
it is merged into the target branch, which says nothing about
implementation. -->

## Observed

Both were found while prototyping Grove in a GitHub agentic workflow on
2026-11-05 (local branch `prototype/gh-aw-handoff`, commit 9481bf4, never
merged), and reproduced on the ledger fixture with Grove at 8a8d064.

1. **A cycle file named without a valid ID.** With the cycle file renamed
   `grove/cycles/261105-proto-unknown-id.yaml` (IDs leave out i, l, o and
   u, so `261105-proto` is no ID):
   - `grove check` prints `No problems in grove/.` and exits 0.
   - `grove changes` prints
     `261015-c0000  Editing an unknown ID prints an error  approved  in cycle `
     with no ID after "in cycle".
   - `grove changes --json` shows no cycle for the change.
   - `grove instructions implement 261105-proto --unattended` prints
     `grove: main holds no approved change 261105-proto`.

   Loading gives such a cycle an empty ID and reports nothing
   (`internal/doc/doc.go` at 8a8d064, the cycles loop of `Load`), while a
   change folder named the same way is reported: "the folder name must
   start with a change ID, such as 261015-k3m9q" (`loadChange`, same file).

2. **A shallow clone.** In `git clone --depth 1` of the fixture, `grove
   check` prints

       grove/changes/261015-c0000-unknown-id-prints-an-error/.grove.yaml: base: f32672a865987f808097717a72b8c43b3488cb75 is not a commit in this repository
       grove/cycles/261105-pr0t0-unknown-id.yaml: planned: d65b1126ee420f58e3f90ce821d97a36cdd756dc is not a commit in this repository

   Both are commits of the repository; the clone did not fetch them
   (`stale` and `stalePlans`, `internal/cli/check.go` at 8a8d064). In the
   prototype, gh-aw checks out one commit by default, and the agent could
   not fetch more from inside its sandbox.

## Expected

1. `grove check` reports the misnamed cycle file and exits non-zero, as it
   does a misnamed change folder, and `grove changes` names no cycle for
   its members.
2. In a shallow clone, the message says the commit is missing from the
   shallow clone and how to fetch it. Elsewhere it is as today.

## Reproduction

1. `just fixture /tmp/f`, then in `/tmp/f`: `grove init`, write an approved
   change and a cycle listing it as
   `grove/cycles/261105-proto-unknown-id.yaml`, commit, run `grove check`
   and `grove changes`.
2. With the cycle named `261105-pr0t0-unknown-id.yaml` and committed,
   `git clone --depth 1 file:///tmp/f /tmp/s`, then `grove check` in
   `/tmp/s`.

## The requirement it breaks

None states either behavior at 8a8d064. The cycles spec's Purpose says a
cycle "is a file `<root>/cycles/<id>-<slug>.yaml`", and "Report a
malformed cycle" lists no misnamed file. The owner agreed to make the
first agreed and to keep the second observed behavior ("Yes to all",
2026-11-05).

## Decisions

The owner accepted every answer below in one round ("Yes to all",
2026-11-05).

1. **The cycles spec gains "Report a cycle file named without an ID".**
   The shallow-clone message stays observed behavior, with no change
   spec, and is shown in the CI guide.
2. **What check says:** "the file name must start with a cycle ID, such
   as 261015-k3m9q", mirroring the change folder's message, and it exits
   non-zero. `grove changes` then names no cycle for its members, matching
   `--json`.
3. **The shallow-clone message**, only when the repository is shallow:
   "base: <sha> is not in this shallow clone; fetch its full history, such
   as with git fetch --unshallow", and the same for `planned`. Otherwise
   the message is unchanged. Git runs only in `internal/gitx`
   (`CLAUDE.md`), which gains the check for a shallow repository.
4. **Docs:** `docs/guides/run-grove-in-ci.md` "What to expect" shows the
   shallow-clone message; the guide already checks out with
   `fetch-depth: 0`. The walkthrough gains a step for each.
5. **Routine defaults:** the `bugfix` schema; two tasks, each starting
   with a test that fails at 8a8d064; `depends_on: [261030-t2pt7,
   261101-mg8hd, 261105-gmpa4]`, since they edit `internal/cli/check.go`
   and the walkthrough files and landing first would make cycle
   261104-z05ey's plan stale; a `Delegated:` line for test layout and
   private helper names. Hiding the cycle in `grove changes` (decision 2)
   needs `internal/cli/changes.go`, added to T-01's files when writing.

## Open questions

None.

## Related documents

- Tasks: [tasks.md](tasks.md)
- Requirements: [specs/cycles/spec.md](specs/cycles/spec.md)
