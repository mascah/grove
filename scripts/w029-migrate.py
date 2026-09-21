#!/usr/bin/env python3
"""W-029's one-time reconciliation. Usage: w029-migrate.py PROJECT GROVE_BINARY

Converts every record and legacy plan/review to a neutral ID in document-date
order, moves the brief, repairs references, writes the mapping page, and
audits. Run from a clean, committed checkout; see the W-029 plan for recovery.
Deleted with schema 2's support, which `convert ID` needs.
"""
import datetime, json, os, posixpath, re, subprocess, sys

ROOT, GROVE = os.path.abspath(sys.argv[1]), os.path.abspath(sys.argv[2])
os.chdir(ROOT)

# Lines whose typed IDs are a fixture's or a quotation's, not this
# repository's records: (old path, substring of the line, IDs kept). Settled in
# rehearsal.
EXCEPTIONS = [
    ('docs/plans/W-001-inspection.md', 'without a closing `---`) makes both', ['W-001']),
    ('docs/plans/W-002-create.md', 'Allocator tests: floor 6 from committed', ['W-003', 'W-005']),
    ('docs/plans/W-002-create.md', '7 from the worktree, Q starts at 1; a counter of 2', ['W-004']),
    ('docs/plans/W-002-create.md', 'normal) with a test that hides a', ['W-030']),
    ('docs/plans/W-002-create.md', 'a persistence-failure test proves no ID', ['W-007']),
    ('docs/plans/W-003-update.md', 'workflow fixture creates', ['W-002']),
    ('docs/plans/W-004-W-005-coordination.md', 'directory is refused. `TestJointWorkflow`', ['W-001']),
    ('docs/plans/W-004-W-005-coordination.md', 'choose the live feature version explicitly, resolve its', ['W-001']),
    ('docs/plans/W-006-workspace-provenance.md', 'res := mustInspect(t, filepath.Join(root, "sub")', ['W-001']),
    ('docs/plans/W-006-workspace-provenance.md', 'res, err := inspect(root,', ['W-001']),
    ('docs/plans/W-007-preserve-updates.md', '[]Field{{"title", "New"}})', ['W-001']),
    ('docs/plans/W-007-preserve-updates.md', 'type: work, title: T, status: proposed, kind: fix', ['W-001']),
    ('docs/plans/W-007-preserve-updates.md', 'nil, "kind", "size")', ['W-001']),
    ('docs/plans/W-007-preserve-updates.md', '_, err := Apply(root, Request{ID:', ['W-001']),
    ('docs/plans/W-008-git-paths.md', 'In an allocation fixture, commit', ['W-001']),
    ('docs/plans/W-010-W-011-agent-handoffs.md', 'Write a pure selection test: requested', ['W-001', 'W-002', 'W-003']),
    ('docs/plans/W-010-W-011-agent-handoffs.md', 'independent; order is', ['W-001', 'W-002', 'W-003']),
    ('docs/plans/W-010-W-011-agent-handoffs.md', 'code := Run([]string{"context",', ['W-001']),
    ('docs/plans/W-010-W-011-agent-handoffs.md', 'reflect.DeepEqual(got.Selected', ['W-001']),
    ('docs/plans/W-010-W-011-agent-handoffs.md', '`root` is a valid local test fixture with', ['W-001']),
    ('grove/plans/P-002-w-029-reconciliation.md', 'for fixture IDs that coincide with real ones', ['W-001']),
    ('docs/reviews/2026-09-19-integrated-cli.md', 'valid Grove records. From main', ['W-001']),
    ('docs/reviews/2026-09-19-integrated-cli.md', 'type: work, title: T, status: proposed, kind: fix', ['W-001']),
    ('docs/reviews/2026-09-19-integrated-cli.md', '--expect <current> --unset kind --unset size` fails', ['W-001']),
    ('docs/reviews/2026-09-19-W-010-dogfood.md', '- Selection: requested `[', ['W-001', 'W-002', 'W-003']),
    ('docs/reviews/2026-09-19-predecessor-work.md', 'context with budget 10000. Those work IDs belong', ['W-002']),
    ('docs/reviews/2026-09-19-shaping-and-runner-evidence.md', '--phase shape --budget 7000`.', ['W-003']),
    ('docs/reviews/2026-09-19-shaping-and-runner-evidence.md', 'records a prior Bench reuse assessment', ['W-003']),
    ('docs/reviews/2026-09-19-board-W-009.md', "proposed under main's title; `b` then shows it active", ['W-001']),
    ('docs/reviews/2026-09-19-board-W-009.md', 'under the feature title with', ['W-002']),
    ('docs/reviews/2026-09-19-board-W-009.md', 'changed since it was selected` stays through', ['W-001']),
    ('docs/reviews/2026-09-20-card-lineage-W-012.md', '`TestHistoryFollowsOneLineOfCommits` builds a real', ['W-001']),
    ('docs/reviews/2026-09-20-W-011-shaping.md', '', ['W-029', 'W-030']),
    ('grove/reviews/R-001-w-019-knowledge-records.md', 'rule by mutation: every one fails a test', ['T-001']),
    ('grove/work/W-007-preserve-updates.md', 'reservation consumed (the next record is', ['W-003']),
    ('grove/questions/Q-001-branch-versions.md', 'the owner selected: "Group under', ['W-003']),
    ('grove/decisions/D-003-allocator-mechanism.md', 'with "I accept', ['D-003']),
    ('grove/work/W-012-card-lineage.md', 'e147e43 09-19 16:07', ['W-006']),
    ('grove/decisions/D-002-sequential-ids.md', '`64a0d76bab9f56dc32e6`', ['W-001']),
    ('grove/decisions/D-002-sequential-ids.md', '`ab03a64356d2ebd9b5e4`', ['Q-001']),
    ('grove/decisions/D-002-sequential-ids.md', '`209c4fa50e8d0b181e86`', ['D-001']),
]
# Files that describe this migration: their path literals stay as written.
LITERAL_PATHS_KEPT = {"grove/work/W-029-migrate-knowledge.md", "grove/plans/P-002-w-029-reconciliation.md",
                      "grove/decisions/D-006-stable-knowledge.md"}
