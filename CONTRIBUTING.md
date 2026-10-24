# Contributing

What Grove does and does not do is set by [`CONSTITUTION.md`](CONSTITUTION.md).
Read it first.

## Setup

Install Go, at the version `go.mod` names, Git, and
[`just`](https://github.com/casey/just).

## Commands

```sh
just build                 # bin/grove
eval "$(just dev-path)"    # build, and put this checkout's bin/ first on PATH
just check                 # gofmt, vet, tests and the line budget; run before every commit
just fixture ~/tmp/ledger  # a fresh fixture repository to walk through
just install               # this checkout's build as your everyday grove, in ~/.local/bin
```

Guides, templates and schemas are embedded, so `just build` picks up an
edit to them. After changing the skill text or a harness row (both in
`internal/cli/init.go`), run `grove init` again in the fixture.

## Tests

`internal/cli/walkthrough_test.go` follows the walkthrough (rule 17), on
repositories built from the fixture. `internal/tui/tui_test.go` covers the
reader. `just check` runs both.

`dev/walkthrough.md` is run by hand: build with `eval "$(just dev-path)"`,
make a fresh fixture with `just fixture`, and follow the steps a change
touches in Claude Code and in Codex.

When a pull request changes instruction text in `internal/instruct`, run
the trials in `dev/trials.md` by hand, in Claude Code and in Codex, and put
their counts in its description.

## Commits and pull requests

Every commit on `main` is a squash-merged pull request whose title is a
Conventional Commit (`type(scope): subject`, `!` before the colon for a
breaking change) and whose description is the body, with a
`BREAKING CHANGE:` footer saying what a user must do. Fixes found in review go
in the same pull request.

- Types: `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `build`, `ci`,
  `chore`, `revert`.
- Scopes, when there is one: `cli`, `instruct`, `doc`, `tui`, `docs`,
  `constitution`, `ci`.
- The description says what was demonstrated and how: on the fixture,
  wherever there is behavior to show.

## Breaking changes

A change is breaking when a project that worked before needs an edit to keep
working, or its CI fails: a `grove.yaml` or `.grove.yaml` field renamed or
removed, a document convention changed (`T-NN`, `Status:`, folder layout), a
command, flag, skill or skill folder renamed or removed, a JSON key changed,
or a new `check` problem reported on documents that were clean before. Before
1.0, a breaking change bumps the minor version.

## Releases

release-please keeps a release pull request open on `main`. The next version
comes from the commit titles since the last release: `fix` bumps the patch
version, `feat` the minor, and a breaking change as above. Merging the release
pull request tags the release, writes `CHANGELOG.md` and creates the GitHub
Release; GoReleaser, in the same workflow, attaches the binaries to it.

The release pull request is opened with the default token, so CI does not run
on it, and the owner merges it as admin. Nothing is versioned by hand. For
release-please to open that pull request, the repository's Actions settings
must allow GitHub Actions to create pull requests.
