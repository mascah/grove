---
id: "G-184"
type: plan
title: "G-169 plan: entrypoint revisions, init --check and the review guide"
status: current
created: "2026-09-25T23:04:16Z"
updated: "2026-09-25T23:04:44Z"
work: ["G-169"]
---

Plan for [G-169](G-169-harness-upgrade-compatibility.md), prepared headless
in `worktree-G-169` from main `6b14141`, record at `sha256:22d79d21…`.
Single implementer, sequential; one independent review on the final
revision.

## Design

- **One owner of the assignment grammar.** The work guide's Inputs already
  state IDs, `--until plan` and the mode; the shaping guide's Inputs state
  topic, IDs and mode. Each gains the one sentence the adapters carried
  alone (anything else is an error; the bound is for the guide). The
  generated adapters and this repository's own adapters keep only the
  argument-as-data rule and defer to the guide's Inputs. Claude's
  `argument-hint` stops listing options (`[OPTIONS]`): it is display
  metadata, and a stale list there would be a second owner.
- **Review guide.** The reviewer definition's body moves to
  `docs/work-review.md`, embedded and printed by `grove guide review`, a
  shipped document under the same rules and in the guides digest. The
  generated `.claude/agents/grove-reviewer.md` keeps its frontmatter
  (`model`, `effort`, `disallowedTools`), states its read-only authority,
  and loads the guide as the skills do. This repository's reviewer becomes a
  development adapter like its skills: no marker, reads
  `docs/work-review.md` as a file. `grove.Reviewer` (the embed of that file)
  is removed; the template is generated in `entrypoints.go`.
- **Entrypoint revision.** An integer, `grove.EntrypointRevision = 2`, with
  `MinEntrypointRevision = 1`, separate from the release version and the
  record schema. Every managed file carries a line
  `grove entrypoint revision 2` (an HTML comment in Markdown, a `#` comment
  in YAML) after the marker. A marked file with no such line predates
  revisions and is revision 1, the templates G-040 introduced: legacy.
  The generated adapters load `grove guide NAME --entrypoint 2`; `guide`
  refuses (exit 1) a revision outside `MIN..CURRENT` with the repair (use
  the grove that wrote it, or rerun `init` and commit), and a grove before
  this change refuses the unknown flag (exit 2). Both are the "fails or
  prints anything other than that guide" the adapter already stops on. A
  plain `grove guide NAME`, which a person and a legacy adapter both run,
  always prints: `guide` cannot tell them apart, so a revision-1 file
  reaches diagnosis through `init --check` and `run`, as below.
- **Diagnosis.** `grove.Diagnose(path, content)` classifies one installed
  file: `custom` (no marker: project-owned, not judged), `current`
  (byte-equal to this binary's template), `compatible` (marked, supported
  revision, other bytes: older or edited), `legacy` (marked, no revision
  line; revision 1, compatible while the minimum is 1), `incompatible`
  (revision unparseable or outside `MIN..CURRENT`). A file's absence is
  `missing`. Content inequality alone never makes a file incompatible.
- **`grove init --check`.** Read-only: plans the entrypoints as `init`
  does, prints `VERDICT PATH (note)` per managed path, writes nothing, and
  exits 1 when any path is `missing` or `incompatible`, 0 otherwise.
  Custom files are listed and not judged.
- **Runner.** `grove run` and the board's `R` (same `Start`) refuse before
  the attempt directory and the provider when the worktree's `grove-work`
  skill or reviewer is `incompatible`, naming the path, the revision and the
  repair. A missing reviewer stays a warning (G-150).
- **Repair path.** `grove init` (unchanged ownership rules) rewrites marked
  files and keeps unmarked ones, configuration, brief and records; the
  result must be committed, since worktrees hold only committed files. An
  existing worktree keeps its branch's copy: run `init` there, or merge the
  target after the refresh is committed, and commit. A session that already
  loaded an adapter or guide keeps what it read: start a new session (or
  reload skills) after refreshing. Documented in the command reference's
  Init and Attempts sections and `init`'s help; the guide says nothing new
  beyond the grammar.

Rejected: comparing against a list of historical template digests (a
growing table, and an edited file would still need a verdict); refusing
legacy files now (they work with this binary); gating `guide` without the
flag (it would break a person's read).

## Steps

1. `docs/work-review.md` from the reviewer body; `GuideFiles["review"]`;
   guides digest covers it; generated reviewer template; dev reviewer
   adapter; grammar sentences into both guides; slimmer adapters (generated
   and dev).
2. `EntrypointRevision`, `MinEntrypointRevision`, revision lines,
   `Diagnose`; `guide --entrypoint N`.
3. `init --check`.
4. Runner refusal for incompatible skill or reviewer.
5. Tests: an old/new fixture pair (the files `init` at `6b14141` writes,
   kept as a txtar-style file under `internal/cli/testdata/` so no nested
   `.claude/` directory exists) through `init --check` (legacy), `init`
   (updated, custom kept), `init --check` (current), `init` again
   (unchanged); a newer revision refused by `guide` and `run`; guide
   refusals; shipped-document checks over the review guide.
6. Documentation: `--help`, `docs/commands.md` (Init, guide, version,
   Attempts), README where adoption mentions the reviewer, CLAUDE.md's
   shipped-document rule, term G-152's meaning, G-169's dangling G-172 link.
7. Verification per CLAUDE.md, including `terminal.py` (the runner's board
   path), then independent review.

Live Claude/Codex session trials (acceptance 4) are out of this
assignment's mandate: none states resource bounds for them. They are
reported as limits, not simulated.
