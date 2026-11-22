# (g) Explore asked to recommend: it should read first, recommend, write
# nothing, and take few of the owner's replies.

setup() { :; }

# A recommendation chooses an option, whatever words or layout surround it.
recommends="\b(keep|keeping|stay with|stick with|move to|moving to|switch to|switching to|use|adopt) [*\`]*(records\.json|json|sqlite)"

runs=(
  "explore => /grove-explore Should ledger keep its records in SQLite instead of records.json? Look into it and recommend."
)
answers=(
  "$recommends => Yes, I'll go with your recommendation."
)

# good checks that the run wrote nothing, that it read ledger/records.py in
# its first turn, that one of its messages recommends an option, whether the
# owner then answered or it finished without asking, and that it took at
# most three of the owner's replies.
good() {
  local changed replies
  changed=$(git diff --name-only main...explore | tr '\n' ' ')
  [ -z "$changed" ] || echo "the run wrote files: $changed"
  awk '/^OWNER:/ {n++} n == 1' "$dir/transcript.txt" | grep -Eq '^TOOL:.*(records\.py|ledger/\*)' || echo "the agent did not read ledger/records.py before its first reply"
  awk '/^AGENT:/ {on = 1} /^(TOOL|RESULT|OWNER)/ {on = 0} on' "$dir/transcript.txt" | grep -Eiq "$recommends" || echo "the agent recommended no option"
  replies=$(($(wc -l < "$dir/turns.txt") - 1))
  [ "$replies" -le 3 ] || echo "the run took $replies of the owner's replies"
}
