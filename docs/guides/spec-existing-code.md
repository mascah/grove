# Specs of the code you already have

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

A new project has no specs. Grove writes a maintained spec only when it
archives a change, so what the code did before you started using Grove is
agreed nowhere. This recipe writes down what one part of the code already
does, as a change with no code to change, and one merge approves and
records it.

It needs nothing beyond proposing a change and building it on its own
branch (see the
[propose spec](https://github.com/mascah/grove/blob/main/grove/specs/propose/spec.md)
and the
[implement spec](https://github.com/mascah/grove/blob/main/grove/specs/implement/spec.md)).

## Steps

**Terminal:** start a branch for it.

```sh
git switch -c record-edit
```

**Chat:** ask for a change that only records what the code does. Name one
command or area, not the whole system.

```text
/grove-propose Record what `ledger edit` does today, as a change that only records existing behavior, with no code to change
```

The agent reads the code and asks you one round of questions. For each
behavior it finds, decide:

- **It is the rule.** The agent writes it as a requirement under ADDED.
  What the code does today goes under "Observed current behavior", with
  its source.
- **It is not what you want, but not now.** Say to keep it as a known
  discrepancy. The requirement says what you want, and the difference is
  written down, to go under the maintained spec's "Known discrepancies"
  when the change is archived.

The change has no tasks. If you decide the code should change now, that is
a task, and the change no longer only records behavior: propose it as an
ordinary change instead.

**Terminal:** commit the draft.

```sh
git add -A && git commit -m "Record what edit does"
```

**Chat:** build it on the same branch.

```text
/grove-implement <change-id>
```

With no tasks, the agent builds nothing. It runs the tests, lists each
scenario with the check that shows it, and archives the change, which
writes `grove/specs/<capability>/spec.md`.

**Terminal:** commit, then look at what merging would agree to.

```sh
git add -A && git commit -m "Archive the record of edit"
grove check
grove review
```

> Merging record-edit into main approves change <change-id> with its code.

Read the maintained spec it writes, then merge. One merge approves the
change and records the spec.

## What to expect

- **The agent leaves out what you might not want.** In the trial, its
  first answer for the traceback `ledger edit 99 --note x` ends in was to
  record it only as observed, with no rule. If you want a rule there, say
  what it is and that the code's behavior is a known discrepancy. The
  maintained spec then holds your rule, and lists the traceback under
  "Known discrepancies".
- **Name the spec if you care what it is called.** The agent picks the
  capability's ID, `ledger-edit` in the trial; give yours in the round.
- **`tasks.md` holds a note and no task.** `grove check` and archiving
  accept it.
- **No tests are added.** Implement checks each scenario by running the
  code once, and reports the one with a known discrepancy as failing. A
  test for recorded behavior is a task, so it is not part of this recipe.
- **Questions you leave open stay open.** In the trial, whether an older
  record with a retired category may be edited went to the maintained
  spec's "Open questions", for its own change later.
- **`grove review` says "with its code"**, though the change has none.
