# W-004 and W-005 coordination plan

Goal: deliver `versions` per [W-004](../../grove/work/W-004-record-versions.md)
and `workspace` per [W-005](../../grove/work/W-005-record-workspace.md), each
work record owning its specification, under the presentation policy in
[Q-001](../../grove/questions/Q-001-branch-versions.md). Schema:
[record model](../record-model.md#on-disk-contract).

Execution: one agent, one worktree `.claude/worktrees/W-004-W-005`, branch
`worktree-W-004-W-005`, sequential: loader over `fs.FS` → Git sources and
inventory → `versions` → W-004 review and closure → `workspace` → joint
fixture → combined review → W-005 closure and documentation. Git 2.50.1,
Go 1.26.2. Reassess if committed-tree reading needs more than `ls-tree` plus
one `cat-file --batch` per branch.

Inspected base: main at `b20d2b0` (merge of `worktree-W-003`), clean tree,
one registered worktree. W-003's interfaces as implemented: `project.Load`
(discovery plus one-checkout validation), `project.ParseRecord`,
`project.Validate`, `project.Revision` (`sha256:` + 64 hex of exact bytes,
the convention `show --json` and `update --expect` use), `repo.CommonDir`
(creates `<common>/grove`, so read-only commands must not call it),
`repo.Git`, and the `update` package's snapshot comparison. `go test ./...`,
`go vet ./...`, gofmt, and `check` (9 records) passed on that base.

## Contract decisions

### Sources (W-004)

- Committed: every `refs/heads/*` tip, read at its observed commit through
  `git ls-tree -r -t -z` and `git cat-file --batch`; the tree is presented to
  the shared loader as an `fs.FS`, so a branch obeys exactly the checkout
  rules (symlinks, misplaced files, schema, graph).
- Live: every registered worktree from `git worktree list --porcelain -z`,
  including detached HEAD, read through `project.Load` at
  `<worktree>/<prefix>`. Bare entries are skipped. Prunable entries, entries
  that cannot report a Git directory, and entries whose common directory is
  another repository are invalid sources.
- Prefix: the selected project's location relative to the worktree root
  (`git rev-parse --show-prefix`), applied to every source. A source without
  `<prefix>grove.yaml` is absent (listed, no records, not an error).
- Live locator: `.` for the main worktree, otherwise the linked worktree's
  administrative name (`<common>/worktrees/<name>`), which Git sanitizes to
  ref-safe characters. Paths remain in JSON and stderr.
- Instability: the worktree inventory is read before and after all live reads;
  any worktree whose path, HEAD, branch, or detached state differs is
  invalid ("changed while being read"). Branch reads are frozen to commits.
- Live change against HEAD: each live record is compared by ID with the
  HEAD commit's validated records: `unchanged`, `modified`, `renamed`
  (same bytes, different path), `added`; records at HEAD missing from live
  files appear as `deleted` rows without content or selector. An absent
  project at HEAD, or an unborn HEAD, is an empty baseline. A HEAD whose
  project does not validate cannot be compared with: the live source stays
  valid and selectable, each of its versions has `change: unknown`, and the
  source carries a `note` saying why (revised after review; discarding a
  valid checkout because an old commit is broken was judged less useful).
- Completeness: any invalid source sets `complete: false` and exit 1 while
  valid sources still print. Inventory failure (no Git, not a repository,
  listing errors) is an error with no view. An ID absent from every valid
  source is `record X not found in any valid source`, exit 1.

### Selector (W-004, consumed by W-005)

```
committed:<ref>@<commit12>:<id>@<rev12>:<binding16>
live:<locator>:<ref|detached>@<commit12>:<id>@<rev12>:<binding16>
```

`commit12` and `rev12` are the first twelve hex digits of the observed commit
and of the record's `sha256:` content revision; they attribute the common
staleness causes. `binding16` is the first sixteen hex digits of the SHA-256
over the NUL-joined identity tuple: repository common directory, prefix,
kind, ref (or `detached`), worktree path and Git directory (live only), full
commit, configuration revision, record ID, record path, full content
revision. The readable part omits the prefix; two projects in one repository
differ only in the digest. Refs
and locators cannot contain `:`; hex cannot contain `@`, so the grammar splits
on `:` and the last `@`. Content equality alone never produces the same
selector in two sources.

### JSON (W-004)

`versions [ID] --json` prints one object:

```
project, repository, prefix, complete,
sources: [{kind, ref|null, commit, worktree?, locator?, detached?, present, valid, config_revision?, note?, diagnostics: [...]}],
records: [{id, versions: [{selector?, kind, ref|null, commit, worktree?, locator?, path, revision?, config_revision, type, title, status, change?, head_path?, source?}]}]
```

Groups are ordered by ID prefix (W, Q, D) then numeric suffix; versions
within a group are committed sources by ref, then live sources with the main
worktree first and linked worktrees by locator. Human output is one row per
version (`ID STATUS SOURCE CHANGE SELECTOR`) on stdout with the source table
and diagnostics on stderr; control characters are escaped in human output and
preserved in JSON.

### Workspace (W-005)

`workspace --source SELECTOR [--json]` re-runs the W-004 inspection for the
selector's record immediately before answering. A live selector succeeds only
when the current live observation reproduces the selector exactly. A
committed selector succeeds only when the branch still points at the observed
commit, exactly one enterable registered worktree has that branch checked
out (a prunable or foreign entry is not a checkout), that checkout's live
`grove.yaml` has the committed configuration bytes, and its live record has
the committed bytes at the committed path. Success
prints the absolute project directory inside the checkout; JSON carries
`checkout, project, record, ref, head, revision, selector`. Distinct
refusals: malformed selector (exit 2); worktree or branch no longer present;
no checkout of the branch; more than one checkout of the branch; branch or
HEAD moved; attached/detached changed; record missing or moved; content
changed (committed selection: refresh and select the live observation);
configuration or project location changed; source invalid. The command
creates nothing and never calls `repo.CommonDir`.

## Ordered steps

1. [x] W-004 active. Loader: `project.LoadFS(fs.FS)`; `Load` wraps
   `os.DirFS`; `Project.Config` retains configuration bytes;
   `repo.Locate` finds the common directory without creating anything.
2. [x] `internal/versions`: tree FS from Git objects, inventory, live
   comparison, selectors, deterministic ordering, `Inspect(root, id)`.
3. [x] CLI `versions [ID] [--json]`, usage, stderr context, exit codes.
4. [x] W-004 fixtures (table below), suites, race, vet, gofmt, `check`,
   Git-state hashes, independent review, fixes, evidence, W-004 done.
5. [x] `internal/versions/workspace.go` (kept in the versions package to
   share fixtures and the inspection): selector parsing, resolution,
   attribution.
6. [x] CLI `workspace --source SELECTOR [--json]`.
7. [x] W-005 fixtures, joint fixture, suites, combined review, fixes,
   evidence, W-005 done; README, model, brief reconciled.

## Acceptance to checks

| Item | Check |
| --- | --- |
| W-004: main/feature statuses and bodies, committed and live labels | `TestInspectMainAndFeature` |
| W-004: branch without checkout; detached, dirty, untracked, deleted, renamed, identical | `TestInspectBranchWithoutCheckoutAndDetached`, `TestInspectLiveChanges` |
| W-004: source-specific configuration, project below root, dependency missing in one source | `TestInspectPrefixAndConfig`, `TestInspectSourceLocalValidation` |
| W-004: invalid YAML, duplicate IDs, inaccessible worktree, changing identity | `TestInspectIncomplete`, `TestInspectUnstable`, `TestVersionsIncompleteAndNotFound` |
| W-004: unborn HEAD, HEAD whose project is invalid | `TestInspectUnbornAndInvalidHEAD` |
| W-004: bytes and revisions agree with BOM/CRLF; selectors differ for identical content | `TestInspectBytesAndSelectors` |
| W-004: paths with spaces, tabs, newlines; repeated reads order consistently | `TestInspectPathsAndRepeatedReads`, `TestVersionsCLI` |
| W-004: refs, index, records, allocator state unchanged; existing suites | `TestVersionsLeavesGitUnchanged`, whole suite |
| W-005: live selection from main, project below root, no Git change | `TestWorkspaceLive` |
| W-005: committed selection resolves; differing live content refuses | `TestWorkspaceCommitted` |
| W-005: branch/HEAD, detached, moved/removed, configuration, record changes refuse | `TestWorkspaceStale` |
| W-005: missing versus ambiguous checkouts; live selection disambiguates | `TestWorkspaceMissingAndAmbiguous` |
| W-005: dirty files survive; paths round-trip through JSON | `TestWorkspaceLive`, `TestWorkspaceCLI` |
| W-005: joint fixture list → select → resolve → `show --project` | `TestJointWorkflow` |

## Integration checks

At each boundary: `gofmt -l .`, `go vet ./...`, `go test ./...`,
`go test -race ./...`, `go run ./cmd/grove check`, and real use of
`versions` and `workspace` against this repository's own worktrees.

## Recovery checkpoints

Resume from the commit list under Progress. Before continuing, run
`git status`, `git log --oneline main..`, and the integration checks; then
compare the unchecked steps above with the diff. W-004 is closed by its own
commit before any W-005 file exists.

## Progress and evidence

### W-004

Commits on `worktree-W-004-W-005`: `5c08930` (plan, W-004 active),
`2b688f7` (loader over `fs.FS`, `repo.Locate`), `b5125fc` (`versions`),
`ca420f5` (review fixes). Base `b20d2b0`.

- `gofmt -l .` clean, `go vet ./...`, `go test ./...`, and
  `go test -race ./...` pass for every package; `go run ./cmd/grove check`
  reports 9 records.
- `internal/versions` fixtures: main/feature statuses and bodies from each
  branch with four ordered versions, distinct selectors for identical bytes,
  the same selectors from either checkout; a branch without a checkout and a
  detached worktree; modified, deleted (row without content or selector),
  added, renamed (with `head_path`), and moved-and-changed records, plus a
  checkout whose HEAD has no project; a project below the repository root
  whose branches configure different record folders, with configuration
  revisions hashing exact `grove.yaml` bytes and a branch without the project
  listed as absent; a dependency present only on main leaving the feature's
  committed and live sources invalid without admitting their records; invalid
  YAML, duplicate IDs, a symlinked configuration, and a prunable worktree
  each attributed to their source with the result incomplete and valid
  sources still contributing; a HEAD move, worktree removal, detachment, and
  a newly registered worktree during reading each marked unstable; BOM plus
  CRLF bytes and revisions agreeing in all four sources; a worktree path with
  newline, tab, and space and a record path with a space round-tripping, with
  three repeated reads ordering identically; an unborn HEAD listing records
  as added; an invalid HEAD project noted with `change: unknown`; a plain
  directory refused.
- `internal/cli` fixtures: usage errors exit 2; table rows and stderr source
  lines; JSON fields for live and committed versions; a control character in
  a worktree path preserved in JSON and escaped in text; incomplete results
  exit 1 with valid rows printed; not-found exit 1; an invalid current
  checkout attributed once as a live source; `TestVersionsLeavesGitUnchanged`
  hashing both checkouts, including `.git`, across text and JSON reads with
  no `.git/grove` created.
- Real use in this worktree: `versions W-004` showed main `proposed` and this
  branch `active`, four rows, no authority chosen; with an uncommitted edit
  the live row became `modified` with a new selector, and every file under
  the repository's `.git` hashed identically before and after.
- Independent review (reviewer agent, `2b688f7..b5125fc`, read-only with
  scratch repositories): no blocking findings. Should-fix, all addressed in
  `ca420f5`: an unborn HEAD made the only live source invalid (now an empty
  baseline); an invalid HEAD project was undocumented and untested (contract
  amended above, test added); the Git-state test named in this plan did not
  exist (added). Nits taken: the test hook left the exported API, commit
  abbreviation is guarded, plan names and the not-found wording match the
  code. Nit recorded: the readable selector omits the prefix. The reviewer
  confirmed identical validation for trees and checkouts, no record leakage
  from invalid sources, no side-effecting Git subcommands, and deterministic
  ordering.
- Limits: no fetch, remotes, tags, or history; separate clones are separate
  repositories; a source can change after it was read, which selectors
  expose at resolution time rather than prevent.

### W-005

Commits: `77cb1ce` (`workspace`, W-005 active), `d232aa2` (review
fixes). Prepared against W-004's closed interface at `a465d88`.

- `gofmt -l .` clean, `go vet ./...`, `go test ./...`, and
  `go test -race ./...` pass for every package; `check` reports 9 records.
- `internal/versions` fixtures (`workspace_test.go`): selector grammar
  accepts every form `versions` prints, including a branch name containing
  `@` and the `.` locator, and rejects fourteen malformed forms; from main
  with the project below the repository root, a live feature selection
  returns the exact checkout, project, record path, branch, HEAD, and
  revision while both checkouts, their staged and unstaged files, and Git
  state hash identically afterwards; main's own live and committed versions
  resolve to main; a matching committed selection resolves to its checkout
  and reports the live selector; an edited live record refuses the committed
  selection with "select the live observation" and the old live selection
  with "changed since it was selected", the fresh live selection resolves,
  and a deleted record refuses both; a HEAD move refuses live and committed
  selections with distinct attribution; detaching, re-attaching, renaming
  the record, changing `grove.yaml`, moving the worktree, removing it, and
  deleting the branch each refuse with their own reason, and restoring the
  state resolves again; an invalid current checkout is reported as an
  invalid source; a branch without any checkout refuses without creating
  one, two forced checkouts of one branch refuse as ambiguous naming both,
  and either live selection resolves.
- `internal/cli` fixtures (`workspace_test.go`): nine usage errors exit 2;
  help works without a project; stdout is exactly the project directory
  with checkout, branch or detached HEAD, record, and revision on stderr;
  JSON fields for a detached worktree whose path contains a newline, with
  the path preserved in JSON and escaped on stderr; all three checkouts hash
  identically afterwards; three refusals exit 1 with nothing on stdout; an
  invalid current checkout still resolves a selection elsewhere; a plain
  directory is refused. `TestJointWorkflow`: `versions --json` lists W-001,
  the live feature version is chosen explicitly, `workspace --json` returns
  its project and revision, and `show --project` there returns the exact
  selected bytes while main's copy and branch are untouched.
- Real use in this worktree: `versions W-005` printed four rows; the live
  selector of this worktree resolved to it, `show --project` with that path
  read W-005 here, and main's committed selection routed to the main
  checkout with its own revision and the main live selector in JSON.
- Combined independent review (reviewer agent, `main..77cb1ce`, read-only
  with scratch repositories): no blocking findings. Should-fix, all
  addressed in the review-fix commit: an absent project at the prefix
  produced "not a valid source:" with no reason; a prunable duplicate
  checkout counted as ambiguous and was named by an empty string; the
  committed route did not compare the checkout's configuration bytes
  (W-005 acceptance said changed configuration refuses; the plan's
  committed condition now says so too). Nits taken: the hand-rolled hex
  encoder is gone; `Resolve` carries a `ponytail:` note on re-running the
  whole inspection. Nit recorded, not verified: no fixture resolves a
  selector from a second Grove project at another prefix in the same
  repository; the prefix is in the binding and would refuse. The reviewer
  confirmed the binding is recomputed at resolution time, attribution
  order, exact `<worktree>/<prefix>` routing, control-character paths,
  only read-only Git subcommands with byte-identical `.git` and checkouts
  across `versions` and `workspace` runs, and every W-005 acceptance
  bullet exercised by a failing-if-broken test.

### Joint verification

`TestJointWorkflow` runs the connected workflow through the CLI: list
versions of W-001, choose the live feature version explicitly, resolve its
workspace, and read the exact selected bytes there with `show --project`,
with main's copy and branch untouched. The same sequence was run by hand in
this repository against W-005 (see W-005 real use). Final checks on the
closing commit: `gofmt -l .` clean, `go vet ./...`, `go test ./...`,
`go test -race ./...`, `check` (9 records), and a documentation link check
over README, model, brief, plan, Q-001, W-004, and W-005.
