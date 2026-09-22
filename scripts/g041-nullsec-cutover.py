#!/usr/bin/env python3
"""G-041's one-time nullsec cutover. Usage:

    g041-nullsec-cutover.py NULLSEC GROVE_BINARY OUT_DIR [--audit]

Initializes the new Grove in a clean, committed nullsec checkout, converts
every predecessor record and legacy plan in a fixed order, sets fields,
moves the brief and plan evidence, repairs references in every tracked text
file, writes the map page and audits. OUT_DIR receives the mapping JSON and
the list of rewritten lines. See the G-041 plan (G-091) for the rules and
recovery. Not product code: deleted once the cutover is merged.
"""
import json, os, posixpath, re, subprocess, sys

ROOT, GROVE, OUT = (os.path.abspath(a) for a in sys.argv[1:4])
os.chdir(ROOT)
OLD = "docs/grove"
TODAY = "2026-09-22"

# Lines whose old IDs are a verbatim quotation, not a reference to rewrite:
# (old path of the file, substring of the line, IDs kept). Settled in rehearsal.
EXCEPTIONS = [
    (OLD + "/history/work/W-005-v4-release.md", "belt hidden until arrival works, close", ["W-005"]),  # the builder's words
]
ACRONYMS = {"npc": "NPC", "eve": "EVE", "ui": "UI", "e2e": "E2E", "ts": "TS", "v1": "V1", "v2": "V2",
            "v3": "V3", "v4": "V4", "webgpu": "WebGPU", "sqlite": "SQLite", "json": "JSON", "pve": "PvE",
            "babylonjs": "Babylon.js", "3d": "3D"}
# A bare old ID, not inside a branch name, filename or path segment.
OLD_ID = re.compile(r"(?:(?<![\w-])|(?<=\bpre-))(?<![^\d]/)((?:W-\d{3}|D-\d{4}))(?![\w]|-(?:[a-z]|T\d))")  # W-019-T3 is a branch
OLD_ANY = re.compile(r"\b(?:W-\d{3}|D-\d{4})\b")
WIKILINK = re.compile(r"\[\[([^\]|\n]+)(?:\|([^\]\n]+))?\]\]")
LINK = re.compile(r"\]\(([^)\s]+)\)")
LABELLED = re.compile(r"\[([^\]\n]*)\]\(([^)\s]+)\)")
BENCH = ("raw/", "projects/", "docs/plan/")


def run(*args):
    done = subprocess.run(args, text=True, capture_output=True)
    if done.returncode:
        sys.exit(f"{' '.join(args)}: exit {done.returncode}\n{done.stderr}")
    return done.stdout


def grove(*args):
    return run(GROVE, "--project", ROOT, *args)


def front(text):
    """(fields, body) of a predecessor file; ({}, text) without frontmatter."""
    if not text.startswith("---\n"):
        return {}, text
    end = text.index("\n---\n", 4)
    fields, key = {}, None
    for line in text[4:end].split("\n"):
        if line.lstrip().startswith("- ") and key:  # a block list item of the previous key
            fields[key] = (fields[key] or []) + [line.lstrip()[2:].strip().strip('"')]
            continue
        key, _, value = line.partition(":")
        value = value.strip()
        if value.startswith("["):
            value = [v.strip().strip('"') for v in value[1:-1].split(",") if v.strip()]
        fields[key.strip()] = value
    return fields, text[end + 5:]


def words(stem, keep_case=False):
    title = " ".join(ACRONYMS.get(w, w) for w in stem.split("-"))
    return title if keep_case else title[:1].upper() + title[1:]


def tracked():
    names = run("git", "ls-files", "-z", "--cached", "--others", "--exclude-standard").split("\0")
    return sorted({n for n in names if n and os.path.isfile(n) and not os.path.islink(n)})


def texts():
    for path in tracked():
        try:
            yield path, open(path, encoding="utf-8").read()
        except UnicodeDecodeError:
            pass


