---
id: "G-082"
type: review
title: "G-040 portable bootstrap review"
status: current
created: "2026-09-22T16:15:58Z"
updated: "2026-09-22T16:28:25Z"
work: ["G-040"]
examined: "66fcccc"
---

## Examined

Branch `worktree-G-040`, base main `ccdc92d`, on 2026-09-22. Two kinds of
evidence, kept apart: the **independent review** by a reviewer subagent that
edited nothing (round one at `731674e` plus two then-uncommitted edits that
became `eb32209`; round two at `66fcccc`), and this session's **trial
evidence** for [G-040](G-040-portable-bootstrap.md)'s acceptance, each
observation labelled real (a fresh harness session), simulation (this
session driving the CLI in a disposable repository) or unexercised. Harness:
Claude Code 2.1.278 (`claude -p` on Claude Opus 5) and codex-cli 0.155.1,
go 1.26.2, git 2.55.0, macOS Darwin 25.6.0. The trial binary was built with
`go build` from a clean clone of the candidate commit (see item 3 below for
why not from the worktree); `grove version` printed
`grove v0.0.0-20260922161140-7ddbf528ccab (7ddbf528ccab…)` for the first
trials and `… 66fcccca0034 (66fcccca…) guides sha256:8e6d8b9d5bb7` after
round one.

## Findings

Round one (independent, 41 tool uses; it ran the suite, `init` against twelve
disposable repositories including symlinked paths, relative and default
`--project`, and thirteen `grove.yaml` shapes through a main-built and a
branch-built binary): no correctness defect in `init`'s plan-then-write logic
or in the `parseConfig` split, whose `check` and `brief --json` output was
byte-identical across all shapes. Findings and dispositions:

1. Consequential: bare `grove init` in the README runs the predecessor's
   `init` when that is first on PATH, silently. Fixed in `66fcccc`: the block
   runs `grove version` first in the same shell and says what a wrong answer
   looks like; the collision itself cannot be guarded from another binary.
2. Consequential: after `init` a target has no statement of how the CLI is
   invoked, while the printed guide forbade assuming PATH. Fixed: both guides
   now defer to the repository's instructions "or, where they say nothing,
   the entrypoint that loaded this guide", and `init` ends with a stderr note
   naming what `AGENTS.md`/`CLAUDE.md` may add.
3. Minor: a binary built in this linked worktree stamped main's revision,
   confidently wrong. Documented in the README (Go recognises only a `.git`
   directory), the evidence binary is built from a clone, and `version` now
   ends with a digest of the embedded guides.
4. Minor: `go install …@COMMIT` resolves only a pushed commit. Documented.
5. Minor: no test covered the printed guides' portability. Fixed: the test
   scans both guides, the four adapters and a policy file for this
   repository's paths, `go run`, and worktree names.
6. Minor: "(`AGENTS.md` here)" and two G-078 parentheticals were this
   repository's. Fixed.
7. Minor: the nested-directory message was ungrammatical and the default
   (no `--project`) path untested. Fixed and tested.
8. Nit: a symlinked parent directory was written through. Fixed: any
   symlinked component of a written path is a conflict, tested.
9. to 11. Nits (README blank line, usage grouping, a needless Builder). Fixed.

Round two (independent, at `66fcccc`; it rebuilt, reran the suite, and
mutation-tested the new assertions: appending `go run` to a guide and
deleting the parent-symlink loop each fail a test): no regression, and
findings 2, 5 and 8 verified by running. Residuals, fixed in `3113280` and
self-checked there, not independently re-reviewed: R1, minor, the record
root bypassed the parent-symlink check, so a symlinked `docs` let `records:
docs/records` be created outside the checkout; now a conflict, with a test
that fails when the guard is removed. R2, minor, the guide named no fallback
for branch names where a repository's instructions say nothing; now "choose
a clear name and state it in the handoff", which is what the three headless
runs did. R3, nit, a README sentence in the wrong paragraph; moved. R4, not
a defect: this record was a skeleton when reviewed; it is filled here.

