# Implementation: Implement suggests a rules line

## Tasks

### T-01 — Implement suggests a rules line

Status: done

`internal/instruct/requires/implement-owner.md` has the agent, when it
finishes a change, suggest one line to add to `context` or under a key of
`rules` in `grove.yaml`, from one correction, quoting it, never several
merged, and add it on the branch only if the owner agrees, with no
`Files:` line. `docs/reference.md` ("grove.yaml") and
`docs/getting-started.md` step 5 say so in plain words;
`dev/walkthrough.md` step 5 suggests no line, and "A line for next time"
is the trial below. `walkthrough_test.go` TestImplementInstructions checks
the owner instructions print the rule and the unattended ones do not.
`just check` passed.

Trial, 2026-10-24, Claude Code headless on the ledger fixture with change
"Add ledger delete ID" (as `dev/trials/e.sh` writes it) approved on main:

- First run, with the rule as first written ("one line ..., quoting the
  correction"): after two corrections (name each test after the
  requirement it checks; run `python3 -m unittest -v`) it suggested one
  `context` value merging both, paraphrased. The rule now says "from one
  correction, quoting it ... never several merged".
- Rerun on a fresh copy, same two corrections: it suggested one line,
  quoting the `-v` correction. Declined in a fork of the session:
  `grove.yaml` unchanged. Agreed in the original session: `grove.yaml` on
  the branch holds the line in `context`, no `Files:` line names
  `grove.yaml`, `grove check`: "No problems in grove/."
- A fresh copy with no correction: "You didn't correct me during this
  change, so I have no `grove.yaml` line to suggest."; `grove.yaml`
  unchanged.

### T-02 — Unattended names a line, never writes it

Status: not started

## Scenarios

## Reconciliation
