# The terminal board

`grove [--project DIR] [--json]`, with no command, opens the board. In this
checkout that is `go run ./cmd/grove`. There is no `board` subcommand and no
work-ID argument; every other command stays noninteractive, and `--help`,
`-h`, and `help` need neither a project nor a terminal.

The board draws on stderr and reads stdin, so both must be terminals, while
stdout may be redirected: `cd "$(go run ./cmd/grove)"`. Without a terminal it
refuses at once with exit 1 and names the noninteractive commands. Text from
records, paths, Git, and an attempt's provider is shown with control
characters escaped. Besides its prompted actions, and the `## Answer` heading
`e` hands the owner's editor, it writes no files, including the framework's
debug logs.

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
none of whose committed current states is on it is marked `not on main`,
unless it has a live attempt, which runs on a branch anyway; the header names
the target. Where the current states diverge, one card sits
in the earliest of their statuses, marked `⑂ 2 states`, and its detail says
which states exist and where. Work whose current state removes its record is
listed under Deleted. Done shows the most recently written cards that fit the
column, newest first, and counts the rest (`+ 23 older · / to search`).
Abandoned is hidden until `a` shows its column, and the shelf row counts it
meanwhile. Neither the bound nor the hiding moves a file. ←/→ and `h` `l`
move to the next column with cards, skipping empty ones, which are still
drawn; past the last column with cards the focus stays.

`b` chooses between the current view and one checkout's own live files; this
changes what is displayed and switches no branch or directory. On a
checkout's board, work with no live record there is listed under Elsewhere
without a status. Questions, decisions, pages, terms, plans and reviews are
never cards; `/` finds them.

