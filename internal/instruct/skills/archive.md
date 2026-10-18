- A change is finished when every task in `tasks.md` has `Status: done` in
  `implementation.md`. Otherwise ask the owner whether it is dropped. A
  dropped change gets `Status: dropped` under the title of
  `implementation.md`, created from its template below if there is none,
  with the owner's reason in their words; it skips the finishing work.
- For a finished change, first do the finishing work "How to archive"
  describes.
- Then run Grove's `archive <change-id>`. It moves the change's folder to
  `changes/archive/`, commits nothing, and refuses a change that is neither
  finished nor dropped.
- Tell the owner what changed and where the change went. Commit only when
  the owner or the project's rules ask. The owner merges this branch, which
  lands the work and the archived change together.
