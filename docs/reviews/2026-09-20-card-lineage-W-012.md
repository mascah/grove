# Card lineage W-012: evidence, 2026-09-20

Branch `worktree-W-012`, worktree `.claude/worktrees/W-012`, base `a0fd23a`
(main). Carried out through `/grove-work W-012` in an interactive session, the
dogfooding run that W-010 was waiting for. One implementer; one separate
reviewer agent that edited nothing. Nothing here was merged or pushed.
[W-012](../../grove/work/W-012-card-lineage.md) owns the outcome and
acceptance; its [plan](../plans/W-012-card-lineage.md) holds the design and the
adjustment made after review.

**What this evidence is not.** Every check below is automated or an agent's
review. Acceptance item 5 also asks for the owner's judgment of the view in a
demo. The owner gave it on 2026-09-20 after running the demo below, and W-012
records it: a step in the right direction that works as described, with the
TUI's information architecture left for later iteration.

## The owner's decision

The record's Next asked whether lineage replaces the version list as what a
card opens to, or sits beside it. Asked once at the start of the session, with
a mock-up of each, the owner selected "beside, history first": the card screen
and its version list stay (Q-001 still needs an explicit version selection),
and History leads the details pane. Everything else was a routine technical
choice inside the record.

## Commits

| Commit | Change |
| --- | --- |
| `5dd8c3e` | Plan, the owner's choice in the record, W-012 set active through `grove update` |
| `e117759` | `versions.HistoryContext`; `Backend.History`, the history read, cache, and History section in `internal/tui`; tests; pseudo-terminal stages |
| `49d2b47` | README: the card's History section |
| `24c1264` | Review fix, round 1: merges that only brought commits in are not rows; merge and date tests |
| `cfcb5d9` | Review fix, round 2: merges diffed against each parent, in date order; one is dropped only when the row beneath holds the same content |
| `b8fc243` | Review fix, round 3: plain `git log --follow`, merges never listed, in date order; a first row where the record's status is not the newest commit's; last reviewed revision |
| `9337a44` | After the final review: that row's wording claims only what is known; two more assertions; README |

## Acceptance, by item

1. **Each commit with date, subject, and status, newest first.**
   `TestHistoryFollowsOneLineOfCommits` builds a real repository where W-001 is
   created, edited, given a status, renamed without a change (to a name full
   of pattern characters), given a status today's schema rejects, and
   finished, with an unrelated commit between: the lineage is exactly those six
   commits, newest first, each with the status the file held there, and the
   unrelated commit is absent; a commit that deleted the file shows `-`.
   `TestHistoryAcrossMergesAndDates` renames and activates the record on a
   branch whose clock is behind, merges it without fast-forward, and gets the
   three commits in order with the rename holding its status; it checks the
   author's date against real Git, and that a hand-resolved merge is no row.
   `TestHistorySaysWhenTheNewestCommitIsNotTheRecord`: the card then gives the
   record's status in a first `here` row.
   `TestHistoryFollowsTheFocusedVersion` checks the rendered rows: date, status,
   seven-character ID, subject, and `?` where no status could be read.
2. **Two branches, two lineages.** The same fixture's `feature` branch gives
   `abandoned: abandon on feature` then `init`, read from either checkout's
   root. In the model test, moving from main's rows to feature's fold replaces
   main's lineage with feature's under a heading naming the branch; places at
   one commit share one read. `TestBoardConnectedWorkflow` does the same walk
   on a real main/feature repository through the real backend, and shows the
   `uncommitted` row for a checkout whose file was edited.
3. **Control characters escaped.** The real-Git fixture carries ESC and a
   bidirectional override in a subject through `HistoryContext` unchanged; the
   model test puts an OSC sequence, BEL, and an override in a subject and a
   CSI sequence in an error, requires their escaped forms on screen, and
   requires that the raw bytes are nowhere in the rendered output.
4. **Cancels cleanly, never blocks keys, no reads for the board.**
   `TestHistoryFollowsTheFocusedVersion` moves around the board, the checkout
   chooser, and the sources screen first: zero history calls.
   `TestHistoryReadYieldsToEveryKey` blocks the backend and shows that moving
   to another version, Enter on a version, `r`, Esc, and `q` each cancel the
   read in flight (its context reports `context.Canceled`), none is refused or
   delayed, a superseded reply is ignored, a refusal returns to a card that
   reads its history again, and nothing starts after `q`.
   `TestHistoryUnderPrefixAndCancellation` kills a blocked `git log` and a
   blocked `git cat-file` and requires the process collected. The
   pseudo-terminal script blocks Git as a card opens and leaves with `q`,
   Ctrl-C, and Esc: the child is dead each time, terminal modes are restored,
   and after Esc the board still takes keys.
5. **Suites.** Below. The owner's demo judgment is recorded in W-012.

## Verification

