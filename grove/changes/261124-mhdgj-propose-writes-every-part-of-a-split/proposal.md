# Propose writes every part of an agreed split

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

Issue #201 reports a propose session in another project. The owner agreed
in one round to split a redesign into four changes, each building on the
last. The agent wrote only the first, and recorded that each later change
would be proposed once the one before it was built:

> The later changes rewrite the same requirements (the window's layout,
> the right pane's rail), so one written earlier would be stale.

So each link of the series cost the owner another round, between builds.
The owner said of it: "I don't love that we didn't write proposals for
these all at once and instead are proposing between each change".

Grove already supports writing the series at once. At acf076d:

- The propose spec's "Offer to split a change" (`grove/specs/propose/spec.md`)
  writes several changes on one branch, each from its own
  `grove instructions propose`, linked with `depends_on`. It does not say
  that every agreed part is written in that session.
- Its "Report a stale requirement", scenario "Rewritten by a change it
  depends on", has `grove check` accept a requirement whose text is the
  one a change it depends on gives it.
- `rewritten` in `internal/cli/check.go` skips a MODIFIED or REMOVED
  requirement absent from the maintained spec at the change's `base`, so a
  change that modifies a requirement only its dependency adds is not
  reported either.

The instructions do not say so. The change-spec guide
(`internal/instruct/schemas/default/documents/change-spec/guide.md` at
acf076d) covers only a capability with a maintained spec and one without,
where it says "Do not use MODIFIED or REMOVED: nothing agreed exists to
modify". It says nothing of a requirement that a change in `depends_on`
adds, modifies or removes. Read literally, the second part of a series
cannot modify a requirement the first adds, and has no text to modify
from, so waiting for the first to be archived is a fair reading. The split
bullet in `internal/instruct/requires/propose.md` does not say to write
every part now either.

## Intended outcome

Once the owner agrees to split, the propose session writes every agreed
part, defers none, and writes each later part against the requirements the
earlier parts leave. The change-spec instructions say how: a requirement
that a change this one depends on adds or modifies is written under
MODIFIED, from the text that change gives it, naming that change, even
where the capability has no maintained spec at `base`. `grove check`
already accepts that text, so no Go changes. The `bugfix` schema uses the
same change-spec guide and gets the same rule.

## Scope

### In scope

- The propose spec: "Offer to split a change" says every agreed part is
  written in the session, and a new requirement says how a change spec is
  written against a change it depends on.
- The propose instructions' split bullet and the change-spec guide.
- The propose instructions test in `internal/cli/walkthrough_test.go`.

### Out of scope

- Go outside that test. `grove check` and `grove archive` already handle a
  requirement given by a dependency.
- Explore. Its instructions do not plan or write a series of changes; they
  hand off to propose, which owns the split.
- `docs/` and `dev/walkthrough.md`. No page describes splitting, so none
  becomes false.
- `/grove-propose --unattended`, which writes one change and names a split
  as an open question ("One request, one change").

## Decisions

The owner accepted each recommended answer of the round with "yes".

1. **Schema `default`, not `bugfix`.** Reason: the change modifies the
   agreed "Offer to split a change", which `bugfix` does not allow.
2. **Two changes to the propose spec.** "Offer to split a change" says
   every agreed part is written in the same session. A new requirement
   says a change spec writes a requirement from the text a change it
   depends on gives it, naming that change, including a requirement only
   that change adds. Reason: those are the two gaps the issue names.
3. **Only the propose instructions and the change-spec guide change.**
   `bugfix` uses the same guide, so both schemas are covered. Explore is
   left out. Reason: explore never plans or writes a series; propose does.
4. **`docs/` is left alone.** Reason: the docs do not describe splitting,
   so nothing in them becomes false.

## Open questions

None must be answered before review.

1. **Routine choices in T-01**, which it decides: the exact wording of the
   instruction text, and which phrases the test checks. Recommended: the
   implementer chooses them, keeping the instruction text under the
   9,000-word budget (8,388 words at acf076d).

## Related documents

- Requirements: [specs/propose/spec.md](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
