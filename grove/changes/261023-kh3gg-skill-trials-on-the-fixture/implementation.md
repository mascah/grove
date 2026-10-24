# Implementation: Skill trials on the fixture

<!-- Progress on change 261023-kh3gg, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Write the skill trials

Status: in progress

The revision of 2026-10-23 has a script run the trials, so `dev/trials.md`
and `CONTRIBUTING.md` must say so; the evidence below is for the hand-run
version.

- `dev/trials.md`: the common setup from `just fixture`, the five counts
  with how each is taken, a `jq … | wc -w` line per harness for the words
  read (formats checked on a local Claude Code transcript and a Codex
  rollout file), the five cases with setup, what to type (`/grove-…`,
  `$grove-…` in Codex) and a good run, and the results table.
- `CLAUDE.md`: layout row for `dev/trials.md`.
- `CONTRIBUTING.md`, "Tests": run the trials when a pull request changes
  instruction text in `internal/instruct`, counts in its description.
- Case (a)'s restore, `git checkout ':/^Set up Grove' -- ledger`, run on a
  scratch fixture: `python3 -m ledger edit 3 --note x` printed
  `error: category 'misc' is not allowed`, exit 2.

### T-02 — Record the baseline

Status: not started

Was blocked on who runs the ten trials by hand; the revision of
2026-10-23 answers it: `just trials` (T-03).

Found while preparing it: in a git worktree, `just build` gives a binary
whose `grove version` prints `grove dev`, with no build info from Git;
the main checkout's build prints `grove 0.4.1-0.20261024005345-8454724ba54d`.
So the commit `dev/trials.md` asks for may have to come from `git rev-parse
HEAD` when the trials run from a worktree.

### T-03 — A script that runs the trials

Status: not started

## Scenarios

None: the change has no change spec.

## Reconciliation

None: the change has no change spec.
