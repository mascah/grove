---
id: "G-105"
type: work
title: "Restore green CI: Linux build, Go patch, and dependency advisories"
status: proposed
created: "2026-09-23T04:44:12Z"
updated: "2026-09-23T04:44:16Z"
kind: fix
size: small
relates_to: ["G-081", "G-045", "G-046"]
---

## Outcome

CI on `main` is green again on both runners, GitHub's Dependabot alert #1 on
`golang.org/x/net` is closed by an upgrade, and `govulncheck` reports no
called vulnerability. The owner asked for this on 2026-09-23 after seeing
"multiple failures in CI" and the push-time notice "GitHub found 1
vulnerability on mascah/grove's default branch (1 moderate)".

## Constraints

**Observed on 2026-09-23 at `main` ff43559** (`gh run list`, `gh run view`,
`gh api repos/mascah/grove/dependabot/alerts`, and local runs):

- CI was green through the G-041 merge (run 35783200988, 2026-09-22 20:53Z)
  and has failed on every `main` push since the G-043 merge (run 35800657182,
  2026-09-23 00:07Z). Three independent causes, in the order they appeared:
  1. **govulncheck, both jobs' first failure.** Four advisories in symbols the
     code calls: GO-2026-6218 (`net/url`) and GO-2026-6088 (`encoding/xml`),
     both fixed in the Go 1.26.6 standard library while `go.mod` pins
     `go 1.26.5` and `setup-go` installs exactly that (a limit G-081 recorded);
     GO-2026-5970 in `golang.org/x/text` v0.24.0, fixed in v0.39.0; and
     GO-2026-5320 in `github.com/yuin/goldmark` v1.7.13, fixed in v1.7.17.
     Reached through `internal/tui` (glamour → x/text, x/net/html) and
     `internal/handoff` (goldmark `ast.Walk`). The same four fail locally with
     `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...`.
  2. **`check (ubuntu-latest)`, since the G-045/G-046 merges** (first seen on
     run 35809197292): `internal/attempt/attempt.go:1012` and
     `attempt_test.go:246` call `syscall.Getsid`, which Go's `syscall` package
     defines on Darwin but not on Linux, so `go vet` fails with
     `undefined: syscall.Getsid`. `GOOS=linux go vet ./internal/attempt/`
     reproduces it on the Mac. G-046's evidence records "Linux not run".
     Once this compile error landed, the govulncheck job stopped reporting
     advisories and fails at "loading packages" instead, which hides cause 1.
  3. **Dependabot alert #1** (GHSA-5cv4-jp36-h3mw, CVE-2026-25680, moderate):
     `golang.org/x/net` < 0.55.0, transitive via
     `glamour → bluemonday → x/net/html`. govulncheck does not flag it because
     no Grove symbol reaches the parser. Dependabot opened PR #5 (x/net to
     0.55.0, a security PR) and PR #6 (grouped: lipgloss 2.0.6, x/ansi 0.11.8,
     goldmark 1.8.6, yaml 3.0.5); both pass the `check` jobs and fail only
     govulncheck on the remaining stdlib and x/text advisories. PR #4 was
     closed by Dependabot as "updatable in another way" when #6 superseded it.
- `golang.org/x/sys` is already an indirect requirement (v0.47.0) and
  `golang.org/x/sys/unix.Getsid(pid int)` exists on both Linux and Darwin.
- go.dev lists Go 1.26.6, 1.26.7 and 1.26.8 as released; 1.27.1 is current.
  `GOTOOLCHAIN` is `auto` on the owner's Mac, so a higher `go` directive
  downloads the toolchain on the next build.

**Trial (observed, disposable clone of ff43559 in the session scratchpad):**
`syscall.Getsid` → `unix.Getsid` in both files with the `golang.org/x/sys/unix`
import, `go mod edit -go=1.26.8`, `go get` of x/net, x/text, goldmark, x/ansi,
yaml and lipgloss at latest, `go mod tidy`. Resulting `go.mod`: Go 1.26.8,
lipgloss 2.0.6, x/ansi 0.11.8, goldmark 1.8.6, yaml 3.0.5, x/sys 0.48.0
(now direct), x/net 0.59.0, x/sync 0.23.0, x/text 0.42.0, ultraviolet
0.0.0-20260811164956. Then `gofmt -l .` empty, `go vet ./...`,
`GOOS=linux go vet ./...`, `go mod tidy -diff`, `go build ./...`,
`go run ./cmd/grove check` (OK: 100 records), `go test -count=1 -timeout 120s
./...` all packages ok, and govulncheck "No vulnerabilities found", exit 0.
Diff: 4 files, 34 insertions, 32 deletions. Not run: the Docker Linux suite
from `AGENTS.md` (vet on Linux was reproduced with `GOOS=linux`; the test
suite ran on macOS only) and a `-race` pass of `internal/attempt`.

**Proposed design (binds nobody):** one branch with two focused commits,
`fix(attempt): use x/sys Getsid so the package builds on Linux` and
`build(deps): Go 1.26.8 and module advisories`, exactly the trial above. Merge
it from `main`'s checkout and push; Dependabot then closes PR #5 and #6 as
already applied, as it did PR #4, and GitHub closes alert #1 on the next
dependency-graph update. Alternative: merge PR #6 then PR #5 first and add only
the Getsid fix, the `go` directive and the x/text bump; same end state, three
merges instead of one. Out of scope: raising the `go` directive to 1.27, a
`go-version: stable` or `check-latest` change in `ci.yml` that would decouple
the runner's toolchain from `go.mod` (G-084 finding 5 recorded that limit;
reopen it in G-081's neighbourhood if patch releases keep breaking the
signal), and CodeQL.

## Acceptance

1. `go vet ./...` and the test suite pass on `ubuntu-latest` and
   `macos-latest`, and the `govulncheck` job passes, on the CI run for the
   `main` push that integrates this work (`gh run list --branch main`).
2. `gh api repos/mascah/grove/dependabot/alerts` shows alert #1 with state
   `fixed` (or `auto_dismissed`), and PR #5 and #6 are closed or merged.
3. `go.mod`'s `go` directive is at least 1.26.6, and
   `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...` on the integrated
   `main` exits 0.
4. `internal/attempt` has no `syscall.Getsid`; `GOOS=linux go vet ./...` is
   clean from a Mac, and the Docker Linux suite in `AGENTS.md` passes
   `internal/attempt`.
5. The owner judges the fix complete when the push-time vulnerability
   notice no longer appears.

## Next

Assign. Small, no plan needed: the trial above is the implementation.
Branch `worktree-G-105` from `main` ff43559 or later. The owner may prefer to
apply the two commits directly on `main`, since the diff is 4 files and CI is
the review; that is their call.
