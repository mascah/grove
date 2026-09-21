---
id: "D-005"
type: decision
title: "Represent terms, plans and reviews as typed sequential-ID records"
status: accepted
created: "2026-09-21T04:37:58Z"
updated: "2026-09-21T04:39:27Z"
relates_to: ["D-002", "D-004", "W-019", "W-029"]
---

## Decision

On 2026-09-20, in an interactive shaping session for
[W-019](../work/W-019-knowledge-artifacts.md), the owner chose, one question at
a time and with the alternatives below in front of them:

- A plan or review is a **record with its own sequential ID** and frontmatter,
  linked to one or several work items by a `work` list on that record. Work's
  attached plans and reviews are derived, never stored a second time.
- A domain term is **one record per term with a sequential `T-NNN` ID**; the
  title is the term. This applies [D-002](D-002-sequential-ids.md)'s single
  identity rule rather than adding a slug identity.
- Each gets **its own type folder under the record root**, following the
  existing folder = type = prefix rule: `terms/` (`T-`), `plans/` (`P-`),
  `reviews/` (`R-`). There is no `artifacts` folder, `artifact` type, or `kind`
  field; "artifacts" stays the brief's umbrella word for plans, reports and
  reviews. The owner asked whether an abstract `artifacts` folder was really
  best; the code showed a type is one table row
  (`internal/create/create.go:26-31`), so a type per kind costs no more.
- This repository's brief and its existing plans and reviews **will be
  migrated** to that layout ("it will be weird to have [two] different layouts"),
  as separate follow-on work,
  [W-029](../work/W-029-migrate-knowledge.md), not inside W-019.

This settles representation only. Field names beyond `work`, each type's
statuses, the brief's configuration key, and schema versioning are W-019's
preparation. A `report` type is not selected yet: nothing here produces one
today, and W-020's handoff is its first consumer.

## Alternatives

- Artifacts identified by project path with a small frontmatter, or only an
  `artifacts` path list on work: smaller, but a rename breaks identity, and a
  review cannot carry the revision it examined (W-019 acceptance 3, W-020).
- One `grove/artifacts/` folder with `kind: plan|report|review`: one counter,
  but mixed browsing, an extra validated field, and one skeleton for all kinds.
- Files beside their work item (`work/W-019/plan.md`): a plan shared by several
  items has no single home.
- Slug term IDs as nullsec has today (`id: damage-application`): terms migrate
  untouched, but every command needs a second ID form and renaming a term
  changes its identity. One glossary file: no per-term links, one conflict
  point across branches.
- Leaving old plans and reviews in `docs/` as plain linked history: no done
  record changes, but two layouts persist.

## Reconsideration

Revisit if browsing `P-`/`R-` numbered files proves harder than today's
`W-019-slug.md` names, if `T-NNN` references make records unreadable in
practice, or if the nullsec cutover (W-023) finds the slug-to-ID mapping for
its terms and wikilinks costly.
