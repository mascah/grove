# Implementation: Serve the project's state as a local page

<!-- Progress on change 261107-9z9tg, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — `grove serve` and the project value

Status: done

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
Answer: "Keep both (Recommended)"

Also changed: `internal/cli/check.go` — `check`'s problem-gathering moved into `diagnose`, which the value uses for the problems `grove check` would report on the target.
Also changed: `internal/cli/run.go` — `latest` now reads through `attempts`, which the value uses for the recorded runs, newest first.
Also changed: `internal/cli/page.html` — a one-line placeholder so the page embeds; T-02 writes the page.
Also changed: `go.mod` — goldmark moves from an indirect to a direct dependency.

### T-02 — The page

Status: done

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
cycles, 1 row (this branch), 12 specs. `serve.go` now also gives each
change its `files` and each cycle its `created`, which the pane and the
timeline read.

Check: open the page on a project you know and judge how each view reads, in light and dark mode; how it looks is a person's call, and the fixture's run shows only that it works.
Answer: "Looks good"

### T-03 — Documentation

Status: done

`docs/watching.md`, "Watching the project", explains starting the server
and each view and the pane in plain words and links the serve spec;
`mkdocs.yml` puts it after "Running unattended". `docs/reference.md`'s
"Terminal commands" has rows for `grove serve` and `grove serve --addr`.
`dev/walkthrough.md` has a "Watching the project" section whose four steps
build the fixture's states and say what each view shows; `TestServe`
follows it. `uvx zensical==0.0.71 build`, as the docs workflow runs it: no
issues.

Also changed: `docs/index.md` — its list of pages gains "Watching the project".
Also changed: `internal/cli/walkthrough_test.go` — `TestServe` checks the walkthrough step's refusal of an address that is not loopback, and of a request for another host.

## Scenarios

- serve "The default address": `internal/cli/walkthrough_test.go` TestServe, passed
- serve "An address given": TestServe, with `--addr 127.0.0.1:0` rather than 8080 so the test takes a free port, passed
- serve "The address is taken": TestServe, passed
- serve "An approved change in no cycle": TestServe, passed
- serve "A marked cycle that is fresh": TestServe, passed
- serve "A member needs a person": TestServe, passed
- serve "A member changed since the plan": TestServe, passed
- serve "A branch squash-merged and overtaken": TestServe, passed
- serve "A build in progress": TestServe, passed
- serve "A draft": TestServe, passed
- serve "The working tree": TestServe, passed
- serve "A teammate's branch": TestServe, passed
- serve "A pushed branch": TestServe, passed
- serve "A timed-out attempt": TestServe, passed
- serve "A change's life": TestServe, passed
- serve "A commit outside the record": TestServe, passed
- serve "A proposal at the target": TestServe, passed
- serve "A cycle file": TestServe, passed
- serve "A document on a branch": TestServe, passed
- serve "Raw HTML in a document": TestServe, passed
- serve "A path outside the root": TestServe, passed
- serve "Readiness at a glance": unverified by a test, since rule 17 keeps browser harnesses out; driven by hand in headless Chrome on the fixture, where K read "grove run would build it" and L named C as needing a person (the count was 3, the fixture also holding stale cycle M); how it reads is T-02's checkpoint
- serve "A refresh while reading": unverified by a test, as above; in headless Chrome the pane kept the document and both columns their scroll positions across a refresh

## Reconciliation

<!-- One line per change spec, once applied to its maintained spec:

    - record-editing: applied to specs/record-editing/spec.md
-->