## Trial evidence per acceptance item

1. **Initialize, validate, shape, hand off** (real and simulation). In a
   disposable repository under the session scratchpad with the clone-built
   binary first on PATH: `init` created nine paths, `check` printed `OK: 0
   records`, `new work` allocated `G-001`. Shaping: `claude -p "/grove-shape
   a --version flag for the program --interaction headless"` (73 s, 10
   turns) read the guide from the binary, listed and inspected the clone,
   made a proposal branch and worktree, wrote a proposed work record and a
   blocking question with options and a recommendation, committed `7c59eb6`,
   and returned the wait, which is the question path G-078 had seen a
   headless session skip. Hand-off: `claude -p "/grove-work G-001
   --interaction headless"` (55 s, 14 turns), the headless work row G-078
   left unexercised, ran `grove guide work`, `context G-001 --interaction
   headless`, `versions`, made branch `work/G-001-greeting` in a worktree,
   found the stub record had no outcome or acceptance, and took the
   missing-decision path: `new question`, `blocks`, a checkpoint in G-001's
   Next with "no plan needed", commit `af6d9dc`, and the wait returned with
   nothing implemented. No predecessor command was used in either run. Its
   first `grove guide work G-001 …` was a usage error it recovered from; the
   adapter now says the command takes no other argument.
2. **Rerun preserves and reports** (simulation). A second `init` printed
   `kept` for the configuration, root and brief and `unchanged` for six
   entrypoints. After a marked edit and a takeover: `updated` restored the
   template exactly and the unmarked file was `kept` with the note. A
   schema 2 `grove.yaml`, a file at the record root, a directory at a managed
   path, a symlinked parent, and a directory below the top each produced a
   conflict line and exit 1 with the tree unchanged. A real upgrade: the
   target set up by the `7ddbf52` build received three `updated` entrypoints
   from the `66fcccc` build and kept the user-owned one.
3. **Unambiguous version and coexistence** (simulation and real). `version`
   from a clone build names the commit; from `go install` into a scratch
   `GOBIN` the same; from `go run` `(devel)`; from this linked worktree the
   enclosing checkout's commit (documented). The predecessor at
   `~/.local/bin/grove` was never replaced and rejects `version`, `guide`
   and, through the entrypoints, everything: real observation in the first
   Codex run, `codex exec --sandbox workspace-write '$grove-shape …
   --interaction headless'` (38 s) ran through a login shell whose PATH put
   the predecessor first; the generated adapter's rule fired, the session
   returned "Blocked by a Grove CLI mismatch: `grove guide shape` fails
   because the installed CLI has no `guide` command", and wrote nothing.
   Documented with the two remedies. Codex also loaded the predecessor's
   globally installed Codex plugin skill first in both runs and, in the
   second, tried its `grove status` before following the generated adapter.
4. **Fresh-session discovery** (real). Claude's init event listed
   `grove-shape` and `grove-work` among the slash commands in the target, and
   both runs above followed the generated adapters. Codex, rerun after the
   target's `AGENTS.md` named the executable (216 s), read the generated
   adapter, ran the named binary's `guide shape`, followed the guide (brief,
   list, versions, branches), created `proposal/version-flag` in a worktree
   under `/private/tmp`, wrote a proposed work record and a blocking
   question, committed `4162790`, and returned the wait. Interactive typing
   of the commands in a target is unexercised.
5. **Package boundary** (source check and test). The printed guides and the
   generated files name no `go run`, worktree name, Go test, `AGENTS.md`
   policy or record of this repository; the test in `init_test.go` enforces
   it. This repository's own adapters stay unmanaged and `init` here reports
   every path `kept`.

## Disposition

Every consequential finding was fixed and re-reviewed; every minor finding
and nit was fixed, the last three self-checked after the second round. No
finding is open. The trial evidence meets acceptance items 1 to 5 with the
limits stated in G-040's Evidence; the owner's judgment of the adoption
experience is not supplied by this record.
