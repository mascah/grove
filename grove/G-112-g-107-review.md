---
id: "G-112"
type: review
title: "G-107 documentation reconciliation review"
status: current
created: "2026-09-23T16:57:58Z"
updated: "2026-09-23T16:59:21Z"
work: ["G-107"]
examined: "6a9f146"
---

## Examined

Round 1: `git diff 768efab d6cc1df`, against [G-107](G-107-current-documentation.md)'s
acceptance and the inventory in plan [G-111](G-111-g-107-docs-plan.md).
Round 2: `git diff d6cc1df 1398458`, the fixes. Reviewer: a fresh `reviewer`
subagent in the implementing Claude Code session (Opus 5.5), read-only, with
no part in writing the diff. It is an agent, not the owner.
The last content commit, `6a9f146`, changes one word the reviewer proposed
in round 2 and was self-checked, not reviewed.

## Findings

Round 1 found nothing high-severity. It verified the claims against code,
`lefthook.yml`, CI, `--help` and the cited records, and found no lost
constraint other than 1 and 3. Its link and anchor check found no problems.

1. Medium: AGENTS.md dropped "paths stay stable" from the old identity rule.
2. Low–medium: the README opening said the board covers "all of the above",
   but it has no create, convert or context path.
3. Low: sibling repositories lost "not *automatically*" in the write scope and
   "Follow their instructions", which conflicted with the brief.
4. Low: the brief said nullsec cut over *after* G-036 closed. G-041 was a
   member, so the cutover came first.
5. Low: the allocation text named the floor scan as `.md` files on "every
   local ref". The code greps every text file on branches, remote-tracking
   refs and tags. It also attached the unrecoverable-reservation notice to
   both notices, though only the initialization notice says that.
6. Low: the archive path `../grove-archive-2026-09-18/` resolved wrongly from
   `grove/`.
7. Low, inherited: "The first three commands read live files" (the first is
   now the board); "Go 1.26 or later" while `go.mod` says 1.26.8.

Round 2: all seven are fixed. One new low issue: `record-model.md:63` said
"record files" where `:321` and the code say any text file.

## Disposition

- Findings 1–7: fixed in `1398458` and confirmed by round 2.
- Round 2 issue: fixed in `6a9f146` as the reviewer proposed, self-checked.
- Two review rounds were used, within the cap of three. No findings are
  open.
- Flagged, not changed, since it is code rather than documentation:
  `grove --help` says the board "Reads only". That needs follow-up work.
