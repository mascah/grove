# Grove on memex: review of the 2026-10-24 to 10-26 sessions

Handoff note. Written 2026-10-26 from the Claude Code transcripts of
`~/GitHub/mascah/memex` (ten sessions under
`~/.claude/projects/-Users-mascah-GitHub-mascah-memex/`, plus the two
worktree folders `-memex-2c149` and `-memex-syssv`). Grove was 0.5.0
throughout. Read `CONSTITUTION.md` before acting on anything here.

## What happened

Five changes went explore → prototype → propose → implement → archive, all
through the single-change path: 261025-jnqjq (vault + CLI), 261025-2c149
(semantic layer), 261025-vvnhn (runtime tuning), 261026-syssv (Mac app,
writing), 261026-epcxe (Mac app, finding). Two prototypes stayed on
`prototype/*` branches; `prototypes/` is in memex's `.git/info/exclude`.
`grove archive` accepted every fold first time. `grove review`, the reader,
plan, cycles and `run` were never used.

The owner agreed with the four findings below and wants to continue in a
new conversation. Decisions so far are marked.

## Decisions (2026-10-26, second session)

- **1.** The CLI says where to look instead of embedding documents.
  Implement and archive list each change document by path for the agent
  to read with its own file tool, and implement stops printing propose's
  templates (it points to `grove instructions propose <id>` for a
  revision). Output then stays the same size however big the change is.
  How the guarantee "build what the target holds" carries over (point at
  working-tree paths and say when they differ from the target, or
  `git show`) is settled in the proposal. The trials count whether agents
  read every listed document. Measured 2026-10-26:
  - Claude Code: a shell result over 30,000 characters becomes a 2 KB
    preview plus a saved file; one `Read` of the file recovers it.
  - Codex 0.164.0 (default config): a shell result over about 10,000
    tokens (about 40 KB of Grove's text) keeps its start and end and
    drops the middle, with "Warning: truncated output (original token
    count: N)"; nothing is saved. Propose's 18.8 KB and 37.6 KB (propose
    printed twice) came through whole; 56 KB was cut. memex's implement
    jnqjq (49 KB), syssv (47 KB) and archive jnqjq (54 KB) would have lost
    their middle, which is the change itself, without the agent able to
    get it back.
