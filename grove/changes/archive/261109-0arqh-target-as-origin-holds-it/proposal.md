# The target is what this checkout and origin both hold

<!-- A proposal describes future work. It becomes approved future work when it
is merged into the target branch, which says nothing about implementation.
There is no status or approval field to maintain. -->

## Problem

- Grove reads the target as the local branch and nothing else.
  `ResolveBranch` in `internal/gitx/gitx.go` at ccaee4d resolves
  `refs/heads/<target>`; `agreed` in `internal/cli/project.go` hands that
  commit to the project index, to the `{{commit}}` written as a change's
  `base` and a cycle's `planned`, to implement's "as main@X holds it", to
  archive and to review; `movedSinceBase` in `internal/cli/check.go` and
  `serve` resolve it the same way. Only `run` looks at origin: it fetches
  the target and refuses when the local branch has commits origin's does
  not (`runAttempt`, `internal/cli/run.go` at ccaee4d). `serve` lists
  origin's branches from the local `refs/remotes/origin` without fetching.
- In a project with a remote, the local branch can hold commits that
  origin's never will as themselves.
  [Issue #140](https://github.com/mascah/grove/issues/140), 2026-11-08,
  reports `grove check` failing permanently with
  `base: <sha> is not a commit in this repository` after "squash or
  rebase merges", a filter-repo cleanup or an import, with no repair but
  editing the YAML by hand. Experiments on the ledger fixture in this
  conversation on 2026-11-08, not committed, showed how the first two
  happen:
  - Two clones of one bare origin. In the first, one commit on the local
    `main` was never pushed; `grove instructions propose` on a branch
    there filled `base:` with that commit; the branch was squash-merged
    into origin's `main`. In the second clone, after a pull, `grove check`
    reported the `base` as not a commit. In the first, `grove check`
    stayed clean, since the commit lingers on the old branch and in the
    reflog, so the proposer never sees it.
  - One clone whose owner committed on the local `main`, proposed there,
    then ran `git pull --rebase` after origin had moved: the proposer's
    `check` clean, a fresh clone's `check` reporting the `base`.
  - A checkout with `origin/main` fetched and no local `main`, as
    `actions/checkout` produces for a pull-request job: `grove check`,
    `grove changes` and `grove instructions propose` each failed with
    `grove: git rev-parse: fatal: Needed a single revision: exit status 128`.
    The CI guide runs `grove check` on `pull_request`.
- Squash-merging branches into the target loses nothing: this repository
  squash-merges every pull request, and every `base` of its 42 archived
  changes and every `planned` of its 9 cycles resolves at ccaee4d. The
  commit is lost only when the local target branch itself held it and was
  rewritten, or when origin's history is rewritten.
- The bookmark `refs/remotes/origin/<target>` moves on fetch, pull and
  push, and never names a commit origin lacks unless origin rewrote
  history; it can only lag behind. The newest commit that both it and the
  local branch hold, their merge base, is therefore a commit origin has,
  and is never ahead of the working tree.

## Intended outcome

Every command reads the target at the newest commit the local branch and
the checkout's `origin/<target>` bookmark both hold, or at whichever of
the two exists when only one does, fetching nothing. A change's `base`
and a cycle's `planned` then name commits origin's history holds,
whatever is later squashed or rebased away on the proposer's machine, and
a checkout with the bookmark and no local branch, such as a pull-request
job's, works. A project with no remote reads the local branch as today.
`run`, and `claim` after it, keep comparing the local tip with the fetched
origin before fast-forwarding.

Approval follows: a change is approved once the target holds it as both
this checkout and origin do. A change merged locally and not yet pushed is
approved by the push.

## Scope

### In scope

- The rule, in the one gitx function every reader of the target calls,
  and its effect on the project index, `base`, `planned`, implement,
  check, archive, review and serve.
- A read of the local tip for run's fast-forward refusal.
- The `config` spec's new requirement, the documentation and the
  walkthrough.

### Out of scope

- Fetching: Grove reads the bookmark as it is (`CONSTITUTION.md` rule 7;
  watching.md, "fetches nothing").
- Remotes not named `origin`, as run and serve already assume.
- A force-pushed or rewritten origin, filter-repo, BFG and imports: the
  repair stays `/grove-propose <change-id>` per change and
  `/grove-plan <cycle-id>` per cycle. No `grove rebase`.
- A message when the local branch is ahead of origin's: the index still
  says "What main says at X", and the documentation explains why X can
  sit behind the local branch.
- Any special case for the false alarm of decision 3.
- Amending `CONSTITUTION.md` rule 4: its own pull request with no code,
  which lands before this change is built.
- `grove instructions implement` and `grove run` naming a missing `base`
  or `planned` in one sentence: change 261108-xmgxy, proposed alongside
  in [pull request #147](https://github.com/mascah/grove/pull/147).

## Decisions

The owner confirmed the scope above and accepted every item of the round
with "yes", 2026-11-09. Each item is as worded there.

1. **The rule.** "With both `refs/heads/main` and `refs/remotes/origin/main`
   present, the target's commit is their merge base; with one, that one;
   with neither, today's error. It governs everything that reads the
   target: the project index, the `{{commit}}` written as `base` and
   `planned`, implement's 'as main@X holds it', check's comparisons,
   archive, review and serve. Run, and claim after it, keep comparing the
   local tip with the fetched origin for their fast-forward refusal. Only
   the remote named `origin` counts, as run and serve already assume."
   Reason: the merge base is a commit origin has and is never ahead of
   the working tree, while `origin/main` alone can be ahead of it after a
   fetch nobody merged.
2. **Approval, and breaking.** "A change merged locally and not pushed is
   not approved to Grove until the push: `grove changes` omits it and
   implement says main holds no such change. A proposer whose main was
   ahead can also see a new check report on documents that were clean. By
   CONTRIBUTING's definition that is breaking, so `feat(cli)!` with a
   `BREAKING CHANGE:` footer saying to push local merges." Reason: a
   project with unpushed local merges sees different output after
   upgrading, and the footer tells its owner what to do.
3. **The false alarm, accepted.** "A proposer whose local main is ahead by
   commits touching a task's named files records the shared commit; once
   those commits reach origin as-is, check reports those files changed
   since `base`, and `/grove-propose <id>` clears it. No special case for
   it." Reason: the alternative, recording the local tip, is the bug.
4. **Constitution first.** "Rule 4's 'the target branch' gains a clause:
   the branch as this checkout and origin both hold it. Its own pull
   request with no code, before the change is built." Reason:
   `CONSTITUTION.md` says changing it is its own pull request that the
   owner reads before any code follows.
5. **Where the requirement lives.** "The `config` spec, which is the
   `grove.yaml` spec and `target` is its field, gains 'Read the target as
   origin holds it' with three scenarios: two clones of one origin, a
   checkout with the bookmark and no local branch, and a project with no
   remote, unchanged." Its Purpose and Boundaries widen when this change
   is archived. Reason: the alternative was a new spec for one
   requirement.
6. **No `depends_on`.** "261108-jpg8n names `gitx.go`, `run.go` and every
   doc file this touches, so whichever lands second re-bases with
   propose." Reason: this change is small, and waiting on claim would
   hold a fix for a reported failure behind a large feature.
7. **One function.** "The rule lives in `gitx.ResolveBranch`, which every
   reader of the target already calls, and run's fast-forward check reads
   the local tip through a separate one-line call. No new message."
   Reason: one place to read the target keeps every command agreeing on
   it; Git runs only in `internal/gitx` (`CLAUDE.md`).
8. **Routine defaults.** The title above; two tasks, T-01 gitx, its
   callers and a test with a bare origin, two clones and a checkout with
   no local branch in `walkthrough_test.go`, T-02 the documentation where
   it says "your main branch's commit", the CI guide noting a
   pull-request job's checkout works, and a walkthrough step; archive's
   "where this branch left the target" unchanged; test layout, private
   helpers and the local-tip helper's name delegated.

## Open questions

None.

## Related documents

- Requirements: [config](specs/config/spec.md)
- Tasks: [tasks.md](tasks.md)
