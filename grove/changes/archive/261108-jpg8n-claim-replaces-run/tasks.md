## T-01 — `grove claim` selects, records and prepares

Result: `grove claim [cycle-id] [--retry]` replaces `grove run`: it
fetches the target and `grove/runs`, refuses while the record on
`grove/runs` says the project is paused, selects the cycle, records the
attempt as a commit on `grove/runs` pushed with a compare-and-swap,
creates `grove/<attempt>` from the fetched target and prints the cycle's
ID; Grove starts no process but Git, and `grove run` is no command.

Files:
- `internal/cli/claim.go`
- `internal/cli/run.go`
- `internal/cli/cli.go`
- `internal/cli/check.go`
- `internal/gitx/gitx.go`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers

Interfaces: `grove claim`; the branch `grove/runs`; the record's fields
`id`, `started`, `cycle`, `changes`, `commit`, `host`, `checkout` and
`branch`.

Constraints: no `--timeout`; `agents` need not list Claude Code; nothing
is written under `.git/grove/`; dependencies stay as rule 16 lists them;
the Go and test line budgets of rule 12 hold.

Acceptance:
- Every scenario of the runs change spec, and the `grove claim` scenarios
  of the cycles change spec, hold on the fixture.
- `TestClaim` follows the walkthrough on disposable repositories with a
  bare `origin` and two clones, and shows one of two claims made at once
  winning while the other finds the project paused.
- `grove run` is reported as an unknown command, and `grove check`'s
  message for a marked stale cycle says `grove claim` skips it.
- `just check` passes.

## T-02 — `serve` reads the records and names `grove claim`

Result: the page's readiness verdicts say `grove claim`, name a claimed
cycle with its host, start time and branch, and its runs table lists the
records on `origin/grove/runs`, or the local `grove/runs`, with each
attempt's host and its branch's latest commit time.

Files:
- `internal/cli/serve.go`
- `internal/cli/page.html`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and private helpers

Acceptance:
- Every scenario of the serve change spec holds on the fixture.
- `TestServe` covers a record on `origin/grove/runs` from another host
  with no local attempt.
- The page line budget of rule 12 holds.

## T-03 — Instruction text names `grove claim`

Result: the plan instructions' verdict, the unattended implement
instructions and the project index say `grove claim` and never
`grove run`.

Files:
- `internal/instruct/requires/plan.md`
- `internal/instruct/requires/implement-unattended.md`
- `internal/instruct/instruct.go`

Acceptance:
- The scenarios of the config, implement and propose change specs hold
  on the fixture.
- `grove instructions plan "x"` and
  `grove instructions implement <change-id> --unattended` on the fixture
  print no `grove run`.
- The word budget of rule 12 holds.

## T-04 — Documentation and the walkthrough

Result: the documentation explains `grove claim`, the caller's one-line
invocation with `timeout`, the branch `grove/runs` and the pause from any
host, and names `grove run` nowhere; the timer guide runs that line from
launchd through `sh -c`; the walkthrough runs the claim steps on the
fixture with a bare `origin` and two clones.

Files:
- `README.md`
- `docs/index.md`
- `docs/getting-started.md`
- `docs/concepts.md`
- `docs/coordinating.md`
- `docs/planning.md`
- `docs/running.md`
- `docs/watching.md`
- `docs/reference.md`
- `docs/guides/run-on-a-timer.md`
- `dev/walkthrough.md`

Acceptance:
- No file under `docs/`, `README.md` or `dev/walkthrough.md` names
  `grove run`, `--timeout` or `.git/grove/`.
- `docs/running.md` shows the invocation
  `K=$(grove claim) && timeout 2h claude -p "/grove-implement $K --unattended"`
  and says what `claim` prints and when it exits non-zero.
- Each documented step is run on the fixture and recorded in
  `dev/walkthrough.md`.
