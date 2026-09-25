---
id: "G-176"
type: review
title: "Review of G-170 release identity"
status: current
created: "2026-09-25T21:37:12Z"
updated: "2026-09-25T21:37:26Z"
work: ["G-170"]
examined: "217c703"
---

## Examined

Three rounds by independent `grove-reviewer` agents, each fresh, dispatched
by the headless `/grove-work G-170` session on `worktree-G-170` from main
`47852e3`, against [G-170](G-170-release-identity.md)'s outcome, constraints
and acceptance and plan [G-174](G-174-plan-for-g-170-shared-build-iden.md).
Round 1 examined `1961905`, round 2 `8f19762`, round 3 `217c703`. They read
the diff and the terms it touches, ran `go vet`, `gofmt -l`, `grove check`,
short and uncached package tests, and built stamped and unstamped binaries
from the linked worktree, a clean clone and a `git archive` tree in
temporary directories outside the checkout; none wrote in the checkout.

## Findings

Round 1 (`1961905`):

1. **Consequential, fixed in `8f19762`.** The stamp, `Build.Revision` and
   `revision unknown` used *revision* for a Git commit, which settled term
   [G-062](G-062-revision.md) reserves for a file's sha256, and G-110 would
   bake the linker name into release tooling. Renamed to
   `github.com/mascah/grove.commit`, `Build.Commit`, `commit unknown`.
2. **Minor, fixed in `8f19762`.** The docs said `(devel)` for a build from
   source; since Go 1.24 a checkout build takes a pseudo-version or tag
   (`+dirty`). Docs and `TestIdentify` rows now use the real shapes.
3. **Minor, fixed in `8f19762`.** The docs said `content` covers everything
   init writes; it leaves out `grove.yaml` and the placeholder brief, now
   said.
4. **Minor, fixed in `8f19762`.** Term [G-152](G-152-shipped-document.md)
   said the guides digest names the shipped copy; it now names both digests.
5. **Minor, fixed in `8f19762` and the handoff.** No artifact-agreement
   rule and no pointer for G-110: the docs state that every artifact of one
   release prints the same line without `vcs` or `modified`, and G-170's
   Next points G-110 at them.
6. **Minor, fixed in `8f19762`.** A personal skill can take precedence over
   the worktree's; the fact says "worktree skill" and the docs say another
   skill is not seen. `TestStampedArchiveBuild` builds with
   `-buildvcs=false` so an enclosing repository cannot lend a commit.

Round 2 (`8f19762`):

1. **Consequential, fixed in `217c703`.** The documented release command
   `-o grove` writes into this repository's `grove/` record directory, so
   `./grove version` failed and a second artifact read `modified`. It now
   builds `bin/grove` (ignored); the author and round 3 ran the block
   verbatim in fresh clones with a clean tree after.
2. **Minor, fixed in `217c703`.** Plan G-174 still named the old stamp;
   updated with the reason.
3. **Minor, fixed in `217c703`.** A misplaced parenthetical in the Attempts
   paragraph; the Go version of an observation.

Round 3 (`217c703`): every round-2 disposition verified, no remaining
findings; two cosmetic notes (line lengths, the plan omitting `go run` from
`(devel)`) left as they are, since `docs/commands.md` owns that account.

Not run by any reviewer: Linux or docker builds, byte-for-byte reproducible
builds (G-110's), a real `grove run` attempt against a live provider.

## Disposition

Every consequential finding is fixed and re-reviewed; round 3 found none
remaining. This review is evidence for the owner's verdict on candidate
`217c703`, not approval.
