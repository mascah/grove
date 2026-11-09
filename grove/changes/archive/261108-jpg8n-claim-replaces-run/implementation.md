# Implementation: grove claim replaces grove run, with attempt records on a branch

<!-- Progress on change 261108-jpg8n, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — `grove claim` selects, records and prepares

Status: done

`internal/cli/run.go` became `internal/cli/claim.go`. `claim` fetches the
target and `grove/runs` from `origin`, using a pattern refspec so that a
project with no record yet fetches without error. It refuses an unclean
checkout, a local target that cannot be fast-forwarded and, without
`--retry`, a latest record whose cycle has a member not archived, naming
that attempt, its host, start time and branch. It selects the given or
oldest eligible cycle as `run` did. It writes the record as a commit with
the empty tree (`hash-object -t tree -w /dev/null`, then `commit-tree`) and
pushes it with `--force-with-lease=refs/heads/grove/runs:<tip just
fetched>`, or, with no `origin`, runs `update-ref` against the tip just
read. When the update is refused and `grove/runs` has moved, it decides
again as a fresh claim; otherwise it says the record was refused,
quoting Git (the wording was set in T-04). Then it creates
`grove/<attempt>` from the fetched target, fast-forwards the local target
and prints the cycle ID alone. With nothing to claim it prints nothing on
stdout and exits 1. The lock, `--timeout`, the Claude Code launcher, the
relay, the process-group stop and `.git/grove/` are gone, and so is the
requirement that `agents` list Claude Code. `cli.go` has
`claim [--retry] [<cycle-id>]`, and `run` falls through to the usage with
status 2. `check.go`'s message for a marked stale cycle says `grove claim`
skips it. `internal/gitx/gitx.go` needed no change.

`TestRun` became `TestClaim` in `internal/cli/walkthrough_test.go`, with a
bare `origin` and two clones, h1 and h2. For two claims at once, h2's
`pre-push` hook runs h1's claim (the test binary run as grove through
`TestMain`) while h2's push is under way. h1's record wins, h2's lease is
refused, and h2 reads again, finds the project paused and creates no
branch. The test also covers:

- `grove run` returning the usage's status 2;
- nothing to claim: exit 1, empty stdout;
- `agents` listing only Codex;
- the local-target and unclean-checkout refusals, the latter with
  `--retry`, the edit kept;
- a push refused by `origin`'s `pre-receive` hook: nothing recorded, no
  remote-tracking ref, no branch;
- the record's fields and empty tree, and no `.git/grove/`;
- the pause on the same host, and its end after one member archived as
  finished and one as dropped;
- a given cycle refused (no such cycle, a member needs a person, paused)
  and claimed with `--retry`, keeping the earlier record and its branch;
- no `origin`: the record on the local `grove/runs`, and the pause ended by
  merging the attempt's branch.

TestFollows, TestReplan and TestMovedFiles name `grove claim`.

`just check` passes: 4303 Go lines, 1941 test lines, 408 page lines, 8614
instruction words. On a fresh fixture, `bin/grove run` exits 2 with the
usage, and `bin/grove claim` prints nothing on stdout, says "nothing to
claim on main at 5334772" and exits 1.

Also changed: `internal/cli/serve.go` — it read `.git/grove/`, which is gone; it reads the local `grove/runs` with `claim`'s `records` until T-02 has it prefer `origin/grove/runs`.
Also changed: `internal/doc/doc.go` — check's `needs_human` message and a comment said `grove run` skips the cycle, which is false once `run` is no command.

### T-02 — `serve` reads the records and names `grove claim`

Status: done

`internal/cli/serve.go` reads the records on `origin/grove/runs`, or on the
local `grove/runs` when there is no `origin/grove/runs`. Each run on the
page is the record plus `moved`: the latest commit time of its branch,
local or on `origin`. A cycle's verdict, after the existing refusals from
`runnable`, is one of these:

- "claimed by attempt X, which H started at T on branch B", when the latest
  record names it and a member is still open;
- otherwise `grove claim would claim it`, with the cycle's ID when it is
  not marked and "; it is not marked unattended";
- while the project is paused after another cycle's attempt, `--retry` in
  the command and "; the project is paused after attempt X of cycle Y",
  with `runs` false.

The pause test is now one helper, `pausing` in `claim.go`, used by both
`claim` and `serve`. `serve`'s branches no longer list `grove/runs`, local
or on `origin`. Without that exclusion, a `grove/runs` branch broke
`/api/project`, because the record branch shares no history with the
target: `git diff: ... no merge base`. `internal/cli/page.html` heads
readiness "What grove claim would do". Its runs table, "Attempts recorded on
grove/runs", shows attempt, cycle, started, host, branch and last commit.

