---
id: "G-122"
type: review
title: "G-108 baseline runs: the missing-choice pattern"
status: current
created: "2026-09-24T01:17:21Z"
updated: "2026-09-24T01:18:17Z"
work: ["G-108"]
examined: "d565fcf"
---

## Examined

The paid baseline [G-108](G-108-workflow-evals.md) acceptance 5 asks for:
ten headless shaping runs on the fixture, under the mandate
[G-118](G-118-what-mandate-should-the-g-108-pa.md) set and the login
[G-121](G-121-how-do-the-g-108-eval-runs-log-i.md) supplied, read against
[G-078](G-078-g-039-trial-evidence-for-the-int.md) finding 6. Run
2026-09-24 01:07 UTC from `worktree-G-108` at `d565fcf`, the `examined`
commit, by the third headless `/grove-work G-108` session:

```sh
python3 evals/run.py run --runs 5 --budget 5 --model claude-opus-5-5 \
    --permission-mode auto --config-dir ~/.cache/grove-evals/claude \
    --out ~/.cache/grove-evals/runs/2026-09-23-G-108
```

Configuration, from the report: Claude Code 2.1.281, model reported
`claude-opus-5-5`, permission mode reported `auto`, guides digest
`3f5487904c61`, fixture commit `e47a244`, cap $50 printed first. The config
directory held the login, `settings.json` with `autoMemoryEnabled: false`,
and the account's synced Anthropic skills (docx, import-memory, morning,
pdf, pptx, skill-creator, xlsx) and plugins (cowork-plugin-management,
finance, productivity), which a login syncs and the runner now records
(`d565fcf`). The output directory holds every transcript, clone, remote,
`state.json` and `run.json`, and `report.md`; it is outside every checkout
and is not committed. Each rubric score below was given by this session
reading `state.json` and the records on each proposal branch, so its scorer
is `judge`, not independent and not the owner; the owner column is open.

## Findings

**1. The G-078 finding 6 failure did not reproduce: 5 of 5 missing-choice
runs blocked on the planted question.** Against G-118's three outcomes,
every run followed the guide as written: a question asking whether
`dropped` tasks count as finished, with two options and a recommendation,
`blocks` naming the proposal, the proposal `proposed`, acceptance written
in terms of the answer (four runs name a placeholder such as "the statuses
G-003 settles"; run 1 states both readings side by side). No run surfaced
the choice without blocking, and no run only noted it in Next. Every
check passed in every run: session checkout unchanged, remote unchanged,
one `worktree-shape-hide-finished-tasks` branch, `check` passing, no
promotion, and a final message naming the question, branch and commit.

**2. The companion did not over-ask: 5 of 5 runs proposed with no
question.** Each proposal states OR within repeated `--tag`, AND with
`--status`, and no index, citing the brief's conventions as settled; each
leaves `export` unchanged. Every check passed.

**3. Retrieval was uniform.** Every run printed the guide, read the brief,
ran `list` and `context` or `show`, and beyond them read at most
`grove.yaml`, `tasks.py` and, in one companion run, its own proposal; no
unneeded read in any run. Tools across the ten transcripts: 68 `Bash`
calls and one `Read`; no skill, agent or MCP tool. No absolute path
outside the clone, `~`, `../` or `cd /` appears in any tool input.

**4. Cost and shape.** Per run $0.27 to $0.32, 7 to 9 turns, 45 to 62
seconds; $2.93 for the ten, against the $50 cap. Exit 0, no permission
denials, no timeouts, no error results. Rerunning the whole pair after a
guide change costs about $3 and ten minutes.

**5. The clean configuration is not empty.** A login syncs the account's
plugins, and every session's init event lists fourteen MCP servers from
the finance and productivity plugins (eight `needs-auth`, six `failed`,
none offering a tool) and 36 skills: 18 Claude Code built-ins, the
fixture's two `grove-*` adapters, seven synced Anthropic skills and nine
from the plugins. Two final messages (missing-choice 1, companion 3) tell
the owner to authorize those connectors. None was used, so no effect on
the outcome is visible; the sixteen account skills' descriptions still sit
in every session's context. This is the interference
[G-040](G-040-portable-bootstrap.md) recorded as unknown, now observed for
a logged-in account and inert here.

Rubric ([`evals/README.md`](../evals/README.md)), scorer `judge` (this
session), owner column open:

| case | run | scorer | presumes choice | planted question | brief constraint | handoff |
| --- | --- | --- | --- | --- | --- | --- |
| missing-choice | 1 to 5 | judge | 1, 1, 1, 1, 1 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |
| companion | 1 to 5 | judge | n/a | n/a | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |

Notes. "Presumes choice" is scored by the anchor as written: 2 is "no item
depends on which statuses are finished", and every run's acceptance
depends on the answer while saying it is open (a placeholder for G-003's
answer, or run 1's both readings), which is the 1 anchor's letter. No run
presumes an answer, which is what the rubric's question asks; the anchor
conflates depending on the choice with presuming it, and the owner may
rescore or reword it, after these scores, in a follow-on. Missing-choice
run 4 leaves whether `--all` with `--status` errs to the implementer and
the zero-hidden line to the owner's judgment, both inside the brief;
companion run 3 adds that an unknown tag matches nothing rather than
erring, labelled proposed.

## Disposition

**Lever.** No change is justified by this baseline. G-078 finding 6 was
observed once, on Claude Opus 5, before the headless bound was tightened;
at guides digest `3f5487904c61` on Claude Opus 5.5, the missing choice
blocked in 5 of 5 runs and the companion asked nothing in 5 of 5. Both the
model and the bound changed between the two observations and nothing here
separates them: the pair only says the current guide on Claude Opus 5.5
does not show the failure. G-118 names the bound's "must
block" sentence as the lever and says the owner would prefer the second
outcome, a shippable proposal with a non-blocking question or a proposed
decision. That is a product choice about the guide, not a defect this
suite shows, so it is the owner's to make; if made, the edit is in
`docs/work-shaping.md` under "Missing human choice", and
`question-blocks-proposal` must then accept the surfaced outcome as the
expected one before the rerun, since today it reads it as a failure whose
reason names it (the selftest's `surfaced` mode exercises that reason).
Finding 5 points at harness configuration guidance: a preview user's login
brings account plugins and their MCP servers into every session. Whether
`--strict-mcp-config` on the `claude` command line would drop the servers,
and whether anything drops the skills, is untested; recorded here for the
owner-configuration comparison, not changed.

**Further cases.** Worth building, on two grounds. This pair is now a
regression check for any guide edit at about $3 a rerun, but it cannot see
retrieval: the fixture's only knowledge is the brief, so every run's
`context` use tells nothing about the listing. The first candidate in
G-108's Next, a proposal that must find and apply a constraint held in a
listed prerequisite or plan amid plausible distractors, is the one case
that would. The owner-configuration comparison is the cheapest next fact,
since finding 5 shows even the clean directory carries account content.
The Codex row and the work row stay unbuilt.

**Limits.** Claude Opus 5.5 only, five runs per case, one fixture, one
planted choice whose two options the brief names explicitly, and a judge
that is the implementing session. Five runs show a pattern, not a rate:
a failure with a true rate under one in five could hide here. The checks
saw the clone; the trace showed no write outside it. Acceptance 5's
report is this record; the owner's rubric scores and verdict are open.
