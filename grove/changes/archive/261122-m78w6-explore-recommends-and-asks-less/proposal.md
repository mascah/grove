# Explore recommends and asks less

## Problem

The owner reviewed every explore session held in Grove and in two other
projects of theirs from 22 October to 22 November 2026: 29 sessions. The
transcripts are on the owner's machine only, so this summarises them.

- The agent's message before each of the owner's replies had a median of
  358 words and a 90th percentile of 731.
- 116 of the owner's 350 replies were three words or fewer ("1", "yes",
  "a"), many of them accepting the first option or a recommendation.
- The owner asked "what do you recommend" or "which is better" at least six
  times, and said: "it's more about what is technically/architecturally
  better, which you haven't explained yet".
- Where sessions broke down, in the owner's words: "You keep burrying these
  questions away from their explanation"; "extreme question fatigue… batches
  of like 28 questions"; "You are not connecting the options clear enough";
  "you are making this impossible to decide… you're taking the fun out of it
  for me"; "I don't know how to make this choice right now as youve
  explained it"; "I don't want to go through another billion question
  interview right now. I want you to take what I said and do the
  investigation and produce recommendations."
- Earlier documents were raised as walls: an agent answered a direction with
  "which the 16 Nov draft lists under 'Not building'", and the owner replied
  "You need to recognize that I may have not previously understood enough."
- Change IDs, rule numbers and decision numbers were named about once a
  turn without saying what they were ("I don't understand what you mean
  about va9qn").

The instructions ask for this, at b6f43a0:

- `internal/instruct/schemas/default/skills/explore.md` lines 26–28: "the
  owner sees the options and open questions before anything is
  recommended".
- `internal/instruct/requires/explore.md` line 5: "your own inference stays
  open, and you say so".
- `explore.md` section 1 interviews intent, then boundaries, then what the
  owner has tried, "one focused question at a time", with no point at which
  it stops asking, and its first reply may not be a recommendation even
  when the owner asked for one.
- `explore.md` line 21 has every finding cited by path and commit in the
  conversation, which puts a findings report before each question.
- `internal/instruct/requires/propose.md` line 44: the round lists every
  decision the agent would otherwise assume, "Leave none out."
- `dev/trials.md` counts the words the owner reads, but none of its six
  cases runs explore.

## Intended outcome

Explore gives a recommendation, with its reason, beside every set of
options, and asks only what would change it. When the owner asks for
research or recommendations, it investigates first and comes back with
findings and a recommendation. An earlier decision that conflicts with what
the owner now wants takes one line, and the owner's direction is followed.

Every skill, when it asks the owner something, puts the question first,
makes it stand on its own, describes options by what the owner gains and
gives up, and says in words what any ID or rule number it names is.

Propose's round asks at most five decisions, in the same message as the
scope, and writes every other choice it makes as an open question with its
recommended answer, for the owner to strike when reviewing the change.

A seventh skill trial runs explore and counts what it costs the owner.

## Scope

### In scope

- The explore instructions (`explore` change spec).
- The opening every skill's instructions share (`instructions` change spec).
- Propose's round (`propose` change spec), in both schemas' propose guides.
- The documentation and walkthrough lines that describe explore's interview
  and the round.
- A skill trial case for explore.
- When archived, the propose spec's Purpose drops "asking one round of every
  decision it would otherwise assume", which this change makes false.

### Out of scope

- Unattended propose: it already writes its round as open questions and
  asks no one.
- Implement, plan, prototype and archive, beyond the shared opening.
- Any change to what explore writes: it still writes nothing unless asked
  (CONSTITUTION.md rule 5).

## Decisions

1. **All the changes reviewed, as one change.** The owner was shown a table
   of the changes above (explore recommends; asks only what would change the
   recommendation and investigates first when asked to; question first,
   options as gains and costs, about 150 words, citations only on request or
   in the proposal, the question standing on its own; plain words for IDs
   and rule numbers; an earlier decision gets one line; propose asks
   decisions that change behavior, scope or data, "cap ~5", and the agent
   makes the rest and lists them in the proposal; the getting-started line;
   an explore trial case) and the offer "One change; I'll make the small
   calls and list them in the proposal for you to strike on the PR." The
   owner answered: "Yes". Reason: the cost the owner pays is the number and
   length of the questions, and each item removes part of it.
2. **This change was written without a round.** By that same "Yes", every
   choice below the table's level is an open question with a recommended
   answer, not a decision.

## Open questions

All can wait for review; strike or change any on the pull request.
`/grove-implement` builds the recommended answer otherwise.

1. **Where the asking rules live.** Recommended: in the opening every
   skill's instructions share, not only in explore, because propose, plan
   and prototype ask the owner too, and the opening says it once. Decides
   T-02.
2. **Explore's question limit.** Recommended: one question a turn; after
   investigating on the owner's request, at most two, each one whose answer
   would change the recommendation. Reason: one question a turn was not the
   problem; questions that changed nothing were. Decides T-01.
3. **Propose's limit.** Recommended: five, the most consequential, in the
   existing order. More than five means the idea is likely more than one
   change, which the split offer in the same message covers. Decides T-03.
4. **One stop, not three.** Recommended: propose states the scope, offers a
   split, asks what must land first and asks the round in one message.
   Reason: each was a separate stop before. Decides T-03.
5. **Where the other choices go.** Recommended: as open questions with the
   recommended answer, its reason and the task it decides, the way
   unattended propose writes its round, routine choices on one open
   question of their own. Reason: Grove's existing rule that a proposal
   lists no assumptions holds, and `/grove-implement` already builds a
   recommended answer and records it for the reviewer. Decides T-03.
6. **Length.** Recommended: "about 150 words unless the owner asks for
   more" in the opening, as guidance, not counted by Grove. Decides T-02.
7. **A new `explore` capability.** Recommended: a maintained spec of its
   own, since no maintained spec at b6f43a0 says how explore converses, and
   the rules belong to neither `prototype` nor `instructions`. Decides T-01.
8. **The trial.** Recommended: case (g), `/grove-explore` given a brief that
   asks for a recommendation; a good run writes nothing, recommends (the
   scripted answer to a message holding a recommendation was given), and
   took at most three of the owner's replies. Decides T-05.

## Related documents

- Requirements: [explore](specs/explore/spec.md),
  [instructions](specs/instructions/spec.md),
  [propose](specs/propose/spec.md)
- Tasks: [tasks.md](tasks.md)
