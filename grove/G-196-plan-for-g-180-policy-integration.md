---
id: "G-196"
type: plan
title: "Plan for G-180 policy-driven integration"
status: current
created: "2026-09-26T03:20:52Z"
updated: "2026-09-26T03:21:19Z"
work: ["G-180"]
---

## Design

Plan for [G-180](G-180-policy-driven-integration.md) under decision
[G-182](G-182-standing-policy-delegation.md). The choices the record leaves
to preparation:

- **Parsing.** `policy:` in `grove.yaml` beside `run:`, read into
  `project.Policy`, absent by default. Keys as the record proposes:
  `budget` (required once `resolve` is present: the policy spends),
  `resolve.budget` (optional; `run.budget` otherwise), `approve.verify`
  (required, non-empty, one shell command each), `approve.max_lines`
  (optional, positive), `approve.never` (optional patterns, relative to the
  project root; `path.Match` syntax plus a trailing `/**` for a directory's
  whole tree), `integrate` (a boolean that needs `approve`). Every problem is
  a `check` diagnostic named `policy.KEY`, so a malformed policy stops every
  command as a malformed `run:` does, and never half-applies.
- **Reading a review's findings.** The review guide asks the reviewer to end
  with exactly one closing line, `Open findings: none` when no finding,
  knowledge findings included, remains open on the commit it examined, and
  `Open findings: N` otherwise; the work guide asks the author to copy the
  last round's closing line into the review record. The policy counts a
  review record whose `work` names the ID, whose status is `current`, whose
  last line starting `Open findings:` is `Open findings: none`, and whose
  `examined` is the candidate or an earlier commit from which only files
  under the record root changed (the handoff's evidence and the review
  record itself). No schema field: prose read by a deterministic rule, as
  the record allows.
- **Attribution.** A delegated approval calls `update.Approve` with the
  verdict `delegated under policy grove.yaml sha256:REV: review G-N examined
  X with no open finding; merged with TARGET at T, verification passed
  (COMMANDS); ATTEMPT` where ATTEMPT names the latest finished attempt whose
  result holds the candidate and its cost, or says none did. It keeps the
  `Verdict on candidate X, DATE:` prefix. A delegated integration appends to
  the done update `Integrated under policy grove.yaml sha256:REV as MERGE on
  TARGET (was B); to reverse it: git revert …`. A delegated resolution's
  feedback begins `delegated under policy grove.yaml sha256:REV, budget $N:`
  before G-178's mandate. The board's Review block says `approved under
  policy` instead of `approved` when the candidate's verdict is delegated.
- **Verified merge.** In a temporary worktree outside every checkout
  (`git worktree add --detach` at the target commit T the prediction used),
  merge the branch tip, run each `verify` command there with `sh -c` in the
  project directory and Git's location variables removed, then remove the
  worktree. Only after every command passed: approve in the branch's
  checkout, then `integrate` in the target's checkout with `Expect: T`,
  which refuses before merging if the target moved since verification.
- **Trigger.** An explicit `grove sweep [--dry-run]` from any checkout, which
  the owner or an external scheduler runs; no resident process and no hook
  in the attempt owner. `--dry-run` prints what would happen to each
  candidate and why and writes nothing; the sweep prints the same plan, then
  one line per act. Without a policy both say nothing is automatic.
- **Per candidate**, each candidate in review on a branch other than the
  target: a group sharing a candidate, an approved one, an open blocking
  question, several branches, or no checkout waits for the owner with the
  reason; one already in the target is skipped. A conflict is resolved
  through `attempt.Resolve` when `resolve` is present, the record holds no
  earlier resolution feedback naming the same target commit (once per
  target movement), and the aggregate budget has room; otherwise it waits.
  A clean merge is approved when `approve` is present and every condition
  holds: no commit after the candidate but records, a qualifying review, no
  changed path matching `never` or outside the project, no binary change,
  and added plus removed lines (the record's own file excluded) within
  `max_lines`; then verified, approved and, with `integrate`, integrated.
- **Out of scope**, as the record says, plus: re-verifying and integrating a
  delegated approval whose integration was refused (it waits for the owner's
  `grove integrate`), and a board key for the sweep.

## Steps

1. `internal/project`: `Policy` and its parsing and diagnostics; tests.
2. `internal/integrate`: `Request.Expect` and `Request.Policy` (the done
   update's attribution and revert line); tests.
3. `internal/attempt`: `Request.Policy` prefixes Resolve's feedback.
4. New `internal/sweep`: `Plan` (evaluation, no writes) and `Run` (acts);
   tests with real repositories and a fake provider on clean, conflicting,
   failing-verification and out-of-policy candidates, and no policy.
5. `internal/cli`: `grove sweep [--dry-run]`, usage and help.
6. `internal/tui`: `approved under policy` in the Review block.
7. Docs: `docs/work-review.md` (closing line), `docs/work-execution.md`
   (copy it into the review record; the sweep among the dispositions),
   `docs/commands.md`, `docs/record-model.md` (`policy:`),
   `docs/board.md`, README if it lists commands; G-059 term only if its
   meaning changes (G-182 says it does not).
8. Verification per `CLAUDE.md`, independent review, handoff.
