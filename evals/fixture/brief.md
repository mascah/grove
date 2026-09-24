# Brief

## Purpose

`tasks` keeps one person's to-do list as plain Markdown files in a Git
repository, so the list works with any editor and its history is Git's.
This is the owner's intent; the owner is the only user.

## Constraints

- A task is one Markdown file under `tasks/`, named `NNN-slug.md`, with YAML
  frontmatter: `status`, `created`, `closed` (a date) and `tags`. It stays readable and editable without the tool.
- A task's status is `open`, `done` (it was completed) or `dropped` (the
  owner decided not to do it).
- A task file never moves and is never renamed once created: history and
  links follow the path.
- No index, cache or database beside the task files; every command reads the
  files.
- `tasks list` output is for people and may change between versions. Scripts
  read `tasks export`, whose JSON is stable.

## Conventions

- A filter flag may repeat. Repeated values of one flag match a task that has
  any of them; different flags must all match. `--status` already works this
  way.
- A command that hides some tasks by default ends its output with how many it
  hid, and offers `--all` to show everything.
- Tags are lowercase words; `tasks add` lowercases them.

## Direction

Finding things comes before moving them: listing and filtering improve before
any sync between machines.
