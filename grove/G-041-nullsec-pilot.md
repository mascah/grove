---
id: "G-041"
type: work
title: "Cut nullsec over to this Grove and uninstall the predecessor"
status: proposed
created: "2026-09-21T00:54:15Z"
updated: "2026-09-22T19:24:43Z"
kind: tooling
size: large
priority: 2
depends_on: ["G-040"]
relates_to: ["G-035", "G-036", "G-064", "G-065"]
formerly: "W-023"
---

## Outcome

Nullsec runs on this Grove, and the predecessor no longer answers anywhere on
the owner's machine. The owner set this outcome in a shaping session on
2026-09-22, in four points: nullsec is bootstrapped with the new CLI; its
records are converted to the same shape G-052 gave this repository; nullsec
development uses the new CLI from then on; and the old CLI is off `PATH` and
interferes with nothing.

Also decided in that session: the first real nullsec change through shape →
implement → review → integrate, and the owner's continue/revise verdict,
belong to [G-036](G-036-interactive-adoption.md), not to this record. They
happen after this cutover, as ordinary nullsec work shaped in nullsec.

## Scope and constraints

**Write scope, granted by assigning this record.** A nullsec branch holding
the migration, the owner's user-level tool and harness configuration listed
under Acceptance 3, and this repository's instructions. `../skills` is
out of scope and stays untouched (owner, 2026-09-22): retiring or archiving it
is a later, separate decision. Merging nullsec and pushing either repository
are the owner's steps unless the assignment says otherwise.

**Observed at nullsec `main` `366b063`, 2026-09-22** (read-only inspection
from this repository):

- Clean `main`; no `active` work. Every local branch is an ancestor of `main`
  except `worktree-W-002-W-003-W-012`, which is 6 commits ahead, dated 2026-09-16,
  and whose three items are recorded as done. The worktrees
  `.claude/worktrees/W-038` (41 untracked captures) and
  `/private/tmp/nullsec-station-foundations` belong to merged branches.
  Branches without `grove.yaml` are absent from `versions` and the board
  (`internal/versions/tree.go`), so old branches add no noise. Deleting them
  is the owner's call.
- Predecessor content: `grove.toml` (its schema 2, root `docs/grove`) and in
  `docs/grove/`:
  - the brief, whose frontmatter includes `focus` and a Next section
  - 43 work records (36 done, 7 proposed; 36 under `history/work/`)
  - 43 decisions (2 `superseded`, under `history/decisions/`)
  - 6 questions with slug IDs (2 `parked`)
  - 6 terms and 7 capabilities
  - 3 `research` evidence files and 4 other evidence documents under
    `history/evidence/`
  - `history/README.md`, whose Bench provenance must be kept
  - 6 `deliveries/*.json` and `.gitkeep` files

  Records link to each other with 202 `[[wikilinks]]`. `docs/plans/` holds
  12 legacy plans plus non-Markdown evidence (`.json`, `.csv`, `.rs`) and one
  Markdown evidence file.
- Predecessor values with no direct equivalent in the record model:
  - work size `bounded` and `spike`, and kind `spike`
  - kind values `agent` and `user`
  - decision status `superseded`, with `superseded_by` and `supersedes`
  - question status `parked`
  - types `capability` and `research`
  - fields `scope` and `applies_to` (both point at capabilities), `plan`,
    `batch`, `batch_reason`, `sources`, `tags`, `url`, `unchanged` and
    `started`
  - date-only `created` and `updated` values

  The record model says decision supersession "can be added when an actual
  replacement needs it". Nullsec has two actual replacements.
- Outside `docs/grove`, 58 distinct old IDs, all naming real records, appear
  about 1,130 times across roughly 150 tracked files: `client/src`, `crates`,
  `art`, `e2e`, `docs/world`, `docs/plans`, README, RETROSPECTIVE and the
  instruction files. About 70 of these are in strings rather than comments:
  vitest titles, Rust assert messages, fitting-comparison headers, a
  `console.log`, and `e2e/scripts/write-retrospective.mjs`, which regenerates
  `docs/RETROSPECTIVE.md`. Separately, 20 `docs/plans/…` path mentions sit in
  11 code files. 294 capture paths such as `art/captures/w039-crew-study/`
  carry old IDs in their names.
- Predecessor hooks on this machine:
  - `~/.local/bin/grove` is a uv tool installed as editable from
    `../skills/cli`.
  - The Claude plugin `grove@mascah` is installed at local scope for
    `../skills`, nullsec, four nullsec worktrees and this repository, where it
    is disabled.
  - nullsec's `.claude/settings.local.json` enables `grove@mascah`, and its
    `.claude/settings.json` enables a `grove@grove-local` that no marketplace
    provides.
  - `~/.codex/config.toml` enables the Codex `grove@grove-local` plugin from
    the `grove-local` marketplace (`../skills`). G-040 saw it load first.
  - Both marketplaces contain only `grove`.
  - nullsec's `CLAUDE.md` and `AGENTS.md` both carry the predecessor's
    `grove:begin` block, and the two files have drifted apart (9 against 17
    smoke tests).
  - No other sibling has a `grove.toml`; Keyborg uses Bench.

**Selected by the owner, 2026-09-22:**

