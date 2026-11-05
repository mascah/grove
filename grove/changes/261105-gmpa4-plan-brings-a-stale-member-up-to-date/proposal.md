# Plan brings a stale member up to date

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

Plan writes only the cycle file: "Its file is the only thing you write:
edit no change's files, and no code"
(`internal/instruct/requires/plan.md` at 6090cb3). It says nothing about a
member that `grove check` reports stale on its base.

Writing the cycle file then hides that report: `movedSinceBase`
(`internal/cli/check.go` at 6090cb3) skips a member of a cycle not done,
leaving it to the plan. Implement then counts a member's named files only
since the cycle's `planned` commit (the implement spec's "A task's file
changed on the target", scenario "A member of a cycle"). So the member is
built on whatever its tasks said before main moved.

Observed in this repository on 2026-11-04, planning cycle 261104-z05ey:
`grove check` reported "T-02's dev/trials/c.sh changed on main since base
1a24167; revise the change with propose" for member 261101-mg8hd, because
its dependency 261029-wtbc0 had landed (e862d60). The drift mattered:
261029-wtbc0 added a `good` function and a "Checked by `good`" paragraph
that T-02's acceptance did not name. The planning session stopped and
sent the owner to a separate `/grove-propose 261101-mg8hd` before it
re-planned. The owner then asked for both to be done in one go, and for
"a proposal to help minimize these stops for a human like this". Both
went into one pull request, #120.

## Intended outcome

Planning a new cycle, the agent runs `grove check` before writing the
cycle file. For each member it reports stale on a requirement, or for a
named file changed since its base, the agent tells the owner that this
edits approved work. It then brings the member up to date on the same
branch, as `grove instructions propose <change-id>` says for such a
change, and goes on with the plan. One merge brings the revision and the
plan, with no second session. The rule is in
[specs/cycles/spec.md](specs/cycles/spec.md).

## Scope

### In scope

- The plan instructions (`internal/instruct/requires/plan.md`) and their
  wording test.
- The cycles spec's requirement "Plan a cycle".
- `docs/planning.md` "Plan a cycle" and the walkthrough's "Plan a cycle"
  step.

### Out of scope

- Go code, and `grove check`'s messages.
- `CONSTITUTION.md`: the revision is approved by the merge that brings it
  (rule 4), and Grove gains no gate (rule 10).
- Re-planning a cycle that already exists. Its members are not reported
  stale on their base by `grove check`, since their plan reports them.
- The propose instructions, and the skill trials.

## Decisions

The owner chose the direction, plan bringing a stale member up to date in
the same session, over having `grove check` skip the merges that archived
a change's own `depends_on` ("Plan re-bases in-session", 2026-11-04), and
accepted every answer below in one round ("Yes to all", 2026-11-04).

1. **Plan revises a stale member without asking first, saying it changes
   approved work, and asks only where bringing it up to date needs a
   decision.** Reason: the point is no extra stop, and the owner reads the
   revision in `grove review` before merging.
2. **Both of check's base reports count: a requirement rewritten since
   base, and a named file changed since base.** Reason: propose brings
   both up to date the same way.
3. **The plan instructions point to `grove instructions propose
   <change-id>` rather than restating propose's rule.** Reason: the rule
   for bringing a change up to date stays in one place, and the plan text
   grows by about 60 words, with the budget at 8,023 of 9,000 words.
4. **The cycles spec's "Plan a cycle" is MODIFIED** to allow that edit,
   with a new scenario "A member stale on its base". Reason: it is the
   requirement that says plan edits no change's files.
5. **Routine defaults:** the `default` schema; one task, touching
   `internal/instruct/requires/plan.md`, `internal/cli/walkthrough_test.go`,
   `docs/planning.md` and `dev/walkthrough.md`; `depends_on:
   [261030-t2pt7, 261101-mg8hd]`, since both edit those files and landing
   first would make cycle 261104-z05ey's plan stale; a `Delegated:` line
   for test layout and private helper names.

## Open questions

None.

## Related documents

- Requirements: [specs/cycles/spec.md](specs/cycles/spec.md)
- Tasks: [tasks.md](tasks.md)
