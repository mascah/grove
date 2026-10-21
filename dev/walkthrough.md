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
it; other shells keep the everyday `grove` from `just install`. `grove version`
(reference, "Terminal commands") shows which one runs: this build prints a
pseudo-version naming the checkout's commit, such as
`grove 0.0.0-<date>-<commit>`, ending `+dirty` when there are
uncommitted changes. After
changing a guide, template or a schema, `just build` is enough; after
changing skill text, run `grove init` again in the fixture. `just fixture ~/tmp/ledger`, `cd ~/tmp/ledger`. For
Codex, use a second fixture and type `$grove-…` where this says `/grove-…`.

## Main line: allow editing older records

1. **Set up** (getting started §1). `grove init`
   > Wrote grove.yaml and the Grove skills for claude and codex; nothing is committed.

   The fixture has no `.claude/`, `CLAUDE.md`, `AGENTS.md` or other sign of
   an agent, so `grove.yaml` lists `claude: skills` and `codex: skills`
   under `agents`. `.claude/skills/` and `.agents/skills/` each hold
   `grove-explore`, `grove-propose`, `grove-apply` and `grove-archive`, and
   every `SKILL.md` has the same body after its name. Then
   `git add -A && git commit -m "Set up Grove"`.
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
   T-01's status and evidence. `python3 -m unittest` passes. With every task
   done, the same session runs `grove instructions archive 261015-k3m9q`,
   adds the requirement to a new `grove/specs/record-editing/spec.md`, with
   `grove/specs/record-editing/.grove.yaml` beside it, runs `grove check`,
   then `grove archive 261015-k3m9q`:
   > Archived 261015-k3m9q to grove/changes/archive/261015-k3m9q-allow-editing-older-records/; nothing is committed.

6. **Merge** (§6). Commit, `git switch main && git merge --no-ff older-edits-impl`. `grove specs`:
   > record-editing  Record editing  grove/specs/record-editing/spec.md

   The spec says "When a user edits an existing record and leaves its
   category unchanged, the system shall save the edit, even if that category
   is no longer allowed for new records." `grove changes` lists nothing. `grove`
   opens the reader on the spec.

7. **Next** (§Next). A fresh `/grove-explore Editing a missing record crashes`
   cites that requirement as agreed and lists no approved changes. Its
   instructions are under 12 KiB.

## Coordinating changes

It starts where the main line's step 6 ends: one archived change and the
`record-editing` spec. Each step names the section of
[Coordinating changes](../docs/coordinating.md) it tests.

1. **A change that waits** ("A change that waits"). From main,
   `git switch -c unknown-id`, and in a fresh session
   `/grove-propose Editing an unknown ID prints an error`. It finds no other
   open change, says so, and leaves `depends_on` out. Commit, then
   `git switch main && git merge --no-ff unknown-id`.
   `git switch -c delete`, and in a fresh session
   `/grove-propose Add ledger delete ID, from the README's planned list`. The
   session asks whether it needs another open change first and names the
   unknown-ID change, since `delete` would use `find`. I say yes.
   `.grove.yaml` has `depends_on: [<unknown-id>]` and `base:` with main's
   full commit. Commit and merge. `grove changes`:
   > \<unknown-id>  Editing an unknown ID prints an error  approved

   > \<delete-id>  Add ledger delete ID  approved  waits on \<unknown-id>

2. **Apply waits** ("A change that waits"). `git switch -c delete-impl`, a
   fresh session, `/grove-apply <delete-id>`. Grove prints
   > grove: \<delete-id> waits on \<unknown-id>, which this working tree has not archived as finished

   The session tells me and builds nothing. `git switch main && git branch -D
   delete-impl && git switch -c unknown-id-impl`, a fresh session,
   `/grove-apply <unknown-id>`, which builds and archives it. Commit and
   merge. `grove changes` lists the delete change without `waits on`, and
   `/grove-apply <delete-id>`, on a new branch `delete-impl`, builds it.
   Leave that branch unmerged for now.

   **Applying every approved change** ("Applying every approved change").
   Before the merge above, from main with both changes approved,
   `git switch -c apply-all`, a fresh session, `/grove-apply`. Its
   instructions list the unknown-ID change first and the delete change
   second. The session runs `grove instructions apply <unknown-id>`, builds
   and archives it, then `grove instructions apply <delete-id>`, which Grove
   now prints since the dependency is archived on this branch, and builds
   and archives that. `grove instructions apply` on the branch then lists
   nothing to build. One commit, one merge; `grove changes` lists neither.
   `git switch main && git branch -D apply-all` to run the steps above
   instead.

