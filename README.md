# Grove

Grove keeps a written record of what your software should do, and helps you
and your coding agent agree on a change before it is built. Specs and changes
are Markdown files in your repository; merging a change into your main branch
approves it. Before proposing, `/grove-prototype <question>` in your
agent's chat tries an idea out on a throwaway branch.

Your agent builds an approved change with you when you type
`/grove-implement <change-id>` in its chat, asking you what the change leaves
open. Or `grove run` does it unattended: `/grove-plan` puts approved
changes in a cycle, and once you mark the cycle `unattended`, `grove run`,
in a checkout you keep for the project, starts one time-limited Claude
Code session that builds it on a new branch, then waits until you merge. See
[Running apply unattended](docs/running.md).

To install the `grove` command, see [Install](docs/install.md).

The documentation is at <https://mascah.github.io/grove/>, built from
[`docs/`](docs/index.md).

## License

MIT. See [LICENSE](LICENSE).
