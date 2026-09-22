---
id: "G-077"
type: review
title: "G-076 independent review of the status filter"
status: current
created: "2026-09-22T15:19:38Z"
updated: "2026-09-22T15:20:03Z"
work: ["G-076"]
examined: "32b7b40"
---

## Examined

## Findings

## Disposition

## Review

Independent reviewer: a fresh Claude Code reviewer agent (Fable 5.1), read-only,
on 2026-09-22, examining branch `worktree-G-076` at candidate `32b7b40`
(base main `ff0f3e9`) against
[G-076](G-076-filter-grove-list-by-status.md)'s constraints, proposed design
and acceptance 1-5. Its evidence: `go test -short -count=1 ./internal/cli`
green; ff0f3e9 and 32b7b40 built as separate binaries print byte-identical
`list` output on this checkout (acceptance 2); `--status bogus`, `--status=`,
a bare `--status`, a blank value, `show G-076 --status active` and
`check --status=open` all exit 2 naming the option with empty stdout and no
`Project:` line (acceptance 3); `--status=review` prints the header alone
(acceptance 4); three mutations of the implementation (filter, list-only
guard, vocabulary check removed) each fail `TestListFiltersByStatus`;
`go vet`, `gofmt -l`, `grove check` (75 records) and the full suite pass.

Verdict: nothing consequential; ship after the one-word documentation fix.

## Findings and dispositions

1. Minor, `docs/record-model.md`: the `list` bullet said a value "outside the
   union of the status vocabularies above" is refused, but the work, question
   and decision vocabularies sit below it under Lifecycle and validation
   boundary. Fixed in `cd1067d`: it now names the type table.
2. Minor, process: at `32b7b40` the record was still `active` without a
   candidate or handoff. That is the handoff step, written after this review.
3. Nit, `internal/cli/cli_test.go`: the row comparison split titles on every
   space, so a title leaking into another column would pass. Fixed in `cd1067d`:
   rows are split on runs of two or more spaces, the tabwriter's column gap.

The tabwriter decision (a filtered table's columns fit its own rows, so it is
narrower than the unfiltered table) was accepted, not a finding: the
unfiltered widths already float with the records present, the record's design
says nothing about padding, and no consumer parses `list` by column offset.

The fixes in `cd1067d` touch one documentation sentence and one test helper and
were self-checked, not re-reviewed: the test still passes and still fails
each of the three mutations above by construction of its comparison.
