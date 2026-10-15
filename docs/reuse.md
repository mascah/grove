# Reused code

Files copied from `~/GitHub/mascah/grove-next` at
`8236a1e59a18f3735943b98d0827297de7940830` (CONSTITUTION.md rule 27), each
stripped to what the walkthrough uses.

| Copied to | From | What was kept |
| --- | --- | --- |
| `internal/instruct/templates/brief.md` | `internal/templates/defaults/brief.md` | Whole file; the title placeholder became "Product brief". |
| `internal/gitx/gitx.go` | `internal/gitx/gitx.go` | `Run`: argv with `-C dir`, Git's standard error in the failure. Dropped the environment cleaning, replace/graft switches, context and wait delay (rule 13). |
