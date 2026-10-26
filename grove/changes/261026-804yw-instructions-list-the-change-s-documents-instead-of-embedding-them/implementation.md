# Implementation: Instructions list the change's documents instead of embedding them

## Tasks

### T-01 — List the change's documents by path

Status: done

Evidence:
- `internal/instruct/instruct.go`: `Instructions.Change` is a list of
  `Listed` (a path, and `Differs`), with no content; `build` lists the
  change's files; `Implement` takes the names of the files that differ;
  `Markdown` prints "The change" as one line to read every document with
  the agent's own file tool, then one bullet per path, a differing one
  with `git show <commit>:<path>`. `Revise` and `Archive` take no files.
- `internal/cli/instructions.go`: for an approved change, `differing`
  compares each of the target's files with the working-tree copy, a
  missing copy counting as differing. It leaves out `implementation.md`,
  which records progress and is not approved text (rule 3), so progress
  written on the branch is never named as differing. `change` keeps its
  signature, since `internal/cli/archive.go`, which no task names, calls
  it.
- `internal/instruct/requires/implement.md`: a document named as
  differing is built as `git show` prints it, unless the owner agreed to
  the revision on this branch. `internal/instruct/requires/archive.md`
  needed no new text (the listing's own line says to read every
  document); both files were tightened to keep the instruction text under
  8,000 words.
- `docs/reference.md` (after the instructions table) and
  `dev/walkthrough.md` ("Revise a draft", now given paths, and a new
  "A large change") say the documents are listed.
- `internal/cli/walkthrough_test.go`: `TestPropose` (revise lists the
  change spec's path, prints no requirement text), `TestImplementInstructions`
  (draft and approved list paths and print no requirement text; after an
  edit of `tasks.md` the unattended instructions name it as differing
  from main with `git show <main>:<path>`, and without one none is named;
  a member of a cycle still gets the cycle file in full), `TestArchive`
  (lists `implementation.md` and the change spec, prints no requirement
  text).
- `just check` passed: 4020 Go lines, 1678 test lines, 7998 instruction
  words.

### T-02 — Implement prints only the implementation template

Status: done

Evidence:
- `internal/instruct/schemas/default/schema.yaml`: implement `writes` is
  `[implementation]`; `bugfix` takes default's skills, so the same.
- `internal/instruct/requires/implement-owner.md`: an agreed edit is made
  by running `grove instructions propose <change-id>` and following it,
  never editing a `Files:` line or `.grove.yaml` by hand (this also covers
  261026-ar7ec's "Revise the change on its branch", in fewer words).
- `docs/getting-started.md` and `dev/walkthrough.md` (main line step 5,
  "Implement asks, and revises with the owner") say so.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` expects
  the implementation template and no proposal, change spec, tasks or
  `.grove.yaml` template, and the propose path; `TestSchemas` checks
  bugfix's implement instructions likewise.
- The on-fixture acceptance with the owner (an agreed revision made after
  running propose) needs an agent session: not run here, see Scenarios.
- `just check` passed: 7970 instruction words.

### T-03 — Trials count listed documents not read

Status: not started

## Scenarios

## Reconciliation
