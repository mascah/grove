---
id: "G-195"
type: work
title: "Coordination-free record IDs"
status: proposed
created: "2026-09-26T02:56:11Z"
updated: "2026-09-26T02:58:10Z"
relates_to: ["G-194", "G-004", "G-006", "G-064"]
---

## Outcome

`grove new` and `grove convert` issue IDs that need no shared state, so
records created in separate clones of one repository, on different
machines, never collide, while every existing numeric ID keeps working
unchanged. The intent is the owner's, selected in
[G-194](G-194-identify-records-by-creation-dat.md) on 2026-09-25: shape on
a laptop and run attempts on an always-on machine without ID collisions.

## Constraints

In scope:

- The form G-194 selects, `G-YYMMDD-xxxxx`: the UTC date of creation, then
  five lowercase Crockford base32 characters (`0-9a-hjkmnp-tv-z`) from
  `crypto/rand`. Proposed: the validator accepts
  `^G-[0-9]{6}-[0-9a-hjkmnp-tv-z]{5}$` and the legacy `^G-[0-9]{3,}$` with
  canonical padding; only the date form is issued. Five characters keep
  the odds of two clones colliding on one day near 3 in a million; four is
  about 1 in 10,000, which is thin for a heavy user, and six is longer than
  the owner wants.
- Proposed: generation retries, a bounded number of times, a tail that the
  loaded project, any local ref or any worktree already holds, using the
  exact-ID search the floor scan already does with one `git grep`. The write
  lock still serializes publication with `update`.
- Delete the counter file, the allocator lock, the floor scan,
  `create.Allocate`, and the initialization and below-floor notices, with
  their tests. `convert` issues through the same generator.
- Where a list orders by ID (`versions` groups, `deps` rows): legacy IDs
  first in numeric order, then date-form IDs lexically, which is
  chronological. The board keeps its `updated` then `created` order.
- Derived slugs cap at 24 characters instead of 32; an explicit `--slug`
  keeps only its character check.
- The documents that own the contract: the record model's identity,
  placement and allocation sections and its clone paragraph; the commands
  document on `new` and `convert`; the shaping guide's sentence on clone
  counters; and this repository's `CLAUDE.md` contract line on the counter.
  Shipped documents keep their linking rules.

Out of scope:

- Renaming or renumbering any existing record; nullsec's records;
  `formerly` and [G-069](G-069-migration-map.md).
- Any fetch, push or remote coordination in Grove.
- Changes to the board's layout beyond what the longer ID forces in columns.
- An entrypoint revision: the entrypoints need nothing new of the binary.

Observed evidence, 2026-09-25 at main `e812672`: the numeric shape lives in
`IDPattern` and `validID` in `internal/project/metadata.go`, the allocator
and `idLine` in `internal/create/create.go`, the ID format in
`internal/update/convert.go`, and `compareIDs` in
`internal/versions/versions.go`. The board sorts by `updated` then
`created`. An attempt ID joins the work ID and a timestamp with `.`, and a
branch joins IDs with `-`; neither parses IDs back, so both accept the new
form. The derived slug cap is in `create.Slug`.

## Acceptance

1. In two disposable clones of one repository with no shared state, `new`
   on each issues date-form IDs, and after merging one into the other
   `check` passes. A test injects the tail generator to show that a tail
   already present in a local ref or worktree is not reused and that the
   retries are bounded.
2. Every record in this repository validates unchanged, `check` passes
   with no rename, and a hand-authored `G-1234` or `G-01` still fails;
   `new` and `convert` never issue a numeric ID.
3. `grove/neutral-ids` and `grove/lock` under the common directory are
   neither read nor written; `new` needs no network and creates only the
   write lock; the read-only commands still create nothing.
4. `versions` and `deps` list legacy IDs before date-form IDs, and the
   latter in date order; the board's order is unchanged.
5. A derived slug is at most 24 characters and the longest generated
   filename is 42 characters; an explicit `--slug` behaves as today.
6. The record model, the commands document, the shaping guide and
   `CLAUDE.md` describe the new form and no longer describe the counter,
   the floor or clone collision checks; `grove guide model` prints the
   updated model; every link resolves.
7. The owner opens the board and a directory listing with a few new
   records and judges the names readable in an actual terminal.

## Next

Assign. The attempt's plan should settle the ID comparator, the retry bound
and how tests inject the tail generator; nothing here waits on another
record. The brief's foundation sentence on sequential IDs changes at the
owner's hand, with this work or before it.
