# Implementation: Propose drafts a change with no one in the session

<!-- Progress on change 261115-z5kj4, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose instructions for a session with no one in it

Status: done

`grove instructions propose --unattended [--schema <name>] "<request>"`,
and the flag as the first word of the one quoted argument, print the
propose instructions with `requires/propose-unattended.md` after
`requires/propose.md` and a task sentence naming the request and no one in
the session. Without the flag the output is byte for byte what main@df86b0b
printed for `grove instructions propose "x"` on the ledger fixture (IDs and
times masked). A change's ID or no request prints "grove: an unattended
propose writes a new change from a request" and exits 1. `just check`
passes at 7,946 instruction words; `TestPropose` in
`internal/cli/walkthrough_test.go` checks each acceptance line.

Check: open question 1 built as recommended: a change's ID, or no request, gets one sentence and exit 1, Grove's status for every command error.
Check: open question 2 built as recommended: one file, `internal/instruct/requires/propose-unattended.md`, appended to what Grove requires of propose and saying which of its lines and the guide's steps it replaces.
Check: the request is printed under the existing "## Owner's description" heading, and the unattended text says "The description below is the request"; a heading of its own would cost Go and words for no behavior.
Check: "prints what it printed at efdb8af" was compared by hand against a build of main; the test asserts only that the flag-less instructions lack the unattended text and keep their existing lines.

Also changed: `.claude/skills/grove-propose/SKILL.md` — `grove init` rewrote it from the new propose description in `skills.yaml`.
Also changed: `.agents/skills/grove-propose/SKILL.md` — the same, for codex.

### T-02 — Docs, reference and walkthrough show the mode

Status: done

`docs/running.md` has "Draft a change with no one in the session": the
`claude -p '/grove-propose --unattended "…"'` line, a
`rules.skills.propose` line that commits, pushes and opens a pull request,
and what a request with no outcome leaves. `docs/reference.md` has a row
for `grove instructions propose --unattended`. `dev/walkthrough.md`
"Running unattended" step 8 drafts "Editing an unknown ID prints an error
and exits 2" and says what Decisions and the open questions hold. The
Purpose of `grove/specs/propose/spec.md` names the mode.

Check: the docs' `rules.skills.propose` line says "With no one in the session, …" because lines under a skill are also given to attended sessions; one line keeps commit, push and pull request together, as the implement rules on the same page do.
Check: walkthrough step 8 runs in `~/tmp/ledger-run2`, the operator's clone that already holds `.claude/settings.local.json`, not in `~/tmp/ledger`, so the headless session has its permissions; it was not run by hand in this session, and T-03's case (f) runs the same request headless.

### T-03 — A skill trial for a draft nobody attends

Status: not started

## Scenarios

## Reconciliation
