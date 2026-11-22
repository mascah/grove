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
`just check` passed: 4,696 Go, 2,559 test and 411 page lines (2,562 test
lines after T-03's scenario checks). One earlier
`just check` run failed in `internal/cli` without my seeing which test;
six runs since have passed.

Check: the page was not seen in a browser, since the browser connection was down in this session; `bin/grove serve` on this repository answered `/api/project` with 261122-7vht2 claimable and its verdict, and the page's script was read, not run. The readiness count and the in-flight mark want a look.
Check: a cycle's verdict when it is claimable but a claim given no piece is paused reads "grove claim <K> would claim it; grove claim given no piece is paused after attempt <attempt> of <cycle or change>", replacing the old "grove claim --retry would claim it; the project is paused after ..."; a cycle in flight now reads "<K> is in flight in attempt ..." where it read "claimed by attempt ...".
Check: one `just check` failure in `internal/cli` went unseen, and later runs passed; `TestServe` binds the fixed address 127.0.0.1:7070, which another session on this machine may have held, but that is a guess.

### T-03 — Docs and walkthrough

Status: done

`docs/running.md`: the intro covers a cycle or a change in no cycle and in
flight; a new "Claim a change you name" gives `grove claim <change-id>`,
each refusal with the sentence Grove prints, in flight, overlap, planned
files and that the check only reduces collisions; "Two checkouts at once"
says each claim needs its own checkout and how a losing claim decides
again; "Claim a cycle you name", "The pause" and "Retry" say the pause is
for a claim given nothing, lasts while any change is in flight, and that
`--retry` claims despite it and despite an overlap; "The record" says a
loose change's record has no `cycle`. `docs/reference.md`: `grove claim
<id>` and `grove claim --retry [<id>]`. `docs/coordinating.md`: the
example gains a loose change with `claimable` and `verdict`, a cycle with
`claimable` and `runs`, and a record with `in_flight`, and says `runs`
and `ready` mean what they did. `docs/watching.md`: readiness lists and
counts each change in no cycle with its verdict, and the attempts mark
those in flight. `docs/concepts.md`: attempt, in flight and paused.
`docs/guides/run-on-a-timer.md`: the timer's line behaves as before and
waits while a piece claimed by name is in flight. `dev/walkthrough.md`:
running step 6 now expects a named claim refused for overlap or claimed,
not paused; a new running step 9 claims changes by name from two clones,
some at once, is refused for a shared file, folder, capability and the
same change, finds a claim given nothing paused, and claims with
`--retry`; Watching steps 2 and 3 expect the count of 4, E's card, and K
in flight.

Checks: running step 9 is `TestClaimByName`, and the changed lines of
running step 6 and Watching steps 2 and 3 are `TestClaim` and `TestServe`;
`just check` passed.

Also changed: `README.md` — it said `grove claim` claims a cycle, which this change makes incomplete.
Check: the walkthrough's steps were not run by hand on the fixture in Claude Code or Codex; their Grove commands are what the tests above run, and the steps that start a session were not changed.
Check: running step 6 cannot say in advance whether the export cycle shares a planned file with the budgets change, since a real `/grove-propose` session picks both changes' files and capabilities, so the step gives both outcomes.
Check: `docs/planning.md` still says to send one change on its own by planning a cycle of one; that remains true, for a timer, and was left as it is.

## Scenarios

- runs "A needs_human left over": internal/cli/walkthrough_test.go TestClaim (cycle N with a needs_human claimed), passed
- runs "A stale plan": TestReplan and TestMovedFiles (`runnable` refuses the stale cycle), passed
- runs "A cycle it follows not done": TestFollows, passed
- runs "No such cycle or change": TestClaim and TestClaimByName (261099-zzzzz), passed
- runs "A change that waits": TestClaimByName (W), passed
- runs "Its dependency finished only on a branch": unverified; no test archives A on its attempt's branch, though the refusal reads only the target's documents, as TestClaimByName's W shows
- runs "A member of a cycle": TestClaimByName (M in K), passed
- runs "An archived change": TestClaimByName (Z), passed
- runs "The same change again" (in flight): TestClaimByName (A from h2), passed
- runs "A file both name": TestClaimByName (C shares ledger/cli.py with A), passed
- runs "A folder one names": TestClaimByName (E's tests/test_cli.py under B's tests/) and internal/doc/doc_test.go TestOverlap, passed
- runs "One capability, different requirements": TestClaimByName (G shares record-editing's spec with F) and TestOverlap, passed; F and G hold the same requirement there, which the check does not read
- runs "A cycle that overlaps a change in flight": TestClaim (cycle M's change shares record-editing's spec with C, in flight in L's attempt), passed
- runs "A cycle with a member in flight": TestServe (K in flight), passed, with neither member archived; the case with A archived is not tested
- runs "What an attempt changed on its branch": unverified; no test revises a change on its attempt's branch, though planned files are read only from the target collection
- runs "The attempt it would overlap is merged": TestClaimByName (G claimed once F is archived, dropped), passed
- runs "Local target has commits of its own": TestClaim, passed
- runs "A marked cycle of two changes, one waiting on the other": TestClaim, passed
- runs "Composed with a harness": unverified; it needs a harness session
- runs "Two marked cycles": TestClaim (K, created first, claimed), passed
- runs "The older marked cycle waits on the newer": TestFollows, passed
- runs "A given cycle not marked": TestCoordinating (`grove claim 261022-p0000`), passed
- runs "A change in no cycle": TestClaimByName (A), passed
- runs "A change composed with a harness": unverified; it needs a harness session
- runs "A given cycle while another is in flight": TestClaimByName (K claimed with five pieces in flight), passed
- runs "A second change that shares nothing": TestClaimByName (B from h1 while A is in flight), passed
- runs "A change stale against its base": TestClaimByName (A, reported by check), passed
- runs "A change with an open question": TestClaimByName (A), passed
- runs "A checkout with no remote": TestClaim, passed
- runs "The branch merged locally": TestClaim, passed
- runs "The first attempt of a project": TestClaim, passed
- runs "A later attempt": TestClaim (`--retry` on M keeps L's record), passed
- runs "Read by hand": TestClaim (`git log -1 --format=%B`), passed
- runs "A loose change's record": TestClaimByName, passed
- runs "Nothing in the working tree": TestClaim, passed
- runs "Pull request not merged yet": TestClaim, passed
- runs "Another host": TestClaim, passed
- runs "Another cycle marked meanwhile": TestClaim, passed
- runs "An older attempt still open": TestClaimByName (paused naming the earlier attempt once F is archived), passed
- runs "A change claimed by name": TestClaimByName (paused naming F's attempt), passed
- runs "Cycle merged": TestClaim, passed
- runs "Blocked change dropped": TestClaim, passed
- runs "A blocked change revised": TestReplan (the plan stays fresh after one squash of the revision and re-plan) with TestClaim's `--retry`, passed; no single test runs both
- runs "The earlier attempt kept": TestClaim, passed
- runs "The same change again" (retry): TestClaimByName (`--retry` on A, the older record no longer in flight), passed
- runs "Despite an overlap": TestClaimByName (`grove claim C --retry`, A still in flight), passed
- runs "A change that waits" (retry): TestClaimByName (`--retry` on W), passed
- runs "Two hosts at once": TestClaim, passed
- runs "Two changes that share nothing": TestClaimByName (H and F at once), passed
- runs "The same change from two hosts": TestClaimByName (A at once), passed
- runs "Two changes that overlap": TestClaimByName (B and E at once), passed
- runs "A push refused for permissions": TestClaim, passed
- changes "One approved, one archived": TestCoordinating, passed
- changes "A dropped change": unverified; this change does not touch it and no test asserts the state `dropped`
- changes "No implementation yet": unverified; this change does not touch it and no test is named for it
- changes "Cycles and attempts": TestServe (K in flight, member A with no `claimable`, the record `in_flight`), passed
- changes "A change in no cycle": TestServe (E), passed
- changes "A change in flight": TestClaimByName (A's verdict), passed; its `ready` is not asserted
- changes "A branch with work": TestServe, passed
- changes "The same value as the page": TestServe, passed
- serve "Nothing in the way": TestServe (E), passed
- serve "A change that waits": TestClaimByName (W's verdict), passed
- serve "A change in flight": TestClaimByName (A's verdict), passed
- serve "A change that overlaps one in flight": TestClaimByName (E's verdict), passed
- serve "A member of a cycle": TestServe (A) and TestClaimByName (M), passed
- serve "An approved change in no cycle": TestServe (E), passed
- serve "The same value as the command": TestServe, passed
- serve "A marked cycle that is fresh": TestServe (K), passed
- serve "A member waits outside the cycle": TestServe (L), passed
- serve "A member changed since the plan": TestServe (M), passed
- serve "A cycle in flight": TestServe (K), passed
- serve "Another piece in flight that shares nothing": TestClaimByName (K, not marked, `claimable` true, `runs` false), passed
- serve "Another piece in flight that shares a file": unverified in the value; the verdict is `busy`'s sentence, which TestClaim shows for a claim of such a cycle
- serve "An attempt from another host": TestServe, passed
- serve "Two attempts, one merged": TestClaimByName (F's record not in flight, no `cycle`), passed
- serve "An attempt retried": TestClaimByName (A's older record), passed
- serve "No record anywhere": TestServe, passed
- serve "Readiness at a glance": unverified; the page was not seen in a browser
- serve "A change in no cycle" (page): unverified; the page was not seen in a browser
- serve "Attempts in flight": unverified; the page was not seen in a browser
- serve "A refresh while reading": unverified; the page was not seen in a browser

## Reconciliation
