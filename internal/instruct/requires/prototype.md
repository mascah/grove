- A prototype answers one question before a change is proposed. Write no
  Grove document and change nothing under the document root.
- Work on a branch `prototype/<slug>`, `<slug>` being a few lowercase words
  of the question joined by hyphens. If the session is on another branch,
  create it from the current commit. Never switch to the target branch or
  commit on it; the prototype branch is never merged into it.
- Add nothing outside the prototype branch, and no dependency beyond what
  the experiment needs.
- What the experiment shows is observed behavior: evidence for the owner,
  not agreement. Whether to propose is the owner's to say.
- Finish: ask the owner whether to commit the experiment on its branch, so
  a proposal can cite its commit. Commit nothing otherwise, and push
  nothing.
- When the owner asks for a proposal, switch back to the branch the session
  started on, then run Grove's `instructions propose` with no argument and
  follow what it prints, here.
