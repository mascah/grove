# Prototype before proposing

## Problem

Between explore, which writes nothing and changes no code
(`internal/instruct/requires/explore.md` at 7ce7b91), and propose, which
writes the change, Grove has no step for trying an idea out. In the
2026-10-22 explore session the owner said: "we want to suss out things as
early as possible BEFORE comitting to changing the specs, because undoing it
later may not be so easy and just adds a lot of wasted time."

The owner prototypes by hand today. In keyborg, on branch
`prototype/accuracy-coaching`, `prototypes/coaching/` and
`docs/work/261021-coaching-prototype/handoff.md` are untracked, and the
draft change `261022-nmwhs` links to that handoff; the link breaks once the
proposal merges without it. The owner's other skill suite has a prototype
skill (`../skills/skills/prototype/SKILL.md`) that frames an experiment as a
bounded question with a verdict and limits.

## Intended outcome

`/grove-prototype <question>` builds a throwaway experiment in the owner's
session that answers one question: how something feels, whether something
is feasible, or how some logic behaves. It settles the question, its limits
and the evidence that would answer it, builds the simplest thing that shows
it on a `prototype/<slug>` branch, runs it, and reports the question, how to
run it, what was observed, a verdict and what it cannot establish. It
writes no Grove document and its branch is never merged (amended rule 5).

Explore suggests a prototype when a question's answer could change what
gets proposed, naming the question; the owner decides. Propose writes a
prototype's findings into the proposal and cites the branch and commit. A
prototype that kills an idea leaves only its branch and the conversation.

The rules are in [specs/prototype/spec.md](specs/prototype/spec.md).

## Scope

### In scope

- The prototype skill: `grove instructions prototype`, its requirements
  and the default schema's guide, the skill `grove init` writes, and
  `rules.skills.prototype`.
- Explore's guide suggesting a prototype; propose's guide using its
  findings.
- Documentation: a docs step that prototypes before proposing; the
  reference.

### Out of scope

- A prototype document kind (`prototype.md`); can be added later if
  proposals fill up with prototype detail.
- Promoting prototype code into the target; the change is implemented from
  its tasks.
- Pushing the prototype branch; the owner does that.

## Decisions

- **A distinct, optional prototype skill between explore and propose.** The
  owner, 2026-10-22: "We need an optional but distinct "prototyping" skill
  in between the explore and propose phase", taking inspiration from their
  own prototype skill.
- **Findings go in the proposal; code stays on a `prototype/<slug>` branch
  cited by commit; a dead end leaves only the branch.** The owner: "yes, go
  with 1". Reason: the proposal is what is reviewed; no new document kind.
- **Explore suggests one when its answer could change the proposal, and the
  owner decides.** The owner chose "1" from: suggest on named triggers;
  always offer; never suggest. They had said "Not all kinds of work will
  require a prototype, but we should suggest building one when it makes
  sense".
- **The constitution allows it** since #19 (rule 5).

### Assumptions to confirm in review

- The skill is named `prototype` (`/grove-prototype`).
- Prototype works on a branch `prototype/<slug>`, creating it from the
  current commit when the session is on another branch; it never switches
  to or commits on the target.
- At the end, prototype asks the owner whether to commit the experiment on
  its branch, so the proposal can cite a commit; it commits nothing
  otherwise, and pushes nothing.
- Prototype adds no dependency to the project beyond what the experiment
  needs, and nothing outside its branch.

## Open questions

None must be answered before review. None blocks a task.

## Related documents

- Requirements: [specs/prototype/spec.md](specs/prototype/spec.md)
- Tasks: [tasks.md](tasks.md)
- Constitution amendment: pull request #19
