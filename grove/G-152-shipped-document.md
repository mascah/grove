---
id: "G-152"
type: term
title: "Shipped document"
status: proposed
created: "2026-09-25T19:06:02Z"
updated: "2026-09-25T19:06:02Z"
---

## Meaning

A document the `grove` binary carries and prints or writes wherever it runs:
the work guide and the shaping guide (`grove guide work|shape`), the record
model (`grove guide model`) and the reviewer agent definition that `grove
init` writes. It is read in projects that hold none of Grove's files and
whose own `G-` IDs are live, so it must read the same, and be true, in every
one of them: it links only within itself or to `https://`, names another
shipped document by the command that prints it, and names no Grove record,
path, invocation or history. This repository's file is the one editable
owner; the binary's copy is what an adopting project reads, and the guides
digest in `grove version` names that copy.

Not shipped: the brief, the command reference, the board guide, AGENTS.md,
this repository's records, and its own `.claude/skills/` adapters, which
read the documents as files. The adapters `init` writes are generated text
under the same rule, not documents.

## Relationships

The rule lives in AGENTS.md's constraints; its reason is
[G-146](G-146-how-should-an-adopting-project-r.md)'s answer, which
[G-151](G-151-strip-grove-repository-pointers.md) proposes extending from
links to mentions. `TestGuideAndVersionNeedNoProject` enforces it. Delivered
by [G-040](G-040-portable-bootstrap.md) (the guides, `init`),
[G-134](G-134-bound-an-attempt-at-its-plan-and.md) (the reviewer definition)
and [G-144](G-144-give-adopting-projects-the-recor.md) (the model).
