# Changes to serve

The changes that change 261118-emse7 proposes to the requirements of
[serve](../../../../specs/serve/spec.md), as it stands at f7e0140. This is
future work: nothing here describes the target branch until the change is
implemented and the maintained spec is reconciled.

## MODIFIED Requirements

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
- THEN the answer is HTML with A's title as its first heading

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

#### Scenario: A change's metadata

- GIVEN the target holds change A
- WHEN the page asks for A's `.grove.yaml`
- THEN the answer is the file's text

#### Scenario: A document on a branch

- GIVEN branch `build-B` holds B's `implementation.md` as the target does
  not
- WHEN the page asks for it with the branch tip's commit as the ref
- THEN the answer is that branch's version

#### Scenario: Raw HTML in a document

- GIVEN a proposal that contains a `<script>` element
- WHEN the page asks for it
- THEN the answer holds no `<script>` element

#### Scenario: A path outside the root

- WHEN the page asks for `../grove.yaml`
- THEN the answer is a refusal, and no file is read
