- Simple work is usually one task. Add another only for a part that could be
  delivered and accepted on its own.
- A change that only records behavior the code already has has no tasks.
- A task states the result wanted and the acceptance expectations someone else
  could check. It refers to requirements instead of repeating them.
- A task names, under `Files:` after `Result:`, the repo-relative paths it
  creates, modifies or tests, one backticked path per bullet, a folder
  ending in `/`; Grove reads them. On an `Interfaces:` line it names the
  functions, types or commands it produces that other tasks or callers use.
  A task says what and where, not how: no code, commands or steps.
- A task may have a `Delegated:` line, right after `Files:`, naming choices
  the implementer may make without asking and paths it may change beyond
  `Files:`, only where the owner agreed in the round. Implementing with the
  owner honors it; implementing unattended ignores it. Grove does not read
  it.
- Never renumber or reuse a task ID. A renamed or reordered task keeps its ID,
  a new task takes the next unused number, and a split task gets new IDs with
  a note naming the one it replaces.
- When the owner wants the README, a changelog or other files outside the
  code updated, that is a task. Applying a change touches nothing no task names.
- No ordering, no execution plan, no status and no checkboxes.
