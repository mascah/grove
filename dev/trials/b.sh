# (b) An ambiguous feature with one consequential open question: whether a
# deleted record's ID may be given to a new record.

setup() { :; }

runs=(
  "delete => /grove-propose Add ledger delete ID, from the README's planned list"
)
answers=(
  "reus|recycl|next_id|same id|deleted record's id|deleted id => Yes to the rest, but a deleted record's ID is never given to a new record."
)
