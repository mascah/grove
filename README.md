# Grove

Grove keeps a written record of what your software should do, and helps you
and your coding agent agree on a change before it is built. Specs and changes
are Markdown files in your repository. Merging a change into your main branch
is what approves it.

Grove has two halves:

- **In your terminal**, the `grove` command sets up a project, lists and
  checks specs and changes, shows what a merge would approve, and opens a
  reader.
- **In your agent's chat** (Claude Code or Codex), Grove's skills explore an
  idea, propose a change, apply it and archive it. Each skill asks `grove`
  for instructions and follows them in the conversation you are already in.

Grove never starts an agent and never commits. You run your agent and you
commit.

## The short version

```text
terminal  grove init                            set up, then commit
terminal  git switch -c my-change
chat      /grove-explore <what you want>        talk it through; nothing is written
chat      /grove-propose                        writes the change
terminal  grove review, then merge              the merge approves the change
chat      /grove-apply <change-id>              builds it, on a new branch
chat      /grove-archive <change-id>            updates the specs, files the change away
terminal  merge
```

In Codex, type `$grove-explore` instead of `/grove-explore`, and so on.

## Docs

- [Getting started](docs/getting-started.md): set up and take one change from
  idea to updated spec.
- [Concepts](docs/concepts.md): specs, changes, approval and archive.
- [Reference](docs/reference.md): every command, skill and setting.
