---
id: "G-092"
type: review
title: "G-041 nullsec cutover review"
status: current
created: "2026-09-22T20:30:09Z"
updated: "2026-09-22T20:34:20Z"
work: ["G-041"]
examined: "4b05447"
---

Two independent readers, both Claude subagents of the implementing session
with read-only instructions; neither edited anything.

## Round 1: rewritten lines, in the rehearsal

Examined: the disposable clone of nullsec `366b063` after the script (the tree
later reproduced live), all 1,368 rewritten lines of that run (the final run,
after the fixes below, rewrote 1,374), and every converted body
against its `formerly` original.

| Finding | Disposition |
| --- | --- |
| A verbatim builder quote in old W-005 (now G-048) ends "close W-005." | Kept as written: exception in the script, listed on the map page |
| G-063 quotes a server log naming `docs/plans/W-020-playtest-catalog.json` | Left rewritten: the file moved, and a rerun prints the new path |
| The brief's Next moved into G-075 was rewritten but introduced as the brief's wording | Lead-in reworded |
| `W-019-T3` is a branch name (merge `a6c37b9`) | Kept: IDs followed by `-T<digit>` are not rewritten |
| `W-021/023/024` and `D-0034/0035` shorthands half-rewritten | Expanded to full IDs before rewriting |
| `history/evidence/…` paths relative to the old root no longer resolved | Rewritten through aliases |
| Predecessor commands such as `grove close W-001` in done records | Rewritten like other IDs, as G-052 did; history, as is the finished study brief `art/design/w039-crew-hull-test.md` |
| Code risk | None: non-comment rewrites are test titles, assert messages, a `console.log` label and trace headers no code compares; the retrospective generator and `docs/RETROSPECTIVE.md` stay consistent |

## Round 2: final, both branches and the machine

Examined: Grove `worktree-G-041` at `db819bf`, nullsec `worktree-grove-cutover`
at `ecbb036`, and the machine state. All 126 converted records' frontmatter and
bodies compared with their originals by script; all 162 tracked Markdown files'
links resolve; all 146 predecessor paths appear on the map page.

| Finding | Severity | Disposition |
| --- | --- | --- |
| G-121's `work` lacked G-073: `W-029-W-031` is a range | should fix | Set to G-072, G-073, G-074 (nullsec `3eb2785`) |
| G-075 has two contradicting Next paragraphs | should fix | Lead-in says G-075's own paragraph supersedes the moved text (`3eb2785`) |
| Three live docs still name predecessor commands | should fix | `section-study-recipe.md`, G-059's Next and `expedition-foundations.md` use new commands (`3eb2785`); done records keep theirs as history |
| README still calls `../skills/` the working CLI and the resume prompt names their Grove CLI | should fix | Corrected on this branch |
| The rollback is referenced but not written | should fix | Written in G-041's Evidence, with the removed Codex configuration verbatim |
| nullsec `main` has no working CLI until its branch merges | minor | Handoff merges nullsec first |
| `sources` and plan-evidence mappings missing from the rules | minor | Plan and map page state them |

No findings for the `superseded` status, frontmatter values, bodies, links,
map completeness, nullsec's Grove section, or the machine state. The reviewer
did not rerun `npm test` or compare the mapping files.

## Round 3: the fixes and the evidence

The same reviewer, at Grove `4b05447` and nullsec `3eb2785`: every disposition
above holds, both `check` runs pass, and the Evidence's commits, uninstall
steps and session results match the files, Git and the session's recorded
tool output. It could not observe the 694-file comparison, the base-commit
smoke failure or `go test` at `db819bf`. Fixed after it: the Codex rollback
also re-adds the plugin (its cache was emptied), the plugin version is named,
the `codex exec` stdin detail and the branch-name exclusion are stated, and
the line counts name their runs. Its note that status was still `active` is
the handoff's own last commit.
