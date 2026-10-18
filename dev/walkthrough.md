# Walkthrough

The working test script. It runs every step of
[Getting started](../docs/getting-started.md) on the ledger fixture
(`dev/fixture.md`), with what to type and what to expect, in Claude Code and
again in Codex. The documentation is the specification (CONSTITUTION.md rule
8); each step here names the page section it tests. Grove's output is quoted
only where it carries meaning.

Setup: in the checkout or worktree under test, `eval "$(just dev-path)"`
in the shell the agent is started from. It builds `bin/grove` and puts that
`bin/` first on PATH, so `grove` there is this build and the agent inherits
it; other shells keep the everyday `grove` from `just install`. After
changing a guide, template or a schema, `just build` is enough; after
changing skill text, run `grove init` again in the fixture. `just fixture ~/tmp/ledger`, `cd ~/tmp/ledger`. For
Codex, use a second fixture and type `$grove-…` where this says `/grove-…`.

## Main line: allow editing older records

1. **Set up** (getting started §1). `grove init`
   > Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.

   `.claude/skills/` and `.agents/skills/` each hold `grove-explore`,
   `grove-propose`, `grove-apply` and `grove-archive`. Under `context` in
   `grove.yaml`, write "Python 3 standard library only; tests run with
   python3 -m unittest." Then `git add -A && git commit -m "Set up Grove"`.
   Running `grove init` again rewrites the skills and leaves `grove.yaml` as
   edited.

2. **Explore** (§2). `git switch -c older-edits`, start `claude`, and type
   `/grove-explore Editing an old record fails because its category was retired`.

   The session reads the code, reproduces
   `error: category 'misc' is not allowed` for `ledger edit 3 --note x`, and
   asks one question at a time. It notices that the README ("every record")
   and the CHANGELOG ("new records") disagree, and asks which is intended. I
   answer that older records keep their category when it is unchanged, and
   leave open what happens if the category is changed. `git status` shows
   nothing new: exploring wrote no file.

3. **Propose** (§3). `/grove-propose` in the same conversation. The session
   creates `grove/changes/261015-k3m9q-allow-editing-older-records/` with
   `.grove.yaml` (its title and created time), `proposal.md` (my decision quoted, the category question open and
   blocking no task, since T-01 leaves the category unchanged),
   `specs/record-editing/spec.md` (an ADDED requirement as an EARS sentence
   with scenarios), and `tasks.md` (T-01 with its acceptance). It runs
   `grove check` until:
   > No problems in grove/.

   Its instructions, `grove instructions propose`, give
   `## What Grove requires of propose` before `## How to propose`.

4. **Review and approve** (§4).
   `git add -A && git commit -m "Propose allowing edits to older records"`,
   then `grove review`:
   > Merging older-edits into main approves change 261015-k3m9q as future work; it does not implement it.

   Then each document, labelled change or maintained spec, and the diff.
   `git switch main && git merge --no-ff older-edits`, then `grove changes`:
   > 261015-k3m9q  Allow editing older records  approved

5. **Apply** (§5). `git switch -c older-edits-impl`, a fresh `claude`,
   `/grove-apply 261015-k3m9q`. The open question blocks no task, so the
   session asks nothing, works T-01 with a test, and writes `implementation.md` with
   T-01's status and evidence. `python3 -m unittest` passes.

6. **Archive** (§6). `/grove-archive 261015-k3m9q`. The session adds the
   requirement to a new `grove/specs/record-editing/spec.md`, with
   `grove/specs/record-editing/.grove.yaml` beside it, runs `grove check`,
   then `grove archive 261015-k3m9q`:
   > Archived 261015-k3m9q to grove/changes/archive/261015-k3m9q-allow-editing-older-records/; nothing is committed.

   Commit, `git switch main && git merge --no-ff older-edits-impl`. `grove specs`:
   > record-editing  Record editing  grove/specs/record-editing/spec.md

   The spec says "When a user edits an existing record and leaves its
   category unchanged, the system shall save the edit, even if that category
   is no longer allowed for new records." `grove changes` lists nothing. `grove`
   opens the reader on the spec.

7. **Next** (§Next). A fresh `/grove-explore Editing a missing record crashes`
   cites that requirement as agreed and lists no approved changes. Its
   instructions are under 12 KiB.

## Other paths

- **Propose from a description** (§3). On a new branch, with no exploration,
  `/grove-propose Editing a missing record should print an error` writes a
  change.
- **Revise a draft** (§3). `/grove-propose <change-id>` in a new conversation
  is given that change's files in full and edits them.
