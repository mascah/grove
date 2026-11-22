# Changes to serve

The changes that change 261122-r8y7c proposes to the requirements of
[serve](../../../../specs/serve/spec.md), as it stands at 15816d7. This is
future work: nothing here describes the target branch until the change is
implemented and the maintained spec is reconciled.

## MODIFIED Requirements

### Requirement: History of the record

When the page asks for `/api/project`, the system shall list, newest
first, each commit on the target's first-parent history that approved a
change (added its folder under `changes/`), revised one (changed a file
in its folder under `changes/` that the commit did not add with the
folder), archived one (added its folder under `changes/archive/`),
planned a cycle (added its file), re-planned one (changed it) or edited a
maintained spec, giving the commit, its date, its subject and each such
mark with the change, cycle or spec it names.

#### Scenario: A change's life

- GIVEN the target's history holds a commit that added change A, a later
  one that changed A's `proposal.md`, a later one that added cycle K
  listing A, and a later one that moved A to `changes/archive/`
- WHEN the page asks for `/api/project`
- THEN those commits appear, newest first, marked archived A, planned K,
  revised A and approved A

#### Scenario: A change approved before its fields moved

- GIVEN a commit that added change A's `.grove.yaml` and `proposal.md`,
  and a later one that moved A's fields into its `proposal.md` and
  deleted its `.grove.yaml`
- WHEN the page asks for `/api/project`
- THEN the first is marked approved A and the later revised A

#### Scenario: A commit outside the record

- GIVEN a commit on the target that changed only `ledger/records.py`
- WHEN the page asks for `/api/project`
- THEN it does not appear in the history

### Requirement: A document rendered

When the page asks for `/api/doc` with a path under the root ending
`.md` or `.yaml` and a ref that is empty, `worktree` or a full commit, the
system shall answer with that file at the target, in the working tree or
at that commit, Markdown rendered to HTML with raw HTML dropped, the
fields a Markdown file starts with between two lines of `---` as text
above the rest, and YAML as text, and refuse any other path or ref.

#### Scenario: A proposal at the target

- GIVEN the target holds change A
- WHEN the page asks for A's `proposal.md` with an empty ref
- THEN the answer is HTML holding A's fields as text, each on its own
  line as the file has them, and below them A's title as the first
  heading

#### Scenario: A cycle file

- GIVEN cycle K's file holds its fields and, below them, notes with a
  heading and a list
- WHEN the page asks for K's file
- THEN the answer is HTML holding K's fields as text, each on its own
  line as the file has them, and below them the notes rendered, the
  heading as a heading and the list as a list

#### Scenario: A cycle file with no notes

- GIVEN cycle K's file holds only its fields
- WHEN the page asks for K's file
- THEN the answer is HTML holding K's fields as text and nothing below

#### Scenario: A document on a branch

- GIVEN branch `build-B` holds B's `implementation.md` as the target does
  not
- WHEN the page asks for it with the branch tip's commit as the ref
- THEN the answer is that branch's version