def sources():
    """(old key, path, type, title, slug, fields) in conversion order."""
    def listing(*folders):
        return sorted(posixpath.join(f, n) for f in folders if os.path.isdir(f)
                      for n in os.listdir(f) if n.endswith(".md"))

    def by_id(path):
        return int(re.match(r"[WD]-(\d+)", posixpath.basename(path)).group(1))

    out, stems = [], {}
    groups = [
        ("decision", sorted(listing(OLD + "/decisions", OLD + "/history/decisions"), key=by_id)),
        ("work", sorted(listing(OLD + "/work", OLD + "/history/work"), key=by_id)),
        ("question", listing(OLD + "/questions")),
        ("term", listing(OLD + "/terms")),
        ("page", listing(OLD + "/capabilities")),
        ("page", listing(OLD + "/history/evidence")),
    ]
    for kind, paths in groups:
        for path in paths:
            text = open(path).read()
            fields, body = front(text)
            stem = posixpath.basename(path)[:-3]
            ident = fields.get("id") or stem
            m = re.match(r"([WD]-\d+)-(.+)", stem)
            rest = m.group(2) if m else stem
            h1 = re.match(r"# (.+)\n", body)
            if h1:
                title = h1.group(1)
            elif m and kind == "page":
                title = m.group(1) + " " + words(rest, True)
            else:
                title = words(rest, kind == "term")
            out.append((ident, path, kind, title, rest, fields))
            stems[stem] = ident
    for path in listing("docs/plans") + listing("docs/plans/evidence"):
        text = open(path).read()
        stem = posixpath.basename(path)[:-3]
        title = re.match(r"# (.+)\n", text).group(1)
        work = re.findall(r"W-\d{3}", stem)
        if path.startswith("docs/plans/evidence/"):
            kind, slug = "page", re.sub(r"^W-\d{3}-", "", stem)
        else:
            first = next(p for i, p, *_ in out if i == work[0])
            kind, slug = "plan", re.sub(r"^W-\d{3}-", "", posixpath.basename(first)[:-3]) + "-plan"
        out.append((path, path, kind, title, slug, {"work": work}))
    return out, stems


def strip_old_front(path, fields):
    text = open(path).read()
    end = text.index("\n---\n", 4) + 5
    head, body = text[:end], text[end:].lstrip("\n")
    if fields:
        body = body[body.index("\n---\n", 4) + 5:]
    open(path, "w").write(head + "\n" + body)


def convert_all(order):
    mapping = []
    for key, path, kind, title, slug, fields in order:
        entry = json.loads(grove("convert", path, "--type", kind, "--title", title, "--slug", slug))
        entry.update(key=key, type=kind, title=title, fields=fields)
        mapping.append(entry)
        if not path.startswith("docs/plans/"):
            strip_old_front(entry["path"], fields if fields.get("type") else {})
        os.remove(path)
    return mapping


