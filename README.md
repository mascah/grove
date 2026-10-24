# Grove

Grove keeps a written record of what your software should do. It helps you
and your coding agent agree on a change before it is built. Specs and
changes are Markdown files in your repository, and merging a change into
your main branch approves it.

One change goes like this:

1. **Propose.** Type `/grove-propose` in your agent's chat. The agent asks
   you the decisions it would otherwise assume, then writes the change.
2. **Review and merge.** Read the change and merge it. The merge approves
   it.
3. **Implement.** Type `/grove-implement <change-id>`. The agent builds the
   change with you, asks you what the change leaves open, and updates the
   specs.
4. **Merge.** Merge the code, the updated specs and the archived change
   together.

A small change can be proposed and built on one branch, so one merge
approves it with its code.

To install the `grove` command, see [Install](docs/install.md). Then
[Getting started](docs/getting-started.md) takes one change through.

## Doing more

Before proposing, `/grove-explore` talks an idea through, and
`/grove-prototype <question>` tries it out on a throwaway branch. Once
several changes are approved, `/grove-plan` puts them in a cycle to build
together. `grove run` builds a cycle unattended, in one time-limited
Claude Code session. See [Running unattended](docs/running.md).

The documentation is at <https://mascah.github.io/grove/>, built from
[`docs/`](docs/index.md).

## License

MIT. See [LICENSE](LICENSE).
