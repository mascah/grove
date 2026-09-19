---
id: "D-003"
type: decision
title: Allocate IDs with flock, a counter file, and a ref scan floor
status: accepted
relates_to: ["D-002", "W-002"]
created: "2026-09-19T15:25:51Z"
updated: "2026-09-19T15:52:00Z"
---

## Acceptance

On 2026-09-19 the owner accepted this mechanism with "I accept D-003".
[W-002](../work/W-002-create-records.md) implements and tests it.

## Mechanism

Implement the [allocation requirements](../../docs/record-model.md#identity-and-dates)
for [W-002](../work/W-002-create-records.md) with three standard-library pieces
under `<git-common-dir>/grove/`:

1. **Lock:** `flock(LOCK_EX)` on `lock`. The kernel releases it when the
   holding process exits, so a crash leaves no stale lock to detect or expire.
2. **Reservations:** `next-ids`, three lines such as `W 3`, written to a
   temp file, fsynced, and renamed while the lock is held. Rename is atomic;
   a torn write cannot produce a half-updated counter.
3. **Floor:** on every allocation, scan the configured record path in every
   local ref (branches, remote-tracking refs, and tags) with one `git grep` for
   `^id: "<prefix>-<digits>"$`, plus the live record trees of every worktree
   from `git worktree list --porcelain`. The issued number is
   `max(counter, highest observed + 1)`.

Missing counter state initializes from the floor and reports it. A counter
below the floor, for example restored from a backup, is corrected upward and
reported. The floor cannot see records that were reserved but never written,
or written then deleted; only the counter covers those, and if both are lost a
number can be reused. The command states this limit in its initialization
message. Ordinary output from the read-only commands never creates this state.

## Alternatives

- **Lock directory via `mkdir` or `O_EXCL` file.** Portable, but a crash
  leaves a stale lock that needs a timeout or PID check, which is guesswork.
  Rejected while all target machines are macOS or Linux.
- **Scan only when the counter is missing.** Cheaper, but a restored old
  counter would silently reissue numbers. Rejected; the scan is one Git
  command on a local repository.
- **Counter only, no scan.** Simplest, but violates the requirement that lost
  state must not restart at 1. Rejected.
- **Scan only, no counter.** Fails the reservation requirement: a command that
  allocated but has not yet written its file is invisible to the next scan.
  Rejected.
- **Filename-based scan.** Files may be renamed; frontmatter is the identity.
  Rejected.

## Evidence

Probed 2026-09-19 in this repository: a process that took the lock and exited
without unlocking released it to a second process immediately, and
`git grep` across all local refs returned exactly the four record IDs.
W-002 implemented the mechanism in `internal/create` the same day; its tests
cover the floor from refs and live worktree files, counter initialization and
below-floor correction, twenty concurrent allocations across two worktrees, a
consumed reservation after a failed creation, and lock release when the holder
is killed. The first real allocation issued W-003 in this repository.

## Reconsideration

Revisit if Windows support is required (no `flock`; use `LockFileEx`), if the
ref scan becomes slow on large repositories (scan only on initialization or
mismatch, keeping the counter authoritative), or if separate clones need a
shared allocation authority.
