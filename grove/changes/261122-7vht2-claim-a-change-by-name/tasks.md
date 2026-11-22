<!-- Each task is a section headed "## T-NN — title". IDs are never
renumbered or reused. -->

## T-01 — Claim a change by name; in flight and overlap for a named claim

Result: `grove claim <id>` records an attempt at a cycle or at a change
in no cycle, refusing what cannot be built, what is in flight and what
overlaps it, while a claim given nothing is paused as long as any change
is in flight.

Files:
- `internal/doc/doc.go`
- `internal/doc/doc_test.go`
- `internal/cli/claim.go`
- `internal/cli/serve.go`
- `internal/cli/cli.go`
- `internal/cli/walkthrough_test.go`

Delegated: private helpers, test names and layout, and the exact wording
of each one-sentence message.

Interfaces: in `internal/doc`, a change's planned files as repo-relative
paths, and whether two sets of planned files overlap; in `internal/cli`,
for the records on `grove/runs` and the target's documents, each change
in flight with its attempt, and what `grove claim` given a piece would
do with it, as one sentence. T-02 uses all of them.

Constraints:
- Built only once the target's `CONSTITUTION.md` holds the amendment of
  pull request #181: rule 11 no longer lists "building two cycles at
  once". If it does not, this task is blocked, saying so.
- A record written before this change reads as it did: its `cycle` and
  `changes` mean what they meant.
- `grove claim` given nothing claims what it claims at 07d319f whenever
  no change is in flight.
- Nothing is read from an attempt's branch.
- `internal/cli/serve.go` changes here only as far as it must to read
  what is in flight where it read the pause; its new fields are T-02's.

Acceptance:
- `grove claim <change-id>`, for an approved change in no cycle that
  waits on nothing, with no change in flight, records an attempt naming
  that change alone and no cycle, creates the attempt's branch and prints
  the change's ID (runs, "A change in no cycle" and "A loose change's
  record").
- An ID that is no buildable piece is refused with its reason: no such
  cycle or change, a change that waits, a member of a cycle not done,
  which is named, an archived change, and the reasons a cycle is refused
  for today (runs, "Refuse a given piece that cannot be built").
- With one piece in flight, a named claim of a piece that overlaps it
  nowhere is recorded; one whose tasks name a file, or a file under a
  folder, the other's name, or that has a change spec for a capability
  the other has one for, is refused naming the shared path and the
  attempt; and the same piece again is refused as in flight (runs, "Start
  one attempt" and "Refuse what is in flight or overlaps it").
- A change stale against its `base`, or with an open question, is
  claimed (runs, "A change stale against its base" and "A change with an
  open question").
- `grove claim` given nothing is paused while any change is in flight,
  naming the newest record in flight, also when a later record's change
  is already archived (runs, "Pause after an attempt").
- `--retry` records an attempt despite a piece in flight or an overlap,
  and still refuses a change that waits (runs, "Retry").
- Two claims at once of pieces that share nothing both record an
  attempt; of the same piece, or of two that overlap, one records and the
  other is refused for what the first recorded (runs, "One claim wins").
- `grove` with no arguments shows `claim` taking a cycle's or a change's
  ID.
- `just check` passes.

## T-02 — The project value and the page

Result: `grove changes --json` and `/api/project` say what `grove claim`
given each piece would do, and which records are in flight, and the page
shows both.

Files:
- `internal/cli/serve.go`
- `internal/cli/page.html`
- `internal/cli/walkthrough_test.go`

Delegated: private helpers, test names and layout, the exact wording of
each verdict, and how the page marks a record in flight.

Constraints:
- A cycle's `runs` is true exactly when it is at 07d319f in a project
  with one attempt at a time: buildable, and no change in flight.
- A change's `ready` means what it means at 07d319f.

Acceptance:
- Each approved change in no cycle not done has `claimable` and
  `verdict`; a member of such a cycle and an archived change have neither
  (changes, "Print the project as data"; serve, "Readiness of each change
  in no cycle").
- A cycle not done has `claimable`, and its verdict names a piece in
  flight it overlaps, or says it is in flight, where today it said
  claimed (serve, "Readiness of each cycle").
- With a piece in flight that shares nothing with buildable cycle K, K's
  `claimable` is true and its `runs` is false (serve, "Another piece in
  flight that shares nothing").
- Each record in `runs` has `in_flight`, true only for the newest record
  naming a change not archived, and a record for a change in no cycle has
  no `cycle` (serve, "Recorded runs").
- `grove changes --json` and `/api/project` are equal but for `read`.
- The page's readiness tab lists each change in no cycle with its
  verdict and counts it, and the in progress tab marks the records in
  flight (serve, "The page").
- `just check` passes.

## T-03 — Docs and walkthrough

Result: the documentation shows claiming a change by name and several
pieces at once, and the walkthrough runs both on the fixture.

Files:
- `docs/running.md`
- `docs/reference.md`
- `docs/coordinating.md`
- `docs/watching.md`
- `docs/concepts.md`
- `docs/guides/run-on-a-timer.md`
- `dev/walkthrough.md`

Acceptance:
- `docs/running.md` shows `grove claim <change-id>` for a change in no
  cycle and what refuses it; says a named claim is refused only for what
  is in flight or overlaps it, with the sentence Grove prints for each,
  and that each claim needs a checkout of its own; says what a change's
  planned files are; and says, under "The pause", "Claim a cycle you
  name", "Two checkouts at once" and "Retry", what holds now: the pause
  is for a claim given nothing and lasts while any change is in flight,
  and `--retry` claims despite it.
- `docs/reference.md` gives `grove claim <id>` for a cycle or a change,
  and `--retry`, as they now behave.
- `docs/coordinating.md`, under "The project as data", shows `claimable`
  and `verdict` on a change in no cycle, `claimable` on a cycle, and
  `in_flight` on a record, and says `runs` and `ready` mean what they
  did.
- `docs/watching.md` says the readiness tab lists each change in no
  cycle and the in progress tab marks attempts in flight.
- `docs/concepts.md` says an attempt is a claim of a cycle or of a change
  in no cycle, what in flight means, and that paused is about a claim
  given nothing.
- `docs/guides/run-on-a-timer.md` says the timer's command behaves as
  before, and waits while a piece claimed by name is in flight.
- `dev/walkthrough.md` has steps, each run by
  `internal/cli/walkthrough_test.go` under T-01 and T-02, that claim a
  change in no cycle by name; claim a second piece that shares nothing
  from another checkout; are refused for a shared file, for a shared
  capability and for the same piece; find a claim given nothing paused;
  and claim with `--retry`.
