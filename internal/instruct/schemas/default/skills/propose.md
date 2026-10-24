**Propose says "We should do this."** It turns a chosen direction into a
change the owner can review: a proposal, the change to the requirements of
each affected capability, and tasks with acceptance expectations. It
describes future work. It does not plan or perform the implementation.

### 1. Start from what is already said

- From this conversation: use what the owner has said here. Do not ask them
  to repeat it.
- From the owner's description below: read the code and the documents in the
  project index that it touches before you say anything.

### 2. Confirm the scope with the owner

- Say in a few sentences what you understand is being proposed: the problem,
  the intended outcome, what is in scope and what is not. Ask the owner to
  confirm or correct it before you write.
- If several options are still open, ask which one.
- If the idea is too big for one change, or has parts the owner could agree
  separately, say so and offer to write them as changes linked with
  `depends_on`.
- If the direction turns out not to be settled, say so. The work then belongs
  back in exploring, not in a proposal that hides the gap.

### 3. Ask one round

- List every decision the documents need that this conversation has not
  settled: which option, which requirements and scenarios, where each task
  goes in the code, what is out of scope, which files outside the code
  change. Number them, each with your recommended answer and its reason,
  worded so that "yes" accepts it, and wait.
- Put first the decisions that change what users see, what is in scope,
  compatibility, or that could lose data, such as an exit status; then
  architecture or interfaces costly to reverse, such as which function
  holds a check; last, on one line of recommended defaults, not a list,
  routine choices inside the agreed approach, such as a test's name and
  place.
- Where small choices are not worth a question while building, such as
  test layout, private helpers or a new test file, you may recommend
  delegating them to the implementer; a task carries them on a `Delegated:`
  line only if the owner agrees.
- A decision the owner leaves open is an open question, naming the task it
  blocks; one the owner answers is a decision, in their words.

### 4. Write the change

Write `design.md` only when a decision needs more room than the proposal's
decisions section. Each task names the files it touches and the interfaces
it produces, as the tasks guide says, so where the change goes in the code
is settled here, not left to whoever builds it. No code, commands or steps,
no sequencing, no estimates.

### 5. Check the change

Read the documents you wrote side by side. Name, for each ADDED or
MODIFIED requirement, the task whose acceptance shows it; a requirement
with none needs an acceptance line or a task. A file a task names that the
scope puts out, or a scenario the requirement's sentence does not cover,
is a contradiction: one side is wrong, and which one is the owner's to say
unless the conversation already settled it.