def set_fields(mapping, ids, stems):
    dropped = {}
    for e in mapping:
        f, kind, sets, gone = e["fields"], e["type"], [], []
        if not f.get("type"):  # a legacy plan or a frontmatter-less document
            if f.get("work"):  # plan evidence is a page, which relates rather than belongs
                field = "work" if kind == "plan" else "relates_to"
                sets.append(field + "=" + json.dumps([ids[w] for w in f["work"]]))
        else:
            status = f.get("status")
            related = []
            for key, value in f.items():
                if key in ("id", "type"):
                    continue
                if key == "status":
                    if kind == "page":
                        gone.append(f"status: {value}")
                    elif status == "parked":
                        note = f"Parked in the predecessor (status `parked`); `open` since the Grove cutover on {TODAY}.\n\n"
                        text = open(e["path"]).read()
                        end = text.index("\n---\n", 4) + 6
                        open(e["path"], "w").write(text[:end] + note + text[end:].lstrip("\n"))
                        gone.append("status: parked")
                    elif status == "done":
                        text = open(e["path"]).read()
                        open(e["path"], "w").write(text.replace("\nstatus: proposed\n", "\nstatus: done\n", 1))
                    elif status not in ("proposed", "open", "current"):
                        sets.append("status=" + status)
                elif key == "kind" and kind == "work":
                    sets.append("kind=" + ("investigation" if value == "spike" else value))
                elif key == "size" and value in ("small", "medium", "large"):
                    sets.append("size=" + value)
                elif key == "priority" and kind == "work":
                    sets.append("priority=" + value)
                elif key in ("depends_on", "members", "blocks"):
                    if value:
                        sets.append(key + "=" + json.dumps([ids[v] for v in value]))
                elif key in ("scope", "applies_to", "supersedes", "superseded_by", "sources"):
                    values = value if isinstance(value, list) else [value]
                    named = [stems.get(v, v) for v in values]  # sources name records by filename stem
                    if key == "sources" and not all(n in ids for n in named):
                        gone.append(f"{key}: {json.dumps(value)}")
                        continue
                    related += [ids[n] for n in named if n and ids[n] not in related and ids[n] != e["id"]]
                else:
                    gone.append(f"{key}: {json.dumps(value) if isinstance(value, list) else value}")
            if related:
                sets.append("relates_to=" + json.dumps(related))
        if sets:
            revision = json.loads(grove("show", e["id"], "--json"))["revision"]
            grove("update", e["id"], "--expect", revision, *[a for s in sets for a in ("--set", s)])
        if gone:
            dropped[e["id"]] = gone
    return dropped


def move_brief(ids, owner):
    fields, body = front(open(OLD + "/brief.md").read())
    parts = re.split(r"(?m)^(?=## )", body)
    nxt = next(p for p in parts if p.startswith("## Next\n"))
    open("grove/brief.md", "w").write("".join(p for p in parts if p is not nxt))
    os.remove(OLD + "/brief.md")
    text = open(owner).read().rstrip("\n")
    moved = nxt[len("## Next\n"):].strip("\n")
    open(owner, "w").write(text + f"\n\nMoved from the brief's Next at the Grove cutover on {TODAY}, with its IDs and links rewritten like the rest:\n\n{moved}\n")
    return [f"{k}: {v}" for k, v in fields.items() if k not in ("type", "id")]


def move_evidence(ids):
    moved = {}
    for folder in ("docs/plans", "docs/plans/evidence"):
        for name in sorted(os.listdir(folder)):
            path = posixpath.join(folder, name)
            if os.path.isfile(path) and not name.endswith(".md"):
                new = "docs/evidence/" + OLD_ID.sub(lambda m: ids[m.group(1)], re.sub(r"^(W-\d{3})-", lambda m: ids[m.group(1)] + "-", name))
                os.makedirs("docs/evidence", exist_ok=True)
                os.rename(path, new)
                moved[path] = new
    return moved


def remove_leftovers():
    removed = []
    for folder, _, files in sorted(os.walk(OLD)):
        for name in sorted(files):
            path = posixpath.join(folder, name)
            removed.append(path)
            os.remove(path)
    for folder in (OLD, "docs/plans"):
        for sub, dirs, _ in sorted(os.walk(folder), reverse=True):
            os.rmdir(sub)  # fails loudly if anything unaccounted for remains
    os.remove("grove.toml")
    return removed


