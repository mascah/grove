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
  files appear as `deleted` rows without content or selector. A HEAD whose
  project does not validate makes the comparison impossible and the live
  source invalid; an absent project at HEAD is an empty baseline.
- Completeness: any invalid source sets `complete: false` and exit 1 while
  valid sources still print. Inventory failure (no Git, not a repository,
  listing errors) is an error with no view. An ID absent from every valid
  source is `record X not found in any inspected source`, exit 1.

### Selector (W-004, consumed by W-005)

```
committed:<ref>@<commit12>:<id>@<rev12>:<binding16>
live:<locator>:<ref|detached>@<commit12>:<id>@<rev12>:<binding16>
```

`commit12` and `rev12` are the first twelve hex digits of the observed commit
and of the record's `sha256:` content revision; they attribute the common
staleness causes. `binding16` is the first sixteen hex digits of the SHA-256
over the NUL-joined identity tuple: repository common directory, prefix,
kind, ref (or `detached`), worktree Git directory (live only), full commit,
configuration revision, record ID, record path, full content revision. Refs
and locators cannot contain `:`; hex cannot contain `@`, so the grammar splits
on `:` and the last `@`. Content equality alone never produces the same
selector in two sources.

### JSON (W-004)

`versions [ID] --json` prints one object:

```
project, repository, prefix, complete,
sources: [{kind, ref|null, commit, worktree?, locator?, detached?, present, valid, config_revision?, diagnostics: [...]}],
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
commit, exactly one registered worktree has that branch checked out, and that
checkout's live record has the committed bytes at the committed path. Success
prints the absolute project directory inside the checkout; JSON carries
`checkout, project, record, ref, head, revision, selector`. Distinct
refusals: malformed selector (exit 2); worktree or branch no longer present;
no checkout of the branch; more than one checkout of the branch; branch or
HEAD moved; attached/detached changed; record missing or moved; content
changed (committed selection: refresh and select the live observation);
configuration or project location changed; source invalid. The command
creates nothing and never calls `repo.CommonDir`.

## Ordered steps

1. [ ] W-004 active. Loader: `project.LoadFS(fs.FS)`; `Load` wraps
   `os.DirFS`; `Project.Config` retains configuration bytes;
   `repo.Locate` finds the common directory without creating anything.
2. [ ] `internal/versions`: tree FS from Git objects, inventory, live
   comparison, selectors, deterministic ordering, `Inspect(root, id, between)`.
3. [ ] CLI `versions [ID] [--json]`, usage, stderr context, exit codes.
4. [ ] W-004 fixtures (table below), suites, race, vet, gofmt, `check`,
   Git-state hashes, independent review, fixes, evidence, W-004 done.
5. [ ] `internal/workspace`: selector parsing, resolution, attribution.
6. [ ] CLI `workspace --source SELECTOR [--json]`.
7. [ ] W-005 fixtures, joint fixture, suites, combined review, fixes,
   evidence, W-005 done; README, model, brief reconciled.

## Acceptance to checks

| Item | Check |
| --- | --- |
| W-004: main/feature statuses and bodies, committed and live labels | `TestInspectMainAndFeature` |
| W-004: branch without checkout; detached, dirty, untracked, deleted, renamed, identical | `TestInspectBranchWithoutCheckout`, `TestInspectLiveChanges` |
| W-004: source-specific configuration, project below root, dependency missing in one source | `TestInspectPrefixAndConfig`, `TestInspectSourceLocalValidation` |
| W-004: invalid YAML, duplicate IDs, inaccessible worktree, changing identity | `TestInspectIncomplete`, `TestInspectUnstable` |
| W-004: bytes and revisions agree with BOM/CRLF; selectors differ for identical content | `TestInspectBytesAndSelectors` |
| W-004: paths with spaces, tabs, newlines; repeated reads order consistently | `TestInspectPaths`, `TestVersionsCLI` |
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

(filled in as steps complete)
