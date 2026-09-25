---
id: "G-148"
type: plan
title: "G-144 plan: ship the record model as grove guide model"
status: current
created: "2026-09-25T04:29:20Z"
updated: "2026-09-25T04:29:23Z"
work: ["G-144"]
---

## Design

As [G-146](G-146-how-should-an-adopting-project-r.md) resolved it, at
`worktree-G-144` `2d5777e`:

- `guides.go` adds `docs/record-model.md` to the embedded `Guides`; `grove
  guide model` prints it verbatim, beside `work` and `shape`, needing no
  project. `version`'s digest hashes it after the two guides, so the digest
  names everything `guide` prints (acceptance 2 asks only for the guides;
  the model is one more list entry, and a schema change then shows in it).
- The model is edited first to hold no `G-` identifier except its ID-format
  examples (`G-001`, `G-1000`, `G-NNN`, the `G-003-inspect-records.md`
  folder example) and no link except in-document anchors and `https://`
  URLs: provenance parentheticals are deleted, "G-NNN owns X" sentences
  become plain statements or are dropped where the model already states the
  contract, the brief is named as `grove brief` prints it, the command
  reference and the shaping guide by name (`grove --help`, `grove guide
  shape`), and repository-specific asides (`main` in this repository,
  `go run ./cmd/grove`, "this repository's old IDs") become generic.
  Headings stay, so README's and records' anchors still resolve.
- The three citations (work-shaping.md's read-in-stages row and its "Use
  only the record types …" line, work-execution.md's read-in-stages row)
  name `grove guide model`.
- `init_test.go`'s portability check covers `guide model`, and a new
  assertion fails on any link in it other than `#` or `https://`, and on any
  `G-` number other than the format examples.
- Reconciled: `docs/commands.md` (Version and guide), the usage text,
  README's command row and record-model row, and AGENTS.md: its record-model
  pointer names `grove guide model`, and one line says a document the binary
  ships links only to other shipped documents, never to a record.

Not changed: the guides' own mentions of Grove's records (G-035, G-038,
G-032), which already say they are Grove's; `grove init` writes nothing new.

## Steps

1. Edit `docs/record-model.md` as above.
2. Embed it, add `guide model`, extend the digest, usage and tests.
3. Point the three citations at `grove guide model`; reconcile docs and
   AGENTS.md.
4. Verify: `go test -short ./internal/cli`, `go vet ./...`, `gofmt -l .`,
   `go run ./cmd/grove check`, `go test -count=1 -timeout 120s ./...`; in a
   disposable `grove init` project, run the built binary's `guide model`
   and confirm it equals the source bytes (acceptance 1 and 2).
5. Independent review through `grove-reviewer` on the final revision, then
   hand off into Review.

Status 2026-09-25: steps 1 to 5 done at `7b841f9`, `3b62bff` and
`fa18712`; review [G-149](G-149-g-144-review.md). Evidence is in G-144.
