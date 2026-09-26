---
id: "G-169"
type: work
title: "Keep installed harness entrypoints compatible with Grove upgrades"
status: done
created: "2026-09-25T21:04:46Z"
updated: "2026-09-26T00:28:06Z"
relates_to: ["G-101", "G-150", "G-152", "G-040", "G-110", "G-170", "G-186"]
candidate: "f2ba8d0f8a13c3de41a6bc119790ba2e5aeb1d04"
approved: "f2ba8d0f8a13c3de41a6bc119790ba2e5aeb1d04"
---

## Outcome

An adopting project can upgrade Grove and continue interactive shaping and work
through Claude and Codex without silently combining incompatible installed
instructions and executable behavior. Workflow changes travel with the binary;
project-owned policy stays with the project.

Owner intent, review and `$grove-shape` conversation 2026-09-25: remove reliance
on Grove's development checkout and the owner's environment, and make the skills
that `init` installs understandable to maintain and distribute. The owner asked
to shape the review's recommendations; the design below was proposed, and
Evidence records what was implemented.

## Observed evidence

At main `05892a2`, [init.go](../internal/cli/init.go) generates skills which
load `grove guide work|shape`, but also encode allowed arguments. G-040's
candidate `4edd800` rejects input other than IDs and interaction mode; today's
template accepts `--until plan`. An older installed adapter can therefore
contradict a newer guide. The full [reviewer](../.claude/agents/grove-reviewer.md)
is copied, rather than loaded at use time.

The review ran a binary built from `git archive HEAD` in a disposable external
Git repository with a stripped environment. After simulated edits to managed
skill and reviewer files, `check` still printed `OK: 0 records`; rerunning
`init` restored both and preserved an unmarked custom skill. This proves current
refresh and ownership behavior, not a real two-release or harness upgrade trial.
[G-150](G-150-launch-attempts-only-where-the-w.md) checks skill presence before
launch; it does not establish compatibility. Existing worktrees retain their
own files after another checkout runs `init`.

## Proposed design and constraints

- Keep one binary-owned workflow. Move evolving assignment grammar into the
  guide. Make the reviewer load binary-owned review instructions, for example
  through `grove guide review`; retain necessary harness metadata and the
  reviewer's read-only authority in its adapter. Preserve argument-as-data
  handling and explicit invocation.
- Give the small entrypoint interface an identifiable compatibility revision,
  separate from package releases and record schema. Specify how legacy files
  without metadata are diagnosed. Content inequality alone is not evidence of
  incompatibility; compatible older entrypoints should remain usable.
- Provide a read-only installation diagnostic, for example `init --check`,
  that distinguishes current, compatible older, incompatible, missing and
  project-owned files. Make supported interactive entrypoints and the runner
  surface incompatibility before work or paid execution begins. Preparation
  must explain how legacy entrypoints reach that diagnosis and recovery.
- Retain explicit `init` refresh and the existing ownership marker. It may
  replace marked files; it must preserve unmarked custom files, configuration,
  brief, records and AGENTS/CLAUDE policy. Never silently repair another
  worktree, commit changes, or claim a custom file is compatible without evidence.
- Explain commit requirements, existing-worktree repair, and restarting or
  reloading sessions which already loaded older instructions. Keep guide access
  usable without a project. Keep this repository's development adapters working.

Alternatives: separately released harness plugins add an installation/version
lifecycle; full copied guides add synchronization work. Neither is proposed.
No new execution provider, daemon, self-update, schema migration or broad
compatibility promise is authorized. [G-101](G-101-attempt-mechanism.md) remains
the runner decision; [G-152](G-152-shipped-document.md) names the shipped-document
boundary. G-040's marker-only approach would be extended, with existing
project ownership preserved. Release policy belongs to a first-release
policy record (named G-172 when this was shaped, a record no branch holds),
not an inferred schema change here.

## Acceptance

