# (f) A draft nobody attends: the request settles the exit status and
# leaves open where the check for an unknown ID goes.

setup() { :; }

# The script makes the branch: Codex's sandbox cannot write .git.
runs=(
  "unknown-id => /grove-propose --unattended Editing an unknown ID prints an error and exits 2"
)
answers=()
# No one is in the session, so nothing answers it.
turns=0

# good checks that the agent's one turn did not end waiting, that unknown-id holds one
# change grove check passes and main none, that its Decisions quote the
# request, and that each open question recommends an answer and names a task.
good() {
  ! waiting || echo "the agent's last message waits for an answer"
  [ -z "$(git ls-tree --name-only main grove/changes/)" ] || echo "main holds a change"
  local new
  new=$(git ls-tree --name-only unknown-id grove/changes/ | grep -v '/archive$' || true)
  [ "$(wc -l <<< "$new" | tr -d ' ')" -eq 1 ] && [ -n "$new" ] || { echo "unknown-id holds not one change but: ${new:-none}"; return; }
  git switch -q unknown-id
  grove check > /dev/null 2>&1 || echo "grove check reports problems: $(grove check 2>&1 | head -3)"
  python3 - "$new/proposal.md" << 'PY'
import re, sys
s = open(sys.argv[1]).read()
def section(name):
    m = re.search(r"^## " + name + r"\s*\n(.*?)(?=^## |\Z)", s, re.M | re.S)
    return m.group(1) if m else None
d, q = section("Decisions"), section("Open questions")
if d is None or not re.search(r"exits? 2", d):
    print("Decisions do not quote the request's exit status")
lead, *items = re.split(r"^(?:\d+\.|-) ", q or "", flags=re.M)
if not items:
    print("the proposal has no open questions")
# A task the section's lead names for every question counts for each.
for i, item in enumerate(items, 1):
    if not re.search(r"recommend", item, re.I) or not re.search(r"T-\d\d", lead + item):
        print(f"open question {i} gives no recommended answer or names no task")
if not re.search(r"records\.py|`find`|ValueError|cli\.py", q or ""):
    print("no open question says where the check for an unknown ID goes")
PY
}
