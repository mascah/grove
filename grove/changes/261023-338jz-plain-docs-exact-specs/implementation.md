# Implementation: Plain docs, exact specs

<!-- Progress on change 261023-338jz, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Move the docs' exact rules into the maintained specs

Status: in progress

### T-02 — The README and the home page lead with the attended path

Status: done

- `README.md`: the four steps (propose, review and merge, implement,
  merge), then one sentence that a small change takes one merge, then
  Install and Getting started. Explore, prototype, cycles, `/grove-plan`
  and `grove run` moved under "Doing more".
- `docs/index.md`: "Two halves" no longer says the command runs a cycle
  or that skills implement without you; "The short version" is now "One
  change", without the plan and run lines, followed by the one-merge
  sentence; cycles, `--unattended` and `grove run` are under "Doing more",
  before "Where to go next".
- Checks: above "Doing more", neither file matches `cycle`, `unattended`
  or `grove run` (grep). Every relative link in both files names a file
  that exists, and "Where to go next" still links each of the 11 pages
  `mkdocs.yml` lists besides the home page.

### T-03 — A plain-language pass, page by page

Status: not started

## Scenarios

## Reconciliation
