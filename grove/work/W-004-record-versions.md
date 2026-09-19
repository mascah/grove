---
id: "W-004"
type: work
title: "Inspect record versions across local branches"
status: done
kind: feature
priority: 2
size: medium
depends_on: []
relates_to: ["Q-001", "W-003"]
created: "2026-09-19T17:49:58Z"
updated: "2026-09-19T20:05:45Z"
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

## Evidence

Done 2026-09-19 on branch `worktree-W-004-W-005` at `ca420f5`, base
`b20d2b0`. The [coordination plan](../../docs/plans/W-004-W-005-coordination.md)
records the finalized JSON, ordering, incomplete-result, and selector
contract, the commits, the fixture list per acceptance item, the suite, race,
vet, gofmt, and `check` results, real use against this repository's two
worktrees, and the independent review. `versions [ID] [--json]` reads every
local branch tip through Git objects and every registered worktree's live
files at this project's prefix, validates each source alone with the
checkout's loader, groups by ID with each source's own status, compares live
records with their HEAD, and prints a selector per version binding repository,
prefix, source, commit, configuration, path, and content revision.

Every acceptance bullet has fixtures: differing statuses and bodies with
committed and live labels; a branch without a checkout; detached, dirty,
untracked, deleted, renamed, and byte-identical records; source-specific
configuration and a project below the repository root; a dependency missing
in one source; invalid YAML, duplicate IDs, an inaccessible worktree, and a
changing HEAD or branch as attributable incomplete results; BOM and CRLF
bytes agreeing with revisions; paths with spaces, tabs, and newlines; repeated
reads ordering identically; and hashed Git metadata, records, and dirty files
unchanged with no coordination state created.

Review disposition: no blocking findings; three should-fix items (unborn
HEAD, undocumented invalid-HEAD comparison, missing Git-state test) were fixed
and covered. A HEAD whose project does not validate keeps the live source
valid with `change: unknown` and a note, a revision of this record's original
proposal. Limits: no fetch, remotes, tags, or history; separate clones are
outside scope; reads are not a simultaneous snapshot, which selectors expose
at resolution time.

## Next

W-005 consumed the selector in the same branch; both were integrated into
main at `5041ae1` on 2026-09-19.
