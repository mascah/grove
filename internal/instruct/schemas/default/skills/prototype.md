**Prototype finds out by building.** A throwaway experiment answers one
question: how something feels, whether something is feasible, or how some
logic behaves. An answer that the idea is wrong is a complete result.

### 1. Settle the question

- Take the question from the owner's description below, or from this
  conversation. If there is none, ask what they want to find out.
- Say in a few sentences the question, its limits (what the experiment
  leaves out) and the evidence that would answer it, and ask the owner to
  confirm or correct them. If the conversation already settles them, say
  them and go on. Ask only what would change the experiment.

### 2. Build the simplest thing that shows it

- The smallest artifact that exposes the answer: a script for a technical
  property, the states and transitions for logic, a runnable screen for how
  something feels.
- Use the project's own tooling where you can. Keep it apart from the
  product's code, such as in `prototypes/<slug>/`, and mark it throwaway.
- No production routes, credentials, writes to real data, or polish that
  adds no evidence.
- One command runs it and shows the state or result that matters.

### 3. Run it

Run representative cases, including those where the idea could fail.
Evidence can be an assertion, a measurement, or what the owner sees: for
how something feels, ask the owner to try it and tell you. Note each mock
or simplification and what it cannot establish. Keep the evidence of a
failure when it changes the verdict.

### 4. Report

Tell the owner the question, the command that runs the experiment, what
was observed, your verdict, and what the experiment cannot establish. A
finding reaches the target only in a proposal; a prototype that kills the
idea leaves only its branch and this conversation.
