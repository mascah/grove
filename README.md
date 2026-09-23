# Grove

A local project workspace for humans and agents: work, questions, decisions,
plans, reviews and knowledge as Markdown files with YAML frontmatter in your
Git repository, read and changed through a CLI and a terminal board, with no
service to run.

This build:

- lists, shows, validates, creates, updates and converts records, with
  sequential IDs shared by every linked worktree;
- shows each record's versions across local branches and worktrees, decides
  which are current by Git ancestry, and finds the checkout holding one;
- assembles staged context for selected work, and carries two shared
  workflows inside the binary: shaping ideas into proposed work, and executing
  assigned work through to a candidate in Review;
- records the owner's approval or feedback on a candidate and merges an
  approved one locally;
- runs one bounded headless agent attempt of a work item as a process that
  outlives the terminal, and lists and stops attempts;
- opens a terminal board, with no command, over all of the above.

Where to go next:

- Try it here: [Use the CLI](#use-the-cli).
- Set it up in another repository:
  [Adopt Grove in another repository](#adopt-grove-in-another-repository).
- Shape or carry out work with an agent:
  [the `grove-shape` skill](#shaping-and-the-grove-shape-skill) and
  [the `grove-work` skill](#work-context-and-the-grove-work-skill).
- The record format and lifecycle: [the record model](docs/record-model.md).
- Why Grove exists and where it is going: [the brief](grove/brief.md).
- Developing Grove itself: [AGENTS.md](AGENTS.md).

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
go run ./cmd/grove approve G-003 "Meets the outcome"    # in the branch's checkout: approved=candidate, verdict appended, committed
go run ./cmd/grove feedback G-003 "Handle the empty case"  # back to active with the text appended; prints where to continue
go run ./cmd/grove integrate G-003 --cleanup            # in the target's checkout: merge, done, then the worktree and branch removed
go run ./cmd/grove run G-003 --budget 5 --permission-mode auto  # one headless attempt in its worktree, outliving this terminal
go run ./cmd/grove attempts G-003                       # running, finished, orphaned or interrupted; grove attempt ATTEMPT for one
go run ./cmd/grove stop G-003.20260922T210000Z          # ends the turn, kills after 15 s, writes the result
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

`approve ID VERDICT` and `feedback ID TEXT` judge a work record in `review`
from a clean checkout of the branch that holds it: `approve` sets `approved`
to the candidate and appends `Verdict on candidate X, DATE: …` to the body;
`feedback` sets `active`, unsets `approved`, keeps the candidate, appends
`Feedback on candidate X, DATE: …` and prints on stderr where to continue
(`/grove-work ID` in that checkout). Each commits the record alone and prints
what `update` prints. Both refuse a checkout whose HEAD lacks the candidate or
whose record has uncommitted changes, and `approve` refuses a tip that changed
any other file after the candidate, since that tip is a new candidate.
`integrate ID [--cleanup]` runs in a clean checkout of the target branch that
`grove.yaml` names: it finds the one branch holding an approved candidate of
ID, merges it with a plain `git merge` (a conflict is aborted and refused,
leaving the target as it was), writes `done` there committed alone, and with
`--cleanup` removes the branch's worktree and the branch, keeping either with
Git's reason when Git refuses, and keeping a worktree that holds ignored
files, which Git would delete. It prints one line per fact as it holds
(`approval:`, `merge:`, `done:`, `cleanup:`); a refusal comes before the merge,
and nothing undoes a merge that happened.

`run ID --budget USD --permission-mode MODE [--model MODEL] [--branch NAME]
[--worktree DIR]` starts one bounded implementation attempt of proposed or
active work as a Grove-owned `claude -p "/grove-work ID --interaction
headless"` process that outlives the terminal: it creates `worktree-ID`
under `.claude/worktrees/` from this checkout's HEAD, or reuses the branch's
registered worktree so a next attempt continues from preserved partial work,
then starts an owner process in its own session that runs the provider there
with `--output-format stream-json`, `--max-budget-usd`, `--permission-mode`
and `--permission-prompts none`, its stdout and stderr written straight to
files under the Git common directory (`.git/grove/attempts/ATTEMPT/`, shared
by every worktree, never committed). Budget and permission mode are required:
Grove sets no default spend or profile. Grove starts one process and never
retries; subagents the provider starts share the budget. `run` refuses a
record that is not proposed or active here or on the branch's worktree (a
candidate in review there awaits the owner's judgment), an open question
that blocks it in either place (the wait the headless guide persists, so
rerunning with nothing changed refuses the same way), uncommitted changes to
the record here, a running or orphaned attempt of the same work, and a
worktree path that is something else. Every refusal comes before a write,
except that what an existing branch holds is checked in its checkout, so a
branch that had no worktree keeps the one `run` made. `attempts [ID]` lists attempts newest first;
`attempt ATTEMPT [--json]` prints one attempt's launch, event counts (parsed
bounded: a line over 1 MiB is counted, not read), the provider's init and
result fields, the result (exit, the worktree's HEAD and whether it holds
uncommitted or untracked changes), the record as the branch holds it, whether
the record on the target changed since launch, and the file paths; no
provider text is printed. Liveness is the owner's file lock, never a pid: `running`
while it is held, `finished` once `result.json` exists, `orphaned` when the
owner is gone but the provider's process group lives, `interrupted` when
nothing is left and no result was written (a machine restart reads so; no
reboot recovery is selected). `stop ATTEMPT` sends SIGINT through the owner
so the provider ends its turn, SIGKILL to its process group after 15 s, and
writes the result; an orphaned attempt is stopped directly and reconciled.
Stop touches neither the worktree nor the record. A result is facts, never
acceptance: the record's own status on the branch, which the headless guide
sets, is the handoff, and a process exit or a `result` event proves nothing
about it. `GROVE_CLAUDE` names another executable, for fakes.

Without `--project`, discovery searches upward for `grove.yaml` and stops at
the current Git checkout boundary. Plain directories also work. Any invalid
record makes the command fail; no partial list or record is printed.

Every `.md` beneath the record root is a record wherever it sits: `new`
gives every type a neutral `G-NNN` ID in a flat `ROOT/G-NNN-slug.md`, and
folders mean nothing. Work moves `proposed`, `active`, `review`, `done`, with
`abandoned` for an explicit human decision; `review` names a `candidate`
commit, and `done` is written on the target after that candidate is merged.
A plan or review names its work in a `work` list, and a page is general
knowledge with no status. `brief: PATH` in `grove.yaml` names the project
brief, which `brief` prints. [The record model](docs/record-model.md) owns
these rules: types, fields, statuses, validation, what `update` refuses, and
what `convert` does and does not rewrite. `schema_version: 3` is the only
schema; inspect an older commit with the CLI in that commit, and use
[G-069](grove/G-069-migration-map.md) to follow an old typed ID (`W-001`).

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
simply absent.

The `CURRENT` column says whether a version is `yes`, current, or `older`.
A version is older when another has different bytes and the record at their
common commit has the older one's bytes, so only the other side changed it
since they split. The common commit is the two commits' merge base, or a
checkout's HEAD for its own uncommitted edit. A revert is a change like any
other. Several current contents are a divergence, shown as they are. A pair
whose common commit cannot be read, or whose several common commits hold
different versions, stays unordered: both current, with a note on stderr.
Reverts carried across merges can make versions older than each other in a
cycle; then a version is current when everything newer than it, through any
chain, is also older than it, and a note says so. A branch whose current state removes the record gets a committed
`deleted` row. Dates, status order, and branch names never decide, so the
answer is the same from every checkout, and no branch is special.
[G-042](grove/G-042-current-view.md) owns this. With an integration target
named in `grove.yaml` (`target: main` here), stderr names it and the `TARGET`
column says whether that branch holds the version's bytes (`yes`, `no`; `-`
without a target). The target labels only and decides nothing.
`--json` adds each version's exact source text, `current`, `older` (why it is
older), and `on_target` (null without a target), each record's `notes`, and
the result's `target` and target `notes`. The command writes nothing: no refs, index, worktrees, records, or
coordination state.

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
two operations in process, and selecting a version never creates a worktree
for a branch without one.
### Adopt Grove in another repository

Build one binary from a named commit and put it on `PATH`. On the owner's
machine it is `~/.local/bin/grove` ([G-041](grove/G-041-nullsec-pilot.md)
records how it replaced the predecessor and how to restore that); rebuilding
is manual, and `grove version` names what is installed:

```sh
go build -o "$HOME/.local/bin/grove" ./cmd/grove   # from a clone at that commit
GOBIN="$HOME/.local/bin" go install github.com/mascah/grove/cmd/grove@COMMIT  # or from the module
grove version
```

`go install …@COMMIT` resolves only a pushed commit.

Codex runs each command through a login shell, so it sees the profile's
`PATH`, not the caller's: put the build directory on the login `PATH` ahead of
any other `grove`, or name the executable in the target's `AGENTS.md` or
`CLAUDE.md`, which the entrypoints defer to for how the CLI is invoked. When
the wrong `grove` answers, the entrypoints stop and say so rather than act.

`grove version` prints the module version and, when the build stamped it, the
VCS revision with `modified` for a dirty tree; `go run` prints `(devel)`, which
means this checkout's files. Build from a primary checkout or a clone: for a
linked worktree that lies inside its repository, Go only recognises the
enclosing checkout's `.git` directory and stamps that checkout's revision and
cleanliness instead (observed with go 1.26.2). The predecessor rejects
`version` as an unknown command, so the line tells the two apart, and the
line ends with a digest of the embedded guides, which names the workflow even
when no revision was stamped. The workflow guides travel inside the binary:
`grove guide work` and `grove guide shape` print them, so the workflow version
is the executable version and no copy is edited elsewhere.

In the target checkout's top directory, in one shell, check which `grove`
answers before `init`: the predecessor also has an `init`, which would write
its own scaffolding instead.

```sh
grove version     # must print "grove v…"; a usage error means the predecessor answered
grove init        # or: grove --project /absolute/path init
grove check
```

`init` ends with a note on stderr naming what the target's `AGENTS.md` or
`CLAUDE.md` should say if the entrypoints' defaults are not wanted: how
`grove` is invoked, and how work and proposal branches are named.

`init` writes `grove.yaml` (`records: grove`, `brief: grove/brief.md`), the
record root, a placeholder brief that states no intent, and the `grove-work`
and `grove-shape` entrypoints for Claude Code (`.claude/skills/`) and Codex
(`.agents/skills/`), each marked as managed. It prints one line per path:
`created`; `kept` for an existing `grove.yaml`, whose own `records` and `brief`
it then follows, for the brief and the record root, and for an entrypoint
without the marker, which is yours; `unchanged`; or `updated` for a marked
entrypoint whose template changed in the binary, which is the managed update.
A `grove.yaml` that is not a schema 3 configuration, or a directory, symlink,
or unreadable file at a managed path, is a conflict: init prints every reason,
writes nothing, and exits 1. It must run at the top of the Git checkout, and it
never reads or writes `AGENTS.md` or `CLAUDE.md`: the entrypoints defer to them
for how the CLI is invoked. Upgrading is building again and rerunning `init`.
Then `/grove-shape TOPIC` develops the brief and proposes work, and
`/grove-work G-001` carries it out, through the guides the binary prints.

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
guide does not repeat. Grove starts an agent only through `run` or the
board's `R`, one assigned work ID per attempt, as the headless form of this
skill; [the dogfooding evidence](grove/G-032-dogfood-review.md) says which
invocations have been exercised.

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

The columns (Proposed, Active, Review, Done) open on the current view: every
work record in its current state across all local branches and checkouts, as
`versions` decides it, the same from any checkout. An old copy on a stale
branch does not hide a later status elsewhere. Each card is a box holding the
record's ID and tag, its title, and its kind, size and priority, or for Done
the date it was last written and its candidate; the focused card has a heavy
border and a `▶` marker, and each column an accent colour that nothing
depends on. A card whose current state is only in a checkout's uncommitted
files is marked `uncommitted`. With a target, a card none of whose committed
current states is on it is marked `not on main`, and the header names the
target. Where the current states diverge, one card sits in the earliest of
their statuses, marked `⑂ 2 states`, and its detail says which states exist
and where. Work whose current state removes its record is listed under
Deleted. Done shows the most recently written cards that fit the column,
newest first, and counts the rest (`+ 23 older · / to search`). Abandoned is
hidden until `a` shows its column, and the shelf row counts it meanwhile.
Neither the bound nor the hiding moves a file.

`b` chooses between the current view and one checkout's own live files, as
before; this changes what is displayed and switches no branch or directory.
On a checkout's board, work with no live record there is listed under
Elsewhere without a status. Questions, decisions, pages, terms, plans and
reviews are never cards; `/` finds them.

Enter on a card opens the record's detail: a boxed header with the ID,
status, title, planning fields, candidate, standing against the target, the
places holding its current state and when it was last written; then its body
rendered from Markdown (headings, emphasis, lists, code, tables) beside a
sidebar of the records linked to it, the timeline of commits that changed
it, and one line per current state. Linked records are listed by role,
derived from fields alone: `plan` and `review` (records whose `work` names
it; a review adds what it `examined` and whether that is the candidate),
`work` for a plan or review, `needs` and `needed by` (`depends_on` either
way), `blocked by` and `blocks` (a question's `blocks`), `part of` and
`member` (`members`), and `related` (`relates_to` either way). Tab moves
focus from the content to the linked records, to the changes, to the
timeline, and back; ↑/↓ and PgUp/PgDn scroll the content or move the cursor.
Enter on a linked record opens its own detail, of any type, and Esc returns;
Enter on a timeline commit shows the record as it was at that commit, and
Esc returns to now.

Work in `review` is a candidate to judge, and its detail is the handoff: the
header adds a Review block (the candidate, whether it is approved, whether
only the record changed since it or the tip is a new candidate, whether the
target holds it, and which checkout each action runs in), the content opens
at its `## Evidence`, and the sidebar lists the candidate's changed files
against the target with their added and removed line counts. Enter on a file
shows its diff in the content pane, escaped like record text with added,
removed and hunk lines coloured, and Esc returns to the content. `a` asks for
a verdict and approves the candidate in the branch's checkout, `f` asks for
feedback and returns the work to `active` there, and `i` confirms the merge
into the target from the target's checkout, then asks whether to remove the
branch's worktree and branch. Each runs the same operation as the command,
one at a time; a result screen shows its facts, or why it was refused, and
the board is re-read. A candidate without a clean checkout of its branch, or
a target without one, is reported instead; the board never creates a
checkout. Below 100 columns the detail shows one pane at a time and Tab cycles
them. A page, term, decision, question, plan or review opens in the same
screen, with the fields its type has.

Work's detail also names its attempts (G-046): how many, and the latest
with its outcome. `R` on proposed or active work asks for a budget in USD,
then a permission mode, both typed each time since neither has a default,
and launches one attempt as `run` does; it runs on the branch the record's
current state stands on when that is not the target, in that branch's
checkout, so after feedback the next attempt continues on the candidate's
branch, and otherwise in a new `worktree-ID`. It is refused up front for
work in review, done or abandoned, work an open question blocks, and work
with an attempt still running or orphaned; `run`'s own refusals follow, and
one more: the record in this checkout changed since the board read it. `A` lists the
attempts of the open work, or on the board every attempt, newest first, and
Enter opens one: its outcome, the facts `attempt` prints, the provider's
final report rendered like a record body, and its recent activity newest
first, one short line per event from the last 1 MiB of its events, so a
flood of output costs one bounded read. The outcome is derived, never
written: `running`, `orphaned`, `interrupted`, or for a finished attempt a
`candidate ready` (only when the attempt ended with its record committed
in review with a candidate, which the owner records at exit), `stopped`, `failed` (no result
event, an error result or a nonzero exit), `waiting on question` (the work's
latest attempt, while an open question blocks it) or `ended without a
handoff`; a clean exit alone is never ready, and a record that says review
without being committed is reported as such. `x` asks, then stops a running or orphaned attempt as
`stop` does, keeping its partial work; `o` opens its work record. The
attempts are files the board only reads: quitting leaves an attempt
running, and the next session shows the same one. While one runs, they are
re-read every 2 s, a running card is tagged `● running`, and when one ends
the board is re-read.

The timeline is the record's history from Git: the commits that changed its
file, newest first, each with its date, the status the record held at that
commit, its short ID, and its subject, following renames. It is the history
of the state the detail shows, named in its heading: in the current view the
first current state's branch or checkout, on a checkout's board that
checkout's HEAD. A checkout whose files differ from its HEAD gets a first
`uncommitted` row. Merges are not listed, so where the record's status is not
the newest listed commit's, a first `here` row gives it and says why. History
is read from Git when a detail is open, never while the board loads, and no
key waits for a read still in progress. It says what happened on one branch
and nothing about whether another branch contains it.
[G-030](grove/G-030-card-lineage.md) owns this.

`v` in a detail opens the record's versions. A version is the record's exact
content; the board read it at every local branch's tip and in every
checkout's files, and lists each differing content once with its own title
and status, current ones first. An older one is marked `older`, and its
details say why. Where several branches and checkouts hold the same content
the row is a fold (`▸ done  same on 4 branches, 4 checkouts`): Enter lists
those places, and selects nothing. Opening a card or a detail selects nothing
either. Moving to one branch or checkout and pressing Enter asks
`workspace`'s resolver about exactly that version's selector; on success the
board closes and prints what `workspace` prints (the project path on stdout,
or its JSON with `--json`; checkout, branch, record, and revision on stderr).
A refusal (the version changed, its checkout is missing or ambiguous, the
record was deleted there) stays on screen with its reason until `r`
refreshes, after which a version must be selected again. The details pane
there begins with the focused version's history. Selecting never creates a
worktree, edits a record, or starts an editor, shell, or agent; only `R`
starts an agent, behind its prompts.

`/` on the board searches every record of the project in its current state,
of every type, including hidden Abandoned work, Done beyond its page, pages
and terms: typing filters by ID, type, status and title, not the body; ↑/↓
move, Enter opens the detail, Esc closes. Letters typed there filter rather
than act; Ctrl-C still interrupts.

Record bodies are rendered by glamour after the same escaping as every other
text, so a body's control sequences show as text; an HTML character reference
such as `&#x1b;` or `&amp;` is never decoded, since Markdown would decode it
after the escaping, and shows as typed in prose and code spans, with an extra
`&amp;` in code blocks and link targets; and of what the renderer emits only
its own styles reach the terminal. Links are shown as text, never as terminal hyperlinks, and a
relative target is shown root-relative (`/G-093-….md`). The render is cached
per record content and width.

Keys: arrows or `h` `j` `k` `l` move; Tab switches between the columns and
Deleted or Elsewhere, between the detail's panes, or between versions and
details; `/` searches; `a` shows or hides Abandoned on the board, and approves in the
detail of work in review, where `f` gives feedback and `i` integrates; `R`
launches an attempt of work, `A` lists attempts, and `x` stops one; `v`
opens a detail's versions; PgUp/PgDn scroll; `s` lists every branch and checkout read with its
diagnostics, which stay reachable while a banner marks an incomplete result;
`r` re-reads; Esc goes back, and quits from the board; `q` quits. Below 100
columns one status column shows at a time; below 40x10 the board asks for
more room. Leaving without a selection prints nothing and exits 0; Ctrl-C
exits 1; a usage error exits 2.

The board draws on stderr and reads stdin, so both must be terminals, while
stdout may be redirected: `cd "$(go run ./cmd/grove)"`. Without a terminal it
refuses at once with exit 1 and names the noninteractive commands. Text from
records, paths, Git, and an attempt's provider is shown with control
characters escaped. Besides its prompted actions it writes no files, including
the framework's debug logs.

Each command's contract, and the evidence and limits behind it, belongs to the
work record that delivered it; `/` on the board or `grove list` finds them.

## Develop Grove

[AGENTS.md](AGENTS.md) holds the development and verification policy. The
board's terminal checks drive the built binary through a pseudo-terminal with
`python3 internal/tui/testdata/terminal.py BINARY` (Unix; skipped without
`python3` and under `-short`). Run `lefthook install` once per clone:
pre-commit formats staged Go files and runs `go vet` and `go mod tidy -diff`;
there is no pre-push hook. `just clean-merged` removes local branches merged
into `main` and their clean worktrees after asking.

GitHub Actions runs gofmt, `go vet`, `go mod tidy -diff`, `go build ./...`,
`go run ./cmd/grove check` and `go test -count=1 -timeout 120s ./...` on
Ubuntu and macOS, and `govulncheck` on Ubuntu, for every push to `main` and
every pull request (`.github/workflows/ci.yml`); it is a signal, not a gate.
Dependabot opens weekly grouped update PRs for Go modules and actions.
