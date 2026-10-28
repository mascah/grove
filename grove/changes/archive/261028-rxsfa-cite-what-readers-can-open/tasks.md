## T-01 — Every skill is told to cite what readers can open

Result: the opening every skill prints holds the principle the `instructions`
change spec requires, after the paragraph on drafts.

Files:
- `internal/instruct/instruct.go`

Interfaces: the `meanings` constant gains one paragraph.

Acceptance:
- `grove instructions explore`, `grove instructions propose` and
  `grove instructions implement <change-id>` each print the sentence agreed
  in the proposal before the skill's own instructions.
- `just check` passes.

## T-02 — A proposal carries a prototype's findings

Result: the propose instructions and the docs say what the MODIFIED
requirement "Propose uses a prototype's findings" says: the proposal stands
without the prototype branch and names it and its commit only as the source.

Files:
- `internal/instruct/requires/propose.md`
- `docs/concepts.md`
- `docs/getting-started.md`
- `dev/walkthrough.md`

Acceptance:
- `grove instructions propose` prints the prototype line as the requirement
  states it, including what to say when the experiment was not committed.
- The docs' and the walkthrough's prototype steps say the proposal states
  the findings and names the branch and commit as their source.
- `just check` passes within the 9,000-word budget.