`TestServe` checks the following:

- K's verdict is "grove claim would claim it" with no runs.
- After a record from host `h2.example` is pushed to `origin/grove/runs`,
  with its branch pushed to `origin` and no local attempt, K is claimed,
  naming the host, time and branch. The run carries its host, cycle,
  branch, start and the branch's commit time.
- With a local copy of `grove/runs` (not listed as a branch) and a second
  record naming L, K says `grove claim --retry would claim it; the project
  is paused after attempt 261106-r0000 of cycle L`, and the runs are newest
  first.

`just check` passes: 4338 Go lines, 1977 test lines, 408 page lines, 8614
instruction words. How the runs table looks in a browser was not looked
at. Its cells are the ones the old table had, apart from the host and last
commit columns.

Also changed: `internal/cli/claim.go` — the pause test moved into `pausing`, so `claim` and `serve` share it.

### T-03 — Instruction text names `grove claim`

Status: done

In `internal/instruct/requires/plan.md`, every `grove run` became
`grove claim`. The finish now says whether `grove claim <cycle-id>` would
claim the cycle. A marked cycle goes to a `grove claim` given no cycle,
and the needs-person obstacle says `grove claim` refuses the cycle. A
member waiting outside the cycle, or a cycle it depends on that is not
done, means `grove claim` will not claim the cycle until that is settled.
"An unattended run will not build" also became "`grove claim` will not
claim". In `internal/instruct/requires/implement-unattended.md`, "only
`grove claim` reads them". The project index in
`internal/instruct/instruct.go` says "`unattended: true` hands one to
`grove claim`". No instruction text names `grove run`.

Tests in `internal/cli/walkthrough_test.go`:

- `TestCycles` wants the new wording and the finish's order, implement,
  then `grove claim <cycle-id>`, then the mark. It also checks that the
  plan instructions name no `grove run`.
- `TestImplementInstructions` checks that the unattended instructions say
  "only `grove claim` reads them" and name no `grove run`. The cycle's
  unattended instructions carry the index line naming `grove claim` and no
  `grove run`.

The config scenario (the project context in the unattended instructions)
is checked by the existing `TestImplementInstructions` case. The propose
scenarios' instruction text names neither command and was unchanged.

`just check` passes: 4338 Go lines, 1977 test lines, 408 page lines, 8613
instruction words. On a fixture, `bin/grove instructions plan "x"` prints
`grove claim` 8 times and `grove run` none. `just trials`
(`CONTRIBUTING.md`, for a pull request changing instruction text) was not
run. It drives Claude Code and Codex headless, so whoever opens the pull
request runs it.

Also changed: `internal/cli/walkthrough_test.go` — its plan and unattended-implement checks matched the `grove run` wording.

### T-04 — Documentation and the walkthrough

Status: done

- `docs/running.md` is rewritten around `grove claim`. It opens with
  `K=$(grove claim) && timeout 2h claude -p "/grove-implement $K --unattended"`
  and says what `claim` prints: the cycle's ID alone on standard output,
  or nothing with a non-zero exit when there is nothing to claim, the
  project is paused or a refusal applies. It has sections on the record
  on `grove/runs` (read with `git log -1 --format=%B origin/grove/runs`),
  two checkouts at once, a cycle you name, no remote, the pause from any
  host, a stopped session (`timeout`, Ctrl-C; Grove does not know) and
  `--retry`.
- `docs/guides/run-on-a-timer.md` runs the same line through `/bin/sh -c`
  from launchd. It notes the `&amp;&amp;` escaping, `timeout` from
  Homebrew's `coreutils`, and that exit code 124 is `timeout`'s.
- `README.md`, `docs/index.md`, `docs/getting-started.md`,
  `docs/concepts.md` (including the Attempt and Paused definitions),
  `docs/coordinating.md`, `docs/planning.md` (the verdict quotes and the
  check messages), `docs/watching.md` (the claimed verdict and the
  attempts table) and `docs/reference.md` (`grove claim`,
  `grove claim <cycle-id>`, `grove claim --retry`; no `--timeout` row)
  name `grove claim`.
- `grep -rn "grove run\b\|--timeout\|\.git/grove" README.md docs
  dev/walkthrough.md` finds nothing.

`dev/walkthrough.md`:

