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

Status: done

Evidence:
- `dev/trials.md` defines "Listed documents not read" as a judged count
  (six counts, four judged) and the results table has its column.
- `dev/trials/run.sh`: the judge prompt lists the files under
  `grove/changes/` at the end of the run and asks for
  `listed_documents_not_read`; each row prints it; every row width has the
  extra column.
- `HARNESSES=claude dev/trials/run.sh c` at e1be131 (plus these edits):
  `| (c) | Claude Code | 0 | 0 | 0 | 0 | 0 | 913 |`; the judge: "Both runs
  read `.grove.yaml`, `proposal.md`, `specs/record-editing/spec.md` and
  `tasks.md` before acting". The same run's owner session wrote
  `` Also changed: `ledger/records.py` — … the owner agreed `` under T-01,
  as 261026-ar7ec intends. Codex was not run.

## Scenarios

- implement "A large change": `internal/cli/walkthrough_test.go` TestImplementInstructions (paths listed, no requirement text), passed; that the agent reads every one is unverified here, and is what the trials' "Listed documents not read" counts
- implement "A document revised on the branch": TestImplementInstructions names an edited `tasks.md` as differing with `git show <main>:<path>`, passed; that the agent builds the owner's agreed revision needs an agent session with the owner, unverified
- implement "A document that differs, unattended": TestImplementInstructions (`--unattended`), passed for the naming; that the agent builds what `git show` prints needs an agent run, unverified
- implement "Archiving": TestArchive (lists `implementation.md` and the change spec, no requirement text), passed
- implement "Templates": TestImplementInstructions and TestSchemas (bugfix), passed
- implement "An agreed revision": unverified; it needs an agent session with the owner (TestImplementInstructions checks the owner's instructions say to run `grove instructions propose <change-id>`)
- propose "Revising a large change": TestPropose (paths listed, no requirement text, the templates still printed), passed

## Reconciliation

- implement: applied to grove/specs/implement/spec.md
- propose: applied to grove/specs/propose/spec.md
