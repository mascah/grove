# Reconcile a maintained spec

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

Archive adds a change's requirements to a maintained spec as the change
wrote them, so a spec grows one change at a time, in the order the
changes arrived. After a while it reads like a history: two requirements
state one rule, once with the owner and once without; another only
matters to someone upgrading from a feature that is gone; related rules
sit far apart. Agents still follow it, but a person can no longer take it
in at a glance.

This recipe rewrites such a spec to say what the project does today, as
one change with no code. It changes the spec's words, not the behavior,
and it keeps the history: the spec's History still links every change
that shaped it, and Git and the archived changes keep the old wording.

It needs nothing beyond proposing a change and building it on its own
branch (see the
[propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md)
and the
[requirements spec](https://github.com/mascah/grove/blob/main/grove/specs/requirements/spec.md)).
To move requirements to narrower specs without rewording them, see
[Split a maintained spec](split-a-maintained-spec.md) instead.

## Before you start

- **Find the spec.** Count each spec's requirements, and open the largest
  in an outline view:

  ```sh
  grep -c '^### Requirement' grove/specs/*/spec.md
  ```

  Look for requirements in pairs ("Do X" and "Do X unattended"), ones that
  only matter when upgrading, and ones that say only that something does
  not happen.
- **Start when no open change touches the spec.** Rewriting renames
  requirements, so a change that modifies one becomes stale and must be
  revised. `grove changes` lists what is open.
- **Spec text only.** If, on the way, you want the code to behave
  differently, that is another change. A requirement you retire may
  leave code behind that no spec agrees to; say so in the proposal and
  remove it later.

## Steps

**Terminal:** start a branch for it.

```sh
git switch -c reconcile-implement-spec
```

**Chat:** ask for a change that rewrites the spec.

```text
/grove-propose Reconcile the `implement` spec with current behavior: rewrite its requirements to say what the code and instructions do today, with no behavior change. In the proposal, map each requirement to where it goes, and list what the code does that no requirement says.
```

Ask the agent to read the code, and any instructions or prompts the
project ships, beside the spec. In the round, decide:

- **Which requirements merge,** and the name of each merged one. Merging
  each pair of a rule with its unattended or other-mode case into one
  requirement is where most of the words go.
- **Which requirements retire,** each with its reason: upgrade-only,
  stating only an absence, or covered by another spec. Search the other
  specs for each retired requirement's name; a Boundaries that points to
  it needs fixing.
- **The drift.** For each thing the code does that no requirement says:
  agree to it in the spec, or leave it as behavior nobody agreed.
- **Which names stay.** A requirement that keeps its name is MODIFIED;
  every other one is REMOVED, with where it went as its reason, and its
  rewrite ADDED.
- **A task for the text around the requirements**, with the new Purpose
  and Boundaries written out in the proposal, and for any other spec that
  names a retired requirement. Archive replaces requirements, not the text
  around them, and implementing touches no file a task does not name.

**Terminal:** commit the draft.

```sh
git add -A && git commit -m "Propose reconciling the implement spec"
grove check
```

**Chat:** build it on the same branch.

```text
/grove-implement <change-id>
```

The agent does the task, lists each scenario with the test that shows
it, then archives the change, which replaces the spec's requirements.

**Terminal:** commit, then read the result.

```sh
git add -A && git commit -m "Archive the reconciled implement spec"
grove check
grove review
```

Read the new spec in an outline view, and the proposal's map beside the
diff. Then merge: one merge approves the change and lands the rewritten
spec.

## What to expect

It was tried on Grove's own `implement` spec
([#112](https://github.com/mascah/grove/pull/112)): 24 requirements and
4,992 words became 8 requirements and 2,713 words, and one merge approved
it with the rewritten spec.

- **Most of the saving is merged pairs.** "Build unattended" restated
  five other requirements for the unattended case, and two more rules
  each had an unattended twin; four requirements for ways in (an approved
  change, a draft, a cycle, a cycle's member) became one.
- **The drift runs both ways.** The instructions did five things no
  requirement said, such as telling the owner how to try each task. Only
  reading the instructions beside the spec found them.
- **A retired requirement may have no code at all.** "No apply skill"
  was already covered by init removing any skill Grove no longer lists.
  Another, "Report rules for apply", left two lines of code that no spec
  now agrees to.
- **Another spec pointed at a retired requirement.** init's Boundaries
  named "No apply skill"; the change's task fixed it.
- **A link inside a requirement breaks when archived.** Archive copies
  requirement text word for word into the maintained spec, where a link
  written relative to the change folder no longer resolves. Name another
  spec in plain words inside a requirement; links in a REMOVED reason are
  fine, since they are not copied.
- **Most scenarios stay unverified.** 17 of 25 said what the agent does;
  a test showed the instructions print each rule, not that an agent
  follows it. Two rules had no test of their text at all.
- **Archive checks the new spec against the change, not the old spec.**
  Nothing compares the rewrite with what it replaces; the proposal's map
  from each old requirement to its new place is what review reads.
- **History needs nothing.** Its lines stay, and archive adds the
  reconcile change's own.
