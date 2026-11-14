## T-01 — One set of implement instructions, ending at a branch

Result: implement asks no one, records each choice as a `Check:` line on a
done task, builds on a drifted file, revises the change on its branch, and
blocks only a costly choice, as the implement change spec's requirements
say; `--unattended` is gone from the command, the skill and the
instructions.

Files:
- `internal/instruct/requires/implement.md`
- `internal/instruct/requires/implement-owner.md`
- `internal/instruct/requires/implement-unattended.md`
- `internal/instruct/requires/implement-cycle.md`
- `internal/instruct/skills.yaml`
- `internal/instruct/instruct.go`
- `internal/instruct/schemas/default/documents/implementation/guide.md`
- `internal/instruct/schemas/default/documents/implementation/template.md`
- `internal/instruct/schemas/default/documents/proposal/guide.md`
- `internal/instruct/schemas/default/documents/tasks/guide.md`
- `internal/cli/cli.go`
- `internal/cli/instructions.go`
- `internal/cli/walkthrough_test.go`

Interfaces: `grove instructions implement <change-id> | <cycle-id>`, with
no flag; `instruct.Implement` and `instruct.Cycle` lose their
`unattended` argument.

Acceptance:
- `just check` passes, with fewer instruction words than 8,613.
- `grove instructions implement <id> --unattended` prints the usage and
  exits 2.
- The implement instructions hold no "with the owner"/"unattended" split,
  no stop and no `Status: to check`; the walkthrough test checks the new
  text.

## T-02 — Docs, walkthrough and specs say the same

Result: the user documentation, the walkthrough and the maintained specs'
non-requirement text describe implement as T-01 builds it, and no
documented step uses `--unattended`.

Files:
- `docs/`
- `dev/walkthrough.md`
- `README.md`
- `grove/specs/implement/spec.md`
- `grove/specs/cycles/spec.md`
- `grove/specs/runs/spec.md`

Acceptance:
- `git grep -- --unattended` finds it only in `CHANGELOG.md` and archived
  changes.

## T-03 — Skill trials expect no stop

Result: trial (d) counts a good run as one that builds on the drifted
`ledger/report.py` without asking and writes a `Check:` line for it;
trial (c) has no `--unattended` run.

Files:
- `dev/trials.md`
- `dev/trials/c.sh`
- `dev/trials/d.sh`
- `dev/trials/run.sh`

Acceptance:
- `just trials` runs every case and its table goes in the pull request.
