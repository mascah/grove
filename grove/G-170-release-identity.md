---
id: "G-170"
type: work
title: "Give every distributed build and attempt attributable release identity"
status: active
created: "2026-09-25T21:04:55Z"
updated: "2026-09-25T21:20:24Z"
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

## Next

Proposed, unassigned. The owner can assign `$grove-work G-170` independently of
G-169. Prepare the shared identity and build-input contract, retaining honest
fallbacks for unstamped builds and historical attempts. G-110 consumes it when
assembling release artifacts.
