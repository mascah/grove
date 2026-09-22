---
id: "G-027"
type: plan
title: "G-023/G-025 agent handoffs implementation plan"
status: current
formerly: "docs/plans/W-010-W-011-agent-handoffs.md"
work: ["G-023", "G-025"]
created: "2026-09-19T21:39:23Z"
updated: "2026-09-21T21:11:16Z"
---

# G-023/G-025 agent handoffs implementation plan

**Historical handoff, superseded for new assignments on 2026-09-20:** G-023 is
done and integrated. G-025's current scope and preparation instruction are in
[its record](G-025-shaping-entrypoint.md), revised under
[G-035](G-035-interactive-adoption.md). Do not execute this
combined task list again or load its predecessor research by default. Retain it
as evidence of the delivered G-023 design and the earlier G-025 proposal.

> For Fable: execute serially in one isolated worktree, using this repository's
> instructions and `superpowers:executing-plans` if available. No automatic merge.
> This is an implementation handoff, not implementation or harness-test evidence.

**Status, 2026-09-20:** the G-023 parts of Tasks 1–5 were carried out on branch
`worktree-W-010`, and the owner then revised the design before integration. The
[revision](#revision-2026-09-20-owner-decisions) below supersedes anything in
this plan that conflicts with it; the interface sections have been rewritten to
describe what is implemented. The task lists are kept as the original handoff
and were not used as a checkpoint: the
[dogfooding evidence](G-032-dogfood-review.md) says what was done
and observed. G-025's tasks are not started.

**Goal:** one short invocation shapes project records or carries selected work
through preparation, implementation, review, and a recoverable handoff.

**Specs:** [G-023](G-023-work-handoffs.md) owns work assignments;
[G-025](G-025-shaping-entrypoint.md) owns authoring instructions.
Read their acceptance and both linked predecessor reviews before implementation.

**Architecture:** `internal/handoff` assembles bounded, read-only context from
one checkout. `grove context` exposes it. Thin Claude/Codex skill adapters load
shared repository guides and request that context. Instructions own judgment;
the command does not authorize work, decide readiness, or execute an agent.

**Stack:** Go 1.26, existing project/repo packages, Goldmark v1.7.13 for Markdown
link extraction only. No provider SDK, prompt service, TUI library, or scheduler.

**Execution:** native serial implementation with independent review after the
context/CLI changes and once on the combined handoff. No fixed provider/model
policy and no automatic controller loop. Use foreground checks and checkpoint
Next/Evidence. Reassess if the chosen harness cannot supply an independent review;
report that limitation rather than pretending self-review is independent.

## Revision 2026-09-20: owner decisions

1. **Repository-local skills, workflow apart from policy.** The local Claude
   and Codex adapters are intentional while the workflow is dogfooded. Plugin
   distribution, global installation, and proving portability in another
   adopted repository are not acceptance requirements. One shared workflow
   (`docs/work-execution.md`) covers interpreting assignments, retrieving
   context, preparation, execution, review, checkpoints, blockers, and handoff.
   Repository development policy (restart history, Go verification, the
   `go run ./cmd/grove` invocation, branch conventions, coexistence with the
   predecessor) lives once in `AGENTS.md`. The workflow never requires the
   predecessor's `grove:work`, `grove:close`, or `grove:shape` skills, and the
   predecessor comparison is review evidence, not required reading. G-025's
   shaping guide should reuse the same split.
2. **Staged retrieval.** `context` reads the configuration, the selected
   records, and explicit includes in full, and lists everything else. The
   workflow starts from the selected record, includes the plan the record names
   when preparing or implementing, and must read blocking questions and
   prerequisites in full before implementing a unit. No filename heuristic
   decides which artifact is current. Plans and reviews stay ordinary linked
   files; there is no attachment schema and no retrieval framework beyond
   `show` and `--include`.
3. **Isolation before the first write.** The workflow inspects branches,
   worktrees, checkpoints, and the intended base, establishes or reuses the
   execution checkout, and verifies the records and inputs there, before it
   writes a plan, checkpoint, question, or record update.
4. **Two CLI defects fixed**: Markdown destination decoding before URL
   interpretation, and file identity tracked for every source and checked
   before another alias is charged.
5. **Evidence kept honest**: source inspection, automated tests, simulated
   agents, real harness execution, and owner acceptance are reported apart.

Out of scope for this revision: plugin distribution, nullsec migration,
attachment reorganization, G-025's implementation, supervised launching.

## Base, ownership, and boundaries

Inspected main `33bffb6`; the separate repair worktree was at `9e8430c` during
preparation. Its presence is not integration proof. Before Go implementation,
verify G-014–G-016 repair/review commits are integrated. Reuse the delivered
G-016 path helpers; do not copy unfinished branch code. The owner has scheduled
G-017 after repairs. Prefer this implementation after that handoff to avoid
concurrent edits to `internal/cli/cli.go`, README, and the brief. G-017 is not a
semantic prerequisite; if scheduling differs, coordinate ownership explicitly.

Keep G-023 and G-025 proposed until their implementation starts. Retain individual
acceptance/evidence even with a shared plan; no release umbrella or new metadata.
Preserve the existing explicit commands and G-017's default TUI invocation.
No agent launch, branch creation by `context`, automatic merge, claim, run schema,
timer, board mutation, global skill installation, or changes to sibling projects.

## Prepared interface decisions

### Entrypoints and authority

Proposed user invocations after implementation:

```text
Claude: /grove-shape <idea or existing work ID>
Claude: /grove-work G-023 G-025
Codex:  $grove-shape <idea or existing work ID>
Codex:  $grove-work G-023 G-025
CLI:    go run ./cmd/grove context G-023 G-025 --json
```

Use `.claude/skills/{grove-work,grove-shape}/SKILL.md` and
`.agents/skills/{grove-work,grove-shape}/SKILL.md`. Keep these four adapters short;
they point to `docs/work-execution.md` or new `docs/work-shaping.md`. Do not copy
the guides into adapters or install/rename the predecessor's `grove:work` plugin.
Use distinct hyphenated names; no symlink installation is needed.

Make these first adapters explicitly invoked: Claude frontmatter
`disable-model-invocation: true`; Codex `agents/openai.yaml` alongside each skill
contains `policy: {allow_implicit_invocation: false}`. This limits discovery,
not the agent's ordinary autonomy within an explicit assignment. A future
headless caller invokes the same skill or loads its shared guide explicitly.
Neither route maintains a second workflow prompt. Claude's current
[programmatic usage documentation](https://code.claude.com/docs/en/headless#create-a-commit)
supports user skill invocation in `-p` prompts; this is documented capability,
not an exercised Grove adapter. Verify the installed harness before claiming
that a particular invocation works.

Make interaction mode explicit in adapter inputs and forward it to `context`
and the shared guide. Examples:

```text
/grove-work G-030 --interaction interactive
claude -p "/grove-work G-030 --interaction headless"
```

These are intended interfaces after implementation, not commands to run during
this shaping task. Ordinary interactive invocations may omit the mode; headless
callers must supply it. Apply the same mode convention to `grove-shape`, alongside
its idea/research mandate. The guides own CLI usage, workflow, bounds and waits;
adapters own argument handling and loading. Thin adapters must not reduce the
preparation, review, recovery or reconciliation responsibilities of the guides.

The harness paths and invocation controls are documented in
[Claude skills](https://code.claude.com/docs/en/skills) and
[Codex skills](https://learn.chatgpt.com/docs/build-skills), inspected 2026-09-19.
No restart skill has yet been exercised in either harness. Do not modify user
settings or enable subagent features to make a test pass.

### Context command

```text
grove [--project DIR] context WORK_ID... [--json]
      [--interaction interactive|headless] [--max-bytes N] [--include PATH]...
```

Default interaction is interactive; this field records caller intent, not a
detected capability or permission. Require at least one explicit work ID and
reject duplicates, non-work IDs, and selectors. IDs are never inferred from a
title or chosen across branches. `--project` selects the only live checkout.
The skill resolves an intended cross-branch version through existing
versions/workspace operations before requesting context in its checkout.

The command is a context operation, not `grove work` or `grove launch`. It emits
the declared selection, a dependency-respecting execution order, observations,
and sources. The reusable guide supplies the actual work instructions. No
restart-specific documentation filename is baked into the generic CLI, and the
adapters include nothing up front (revised 2026-09-20; the first design had
them include the brief, AGENTS, record model, and guide). A project without
those files still supports ordinary context inspection.

Output facts and content (inclusion policy revised 2026-09-20, format 2):

- Root, selected IDs in requested order, execution order,
  interaction mode, and Git checkout/common-directory/full ref/HEAD if present.
  Detached HEAD has no invented branch. Outside Git use an explicit absent Git
  observation; an error inside a detected repository is not silently non-Git.
- **Full sources with exact revisions:** `grove.yaml`, the selected records, and
  explicit `--include PATH` files. Nothing else is read in full by default.
- **Record observations:** transitive `depends_on` records, questions blocking
  any selected/prerequisite work (resolved ones too), direct `relates_to` and
  `members` records of selected work, and records a selected body links to.
  Each row has ID, type, title, status, path, the revision the loader saw, the
  roles that list it, and whether its source is included. Relationships are
  not expanded recursively and never add to the selection.
- List prerequisite observations: whether selected, status, and question status.
  `done` never sets an integrated/ready boolean. Open blocking questions and
  undelivered external dependencies are visible even though context can be read.
  An abandoned prerequisite is not delivered. Member inclusion is not ordering.
- **Link observations:** every direct Markdown link in a selected work body,
  with the destination as written, the project path it resolves to, and a
  reason. Links are never opened: a listed path is not validated, not even for
  existence. Links inside included documents are not parsed.
- Explicit project-relative `--include PATH` files are required UTF-8
  regular-file sources of any type, included once however they are spelled. No
  URL, absolute path, Git metadata, sibling repository, or shell expansion is
  accepted. Adapter arguments are data, never interpolated into shell snippets.
- Retrieval uses existing commands: `show ID` for a listed record, `--include`
  for a listed file. The workflow decides which plan is current from the
  record's own prose, never from a filename.

Use stable Kahn ordering over selected IDs: edges reflect prerequisite reachability,
including paths through unselected prerequisites; ties follow requested order.
Do not add prerequisite IDs to the execution selection. Any unsupported gate or
scope ambiguity is for the guide to investigate, not an auto-acquisition policy.

An absent plan is not a structural error or readiness conclusion. Preserve the
record's exact prose and require the guide to inspect whether planning is needed.
A missing `--include` is an error naming the path. A link to a missing document
is not discovered, because links are not opened (revised 2026-09-20).

### Link and read boundary

Parse Markdown with Goldmark; inspect `*ast.Link`, not regex matches, so inline
code, fences, examples, images, and HTML do not become filesystem requests.
Support inline and resolved reference-style links, escaped destinations and
angle-wrapped destinations. Goldmark keeps the destination as written, so apply
Markdown backslash and entity decoding first, then decode URL paths once before
classification; report the original destination (fixed 2026-09-20). Resolve relative to the referring record, normalize separators
without trimming real filename whitespace, and reject NUL/invalid escapes.

Links outside the selected project, external URLs, and fragment-only links
are labelled references with no path, never followed or fetched. Every parsed
link is listed with a reason. Links inside included
documents are deliberately not traversed: output a fixed scope notice and the
guide must inspect those documents for additional required evidence. Use
`--include` to add a necessary local source; consult sibling projects only through
their own instructions/CLI. Completeness means this bounded inclusion contract,
not every fact an implementer could need. Fragments do not filter file content;
record them as references, without claiming anchor validation.

Reject symlink/nonregular source paths and traversal into `.git`, including a
linked-worktree Git file. Ordinary `../../docs/plan.md` links may normalize inside
the root; `--include` paths must themselves be clean project-relative paths.
Use Go `os.OpenRoot` for confinement, component checks and opened-file checks.
On supported Unix systems open sources with no-follow/nonblocking flags, then
require a regular descriptor before reading; a file swapped for a FIFO must not
hang. An escape through a concurrently replaced parent must never read outside
the root. Keep these helpers private to handoff; no unrelated reader rewrite.

The tagged [Goldmark module](https://raw.githubusercontent.com/yuin/goldmark/v1.7.13/go.mod)
declares Go 1.22; its [AST](https://raw.githubusercontent.com/yuin/goldmark/v1.7.13/ast/inline.go)
exposes link destinations. Pin v1.7.13 rather than writing a partial Markdown
parser. Only parse, never render HTML or execute embedded content.

Default source budget: 262144 bytes; `--max-bytes` accepts integers 1–8388608.
Count unique included source bytes, including config; this is not a token or
serialized-output limit. Every source registers its file identity, and an
already-included file is recognised before another name for it is read or
charged (fixed 2026-09-20). Bound each read before allocation and while reading;
report the source exceeding the remaining limit, with rerun guidance. Reject
invalid UTF-8. No silent truncation, eviction, summaries, or successful partial
bundles. Required-source errors produce no result stdout in either format.

Buffer the full result. Before return, reload and compare the config and entire
record set (including added/deleted records affecting blockers), recheck artifact
bytes/identity, root directory identity and Git checkout identity. Compare the
held root with the current root pathname so replacing/moving that directory
cannot return a now-misbound path. Refuse observed changes with refresh
guidance. Use no write lock or allocator helper. These are optimistic read checks,
not an atomic filesystem snapshot or a lease; execute/resume revalidates inputs.

Default text output is readable context, with source text clearly marked as data
and terminal controls escaped while retaining newlines/tabs. JSON contains exact
UTF-8 source strings and the same revisions. Source hashes always use original
bytes. Exit 0 means context assembled, not ready/accepted/authorized; exit 1 means
read/validation/change/output failure; usage is 2. Help requires no project/Git.

## Task 1: Assemble selection, relationships, and provenance (G-023)

Files: new `internal/handoff/context.go`, `selection.go`, `git.go`, and their
tests. Reuse `project.Load`, `project.Revision`, and the integrated repo path
helpers. Keep helpers independently testable; proposed exported interface:

```go
type Options struct {
    Interaction string
    MaxBytes    int
    Include     []string
}
type Source struct {
    Path, Revision, Content string
    Reasons                 []string
}
type Record struct {
    ID, Path, Type, Title, Status, Revision string
    Roles                                  []string
    Selected, Included                     bool
    Source                                 string // sources[] path holding it
}
type Requirement struct {
    Work, Prerequisite, Status string
    Selected                  bool
}
type Question struct {
    ID, Status string
    Blocks     []string
}
type Reference struct { From, Target, Path, Reason string }
type Git struct { Checkout, CommonDir, Ref, Head string }
type Bundle struct {
    FormatVersion                 int // JSON output version (2), not record schema
    Root, Interaction, ScopeNotice string
    Selected, Order               []string
    Git                           *Git
    Records                       []Record
    Requirements                  []Requirement
    Questions                     []Question
    Sources                       []Source
    References                    []Reference
    SourceBytes                   int
}
func Build(ctx context.Context, root string, ids []string, opts Options) (*Bundle, error)
func Text(bundle *Bundle) []byte
```

Define snake_case JSON tags, emit arrays as arrays, and document the final output
shape in README. Source order is normalized path order; sort observations and
references by their identity fields, retaining requested/order lists separately.
No timestamp/random bundle ID that makes unchanged output differ.

- [ ] Write a pure selection test: requested `[W-002,W-001,W-003]`, W-002 depends
  on W-001, W-003 independent; order is `[W-001,W-002,W-003]`. Add a dependency
  through an unselected node and prove it is context only. Test member order,
  shared prerequisites, duplicate/missing IDs, unrelated cycles caught by Load,
  open/resolved blocking questions, and abandoned/done external dependencies.
- [ ] Run `go test ./internal/handoff -run Selection -count=1`; observe the
  intended failing assertions before implementing deterministic selection.
- [ ] Capture Git identity with read-only commands and exact path terminators.
  In Go use argv slices and the delivered path helper; never concatenate shell
  commands or call `repo.CommonDir`. Disable Git optional index refresh for any
  status read; no status read is needed merely to return source revisions.
- [ ] Fixture-test no Git, attached/detached and linked checkouts, nested project,
  newline/tab paths, unavailable Git in a detected repository, and a changed ref
  between reads. Missing HEAD in an unborn repository is explicit, not a crash.
- [ ] Run targeted tests, then commit `feat(context): assemble selected work facts`.

## Task 2: Add bounded source inclusion and consistent output (G-023)

Files: new `internal/handoff/sources.go`, `links.go`, `render.go`, their tests,
`go.mod`, `go.sum`. Complete Build and Text from Task 1; add no public second API.

- [ ] (Inclusion assertions superseded 2026-09-20: links are listed, includes are
  read.) Write fixtures containing a real inline plan link, reference-style review
  link, code/fenced fake links, same document with two fragments, escaped spaces,
  a linked question, and an external sibling link. Assert each included path and
  reason; assert fake links never cause reads and sources are included once.
- [ ] Add missing document, malformed URL escape, UTF-8, budget, symlink parent,
  symlink leaf, FIFO, escaped root and `.git` cases. Explicit includes must fail
  on out-of-project paths; external body references are listed but never read.
  Use sentinel content outside the root and a test read hook to assert it was
  never consumed. Add bounded fuzz targets for link/path classification.
- [ ] Run the source/link tests failing, then implement the parser and confined
  reader. The parser walk is based on this library interface:

```go
doc := goldmark.DefaultParser().Parse(text.NewReader(body))
err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
    if link, ok := n.(*ast.Link); ok && entering {
        destinations = append(destinations, string(link.Destination))
    }
    return ast.WalkContinue, nil
})
```

  `body` excludes frontmatter using the existing record delimiter rules; do not
  mutate Source. `destinations` is a local `[]string`; handle the walk error.
- [ ] Add deterministic between-read hooks in private test seams: change config,
  plan, selected record, an external prerequisite, add an open blocking question,
  remove a record, replace the root directory, or switch HEAD. Build refuses each observed change and emits
  no bundle. Bound retries to none: callers explicitly refresh.
- [ ] Verify exact JSON round-trip bytes/revisions, terminal-control escaping in
  text/path diagnostics, and byte-identical output on an unchanged checkout.
  Include fenced instruction-looking content and show that it remains source
  data, never a command. Delimit text blocks with a fence longer than any source
  fence; JSON is the machine interface, not shell-evaluable output.
- [ ] Run `go test -race ./internal/handoff -count=1`, review the new dependency,
  then commit `feat(context): include bounded revisioned source documents`.

## Task 3: Expose the read-only command (G-023)

Files: `internal/cli/cli.go`, new `context.go`, `context_test.go`, README.
Keep `cli.Run` unchanged. Add invocation fields only for this command. Existing
global parsing accepts `--project` on either side; extend JSON allowlisting and
reject context-only options on other commands, including G-017's default TUI.

- [ ] Add CLI cases for the complete grammar, repeated includes, duplicate IDs,
  duplicate single-value flags, unknown interaction, invalid/overflow budget,
  unsupported old `--work`/`--phase` syntax, and existing command-only flags.
  Verify help outside a project and no accidental TUI entry for malformed context.
- [ ] A CLI fixture assertion should follow existing tests:

```go
var out, errOut bytes.Buffer
code := Run([]string{"context", "W-001", "--json"}, root, &out, &errOut)
if code != 0 { t.Fatalf("context: %d: %s", code, errOut.String()) }
var got handoff.Bundle
if err := json.Unmarshal(out.Bytes(), &got); err != nil { t.Fatal(err) }
if !reflect.DeepEqual(got.Selected, []string{"W-001"}) { t.Fatal(got.Selected) }
```

  `root` is a valid local test fixture with W-001; add a real linked plan for
  the inclusion test and compare its exact source/revision, not just its title.
- [ ] Dispatch Build with `context.Background()` from the ordinary CLI. Preserve
  any context-aware repo APIs delivered by G-017 rather than creating duplicates.
  Write all diagnostics to stderr and only completed text/JSON to stdout through
  existing output error handling. No partial stdout for missing/oversized sources.
- [ ] Hash files in the project, Git common directory, and linked checkout before
  and after success/refusal/blocked-context reads, with no preexisting `grove`
  coordination folder. Assert no changed files or new state. Test failing writers.
- [ ] Run targeted CLI/handoff tests and obtain an independent boundary review;
  commit `feat(cli): expose read-only work context`. Fix consequential findings
  before the guides depend on the output contract.

## Task 4: Deliver the shared guides and thin adapters (G-023/G-025)

Files: update `docs/work-execution.md`; add `docs/work-shaping.md`; the four
SKILL.md files and two Codex `agents/openai.yaml` files listed above; README and
a short AGENTS discovery pointer. Do not turn the plan itself into a skill.

- [ ] Author the shared guides from each record's acceptance and the predecessor
  responsibility mapping. Work guide: inspect actual state and existing owners,
  assemble context, resolve required omissions, prepare missing plan, choose
  isolated execution base, implement serially, review, reconcile and hand off.
  The initial supported path is native serial execution, with independent
  review at consequential boundaries and at the final combined revision.
  Fixed model assignments and the old size-triggered controller are deferred.
- [ ] Specify bounded correction: at most three fix/review rounds per review
  gate; after exhaustion preserve changes/findings and hand off. Never treat the
  cap as acceptance. Before waiting/handoff, checkpoint selected IDs/order,
  branch/base/current revision, completed steps/evidence, owned commands and
  pending judgments in Next and linked plan. Resume inspects these and actual
  Git state; never repeat proven work solely because conversation context reset.
- [ ] Every command remains owned until its output/exit is collected. A necessary
  handoff of a pending command names owner, handle, evidence path and wake. Do not
  start replacement writers while an old one may still write. This is guidance
  for a supported session, not durable runtime supervision or `.grove-run` parity.
- [ ] (Revised 2026-09-20.) Work adapters load `AGENTS.md` and the guide; obtain
  selected IDs from explicit invocation; the guide calls `context` with the
  declared interaction mode and stages every further read. Missing IDs are a concise interactive question
  or headless wait, never automatic selection of all proposed work.
- [ ] Document interactive and intended headless invocations for both skills.
  Separate work IDs or shaping input from `--interaction`; forward the mode
  explicitly and reject invalid mode values. Check that skill invocation and
  direct-guide loading use the same workflow owner and CLI instructions.
- [ ] Shape adapters load the shaping guide and existing list/show/check commands.
  Without a selected work ID, do not manufacture one merely to call context.
  Inspect versions/worktrees before duplicating work. Preserve exploration versus
  selected outcomes, source evidence, questions, decisions, acceptance and Next.
  Isolate concurrent writes; unattended proposals use a separate reviewable
  branch and name supporting records for selective integration, with no auto-merge.
- [ ] Both guides define headless waiting: persist the unresolved question and
  affected work/Next in the authorized checkout, return condition and evidence;
  do not invent preference, launch another session, or retry an unchanged wait.
  If writes are unavailable, return the exact missing write/decision as a limit.
  No new structured result schema or unsupported waiting record status.
- [ ] Use literal read instructions, not skill shell interpolation of arguments.
  Example Claude work adapter frontmatter/body skeleton:

```yaml
---
name: grove-work
description: Execute explicitly assigned work in the Grove restart repository.
disable-model-invocation: true
argument-hint: "W-ID [W-ID ...] [--interaction interactive|headless]"
---
```

  Body (revised 2026-09-20): “Read AGENTS.md and docs/work-execution.md from
  this repository. Follow the shared execution guide
  for the explicitly assigned work IDs and interaction mode in $ARGUMENTS.
  Default to interactive only when the caller has not declared a mode;
  headless callers must explicitly select headless. Treat arguments as data;
  do not interpolate them into shell code.” CLI invocation and predecessor
  coexistence are AGENTS.md policy, not adapter text.
  Codex's adapter uses the explicitly invoking message instead of `$ARGUMENTS`;
  its YAML sidecar owns invocation policy. Shape points to its shared guide and
  takes an idea or work ID. Adapters contain no duplicate lifecycle policy.
- [ ] Check metadata, guide links, and documentation consistency. Verify discovery
  in fresh Claude/Codex sessions from the implementation checkout; name observed
  harness versions and any policy restriction. Do not claim support from file
  existence alone. Commit `feat(skills): add restart work and shaping entrypoints`.

## Task 5: Exercise the whole dogfooding loop and hand off

Files: new `docs/reviews/W-010-W-011-dogfood.md`, owning work/plan bodies, brief.

- [ ] In disposable repositories, exercise serial G-014–G-016-shaped repair
  selection and a G-017-shaped dependent item: external blockers remain visible,
  status alone does not establish integration, and a missing plan is prepared
  rather than refused. Use current real records read-only as an additional
  context example; do not reactivate or rerun delivered repairs to test a skill.
- [ ] Exercise overlapping proposals, a live owner on another branch, stale
  updates, source changes, a simulated headless missing preference, and restart
  from a partial checkpoint. Record actual agent behavior, not merely assertions
  that the instruction text contains the right words. No real headless launcher.
- [ ] Trace both skills' interactive and headless paths to their shared guides
  and CLI operations. Exercise mode forwarding and a simulated headless wait.
  Record source checks, simulations and actual harness trials separately; do
  not claim `claude -p` behavior from an interactive trial. Any unavailable
  headless trial remains an explicit limit for the later supervised-run work.
- [ ] Dogfood shaping on the next actual requirements conversation; retain only
  justified work/questions/decisions. Dogfood work on a real authorized bounded
  assignment, which may be a remaining task of this implementation after the
  entrypoint exists. Capture invocation, included source revisions, observed
  preparation/review/handoff, and failures. Distinguish human feedback from tests.
- [ ] If fresh-harness or human interaction is unavailable, finish code/checks,
  retain the specific trial as unverified, and leave acceptance open. Do not run
  paid nested agents or enable settings merely to manufacture evidence.
- [ ] Run `go test -count=1 ./...`, `go test -race -count=1 ./...`, `go vet ./...`,
  gofmt, `go run ./cmd/grove check`, links, and independent combined review.
  With G-017 integrated, verify the new explicit command does not disturb bare
  invocation/help/nonterminal behavior using its existing tests.
- [ ] Reconcile each unit independently through current CLI/body edits. Return
  branch/worktree, commits, actual acceptance evidence, remaining trials and
  exact integration action; no merge/push. Commit the evidence with the work.

## Acceptance trace and next investment

G-023 acceptance 1–3: Tasks 1–3 and fixtures in Task 5; 4, 6–8: Task 4 plus
observed scenarios in Task 5; 5: real assignment trial. G-025 acceptance 1 and 5:
Task 4 discovery trials; 2–4 and 6: guide behavior and Task 5 fixtures/real conversation.

Review focus: context mistaken for authority; a link escaping project ownership;
an omitted/changed prerequisite or plan; record Done mistaken for integration;
and file-existence tests mistaken for demonstrated harness behavior. The owning
tasks above exercise each boundary. After this loop works, shape one manually
launched supervised headless attempt before schedules or board launch controls.
