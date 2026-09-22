---
id: "G-050"
type: review
title: "G-025 shaping entrypoint: evidence"
status: current
formerly: "docs/reviews/2026-09-20-W-011-shaping.md"
work: ["G-025"]
created: "2026-09-21T03:02:56Z"
updated: "2026-09-21T21:11:16Z"
---

# G-025 shaping entrypoint: evidence

2026-09-20, branch `worktree-W-011` from `main` `70f19c5`.
[G-025](G-025-shaping-entrypoint.md) owns acceptance;
its [plan](G-049-shaping-entrypoint-plan.md) says what was to be shown.
Three kinds of evidence are kept apart below: source checks, observed agent
behavior in real harnesses, and simulation with the CLI. What nobody exercised
is listed at the end. The trials ran against the guide and adapters at
`5366ddd`; the one guide change they caused is named where it arose.

## Source checks

- `.claude/skills/grove-shape/SKILL.md` and `.agents/skills/grove-shape/SKILL.md`
  differ only in the Claude frontmatter (`disable-model-invocation`,
  `argument-hint`) and in how the request arrives (`$ARGUMENTS` against "the
  message that invoked this skill"). Both name `AGENTS.md` and
  `docs/work-shaping.md` and nothing else.
  `.agents/skills/grove-shape/agents/openai.yaml` is a copy of `grove-work`'s
  (`allow_implicit_invocation: false`).
- `go run ./cmd/grove check`: OK, 29 records. Every relative link
  in the nine Markdown files changed since `70f19c5` (`git diff --name-only
  70f19c5 -- '*.md'`) was resolved against the file system with a short Python
  script: 0 broken. The one in-page anchor, `#headless-shaping`, matches its
  heading. The reviewer repeated this independently with the same result.
- No Go source changed; the Go suite was not a relevant check and was not used
  as evidence for this work.

## Observed agent behavior

Disposable standalone clones of `5366ddd` under a session scratch directory,
`origin` removed, their branch renamed `main`, each with its own ID counter
(this repository's `.git/grove/next-ids` stayed at `W 29`, `D 5`). Fixture,
written for the trial by the author of the guide: on `main`, W-029 "Keep
finished work out of the default list", a rough proposal with no acceptance;
only on branch `worktree-W-030`, W-030 "Filter list output by status", whose
acceptance says output is unchanged without the flag. The topic given to both
harnesses states an unmade product choice: "I want finished work out of my way
in grove list. I have not decided whether done records should be hidden by
default or only when I pass a flag." Each clone was inspected with Git
afterwards rather than trusting the report.

- **Discovery, Claude Code 2.1.278:**
  `claude -p "/grove-shape --interaction headless"`, allowed tools Read, Glob,
  Grep. The init event lists `grove-shape` beside `grove-work` as skill and
  slash command. It read `docs/work-shaping.md` (through a read-only `cat`,
  which this harness permitted), returned the no-topic wait, and chose nothing
  from the backlog. 2 turns, nothing written.
- **Discovery, codex-cli 0.155.1:**
  `codex exec --sandbox read-only '$grove-shape --interaction headless'` read
  `.agents/skills/grove-shape/SKILL.md` and the guide and returned the same
  wait.
- **Headless shaping, Claude:** same topic with `--interaction headless`,
  tools Bash, Read, Edit, Write, Glob, Grep; 15 turns, 107 s. Order in the
  stream log: guide; W-029, `list`, `versions W-029`, worktrees and branches;
  W-030 read from its branch with `versions W-030` and `git show`; the `list`
  code in `internal/cli/cli.go` and the record model; only then
  `git worktree add -b worktree-shape-list-finished-work … main`, `new
  question`, body edits, revision-checked `update Q-002 … blocks=["W-029"]`,
  `check`, commit `53f174f`. Inspection: `main` clean and still `a0f46af`; the
  diff from `main` is W-029 and the new Q-002 only; W-029 still `proposed`; no
  code touched; no duplicate work record. Q-002 holds both options, evidence,
  and a recommendation marked as one. It reported that `relates_to` refused
  W-030 because that record is not in the checkout, and named it in prose
  instead. **The guide now says so** (step 2). It also observed that hiding by
  default would affect this guide's own use of `list`.
- **Headless shaping, Codex:** `codex exec --sandbox workspace-write` with the
  same request. Same order: guide, brief, branches and worktrees, `list`,
  `versions`, W-029, the code, W-030 from its branch, `context W-029
  --interaction headless`, then `git worktree add -b
  worktree-shape-list-finished … main`, `new question --slug`, `update`,
  `check`, commit `8ca85d4`. Inspection: `main` clean at `a0f46af`; diff is
  W-029 and Q-002 only; Q-002 `open` with `blocks: ["W-029"]`; W-029
  `proposed`. It retitled W-029 and set `relates_to` to G-042 and G-043, a
  judgment the guide permits and the owner can reject in review. It set
  `GOCACHE` under `/private/tmp` because the sandbox refused the default.
- **Unchanged wait, Claude:** the same command again in the same clone.
  4 turns, 29 s, read-only commands, no new branch, question, or commit
  (branches and worktrees identical afterwards). It returned the same wait on
  Q-002. Not repeated with Codex.

What these show: both harnesses discover the adapter from an explicit headless
invocation and reach the same guide; a fresh agent refined the overlapping
proposal instead of duplicating it, found the other-branch work through
`versions`, turned the unmade choice into a durable question and a returned
wait on an isolated branch, and promoted, implemented, launched, and merged
nothing. What they do not show: the fixture and topic were written by the
guide's author, the topic made the missing choice explicit, and one run per
harness is not a rate.

## Simulation with the CLI (no agent)

In the fixture clone, with a built binary and an absolute `--project`:

- Stale revision: after appending to W-029's body, `update W-029 --expect OLD
  --set priority=2` exits 1 with "W-029 changed since the expected revision;
  its current revision is sha256:d3c4a84e…". The file was not changed by the
  refusal. No agent was observed meeting this refusal.
- Other-branch relationship: with the current revision,
  `--set 'relates_to=["W-030"]'` on `main` exits 1 with "relates_to: unresolved
  target W-030".

## Not exercised

- **Acceptance 2, a real requirements conversation.** This session had the
  guide's author and no requirements to discuss; inventing one would
  manufacture records. It needs the owner running `/grove-shape TOPIC` on a
  real idea in a fresh interactive session, and the owner's judgment of the
  result.
- Interactive discovery: `/grove-shape` typed in a Claude Code session and
  `$grove-shape` in the Codex TUI; whether `agents/openai.yaml` prevents
  implicit invocation.
- Interactive write placement (step 4): asking before using another
  assignment's checkout, and committing only on agreement.
- A headless run that meets another session's uncommitted version, a stale
  revision, or an unwritable branch; Codex's unchanged-wait rerun.
- Moving proposals between separate clones with colliding IDs.

## Owner's disposition

2026-09-20: the owner closed G-025 without the acceptance-2 conversation, for
the reasons in the record, and will exercise `/grove-shape G-037` afterwards.
Everything under "Not exercised" above was still unexercised at closure.

## Independent review

Three rounds, the first on `218d03e` against `70f19c5`, by a fresh read-only reviewer agent
inside the implementing Claude session: independent of the author's context,
not a person and not another harness. It confirmed that Outcome, Selected scope
and Acceptance are unchanged, that every CLI form and status in the guide
matches `--help` and the record model, that the quoted refusals are real
strings, and that the evidence does not overclaim. Findings:

| Finding | Disposition |
| --- | --- |
| Consequential: the record's Next told the owner to shape "on this branch", which the guide's step 4 forbids for an execution checkout | Fixed: Next now says merge first, then shape from the main checkout |
| Minor: link-check provenance pointed at the record, which held no command | Fixed above |
| Minor: this section referred to findings that did not exist yet | Fixed: findings live here |
| Minor: headless base rule did not say what to do when the default base lacks the records | Fixed in step 4, parallel to the work guide's step 3 |
| Minor: G-023's Next still says G-025 is not started | Not changed: G-023 is a done record outside this assignment, and the brief and G-025 own current state |

Round 2, the same reviewer on `218d03e..8ae82a9`: the first fix resolved its
objection and the headless base rule is a faithful parallel of the work guide,
but the record's new Next (merge, then shape) inverted the brief's (shape, then
integrate). Consequential; fixed by reordering the brief's Next. Two minor
wording points, naming the owner as the actor inside the Next sentence and this
section's round count, were also fixed.

Round 3, on `8ae82a9..37546a6`: no objections; the two Next statements agree,
`check` passes and links resolve. The reviewer accepted `37546a6`. Only this
paragraph and the round count were written afterwards.
