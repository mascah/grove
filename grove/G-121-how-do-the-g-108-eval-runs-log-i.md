---
id: "G-121"
type: question
title: "How do the G-108 eval runs log in to the clean config directory?"
status: open
created: "2026-09-24T00:55:29Z"
updated: "2026-09-24T00:55:54Z"
blocks: ["G-108"]
---

## Question

[G-118](G-118-what-mandate-should-the-g-108-pa.md) is resolved: opus 5.5,
5 runs per case, $5 per run, permission mode `auto`, the fixture approved as
built. Its last item, where the clean `CLAUDE_CONFIG_DIR` lives and how it
logs in, has no answer, and it is the one input a headless session cannot
supply: the runner refuses the owner's own configuration on purpose, and a
new directory has no login.

Observed 2026-09-23 on `worktree-G-108` at `433e338`, Claude Code 2.1.281:
`CLAUDE_CONFIG_DIR=~/.cache/grove-evals/claude claude -p` with `--model
claude-opus-5-5 --permission-mode auto` and a $0.10 cap returned `Not
logged in · Please run /login` at no cost. `/login` needs a browser, so the
runs cannot start until the owner does one of:

1. **Log in once** (recommended, matches G-118's recommendation and what a
   preview user gets): run `CLAUDE_CONFIG_DIR=~/.cache/grove-evals/claude
   claude`, then `/login`, then quit. The directory then holds only that
   login; the runner still refuses it if a `CLAUDE.md`, skills, plugins or
   settings appear there.
2. **Export a credential** in the shell that starts the runs:
   `ANTHROPIC_API_KEY` (API billing, not the subscription) or
   `CLAUDE_CODE_OAUTH_TOKEN` from `claude setup-token`; both pass through the
   runner's scrubbed environment.

Then assign `/grove-work G-108` again on `worktree-G-108`. The successor
runs, with exactly G-118's values:

```sh
python3 evals/run.py run --runs 5 --budget 5 --model claude-opus-5-5 \
    --permission-mode auto --config-dir ~/.cache/grove-evals/claude --out DIR
```

Cap: 2 cases × 5 runs × $5 = $50, printed before the first run.

Who can answer: the owner, by doing one of the two and setting this record
`resolved`, naming which. A different directory is fine; say which.

## Next

Waits for the owner.
