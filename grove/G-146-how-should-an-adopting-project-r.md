---
id: "G-146"
type: question
title: "How should an adopting project read the record model the guides cite?"
status: open
created: "2026-09-25T04:10:06Z"
updated: "2026-09-25T04:10:26Z"
blocks: ["G-144"]
---

## Question

[G-144](G-144-give-adopting-projects-the-recor.md) needs the owner to
choose its design, which the record leaves open ("Choosing one is the
owner's"). No plan is written until this is resolved.

Observed at `worktree-G-144` `670ca9c`:

- The guides cite "the record model" three times:
  [work-shaping.md](../docs/work-shaping.md) lines 60 and 146,
  [work-execution.md](../docs/work-execution.md) line 87. The reviewer
  definition does not.
- `guides.go` embeds only the two guides and the reviewer definition;
  `grove guide` accepts `work` or `shape`; the `version` digest hashes the
  two guides. `grove init` writes no model and no pointer.
- [docs/record-model.md](../docs/record-model.md) is 542 lines and links
  35 of this repository's own `G-` records (G-064, G-065, G-052 and
  others) as the history of each rule. An adopting project has its own
  `G-` IDs, so `grove show G-064` there reads an unrelated record.

Options, recommendation first:

1. **Embed and print it.** `guides.go` embeds `docs/record-model.md`,
   `grove guide model` prints it, and the three citations name that
   command. The copy is the binary's own, so it cannot drift (acceptance
   2), and the guides digest changes with the citation edits. Recommended.
   Sub-choice for Grove's own `G-` links in the printed copy:
   - Recommended: print it verbatim under a one-line preamble saying its
     `G-` links name records in Grove's repository
     (github.com/mascah/grove), not this project's.
   - Or rewrite those links to repository URLs at print time.
   - Or move the history links out of the model so it holds none.
   Spelling of the command: `grove guide model` recommended, or
   `grove guide record-model`.
2. **`grove init` writes a copy** into the project (for example
   `docs/grove-record-model.md`). Readable as a file, but it goes stale
   when the binary changes, so acceptance 2 needs a staleness check, and
   it adds a file to every adopting project.
3. **Say it is unavailable.** The guides say the model exists only in
   Grove's repository and a session relies on `grove check` and CLI
   refusals. Least code, but acceptance 1 ("can read the record model")
   is then not met as written and G-144's outcome would change.

Who can answer: the owner, by setting this record `resolved` with the
option (and sub-choice and spelling, if option 1) in its Answer; any item
left blank takes the recommendation. Then launch
`/grove-work G-144 --until plan` on `worktree-G-144`.

## Next
