# Reused code

Files copied from `~/GitHub/mascah/grove-next` at
`8236a1e59a18f3735943b98d0827297de7940830` (CONSTITUTION.md rule 27), each
stripped to what the walkthrough uses.

| Copied to | From | What was kept |
| --- | --- | --- |
| `internal/gitx/gitx.go` | `internal/gitx/gitx.go` | `Run`: argv with `-C dir`, Git's standard error in the failure. Dropped the environment cleaning, replace/graft switches, context and wait delay (rule 13). |
| `internal/doc/frontmatter.go` | `internal/doc/frontmatter.go`, `internal/doc/document.go` | Frontmatter split, YAML error line, duplicate and unknown keys, required fields as diagnostics. Dropped kinds, typed schemas, timestamps, manifests, extensions, `Format` validation. |
| `internal/doc/doc.go` | `internal/doc/doc.go`, `internal/doc/ids.go`, `internal/doc/collection.go` | `Diagnostic`; `NewID` (YYMMDD-xxxxx, Crockford base32); the collection layout, as globs. Dropped severities, task parsing, implementation notes, review blockers, duplicate detection. |
| `internal/gitx/tree.go` | `internal/gitx/tree.go` | `TreeFS` over `ls-tree` and one `cat-file --batch`. Dropped path validation, reply checks and `HasFile`. |
| `internal/gitx/gitx.go` (`ResolveBranch`) | `internal/gitx/repo.go` | Branch tip. |
| `internal/instruct/guides/explore.md` | `internal/instruct/guides/explore.md` | The interview rules of section 1, section 3 and the finish, rewritten for a conversation that saves no document. |
| `internal/instruct/templates/spec.md` | `internal/templates/defaults/spec.md` | Whole file; the title placeholder became `{{title}}`. |
| `internal/instruct/templates/{proposal,change-spec,tasks}.md` | `internal/templates/defaults/` | Whole files; `text/template` fields became plain placeholders, the spec-ID conditionals became `<capability>`. |
| `internal/instruct/guides/propose.md`; sections 3 to 5 as `guides/{proposal,change-spec,tasks}.md` | `internal/instruct/guides/propose.md` | Sections 2, 6 and 7, with section 1 rewritten for proposing from the conversation, a description or a draft; `check` is Grove's `check`. Each document's section became that kind's guide in the schema. |
| `internal/tui/render.go` | `internal/render/render.go` | `Markdown`, `Sanitize` and the output filter that lets only the style's own SGR sequences through. Dropped `Plain`, truncation and the wrap fallback, which needed `charmbracelet/x/ansi` directly. |
