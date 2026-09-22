---
id: "G-091"
type: plan
title: "G-041 nullsec cutover: mapping, conversion, uninstall"
status: current
created: "2026-09-22T19:39:04Z"
updated: "2026-09-22T19:40:45Z"
work: ["G-041"]
---

## Design

**Bases.** This repository: `worktree-G-041` in `.claude/worktrees/G-041`
from `main` `77b4111`. Nullsec: branch `worktree-grove-cutover` in a linked
worktree `.claude/worktrees/grove-cutover` from nullsec `main` `366b063`, the
commit G-041 inspected, which is still nullsec's `main`. Nullsec has no record
for this work, so the branch is named for the change.

**Inventory corrections at `366b063`.** G-041's inventory holds, with four
additions: no predecessor record has a `title` field; the three `research`
files have no `id`; `history/README.md` and four evidence documents have no
frontmatter; and `docs/plans/evidence/` holds two Markdown files
(`W-035-baseline.md`, `W-036-fitting-comparison.md`), not one. The whole
predecessor tree entered nullsec in one commit on 2026-09-14 (imported from
Bench), so Git dates say nothing about when a record was written.

**Numbering order.** One sequence, old ID order within groups: decisions
`D-0001` to `D-0043` (so `D-00NN` becomes `G-0NN`, a convenience, not a
rule), work `W-001` to `W-040` then `W-901` to `W-903` (offset 43: `W-001` is
`G-044`), questions, terms, capabilities, the seven evidence documents, the
twelve legacy plans, the two Markdown plan evidence files (each group by old
path, bytewise), and the map page last. Date order, which G-052 used, has no
data here: dates are day-only or absent, and Git has only the import date.

**Conversion.** `grove init` first (`grove.yaml`, `grove/`, the placeholder
brief, six entrypoints). Then `convert PATH --type T --title TITLE --slug S`
per source; the originals are deleted after conversion. Then the old
frontmatter, which `convert` keeps as body bytes, is removed from each body.
Nothing else in a body changes except the reference repair below.

| Predecessor | Becomes |
| --- | --- |
| `work`, `decision`, `question`, `term` | the same type |
| `capability`, `research`, frontmatter-less evidence, plan evidence `.md` | `page`; plan evidence relates to its work, since a page has no `work` |
| `docs/plans/*.md` | `plan`, `work` from the IDs in the filename (`W-029-W-031` is the range G-072 to G-074, set after review) |
| `history/README.md` | removed; its Bench provenance is quoted in the map page |
| `brief.md` | `grove/brief.md` (a move over the placeholder), frontmatter removed |
| `deliveries/*.json`, `.gitkeep`, `grove.toml` | removed, listed in the map page |
| `docs/plans/**` non-Markdown evidence | moved to `docs/evidence/`, old ID prefix in the name rewritten |

Titles: a document's `# ` heading where it has one (the four evidence
documents and the plans); otherwise the filename without the ID prefix, with
hyphens as spaces, the first letter capitalised except for terms, whose title
is the term, and a short acronym table (`npc` NPC, `eve` EVE, `ui` UI, `e2e`
E2E, `v1` V1 and so on). Slugs: the old filename without its ID prefix; a plan
takes its first work's slug plus `-plan`. The rehearsal prints every title.

| Field | Mapping |
| --- | --- |
| `status` | kept; `parked` becomes `open` with a first body line saying so; `superseded` decisions keep it (below) |
| work `done` | written directly into the frontmatter, without `candidate`: `update` writes done only with a candidate, and these are historical Done (record model), never backfilled |
| `kind` | kept; `spike` becomes `investigation`; research `agent`/`user` dropped |
| `size` | `small`/`large` kept; `bounded` and `spike` dropped (no equivalent, not synthesized) |
| `priority`, `depends_on`, `members`, `blocks` | kept, IDs mapped |
| `scope`, `applies_to`, `supersedes`, `superseded_by` | merged into `relates_to` in that order, deduplicated |
| `plan` | dropped: the plan record's `work` names the work instead |
| `created`, `updated` (day-only) | dropped; the map page shows the old day. `update` writes `updated` when it sets fields |
| `id`, `type` | replaced by the new record's |
| `sources` | `relates_to` when every entry names a nullsec record (four decisions); otherwise dropped |
| `focus` (brief), `project`, `url`, `tags`, `started`, `unchanged`, `batch`, `batch_reason` | dropped |

Every dropped value is printed per record in the map page, so the page plus
Git history at `formerly` accounts for all of it.

