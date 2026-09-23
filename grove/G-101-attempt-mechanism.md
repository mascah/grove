---
id: "G-101"
type: decision
title: "Run attempts as a Grove-owned claude -p process"
status: accepted
created: "2026-09-23T02:33:53Z"
updated: "2026-09-23T02:34:39Z"
---

## Decision

An implementation attempt that Grove starts (`grove run`, [G-045](G-045-durable-attempt.md))
is a Grove-owned `claude -p "/grove-work ID --interaction headless"`
process under an on-demand owner process, with structured `stream-json`
events written to files, a required dollar budget and permission mode,
`--permission-prompts none`, a session id Grove generates, and Grove-owned
worktree, identity, duplicate-start refusal, Stop and owner-loss handling.
Ordinary CLI access stays daemon-free: nothing resident is required to
launch, inspect or stop an attempt. The selection and its evidence are in
plan [G-100](G-100-g-045-durable-attempt-plan.md); capabilities are pinned
to Claude Code 2.1.280 and recorded per attempt from the provider's init
event.

## Alternatives

- **Native background sessions** (`claude --bg`, `agents`, `attach`,
  `logs`, `stop`, `rm`): interactive sessions hosted by Claude's supervisor
  daemon. Claude Code rejects `--bg` with `-p`; `logs` is terminal text,
  not events; `--max-budget-usd` works only with `--print`; permission
  prompts wait for an attach; Claude moves the session into its own
  worktree; state under `~/.claude/jobs/` is not a stable interface. The
  right tool for a person's parallel interactive work, not for a bounded
  assignment whose result Grove reconciles.
- **A tmux host**: persistence and operator access for the same process,
  plus a dependency and a socket namespace, with no attempt bookkeeping. It
  can host the owner later if operator attachment is wanted.
- **The predecessor's `claude-p.sh` adapter**: reused an unvalidated
  worktree, skipped permissions, and filled in a tested HEAD it had not
  tested ([G-024](G-024-predecessor-work-review.md)); its pattern of an
  owned process with raw capture is kept, the rest is not.

## Reconsideration

Reconsider when native sessions gain print-mode events, budgets and a
stable state interface; when scheduling or batches need a resident
service; or when reboot recovery becomes a selected requirement.
