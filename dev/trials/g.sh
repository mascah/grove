# (g) Explore asked to recommend: it should read first, recommend, write
# nothing, and take few of the owner's replies.

setup() { :; }

runs=(
  "explore => /grove-explore Should ledger keep its records in SQLite instead of records.json? Look into it and recommend."
)
answers=(
  "recommend|short answer|\bi('d| would) (keep|go|stay|suggest|choose|pick|use) => Yes, I'll go with your recommendation."
)

# good checks that the run wrote nothing, that it recommended, either getting
# the scripted answer to a message holding a recommendation or ending on
# one, and that it took at most three of the owner's replies.
good() {
  local changed replies
  changed=$(git diff --name-only main...explore | tr '\n' ' ')
  [ -z "$changed" ] || echo "the run wrote files: $changed"
  [ -n "$(answered "I'll go with your recommendation")" ] || grep -qi recommend <<< "$last" || echo "the agent recommended nothing"
  replies=$(($(wc -l < "$dir/turns.txt") - 1))
  [ "$replies" -le 3 ] || echo "the run took $replies of the owner's replies"
}
