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

Status: done

- Every page in `docs/`, guides included: long conditional sentences
  split, one idea per sentence; lists of conditions are bullet lists
  (such as `docs/running.md`, "An ordinary run", step 2, and the
  `grove check` reports in `docs/planning.md`). Headings, quoted output,
  code blocks and the spec links are unchanged; no fact was dropped.
- `docs/planning.md` no longer has "A change that waits on one outside
  the cycle keeps `grove run` from the cycle…"; it says it in shorter
  sentences.
- `docs/concepts.md`: "Two halves" names `grove run` as the one command
  that starts an agent (done in T-01). The glossary now defines approved,
  archived as finished, attempt, base, blocked, build order, cycle and
  done, draft, dropped, needs a person, paused, planned, stale,
  unattended and waits on, and says it holds where a page words a term
  differently. "What Grove requires at each step" is a list.
- Wording that differed from the code, fixed on the way: the pause lasts
  until the cycle's changes are archived on the target, as finished or
  dropped, not until "merged or dropped" (`docs/running.md`); a tool
  result that is not an error gets no row (`docs/guides/run-on-a-timer.md`);
  `grove review` labels files in `cycles/` `cycle`
  (`docs/guides/glossary-and-decisions.md`).
- `dev/walkthrough.md` no longer calls the documentation the
  specification; it runs the documentation's steps, as rule 8 says. No
  section heading changed, so each step still names its section as it is
  called.
- `README.md` and `mkdocs.yml` needed no further edit.
- Checks: `just check` passes (walkthrough_test included); `grove check`
  clean; `uvx zensical==0.0.71 build` (the docs workflow's build) on a
  copy of `docs/` and `mkdocs.yml`: "No issues found"; every relative link
  and every GitHub spec link in `README.md` and `docs/` resolves (the
  `changes` and `init` specs once archived). Not run in this session: the
  walkthrough by hand in Claude Code and Codex; the docs' steps and quoted
  output it runs are unchanged.

## Scenarios

Each run on a fresh ledger fixture (`just fixture`) with this branch's
`bin/grove`, 2026-10-24:

- changes "A change that waits": `grove changes` printed `261024-aaaaa  Change A  approved` and `261024-bbbbb  Change B  approved  waits on 261024-aaaaa`, passed
- changes "Its dependency archived": after archiving A, only B listed, with no `waits on`, passed
- changes "A draft": an uncommitted change D had no line, passed
- changes "One approved, one archived": `--json` gave B `approved` with no `waits_on`, A `archived` at `changes/archive/261024-aaaaa-a/proposal.md` with tasks `done`, passed
- changes "A dropped change" (`--json`): C archived with `Status: dropped` had `state` `dropped`, passed
- changes "No implementation yet": tasks with no `status`, no `description` key, passed
- changes "A task not done": `grove archive` refused, "T-02 is not done in implementation.md; …", folder not moved, passed
- changes "Every task done": "Archived 261024-aaaaa to grove/changes/archive/261024-aaaaa-a/; nothing is committed.", move uncommitted in `git status`, passed
- changes "A dropped change" (archive): C moved to `changes/archive/`, passed
- init "A repository with AGENTS.md": `agents` held `codex: skills` only, "Wrote grove.yaml and the Grove skills for codex; nothing is committed.", passed
- init "No marks": `claude: skills` and `codex: skills`, passed
- init "grove.yaml already there": "Kept grove.yaml and wrote the Grove skills for cursor; …", file unchanged, only `.agents/skills/` written, passed
- init "Claude Code as commands": six files in `.claude/commands/grove/`, each running `grove instructions`, passed
- init "Two harnesses share a folder": one `.agents/skills/grove-<skill>/` per skill, passed
- init "A form a harness lacks": `harness "codex" has no form "commands"; it has skills`, no skill file written, passed
- init "Claude Code switched to commands": `.claude/skills/grove-*` removed and named, `.claude/skills/mine/` kept, passed
- init "A skill Grove no longer has": `.agents/skills/grove-old` removed and named, passed
- propose "No such change": `.grove.yaml:4: depends_on: 260101-zzzzz is no change, open or archived`, exit 1, passed
- propose "A dropped change": "depends_on: 261024-ccccc was dropped, so this change could never be applied", passed
- propose "A loop": each change's `depends_on` reported, "waits on this change in turn, so neither can be applied", passed
- propose "A dependency archived as finished": `No problems in grove/.`, passed
- propose "Rewritten by another change": `spec.md:5: requirement "Keep a category" was rewritten since base 972422f; revise the change with propose`, passed
- propose "Removed by another change": "… was removed since base 972422f …", passed
- propose "Rewritten by a change it depends on": `No problems in grove/.`, passed
- propose "Already reconciled": `No problems in grove/.`, passed

## Reconciliation
