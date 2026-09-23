---
id: "G-105"
type: work
title: "Restore green CI: Linux build, Go patch, and dependency advisories"
status: review
created: "2026-09-23T04:44:12Z"
updated: "2026-09-23T05:00:10Z"
kind: fix
size: small
relates_to: ["G-081", "G-045", "G-046"]
candidate: "e30f90c6572edf3493e9f2c24504b6477e64053e"
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

## Evidence

Implemented 2026-09-23 on `worktree-G-105` from `main` 8310af1 (the record
at `sha256:ad7123fd…`; no plan, as Next said). Commits after the `active`
update e7dd366:

- `8226c34 fix(attempt): use x/sys Getsid so the package builds on Linux`:
  both call sites use `unix.Getsid`; x/sys becomes a direct requirement.
- `b8f7bb3 build(deps): Go 1.26.8 and module advisories`: exactly the
  trial's `go.mod` (Go 1.26.8, lipgloss 2.0.6, x/ansi 0.11.8, goldmark
  1.8.6, yaml 3.0.5, x/sys 0.48.0, x/net 0.59.0, x/sync 0.23.0, x/text
  0.42.0, ultraviolet 20260811164956).
- `64f0a0c test: make three Linux-only failures independent of wrap and
  timing`, and `1be9030` from the review. The Docker Linux suite, which the
  trial had not run, showed that the Getsid compile error had hidden three
  more ubuntu failures, all in tests: `TestBoardReviewWorkflow` wanted a path
  on the card's line where a short `/tmp` wraps it (now read unwrapped);
  `terminal.py`'s `attempt_lifecycle` advanced its mark past a frame that
  already held `step 19999`; and `TestStopAfterReconnect`'s fake writes its
  start mark before its INT trap, so a started provider can die of SIGINT
  (exit -1, signal `interrupt`) instead of exiting 130, and the check now
  accepts both. No product code changed beyond Getsid.

Verification at `1be9030` on macOS (go1.26.8): `gofmt -l .` empty,
`go vet ./...`, `GOOS=linux go vet ./...`, `go mod tidy -diff`,
`go build ./...`, `go run ./cmd/grove check` OK, `go test -count=1
-timeout 120s ./...` all ok, `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0
./...` "No vulnerabilities found", exit 0; no `syscall.Getsid` under
`internal/`. Linux, `docker run … golang:1.26` (go1.26.8 linux/arm64, the
`AGENTS.md` command): `go vet ./...` and the full suite passed twice at the
tree committed as `64f0a0c`, `TestStopAfterReconnect -count=20` passed, and
`internal/cli` and `internal/tui` passed again at `1be9030`.

Against acceptance: 3 and 4 are met on the branch (`go 1.26.8`, govulncheck
exit 0, no `syscall.Getsid`, Linux vet clean, Docker suite passes
`internal/attempt`). 1, 2 and 5 can only be observed after the merge is
pushed to `main`: the CI run, Dependabot closing PR #5 and #6 and alert #1,
and the owner's push-time notice. Not run: CI itself (nothing was pushed), a
Linux amd64 suite (the reviewer vetted amd64; the suite ran on arm64), and
`-race` on `internal/attempt`.

Review: [G-106](G-106-g-105-green-ci-review.md), independent, examined
`64f0a0c`; nothing blocking, three notes, two fixed in `1be9030`. Open, out
of scope: with a checkout path long enough to wrap the board's review card
three rows deep, the card drops the row naming the integrate target (G-043
and G-044's card).

## Next

In review. Candidate is the commit that set this record's evidence; the
branch is `worktree-G-105`. To judge and integrate:

```sh
# in /Users/mascah/GitHub/mascah/grove/.claude/worktrees/worktree-G-105
go run ./cmd/grove approve G-105 "VERDICT"
# in main's checkout, /Users/mascah/GitHub/mascah/grove
go run ./cmd/grove integrate G-105
git push
gh run list --branch main --limit 1
gh api repos/mascah/grove/dependabot/alerts --jq '.[] | [.number, .state] | @tsv'
gh pr list --state all --limit 5
```

Acceptance 1, 2 and 5 are read from those last three after the push.
