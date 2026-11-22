# Changes to explore

<!-- The changes that change 261122-m78w6 proposes to the requirements of
`explore`. This is future work: nothing here describes the target branch
until the change is implemented and the maintained spec is reconciled. -->

There is no maintained spec for `explore` at b6f43a0. Explore's suggestion
of a prototype is in `grove/specs/prototype/spec.md`, and of `bugfix` in
`grove/specs/schemas/spec.md`; this capability is how explore converses
with the owner.

## ADDED Requirements

### Requirement: Recommend beside every set of options

When the explore instructions have the agent offer the owner options, they
shall have it say which option it recommends and why, and treat the
owner's answer, including a "yes" to the recommendation, as settling the
question.

#### Scenario: Two ways to keep an older record editable

- GIVEN the owner explores "Editing an old record fails because its
  category was retired", and the agent finds two ways forward: keep an
  older record's category when the edit leaves it unchanged, or let any
  category through on edit
- WHEN it offers them
- THEN it names the one it recommends with its reason, and the owner's
  "yes" settles it as the owner's choice

### Requirement: Ask only what would change the recommendation

While exploring, the explore instructions shall have the agent ask the
owner at most one question a turn, or two in a first reply that follows an
investigation, only a question whose answer would change what it
recommends, and, once no remaining question would, stop asking and finish.

#### Scenario: A question that changes nothing

- GIVEN the agent recommends keeping an older record's category when it is
  unchanged, and whether the README says "every record" or "older records"
  would not change that recommendation
- WHEN it decides what to ask next
- THEN it does not ask the owner about the README's wording in explore,
  and it finishes, saying what it found, what the owner decided, its
  recommendation and what is still open

### Requirement: Investigate first when asked to

When the owner's description asks the agent to research, investigate or
recommend, the explore instructions shall have it read the code and the
documents that bear on the description before it asks anything, and make
its first reply its findings and its recommendation, with at most two
questions, each one whose answer would change the recommendation.

#### Scenario: Look into it and recommend

- GIVEN the owner runs `/grove-explore Should ledger keep its records in
  SQLite instead of records.json? Look into it and recommend.`
- WHEN the agent replies for the first time
- THEN the reply follows its reading of `ledger/records.py` and the
  README, gives its findings and the option it recommends with its reason,
  and asks at most two questions

#### Scenario: A description that asks nothing to be looked into

- GIVEN the owner runs `/grove-explore Users want to rename a record's
  category`
- WHEN the agent replies for the first time
- THEN the reply is one question, or a short account of what it will look
  up, and not a plan, a spec or a list of requirements

### Requirement: An earlier decision takes one line

If what the owner now wants conflicts with an earlier decision in a
maintained spec, an archived change or another project document, then the
explore instructions shall have the agent say so in one line, saying in
words what the earlier decision was, and continue with what the owner now
wants without asking whether to reverse it.

#### Scenario: Records move out of records.json

- GIVEN an archived change decided that records stay in `records.json`,
  and the owner now wants them in SQLite
- WHEN the agent mentions the earlier decision
- THEN it says in one line that this reverses the decision to keep
  records in `records.json`, and its next question or recommendation is
  about SQLite

## Observed current behavior

What the explore instructions say at b6f43a0:

- They forbid a recommendation before the owner has seen the options:
  "the owner sees the options and open questions before anything is
  recommended" (`internal/instruct/schemas/default/skills/explore.md`,
  lines 26–28), and "your own inference stays open, and you say so"
  (`internal/instruct/requires/explore.md`, line 5).
- They interview intent, then boundaries, then what the owner has tried,
  "one focused question at a time", and name no point at which to stop
  asking (`explore.md`, section 1).
- Their first reply "is a question or a short account of what you will look
  up", whatever the owner asked for (`explore.md`, line 16).
- They have every finding cited by "file path and commit" in the
  conversation (`explore.md`, line 21), and where documents or code
  disagree, ask which is intended (line 22).
