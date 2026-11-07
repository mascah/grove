# Tasks: Serve the project's state as a local page

## T-01 — `grove serve` and the project value

Result: `grove serve` listens on loopback, answers `/api/project` with the
value the serve spec describes and `/api/doc` with a document, and the
end-to-end test checks both on the fixture.

Files:
- `internal/cli/serve.go`
- `internal/cli/cli.go`
- `internal/cli/walkthrough_test.go`

Delegated: test layout and helpers in `internal/cli/walkthrough_test.go`;
the JSON field names.

Interfaces: `grove serve [--addr <host:port>]`; `GET /api/project`;
`GET /api/doc?ref=<empty|worktree|commit>&path=<file under the root>`.

Constraints: the verdicts come from `runnable` and `changedSincePlan` as
`run` and `check` have them, exporting nothing; Git runs only through
`internal/gitx`; a branch's work is told by blob against the target's
history; Markdown renders through goldmark with raw HTML dropped; the
total stays within rule 12's 5,000 Go lines.

Acceptance:
- "Start a local server": the default and a given address print their URL
  and serve `/`; a taken address is reported and exits non-zero.
- "The project value", "Readiness of each cycle", "Work in progress by
  branch", "Recorded runs" and "History of the record": the test builds
  the fixture through the walkthrough's step, with a fresh marked cycle, a
  cycle whose member needs a person, a stale cycle, a build branch, a
  draft branch, a squash-merged branch, an origin-only branch, a pushed
  branch, an uncommitted edit and an attempt record, and checks each
  scenario's value.
- "A document rendered": a proposal comes back as HTML, a cycle file as
  text, a branch's file at its tip, a `<script>` is dropped, and a path
  outside the root is refused.

## T-02 — The page

Result: the page at `/` shows readiness, in progress, graph, history and
specs, with the document pane, from the value, within 500 lines of HTML
and script that `just check` counts.

Files:
- `internal/cli/page.html`
- `internal/cli/serve.go`
- `justfile`

Interfaces: `GET /` serves the embedded page.

Constraints: no build step, no framework, no external resource; the page
reads only `/api/project` and `/api/doc`.

Acceptance:
- "The page": every view and the pane behave as its scenarios say on the
  fixture, checked by hand through the walkthrough's step, with no console
  error, in light and dark mode.
- `just check` prints the page's line count and fails above 500.

## T-03 — Documentation

Result: a reader can find `grove serve` and learn what each view shows and
how to read it.

Files:
- `docs/watching.md`
- `mkdocs.yml`
- `docs/reference.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/watching.md` explains each view and the pane in plain words,
  links the serve spec, and sits in the nav after "Running unattended".
- The reference's "Terminal commands" has a row for `grove serve` and
  `--addr`.
- `dev/walkthrough.md` has a step that runs `grove serve` on the fixture
  and says what each view shows, which T-01's test follows.
- The docs site builds.
