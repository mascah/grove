---
id: "G-144"
type: work
title: "Give adopting projects the record model the guides cite"
status: proposed
created: "2026-09-25T03:53:22Z"
updated: "2026-09-25T03:53:35Z"
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

## Next

Captured 2026-09-24 from G-143; not assigned. Shape the choice of design
with the owner before assigning.

Checkpoint 2026-09-24, `/grove-work G-144 --until plan --interaction
headless` on `worktree-G-144`, base `main` `670ca9c`, record revision
`sha256:8276cf78a97f`: waiting on question
[G-146](G-146-how-should-an-adopting-project-r.md), which blocks this
record and holds the options, evidence and recommendation. No plan is
written, since the design is the owner's choice; status stays `proposed`.
No commands are owned. Once G-146 is resolved, rerun
`/grove-work G-144 --until plan` on this branch to write the plan against
the answer, or drop `--until plan` to go on to implement it.

Resolved 2026-09-24: G-146 selected option 1 with the model holding no
`G-` links. Next: `/grove-work G-144 --until plan` on this branch, or drop
`--until plan` to implement.
