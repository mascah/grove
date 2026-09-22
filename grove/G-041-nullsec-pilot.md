---
id: "G-041"
type: work
title: "Cut nullsec over to this Grove and uninstall the predecessor"
status: active
created: "2026-09-21T00:54:15Z"
updated: "2026-09-22T19:40:45Z"
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

## Evidence

Plan [G-091](G-091-nullsec-cutover-plan.md), review
[G-092](G-092-nullsec-cutover-review.md). Grove branch `worktree-G-041` from
`main` `77b4111`: plan `499f11e`, active `2c349bc`, `superseded` status
`1452c6f`, the migration script `5ef2978` (removed in `db819bf`; read it with
`git show 5ef2978:scripts/g041-nullsec-cutover.py`), docs `1dc3922`, then the
review fixes and this evidence. Nullsec branch `worktree-grove-cutover` in
`nullsec/.claude/worktrees/grove-cutover` from nullsec `main` `366b063`:
`5381966` `grove init`, `473529b` conversion and rewrite, `ecbb036`
instructions and settings, `3eb2785` review fixes.

Owner choices during the session, 2026-09-22: install from the local branch
commit, since pushing is the owner's step ("Local branch commit"), and
uninstall in this session once the live migration passed ("Yes, this session").
Decided in the plan: a `superseded` decision status (the record model's
anticipated supersession; `accepted` or `rejected` would misstate the two
replaced decisions), numbering by old ID within groups because nullsec's dates
are day-only and Git has only the 2026-09-14 import, and historical `done`
written directly without `candidate`, since `update` writes done only with one
and AGENTS.md forbids backfilling.

Against each acceptance item:

1. **Rehearsal.** In a `git clone --no-hardlinks` of nullsec `366b063` under
   the session scratchpad, reached by absolute path, the script ran with a
   binary built from a clone of `1452c6f`: `check` `OK: 127 records` (126
   converted plus the map page `G-127`), 0 broken links (the same 0 as before
   the migration), and 5 remaining old-ID lines outside branch and worktree names (which the
   map page's branch rule covers), each listed on the map page
   (a verbatim builder quote, two `.grove-run/` artifact names never tracked,
   the `W-019-T3` branch, one folder named as it was). Every predecessor file
   is converted, moved (three non-Markdown evidence files to
   `docs/evidence/`), or removed with a reason on the map page. Recovery
   (`git reset --hard`, `git clean -fd`, deleting the clone's
   `.git/grove/neutral-ids`) reproduced a byte-identical `mapping.json` and an
   identical diff apart from timestamps, twice.
2. **Live migration.** The same script and binary in the nullsec worktree
   produced a byte-identical `mapping.json` and the same 694 files, hashed with
   timestamps normalised, as the rehearsal. Nullsec now has `grove.yaml`, the
   six entrypoints from `init`, 127 records under `grove/` with the map page,
   and no `grove.toml`, `docs/grove/` or `docs/plans/`. Its `CLAUDE.md` and
   `AGENTS.md` carry a `## Grove` section on the new CLI in place of the
   `grove:begin` block. `npm test` at `ecbb036`: exit 0, cargo workspace ok,
   405 vitest, 17 of 17 smoke. Earlier runs failed the smoke test
   "pointer round trip" in two full-suite runs at `ecbb036`; it passed 5 of 5 alone at
   both `366b063` and `ecbb036`, and a full smoke run at `366b063` also failed
   one test in one of two runs, so it is an existing flake, not the migration
   (load averages 3 to 4). The string changes were read line by line by the
   first reviewer (G-092, round 1): test titles, assert messages and comments,
   none compared by code. The review fixes `3eb2785` touch Markdown only;
   `grove check` passes after them.
3. **The predecessor is gone.** `uv tool uninstall grove`; `grove@mascah`
   uninstalled at local scope in `skills`, `nullsec`, this repository and the
   four vanished nullsec worktrees W-035, W-036, W-037 and W-039 (each through
   an empty directory made for the command and removed after); the `mascah`
   marketplace removed; nullsec's `settings.local.json` now enables none; the
   orphaned caches `~/.claude/plugins/cache/{mascah,grove-local}` removed;
   `codex plugin remove grove@grove-local` and `codex plugin marketplace
   remove grove-local` removed both sections from `~/.codex/config.toml`, and
   the empty Codex cache folder was removed. Nullsec's tracked
   `.claude/settings.json`, which only enabled `grove@grove-local`, is deleted
   on the branch. `~/.local/bin/grove` is built from a clone of `5ef2978`
   (same Go code as `1452c6f`) and prints `grove
   v0.0.0-20260922195818-5ef2978c9bf8 (5ef2978c…) guides sha256:99621fb54f05`
   from a login shell (`zsh -lc`, only that `grove` on `PATH`), from Claude's
   Bash tool, and from `codex exec`. `../skills` has no working-tree change.
4. **Fresh sessions**, all in the nullsec worktree, observed from their tool
   calls: a `claude -p` session ran `grove version`, `grove brief`, `grove
   check` (`OK: 127 records`) and `grove list`; `/grove-work --interaction
   headless` and `/grove-shape --interaction headless` ran `grove guide work`
   and `grove guide shape` and returned the guides' no-assignment wait; their
   skill list held `grove-shape` and `grove-work` and no `grove:*` skill. The
   entrypoints are `disable-model-invocation`, so a model asked to list its
   skills does not see them; a typed invocation does. `codex exec -s
   read-only` with `$grove-work` and `$grove-shape` ran `grove guide work` and
   `grove guide shape` and returned the same waits; `codex exec` waits on
   stdin until closed, so it ran with `< /dev/null` after a first attempt
   timed out. No session wrote a file.
5. **This repository.** `AGENTS.md` says the installed `grove` is a build of
   this CLI that can lag the checkout, keeps `go run ./cmd/grove` for
   development here, and reads nullsec through the installed CLI and
   `../skills` as files; the README's adoption section, introduction and
   reset notes match. The rollback is below.
6. **The owner's judgment** of the converted tree and its board is open.

**Rollback.** Revert nullsec's merge of `worktree-grove-cutover` (restoring
`grove.toml`, `docs/grove/` and the tracked `grove@grove-local` setting), then:

