# W-007 source-preserving update repair plan

> For Fable: implement one task at a time in an isolated worktree using the
> repository instructions and an inline execution workflow (for Codex,
> `superpowers:executing-plans`). Do not dispatch or merge automatically.

Review base: `9b7f730`; include the review/shaping documents in the execution
base. Read the linked work record first; it owns outcome, constraints, and
acceptance. This is a proposed technical repair plan, not implementation evidence.
Go 1.26+, current standard library and existing YAML dependency only.

**Goal:** satisfy the existing byte-preservation and observed-configuration
change contracts for all demonstrated valid requests.

**Architecture:** retain `update.Apply`, its shared lock and publication path.
Keep lexical edit planning in `edit.go` (extract a private span planner only if
needed); use `Project.Config` already supplied by the loader for exact comparisons.
No YAML reserialization or parser replacement.

**Spec:** [W-007](../../grove/work/W-007-preserve-updates.md).

## Global constraints and review focus

Preserve exact unrelated bytes, revision/no-op/date behavior, file permissions,
reservation consumption, lock order, and applied-state errors. Review: comments
between colon and value, explicit scalar keys, adjacent trailing flow removals
plus `updated` insertion, separator-adjacent comments, and valid configuration
formatting edits. Each is exercised below.

## Task 1: Preserve comments and support accepted entry forms

Files: `internal/update/edit.go`, `edit_test.go`, `update_test.go`.
Consumes accepted yaml.Node positions and `change` requests; produces the same
`Edit(source, changes) ([]byte, error)` output with exact byte preservation.

- [ ] Add exact expected-output fixtures for `title: # retain\n  First` updated
  to `New`. Expected entry is `title: # retain\n  "New"`, keeping all other
  bytes except insertion/change of `updated` at Apply level. Include a standalone
  comment between key and value and CRLF/BOM variants. This assertion currently
  fails at Apply level:

```go
root := gitProject(t)
source := strings.Replace(work, "title: First", "title: # retain\n  First", 1)
write(t, root, "grove/work/W-001-first.md", source)
apply(t, root, "W-001", []Field{{"title", "New"}})
if !strings.Contains(read(t, root, "grove/work/W-001-first.md"), "# retain") {
    t.Fatal("comment outside replaced value lost")
}
```

Add explicit `? status\n: proposed` set and explicit optional-key unset cases,
including quoted/tagged keys and multiline values already accepted by the reader.
- [ ] Run `go test -count=1 -run 'TestEdit|TestUpdate' ./internal/update` and
  capture the new fixture failures before changing the editor.
- [ ] Separate key/colon/value token ranges. For a later-line value, replace its
  token span at its existing indentation; remove continuation bytes belonging
  to that value only. Preserve comments before its start. For explicit keys,
  locate the separate colon indicator bounded by that entry's value and next
  entry; do not guess from the key text or strip the explicit `?` marker.
  Keep the overlap/bounds guards and parse/unchanged-field checks.
- [ ] Run exact-byte and Apply regressions. Include old tagged/anchored fixtures;
  no reader form is removed. Commit `fix(update): preserve comments around YAML values`.

## Task 2: Plan flow separators across all requested changes

Files: `internal/update/edit.go`, `edit_test.go`, `update_test.go`.
Consumes the full ordered `[]change`; returns disjoint span edits.

- [ ] Add the current failing fixture:

```go
root := gitProject(t)
write(t, root, "grove/work/W-001-first.md",
    "---\n{id: W-001, type: work, title: T, status: proposed, kind: fix, size: small}\n---\nBody\n")
apply(t, root, "W-001", nil, "kind", "size")
```

Expected output retains required entries and body, removes both optional fields,
and appends `updated` with one valid separator. Table-test reversed request order,
existing/absent `updated`, first/middle/last optional fields, trailing comma,
multiline flow mapping, set-plus-unset, and entry/standalone comments. Assert the
exact expected bytes, including removal of an unset entry's inline comment.
- [ ] Confirm `frontmatter: overlapping edits` on the base for the final-two case.
- [ ] Plan removed runs and retained neighbors together:

```text
Mark removed entries, retained entries, and appends for the whole request.
Assign each separator to at most one removal run.
Use retained neighbors/mapping boundaries to choose the surviving separators.
Preserve standalone comments; remove only removed-entry inline comments.
Plan appends against the resulting retained boundary, then apply disjoint spans
in descending byte order and validate with the existing parser/guard path.
```

Do not disable overlap detection or rewrite the complete flow mapping.
- [ ] Run editor/Apply tests and commit `fix(update): coordinate flow entry removals`.

## Task 3: Compare exact configuration before publication

Files: `internal/update/update.go`, `update_test.go`,
`internal/create/create.go`, `create_test.go`.
Consumes the immutable loaded `Project.Config`; exported APIs remain unchanged.

- [ ] Retain this deterministic failing regression in the update package:

```go
root := gitProject(t)
_, err := Apply(root, Request{ID: "W-001",
    Expect: revision(t, root, "grove/work/W-001-first.md"),
    Set: []Field{{"status", "active"}}}, now, func(step string) error {
    if step == "compare" {
        return os.WriteFile(filepath.Join(root, "grove.yaml"),
            []byte("# concurrent edit\nschema_version: 1\nrecords: grove\n"), 0644)
    }
    return nil
})
if err == nil { t.Fatal("published despite changed configuration") }
```

Also assert original target bytes and no temporary record residue. For creation,
load `p`, change its config bytes without changing parsed meaning, then call
`create.New(p, ...)`: the immutable allocation input differs from the under-lock
reload. Require refusal, no new record, and the consumed counter reservation.
This requires no new exported failure hook or timing-sensitive test.
- [ ] Run both base-failing cases. Add `!bytes.Equal(p.Config, snapshot.Config)`
  to update's `same` comparison and `!bytes.Equal(current.Config, p.Config)`
  to creation's under-lock comparison. Report configuration changed, preserving
  existing reserved-but-not-created diagnostics and lock order.
- [ ] Run all create/update tests, including concurrent and fault-injection
  fixtures. Commit `fix(records): refuse configuration changes before publication`.

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

## Implementation notes, 2026-09-19

Implemented on branch `worktree-W-006-W-008`; the
[work record](../../grove/work/W-007-preserve-updates.md) owns the evidence.
Bounded adjustments to this plan, and why:

- Tasks 1 and 2 share one commit (`2cc7814`). Both rewrite the same functions
  of `edit.go` (`edit`, the entry offsets, and the flow planner that replaces
  the flow half of `append`), so separate commits would not each build.
- A later-line value with no comment between its colon and itself keeps the
  W-003 behaviour and collapses onto the key's line; existing fixtures depend
  on it. Only when a comment is present is the value replaced in place. A
  block sequence at its key's indentation is the one value that may start
  there, so its replacement is indented two further columns to stay valid.
- Flow separators: an entry takes the comma that follows it; the run of
  removed entries ending the mapping takes the comma before it when only
  whitespace separates them. A comma hidden behind a comment is left as a
  trailing comma, which YAML accepts, rather than deleting the comment. An
  append placed after such a run reuses the surviving comma.
- A removed flow entry's same-line comment goes with it, including one that
  follows a trailing comma after the last entry; standalone comment lines stay.
- Refused safely, as a truly ambiguous span: a flow separator on a later line
  than its entry's value (`{a: 1\n, b: 2}`, removing `a`).
- Review found `keyStartsAt`'s plain-key branch always true (from W-003).
  `dedac73` makes the guard compare the source with the parsed key.