SLUGS = {"P-001": "flexible-records-plan", "P-002": "reconciliation-plan",
         "R-001": "knowledge-records-review", "R-002": "flexible-records-review"}
REVIEW_WORK = {"docs/reviews/2026-09-19-integrated-cli.md": ["W-003", "W-004", "W-005", "W-006", "W-007", "W-008"],
               "docs/reviews/2026-09-19-predecessor-work.md": ["W-010"],
               "docs/reviews/2026-09-19-shaping-and-runner-evidence.md": ["W-010"]}
REVIEW_RELATES = {"docs/reviews/2026-09-20-direction-evaluation.md": ["D-004"]}
# Not part of a branch name, filename or path; "W-010/W-011" is two IDs.
TYPED = re.compile(r"(?<![\w-])(?<![^\d]/)([WQDTPR]-\d{3,})(?!\w|-(?!shaped))")
TYPED_ANY = re.compile(r"[WQDTPR]-\d{3,}")


def run(*args, **kw):
    return subprocess.run(args, check=True, text=True, capture_output=True, **kw).stdout


def grove(*args):
    return run(GROVE, "--project", ROOT, *args)


def utc(stamp):
    return datetime.datetime.fromisoformat(stamp.replace("Z", "+00:00")).astimezone(datetime.timezone.utc)


def markdown_files():
    for folder, dirs, files in os.walk("."):
        dirs[:] = [d for d in dirs if d != ".git" and posixpath.join(folder, d) != "./.claude/worktrees"]
        for name in files:
            if name.endswith(".md"):
                yield posixpath.normpath(posixpath.join(folder, name))


def sources():
    found = []
    for path in markdown_files():
        text = open(path).read()
        if path.startswith("grove/"):
            ident = re.search(r'^id: "([^"]+)"', text, re.M).group(1)
            created = re.search(r'^created: "([^"]+)"', text, re.M).group(1)
            found.append((utc(created), ident, path, None))
        elif path.startswith(("docs/plans/", "docs/reviews/")):
            added = run("git", "log", "--diff-filter=A", "--follow", "--format=%aI", "--", path).split()[-1]
            found.append((utc(added), path, path, text))
    return sorted(found, key=lambda s: (s[0], s[1].encode()))


