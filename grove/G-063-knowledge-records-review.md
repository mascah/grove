---
id: "G-063"
type: review
title: "Independent review of G-037 knowledge records"
status: current
created: "2026-09-21T05:59:58Z"
updated: "2026-09-21T05:59:59Z"
work: ["G-037"]
examined: "fc9bef1563a27fb6d9b22d7dafe1cd08f8e06a22"
relates_to: ["G-051"]
formerly: "R-001"
---

## Examined

[G-037](G-037-knowledge-artifacts.md) against its acceptance,
[G-051](G-051-typed-knowledge-records.md) and its
[plan](G-053-knowledge-artifacts-plan.md). An independent reviewer
agent in the implementing session's harness, read-only in the execution
checkout, with its own binary and fixtures in fresh temporary repositories;
it edited nothing under review. Three passes over `d9fc2a5..HEAD`: `e8654fc`,
then fix round 1 at `46ed102`, then fix round 2 and the combined revision at
`fc9bef1`, which `examined` names. This record and G-037's reconciliation were
written after that commit and were not reviewed.

Reviewer's verification at `fc9bef1`, run serially: `gofmt -l .` and
`go vet ./...` clean; `go test -count=1 ./...` and
`go test -race -count=1 -p 1 -timeout 30m ./...` all eight packages ok;
`check` reports 40 records. Its probes compared schema-1 behavior against a
binary built from the base, tried twenty brief paths, and checked each new
rule by mutation: every one fails a test when removed, except the `T-001`
word in the ID hint.

## Findings

Round 1, at `e8654fc`:

1. Consequential: `new term` with an existing title wrote the file and left
   the project unloadable. Fixed: refused before an ID is reserved.
2. The symlinked-parent guard for the brief had no test. Fixed.
3. On a case-insensitive filesystem `brief: grove/Work/x.md` made a record
   serve as the brief. Fixed: type folders are compared without case.
4. The record model implied committed sources check an in-root brief's
   existence. Fixed: they never do.
5. Schema 1's misplaced-record message had changed. Fixed: byte-identical to
   the base again.
6. The ID hint named only `W`, `Q`, `D`. Fixed.
7. No test read a committed schema-2 tree. Fixed, through `versions` and
   `workspace`.
8. The plan said eight terms; nine exist. Candidate was added because review,
   approval and integration are defined by it; the plan now says so.

Round 2, at `46ed102`:

- A. A term defined by another process between `new`'s load and its turn at
  the write lock could still be written. Fixed: rechecked under the lock,
  reported as "reserved but not created".
- B. A brief under the record root spelled with another case than the disk
  was reported as misplaced. Fixed: the exemption ignores case.

Final, at `fc9bef1`: no consequential finding open. Two accepted limits: a
writer that is not Grove can still add a colliding term inside the lock
window, as it always could add a duplicate ID; and on a case-sensitive
filesystem a second file differing from the brief only by case loses its
misplaced-file warning. No record can be hidden by either. One performance
note, not acted on: `compareIDs` rebuilds the prefix string per comparison.

Not checked by the reviewer: the board in a real terminal (no TUI file
changed), and a true two-process race on `new term` (exercised through
`create.New` with a stale project instead).

## Disposition

All findings fixed in `46ed102` and `fc9bef1`, each with a regression test
except finding 6. This is review evidence, not approval: the owner's judgment
of G-037, including the nine terms, is separate and outstanding.

Observed during verification and not caused by this work (`internal/repo` has
no change on the branch): under `-race`, or a fully parallel `go test ./...`,
tests in `internal/versions`, `internal/cli` and `internal/tui` intermittently
fail "files in the repository changed" on `.git/objects/maintenance.lock`, and
`internal/versions` hung twice in roughly thirty runs, until the test
timeout. The captured dump shows `repo.GitContext` blocked in
`exec.Cmd.Wait` (`internal/repo/repo.go:124`) after `git` exited, its output
pipe still held open, consistent with Git's detached auto-maintenance in the
fixtures. The base `d9fc2a5` shows the same family of failure
(`TestWorkspaceStale`, `maintenance.lock`). It deserves its own work item.
