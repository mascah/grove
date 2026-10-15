# Walkthrough

This page is the specification (CONSTITUTION.md rule 8). It follows one small
change on the ledger fixture (`docs/fixture.md`) from `init` to the maintained
spec that describes the new behavior. Steps are what I do. Grove's output is
quoted only where it carries meaning. The commands named here are the whole
command surface; anything not named is not built.

Setup: `just build`, `G=~/GitHub/mascah/grove-next2/bin/grove-next2`,
`just fixture ~/tmp/ledger`, `cd ~/tmp/ledger`.

## Main line: allow editing older records

1. **Init.** `$G init`
   > Wrote grove.yaml, grove/brief.md, .claude/skills/grove-implement/SKILL.md and a .gitignore entry for .grove/; nothing is committed.

   *(Optional)* edit `grove/brief.md`. A brief still identical to the sample is
   not given to sessions. Then `git add -A && git commit -m "Set up Grove"`.

2. **Explore.** `$G explore "Allow editing older records"`
   > Workspace .grove/worktrees/261014-n6va0-allow-editing-older-records on branch grove/261014-n6va0-allow-editing-older-records, from main@1a2b3c4.

   Claude Code opens in the workspace with the instructions already sent. The
   session reads the code, reproduces the failure, saves what it observed in
   `grove/explorations/261014-n6va0-allow-editing-older-records.md`, and asks
   me before it lists options. I answer. My answers are saved as quotes; what I
   leave undecided stays under open questions.

3. *(Optional)* **Baseline.** I ask for a baseline spec. The session writes
   `grove/specs/record-editing/spec.md` in the workspace: the rules I say are
   intended, and, where the code differs, "Known discrepancies"; then it runs
   `$G check` as in step 4. It becomes the
   maintained spec when the branch is merged in step 5. Without a baseline, the
   change spec adds requirements to a capability that has no maintained spec,
   and reconciliation in step 6 creates the spec.

4. **Propose.** I say "write the proposal". The session runs
   `$G instructions propose 261014-n6va0` and follows it in the same
   conversation. It creates `grove/changes/261014-k3m9q-allow-editing-older-records/`
   with `proposal.md`, `specs/record-editing/spec.md` (ADDED requirements as
   EARS sentences with scenarios) and `tasks.md` (T-01 with its acceptance). It
   asks rather than assumes, and writes undecided points as open questions, not
   decisions. It runs `$G check` in the workspace and fixes each problem it
   prints, until:
   > No problems in grove/.

   I `/exit`.
   > Draft: change 261014-k3m9q in .grove/worktrees/261014-n6va0-allow-editing-older-records, 4 files uncommitted.
   > Continue: $G explore 261014-n6va0 --resume

5. **Commit and merge.** In the workspace,
   `git add -A && git commit -m "Propose allowing edits to older records"`.
   *(Optional)* in main, `$G review 261014-k3m9q`:
   > Merging grove/261014-n6va0-allow-editing-older-records into main approves change 261014-k3m9q as future work; it does not implement it.

   followed by each document the merge adds or changes, labelled change,
   exploration or maintained spec, and the diff. Then in main,
   `git merge --no-ff grove/261014-n6va0-allow-editing-older-records` (a squash
   merge works too). The merge is the approval (rule 4). `$G changes`:
   > 261014-k3m9q  Allow editing older records  approved

6. **Implement.** In main, `git switch -c older-edits`, start my own `claude`,
   and run `/grove-implement 261014-k3m9q`. The skill runs
   `$G instructions implement 261014-k3m9q` and follows the result; Grove
   creates no workspace and launches nothing (rule 2). The session works T-01
   with a test, writes `implementation.md` in the change folder (each task's
   status and evidence, and which change specs it reconciled), reconciles the
   change spec into `grove/specs/record-editing/spec.md`, and commits. I
   `git switch main && git merge older-edits`.

