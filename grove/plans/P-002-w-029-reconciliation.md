---
id: "P-002"
type: plan
title: "W-029 reconciliation: mapping, conversion, reference repair, schema 1/2 removal"
status: current
created: "2026-09-21T21:00:03Z"
updated: "2026-09-21T21:00:03Z"
work: ["W-029"]
---

## Design

Branch `worktree-W-029` from main `70de539`, which contains W-030. Every local
branch is an ancestor of main, so no unmerged branch holds old-layout work;
reintegration is rehearsed with a synthetic branch.

**Inventory at `70de539`.** 46 records with this plan (26 work, 6 decisions,
1 question, 9 terms, 2 plans, 2 reviews), 13 legacy plans in `docs/plans/`,
9 legacy reviews in `docs/reviews/`, the brief, and 3 spent handoff prompts in
`docs/prompts/`. The prompts are `.txt` evidence that W-010 links and its
dogfood review already calls no longer needed: they stay where they are as
historical literals, listed in the mapping page as kept, not converted.
README, AGENTS.md, the adapters and the four `docs/*.md` guides keep their
homes.

**Order.** One sequence over all 68 sources by document date, oldest first:
a record's `created`, a legacy document's first-commit author date in UTC.
Ties break on the old ID or path, bytewise. The order conveys nothing. The
mapping page is created last with `new page`.

**Legacy documents.** `convert PATH --type plan|review --title H1`, slug from
the old filename without its date or ID prefix. `work` is the work IDs in the
filename; a review without one takes the work records that link it
(`integrated-cli`: W-003 to W-008; `predecessor-work` and
`shaping-and-runner-evidence`: W-010). `direction-evaluation` names no work
and gets `relates_to` D-004, which links it. Status stays the type's first,
`current`: judging a plan superseded is not mechanical. No `examined` and no
dates are set: those documents never declared one examined commit, and their
commit tables stay in the body. R-001 and R-002 keep theirs through `convert`.
The originals and the empty folders are removed in the same commit.

**Script.** One Python script, `scripts/w029-migrate.py`, committed with the
migration and deleted with schema 2's support, which ends its usefulness. It
(1) edits `schema_version` to 3 and checks, (2) runs `convert` in order and
collects the JSON lines, (3) sets `work`/`relates_to` on converted documents,
(4) moves the brief to `grove/brief.md` and edits `brief:`, (5) repairs
references, (6) writes the mapping page body, (7) runs `check` and the audit.

**Reference repair**, in this order, over every tracked `.md` file (the prompts are left alone):

1. Markdown link destinations are resolved against the linking file's old
   path, mapped if the target moved, and made relative to the linking file's
   new path, fragment kept. This also repairs unmoved targets linked from
   moved files.
2. Root-relative old paths in prose and inline code become new paths.
3. Bare typed IDs that are in the mapping become neutral IDs, except inside a
   `worktree-…` branch name, a `formerly` value, the mapping page, a URL, and
   an explicit exception list. The script prints every occurrence it skips and
   every typed ID with no mapping (fixture IDs such as `Q-002`). Before
   publication I read each done record's and review's rewritten evidence lines
   for fixture IDs that coincide with real ones (a test's own `W-001`) and add
   them to the exception list. Exceptions are recorded in the mapping page.
4. Directory conventions in current instructions (`grove/work/`,
   `docs/plans/`, `worktree-W-012`, typed `new` wording) are hand edits to
   README, AGENTS.md, the guides and adapters, not script output.

Go sources change in the schema-removal phase, not here; a Go comment that
names a migrated record gets its new ID there.

**Rehearsal and recovery.** The script runs first in a disposable clone by
absolute `--project` path; its mapping is the prepared mapping, and the real
run must reproduce it exactly because both counters start unset. If the real
run fails: `git reset --hard` and `git clean` in this worktree only, then
delete the common directory's `grove/neutral-ids`, which is safe exactly while
no `G-` record exists in any ref or worktree (the floor scan would restore a
higher number otherwise). A rerun over converted sources is refused by
`convert` and reserves nothing. Rollback after the commit is `git revert` of
the migration commit; consumed numbers stay consumed.

**Reintegration.** In the clone, merge a branch that edits an old-layout
record and adds a typed-ID one. Expected: Git reports the edit as a
rename/modify it carries to the new path or a conflict, and `check` refuses
any restored typed-ID file, by `formerly` before the schema removal and by the
ID grammar after it. Lineage limit to document: `history` follows the rename;
an old commit is read with that commit's CLI.

**Schema 1/2 removal.** `schema_version` must be 3; the number is not
renumbered. `TypeInfo` loses `Prefix`, `Folder` and `Schema`; an ID is
`G-NNN` only; `next-ids` handling, folder rules, schema-gated wording and the
board's `W-` fallback go. `convert ID` goes with typed IDs, its only input;
`convert PATH`, `formerly` and its duplicate check stay. Test fixtures move to
neutral IDs and flat paths; old/new-checkout tests are deleted.
`docs/record-model.md` becomes one schema's description. `versions` will
report a schema-2 branch as unreadable: the merged local branches are the
owner's to delete, and the handoff says so.

## Steps

1. [ ] Commit this plan; set W-029 active.
2. [ ] Write the script; rehearse in a disposable clone; read the skip and
   unmapped reports; settle the exception list; rehearse reintegration and
   recovery.
3. [ ] Run it here; hand-edit current instructions; `check`, link audit,
   obsolete-ID audit; exercise context, show, board, history on migrated
   proposed and done work. Commit. Independent review of the migration.
4. [ ] Remove schemas 1 and 2 and the script; fix tests and docs; full
   verification. Independent review of the combined revision.
5. [ ] Reconcile W-029, the brief and the documented limits; hand off for the owner's
   judgment of the tree and board (acceptance 6).