- **A question that blocks a task** (§5). On a new branch,
  `/grove-propose Editing a missing record should print an error`; when
  asked whether it should exit 2 like other errors, defer it. The proposal
  names T-01 as blocked by that question. Merge it, then, on another branch,
  run apply headless: `claude -p "/grove-apply <change-id>"`. It builds
  nothing, `implementation.md` says `Status: blocked` under T-01 with the
  question, and its last message names the question. `grove archive
  <change-id>` refuses.
- **A title that drifted** (reference, "Files"). Change `title` in a
  change's `.grove.yaml` and run `grove check`:
  > grove/changes/<change-id>-<slug>/proposal.md:1: the first heading should be "# <new title>", the title in .grove.yaml
- **A change without its proposal** (reference, "Files"). After step 3,
  delete the change's `proposal.md` and run `grove check`. It exits non-zero:
  > grove/changes/<change-id>-<slug>/proposal.md: missing; it is the change's main document

  `grove instructions propose <change-id>` still finds the change, so a
  session can write the proposal again.
- **A mistake in grove.yaml** (reference, "grove.yaml"). Add `derp: [x]`
  under `rules` and a top-level `rule: {}` in `grove.yaml`. `grove check`
  reports both and exits non-zero:
  > grove.yaml:<line>: unknown rules key "derp"; kinds are proposal, change-spec, tasks, implementation, spec, and skills are explore, propose, apply, archive

  > grove.yaml:<line>: unknown field "rule"; grove.yaml takes root, target, agents, context and rules

  `grove changes` still works.
- **Commit as you go** (reference, `rules`). With
  `apply: [Commit after each task, naming its task ID.]` under `rules`
  in `grove.yaml`, `/grove-apply` commits T-01 with its ID in the
  message.
- **Drop a change** (concepts, "The loop"). `/grove-archive <change-id>` on
  an approved change nobody will build asks why, writes `Status: dropped`
  and the reason in `implementation.md`, and moves it without touching the
  specs. Before that, `grove archive <change-id>` says
  > T-01 is not done in implementation.md; archive a change when every task is done or the file says Status: dropped

## Decided

Owner, asked whether to add each:

- A review session or `instructions review <id>`, 2026-10-14: "No".
- A top-level `implement <change-id>`, 2026-10-14: "No".
- Whether a session can confirm the Grove files it wrote are valid,
  2026-10-14: "yes"; how: "by calling the CLI?".
- Bubble Tea and Glamour, 2026-10-14: "We need bubble tea and glamour".
- Codex reads repository skills from `.agents/skills/`, typed
  `$grove-<skill>`, 2026-10-15: "I know this is how it works. We don't
  need to verify this."
- The rest of this page's shape is decisions D-33 to D-43 of 2026-10-15.

## Open questions

- Does the reader showing `implementation.md` count as reading it under rule
  3? Owner: "No, but im not certain." Until settled, the reader shows it.

## Assumptions to confirm

Agent inference, not the owner's decisions (rule 25).

- `grove.yaml` fields and defaults: `root: grove`, `target: main`,
  `agents: [claude, codex]`, `context`, and `rules` keyed by `proposal`,
  `change-spec`, `tasks`, `implementation` and `spec`, or by a skill.
- `init` never overwrites `grove.yaml` and always rewrites the skills.
  Its second run says "Kept grove.yaml and wrote the Grove skills for claude
  and codex; nothing is committed." A field left out of `grove.yaml` takes
  its default; an agent other than `claude` or `codex` stops `init`.
- A skill runs `grove` from the PATH of the agent's shell. Claude
  Code's skills pass the user's words as `$ARGUMENTS`; Codex has no such
  substitution, so its skills name them `<args>` and say what that stands
  for.
- `review` takes no argument and compares the current branch's last commit
  with the target, so uncommitted files are not shown. A change counts as
  approved by the merge when the target holds no change with its ID; with
  none, the first line is "Merging <branch> into main approves no new
  change."
- `instructions propose` tells a change ID (`YYMMDD-xxxxx`) from a
  description by its form.
- `archive` only moves the folder; the agent does the reconciliation, from
  `instructions archive`. Both read the change from the working tree, where
  `implementation.md` is, and neither checks that the change is approved.
- Every skill commits only when the owner or its `rules` ask; otherwise
  the owner commits after archiving (getting started §6).
- A blocked task's question is answered by revising the change with
  propose and merging it, not by the owner telling the applying session.
- Apply leaves files no task names, such as the README or a changelog,
  as they are; propose makes such updates tasks when the owner wants them.
- `changes` lists approved changes that are not archived; drafts on the
  current branch are shown by `review`, not by `changes`.
- IDs use the UTC date.
- Explore's instructions offer no document to fill in.
- A change's ID is the start of its folder name. `check` reports a change
  folder whose name does not start with an ID.
- A new change's folder name is its ID and a slug the agent chooses from the
  title; Grove allocates the ID and `created` timestamp.
