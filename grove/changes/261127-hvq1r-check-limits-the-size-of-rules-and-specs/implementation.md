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
Check: the messages print plain numbers (`3000`, not `3,000`), and a spec report names the change spec's path with no line.

## Scenarios

## Reconciliation
