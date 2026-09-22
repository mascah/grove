---
id: "G-081"
type: work
title: "Run secure CI and Dependabot on GitHub"
status: proposed
created: "2026-09-22T16:10:10Z"
updated: "2026-09-22T16:18:32Z"
kind: tooling
size: medium
relates_to: ["G-040", "G-044"]
---

## Outcome

Every push to `main` and every pull request on the public
`github.com/mascah/grove` repository gets the same checks the owner runs
locally, run by GitHub Actions under a hardened configuration, and dependency
and action updates arrive as Dependabot pull requests. The owner asked for
this on 2026-09-22,
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
- The Claude Code GitHub App and its two workflows in
  [PR #1](https://github.com/mascah/grove/pull/1) are dropped for now: the
  owner will close the PR unmerged and uninstall the App. They first chose a
  hardened same-repo, open-only review, then on reading the risk assessment
  below chose to drop it, because the App holds standing write access to a
  public repository and an `@claude` on an outsider's issue or PR reads that
  untrusted text, while every change here is integrated and reviewed locally.
  Revisit when Grove's own pull-request path exists.
- A LICENSE file is out of scope. The owner chose to leave the repository
  unlicensed for now; this is an open item, not a question record.

**Observed on 2026-09-22 at `main` ccdc92d:**

- No `.github` directory exists on `main`. PR #1 adds `claude.yml` and
  `claude-code-review.yml` and nothing else; the `CLAUDE_CODE_OAUTH_TOKEN`
  secret is set (`gh secret list`) and becomes unused once the App goes.
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
- The risk assessment the owner decided on, from
  `anthropics/claude-code-action` at `main` on 2026-09-22: only write-access
  actors can trigger it (`docs/security.md`), so here only the owner; its
  default Bash allowlist is `git add`, `git commit`, a push wrapper and
  `git rm` (`src/modes/tag/index.ts`), so a prompt injection cannot reach the
  network or the OAuth token through Claude's tools; it cannot approve PRs or
  submit reviews (`docs/capabilities-and-limitations.md`). Its reach is
  therefore commits and pushes to branches and comments, through the App's
  standing write access, whenever the owner tags it on untrusted content. PR
  #1 as generated also used floating action tags, reviewed on every push and
  on fork PRs, and set no `timeout-minutes` or `concurrency`.
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
- Not carried over from `xfer`, and why: labeler, stale, PR-title, and
  CODEOWNERS serve a team; release-please and any release build belong to
  G-040; the scheduled security audits that open issues are covered for Go by
  `govulncheck` plus Dependabot alerts; the Claude review and codebase-review
  workflows are dropped by the decision above. If the App returns, the
  hardening that was shaped and then dropped is: `claude-code-review.yml` on
  `types: [opened, ready_for_review]` guarded by
  `if: github.event.pull_request.head.repo.full_name == github.repository`,
  SHA-pinned, with timeouts and concurrency, `contents: read`, and
  `include_comments_by_actor` limited to the owner.

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
4. PR #1 is closed unmerged, the Claude Code GitHub App is no longer
   installed on the repository, and the `CLAUDE_CODE_OAUTH_TOKEN` secret is
   removed; `gh pr view 1`, the repository's installed-apps settings, and
   `gh secret list` show it. These are the owner's own steps.
5. `README.md`'s testing section says that CI runs the same checks, so a
   contributor reading it does not learn two contracts.
6. The owner judges the setup useful and not noisy after the first week of
   pushes and Dependabot PRs; that judgment is theirs and is recorded here.

## Next

Assign. No plan is needed: the design above and `xfer`'s `.github` are enough
for one implementation. The owner's own steps, in any order: close PR #1
unmerged (`gh pr close 1`), uninstall the Claude Code GitHub App and delete
the `CLAUDE_CODE_OAUTH_TOKEN` secret, change the repository settings in
acceptance 3, and merge the result. When G-040 selects a distribution
mechanism, its release workflow should reuse this work's hardening pattern
rather than restart it.
