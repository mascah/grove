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

Check: `bin/grove check` on this branch reports one problem: "grove/cycles/261118-5b98a-cycle-files-in-markdown.md: planned: 261118-emse7 changed on main since this plan, so grove claim skips this cycle". This change's T-02 names the folder `grove/cycles/`, and the two commits on main since that cycle's `planned` (c9ef062, dc8d335) edited the cycle file under its old `.yaml` name, which the staleness rule excuses only as the cycle file's own path, now `.md`. It is reported against main's copy of this change, so it stays until the merge that archives 261118-emse7 reaches main, after which no open member names `grove/cycles/` and the report goes. No code was added for it: a cycle file's old name is a state rule 13 leaves unguarded, and only a member naming `grove/cycles/` meets it.
Check: until main holds `.md` cycle files, this branch's `bin/grove` finds no cycle on main: `bin/grove instructions implement 261118-5b98a` says "main holds no approved change 261118-5b98a", and member instructions print no `## The cycle`. The `grove` built from main still reads main's `.yaml` cycle for the cycle's own instructions.

## Scenarios

## Reconciliation
