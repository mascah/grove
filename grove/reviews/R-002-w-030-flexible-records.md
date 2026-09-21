---
id: "R-002"
type: review
title: "Independent review of W-030 flexible records"
status: current
created: "2026-09-21T16:32:01Z"
updated: "2026-09-21T16:32:02Z"
work: ["W-030"]
examined: "94515975deb6b6e7a4fad8a0b2b04c2de9ab4d90"
---

## Examined

Round 1: commit `9451597` on `worktree-W-030` against base `bd6debe`, by an
independent reviewer agent that wrote none of the code and edited nothing. It
read [W-030](../work/W-030-flexible-records.md),
[D-006](../decisions/D-006-stable-knowledge.md),
[P-001](../plans/P-001-w-030-flexible-records.md) and the record model, and
probed throwaway repositories with a binary from each of `bd6debe` and
`9451597`. Its own runs: `go vet`, `gofmt -l` and `grove check` clean;
`go test -count=1 ./...` and the race run each failed once and passed on
retry, on `TestContextLeavesEverythingUnchanged` (Git's background maintenance
lock racing the fixture snapshot; the test predates W-030) and the connected
terminal check under `-race` (timing; also seen at the base per R-001).

## Findings

Consequential:

1. `internal/tui/model.go` `isWork` took the type of the first version with a
   record, and committed sources sort first. With schema 3 letting sources
   disagree, work reclassified from a page in the live checkout vanished from
   the board and a page reclassified from work was shelved. The added board
   test gave every version one type, so it could not fail.
2. `convert` compared `formerly` byte-exactly, so on a case-insensitive
   filesystem `docs/plans-old.md` and `docs/Plans-old.md` became two records of
   one file with `check` passing: the rerun duplicate acceptance 7 forbids.

Minor: (3) `update --set formerly=` on schemas 1/2 changed wording from the
unknown-field message, and `new bogus` offered `page` to a schema-1 project;
(4) P-001 described a write order the code inverts; (5) `convert` refuses the
brief and the docs did not say so; (6) `--slug` had different rules on `new`
and `convert`; (7) a partial conversion printed no mapping line.

Notes: a document's BOM was embedded mid-file; `.MD` and dot folders are
silent rules; a refusal after reservation consumes an ID; `history.go` parses
with schema 3 unconditionally (harmless today). Confirmed sound: the allocator
split under a 12-way mixed-binary concurrency probe, nested and ref recovery,
the page boundary, reclassification guards, path-escape and symlink refusals,
CRLF/BOM and block-list handling in record conversion.

Acceptance as judged at `9451597`: 1-5 met, 6 not met (finding 1), 7 partial
(finding 2).

## Disposition

All seven findings fixed with regressions in the commit after `9451597`:
`isWork` follows the board source's own record and the board test now has
sources that disagree; `formerly` is compared without case in `convert` and in
validation; schema-1/2 wording restored for `formerly` and `new`; P-001 and
the record model corrected (write order, brief, case rule, gaps, dot folders,
partial mapping, the deleted-neutral-work limit); one slug rule through
`create.ValidSlug`; a partial conversion returns and prints its mapping with
exit 1; a document's leading BOM is dropped. Left as is: `.MD` is not a record
extension and `history.go`'s permissive parse, both documented or harmless.
The pre-existing maintenance-lock flake is outside W-030 and is reported, not
fixed. Round 2 is recorded below when it returns.