def repair(paths, stems, ids, report, changed):
    back = {}
    for old, new in paths.items():  # the first entry is the file's own path; later ones are aliases
        back.setdefault(new, old)
    by_length = sorted(paths, key=len, reverse=True)
    exempt_paths = {"grove.yaml"}

    def link(match, old_file, new_file):
        label, dest = match.groups()
        if re.match(r"[a-zA-Z][a-zA-Z0-9+.-]*:|#|/", dest):
            return match.group(0)
        target, _, fragment = dest.partition("#")
        resolved = posixpath.normpath(posixpath.join(posixpath.dirname(old_file), target))
        if resolved not in paths and old_file == new_file:
            return match.group(0)
        relative = posixpath.relpath(paths.get(resolved, resolved), posixpath.dirname(new_file) or ".")
        if label in (target, posixpath.basename(target), "`" + posixpath.basename(target) + "`"):
            label = label.replace(posixpath.basename(target), posixpath.basename(paths.get(resolved, resolved)))
        return "[" + label + "](" + relative + ("#" + fragment if fragment else "") + ")"

    def wiki(match, new_file):
        target, label = match.group(1).strip(), match.group(2)
        ident = stems.get(target)
        if ident is None:
            return match.group(0)
        path = new_path_of[ident]
        text = label or (ids[ident] if re.match(r"[WD]-\d", target) else target)
        return f"[{text}]({posixpath.relpath(path, posixpath.dirname(new_file) or '.')})"

    new_path_of = {k: paths[p] for k, p in old_path_of.items()}
    for new_file, text in texts():
        if new_file in exempt_paths:
            continue
        old_file = back.get(new_file, new_file)
        keep = [(needle, kept) for path, needle, kept in EXCEPTIONS if path == old_file]
        lines = text.split("\n")
        for n, line in enumerate(lines):
            if line.startswith("formerly:"):
                continue
            before = line
            if new_file.endswith(".md"):
                line = LABELLED.sub(lambda m: link(m, old_file, new_file), line)
                line = WIKILINK.sub(lambda m: wiki(m, new_file), line)  # already relative to new_file
            for old in by_length:
                line = line.replace(old, paths[old])
            # "W-021/023/024" and "D-0034/0035" abbreviate the later IDs
            line = re.sub(r"\b([WD]-)(\d{3,4})((?:/\d{3,4})+)\b",
                          lambda m: m.group(1) + m.group(2) + "".join("/" + m.group(1) + n for n in m.group(3)[1:].split("/")), line)
            kept = {i for needle, some in keep if needle in line for i in some}
            if kept & set(OLD_ID.findall(line)):
                report.append(f"{new_file}:{n + 1}: {', '.join(sorted(kept))}")
            line = OLD_ID.sub(lambda m: m.group(1) if m.group(1) in kept else ids.get(m.group(1), m.group(1)), line)
            if line != before:
                changed.append(f"{new_file}:{n + 1}: {line.strip()}")
            lines[n] = line
        new_text = "\n".join(lines)
        if new_text != text:
            open(new_file, "w").write(new_text)


GROVE_SECTION = """## Grove
Project knowledge lives in `grove/` as Grove records: `grove.yaml` configures it and `grove/brief.md` is the brief. Use the `grove` CLI on `PATH`; `grove version` must print `grove v…`. Read with `grove brief`, `grove list`, `grove show G-NNN` and `grove context G-NNN`, or `grove` alone for the board. Create records with `grove new` and change fields with `grove update`; never hand-number an ID.
- Shape ideas into proposed work with the `grove-shape` entrypoint and carry assigned work IDs with `grove-work` (`/grove-shape`, `/grove-work` in Claude Code; `$grove-shape`, `$grove-work` in Codex). `grove guide shape` and `grove guide work` print the workflows they follow.
- Implement in a linked worktree under `.claude/worktrees/` on branch `worktree-G-NNN` (`worktree-G-NNN-G-MMM` for several IDs), never on main. A headless proposal branch is `worktree-shape-SLUG`.
- [MAP_ID](MAP_PATH) maps the predecessor's IDs and paths to these records.
"""


def hand_edits(map_id, map_path):
    section = GROVE_SECTION.replace("MAP_ID", map_id).replace("MAP_PATH", map_path)
    for name in ("CLAUDE.md", "AGENTS.md"):
        text = open(name).read()
        new = re.sub(r"<!-- grove:begin -->\n.*?<!-- grove:end -->\n", section, text, flags=re.S)
        assert new != text, name
        open(name, "w").write(new)
    settings = json.load(open(".claude/settings.json"))
    assert settings == {"enabledPlugins": {"grove@grove-local": True}}, settings
    os.remove(".claude/settings.json")


