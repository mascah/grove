# Run Grove in CI

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

Two jobs a workflow can do without an agent:

- On each pull request, `grove check` fails the build when the branch's
  changes or specs have a problem, before anyone approves them by merging.
- On each push to your main branch, `grove changes --json` keeps one GitHub
  issue per change, as [Push to a tracker](push-to-a-tracker.md) does from
  chat. Grove pushes nothing; the workflow does, with `gh`.

## What to copy

Save this as `.github/workflows/grove.yml`. It assumes your target branch is
`main`; if yours differs, change `main` under `branches` and in the
`git branch` step.

```yaml
name: Grove
on:
  pull_request:
  push:
    branches: [main]
env:
  GROVE_VERSION: v0.1.0
jobs:
  grove:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      issues: write
    concurrency: grove-${{ github.ref }}
    steps:
      - uses: actions/checkout@v7
        with:
          fetch-depth: 0
      - run: curl -fsSL "https://raw.githubusercontent.com/mascah/grove/$GROVE_VERSION/install.sh" | sh -s -- "$GROVE_VERSION"
      - run: echo "$HOME/.local/bin" >> "$GITHUB_PATH"
      # A pull request's checkout has no local main; grove check compares with it.
      - if: github.event_name == 'pull_request'
        run: git branch main origin/main
      - run: grove check
      - name: One issue per change
        if: github.event_name == 'push'
        env:
          GH_TOKEN: ${{ github.token }}
        run: |
          gh label create grove --force
          gh issue list --label grove --state all --limit 500 --json number,state,body > issues.json
          grove changes --json > changes.json
          root=$(jq -r .root changes.json)
          target=$(jq -r .target changes.json)
          jq -r '.changes[] | [.id, .state, .path, .title] | @tsv' changes.json |
          while IFS=$'\t' read -r id state path title; do
            marker="<!-- grove:$id -->"
            { git show "$target:$root/$path"; echo; echo "$marker"; } > body.md
            issue=$(jq -r --arg m "$marker" '[.[] | select(.body | contains($m))][0] // empty | "\(.number) \(.state)"' issues.json)
            if [ -z "$issue" ]; then
              url=$(gh issue create --title "$title" --body-file body.md --label grove)
              number=${url##*/} open=OPEN
            else
              read -r number open <<< "$issue"
              gh issue edit "$number" --title "$title" --body-file body.md
            fi
            if [ "$state" = approved ] && [ "$open" = CLOSED ]; then
              gh issue reopen "$number"
            elif [ "$state" != approved ] && [ "$open" = OPEN ]; then
              gh issue close "$number"
            fi
          done
```

The install script's URL names the version too, so a later change to the
script cannot change what an old pin runs. If your project pins Grove in
`mise.toml` ([Install](../install.md)), replace the two install steps with
`- uses: jdx/mise-action@v5` and drop `GROVE_VERSION`.

To only check, delete the last step and the `issues: write` permission. To
only keep issues, delete the `git branch` and `grove check` steps.

## What to expect

A pull request whose changes or specs have a problem fails with the lines
`grove check` prints, such as
`grove/changes/261020-cr2tr-monthly-budgets-in-report/.grove.yaml:5: unknown field "cycle"; ...`.

After a push to main, each change has one issue labelled `grove`, its title
the change's title and its body the change's main document ending with
`<!-- grove:<id> -->`. The issue of an archived or dropped change is closed,
and one closed by hand while its change is still approved is reopened.
Issues made by the chat recipe carry the same marker, so the workflow takes
them over rather than creating second ones.

Each push rewrites every issue's body, so edit the change, not the issue.
Upgrade Grove by changing the pinned version in a pull request of its own:
a release with breaking changes reaches the project only then.
