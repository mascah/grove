---
id: "G-151"
type: work
title: "Strip Grove-repository pointers from the shipped guides and model"
status: active
created: "2026-09-25T19:06:01Z"
updated: "2026-09-25T19:22:18Z"
kind: fix
priority: 1
size: small
relates_to: ["G-110", "G-144", "G-146", "G-149", "G-152"]
---

## Outcome

The three documents the binary ships, the work guide, the shaping guide and
the record model (`grove guide work|shape|model`), point a session in an
adopting project at nothing that project lacks: no Grove record by number,
no Grove document the binary does not carry, no predecessor; and a test keeps
it so. Owner intent, review conversation 2026-09-25: Grove is used in other
projects as an installed CLI without this repository, and pointers into it
are part of the early self-dogfooding to untangle.

## Constraints

Observed at `main` `001b271`, after that commit removed the last path links
and [G-144](G-144-give-adopting-projects-the-recor.md) removed the model's
`G-` links:

- [work-execution.md](../docs/work-execution.md) line 26: "as Grove's own
  records G-035 and G-038 selected and implemented"; line 465: which
  invocation rows were exercised, and what the workflow kept from the
  predecessor's `/work`, "are recorded in Grove's own repository (G-032)".
  [work-shaping.md](../docs/work-shaping.md) line 264: the same for G-050.
  [record-model.md](../docs/record-model.md) lines 8 and 408 name "Grove's
  command reference", which is `docs/commands.md` here and is not shipped;
  line 244 names "the predecessor skills CLI's schema numbering", a tool an
  adopting project has never seen ([G-149](G-149-g-144-review.md) finding
  3). Every other `G-` number in the three documents is an example: G-030
  and G-031 in the work guide, G-037 in the shaping guide, G-001, G-003 and
  G-1000 in the model.
- [G-146](G-146-how-should-an-adopting-project-r.md)'s answer gives the
  reason: a `G-` ID is a live identifier in every adopting project, so a
  shipped sentence saying "G-035 selected it" is a false statement in that
  project's vocabulary, read by a session that takes IDs literally. G-149
  finding 4 noted that the reason covers unlinked mentions too; G-144 left
  widening the rule to the owner, who selects it through this record.
  [G-143](G-143-g-135-codex-eval-row-pattern-pla.md) finding 2 shows what a
  dangling pointer costs: five of nine Codex sessions searched the owner's
  home for "the record model".
- `TestGuideAndVersionNeedNoProject` (`internal/cli/init_test.go`) now
  checks that no shipped document links outside itself or `https://` and
  that the model names only its example IDs; the guides' `G-` numbers are
  unchecked, and G-149 finding 7 noted that an unlinked mention of an
  allowed number would pass.
- The README's owner table already points this repository's readers at
  G-032 and G-050 for what was exercised, so the guides lose nothing by
  dropping those sentences.

Proposed design, labelled proposed:

1. Rewrite or delete the six passages: the Lifecycle sentence keeps the
   lifecycle and drops the provenance; the two history sentences go; the
   model's two command-reference clauses go, leaving `grove --help` as the
   shipped account of usage; the predecessor clause goes. Nothing shipped
   links to `github.com` yet: the preview's public documentation
   ([G-110](G-110-external-preview.md)) is where a URL to the command
   reference becomes stable, and adding one then is one edit.
2. The test names each shipped document's allowed example IDs and fails on
   any other `G-` number, and the portability denylist gains "Grove's own
   repository" and "command reference".
3. AGENTS.md's rule says "names no record" as well as links none. The
   reviewer definition, also shipped, names none today and joins the same
   check.

