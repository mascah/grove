---
id: "G-076"
type: work
title: "Filter grove list by status"
status: review
created: "2026-09-22T15:01:38Z"
updated: "2026-09-22T15:20:38Z"
kind: feature
size: small
relates_to: ["G-039"]
candidate: "840a78b"
---

## Outcome

`grove list` can show only the records in a chosen status, so a person or an
agent finds open work in one command instead of piping the whole table
through `grep`. The owner asked for this on 2026-09-22, in the `/grove-shape`
session that wrote this record, as the real change for the
[G-039](G-039-interactive-loop.md) interactive-loop trial (its plan G-075,
step 1, exists only on branch `worktree-G-039` until that work is
integrated, so it is not linked here).

## Constraints

Observed at main `5f07c94`, 2026-09-22:

- `grove list --status proposed` exits 2 with `unknown option --status`.
  `list` takes no option of its own; it prints every record of the checkout
  as `ID TYPE STATUS TITLE`, a page's status as `-`
  (`internal/cli/cli.go:156`). Options are one flat table in `parseArgs`
  in the same file, with `--include` already repeatable.
- The record model gives each type a closed status vocabulary
  (`docs/record-model.md`, "Statuses"), and the loader rejects a value
  outside it. This checkout holds 7 distinct values across 74 records; the
  same word can belong to several types (`proposed` is work, question,
  decision and term; `current` is plan and review).
- The shaping and work guides read the unfiltered `list` at the start of a
  session, and the board has its own fixed status columns
  (`internal/tui/model.go:37`); neither depends on this change.

The owner decided on 2026-09-22, in that session, among three offered scopes:
**status only, and `list` without the flag keeps printing everything.** Out
of scope, each a separate proposal if wanted: a `--type` filter, hiding
finished work by default, filtering the board, and `--json` for `list`.

Proposed design, binding nobody: a repeatable `--status VALUE` (also
`--status=VALUE`), accepted only after `list`; a record is shown when its
status equals any given value exactly; a value outside the union of the
model's status vocabulary is a usage error (exit 2) rather than an empty
result, since the vocabulary is closed; matching nothing prints the header
only and exits 0. The usage text, the README's `list` line and the record
model's `list` contract say so.

## Acceptance

1. `grove list --status proposed` prints the same header and rows, in the
   same order and format, as `grove list` minus every row whose status is not
   `proposed`; `--status active --status review` unions the two.
2. `grove list` without the option prints exactly what it printed before.
3. A value outside the model's vocabulary, an empty value, and `--status` on
   any other command are refused as usage errors (exit 2, message naming the
   option) before any project is read.
4. Filtering to a status no record holds prints the header and exits 0.
5. Usage text, README and the record model describe the option and agree
   with the behavior; `go vet`, `gofmt -l`, `grove check` and the suite pass
   within the AGENTS.md budget.
6. The owner runs it on this repository and judges from the output that it
   answers "what is open" without further piping.

## Evidence

Implemented 2026-09-22 on branch `worktree-G-076` (worktree
`.claude/worktrees/G-076`) from main `ff0f3e9`, starting from this record at
revision `sha256:b5c9464e…` with no plan record (none needed: the proposed
design above fixes every choice). Commits: `32b7b40` the feature, `cd1067d`
the review fixes, `a28f3d5` the review record; the candidate is the commit
that adds this section.

Behavior against the acceptance:

1. `list --status proposed` prints the header and the `proposed` rows in the
   unfiltered order; `--status active --status review` unions them
   (`TestListFiltersByStatus`, `internal/cli/cli_test.go`). Decision: the
   tabwriter fits column widths to the rows it prints, so a filtered table
   is narrower than the unfiltered one; "same format" is read as the same
   four columns, since the unfiltered widths already float with the records
   present. The reviewer accepted that reading (G-077).
2. `list` without the option is unchanged: the reviewer built `ff0f3e9` and
   `32b7b40` and compared their `list` output on this checkout byte for byte.
3. A value outside the union of the type table's vocabularies, an empty or
   blank value, a bare `--status`, and `--status` on `show` or `check` exit 2
   with a message naming the option, empty stdout, and no `Project:` line,
   so no project was read (same test, and the reviewer's manual runs).
4. `list --status=review` on this checkout printed the header alone, exit 0.
5. Usage text (`internal/cli/cli.go`), README and `docs/record-model.md`
   describe the option; verification at `a28f3d5`: `go vet ./...` clean,
   `gofmt -l .` empty, `grove check` → `OK: 76 records`,
   `go test -count=1 -timeout 120s ./...` all green. `internal/versions`
   takes 7.3 s, above the five-second budget, before and after this change
   (its Git process count, not these tests; `internal/cli` is 1.8 s in
   `-short`).
6. Owner's judgment: pending, see Next.

Review: [G-077](G-077-g-076-independent-review-of-the.md), independent,
examined `32b7b40`: nothing consequential; two minor findings and one nit,
each fixed in `cd1067d` or being this handoff. Limits: the `cd1067d` fixes
(one documentation sentence, one test helper) were self-checked, not
re-reviewed.

## Next

Candidate awaits the owner's judgment (acceptance 6). Demo on this repository:

```sh
cd .claude/worktrees/G-076
go run ./cmd/grove list --status active --status review
go run ./cmd/grove list --status proposed
go run ./cmd/grove list --status bogus   # exit 2
```

Integrate on approval, from the main checkout:

```sh
git merge --ff-only worktree-G-076
go run ./cmd/grove show G-076 --json    # take the revision
go run ./cmd/grove update G-076 --expect REVISION --set status=done
git commit -am "docs(G-076): mark done on the owner's acceptance and merge"
```

Feedback instead: write it here and set `status=active` on the branch.
