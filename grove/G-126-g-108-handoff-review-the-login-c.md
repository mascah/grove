---
id: "G-126"
type: review
title: "G-108 handoff review: the login change and G-122"
status: current
created: "2026-09-24T01:25:14Z"
updated: "2026-09-24T01:30:26Z"
work: ["G-108"]
examined: "c3386ae"
---

## Examined

An independent review of everything on `worktree-G-108` after
[G-119](G-119-g-108-eval-skeleton-review.md)'s `1f03a12`: the runner's
login change (`433e338`, `d565fcf`), the records of the paid runs
([G-122](G-122-g-108-baseline-runs-the-missing.md),
[G-108](G-108-workflow-evals.md),
[G-115](G-115-g-108-eval-skeleton-plan.md)) and the fixes each round
produced. A separate reviewer agent (Claude Code 2.1.281 subagent,
read-only, told not to edit, commit or run the real `claude`) examined
`fdf5d4e`, then `b6b3c7a`, then `c3386ae`; `examined` is the last. Its
evidence: `python3 evals/run.py selftest`, importing `evals/run.py` and
calling `run()` against temporary config directories with a nonexistent
`--claude`, so no spend, and reading every `run.json`, `report.md` and
transcript under the runs' output directory, plus the proposal branches
of missing-choice 2 and 4 and companion 3.

## Findings

Round 1, on `fdf5d4e`, seven findings:

1. High: the runner accepted `autoMemoryEnabled: true` (memory is written
   under the config directory and a reused `--out` could carry it between
   runs) and an authored skill under `skills/synced/ID/` that the sync's
   manifest does not name; accepting synced content at all was the
   session's decision, not G-121's, and unrecorded.
2. Low: a manifest entry without `name` and a `settings.json` that is a
   list crashed before spending instead of refusing.
3. Medium: the selftest never executed the `surfaced, not blocking`
   reason that G-122's Disposition relies on, and compared only which
   checks failed, not why.
4. Medium: G-122 scored "presumes choice" 2 for every run, but the
   rubric's 2 anchor says no item depends on the answer, and every run's
   acceptance does, while saying it is open: the 1 anchor's letter.
5. Medium: Disposition claims beyond the evidence: "blocking ten times
   out of ten across both cases" when the companion cannot block, a causal
   "the sentence produced" when the model changed too, MCP server
   descriptions "in context" when the servers offered no tool, a flat
   "36 skills" that mixes 18 built-ins and two fixture adapters with the
   sixteen account skills, and an untested claim about
   `--strict-mcp-config`.
6. Low: the handoff lacked the third session's starting revisions, a
   `check` result for the new records, and a named review of the code
   after `1f03a12`.
7. Low: "second headless session" for the third; a retrieval sentence
   wrong for two companion runs; only one of two connector mentions named.

Verified in round 1 against the output directory: every cost, turn count,
duration, check result, exit and denial count in G-122; 68 `Bash` and one
`Read`; fourteen MCP servers and 36 skills per session; no absolute path
outside the clone in any tool input (a regex over `/Users`, `/tmp`,
`/private`, `/var`, `/etc`, `/home`, `~`, `../` and `cd /`); the guides
digest; the rubric scores of the three spot-checked runs except finding 4.

Round 2, on `b6b3c7a`: all seven hold. New: A, medium, a directory with no
`settings.json` or without the key was accepted although memory is on by
default (the probe's own `projects/…/memory` directory predates the
owner's settings file); B, low, a stray file under `synced/`, a manifest
that is not JSON or whose list holds strings crashed before spending; C,
low, a refusal named `skills/synced/mine` without the ID; D, unverified,
a skill directory at the ID level with no manifest is accepted and it is
unknown whether Claude loads one from there; E, wording, "current default
model" for a model the runs requested explicitly.

Round 3, on `c3386ae`: A, B, C and E hold (a missing settings file, a
missing key and a string `"false"` are refused; a stray file, invalid
JSON, a list and string entries in a manifest are refused by name; the
real login directory is still accepted with the same settings, skills and
plugins the paid runs recorded). New: 1, medium, the round 2 fix masked
the selftest's key-refusal cases, since a settings file lacking
`autoMemoryEnabled` is refused by the memory check whatever the key check
does, and a mutant with the key check removed passed the selftest; 2, low,
a `manifest.json` that is a directory crashed before spending.

## Disposition

1. Fixed `b6b3c7a`: `autoMemoryEnabled` true refused, unlisted entries
   under a synced ID refused; G-108's Evidence records the decision.
   Round 2's A fixed `c3386ae`: the key must be present and false, the
   README says so.
2. Fixed `b6b3c7a` (`.get("name")`, a non-object settings file refused);
   round 2's B fixed `c3386ae` (non-directories and unreadable manifests
   refused by name), C likewise.
3. Fixed `b6b3c7a`: a `surfaced` fake mode, reason-prefix assertions for
   `bad` and `surfaced`, and nine refusal shapes tried by `c3386ae`.
4. Fixed `b6b3c7a`: scored 1 by the anchor as written, with the reading
   explained and a rescoring or rewording left to the owner after the
   fact.
5. Fixed `b6b3c7a`, E fixed `c3386ae`: each sentence now states what the
   runs show and labels the confound and the untested flag.
6. Fixed `b6b3c7a`: starting revisions, verification and this record
   named in G-108's Evidence; the candidate is set by the status change
   that follows the commit holding this record.
7. Fixed `b6b3c7a`.

Round 3's two findings were fixed after the cap, in the commit holding
this record, and are self-checked only: the refusal cases carry
`autoMemoryEnabled: false` where the key check is under test, each asserts
its own reason, each runs against a fresh output directory, and `OSError`
joins the caught manifest errors; `selftest: ok`, and the reviewer's
mutant (the key check removed) and a second (the memory check removed)
both fail the selftest with `settings.json accepted`. No fourth
independent round examined that commit.

D stays open as unverified and unchanged: the runner reports what a
manifest names and refuses what it does not; a skill at the ID level
without a manifest is a shape no login has produced here. The runs at
`d565fcf` were made with the looser check; the reviewer confirmed in
rounds 2 and 3 that the real login directory passes the tightened check
with the same recorded settings, skills and plugins, so those runs are
representative of the runner as handed off.