3. **A change that waits on nothing** ("A change that waits"). From main,
   `git switch -c budgets`, a fresh session, `/grove-propose Monthly
   budgets in report`. It needs no other change first, and its
   `.grove.yaml` has no `depends_on`. Commit and merge: two changes are
   approved, delete and budgets, the unknown-ID change being archived.
   `grove changes` lists both without `waits on`.

4. **When the plan is stale** ("When the plan is stale"). From main,
   `git switch -c recategorize`, a fresh session, `/grove-propose An older
   record's category can be changed, only to a current category`. When it
   offers a new requirement, I ask it to fold the rule into the main line's
   requirement instead, so its change spec MODIFIES that requirement.
   Commit and merge. Then
   `git switch main && git switch -c show-saved`, a fresh session,
   `/grove-propose Editing a record prints what was saved`, which MODIFIES
   the same requirement from the same text. Commit and merge. Apply
   `recategorize` on `recategorize-impl`; the `grove check` its archive runs
   after reconciling the spec reports show-saved, not recategorize, and
   exits non-zero:
   > grove/changes/\<show-saved>/specs/record-editing/spec.md:\<line>: requirement "\<name>" was rewritten since base \<sha7>; revise the change with propose

   The session archives recategorize anyway, since the report is not about
   it. Commit and merge; `grove check` on main reports the same line. Set
   show-saved's `depends_on` to `[<recategorize-id>]` in the working tree:
   `grove check` is clean, since show-saved then builds on the change that
   rewrote the requirement. Discard that edit.
   `git switch -c show-saved-revise`, a fresh session,
   `/grove-propose <show-saved>`, asked to bring it up to date with main.
   The session adds recategorize's scenario to show-saved's MODIFIED text
   and sets `base` to main's commit. `grove check` is clean. Commit and
   merge.

5. **The project as data** ("The project as data"). `grove changes --json`
   shows the unknown-ID and recategorize changes as `archived` with `T-01` `done`, and the delete change with no
   `waits_on`, since the unknown-ID change is archived.

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
  > grove/changes/<change-id>-<slug>/proposal.md: missing; it is the change's main document in schema default

  `grove instructions propose <change-id>` still finds the change, so a
  session can write the proposal again.
- **Your own tasks template** (reference, "Change schemas"). On a new
  branch, copy the tasks template from `grove instructions propose "x"` (the
  `<file>` whose path ends in `tasks.md`) to
  `grove/schemas/default/documents/tasks/template.md` with `{{change}}` where it
  shows the new change's ID, and add a line
  `Rollback: [how to undo this task]` to its task. `grove instructions
  propose "x"` now gives that template. `/grove-propose Editing a missing
  record should print an error` writes `tasks.md` with a `Rollback:` line in
  each task, and `grove check` is clean.
- **A bug** (reference, "Change schemas"). On a new branch, write
  `grove/schemas/bugfix/` from the reference's example. In a fresh
  session, `/grove-propose Editing a missing record crashes. It's a bug.`
  The session sees `bugfix` under "Change schemas", confirms with me that
  this is a bug, runs `grove instructions propose --schema bugfix "..."`,
  and follows that: it reproduces the `StopIteration` traceback of
  `ledger edit 99 --note x`, and finds no maintained requirement for a
  missing record, so it asks whether the corrected behavior should be
  agreed. I say yes. It writes `proposal.md` with Observed, Expected,
  Reproduction and The requirement it breaks, a change spec adding "editing
  a missing record prints an error" to `record-editing`, `tasks.md` whose
  T-01 starts with a test that reproduces the crash, and `.grove.yaml` with
  `schema: bugfix`. Commit; `grove review` labels the schema's files
  `change schema`, not `maintained spec`. Merge it, then on another branch
  `/grove-apply <change-id>` builds it (its instructions name schema
  `bugfix`) and `/grove-archive <change-id>` adds the requirement to
  `grove/specs/record-editing/spec.md` and archives the change.
- **A mistake in a schema** (reference, "Change schemas"). After "A bug",
  rename `grove/schemas/bugfix/documents/proposal/template.md` to
  `documents/propsal/template.md` and add `instrution: x` to its `schema.yaml`.
  `grove check` reports both and exits non-zero:
  > grove/schemas/bugfix/schema.yaml:2: unknown field "instrution"; schema.yaml takes description, documents and skills

  > grove/schemas/bugfix/documents/propsal/template.md: schema bugfix does not use this file

  `grove changes` still works.
