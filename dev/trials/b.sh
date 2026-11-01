# (b) An ambiguous feature with one consequential open question: whether a
# deleted record's ID may be given to a new record.

setup() { :; }

runs=(
  "delete => /grove-propose Add ledger delete ID, from the README's planned list"
)
answers=(
  "recycl|next_id|same id|deleted (record'?s? )?ids?\b|\bids? of (an? |the )?deleted|\bids? ((be|are|get|is|never) )*re-?us|re-?us[a-z]* (an? |the |its |their |old |deleted |freed |that |highest |max(imum)? )*ids?\b => Yes to the rest, but a deleted record's ID is never given to a new record."
)

# good checks that the agent asked about a deleted record's ID, so got its
# scripted answer, and changed nothing outside grove/.
good() {
  [ -n "$(answered "a deleted record's ID is never given to a new record")" ] || echo "the agent did not ask whether a deleted record's ID may be given to a new record"
  local changed
  changed=$(git diff --name-only main...delete -- . ':!grove' | tr '\n' ' ')
  [ -z "$changed" ] || echo "the run changed files outside grove/: $changed"
}
