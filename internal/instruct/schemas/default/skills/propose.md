**Propose says "We should do this."** It turns a chosen direction into a
change the owner can review: a proposal, the change to the requirements of
each affected capability, and tasks with acceptance expectations. It does
not plan or perform the implementation.

### 1. Start from what is already said

- From this conversation: use what the owner has said here. Do not ask them
  to repeat it.
- From the owner's description below: read the code and the documents in the
  project index that it touches before you say anything.

### 2. Confirm the scope and ask one round, in one message

- Say in a few sentences what you understand is being proposed: the problem,
  the intended outcome, what is in scope and what is not.
- If the direction turns out not to be settled, say so. The work then belongs
  back in exploring, not in a proposal that hides the gap.
- The decisions the documents need that this conversation has not settled
  are which option, which requirements and scenarios, where each task goes
  in the code, what is out of scope and which files outside the code
  change. Ask at most five that matter most, those that change what users see,
  what is in scope, compatibility, or that could lose data, such as an exit
  status, before architecture or interfaces costly to reverse, such as
  which function holds a check. Routine choices inside the agreed approach,
  such as a test's name and place, are not asked: they are open questions.
- More than five that matter means the idea is likely more than one change;
  offer the split in the same message.
- Where small choices are not worth a question while building, such as
  test layout, private helpers or a new test file, you may recommend
  delegating them to the implementer; a task carries them on a `Delegated:`
  line only if the owner agrees.

### 3. Write the change

Write every decision in the proposal, however much room it needs. Each
task names the files it touches and the interfaces
it produces, as the tasks guide says, so where the change goes in the code
is settled here, not left to whoever builds it. No code, commands or steps,
no sequencing, no estimates.

### 4. Check the change

Read the documents you wrote side by side. Name, for each ADDED or
MODIFIED requirement, the task whose acceptance shows it; a requirement
with none needs an acceptance line or a task. A file a task names that the
scope puts out, or a scenario the requirement's sentence does not cover,
is a contradiction: one side is wrong, and which one is the owner's to say
unless the conversation already settled it.