1. In a disposable adopting project, changing only the installed binary changes
   the work, shaping and review instructions used by a fresh session without
   rewriting compatible entrypoints. Evolving argument rules have one owner.
2. Read-only diagnosis covers current, compatible older, incompatible, legacy,
   missing and custom entrypoints without writes. Incompatible supported
   invocations stop actionably; runner refusal occurs before paid execution.
3. Explicit refresh is idempotent, preserves project-owned content, and gives a
   usable path for a committed older checkout and an existing worktree. Tests
   exercise an actual old/new fixture pair, not only repeated fresh `init`.
4. Evidence distinguishes deterministic checks from fresh Claude/Codex session
   discovery and interactive invocation. Live trials need an explicit resource
   mandate; unavailable trials are reported as limits, not simulated successes.
5. Command documentation and shipped instructions explain the contract and
   repair path; repository verification and shipped-document checks pass.

## Evidence

Implemented 2026-09-25 by a headless `/grove-work G-169` session on
`worktree-G-169` from `main` `6b14141`, record at `sha256:22d79d21…` and plan
[G-184](G-184-g-169-plan-entrypoint-revisions.md) at `sha256:b8613345…` when
implementation started. Commits: `0cb2251` plan, `dca444b` active,
`bde8bd8` implementation, `2b54de0` records, `d6c53df` and `549e477` round-1
fixes, `c5c5ffe` and `5a163a1` round-2 fixes, then the evidence commit named
in `candidate`, which changes only records.

**What changed.** The assignment grammar and the review brief moved into the
binary: the work guide's Inputs own the grammar, and the reviewer's brief is
`docs/work-review.md`, printed by `grove guide review`. Generated entrypoints
keep only the argument-as-data rule, harness metadata and (for the reviewer)
its read-only authority, and load their guide with `grove guide NAME
--entrypoint 2`. Each managed file names its
[entrypoint revision](G-186-entrypoint-revision.md), an integer apart from
the release version and schema. This binary writes and serves revision 2. A
marked file with no revision line is revision 1, `legacy`, and is not served
(review round 1 found that its generations disagree). `init --check` reports
each managed path as current, compatible, legacy, incompatible, missing,
custom or conflict, writes nothing, and exits 1 on legacy, incompatible,
missing or conflict. `guide --entrypoint N` refuses an unserved N. `run` and
the board's `R` refuse a marked skill or reviewer of an unserved revision,
before any branch for a new one and before the attempt directory for an
existing one. The work and shaping guides tell a session loaded through a
revision-less managed skill to stop and name the repair. This repository's
reviewer became a development adapter like its skills (no marker, reads
`docs/work-review.md`), and all its adapters dropped the grammar.
[G-152](G-152-shipped-document.md) now names the review guide. The dangling
G-172 link was replaced by text: no branch holds that record.

**Against acceptance.**

1. Deterministic, verified by hand: a binary built from `bde8bd8` with a line
   appended to `docs/work-review.md` reported entrypoints written by the
   unedited binary as all `current`, and its `guide review --entrypoint 2`
   printed the edited text. Round 1 repeated this for the work guide.
   `TestGuideServesEntrypointRevisions` asserts every generated Markdown
   entrypoint loads its guide with its revision. Evolving grammar is owned by
   the guides' Inputs alone. A legacy install keeps its own grammar until
   refreshed, and is not served.
2. `TestInitCheckDiagnosesAnOldInstallAndInitRefreshesIt` covers legacy,
   custom, current, compatible (other text at revision 2), incompatible (99,
   0, `two`) and missing, and asserts `--check` writes nothing.
   `TestGuideServesEntrypointRevisions` covers the `guide` refusal.
   `TestRefusals` covers a newer and a revision-less skill in HEAD (no branch
   made) and an unserved reviewer on an existing branch, all before any
   attempt. By hand with a binary from `549e477`, in a project initialized by
   `a84d01d^`'s binary (the pre-G-134 generation), `init --check` gave six
   legacy and a missing reviewer (exit 1). `run --until plan` was refused
   before any branch. After `init` and a commit, the check exited 0. A
   binary from `6b14141` refuses `guide work --entrypoint 2` (exit 2), which
   stops a new entrypoint running under an old `grove`.
