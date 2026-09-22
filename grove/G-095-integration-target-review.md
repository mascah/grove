---
id: "G-095"
type: review
title: "Review of G-042 integration target"
status: current
created: "2026-09-22T22:48:37Z"
updated: "2026-09-22T22:48:49Z"
work: ["G-042"]
examined: "3851182"
---

## Examined

The same independent reader as [G-094](G-094-current-view-review.md) (a
Claude reviewer subagent of the implementing session, with no edit rights over
the branch) reviewed the integration target that the owner added to
[G-042](G-042-current-view.md) after candidate `046150e`. It worked in two
rounds, reading the code and running scratch tests in a copy of the worktree:

- Round 1 examined `git diff 046150e..00ddefe`, against the plan's
  [revised decision](G-093-current-view-plan.md#revised-decision-integration-target-2026-09-22).
- Round 2 examined `3851182`, the fix. No defects remained.

## Findings

Round 1 confirmed one defect by running it. A target branch with no
`grove.yaml` was dropped with the note "could not be read". That is exactly
the adoption case, where the branch adding Grove names a `main` that has no
project yet. The note was false, and the target was lost.

It found no defect in the other areas it checked:

- `OnTarget` for deletions, absence and one-record inspections;
- invocation independence;
- the tag and label rules, including divergence, uncommitted and the Deleted
  shelf;
- escaping of the target text on every surface;
- the claim that the target is never passed to Git.

Suggestions:

- Rebuilding the installed binary after the merge: old builds reject the new key.
- Refusing likely typos in the value.
- Filling test gaps for deletions against the target, a divergence main holds
  one side of, and the shelf tag.
- The version package's test time.

It also noted a subtlety that matches the plan: on the target means the same
bytes.

Round 2: none.

## Disposition

- **The adoption defect:** fixed in `3851182`. `TestCurrentViewTargetBeforeAdoption`,
  which is the reviewer's scenario, fails at `00ddefe`.
- **Typos:** a value with surrounding spaces or a `refs/` prefix is refused
  (`TestTarget`).
- **Test gaps:** filled in `3851182`, in the versions and TUI target tests.
- **Rebuild:** named in G-042's integration steps.
- **Test time:** the versions package's short tests take 4.1 to 4.4 s alone.
