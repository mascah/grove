# Changes to propose

The maintained spec is [grove/specs/propose/spec.md](../../../../specs/propose/spec.md).

These rules are stated today only in `docs/coordinating.md` ("A change
that waits", "When the plan is stale") and `docs/reference.md` (`grove
check`). The code already does what they say; this change makes them
agreed.

## ADDED Requirements

### Requirement: Report a dependency that can never be met

When the owner runs `grove check`, the system shall report, for each open
change, each ID in its `depends_on` that names no change, open or
archived, that names a change archived as dropped, or that names an open
change which waits on this one in turn, directly or through other open
changes, at the `depends_on` line of the change's `.grove.yaml`.

#### Scenario: No such change

- GIVEN open change B's `.grove.yaml` holds `depends_on: [260101-zzzzz]`,
  and no change, open or archived, has that ID
- WHEN the owner runs `grove check`
- THEN it reports B's `.grove.yaml` at its `depends_on` line, saying
  260101-zzzzz is no change, open or archived, and exits non-zero

#### Scenario: A dropped change

- GIVEN open change B depends on A, and A is archived with
  `implementation.md` saying `Status: dropped`
- WHEN the owner runs `grove check`
- THEN it reports B's `depends_on`, saying A was dropped, so B could never
  be applied

#### Scenario: A loop

- GIVEN open change A depends on B, and open change B depends on A
- WHEN the owner runs `grove check`
- THEN it reports the `depends_on` of each, saying the other waits on it in
  turn

#### Scenario: A dependency archived as finished

- GIVEN open change B depends on A, and A is archived with every task
  `Status: done`
- WHEN the owner runs `grove check`
- THEN B's `depends_on` is not reported

### Requirement: Report a stale requirement

When the owner runs `grove check`, the system shall report each
requirement an open change MODIFIES or REMOVES whose text in the working
tree's maintained spec is neither its text at the change's `base`, nor
the text this change gives it, nor the text a change it depends on gives
it, naming the change spec's file and line, saying whether the
requirement was rewritten or removed since `base`, and saying to revise
the change with propose.

#### Scenario: Rewritten by another change

- GIVEN change S, proposed at `base` 4f1c2a9, MODIFIES "Keep a category"
  of `record-editing`, and change T, which S does not depend on, rewrote
  that requirement in the maintained spec and was archived
- WHEN the owner runs `grove check`
- THEN it reports S's change spec at the requirement's line, saying "Keep
  a category" was rewritten since base 4f1c2a9, and to revise the change
  with propose

#### Scenario: Removed by another change

- GIVEN change S MODIFIES "Keep a category", and another change it does
  not depend on removed that requirement from the maintained spec
- WHEN the owner runs `grove check`
- THEN it reports that "Keep a category" was removed since S's `base`

#### Scenario: Rewritten by a change it depends on

- GIVEN change S depends on T, and T's archive gave "Keep a category" the
  text T's change spec gives it
- WHEN the owner runs `grove check`
- THEN S is not reported

#### Scenario: Already reconciled

- GIVEN the maintained spec gives "Keep a category" the text S's own
  change spec gives it, as archiving S leaves it
- WHEN the owner runs `grove check`
- THEN S is not reported

## Observed current behavior

- `grove check` reports a `depends_on` naming no change, open or archived,
  a dropped change, or an open change that reaches back to this one,
  with the messages "is no change, open or archived", "was dropped, so
  this change could never be applied" and "waits on this change in turn,
  so neither can be applied" (`internal/doc/doc.go`, `dependencies`, at
  0c6f419).
- `grove check` compares each MODIFIED or REMOVED requirement with the
  maintained spec in the working tree, at `base`, and as the change and
  the changes it depends on give it, and says "rewritten" or "removed"
  (`internal/cli/check.go`, `stale`, at 0c6f419).
- `docs/coordinating.md` and the `grove check` row of `docs/reference.md`
  state both rules at 0c6f419; no maintained spec does.
