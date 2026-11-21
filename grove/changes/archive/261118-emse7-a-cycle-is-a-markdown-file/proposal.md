# A cycle is a Markdown file with plan's notes

## Problem

A cycle file holds seven fields and nothing else: `title`, `created`,
`changes`, `planned`, and optionally `goal`, `depends_on` and `unattended`
(`internal/doc/doc.go` at f7e0140, where `Load` reads `cycles/*.yaml`;
`grove/specs/cycles/spec.md` at f7e0140, "Report a malformed cycle"). `goal`
is its only prose, one line, and plan writes it only if the owner gives one
(`internal/instruct/schemas/default/skills/plan.md` at f7e0140).

What a plan learns has nowhere to go. Plan ends by telling the owner the
members in build order and what each waits on
(`internal/instruct/requires/plan.md` at f7e0140), and writes none of what
it found. In another project of the owner's, private, the plans' pull
request descriptions carried it instead. One, planning three cycles over
ten changes, held why the cycles were cut where they were, four pairs of
changes that rewrite the same requirement, three changes left unplanned
and why, and a prototype's finding that unblocks a task. Two others said
which members were stale on their base, that every reference in their
tasks still held at the new base, and what was edited. None of it is in
the cycle files those pull requests added.

The owner, on that project: "I kept the conversation that committed the
plans open to remember some context about the plan."

Whoever builds the cycle never sees a pull request description. They do
see the cycle file: the implement instructions print it in full for the
cycle and for each member (`internal/instruct/instruct.go` at f7e0140,
`Cycle` and `InCycle`).

## Intended outcome

A cycle is a Markdown file, `<root>/cycles/<id>-<slug>.md`. It starts with
the same seven fields, as YAML between two lines of `---`, and below them
holds plan's notes: what the plan found on the target at `planned`, what
whoever builds the cycle should know above any one change, and why these
changes are planned together in this order and what was left out. The
notes are for reading. Grove parses none of them, gates on none of them,
and prints them to the builder with the file as it does today.

Only plan writes the notes, and re-planning brings them up to date. They
change no member: where the notes and a member's documents disagree, the
member's hold, and the implement instructions say so.

Grove reads no cycle from a `.yaml` file. `grove check` reports each one
left under `cycles/` and says how to convert it. `grove serve` shows a
cycle's fields as text and its notes rendered.

## Scope

### In scope

- The cycle file's name and shape: `.md`, its fields between two lines of
  `---`, the notes below. The fields, their checks and their messages stay
  as they are.
- `grove check` reporting a cycle file that does not start with its
  fields, and a `.yaml` file left under `cycles/`.
- The cycle template and the plan instructions: what the notes hold, that
  there may be none, and that re-planning brings them up to date.
- The implement instructions saying, with a cycle file that holds notes,
  that a member's documents hold over them.
- `grove serve` answering a cycle file with its fields as text and its
  notes rendered.
- The maintained specs, `docs/`, `dev/walkthrough.md`, the tests, and
  Grove's own ten cycle files, converted with no notes.

### Out of scope

- A document above cycles, for reasoning that belongs to no single cycle.
  Each cycle says why it follows the one before and what it left out.
