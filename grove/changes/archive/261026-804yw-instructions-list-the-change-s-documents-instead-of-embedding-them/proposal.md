# Instructions list the change's documents instead of embedding them

## Problem

The review of Grove 0.5.0 on memex (`dev/memex-review-2026-10-26.md`,
finding 1) observed:

- Claude Code turns a shell result over 30,000 characters into a 2 KB
  preview and a saved file. All five `grove instructions implement` runs
  and two of five `grove instructions archive` runs tripped it, each
  costing one or two extra turns; by the third session agents redirected
  the output to a file before reading it.
- Codex 0.164.0, measured 2026-10-26, keeps the start and end of a shell
  result over about 10,000 tokens (about 40 KB of Grove's text), drops the
  middle with "Warning: truncated output", and saves nothing. memex's
  implement jnqjq (49 KB), syssv (47 KB) and archive jnqjq (54 KB) would
  have lost their middle, which is the change itself, with no way to get
  it back.
- Of those sizes, the change's own documents were 11 to 45 KB, and about
  9 KB of implement's 17 KB fixed part was propose's templates
  (`.grove.yaml`, proposal, change spec, tasks), which implement needs only
  when revising.

At ea7fbaf, `build` in `internal/instruct/instruct.go` puts every file of
the change into `Instructions.Change`, and `Markdown` prints each in full
under "The change"; the default schema's implement `writes` lists
proposal, change-spec, tasks and implementation
(`internal/instruct/schemas/default/schema.yaml`). Rule 12's word budget
bounds only the fixed text.

## Intended outcome

`grove instructions implement`, `grove instructions archive` and
`grove instructions propose <change-id>` name each document of the change
by path, for the agent to read with its own file tool, instead of printing
its text. Where the working tree's copy of an approved change's document
differs from the target's, implement says so and gives the `git show`
command for the approved text. Implement prints only the implementation
template, and points to `grove instructions propose <change-id>` for any
revision. The printed output no longer grows with the change's documents.

## Scope

### In scope

- Implement (with the owner and unattended, for a change or a member of a
  cycle), archive and propose's revise path.
- Implement's documents: the implementation template only.
- The implement instructions' revision text.
- The docs and walkthrough steps that describe the embedded change.
- A trials count of listed documents an agent did not read.

### Out of scope

- The cycle file, which stays printed in full (round, below).
- The project index, which already lists documents by path.
- Harness settings such as Claude Code's output limit: they fix one harness
  only.

## Decisions

- "The CLI says where to look instead of embedding documents. Implement
  and archive list each change document by path for the agent to read with
  its own file tool, and implement stops printing propose's templates (it
  points to `grove instructions propose <id>` for a revision). Output then
  stays the same size however big the change is." (owner,
  `dev/memex-review-2026-10-26.md`, decision 1)
- How "build what the target holds" carries over, settled in the round
  ("yes"): implement lists each document by its working-tree path and says
  to read every one before building. For an approved change, each document
  whose working-tree copy differs from the target's is named, with
  `git show <commit>:<path>` as the approved text. With the owner, a
  difference that is a revision the owner agreed on this branch is built
  as revised; unattended builds the target's text.
- Propose's revise path stops embedding too: all three go through `build`,
  and revising a 31 KB change would overflow the same way (round, "yes").
- Implement's `writes` in the default schema becomes `[implementation]`;
  the implement instructions point to `grove instructions propose
  <change-id>` before any revision; archive keeps `[spec, implementation]`
  (round, "yes").
- The cycle file stays embedded: it is a few hundred bytes, and the
  decision named change documents (round, "yes").
- "The trials count whether agents read every listed document." (owner,
  decision 1) As a judged count, "Listed documents not read" (round,
  "yes").
- This change depends on 261026-ar7ec, which edits the same implement
  instructions (owner, in this conversation: "yes").
- The instruction text stays at or under 8,000 words (rule 12) by
  tightening the text this change touches (round, "yes").

## Open questions

None.

## Related documents

- Requirements: [implement](specs/implement/spec.md),
  [propose](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
