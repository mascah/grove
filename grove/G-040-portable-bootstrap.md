---
id: "G-040"
type: work
title: "Bootstrap projects with portable Grove workflows"
status: active
created: "2026-09-21T00:54:15Z"
updated: "2026-09-22T16:06:20Z"
kind: feature
size: medium
priority: 2
depends_on: ["G-039"]
relates_to: ["G-035", "G-036", "G-025", "G-037", "G-038", "G-064", "G-065"]
formerly: "W-022"
---

## Outcome

An existing repository can adopt a versioned Grove CLI and shared workflows,
then shape and execute work through thin Claude/Codex entrypoints without a
dependency on this repository's development-only instructions.

## Scope and bounds

Provide minimal setup with valid `grove.yaml`, discoverable brief/knowledge
locations and harness entrypoints. Create directories as needed and preserve
existing AGENTS/CLAUDE instructions. Setup can leave a clearly incomplete brief
for an interactive shaping session; it must not invent project intent.

The target default uses G-065's one root, flat creation, neutral new IDs and
general knowledge pages, as selected in
[G-064](G-064-stable-knowledge.md). Do not pre-create type folders,
require knowledge classification, add per-type routing/prefix settings, or move
existing records when setup is rerun. Keep shared workflow
instructions versioned with Grove; adapters load one owner rather than divergent
editable copies. Choose packaging/update behavior during preparation. A TUI
setup wizard is optional later, not a prerequisite for this outcome.

## Acceptance

1. A disposable existing repository can initialize, validate, shape and hand
   off selected work through the new CLI without predecessor initialization.
2. Re-running setup preserves user content and makes managed updates explicit;
   conflicting configuration or instruction blocks receive actionable outcomes.
3. The actual executable and workflow version are unambiguous. Document staged
   coexistence with the predecessor's installed `grove` without replacing it
   globally as a side effect of setup or testing.
4. Exercise fresh-session discovery in Claude and Codex where available; report
   real behavior separately from file existence and simulated harness calls.
5. Package boundaries separate portable workflow from this repo's Go tests,
   worktree naming, `go run` invocation and development history.

## Preparation and next

Assigned alone (`/grove-work G-040`, interactive) on 2026-09-22; branch
`worktree-G-040` in `.claude/worktrees/G-040`, base main `ccdc92d`, this
record at `sha256:5d171634…` and the plan
[G-080](G-080-portable-bootstrap-plan.md) at `sha256:c1bfd3e6…` when
implementation started. Live sibling installation and migration belong to
G-041, not this assignment.

## Evidence

Commits: `73b2044` plan, `b8f7cf5` active, `7ddbf52` the implementation,
`731674e` and `eb32209` fixes from the first trials, `66fcccc` and
`3113280` the independent review's two rounds, then the evidence commit,
which is the candidate and changes only records. Review record:
[G-082](G-082-portable-bootstrap-review.md), with the independent
findings, their dispositions, and the trial evidence per acceptance item.

Changed behaviour, against each acceptance item:

1. `grove init` sets up a Git checkout's top with `grove.yaml`
   (`schema_version: 3`, `records: grove`, `brief: grove/brief.md`), the
   record root, a placeholder brief that states no intent, and six marked
   entrypoints: `grove-work` and `grove-shape` for Claude Code under
   `.claude/skills/` and for Codex under `.agents/skills/` with their
   policies. A disposable repository was initialized, validated, shaped
   headless by Claude and by Codex, and handed off headless by Claude,
   through those entrypoints and the new CLI only (G-082, items 1 and 4).
2. A rerun keeps the configuration, root, brief and any unmarked file,
   reports `unchanged` or `updated` for marked ones, and on any conflict
   prints every reason, writes nothing, and exits 1 (G-082, item 2).
3. `grove version` prints the module version, the VCS revision when
   stamped, and a digest of the embedded guides; `grove guide work|shape`
   prints the guides the binary carries, so the workflow version is the
   executable version. The README's "Adopt Grove in another repository"
   documents building from a clone or `go install` from a pushed commit,
   coexistence with the predecessor on PATH, the Codex login-shell case, and
   the linked-worktree stamping limit; the predecessor stays installed.
4. Fresh `claude -p` and `codex exec` sessions in the target discovered and
   followed the generated entrypoints; the first Codex run stopped on the
   predecessor answering, as the adapter says to (G-082, items 3 and 4).
5. The guides no longer link into this repository; the printed guides and
   generated files are scanned by a test for `go run`, worktree names and
   this repository's paths (G-082, item 5).

Decisions taken inside the outcome, recorded in G-080: the binary is the one
owner of the workflow (a root package embeds `docs/`), `init` never touches
`AGENTS.md` or `CLAUDE.md`, adapters carry a marker rather than a version,
and distribution is a build from a named commit with no installer or
self-update. The plan's one adjustment: the guides also defer to the loading
entrypoint where a repository has no instructions (G-082 finding 2).

Verification at `3113280`, the last code change, uncached: `gofmt -l .` clean, `go vet ./...` ok,
`go run ./cmd/grove check` → `OK: 80 records`, `go test -count=1 -timeout
120s ./...` ok in every package (`internal/versions` 7.2 s, the known Git
process ceiling). Every link written here and in G-082 resolves in this
checkout.

Limits: interactive typing of `/grove-work` and `$grove-shape` in a target
is unexercised; the headless runs used Claude Opus 5 and codex-cli 0.155.1
in this session's sandbox; the Codex sessions load the predecessor's global
Codex plugin first, which `init` cannot control; a binary built in a linked
worktree names the enclosing checkout's commit. The disposable repositories
and their worktrees are under the session scratchpad and are not part of
this branch.

## Next

In Review. The candidate is the evidence commit named in `candidate`; the
branch tip adds only this status change. To judge it, from a checkout of
`worktree-G-040`:

```sh
go run ./cmd/grove context G-040 --include grove/G-080-portable-bootstrap-plan.md
go run ./cmd/grove show G-082
git diff --stat ccdc92d..HEAD
go test -count=1 -timeout 120s ./...
# demo in a disposable repository, with a binary from a clone of this commit:
git clone -q --no-hardlinks . /tmp/grove-build && (cd /tmp/grove-build && go build -o /tmp/grove-bin/grove ./cmd/grove)
mkdir -p /tmp/grove-target && cd /tmp/grove-target && git init -q && PATH=/tmp/grove-bin:$PATH grove version && PATH=/tmp/grove-bin:$PATH grove init && PATH=/tmp/grove-bin:$PATH grove check
```

Approve: on `main`, `git merge --ff-only worktree-G-040`, then there quote
the verdict in this record and run
`go run ./cmd/grove update G-040 --expect REVISION --set status=done`
(revision from `show G-040 --json` after the edit), and commit both.
Feedback: `--set status=active` on the branch with the feedback here. Then
G-041, which depends on this record, can be prepared.