- **2.** Quote the title in the templates (unchanged decision).
- **3.** `Files:` takes one backticked path per bullet; the comma line stays
  valid, so not breaking. A trailing `/` names a folder and everything
  under it. `tasks.md` is frozen during implement: only a propose revision
  edits `Files:`. A file the owner agrees to change beyond
  `Files:`/`Delegated:` goes on an `Also changed:` line under that task in
  `implementation.md`, with the reason; unattended never adds one. Grove
  reports planned against actual in `grove review` (rule 6), from git
  (the branch's changed files against `Files:` and `Also changed:`): each
  changed file as planned, agreed (with the reason) or named nowhere. Not
  a `dev/` script: Grove is a standalone CLI. Review reading `Also
  changed:` needs a rule 3 amendment first, in its own pull request
  (owner agreed). Evidence from memex: `Files:` was
  edited during implement in 3 of 5 changes (7 entries, `Cargo.lock`
  three times), and 3 of 5 changed files no list named even at archive,
  so an agent-kept list of changed files is not reliable.
- **4.** Re-basing with propose's revise path is the normal step after a
  dependency lands; the message, the implement instructions and the docs
  say so. No new logic.

Proposed order: 2 (bugfix schema, small); the rule 3 amendment PR; 3;
1; 4.

## 1. Implement and archive instructions overflow the harness (needs understanding first)

**Symptom.** Claude Code truncates a Bash tool result above 30,000
characters, shows a 2 KB preview and saves the rest to
`tool-results/<id>.txt`. All five `grove instructions implement` runs and two
of five `instructions archive` runs tripped it. Each cost one or two extra
turns (`cat` the saved file, which overflowed again, then `Read`); one
archive was read only in `sed` slices. By the third session the agents
redirected the output to a scratch file before reading it.

**Why the word budget did not catch it.** `just check` counts words in
`internal/instruct` (8,000, rule 12) and #73 "fewer words" trimmed them.
That bounds the fixed text. The printed output also embeds the whole change
as the target holds it, and that is the bulk:

| run | total | fixed text + templates | project index | the change |
| --- | --- | --- | --- | --- |
| implement jnqjq | 49 KB | 16 KB | 0 KB | 31 KB |
| implement 2c149 | 37 KB | 17 KB | 1 KB | 18 KB |
| implement vvnhn | 36 KB | 17 KB | 1 KB | 17 KB |
| implement syssv | 47 KB | 17 KB | 2 KB | 27 KB |
| implement epcxe | 31 KB | 17 KB | 1 KB | 11 KB |
| archive jnqjq | 54 KB | 8 KB | 0 KB | 45 KB |

Of the 17 KB fixed part of implement, about 9 KB is the propose-time
templates (`.grove.yaml`, proposal, change spec, tasks), which implement
needs only when revising a change. The archive run embeds the change plus
its `implementation.md` (7.7 KB of task notes and 3.4 KB of scenarios for
jnqjq). With 17 KB fixed, any change over about 13 KB overflows; memex's
changes ran 11 to 31 KB. Propose (about 20 KB) and explore/prototype (4 to
6 KB) never overflowed.

**Options, none chosen yet.**

- Say so in the thin skill: "the output may be long; if the harness
  truncates it, read the saved file in full before acting." Zero Grove
  change, keeps the extra turn.
- Stop embedding the change in implement/archive and print, per document,
  the `git show <target-commit>:<path>` command instead (the instructions
  already do this for the project index). The agent reads what it needs.
  Rule 1 and 9 hold. Trade: the agent may skip a document.
- Drop the propose templates from implement and point at
  `grove instructions propose <id>` for a revision. Saves ~9 KB; large
  changes still overflow.
- Write the instructions to a file and print its path. Check rule 9
  ("run `<binary> instructions <skill>` and follow the result").
- Harness side: Claude Code reads `BASH_MAX_OUTPUT_LENGTH` (verify the exact
  variable and whether a project `.claude/settings.json` `env` can set it).
  Fixes only Claude.

Open question for the owner: is the embedded change a deliberate guarantee
(the agent builds exactly what the target holds) or a convenience?

## 2. `.grove.yaml` template prints an unquoted title (decision: fix)

`internal/instruct/metadata/change.grove.yaml` has `title: <title>`. Two of
five proposals had a colon in the title and failed `grove check` with
"invalid YAML: mapping values are not allowed in this context". Fix: quote
it in the template (`title: "<title>"`). Check whether the cycle template
and any doc or walkthrough step show an unquoted title too.

## 3. `Files:` rejects a folder path with a trailing slash (decision: fix, and question the format)

Three occurrences: `crates/memex-core/tests/fixtures/vault/` (jnqjq T-01
and T-02) and `app/src-tauri/icons/` (syssv T-02). The parser
(`internal/doc/doc.go`, `filesLine`, `quoted`, `fs.ValidPath`) handles
wrapped lines fine; `fs.ValidPath` rejects the trailing slash and the
message "should list repo-relative paths in backticks, separated by commas"
does not say why. Minimal fix: strip a trailing `/` before validating, or
name the rule in the tasks guide and the message.

The owner also asks whether a comma list on one line is the right
representation. Considerations:

- The agents wrapped it at 80 columns every time; the parser copes, but
  readers see a run-on line.
- One line per path (a bullet list under `Files:`) reads better in Markdown
  and diffs better; the parser change is small (`quoted` already pulls
  every backticked path from the continuation).
- A folder entry means "everything under it"; `check`'s moved-file diff
  (`git diff --name-only base target -- <paths>`) already treats a folder
  path as a prefix, so folders are legitimate entries and the guide should
  say so.
- Whatever the shape, it is read by `check`, `instructions implement`, `run`
  (rule 3), the tasks guide, the trials and the walkthrough; change all of
  them together.

## 4. "changed on main since base" after a dependency lands (owner unsure)

`movedSinceBase` in `internal/cli/check.go` diffs each task's files between
the change's `base` and the target. A change with `depends_on` is proposed
before its dependency lands, so every file the dependency creates is
reported as "changed on main since base <base>; revise the change with
propose" from the moment it merges until the change is archived.

Seen on 2c149 (4 files) and epcxe (7 files, 11 notices). The implement
instructions already cover it: "tell the owner what changed, ask whether to
build on the file as it is now or revise the change with propose". Both
times the owner said build on. The notices then kept firing through
`check` and `archive` (archive ignores them, rule 3); one agent filtered
them with `grep -v`.

Options:

- Leave it; it is rule 3 as written. Make the message say "or re-base with
  propose once <dependency> has landed" so it reads as expected rather than
  as an error.
- Have propose's revise path (which already sets `base` to the target
  commit) be the normal step after a dependency lands, and say so in the
  implement instructions and docs.
- Silence the notice when the file's only changes between `base` and the
  target come from the commits that archived changes named in `depends_on`.
  Needs `git log` per file and a way to recognise those commits; probably
  more code than rule 13 allows without a reproduced fixture case.
- Rule 10 already has this exemption for cycles ("other than by the merges
  that archived the changes of the cycles it depends on"); extending it to
  changes would be a constitution amendment.

## Smaller observations, no decision needed

- After committing two proposals straight onto `main`, the agent still
  called them drafts until the owner asked "aren't these already on main?".
  The propose finish text assumes a branch; one sentence for the on-target
  case would close it.
- `grove check` exits non-zero on the moved-file notices, which breaks
  `&&` chains agents write; they worked around it.
- Agents mint IDs by running `grove instructions propose "x"` just to see
  the template. Harmless.
- Every revision made during implement (jiff dependency, cli change spec,
  Open Vault task, extra `Files:` entries) followed "tell the owner, edit
  once agreed". `implementation.md` statuses and the scenario list were
  kept in every change.

## Where the evidence is

- Condensed per-session timelines and the full text of every `grove`
  command's output are in this session's scratchpad
  (`/private/tmp/claude-501/-Users-mascah-GitHub-mascah-grove/5bec941f-783c-4c9a-9d0c-41cf542d5168/scratchpad/`,
  files `*.c2`, `grove-cmds.txt`, `condense2.py`, `dumpgrove.py`). The
  scratchpad is per session and may be gone; the scripts are trivial to
  rewrite from the transcripts.
- Oversized outputs saved by the harness:
  `9928f2f3-*/tool-results/b1f88es53.txt` (implement jnqjq),
  `b2dd1hzfu.txt` (archive jnqjq), `7209e9f9-*/tool-results/bh9ya3f22.txt`
  (implement syssv), `e02d76cc-*/tool-results/b6b310t8y.txt` (implement
  epcxe).