**Superseded decisions: add `superseded` to decision statuses.** The record
model already anticipates it ("Supersession can be added when an actual
replacement needs it"), and nullsec has two actual replacements. Both other
mappings misstate a decision's standing: `accepted` says it is in force, and
`rejected` says it was never adopted. The change is one status in the type
table, the record model's two decision lines, and a test; the replacement
link is `relates_to` plus the decision's own body, as for plans, so no new
field. The installed binary must include it, so it is built from this branch.

**Reference repair**, script output over every tracked text file of nullsec
(binary files skipped):

1. Wikilinks in the converted Markdown: `[[OLD-STEM]]` becomes
   `[G-0NN](G-0NN-slug.md)`, `[[slug]]` becomes `[slug](G-0NN-slug.md)` so
   prose such as "an [[intent]]s" still reads, and `[[target|label]]` keeps
   its label. Wikilinks into the Bench vault (`raw/…`, `projects/…`,
   `docs/plan/…`) are kept as written and listed. Wikilink look-alikes in code
   (`[[package]]`, array literals) are never touched: only converted records
   are rewritten this way.
2. Markdown link destinations resolved against the linking file's old path,
   mapped if the target moved, made relative to the new path, fragment kept.
3. Old paths in prose and code become new paths, including the pre-`history/`
   paths that closed records had before the predecessor moved them.
4. Bare old IDs (`W-NNN`, `D-NNNN`) become new IDs, except in branch and
   worktree names, `formerly`, the map page, the Bench wikilinks, and an
   explicit exception list of verbatim quotations (the builder's words, commit
   subjects). Lowercase `wNNN` capture and intermediate names in `art/` are not
   IDs to the script and stay; the map page explains them.
5. Hand edits: the `grove:begin` blocks in `CLAUDE.md` and `AGENTS.md` become
   one short `## Grove` section naming `grove` on `PATH`, `grove brief`,
   `grove context`, and the two entrypoints. Their other drift stays.
   `.claude/settings.json`, which only enables the predecessor plugin, is
   deleted. The brief's Next moves to the end of the owning work's (old
   `W-032`) body under `## Next`, with a line saying where it came from.

**Script.** One Python script, `scripts/g041-nullsec-cutover.py` in this
repository's branch, not nullsec's and not product code; deleted before the
candidate, so it stays at a named commit. It takes the nullsec checkout and
the binary, runs `init`, converts in order, strips old frontmatter, sets
fields with `update` (status `done` written directly), moves the brief and the
evidence, repairs references, writes the map page with `new page`, runs
`check`, and audits: every Markdown link in the repository resolves, and each
remaining old-ID mention is on the exception list. It writes the mapping as
JSON outside the checkout for comparison.

**Rehearsal and recovery.** In a clone (`git clone --no-hardlinks`, its own
common directory, so its counter is not nullsec's) under the session
scratchpad, reached by absolute `--project`. An independent reader examines the
rewritten ID lines for verbatim quotations and for strings code compares, and
the exceptions are settled there. Recovery: `git reset --hard`, `git clean
-fd`, delete the clone's `.git/grove/neutral-ids`, rerun; the mapping JSON
must be byte-identical. A live failure recovers the same way in the nullsec
worktree only, which is safe while no `G-` record exists in any nullsec ref.

**Live.** The same script and binary in the nullsec worktree; the mapping
must equal the rehearsal's. `npm test` there, then read the diff of every
non-Markdown file. Commits: `init`, the conversion and repair, the hand
edits.

**Machine** (acceptance 3). Build from a clone of a named commit of this
branch (a linked worktree stamps the enclosing checkout) into
`~/.local/bin/grove` after `uv tool uninstall grove` frees the path; uninstall
`grove@mascah` in every local scope, remove the `mascah` marketplace and
nullsec's `settings.local.json` entry; remove `grove@grove-local` and the
`grove-local` marketplace from `~/.codex/config.toml`. `../skills` untouched.
The installed build is replaced from `main` after the merge; that rebuild is
the integrator's step.

**Fresh sessions** (acceptance 4) run in the nullsec worktree, since nullsec
`main` has no `grove.yaml` until the owner merges: headless `claude -p` and
`codex exec` asked to report `grove version`, the brief's first heading and the
record count through the CLI, and which `grove-*` skills they can see, changing
nothing. Observed output is recorded apart from file existence.

**Adjustments from evidence.** The rehearsal and the reader added: block
YAML lists in the parser; `pre-W-036` rewritten; `W-019-T3` kept as a branch;
`W-021/023/024` shorthands expanded; paths relative to the old root
(`history/evidence/…`) rewritten; link labels that spelled the old filename
renamed; the brief-Next lead-in reworded. The final review added G-121's range,
a lead-in saying G-075's own Next supersedes the moved text, and three live
instructions (`art/design/section-study-recipe.md`, G-059's Next,
`docs/world/expedition-foundations.md`) moved to new CLI commands, which the
plan had limited to the two instruction files.

## Steps

1. [x] Commit this plan (`499f11e`); set G-041 active (`2c349bc`).
2. [x] `superseded` decision status: type table, record model, test (`1452c6f`).
3. [x] Script (`5ef2978`); rehearsal in a disposable clone; recovery by reset,
   clean and deleting `neutral-ids` reproduced a byte-identical mapping and an
   identical tree apart from timestamps; an independent reader read all
   rewritten lines; exceptions settled.
4. [x] Live on `worktree-grove-cutover` (`5381966`, `473529b`, `ecbb036`, review
   fixes `3eb2785`); mapping and tree identical to the rehearsal; `npm test`
   passes.
5. [x] Machine: binary from `5ef2978` installed; predecessor removed; `grove
   version` from a login shell, Claude's Bash tool and Codex's shell.
6. [x] Fresh Claude and Codex sessions in the nullsec worktree.
7. [x] This repository's `AGENTS.md` and README (`1dc3922` and after review);
   rollback in G-041.
8. [x] Independent review (G-092); evidence and handoff into Review.
