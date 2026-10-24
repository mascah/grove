#!/usr/bin/env bash
# Runs the skill trials of dev/trials.md (CONSTITUTION.md rule 17): each
# case, or those named, in Claude Code and in Codex, on a fresh fixture,
# with this checkout's build of grove. Prints the results table.
#
#   dev/trials/run.sh [a b c d e]
#
# HARNESSES (default "claude codex") picks the harnesses, TRIALS_OUT the
# output folder, TURNS (default 8) the owner replies a run may get, and
# SETUP_ONLY=1 stops each trial once its fixture is set up.
set -euo pipefail

here=$(cd "$(dirname "$0")" && pwd)
checkout=$(cd "$here/../.." && pwd)
out=${TRIALS_OUT:-${TMPDIR:-/tmp}/grove-trials/$(date +%Y%m%d-%H%M%S)}
turns=${TURNS:-8}
cases=("$@")
[ ${#cases[@]} -gt 0 ] || cases=(a b c d e)
fallback="Yes, go with what you recommend."
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

# snapshot lists every file of the fixture with its hash, to tell what a run changed.
# A file the run moved or deleted is left out, so it shows as changed.
snapshot() { git ls-files -co --exclude-standard -z | xargs -0 shasum 2> /dev/null | sort -k2 || true; }

# turn sends msg to the agent, resuming session sid when set, and sets sid,
# last (the agent's last message) and words (the words of its messages).
turn() {
  local msg=$1 log=$dir/turn-$((n += 1)).jsonl
  echo "OWNER: $msg" >> "$dir/transcript.txt"
  if [ "$harness" = claude ]; then
    claude -p "$msg" ${sid:+--resume "$sid"} --setting-sources project,local --permission-mode acceptEdits \
      --allowedTools "$claude_tools" --output-format stream-json --verbose < /dev/null > "$log" 2>> "$dir/stderr.txt" || true
    sid=$(jq -r 'select(.type=="result") | .session_id' "$log" | tail -1)
    last=$(jq -rs '[.[] | select(.type=="result") | .result] | last // ""' "$log")
    jq -r 'select(.type=="assistant" and .parent_tool_use_id==null) | .message.content[]? |
      if .type=="text" then "AGENT: " + .text
      elif .type=="tool_use" then "TOOL: " + .name + " " + ((.input.command // .input.file_path // .input.skill // .input.pattern // "") | tostring)
      else empty end' "$log" >> "$dir/transcript.txt"
    words=$((words + $(jq -r 'select(.type=="assistant" and .parent_tool_use_id==null) | .message.content[]? | select(.type=="text") | .text' "$log" | wc -w)))
  else
    local cfg=(--json -c allow_login_shell=false)
    if [ -z "$sid" ]; then
      codex exec "${cfg[@]}" --sandbox workspace-write "$msg" < /dev/null > "$log" 2>> "$dir/stderr.txt" || true
      sid=$(jq -r 'select(.type=="thread.started") | .thread_id' "$log")
    else
      codex exec resume "${cfg[@]}" -c sandbox_mode='"workspace-write"' "$sid" "$msg" < /dev/null > "$log" 2>> "$dir/stderr.txt" || true
    fi
    last=$(jq -rs '[.[] | select(.type=="item.completed" and .item.type=="agent_message") | .item.text] | last // ""' "$log")
    jq -r 'select(.type=="item.completed") | .item |
      if .type=="agent_message" then "AGENT: " + .text
      elif .type=="command_execution" then "TOOL: " + .command + " (exit \(.exit_code))"
      elif .type=="file_change" then "TOOL: edit " + ([.changes[]?.path] | join(", "))
      else empty end' "$log" >> "$dir/transcript.txt"
    words=$((words + $(jq -r 'select(.type=="item.completed" and .item.type=="agent_message") | .item.text' "$log" | wc -w)))
  fi
  # A turn the harness failed, such as at a usage limit, fails the trial rather than scoring it.
  failed=$(jq -rs '[.[] | select(.type=="turn.failed" or (.type=="result" and .is_error==true)) | .error.message // .result // "error"] | first // ""' "$log")
  [ -z "$failed" ] || { echo "| ($c) | $h | failed: $failed | | | | |" > "$dir/row.md"; return 1; }
}

# reply sets answer to the case's first unused answer whose pattern matches
# the agent's last message, or the fallback.
reply() {
  local i
  answer=$fallback
  for i in "${!answers[@]}"; do
    if [ -z "${used[$i]:-}" ] && grep -Eiq -- "${answers[$i]%% => *}" <<< "$last"; then
      used[$i]=1 answer=${answers[$i]#* => }
      return
    fi
  done
}

# waiting asks a small model whether the agent's last message waits for the
# owner's answer, as against reporting what it did or offering more work.
waiting() {
  local verdict
  verdict=$(claude -p "Here is the last message an AI coding agent sent the person it works for. Does it stop to wait for that person to answer or decide something before it can go on with the task it was given? Offers of optional further work, such as committing or a next step, do not count. Reply with only yes or no.

$last" --model haiku --setting-sources project,local --tools "" --output-format json < /dev/null 2>> "$dir/stderr.txt" | jq -r '.result')
  echo "turn $n: waiting: $verdict" >> "$dir/waits.txt"
  grep -qi '^yes' <<< "$verdict"
}

# changes lists the change folders, open and archived, so a change a run
# proposes and archives is found too.
changes() { ls grove/changes grove/changes/archive 2> /dev/null | grep -E '^[0-9]{6}-' | sort || true; }

# allowed lists the paths the Files: and Delegated: lines of every tasks.md name.
allowed() {
  find grove/changes -name tasks.md -exec awk '/^(Files|Delegated):/{on=1} /^$/{on=0} on' {} + |
    grep -o '`[^`]*`' | tr -d '`' | sort -u
}

# trial runs case c in harness h and writes its row of the table.
trial() {
  local c=$1 h=$2
  dir=$out/$c-$h harness=$h n=0 words=0
  local repo=$dir/repo
  mkdir -p "$dir"
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
  # A case sets check_ok to the problems its setup means grove check to report.
  if ! grove check > "$dir/setup-check.txt" 2> /dev/null && grep -vqE "${check_ok:-^$}" "$dir/setup-check.txt"; then
    echo "($c) $h: the setup is not clean: $(cat "$dir/setup-check.txt")" >&2
    return 1
  fi
  [ -z "${SETUP_ONLY:-}" ] || { echo "| ($c) | $h | setup only | | | | |" > "$dir/row.md"; return 0; }
  local outside=() run branch prompt before after had new
  for run in "${runs[@]}"; do
    branch=${run%% => *} prompt=${run#* => }
    [ "$(git branch --show-current)" = "$branch" ] || { git diff --quiet HEAD && [ -z "$(git ls-files -o --exclude-standard)" ] || commit "After the last run"; git switch -q "$branch" 2> /dev/null || git switch -q -c "$branch" main; }
    had=$(changes)
    prompt=${prompt//\{id\}/$id}
    [ "$h" = codex ] && prompt=${prompt/\/grove-/\$grove-}
    before=$(snapshot) sid="" last=""
    turn "$prompt"
    local k=0
    # An unattended run has no owner to answer it.
    [[ $prompt == *--unattended* ]] || while [ $((k += 1)) -le "$turns" ] && waiting; do reply && turn "$answer"; done
    after=$(snapshot)
    while read -r p; do
      [[ $p == grove/* ]] || grep -qxF "$p" <(allowed) || outside+=("$p")
    done < <(diff <(echo "$before") <(echo "$after") | sed -n 's/^[<>] [0-9a-f]*  //p' | sort -u)
    new=$(comm -13 <(echo "$had") <(changes) | head -1 || true)
    [ -z "$new" ] || id=${new:0:12}
  done
  printf '%s\n' "${outside[@]:-}" | sed '/^$/d' | sort -u > "$dir/outside.txt"
  judge "$c" > "$dir/judge.json" || true
  local j=$dir/judge.json name
  name=$([ "$h" = claude ] && echo "Claude Code" || echo Codex)
  printf '| (%s) | %s | %s | %s | %s | %s | %s |\n' "$c" "$name" \
    "$(jq -r '.unneeded_questions.count // "?"' "$j" 2> /dev/null || echo '?')" \
    "$(jq -r '.missed_decisions.count // "?"' "$j" 2> /dev/null || echo '?')" \
    "$(wc -l < "$dir/outside.txt" | tr -d ' ')" \
    "$(jq -r '.done_claims_without_evidence.count // "?"' "$j" 2> /dev/null || echo '?')" \
    "$words" > "$dir/row.md"
}

# judge asks a model for the three judged counts of the run in $dir.
judge() {
  local prompt
  prompt="You judge one run of Grove's skill trials. Below are how each count is defined, the case, and the run's transcript: OWNER lines are what the tester typed, AGENT lines the agent's messages, TOOL lines its tool calls. Count only from the transcript.

$(awk '/^## The counts/{on=1;next} /^## /{on=0} on' "$checkout/dev/trials.md")

$(awk -v h="### ($1)" 'index($0,h)==1{on=1} on&&/^##/&&index($0,h)!=1{exit} on' "$checkout/dev/trials.md")

Transcript:

$(cat "$dir/transcript.txt")

Reply with only this JSON, each why one sentence naming what you counted:
{\"unneeded_questions\": {\"count\": 0, \"why\": \"\"}, \"missed_decisions\": {\"count\": 0, \"why\": \"\"}, \"done_claims_without_evidence\": {\"count\": 0, \"why\": \"\"}}"
  claude -p "$prompt" --setting-sources project,local --tools "" --output-format json < /dev/null 2>> "$dir/stderr.txt" |
    jq -r '.result' | python3 -c 'import sys, json; t = sys.stdin.read(); print(json.dumps(json.loads(t[t.index("{"):t.rindex("}") + 1]), indent=2))'
}

mkdir -p "$out"
for c in "${cases[@]}"; do
  for h in ${HARNESSES:-claude codex}; do
    (trial "$c" "$h") &
  done
done
wait
echo "Trials at $commit_of ($out)"
echo
echo "| Case | Harness | Unneeded questions | Missed decisions | Edits outside Files/Delegated | Done claims without evidence | Words read |"
echo "| --- | --- | --- | --- | --- | --- | --- |"
for c in "${cases[@]}"; do for h in ${HARNESSES:-claude codex}; do cat "$out/$c-$h/row.md" 2> /dev/null || echo "| ($c) | $h | failed | | | | |"; done; done
