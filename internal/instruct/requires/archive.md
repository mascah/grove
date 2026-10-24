- First run Grove's `check`. If it reports a requirement of this change as
  rewritten or removed since its base, stop and report it: the change must
  be revised with propose and merged again first.
- A change is finished when every task in `tasks.md` has `Status: done` in
  `implementation.md`. Otherwise ask the owner whether it is dropped. A
  dropped change gets `Status: dropped` under the title of
  `implementation.md`, created from its template below if there is none,
  with the owner's reason in their words; it skips the finishing work.
- For a finished change, first tell the owner each scenario the
  `## Scenarios` section of `implementation.md` marks `unverified`; it does
  not stop the archive. Then do the finishing work "How to archive"
  describes.
- Then run Grove's `archive <change-id>`. It moves the change's folder to
  `changes/archive/` and refuses a change that is neither finished nor
  dropped. For a finished change it also refuses, printing
  each requirement as `path:line: message`, when a maintained spec does not
  hold what this branch and the change give it; fix the spec as the
  messages say and run `archive` again.
- Tell the owner what changed and where the change went. The owner merges
  this branch, which lands the work and the archived change together.
