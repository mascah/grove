# W-011 plan: shaping guide and `grove-shape` adapters

Prepared 2026-09-20 on `worktree-W-011` from `main` `70f19c5`, record revision
`sha256:c38704a1…`. [W-011](../../grove/work/W-011-shaping-entrypoint.md) owns
outcome, scope and acceptance; this is implementation steps only. It replaces
the W-011 portion of the superseded
[shared handoff plan](W-010-W-011-agent-handoffs.md). No Go code changes: the
guide uses today's `list`, `show`, `check`, `new`, `update`, `versions`,
`workspace` and `context`.

## Deliverables

1. `docs/work-shaping.md`: the one shaping guide, written as workflow and not
   repository policy, in the shape of [work-execution.md](../work-execution.md):
   - Inputs: a topic in the person's words and/or record IDs to refine, plus the
     same `--interaction` contract. A shaping invocation is a mandate to write
     proposals and knowledge, never to implement, promote status, or merge.
   - Authority: intent (what the person said), observed evidence (what was
     inspected, with where), proposed design (the agent's), and decisions (only
     with attributable authority) stay distinguishable in what is written.
   - Staged reading: open exploration starts with the direction document and
     `list`; selected work uses `context IDs`; code, related records, branches
     and worktrees are read when the conversation reaches them.
   - Duplicate check before any `new`: `list`, `versions` (other branches and
     live worktrees), and a text search of record bodies; refine an existing
     proposal rather than create an overlapping one, and do not edit a record
     whose newer version is another session's unfinished work.
   - Writing: work via `new work` then body (Outcome, constraints/scope,
     Acceptance, Next) and revision-checked `update` for fields; a stale
     `--expect` refusal means reread, not retry blindly. Questions only for real
     unresolved human choices, with `blocks`; routine technical unknowns are
     investigated. Decisions only for choices a named person actually made in
     the session or a linked source; otherwise `proposed` or a question.
     Supporting knowledge links existing ordinary documents; no new record
     types, fields, statuses or files under the record root (W-019/W-020).
   - Where writes go. Interactive: the session's checkout, stated before the
     first write, unless it is another assignment's execution checkout or holds
     someone's uncommitted record edits; then ask. Commit only shaping's own
     files, when the person agrees. Headless: an isolated new proposal branch
     and worktree, always committed, never merged.
   - Headless adaptation: explicit mandate and declared human availability;
     missing human choice becomes an open question and a returned wait;
     unchanged wait is returned, not reworked; supporting knowledge is listed
     for selective integration; no launcher, timer or supervision.
   - Return: records created/changed with revisions, questions open, authority
     of each decision, the checkout/branch/commit (or "uncommitted") that holds
     them and where today's checkout-scoped board and `versions` will show
     them, `check` result, and that nothing was assigned.
2. `.claude/skills/grove-shape/SKILL.md` and `.agents/skills/grove-shape/`
   (`SKILL.md`, `agents/openai.yaml`): copies of the `grove-work` adapter
   pattern, explicit invocation only, pointing at `AGENTS.md` and the guide.
3. Policy and contract owners: a short AGENTS.md shaping paragraph reusing the
   Assigned-work invocation bullets; a README paragraph beside the `grove-work`
   one; the brief's Observed state and Next; W-011's Next and plan link.

## Evidence (kept apart by kind, in `docs/reviews/2026-09-20-W-011-shaping.md`)

- **Source checks:** adapters differ only in harness frontmatter/argument
  wording; `check` passes; every relative link in changed files resolves.
- **Observed agent behavior, disposable standalone clones** (origin removed,
  own counters, absolute `--project`/cwd, inspected directly afterwards):
  discovery with read-only tools via `claude -p "/grove-shape --interaction
  headless"` and `codex exec --sandbox read-only '$grove-shape …'`; then one
  headless shaping trial per harness on a fixture with an overlapping proposal
  on `main`, related work only on another branch, and a topic containing an
  unmade product choice. Expected: no duplicate, other-branch work found
  through `versions`, a question with `blocks`, a proposal branch, a wait, no
  status promotion or implementation. A rerun with nothing changed returns the
  same wait.
- **Simulation:** stale revision shown with the CLI directly
  (`update --expect OLD` refused after a body edit), labelled as such.
- **Not producible by this session:** acceptance 2 needs the owner's real
  requirements conversation through `/grove-shape` in a fresh session, and
  interactive Codex TUI discovery. W-011 stays active until that is recorded;
  unexercised rows are stated in the review file.

## Review

Documentation and instructions only: one independent reviewer (fresh agent,
read-only) on the final combined revision against the record's acceptance and
scope; at most three fix rounds.