- "Running unattended" is rewritten as seven steps on a bare `origin`,
  `~/tmp/ledger.git`, with two clones, `~/tmp/ledger-run` and
  `~/tmp/ledger-run2`: set up, claim a cycle and read its record, compose
  with a harness, the pause from two checkouts and a race, refusals,
  claim a cycle you name, and no remote. With a bare `origin`, the
  attempt's branch is merged into `main` by hand where the pull request
  would be.
- "Watching the project" step 3 has a second clone claim K and push its
  attempt's branch, in place of a hand-written `.git/grove/` record.
- The planning section's verdicts, check messages and the "Delegated"
  path name `grove claim`.

The walkthrough's claim steps were run with `bin/grove` on fresh fixtures
by a script. The sessions were stood in for by writing and archiving the
changes by hand, as the tests do. Every quoted message matched:

- nothing to claim, exit 1;
- the cycle ID alone, the attempt's branch at `main`, a clean
  `git status`;
- the record by `git log -1 --format=%B origin/grove/runs`, with an empty
  `ls-tree`;
- paused from the other clone, naming the attempt, host, time and branch;
- after the merge, two claims started together with `&`: one printed the
  budgets cycle's ID and the other said paused after that attempt, with
  two records on `origin`;
- the local-commits refusal, and the uncommitted refusal with `--retry`,
  the edit kept;
- a `pre-receive` hook refusing: no branch, still two records;
- a given cycle that needs a person, paused without `--retry`, claimed
  with it, the earlier records kept under it;
- no such cycle, a cycle waiting on a cycle, and an unmarked cycle
  claimed while another is marked;
- no remote: the record on the local `grove/runs` with `main`'s commit,
  paused, then nothing to claim after the merge, where
  `K=$(grove claim) && …` exits 1.

Not run here: the Claude Code sessions (walkthrough step 3, the
`/grove-propose` and `/grove-plan` sessions in step 6, and the Codex
pass), the launchd timer guide, and the page in a browser for "Watching
the project" step 3. Each needs a harness session or the owner's machine
set up for it. The script is not committed. It used shell stand-ins for
the agent's edits.

`just check` passes: 4338 Go lines, 1977 test lines, 408 page lines, 8613
instruction words.

Also changed: `internal/cli/claim.go` — running the walkthrough's refusal step showed a push `origin` refused printed bare Git text; it now says "the record was refused, so nothing is recorded:" before quoting Git, as the runs spec's "A push refused for permissions" asks.
Also changed: `internal/cli/walkthrough_test.go` — `TestClaim` matches that message.

## Scenarios

Each line names the test in `internal/cli/walkthrough_test.go`, or the
walkthrough step run by hand, that shows the scenario. Every test named
passed under `just check`.

