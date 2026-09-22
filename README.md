# Grove

A local project workspace for humans and agents, built around a CLI and durable files.

**Status: the first Go CLI lists, shows, validates, creates, and updates local
project records, shows each record's versions across local branches, and
locates the checkout holding a selected version. Run without a command, it
opens a read-only terminal Kanban board over the same operations. `context`
assembles staged context for selected work, the `grove-work` skill carries it
out, and the `grove-shape` skill shapes proposals.**

The selected next milestone is a complete interactive shape → implement →
review → integrate loop on real nullsec work. Start with
[G-036's adoption roadmap](grove/G-047-adoption-roadmap-plan.md) for the ordered
work. Its proposed capabilities are not commands available in this build:
the current board remains checkout-scoped. Work has the Review status: an
implementation ends with its work record in Review, naming its `candidate`
commit, and `done` is written where that candidate was merged.

Start with [the restart brief](grove/brief.md) and
[the accepted record model](docs/record-model.md). The brief records the selected
direction and next investment. The
[direction evaluation](grove/G-048-direction-evaluation-review.md) retains
research and evidence from Bench, the sibling skills and the current prototype. [grove.yaml](grove.yaml)
configures the record tree. The completed first implementation is
[CLI inspection](grove/G-003-inspect-records.md).
The installed `grove` still belongs to the sibling skills project; it does not
read this new format.

## Use the CLI

Requires Go 1.26 or later. Run from this repository:

```sh
go run ./cmd/grove                 # the terminal board; needs a terminal
go run ./cmd/grove list
go run ./cmd/grove list --status active --status review  # only records in those statuses
go run ./cmd/grove show G-003
go run ./cmd/grove check
go run ./cmd/grove new work "Title of the work" --slug short-name
go run ./cmd/grove new term "Attempt"  # also question, decision, plan and review
go run ./cmd/grove new page "Notes"     # general knowledge, no status
go run ./cmd/grove convert notes/old-plan.md --type plan --title "Old plan"  # prints the mapping
go run ./cmd/grove brief               # the brief grove.yaml names
go run ./cmd/grove show G-003 --json
go run ./cmd/grove update G-003 --expect sha256:HEX --set status=active --unset size
go run ./cmd/grove update G-003 --set status=done --commit  # at a shell: no lookup, one commit
go run ./cmd/grove versions G-003 --json
go run ./cmd/grove workspace --source SELECTOR --json
go run ./cmd/grove --project "$(go run ./cmd/grove workspace --source SELECTOR)" show G-003
```

The first three commands read live files without modifying them; `new` adds
one file and prints its path. `list` shows ID, type,
status, and title, and with `--status VALUE` (repeatable) only the records
in any given status, refusing a value outside the record model's status
vocabulary; `show` prints the exact Markdown source, or with `--json`
one object holding the path, a `sha256:` content revision, and the source;
`check` validates metadata and relationships. `update` changes frontmatter
fields of one record, keeps every other byte of the file, sets `updated`, and
prints the resulting revision. `--expect REVISION` is optional: it refuses a
file that no longer hashes to it, which an agent session passes because its
read may be old, while a person at a shell omits it. `--commit` then commits
that one file with a generated message such as `docs(G-003): set
status=done`, leaving every other path as it was, and adds `commit` to the
result. Lists are JSON arrays such as `'["G-003"]'`; `--unset` removes an
optional field. Project/file context and errors go to stderr, so
stdout can be redirected. Exit codes are 0 for success, 1 for inspection/output
errors, and 2 for invalid command usage.

Without `--project`, discovery searches upward for `grove.yaml` and stops at
the current Git checkout boundary. Plain directories also work. Any invalid
record makes the command fail; no partial list or record is printed.

Every `.md` beneath the record root is a record wherever it sits, and folders
mean nothing: `new` gives every type a neutral `G-NNN` ID in a flat
`ROOT/G-NNN-slug.md`. Work, questions, decisions, terms, plans and reviews keep
their own rules. A plan or review names its work in a `work` list, set with
`update`, and `context G-NNN` lists the plans and reviews attached to the
selected work without reading them; a review can record the Git commit it
`examined`. Work moves `proposed`, `active`, `review`, `done`, with
`abandoned` for an explicit human decision. `review` requires `candidate`, the
commit offered for judgment, and `update` writes `done` only with a candidate
that the checkout's HEAD contains: Done means accepted and merged, so the
integrator writes it in the target's checkout after the merge, where the
check holds it to the code that landed. A `done` record without a candidate
predates that meaning. `new page "Title"` creates general knowledge with a title and no
status; pages are never work cards and `context` reads one only through
`--include`. `update --set type=...` reclassifies in place, keeping ID and path.
`convert` turns a Markdown document outside the record root into a record and
prints the old-to-new mapping. `brief: PATH` in `grove.yaml` names the project
brief, which `brief` prints and `check` requires to exist.

`schema_version: 3` is the only schema: Grove keeps no backward compatibility
before its first release. [G-052](grove/G-052-migrate-knowledge.md) converted
this repository from typed IDs (`W-001`) in type folders, and
[G-069](grove/G-069-migration-map.md) maps every old ID and path to its
counterpart. Inspect an older commit with the CLI in that commit; `versions`
and the board report a branch that predates the conversion as a source they
cannot inspect. The
[record model](docs/record-model.md#identity-and-placement-apart-from-classification)
has the page boundary and conversion's limits.

`new` and `update` require a Git checkout. `new` takes the next number
from a counter under the repository's common Git directory,
shared by every linked worktree, and floors it by the highest ID on any local
ref or worktree. Never number new records by hand. Both commands serialize
through a write lock in that same directory; `update` refuses a stale
`--expect`, an invalid project, or any change it observes while preparing the
write, and reports when a failure happened after the file was replaced,
including a `--commit` that Git refused: the file then holds the update,
nothing was committed, and the message says whether it was staged.

`versions [ID]` shows one row per version of each record: its committed
version on every local branch tip and its live version in every registered
worktree, grouped by ID with each source's own status, so main can see a
feature branch's progress without switching or merging. Live rows say how the
file compares with that checkout's HEAD (`unchanged`, `modified`, `renamed`,
`added`, `deleted`). Every row ends with a selector that binds the repository,
source, commit, configuration, record path, and content revision; identical
bytes in two sources get two selectors. Sources are listed on stderr with any
diagnostics; an invalid or unreadable source makes the result incomplete and
the exit code 1 while valid sources still print. A live project must be
inside its registered checkout: a project location reached through a symlink,
or belonging to another nested repository, is invalid, while a missing one is
simply absent. `--json` adds each version's
exact source text. No status is chosen as authoritative, and the command
writes nothing: no refs, index, worktrees, records, or coordination state.

`workspace --source SELECTOR` takes one selector from `versions`, checks that
the version is still exactly what was selected, and prints the absolute
project directory of the existing checkout holding it, ready for `--project`.
A live selection resolves to its worktree. A committed selection resolves to
the one worktree that has that branch checked out, provided its live copy of
the record still has the committed bytes; otherwise, or when no or several
worktrees hold the branch, it refuses and says to run `versions` again and
select a live version. Any change since selection (branch or HEAD moved,
worktree moved or removed, attached or detached state, configuration, record
path, or content) refuses with the reason, and the selected checkout is
re-read once more just before its path is returned. The command never creates a
worktree, switches a branch, launches anything, claims ownership, or edits a
record; the directory it prints is a location, not write authority, and a
later `update` performs its own revision check. The board below uses these
two operations in process; creating a worktree for a branch without one
remains future work.

### Work context and the `grove-work` skill

`grove [--project DIR] context WORK_ID... [--json] [--interaction interactive|headless] [--max-bytes N] [--include PATH]...`
assembles context for explicitly selected work in one checkout, so an
assignment is a few IDs rather than a composed prompt. It prints the IDs
ordered prerequisites first (ties keep the requested order; prerequisites are
never added to the selection), Git identity, and two kinds of content that it
keeps apart:

- **Sources, read in full**, each once with its `sha256:` revision:
  `grove.yaml`, the selected work records, and every `--include PATH`. That is
  the whole default. An include is a required project-relative file of any
  type; one file reached by several names (letter case, a hard link) is
  included and charged once.
- **Observations, listed and not read**: the transitive `depends_on`
  prerequisites, questions blocking the selected work or a prerequisite
  (resolved ones too), and the records the selected work names in `relates_to`
  or `members` or links to, each with title, status, path, revision, why it is
  listed, and whether its source is included (`source` is the `sources[]` path
  that holds it, which differs from `path` when it was included under another
  name). Every link in the selected
  records' bodies is listed with the project path it resolves to. Links are
  taken from parsed Markdown, so code, fences, images, and HTML are never
  links; Markdown escapes and entities are decoded before the URL is. A link is
  never opened, so a listed path is not checked, not even for existence. URLs,
  fragment-only links, other projects, absolute paths, and Git metadata are
  listed with that reason and no path.

Retrieval is staged with existing commands: `show ID` prints a listed record,
and `--include PATH` adds a listed file, such as the plan the record names as
current, with its revision. The command picks no plan or review by itself, and
it never follows links inside an included document. A listed record's revision
says which version was seen; a listing is not a reading of its constraints.

It refuses rather than degrade: unknown, duplicate, or non-work IDs, a malformed
link destination, an included file that is missing, a symlink, not regular, or
not UTF-8, a source that does not fit `--max-bytes` (default 262144, counting
source bytes, at most 8388608), or a checkout that changed while it was read.
Nothing is truncated or summarized to fit. A refusal prints nothing
to stdout. Exit 0 means context was assembled, never that work is ready,
authorized, or integrated: a prerequisite's `done` is its recorded status, and
only one with a `candidate` claims a merge, where it was written.
`--interaction` records whether a person can answer (default `interactive`);
it is the caller's declaration, passed through for the guide to act on. The
command writes nothing and starts nothing, and needs Git only inside a
repository. Text output escapes terminal controls and fences each source;
`--json` has the exact source strings (`format_version` 2; version 1 read
prerequisites, related records, and linked documents in full):

```text
{format_version, root, interaction, selected[], order[],
 git: null | {checkout, common_dir, ref, head},
 records[]: {id, path, type, title, status, revision, roles[], selected, included, source},
 requirements[]: {work, prerequisite, status, selected},
 questions[]: {id, status, blocks[]},
 sources[]: {path, revision, reasons[], content},
 references[]: {from, target, path, reason},
 scope_notice, source_bytes, max_bytes}
```

The [work guide](docs/work-execution.md) is the workflow that uses it: staged
reading, isolation before the first write, preparation, review with bounded fix
rounds, checkpoints, and what to do when a human decision is missing,
interactive or headless. The `grove-work` skill is a thin adapter to that guide
for [Claude](.claude/skills/grove-work/SKILL.md) (`/grove-work G-030`) and
[Codex](.agents/skills/grove-work/SKILL.md) (`$grove-work G-030`), invoked
explicitly and kept in this repository on purpose while it is dogfooded;
[AGENTS.md](AGENTS.md) holds this repository's development policy, which the
guide does not repeat. Grove launches no agent;
[the dogfooding evidence](grove/G-032-dogfood-review.md) says which
invocations have actually been exercised.

### Shaping and the `grove-shape` skill

The [shaping guide](docs/work-shaping.md) is the workflow before an assignment
exists: discuss an idea or existing records, check other branches and
worktrees for overlap, and write proposed work, real human questions, and
attributable decisions with today's `new`, `update`, `versions`, and `check`.
It keeps intent, observed evidence, proposed design, and decisions apart, and
it never assigns, implements, promotes status, or merges. The `grove-shape`
skill is the same kind of thin adapter for
[Claude](.claude/skills/grove-shape/SKILL.md) (`/grove-shape TOPIC`) and
[Codex](.agents/skills/grove-shape/SKILL.md) (`$grove-shape TOPIC`);
[its evidence](grove/G-050-shaping-review.md) says what has been
exercised.

### The terminal board

`grove [--project DIR] [--json]`, with no command, opens the board. In this
checkout that is `go run ./cmd/grove`. There is no `board` subcommand and no
work-ID argument; every command above stays noninteractive, and `--help`,
`-h`, and `help` need neither a project nor a terminal.

The columns (Proposed, Active, Review, Done, Abandoned) show the live work records of
one checkout, named in the header: at first the checkout the command ran in.
`b` chooses another checkout's live files as the board; this changes what is
displayed and switches no branch or directory. Work with no live record in
that checkout is listed under Elsewhere without a status. Questions and
decisions are not on the board. No status is combined across branches.

Enter on a card opens that record's versions. A version is the record's exact
content; the board read it at every local branch's tip and in every checkout's
files, and lists each differing content once with its own title and status.
Where several branches and checkouts hold the same content the row is a fold
(`▸ done  same on 4 branches, 4 checkouts`): Enter lists those places, and
selects nothing. Opening a card selects nothing either. Moving to one branch
or checkout and pressing Enter asks `workspace`'s resolver about exactly that
version's selector; on success the board closes and prints
what `workspace` prints (the project path on stdout, or its JSON with
`--json`; checkout, branch, record, and revision on stderr). A refusal (the
version changed, its checkout is missing or ambiguous, the record was deleted
there) stays on screen with its reason until `r` refreshes, after which a
version must be selected again. The board never creates a worktree, edits a
record, or starts an editor, shell, or agent.

The card's details begin with History: the commits that changed the record's
file, newest first, each with its date, the status the record held at that
commit, its short ID, and its subject, following renames. It is the history of
whichever row has focus, named in its heading: the board's checkout while the
card's first line has focus, otherwise that branch's tip or that checkout's
HEAD. A checkout whose files differ from its HEAD gets a first `uncommitted`
row. Merges are not listed, so where the record's status is not the newest
listed commit's, a first `here` row gives it and says why. History is read from
Git when a card is open, never while the board loads, and no key waits for a
read still in progress. It says what happened on one branch and nothing about
whether another branch contains it. [G-030](grove/G-030-card-lineage.md)
owns this.

Keys: arrows or `h` `j` `k` `l` move; Tab switches between the columns and
Elsewhere, or between versions and details; PgUp/PgDn scroll details; `s`
lists every branch and checkout read with its diagnostics, which stay reachable while a banner
marks an incomplete result; `r` re-reads; Esc goes back, and quits from the
board; `q` quits. Below 100 columns one status column shows at a time; below
40x10 the board asks for more room. Leaving without a selection prints nothing
and exits 0; Ctrl-C exits 1; a usage error exits 2.

The board draws on stderr and reads stdin, so both must be terminals, while
stdout may be redirected: `cd "$(go run ./cmd/grove)"`. Without a terminal it
refuses at once with exit 1 and names the noninteractive commands. Text from
records, paths, and Git is shown with control characters escaped. It writes no
files, including the framework's debug logs.

Use `go test -short ./...` while iterating and `go test -count=1 -timeout 120s
./...` plus `go vet ./...` as evidence; `-race` is per package only, since a
whole-suite race run hangs in the Go toolchain on macOS. `internal/tui` also
drives the built binary through a pseudo-terminal with
`python3 internal/tui/testdata/terminal.py BINARY` (Unix; skipped without
`python3` and under `-short`).
Run `lefthook install` once per clone: pre-commit formats staged Go files and
runs `go vet` and `go mod tidy -diff`; pre-push runs `go test ./...`.
Agent execution remains future work. The
[integrated CLI review](grove/G-022-integrated-cli-review.md) found
workspace-provenance, update-preservation, and Git-path defects; G-014 through
G-016 repair them, with [evidence and remaining limits](grove/G-028-repairs-review.md).
[G-017](grove/G-017-terminal-picker.md) owns the board's contract and its
[evidence](grove/G-029-board-review.md), including the owner's judgment
from a demo, which automated checks do not supply.

The intended experience combines linked work, questions, research, project
knowledge, and evidence. A CLI serves agents and humans; a TUI can make the
project visible and eventually launch agent sessions. An optional local web UI
can operate on the same records. Core project operations require no hosted
service or persistent daemon.

## Local restart

On 2026-09-18, the first Grove application checkout was preserved intact at
`../grove-archive-2026-09-18/` and this fresh repository replaced `../grove/`.
The archived checkout retains its Git history and ignored local data. Its last
commit was `be40e46`. Existing deployments and external data were not changed.
The archive's instructions and service-based planning authority belong to that
old application; they do not govern this restart.

`../skills/` remains the working Grove skill suite and CLI. `../nullsec/` remains
the real project providing workflow evidence. Neither is migrated by this reset.

## Resume this conversation

Give a new agent this prompt:

> Read AGENTS.md, grove/brief.md, and docs/record-model.md. Continue the file-backed
> Grove CLI and interactive workspace described there. The old application was
> archived. The Go list/show/check/new/update CLI and its records now work
> locally, as do versions and workspace; G-003, G-007 and G-009 to G-011 record
> verification. Find the next
> action in G-036 and the work records' Next. Treat the brief's
> remaining proposals as proposals. Inspect
> the sibling skills and nullsec projects through their Grove CLI when evidence
> is needed. Preserve this direction and update the brief as choices settle.
