# Implementation: A cycle is a Markdown file with plan's notes

<!-- Progress on change 261118-emse7, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Grove reads, writes and shows a cycle as Markdown

Status: done

- `internal/doc`: `Load` reads `cycles/*.md`; `doc.Front` splits a file
  into its fields (the opening `---` line kept, so the YAML parser's lines
  are the file's) and its trimmed notes, `doc.Cycle.Notes`. A `.md` that
  does not start with its fields is reported, and each `.yaml` under
  `cycles/` is a collection diagnostic; no cycle is read from it.
- `internal/instruct`: `metadata/cycle.yaml` became `metadata/cycle.md`,
  fields between two lines of `---` and a comment on the notes; plan names
  `<root>/cycles/<id>-<slug>.md`; `requires/plan.md` and the default plan
  guide say what the notes hold, that there may be none, that they change
  no member, and that re-planning brings them up to date. `Markdown()`
  prints the sentence that a member's documents hold over the notes after
  the cycle file, for the cycle's session and each member's, only when
  `Instructions.Notes` is set (from `Cycle` and `InCycle`).
- `internal/cli/serve.go`: `/api/doc` answers a Markdown file that starts
  with fields as `<pre><code>` of the fields, then the notes rendered.
- Tests: `internal/doc/doc_test.go` `TestFront` (fields and notes, fields
  only, no opening `---`, no closing `---`); `internal/cli/walkthrough_test.go`
  `TestCycles` (`.md` plan path and template, notes text, a cycle left as
  YAML open and done, converted and claimed, notes not checked, a file
  without fields, line numbers of the file), `TestImplement` (notes
  sentence absent without notes, present for the cycle and `--runner`
  member with notes), `TestServe` (fields as text, notes rendered,
  `.grove.yaml` as text), and every other test's cycle files as `.md`.
- Seen on the fixture with `bin/grove`: `grove instructions plan "0.4
  release"` names `grove/cycles/<id>-<slug>.md` and prints the template;
  `grove check` reports a `.yaml` cycle and a `.md` with a heading for
  fields, exit 1.
- `just check`: passed; 8,143 instruction words (at most 8,147).

Check: the sentence that a member's documents hold over the notes is printed only for a cycle file that holds notes (open question 1), as the implement change spec's "A cycle with no notes" scenario says; it reads "The notes below the cycle's fields are its plan's advice: where they and a member's documents disagree, the member's documents hold." and is a Go string in `Markdown()`, like the other text Grove writes around a printed file.
Check: in place of notes the template holds an HTML comment listing the three kinds of note and saying to delete it, leaving nothing below the fields, when there is nothing to say (open question 2); `requires/plan.md` adds that there may be none and that the notes change no member.
Check: the `.yaml` message names the new file without backticks, "rename it to 261020-m4k8t-0-4-release.md and put its fields between two lines of ---", since no other check message quotes a path in backticks.
Check: after review, `doc.Front` closes the fields only at a line that is exactly `---`, in either line ending, so an indented `---` inside a YAML block scalar stays a field value; and the notes keep their leading indentation, dropping only blank lines above them, so notes opening with an indented code block render as one; `grove serve` drops the opening line by its length, so a CRLF file shows no `---` among its fields. TestFront and TestServe cover each.
Check: `doc.Front` is exported so `grove serve` splits a cycle file the same way `Load` does; any `.md` under the root that starts with a line of `---` is shown so, which today is only a cycle file.
Check: a `.md` cycle file that does not start with its fields is still listed as a cycle with no fields, so its problem shows under it; with no members it counts as done and is claimed by nothing.

### T-02 — Docs, walkthrough and Grove's own cycle files

Status: done

- `docs/planning.md` shows a cycle file as plan writes it, fields between
  two lines of `---` and notes below, says what the notes hold, that only
  plan writes them, that re-planning brings them up to date, and that a
  change's documents hold over them; its examples name `.md` files with
  the file's line numbers, and its malformed list adds a file that does
  not start with its fields.
- `docs/reference.md` gives `cycles/<id>-<slug>.md`, the fields between two
  lines of `---`, the notes below, and one paragraph for a project still
  holding `.yaml` cycle files: what `grove check` reports and how to
  convert one.
- `docs/watching.md`: the pane shows a cycle's fields as text and its notes
  rendered.
- `dev/walkthrough.md`: planning step 1 writes a `.md` cycle with notes;
  step 3 one with none; step 4 notes not checked, a `.yaml` cycle left over
  and converted back, and a file that does not start with its fields;
  step 5 re-planning brings the notes up to date; step 8 the implement
  instructions for a cycle whose file holds notes; the serve pane step;
  every other step names `.md` cycle files.
- `grove/cycles/`: all eleven files (the proposal counted ten at f7e0140;
  261118-5b98a was planned since) renamed to `.md` with a line of `---`
  above and below their fields, nothing else changed and no notes.
- `grep` finds no cycle file ending `.yaml` under `docs/` or in
  `dev/walkthrough.md` except where `docs/reference.md` and planning step 4
  describe one left over.

Check: after review, staleness counts a commit that edited the cycle file under its old `.yaml` name as editing the cycle file, in `internal/cli/check.go`. Without it, converting a cycle one of whose members names `grove/cycles/` made it stale, as this branch's own cycle 261118-5b98a was (its plan and mark, c9ef062 and dc8d335, edited the `.yaml`), against the cycles scenario "Converted". The earlier note here said rule 13 left this unguarded; the review showed it is a documented step, the conversion, so it is handled and tested by TestConvertedCycle.
Check: until main holds `.md` cycle files, this branch's `bin/grove` finds no cycle on main: `bin/grove instructions implement 261118-5b98a` says "main holds no approved change 261118-5b98a", and member instructions print no `## The cycle`. The `grove` built from main still reads main's `.yaml` cycle for the cycle's own instructions.

## Scenarios

- cycles "An open cycle left as YAML": internal/cli/walkthrough_test.go TestCycles ("Leave 0.4 as YAML": check's message and exit 1, `grove changes` with no "in cycle", nothing claimable), passed
- cycles "A done cycle left as YAML": internal/cli/walkthrough_test.go TestCycles (the done 0.4 cycle renamed to `.yaml` at the end), passed
- cycles "Converted": internal/cli/walkthrough_test.go TestCycles ("Convert 0.4": check clean, `planned` unchanged, `eligible` claims it), passed; TestConvertedCycle (a member naming `grove/cycles/`, the `.yaml` plan and mark commits before the conversion), passed
- cycles "Two approved changes": internal/cli/walkthrough_test.go TestCycles (plan instructions name `grove/cycles/<id>-<slug>.md` and a template with fields between two lines of `---`), passed; the agent writing the file is unverified, as `just trials` was not run
- cycles "The agent titles the cycle": internal/cli/walkthrough_test.go TestCycles ("write its title yourself from its members, without asking"), passed; agent behavior unverified
- cycles "A draft change named": unverified; agent behavior, instruction text unchanged by this change
- cycles "A member stale on its base": internal/cli/walkthrough_test.go TestCycles (the bring-up-to-date text), passed; agent behavior unverified
- cycles "What the plan found is noted": internal/cli/walkthrough_test.go TestCycles (template comment and notes text printed), passed; whether an agent writes the note is unverified
- cycles "A change left out is noted": as above, instruction text passed; agent behavior unverified
- cycles "Nothing to note": internal/cli/walkthrough_test.go TestCycles ("with nothing to say, delete this comment"), passed; agent behavior unverified
- cycles "Re-plan a stale cycle": internal/cli/walkthrough_test.go TestCycles and TestClaim (re-plan sets `planned`, fresh again), passed
- cycles "Re-plan after a member's file moved": internal/cli/walkthrough_test.go TestMovedFiles, passed
- cycles "A note the target made false": internal/cli/walkthrough_test.go TestCycles (new and revise plan instructions say to bring the notes up to date, nothing the target made false), passed; agent behavior unverified
- cycles "One change in two cycles": internal/cli/walkthrough_test.go TestCycles (`261099-zzzzz-later.md:4`), passed
- cycles "A needs_human left over": internal/cli/walkthrough_test.go TestCycles (`:6: unknown field "needs_human"`), passed
- cycles "A dependency on no cycle": internal/cli/walkthrough_test.go TestCycles (`261099-y0000-next.md:6`), passed
- cycles "Cycles that wait on each other": internal/cli/walkthrough_test.go TestCycles, passed
- cycles "No fields at the top": internal/cli/walkthrough_test.go TestCycles (`261099-zzzzx-heading.md`), passed; internal/doc/doc_test.go TestFront, passed
- cycles "The line is the file's": internal/cli/walkthrough_test.go TestCycles (`needs_human` reported on line 6 with `---` on line 1), passed
- cycles "Notes are not checked": internal/cli/walkthrough_test.go TestCycles (notes holding `needs_human: [C]`, check clean of the cycle), passed
- cycles "A letter IDs leave out": internal/cli/walkthrough_test.go TestMisnamedCycle (`261105-proto-unknown-id-cycle.md`), passed
- implement "A note a member contradicts": internal/cli/walkthrough_test.go TestImplementInstructions (cycle and `--runner` member instructions print the file with its notes and the sentence), passed; building B's T-01 in `ledger/report.py` is agent behavior, unverified
- implement "A cycle with no notes": internal/cli/walkthrough_test.go TestImplementInstructions (no sentence for either), passed
- serve "A proposal at the target": internal/cli/walkthrough_test.go TestServe, passed
- serve "A cycle file": internal/cli/walkthrough_test.go TestServe (cycle L: fields in `<pre><code>`, `<h2>Found</h2>` and a list), passed
- serve "A cycle file with no notes": internal/cli/walkthrough_test.go TestServe (cycle K's answer ends with its fields), passed
- serve "A change's metadata": internal/cli/walkthrough_test.go TestServe (`.grove.yaml` answered as text), passed
- serve "A document on a branch": internal/cli/walkthrough_test.go TestServe, passed
- serve "Raw HTML in a document": internal/cli/walkthrough_test.go TestServe, passed
- serve "A path outside the root": internal/cli/walkthrough_test.go TestServe, passed

## Reconciliation

- cycles: applied to grove/specs/cycles/spec.md, with the Purpose sentence on the cycle file as the change spec gives it
- implement: applied to grove/specs/implement/spec.md
- serve: applied to grove/specs/serve/spec.md
