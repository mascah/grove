#!/usr/bin/env python3
"""tasks: one person's to-do list as Markdown files under tasks/."""
import argparse, datetime, json, os, re, sys

DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "tasks")
STATUSES = ("open", "done", "dropped")


def load():
    tasks = []
    for name in sorted(os.listdir(DIR)):
        if not name.endswith(".md"):
            continue
        text = open(os.path.join(DIR, name)).read()
        front, _, body = text.partition("\n---\n")
        task = {"id": name[:3], "file": name, "title": body.strip().splitlines()[0].lstrip("# ")}
        for line in front.splitlines()[1:]:
            key, _, value = line.partition(": ")
            task[key] = json.loads(value) if value.startswith("[") else value
        tasks.append(task)
    return tasks


def write(task, title):
    lines = ["---"] + [f"{k}: {json.dumps(task[k]) if k == 'tags' else task[k]}" for k in ("status", "created", "closed", "tags")]
    with open(os.path.join(DIR, task["file"]), "w") as f:
        f.write("\n".join(lines) + f"\n---\n# {title}\n")


def add(args):
    tasks = load()
    number = max([int(t["id"]) for t in tasks] + [0]) + 1
    slug = re.sub(r"[^a-z0-9]+", "-", args.title.lower()).strip("-")
    task = {"file": f"{number:03d}-{slug}.md", "status": "open", "created": datetime.date.today().isoformat(), "closed": "", "tags": sorted({t.lower() for t in args.tag})}
    write(task, args.title)
    print(f"{number:03d}")


def close(args, status):
    for task in load():
        if task["id"] == args.id:
            task["status"], task["closed"] = status, datetime.date.today().isoformat()
            write(task, task["title"])
            return
    sys.exit(f"no task {args.id}")


def list_(args):
    for task in load():
        if args.status and task["status"] not in args.status:
            continue
        tags = " ".join("#" + t for t in task["tags"])
        print(f"{task['id']}  {task['status']:<7}  {task['title']}  {tags}".rstrip())


def export(args):
    json.dump(load(), sys.stdout, indent=2)
    print()


def main():
    parser = argparse.ArgumentParser(prog="tasks")
    sub = parser.add_subparsers(dest="command", required=True)
    p = sub.add_parser("add"); p.add_argument("title"); p.add_argument("--tag", action="append", default=[]); p.set_defaults(run=add)
    p = sub.add_parser("list"); p.add_argument("--status", action="append", choices=STATUSES); p.set_defaults(run=list_)
    p = sub.add_parser("done"); p.add_argument("id"); p.set_defaults(run=lambda a: close(a, "done"))
    p = sub.add_parser("drop"); p.add_argument("id"); p.set_defaults(run=lambda a: close(a, "dropped"))
    sub.add_parser("export").set_defaults(run=export)
    args = parser.parse_args()
    args.run(args)


if __name__ == "__main__":
    main()
