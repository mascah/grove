# Tasks name their files

## Problem

Two observations from the owner's use of Grove, 2026-10-22 and 2026-10-23.

- **Building touched files no task named.** Implementing change
  261022-1grb0 (session `bbeed28a`), the agent edited `docs/concepts.md`
  and `docs/getting-started.md`, which no task named, and said so
  afterwards. The implement and apply guides
  (`internal/instruct/schemas/default/skills/{implement,apply}.md` at
  b143991) say files no task names stay as they are, but a task names
  files only when it is about docs: the propose guide
  (`internal/instruct/schemas/default/skills/propose.md`) says "No
  implementation planning: no file-by-file plan", a line commit 6e2b727
  added without a recorded reason, and the tasks guide says a task "refers
  to requirements instead of repeating them".
- **Where the fix goes was left to the implementer.** Exploring on the
  ledger fixture (session `d50a0a86`), the agent closed with "Where the fix
  goes in the code ... these are details for whoever implements it, not
  questions for you". The proposal carried it as an assumption, and the
  implement session built on it unasked.

Grove detects a change spec that went stale against the maintained spec
since the change's `base` (`internal/cli/check.go`, `stale`, at b143991)
and a cycle whose member folders moved since `planned` (`changedSincePlan`).
Nothing detects code moving under a change, because no document names the
code a change touches.

## Intended outcome

- Each task in `tasks.md` names the files it creates, modifies and tests
  on a `Files:` line, and the interfaces it produces on an `Interfaces:`
  line, without code or steps.
- `grove check` reports a task whose named file changed on the target
  since the change's `base`, for a change in no open cycle, and reports a
  cycle as stale when a member's named file changed since `planned`.
  `grove run` skips such a cycle as it skips any stale cycle, and
  `/grove-implement` tells the owner and asks before building the task.
- Implement and apply touch only the files the tasks name. A file found to
  need a change while building is added to the task with the owner's
  agreement under implement, and left untouched and named in the task's
  evidence under apply.

## Scope

### In scope

- The `Files:` and `Interfaces:` lines in a task; the tasks guide and
  template, the propose, implement and apply guides, and the bugfix
  example guide in `docs/reference.md`.
- Grove reading the `Files:` line; the `check` report against `base`; the
  cycle stale rule widened to named files; `instructions implement` telling
  the owner.
- The documentation and walkthrough steps for each.

### Out of scope

- Code snippets, commands or bite-sized steps in tasks: the superpowers
  plan layer that is written minutes before execution. A task says what and
  where, not how.
- `grove review` reporting files a branch changed that no task names.
- A new gate beyond the cycle stale rule (`CONSTITUTION.md` rule 10, as
  amended by pull request #28).
- Propose asking a round before writing: change 261023-fca8x.

## Decisions

- **File-level thinking happens at propose, inside tasks.** The owner,
  2026-10-23, offered tasks naming files, module-level decisions in
  `design.md`, a plan written at implement time, or leaving it: "1." and "I
  lean this way, because I'm starting to imagine whether we should track a
  hash or SHA of the files a change expected to change so grove can easily
  check". Reason: the owner wants "to avoid wasting development time
  building the wrong thing", and a named file is what a check can watch.
- **Files and interfaces, no code.** The owner, 2026-10-23, offered files
  only, files and interfaces, or the full superpowers plan with steps and
  code, and took the recommendation: "yes". Reason: files give the check
  its paths and interfaces settle where a fix goes, while code steps are
  the layer that goes stale soonest under approval by merge.
- **Two baselines, by how the change is built.** The owner, 2026-10-23:
  "1". A change built alone is checked against its `base`, and implement
  asks; a member of a cycle is checked against the cycle's `planned`,
  re-planning being the one softening, so members built in one run never
  make each other stale. Reason: the owner's "we would likely need to
  allow items planned to change together in a cycle to soften up the
  rules"; `planned` already is that shared baseline.
- **Git is the hash.** The owner's "track a hash or SHA of the files" is
  met by `git diff --name-only` between `base` or `planned` and the target,
  as `changedSincePlan` does today; no new field is recorded. The owner
  accepted this in the exploration's summary.
- **The round of 2026-10-23, accepted with "yes":** a `Files:` line after
  `Result:` listing repo-relative paths in backticks, comma-separated, the
  only line Grove parses, and an `Interfaces:` line of free text; an
  unnamed file found while building is added to the task with the owner's
  agreement under implement, and left untouched and named in the task's
  evidence under apply; `check` compares against `base` only a change in no
  cycle that is not done; the bugfix example guide drops its "no
  file-by-file plan" line and names its files; no dependency between this
  change and 261023-fca8x, both adding to a new `propose` spec.
- **The constitution allows it** since pull request #28 (rules 3 and 10).

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md) (no
  maintained spec at b143991), [specs/cycles/spec.md](specs/cycles/spec.md)
  (maintained spec: `grove/specs/cycles/spec.md`),
  [specs/implement/spec.md](specs/implement/spec.md) (maintained spec:
  `grove/specs/implement/spec.md`)
- Tasks: [tasks.md](tasks.md)
- Constitution amendment: pull request #28
