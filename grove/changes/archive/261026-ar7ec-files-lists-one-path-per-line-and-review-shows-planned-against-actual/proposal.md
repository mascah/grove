# Files lists one path per line, and review shows planned against actual

## Problem

The review of Grove 0.5.0 on memex (`dev/memex-review-2026-10-26.md`,
finding 3) observed:

- `grove check` rejected a folder path with a trailing slash three times
  (`crates/memex-core/tests/fixtures/vault/` in jnqjq T-01 and T-02,
  `app/src-tauri/icons/` in syssv T-02). `fs.ValidPath` in
  `internal/doc/doc.go` at ea7fbaf rejects the slash, and the message,
  "should list repo-relative paths in backticks, separated by commas", does
  not say why.
- Agents wrapped the comma line at 80 columns every time: the parser copes,
  readers see a run-on line.
- `Files:` was edited during implement in 3 of 5 changes (7 entries,
  `Cargo.lock` three times), and 3 of 5 changed files no list named even at
  archive. An agent-kept list of changed files is not reliable, and nothing
  shows the owner what a branch changed beyond what was planned.

At ea7fbaf, `naming` in `internal/cli/check.go` compares a named file with
the changed paths exactly, so a folder entry, were it accepted, would match
no changed file.

## Intended outcome

- A task's `Files:` line may list one backticked path per bullet. The comma
  line stays valid. A path ending in `/` names a folder and every file
  under it, in every place Grove reads a task's named files.
- A malformed `Files:` line is reported with both forms named.
- During implement, `tasks.md`'s `Files:` lines stay as approved. A file
  the owner agrees to change beyond a task's `Files:` and `Delegated:`
  lines is recorded on an `Also changed:` line under that task in
  `implementation.md`, with the reason. Unattended never writes one.
- `grove review` lists each file the branch changes outside the document
  root as planned, agreed (with the reason) or named nowhere, from Git.

## Scope

### In scope

- The bullet form and folder entries of `Files:`, in `grove check`,
  `grove instructions implement`, `grove run` (through the implement
  instructions) and a cycle's staleness.
- The `Files:` message.
- `Also changed:` lines in `implementation.md`, written by implement with
  the owner and read by `grove review`.
- Implement editing no `Files:` line.
- The tasks guide and template, the implementation guide, the implement
  instructions, the docs, the walkthrough and the trials.

### Out of scope

- The rule 3 amendment letting Grove read `Also changed:`: its own pull
  request, with no code, before this change is approved.
- Embedding or listing the change in the implement instructions: change
  261026-804yw.
- `grove check` exiting non-zero on moved-file notices.

## Decisions

- "`Files:` takes one backticked path per bullet; the comma line stays
  valid, so not breaking. A trailing `/` names a folder and everything
  under it." (owner, `dev/memex-review-2026-10-26.md`, decision 3)
- "`tasks.md` is frozen during implement: only a propose revision edits
  `Files:`." Accepted in the round: implement writes or edits no `Files:`
  line, on an existing task or a new one; its agreed revisions may still
  edit the proposal, the change specs and the other lines of a task; a
  revision that changes a task's files goes through
  `grove instructions propose <change-id>` on the branch ("yes"). Reason: the
  approved plan stays what review compares against.
- "A file the owner agrees to change beyond `Files:`/`Delegated:` goes on
  an `Also changed:` line under that task in `implementation.md`, with the
  reason; unattended never adds one." Accepted in the round: one line per
  file, `` Also changed: `path` — reason ``; Grove reads the backticked
  path, the rest is the reason ("yes").
- "Grove reports planned against actual in `grove review` (rule 6), from
  git ... Not a `dev/` script: Grove is a standalone CLI." Accepted in the
  round: review lists every file the branch changes outside the document
  root against `<target>...HEAD`, each `planned` with the task, `agreed`
  with the task and reason, or `named nowhere`, counting the changes the
  branch holds open and those it holds archived that the target does not;
  it prints no list when there is no such change, and still exits 0
  ("yes").
- A folder entry matches every file under it in `check`'s moved-file
  report, implement's changed-on-target line, a cycle's staleness and
  review (round, "yes").
- The malformed-line message names both forms in one sentence (round,
  "yes").
- "Review reading `Also changed:` needs a rule 3 amendment first, in its
  own pull request (owner agreed)." This change is not approved until that
  amendment is in the target; `depends_on` names changes only, so it does
  not name it (round, "yes").
- The trials' count of edits outside `Files:`/`Delegated:` accepts the
  paths an `Also changed:` line names and parses the bullet form (round,
  "yes").
- The instruction text stays at or under 8,000 words (rule 12) by
  tightening the text this change touches (round, "yes").

## Open questions

None.

## Related documents

- Requirements: [propose](specs/propose/spec.md),
  [implement](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