- Anyone but plan writing the notes. Implement edits no cycle file, and
  propose, re-planning the cycle of a member it revised, still changes
  only `planned` (`grove/specs/propose/spec.md`, "Re-plan the cycle of a
  revised member", unchanged).
- Plan's verdict: the build order, the obstacles to `grove claim` and what
  the runner's checkout needs stay spoken, not written.
- Grove reading, checking or gating on anything in the notes
  (`CONSTITUTION.md` rule 10).
- Reading a `.yaml` cycle for a release, or converting one: there is no
  fallback and no command for it (`CONSTITUTION.md` rule 13).
- Cycle files in other projects: each project converts its own, as
  `docs/reference.md` will say.
- A new key in `/api/project`: the page reads the notes with the file.
- How `grove serve`'s history labels the commit that converts a project's
  cycle files.
- Grove naming requirements two changes both change: change 261118-40s1a.

## Decisions

The owner decided 1 to 5 while exploring, and answered the round before
writing, 6 to 20, with "yes" to each recommended answer.

1. The notes hold what planning found on the target, guidance to the
   builder above any one change, and some of why the cycle is shaped as it
   is. Asked what they went back to the planning conversation for, the
   owner said "mostly b and d, a bit of a", of: (a) why this cycle is
   shaped this way, (b) what planning found on the target, (c) what the
   run needs, (d) guidance to the builder above any one change.
2. The notes are for reading; no command decides anything from them.
   Asked whether Grove or a script should act on any of it, the owner said
   "mostly a, a bit of c", of: (a) all for reading, (b) addressed per
   member, (c) some of it should be Grove's own check. The check is change
   261118-40s1a.
3. Reasoning that spans cycles goes in the cycle it explains, not in a
   document above them: "mostly a", of: (a) in the cycle it explains, (b)
   one document for the planning session, (c) it stays in the pull request
   description. Reason: no new kind of document.
4. Only plan writes the notes: "yes only plan writes it". Reason: the
   staleness rule excuses a commit that also edits the cycle file
   (`grove/specs/cycles/spec.md`, "Count a member changed with its plan"),
   so a file edited during a build would excuse changes no plan checked.
5. The file is Markdown with its fields as YAML at the top, chosen over an
   optional `notes` field in the `.yaml`, a `.md` beside the `.yaml`, and
   a folder like a change's: "I think I prefer option 2." Reason: the data
   stays YAML read by the same parser, the prose is Markdown that
   `grove serve` and a code host render, and a mistake in the prose cannot
   make the fields unreadable.
6. No fallback for old files, as with `needs_human` (change 261114-yx2b7):
   the owner chose "a", Grove reads only `.md` and `grove check` reports
   each `.yaml`. Every `.yaml` under `cycles/`, done or not, is reported
   and `grove check` exits non-zero; Grove reads no cycle from it, so
   `grove changes`, `grove claim` and `grove serve` show nothing for it.
7. The fields are the same seven, unchanged, `goal` among them: it is
   still the one line that lists show.
8. The body is free Markdown that Grove never parses. It may be empty, and
   it has no title heading.
9. Plan writes notes only when there is something to say, of three kinds:
   what it found on the target at `planned`; what the builder should know
   above any one change; why these changes in this order, and what was
   left out. It never writes there anything that changes a member's
   requirements or tasks: that is a revision of the member.
10. `/grove-plan <cycle-id>` brings the notes up to date with the new
    `planned`. Propose, re-planning the cycle of a member it revised,
    still changes only `planned`. Reason: the propose spec stays as it is.
11. The implement instructions say that where the notes and a member's
    documents disagree, the member's hold. Reason: a note gone stale is
    then harmless.
12. `grove serve`'s pane shows a cycle's fields as text, as it shows YAML
    now, and renders the notes. `/api/project` gains no key.
13. Grove's own ten cycle files are converted in this change, with no
    notes. The owner converts the other project's; `docs/reference.md`
    says how.
14. No amendment to `CONSTITUTION.md`: rule 7 says "Changes, specs and
    cycles are files under one configured root" and names no extension.
15. The split of fields from notes sits in `internal/doc`, before the
    existing parser, and a problem is reported by its line in the file.
    The template becomes `internal/instruct/metadata/cycle.md`.
16. The two new messages, one sentence each (`CONSTITUTION.md` rule 14):
    "a cycle file is Markdown now; rename it to `<name>.md` and put its
    fields between two lines of ---" and "a cycle file starts with its
    fields between two lines of ---".
17. Tests follow new walkthrough steps in
    `internal/cli/walkthrough_test.go`, with one direct test of the split
    in `internal/doc`.
18. This change adds at most 150 words of instruction text to the 7,997
    counted at f7e0140 (`CONSTITUTION.md` rule 12).
19. Breaking: the commit is `feat!` with a `BREAKING CHANGE` footer saying
    what a project holding `.yaml` cycle files must do.
20. Test layout and private helpers are the implementer's to choose.

## Open questions

Neither needs an answer before review.

- Whether the sentence of decision 11 is printed for a cycle whose file
  holds no notes. The owner agreed the sentence, not when it shows. It
  decides how T-01 is built; the change specs say only for a file that
  holds notes, which `/grove-implement` builds and records for the
  reviewer.
- What the template says in place of notes, so that plan writes none when
  there is nothing to say. It decides how T-01 is built; `/grove-implement`
  writes its best wording within decision 18 and records it.

## Related documents

- Requirements: [cycles](specs/cycles/spec.md), [serve](specs/serve/spec.md),
  [implement](specs/implement/spec.md)
- Tasks: [tasks.md](tasks.md)
