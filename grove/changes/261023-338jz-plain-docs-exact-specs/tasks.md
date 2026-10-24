## T-01 — Move the docs' exact rules into the maintained specs

Result: each exact rule `docs/` states is held by a maintained
requirement, and the page that states it explains it in plain words and
links the spec, as rule 8, amended in #78, says.

Files: `docs/index.md`, `docs/install.md`, `docs/getting-started.md`,
`docs/concepts.md`, `docs/coordinating.md`, `docs/planning.md`,
`docs/running.md`, `docs/reference.md`, `grove/specs/config/spec.md`

Interfaces: none; it is documentation and specs.

Constraints:
- Audit each page listed, as the dependencies left it. For each exact
  rule, find the maintained requirement that holds it, or the
  requirement this change ADDs for it.
- A rule held nowhere gets a requirement in this change before it leaves
  the page. The requirements in
  [specs/propose/spec.md](specs/propose/spec.md),
  [specs/changes/spec.md](specs/changes/spec.md) and
  [specs/init/spec.md](specs/init/spec.md) are the ones found while
  proposing. Any other goes in a revision of this change, written with
  full EARS and scenarios and agreed with the owner while implementing.
  Those not yet held at 0c6f419 include: `grove.yaml`'s fields other than
  `context`, and `grove check`'s report of an unknown one (`config`); the
  change schemas, a project's own files and schemas, and `grove check`'s
  schema problems; the folder layout and main-document checks; `grove
  review`'s labels; `grove specs`, `grove version`, the reader and IDs.
- A page keeps its steps, examples and quoted messages; the walkthrough
  runs them (rule 8). It keeps a restated rule only where a reader needs
  it to act, such as the selection rule in `docs/running.md`, "An
  ordinary run", step 2, which `grove/specs/runs/spec.md` holds;
  otherwise it gives one plain sentence and the link. The implementer
  proposes which restatements stay, and the owner agrees in the session.
- The `config` spec's Boundaries stops sending the reader to
  `docs/reference.md` for the rest of `grove.yaml`.
- This change edits no Go code. A scenario the code does not meet is
  raised with the owner before archiving.

Acceptance:
- For each page, a list in `implementation.md` names each rule the page
  states, the requirement, by spec and name, that holds it, whether the
  page kept, shortened or dropped it, and any difference fixed.
- Each rule a page states is held by a maintained requirement, and the
  page links it. Where the page and the spec differ, the page is fixed to
  match the spec.
- On the ledger fixture, each scenario of the requirements this change
  ADDs holds as written.
- `grove check` on Grove's own repository stays clean, and archiving
  this change folds the `changes` and `init` capabilities into
  `grove/specs/`.

## T-02 — The README and the home page lead with the attended path

Result: `README.md` and `docs/index.md` show one change taken from
propose, through review and merge, then implement, to merge, before
anything else; cycles, `/grove-plan` and `grove run` come after, under a
"Doing more" heading.

Files: `README.md`, `docs/index.md`

Interfaces: none.

Acceptance:
- Above "Doing more", neither file names a cycle, unattended work or
  `grove run`. README's second paragraph no longer brings in unattended
  runs.
- The steps in the order the owner decided: propose, review and merge,
  implement, merge.
- Every link in both files resolves, and `docs/index.md`'s "Where to go
  next" still names every page in `mkdocs.yml`.

## T-03 — A plain-language pass, page by page

Result: every page in `docs/` says one idea per sentence and breaks up
long conditional sentences, each term it gives a meaning is defined in
`docs/concepts.md`, "Glossary", which holds where a page differs, and
`dev/walkthrough.md` still runs every step the docs give.

Files: `README.md`, `docs/index.md`, `docs/install.md`,
`docs/getting-started.md`, `docs/concepts.md`, `docs/coordinating.md`,
`docs/planning.md`, `docs/running.md`, `docs/reference.md`,
`docs/guides/glossary-and-decisions.md`,
`docs/guides/push-to-a-tracker.md`,
`docs/guides/run-apply-on-a-timer.md`, `docs/guides/run-on-a-timer.md`,
`docs/guides/run-grove-in-ci.md`, `mkdocs.yml`, `dev/walkthrough.md`

Interfaces: none.

Constraints:
- Quoted command output stays exactly as Grove prints it.
- A page that the dependencies renamed or moved is edited at its new
  path.

Acceptance:
- `docs/planning.md` no longer has "A change that waits on one outside
  the cycle keeps `grove run` from the cycle until that one is archived
  as finished"; what it said is in shorter sentences.
- `docs/concepts.md`, "Two halves", no longer says Grove starts no agent
  and makes no commits; it names `grove run` as the one command that
  starts an agent.
- Each term the pages give a meaning, such as cycle, unattended, attempt,
  base, stale and dropped, is defined in the glossary. Another page may
  say what it means where a reader needs it, agreeing with the glossary,
  and otherwise uses it or links it.
- `dev/walkthrough.md` no longer calls the documentation the
  specification; it says, as amended rule 8 does, that it runs the
  documentation's steps. Each of its steps names the page section it
  tests as that section is called after the pass. Run by hand on a fresh fixture, as
  `CONTRIBUTING.md` says, its steps still pass, and `just check` passes.
- The docs site builds from `mkdocs.yml` with no broken link.
