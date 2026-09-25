---
id: "G-170"
type: work
title: "Give every distributed build and attempt attributable release identity"
status: active
created: "2026-09-25T21:04:55Z"
updated: "2026-09-25T21:37:56Z"
relates_to: ["G-062", "G-152", "G-169"]
---

## Outcome

Every distributed Grove binary and recorded attempt can be attributed to its
release and shipped instruction content, including a build from a source archive
without Git metadata. One implementation owns build identity.

Owner intent, review and shaping conversation 2026-09-25: prepare a dependable
versioning and distribution process for installed Grove. The following design
is proposed, not a selected version number or release policy.

## Observed evidence

At main `05892a2`, [versionLine](../internal/cli/init.go) uses Go build metadata
and a digest over work, shape and model. It excludes generated entrypoints and
the embedded reviewer. [Attempts](../internal/attempt/attempt.go) independently
format module version and commit, without that guide digest.

The review extracted tracked HEAD with `git archive`, built successfully without
`.git`, and observed `grove (devel) guides sha256:67dab310de86`. Its guides and
`init` worked in a disposable project. This is a source-archive smoke test, not
Homebrew installation, cross-platform or reproducible-build evidence.

## Proposed design and constraints

- Provide explicit release version and source revision metadata for release
  builds, with a shared representation used by `version` and attempt records.
  Retain useful development-build fallbacks, and never present absent metadata
  as a known release. Keep ordinary checkout builds and `go install` usable.
- Identify all binary-owned workflow content, including reviewer instructions
  and generated entrypoints. Existing guide digests may remain separately
  useful; document exactly what each identity covers. Metadata must not imply
  that a project's custom or older installed files equal the binary templates.
- Preserve the actual reviewer-file identity already captured at launch and
  capture entrypoint provenance as needed to explain mixed installations.
  Distinguish the launching binary from the executable an agent resolves through
  PATH or project policy; a launcher stamp alone does not prove that equality.
- Keep release version, record schema and entrypoint compatibility distinct.
  Do not add independent release streams for guides or per-release churn in
  otherwise unchanged project adapters. Coordinate with
  [G-169](G-169-harness-upgrade-compatibility.md) without requiring it to finish
  before preparing this independent work.
- Specify release build inputs so G-110 can stamp both source-based and packaged
  builds. No timestamp or machine path should gratuitously prevent reproducible
  artifacts. Packaging and byte-for-byte reproducibility claims belong to G-110.

## Acceptance

1. An explicitly stamped source-archive build reports the requested release and
   source revision without `.git`; ordinary development builds remain usable and
   visibly distinguish known metadata from missing metadata.
2. `version` and a fake-provider attempt agree on the shared executable identity.
   Shipped-content changes affect the relevant digest; installed custom/older
   files retain their own provenance rather than borrowing the template identity.
3. Tests cover release, development and missing-metadata cases and verify the
   identity is usable outside Grove's development checkout. Any claim about the
   agent's resolved Grove executable is checked or explicitly marked unknown.
4. The build contract, command documentation and attempt representation are
   reconciled; applicable repository checks pass. G-110 receives exact stamping
   and verification instructions, with no release published by this work.

## Evidence

Implemented headless on `worktree-G-170` from main `47852e3`, starting from
this record at `sha256:1f4e92a9…` and plan
[G-174](G-174-plan-for-g-170-shared-build-iden.md) at `sha256:af39d944…`
(committed `bef952c`). Code and docs: `1961905`, `8f19762`, `217c703`; the
candidate commit adds only this record's evidence and review
[G-176](G-176-review-of-g-170-release-identity.md).

Decisions, routine within the outcome:

- The root package `grove` owns identity (`identity.go`), since it embeds
  the shipped content and sits below both `internal/cli` and
  `internal/attempt`, which the board also calls. `init`'s entrypoint
  templates moved there verbatim (`entrypoints.go`); `init`'s bytes are
  unchanged.
