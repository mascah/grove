## T-01 — Plan cycles

Result: cycle files, `/grove-plan`, the `unattended` mark, `grove check`'s
cycle reports and `grove changes` and `grove review` showing cycles behave
as [specs/cycles/spec.md](specs/cycles/spec.md) says, and the
documentation specifies them.

Constraints: within the rule 12 budget; messages are one sentence (rule
14). `phase-3-cycles` may be reused where it fits; it is not the spec.

Acceptance:
- A docs page shows planning a cycle from a description and from one
  change ID, plan asking about a change that needs a person, merging the
  plan, and marking the cycle `unattended` in a later merge;
  `docs/reference.md` lists the skill, `grove instructions plan`, the cycle
  file's fields and `rules.skills.plan`. `dev/walkthrough.md` runs those
  steps (rule 8).
- On the ledger fixture, in Claude Code: `/grove-plan` writes a cycle of
  two approved changes, asks about one whose task needs a human eye,
  records the owner's answer, and edits no change's files.
- `internal/cli/walkthrough_test.go` covers `grove check` reporting a
  marked cycle with `needs_human`, a marked stale cycle, a malformed cycle
  file (one without `created` among them) and one change in two open
  cycles; `just check` passes.

## T-02 — Run, apply and implement take cycles

Result: `grove run`, `/grove-apply <cycle-id>` and `/grove-implement
<cycle-id>` behave as [specs/runs/spec.md](specs/runs/spec.md),
[specs/apply/spec.md](specs/apply/spec.md) and
[specs/implement/spec.md](specs/implement/spec.md) say, and apply's
change-ID and no-ID forms are gone.

Constraints: the pull request title marks the change breaking, with a
`BREAKING CHANGE:` footer saying how to move from `/grove-apply` to
`/grove-implement` or a cycle.

Acceptance:
- `docs/running.md`, `docs/coordinating.md`, `docs/reference.md`,
  `docs/index.md`, `README.md`, `docs/guides/run-apply-on-a-timer.md` and
  `docs/guides/run-grove-in-ci.md` describe the runner building only
  marked cycles and type no `/grove-apply` without a cycle ID;
  `docs/reference.md` lists `grove instructions apply <cycle-id>
  <change-id>`.
  `dev/walkthrough.md` runs the changed steps.
- On the fixture with a remote and Claude Code: an approved change in no
  cycle is not built by `grove run`; after planning it as a cycle of one
  and marking it, `grove run` builds it, and the project is paused until
  the pull request is merged.
- `internal/cli/walkthrough_test.go`, with a stand-in `claude`, covers: no
  eligible cycle, an unmarked cycle, a stale marked cycle, two marked
  cycles planned the same day selecting the earlier `created`, the pause
  per cycle, `grove instructions apply` refusing a change ID and no ID,
  and `grove instructions apply <cycle-id> <change-id>` printing the
  member's instructions with the cycle's plan and refusing a change
  outside the cycle and a member waiting on one not archived; `just check`
  passes.
