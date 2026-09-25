---
id: "G-180"
type: work
title: "Resolve, approve and integrate candidates under an explicit owner policy"
status: proposed
created: "2026-09-25T21:39:28Z"
updated: "2026-09-25T21:56:44Z"
kind: feature
depends_on: ["G-177", "G-178"]
relates_to: ["G-044", "G-058", "G-059", "G-060", "G-101", "G-134", "G-140", "G-142", "G-162", "G-163", "G-179", "G-182"]
---

## Outcome

Where the owner has written a standing policy, Grove starts a resolution
attempt for a candidate that conflicts with the target, and approves and
integrates a candidate that meets the policy's conditions after its
independent review, with each act attributed to the policy and its
evidence, so a set of independently implemented items can go from review to
the target while the owner is away, and the owner judges only what the
policy leaves to them.

Decision [G-182](G-182-standing-policy-delegation.md), the owner on
2026-09-25 answering [G-179](G-179-standing-policy-question.md): option 4,
automatic resolution plus delegated approval and integration under a
narrow written policy the owner extends. The initial policy below is
proposed at the owner's request and binds nobody until the owner writes it
into `grove.yaml`.

## Constraints

### Observed evidence

At main `47852e3`, `approved` is a commit that must equal `candidate`, and
the verdict is a body paragraph (record model, work lifecycle); nothing in
the schema distinguishes the owner's verdict from a delegated one, and
`check` accepts any equal value, which [G-044](G-044-review-integration.md)
notes as the hand-set escape that skips the verdict and tip checks. The
independent reviewer (`.claude/agents/grove-reviewer.md`,
[G-134](G-134-bound-an-attempt-at-its-plan-and.md)) returns findings and
limits and is "never approval"; the [work guide](../docs/work-execution.md)
dispatches it per gate and caps fix rounds at three. `integrate` merges,
then writes done, and runs no verification; approval needs the branch's
clean checkout and integration the target's. Each attempt records its
model, effort, reviewer definition hash and cost
([G-134](G-134-bound-an-attempt-at-its-plan-and.md),
[G-140](G-140-default-an-attempt-s-budget-mode.md)), and nothing resident runs between
attempts ([G-101](G-101-attempt-mechanism.md)). G-179's evidence lists the
authority and spend constraints.

### Proposed design and scope

- A `policy:` mapping in `grove.yaml` beside `run:`, absent by default.
  Proposed initial policy for this repository, narrow on purpose:

  ```yaml
  policy:                 # standing delegation (G-182); absent means nothing automatic
    budget: 30            # dollars, aggregate for every automatic act in one sweep
    resolve:
      budget: 10          # one attempt per candidate per target movement; run: defaults otherwise
    approve:
      verify:             # must pass on the merged result, in a temporary worktree
        - go test -count=1 -timeout 120s ./...
        - go vet ./...
        - go run ./cmd/grove check
      max_lines: 300      # added plus removed, the record's own file excluded
      never:              # a change to any of these always waits for the owner
        - grove.yaml
        - .github/**
        - .claude/**
        - .agents/**
        - go.mod
        - go.sum
        - lefthook.yml
    integrate: true       # merge and write done after a delegated approval
  ```

  `never` is the trust boundary: files that change what the automation
  itself does or what CI runs. `max_lines` bounds the blast radius. `verify`
  is this repository's final evidence from its agent instructions. The
  budgets are guesses to adjust. Likely extensions, deliberately absent
  from the initial policy: a `paths:` allowlist, a `kinds:` list (for
  example only `fix` and `tooling`), a `sizes:` list, and a higher
  `max_lines`.
- Conditions the policy cannot switch off, so that its narrowness holds by
  construction: the record is in review with its candidate unchanged since
  the review; a review record examined that candidate and reports no open
  finding and no knowledge finding; no open question blocks the work; the
  merge is clean (a conflict goes to resolution, never to approval); the
  merged result passed `verify`. How "no open finding" is read from a
  review record is chosen in preparation, since the schema does not
  structure findings: a conventional closing line the reviewer definition
  writes, or a review field.
- Attribution: a delegated approval writes `approved` and a verdict
  paragraph that keeps the `Verdict on candidate X, DATE:` prefix
  `integrate` already quotes, with text beginning `delegated under policy`
  and naming the `grove.yaml` revision, the review record, the attempt and
  the verification result, so it is distinguishable from the owner's own
  verdict wherever the record is shown; the merge and `done` are the
  ordinary `integrate` acts. Introduce a field only if prose attribution
  proves insufficient.
- Verify before the target moves: make the merge in a temporary worktree,
  run the policy's commands there, and only then merge on the target. A
  failure leaves the target unchanged and the candidate in review with the
  reason, exactly as a refused `integrate` does.
- Trigger: the owner process of a finishing attempt, or an explicit sweep
  command the owner or an external scheduler runs. Grove starts no resident
  service (G-101); an unchanged wait does not retry.
- Everything the policy does not name waits for the owner as today.
  Reconcile the brief's lifecycle paragraph and the [approval](G-059-approval.md)
  term if the answer changes their meaning; the existing reviewer
  definition suffices initially, with its findings consumed by the policy.
- Out of scope: choosing the policy's conditions for the owner; a second
  reviewer definition; executing several selected items
  ([G-162](G-162-bounded-work-selection.md)); a shared candidate for a chain
  ([G-163](G-163-selected-work-review-boundary.md)).

## Acceptance

1. With no policy, nothing changes. With one, every automatic act is visible
   in advance, as what would happen to each candidate and why, and
   attributable afterwards, with the policy, evidence, attempt and cost in
   the record and on the board.
2. A candidate outside the policy waits for the owner. One inside it is
   integrated only after the merged result passed the policy's verification;
   a failure leaves the target and the record unchanged with the reason.
3. Automatic resolution runs at most once per target movement within the
   aggregate budget; a repeated conflict waits with the reason.
4. A delegated verdict is distinguishable from the owner's in the record and
   on the board, and reversible: the record names the merge to revert.
5. Exercised with fake providers on clean, conflicting, failing-verification
   and out-of-policy candidates, then a separately bounded real-provider
   trial on a disposable project before use in this repository.

## Next

G-179 is answered ([G-182](G-182-standing-policy-delegation.md)) and this
scope is refined to it, with the initial policy proposed above. Preparation
settles the policy's parsing, how a review's findings are read, the
attribution text, the verified merge and the sweep's trigger; then
assignment through `$grove-work G-180`, after G-177 and G-178, which it
depends on. No policy is in `grove.yaml` yet: writing one is the owner's
act. No spend, merge or implementation is authorized here.
