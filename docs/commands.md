# Command reference

`grove --help` gives every command's usage. This document owns what it does
not say about the commands below. The commands over records (`list`, `show`,
`brief`, `check`, `new`, `update`, `convert`), project discovery and exit
codes belong to the record model: its
[reading and writing](record-model.md#reading-and-writing-records),
[configuration and discovery](record-model.md#configuration-and-discovery),
[conversion](record-model.md#identity-and-placement-apart-from-classification)
and [brief](record-model.md#knowledge-records-and-the-brief) sections. `approve`,
`feedback` and `integrate` belong to its
[work lifecycle](record-model.md#work-lifecycle). [The board](board.md) has
its own document. Each command's acceptance, evidence and limits belong to
the work record that delivered it.

## Versions

`versions [ID] [--json]` shows one row per version of each record: its
committed version on every local branch tip and its live version in every
registered worktree, grouped by ID with each source's own status, so main can
see a feature branch's progress without switching or merging. Live rows say
how the file compares with that checkout's HEAD (`unchanged`, `modified`,
`renamed`, `added`, `deleted`). Every row ends with a selector that binds the
repository, source, commit, configuration, record path, and content revision;
identical bytes in two sources get two selectors. Sources are listed on
stderr with any diagnostics; an invalid or unreadable source makes the result
incomplete and the exit code 1 while valid sources still print. A live
project must be inside its registered checkout: a project location reached
through a symlink, or belonging to another nested repository, is invalid,
while a missing one is simply absent.

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
chain, is also older than it, and a note says so. A branch whose current
state removes the record gets a committed `deleted` row. Dates, status order,
and branch names never decide, so the answer is the same from every checkout,
and no branch is special. [G-042](../grove/G-042-current-view.md) owns this.

With an integration target named in `grove.yaml` (`target: main` here),
stderr names it and the `TARGET` column says whether that branch holds the
version's bytes (`yes`, `no`; `-` without a target). The target labels only
and decides nothing. `--json` adds each version's exact source text,
`current`, `older` (why it is older), and `on_target` (null without a
target), each record's `notes`, and the result's `target` and target
`notes`. The command writes nothing: no refs, index, worktrees, records, or
coordination state.

## Workspace

`workspace --source SELECTOR [--json]` takes one selector from `versions`,
checks that the version is still exactly what was selected, and prints the
absolute project directory of the existing checkout holding it, ready for
`--project`:

```sh
grove --project "$(grove workspace --source SELECTOR)" show G-003
```

A live selection resolves to its worktree. A committed selection resolves to
the one worktree that has that branch checked out, provided its live copy of
the record still has the committed bytes; otherwise, or when no or several
worktrees hold the branch, it refuses and says to run `versions` again and
select a live version. Any change since selection (branch or HEAD moved,
worktree moved or removed, attached or detached state, configuration, record
path, or content) refuses with the reason, and the selected checkout is
re-read once more just before its path is returned. The command never creates
a worktree, switches a branch, launches anything, claims ownership, or edits a
record; the directory it prints is a location, not write authority, and a
later `update` performs its own revision check. The board uses `versions` and
`workspace` in process, and selecting a version never creates a worktree for
a branch without one.

## Context

`context WORK_ID... [--json] [--interaction interactive|headless]
[--max-bytes N] [--include PATH]...` assembles context for explicitly selected
work in one checkout, so an assignment is a few IDs rather than a composed
prompt. [The work guide](work-execution.md) is the workflow that uses it. It
prints the IDs ordered prerequisites first (ties keep the requested order;
prerequisites are never added to the selection), Git identity, and two kinds
of content that it keeps apart:

- **Sources, read in full**, each once with its `sha256:` revision:
  `grove.yaml`, the selected work records, and every `--include PATH`. That is
  the whole default. An include is a required project-relative file of any
  type; one file reached by several names (letter case, a hard link) is
  included and charged once.
- **Observations, listed and not read**: the transitive `depends_on`
  prerequisites, questions blocking the selected work or a prerequisite
  (resolved ones too), plans and reviews whose `work` names a selected ID, and
  the records the selected work names in `relates_to` or `members` or links
  to, each with title, status, path, revision, why it is listed, and whether
  its source is included (`source` is the `sources[]` path that holds it,
  which differs from `path` when it was included under another name). Every
  link in the selected records' bodies is listed with the project path it
  resolves to. Links are taken from parsed Markdown, so code, fences, images,
  and HTML are never links; Markdown escapes and entities are decoded before
  the URL is. A link is never opened, so a listed path is not checked, not
  even for existence. URLs, fragment-only links, other projects, absolute
  paths, and Git metadata are listed with that reason and no path.

Retrieval is staged with existing commands: `show ID` prints a listed record,
and `--include PATH` adds a listed file, such as the plan the record names as
current, with its revision. The command picks no plan or review by itself, and
it never follows links inside an included document. A listed record's revision
says which version was seen; a listing is not a reading of its constraints.

It refuses rather than degrade: unknown, duplicate, or non-work IDs, a
malformed link destination, an included file that is missing, a symlink, not
regular, or not UTF-8, a source that does not fit `--max-bytes` (default
262144, counting source bytes, at most 8388608), or a checkout that changed
while it was read. Nothing is truncated or summarized to fit. A refusal
prints nothing to stdout. Exit 0 means context was assembled, never that work
is ready, authorized, or integrated: a prerequisite's `done` is its recorded
status, and only one with a `candidate` claims a merge, where it was written.
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

## Dependencies

`deps [WORK_ID...] [--json]` shows how work depends on other work, from one
checkout's records, so the owner can choose what to do next or assign
together without asking an agent to reconstruct the order. Only `depends_on`
orders: `members`, `priority` and question `blocks` never do, and an absent
`depends_on` means no declared prerequisite, not readiness.

Without IDs it lists the unfinished work (proposed, active, review), one row
each with `GROUP`, `LAYER`, `NEEDS` (its `depends_on`), `UNLOCKS` (the
unfinished work that names it) and `DELIVERY` (below). Work in one group needs other work in it,
directly or through anything else; a separate group is unrelated. A layer is
one more than the deepest work of its group that it needs, so layer 0 needs
no unfinished work. Equal layers have no declared order, which is not
evidence that the work can proceed in parallel. The prerequisites of that
work which are not unfinished (done, abandoned) are listed after it.

With IDs it previews that selection: the IDs as given, their order as
[`context`](#context) prints it, every transitive prerequisite outside the
selection with the selected work that needs it, listed and never added, the
open questions blocking any of them, and the pairs with no declared order.
Unknown, duplicate and non-work IDs are refused as `context` refuses them.
The preview adds no work, starts nothing and authorizes nothing.

Each item's delivery is a fact, not its status alone: proposed and active
work awaits implementation; a candidate in review or done says whether HEAD
of this checkout and the target branch contain it, by Git ancestry, or that
Git cannot read it here; done without a candidate says its delivery is
unrecorded; abandoned work will not be delivered, and a note names the work
that still needs it. A candidate in review also says what merging it into
the target's current tip would do, performed with `git merge-tree` in
objects only: `integrated`, a fast-forward, a clean merge although the
target moved since the branch left it, or a conflict in named files, with
the target commit it read, stale once the target moves. A selection with
two or more candidates in review not on the target merges them in its
order, each onto the ones before, and a note names where the first conflict
lands and in which files; Grove states that order but never chooses one, and
a clean order is not evidence that the changes work together. Git before
2.38 cannot perform such a merge, and the delivery says the merge could not
be predicted. Notes also report what the current view (see
[Versions](#versions)) says about each listed record: an older version here,
diverging current versions, other `depends_on` elsewhere (never merged into
this checkout's), uncommitted changes here, a record that changed while it
was read, and an incomplete inspection, which prints everything, says so on
stderr and exits 1. `--json` prints `{checkout: {root, ref, head}, target,
selected, order, items[]: {id, title, status, revision, candidate, outside,
layer, group, needs, unlocks, needed_by, delivery, merge}, questions[]: {id,
title, blocks}, notes, merge_order}`, where `merge` is `{target, commit,
outcome, conflicts}` or null and `merge_order` lists `{id, target, commit,
outcome, conflicts}` up to the first conflict or is null, with `selected`
null for the overview. The command writes
nothing but the objects a merge writes, and no ref.
[G-161](../grove/G-161-dependency-view.md) owns it, and
[G-177](../grove/G-177-merge-prediction.md) the merge prediction; the board's
`g` shows the same interpretation ([Dependencies](board.md#dependencies)).

## Attempts

Grove starts an agent only through `run` or the board's `R`, one assigned
work, or one explicit selection of work, per attempt. `run ID...
[--dry-run | --expect DIGEST] [--budget USD] [--permission-mode MODE]
[--until plan] [--model MODEL] [--effort LEVEL] [--branch NAME]
[--worktree DIR]` starts one bounded implementation
attempt of proposed or active work as a Grove-owned `claude -p "/grove-work
ID... --interaction headless"` process that outlives the terminal
([G-101](../grove/G-101-attempt-mechanism.md),
[G-045](../grove/G-045-durable-attempt.md)), the IDs passed as given. It creates `worktree-` plus the IDs joined by `-`
(`worktree-G-030`, `worktree-G-030-G-031`) under
`.claude/worktrees/` from this checkout's HEAD, or reuses the branch's
registered worktree so a next attempt continues from preserved partial work,
then starts an owner process in its own session that runs the provider there
with `--output-format stream-json`, `--max-budget-usd`, `--permission-mode`
and `--permission-prompts none`, its stdout and stderr written straight to
files under the Git common directory (`.git/grove/attempts/ATTEMPT/`, shared
by every worktree, never committed). Budget and permission mode are required,
from the flags or from the launching checkout's `grove.yaml` `run:` defaults
([record model](record-model.md#configuration-and-discovery),
[G-140](../grove/G-140-default-an-attempt-s-budget-mode.md)), which may also
set the model and effort; a flag overrides its default for one launch, and
`attempt.json` records the resolved values alike. Without either source `run`
is refused as a usage error: Grove itself sets no default spend or profile.
Grove starts one process and never retries; subagents the provider starts
share the budget.

Several IDs are one selection
([G-162](../grove/G-162-bounded-work-selection.md),
[G-188](../grove/G-188-selected-work-shared-candidate.md)): still one
process, one worktree and one budget over all of it, never one process per
ID. `run` orders the members as [`deps`](#dependencies) does and adds
nothing: a prerequisite outside the selection is listed with its delivery at
the base (the launching HEAD, or the reused branch's tip), never implemented.
A member **waits**, and the agent does not start it, when an open question
blocks it, when an outside prerequisite is not delivered at the base
(proposed, active, review or abandoned, or done with a candidate the base
lacks; done without a candidate is reported as unrecorded delivery), or when
a selected prerequisite waits. One ID is a selection of one, so work whose
prerequisite is undelivered is refused too. Bounded by `--until plan`, only
an open question stops a member: a plan needs its prerequisites named, not
delivered. The agent implements the members one at a
time in that order; a new question, an outside blocker or a failure stops
that member and every member that needs it while the rest continue, and
budget exhaustion or `stop` ends everything with what is committed kept. The
complete members enter review together on one shared candidate, judged per
member and integrated as a group ([Work lifecycle](record-model.md#work-lifecycle));
a started member left incomplete holds the whole branch out of review,
since integrating it would carry the unfinished code, while one that never
started keeps its wait and does not. `--dry-run` checks everything a launch
checks and prints the assignment without writing or starting anything: the
IDs, order, each member's status, revision and whether it can start or
waits and why, the outside prerequisites, the base, worktree, bounds, review
boundary, continuation policy, and a digest, sha256 over the IDs as given,
each member's revision, the base and the resolved options. `--expect DIGEST`
refuses a launch whose assignment no longer has that digest and prints what
it would run now. `attempt.json` records the selection with its digest and
member revisions; an attempt from before selections reads as a selection of
its one work.

Three options shape one launch, recorded in `attempt.json` and reported as
the attempt's `Requested:` fact ([G-134](../grove/G-134-bound-an-attempt-at-its-plan-and.md)).
`--until plan` adds the bound to the assignment (`/grove-work ID --until plan
--interaction headless`): the attempt stops at a committed plan with the
record's status as it found it and the continuation in its Next, as the work
guide's step 4 says, and launching again without the bound, after reading the
plan, is the implementation. `--model` and `--effort` pass through to the
provider, which owns their values (`claude --help`); a preparation attempt and
the implementation after it can differ in both. The launch also records the
sha256 of the worktree's `.claude/agents/grove-reviewer.md`, the reviewer
definition `init` writes and step 6 reviews through, or `none` where there is
none, which the launch warns of: without it the attempt has no independent
reviewer, and work whose record requires one stays active. It records the
sha256 of the worktree's `.claude/skills/grove-work/SKILL.md` as `skill`,
and in `differs_from_template` which of that skill and the reviewer are not
byte for byte the launching `grove`'s templates: a custom, older or newer
file keeps its own digest and never borrows the template's identity, and the
attempt's `Entrypoints:` fact says which (an attempt from before this says
they were not recorded). A skill the harness prefers over the worktree's,
such as a personal one of the same name, is not seen. `grove_version` is the
launching `grove`'s
[version](#version-and-guide) line. The `grove` the agent itself runs is not
recorded: `PATH` or the project's instructions choose it, and it need not be
the launching one. An attempt started
by hand, such as an interactive `/grove-work`, writes no attempt files: it is
visible only as its branch, its worktree and the checkpoint in the work's
Next.

`grove --help` lists what `run` refuses. The worktree holds only what is
committed, so `run` refuses a worktree without `.claude/skills/grove-work/SKILL.md`,
the skill the prompt names, which `init` writes and you commit
([G-150](../grove/G-150-launch-attempts-only-where-the-w.md)); a new branch
is checked in HEAD before it is created. It likewise refuses a worktree
whose marked `grove-work` skill or reviewer definition carries an
[entrypoint revision](#entrypoint-revisions) the launching `grove` does not
serve, legacy included, which would stop or contradict the guide after the
spend began
([G-169](../grove/G-169-harness-upgrade-compatibility.md)). An open
question that blocks the work is the wait the headless guide persists, so
rerunning with nothing changed refuses the same way: a selection none of
whose members can start is refused, naming each wait, whether the launching
checkout or the reused branch holds it. The launch prints the selection as
the preview does. A running or orphaned attempt whose selection shares any
member refuses the launch, so overlapping selections cannot both own a
record. After feedback reopens a group, a selection that leaves out a member
still sharing the candidate on that branch is refused: the group runs again
together, on the same branch (`--branch` names it when the IDs are given in
another order than the branch's name). Every refusal comes before a write, except
that what an existing branch holds is checked in its checkout, so a branch
that had no worktree keeps the one `run` made.

`attempts [ID]` lists attempts newest first, or those whose selection
includes ID, with every member in the WORK column. `attempt ATTEMPT [--json]`
prints one attempt's launch, what it requested, event counts (parsed bounded:
a line over 1 MiB is counted, not read), the provider's init fields with the
model that actually ran, the result event's fields with its cost split by
model where the provider reports one (a subagent on another model shows
apart), the result
(exit, the worktree's HEAD and whether it holds uncommitted or untracked
changes), the record as the branch holds it, whether the record on the target
changed since launch, and the file paths; no provider text is printed. For a
selection it prints each member at launch and, once ended, one line per member
as the branch holds it: awaiting judgment (review, with the candidate),
active with its checkpoint in its Next, not started and why, or waiting on a
question; and whether each member's record on the target changed.
Liveness is the owner's file lock, never a pid: `running` while it is held,
`finished` once `result.json` exists, `orphaned` when the owner is gone but
the provider's process group lives, `interrupted` when nothing is left and no
result was written (a machine restart reads so; no reboot recovery is
selected).

`stop ATTEMPT` sends SIGINT through the owner so the provider ends its turn,
SIGKILL to its process group after 15 s, and writes the result; an orphaned
attempt is stopped directly and reconciled. Stop touches neither the worktree
nor the record. A result is facts, never acceptance: the record's own status
on the branch, which the headless guide sets, is the handoff, and a process
exit or a `result` event proves nothing about it. `GROVE_CLAUDE` names
another executable, for fakes.

## Init

`init` runs at the top of a Git checkout (or `--project /absolute/path`). It
writes `grove.yaml` (`records: grove`, `brief: grove/brief.md`, and no `run:`
launch defaults), the record
root, a placeholder brief that states no intent, the `grove-work` and
`grove-shape` entrypoints for Claude Code (`.claude/skills/`) and Codex
(`.agents/skills/`), and Claude Code's `grove-reviewer` agent definition
(`.claude/agents/grove-reviewer.md`), each marked as managed. The reviewer is
read-only, inherits the session's model at `high` effort, and loads the
review guide the work guide's step 6 dispatches (`grove guide review`), as
the skills load theirs. Codex has no
equivalent, so a Codex session reports its missing independent reviewer as
step 6 says. It prints one line per path:
`created`; `kept` for an existing `grove.yaml`, whose own `records` and
`brief` it then follows, for the brief and the record root, and for an
entrypoint without the marker, which is yours; `unchanged`; or `updated` for
a marked entrypoint whose template changed in the binary, which is the
managed update. A `grove.yaml` that is not a schema 3 configuration, or a
directory, symlink, or unreadable file at a managed path, is a conflict: init
prints every reason, writes nothing, and exits 1. It never reads or writes
`AGENTS.md` or `CLAUDE.md`: the entrypoints defer to them for how the CLI is
invoked, and `init` ends with a note on stderr naming what they should say if
the entrypoints' defaults are not wanted: how `grove` is invoked, and how
work and proposal branches are named, and to commit what it wrote: an
attempt's worktree holds only committed files, and `run` refuses one without
the `grove-work` skill.

### Entrypoint revisions

The files `init` writes are thin: the work, shaping and review instructions,
including what an assignment or a shaping request may hold, are the guides
the binary prints, so installing another `grove` changes what a fresh
session follows without rewriting them. What an entrypoint needs of the
binary is its *entrypoint revision*, an integer apart from the release
version and the record schema, written in each managed file as
`grove entrypoint revision N` and passed as
`grove guide NAME --entrypoint N`. From revision 2 an entrypoint holds
nothing the guides evolve (the assignment grammar and the review brief are
the guides'), so a guide change never needs a new revision; one changes only
when an entrypoint needs something an older `grove` lacks, or a newer one
stops serving what an older entrypoint asks. This `grove` writes and serves
revision 2. What `init` wrote before revisions existed, marked files with no
revision line, is revision 1, `legacy`: each generation kept its own
assignment grammar and review brief (the earliest rejects `--until plan`,
and has no reviewer), and nothing in the file says which, so it is not
served. Refresh it with `init`.

`guide --entrypoint N` refuses a revision it does not serve, printing
nothing and saying how to repair it, and a `grove` from before revisions
refuses the option; either way the entrypoint stops at its first command,
before any work. A plain `guide` always prints, since a person reads it the
same way, and a legacy skill loads its guide that way; so the work and
shaping guides' Inputs tell a session loaded through a managed skill with no
revision line to stop and name the repair. A legacy reviewer loads no guide
and is reached only through `init --check` and `run`.

`init --check` diagnoses each managed path and writes nothing, one line
each: `current` (this `grove`'s template), `compatible` (a served revision
in other text, older or edited), `legacy` (no revision line), `incompatible`
(a stated revision this `grove` does not serve, newer or older, or not a
number), `missing`, `custom` (no marker: yours, listed and never judged), or
`conflict` (not a file `init` could replace). It exits 1 when any path is
legacy, incompatible, missing or a conflict. Different text alone never
makes a file incompatible.

Upgrading is installing the new `grove`, then in the target's checkout
`init --check`, `init`, and committing what it updated; `init` rewrites only
marked files and keeps unmarked ones, the configuration, the brief and the
records. A linked worktree holds its own branch's copy, which `init` in
another checkout never touches: run `init` in that worktree and commit it
there, or merge the target after the refresh is committed. A session that
already loaded an entrypoint or a guide keeps what it read, so start a new
session (or reload skills) after refreshing. `run` checks the worktree it
launches in against the launching `grove`, and refuses an unserved revision
before any spend; the `grove` the session itself runs (`PATH` or the
project's instructions) is not checked, so where it differs, the session's
first `guide` is what stops it.

Check which `grove` answers before `init`: the predecessor also has an
`init`, which would write its own scaffolding instead. Codex runs each
command through a login shell, so it sees the profile's `PATH`, not the
caller's: put the build directory on the login `PATH` ahead of any other
`grove`, or name the executable in the target's `AGENTS.md` or `CLAUDE.md`.
When the wrong `grove` answers, the entrypoints stop and say so rather than
act.

## Version and guide

`version` names the executable and the workflow content it ships, in one
line that `attempt.json` also records as `grove_version`:

```text
grove VERSION (COMMIT[, vcs OTHER][, modified]) guides sha256:GUIDES content sha256:CONTENT
```

- `VERSION` is the release stamp, else the module version Go recorded: for
  `go build` in a checkout, a tag or pseudo-version from that checkout, with
  `+dirty` for uncommitted changes; for `go install …@VERSION`, that tag or
  pseudo-version; `(devel)` for `go run` and a tree without `.git`; and
  `(version unknown)` without build information.
- `COMMIT` is the stamped commit, else Go's `vcs.revision`, else
  `commit unknown`, as `go run`, a build from a source archive without
  `.git` and a `go install` print it. Where Go recorded a different commit
  from the stamp, `vcs OTHER` follows it. `modified` is Go's `vcs.modified`:
  the checkout Go found had uncommitted changes, stamped or not. Build from a
  primary checkout or a clone: for a linked worktree that lies inside its
  repository, Go only recognises the enclosing checkout's `.git` directory
  and takes that checkout's pseudo-version, commit and cleanliness instead
  (observed with go 1.26.2 and 1.26.8), which a stamp then shows as `vcs OTHER`; the
  `content` digest still names the worktree's own files.
- `guides` digests the work, shaping and review guides and the record model,
  which `guide` prints. `content` digests the workflow the binary ships into
  a session or a project: those four and the harness entrypoints `init`
  writes, the reviewer definition among them. It leaves out the `grove.yaml` and
  placeholder brief `init` creates once and never rewrites. Neither digest
  describes a project's installed files, which may be custom or from another
  release; an attempt records those itself ([Attempts](#attempts)).

The predecessor rejects `version` as an unknown command, so the line tells
the two apart.

A release build stamps the version and the full commit through the linker,
with `-trimpath` so no machine path enters the binary; Go embeds no build
time. The same command serves a clone and an extracted source archive:

```sh
go build -trimpath -ldflags "-X github.com/mascah/grove.version=v0.1.0 -X github.com/mascah/grove.commit=$(git rev-parse HEAD)" -o bin/grove ./cmd/grove
bin/grove version        # grove v0.1.0 (<that commit>) guides sha256:… content sha256:…
go version -m bin/grove  # Go's own record: -trimpath=true, and vcs.* from a clone
```

`bin/` is ignored, so an artifact leaves the tree clean for the next one;
`-o grove` would write into this repository's `grove/` record directory.
From an archive, pass the commit the archive was made from instead of
`git rev-parse`. With `-trimpath`, Go leaves `-ldflags` out of the build
settings, so `version` is where the stamp is read. A stamp is a claim the
builder makes: `version` prints it as given, and only the builder's inputs
make it true. Without a stamp, the line falls back as above and never
presents a release. Every artifact of one release, whatever its platform,
prints the same version, commit, `guides` and `content`, with neither `vcs`
nor `modified`; any difference means one was built from other inputs.

`go install …@COMMIT` resolves only a pushed commit, and rebuilding an
installed binary is manual.

`guide work`, `guide shape` and `guide review` print the
[work](work-execution.md), [shaping](work-shaping.md) and
[review](work-review.md) guides the binary carries, so the workflow version
is the executable version and no copy is edited elsewhere. `--entrypoint N`
is how an entrypoint `init` wrote asks for one
([Entrypoint revisions](#entrypoint-revisions)). `guide model`
prints the [record model](record-model.md) the guides cite, the contract that
binary validates, so a session in any project reads it without Grove's
repository. Since it ships into projects whose own `G-` IDs are live, it links
to no record.
