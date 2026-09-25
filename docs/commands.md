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

## Attempts

Grove starts an agent only through `run` or the board's `R`, one assigned
work ID per attempt. `run ID [--budget USD] [--permission-mode MODE] [--until
plan] [--model MODEL] [--effort LEVEL] [--branch NAME] [--worktree DIR]`
starts one bounded implementation
attempt of proposed or active work as a Grove-owned `claude -p "/grove-work
ID --interaction headless"` process that outlives the terminal
([G-101](../grove/G-101-attempt-mechanism.md),
[G-045](../grove/G-045-durable-attempt.md)). It creates `worktree-ID` under
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
reviewer, and work whose record requires one stays active. An attempt started
by hand, such as an interactive `/grove-work`, writes no attempt files: it is
visible only as its branch, its worktree and the checkpoint in the work's
Next.

`grove --help` lists what `run` refuses. The worktree holds only what is
committed, so `run` refuses a worktree without `.claude/skills/grove-work/SKILL.md`,
the skill the prompt names, which `init` writes and you commit
([G-150](../grove/G-150-launch-attempts-only-where-the-w.md)); a new branch
is checked in HEAD before it is created. An open question that blocks the
work is the wait the headless guide persists, so rerunning with nothing
changed refuses the same way. Every refusal comes before a write, except
that what an existing branch holds is checked in its checkout, so a branch
that had no worktree keeps the one `run` made.

`attempts [ID]` lists attempts newest first. `attempt ATTEMPT [--json]`
prints one attempt's launch, what it requested, event counts (parsed bounded:
a line over 1 MiB is counted, not read), the provider's init fields with the
model that actually ran, the result event's fields with its cost split by
model where the provider reports one (a subagent on another model shows
apart), the result
(exit, the worktree's HEAD and whether it holds uncommitted or untracked
changes), the record as the branch holds it, whether the record on the target
changed since launch, and the file paths; no provider text is printed.
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
read-only, inherits the session's model at `high` effort, and carries the
standard review brief the work guide's step 6 dispatches; its text is Grove's
own `.claude/agents/grove-reviewer.md`, embedded in the binary. Codex has no
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
the `grove-work` skill. Upgrading is building again, rerunning `init` and
committing what it updated.

Check which `grove` answers before `init`: the predecessor also has an
`init`, which would write its own scaffolding instead. Codex runs each
command through a login shell, so it sees the profile's `PATH`, not the
caller's: put the build directory on the login `PATH` ahead of any other
`grove`, or name the executable in the target's `AGENTS.md` or `CLAUDE.md`.
When the wrong `grove` answers, the entrypoints stop and say so rather than
act.

## Version and guide

`version` prints the module version and, when the build stamped it, the VCS
revision with `modified` for a dirty tree; `go run` prints `(devel)`, which
means this checkout's files. Build from a primary checkout or a clone: for a
linked worktree that lies inside its repository, Go only recognises the
enclosing checkout's `.git` directory and stamps that checkout's revision and
cleanliness instead (observed with go 1.26.2). The predecessor rejects
`version` as an unknown command, so the line tells the two apart, and the
line ends with a digest of the embedded guides and record model, which names
the workflow even when no revision was stamped.

`go install …@COMMIT` resolves only a pushed commit, and rebuilding an
installed binary is manual.

`guide work` and `guide shape` print the [work](work-execution.md) and
[shaping](work-shaping.md) guides the binary carries, so the workflow version
is the executable version and no copy is edited elsewhere. `guide model`
prints the [record model](record-model.md) the guides cite, the contract that
binary validates, so a session in any project reads it without Grove's
repository. Since it ships into projects whose own `G-` IDs are live, it links
to no record.
