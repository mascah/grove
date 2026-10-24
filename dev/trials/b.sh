# (b) An ambiguous feature with one consequential open question: whether a
# deleted record's ID may be given to a new record.

setup() { :; }

runs=(
  "delete => /grove-propose Add ledger delete ID, from the README's planned list"
)
answers=(
  "recycl|next_id|same id|deleted (record'?s? )?ids?\b|\bids? of (an? |the )?deleted|\bids? ((be|are|get|is|never) )*re-?us|re-?us[a-z]* (an? |the |its |their |old |deleted |freed |that |highest |max(imum)? )*ids?\b => Yes to the rest, but a deleted record's ID is never given to a new record."
)
