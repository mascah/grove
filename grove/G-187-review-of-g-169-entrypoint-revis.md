---
id: "G-187"
type: review
title: "Review of G-169 entrypoint revisions"
status: current
created: "2026-09-25T23:20:24Z"
updated: "2026-09-25T23:29:11Z"
work: ["G-169"]
examined: "5a163a1977a8e6ec82f81e11ce9276c2b0ba94f7"
---

## Examined

Three rounds by fresh, independent `grove-reviewer` agents, 2026-09-25,
dispatched by the headless `/grove-work G-169` session on `worktree-G-169`
(base `main` `6b14141`), each told the record, the plan
[G-184](G-184-g-169-plan-entrypoint-revisions.md), the commands it could run
and, from round 2, the earlier findings and their dispositions. Each built
binaries from this and older commits into temporary directories and ran them
in disposable repositories; none wrote to the checkout.

- Round 1: `6b14141..2b54de0` (implementation `bde8bd8`).
- Round 2: `6b14141..549e477` (fixes `d6c53df`, `549e477`).
- Round 3: `6b14141..5a163a1` (fixes `c5c5ffe`, `5a163a1`). `examined` is
  `5a163a1`, the last commit to change anything but records.

## Findings

Round 1:

1. Revision-less managed files span generations: before G-134 the work skill
   rejected `--until plan` and there was no reviewer, and the pilot adopter
   holds that generation. Serving them as revision 1 let `run --until plan`
   launch into a contradiction after spend. Verified with a binary from
   `a84d01d^`.
2. No rule said when revision 1 stops being served, though its files carry a
   grammar and review brief the next guide change would contradict.
3. "Entrypoint revision" was a shared concept with no term, and the usage
   placeholder `REVISION` overloaded settled term [G-062](G-062-revision.md).
4. `run` checks against the launching `grove`, not the one the session runs;
   the docs claimed more.
5. Minor: the adapters' repair hint named `init --check`, which an old
   `grove` lacks; the data rule read toward passing the bound to commands;
   the record's Next was stale; the plan's shaping-guide step was not done.

Round 2: all round-1 findings resolved except the stale Next (deferred to
handoff). New: (1) a legacy skill's interactive session was still never
stopped, since plain `guide` always prints; (2) the plan's Design
contradicted its adjustment; (3) the `init --check` legacy note claimed a
review brief for every path; (4) [G-186](G-186-entrypoint-revision.md) said
all three commands "refuse a marked file".

Round 3: all round-2 findings resolved; the new guide sentence satisfies the
shipped-document test and cannot fire for a revision-2 skill (verified), this
repository's unmarked adapters (verified) or a person. New: two leftover
sentences in the plan. Nothing else consequential.

## Disposition

- R1.1 and R1.2: fixed in `d6c53df`. `MinEntrypointRevision = 2`: legacy is
  not served, `init --check` exits 1 on it, and `run` refuses it in HEAD and
  in an existing worktree. From revision 2 an entrypoint holds nothing the
  guides evolve. Tests: the legacy and newer cases in `TestRefusals` and
  `TestInitCheckDiagnosesAnOldInstallAndInitRefreshesIt`.
- R1.3: term G-186 (proposed); usage writes `--entrypoint N`.
- R1.4: stated as a limit in the command reference's Entrypoint revisions.
- R1.5: adapter wording fixed in `d6c53df`; the plan records the shaping
  guide; the Next was reconciled at handoff.
- R2.1: fixed in `c5c5ffe`: the work and shaping guides' Inputs tell a
  session loaded through a managed skill with no revision line to stop and
  name the repair. A legacy reviewer loads no guide and stays reachable only
  through `init --check` and `run`, as the command reference says.
- R2.2 to R2.4 and R3.1: record and message wording, fixed in `c5c5ffe`,
  `5a163a1` and the evidence commit.

Open limit named by round 3: whether a harness strips HTML comments from a
skill's body before the model reads it is unobserved; if it does, the R2.1
sentence fails open for a legacy skill as before, and it still cannot misfire
for a revision-2 one. Round 3 observed that an agent definition's comment
reached the reviewer.
