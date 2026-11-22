# A change's fields lead its main document

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

A change keeps its fields in a hidden file beside its proposal, while a
cycle keeps them at the top of its own Markdown file.

- A change's title, created time, description, schema, `depends_on` and
  `base` are in `changes/<id>-<slug>/.grove.yaml`
  (`internal/instruct/metadata/change.grove.yaml` and `internal/doc/meta.go`
  at 15816d7). Grove finds a change by globbing for that file
  (`internal/doc/doc.go` at 15816d7, `changes/*/.grove.yaml`).
- Change 261118-emse7 made a cycle one Markdown file, its fields between
  two `---` lines at the top. Its decision 5 gives the reason: "the data
  stays YAML read by the same parser, the prose is Markdown that
  `grove serve` and a code host render, and a mistake in the prose cannot
  make the fields unreadable." `internal/doc/meta.go` at 15816d7 has `Front`
  to split such a file.
- The title is written twice, in `.grove.yaml` and as the proposal's first
  heading, in two files. `grove check` compares them.
- A dotfile is easy to miss in a listing. Tools that read Markdown with
  YAML at the top skip it entirely. For example, the Open Knowledge Format
  v0.3 (github.com/open-knowledge-format/open-knowledge-format, `SPEC.md`
  §3.2) leaves out files whose names start with `.`.

## Intended outcome

A change's fields sit between two `---` lines at the top of its main
document, `proposal.md` in both built-in schemas. The fields and their
meanings are unchanged, and `.grove.yaml` no longer exists.

- Grove reads the fields with the parser it uses for cycles. It finds a
  change by the one Markdown file directly in its folder that starts with
  `---`.
- `grove check` reports any `.grove.yaml` left in a change folder.
- The reader and `grove serve` show the fields as text above the rendered
  proposal.
- Grove's own changes, open and archived, are converted. Other projects
  convert theirs from a prompt in the docs.

## Scope

### In scope

- Where Grove reads a change's fields, how it finds a change, and the
  template propose writes.
- `metadata` in a change schema: `frontmatter` replaces
  `changes/<change>/.grove.yaml`.
- `grove check` reporting a leftover `.grove.yaml` in a change folder, or
  a folder with more than one main document.
- The reader, `grove serve`'s rendered document and its history of the
  record.
- The instructions that say not to edit `.grove.yaml` by hand.
- Converting every change under `grove/changes/`, this one included.
- The docs, the walkthrough, the trial scripts, and the push-to-a-tracker
  guide's issue body.

### Out of scope

- A new field, a renamed field, or a field with a new meaning. In
  particular, OKF's `type` is not added: no documented step needs it.
- Reading `.grove.yaml` as a fallback.
- Converting other projects' changes, such as tend's. Their owners convert
  them from the docs' prompt.
- Cycles and maintained specs. Cycles already work this way, and
  maintained specs have no fields. The report of a leftover `.grove.yaml`
  in a spec folder stays.

## Decisions

The owner accepted each recommended answer of the round with "yes".

1. **No fallback.** Grove reads no change from a folder that holds only
   `.grove.yaml`. `grove check` reports each such file and exits non-zero,
   as 261118-emse7 did for cycles in its decision 6. Until another project
   is converted, its changes do not show in Grove. Reason: a fallback is a
   second format to keep working forever, and converting is mechanical.
2. **How Grove finds a change.** A change's main document is the one
   Markdown file directly in its folder whose first line is `---`, and
   `grove check` reports a folder with more than one. In a schema,
   `metadata: frontmatter` on the main kind replaces
   `metadata: changes/<change>/.grove.yaml`. Reason: a project's own schema
   stays free to name its main document something other than
   `proposal.md`. Fixing the name would have been simpler, but it breaks
   such a schema.
3. **Convert all.** Every change under `grove/changes/`, open and archived,
   is converted in this change. Reason: Grove still reads archived changes
   for `depends_on`, history and state, and decision 1 gives no fallback.
4. **Keep the title heading.** The proposal's first heading stays
   `# <title>`, and `grove check` still compares it with the `title` field.
   Reason: the reader, `grove serve` and tracker issues render the body,
   and the title should show there.
5. **Display.** The reader and `grove serve` show the fields as text above
   the rendered body, as they do for cycles. `/api/project` gains no key,
   and a change's `files` no longer lists `.grove.yaml`.
6. **Hand edits.** The implement and propose instructions say never to
   edit the fields above the second `---` by hand, where they now say
   `.grove.yaml`. `base` still comes only from the template.
7. **Tracker guide.** The issue body is the proposal without its fields.
8. **No `depends_on`.** This change does not wait on 261122-m78w6 or on
   261122-z1kh0, although it shares files with both. Reason: they change
   different text. Whichever is built second builds on the files as the
   target then holds them.
9. **Routine choices.** Three tasks: Grove's reading and writing, the
   conversion of Grove's own changes, and the docs. There is one change
   spec for each capability whose requirements name `.grove.yaml`:
   changes, maintained-specs, propose, schemas and serve. Test layout and
   exact wording are delegated.

## Open questions

These can wait for review. `/grove-implement` builds the recommended
answer otherwise.

1. **Where the fields' template comes from.** Recommended: Grove prints
   the main kind's template with its own fields block, the lines of
   `internal/instruct/metadata/change.grove.yaml` between two `---` lines,
   put before it. The schema's template stays only the body. Reason: the
   fields, `base` in particular, stay Grove's and are filled by Grove, so
   a project that overrides the proposal template cannot drop them.
   Decides T-01.
2. **The conversion commit in the history.** Recommended: accept that
   `grove serve`'s history marks the commit converting Grove's changes as
   revising each open change it touches, with no rule to hide it. Reason:
   it happens once, and it did revise those files. Decides T-02.

## Related documents

- Requirements:
  [specs/changes/spec.md](specs/changes/spec.md),
  [specs/maintained-specs/spec.md](specs/maintained-specs/spec.md),
  [specs/propose/spec.md](specs/propose/spec.md),
  [specs/schemas/spec.md](specs/schemas/spec.md),
  [specs/serve/spec.md](specs/serve/spec.md)
- Tasks: [tasks.md](tasks.md)
