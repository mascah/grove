# A word budget for grove.yaml's context and rules

## Problem

Rule 12 of `CONSTITUTION.md` limits Grove's own instruction text to 9,000
words, counted by `just check`, because every word is read on every run.
A project's own text has no such limit:

- `context` in `grove.yaml` is printed under "Project context" in the
  instructions of every skill (`grove/specs/config/spec.md`, "Give context
  to every skill"; `internal/instruct/instruct.go:646` at dc8d335).
- `rules` lines are printed with every instruction for their skill
  (`rules.skills`, `rules.delivery`) or document kind (`rules.documents`)
  (`internal/instruct/instruct.go:650`, `:652` and `:669` at dc8d335).
- Nothing measures either. `docs/reference.md` at dc8d335 says only
  "Short rules work better than long ones."

Issue [#109](https://github.com/mascah/grove/issues/109) reports that
`context` grows the way `CLAUDE.md` files do, one line per thing an agent
once got wrong, and that its author keeps such a budget on their own
instruction files with a hook modeled on rule 12 (kellypj24/claude-toolkit,
`hooks/word-budget.sh`): their global `CLAUDE.md` is at 386 of its 400
words, so the next rule means cutting one.

## Intended outcome

A project owner who wants a budget on `context` and `rules` can copy a
tried recipe from a new guide: one line in the project's own checks, such
as `just check` or the CI job of [Run Grove in CI](../../../docs/guides/run-grove-in-ci.md),
that counts the words of `context` and of every `rules` line in
`grove.yaml` and fails above a number the owner chooses. Grove itself
does not change: no field, no message, no `grove check` problem. The
reference's line on short rules links the guide.

This is option 1 of the issue. Which option is open question 1.

## Scope

### In scope

- `docs/guides/word-budget.md`, trialed on the ledger fixture, with the
  trial's result in the implementing pull request.
- Links to it from `docs/index.md`'s guides list, `mkdocs.yml`'s Guides
  section and `docs/reference.md` after "Short rules work better than long
  ones."

### Out of scope

- A `budget` field in `grove.yaml` and a `grove check` problem reporting
  `context is N words, budget M` (option 2 of the issue): it would be a new
  `check` problem, breaking under `CONTRIBUTING.md` "Breaking changes", and
  needs a change spec to `config`. See open question 1.
- A default budget Grove states: see open question 2.
- Counting the words of a whole skill's instructions, Grove's text
  included, or of files Grove does not read, such as `CLAUDE.md`.
- Any change to Grove's own budget (rule 12).

## Decisions

None made beyond the request. The request (issue #109) names three
options and says "I'd lean toward 1 unless you want Grove to state a
default. Happy to propose whichever you pick." This draft follows that
lean; the owner's pick is open question 1.

## Open questions

This change was drafted with no one in the session. Each question below
carries the recommended answer the documents are written to; merging this
change as written accepts every one. None must be answered before review.
Answering one otherwise means revising the change with
`/grove-propose 261118-7rgct` before it is merged.

1. **Option 1, a guide only, and no change to Grove?** Recommended: yes.
   Rule 8 says a field or message exists only if a documentation step needs
   it, and a guide using `grove.yaml` as it is needs neither; a budget is
   the project's choice of number, so it belongs in the project's checks,
   as rule 12's does in Grove's `just check`; and option 2's new `check`
   problem is breaking. Choosing option 2 replaces this change with one
   holding a change spec to `config` and tasks in `internal/cli/check.go`
   and `internal/doc`; choosing option 3 drops it. Decides T-01 and
   whether the change exists.
2. **The guide states no default budget, and its example's number is
   marked as one to replace?** Recommended: yes. The issue makes a default
   the reason to prefer option 2; a guide that named a number would state
   one without a rule or a check behind it. Decides T-01.
3. **One total: the words of `context` plus every string under `rules`
   (`documents`, `skills` and `delivery`), not counted per skill?**
   Recommended: yes. One number is what the owner keeps under a budget;
   per-skill totals would mean listing which document kinds each skill
   writes, which only Grove's schemas know. The total is an upper bound on
   what any one skill is given beyond Grove's text. Decides T-01.
4. **The guide counts with `yq` and `wc -w`, in a line the project adds to
   its own checks, shown for `just check` and for the workflow of Run
   Grove in CI?** Recommended: yes. `yq` reads YAML as Grove does, so the
   keys of `rules` are not counted as words, and Run Grove in CI already
   relies on tools of the user's (`jq`, `gh`). Decides T-01.
5. **The trial runs on the ledger fixture with a `context` and a `rules`
   line added, showing the check pass under the budget and fail, naming
   the count and the budget, once a line takes it over?** Recommended: yes.
   It is the smallest trial that shows both outcomes. Decides T-01's
   acceptance.

Routine defaults, recommended together: the guide is
`docs/guides/word-budget.md`, titled "Keep a word budget on context and
rules", opening as the other guides do; it is listed after "Running Grove
in CI" in `docs/index.md` and after "Run Grove in CI" in `mkdocs.yml`; the
trial's result goes in the implementing pull request's description.

## Related documents

- Requirements: none. The change alters no requirement: a guide is not
  specification (rule 8), and `grove/specs/config/spec.md` holds as it is.
- Tasks: [tasks.md](tasks.md)
