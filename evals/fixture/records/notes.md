## Decision

Owner, 2026-08-21: the owner writes notes in a task file, below its title
line, with any editor. Every command that rewrites a task file keeps
everything below the title exactly as it was.

Observed the same day: `write` in `tasks.py` rewrites the whole file from
the frontmatter and the title alone, so `tasks done` and `tasks drop`
erased the notes of the task they closed. Until `write` keeps the rest of
the file, no new command may rewrite a task file through it.

## Alternatives

- Notes in a separate file beside each task: two files per task, and the
  brief keeps a task to one file.
- Notes in a frontmatter field: unreadable in an editor for anything
  longer than a line.

## Reconsideration

When a command must change a task's body itself.
