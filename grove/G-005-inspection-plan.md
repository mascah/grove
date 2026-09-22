---
id: "G-005"
type: plan
title: "G-003 inspection implementation plan"
status: current
formerly: "docs/plans/W-001-inspection.md"
work: ["G-003"]
created: "2026-09-19T15:19:41Z"
updated: "2026-09-21T21:11:15Z"
---

# G-003 inspection implementation plan

Goal: deliver the read-only Go CLI described by [G-003](G-003-inspect-records.md).
Spec: [record model](../docs/record-model.md), including its first-reader proposals.

Execution: single agent, sequential reader → graph validation → command integration.
Reason: commands consume one shared loader and validation result.
Delegation: none; review the completed diff against acceptance.
Runtime: interactive session in the current checkout; preserve the existing
uncommitted documentation and records. No concurrent implementation or installation.
Reassess: only if the outcome requires scope beyond inspection.

Baseline: `fc9bb26` plus the accepted, uncommitted documentation and four records.
Go 1.26.2 is installed; no application code or tests exist at preparation time.

## Implementation decisions

- Adopt the documented discovery, type-folder, optional-field, graph, and command
  proposals for this first reader; no readiness/acceptance inference from prose.
- Standard-library command parsing and tabular output; no CLI framework.
- Go module `github.com/mascah/grove`, Go 1.26; YAML via the YAML organization's
  stable v3 API. Pin the resolved patch version in `go.mod`/`go.sum`.
  [Dependency documentation](https://pkg.go.dev/go.yaml.in/yaml/v3).
- Decode YAML nodes to enforce schema types without scalar coercion. Reject
  aliases, merge keys, custom tags, duplicate/unknown keys, null optionals, and
  extra YAML documents. Support ordinary block/flow lists and quoted/block text.
- Require opening/closing `---` lines for frontmatter; tolerate UTF-8 BOM and
  CRLF. Keep the original file bytes for `show`. Reject invalid UTF-8.
- Absent optional fields remain unspecified; never synthesize priority or size.
- Resolve `--project` relative to invocation directory; otherwise find the nearest
  config, stopping at a `.git` file/directory boundary (including linked worktrees).
  Do not depend on a Git executable, run hooks, or read shared counter state.
- Reject symlinked config/root components and entries in the record tree. Inspect
  only regular `.md` files under the three type folders; report misplaced Markdown.
- Return stable path/field diagnostics. A bad config stops traversal; independent
  record errors accumulate. Relationships resolve only after all records load.
- Detect member and dependency cycles independently. Shared/nested membership is
  valid; only unresolved/wrong-type targets, duplicates, self-links, or cycles fail.
- `list` orders dated records first by creation time, then ID prefix and numeric
  suffix (without an integer-size limit); undated records last. No priority sort.
- CLI syntax: `grove [--project DIR] list|show ID|check`; also allow `--project`
  after the command. Context goes to stderr; list/source/check result to stdout.
  Exit 0 for success/help, 1 for project/record/output errors, 2 for invalid usage.
- No partial stdout on an invalid project. `show` writes exact source bytes;
  project and relative file paths appear on stderr. No automatic writes or repair.

## Files and interfaces

- `cmd/grove/main.go`: process entry point, delegates to command runner.
- `internal/cli/cli.go`: `Run(args []string, cwd string, out, errOut io.Writer) int`.
- `internal/project/project.go`: `Load(cwd, explicit string) (*Project, []Diagnostic)`,
  discovery, configuration, record-tree traversal, deterministic ordering.
- `internal/project/metadata.go`: YAML/frontmatter parsing and field validation;
  `Record` retains identity/type/title/status, optional planning fields,
  relationships/dates, relative path, and original source bytes.
- `internal/project/graph.go`: project identity index, references, separate cycles.
- `internal/project/*_test.go`, `internal/cli/cli_test.go`: fixtures in temporary
  directories; no test writes to the real record tree.
- README, AGENTS, model, brief, G-003: invocation instructions and verified state.

## Steps and verification

- [x] Reader: write discovery/config/frontmatter tests first; run
  `go test ./internal/project` and confirm missing reader failure. Implement
  configuration, strict metadata, traversal, and ordering; rerun to green.
- [x] Graph: write tests for duplicate IDs, wrong/missing targets, self-links,
  duplicate edges, independent cycles, and valid shared membership. Confirm
  failure before adding graph checks; rerun the full suite to green.
- [x] Commands: write tests for all three commands, exact-source display,
  usage/errors, project overrides, empty projects, ordering, diagnostics, and
  content-hash preservation for valid/invalid projects. Confirm failure before
  implementing `Run` and `main`; rerun the suite to green.
- [x] Review and verify: run `go test ./...`, `go test -race ./...`, `go vet ./...`,
  `gofmt` checks, and build a temporary binary. Run list/show/check against Grove's
  records and an actual linked-worktree fixture. Compare project file hashes
  before/after commands. Check doc links and reconcile G-003 and the brief.

Use writable temporary Go caches if the sandbox excludes the standard locations;
this is a verification-environment override, not a product requirement.

## Review focus

- A malformed neighbor must make `show` fail, rather than hide invalid project state.
- A nested Git checkout without config must not read its parent's project.
- Numeric-looking strings, nulls, aliases, duplicate keys and timestamp coercion
  must produce field errors instead of silently changing meaning.
- Renamed files, nested type folders, and missing optional metadata remain valid.
- Symlinked roots/files and misplaced Markdown cannot silently contribute data.

All five cases belong in the reader/command fixture tests above. This reader
observes live files without promising an atomic snapshot of concurrent edits.
Shared allocation, safe mutations, UI, and cross-branch resolution remain deferred.

## Progress and evidence

Prepared against the baseline above. No application tests existed to run before
implementation. Completed 2026-09-19; G-003 is reconciled and closed.

- `go build ./...`, `go vet ./...`, `gofmt -l .` (clean), `go test ./...`, and
  `go test -race ./...` pass for `internal/cli` and `internal/project`.
- `list`, `show G-003`, and `check` run against Grove's four records; `check`
  reports `OK: 4 records`. SHA-1 hashes of `grove.yaml` and every record file
  are identical before and after the three commands.
- Linked-worktree fixture (temporary repo, main plus a `feature` worktree):
  `list` from the worktree shows its uncommitted `status: done` edit while
  `--project ../main` shows main's `active`; a malformed neighbor
  (`W-002-bad.md` without a closing `---`) makes both `list` and `show W-001`
  exit 1 with a path diagnostic; discovery from a nested directory resolves the
  worktree root; the main checkout's files and `git worktree list` are unchanged.
  The fixture was removed.
- Not verified: concurrent writers, atomic snapshots, cross-branch aggregation.