7. **See the result.** `$G specs`:
   > record-editing  Record editing  grove/specs/record-editing/spec.md

   The spec now says "When a user edits an existing record and leaves its
   category unchanged, the system shall save the edit, even if that category is
   no longer allowed for new records." `$G changes` still says `approved`:
   Grove does not read `implementation.md` (rule 3); people do. `$G` opens the
   reader on the specs and changes, `implementation.md` shown as written. A
   fresh `$G explore "Editing a missing record crashes"` cites that requirement
   as agreed behavior, lists 261014-k3m9q as approved work, and keeps what the
   code does apart from both.

## Leaving and coming back

- `/exit` at any point prints the two lines of step 4 (naming the exploration
  when there is no change yet). Nothing is committed or approved.
- `$G workspaces` lists each workspace and what is in it:
  > 261014-n6va0  Allow editing older records  change 261014-k3m9q draft, 4 files uncommitted
- `$G explore 261014-n6va0 --resume` reopens the latest conversation Grove
  started for 261014-n6va0, never another one run in that folder.
- `$G explore 261014-n6va0` starts a new conversation in the same workspace,
  given the exploration and its change in full. It becomes the one `--resume`
  reopens.

## A second change on the same capability

The fresh session of step 7. Its instructions are under 12 KiB: the brief and
the new exploration in full, and `record-editing`, 261014-k3m9q and the first
exploration listed with the revision to read. The session reads the spec,
reproduces the `StopIteration`, and asks whether a missing id is an error. I
answer and say "write the proposal". The new change's
`specs/record-editing/spec.md` adds one requirement to the maintained spec.
Commit in the workspace, merge, then `/grove-implement <change-id>` on a branch
and merge. `$G specs` still lists one spec, now with both requirements;
`$G changes` lists both changes as approved.

## Command surface

`init`, `explore "<description>"`, `explore <id> [--resume]`, `workspaces`,
`instructions propose <id>`, `check`, `instructions implement <change-id>`,
`review <change-id>`, `specs`, `changes`, and no arguments for the reader. The
one skill is `/grove-implement`.

## Decided

Owner, 2026-10-14, asked whether to add each:

- A review session or `instructions review <id>`: "No".
- A top-level `implement <change-id>` beside `instructions implement`: "No".
- Whether a session can confirm that the Grove files it wrote are valid
  before I commit: "yes"; how: "by calling the CLI?".
- `--resume` must restore the conversation the exploration was started in,
  never an arbitrary later one: "This is not ok in any scenario … It is not
  optional." With several, it reopens "The latest one grove started for that
  id" (2026-10-15).
- Whether the modules Bubble Tea and Glamour pull in count against rule 16:
  "We need bubble tea and glamour".

## Open questions

- Does the reader showing `implementation.md` as written count as reading it
  under rule 3? Owner: "No, but im not certain." Until settled, it shows it.
- Should `changes` keep saying `approved` for a change that has been
  implemented and reconciled? Owner: "Im not sure." Until settled, it does,
  because no code reads `implementation.md` (rule 3).

## Assumptions to confirm

Agent inference, not the owner's decisions (rule 25).

- Ids are `YYMMDD-xxxxx`, names are `<id>-<slug>`, frontmatter carries
  `created` and an optional `description` shown in indexes (from the
  2026-10-14 feedback, a lean under rule 22).
- The harness is `claude`. Grove starts each conversation with
  `--session-id <uuid>`, records the ID in the worktree's own Git directory,
  and `--resume` runs `claude --resume <uuid>`.
- Workspaces live in `.grove/worktrees/` on branches `grove/<id>-<slug>`.
  `grove.yaml` holds only the root (`grove`) and the target branch (`main`).
- Grove commits nothing; I commit and merge with ordinary Git.
- The command is `check`. It reads the Grove documents in the working tree it
  runs in, prints each problem as `path:line: message`, and exits non-zero
  when there is one.
- `go.mod` requires directly only the four modules rule 16 names; what
  Bubble Tea and Glamour pull in is listed there as indirect.