- config "Building unattended": TestImplementInstructions (the Project context before what Grove requires of implement, unattended), passed
- config "No context": TestPrototype (propose instructions with no Project context), passed
- cycles "Nothing in the way": unverified; it is what the planning agent says. TestCycles checks that the plan instructions tell it to name `/grove-implement`, then `grove claim <cycle-id>`, then the mark, in that order, passed
- cycles "A member needs a person": unverified, the agent's words; TestCycles checks the instruction text, passed
- cycles "A member waits on a change outside the cycle": unverified, the agent's words; TestCycles checks the instruction text, passed
- cycles "The cycle waits on a cycle": unverified, the agent's words; TestCycles checks "`grove claim` will not claim the cycle until it is", passed
- cycles "The runner needs something set up": unverified, the agent's words
- cycles "Planned but not marked": TestClaim (nothing to claim with K unmarked), passed
- cycles "The cycle it follows is not done": TestFollows (`grove claim K` refused, no branch), passed
- cycles "The cycle it follows is done": TestFollows (K selected once L is archived), passed
- cycles "Built with the owner": TestFollows (`instructions implement K` lists its member), passed
- cycles "Marked with a change that needs a person": TestCycles (check's `unattended` line saying `grove claim` skips the cycle), passed
- cycles "Marked and stale": TestReplan, passed
- cycles "Marked and a member's file moved": TestMovedFiles (check, and no cycle selected), passed
- cycles "Revised and re-planned in one squash merge": TestReplan, passed
- cycles "Revised and re-planned in one merge commit": TestReplan, passed
- cycles "Marked after a member changed": TestReplan, passed
- cycles "Revision and re-plan as two commits on the target": TestReplan (stale, not selected), passed; `grove claim K` itself not run on it
- cycles "A member lands in its own squash merge": TestFollows, passed
- cycles "Then a commit archiving nothing": TestFollows, passed
- cycles "The cycle it follows lands in a squash merge": TestFollows, passed
- cycles "Stacked pull requests": TestFollows (L's changes landed as two squashes), passed
- cycles "Through another cycle": TestFollows (M's change landing in a merge commit), passed
- cycles "Any other commit still counts": TestFollows, passed
- implement "An approved change": TestImplementInstructions (instructions), passed; the agent's build unverified
- implement "A draft": TestImplementInstructions, passed
- implement "A change that waits": TestCoordinating, passed
- implement "A cycle": TestImplementInstructions and TestCoordinating (the build order, "only `grove claim` reads them"), passed; the agent's build unverified
- implement "A member left out": TestCoordinating, passed
- propose "A blocked member answered": unverified, the agent's edit; TestCycles checks that propose prints the member's cycle file, passed
- propose "A member that no longer needs a person": unverified, the agent's edit; the claim after the squash merge is TestReplan's squash case, passed
- propose "A change in no cycle": TestCycles (propose prints no cycle file), passed
- runs "Two hosts at once": TestClaim (h2's pre-push hook runs h1's claim; h2 is paused and creates no branch), passed; by hand in walkthrough step 4 with `&`
- runs "A push refused for permissions": TestClaim (`pre-receive` hook: nothing on origin, no remote-tracking ref, no branch), passed; walkthrough step 5
- runs "Local target has commits of its own": TestClaim, passed
- runs "A marked cycle of two changes, one waiting on the other": TestClaim (the record names K, A then B, the target commit, the host and the branch), passed
- runs "Composed with a harness": TestClaim (the ID alone on stdout, on the attempt's branch), passed; the harness session itself unverified, no session run
- runs "Two marked cycles": TestClaim (K, created first, over L, the lower ID), passed
- runs "The older marked cycle waits on the newer": TestFollows (L selected while K waits on it), passed
- runs "A given cycle not marked": TestClaim (`--retry` M while L is marked), passed; walkthrough step 6
- runs "A given cycle while paused": TestClaim, passed
- runs "A checkout with no remote": TestClaim, passed; walkthrough step 7
- runs "The branch merged locally": TestClaim, passed; walkthrough step 7
- runs "A member needs a person": TestClaim, passed
- runs "A stale plan": unverified through `grove claim K` itself. TestReplan and TestMovedFiles show the same `runnable` check refusing the stale cycle, passed
- runs "A cycle it follows not done": TestFollows, passed
- runs "No such cycle": TestClaim, passed
- runs "Every change archived": TestClaim (the last claim with no remote), passed
- runs "An approved change in no cycle": TestClaim, passed
- runs "A planned cycle not marked": TestClaim, passed
- runs "A marked cycle that went stale": TestClaim, passed
- runs "Composed with a harness, nothing to claim": walkthrough step 7 by hand (`K=$(grove claim) && …` exits 1, nothing started), passed
- runs "Retrying after an interrupted edit": TestClaim (`--retry` refused, the edit kept), passed
- runs "The first attempt of a project": TestClaim (one commit, empty tree, no `.git/grove/`), passed
- runs "A later attempt": TestClaim (two records, the older as it was), passed
- runs "Read by hand": walkthrough step 2 (`git log -1 --format=%B origin/grove/runs`), passed
- runs "Nothing in the working tree": TestClaim (`git status` clean, HEAD at the fetched target), passed
- runs "Pull request not merged yet": TestClaim, passed
- runs "Another host": TestClaim (h2, a separate clone), passed
- runs "Another cycle marked meanwhile": TestClaim, passed
- runs "Cycle merged": TestClaim (L claimed after K's members are archived), passed
- runs "Blocked change dropped": TestClaim (A finished, B dropped, not paused), passed
- runs "A blocked change revised": TestReplan (the squash re-plan is fresh) and TestClaim (`--retry` records anew), passed
- runs "The earlier attempt kept": TestClaim (the earlier record and branch kept after `--retry`), passed
- serve "A marked cycle that is fresh": TestServe, passed
- serve "A member needs a person": TestServe, passed
- serve "A member changed since the plan": TestServe, passed
- serve "A claimed cycle": TestServe, passed
- serve "An attempt from another host": TestServe, passed
- serve "No record anywhere": TestServe, passed
- serve "Readiness at a glance": unverified; how the page reads needs a browser
- serve "A refresh while reading": unverified; needs a browser

## Reconciliation

- config: applied to specs/config/spec.md
- cycles: applied to specs/cycles/spec.md
- implement: applied to specs/implement/spec.md
- propose: applied to specs/propose/spec.md
- runs: applied to specs/runs/spec.md
- serve: applied to specs/serve/spec.md
