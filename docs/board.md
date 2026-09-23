# The terminal board

`grove [--project DIR] [--json]`, with no command, opens the board. In this
checkout that is `go run ./cmd/grove`. There is no `board` subcommand and no
work-ID argument; every other command stays noninteractive, and `--help`,
`-h`, and `help` need neither a project nor a terminal.

The board draws on stderr and reads stdin, so both must be terminals, while
stdout may be redirected: `cd "$(go run ./cmd/grove)"`. Without a terminal it
refuses at once with exit 1 and names the noninteractive commands. Text from
records, paths, Git, and an attempt's provider is shown with control
characters escaped. Besides its prompted actions it writes no files, including
the framework's debug logs.

## Columns and cards

The columns (Proposed, Active, Review, Done) open on the current view: every
work record in its current state across all local branches and checkouts, as
[`versions`](commands.md#versions) decides it, the same from any checkout. An
old copy on a stale branch does not hide a later status elsewhere. Each card
is a box holding the record's ID and tag, its title, and its kind, size and
priority, or for Done the date it was last written and its candidate; the
focused card has a heavy border and a `▶` marker, and each column an accent
colour that nothing depends on. A card whose current state is only in a
checkout's uncommitted files is marked `uncommitted`. With a target, a card
none of whose committed current states is on it is marked `not on main`, and
the header names the target. Where the current states diverge, one card sits
in the earliest of their statuses, marked `⑂ 2 states`, and its detail says
which states exist and where. Work whose current state removes its record is
listed under Deleted. Done shows the most recently written cards that fit the
column, newest first, and counts the rest (`+ 23 older · / to search`).
Abandoned is hidden until `a` shows its column, and the shelf row counts it
meanwhile. Neither the bound nor the hiding moves a file.

`b` chooses between the current view and one checkout's own live files; this
changes what is displayed and switches no branch or directory. On a
checkout's board, work with no live record there is listed under Elsewhere
without a status. Questions, decisions, pages, terms, plans and reviews are
never cards; `/` finds them.

## Record detail

Enter on a card opens the record's detail: a boxed header with the ID,
status, title, planning fields, candidate, standing against the target, the
places holding its current state and when it was last written; then its body
rendered from Markdown (headings, emphasis, lists, code, tables) beside a
sidebar of the records linked to it, the timeline of commits that changed
it, and one line per current state. Linked records are listed by role,
derived from fields alone: `plan` and `review` (records whose `work` names
it; a review adds what it `examined` and whether that is the candidate),
`work` for a plan or review, `needs` and `needed by` (`depends_on` either
way), `blocked by` and `blocks` (a question's `blocks`), `part of` and
`member` (`members`), and `related` (`relates_to` either way). Tab moves
focus from the content to the linked records, to the changes, to the
timeline, and back; ↑/↓ and PgUp/PgDn scroll the content or move the cursor.
Enter on a linked record opens its own detail, of any type, and Esc returns;
Enter on a timeline commit shows the record as it was at that commit, and
Esc returns to now. Below 100 columns the detail shows one pane at a time and
Tab cycles them. A page, term, decision, question, plan or review opens in
the same screen, with the fields its type has.

## Judging a candidate

Work in `review` is a candidate to judge, and its detail is the handoff: the
header adds a Review block (the candidate, whether it is approved, whether
only the record changed since it or the tip is a new candidate, whether the
target holds it, and which checkout each action runs in), the content opens
at its `## Evidence`, and the sidebar lists the candidate's changed files
against the target with their added and removed line counts. Enter on a file
shows its diff in the content pane, escaped like record text with added,
removed and hunk lines coloured, and Esc returns to the content. `a` asks for
a verdict and approves the candidate in the branch's checkout, `f` asks for
feedback and returns the work to `active` there, and `i` confirms the merge
into the target from the target's checkout, then asks whether to remove the
branch's worktree and branch. Each runs the same operation as the command
([Work lifecycle](record-model.md#work-lifecycle)), one at a time; a result
screen shows its facts, or why it was refused, and the board is re-read. A
candidate without a clean checkout of its branch, or a target without one, is
reported instead; the board never creates a checkout.

## Attempts

Work's detail also names its attempts
([G-046](../grove/G-046-managed-runs.md)): how many, and the latest with its
outcome. `R` on proposed or active work asks for a budget in USD, then a
permission mode, both typed each time since neither has a default, and
launches one attempt as [`run`](commands.md#attempts) does; it runs on the
branch the record's current state stands on when that is not the target, in
that branch's checkout, so after feedback the next attempt continues on the
candidate's branch, and otherwise in a new `worktree-ID`. It is refused up
front for work in review, done or abandoned, work an open question blocks,
and work with an attempt still running or orphaned; `run`'s own refusals
follow, and one more: the record in this checkout changed since the board
read it.

`A` lists the attempts of the open work, or on the board every attempt,
newest first, and Enter opens one: its outcome, the facts `attempt` prints,
the provider's final report rendered like a record body, and its recent
activity newest first, one short line per event from the last 1 MiB of its
events, so a flood of output costs one bounded read. The outcome is derived,
never written: `running`, `orphaned`, `interrupted`, or for a finished
attempt a `candidate ready` (only when the attempt ended with its record
committed in review with a candidate, which the owner records at exit),
`stopped`, `failed` (no result event, an error result or a nonzero exit),
`waiting on question` (the work's latest attempt, while an open question
blocks it) or `ended without a handoff`; a clean exit alone is never ready,
and a record that says review without being committed is reported as such.
`x` asks, then stops a running or orphaned attempt as `stop` does, keeping its
partial work; `o` opens its work record. The attempts are files the board
only reads: quitting leaves an attempt running, and the next session shows
the same one. While one runs, they are re-read every 2 s, a running card is
tagged `● running`, and when one ends the board is re-read.

## Timeline

The timeline is the record's history from Git: the commits that changed its
file, newest first, each with its date, the status the record held at that
commit, its short ID, and its subject, following renames. It is the history
of the state the detail shows, named in its heading: in the current view the
first current state's branch or checkout, on a checkout's board that
checkout's HEAD. A checkout whose files differ from its HEAD gets a first
`uncommitted` row. Merges are not listed, so where the record's status is not
the newest listed commit's, a first `here` row gives it and says why. History
is read from Git when a detail is open, never while the board loads, and no
key waits for a read still in progress. It says what happened on one branch
and nothing about whether another branch contains it.
[G-030](../grove/G-030-card-lineage.md) owns this.

## Versions and selecting a checkout

`v` in a detail opens the record's versions. A version is the record's exact
content; the board read it at every local branch's tip and in every
checkout's files, and lists each differing content once with its own title
and status, current ones first. An older one is marked `older`, and its
details say why. Where several branches and checkouts hold the same content
the row is a fold (`▸ done  same on 4 branches, 4 checkouts`): Enter lists
those places, and selects nothing. Opening a card or a detail selects nothing
either. Moving to one branch or checkout and pressing Enter asks
[`workspace`](commands.md#workspace)'s resolver about exactly that version's
selector; on success the board closes and prints what `workspace` prints (the
project path on stdout, or its JSON with `--json`; checkout, branch, record,
and revision on stderr). A refusal (the version changed, its checkout is
missing or ambiguous, the record was deleted there) stays on screen with its
reason until `r` refreshes, after which a version must be selected again. The
details pane there begins with the focused version's history. Selecting never
creates a worktree, edits a record, or starts an editor, shell, or agent;
only `R` starts an agent, behind its prompts.

## Search

`/` on the board searches every record of the project in its current state,
of every type, including hidden Abandoned work, Done beyond its page, pages
and terms: typing filters by ID, type, status and title, not the body; ↑/↓
move, Enter opens the detail, Esc closes. Letters typed there filter rather
than act; Ctrl-C still interrupts.

## Rendering

Record bodies are rendered by glamour after the same escaping as every other
text, so a body's control sequences show as text; an HTML character reference
such as `&#x1b;` or `&amp;` is never decoded, since Markdown would decode it
after the escaping, and shows as typed in prose and code spans, with an extra
`&amp;` in code blocks and link targets; and of what the renderer emits only
its own styles reach the terminal. Links are shown as text, never as terminal
hyperlinks, and a relative target is shown root-relative (`/G-093-….md`). The
render is cached per record content and width.

## Keys

| Key | Where | Action |
| --- | --- | --- |
| arrows, `j` `k` | everywhere but search, where letters are typed | Move |
| `h` `l` | board | Move between columns |
| Tab | board | Switch between the columns and Deleted or Elsewhere |
| Tab | detail | Cycle the content, linked records, changes and timeline |
| Tab | versions | Switch between versions and details |
| Enter | anywhere with a cursor | Open the card, record, file diff, commit, attempt or fold; in versions, select that checkout |
| PgUp/PgDn | detail, versions, search, one attempt, sources | Scroll |
| `/` | board | Search every record |
| `a` | board | Show or hide Abandoned |
| `a`, `f`, `i` | detail of work in review | Approve, give feedback, integrate |
| `R` | detail of proposed or active work | Launch an attempt |
| `A` | board or any detail | List attempts: of that work, or every attempt |
| `x` | attempts list or one attempt | Stop the attempt |
| `o` | attempts list or one attempt | Open its work record |
| `v` | detail | Open the record's versions |
| `b` | board | Choose the current view or one checkout's board |
| `s` | everywhere | List every branch and checkout read, with diagnostics; reachable while a banner marks an incomplete result |
| `r` | everywhere | Re-read |
| Esc | everywhere | Go back; quits from the board |
| `q` | everywhere | Quit |

Below 100 columns one status column shows at a time; below 40x10 the board
asks for more room. Leaving without a selection prints nothing and exits 0;
Ctrl-C exits 1; a usage error exits 2.
