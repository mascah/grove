# Prototype

## Purpose

`/grove-prototype <question>` answers one question before a change is
proposed, with a throwaway experiment on a prototype branch in the owner's
own session; explore suggests one when the answer would change the
proposal. The question may be how something feels, whether something is
feasible, or how some logic behaves. Its findings reach the target only in
a proposal. "The owner" is the person in whose session the skill runs.

## Boundaries

This covers the prototype instructions, explore's suggestion of a
prototype, and propose's use of a prototype's findings. The rest of
exploring and proposing belongs to their own instructions. Promoting
prototype code into the target is out: a change is implemented from its
tasks. Pushing the prototype branch is the owner's. Prototype saves no
document of its own and its branch is never merged (`CONSTITUTION.md` rule
5).

## Requirements

### Requirement: Prototype a question

When the owner runs `/grove-prototype <question>`, the system shall print
instructions for agreeing the question, its limits and the evidence that
would answer it with the owner, building the simplest throwaway experiment
that shows it on a `prototype/<slug>` branch, running it on representative
cases, and reporting the question, how to run it, what was observed, a
verdict and what it cannot establish.

#### Scenario: How a screen feels

- GIVEN the owner asks whether a highlight-then-retype flow is pleasant to
  use
- WHEN the agent prototypes it
- THEN a runnable experiment is on a `prototype/<slug>` branch, and the
  report gives the command to run it, what was observed, a verdict and its
  limits

### Requirement: A prototype writes no Grove document

The prototype instructions shall write no Grove document, change nothing on
the target branch, and commit only on the prototype branch when the owner
asks.

#### Scenario: Started on the target

- GIVEN the owner's session is on the target branch
- WHEN the owner runs `/grove-prototype <question>`
- THEN the experiment is built on a new `prototype/<slug>` branch, and the
  target branch and the document root are unchanged

### Requirement: Explore suggests a prototype

When, while exploring, a question comes up whose answer a throwaway
experiment could show and that answer could change what is proposed, the
explore instructions shall have the agent name the question and ask the
owner whether to prototype it, and build nothing itself.

#### Scenario: Feel is unknown

- GIVEN the owner is exploring a new practice screen and does not know
  whether it will feel engaging
- WHEN the agent recognizes that
- THEN it names the question a prototype would answer and asks whether to
  run `/grove-prototype`

#### Scenario: A small, clear change

- GIVEN the owner is exploring a one-line message fix
- WHEN the interview ends
- THEN no prototype is suggested

### Requirement: Propose uses a prototype's findings

When the conversation holds a prototype's findings, the propose
instructions shall have the agent write them into the proposal, citing the
prototype branch and its commit, and link no file outside the change and
the target.

#### Scenario: Proposing after a prototype

- GIVEN a prototype on `prototype/coaching` at commit `abc1234` found that
  retyping one highlight felt right and three felt tedious
- WHEN the owner proposes from the conversation
- THEN the proposal states that finding and cites `prototype/coaching` at
  `abc1234`

## Known discrepancies

None recorded.

## Open questions

None.