def convert_all():
    mapping = []  # {from, from_path, id, path, date}
    for date, key, path, text in sources():
        if text is None:
            args = ["convert", key] + (["--slug", SLUGS[key]] if key in SLUGS else [])
        else:
            kind = "plan" if path.startswith("docs/plans/") else "review"
            title = re.search(r"^# (.+)$", text, re.M).group(1)
            stem = re.sub(r"^\d{4}-\d\d-\d\d-", "", posixpath.basename(path)[:-3])
            slug = re.sub(r"-+", "-", TYPED_ANY.sub("", stem)).strip("-") + "-" + kind
            args = ["convert", path, "--type", kind, "--title", title, "--slug", slug]
        entry = json.loads(grove(*args))
        entry["date"] = date.strftime("%Y-%m-%dT%H:%M:%SZ")
        mapping.append(entry)
        if text is not None:
            os.remove(path)
    return mapping


def set_fields(mapping, ids):
    for entry in mapping:
        old = entry["from"]
        if not old.startswith("docs/"):
            continue
        work = REVIEW_WORK.get(old) or TYPED_ANY.findall(posixpath.basename(old))
        field, values = ("work", work) if work else ("relates_to", REVIEW_RELATES[old])
        revision = json.loads(grove("show", entry["id"], "--json"))["revision"]
        grove("update", entry["id"], "--expect", revision, "--set", field + "=" + json.dumps([ids[v] for v in values]))


def repair(paths, ids, report):
    back = {new: old for old, new in paths.items()}
    by_length = sorted(paths, key=len, reverse=True)

    def link(match, old_file, new_file):
        dest = match.group(1)
        if re.match(r"[a-zA-Z][a-zA-Z0-9+.-]*:|#", dest):
            return match.group(0)
        target, _, fragment = dest.partition("#")
        resolved = posixpath.normpath(posixpath.join(posixpath.dirname(old_file), target))
        relative = posixpath.relpath(paths.get(resolved, resolved), posixpath.dirname(new_file) or ".")
        if target.endswith("/"):
            relative += "/"
        return "](" + relative + ("#" + fragment if fragment else "") + ")"

    for new_file in markdown_files():
        old_file = back.get(new_file, new_file)
        skip = [(needle, kept) for path, needle, kept in EXCEPTIONS if path == old_file]
        lines = open(new_file).read().split("\n")
        for n, line in enumerate(lines):
            if line.startswith("formerly:"):
                continue
            line = re.sub(r"\]\(([^)\s]+)\)", lambda m: link(m, old_file, new_file), line)
            if old_file not in LITERAL_PATHS_KEPT:
                for old in by_length:
                    line = line.replace(old, paths[old])
            kept = {i for needle, some in skip if needle in line for i in some}
            if kept & set(TYPED.findall(line)):
                report.append(f"{new_file}:{n + 1}: {', '.join(sorted(kept))}")
            line = TYPED.sub(lambda m: m.group(1) if m.group(1) in kept else ids.get(m.group(1), m.group(1)), line)
            lines[n] = line
        open(new_file, "w").write("\n".join(lines))