- **A mistake in grove.yaml** (reference, "grove.yaml"). Add `derp: [x]`
  under `rules`, `propsal: [x]` under `rules: documents:`, and a top-level
  `rule: {}` in `grove.yaml`. `grove check` reports all three and exits
  non-zero:
  > grove.yaml:<line>: unknown rules key "derp"; rules takes documents and skills

  > grove.yaml:<line>: unknown rules.documents key "propsal"; kinds are proposal, change-spec, tasks, implementation, spec

  > grove.yaml:<line>: unknown field "rule"; grove.yaml takes root, target, agents, schema and rules

  `grove changes` still works.
- **Agents from the repository** (reference, "Harnesses"). On a fresh
  fixture, `mkdir -p .claude .opencode .agents/skills/notes`, then
  `grove init`:
  > Wrote grove.yaml and the Grove skills for claude and opencode; nothing is committed.

  `grove.yaml` lists `claude: skills` and `opencode: skills`;
  `.claude/skills/` and `.agents/skills/` hold the four skills, the same
  file in each.
- **Commands** (reference, "Harnesses"). Continuing, add an `AGENTS.md`,
  set `claude: commands` in `grove.yaml`, and run `grove init`. Later runs
  read only `grove.yaml`, so codex stays out:
  > Kept grove.yaml and wrote the Grove skills for opencode and commands for claude, and removed grove-apply, grove-archive, grove-explore and grove-propose from .claude/skills; nothing is committed.

  `.claude/commands/grove/` holds `explore.md`, `propose.md`, `apply.md` and
  `archive.md`. In Claude Code, `/grove:propose Editing a missing record
  should print an error` runs `grove instructions propose "Editing a
  missing record should print an error"`.
- **Drop a harness** (reference, "Harnesses"). Continuing, delete the
  `opencode` line and run `grove init`:
  > Kept grove.yaml and wrote the Grove commands for claude, and removed grove-apply, grove-archive, grove-explore and grove-propose from .agents/skills; nothing is committed.

  A skill of my own, `.agents/skills/notes/`, is still there.
- **A harness Grove does not have** (reference, "grove.yaml"). Set
  `agents` to `codex: commands` and `zed: skills`. `grove check` reports
  both and exits non-zero:
  > grove.yaml:\<line>: harness "codex" has no form "commands"; it has skills

  > grove.yaml:\<line>: unknown harness "zed"; agents takes claude, codex, cursor, hermes, opencode and pi

  `grove init` stops at the first, writing nothing.
