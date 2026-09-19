---
id: "W-004"
type: work
title: "Inspect record versions across local branches"
status: proposed
kind: feature
priority: 2
size: medium
depends_on: []
relates_to: ["Q-001", "W-003"]
created: "2026-09-19T17:49:58Z"
updated: "2026-09-19T17:54:10Z"
---

## Outcome

From one checkout, inspect the records on local branch tips and in registered
worktrees, with enough source information to choose which version to act on.
Main can show Fable's branch progress without switching branches or merging
records. Q-001's grouping and explicit-selection policy is accepted; the
implementation contract below remains proposed.

## Why now

The owner selected cross-branch coordination as the next experience. Priority 2
reflects that direction; medium reflects committed-tree loading, source-local
validation, and live/committed comparison. Safe mutation in W-003 is useful
alongside this feature but is not a read-only inspection prerequisite.

## Constraints

- Preserve existing `list`, `show`, and `check` behavior. New cross-branch reads
  change no records, configuration, refs, index, worktrees, or allocator state.
- Use this repository's accepted record model, including configured record roots
  and source-local relationship validation. No schema migration or new records.
- Q-001 owns version policy. No automatic status reconciliation, integration
  inference, timestamp precedence, or inferred execution ownership.
- Group by record ID and show each branch's status. Preserve explicit source
  choices for opening a workspace; no group-level default grants an editing
  destination, even when observations contain identical bytes.
- First scope: local branch tips and registered live worktrees in the same
  repository. No fetch, remote refs, tags, historical search, or persistent index.

## Design proposal

Proposed CLI surface: `grove versions [ID] [--json]`, with existing `--project`
discovery. Without ID, list grouped observations for all three record types;
with ID, restrict presentation to that identity after validating each source.
Plain directories fail with a clear Git-required diagnostic.

Retain a committed snapshot for every local branch tip and a live observation
for every registered checkout, including detached HEAD. Read each source at
the selected project's repository-relative location and use its own
`grove.yaml`. Missing configuration means the project is absent in that source;
an existing invalid configuration or unreadable path is an error. A configured
but missing record root follows the accepted model's validation failure.

Freeze branch reads to observed commit IDs. Record worktree identity before and
after reading; if branch/HEAD changes, report an unstable source. This does not
promise a simultaneous snapshot across all live files or branches. Use Git's
NUL-delimited path formats and full refs; do not parse display-oriented paths.

Each observation needs record ID, source kind (committed/live), full branch ref
when attached, observed commit, worktree identity/path for live data, relative
project and record paths, configuration revision, content revision, and validated
metadata. JSON also
provides exact source text and source diagnostics. The record content revision
must use the same byte-hash convention as W-003; a version selector additionally
binds source identity and cannot be just the record hash or a transient row
number. Finalize and document this selector format before implementation.

Compare live records against that checkout's observed HEAD, including untracked
records and records deleted from live files. Identical bytes can be grouped for
display while retaining distinct source identities. Do not present a committed
record as currently present in a worktree when its live file is absent.

Validate each entire source before admitting its records as valid observations.
Show diagnostics for invalid/inaccessible sources alongside valid sources,
mark the aggregate incomplete, and exit 1. JSON carries that completeness
status. Failure to discover the source inventory also fails; it cannot produce
a successful empty view. An ID absent from all successfully inspected sources
is an explicit not-found result. Define exact output fields and deterministic
ordering in the implementation preparation, preserving numeric ID order.

## Acceptance

- A main/feature fixture shows differing statuses and bodies for one ID, with
  correct committed and live source labels; no status is selected as authority.
- A branch without a checkout contributes its committed records. Detached,
  dirty, untracked, deleted, renamed, and byte-identical records are represented
  without losing provenance or confusing absence with read failure.
- Source-specific configuration and a project below the repository root work.
  A dependency missing in one source stays invalid even if another contains it.
- Invalid YAML, duplicate IDs within one source, inaccessible worktrees, and
  changing branch/HEAD identity produce attributable incomplete results.
- Full source text and content revisions agree byte-for-byte, including BOM and
  CRLF. Source selectors distinguish identical content in different contexts.
- Paths containing spaces, tabs, and newlines round-trip in machine output;
  human output escapes controls. Repeated unchanged reads order consistently.
- Ref, index, record, and allocator-state checks prove reads leave them unchanged;
  the existing local inspection suite still passes.

## Preparation and execution boundary

Inspected base: `ee69c42`. Likely interfaces: source loading and graph validation
in `internal/project`, a new Git-source package, and `internal/cli`. Extract
shared validation without giving the reader write responsibilities. Do not reuse
the allocator's ID-prefilter scan as a full parser.

Recommend a single Fable agent in an isolated worktree. W-003 may touch the same
loader and CLI, so separate work IDs do not prove safe parallel implementation.
Agree on revision representation and inspect the integration base before
dispatch. Required verification: relevant fixtures, full Go suite, race suite,
vet, and independent review of source identity and incomplete-result handling.

## Next

Finalize the JSON/selector contract and ordered implementation plan using
Q-001's accepted policy. W-005 consumes this source identity. Do not dispatch
this draft as an already approved specification.
