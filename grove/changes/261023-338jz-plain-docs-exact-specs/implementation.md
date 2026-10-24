# Implementation: Plain docs, exact specs

<!-- Progress on change 261023-338jz, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Pages link the specs and agree with them

Status: done

Built to the revision of 2026-10-24 (proposal, "Docs teach; they need not
map onto the specs"). Teaching prose is kept by default. Links are full
GitHub links to `grove/specs/<id>/spec.md` on `main`; the `changes` and
`init` links resolve once this change is archived and merged.

Per page:

- `docs/index.md`: no link added; it points to the pages that link.
  Nothing differed.
- `docs/install.md`: no spec (the owner agreed `install.sh` stays
  unspecified). Nothing differed.
- `docs/getting-started.md`: links to the init, config, prototype,
  propose, implement and requirements specs. Nothing shortened. Nothing
  differed.
- `docs/concepts.md`: links to the maintained-specs and schemas specs.
  Fixed "Two halves", which said Grove starts no agent; it now names
  `grove run` as the one command that starts one.
- `docs/coordinating.md`: links to the propose spec (unmet `depends_on`,
  stale requirement), the implement spec (building a cycle) and the
  changes spec (the JSON). Nothing shortened. Nothing differed.
- `docs/planning.md`: link to the cycles spec. Nothing differed.
- `docs/running.md`: link to the runs spec. Fixed step 2 of "An ordinary
  run", which said a cycle is selected when "its changes are as its plan
  read them", leaving out the files its tasks name; it now says its plan
  is not stale, linking where staleness is explained. The 5-second grace
  and the attempt record's keys stay: the code does them.
- `docs/reference.md`: shortened the rows of `grove check`, `grove
  archive`, `grove run`, `grove run <cycle-id>` and the `grove
  instructions` rows for prototype, propose, plan, implement and archive,
  which copied spec detail, to plain sentences with links to the
  requirements, propose, cycles, implement, prototype, changes, init and
  runs specs. Links added under "Harnesses" (init), "grove.yaml" (config),
  "Change schemas" (schemas) and "Files" (maintained-specs). Fixed the
  template values: `{{base}}` is the target's commit shortened to seven
  characters, `{{commit}}`, which the page did not name, is the full
  commit, and `{{title}}` is filled once a change has one
  (`internal/instruct/instruct.go`, `fill`). The `grove run` row no longer
  restates the fetch, which left out the fast-forward of the local
  target; it links the runs spec.
- `grove/specs/config/spec.md`: Boundaries now says the rest of
  `grove.yaml` has no maintained spec yet and `docs/reference.md`
  explains it.

Checks: `grove check` clean; no test reads `docs/`.

### T-02 — The README and the home page lead with the attended path

Status: done

- `README.md`: the four steps (propose, review and merge, implement,
  merge), then one sentence that a small change takes one merge, then
  Install and Getting started. Explore, prototype, cycles, `/grove-plan`
  and `grove run` moved under "Doing more".
- `docs/index.md`: "Two halves" no longer says the command runs a cycle
  or that skills implement without you; "The short version" is now "One
  change", without the plan and run lines, followed by the one-merge
  sentence; cycles, `--unattended` and `grove run` are under "Doing more",
  before "Where to go next".
- Checks: above "Doing more", neither file matches `cycle`, `unattended`
  or `grove run` (grep). Every relative link in both files names a file
  that exists, and "Where to go next" still links each of the 11 pages
  `mkdocs.yml` lists besides the home page.

### T-03 — A plain-language pass, page by page

Status: not started

## Scenarios

## Reconciliation
