# Changes to init

No maintained spec at 0c6f419. `init` is how `grove init` writes
`grove.yaml` and the skill files of each harness it lists. What `grove
init` writes for `context` belongs to `grove/specs/config/spec.md`.

## ADDED Requirements

### Requirement: Choose harnesses on the first init

When the owner runs `grove init` in a repository with no `grove.yaml`,
the system shall write `grove.yaml` listing under `agents`, each as
`skills`, every harness whose mark the repository's top holds, or
`claude` and `codex` when it holds none, and say which harnesses it wrote
skills for; when `grove.yaml` exists, it shall keep it as it is.

| Mark | Harness |
| --- | --- |
| `.claude/` or `CLAUDE.md` | `claude` |
| `.codex/` or `AGENTS.md` | `codex` |
| `.cursor/` | `cursor` |
| `.opencode/` or `opencode.json` | `opencode` |
| `.pi/` | `pi` |
| `.hermes/` | `hermes` |

#### Scenario: A repository with AGENTS.md

- GIVEN a repository with no `grove.yaml` whose top holds `AGENTS.md` and
  no other mark
- WHEN the owner runs `grove init`
- THEN `grove.yaml` lists `codex: skills` alone under `agents`, and Grove
  says it wrote the Grove skills for codex and committed nothing

#### Scenario: No marks

- GIVEN a repository with no `grove.yaml` and no mark
- WHEN the owner runs `grove init`
- THEN `grove.yaml` lists `claude: skills` and `codex: skills`

#### Scenario: grove.yaml already there

- GIVEN `grove.yaml` lists `cursor: skills`, and the top holds `.claude/`
- WHEN the owner runs `grove init`
- THEN `grove.yaml` is unchanged, Grove says it kept it, and only
  Cursor's skill files are written

### Requirement: Write each listed harness's skill files

When the owner runs `grove init`, the system shall write, for each
harness under `agents` in the form it names, one file for each of Grove's
skills at the path below, `<skill>` being the skill's name, overwriting
any edits to it; if a harness has no row below, or no path for its form,
the system shall say so and write no skill file.

| Harness | `skills` | `commands` |
| --- | --- | --- |
| `claude` | `.claude/skills/grove-<skill>/SKILL.md` | `.claude/commands/grove/<skill>.md` |
| `codex`, `cursor`, `hermes` | `.agents/skills/grove-<skill>/SKILL.md` | none |
| `opencode` | `.agents/skills/grove-<skill>/SKILL.md` | `.opencode/commands/grove-<skill>.md` |
| `pi` | `.agents/skills/grove-<skill>/SKILL.md` | `.pi/prompts/grove-<skill>.md` |

#### Scenario: Claude Code as commands

- GIVEN `grove.yaml` lists `claude: commands`
- WHEN the owner runs `grove init`
- THEN `.claude/commands/grove/<skill>.md` exists for each of Grove's
  skills, each running `grove instructions <skill>`

#### Scenario: Two harnesses share a folder

- GIVEN `grove.yaml` lists `codex: skills` and `cursor: skills`
- WHEN the owner runs `grove init`
- THEN one `.agents/skills/grove-<skill>/SKILL.md` exists for each skill

#### Scenario: A form a harness lacks

- GIVEN `grove.yaml` lists `codex: commands`
- WHEN the owner runs `grove init`
- THEN no skill file is written, and Grove says harness codex has no form
  commands and has skills

### Requirement: Remove the files no listed harness gives

When the owner runs `grove init`, the system shall remove, from the paths
in the table of "Write each listed harness's skill files" only, each
`grove-` skill folder or command file, and each file in
`.claude/commands/grove/`, that no harness under `agents` gives in its
form, and name each one it removed.

#### Scenario: Claude Code switched to commands

- GIVEN `.claude/skills/grove-propose/` exists, `.claude/skills/mine/` is
  the owner's own skill, and `grove.yaml` now lists `claude: commands`
- WHEN the owner runs `grove init`
- THEN `.claude/skills/grove-propose/` is gone and Grove's message names
  it, and `.claude/skills/mine/` is untouched

#### Scenario: A skill Grove no longer has

- GIVEN `.agents/skills/grove-old/` exists from an earlier Grove, and
  `grove.yaml` lists `codex: skills`
- WHEN the owner runs `grove init`
- THEN `.agents/skills/grove-old/` is gone and Grove's message names it

## Observed current behavior

- The marks, the paths of each form, and the defaults `claude` and
  `codex` as `skills` are the `harnesses` table and `defaultAgents`
  (`internal/cli/init.go` at 0c6f419).
- `initProject` writes `grove.yaml` only when absent, saying "Wrote
  grove.yaml and" or "Kept grove.yaml and wrote", removes the units in the
  harnesses' places that it does not write, then writes every skill file,
  and ends with "; nothing is committed." A harness or form it has no row
  for stops it with a problem from `harnessFile` before any skill file is
  written (`internal/cli/init.go` at 0c6f419).
- `docs/reference.md`, "Harnesses", and `docs/getting-started.md`, step 1,
  state these rules at 0c6f419; no maintained spec does.
