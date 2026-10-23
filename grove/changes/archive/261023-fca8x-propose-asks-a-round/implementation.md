# Implementation: Propose asks a round

<!-- Progress on change 261023-fca8x, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Propose asks a round

Status: done

- `internal/instruct/requires/propose.md`: for a new change, once the
  scope is confirmed and before any document is written, one numbered
  round of every decision the agent would otherwise assume, each with a
  recommended answer worded so that "yes" accepts it; nothing to assume is
  said and written; each answer is a decision quoting the owner; a
  deferred question stays an open question, and nothing inferred is written
  as an assumption. In every schema, since it is what Grove requires.
- `schemas/default/skills/propose.md`: step 2 keeps confirming the scope;
  a new step 3, "Ask one round", lists what the round covers; writing is
  now step 4.
- The proposal guide and template drop "Assumptions to confirm in review":
  what is inferred is asked in the round or left as an open question.
- Docs: `docs/concepts.md` ("The loop", Propose), `docs/reference.md`
  (what a schema cannot change, and the bugfix example asking in the
  round), `dev/walkthrough.md` main-line step 3 proposing through a round.
- The explore guide keeps "One focused question at a time", unchanged.
- `internal/cli/walkthrough_test.go`, `TestPropose`: the propose
  instructions print the round and no "Assumptions to confirm". `just
  check` passes (3396 Go lines, 1140 test lines); `grove check` is clean.
- Not run: the acceptance's hand check on the ledger fixture in Claude Code
  (`/grove-propose` from an explore conversation); this unattended session
  has no owner to answer the round.
- A file no task names that looks like it should change, left as it is:
  `docs/getting-started.md` line 94, "It asks rather than assumes", does
  not mention the round.

## Reconciliation

- propose: applied to specs/propose/spec.md, which change 261023-6j5bw
  created on this branch
