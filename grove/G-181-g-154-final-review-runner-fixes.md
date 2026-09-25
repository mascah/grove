---
id: "G-181"
type: review
title: "G-154 final review: runner fixes and the without-row handoff"
status: current
created: "2026-09-25T21:40:39Z"
updated: "2026-09-25T21:40:58Z"
work: ["G-154"]
examined: "c45af99"
---

## Examined

The final gate of [G-154](G-154-listed-constraint-eval.md) (plan
[G-155](G-155-g-154-listed-and-code-constraint.md) step 6), on
`worktree-G-154` from main `b684951`, by two fresh `grove-reviewer` agents
dispatched from the headless `/grove-work G-154` session of 2026-09-25,
read-only, writing only under `/tmp`:

- Round 1 examined `b684951..9a226b2`, the whole branch, most closely what
  the earlier review [G-159](G-159-g-154-runner-cases-review.md) (at
  `47223ed`) never saw: `c895666`, `40e1263`, and the amendments after
  [G-173](G-173-what-should-g-154-s-with-row-bec.md)'s answer "Drop the with
  row".
- Round 2 examined `9a226b2..c45af99`, the fixes, and the whole branch's
  readiness. `examined` is `c45af99`; `144329d` after it rewords G-160
  finding 3 and its Disposition as round 2 finding 1 asks, self-checked
  only.

Both ran `python3 evals/run.py selftest` (`selftest: ok`) and `go run
./cmd/grove check` (OK), reran `retrieval()` on the ten retained `without`
transcripts under `~/.cache/grove-evals/runs/2026-09-25-G-154-without/`,
and probed the runner with copies that revert each fix. No paid eval ran,
and no Codex.

## Findings

Round 1:

1. **Minor, fixed in `7026cbe`.** Since `40e1263` every reader argument
   went through `glob.glob` joined to the clone path, so a `[`, `*` or `?`
   in `--out` hid every read. Now `glob.escape(clone)`; a selftest clone at
   `odd[1]` fails without it.
2. **Minor, not changed.** The `for`-loop test finds a reader and `$VAR`
   anywhere in the command, not only in the loop body, so `do echo "cat
   $f"` counts as a read. The README says "in the same command"; no
   retained transcript has that shape, and the facts are never scored.
3. **Minor, fixed in G-160.** G-001, a record no case needs, was also read
   through `grove show` in four runs, which no fact counts (the unneeded
   reads count files only): 8 of 10 runs read it. G-160 finding 4 says so
   and that a comparison on unneeded reads must add those.
4. **Minor, fixed in `7026cbe`.** `c895666` excluded only absolute
   `--include` paths, but `context` refuses whatever `fs.ValidPath`
   rejects. The runner now excludes a path with an empty, `.` or `..`
   element; a selftest assertion on `./grove/…` fails with the old rule;
   the README says so.
5. **Minor, fixed.** Stale text after G-173: G-160 finding 3's last
   sentence and its Disposition's `with`-row sentence, and G-173's Next.
   G-154's Next is rewritten at the handoff.

Round 1 also found: acceptance 1, 2 and 4 met; 3 met by G-160 after the
amendment; the amendments to acceptance 2 and 3 do what G-173's answer
says and nothing more, and the Outcome's clause on G-153's search was
removed with them, which the owner should see; no decision record is
needed for G-173 (cheap to reverse, attributed in G-154's Outcome) and no
term is introduced or contradicted; the branch merges into main cleanly.

Round 2:

1. **Low, fixed in `144329d`.** G-160 said every run reads every record,
   against its own 8-of-10 count for G-001: now "nearly every record".
2. **Informational.** `context` also refuses a `.git` element and
   symlinks, which the runner's include rule does not copy; neither can be
   a record read.

Round 2 verified each round 1 fix against copies with it reverted, found
`retrieval()` at `9a226b2` and `c45af99` identical on all ten runs, and
reproduced G-160 finding 4's recomputed facts.

## Disposition

No consequential finding remains. Two informational limits stay open:
round 1 finding 2 and round 2 finding 2. The branch is ready for handoff
into Review; this record is evidence, not approval.
