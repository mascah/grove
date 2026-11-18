## T-01 — A guide to a word budget on context and rules, trialed on the fixture

Result: `docs/guides/word-budget.md` is a recipe that adds to a project's
own checks a count of the words of `context` and of every `rules` line in
`grove.yaml`, failing above a budget the owner chooses; it is trialed on
the ledger fixture and linked from the guides list, the docs site and the
reference.

Files:
- `docs/guides/word-budget.md`
- `docs/index.md`
- `docs/reference.md`
- `mkdocs.yml`

Interfaces: none.

Constraints: the guide uses only Grove's existing features and changes
nothing Grove reads or prints (rule 8); it states no default budget (open
question 2); it links `docs/reference.md` and the config spec rather than
restating where `context` and `rules` are printed.

Acceptance:
- The guide opens as the other guides do, saying it is a tried recipe and
  not part of what Grove does. It says why a budget helps: `context`
  reaches every skill and a `rules` line every instruction of its skill or
  document kind, as rule 12's budget does for Grove's own text.
- It gives one count, the words of `context` plus every string under
  `rules.documents`, `rules.skills` and `rules.delivery`, with no key of
  `rules` counted, and shows it as a line in a `justfile` recipe and as a
  step in the workflow of [Run Grove in CI](../../../docs/guides/run-grove-in-ci.md),
  each failing with a message naming the count and the budget.
- On a fresh ledger fixture whose `grove.yaml` is given a `context` and
  one `rules.skills.implement` line, following only the guide: the count
  equals the words of those two texts, the check passes under the budget,
  and it fails, naming the count and the budget, once a line takes the
  total over it.
- The implementing pull request's description says what the trial did and
  showed, and the guide has what the trial taught.
- `docs/index.md`'s guides list and the Guides section of `mkdocs.yml` link
  the guide, and `docs/reference.md` links it from the sentence "Short
  rules work better than long ones."
- `just check` and the docs site build pass.
