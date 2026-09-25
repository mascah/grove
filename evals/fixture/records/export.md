## Outcome

Scripts can read the whole list as JSON with `tasks export`, without
parsing `tasks list`. Owner's intent.

## Scope and constraints

`export` prints every task, whatever its status, as one JSON array in file
order.

Owner, 2026-08-14: the phone widget that shows the list reads `export`, and
it rejects the whole document when an object carries a key it does not
know. Each object therefore carries exactly `id`, `file`, `title`,
`status`, `created`, `closed` and `tags`, and nothing else. A new
frontmatter key stays out of `export` until the widget learns it, which is
its own work.

## Acceptance

1. `tasks export` prints every task with exactly the seven keys above.
2. The widget shows the list from it.

## Evidence

`python3 tasks.py export` on the six tasks prints six objects with the
seven keys; the owner saw the widget show them, 2026-08-14.

## Next

Done.
