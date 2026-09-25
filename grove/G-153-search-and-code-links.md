---
id: "G-153"
type: work
title: "Search record bodies and list the records that describe the code a change touches"
status: active
created: "2026-09-25T19:15:12Z"
updated: "2026-09-25T20:41:04Z"
relates_to: ["G-042", "G-065", "G-108", "G-114", "G-146", "G-151", "G-152", "G-154"]
kind: feature
size: medium
---

## Outcome

A project that keeps its knowledge in Grove can find a record by what its
body says, not only by its title; a session about to change a file, and a
reviewer judging a candidate, can list the records that describe that file;
and the board's review of a candidate shows, beside each file it changes,
the records that describe it. The value is for the adopting project and its
agents: this repository's records serve below only as coverage evidence.

Owner intent, shaping conversation 2026-09-25: after reading a source-based
assessment of [Mex](https://github.com/mex-memory/mex) at commit
`57f565c`, the owner agreed to borrow body search, explicit code-to-record
links and change-aware knowledge checks at review, evaluated against Grove's
existing retrieval before any index or code graph, and corrected the
framing from Grove's own development to the ability an adopting project
gets. [G-154](G-154-listed-constraint-eval.md) measures the effect.

## Constraints

Observed at main `001b271`, 2026-09-25, in this checkout:

- Board search (`/`) matches ID, type, status and title, never the body
  ([search.go](../internal/tui/search.go), [board.md](../docs/board.md#search)).
  Every record's bytes are already in memory: the loader keeps `Source` on
  each record ([metadata.go](../internal/project/metadata.go)).
- `context` takes the selected records' real Markdown links from parsed
  Markdown, resolves each to a project path and lists it with a reason; code
  spans, fences, images and HTML are never links, by design
  ([sources.go](../internal/handoff/sources.go),
  [commands.md](../docs/commands.md#context)). Nothing lists the reverse:
  the records whose links or text name a given path.
- The review view reads a candidate's changed files against its merge base
  on demand (`ChangesContext` in [changes.go](../internal/versions/changes.go))
  and shows them with diffs; nothing relates a changed file to a record.
- The reviewer definition ([grove-reviewer.md](../.claude/agents/grove-reviewer.md))
  finds knowledge conflicts "with the project's `grove list` and a text
  search of the record bodies"; the [work guide](../docs/work-execution.md)
  step 6 and the [shaping guide](../docs/work-shaping.md) step 2 ask for the
  same search. Today that is `grep` over one checkout, outside Grove.
- Scale, measured with a binary built at `001b271` and `/usr/bin/time`:
  146 records, 1.05 MB; `list` 0.36 s; `versions` 0.10 s over two branches
  and two worktrees; `grep` over every body 0.02 s. nullsec, read with the
  installed `grove` from its checkout: 127 records, 43 of them decisions.
- Coverage here: records whose body links a code file are 9 work, 2 review,
  1 plan and 1 question; no term, decision or page does. Linking became a
  habit only recently; older records name paths in code spans (`internal/tui`
  37 times, `repo.Command` 8 times). In nullsec no record links code, 38
  records name a `.rs` file in a code span, 4 decisions name a file or a
  `::` path, and 8 decisions carry an "Implemented by" or amendment
  paragraph (nullsec G-016: the warp curve in `crates/sim/src/warp_profile.rs`;
  nullsec G-033: `ws.rs` and `store::active_operation`). Code spans are the
  adopter's form; links are this repository's.
- Retrospective over the last five merges into main, joining each merge's
  changed non-record files with the records as they stood before it (script
  run in the shaping session; output in the conversation): G-144 (`c294ab5`)
  4 files, 3 records by link, 7 by code span; G-135 (`c6c5a4d`) 1, 1, 5;
  G-125 (`9f94493`) 9, 9, 20; G-114 (`38443ee`) 0 code files; G-123
  (`1d36a99`) 6, 5, 5. Older merges show no link pairs and up to 25 span
  pairs. One sampled flag was intact at file level: G-144 changed
  `versionLine` in `init.go`, and G-140 links that file for `defaultConfig`.
  That is the file-granularity limit, and the measurement to keep taking.
- The current view ([current.go](../internal/versions/current.go), G-042)
  decides which observation of a record is current from Git ancestry; a
  record with two current contents is genuine divergence. A search result
  must show both and never let match order stand in for that decision.
- A shipped document may name a command but no Grove record or path
  (AGENTS.md, [G-146](G-146-how-should-an-adopting-project-r.md);
  [G-152](G-152-shipped-document.md), proposed). One-line list output escapes
  control characters ([record model](../docs/record-model.md)); a printed
  snippet must too.
- [G-151](G-151-strip-grove-repository-pointers.md), proposed on main at
  `f4a21d3` during this shaping, edits the same shipped guides and reviewer
  definition; whichever integrates second reconciles.
- G-114 kept "`context` supplying terms or decisions automatically" and any
  new field or index out of scope, to be reconsidered only if evidence shows
  links present and unread. This record adds a lookup, not a reading, and
  keeps that line.
- Mex, for what is borrowed and what is not: its reverse lookup returns the
  matched symbols so a caller can say why a page was found; its grounding
  keeps lifecycle canonical in Markdown and health (`fresh`, `changed`,
  `missing`, `ambiguous`) derived per checkout, never written back; its Wiki
  search orders hits by deterministic rules, not a blended score. Borrowed:
  a reason per hit and health that is a listing at review time. Not
  borrowed: the code graph, symbol identity fingerprints, body hashes,
  SQLite, and its Inbox, Relay and Workstream objects.

**Proposed design**, binding nobody until a plan settles it:

- One matcher over loaded records with four tiers, each hit carrying its
  tier as the reason and the matching line as a snippet: the ID, type,
  status and title match the board has today; a Markdown link in the body
  whose resolved project path equals the query or lies under it as a
  directory; a code span in the body equal to the query or to a trailing
  run of its path components, so `ws.rs` names `crates/server/src/ws.rs`
  and lists every file it could mean; and body text containing the query,
  case-insensitive. Within a tier, the existing inspection order. No score.
- Three surfaces. The board's `/` searches the current view, showing each
  current content of a divergent record with its source. `grove search
  QUERY [--json]` searches the selected checkout, like `list`, and prints
  ID, type, status, title, reason and snippet, escaped. The board's review
  detail lists beside each changed file the records that link it or name
  it in a code span, from data the board already holds, with no Git process
  beyond the changes read the view already makes.
- Guide text, a sentence each: the work guide's step 5, before a unit,
  search the paths the plan names and read what governs the unit; its step
  6 and the reviewer definition, search the diff's paths and treat each hit
  as a claim to recheck, stale ones being findings for the author; the
  shaping guide's step 2 names the command for the text search it asks for.
- Nothing stored: no field, no health status, no index file. What a change
  did to a described file is a comparison a reader makes at review, as the
  record model already says of `examined`.

Out of scope: SQLite or any on-disk index; a code graph, Tree-sitter or
`go/parser`; symbol identity or content hashes; a frontmatter field for an
observed commit, which stays a body convention; a stored drift or health
status; semantic search or embeddings and wikilink syntax (both kept out by
[G-065](G-065-flexible-records.md)); changing what `context` reads in full;
and converting nullsec's code mentions into links, which is that project's
own choice. Reconsider symbol-level identity when, over several real
reviews, most flagged records had claims the change left intact; reconsider
an index when search over the current view passes about a second or hit
lists need ranking to be readable.

## Scope at preparation

Narrowed 2026-09-25 by the headless `/grove-work G-153` session under the
owner's rule in Next: "if agents already find and apply the listed
constraint without search, preparation narrows this record to the board
search and the review listing". G-154's `without` row (review G-160 on
`worktree-G-154` at `5ff6eb3`, examined `eeec725`) found both constraints
applied in 10 of 10 runs without search, since at fixtures of five and two
records every run reads every record, and named no lever. So `grove
search`, its `--help` and commands.md entry, and the guide and reviewer
sentences are out of this record: nothing yet shows an agent needs them.
They are the owner's to reshape as their own work once a case large enough
that reading every record costs more than choosing shows the listing
missing a constraint (G-160 Disposition). Consequence for G-154: its
`with` row needs a guides or CLI change to compare; this record now makes
none, so that row compares nothing new until such work lands.

## Acceptance

1. On the board, `/` finds a record by a word that occurs only in its body,
   shows why each hit matched and the matching line, and lists title hits
   before link, code-span and text hits. A query that is a project path
   lists the records whose links resolve to it or under it, and those naming
   it in a code span, each with its tier. A divergent record shows every
   current content. [board.md](../docs/board.md) documents it.
2. Narrowed out at preparation (Scope at preparation): `grove search QUERY
   [--json]` in a checkout printing the same tiers.
3. The board's review detail lists, beside each file the candidate changes,
   the records that link it or name it in a code span, with the reason, and
   says when no record names a file; it stores nothing and starts no Git
   process the view did not already start. The owner judges the layout in a
   real terminal, and the terminal lifecycle checks pass.
4. Narrowed out at preparation: the work guide, the shaping guide and the
   reviewer definition naming the command.
5. Evidence records two observations with their commands: on this
   repository, the review listing for the five merges in Constraints
   reproduces the link pairs there; on nullsec, read from its checkout
   without writing, the board search's matcher for a `.rs` path names the
   decisions the Constraints cite.
6. The checks in AGENTS.md pass, no package exceeds five seconds, and no
   file is written anywhere by search or the review listing.

## Evidence

Headless `/grove-work G-153` of 2026-09-25, G-153 alone, on
`worktree-G-153` (`.claude/worktrees/worktree-G-153`), base main
`3f2b923`. Started from this record at `sha256:18a4d3a4…`; narrowed and
planned in `aab404d` (record `sha256:8f6ea855…`, plan
[G-164](G-164-g-153-board-body-search-and-revi.md) `sha256:bb752dc4…`);
implementation `702b576`, review fixes `7f00b39`.

What changed, per acceptance item:

1. Board `/` matches in four tiers, first that applies: `title` (ID, type,
   status, title, as before), `link` (a body link resolving, through
   context's own `resolve`, to the path or under it), `code span` (the span
   is the path or its last components), `text` (a body line, any case).
   Hits sort by tier, then inspection order; with a query each hit has a
   second row with its tier and the escaped line. In the current view each
   current state (`currentStates`, G-042's `Older`) is its own hit naming
   where it is held; a checkout's board searches its own copy. The count is
   of records. [board.md](../docs/board.md#search) documents it. Tests:
   `TestSearchMatchesBodiesByTier` (body word, tier order, path and
   directory query, span by file name, divergent states, hostile line
   escaped, checkout board), and the existing `TestSearchReachesEveryRecord`.
2. Narrowed out (Scope at preparation).
3. The review detail's Changes section has, under each file, `described by
   ID tier, …` or `no record names it`: link or code span only, the open
   record excluded, each ID once, renames by either side, the prefix
   stripped as Git gives it, a file outside the project by code span only.
   It reads the loaded records: no Git process (the test asserts the only
   read is the one changes read), nothing stored (`internal/handoff`
   `Mentions` is pure; the model's cache is in memory and dropped per
   result). Test: `TestReviewListsRecordsDescribingEachFile`, with and
   without a prefix. **The owner's layout judgment in a real terminal is
   open.**
4. Narrowed out.
5. Observations, by throwaway harnesses in `internal/tui` (a `_test.go`
   calling `pathTier` or `match` with `handoff.Mentions`, run then deleted,
   never committed):
   - This repository: for each merge, `git archive MERGE^1 grove.yaml grove`
     extracted under `/tmp`, loaded with `project.LoadFS`, against `git
     diff --name-only MERGE^1 MERGE` outside `grove/`. Code-file link pairs
     reproduce Constraints exactly: G-144 `c294ab5` 4 code files, 3 link
     pairs (guides.go G-107, cli.go and init.go G-140); G-135 `c6c5a4d` 1,
     1; G-125 `9f94493` 9, 9; G-114 `38443ee` 0 code files; G-123 `1d36a99`
     6, 5. Span pairs are more than Constraints counted (12, 4, 32, 0, 10
     against 7, 5, 20, -, 5), since a span naming a file's last components
     (`model.go`) counts for every file it could mean. The reviewer
     reproduced the link pairs independently.
   - nullsec, loaded read-only with `project.LoadFS(os.DirFS(...))` at
     `3eb2785`, `git status` clean after: `match` for
     `crates/sim/src/warp_profile.rs` lists decision G-016, page G-101 and
     plan G-122 by code span; for `crates/server/src/ws.rs`, decision G-033
     and works G-055, G-066, G-068 and plans G-113, G-117, G-124 by code
     span, then pages G-102, G-103, G-104 by text. The installed `grove show`
     from nullsec's checkout confirms both decisions name the files.
6. At `7f00b39`: `go vet ./...` clean, `gofmt -l .` empty, `go run
   ./cmd/grove check` OK; `go test -count=1 -timeout 120s ./...` passed
   at `702b576`, and tui and handoff again uncached at `7f00b39`;
   `python3 internal/tui/testdata/terminal.py` on a binary built at
   `7f00b39`: all 11 checks ok. Under `-short` tui 0.7 s and handoff
   1.0 s; the full tui run is 11.7 s, as on main at `3f2b923` (11.7 s, the
   pty test, skipped under `-short`); cli, integrate, update and versions
   exceed five seconds uncached on main as well, untouched here. Parsing
   every body here takes 17 ms once, a cached search pass 5 ms.

Decisions: the matcher's tiers live in `internal/tui/search.go`, the body
walk in `internal/handoff` beside `resolve`, so links mean exactly what
`context` lists; a span naming last components lists every file it could
mean (the record's design); the review listing shows plain rows, not
selectable entries, so Tab order is unchanged.

Review: [G-167](G-167-g-153-board-search-and-review-li.md).

Limits: the owner's layout judgment (acceptance 3); span matching is at file
granularity and over-reports common file names; G-154's `with` row has no
product change to measure (Scope at preparation); no Linux run.

## Next