At `9337a44`, uncached, on the owner's Mac (Git 2.50.1, Go 1.26):
`gofmt -l .` prints nothing; `go vet ./...` passes; `go test ./... -count=1`
and `go test -race ./... -count=1` pass in every package, including
`internal/tui`'s pseudo-terminal run against a freshly built (and, under
`-race`, race-enabled) binary; `go run ./cmd/grove check` reports 17 valid
records. The longest packages were `internal/versions` (50 s, 56 s under
`-race`) and `internal/tui` (16 s, 41 s).

## Independent review

One reviewer agent, which ran its own fixtures and mutations and changed
nothing, reviewed `a0fd23a..49d2b47`, then each fix.

The first design asked Git to diff merges so that a merge row would carry a
status. Each round found a confirmed, consequential defect of that choice, and
the third fix removed the choice instead of patching it again:

| Round | Finding | Disposition |
| --- | --- | --- |
| 1, `49d2b47` | Every merge that brought a branch in became a row, repeating its commits and naming another branch: nine of thirteen work records on main. No fixture had a merge and no test read a date from Git (two planted mutations survived). | `24c1264`: drop a merge whose content a listed commit also has; merge and date tests |
| 2, `24c1264` | That filter dropped a conflict resolved by keeping one side whole, so the first row could show a status the record does not hold; Git itself omitted such a merge when the kept side was the first parent's; a merge carrying another merge's content stayed. | `cfcb5d9`: merges diffed against each parent, date order, drop a merge only when the row beneath holds the same content |
| 3, `cfcb5d9` | Under any merge diff, since the first commit, `--follow` takes a rename seen from a merge's other parent for the record's own: on this repository's own workflow (rename and activate on a branch, merge without fast-forward) the activating commit vanished and the rename showed as a deletion. The adjacent-row filter also kept every merge when timestamps tie, which ten seconds of main's history already do. | `b8fc243`: plain `git log --follow --no-merges --date-order`, the record's own example; merges are never rows; a first row where the record's status is not the newest commit's |

The final round, on `b8fc243`, reported **nothing consequential remaining in
the read**: its four rename and merge fixtures, including the two that broke
the earlier designs (rename then merge; rename and edit against an edit of the
old path, in conflict; rename on both sides; delete against edit), match plain
`git log --follow` commit for commit with no commit lost, duplicated, or above
its descendant; every tracked record and document on main matches it with no
merge row; and nothing of the removed merge logic is left in code or text. It
recorded one point below that bar, confirmed: the first row then read "a merge
left the record this way", which is wrong when a clean merge combined a status
set by an older listed commit with a newer commit from the other line. The row
never showed a wrong status and named no branch. Taken in `9337a44`, after the
last review and so **not independently reviewed**: the row now reads
`here  <status>  the record's status here; the newest commit below differs
because merges are not listed`, and the two conditions it found unpinned
(newest rather than oldest commit; not for a checkout with uncommitted
changes) are asserted. That commit changes wording and tests only.

That is three fix and review rounds, the workflow's cap. It ended with no open
consequential finding, not by exhausting the cap with one outstanding.

Smaller points taken: a redundant truncation, the README's "any key replaces a
read" (now "no key waits for"), `-` for a commit that deleted the file, and
pins for it and for `--date-order`. Declined, with reasons: a status holding
control bytes misaligns its column (still escaped, cosmetic); the read also
runs below 100 columns where details wait behind Tab (it is then ready when Tab
is pressed); the scroll clamp after a reply has no test of its own. Removing
`--no-merges` alone leaves the suite green, because on Git 2.50.1 `--follow`
lists no merges anyway; the flag states the intent.

Attacks the reviewer reported clean, none of which the later fixes touched:
every history string reaches the screen escaped and every rendered row is
exactly the pane's width at five widths; nineteen hostile Git configuration
settings left the log output byte-identical; odd subjects, SHA-256
repositories, delete-then-re-add, and pattern-like names parse; the mutations
it planted in the model (the stale-reply, superseding, quitting, `HeadPath`,
and board-checkout guards) all failed the suite.

## Limits

- Merges are never rows. Content that a merge itself gave the record, by a
  hand resolution or a clean three-way combination, has no row of its own; the
  `here` row covers only a status that differs from the newest listed
  commit's, and only for a committed or unchanged version. It cannot say which
  commit or merge set that status.
- History is uncapped: one blob lookup per commit. A full read of this
  repository's longest file history (94 commits) took 0.07 s in review.
- Only Git 2.50.1 was run.
- Rows show the author's date in the local time zone, and rename detection is
  Git's (`--follow`): a rename combined with a large rewrite ends the lineage
  there.
- Below 100 columns the History section is behind Tab with the rest of the
  details.

## Demo

```sh
cd .claude/worktrees/W-012 && go run ./cmd/grove
```

Right to the Active column, Enter on W-012 or W-010: History leads the details
pane for this checkout. Down to a fold, Enter, and down through its places to
see each branch's or checkout's lineage; Tab then PgDn scrolls a long one.