- The target matches this repository: `grove.yaml` from `grove init`,
  `schema_version: 3`, root `grove/`, brief `grove/brief.md`, neutral `G-`
  IDs flat under the root, the `grove-work` and `grove-shape` entrypoints, and
  a migration map page like [G-069](G-069-migration-map.md).
- Old IDs are rewritten everywhere in tracked text, code, tests and the
  retrospective generator included, as G-052 did here. The owner chose this
  after comparing its cost with a records-and-docs-only rewrite. What cannot
  change stays and is resolved by the map page: commit messages, branch
  names and art capture file and folder names. Exceptions (verbatim quotes,
  anything that is not a nullsec record) are listed in the map page.
- Uninstall only: remove every predecessor hook listed above from this
  machine. Leave `../skills` untouched.

**Proposed design, for preparation to confirm in a plan record:**

- Follow the G-052 method ([G-068](G-068-reconciliation-plan.md)): a one-off
  script that is not committed to either repository's product code. It runs
  `convert` over each predecessor file (the new root is outside `docs/grove`,
  so `convert` accepts them and writes `formerly`), strips the old frontmatter
  from the body, sets fields with `update`, repairs wikilinks, Markdown links,
  old paths and bare IDs from the collected mapping, then writes the map page.
- Capabilities and research/evidence documents become pages. Legacy plans
  become plan records with `work`. `scope`, `applies_to` and supersession
  links become `relates_to`. `parked` becomes `open` with a body note. Kind
  `spike` becomes `investigation`. Values with no equivalent are dropped
  rather than synthesized, and the map page records them. Deliveries,
  `.gitkeep` files and `grove.toml` are removed, and the map page lists them.
  Non-Markdown evidence stays an ordinary file, kept in place or moved, as
  long as no parallel layout of the records remains; the map page says which.
  The brief loses its frontmatter and Next, whose content moves to the owning
  work's Next, as this repository's brief did.
- For the two superseded decisions, preparation either adds a `superseded`
  decision status (a CLI and record-model change that the record model
  already anticipates) or chooses another mapping. The plan says which, and why.
- The `grove:begin` blocks in nullsec's `CLAUDE.md` and `AGENTS.md` become a
  short section on the new CLI and its entrypoints. `init` never edits those
  files. Any other change to them is outside this record.
- Build the binary from a named, pushed commit of this repository, stamped,
  into `~/.local/bin/grove` once the uv tool frees that path. Rebuilding stays
  manual, and `grove version` names what is installed.
- Sequence: rehearse in a disposable clone of nullsec reached by absolute
  `--project` path, never a linked worktree, since `new` and `convert` there
  would advance nullsec's shared counter. Then run live on a nullsec branch,
  where the result must reproduce the rehearsal mapping. Then uninstall the
  predecessor and install the binary, and then run the fresh-session checks.

**Limits to keep visible:** nullsec's `G-` numbers are independent of this
repository's, so this repository's G-041 and nullsec's G-041 will be
different records. Conversations spanning both need to name the repository.
Keyborg is not migrated alongside nullsec.

## Acceptance

1. **Rehearsal.** In a disposable clone, the migration converts or accounts
   for every predecessor file listed above: each is converted, with its old
   ID or path mapped to its new one, or removed on purpose with a reason.
   `grove check` passes, every Markdown link in the repository resolves, and
   an audit shows that each remaining old-ID mention is a listed exception.
   Recovery (reset, clean, delete the clone's `grove/neutral-ids`) reproduces
   the identical mapping.
2. **Live migration.** A nullsec branch reproduces the rehearsal mapping.
   After it, nullsec has `grove.yaml` and the entrypoints from `init`, the
   converted records under `grove/` with the map page, and no `grove.toml`,
   `docs/grove/` or legacy plan Markdown. Its instruction files describe the
   new CLI, and nullsec's `npm test` passes. The diff has been read for the
   string changes.
3. **The predecessor is gone.** The uv tool is uninstalled. `grove@mascah`
   is uninstalled in every local scope, and the `mascah` Claude marketplace
   is removed. Nullsec's plugin entries are removed. The Codex
   `grove@grove-local` plugin and `grove-local` marketplace are removed from
   `~/.codex/config.toml`. `grove version` prints the installed build from a
   login shell, from Claude's Bash tool, and from Codex's shell. `../skills`
   is unchanged.
4. **Fresh sessions.** Fresh Claude and Codex sessions in nullsec find the
   brief and records through the new CLI and load the `grove-shape` and
   `grove-work` entrypoints without the predecessor answering. The report
   keeps observed behaviour separate from file existence.
5. **This repository.** Its `AGENTS.md` (the installed-`grove` and
   `../nullsec` retrieval lines) and the README's coexistence section match
   the new state; `go run ./cmd/grove` stays the way to develop here. The
   rollback is written down: revert nullsec's merge, then
   `uv tool install --editable ../skills/cli` and reinstall the plugins.
6. **The owner's judgment.** The owner judges the converted nullsec tree and
   its board in an actual terminal. No check substitutes for this.

## Next

Assign with `/grove-work G-041`. Preparation writes a plan record with the
field mapping, the superseded-decision choice and the script outline, then
runs the rehearsal. After this is done, G-036's Next takes the first real
nullsec change.
