# Implementation: Check limits the size of rules and specs

<!-- Progress on change 261127-hvq1r, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — `grove check` reports a rule or a spec a change grows past its limit

Status: done

`doc.Requirement.RuleWords` counts the words of a rule, its text before its
first `#### Scenario:` heading outside fenced code
(`internal/doc/requirements.go`). `sizes` in `internal/cli/check.go` adds
the three reports to each open change, with the limits as the constants
`ruleLimit` and `specLimit`; it reads only what `atBases` already read, so
`grove check` starts no more Git processes. `requirement` now returns the
whole `doc.Requirement` so `sizes` can count a rule's words. The
requirements spec's Purpose starts with the proposal's sentence.

Checks: `internal/doc` `TestRuleWords` (a rule of exactly 60 words, a rule
with a fenced code block holding a scenario heading, a requirement with no
scenario); `internal/cli` `TestSizes` (every scenario of the three
requirements); `TestGitProcesses` unchanged; `bin/grove check` in this
repository with this change open prints `No problems in grove/.`;
`just check` passes.

Check: a MODIFIED requirement's earlier text is the last text a change it depends on gives it, else what `base` holds, as the proposal recommends; the spec's "before" count is the spec at `base`, so a dependency's other additions to the spec are not counted in either total. No test covers a dependency.
Found: `TestServe` in `internal/cli/walkthrough_test.go` fails intermittently at its first `serving` call ("grove serve printed \"\"", expecting `http://127.0.0.1:7070/`), on main at 9d55895 as on this branch: 4 of 5 runs of `go test ./internal/cli -run TestServe -count=5` failed on both, with nothing listening on port 7070; seen when `just check` failed once on the untouched baseline.
Check: the messages print plain numbers (`3000`, not `3,000`), and a spec report names the change spec's path with no line.

### T-02 — The change-spec guide, the reference and the walkthrough

Status: done

The change-spec guide gains a bullet giving both limits and saying
`grove check` reports a change that exceeds either (instruction text 8,492
words of 9,000). `docs/reference.md`'s `grove check` row names the reports
and the limits, linking the requirements spec and the guides to splitting
and reconciling a maintained spec. `dev/walkthrough.md` gains "A rule over
its limit" under Other paths, which `TestSizes` follows: shortening the
rule clears its report, and `grove instructions propose` prints the
guide's line.

Checks: `TestSizes`, passed; the step run by hand on a fresh fixture with
`bin/grove` (a rule of 61 words reported at line 5, exit 1; one word cut,
`No problems in grove/.`); `just check` passes.

## Scenarios

- requirements "A rule of 61 words": internal/cli/walkthrough_test.go TestSizes ("Export"), passed
- requirements "A rule of 60 words with long scenarios": TestSizes ("Import"), passed
- requirements "A long rule no open change names": TestSizes ("Untouched", 90 words), passed
- requirements "A long rule made longer": TestSizes ("Long", 80 to 95), passed
- requirements "A long rule made shorter": TestSizes ("Cut", 80 to 70), passed
- requirements "A short rule taken over the limit": TestSizes ("Short", 50 to 65), passed
- requirements "Only the scenarios change": TestSizes ("Same"), passed
- requirements "Adding to a spec already over": TestSizes (`over`, 3000 to 3200), passed
- requirements "Adding and removing": TestSizes (`mixed`), passed
- requirements "Taking a spec over": TestSizes (`near`, 2400 to 2600), passed
- requirements "A spec that stays under": TestSizes (`under`), passed
- requirements "A new capability": TestSizes (`exports`, 0 to 2600), passed
- requirements "The propose instructions": TestSizes (the guide's line in `grove instructions propose`), passed

## Reconciliation

- requirements: applied to specs/requirements/spec.md
