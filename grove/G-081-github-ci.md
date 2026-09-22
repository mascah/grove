---
id: "G-081"
type: work
title: "Run secure CI, Dependabot, and Claude review on GitHub"
status: proposed
created: "2026-09-22T16:10:10Z"
updated: "2026-09-22T16:11:18Z"
kind: tooling
size: medium
relates_to: ["G-040", "G-044"]
---

## Outcome

Every push to `main` and every pull request on the public
`github.com/mascah/grove` repository gets the same checks the owner runs
locally, run by GitHub Actions under a hardened configuration; dependency and
action updates arrive as Dependabot pull requests; and the Claude Code
workflows from [PR #1](https://github.com/mascah/grove/pull/1) run with
bounded triggers and permissions. The owner asked for this on 2026-09-22,
citing the comprehensive CI suites of their professional projects (the sibling
`xfer` repository's `.github` is the example) and the repository being public.

## Constraints

**Selected by the owner, 2026-09-22, in the shaping session:**

- CI is a signal, not a gate. It runs on pushes to `main` and on pull
  requests; no ruleset requires the checks yet. The owner keeps integrating
  locally with fast-forward merges and direct pushes for now, and wants pull
  requests into `main` supported as a second path. Grove itself has no support
  for the pull-request path; that is a gap for [G-044](G-044-review-integration.md)'s
  neighbourhood, not this work.
- The automatic Claude review runs only for pull requests from this
  repository, only when a PR is opened or marked ready for review, never on
  every push and never for fork PRs. The `@claude` mention workflow keeps the
  action's default of write-access actors only.
- A LICENSE file is out of scope. The owner chose to leave the repository
  unlicensed for now; this is an open item, not a question record.

**Observed on 2026-09-22 at `main` ccdc92d:**

- No `.github` directory exists on `main`. PR #1 adds `claude.yml` and
  `claude-code-review.yml` and nothing else; the `CLAUDE_CODE_OAUTH_TOKEN`
  secret is set (`gh secret list`).
- Repository settings (`gh api repos/mascah/grove`, `.../actions/permissions`,
  `.../branches/main/protection`): public, no branch protection or rulesets,
  default workflow token read-only, any action allowed, SHA pinning not
  required, Dependabot alerts and security updates disabled, secret scanning
  and push protection enabled.
- The local evidence contract in `AGENTS.md` and `lefthook.yml` is `gofmt -l .`,
  `go vet ./...`, `go mod tidy -diff`, `go run ./cmd/grove check`, and
  `go test -count=1 -timeout 120s ./...`. That suite passed in about 8 s wall
  on the owner's Mac (`internal/versions` 7.5 s, `internal/tui` 5.4 s); it is
  Git-process-bound, so a Linux runner's cheaper fork should keep it inside the
  120 s timeout. `-race` across the suite hangs on macOS and is not part of the
  contract.
- Windows does not build: `internal/repo/repo.go:168` uses `syscall.Flock`.
  Linux and macOS build. `internal/tui`'s terminal test needs `python3` and a
  pseudo-terminal, which Ubuntu and macOS hosted runners provide.
- `go.mod` declares `go 1.26.0`; `actions/setup-go` reads it with
  `go-version-file`.
- PR #1 as generated: both workflows use floating action tags, the review
  triggers on `synchronize` (every push) and on fork PRs, and neither has
  `timeout-minutes` or `concurrency`. The action's security documentation
  (`anthropics/claude-code-action/docs/security.md`) says it checks write
  access on the triggering actor by default, warns against checking out an
  untrusted ref at the workspace root, and that `pull_request_target` runs
  with the base repository's secrets. PR #1's body claims Claude can create
  branches and commits while the job grants `contents: read`; whether the
  Claude GitHub App's own token bypasses the job permissions is unverified.
- No existing record mentions CI, GitHub Actions, or Dependabot.
  [G-040](G-040-portable-bootstrap.md) owns selecting a distribution and
  upgrade mechanism, so a release workflow (goreleaser or similar) belongs to
  it, not here; a `worktree-G-040` branch exists at b8f7cf5.

**Proposed design (binds nobody; the assignee selects the details):**

- One `ci.yml` on `pull_request` and `push` to `main`, no `paths-ignore`,
  because `grove check` validates the records under `grove/` too. Jobs: the
  five local checks above plus `go build ./...`, on an `ubuntu-latest` and
  `macos-latest` matrix, `-short` off, `GOFLAGS` untouched. Add
  `govulncheck ./...` as a separate job so a new advisory fails visibly
  without blocking the build job's signal. Optionally CodeQL for Go, which
  is free on public repositories.
- Hardening in every workflow: top-level `permissions: contents: read` with
  per-job additions only where a job writes; every `uses:` pinned to a full
  commit SHA with the version in a comment; `timeout-minutes` on every job;
  `concurrency` keyed on the ref with `cancel-in-progress`; no
  `pull_request_target` anywhere.
- Repository settings, changed by the owner or through `gh api`, recorded as
  evidence rather than files: require SHA pinning for actions, narrow allowed
  actions to GitHub-authored and verified-creator actions plus
  `anthropics/*`, enable Dependabot alerts and security updates.
- `dependabot.yml` for `gomod` and `github-actions`, weekly, each grouped into
  one PR like `xfer`'s, with the `dependencies` label. No Docker, npm, or
  Terraform ecosystems exist here.
- PR #1: keep both workflows, but restrict `claude-code-review.yml` to
  `types: [opened, ready_for_review]` with
  `if: github.event.pull_request.head.repo.full_name == github.repository`,
  pin the action and checkout to SHAs, add timeouts and concurrency, and
  merge it once hardened. The `xfer` review workflow's prompt shape
  (`gh pr diff`, one `gh pr comment`) is a reference, not a requirement:
  PR #1 uses the `code-review` plugin with inline comments instead.
- Not carried over from `xfer`, and why: labeler, stale, PR-title, and
  CODEOWNERS serve a team; release-please and any release build belong to
  G-040; the scheduled security audits that open issues are covered for Go by
  `govulncheck` plus Dependabot alerts; the scheduled codebase reviews are a
  later candidate once the loop has evidence they pay for themselves.

## Acceptance

1. A pull request from this repository and a push to `main` each run the
   CI workflow, and its checks are the local contract exactly: a change that
   `gofmt -l`, `go vet`, `go mod tidy -diff`, `grove check`, or the test
   suite would reject locally is red on GitHub, and `main` at the candidate
   is green on both runner operating systems.
2. Every workflow file passes a review against the hardening list above:
   read-only default token, SHA-pinned actions, timeouts, concurrency, no
   `pull_request_target`. The evidence names each file and each pin.
3. The repository settings for SHA pinning, allowed actions, and Dependabot
   alerts are on, shown by `gh api` output in the evidence, and Dependabot
   has opened or would open grouped `gomod` and `github-actions` PRs.
4. The Claude review workflow does not run for fork PRs or for pushes to an
   open PR, and the `@claude` workflow runs only for write-access actors;
   the evidence shows the trigger conditions and, where a real run exists,
   the run. The unverified claim about the App token's write reach is either
   verified with a citation or the workflow's permissions are documented as
   the only guard.
5. `README.md`'s testing section says that CI runs the same checks, so a
   contributor reading it does not learn two contracts.
6. The owner judges the setup useful and not noisy after the first week of
   pushes and Dependabot PRs; that judgment is theirs and is recorded here.

## Next

Assign. No plan is needed: the design above and `xfer`'s `.github` are enough
for one implementation. The assignee should start from PR #1's branch or
cherry-pick its two files, since merging it as-is would put unpinned,
every-push workflows on `main`. The owner's own steps are the repository
settings and merging the result; both stay theirs. When G-040 selects a
distribution mechanism, its release workflow should reuse this work's
hardening pattern rather than restart it.
