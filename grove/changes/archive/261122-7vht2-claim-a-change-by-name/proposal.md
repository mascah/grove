# Claim a change by name, and more than one piece of work at once

## Problem

`grove claim` claims only a cycle, and one attempt pauses the whole
project.

- An approved change that no cycle lists is never claimed
  (`grove/specs/runs/spec.md` at 07d319f, "Nothing to apply", scenario "An
  approved change in no cycle"), and `grove claim <id>` looks the ID up
  among cycles only (`internal/cli/claim.go` at 07d319f).
- While the latest record on `grove/runs` names a cycle with a member not
  archived, every claim is refused, named or not, from every checkout
  (runs spec, "Pause after an attempt" and the scenario "A given cycle
  while paused").
- `CONSTITUTION.md` at 07d319f says the same: rule 10 has `claim` claim
  "no cycle while the latest record on `grove/runs` names a cycle not yet
  done", and rule 11 puts "building two cycles at once" out of scope.

So before anything is built, someone plans a cycle and merges it, and for
a timer also merges the `unattended` mark; and nothing else starts until
that cycle's pull request is merged. Both merges schedule work that was
already approved.

tend, the runner the owner builds on Grove, is moving to a mode a project
opts into where agents plan and build without a person
([direction.md](https://github.com/mascah/tend/blob/main/docs/work/261121-autonomous-direction/direction.md),
2026-11-21). Its decision 8 reads: "Approval stays a merge; scheduling
stops being one. An approved change is ready when everything it waits on
is archived as finished on the target (Grove's existing meaning; finished
on an attempt's branch releases nothing) and nothing in flight is planned
to touch the same requirements or files. That overlap check only reduces
collisions, since an implementer may touch files its plan did not name;
the merge train (decision 9) is the guarantee. How many pieces run at
once, and which first, is a setting in the flow file." This change is
Grove's part of it.

What already fits: an attempt's record lists its changes, and the pause
reads only that list (`internal/cli/claim.go` at 07d319f, `attempt` and
`pausing`); a change's planned files are its tasks' `Files:` lines
(`CONSTITUTION.md` rule 3) and the capabilities it has change specs for
(`internal/doc/doc.go` at 07d319f, `Capabilities`); and each claim is
already pushed so that it succeeds only while `grove/runs` is where the
claim read it.

## Intended outcome

`grove claim <id>` takes a cycle or an approved change that no cycle not
done lists. It is refused only when what it names cannot be built, is in
flight, or shares a planned file with something in flight, so pieces of
work that share nothing can be claimed at once, each in its own checkout.
A change is in flight while a record names it and the target does not
hold it archived. `--retry` claims despite what is in flight.

`grove claim` given nothing does what it does today: it takes the oldest
cycle marked `unattended`, and waits while anything is in flight.

`grove changes --json` and the page `grove serve` shows say, for each
cycle and each change in no cycle, whether `grove claim <id>` would
claim it and why not, and which recorded attempts are in flight.

Grove still chooses nothing: which piece is claimed next, and how many
attempts run at once, are the caller's.

## Scope

### In scope

- `grove claim <change-id>` for an approved change in no cycle not done,
  and its refusals.
- A named claim refused for being in flight or for overlap, in place of
  the pause; the pause, for a claim given nothing, counted over every
  record; `--retry` for both.
- The attempt record for a change in no cycle.
- `claimable` and `verdict` for such a change, `claimable` for a cycle,
  and `in_flight` for a record, in `grove changes --json` and
  `/api/project`; the page's readiness and in progress tabs.
- The runs, changes and serve specs, the user documentation and the
  walkthrough. Beyond the files decision 13 lists, `docs/watching.md`
  and `docs/concepts.md` describe readiness, attempts and the pause, and
  T-03 brings them up to date.

### Out of scope

- `CONSTITUTION.md`: rules 2, 3, 4, 10 and 11 are amended in pull request
  #181, which holds no code. This change is built only after that merges.
- A cycle file, `/grove-plan` and the `unattended` mark: unchanged. A
  cycle's stale plan still stops any claim of it.
- A `grove claim` given nothing picking a change in no cycle: it never
  does. Which piece goes first is the caller's.
- A command that releases a claim, and a claim expiring by time: rule 11
  keeps recovery and expiry out of scope.
- Exit statuses that tell a refusal from a failure, a promise that the
  JSON only gains fields, and a `grove ready` command.
- A setting in `grove.yaml` for any of this.
- `grove changes`' text lines, and `grove check`'s reports.
- Reconciling a requirement rewritten on the target: change
  261122-6ndwb, which this one depends on.
- Removing the History line from the fold so that two changes on one
  capability stop colliding in Git.
- Worktrees, and everything in tend: its count of pieces at once, its
  merge train, and when it retries.

## Decisions

The owner settled the first six in an explore on 2026-11-21, one question
at a time, and answered the round before writing with "yes" to each
recommended answer.

1. What is claimed: a change, with the cycle kept as an optional group.
   The owner: "1, keep the cycle as an optional group". A cycle's open
   members are still claimed as one attempt, and a cycle stays the one
   place a change builds on another's unmerged work. Reason: across
   attempts, finished on an attempt's branch releases nothing, so without
   the cycle every dependency link costs a merge; and 11 of the 12 cycles
   in `grove/cycles/` at 07d319f list two to five changes.
2. No setting. The owner: "2, no key". The difference is in how `claim`
   is called: given nothing it keeps the pause, and given an ID it is
   refused only for being in flight or for overlap. Reason: a key would
   give the project two modes for every spec, verdict and test to
   describe; a project opts in through its runner's configuration, where
   the owner has put how work flows. The pause still protects the
   one-checkout timer of `docs/guides/run-on-a-timer.md`.
3. Staleness against `base` stops no claim of a change in no cycle. Asked
   who reconciles a change whose requirement was rewritten, the owner
   said: "during the build". Reason: two changes that overlap are built
   one after the other, so the second is stale as soon as the first
   lands, and refusing it would bring back a merge that only schedules.
   Change 261122-6ndwb has implement reconcile it, which is why this
   change names it in `depends_on`.
4. Overlap is by capability, not by requirement. The owner: "2, the same
   capability". A change's planned files are its tasks' `Files:` lines
   and the maintained spec of each capability it has a change spec for.
   Reason: every fold appends a line at the end of the maintained spec's
   History (`internal/instruct/schemas/default/skills/archive.md` at
   07d319f), so two changes on one capability collide in Git whatever
   requirements they touch; and one rule about files needs no separate
   comparison of requirements.
5. What lets go of a claim. The owner: "1, only what exists today". A
   change leaves flight when the target holds it archived, as finished or
   dropped, or when it is claimed again with `--retry`. Reason:
   `CONSTITUTION.md` rule 11 keeps recovery and expiry out of scope, and
   no documented step needs a release (rule 8).
6. What a runner reads and calls. The owner: "1, extend what exists, no
   extras." No new command, no new exit status, no promise about the
   JSON.
7. A `grove claim` given nothing is paused while any change is in flight,
   not only a change of the latest record; `--retry` still claims
   regardless. Reason: with several attempts open, an older one still in
   the timer's checkout would otherwise go unnoticed.
8. In the JSON, a cycle's `runs` keeps today's meaning, false while
   anything is in flight; `claimable` is new, on each cycle not done and
   each change in no cycle, and says whether `grove claim <id>` would
   record an attempt now; such a change also has `verdict`; each record
   has `in_flight`. Reason: today's tend picks a cycle by `runs`, so it
   keeps building one at a time with no change to it.
9. `grove claim <change-id>` for a member of a cycle not done is refused,
   naming the cycle to claim.
10. A cycle's stale plan still stops both forms of `claim`.
11. `--retry` on a named claim skips only the two new refusals. A change
    that waits, a member of a cycle, and a cycle that waits or is stale
    are still refused.
12. Planned files and the overlap test go in `internal/doc`; in flight
    and the refusals in `internal/cli/claim.go`, in place of `pausing`;
    the new fields in `internal/cli/serve.go`; the page in
    `internal/cli/page.html`.
13. Routine, accepted together: a record for a change in no cycle leaves
    `cycle` out; open questions and blocked tasks stop no claim, as for a
    cycle's members today; overlap reads only the fetched target's
    documents, never an attempt's branch; `grove changes`' text lines are
    unchanged; the page's readiness tab lists each change in no cycle
    with its verdict and its runs table marks records in flight;
    `docs/running.md`, `docs/reference.md`, `docs/coordinating.md`, the
    timer guide and `dev/walkthrough.md` are updated; private helpers,
    test names and layout, and the exact wording of one-sentence messages
    are the implementer's to choose.

## Open questions

Can wait: whether this fits the budget of `CONSTITUTION.md` rule 12
without raising it. `just check` at 07d319f counts 4,544 of 5,000 Go
lines, 2,360 of 3,000 test lines and 410 of 500 page lines. If T-01 or
T-02 cannot pass `just check` for the budget alone, that task is blocked
with the count, since raising the budget is the owner's change to the
constitution.

## Related documents

- Requirements: [runs](specs/runs/spec.md),
  [changes](specs/changes/spec.md), [serve](specs/serve/spec.md)
- Tasks: [tasks.md](tasks.md)
