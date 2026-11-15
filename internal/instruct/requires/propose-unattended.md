With no one in the session, these replace every line above and step below
that has you confirm the scope, ask the owner, wait, or offer
`/grove-implement`:

- Ask and wait for no one. Work on a branch other than the target.
- The description below is the request. If it does not say what outcome
  is wanted, write no document; say that it names none and what would.
- Write one change, in the schema these instructions follow, with no
  `depends_on` the request does not name. A split you would offer, an
  open change that may need to land first, or another listed schema, such
  as `bugfix` for wrong behavior, is an open question where there is one.
- Decisions quote only what the request says, or say that none was made
  beyond it. Write each decision of the round as an open question, in its
  order, with your recommended answer, its reason and the task it decides.
  Write no `Delegated:` line.
- What the check of the change would ask is an open question too, naming
  the task; only such a question says it must be answered before review;
  change neither side.
- Finish: name the files you wrote, the qualified task IDs and the open
  questions, and say that merging this branch approves the change. Offer
  no `/grove-implement` on this branch: a draft nobody attended is
  approved alone, never with its code.
