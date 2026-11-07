# Implementation: Serve the project's state as a local page

<!-- Progress on change 261107-9z9tg, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — `grove serve` and the project value

Status: to check

`internal/cli/serve.go` listens on `--addr` (default `127.0.0.1:7070`),
prints the URL and answers `/`, `/api/project` and `/api/doc`. The value's
cycle verdicts come from `runnable` and `changedSincePlan` as they are,
members in build order from `buildable`; branches are told by blob against
one `git rev-list --objects` of the target under the root, each tip's scan
kept per target commit; history is one first-parent `git log
--name-status`; documents render through goldmark, which drops raw HTML.
`internal/cli/cli.go` parses `serve [--addr]` and lists it in the usage.
`TestServe` in `internal/cli/walkthrough_test.go` builds the fixture with
every state the spec names and checks the value and documents: passed.
`just check`: passed, 4,535 Go lines and 1,974 test lines. On Grove's own
repository: 0 branch rows, first request 0.9 s, then 0.1 s.

Check: two refusals the spec does not name, built from rule 7's "on loopback": `--addr` with a host that is not loopback is refused ("--addr must be a loopback address and port, such as 127.0.0.1:7070"), and a request whose Host header is not a loopback name gets 403, so a site the browser visits cannot read the project through a DNS name it points at 127.0.0.1. Keep both, or drop either?

Also changed: `internal/cli/check.go` — `check`'s problem-gathering moved into `diagnose`, which the value uses for the problems `grove check` would report on the target.
Also changed: `internal/cli/run.go` — `latest` now reads through `attempts`, which the value uses for the recorded runs, newest first.
Also changed: `internal/cli/page.html` — a one-line placeholder so the page embeds; T-02 writes the page.
Also changed: `go.mod` — goldmark moves from an indirect to a direct dependency.

### T-02 — The page

Status: to check

`internal/cli/page.html`, 408 lines of HTML and script with no external
resource, built from the prototype's page on the product's field names:
readiness, in progress, graph, history and specs tabs with counts, a
document pane rendering Markdown and showing YAML, relative links opening
in the pane, and a refresh every 10 seconds that a checkbox stops.
`justfile`: `just check` prints the page's line count and fails above 500.
Driven through headless Chrome (Playwright) on the ledger fixture built
with the states `TestServe` builds, in light and dark mode: no console
error; readiness counts 3 open cycles, with the marked cycle "grove run
would build it", the export cycle naming its member that needs a person
and the budgets cycle stale; in progress counts 4 rows; a proposal opened
in the pane, its relative link to the spec opened in the pane, and after a
refresh the pane held the same document at the same scroll position, as
did the main column; the draft's proposal showed no `<script>`; unchecking
the box stopped the refresh. On Grove's own repository: no error, 0 open
cycles, 1 row (this branch), 12 specs.

Check: open the page on a project you know and judge how each view reads, in light and dark mode; how it looks is a person's call, and the fixture's run shows only that it works.

Also changed: `internal/cli/serve.go` — the value gives each change its `files` and each cycle its `created`, which the pane and the timeline need.

### T-03 — Documentation

Status: not started

## Scenarios

<!-- Before finishing, one line per scenario of the change specs: the check
that shows it and its result, or `unverified` and why:

    - record-editing "Update the note on a legacy record": tests/test_cli.py test_edit_legacy, passed
-->

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
