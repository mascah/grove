# Implementation: Skill trials on the fixture

<!-- Progress on change 261023-kh3gg, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Write the skill trials

Status: done

Rebuilt for the revision of 2026-10-23: `dev/trials.md` says how `just
trials` runs each case (the common setup, the headless sessions, the
scripted answers and the waiting check), which counts are computed and
which judged, each case's setup as the script writes it, its runs, its
owner's answers and a good run, and the table it prints. `CLAUDE.md`'s
row names `dev/trials/` too, and `CONTRIBUTING.md`, "Tests", says to run
`just trials` and put its table in the pull request's description. The
judge reads "The counts" and the case's section of `dev/trials.md`
(`run.sh`, `judge`), checked by the (b) runs above.

Evidence of the hand-run version, before the revision:

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

Status: done

- `dev/trials/run.sh`, with `dev/trials/a.sh` to `e.sh` (each case's
  setup, runs and answers), and `just trials [case ...]` in `justfile`.
  Bash, jq and python3; no Go, so rule 12's count is unchanged.
- `SETUP_ONLY=1 HARNESSES=claude dev/trials/run.sh`: every setup's
  `grove check` printed `No problems in grove/.`, except (d), whose only
  line is T-01's `ledger/report.py` changed since base, as it means to;
  (a)'s `python3 -m unittest` failed and its `grove specs` listed
  `record-editing`; (e) left `ledger/cli.py`, `tests/test_cli.py` and
  `implementation.md` uncommitted on `delete-impl`.
- `dev/trials/run.sh b` in both harnesses: each folder held `repo/`, the
  turns, `transcript.txt`, `waits.txt`, `outside.txt` and a `judge.json`
  with a reason per count, and the table printed both rows.
- Found while trying it, and fixed in the script: Codex runs commands in
  a login shell whose profile put the everyday `grove` first on PATH
  (`-c allow_login_shell=false`); Codex's sandbox cannot write `.git`, so
  the script makes the branches; `codex exec` and `claude -p` read stdin
  when it is open, which stalled the loop; a round that ends "Yes to all,
  or tell me which numbers to change." has no question mark, so a small
  model decides whether the agent waits.

## Scenarios

None: the change has no change spec.

## Reconciliation

None: the change has no change spec.
