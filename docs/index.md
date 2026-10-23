# Grove

Grove keeps a written record of what your software should do, and helps you
and your coding agent agree on a change before it is built. Specs and changes
are Markdown files in your repository. Merging a change into your main branch
is what approves it.

Grove has two halves:

- **In your terminal**, the `grove` command sets up a project, lists and
  checks specs and changes, shows what a merge would approve, opens a
  reader, and can run apply unattended.
- **In your agent's chat** (Claude Code, Codex, Cursor, OpenCode, Pi or
  Hermes), Grove's skills explore an idea, prototype it on a throwaway branch,
  propose a change, implement it
  with you or apply it without you, and archive it. Each skill asks `grove` for instructions and follows them in
  the conversation you are already in.

Grove never commits. It starts an agent only when you run `grove run`,
which launches one Claude Code session to build what is approved; otherwise
you run your agent. You commit, or your agent does when a rule asks.

## The short version

```text
terminal  grove init                            set up, then commit
terminal  git switch -c my-change
chat      /grove-explore <what you want>        talk it through; nothing is written
chat      /grove-prototype <question>           optional: try it on a throwaway branch
chat      /grove-propose                        writes the change
terminal  grove review, then merge              the merge approves the change
chat      /grove-implement <change-id>          builds it with you, updates the specs, files it away
terminal  merge
```

Or, once several changes are approved, `/grove-plan` puts them in a cycle,
`/grove-apply <cycle-id>` builds them all in order, and you merge once.
Mark the cycle `unattended` and `grove run` builds it for you.

Some agents type skills another way, such as `$grove-explore`; the
[Reference](reference.md) lists how each is typed.

## Where to go next

- [Install](install.md) the `grove` command.
- [Getting started](getting-started.md): set up and take one change from
  idea to updated spec.
- [Concepts](concepts.md): the problem Grove solves, and what specs,
  changes, approval and archive mean.
- [Coordinating changes](coordinating.md): changes that wait on others,
  stale plans, and the project as data.
- [Planning cycles](planning.md): group approved changes, say which need
  a person, and mark a cycle for `grove run`.
- [Running apply unattended](running.md): `grove run` builds a cycle you
  marked in one Claude Code session, on its own or from a timer.
- Guides: recipes others have tried, such as
  [pushing changes to a tracker](guides/push-to-a-tracker.md),
  [running Grove in CI](guides/run-grove-in-ci.md),
  [running apply on a timer](guides/run-apply-on-a-timer.md) and
  [keeping a glossary and decision records](guides/glossary-and-decisions.md).
- [Reference](reference.md): every command, skill and setting.