Alternatives considered: shipping `docs/commands.md` as `grove guide
commands` after removing its own record links (consistent, but a fourth
embedded document for two clauses; open to G-110 if preview users need it),
and linking `docs/commands.md` on GitHub by `https://` (allowed by the rule,
but pinned to `main` and so not the binary's version).

## Acceptance

1. `grove guide work`, `guide shape` and `guide model` from a binary built
   at the candidate contain no `G-` number except the listed examples, and
   no mention of Grove's repository, its command reference or the
   predecessor, and read as complete sentences where text was removed.
2. `TestGuideAndVersionNeedNoProject` fails when a `G-` mention or either
   phrase is appended to any of the three documents, mutation-checked and
   recorded in Evidence, and passes at the candidate.
3. AGENTS.md's shipped-document rule covers mentions; README and
   `docs/commands.md` still say where this repository's own history lives;
   Evidence records the new guides digest from `grove version`.
4. `go vet`, `gofmt`, `grove check` and the full uncached suite pass.

## Evidence

Branch `worktree-G-151`, base `main` `28aaf95`, code reviewed at `02a0e08`;
the candidate is the commit recording this Evidence, which adds only this
record and G-156 to it. Started
from this record at `sha256:47d9e7e0aa18` with no plan (the record said none
was needed). Commits: `114db6b` (text, test, rule), `be4efcb` and `02a0e08`
(review fixes).

1. Six passages rewritten or removed as proposed: the Lifecycle sentence
   keeps the lifecycle without G-035/G-038; the work guide's G-032 and the
   shaping guide's G-050 history sentences are gone, each paragraph ending
   on its previous sentence; the model's two command-reference clauses are
   gone (`grove --help` is the shipped account of usage) and the
   predecessor clause reads "unrelated to any other tool's schema numbering
   or configuration". Beyond the six, the model's intro "How each rule was
   chosen is in Grove's own records and Git history" became "It states each
   rule, not how the rule was chosen", the same pointer in other words. A
   binary built at `02a0e08` prints only G-030 and G-031 (`guide work`),
   G-037 (`guide shape`) and G-001, G-003 and G-1000 (`guide model`), and
   none of "Grove's own repository/records", "command reference" or
   "predecessor", case-insensitively.
2. `TestGuideAndVersionNeedNoProject` checks the three guides and the
   reviewer definition: links only to `#` or `https://`; `G-` numbers only
   from each document's own example list (none for the reviewer); and,
   after lower-casing, collapsing whitespace and mapping U+2019, none of
   "grove's own repository", "grove's repository", "grove's own records",
   "command reference", "predecessor". Mutation at `114db6b`: appending
   "See G-035.", "See G-001.", and each phrase to each of the four documents
   failed the test in every case except "See G-001." on the model, an
   allowed example (G-149 finding 7's ceiling, which now holds per
   document). At `be4efcb` and `02a0e08`: "Command reference", "The
   Predecessor tool", "Grove's repository", "See Grove's\nrepository
   here.", "Grove's command\nreference." and "See Grove’s repository." each
   fail. Passes at the candidate.
3. AGENTS.md's rule: a shipped document "never links a record or names one
   beyond its own example IDs, and never names Grove's repository or the
   predecessor" (G-146, G-151). README lines 77-78 still point at G-032 and
   G-050, and `docs/commands.md` is unchanged. `grove version` guides digest
   at the candidate: `sha256:67dab310de86`. The `guides.go` comment and the
   [G-152](G-152-shipped-document.md) term's Relationships were reconciled.
4. At `02a0e08`: `go vet ./...` clean, `gofmt -l .` empty, `grove check` OK
   (150 records), `go test -count=1 -timeout 120s ./...` all ten packages
   ok.

Decision taken: the model no longer names the command reference, which
reverses the part of [G-146](G-146-how-should-an-adopting-project-r.md)'s
Answer that said to cite it by name; this record's design is the owner's
selection, and G-146 stays as answered.

Review: [G-156](G-156-g-151-review-shipped-guides-and.md), three rounds,
examined `02a0e08`, nothing open beyond two accepted notes.

Limits: a phrase with Markdown between its words, or an allowed example ID
used as a provenance mention, still passes. A binary built inside this
nested worktree stamps the enclosing `main` checkout's revision in `grove
version` (Go finds the parent `.git` directory, not the worktree's `.git`
file); the guides digest is this checkout's. Not in scope.

## Next

In review; the candidate differs from the reviewed `02a0e08` only in
records. When [G-110](G-110-external-preview.md) adds an
`https://` link to the public command reference, the test's "command
reference" entry needs changing with it. Integrator:

```sh
grove approve G-151 "VERDICT"   # in worktree-G-151
grove integrate G-151           # in the main checkout
```
