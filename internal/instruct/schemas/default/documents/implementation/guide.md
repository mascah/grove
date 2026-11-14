Write `implementation.md` in the change folder and keep it current: each
task's status (not started, in progress, done, blocked) and its evidence
(what changed, the checks run and their result). A `Check:` line under a
task tells the reviewer a choice made or what could not be seen, and why,
or a blocked task's question; an `Answer:` line quotes the owner's answer.
It is not part of what was
approved and does not amend it. Recording progress never edits `tasks.md`:
a task's status lives here, under a heading with its ID. Grove's `archive`
reads only the `Status:` lines, so keep them exactly as shown. Under a
task, `` Also changed: `<path>` — <why> `` records a file the implementing
agent changed beyond its `Files:` and `Delegated:` lines, and why, one per
file.
Before finishing, list under `## Scenarios` each scenario of the change
specs with the check that shows it and its result, or `unverified` and why.