def new_page():
    grove("new", "page", "Predecessor identity and path migration map", "--slug", "migration-map")
    return next(p for p in tracked() if re.match(r"grove/G-\d+-migration-map\.md$", p))


def mapping_page(page, mapping, dates, evidence, removed, dropped, brief_dropped, leftovers, bench):
    head = open(page).read().split("\n---\n")[0] + "\n---\n"
    rows = "\n".join(
        f"| {e['fields'].get('id', '-')} | `{e['from_path']}` | {e['id']} | `{e['path']}` | {dates.get(e['from_path'], '-')} |"
        for e in mapping)
    moved = "\n".join(f"| - | `{old}` | - | `{new}` | - |" for old, new in evidence.items())
    gone = "\n".join(f"- {i}: " + "; ".join(f"`{d}`" for d in ds) for i, ds in dropped.items())
    excepted = "\n".join("- `" + r.replace("`", "") + "`" for r in leftovers) or "- None."
    benches = "\n".join("- `" + b + "`" for b in bench) or "- None."
    removed_list = "\n".join(f"- `{p}`" for p in removed + ["grove.toml", ".claude/settings.json"])
    open(page, "w").write(head + f"""
The Grove cutover (Grove's G-041) converted every record of the predecessor
tool, every legacy plan and the Markdown plan evidence of this repository into
Grove records with neutral IDs directly under `grove/` on {TODAY}. This page is
the one durable mapping from an old ID or path to its counterpart. To find what
an old ID became, look it up here or search `formerly:`; to read an old commit,
use that commit's files. These `G-` numbers belong to this repository and are
independent of Grove's own repository, whose records use the same prefix.

Numbers follow old ID order within groups: decisions (`D-00NN` became `G-0NN`),
work (`W-0NN` became the number 43 higher), questions, terms, capabilities,
evidence, legacy plans and plan evidence. The order conveys no authority or
priority. "Old date" is the predecessor's day-only `created`, else its
`updated`; neither was kept, since the whole tree entered this repository in
one commit on 2026-09-14 and the record model invents no dates.

| Old ID | Old path | ID | Path | Old date |
| --- | --- | --- | --- | --- |
{rows}
| - | `{OLD}/brief.md` | - | `grove/brief.md` | the brief, not a record |
{moved}

## Conversion rules

- Capabilities, research and other evidence became pages; legacy plans became
  plans with `work` set from the IDs in their filenames.
- `scope`, `applies_to`, `supersedes` and `superseded_by` became `relates_to`.
  The two replaced decisions keep the status `superseded`.
- Work kind `spike` became `investigation`. Sizes `bounded` and `spike` have
  no equivalent and were dropped. `parked` questions became `open` with a
  first line saying so. Done work is historical Done: it has no `candidate`.
- The brief's Next moved to the end of the owning work's Next; its
  frontmatter was dropped (below).
- Wikilinks became Markdown links; bare old IDs and old paths were rewritten
  in every tracked text file, code and tests included.

## Dropped values

Values with no equivalent, per new ID; the original is at `formerly` in the
commit before the cutover.

- brief: {"; ".join(f"`{d}`" for d in brief_dropped)}
{gone}

## Removed

The predecessor's delivery receipts, placeholder files and configuration,
and the plugin setting that enabled it:

{removed_list}

`history/README.md` recorded where the records came from, which stays true:

> Prior records: Bench vault at `~/GitHub/mascah/bench/projects/nullsec` and
> `~/GitHub/mascah/bench/docs/plan/nullsec`. Last Bench-managed commit of this
> repo: `d33dcf429232d8d5b8a22a3fb15575d440b26c84`. Bench run records for this
> repo were under `docs/plan/` at that commit.

## Kept as written

Old identifiers remain in `formerly` fields, in this page, in Git branch names
such as `worktree-W-038` and in commit messages, which name what existed then.
Art capture and intermediate names such as `art/captures/w039-crew-study/`
carry the old work number in lowercase (`w039` is W-039) and were not renamed.
Wikilinks into the Bench vault stay as written:

{benches}

Every other old ID and path was rewritten. These lines keep one on purpose:
a branch name, a file name of the predecessor's run artifacts that were never
tracked here, a folder named as it was when the text was written, or a
verbatim quotation (shown without backticks):

{excepted}
""")
    return page


