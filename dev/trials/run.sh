#!/usr/bin/env bash
# Runs the skill trials of dev/trials.md (CONSTITUTION.md rule 17): each
# case, or those named, in Claude Code and in Codex, on a fresh fixture,
# with this checkout's build of grove. Prints the results table.
#
#   dev/trials/run.sh [a b c d e]
#
# HARNESSES (default "claude codex") picks the harnesses, RUNS (default 3)
# the runs of each case and harness, JUDGINGS (default 3) the judgings of
# each run, CLAUDE_MODEL and CODEX_MODEL the models under test, TRIALS_OUT the
# output folder, TURNS (default 8) the owner replies a run may get, and
# SETUP_ONLY=1 stops each trial once its fixture is set up.
set -euo pipefail

here=$(cd "$(dirname "$0")" && pwd)
checkout=$(cd "$here/../.." && pwd)
out=${TRIALS_OUT:-${TMPDIR:-/tmp}/grove-trials/$(date +%Y%m%d-%H%M%S)}
turns=${TURNS:-8}
runs_each=${RUNS:-3}
judgings=${JUDGINGS:-3}
harnesses=${HARNESSES:-claude codex}
cases=("$@")
[ ${#cases[@]} -gt 0 ] || cases=(a b c d e)
fallback="Yes, go with what you recommend."
# Changing the agent under test never changes the grader: the judge and the
# waiting check are fixed here. Every model runs at medium effort.
claude_model=${CLAUDE_MODEL:-claude-opus-5-5}
codex_model=${CODEX_MODEL:-gpt-6.1-sol}
judge_model=claude-opus-5-5
waiting_model=claude-sonnet-5-5
codex_auth=${CODEX_HOME:-$HOME/.codex}/auth.json
claude_tools="Read Glob Grep Edit Write Skill TodoWrite Bash(grove:*) Bash(git:*) Bash(python3:*) Bash(ls:*) Bash(cat:*) Bash(grep:*) Bash(rg:*) Bash(head:*) Bash(tail:*) Bash(sed:*) Bash(find:*) Bash(wc:*) Bash(diff:*) Bash(mkdir:*) Bash(mv:*)"

(cd "$checkout" && go build -o bin/grove ./cmd/grove)
export PATH="$checkout/bin:$PATH"
commit_of=$(git -C "$checkout" rev-parse --short HEAD)$(git -C "$checkout" diff --quiet HEAD || echo +dirty)

# Helpers for the case files, run in the fixture.

# put writes stdin to path, with @BASE@ as the current commit.
put() { mkdir -p "$(dirname "$1")"; sed "s/@BASE@/$(git rev-parse HEAD)/g" > "$1"; }
commit() { git add -A && git commit -qm "$1"; }
# spec writes a maintained spec from stdin, with the .grove.yaml a Grove
# before 261023-d6r79 needs to find it, unless this Grove no longer reads one.
spec() {
  put "grove/specs/$1/spec.md"
  printf 'title: %s\n' "$2" > "grove/specs/$1/.grove.yaml"
  local out
  out=$(grove check 2> /dev/null || true)
  if grep -q "specs/$1/.grove.yaml: no longer read" <<< "$out"; then rm "grove/specs/$1/.grove.yaml"; fi
}
# edit replaces old with new, both literal, once in file.
edit() { python3 -c 'import sys; p,o,n=sys.argv[1:]; s=open(p).read(); assert o in s, (p,o); open(p,"w").write(s.replace(o,n,1))' "$@"; }

# snapshot lists every file of the fixture with its hash. A file moved or
# deleted is left out, so it shows as changed.
snapshot() { git ls-files -co --exclude-standard -z | xargs -0 shasum 2> /dev/null | sort -k2 || true; }

# Helpers for a case's good, which reads what each turn did: turns.txt has a
# line "<turn>\t<branch>\t<owner's message>" per turn, snapshot-<turn>.txt
# the fixture after it, and start-<turn>.txt the fixture before a run's first.

# answered prints the turn whose owner's message holds $1, such as a scripted answer.
answered() { awk -F'\t' -v a="$1" 'index($3, a) {print $1; exit}' "$dir/turns.txt"; }
# kept reports whether path $1 stayed as it was at the start of branch $2's
# run after each of its turns before turn $3, or after all of them.
kept() {
  local t first="" h
  for t in $(awk -F'\t' -v b="$2" '$2 == b {print $1}' "$dir/turns.txt"); do
    [ -n "$first" ] || { first=$t h=$(awk -v p="$1" '$2 == p {print $1}' "$dir/start-$t.txt"); }
    [ -z "${3:-}" ] || [ "$t" -lt "$3" ] || break
    [ "$(awk -v p="$1" '$2 == p {print $1}' "$dir/snapshot-$t.txt")" = "$h" ] || return 1
  done
}
# task prints task $1's section of the implementation.md on stdin.
task() { awk -v t="### $1 " 'index($0, t) == 1 {on = 1; next} on && /^##/ {exit} on'; }

# save commits what the last run left, if anything.
save() { git diff --quiet HEAD && [ -z "$(git ls-files -o --exclude-standard)" ] || commit "After the last run"; }

# result, a jq function, prints a tool's result for the transcript: an error
# mark, then its first 20 lines, or all of a grove command's, since the judge
# checks what Grove's instructions listed.
result='def result($out; $err; $cmd):
  ($out | rtrimstr("\n") | split("\n")) as $l
  | (if $cmd | test("(^|[^A-Za-z0-9_/.-])grove( |$)") then $l else $l[:20] end) as $s
  | "RESULT" + (if $err then " (error)" else "" end) + ":", ($s[] | "  " + .),
    (if ($l | length) > ($s | length) then "  … \(($l | length) - ($s | length)) more lines" else empty end);'

# turn sends msg to the agent, resuming session sid when set, and sets sid,
# last (the agent's last message), said (all its messages this turn) and
# words (the words of its messages).
turn() {
  local msg=$1 log=$dir/turn-$((n += 1)).jsonl
  echo "OWNER: $msg" >> "$dir/transcript.txt"
  if [ "$harness" = claude ]; then
    claude -p "$msg" ${sid:+--resume "$sid"} --model "$claude_model" --effort medium --setting-sources project,local --permission-mode acceptEdits \
      --allowedTools "$claude_tools" --output-format stream-json --verbose < /dev/null > "$log" 2>> "$dir/stderr.txt" || true
    sid=$(jq -r 'select(.type=="result") | .session_id' "$log" | tail -1)
    last=$(jq -rs '[.[] | select(.type=="result") | .result] | last // ""' "$log")
    jq -rs "$result"'
      (map(.message.content?[]? | objects | select(.type=="tool_use") | {key: .id, value: (.input.command // "" | tostring)}) | from_entries) as $cmd
      | .[] | select(.parent_tool_use_id==null)
      | if .type=="assistant" then .message.content[]? |
          if .type=="text" then "AGENT: " + .text
          elif .type=="tool_use" then "TOOL: " + .name + " " + ((.input.command // .input.file_path // .input.skill // .input.pattern // "") | tostring)
          else empty end
        elif .type=="user" then .message.content[]? | objects | select(.type=="tool_result")
          | result(if (.content | type)=="string" then .content else [.content[]? | .text? // empty] | join("\n") end; .is_error==true; $cmd[.tool_use_id] // "")
        else empty end' "$log" >> "$dir/transcript.txt"
    said=$(jq -r 'select(.type=="assistant" and .parent_tool_use_id==null) | .message.content[]? | select(.type=="text") | .text' "$log")
  else
    # HOME and CODEX_HOME are the trial's, so none of the user's Codex config,
    # plugins, hooks, memories or skills load.
    local cfg=(--json -m "$codex_model" -c model_reasoning_effort='"medium"' -c allow_login_shell=false)
    if [ -z "$sid" ]; then
      HOME=$dir CODEX_HOME=$dir/codex codex exec "${cfg[@]}" --sandbox workspace-write "$msg" < /dev/null > "$log" 2>> "$dir/stderr.txt" || true
      sid=$(jq -r 'select(.type=="thread.started") | .thread_id' "$log")
    else
      HOME=$dir CODEX_HOME=$dir/codex codex exec resume "${cfg[@]}" -c sandbox_mode='"workspace-write"' "$sid" "$msg" < /dev/null > "$log" 2>> "$dir/stderr.txt" || true
    fi
    last=$(jq -rs '[.[] | select(.type=="item.completed" and .item.type=="agent_message") | .item.text] | last // ""' "$log")
    jq -r "$result"'select(.type=="item.completed") | .item |
      if .type=="agent_message" then "AGENT: " + .text
      elif .type=="command_execution" then "TOOL: " + .command, result(.aggregated_output // ""; .exit_code != 0; .command)
      elif .type=="file_change" then "TOOL: edit " + ([.changes[]?.path] | join(", ")), result(.status // ""; .status=="failed"; "")
      else empty end' "$log" >> "$dir/transcript.txt"
    said=$(jq -r 'select(.type=="item.completed" and .item.type=="agent_message") | .item.text' "$log")
  fi
  words=$((words + $(wc -w <<< "$said")))
  printf '%s\t%s\t%s\n' "$n" "$(git branch --show-current)" "$msg" >> "$dir/turns.txt"
  snapshot > "$dir/snapshot-$n.txt"
  # A turn the harness failed, such as at a usage limit, fails the trial rather than scoring it.
  failed=$(jq -rs '[.[] | select(.type=="turn.failed" or (.type=="result" and .is_error==true)) | .error.message // .result // "error"] | first // ""' "$log")
  [ -z "$failed" ] || { echo "$failed" > "$dir/failed.txt"; return 1; }
}

# reply sets answer to the case's first unused answer whose pattern matches
# any of the agent's messages this turn, or the fallback: Codex can split a
# turn into several messages, the question in an earlier one.
reply() {
  local i
  answer=$fallback
  for i in "${!answers[@]}"; do
    if [ -z "${used[$i]:-}" ] && grep -Eiq -- "${answers[$i]%% => *}" <<< "$said"; then
      used[$i]=1 answer=${answers[$i]#* => }
      return
    fi
  done
}

# waiting asks a model whether the agent's last message waits for the
# owner's answer, as against reporting what it did and leaving the next step
# to the owner, which a scripted reply would otherwise take for them.
waiting() {
  local verdict
  verdict=$(claude -p "Here is the last message an AI coding agent sent the person it works for. Does it stop to wait for that person to answer or decide something before it can go on with the task it was given? A step the agent says it will take itself once the person answers or gives feedback counts. Next steps it leaves for the person to take, such as committing, merging, reviewing or running another command or skill, do not count, and neither do offers of optional further work, such as a commit. Reply with only yes or no.

$last" --model "$waiting_model" --effort medium --setting-sources project,local --tools "" --output-format json < /dev/null 2>> "$dir/stderr.txt" | jq -r '.result')
  echo "turn $n: waiting: $verdict" >> "$dir/waits.txt"
  grep -qi '^yes' <<< "$verdict"
}

# outside lists, as "branch: path", each file outside grove/ that the current
# branch changes against main and that grove review marks neither planned nor
# agreed: named nowhere, or not listed, as on a branch holding no change.
outside() {
  local b marked
  b=$(git branch --show-current)
  marked=$(grove review | sed -nE '/^diff --git/q; s/^  (planned|agreed) .*  ([^ ]+)$/\2/p' | sort -u)
  git diff --name-only --no-renames main...HEAD | { grep -v '^grove/' || true; } | sort -u |
    comm -23 - <(echo "$marked") | sed "s|^|$b: |"
}

# trial runs case c in harness h, as run r of the k, and leaves its counts
# in its folder.
trial() {
  local c=$1 h=$2 r=$3
  dir=$out/$c-$h-$r harness=$h n=0 words=0
  local repo=$dir/repo
  mkdir -p "$dir"
  [ "$h" = claude ] || { mkdir -p "$dir/codex" && ln -s "$codex_auth" "$dir/codex/auth.json"; }
  (cd "$checkout" && just fixture "$repo" > /dev/null)
  cd "$repo"
  git config user.name "Fixture Owner" && git config user.email fixture@example.com
  grove init > /dev/null
  edit grove.yaml 'context: ""' 'context: "Run python3 -m unittest before saying a task is done."'
  commit "Set up Grove"
  local id="" check_ok="" runs=() answers=() used=()
  # shellcheck source=/dev/null
  source "$here/$c.sh"
  setup
  # The judge reads the fixture as the setup left it, uncommitted work
  # included, before any run wrote to it.
  mkdir "$dir/setup" && git ls-files -co --exclude-standard -z | tar -c --null -T - -f - | tar -x -C "$dir/setup"
  # A case sets check_ok to the problems its setup means grove check to report.
  if ! grove check > "$dir/setup-check.txt" 2> /dev/null && grep -vqE "${check_ok:-^$}" "$dir/setup-check.txt"; then
    echo "the setup is not clean: $(cat "$dir/setup-check.txt")" > "$dir/failed.txt"
    return 1
  fi
  [ -z "${SETUP_ONLY:-}" ] || { touch "$dir/setup-only"; return 0; }
  local run branch prompt had new branches=()
  for run in "${runs[@]}"; do
    branch=${run%% => *} prompt=${run#* => }
    [ "$(git branch --show-current)" = "$branch" ] || { save; git switch -q "$branch" 2> /dev/null || git switch -q -c "$branch" main; }
    [[ " ${branches[*]} " == *" $branch "* ]] || branches+=("$branch")
    had=$(ls grove/changes 2> /dev/null || true)
    if [[ $prompt == *"{id}"* && -z $id ]]; then
      echo "no change for \"$prompt\"; the run before it proposed none" > "$dir/failed.txt"
      return 1
    fi
    prompt=${prompt//\{id\}/$id}
    [ "$h" = codex ] && prompt=${prompt/\/grove-/\$grove-}
    sid="" last=""
    snapshot > "$dir/start-$((n + 1)).txt"
    turn "$prompt"
    local k=0
    # An unattended run has no owner to answer it.
    [[ $prompt == *--unattended* ]] || while [ $((k += 1)) -le "$turns" ] && waiting; do reply && turn "$answer"; done
    new=$(comm -13 <(echo "$had") <(ls grove/changes 2> /dev/null) | grep -v '^archive$' | head -1 || true)
    [ -z "$new" ] || id=${new:0:12}
  done
  save
  local changes
  changes=$(find grove/changes -type f 2> /dev/null | sort)
  # A case's good prints one line per check of its "a good run" that failed.
  good > "$dir/good.txt" 2>&1 || true
  # Each judging reads the saved transcript; no agent session starts after the last turn.
  local i
  for ((i = 1; i <= judgings; i++)); do (cd "$dir/setup" && judge "$c" "$changes") > "$dir/judge-$i.json" || true; done
  for branch in "${branches[@]}"; do git switch -q "$branch" && outside; done > "$dir/outside.txt"
  echo "$words" > "$dir/words.txt"
}

# judge asks a model for the four judged counts of the run in $dir, of case
# $1, whose change documents at the run's end are $2.
judge() {
  local prompt
  prompt="You judge one run of Grove's skill trials. Below are how each count is defined, the case, and the run's transcript: OWNER lines are what the tester typed, AGENT lines the agent's messages, TOOL lines its tool calls, and RESULT lines what each returned, its first 20 lines, or all of a grove command's. Your working folder holds the case's fixture as its setup left it, before the run: none of what the run wrote is there. Read it with Read, Grep and Glob to check what the repository already answered when the agent asked.

$(awk '/^## The counts/{on=1;next} /^## /{on=0} on' "$checkout/dev/trials.md")

$(awk -v h="### ($1)" 'index($0,h)==1{on=1} on&&/^##/&&index($0,h)!=1{exit} on' "$checkout/dev/trials.md")

The files under grove/changes/ at the end of the run, which the instructions of a change list as its documents:

$2

Transcript:

$(cat "$dir/transcript.txt")

Reply with only this JSON, each why one sentence naming what you counted:
{\"unneeded_questions\": {\"count\": 0, \"why\": \"\"}, \"missed_decisions\": {\"count\": 0, \"why\": \"\"}, \"done_claims_without_evidence\": {\"count\": 0, \"why\": \"\"}, \"listed_documents_not_read\": {\"count\": 0, \"why\": \"\"}}"
  claude -p "$prompt" --model "$judge_model" --effort medium --setting-sources project,local \
    --tools Read,Grep,Glob --allowedTools "Read Grep Glob" --output-format json < /dev/null 2>> "$dir/stderr.txt" |
    jq -r '.result' | python3 -c 'import sys, json; t = sys.stdin.read(); print(json.dumps(json.loads(t[t.index("{"):t.rindex("}") + 1]), indent=2))'
}

mkdir -p "$out"
for c in "${cases[@]}"; do
  for h in $harnesses; do
    # A case and harness's runs go one after another. Each runs in the
    # background and is waited for, not run under ||, so set -e holds in it.
    (for ((r = 1; r <= runs_each; r++)); do (trial "$c" "$h" "$r") & wait $! || true; done) &
  done
done
wait
version() { "$1" --version 2> /dev/null | head -1 || echo "no $1"; }
echo "Trials at $commit_of ($out): Claude Code $claude_model and Codex $codex_model at medium effort, judge $judge_model at medium effort, waiting check $waiting_model at medium effort; $(version claude), $(version codex)"
echo
# One row per case and harness: how many runs passed their case's good, and
# each count the median of its runs, and at more than one run their min–max;
# a run's judged count is the median of its judgings. Below the table, each
# failed run, each check a run failed, and each count its judgings disagree
# on. It exits 0 whatever it shows: a record, not a gate.
python3 - "$out" "$runs_each" "$judgings" "$harnesses" "${cases[@]}" << 'EOF'
import json, os, statistics, sys

out, k, j, harnesses, cases = sys.argv[1], int(sys.argv[2]), int(sys.argv[3]), sys.argv[4].split(), sys.argv[5:]
judged = {"unneeded_questions": "Unneeded questions", "missed_decisions": "Missed decisions",
          "done_claims_without_evidence": "Done claims without evidence", "listed_documents_not_read": "Listed documents not read"}
columns = ["unneeded_questions", "missed_decisions", "outside", "done_claims_without_evidence", "listed_documents_not_read", "words"]


def judging(path, key):
    try:
        return int(json.load(open(path))[key]["count"])
    except (OSError, ValueError, KeyError, TypeError):
        return None


def spread(values):
    if not values:
        return "?"
    m = f"{statistics.median(values):g}"
    return m if k == 1 else f"{m} ({min(values):g}–{max(values):g})"


print("| Case | Harness | Good run | Unneeded questions | Missed decisions | Edits outside Files/Also changed | Done claims without evidence | Listed documents not read | Words read |")
print("| --- | --- | --- | --- | --- | --- | --- | --- | --- |")
notes = []
for c in cases:
    for h in harnesses:
        name = {"claude": "Claude Code", "codex": "Codex"}.get(h, h)
        counts = {col: [] for col in columns}
        setup_only, passed = False, 0
        for r in range(1, k + 1):
            d, where = f"{out}/{c}-{h}-{r}", f"({c}) {name} run {r}"
            if os.path.exists(f"{d}/setup-only"):
                setup_only = True
                continue
            if not os.path.exists(f"{d}/words.txt"):
                failed = open(f"{d}/failed.txt").read().strip() if os.path.exists(f"{d}/failed.txt") else "it did not finish"
                notes.append(f"{where} failed: {failed}")
                continue
            for key, label in judged.items():
                values = [v for v in (judging(f"{d}/judge-{i}.json", key) for i in range(1, j + 1)) if v is not None]
                if values:
                    counts[key].append(statistics.median(values))
                if len(set(values)) > 1:
                    notes.append(f"{where}: {label} judged {', '.join(map(str, values))}")
            failed_checks = [line.strip() for line in open(f"{d}/good.txt") if line.strip()]
            passed += not failed_checks
            notes += [f"{where}: {line}" for line in failed_checks]
            counts["outside"].append(sum(1 for line in open(f"{d}/outside.txt") if line.strip()))
            counts["words"].append(int(open(f"{d}/words.txt").read()))
        if setup_only:
            cells = ["setup only"] + [""] * 6
        elif not counts["words"]:
            cells = ["failed"] + [""] * 6
        else:
            cells = [f"{passed}/{k}"] + [spread(counts[col]) for col in columns]
        print(f"| ({c}) | {name} | " + " | ".join(cells) + " |")
if notes:
    print()
    print("\n".join(notes))
EOF