def mapping_page(mapping, report):
    grove("new", "page", "Identity and path migration map", "--slug", "migration-map")
    page = next(p for p in markdown_files() if p.endswith("-migration-map.md"))
    front = open(page).read().split("\n---\n")[0] + "\n---\n"
    rows = "\n".join(f"| {e['from'] if not e['from'].startswith('docs/') else '-'} | `{e['from_path']}` | {e['id']} | `{e['path']}` | {e['date']} |" for e in mapping)
    excepted = "\n".join("- `" + r + "`" for r in report) or "- None."
    open(page, "w").write(front + f"""
W-029 reconciled every record, legacy plan and legacy review into neutral IDs
directly under `grove/` on 2026-09-21. This page is the one durable mapping
from an old ID or path to its current counterpart. To read an old commit, use
the old ID or path there with that commit's own CLI (`go run ./cmd/grove`);
to find what it became, look it up here or search `formerly:`. Numbers follow
document date (a record's `created`, a legacy document's first commit), ties
broken on the old ID or path; the order conveys no authority or priority.

| Old ID | Old path | ID | Path | Document date |
| --- | --- | --- | --- | --- |
{rows}
| - | `docs/restart-brief.md` | - | `grove/brief.md` | the brief, not a record |

## Kept in place

- `docs/prompts/*.txt`: three spent handoff prompts, historical evidence that
  the handoff work links. Their old IDs and paths are literals of their time.
- `README.md`, `AGENTS.md`, the skill adapters, `docs/record-model.md`,
  `docs/work-execution.md` and `docs/work-shaping.md`: functional homes, with
  references updated.

## Historical literals

Old identifiers remain in `formerly` fields, in this table, in Git branch
names such as `worktree-W-012` and in commit messages, which name what existed
then. The typed IDs named on these lines were left as written: they belong to
test fixtures, disposable trial clones or the predecessor project rather than
to this repository's records, or they sit in a verbatim quotation (the
owner's words, a commit subject, the 2026-09-19 renumbering table). Every
other typed ID that had a counterpart was rewritten, command transcripts
included:

{excepted}

A typed ID with no row in the table (`W-014` to `W-017`, `W-090`, `Q-002` and
the like) never named a record here: it is a fixture's or the predecessor's.
Records and evidence dated before this migration also name the folders of
their time (`grove/work/`, `docs/plans/`, `docs/reviews/`) where they describe
what was then true; nothing lives there now.
""")
    return page


def audit():
    broken, leftovers = [], []
    for path in markdown_files():
        for n, line in enumerate(open(path).read().split("\n"), 1):
            for dest in re.findall(r"\]\(([^)\s]+)\)", line):
                target = dest.partition("#")[0]
                there = posixpath.normpath(posixpath.join(posixpath.dirname(path), target))
                if target and not re.match(r"[a-zA-Z][a-zA-Z0-9+.-]*:", dest) and \
                        not there.startswith("../") and not os.path.exists(there):  # siblings are not checked
                    broken.append(f"{path}:{n}: {dest}")
            if path.endswith("-migration-map.md") or line.startswith("formerly:"):
                continue
            line = re.sub(r"worktree-[\w-]+|worktrees/[\w-]+", "", line)  # branch and worktree names
            if TYPED_ANY.search(line) or re.search(r"docs/(plans|reviews|restart-brief)|grove/(work|decisions|questions|terms|plans|reviews)/", line):
                leftovers.append(f"{path}:{n}: {line.strip()}")
    return broken, leftovers


def finish():
    broken, leftovers = audit()
    print(f"{len(broken)} broken links; {len(leftovers)} lines with typed IDs or old paths")
    for line in broken:
        print("BROKEN", line)
    for line in leftovers:
        print("OLD", line)


if sys.argv[3:] == ["--audit"]:  # after the hand edits
    print(grove("check"), end="")
    finish()
    sys.exit()
config = open("grove.yaml").read()
open("grove.yaml", "w").write(config.replace("schema_version: 2", "schema_version: 3"))
grove("check")
mapping = convert_all()
ids = {e["from"]: e["id"] for e in mapping if not e["from"].startswith("docs/")}
set_fields(mapping, ids)
os.rename("docs/restart-brief.md", "grove/brief.md")
config = open("grove.yaml").read()
open("grove.yaml", "w").write(config.replace("brief: docs/restart-brief.md", "brief: grove/brief.md"))
for folder in ["docs/plans", "docs/reviews"] + ["grove/" + d for d in ("work", "questions", "decisions", "terms", "plans", "reviews")]:
    os.rmdir(folder)  # fails loudly if anything unaccounted for remains
paths = {e["from_path"]: e["path"] for e in mapping}
paths["docs/restart-brief.md"] = "grove/brief.md"
report = []
repair(paths, ids, report)
page = mapping_page(mapping, report)
print(grove("check"), end="")
json.dump(mapping, open(os.path.join(os.environ["W029_OUT"], "w029-mapping.json"), "w"), indent=1)
unused = [e for e in EXCEPTIONS if not any(r.startswith(paths.get(e[0], e[0]) + ":") for r in report)]
assert not unused, unused
print(f"{len(mapping)} converted; map {page}")
finish()
