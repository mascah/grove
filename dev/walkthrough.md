# Walkthrough

The working test script. It runs every step of
[Getting started](../docs/getting-started.md) on the ledger fixture
(`dev/fixture.md`), with what to type and what to expect, in Claude Code and
again in Codex. The documentation is the specification (CONSTITUTION.md rule
8); each step here names the page section it tests. Grove's output is quoted
only where it carries meaning.

Setup: `just build`, `G=~/GitHub/mascah/grove-next2/bin/grove-next2`,
`just fixture ~/tmp/ledger`, `cd ~/tmp/ledger`. For Codex, use a second
fixture and type `$grove-…` where this says `/grove-…`.

## Main line: allow editing older records

1. **Set up** (getting started §1). `$G init`
   > Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.

   `.claude/skills/` and `.agents/skills/` each hold `grove-explore`,
   `grove-propose`, `grove-implement` and `grove-archive`. Under `context` in
   `grove.yaml`, write "Python 3 standard library only; tests run with
   python3 -m unittest." Then `git add -A && git commit -m "Set up Grove"`.
   Running `$G init` again rewrites the skills and leaves `grove.yaml` as
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

3. *(Optional)* **Baseline** (§2). I ask for a baseline spec. The session
   writes `grove/specs/record-editing/spec.md` with the rules I confirm, and
   "Known discrepancies" where the code differs. Then it runs `$G check`.

4. **Propose** (§3). `/grove-propose` in the same conversation. The session
   creates `grove/changes/261015-k3m9q-allow-editing-older-records/` with
   `proposal.md` (my decision quoted, the category question open),
   `specs/record-editing/spec.md` (an ADDED requirement as an EARS sentence
   with scenarios), and `tasks.md` (T-01 with its acceptance). It runs
   `$G check` until:
   > No problems in grove/.

5. **Review and approve** (§4).
   `git add -A && git commit -m "Propose allowing edits to older records"`,
   then `$G review`:
   > Merging older-edits into main approves change 261015-k3m9q as future work; it does not implement it.

   Then each document, labelled change or maintained spec, and the diff.
   `git switch main && git merge --no-ff older-edits`, then `$G changes`:
   > 261015-k3m9q  Allow editing older records  approved

6. **Implement** (§5). `git switch -c older-edits-impl`, a fresh `claude`,
   `/grove-implement 261015-k3m9q`. The session asks the open question before
   building, works T-01 with a test, and writes `implementation.md` with
   T-01's status and evidence. `python3 -m unittest` passes.

7. **Archive** (§6). `/grove-archive 261015-k3m9q`. The session adds the
   requirement to `grove/specs/record-editing/spec.md`, creating it if step 3
   was skipped, runs `$G check`, then `$G archive 261015-k3m9q`:
   > Archived 261015-k3m9q to grove/changes/archive/261015-k3m9q-allow-editing-older-records/; nothing is committed.

   Commit, `git switch main && git merge --no-ff older-edits-impl`. `$G specs`:
   > record-editing  Record editing  grove/specs/record-editing/spec.md

   The spec says "When a user edits an existing record and leaves its
   category unchanged, the system shall save the edit, even if that category
   is no longer allowed for new records." `$G changes` lists nothing. `$G`
   opens the reader on the spec.

8. **Next** (§Next). A fresh `/grove-explore Editing a missing record crashes`
   cites that requirement as agreed and lists no approved changes. Its
   instructions are under 12 KiB.

## Other paths

- **Propose from a description** (§3). On a new branch, with no exploration,
  `/grove-propose Editing a missing record should print an error` writes a
  change.
- **Revise a draft** (§3). `/grove-propose <change-id>` in a new conversation
  is given that change's files in full and edits them.
- **Drop a change** (concepts, "The loop"). `$G archive <change-id>` on an
  approved change nobody will build moves it without touching the specs.

## Decided

Owner, asked whether to add each:

- A review session or `instructions review <id>`, 2026-10-14: "No".
- A top-level `implement <change-id>`, 2026-10-14: "No".
- Whether a session can confirm the Grove files it wrote are valid,
  2026-10-14: "yes"; how: "by calling the CLI?".
- Bubble Tea and Glamour, 2026-10-14: "We need bubble tea and glamour".
- Codex reads repository skills from `.agents/skills/`, typed
  `$grove-<workflow>`, 2026-10-15: "I know this is how it works. We don't
  need to verify this."
- The rest of this page's shape is decisions D-33 to D-43 of 2026-10-15.

## Open questions

- Does the reader showing `implementation.md` count as reading it under rule
  3? Owner: "No, but im not certain." Until settled, the reader shows it.

## Assumptions to confirm

Agent inference, not the owner's decisions (rule 25).

- `grove.yaml` fields and defaults: `root: grove`, `target: main`,
  `agents: [claude, codex]`, `context`, and `rules` keyed by `proposal`,
  `change-spec`, `tasks`, `implementation` and `spec`.
- `init` never overwrites `grove.yaml` and always rewrites the skills.
  Its second run says "Kept grove.yaml and wrote the Grove skills for claude
  and codex; nothing is committed." A field left out of `grove.yaml` takes
  its default; an agent other than `claude` or `codex` stops `init`.
- A skill runs the binary by the absolute path `init` was run from. Claude
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
  `instructions archive`.
- `changes` lists approved changes that are not archived; drafts on the
  current branch are shown by `review`, not by `changes`.
- IDs use the UTC date.
- Explore's instructions include the maintained spec's template and guide,
  for a baseline the owner asks for (getting started §2); they offer no
  other document.
- A new change's folder name is its ID and a slug the agent chooses from the
  title; Grove allocates the ID and `created` timestamp.
- The docs call the binary `grove`; this trial builds `bin/grove-next2`.
