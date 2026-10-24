# <title>

<!-- A maintained spec states the agreed behavior expected of the target branch.
It is not a report of what the code does today or of what is deployed: where
they differ, keep the rule and record the difference under "Known discrepancies". -->

## Purpose

<!-- One short paragraph: the enduring area of behavior this covers, who relies
on it, and any term a reader needs defined. Its first sentence is the spec's
description in lists, so make it one line saying what the spec covers. -->

## Boundaries

<!-- What belongs to this capability, what is deliberately outside it, and which
neighbouring spec or system owns the rest. -->

## Requirements

<!-- One block per rule. State the rule as a single EARS sentence, then pin its
edges with scenarios. Pick the pattern that fits:

    When <event>, the system shall <response>.
    While <state>, the system shall <response>.
    If <unwanted condition>, then the system shall <response>.
    Where <optional feature is present>, the system shall <response>.
    The system shall <response>.

Example:

    ### Requirement: Retain previously valid values

    When a user submits an otherwise valid edit that leaves an existing
    record's previously valid category unchanged, the system shall save the
    edit.

    #### Scenario: Update the note on a legacy record

    - GIVEN a record whose category was valid when created but is no longer
      permitted for new records
    - WHEN the user changes only its note and saves
    - THEN the note is saved and the category remains unchanged

Describe behavior a user or caller can observe. Mention implementation only
where it explains a constraint. -->

### Requirement: [short name of the rule]

When [event], the system shall [response].

#### Scenario: [a concrete case at the boundary]

- GIVEN [starting state]
- WHEN [action]
- THEN [observable result]

## Known discrepancies

<!-- Where the target branch is known to differ from a requirement above: name
the requirement, say what was observed, and give the evidence (file path and
revision, test, issue). A defect does not change the agreed rule. -->

None recorded.

## Open questions

<!-- Unresolved points that affect a requirement, each phrased so that someone
can answer it. -->

None.
