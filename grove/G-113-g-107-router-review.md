---
id: "G-113"
type: review
title: "G-107 router restructure review"
status: current
created: "2026-09-23T18:34:50Z"
updated: "2026-09-23T18:35:06Z"
work: ["G-107"]
examined: "833ac1b"
---

## Examined

The second attempt at [G-107](G-107-current-documentation.md), after the
owner's feedback on candidate `71a650e`, against its revised acceptance 3 to
5 and the revision section of plan [G-111](G-111-g-107-docs-plan.md).
Reviewer: a fresh `reviewer` subagent in the implementing Claude Code session
(Opus 5.5), read-only, with no part in writing the diff. It is an agent, not
the owner. Three rounds, which is the cap:

- Round 1: `git diff 4159e79 49c8f5b`.
- Round 2: `git diff 49c8f5b e07f98e`, the fixes.
- Round 3: `git diff e07f98e 833ac1b`, the second fixes.

Each round checked links and anchors with a script and ran
`go run ./cmd/grove check`. Round 1 compared the old README paragraph by
paragraph with the new owners, and checked the old AGENTS.md rule by rule.
It also checked the board's key table against `internal/tui`.

## Findings

Round 1 found nothing high-severity.

1. Medium: AGENTS.md said "change status or fields **only** with
   `update`". That contradicted `approve`, `feedback` and `integrate`, and
   the record model's hand edits.
2. Medium: staged retrieval had lost its second stage, and the statement
   that `context` writes nothing.
3. Medium: several rules named no owning record, although revised
   acceptance 3 asks for one on each line.
4. Medium: `docs/record-model.md` still said the README describes the
   commands, which sent a reader in a circle.
5. Low: three rows of the board's key table overstated where a key works.
6. Low: the development tooling was only in the README, where AGENTS.md no
   longer sends anyone.
7. Low: the plan promised links to the adapters, and the README did not
   give them.
8. Low: small facts had lost their owner: that there is no pre-push hook,
   that rebuilding is manual, and that Grove starts an agent only through
   `run` or `R`.
9. Low: the `go install` caveat sat apart from the command it limits.
10. Low: one link sent four different subjects to a single record-model
    anchor.
11. Low: an ambiguous "it" in AGENTS.md.
12. Low: `run`'s refusals duplicated `grove --help`.

Round 2 confirmed all twelve. It raised new low points:

- the key rows were still slightly wrong;
- the brief repeated the installed-binary policy;
- two lines in `docs/commands.md` were unwrapped;
- a README comment was awkward;
- G-041 was a weak owner for the sibling write-scope rule;
- the race rule could name G-081;
- the record model still summarised `versions` and `workspace` beside
  `docs/commands.md`.

Round 3 confirmed every round-2 fix and found no new issue.

## Disposition

- Findings 1 to 12: fixed in `e07f98e`, and confirmed in round 2.
  - For finding 3, AGENTS.md now names the owning record wherever a record
    holds the reasons: G-065, G-052, G-017, G-041, G-071, G-081 and the
    brief.
  - The Constraints header says that a rule naming no record is AGENTS.md's
    own policy. That covers the commit, verification and sibling
    write-scope rules, which no record owns.
  - For finding 8, the `--commit` message example was not restored.
    `--help` and the record model already say "generated message".
- The round-2 points: fixed in `833ac1b`, and confirmed in round 3.
- Open findings: none. The owner has not judged the result.
