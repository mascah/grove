---
id: "G-026"
type: review
title: "Shaping and headless-run evidence"
status: current
formerly: "docs/reviews/2026-09-19-shaping-and-runner-evidence.md"
work: ["G-023", "G-025"]
created: "2026-09-19T21:24:18Z"
updated: "2026-09-21T21:23:30Z"
---

# Shaping and headless-run evidence

Source review for G-023 and G-025, 2026-09-19. The restart brief owns selected
direction; the work records own proposed outcomes. This report is evidence and
engineering recommendations, not an approved runner specification.

## Predecessor planning instructions

Inspected sibling skills at `ec87bb2eb19a350634c4976cc5d215a57deb58de`:
[explore](../../skills/skills/explore/SKILL.md),
[shape](../../skills/skills/shape/SKILL.md), and
[planning](../../skills/references/planning.md). Retrieved its knowledge with
its own `grove status`, `grove find exploration`, `grove find supervisor`, and
`grove context --work W-003 --phase shape --budget 7000`. W-003 here denotes
predecessor work, not the restart's updater. No workflows were invoked.

- Explore is a human conversation before outcome selection. It contributes
  ideas and records settled points; it cannot simply become an unattended job
  that supplies the human's preferences itself.
- Shape investigates facts, tests assumptions against concrete situations,
  selects or prepares work within authority, and records outcomes, constraints,
  acceptance, design, relationships, and Next. It resolves technical unknowns
  before asking for product preferences. This can support bounded unattended
  research when its mandate and allowed writes are explicit.
- Planning treats missing plans as agent work within an authorized outcome,
  separates delegation from runtime, and requires concrete waiting conditions
  for headless callers. An unchanged human wait must not spend another model run.
- The predecessor's W-003 records a prior Bench reuse assessment: extract
  execution semantics rather than adopt its bundle/seat/ledger pipeline. That
  conclusion is inherited evidence, not a fresh benchmark or restart decision.
  Its executor cancels on driver signals; it does not itself establish the
  newly selected requirement that Grove runs survive the TUI exiting.

Do not import predecessor archive/claim commands, question deletion, schema,
fixed model assignments, or readiness labels. Restart questions retain identity;
the accepted record model and this repository's CLI remain authoritative.

## Bench worker and tmux

Inspected Bench at `e39848aabcc9934477440f5b68b9c7afa649aacb`:
[worker](../../bench/scripts/bench_worker.py),
[tmux adapter](../../bench/scripts/bench_tmux.py), relevant portions of
[loop](../../bench/scripts/bench_loop.py), and
[worker test cases](../../bench/scripts/test_worker.py).
Its existing untracked archive was preserved. No worker, model, scheduler, or
tmux session was started; tests were inspected, not executed.

The worker builds a `claude -p` invocation with JSON streaming, writes process
identity before releasing a gated child, drains stdout and stderr separately,
saves raw bytes before parsing, and derives a readable log. It separately writes
the provider's terminal result and an atomic exit record. Renderer exceptions
are caught so malformed presentation need not destroy raw capture. Identity
includes PID and process start time; a bare PID is insufficient for ownership.
The tmux adapter uses an isolated socket and detached panes for worker wrappers.

These are useful patterns, with limits that a new implementation must address:

- The worker handles HUP/TERM/INT by killing its child group. A wrapper directly
  coupled to TUI signals will not satisfy continued execution; the owner and
  stream consumer must live independently of the UI.
- Its `completed` exit outcome means zero process exit plus a result event;
  the result's `is_error` is not part of that classification. The surrounding
  loop also reads provider result fields. Never translate that one label into
  work acceptance or a Done card.
- Raw stdout is flushed, but not fsynced per event. The record helper fsyncs its
  temporary file before rename, without directory fsync. Do not promise lossless
  power-failure durability from these operations.
- A never-terminated event can grow the worker's buffer. Its renderer writes
  provider text to terminal output without Grove's control-character escaping.
  Require bounded parsing and safe display; raw capture and rendered output
  serve different purposes.
- Ignoring renderer exceptions does not establish independence from a blocked
  output sink. A closed or slow TUI must not stop stream draining. Persist and
  replay through the run owner rather than pipe the provider directly into UI.
- Reconnecting a viewer and resuming a provider conversation are different
  actions. The first must not issue another model call or duplicate execution.

Official Claude documentation confirms newline-delimited `stream-json`, with
`--verbose` and optional partial-message events for token updates. The adapter
must validate installed-version behavior with fixtures and a bounded real trial;
stream events are observations, not Grove acceptance evidence by themselves.
[Claude programmatic usage](https://code.claude.com/docs/en/headless),
[CLI reference](https://code.claude.com/docs/en/cli-reference).
The official tmux guide describes detach/reattach while programs continue.
That establishes hosting feasibility, not Grove ownership/recovery semantics.
[tmux guide](https://github.com/tmux/tmux/wiki/Getting-Started).

## Recommended boundary for later implementation

The proposed end-user sequence is: discuss or investigate an idea, retain a
proposal, inspect it on the board, choose manual activation or Start
implementation, observe the attempt, then review its evidence and integration
handoff. Launch remains an explicit action bound to the selected version and
workspace. A failed launch must not fabricate a running session or completed
status transition; define recoverable ordering of ownership, launch, and the
revision-checked record update before implementation.

Work lifecycle and process state differ: active work can be waiting on a human,
have a failed attempt, or have no agent because a person owns it. Stopping an
attempt preserves work and evidence; it does not abandon the work record.
Display run activity alongside the source-specific status, never by inventing
an aggregate branch status. No new frontmatter lifecycle values are implied.
Implementation on a new branch may leave main's version proposed until records
are integrated; that is not a reason to rewrite main's record behind the user.

Keep instructions and assignment assembly useful without a process launcher.
Use one assignment shape for interactive and headless execution, with explicit
interaction mode: interactive can ask; headless persists a question and returns
a concrete wait when a human decision is required. Do not implement that as an
infinite retry loop or silently answer on the owner's behalf.

A future optional run owner should retain the input snapshot and work/source
binding, a unique attempt identity, workspace ownership, process identity,
raw stream, derived activity, provider result, and validated Grove handoff.
It survives UI loss. The TUI observes/replays its files and explicitly requests
Stop; reopening only reconnects. Host reboot or owner death is interrupted or
unknown until reconciled, never automatic success or permission to relaunch.

Prefer an on-demand per-run owner over requiring a resident service for core
Grove use. Compare detached subprocess hosting with a dedicated tmux host in a
bounded lifecycle probe before selecting one. Tmux can supply persistence and
operator access; attempt bookkeeping, duplicate-start prevention, stop ownership,
and result reconciliation still belong to Grove. A long-lived machine service
is a later option if actual scheduling/recovery needs justify it.

Probe success must include: two competing starts; UI exit/terminal loss while
stdout and stderr continue; reopening and replay without new execution; partial,
unknown and oversized events; provider failure despite a result; owner death;
Stop after reconnect; stale source revisions; and preservation of partial work.
Use fake processes first, then a deliberately bounded real provider trial.

Scheduling is another caller of the same run contract. It must define mandate,
trigger identity, duplicate/overlap policy, budget, workspace, allowed writes,
and unchanged-wait behavior. Start with one manually triggered research run
before introducing a timer. Research can produce proposed work or questions;
it must not silently authorize implementing its own proposals.