3. The fixture `internal/cli/testdata/entrypoints-revision-1.txtar` holds the
   seven files `init` wrote at `6b14141`. Round 1 regenerated it and found it
   byte-identical. The test above runs it through `init --check` (legacy),
   `init` (six updated, the custom file kept), `init --check` (current) and
   `init` again (unchanged, no hash changed). By hand with real binaries: the
   old `init`, then the new one updated all seven and kept `grove.yaml`, the
   record root and the brief. A second run left the tree clean. A linked
   worktree from the old commit stayed legacy until `init` ran there, then was
   current. The command reference's Entrypoint revisions section gives the
   repair path: commit, repair in an existing worktree, restart sessions.
4. Only deterministic checks and scripted binary runs were performed. No
   fresh Claude or Codex session was started: the assignment stated no
   resource bounds for live trials. Unverified, therefore: that a fresh
   session discovers the revision-2 skills and follows `--entrypoint` loading
   (the invocation text differs from G-040's trial). Also unverified: that a
   harness keeps a skill's HTML-comment marker visible, which the legacy stop
   depends on (review round 3 saw an agent definition's comment reach the
   reviewer). A current `/grove-work` session loads through this repository's
   development adapter, not a generated one.
5. `docs/commands.md` (Attempts, Init, Entrypoint revisions, Version and
   guide), the README's adoption block, `grove --help`, CLAUDE.md and the work
   and shaping guides describe the contract and repair. The review guide
   passes the shipped-document test.

**Verification** at `c5c5ffe`'s tree, whose code equals `5a163a1`: `go vet
./...` was clean, `gofmt -l .` was empty, and `go run ./cmd/grove check`
reported `OK: 179 records`. `go test -count=1 -timeout 120s ./...` passed
every package. `python3 internal/tui/testdata/terminal.py` passed 12 of 12.
Inherited limit: `go test -count=1 ./internal/attempt` alone takes 5.5 to
5.6 s against 5.3 to 5.4 s at base `6b14141` under the same load; it was
already over five seconds (G-150).

**Review.** [G-187](G-187-review-of-g-169-entrypoint-revis.md), three
independent `grove-reviewer` rounds. Round 1's legacy finding changed the
design (plan adjustment). Every consequential finding was fixed, and the
last round found nothing consequential.

**Decisions taken inside the outcome** (G-184 records why): an integer
revision rather than a digest table; legacy unserved rather than guessed;
`guide` without the flag always prints, for people; the runner checks
against the launching `grove` only, as stated in the docs. New term
[G-186](G-186-entrypoint-revision.md), proposed.

## Next

In review 2026-09-25 on `worktree-G-169` (base `6b14141`). The candidate is
the evidence commit named in `candidate`; the branch tip adds only this
status change. For the owner:

- Judge the design choice review round 1 forced: every existing install
  (this repository's pilot adopter included) is `legacy` and must run
  `grove init` and commit before `run` will launch there, and its sessions
  are told to stop. The alternative, serving the `6b14141` generation by
  content, was rejected as a digest table (G-184).
- Consider settling term G-186.
- Live Claude/Codex discovery of the revision-2 skills needs a resource
  mandate if wanted before release work (G-110).
- Demo, from this worktree:

  ```sh
  B=$(mktemp -d); go build -o $B/grove ./cmd/grove
  P=$(mktemp -d); git -C $P init -q; $B/grove --project $P init
  $B/grove --project $P init --check; $B/grove guide work --entrypoint 1
  ```

Integration, as given:

```sh
go run ./cmd/grove approve G-169 "VERDICT"   # in this worktree
go run ./cmd/grove integrate G-169           # in the main checkout
```

Verdict on candidate f2ba8d0, 2026-09-26: approved with the requirement to run grove init in existing projects
