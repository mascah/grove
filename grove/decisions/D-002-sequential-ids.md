---
id: "D-002"
type: decision
title: Use shared sequential IDs and short filenames
status: accepted
relates_to: ["D-001", "W-001"]
created: "2026-09-19T14:32:32Z"
updated: "2026-09-19T14:32:32Z"
---

## Acceptance and rationale

On 2026-09-19, the owner found the initial timestamp/random-ID/slug filenames
too long to read without horizontal scrolling. They proposed using Git's shared
metadata to prevent branches allocating duplicate numbers, then accepted moving
forward with that approach.

This revises the identity and filename choices in [D-001](D-001-starter-defaults.md).
The [record model](../../docs/record-model.md#identity-and-dates) owns the current
contract. Sequential IDs are actual identities; there is no hidden random ID.
Timestamps remain useful metadata and no longer appear in generated filenames.

Git documents the shared metadata directory used by linked worktrees in its
[worktree reference](https://git-scm.com/docs/git-worktree#_details). That supports
the proposed allocation mechanism; no allocator or concurrent-allocation
experiment has been implemented in this restart. Sharing a directory alone does
not serialize writers: the creation command must implement locking and durable
reservation before issuing IDs.

## Starter-record migration

The three starter records were renumbered before the new CLI existed. Their
creation dates and content were preserved, modification dates advanced, and
relationships and repository links updated. These old IDs are historical
references, not supported lookup aliases:

| Previous ID | Current ID |
| --- | --- |
| `64a0d76bab9f56dc32e6` | `W-001` |
| `ab03a64356d2ebd9b5e4` | `Q-001` |
| `209c4fa50e8d0b181e86` | `D-001` |

The IDs here were assigned manually in the initial record set. Future allocator
initialization must account for them; no shared counter state has been created.

## Limits and reconsideration

The coordination scope is one local repository and its linked worktrees.
Separate clones, imports, direct ID authoring, and lost counter state need
explicit handling as specified in the model. Revisit the allocation strategy
if independently edited clones become a common workflow. Reader implementation
can proceed without solving allocation recovery first.
