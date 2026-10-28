# Split a maintained spec

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

A maintained spec named for a whole program or screen, such as `app`,
gains requirements with every feature, because every feature belongs to
it. An agent working on one part of it reads all of it, and changes built
side by side all edit the same file. This recipe moves requirements out of
such a spec into narrower specs, one per area, as a change with no code.
Nothing is reworded, so nothing agreed changes; archive checks the moves,
and each spec's History keeps the reasons for its requirements one link
away.

It needs nothing beyond proposing a change and building it on its own
branch (see the
[propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md)
and the
[requirements spec](https://github.com/mascah/grove/blob/main/grove/specs/requirements/spec.md)).

## Before you start

- **Find the areas.** Group the requirements that changes have touched
  together. The spec's History links each change, and each archived
  change's `specs/<capability>/spec.md` names the requirements it added or
  modified.
- **Name the new specs with a shared prefix**, such as `app-editor` and
  `app-graph`, so `grove specs` lists them together.
- **Keep the old spec** for what is left, such as the window and settings.
  Its Boundaries can then say where each area went, and a change that still
  adds to it lands in a spec that exists.
- **Start when no open change touches the spec.** A change that modifies or
  removes a requirement you move becomes stale and must be revised; one
  that only adds to the old spec gets no warning.
- **Move only.** Rewording, merging or reordering requirements changes what
  is agreed. Do that in a later change, where review can see it.

## Steps

**Terminal:** start a branch for it.

```sh
git switch -c split-app-spec
```

**Chat:** ask for a change that only moves requirements.

```text
/grove-propose Split the `app` spec by area into app-* specs, moving requirements word for word, with no behavior change
```

In the round, decide:

- **Which requirement goes where**, and which stay.
- **Each new spec's title, Purpose and Boundaries**, written out in the
  proposal, so you review them before archive writes them.
- **A task for the old spec's Purpose, Boundaries and known
  discrepancies.** Archive moves requirements, not the text around them,
  and implementing touches no file a task does not name.
- **History:** each new spec lists the changes behind the requirements it
  receives, then this change; the old spec keeps its lines.
- **Optionally, a rule in `grove.yaml`** so new requirements go to the
  right spec, as a task in the same change:

  ```yaml
  rules:
    documents:
      change-spec:
        - Put an app requirement in the `app-*` spec for its area; a new area of the app gets its own `app-<area>` spec.
  ```

The change spec for the old spec lists each moved requirement under
REMOVED, with a reason such as "Moved word for word to `app-graph`." Each
new spec's change spec lists its requirements under ADDED. Ask the agent to
copy each requirement block exactly, scenarios and all, not to write it
again.

**Terminal:** commit the draft.

```sh
git add -A && git commit -m "Propose splitting the app spec"
grove check
```

**Chat:** build it on the same branch.

```text
/grove-implement <change-id>
```

The agent does the tasks, then archives the change, which writes the new
specs and takes the moved requirements out of the old one.

**Terminal:** commit, then check that every requirement moved unchanged.

```sh
git add -A && git commit -m "Archive the split of the app spec"
git diff --color-moved main -- grove/specs
grove review
```

In the diff, requirement lines show as moved. Only the new specs'
headings, Purpose, Boundaries, Known discrepancies and History show as
added. Then merge: one merge approves the change and lands the new specs.

## What to expect

It was tried on a project whose `app` spec had 37 requirements and 7,168
words: one change moved 32 of them into six `app-*` specs, and one merge
approved it with the new specs. The messages below come from a trial on a
copy of that project.

- **Archive writes the new specs from the proposal.** Each new spec's
  Purpose and Boundaries came out as the proposal wrote them, and
  `--color-moved` showed every requirement line as moved.
- **Archive checks the specs against the change.** One word changed in a
  moved requirement gave `requirement "Graph hover" is not as <change-id>
  gives it`, and a requirement left in the old spec gave `requirement
  "Graph view" is still here though <change-id> removes it`.
- **Archive does not check the change against the old spec.** If the
  change's ADDED text differs from the requirement it moves, archive
  accepts it. `git diff --color-moved` shows it; GitHub's pull request diff
  does not mark moved lines.
- **A change open since before the split goes stale.** For one that
  modified a moved requirement, `grove check` gave `requirement "Graph
  hover" was removed since base <commit>; revise the change with propose`.
  The message does not say where the requirement went; the old spec's
  Boundaries does.
- **Archived changes keep citing the old spec.** Their change specs and
  proposals name its path as it was at each change's `base`. Leave them
  as they are.
- **History links work from any spec.** Each line links
  `../../changes/archive/<folder>/`, and archive does not compare History,
  so lines copied into a new spec need no change.
- **The moves need no task.** Archive does them from the change specs.
- **Implement still lists every moved scenario.** With no code to check,
  it gave each the check recorded by the archived change that added it.
  Saying so in `tasks.md` tells it where to look.
- **`grove specs` shows the new specs once your target branch has them.**
  It and the project index read your local `main`, not the branch you are
  on, so update `main` after the merge.
