---
id: "G-168"
type: review
title: "G-161 deps: first review gate on the shared model and command"
status: current
created: "2026-09-25T21:00:50Z"
updated: "2026-09-25T21:01:02Z"
work: ["G-161"]
examined: "c92d16f"
---

## Examined

The layout-independent part of [G-161](G-161-dependency-view.md), plan
[G-165](G-165-g-161-dependency-view-plan.md) steps 1–3, on `worktree-G-161`:
round 1 examined `3d132c6..4d65a57`, round 2 `3d132c6..c92d16f`. Two fresh
`grove-reviewer` agents, read-only, each running `go build`, `go vet`,
`gofmt -l`, `go test -short ./internal/deps ./internal/cli ./internal/handoff`,
`grove check` and `grove deps` against this repository. The board (step 4)
waits on [G-166](G-166-g-161-dependency-layout.md) and was out of scope.

## Findings

Round 1 (no correctness bug in the preview):

1. Medium: the overview printed no delivery for its rows, so a review
   candidate off HEAD and the target was visible only in `--json` or a
   preview.
2. Medium: `Compare` assumed the root checkout's records, which would give
   false "changed while it was being read" notes for a board View built from
   another source.
3. Low: `depends_on` elsewhere was compared order-sensitively.
4. Low: `deps` refusals named `context`.
5. Low: the overview lists only direct prerequisites that are not
   unfinished, so an abandoned X beneath a done D is not shown.

Round 2 (nothing consequential): the plan's Binding bullet named a
current-view overview without saying what source `Compare` gets there, and
its command contract listed stale columns.

## Disposition

1–4 fixed in `c92d16f` with regressions (`TestDepsCLI` asserts the review
row's delivery; case E in `TestCompareDescribesOtherVersions`;
`TestDepsUsage` asserts no `context` in refusals); round 2 confirmed each.
5 rejected, and round 2 agreed: a done D was delivered without X, so X
constrains nothing downstream of D; an unfinished X is its own row. Round 2's
two plan findings are fixed in the plan text; they touched no behaviour.
