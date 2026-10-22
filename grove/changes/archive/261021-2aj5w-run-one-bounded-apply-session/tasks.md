## T-01 — Run one bounded apply session

Result: `grove run` and `grove run --retry` behave as
[specs/runs/spec.md](specs/runs/spec.md) says, and the documentation
specifies them.

Constraints: no code before the constitution amendment, pull request #17,
is merged. Standard library only (rule 16). Within the rule 12 budget.
Messages are one sentence (rule 14).

Acceptance:
- `docs/reference.md` lists `grove run` with `--retry` and `--timeout`.
  A docs page walks through an ordinary run, the pause, an interruption
  and a retry. It says the operator prepares the checkout, Claude Code's
  permission settings and the `gh` login, and shows the `rules.skills.apply`
  lines that make the agent commit, push and open a pull request.
  `dev/walkthrough.md` runs those steps (rule 8).
- On the ledger fixture, in a disposable repository with a remote, with
  Claude Code and those rules: two dependent changes are applied in one
  attempt. The output
  and the agent's pull request can be inspected. `grove run` before the
  merge launches nothing; after the merge it considers new work.
- On the fixture: an attempt is cancelled after a visible edit. No process
  in its process group survives, and the branch, edit and output remain.
  `grove run --retry` refuses while the edit is uncommitted, then starts
  once the operator has dealt with it. `--timeout` is exercised with a short
  limit.
- `internal/cli/walkthrough_test.go` covers the steps that need no real
  agent, with a stand-in `claude` on `PATH`: nothing to apply, a run in
  progress, paused, unclean checkout, a local target that cannot be
  fast-forwarded, failed to start, timeout with a child that ignores the
  request to stop.
- `CLAUDE.md`'s layout table names any new package (rule 15), and
  `just check` passes.

## T-02 — Run it on a launchd timer

Result: a guide in `docs/guides/` shows a launchd job running `grove run`
on an interval in a dedicated checkout, as a trialed recipe (rule 8).

Acceptance:
- The guide gives the job definition, the Claude Code settings an
  unattended session needs, and where to find each attempt's output.
- It was tried on a Mac against a fixture repository with a remote. The
  timer started one attempt, and firings while it ran or while the project
  was paused launched nothing.