- **A dependency that can never be met** (coordinating, "A change that
  waits"). After coordinating step 2, set the delete change's `depends_on`
  to `[260101-zzzzz]`, and `grove check` reports it and exits non-zero:
  > grove/changes/\<delete-id>-\<slug>/.grove.yaml:\<line>: depends_on: 260101-zzzzz is no change, open or archived

  Then propose another change on a branch with `depends_on: [<delete-id>]`,
  and set the delete change's to `[<that change>]`. `grove check` reports
  each with its `path:line`:
  > grove/changes/\<delete-id>-\<slug>/.grove.yaml:\<line>: depends_on: \<other-id> waits on this change in turn, so neither can be applied

  `grove changes` still works. Discard the edits.
- **A change left out of a batch** (coordinating, "Applying every approved
  change"). With an approved change whose `depends_on` names
  `260101-zzzzz`, `grove instructions apply` ends its list with
  > Leave out \<id>, "\<title>": it waits on 260101-zzzzz, which is neither being built nor archived as finished. Grove's `check` says why.
- **Apply unattended** (coordinating, "Applying every approved change").
  From main with two changes approved, `claude -p "/grove-apply"` with
  `apply: [Commit after each task, naming its task ID.]` under
  `rules: skills:` in `grove.yaml`. It creates a branch, builds and
  archives both, commits, asks nothing, and its last message says which
  changes were archived.
- **Commit as you go** (reference, `rules`). With
  `apply: [Commit after each task, naming its task ID.]` under
  `rules: skills:` in `grove.yaml`, `/grove-apply` commits T-01 with its ID in the
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
- The wording of phase 3's `check` messages and the JSON's `root`-relative
  paths. Owner's lean, from Remote Control, 2026-10-20: "the wording is
  fine for now".

## Assumptions to confirm

Agent inference, not the owner's decisions (`CLAUDE.md`).

- `grove.yaml` fields and defaults: `root: grove`, `target: main`,
  `agents` (left out or empty, `claude` and `codex` as skills), `schema: default` and `rules`, with `documents` keyed by
  `proposal`, `change-spec`, `tasks`, `implementation` and `spec`, and
  `skills` keyed by a skill.
- `init` never overwrites `grove.yaml` and always rewrites the skills.
  Its second run says "Kept grove.yaml and wrote the Grove skills for claude
  and codex; nothing is committed." A field left out of `grove.yaml` takes
  its default. A harness Grove has no row for, or a form it does not
  have, stops `init` before it writes anything, with `check`'s words.
- `agents` takes only the map of harness to form. The list form,
  `agents: [claude, codex]`, which no release wrote, is a YAML error, not a
  `check` problem; v0.1.0 is the first release, so the change is not marked
  breaking.
- `init` removes Grove's files it no longer writes before it writes the
  rest, so a skill folder linked to another, such as `.claude/skills` to
  `.agents/skills`, keeps its skills.
- `init`'s message groups harnesses by form, skills first ("the Grove
  skills for opencode and commands for claude"), and names what it removed
  by folder ("grove-apply, … from .claude/skills"). It leaves a folder it
  emptied, such as `.claude/commands/grove/`, in place; Git ignores it.
- Detection looks only at the top of the repository, and a `CLAUDE.md` or
  `AGENTS.md` there counts even when it is a link to the other.
- A skill runs `grove` from the PATH of the agent's shell. Every
  `SKILL.md` has one body: it names the user's words `<args>` and says
  that is what the user typed after the skill's name, or nothing (Claude
  Code appends them as `ARGUMENTS:`). A command passes them as
  `$ARGUMENTS`, and its file has only a `description` above its body.
- OpenCode's skills are not typed: Grove writes no `opencode/slash`
  metadata, since the skill files are shared with other harnesses; typing
  is what OpenCode's commands are for.
- Pi's one-time step is trusting the project when Pi asks, or `/trust`.
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
- The fixture's "Planned for 0.4" leaves out the direction note's
  `imports/2024.csv` and its import item; no step needs them.
- The project index lists each approved change's capabilities (the
  `<capability>` folders of its files) and what it waits on, rather than
  dependencies Grove works out as candidates. Propose names the likely ones
  from it.
- `grove changes` writes `waits on <id>, <id>` after `approved`, and
  nothing for a change whose dependencies are all archived as finished.
- A change with no dependency leaves `depends_on` out; `depends_on: []` is
  reported like any other empty field.
- `base` comes from `{{commit}}`, a placeholder only Grove's `.grove.yaml`
  template fills; `{{base}}` stays the short commit in schema templates.
- The docs page for coordinating changes is `docs/coordinating.md`.
- `check` and `instructions apply` look for what `depends_on` names among
  the working tree's changes, open and archived, and not on the target; a
  branch holds the target's archive as of where it started. `changes` reads
  the target.
- The stale check compares a requirement's text below its heading, trimmed,
  by name: the working tree's maintained spec against what it was at
  `base`, what the change itself gives it, and what each change it depends
  on, directly or through others, gives it. One missing from the working
  tree was removed. A requirement the maintained spec did not have at
  `base` is not compared. A `base` that is not a commit is reported.
- `instructions apply` with no ID orders the target's approved changes that
  the working tree has not archived: each after its dependencies, otherwise
  by `created`. It prints the list and the loop; each change's own
  instructions come from `instructions apply <id>`, run by the session.
- Apply ends by running the archive instructions; `/grove-archive` stays
  for dropping a change and for archiving by hand.
- `grove changes --json` gives paths under the root, with `root` at the
  top, rather than from the repository's top; `state` is worked out when
  the project is read; a field with no value is left out; a change's
  capabilities and problems are not in it.
- `grove version` prints one line, `grove <version>`, without a leading
  `v`, as `grove 0.1.0`; a release build sets it, `go install` of a tag
  takes the tag, and a build with no version recorded says `dev`.
- `install.sh` takes a tag as the release names it (`v0.1.0`), not a bare
  version, and replaces a `grove` already in the folder without asking. It
  says "Installed grove in <dir>." and, when that folder is not on PATH,
  "<dir> is not on your PATH; add it to run grove."; each failure is one
  line starting "grove: ".
