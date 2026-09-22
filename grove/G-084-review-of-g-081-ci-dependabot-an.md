---
id: "G-084"
type: review
title: "Review of G-081 CI, Dependabot, and Go 1.26.5"
status: current
created: "2026-09-22T16:38:25Z"
updated: "2026-09-22T16:41:41Z"
work: ["G-081"]
examined: "64ad5d3"
---

## Examined

Branch `worktree-G-081`, base `main` 6c11ad8. Round 1 examined d4f5360
(`git diff main..d4f5360`: `go.mod` to Go 1.26.5, `.github/workflows/ci.yml`,
`.github/dependabot.yml`, README). Round 2 examined 64ad5d3, the fix commit,
which is the `examined` field; the candidate differs from it only by this
record and [G-081](G-081-github-ci.md)'s evidence.

## Review

Independent reviewer: a fresh Claude Code reviewer agent (Fable 5.1),
read-only, working from the worktree with `gh api` and scratchpad clones. It
verified both action SHAs against their tags, the hardening list line by line,
the check set against `AGENTS.md` and `lefthook.yml`, the gofmt guard against
an unformatted file, and the suite plus `grove check` from a credential-free
shallow clone (`GIT_CONFIG_GLOBAL=/dev/null`), all green. It read
`actions/setup-go`'s installer at the pinned SHA to confirm `go 1.26.5`
installs exactly 1.26.5, and GO-2026-4970's affected symbols against
`os.OpenRoot` calls in `internal/handoff`.

## Findings

Round 1, no blocking findings:

1. should-fix: `cancel-in-progress: true` also cancelled `main` runs, so a
   push during another push's run left the earlier commit without a verdict.
2. should-fix: govulncheck is symbol-level; nine known advisories in required
   modules the code does not call leave the job green, so "a new advisory
   fails visibly" holds only for called symbols.
3. should-fix: acceptance 3's repository settings were still off at review
   time (`allowed_actions: all`, `sha_pinning_required: false`,
   vulnerability alerts 404).
4. nit: `gofmt -l` writes a parse error to stderr and exits 2 with empty
   stdout, so the guard passed on a lone syntax error.
5. nit: `go-version-file: go.mod` installs exactly 1.26.5 while go.dev lists
   up to 1.26.8; nothing bumps the `go` directive automatically.
6. nit: the README sentence did not name `grove check` and implied
   govulncheck ran on both operating systems.
7. nit: `go build ./...` adds no coverage beyond vet, test, and `grove check`.
8. nit: `actions/checkout` persisted the job token although no job pushes.

Round 2 (64ad5d3), no blocking findings, no regressions: fixes 1, 4, 6, 8
confirmed, each reproduced (the concurrency expression is a documented
boolean form; the guard now fails on a syntax error, an unformatted file, and
passes clean; the README matches clause by clause; a credential-free clone
runs the suite and `grove check`). One residual nit accepted: with
`cancel-in-progress` false on `main`, a third push inside one run still
cancels the pending middle one; a per-SHA group on push would close it.

## Disposition

- 1, 4, 6, 8: fixed in 64ad5d3.
- 2, 5: recorded as limits in G-081's evidence; Dependabot alerts (acceptance
  3) cover uncalled advisories, and the toolchain moves when go.mod is edited.
- 3: the owner's steps, with runnable commands in G-081's Next.
- 7: kept; the record's design asks for `go build` by name.
- Round 2 residual nit: accepted as is; noted in G-081 for a later change.