The header gives the time of the last read. `r` re-reads everything, and
the board also re-reads when its terminal window regains focus, unless a
read or an action is under way or a detail left at a timeline commit or a
diff is open, even beneath another screen; a terminal that does not report
focus sends nothing, so there only `r` does. While an attempt runs, the
board also follows the branch tips (see [Attempts](#attempts)). Otherwise it
starts no process except on focus or a key, and an uncommitted edit in a
checkout's files shows only after a re-read.

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
Enter on a linked record opens its own detail, of any type, and Esc returns.
A record already on the path is returned to instead of opened again, so A,
then B from A's sidebar, then A from B's is A alone, one Esc from the board. From
16 rows, a row at the top of the header shows the path, such as `board ›
G-108 › G-115`, losing its start when too long; a record `o` opened from an attempt
has `attempts` before it, and Esc from it returns to the attempt.
Enter on a timeline commit shows the record as it was at that commit, and
Esc returns to now. Below 100 columns the detail shows one pane at a time and
Tab cycles them. From 100 columns `w` hides the sidebar so the content takes
the width and a mouse selection takes no sidebar text, with Tab cycling as
below 100 columns; `w` shows it again, and the choice lasts for the session.
A page, term, decision, question, plan or review opens in the same screen,
with the fields its type has.

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
state and time. `R` on proposed or active work asks for a budget in USD, then a
permission mode, both typed each time since neither has a default; then, where
Enter alone asks for nothing, the bound (`plan` stops the attempt at a
committed plan), a model and an effort
([G-134](../grove/G-134-bound-an-attempt-at-its-plan-and.md)); and launches
one attempt as [`run`](commands.md#attempts) does; it runs on the
branch the record's current state stands on when that is not the target, in
that branch's checkout, so after feedback the next attempt continues on the
candidate's branch, and otherwise in a new `worktree-ID`. It is refused up
front for work in review, done or abandoned, work an open question blocks,
and work with an attempt still running or orphaned; `run`'s own refusals
follow, and one more: the record in this checkout changed since the board
read it.

`A` lists the attempts of the open work, or on the board every attempt
([G-109](../grove/G-109-attempts-usability.md)), in three groups, newest first
within each: **Needs you**, **Running** and **Settled**. Each row gives the
work's ID and title, a short state and how long it has run or how long ago it
ended. The title goes below 60 columns, and the state is cut last. Only the
latest attempt of work that is still proposed, active or in review needs you:
a candidate to judge, a question to answer, a plan to read, a failure, an
interruption, an end without a handoff, or feedback given or a question
answered that awaits the next launch. An attempt bounded at its plan that
ended cleanly, with no question and no candidate, its record committed and
nothing left uncommitted, is `plan ready: read it, then R`: its work's detail lists the plan, and `R` there without the
bound launches the implementation on the same branch, which is the owner's
approval of that plan. An
orphan always needs you, since its process runs unowned. A stopped attempt,
an earlier attempt of the same work and any attempt of work now done or
abandoned are settled, and each says why, such as `done: candidate 1614e89`
or `candidate 71a650e, superseded`.

Enter opens one attempt. At the top are the work's ID and title, a coloured
state and the run's configuration: attempt, model and provider version,
what the launch asked for (bound, model, effort, and the reviewer
definition's digest or its absence), budget with what it cost and, where
more than one model spent it, the split, permission mode, branch and base,
start and end.
What the run used follows: turns, tokens, context size against the model's
window, subagents, compactions, tool calls and tool errors. Then come `State`, the
outcome in a sentence with why it is settled, and `Next`, the keys that act
on it. `d` shows or hides the details: the facts `attempt` prints, the raw
`events.jsonl` and `stderr` paths included. Last come the provider's final
report, rendered like a record body, and the activity, newest first. From
100 columns they sit side by side, the report on the left. Each activity row
has the event's local time where the provider gave one, and consecutive
identical rows are one row with a count, such as `system: thinking_tokens
(×14)`.

The activity and the figures come from one bounded read of the last 1 MiB of
the events, so a flood of output costs the same. A count from a window that
began inside the file is shown as `≥N`, and an unknown figure as `–`, never
as 0. Tokens are the result event's totals once the run has ended; while it
runs, only the input tokens are summed, since Claude reports a message's
output before writing it. Turns are the sum of the result events' own counts,
one per query of a resumed session, and until the run ends with one they are
top-level messages shown as `≈N`. A read that fails keeps what the last good
one showed, under the failure. The figures are provider-neutral fields
(`attempt.Metrics`) that the reader of Claude's stream fills.

The outcome is derived, never written: `running`, `orphaned`,
`interrupted`, or for a finished attempt a `candidate ready` (only when the
attempt ended with its record committed in review with a candidate, which
the owner records at exit), `stopped`, `failed` (no result event, an error
result or a nonzero exit), `waiting on question` (the work's latest attempt,
while an open question blocks it), `answered since` (the work's latest
attempt, when a question that blocks it, not asked after the attempt ended,
is resolved and was last written from the second it ended on) or `ended
without a handoff`. A clean exit
alone is never ready, and a record that says review without being committed
is reported as such. `x` asks, then stops a running or orphaned attempt as
`stop` does, keeping its partial work, and `o` opens its work record. The
board only reads the attempts' files: quitting leaves an attempt running,
and the next session shows the same one. While one runs, they are re-read
every 2 s with one listing of the branch tips, a running card is tagged
`● running` (`● orphaned` for an orphan) and has a border in a colour no
column uses, and the board is re-read when one ends or a branch tip has moved
since the last read, so a status an attempt commits moves its card while it
runs; a moved tip waits while a detail left at a timeline commit or a diff
is open.

## Answering a question

An attempt that needs a decision writes a question that `blocks` its work
and ends ([work guide](work-execution.md#when-a-human-decision-is-missing)).
`e` on an open question's detail, or on the attempts list or one attempt
waiting on a question, which opens that question's detail above the attempt,
answers it ([G-125](../grove/G-125-answer-a-blocking-question-from.md)). The
detail's header names where `e` writes, or why it cannot: the one checkout on
the branch the question's current state stands on, whose copy is still what
the board read. No such checkout, two on the branch, or a file changed since
the read is refused with the reason, and nothing is written.

Otherwise the board appends a `## Answer` heading when the body has none,
suspends itself, and runs `$VISUAL`, else `$EDITOR`, else `vi`, as Git does,
on the question's file in that checkout, on the terminal itself. When the
editor exits the board resumes. An editor that saved nothing gets the heading
taken back, also when the session ends meanwhile, and one that failed is
reported; either way nothing more is written. After an edit the board
re-reads, and a prompt asks to resolve the question and commit it there: `y`
runs `update --set status=resolved --commit` at the revision the editor left,
so one commit on that branch holds the answer and the status, and the result
names the work to launch next with `R`; `n` or Esc leaves the edit
uncommitted in that checkout and says so, and `e` then reopens it and offers
the resolve again, saved or not. The attempt then shows `question answered:
R again`, which is read from the question's `updated`, as `update` sets it: a
status changed by hand without it leaves the attempt `ended without a
handoff`, and a later update of an answered question can mark a later attempt
answered too. Only questions are edited, and only in the owner's editor; the
board has no text editing of its own.

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
only `R` starts an agent, behind its prompts, and only `e` starts the
owner's editor, on an open question (see [Answering a
question](#answering-a-question)).

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
| ↑/↓ | everywhere | Move |
| `j` `k` | everywhere but search, where letters are typed | Move |
| ←/→, `h` `l` | board | Move to the next column with cards |
| Tab | board | Switch between the columns and Deleted or Elsewhere |
| Tab | detail | Cycle the content, linked records, changes and timeline |
| Tab | versions | Switch between versions and details |
| Enter | anywhere with a cursor | Open the card, record, file diff, commit, attempt or fold; in versions, select that checkout |
| PgUp/PgDn | detail, versions, search, one attempt, sources, result | Scroll |
| `/` | board | Search every record |
| `a` | board | Show or hide Abandoned |
| `a`, `f`, `i` | detail of work in review | Approve, give feedback, integrate |
| `R` | detail of proposed or active work | Launch an attempt |
| `e` | detail of an open question; attempts list or one attempt waiting on a question | Answer it in your editor, then resolve and commit it on its branch |
| `A` | board or any detail | List attempts: of that work, or every attempt |
| `x` | attempts list or one attempt | Stop the attempt |
| `o` | attempts list or one attempt | Open its work record |
| `d` | one attempt | Show or hide its details |
| `v` | detail | Open the record's versions |
| `w` | detail, from 100 columns | Hide or show the sidebar |
| `b` | board | Choose the current view or one checkout's board |
| `s` | everywhere | List every branch and checkout read, with diagnostics; reachable while a banner marks an incomplete result |
| `r` | everywhere | Re-read |
| Esc | everywhere | Go back; quits from the board |
| `q` | everywhere | Quit |

Below 100 columns one status column shows at a time; below 40x10 the board
asks for more room. Leaving without a selection prints nothing and exits 0;
Ctrl-C exits 1; a usage error exits 2.
