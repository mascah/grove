# Reused code

Files copied from `~/GitHub/mascah/grove-next` at
`8236a1e59a18f3735943b98d0827297de7940830` (CONSTITUTION.md rule 27), each
stripped to what the walkthrough uses.

| Copied to | From | What was kept |
| --- | --- | --- |
| `internal/instruct/templates/brief.md` | `internal/templates/defaults/brief.md` | Whole file; the title placeholder became "Product brief". |
| `internal/gitx/gitx.go` | `internal/gitx/gitx.go` | `Run`: argv with `-C dir`, Git's standard error in the failure. Dropped the environment cleaning, replace/graft switches, context and wait delay (rule 13). |
| `internal/doc/frontmatter.go` | `internal/doc/frontmatter.go`, `internal/doc/document.go` | Frontmatter split, YAML error line, duplicate and unknown keys, required fields as diagnostics. Dropped kinds, typed schemas, timestamps, manifests, extensions, `Format` validation. |
| `internal/doc/doc.go` | `internal/doc/doc.go`, `internal/doc/ids.go`, `internal/doc/collection.go` | `Diagnostic`; `NewID` (YYMMDD-xxxxx, Crockford base32); the collection layout, as globs. Dropped severities, task parsing, implementation notes, review blockers, duplicate detection. |
| `internal/gitx/worktree.go` | `internal/gitx/worktree.go` | NUL porcelain worktree list; `AddWorktree`. Dropped bare/detached/prunable fields and the hook override. |
| `internal/gitx/tree.go` | `internal/gitx/tree.go` | `TreeFS` over `ls-tree` and one `cat-file --batch`. Dropped path validation, reply checks and `HasFile`. |
| `internal/gitx/gitx.go` (`ResolveBranch`, `Uncommitted`) | `internal/gitx/repo.go`, `internal/gitx/status.go` | Branch tip; NUL status without renames, counted. |
| `internal/instruct/templates/exploration.md` | `internal/templates/defaults/exploration.md` | Whole file; `text/template` fields became plain placeholders. |
| `internal/instruct/guides/explore.md` | `internal/instruct/guides/explore.md` | Sections 1 to 3 and the finish; `check` removed (not in the walkthrough). |
