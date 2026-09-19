# W-006 workspace provenance implementation plan

> For Fable: implement one task at a time in an isolated worktree using the
> repository instructions and an inline execution workflow (for Codex,
> `superpowers:executing-plans`). Do not dispatch or merge automatically.

Review base: `9b7f730`; include the review/shaping documents in the execution
base. Read the linked work record first; it owns outcome, constraints, and
acceptance. This is a proposed technical repair plan, not implementation evidence.
Go 1.26+, current standard library and existing YAML dependency only.

**Goal:** refuse foreign project routing and workspace selections invalidated
by an observed checkout change.

**Architecture:** retain `versions.Inspect(root, id)` and
`versions.Resolve(root, selector)` and their outputs. Factor selected live-source
validation into `internal/versions/live.go`, reused by inspection and final
resolution. Keep Git discovery in `internal/repo` and schema validation in
`internal/project`; do not give either reader write responsibilities.

**Spec:** [W-006](../../grove/work/W-006-workspace-provenance.md).

## Global constraints and review focus

No selector/schema/CLI changes, new dependency, Git/filesystem writes, UI,
worktree creation, or automatic reselection. Preserve incomplete results and
source-local graph validation. Review: intermediate prefix symlinks, another
repository at an ordinary prefix directory, prunable entries after first read,
committed-route ambiguity arising during resolution, and mutations of an early
source while later sources are inspected. The tasks below cover each case.

## Task 1: Validate the entire live project location

Files: `internal/versions/live.go` (new), `versions.go`, `versions_test.go`.
Consumes `repo.Locate`, worktree inventory, `project.Load`, and selected prefix.
Produce an unexported helper that takes expected repository, worktree entry and
prefix, and returns a validated `Source` or its absent/invalid diagnostics;
keep exported signatures unchanged. Reuse the helper again in Task 2.

- [ ] Add a fixture using existing test helpers:

```go
root := repoFixture(t)
// Move grove.yaml and grove/ into root/sub; commit that nested project.
os.Mkdir(filepath.Join(root, "sub"), 0755)
os.Rename(filepath.Join(root, "grove.yaml"), filepath.Join(root, "sub/grove.yaml"))
os.Rename(filepath.Join(root, "grove"), filepath.Join(root, "sub/grove"))
commit(t, root, "nested project")
wt := addWorktree(t, root, "feature", "", "-b", "feature")
outside := repoFixture(t)
os.RemoveAll(filepath.Join(wt, "sub"))
os.Symlink(outside, filepath.Join(wt, "sub"))
res := mustInspect(t, filepath.Join(root, "sub"), "W-001")
if res.Complete || source(t, res, "live", "feature").Valid {
    t.Fatal("external project admitted as a feature source")
}
```

Check every setup error in the retained test. Add table cases for a middle
prefix symlink, an internal symlink, a dangling symlink, a regular file, a
nested foreign Git repository without a symlink, and a missing directory.
Missing must be absent; the others invalid; main's valid records remain visible.
- [ ] Run `go test -count=1 -run TestInspectProjectLocation ./internal/versions`
  and retain the foreign-project failure before implementation.
- [ ] Implement the helper with this order:

```text
Verify registered worktree Git/common identity.
Walk prefix components with Lstat from that worktree (no EvalSymlinks shortcut).
Missing component => absent; non-directory/symlink/read error => invalid.
Verify Git identity from the actual project directory equals the expected
common directory, worktree Git directory and prefix.
Lstat grove.yaml: missing => absent; otherwise load the entire source.
Retain exact configuration/content revisions and existing source diagnostics.
```

A configured path may not escape the selected checkout even if file bytes match.
Keep unborn HEAD and invalid-HEAD comparison handling from the existing reader.
- [ ] Run the new tests plus `TestInspectPrefixAndConfig`,
  `TestInspectIncomplete`, `TestInspectSourceLocalValidation` and commit
  `fix(versions): verify live project ownership through its prefix`.

## Task 2: Reject observed instability before returning a workspace

Files: `versions.go`, `workspace.go`, `live.go`, both versions test files,
`internal/cli/workspace_test.go`.
Consumes Task 1's helper and existing exact selector generation; outputs remain
`(*Result, error)` and `(*Workspace, error)`.

- [ ] Retain this deterministic base-failing regression:

```go
func TestInspectPrunableDuringRead(t *testing.T) {
    root := repoFixture(t)
    wt := addWorktree(t, root, "feature", "", "-b", "feature")
    res, err := inspect(root, "W-001", func() {
        if err := os.RemoveAll(wt); err != nil { t.Fatal(err) }
    })
    if err != nil { t.Fatal(err) }
    if res.Complete || source(t, res, "live", "feature").Valid {
        t.Fatal("deleted checkout remained selectable")
    }
}
```

Add a package-private resolve hook, matching the existing `inspect` pattern,
to change the selected target before the final check. Cover config/content/path,
prunable/foreign ownership, branch/HEAD, and a new forced duplicate checkout for
a committed route. A live explicit selection must still disambiguate duplicates.
- [ ] Run failing tests before edits. At CLI level, a temporary Git wrapper can
  delete the disposable feature directory before the second `worktree list`,
  then exec real Git; require exit 1, no success stdout, attributable diagnostics.
  This reproduced exit 0 at the reviewed base. Do not use timing sleeps or delete
  any checkout outside the fixture.
- [ ] Compare administrative validity as well as path/ref/HEAD at the second
  inventory. Re-enter affected sources to detect foreign/inaccessible paths.
  In `Resolve`, perform final target validation using the helper and recompute
  its selector; for committed sources re-read the branch tip and enterable
  checkout mapping too. Reject any mismatch; do not return cached target bytes.
- [ ] Test unrelated invalid sources do not block a healthy explicit live
  selection; retain no-checkout versus ambiguous diagnostics. Hash all fixture
  Git state and record/config/dirty files around both success and refusal.
  Run full `internal/versions` and relevant CLI tests; commit
  `fix(workspace): revalidate the selected checkout before returning`.

## Verification and return

- [ ] Run targeted package/CLI regressions with `-count=1`, then
  `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  `gofmt -l cmd internal`, and `go run ./cmd/grove check`.
- [ ] Obtain an independent review of the final diff, with the original
  reproducer run against both base and candidate. Preserve concrete failing
  and passing outputs; test names or a done status alone are not evidence.
- [ ] Reconcile the work record, README/model if behavior wording changes,
  and the restart brief's concrete next action. Change fields with
  `go run ./cmd/grove update`; keep evidence in the record/plan bodies.
- [ ] Return focused Conventional Commits, tested candidate revision, review
  disposition, and remaining limits. Do not merge or remove old worktrees.
