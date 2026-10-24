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

Status: done

`internal/instruct/requires/implement-unattended.md` lets the report name
one line per change to add to `context` or under a key of `rules`, with
what was found, and says not to add it to `grove.yaml`.
`docs/getting-started.md` step 5 says `--unattended` never edits
`grove.yaml` but may name such a line in its report; `dev/walkthrough.md`
"Implement a cycle unattended, headless" says `grove.yaml` is unchanged
and a suggested line is in the last message. TestImplementInstructions
checks the unattended instructions print the rule and the owner ones do
not.

Trial, 2026-10-24, on the ledger fixture with `context: ""` and cycle
261024-k0000 holding "Add ledger delete ID":
`claude -p "/grove-implement 261024-k0000 --unattended"` built and
archived the change on a new branch, asked nothing, left `grove.yaml` as
main has it, and its last message ended with a "Possible `grove.yaml`
addition (not applied)", a `context` line naming
`python3 -m unittest`.

## Scenarios

- implement "The owner agrees": T-01 trial, Claude Code on the ledger
  fixture, passed (the line is in `context`, not under `rules: skills:
  implement:`, and quotes the `-v` correction rather than the test-naming
  one). TestImplementInstructions shows the rule is in the owner's
  instructions.
- implement "The owner declines": T-01 trial, forked session, passed.
- implement "Nothing for next time": T-01 trial with no correction,
  passed; and the walkthrough's main line, by hand, not run here.
- implement "Two corrections": T-01 rerun, one line suggested, passed.
- implement "A test command nobody wrote down": T-02 trial, passed.
  TestImplementInstructions shows the rule is only in the unattended
  instructions.

## Reconciliation

- implement: applied to specs/implement/spec.md
