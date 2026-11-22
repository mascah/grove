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

Status: not started

### T-03 — Docs and walkthrough

Status: not started

## Scenarios

## Reconciliation
