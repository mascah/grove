<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Grove reads, gates and asks nothing about needs_human

Result: no Go code, page, template or instruction text reads, writes, asks
about or reports `needs_human`, and plan titles the cycle itself.

Files:
- `internal/doc/doc.go`
- `internal/doc/meta.go`
- `internal/cli/claim.go`
- `internal/cli/serve.go`
- `internal/cli/page.html`
- `internal/instruct/instruct.go`
- `internal/instruct/metadata/cycle.yaml`
- `internal/instruct/requires/plan.md`
- `internal/instruct/requires/propose.md`
- `internal/instruct/requires/implement-owner.md`
- `internal/instruct/requires/implement-unattended.md`
- `internal/instruct/schemas/default/skills/plan.md`
- `internal/instruct/skills.yaml`
- `internal/cli/walkthrough_test.go`

Interfaces: `doc.Meta` has no `NeedsHuman`; the cycle object `grove serve`
gives has no `needs_human` key.

Acceptance:
- `grove check` on a cycle file holding `needs_human` reports that line as
  an unknown field (cycles, "Report a malformed cycle"), and reports
  nothing else for a marked cycle naming a member there.
- `grove claim K` claims a buildable cycle whose file still holds
  `needs_human` (runs, "Refuse a given cycle that cannot be built").
- `/api/project` gives no `needs_human` key and no needs-a-person verdict;
  readiness for a cycle with a member waiting outside it names that member
  (serve, "Readiness of each cycle", "The page").
- The plan and propose instructions ask nothing about a change needing a
  person, plan's instructions have the agent write the cycle's title
  itself, and the implement instructions name no `needs_human` (cycles,
  "Plan a cycle", "Say whether the cycle can run unattended"; propose,
  "Re-plan the cycle of a revised member"; implement, "Build a change, a
  draft or a cycle").
- `just check` passes and prints fewer instruction words than 8,613.

## T-02 — Docs, walkthrough and Grove's own cycle file

Result: the documentation, the walkthrough and Grove's own documents say
nothing of `needs_human` but what a project that still has it must do.

Files:
- `docs/concepts.md`
- `docs/coordinating.md`
- `docs/planning.md`
- `docs/reference.md`
- `docs/running.md`
- `docs/watching.md`
- `dev/walkthrough.md`
- `grove/cycles/261023-9kfxn-smaller-and-plainer.yaml`

Acceptance:
- No file outside `grove/changes/archive/` and `CHANGELOG.md` names
  `needs_human` or "needs a person", except where `docs/reference.md`
  says a cycle file holding `needs_human` gets an unknown-field report and
  the key is deleted.
- `grove check` on this repository is clean.