def audit():
    broken, leftovers = [], []
    for path, text in texts():
        for n, line in enumerate(text.split("\n"), 1):
            if path.endswith(".md"):
                for dest in LINK.findall(line):
                    target = dest.partition("#")[0]
                    there = posixpath.normpath(posixpath.join(posixpath.dirname(path), target))
                    if target and not re.match(r"[a-zA-Z][a-zA-Z0-9+.-]*:|/", dest) and not os.path.exists(there):
                        broken.append(f"{path}:{n}: {dest}")
            if path.endswith("-migration-map.md") or line.startswith("formerly:"):
                continue
            named = re.sub(r"worktree-[\w-]+|worktrees/[\w-]+|shape/[\w-]+", "", line)  # branch names
            if OLD_ANY.search(named) or re.search(r"docs/grove\b|docs/plans\b|grove\.toml", named):
                leftovers.append(f"{path}:{n}: {line.strip()}")
    return broken, leftovers


def finish():
    broken, leftovers = audit()
    print(f"{len(broken)} broken links; {len(leftovers)} lines with old IDs or paths")
    for line in broken:
        print("BROKEN", line)
    for line in leftovers:
        print("OLD", line)


if sys.argv[4:] == ["--audit"]:
    finish()
    sys.exit()
if sys.argv[4:] == ["--baseline"]:  # before any change: links already broken
    print("\n".join(audit()[0]))
    sys.exit()

print(grove("init"), end="")
order, stems = sources()
dates = {path: f.get("created") or f.get("updated") or "-" for _, path, _, _, _, f in order}
old_path_of = {key: path for key, path, *_ in order}
mapping = convert_all(order)
ids = {e["key"]: e["id"] for e in mapping}
dropped = set_fields(mapping, ids, stems)
owner = next(e["path"] for e in mapping if e["key"] == "W-032")
brief_dropped = move_brief(ids, owner)
evidence = move_evidence(ids)
removed = remove_leftovers()
paths = {e["from_path"]: e["path"] for e in mapping}
for e in mapping:  # closed records were linked by the path they had before history/
    if "/history/" in e["from_path"]:
        paths[e["from_path"].replace("/history/", "/")] = e["path"]
for old, new in list(paths.items()):  # records named paths relative to the old root
    if old.startswith(OLD + "/history/"):
        paths[old[len(OLD) + 1:]] = new
paths[OLD + "/brief.md"] = "grove/brief.md"
paths.update(evidence)
report, changed = [], []
bench = []
for path, text in texts():
    for n, line in enumerate(text.split("\n"), 1):
        for m in WIKILINK.finditer(line):
            if path.endswith(".md") and m.group(1).startswith(BENCH):
                bench.append(f"{path}:{n}: {m.group(0)}")
repair(paths, stems, ids, report, changed)
page = new_page()
hand_edits(posixpath.basename(page).split("-migration")[0], posixpath.relpath(page, "."))
leftovers = audit()[1]  # the page is still empty, so it is not among them
mapping_page(page, mapping, dates, evidence, removed, dropped, brief_dropped, leftovers, bench)
print(grove("check"), end="")
json.dump([{k: e[k] for k in ("key", "from_path", "id", "path", "type", "title")} for e in mapping],
          open(os.path.join(OUT, "mapping.json"), "w"), indent=1)
open(os.path.join(OUT, "changed.txt"), "w").write("\n".join(changed) + "\n")
unused = [x for x in EXCEPTIONS if not any(r.startswith(paths.get(x[0], x[0]) + ":") for r in report)]
assert not unused, unused
print(f"{len(mapping)} converted; map {page}; {len(changed)} lines rewritten")
finish()
