---
id: "G-144"
type: work
title: "Give adopting projects the record model the guides cite"
status: active
created: "2026-09-25T03:53:22Z"
updated: "2026-09-25T04:29:39Z"
relates_to: ["G-135", "G-143"]
---

## Outcome

A session following Grove's guides in any project that ran `grove init`
can read the record model the guides cite from that project or the `grove`
binary, without looking outside the project.

## Constraints

Observed at `worktree-G-135` `f37be4e`:

- [docs/work-shaping.md](../docs/work-shaping.md) (its read-in-stages
  table, and "Use only the record types, fields, and statuses the record
  model documents") and [docs/work-execution.md](../docs/work-execution.md)
  (its read-in-stages table) send a session to "the record model" when a
  field's meaning or allowed value matters or the CLI refuses a change.
- The record model is [docs/record-model.md](../docs/record-model.md) in
  this repository only. The binary embeds the two guides (`grove guide
  work|shape`) and the reviewer definition, not the record model, and
  `grove init` writes neither the model nor a pointer to it, so an adopting project has
  no copy.
- [G-143](G-143-g-135-codex-eval-row-pattern-pla.md) finding 2: in 5 of 9
  Codex shaping runs on the eval fixture, which is such a project, the
  session hunted for it, listing the owner's home directory, searching
  other repositories, and reading `docs/record-model.md` from the owner's
  own Grove checkout. Claude's runs did without it. Either way the guide's
  instruction cannot be followed as written outside this repository, and a
  session that does find a copy elsewhere may read a version that does not
  match the binary.

Selected design ([G-146](G-146-how-should-an-adopting-project-r.md),
resolved 2026-09-24): embed `docs/record-model.md` in the binary, print it
with `grove guide model`, have the three citations name that command, and
first edit the model to hold no `G-` link and no path link into this
repository, so it ships verbatim with no preamble. A `G-` ID is a live
identifier in every adopting project, so the model cannot carry Grove's
own. G-146 holds the alternatives, the reasons and what the edit covers.

## Acceptance

1. In a disposable project made by `grove init`, a session can read the
   record model the installed `grove` implements without leaving the
   project, and the guides say how.
2. The copy it reads cannot differ from what that binary validates, or the
   difference is detected; the guides digest reflects any guide change.
3. `go run ./cmd/grove check` and the repository's verification pass.

## Evidence

Implemented on `worktree-G-144`, base `main` `670ca9c` (branch point of the
shaping commits; G-146 resolved at `2d5777e`), from record revision
`sha256:1bc032d9b226` and plan [G-148](G-148-g-144-plan.md) as committed at
`3a57f52`. Implementation commits `7b841f9`, `3b62bff`, `fa18712`.

What changed:

- `docs/record-model.md` holds no `G-` link or repository path: provenance
  parentheticals are gone, "G-NNN owns X" sentences are plain statements
  (`update`'s failure reporting and the `flock` lock are now stated where
  the links used to point), the brief, command reference and shaping guide
  are named by command, and `main in this repository` and `go run` asides
  are generic. The only `G-` numbers left are the format examples `G-001`,
  `G-003`, `G-1000`. Headings are unchanged, so existing anchors resolve.
- `guides.go` embeds it; `grove guide model` prints it with no project; the
  `version` digest hashes work, shape and model; bad names exit 2 with
  "guide requires one argument, work, shape or model".
- The three citations name the command: `docs/work-shaping.md` lines 60 and
  146, `docs/work-execution.md` line 87.
- Reconciled: `docs/commands.md` (Version and guide), the usage text,
  README's command and record-model rows, AGENTS.md's pointer, and a new
  AGENTS.md constraint: a document the binary ships links only to other
  shipped documents, never to a record (G-146).
- `internal/cli/init_test.go`: `guide model` joins the prefix and
  portability checks, and a new assertion fails on any link other than `#`
  or `https://` and on any `G-` number but the examples. Mutation-checked:
  appending a G-064 link made it fail with both messages.

Against acceptance, at `fa18712`:

1. A binary built from `fa18712` ran `grove init` in a disposable Git
   project under `$TMPDIR` (`check`: `OK: 0 records`); `grove guide model`
   printed bytes `cmp`-identical to `docs/record-model.md`, and `guide
   work` and `guide shape` name `grove guide model`.
2. The printed model is the embedded file from the binary's own commit, so
   it cannot differ; the digest covers both guides and the model and
   changed with the citation edits: the same hash over `main` `670ca9c`'s
   two guides is `0c163c41a0f2`, and the built binary at `fa18712` prints
   `guides sha256:cd1491da3984`, which equals the hash of the three files.
3. `go vet ./...` clean, `gofmt -l .` empty, `go run ./cmd/grove check`
   `OK: 143 records`, `go test -count=1 -timeout 120s ./...` all ok.

Review: [G-149](G-149-g-144-review.md), two independent `grove-reviewer`
rounds, `examined` `3b62bff`; `fa18712` applies round 2's two wording
notes to the model only. No open blocking or should-fix finding.

Limits: the TUI terminal script was not run (no TUI change). The installed
`~/.local/bin/grove` is not rebuilt. `worktree-G-140` also edits
`docs/record-model.md` (a `run:` configuration key); whichever merges second
reconciles the model, keeping it free of `G-` links.

## Next

Captured 2026-09-24 from G-143; shaped through G-146, resolved 2026-09-24,
which selected the embedded model with no `G-` links.

In review 2026-09-25 on `worktree-G-144`: candidate is the evidence commit
this record names in `candidate`. For the owner:

- Judge the model edit: `git diff 2d5777e fa18712 -- docs/record-model.md`,
  and read it as an adopting project would with `go run ./cmd/grove guide
  model`.
- Open for the owner, not done here (G-149 finding 4): the guides still
  mention Grove's own records as plain text (G-035 and G-038 in
  `docs/work-execution.md`'s Lifecycle, G-032 in its Invocation, G-050 in
  `docs/work-shaping.md`). The new AGENTS.md rule forbids links only; if the
  G-146 reasoning should cover mentions too, capture that as new work.

Integration, as given:

```sh
go run ./cmd/grove approve G-144 "VERDICT"   # in this worktree
go run ./cmd/grove integrate G-144           # in the main checkout
```
