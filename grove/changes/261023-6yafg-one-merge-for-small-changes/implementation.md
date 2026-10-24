# Implementation: One merge for small changes

<!-- Progress on change 261023-6yafg, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Implement builds a draft with the owner

Status: done

- `internal/cli/instructions.go`: `instructionsChange` reads a change the
  target does not hold from the working tree, with the owner only; with
  `--unattended` it is refused as an unknown change is ("main holds no
  approved change <id>"). A draft that waits is refused as any change is.
  A change the target holds is read from the target, as before.
- `internal/instruct/instruct.go`: `Implement` calls a change the target
  does not hold a drafted change, "as this working tree holds it, with the
  owner", in its task line.
- `internal/instruct/requires/implement.md`: a rule for a change the line
  above calls a draft: build it as the working tree holds it; merging this
  branch approves it with its code, so its proposal, change specs and tasks
  are edited only with the owner's agreement, as when revising a change on
  its branch. The instructions for an approved change gain this one
  conditional bullet; what they require of an approved change is unchanged,
  and their task line and every other rule are as before.
- `internal/instruct/skills.yaml`: implement's description says "an
  approved or drafted change"; `grove init` rewrote
  `.claude/skills/grove-implement/SKILL.md` and
  `.agents/skills/grove-implement/SKILL.md`.
- `internal/cli/walkthrough_test.go`: `TestImplementInstructions` checks a
  change proposed and committed on a branch, not merged: implement
  instructions call it a draft, say merging approves it with its code and
  that its documents are edited only with the owner's agreement; with
  `--unattended` it is refused with "main holds no approved change <id>";
  a draft whose `depends_on` names the unfinished first draft is refused as
  waiting. The approved change's checks are as before. `TestInit` checks the
  description in both skill files.
