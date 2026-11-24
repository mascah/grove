<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Propose writes every part of a split, against its dependency's text

Result: `grove instructions propose` has the agent write every agreed part
of a split in the session, and its change-spec guide says how to write a
requirement that a change this one depends on adds or modifies, as "Offer
to split a change" and "Write a change spec against a change it depends
on" say.

Files:
- `internal/instruct/requires/propose.md`
- `internal/instruct/schemas/default/documents/change-spec/guide.md`
- `internal/cli/walkthrough_test.go`

Constraints:
- No Go outside the test. `grove check` already accepts the text a
  dependency gives a requirement.
- The instruction text stays under the 9,000-word budget `just check`
  counts.

Acceptance:
- The output of `grove instructions propose "<description>"` says that,
  once the owner agrees to split, every agreed part is written in this
  session, each later part against the requirements the earlier parts
  leave, and none deferred because an earlier part rewrites the same
  requirements. The propose instructions test checks for that text.
- The same output's change-spec guide says that a requirement a change in
  `depends_on` adds or modifies goes under MODIFIED, from the text that
  change gives it, naming that change, even where the capability has no
  maintained spec at `base`; and `grove instructions propose --schema
  bugfix "<description>"` prints the same guide text. The test checks it.
- On a disposable repository, change A ADDS a requirement to a capability
  with no maintained spec, change B depends on A and MODIFIES that
  requirement from A's text, and `grove check` reports neither, as the
  scenario "A requirement only the dependency adds" says.
- `just check` passes.
