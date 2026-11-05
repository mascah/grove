## T-01 — Implement stops at checkpoints and hands off in the branch

Result: the implement instructions, the `implementation.md` guide and
template, the docs and the walkthrough say what "Stop at a checkpoint",
"Take a checkpoint's answer" and "Settle what the change does not" in
[specs/implement/spec.md](specs/implement/spec.md) say, as the proposal's
decisions give them.

Files:
- `internal/instruct/requires/implement.md`
- `internal/instruct/requires/implement-owner.md`
- `internal/instruct/requires/implement-unattended.md`
- `internal/instruct/requires/implement-cycle.md`
- `internal/instruct/schemas/default/documents/implementation/guide.md`
- `internal/instruct/schemas/default/documents/implementation/template.md`
- `internal/cli/walkthrough_test.go`
- `docs/concepts.md`
- `docs/getting-started.md`
- `docs/coordinating.md`
- `dev/walkthrough.md`

Delegated: test layout and private helper names

Interfaces: `Status: to check`, and the `Check:` and `Answer:` lines
under a task in `implementation.md`.

Constraints: rule 12's instruction-text budget; rule 3, so the agent,
not the owner, writes `implementation.md`; no Go code beyond the test.

Acceptance:
- `grove instructions implement <change-id>` no longer says to tell the
  owner how to try each task and take their feedback before the next; it
  says to build every task that waits on no checkpoint, then stop once
  and put each checkpoint to the owner, and still to stop at once when
  the change's approach is wrong.
- With and without `--unattended`, the instructions say what a
  checkpoint is, when to build a guess (`Status: to check`) and when to
  build nothing (`Status: blocked`), to write a `Check:` line, to leave a
  task that builds on a checkpoint unbuilt naming it, to take an answer
  as "Take a checkpoint's answer" says, and to list checkpoints first in
  the report, the cycle's report included. The wording test checks the
  new text and the absence of the old.
- The `implementation.md` guide and template show `Status: to check`,
  `Check:` and `Answer:`.
- `grove archive` on a change with a task `Status: to check` refuses it,
  as it refuses any task not done; a test shows it.
- `docs/concepts.md`, `docs/getting-started.md` and
  `docs/coordinating.md` describe checkpoints, how a later session on
  the branch picks them up, with the owner or unattended, and that an
  answer changing the change still needs `/grove-propose`.
- A walkthrough step builds a change unattended to a `Status: to check`
  task, then answers it in a fresh `/grove-implement` on the same branch,
  and the change is archived.
- `just check` passes.
