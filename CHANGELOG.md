# Changelog

## [0.4.0](https://github.com/mascah/grove/compare/v0.3.0...v0.4.0) (2026-10-23)


### ⚠ BREAKING CHANGES

* /grove-apply <change-id> and /grove-apply with no ID are gone. Build one change with /grove-implement <change-id>; to build unattended, plan a cycle with /grove-plan, mark it unattended: true, and run /grove-apply <cycle-id> or grove run, which now builds only marked cycles. grove check needs the main branch locally when a marked cycle exists. Re-run grove init to add the grove-plan skill.
* **instruct:** implement writes proposal, change-spec and tasks. A project schema whose documents leave out one of those kinds must set skills.implement.writes in its schema.yaml, as it already does for propose.

### Features

* cycles plan the changes grove run builds ([#24](https://github.com/mascah/grove/issues/24)) ([8c0a891](https://github.com/mascah/grove/commit/8c0a891541139e708e6187c0c4337024dec45423))
* **instruct:** /grove-implement builds a change with the owner ([#21](https://github.com/mascah/grove/issues/21)) ([b4e3223](https://github.com/mascah/grove/commit/b4e322321c1cae6d35acee81b510be00243f562e))
* plan's verdict, runs without a remote, streamed output ([#27](https://github.com/mascah/grove/issues/27)) ([44f4a6b](https://github.com/mascah/grove/commit/44f4a6b953593f0ef628e89a425c1233088d8693))
* tasks name their files, check notices when they move, and propose asks a round ([#31](https://github.com/mascah/grove/issues/31)) ([0d3a90e](https://github.com/mascah/grove/commit/0d3a90e2eb2b1ef6ced6171c7527469bbc61ac9b))


### Bug Fixes

* **instruct:** implement reopens revised tasks and stops on a stale change ([#23](https://github.com/mascah/grove/issues/23)) ([90565d4](https://github.com/mascah/grove/commit/90565d4f73d28a2d3f55c82176e859334b6fcefd))

## [0.3.0](https://github.com/mascah/grove/compare/v0.2.0...v0.3.0) (2026-10-22)


### ⚠ BREAKING CHANGES

* **cli:** `grove check` reports a change that rewrites a requirement another change in the same working tree has rewritten, where it compared with the target before; the stale message no longer names the target. The apply skill archives a finished change instead of handing off to `/grove-archive`.

### Features

* **cli:** apply every approved change in one session, archiving each ([#15](https://github.com/mascah/grove/issues/15)) ([1e3e819](https://github.com/mascah/grove/commit/1e3e8192651562fc06073fdf9e32038a3b22f80c))
* **cli:** grove run launches one bounded apply session ([#18](https://github.com/mascah/grove/issues/18)) ([e80aa4e](https://github.com/mascah/grove/commit/e80aa4e4db78736470da793b22a9123c79ed6218))

## [0.2.0](https://github.com/mascah/grove/compare/v0.1.0...v0.2.0) (2026-10-20)


### ⚠ BREAKING CHANGES

* **cli:** grove.yaml no longer has context. grove check reports it as an unknown field; delete it and keep project facts in CLAUDE.md or AGENTS.md, or in rules.

### Features

* **cli:** drop context from grove.yaml ([#12](https://github.com/mascah/grove/issues/12)) ([e3826d0](https://github.com/mascah/grove/commit/e3826d0e8ee1cfd002c3366590129e2f0d80f4e7))

## 0.1.0 (2026-10-20)


### Features

* **cli:** any harness, as skills or commands; init detects them and prunes ([#7](https://github.com/mascah/grove/issues/7)) ([131da7e](https://github.com/mascah/grove/commit/131da7e8023a3af28d34ace9f15abababbd33115))
* **cli:** grove version ([#4](https://github.com/mascah/grove/issues/4)) ([34eae85](https://github.com/mascah/grove/commit/34eae8599b9c28eee837429801f6b6da46accf60))
* install script ([#5](https://github.com/mascah/grove/issues/5)) ([b6a86ad](https://github.com/mascah/grove/commit/b6a86ad03512122b3df73642e23df136f19979a4))