- The stamp is `-X github.com/mascah/grove.version=…` and
  `-X github.com/mascah/grove.commit=…`: *commit*, not *revision*, which
  settled term [G-062](G-062-revision.md) reserves for a file's sha256.
- The `guides` digest keeps its meaning (`sha256:67dab310de86` at base and
  candidate); a new `content` digest covers the three guides, the reviewer
  and every other `init` entrypoint template, each hashed with name and
  length. `grove.yaml` and the placeholder brief are project-owned and out.
- The agent's own `grove` is marked not recorded rather than probed: `PATH`
  or project policy chooses it and a launch-time probe would prove neither.

Against acceptance:

1. `TestStampedArchiveBuild` (skipped under `-short`) copies the tracked
   tree without `.git`, builds it stamped and unstamped, and outside any
   checkout gets `grove v9.9.9 (<commit>) guides … content …` and
   `grove (devel) (commit unknown) …`; unknown version or commit is always
   said. `TestIdentify` covers release, release without commit, stamp that
   differs from Go's commit (`vcs OTHER`), dirty tree, development
   pseudo-version, `go run`/archive, `go install`, no build information.
2. `version` and `attempt.json`'s `grove_version` both print
   `grove.Identity().String()` (asserted in `TestGuideAndVersionNeedNoProject`
   and `TestRunToResult` with the fake provider).
   `TestIdentityDigestsCoverTheirContent`: a reviewer or skill template
   change moves `content`, not `guides`. The launch records the worktree
   skill's sha256 (`skill`) beside the reviewer's and lists in
   `differs_from_template` which are not the launching templates; the
   `Entrypoints:` fact says so, and "not recorded at launch" for older
   attempts (`TestInputsChanged`, `TestFactsOfAnEarlierLaunch`).
3. The stamped archive build runs `version`, `init`, `check` and
   `guide work` in a temporary Git repository with no Grove source. The
   `Agent's grove:` fact states it is not recorded.
4. [Version and guide](../docs/commands.md#version-and-guide) is the build
   contract: line format, fallbacks, what each digest covers, the release
   command into `bin/`, `go version -m` (which omits `-ldflags` under
   `-trimpath`), and the rule that every artifact of one release prints
   the same line without `vcs` or `modified`.
   [Attempts](../docs/commands.md#attempts) covers the new fields; README's
   install check and term [G-152](G-152-shipped-document.md) are
   reconciled. No release was published.

Verification at `217c703` (Go 1.26.8, darwin/arm64): `go vet ./...` clean,
`gofmt -l .` empty, `go run ./cmd/grove check` OK 160 records,
`go test -count=1 -timeout 120s ./...` all ok. `go test -short` per touched
package under five seconds. The documented release block, run verbatim in
a fresh clone, printed `grove v0.1.0 (<HEAD>) guides sha256:67dab310de86
content sha256:9fb4892e0afb`, `-trimpath=true`, and left `git status` clean
(author at `8f19762`, round 3 at `217c703`).

Review: [G-176](G-176-review-of-g-170-release-identity.md), three rounds,
two consequential findings (the *revision* term; `-o grove` landing in the
record directory), both fixed; round 3 found none remaining.

Limits: no Linux or docker build, no byte-for-byte reproducibility check
(G-110's), no live `grove run`. A linked worktree inside its repository
builds with the enclosing checkout's pseudo-version and commit, shown as
`vcs OTHER` beside a stamp; build releases from a clone. Entrypoint
compatibility revisions remain [G-169](G-169-harness-upgrade-compatibility.md)'s.

## Next

In review on `worktree-G-170` (base `47852e3`); the owner judges the
candidate this status change names. G-110 stamps release artifacts exactly
as [Version and guide](../docs/commands.md#version-and-guide) says, choosing
the version itself.

```sh
grove approve G-170 "VERDICT"   # in this worktree, clean
grove integrate G-170           # in the main checkout, clean
```

Or `grove feedback G-170 "TEXT"` here to return it to active.
