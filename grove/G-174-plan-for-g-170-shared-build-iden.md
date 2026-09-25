---
id: "G-174"
type: plan
title: "Plan for G-170: shared build identity and release stamping"
status: current
created: "2026-09-25T21:19:48Z"
updated: "2026-09-25T21:20:14Z"
work: ["G-170"]
---

Plan for [G-170](G-170-release-identity.md), prepared headless in
`worktree-G-170` from main `47852e3`. Single implementer, sequential; one
independent review on the final revision.

## Design

- **One owner.** The root package `grove`, which already embeds the guides
  and the reviewer, owns build identity in `identity.go`: `Identity()`
  returns a `Build` whose `String()` is the `version` line. `internal/cli`
  (`version`) and `internal/attempt` (`attempt.json`'s `grove_version`) both
  call it, and `attempt`'s own `groveVersion` is deleted. It must sit below
  `attempt`, which `cli` and the board both call.
- **Release inputs.** Two unexported variables set by the linker:
  `-X github.com/mascah/grove.version=vX.Y.Z` and
  `-X github.com/mascah/grove.commit=<40-hex commit>`. A stamped value wins
  over Go's build information; without one, the module version
  (a checkout's pseudo-version or tag, `(devel)` without `.git`) and `vcs.revision` /
  `vcs.modified` remain the fallback. A missing commit prints
  `commit unknown` rather than nothing, so an unstamped archive build never
  reads as a known release. `vcs.modified` still marks a stamped build from a
  dirty tree. No timestamp or path enters the line; release builds add
  `-trimpath`.
- **Content identity.** The existing `guides sha256:` digest keeps its
  meaning (work, shape and model, concatenated). A new `content sha256:`
  digest covers everything the binary ships into sessions and projects: those
  three guides, the reviewer definition and every `init` entrypoint template,
  each hashed with its name and length. The entrypoint templates move with
  their marker from `internal/cli/init.go` into the root package
  (`entrypoints.go`, `Entrypoints()`, `ManagedMarker`) so identity and the
  launcher can see them; `init` behaviour and bytes are unchanged.
- **Installed provenance.** The launch keeps the reviewer file's sha256 and
  adds the worktree's `grove-work` skill sha256, plus which of the two differ
  from the launching binary's templates (`differs_from_template`). Facts print
  them on an `Entrypoints:` line, and, for attempts recorded before this, say
  they were not recorded. The `Started:` fact names the launching grove, and
  a line states that the grove the agent runs is not recorded: PATH or the
  project's instructions choose it.
- Adjusted after review of `1961905`: the stamp and the line say *commit*,
  since G-062 settles *revision* as a file's sha256; a stamp that differs
  from Go's `vcs.revision` adds `vcs OTHER`; release builds write `bin/`,
  since `-o grove` lands in this repository's record directory.
- Not in scope: `version --json`, per-file template digests, running PATH's
  grove at launch, release workflows or packaging (G-110), entrypoint
  compatibility revisions (G-169), or a selected version number or release
  policy.

## Steps

1. Move entrypoint templates to `entrypoints.go`; `init` uses them. Existing
   init tests pass unchanged apart from the marker's new name.
2. `identity.go` with a pure `identify(info, ok, version, commit)` and
   table tests: release stamped, stamped without commit, development with
   VCS, development modified, missing VCS, no build information; content
   digest changes when any shipped input changes and guides digest ignores
   reviewer and entrypoints.
3. `cli` `version` and `attempt` launch use `Identity()`; attempt records the
   skill digest and template differences; facts updated; tests assert
   `version` output equals the attempt's `grove_version` (fake provider).
4. A build test (skipped under `-short`: it builds): copy the tracked tree
   without `.git`, build with the release flags and `-trimpath`, and outside
   the checkout check `version`, `guide work` and `init` in a disposable
   repository; also build the same copy unstamped and check
   `commit unknown`.
5. Docs: `docs/commands.md` Version and run sections (build contract and
   exact stamping and verification commands for G-110, what each digest
   covers, attempt fields); README install line if it changes.
6. Verification per AGENTS.md, independent review, handoff.