```sh
rm ~/.local/bin/grove     # the Go build occupies the path the uv tool links
uv tool install --editable ~/GitHub/mascah/skills/cli
claude plugin marketplace add mascah/skills
(cd ~/GitHub/mascah/skills && claude plugin install grove@mascah --scope local)
(cd ~/GitHub/mascah/nullsec && claude plugin install grove@mascah --scope local)
```

and restore these two sections of `~/.codex/config.toml`, removed verbatim:

```toml
[marketplaces.grove-local]
source_type = "local"
source = "/Users/mascah/GitHub/mascah/skills"

[plugins."grove@grove-local"]
enabled = true
```

then run `codex plugin add grove@grove-local`, since `codex plugin remove`
also emptied its cache. The removed Claude plugin was `grove` 0.2.0 at skills
commit `3877d116`; `marketplace add` installs whatever `mascah/skills` holds
now, so compare the version.

Before the cutover `grove@mascah` was also installed, disabled, at local scope
for this repository and for four nullsec worktrees that no longer exist; those
need no restoring.

**Limits.** Nullsec `main` has no working `grove` until its branch merges: its
committed instructions name the uninstalled predecessor. Nullsec's
`.git/grove/claims.json`, the predecessor's local claim state, is left in
place; the new CLI never reads it. The installed binary comes from an unpushed
commit and should be rebuilt from `main` after the merge. Nullsec's old
branches and worktrees are untouched. Grove verification at `db819bf`:
`gofmt -l .` clean, `go vet ./...` ok, `check` `OK: 87 records`, `go test
-count=1 -timeout 120s ./...` ok in every package; later commits change only
records and the README.

## Next

In Review. To judge it:

```sh
cd ~/GitHub/mascah/grove/.claude/worktrees/G-041
go run ./cmd/grove context G-041 --include grove/G-091-nullsec-cutover-plan.md
go run ./cmd/grove show G-092
git diff --stat CANDIDATE HEAD     # only this record
cd ~/GitHub/mascah/nullsec/.claude/worktrees/grove-cutover
grove check && grove show G-127 | less
grove                              # the board, in a real terminal (acceptance 6)
```

Approve and integrate, nullsec first, since its `main` needs the new layout:

```sh
cd ~/GitHub/mascah/nullsec && git merge worktree-grove-cutover && grove check
cd ~/GitHub/mascah/grove && git merge worktree-G-041
# quote the verdict in this record's Evidence, then:
go run ./cmd/grove update G-041 --set status=done --commit
git push    # both repositories, when wanted
rm -rf /tmp/grove-install && git clone -q --no-hardlinks . /tmp/grove-install && (cd /tmp/grove-install && go build -o ~/.local/bin/grove ./cmd/grove) && grove version
```

Feedback: write it here and `--set status=active` on `worktree-G-041`.
Rejection: `status=abandoned` with the reasons, and the rollback above. After
done, G-036's Next takes the first real nullsec change, shaped in nullsec.
