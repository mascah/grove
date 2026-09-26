---
id: "G-186"
type: term
title: "Entrypoint revision"
status: proposed
created: "2026-09-25T23:19:59Z"
updated: "2026-09-25T23:29:12Z"
relates_to: ["G-169", "G-062", "G-152", "G-170"]
---

## Meaning

An integer naming what an entrypoint `grove init` writes (the `grove-work`
and `grove-shape` skills, their Codex policies, the `grove-reviewer` agent
definition) needs of the `grove` it calls: which guides it loads and how.
Each such file states it as `grove entrypoint revision N` and passes it as
`grove guide NAME --entrypoint N`. A binary serves a range of revisions:
`guide` refuses an `--entrypoint` outside it, `init --check` reports a
marked file outside it and exits 1, and `run` refuses to launch with one. A
file with the managed marker and no revision line predates revisions and is
revision 1, legacy.

Not a [revision](G-062-revision.md): that is one file's `sha256:` content
identity, and two files of one entrypoint revision may differ in every
byte. Not the release version `grove version` prints, nor the record
schema's `schema_version`: a release may keep the entrypoint revision, and
a revision changes only when what entrypoints need of the binary changes,
never for new wording.

## Relationships

Introduced by [G-169](G-169-harness-upgrade-compatibility.md), whose plan
[G-184](G-184-g-169-plan-entrypoint-revisions.md) records the design; the
command reference's Entrypoint revisions section owns the behaviour.
[G-170](G-170-release-identity.md) keeps release version, schema and
entrypoint compatibility distinct. The files it names are generated text
under the [shipped document](G-152-shipped-document.md) rule.
