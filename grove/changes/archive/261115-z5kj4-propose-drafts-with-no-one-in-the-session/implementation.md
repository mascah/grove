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
passes at 7,953 instruction words; `TestPropose` in
`internal/cli/walkthrough_test.go` checks each acceptance line.

Check: open question 1 built as recommended: a change's ID, or no request, gets one sentence and exit 1, Grove's status for every command error.
Check: open question 2 built as recommended: one file, `internal/instruct/requires/propose-unattended.md`, appended to what Grove requires of propose and saying which of its lines and the guide's steps it replaces.
Check: the request is printed under the existing "## Owner's description" heading, and the unattended text says "The description below is the request"; a heading of its own would cost Go and words for no behavior.
Check: after case (f)'s runs, the owner asked for two rewordings: only a question from the check of the change says it must be answered before review, since both agents put that on every open question; and a split, a change to land first or another schema is an open question only "where there is one", since Claude Code wrote a dependencies question with nothing to recommend. The full trial run at 413e0fe includes them.
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

Status: done

`dev/trials/f.sh` and case (f) in `dev/trials.md` add the sixth case;
`dev/trials/run.sh` runs it by default. It sets `turns=0`, so no reply
reaches the session, and its `good` asks the waiting check once of the last
message. One run per harness, one judging (`RUNS=1 JUDGINGS=1 just trials
f`): Codex passed; Claude Code failed two checks, its open question on the
schema naming "the whole change" rather than a task, and its question on
dependencies giving no recommended answer. Both judged counts were zero
but words read (Claude Code 571 in an earlier run, Codex 96).

Check: the full `just trials` run at 413e0fe, three runs per case and harness, is in the pull request. Case (f): Codex 3/3 good runs; Claude Code 0/3, each run failing two open-question checks: the schema question names no task, since it decides the whole change, and one question (routine defaults, other open changes) names no task or recommends nothing. Claude Code's unneeded questions were 1 per run, each an open question the repository already answers.
Check: case (e) failed every run in both harnesses on a defect of its own: its `good` reads `implementation.md` at the change's folder, which `/grove-implement` moves to the archive when every task is done, as case (c)'s `good` already allows. This change does not touch `dev/trials/e.sh`.
Check: case (f) makes no reply even when the waiting check says the agent waits. A first run let the runner answer "Yes, go with what you recommend.", and the agent turned its open questions into decisions; no one is in an unattended session to answer.
Check: both agents marked every open question "must be answered before review", taking the line meant for what the check of the change would ask as meant for all of them. The instruction text has 8 words of budget left in this change; rewording is left to review.
Also changed: `CLAUDE.md` — its layout table said the trials have five cases.

## Scenarios

- propose "A request from a tracker issue": `dev/trials/f.sh` case (f), one run per harness: both asked nothing, wrote one change on `unknown-id` that `grove check` passes, and named files, task IDs and open questions; Codex passed `good`, Claude Code failed two open-question checks (see T-03)
- propose "The flag inside the quotes": `internal/cli/walkthrough_test.go` TestPropose, passed
- propose "A schema the caller names": `internal/cli/walkthrough_test.go` TestPropose, passed
- propose "Without the flag": `internal/cli/walkthrough_test.go` TestPropose, passed; the output is also byte for byte main@df86b0b's, compared by hand
- propose "The request settles one decision": case (f), one run per harness: both Decisions quote the request for the exit status and an open question asks where the check goes; Claude Code recommended `find` raising `ValueError`, Codex recommended the check in `ledger/cli.py`, so the scenario's recommendation held in one of two runs
- propose "The request settles nothing more": unverified; no run used "Add ledger delete ID" unattended, and TestPropose checks only that the instructions say so
- propose "A request to look into something": a headless Claude Code run of `/grove-propose --unattended Look into the report command` on the fixture wrote no change folder and said the request names no outcome and what would; passed, and no test covers it
- propose "Two parts": unverified; not run, and TestPropose checks only the instruction text
- propose "An open change touches the same code": unverified; not run, and TestPropose checks only the instruction text
- propose "A change's ID": `internal/cli/walkthrough_test.go` TestPropose, passed
- propose "No request": `internal/cli/walkthrough_test.go` TestPropose, passed
- propose "A requirement no task covers": `internal/cli/walkthrough_test.go` TestPropose checks the instruction text, passed; the agent's behavior is unchanged by this change and not run
- propose "Scope and a task disagree": `internal/cli/walkthrough_test.go` TestPropose checks the instruction text, passed; not run
- propose "Scope and a task disagree with no one to ask": TestPropose checks the instruction text, passed; the agent's behavior is unverified, since no run had such a contradiction
- propose "A new change on its branch": `internal/cli/walkthrough_test.go` TestPropose, passed
- propose "Revising an approved change": `internal/cli/walkthrough_test.go` TestPropose, passed
- propose "A draft nobody attended": TestPropose checks no `/grove-implement` offer, passed; in case (f) Claude Code's last message said merging approves the change and named no `/grove-implement`
- schemas "Proposing a bug fix": `internal/cli/walkthrough_test.go` TestSchemas, passed
- schemas "Exploring a bug": `internal/cli/walkthrough_test.go` TestSchemas, passed
- schemas "A bug reported with no one in the session": unverified as written, since no run used "Editing a missing record crashes"; in case (f), whose request describes wrong behavior too, both runs kept `schema: default` and asked about `bugfix` as an open question

## Reconciliation

- propose: applied to grove/specs/propose/spec.md
- schemas: applied to grove/specs/schemas/spec.md
