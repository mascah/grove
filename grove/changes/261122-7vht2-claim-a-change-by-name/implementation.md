# Implementation: Claim a change by name, and more than one piece of work at once

<!-- Progress on change 261122-7vht2, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Claim a change by name; in flight and overlap for a named claim

Status: done

`internal/doc`: `Change.Planned(root)` (tasks' `Files:` paths plus
`<root>/specs/<capability>/spec.md` for each capability the change has a
change spec for) and `Overlap(a, b)` (a shared path, either side naming a
folder above the other). `internal/cli/claim.go`: `pausing` is gone;
`inFlight` maps each change a record names and the target does not hold
archived to its newest record, `newest` gives the newest record in flight
(the pause for a claim given nothing), `piece` looks an ID up as a cycle
or a loose change and gives each refusal of "Refuse a given piece that
cannot be built", and `busy` gives "in flight" and "overlaps"; `--retry`
skips only `busy` and the pause. A record for a loose change leaves out
`cycle`. `serve.go` reads the pause through `newest(recs, inFlight(...))`
and nothing more. `cli.go`: the usage names `<change-id>`.

Checks: `internal/doc/doc_test.go` `TestOverlap`; `internal/cli/walkthrough_test.go`
`TestClaimByName` (new: same change, folder overlap and nothing shared
at the same moment from two clones; shared file; shared capability; a
change that waits, also with `--retry`; a member of a cycle; an archived
change; no such ID; the pause naming the newest record in flight, also
after a later record's change is archived; `--retry` after the ID despite
an overlap; a loose change's record has no `cycle`; the claimed change is
stale against its base and has an open question), `TestClaim` (updated:
"no cycle or change", and a named cycle now refused for sharing
record-editing's spec with the cycle in flight instead of the pause), and
the usage line. `just check` passed: 4,674 Go, 2,531 test lines.

Check: the target's `CONSTITUTION.md` at baea991 holds pull request #181's amendment; rule 11 no longer lists building two cycles at once, so the task is not blocked.
Check: `internal/cli/walkthrough_test.go` changed on main since base 07d319f (change 261122-6ndwb added implement's rewritten-requirement test); the claim tests are built on the file as it is now, with nothing of that diff touched.
Check: `--retry` is accepted before or after the ID, since the runs spec's scenarios write `grove claim A --retry` and Go's flag package stops at the first argument; claim's arguments are read without it.
Check: an overlap on a capability names the maintained spec's path, such as `grove/specs/record-editing/spec.md`, as the shared path, rather than the words "the maintained spec of record-editing"; the message wording was delegated.
Check: a change that is both a member of a cycle not done and waiting is refused as a member of the cycle, naming the cycle to claim, since claiming the cycle is what the operator can do.

### T-02 — The project value and the page

Status: done

`internal/cli/serve.go`: the records are read before the changes, so each
loose change gets `claimable` and `verdict` from the same `piece` and
`busy` that `grove claim` runs; a cycle not done gets `claimable` from
`runnable` then `busy`, its verdict naming the attempt it is in flight in
or the path it shares, and `runs` is `claimable` with no change in
flight; each record gets `in_flight`. `ready` is unchanged.
`internal/cli/page.html`: the readiness count adds the loose changes, each
loose change shows as a card with its verdict over its row, a cycle's
colour follows `claimable`, and the runs table names a loose change's ID
where it has no cycle and marks a record in flight with a chip.

Checks: `TestServe` (updated: loose change `e` claimable with its
verdict, member `a` and archived `f` with neither; K in flight in
attempt 261105-r0000, `claimable` and `runs` false; records `in_flight`;
`changes --json` equal to `/api/project` but for `read`), and
`TestClaimByName`'s project value (cycle K, sharing nothing with four
changes in flight, `claimable` true and `runs` false; `e`'s verdict names
the shared `tests/test_cli.py`; `a` in flight; `g` claimable; member `m`
with no verdict; each change's newest record in flight but dropped `f`'s,
and the older of `a`'s two records not; no record has `cycle`).
`just check` passed: 4,696 Go, 2,559 test and 411 page lines. One earlier
`just check` run failed in `internal/cli` without my seeing which test;
six runs since have passed.

Check: the page was not seen in a browser, since the browser connection was down in this session; `bin/grove serve` on this repository answered `/api/project` with 261122-7vht2 claimable and its verdict, and the page's script was read, not run. The readiness count and the in-flight mark want a look.
Check: a cycle's verdict when it is claimable but a claim given no piece is paused reads "grove claim <K> would claim it; grove claim given no piece is paused after attempt <attempt> of <cycle or change>", replacing the old "grove claim --retry would claim it; the project is paused after ..."; a cycle in flight now reads "<K> is in flight in attempt ..." where it read "claimed by attempt ...".
Check: one `just check` failure in `internal/cli` went unseen, and later runs passed; `TestServe` binds the fixed address 127.0.0.1:7070, which another session on this machine may have held, but that is a guess.

### T-03 — Docs and walkthrough

Status: not started

## Scenarios

## Reconciliation
