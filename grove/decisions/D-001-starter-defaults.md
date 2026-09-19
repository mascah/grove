---
id: "D-001"
type: decision
title: Adopt the starter record file defaults
status: accepted
relates_to: ["W-001", "D-002"]
created: "2026-09-19T14:08:40Z"
updated: "2026-09-19T14:32:32Z"
---

## Acceptance

On 2026-09-19, the owner selected root-level `grove/` as the configurable default,
then accepted the drafted file defaults with "YeaI accept those defaults".

The [record model](../../docs/record-model.md#on-disk-contract) owns those defaults
and their rationale. This is the acceptance receipt, not a second editable schema.
The [restart brief](../../docs/restart-brief.md) remains the product direction owner.
Supporting reader proposals are distinguished from accepted defaults in the model;
this decision does not assert that the CLI or those proposals are implemented.

## Reconsideration

The initial random-ID and timestamp-prefixed filename trial proved cumbersome.
[D-002](D-002-sequential-ids.md) replaces those choices with sequential IDs and
short filenames. The other defaults remain accepted; the record model reflects
the current contract.

Revisit when actual browsing or ID entry is awkward, a schema change needs
migration, or the first reader exposes a concrete inconsistency. Record a
replacement choice explicitly instead of silently rewriting this acceptance.
