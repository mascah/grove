# The project as JSON

## Problem

- `grove changes --json` prints the target's changes and nothing else
  (changes spec at 8174594, "Print the project as data"). `grove serve`
  computes far more for its page: each cycle with the verdict `grove run`
  would give it and its members in build order, the branches holding work
  in progress with each task's status there, the recorded attempts, the
  history of the record and the problems `grove check` would report
  (serve spec at 8174594, "The project value"). Only a browser on the same
  machine can read that, at `/api/project`, while a server runs.
- The guides "Push to a tracker" and "Run Grove in CI" at 8174594 project
  the changes into GitHub issues from `grove changes --json`. Nothing lets
  a script project cycles, who is building what, or the attempts, so a
  tracker can mirror approval but not progress or claims.
- The owner wants trackers to be a projection of Grove's record for
  people, never a source Grove reads (this conversation, 2026-11-08). With
  [261108-jpg8n](../261108-jpg8n-claim-replaces-run/proposal.md), the
  attempt records live on the branch `grove/runs` and are the same from
  every host, so one export of the project can say who claimed what.

## Intended outcome

`grove changes --json` prints the one JSON object `grove serve` answers
`/api/project` with: the keys it prints today, unchanged in shape, and
beside them the time of the read, the maintained specs, each change's
cycle and readiness and files, the cycles with their verdicts, the
branches with work, the recorded attempts, the history of the record and
the problems. The serve spec names that object as what `grove changes
--json` prints. One command, run in any checkout after a fetch, feeds any
tracker projection as a script, and Grove writes to no tracker.

## Scope

### In scope

- `grove changes --json` printing the project value, and the serve spec
  saying `/api/project` answers with the same object.
- Documentation of the value where "The project as data" is explained.

### Out of scope

- A new command. Rule 8 wants a command only where a documentation step
  needs it, and the step here is a flag that exists.
- Changing the shape of any key `grove changes --json` prints today.
- Fetching from `origin`; the owner fetches, as for `serve`.
- Grove writing to a tracker, and a tracker-projection guide, until a
  recipe has been trialed (rule 8).

## Decisions

The owner accepted the round in this conversation on 2026-11-08 with
"yes".

1. **One value, no new command.** "`grove changes --json` gains `cycles`,
   `branches` and `runs` with the shapes `serve`'s `/api/project` has
   today, existing keys unchanged, and the serve spec names them as one
   value. No new command." The scope the owner confirmed in the same
   round says the command carries "the whole project value `serve` already
   computes", so the time of the read, the specs, the history and the
   problems come with it. Reason: the value exists; printing it is the
   whole change, and two shapes of the same facts would drift.
2. **A change of its own, depending on 261108-jpg8n.** Reason: it can be
   agreed and built alone, and its `runs` key reads the branch that change
   introduces.
3. **Routine defaults**: test layout and private helpers are delegated to
   the implementer.

## Open questions

None.

## Related documents

- Requirements: [changes](specs/changes/spec.md), [serve](specs/serve/spec.md)
- Tasks: [tasks.md](tasks.md)
