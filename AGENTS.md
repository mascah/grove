# Working on Grove

Read `docs/restart-brief.md` first. This repository is a fresh product restart,
currently containing documentation only. The brief distinguishes selected
direction, observed evidence, and proposed design; preserve those distinctions.

- Keep the product useful through ordinary local files and a CLI without a
  required running service. The core record model in `docs/record-model.md`
  is accepted: Markdown with YAML frontmatter. Use Go for the first CLI.
  Random IDs are accepted as a trial; their exact format, storage layout/versioning,
  UI framework, and runner contract remain open.
- Keep deterministic validation and state changes in software where useful;
  do not assume software can replace judgment instructions or prove acceptance.
- Record settled choices and the concrete next action in the brief while it
  remains small. Do not create a second editable account of the same direction.
- `../skills/` and `../nullsec/` are evidence and potential compatibility targets,
  not automatically part of an implementation's write scope. Follow their
  instructions; retrieve their Grove knowledge through `grove status`,
  `grove find`, and `grove context` from the relevant repository.
- The installed `grove` command currently belongs to the sibling skills project.
  No new Grove executable or configuration exists here yet.
- The archived application's service authority, architecture, credentials,
  deployment procedures, and backlog are historical. Do not revive them as
  requirements for this project or copy private local data into this repository.
- Use focused Conventional Commits. Preserve unrelated work and isolate
  concurrent implementation in separate worktrees.
- Verify claims against actual results. Documentation-only changes need link
  and consistency checks; do not invent an application test suite before code
  exists.
