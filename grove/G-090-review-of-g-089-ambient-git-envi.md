---
id: "G-090"
type: review
title: "Review of G-089 ambient Git environment guard"
status: current
created: "2026-09-22T17:25:37Z"
updated: "2026-09-22T17:28:10Z"
work: ["G-089"]
examined: "ca6a599"
---

## Examined

Branch `worktree-G-089`, base `main` adce678. Round 1 examined 7924765
(`git diff main..7924765`); round 2 examined ca6a599, the fix commit, which
is the `examined` field. The candidate differs from it only by this record
and [G-089](G-089-ignore-ambient-git-environment-w.md)'s evidence.

## Review

Independent reviewer: a fresh Claude Code reviewer agent (Fable 5.1),
read-only in the worktree, scratch copies and scratch repositories in the
session scratchpad. It reproduced the incident's mechanism with scratch
hooks (a linked-worktree `pre-push` sees `GIT_DIR`; the main checkout's does
not; `pre-commit` sees `GIT_INDEX_FILE` too), proved the regression test
fails without the guard by deleting the one `cmd.Env` line in a copy of the
module, ran the whole suite with `GIT_DIR` and `GIT_WORK_TREE` exported to a
decoy repository and diffed the decoy's refs, worktrees, and config before
and after, tested the candidate variables that were not dropped
(`GIT_CONFIG_PARAMETERS`, `GIT_CONFIG_COUNT`, `GIT_PREFIX`,
`GIT_CEILING_DIRECTORIES` do not move `git -C dir`), and ran the branch's
pre-push line under lefthook 2.1.8 with a real push from a linked worktree
to a scratch bare remote, showing all seven variables gone and
`GIT_COMMITTER_DATE` and `GIT_TRACE` kept.

## Findings

Round 1, no blocking findings:

1. should-fix: the hook unset three of the seven variables while
   `AGENTS.md` said it scrubs them all.
2. should-fix: the Docker reproduction sentence was welded onto the G-089
   rule in `AGENTS.md`, where it does not belong.
3. nit: one `GROVE` spawn in `terminal.py` lacked `env=clean_env()`.
4. nit: pre-commit jobs inherit the variables too, unscrubbed and unexplained.
5. nit, process: the record was still `active` without a candidate.

Round 2 (ca6a599), no blocking findings: 1 to 4 confirmed fixed, the
seven-flag `env -u` line verified under lefthook, the decoy suite rerun
identical, `internal/tui` passing with the decoy exported and no `GOFLAGS`.
Two style nits remain: the 128-character `AGENTS.md` Docker line and the
119-character `lefthook.yml` comment.

## Disposition

- 1 to 4: fixed in ca6a599.
- 5: closed by the handoff that carries this record.
- Round 2 style nits: left as they are.
