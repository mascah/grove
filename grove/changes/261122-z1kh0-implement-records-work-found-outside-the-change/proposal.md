# Implement records work it finds outside the change

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

An implementing agent often sees work the change does not ask for: a bug
next to the code it is changing, a document that is already false, a
missing test. Grove gives that work nowhere to go.

- The implement instructions say "Build what it says, nothing more"
  (`internal/instruct/requires/implement.md` at 15816d7). Their lines for
  `implementation.md` are `Status:`, `Check:` (a choice for the reviewer),
  `Answer:` and `Also changed:` (a file the task needed or made false).
  None of them is for work outside the change.
- Their finish line has the report list every `Check:` line, each blocked
  task's question and each unverified scenario, "then what was done and
  what is left". Anything else the agent saw is in that chat report at
  most, and it is lost when the session ends. `docs/concepts.md` at 15816d7
  opens with this problem: "the reasons for each change are settled in a
  conversation that ends when the session does."
- Agents fit such findings into `Check:` lines today. In
  `grove/changes/archive/261122-7vht2-claim-a-change-by-name/implementation.md`
  T-03 says: "Check: `docs/planning.md` still says to send one change on
  its own by planning a cycle of one; that remains true, for a timer, and
  was left as it is." A `Check:` line asks the reviewer to judge this
  change, not to plan another.

The beads issue tracker (github.com/gastownhall/beads,
`docs/core-concepts/dependencies.md`) gives found work its own link,
`discovered-from`, and tells agents to file it rather than keep it in
Markdown TODO lists. This change takes the idea in Grove's shape: a line
in the file the agent already writes, and the existing tracker guide.

## Intended outcome

When the agent sees work outside the change that no task needs, it does
not do it. It writes a `Found:` line under the task it was building,
saying what is wrong or missing, where, and how it saw it. Its report
lists those lines after everything the reviewer must judge. The lines are
merged with the code and stay in the archived change. The push-to-a-tracker
guide turns each one into an open issue, so it reaches whoever plans next.
Grove itself reads nothing from them.

## Scope

### In scope

- The implement instructions, and the implementation template and guide:
  when to write a `Found:` line, its form, and where the report lists it.
- The push-to-a-tracker guide: one issue per `Found:` line.
- The docs and the walkthrough that describe what implement writes and
  reports.

### Out of scope

- Grove reading `Found:` lines. `grove changes`, its `--json`,
  `grove serve` and `grove review` stay as they are.
- Explore or propose reading `Found:` lines. They reach the next session
  through the tracker, or through the owner.
- Grove writing to a tracker itself (`CONSTITUTION.md` rule 11).
- Fixing a found item in the same change. That is work the change does
  not ask for. A file the task needs is still an `Also changed:` line, and
  a wrong approach is still a revision of the change.

## Decisions

The owner accepted each recommended answer of the round with "yes".

1. **Form and place.** A `Found:` line goes under the task being built when
   the agent saw it: `Found: <what is wrong or missing, where, and how it
   was seen>`. It is not a section of its own. Reason: the task gives the
   line its context, and the template already holds each task's lines.
2. **What counts.** Anything outside the change that the agent did not
   change: a bug, a false document, a missing test, or code that differs
   from a maintained spec. The agent changes nothing for it. Something a
   task needs stays an `Also changed:` line or a revision, as today.
   Reason: the line records work for later. It must not become a way to
   widen the change.
3. **Grove reads nothing from it.** It adds no field to `grove changes`, its
   `--json`, `grove serve` or `grove review`, and `CONSTITUTION.md` is not
   amended. Reason: rule 3 lists what Grove reads from `implementation.md`,
   and nothing documented needs Grove to read this line.
4. **Report order.** The report lists `Found:` lines after every `Check:`
   line, blocked task's question and unverified scenario, and before what
   was done. Reason: the reviewer judges this change first. Found work is
   for later.
5. **Tracker guide.** The guide's agent opens one issue for each `Found:`
   line in the target's `implementation.md` files.
   - The issue is labelled `grove-found`.
   - Its body is the line and a link to its change.
   - The marker `<!-- grove-found:<change-id>:<T-NN>:<n> -->` keeps the
     agent from opening a second issue. `n` is the line's place among that
     task's `Found:` lines.
   - The agent never closes these issues; the owner does.

   Reason: found work outlives its change, whose own issue the guide
   closes when the change is archived.
6. **No `depends_on`.** This change does not wait on 261122-m78w6, although
   both change `docs/getting-started.md`, `docs/concepts.md`,
   `docs/reference.md`, `dev/walkthrough.md` and
   `internal/cli/walkthrough_test.go`. Reason: they change different text.
   Whichever is built second builds on the files as the target then holds
   them.
7. **Routine choices.** Two tasks. The walkthrough step is the fixture's
   traceback for an unknown ID, seen while building `ledger delete`. The
   only change spec is `implement`. A glossary entry is added for "Found
   line". Test layout and exact wording are delegated.

## Open questions

None.

## Related documents

- Requirements: [specs/implement/spec.md](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
