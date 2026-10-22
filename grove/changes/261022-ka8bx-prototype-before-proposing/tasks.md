## T-01 — The prototype skill

Result: `/grove-prototype`, explore's suggestion and propose's use of
findings behave as [specs/prototype/spec.md](specs/prototype/spec.md) says,
and the documentation specifies them.

Constraints: within the rule 12 budget; messages are one sentence (rule
14).

Acceptance:
- `grove init` writes the prototype skill for each harness under `agents`;
  `rules.skills.prototype` is accepted.
- `docs/getting-started.md` or a docs page shows exploring, prototyping and
  proposing from the findings; `docs/reference.md` lists the skill and
  `grove instructions prototype "<question>"`; `docs/concepts.md` places it
  between explore and propose. `dev/walkthrough.md` runs those steps (rule
  8).
- On the ledger fixture, in Claude Code: an explore about how a new report
  reads suggests a prototype; `/grove-prototype` builds it on a
  `prototype/<slug>` branch and reports a verdict; the proposal that
  follows cites the branch and commit, and the target and document root are
  unchanged by the prototype.
- `just check` passes.
