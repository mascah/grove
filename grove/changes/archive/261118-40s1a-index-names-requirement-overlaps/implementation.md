# Implementation: The project index names requirements two changes both change

<!-- Progress on change 261118-40s1a, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — The index names a shared requirement, and plan tells the owner

Status: done

`doc.Share` and `doc.Change.Shares` (`json:"-"`, so `grove changes --json`
is unchanged); `doc.Load` fills it in `Collection.shares` from the ADDED,
MODIFIED and REMOVED blocks of each open change's capability files, by
capability and requirement name, leaving out a pair where one depends on
the other through `depends_on`, directly or through other changes, open or
archived. The index prints one line per share under the change. The plan
requirements gain one bullet (56 words; `just check` prints 8199
instruction words, 8143 when 261118-emse7 was archived).

Checks: `internal/doc/doc_test.go` `TestShares` (a MODIFIED/REMOVED pair,
a pair ordered through a third change, the same name in two capabilities),
passed. `internal/cli/walkthrough_test.go` `TestCoordinating`: with
recategorize and show-saved approved, `instructions plan "0.5 release"`
and `instructions propose` name each under the other; a cycle holding
both checks clean and `grove claim` claims it; with show-saved's
`depends_on: [recategorize]` committed, or recategorize archived, no line;
passed. By hand with `bin/grove` on a fresh fixture with two changes
modifying "Keep a category": the index shows the line under each, `grove
check` and `grove changes` print what they did before.

Check: the index line reads `also changed by <change-id>: <capability> "<requirement>"`, the tasks' suggested words, since the proposal's open question on the wording had no answer; the requirement name is printed with Go's `%q`, as `grove check` prints requirement names.
Check: 261118-emse7, built and archived on this branch (b26f250), changed `internal/doc/doc.go`, `internal/instruct/instruct.go`, `internal/instruct/requires/plan.md` and `internal/cli/walkthrough_test.go` since this change's base f7e0140; this builds on them as they are now. The new plan bullet sits beside emse7's bullet on notes; it does not ask plan to write the pair in the cycle's notes, which the change does not say.
Check: in the walkthrough test every change `propose` writes holds an ADDED "Edit older records" block, so delete and budgets now show an `also changed by` line for it too; that is the behavior the requirement asks for, and the test asserts only the "Keep a category" lines.

### T-02 — Docs and walkthrough

Status: done

`docs/planning.md`, "Plan a cycle": the agent tells you when two changes
it proposes change the same requirement, that the one built second is
stale on it once the first is archived, and that you choose; Grove
refuses neither. `docs/coordinating.md`, "When the plan is stale": the
project index names the pair, under each, before either is built, with
the line it prints and a link to the cycles spec's requirement.
`dev/walkthrough.md`, Coordinating changes step 4: with recategorize and
show-saved approved, `grove instructions plan "0.5 release"` and `grove
instructions propose` show the line under each, the plan requirements say
what to tell the owner, a cycle holding both checks clean and is claimed,
and the line goes with `depends_on` committed or recategorize archived;
its note on the project index says it names such a pair. The steps are
those `TestCoordinating` follows (T-01), passed.

Check: the walkthrough step for the pair sits in Coordinating changes step 4, where the fixture's two changes modifying "Keep a category" are both approved, not in Planning cycles, whose start has recategorize already archived; the tasks name no section.

## Scenarios

- cycles "Two changes modify one requirement": `internal/cli/walkthrough_test.go` `TestCoordinating` (plan and propose instructions name each under the other) and `internal/doc/doc_test.go` `TestShares`, passed
- cycles "One removes what the other modifies": `internal/doc/doc_test.go` `TestShares` (261001-aaaaa MODIFIED, 261001-bbbbb REMOVED), passed
- cycles "One depends on the other": `internal/cli/walkthrough_test.go` `TestCoordinating` (show-saved's `depends_on: [recategorize]` committed on main), passed
- cycles "Through another change": `internal/doc/doc_test.go` `TestShares` (261001-ccccc depends on 261001-ddddd, which depends on 261001-aaaaa), passed
- cycles "The same name in two capabilities": `internal/doc/doc_test.go` `TestShares` (261001-eeeee adds "Keep a category" to `reports`), passed
- cycles "One of the two is archived": `internal/cli/walkthrough_test.go` `TestCoordinating` (after recategorize is archived), passed
- cycles "The owner moves one to a later cycle": unverified; it is what a plan session does with the owner, and only the instruction text is tested (`TestCoordinating` asserts the plan requirements say to tell the owner and ask whether both stay)
- cycles "The owner keeps both": `internal/cli/walkthrough_test.go` `TestCoordinating` (a cycle listing both checks clean and `grove claim` claims it), passed; the session's listing both is unverified, as above
- cycles "No pair among the members": unverified; whether the agent stays silent is a session's behavior, judged by a person or a skill trial, not run here

## Reconciliation

- cycles: applied to grove/specs/cycles/spec.md
