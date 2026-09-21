---
id: "G-033"
type: plan
title: "G-030 card lineage plan"
status: current
formerly: "docs/plans/W-012-card-lineage.md"
work: ["G-030"]
updated: "2026-09-21T21:11:16Z"
---

# G-030 card lineage plan

**Goal:** an open card's details pane leads with the focused version's history:
each commit that touched the record's file, with its date, the record's status
at that commit, and the subject, newest first.

**Spec:** [G-030](G-030-card-lineage.md). On 2026-09-20 the
owner selected "beside, history first": the card screen and its version list
stay as they are (G-002 still needs an explicit version selection), and History
is the first section of the details pane. No new screen and no new key.

**Base:** `a0fd23a` (main), branch `worktree-W-012`.

## Design

Which history is shown follows the focus, so acceptance 2 falls out of moving
between rows:

| Focus | History of |
| --- | --- |
| ID header (where a card opens) | the board's checkout's version; for an Elsewhere card, the group's first version |
| A fold | its first place; the heading names it |
| A version row | that branch tip or checkout HEAD |

The heading always names the branch or checkout, and says nothing about any
other branch.

**The read** (`internal/versions/history.go`, `HistoryContext(ctx, root,
commit, path)`), two Git processes per read, both under the caller's context:

1. `git log --follow --raw --no-abbrev --diff-merges=first-parent
   --no-show-signature --no-color --format=%x00%H %at %s COMMIT --
   :(literal)PATH` (superseded, see the adjustment below), run in the project root so the project-relative record
   path resolves under a nested prefix. `--raw` gives the record's blob ID at
   each commit, so renames need no path handling and no path is ever parsed
   from Git's quoted output. A commit line starts with NUL, which a subject
   cannot hold; a raw line starts with `:`.
2. One `git cat-file --batch` for those blobs, through the existing `objects`
   reader (IDs only, cached by ID). The status comes from
   `project.ParseRecord`, which reports the field even for a historical record
   that today's validation would reject; an unreadable status shows as `?`.

Adjusted after review on 2026-09-20. The first design asked Git to diff
merges (`--diff-merges=first-parent`, then `separate`) so that a merge would
carry a blob and a status. Three review rounds each found a defect of that
choice: merges that brought a branch in became rows naming another branch; a
filter for those dropped a conflict resolved by keeping one side; an
adjacent-row filter depended on an order that tied timestamps break; and, under
any merge diff, `--follow` takes a rename seen from a merge's other parent for
the record's own, loses the commits after it, and shows the rename as a
deletion. The read is now what the record's own example used, plain
`git log --follow`, stated as `--no-merges --date-order` with a NUL-separated
format: merges are never rows, and no commit is listed above one made from it.
The cost is that content a merge itself gave the record has no row. The view
covers the part of that a reader could be misled by: when a committed or
unchanged version's status is not the newest listed commit's, a first `here`
row gives the record's status and says merges are not listed. It does not say
a merge set it: a clean merge can combine a status set by an older listed
commit with a newer commit from the other line. A commit that deleted the file
shows `-`.

A live version follows its path at HEAD (`HeadPath` when the checkout renamed
it) from the checkout's HEAD commit. A change other than `unchanged` is a first
"uncommitted" row built from the classification the board already has, with no
read. An `added` record has no committed history and starts no read.

**The model** (`internal/tui`): `Backend.History` is a third effect. After any
message that starts no other command, the model starts a history read when the
card screen is showing, nothing else is pending, and the focused history is not
already held. Reads are keyed by commit and path, so places at one commit share
one read, and results are kept until the next inspection replaces the result.
The one-read rule holds: a history read is `pending`, a newer one cancels and
outdates it, and so do refresh, a workspace selection, Esc, and quit. A history
read never blocks a key: refresh and selection supersede it rather than refuse.
The board screen never starts one. Subjects reach the screen only through
`wrap`, which escapes them as it escapes record text.

## Tasks

1. `versions.HistoryContext` with a real-Git fixture test: create, edit,
   rename, status changes, a control character in a subject, a second branch
   with its own commits, cancellation. (Acceptance 1, 2, 3 at the read.)
2. Model and view: `Backend.History`, the pending kind, the cache, the History
   section, the uncommitted row. Model tests: opening a card starts one read and
   the board load none; moving rows supersedes; a late reply is ignored; Enter
   and `r` supersede; Esc cancels; escaped subject on screen. (Acceptance 2, 3,
   4.)
3. Wire `versions.HistoryContext` in `run.go`; extend the connected pseudo-
   terminal workflow test to see a history row. README and record model only if
   their contract changes (the README's board description does).
4. Verification per `AGENTS.md`, independent review, evidence in
   `docs/reviews/`, record reconciliation. The owner's demo judgment stays
   separate.
