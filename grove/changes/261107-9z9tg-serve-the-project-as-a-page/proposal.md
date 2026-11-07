# Serve the project's state as a local page

## Problem

An owner who wants to know where a project stands reads several things
and puts them together by hand. `grove changes` and `grove changes --json`
list the approved changes and their cycles
([changes](../../specs/changes/spec.md) at e89dd0d, "Print the project as
data"), `grove check` reports a stale cycle only as a problem, and `grove
run` says why it cannot build a cycle only when it is run
([runs](../../specs/runs/spec.md) at e89dd0d). None of them shows which
branches hold work in progress, how far each change's tasks are, which
attempts `grove run` recorded in this checkout, or what the target's
history did to the record. The reader that `grove` opens shows the
maintained specs and approved changes and nothing of the rest
(`docs/reference.md` at e89dd0d, "Terminal commands"). The owner reads
`git branch`, `git log` and three commands instead.

## What the prototype found

A throwaway experiment on the branch `prototype/serve-page`, commits
d6e9fb1 and 545ae93, not pushed, asked whether one hand-written HTML page,
embedded in a Go `net/http` server run from the checkout with no build
step, shows readiness, work in progress by branch, history, the dependency
graph and a change's documents well enough that a front-end framework is
not needed. It was run on Grove's own repository (41 archived changes, 9
done cycles, 134 unmerged branches) and on a synthetic ledger fixture (12
changes, 3 cycles, a draft, a run branch, a teammate's branch only on
origin, an uncommitted edit), and driven through headless Chrome.

What was observed:

- All five views render on both repositories with no console errors, in
  light and dark mode. Readiness shows each open cycle with the verdict
  `grove run` would give it and its members in build order; the fixture's
  marked cycle reads "grove run would build it", the next names the member
  that needs a person and marks two stale ones.
- Telling a branch's work by content was wrong: comparing a branch's files
  with the target's current copies listed 84 of Grove's branches as in
  progress when every one had been squash-merged and overtaken. The owner
  said: "the in progress tab on the grove project is weird there's a ton
  of items on the page but not much actually in progress." Telling it by
  blob, with one `git rev-list --objects` over the target's history under
  the root, so that a file counts only when the target never held that
  version, listed 0 of Grove's branches and kept the fixture's draft,
  three builds in progress, the teammate's branch and the uncommitted
  edit. A local branch and its origin twin at the same tip collapsed to
  one row.
- History from one first-parent `git log --name-status` pass marks each
  commit that approved, revised or archived a change, planned or
  re-planned a cycle, or edited a maintained spec. A "touched" mark, for
  a target commit to a file a task names while the change is open, fired
  on typo fixes too.
- Documents as plain text were not readable enough; rendered through
  goldmark, which Glamour already depends on, with raw HTML dropped, they
  were, and relative links between specs could open in the pane.
- A page that rebuilt its document pane on every refresh without fetching
  again showed "loading…" indefinitely; keeping the loaded document across
  refreshes fixed it.

| Measure | Value |
| --- | --- |
| Grove repository, first request | 0.9 s |
| Grove repository, cached by branch tip | 80 ms |
| Fixture, first request | 0.2 s |
| Go in the prototype | 713 lines, about 150 repeating `run` and `check` |
| HTML and script | 415 lines |

Verdict: a plain page served from the checkout covers the scope, and a
product version that shares `run`'s and `check`'s logic fits within the
budget of rule 12.

## Intended outcome

An owner runs `grove serve` in a checkout and opens one local URL. The
page shows what the target branch holds and what the checkout knows
around it: what `grove run` would do with each cycle and why, which
branches and the working tree hold a draft or a change being built and
how far its tasks are, the attempts recorded in this checkout, the
dependency graph of changes and cycles, the history of the record on the
target, the maintained specs, and any of these documents rendered in a
side pane. It refreshes itself every 10 seconds, reads Git and files, and
writes nothing.

## Scope

### In scope

- `grove serve [--addr <host:port>]`: a local HTTP server on loopback,
  run from a checkout, printing its URL and running until interrupted.
- One JSON value the page reads, built from the target branch, the
  checkout's branches, the working tree and `.git/grove/`: specs, changes,
  cycles with their verdicts, branches with work the target lacks, run
  attempts, history, problems.
- The page: readiness, in progress, graph, history and specs views, a
  document pane rendering Markdown and showing YAML, polling every 10
  seconds.
- A documentation page, a reference row, a walkthrough step and the
  end-to-end test that follows it.

### Out of scope

- Anything hosted, any authentication, any build step or front-end
  framework.
- Fetching from origin or pushing; the page shows what the checkout has.
- Editing a document from the page, or any write at all.
- Extending `grove changes --json`, which keeps printing what it prints
  today; scripts wanting cycles, branches or history can read the server's
  value.
- A tmux dashboard, which rule 11 keeps out.
- The constitution amendment this change needs (rules 7, 12 and 16),
  which is its own pull request with no code, merged before any code of
  this change is written; see Decisions.

## Decisions

Each was put to the owner as a numbered recommendation on 2026-11-07 in
the conversation that proposed this change, and the owner answered
"yes" to the round.

1. **The constitution changes first, in its own pull request.** Rule 7
   gains: "`serve` is a local HTTP server on loopback that the owner runs
   from a checkout; it reads Git and the working tree, holds nothing and
   writes nothing, and is not a service." Rule 12 gains a budget of 500
   lines for the page's HTML and script, counted by `just check`. Rule 16
   adds `github.com/yuin/goldmark`, already in the module graph under
   Glamour. Rule 11 stays as it is. Reason: rule 7 says "no service" and
   rule 16 lists every dependency; the constitution says its own change
   is a pull request with no code that the owner reads before any code
   follows.
2. **One change, not two.** `grove changes --json` keeps printing what it
   prints today; the page's value is documented by the `serve` spec and
   built by one Go function and one struct that extends the changes value.
   Reason: splitting only pays if `changes --json` should gain the new
   keys for scripts, which the owner declined for now; it is reversible
   later.
3. **All five views as prototyped, minus the history's "touched" mark.**
   The marks that stay are created, approved, revised, archived, planned,
   re-planned and spec. Reason: "touched" fired for any target commit to a
   file a task names while the change is open, typo fixes included.
4. **Work in progress is told by blob.** A branch not merged into the
   target, or the working tree, counts when a file under the root has a
   blob the target's history never held; a change is a draft when the
   target lacks it and in progress when the target approves it; a branch
   whose only such changes are archived on the target is left out, with
   no toggle to show it; a local branch and its origin twin at the same
   tip are one row. Reason: the prototype's content comparison listed 84
   landed branches on Grove's repository; by blob it listed none and lost
   nothing on the fixture.
5. **Documents render as Markdown** through goldmark with raw HTML
   dropped, YAML as text, relative links between documents opening in the
   pane, served only for `.md` and `.yaml` files under the root at the
   target, at a branch tip, or in the working tree. Reason: plain text was
   not readable enough in the prototype; goldmark is already in the module
   graph; dropping raw HTML keeps a teammate's file from scripting the
   page.
6. **The command is `grove serve [--addr <host:port>]`**, default
   `127.0.0.1:7070`, printing the URL and running until interrupted, never
   opening a browser and never fetching from origin. Reason: a fixed
   default makes the URL bookmarkable; loopback keeps it local; opening a
   browser and fetching are lines the owner can type.
7. **The code lives in `internal/cli/serve.go`** with the page embedded
   beside it as `internal/cli/page.html`, so the eligibility and stale
   logic in `run.go` and `check.go` is reused without exporting it, and
   the branch scan and history go in the same file; Git keeps running only
   through `internal/gitx`. Reason: one file and no new exports; the
   prototype's 150 repeated lines are the lines this saves.
8. **Documentation is a new page `docs/watching.md`**, "Watching the
   project", in the nav after "Running unattended", with a row in the
   reference and a walkthrough step that the end-to-end test follows on
   the fixture. Reason: rule 8 wants a documentation step for every
   command, and rule 17 wants the test to follow the walkthrough.
9. **Delegated to the implementer:** test layout and helpers inside
   `internal/cli/walkthrough_test.go`, and the JSON field names.
10. **Routine defaults,** accepted in one line: spec ID `serve`; the page
    polls every 10 seconds with a checkbox to stop; the value holds the
    root, target, commit, time of the read, specs, changes, cycles,
    branches, runs, history and problems; branch rows sort newest first;
    attempt records come from `.git/grove/*.json` as `run` writes them.

The change depends on no other: the project index at e89dd0d lists no
approved change.

## Open questions

None. Nothing was deferred in the round.

## Related documents

- Requirements: [serve](specs/serve/spec.md)
- Tasks: [tasks.md](tasks.md)
