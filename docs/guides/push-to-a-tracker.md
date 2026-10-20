# Push to a tracker

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

Some teams plan from a tracker, such as GitHub issues, rather than from the
repository. This recipe keeps one issue per change, so the tracker shows
what is approved and what is done, while Git stays the source of truth.

Grove pushes nothing to GitHub or any other tracker. Your agent does, from
`grove changes --json` (see [Coordinating changes](../coordinating.md),
"The project as data"), with the tools it already has. Each run replaces
issue bodies whole, so edit the change, not the issue, and discuss in the
issue's comments.

## What to copy

**Chat**, on your main branch after a merge:

```text
Run `grove changes --json`. For each change, create or update one GitHub
issue with the label `grove`: the title is the change's title, and the body
is the file at <root>/<path> as main holds it (`git show main:<root>/<path>`),
ending with the line <!-- grove:<id> -->. Find an existing issue by that
marker, open or closed, and never create a second. Close the issue of each
archived or dropped change, and reopen one whose change is approved. Then
tell me what you created, updated and closed.
```

With GitHub's `gh`, the agent's commands look like this:

```sh
gh label create grove --force
gh issue list --label grove --state all --limit 500 --json number,title,body,state
git show main:grove/changes/261019-r2w7k-add-ledger-delete/proposal.md > body.md
echo '<!-- grove:261019-r2w7k -->' >> body.md
gh issue create --title "Add ledger delete ID" --body-file body.md --label grove
gh issue edit 12 --body-file body.md
gh issue close 9
```

For another tracker, name it in the prompt instead of GitHub; the agent
uses whatever tools it has for that tracker.

## What to expect

The first run creates one issue per change, each labelled `grove`, its body
the change's proposal ending with `<!-- grove:<id> -->`. The issues of
archived and dropped changes are closed.

Run the prompt again after each merge into your main branch, or let a
workflow do it on each push: see [Run Grove in CI](run-grove-in-ci.md). An
issue whose change was archived since the last run closes, and the markers
keep the agent from creating a second issue for a change it has seen.
