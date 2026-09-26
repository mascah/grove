---
id: "G-194"
type: decision
title: "Identify records by creation date and a random tail, without a counter"
status: accepted
created: "2026-09-26T02:56:11Z"
updated: "2026-09-26T02:58:25Z"
relates_to: ["G-004", "G-006", "G-064", "G-001", "G-051", "G-195"]
---

## Decision and authority

On 2026-09-25, in an interactive shaping session, the owner selected
coordination-free record IDs after an exploration of running Grove from two
machines: a laptop for shaping and an always-on Mac mini for attempts.
Presented with the designs under Alternatives, they answered "yeah I agree
with the recommendation B with the date form. Some sort of slug might be
nice but the length would need to be constrained", then chose, with four
options in front of them, to keep the slug in the filename only with its
derived length capped at 24 characters.

Selected:

- An ID Grove issues is `G-`, the UTC date of creation as six digits, a
  hyphen, and a random tail of lowercase Crockford base32 characters,
  generated without any shared state: `G-250925-7k2qm`. The tail's length is
  a routine choice of [G-195](G-195-coordination-free-record-ids.md); five is
  proposed there.
- The ID is the identity and never changes, and it is not derived from the
  title. The slug stays in the filename only, as today, with the derived slug
  capped at 24 characters so the longest generated filename stays within one
  character of today's 41.
- Existing numeric IDs remain valid identities, are never renumbered, and
  are never issued again. No record is renamed: this is not a second
  reconciliation like [G-052](G-052-migrate-knowledge.md).
- Legacy and date-form IDs coexist in the validator; `new` and `convert`
  issue only the date form.

This revises [G-004](G-004-sequential-ids.md), sequential IDs coordinated
through Git's common directory, and retires
[G-006](G-006-allocator-mechanism.md)'s allocator for future creation: the
counter, the allocator lock and the floor scan go with it.
[G-064](G-064-stable-knowledge.md)'s "one neutral sequential ID namespace"
becomes one neutral namespace without a sequence; its stable identity,
stable placement and short-slug rules stand. The brief's foundation
sentence on stable sequential IDs and clone collision checks needs the
matching change; the brief owns direction, so that edit is the owner's.

## Evidence

Observed 2026-09-25 at main `e812672` with a build of this checkout:

- Two local clones each initialized their own counter at 192 and both
  issued G-192 for different records. After one fetched the other, its next
  ID was G-193, because the floor scan reads remote-tracking refs. Git merged
  the two G-192 files without a conflict, since the filenames differ; only
  `grove check` reported the duplicate.
- The counter lives at `grove/neutral-ids` under the Git common directory
  (`internal/create/create.go`), so clones share nothing. Grove runs no fetch
  or push anywhere, `versions` and the board read `refs/heads/` only, and
  attempts live under the launching machine's common directory.
- The two-machine workflow makes the collision the default outcome rather
  than an edge case: an attempt on the mini mints plan and review IDs on an
  unpushed branch while the laptop shapes.
- The record model forbids renumbering, so a collision has no sanctioned
  repair.
- A compare-and-swap counter on a remote ref, pushed with
  `--force-with-lease=refs/grove/ids:<seen>`, was verified against a bare
  repository: a stale claim is rejected, provided each claim's blob is
  unique, because identical content hashes identically and the loser's push
  becomes a silent no-op.
- [G-001](G-001-starter-defaults.md)'s first trial used a 20-character
  random ID with a timestamp-prefixed filename, which the owner found too
  long to read (G-004). The date form is 14 characters with a five-character
  tail.
- Collision odds for a five-character tail: two clones each creating ten
  unseen records on one day collide with probability about 3 in a million
  per day. A four-character tail is about 1 in 10,000 per day; six
  characters about 1 in 100 million.
- mex (`https://github.com/mex-memory/mex`) mints 26-character ULIDs and
  states there is no cross-clone lock, so conflicts surface later in Git.
  Grove takes the coordination-free principle, not the length.

## Alternatives

- **Sequential IDs with a remote arbiter.** A counter ref on a Git remote
  updated by compare-and-swap, or the always-on machine as the allocator
  over ssh. Verified workable. Rejected: sequential numbering is a property
  of a central service; every clone that creates records needs push access
  to the arbiter, which a contributor on a fork lacks; `new` needs the
  network; an offline mode reopens the risk; and a missed coordination has
  no repair.
- **ULIDs as mex issues them.** Coordination-free and time-sortable, at 26
  characters. Rejected for the length G-004 already rejected; the date form
  keeps the ordering property readably.
- **A short random tail without the date.** Shortest, but a directory
  listing loses chronological order and the name says nothing. Rejected.
- **The slug inside the ID.** Most readable, but derived from the title,
  long in references, commits and branch names, and open to same-day
  same-slug collisions. Rejected; G-004 and
  [G-051](G-051-typed-knowledge-records.md) avoided a slug identity for the same
  reasons.
- **Per-clone prefixes or reserved blocks.** Typed and per-clone prefixes
  were removed by G-064; blocks still need an arbiter. Rejected.
- **Renumbering existing records to the new form.** 185 records and about
  4,500 references here, plus nullsec's, and the never-rename rule protects
  identity. Rejected; coexistence costs one alternation in the validator.

## Reconsideration

Reopen if a collision is observed in practice (lengthen the tail), if an
adopter needs sequential numbers enough to run an arbiter, if the first
release wants to fix one ID form and retire the legacy pattern, or if
browsing shows the date-form names awkward in the board or in editor file
pickers.
