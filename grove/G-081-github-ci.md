---
id: "G-081"
type: work
title: "Run secure CI and Dependabot on GitHub"
status: review
created: "2026-09-22T16:10:10Z"
updated: "2026-09-22T17:18:39Z"
kind: tooling
size: medium
relates_to: ["G-040", "G-044"]
candidate: "81cd0c0"
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

## Evidence

Implemented 2026-09-22 on `worktree-G-081` from `main` 6c11ad8, starting from
this record at revision ba0e9175 with no plan (the record said none was
needed). Commits: d8aae4e (`go.mod` to Go 1.26.5), d4f5360 (CI workflow,
Dependabot, README), 64ad5d3 (review fixes), then the commit carrying this
evidence and [G-084](G-084-review-of-g-081-ci-dependabot-an.md), which is the
`candidate`.

**Against the acceptance:**

1. `.github/workflows/ci.yml` runs on `pull_request` and `push` to `main`.
   The `check` job on `ubuntu-latest` and `macos-latest` runs, in order,
   `test -z "$(gofmt -l . 2>&1 | tee /dev/stderr)"`, `go vet ./...`,
   `go mod tidy -diff`, `go build ./...`, `go run ./cmd/grove check`, and
   `go test -count=1 -timeout 120s ./...`. The gofmt guard was shown to exit
   1 on an unformatted file and on a syntax error, and 0 on the clean tree.
   Not verified: a run on GitHub, since this session does not push; green on
   both runners is the integrator's observation after the push (Next).
2. Hardening, one workflow file: `permissions: contents: read` at top level
   and no job addition; `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1`
   (v7.0.1) and `actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e`
   (v7.0.0), both resolved from their tags with `gh api` and re-verified by
   the reviewer; `timeout-minutes` 15 on `check` and 10 on `govulncheck`;
   `concurrency` on `${{ github.workflow }}-${{ github.ref }}` with
   `cancel-in-progress` only for pull requests, so each push to `main` keeps
   its result; no `pull_request_target`; `persist-credentials: false` on
   both checkouts. `actionlint` v1.7.12 reports nothing.
3. Settings are unchanged and remain the owner's step; observed at review:
   `allowed_actions: all`, `sha_pinning_required: false`, vulnerability
   alerts 404, automated security fixes disabled. `.github/dependabot.yml`
   groups `gomod` and `github-actions` weekly, one PR each; it omits
   `labels`, relying on Dependabot's default `dependencies` label, since the
   repository has no such label and a listed label that does not exist is
   ignored. The first Dependabot PR confirms both.
4. Observed: `gh pr view 1` is `CLOSED`, `gh secret list` is empty. The App
   installation cannot be read with the user token; the owner confirms it in
   the installed-apps settings.
5. README's testing paragraph now says CI runs the same checks plus
   `grove check` and `go build` on both systems and `govulncheck` on Ubuntu,
   as a signal, and that Dependabot opens grouped PRs.
6. The owner's judgment after the first week; not yet given.

**Decisions:**

- `go.mod` now requires Go 1.26.5. At 1.26.2, govulncheck exits 3 on
  GO-2026-4970 (an `os.Root` escape through a symlink plus trailing slash),
  reached from `readConfined` in `internal/handoff`, the confined reader
  behind `grove context`. The toolchain is the fix; `GOTOOLCHAIN=auto`
  downloads it on the next build. Limit: `setup-go` installs exactly the
  `go` directive, so CI's toolchain moves only when `go.mod` is edited.
- govulncheck runs as `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0`,
  pinned through the module proxy rather than a third-party action. It is
  symbol-level: advisories in modules the code does not call leave it green
  (nine such at the candidate), which Dependabot alerts cover once enabled.
- Not added: CodeQL (optional in the design; add when the owner wants a
  second scanner), an `anthropics/*` allow pattern (no Anthropic action is
  used now), and a `labels` key.

**Verification** at the Go tree of 64ad5d3 (identical to d4f5360), on
macOS with Go 1.26.5: `gofmt -l .` empty, `go vet ./...`, `go mod tidy
-diff`, `go build ./...`, `go run ./cmd/grove check` (`OK: 80 records`;
81 with G-084), `go test -count=1 -timeout 120s ./...` all packages ok,
`internal/versions` 7.3 s; `govulncheck` 0 called vulnerabilities, exit 0;
`actionlint` clean. The reviewer repeated the suite and `grove check` from
a credential-free shallow clone of 64ad5d3, all green.

**Review:** G-084, two rounds, no blocking findings. Round 1's should-fix
items were the `main` cancellation (fixed), govulncheck's symbol-level scope
(recorded above), and the settings still off (owner's step); nits fixed:
gofmt parse errors, README wording, persisted credentials; kept: `go build`
by the design's request. Round 2 confirmed the fixes with one accepted
residual: a third push to `main` inside one run still cancels the pending
middle run; a per-SHA group on push would close that if it ever matters.

**Second candidate, 2026-09-22.** The owner asked for a pull request, so
`worktree-G-081` was pushed and PR #2 opened. That push ran lefthook's
pre-push hook, whose test fixtures inherited the `GIT_DIR` Git exports to
hooks and committed into this repository: the branch tip, HEAD, four stray
branches, two temp worktrees, and `core.bare` were damaged and the fixture
tip reached origin. The owner repaired it by hand; the cause and fix are
[G-089](G-089-ignore-ambient-git-environment-w.md). The repaired PR then ran
CI for the first time: `check (macos-latest)` and `govulncheck` passed,
`check (ubuntu-latest)` failed on two tests with macOS assumptions, fixed in
d367972 and shown passing on Linux with the whole suite in
`docker run golang:1.26` (git 2.47.3, go1.26.8) before pushing again. Until
G-089 is merged, every push from a linked worktree must use `--no-verify`.

## Next

In Review, second candidate. PR #2 is open from `worktree-G-081`; its CI run
on this candidate is the evidence for acceptance 1 (`gh pr checks 2`). The
integrator's steps, from the main checkout:

1. Merge: `main` moved past the base, so `git merge worktree-G-081` makes a
   merge commit, as today's G-079 and G-040 merges did; then
   `git push --no-verify origin main` and `gh run watch` until green. Merge
   [G-089](G-089-ignore-ambient-git-environment-w.md) next, so the hook is
   safe again.
2. Settings (acceptance 3), then read them back:

   ```sh
   gh api -X PUT repos/mascah/grove/actions/permissions -F enabled=true -f allowed_actions=selected -F sha_pinning_required=true
   printf '{"github_owned_allowed":true,"verified_allowed":true,"patterns_allowed":[]}' | gh api -X PUT repos/mascah/grove/actions/permissions/selected-actions --input -
   gh api -X PUT repos/mascah/grove/vulnerability-alerts
   gh api -X PUT repos/mascah/grove/automated-security-fixes
   gh api repos/mascah/grove/actions/permissions; gh api repos/mascah/grove/actions/permissions/selected-actions; gh api -i repos/mascah/grove/vulnerability-alerts | head -1
   ```

3. Confirm the Claude Code GitHub App is gone at
   https://github.com/settings/installations (acceptance 4).
4. Write the verdict into this record with the `gh api` output, then on
   `main`: `go run ./cmd/grove update G-081 --expect REVISION --set status=done`
   and commit. Acceptance 6 is recorded after the first week.
