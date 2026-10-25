# Changelog

## [0.5.0](https://github.com/mascah/grove/compare/v0.4.0...v0.5.0) (2026-10-25)


### ⚠ BREAKING CHANGES

* **cli:** `grove check` reports a stale cycle that is not marked `unattended`. Re-plan each cycle it names with `/grove-plan <cycle-id>`.
* a maintained spec no longer has a `.grove.yaml`; `grove check` reports each `grove/specs/<id>/.grove.yaml` until it is deleted. Make sure each spec's first heading is its title and its Purpose's first sentence describes it.
* the apply skill and `grove instructions apply` are removed. Replace `/grove-apply <cycle-id>` with `/grove-implement <cycle-id> --unattended` in prompts, scripts and timers. Move lines under `rules.skills.apply` to `rules.skills.implement`, then run `grove init` again. If a project schema names `apply` under `skills` or overrides `skills/apply.md`, `grove check` reports it until that text moves to implement's guide.
* **cli:** `grove check` now reports requirements with no "shall" or no scenario, and change headings in maintained specs, on documents that were clean before. Add the missing "shall" or scenario, and fold any leftover ADDED, MODIFIED or REMOVED heading into `## Requirements`.

### Features

* **cli:** a cycle can depend on cycles ([#69](https://github.com/mascah/grove/issues/69)) ([8454724](https://github.com/mascah/grove/commit/8454724ba54d0f63c54283c7e333a1896b5b6085))
* **cli:** check requirements and the archive fold ([#57](https://github.com/mascah/grove/issues/57)) ([f749c46](https://github.com/mascah/grove/commit/f749c46a7671911d9136c06f77c9e9ab468f0b04))
* **cli:** grove run builds a cycle it is given ([#35](https://github.com/mascah/grove/issues/35)) ([8aacbb7](https://github.com/mascah/grove/commit/8aacbb78dcff4aa2fae4f165c5b8a3777933a54a))
* **cli:** one merge for small changes ([#64](https://github.com/mascah/grove/issues/64)) ([29c5c7e](https://github.com/mascah/grove/commit/29c5c7e4b8513d57674c43a3ebbb50e3eb34bf55))
* **cli:** restore context in grove.yaml ([#37](https://github.com/mascah/grove/issues/37)) ([ed415d2](https://github.com/mascah/grove/commit/ed415d23c7c97240ba6208fc008f9ba7454aceb0))
* **cli:** revising a change re-plans its cycle ([#67](https://github.com/mascah/grove/issues/67)) ([9308040](https://github.com/mascah/grove/commit/9308040009674b9c0095f9844605af5a99f8c556))
* **cli:** styled grove run output ([#36](https://github.com/mascah/grove/issues/36)) ([1441b53](https://github.com/mascah/grove/commit/1441b53ac4a50a2026d95f0be405cca487736383))
* **dev:** skill trials on the fixture, run by just trials ([#70](https://github.com/mascah/grove/issues/70)) ([ba80a70](https://github.com/mascah/grove/commit/ba80a704f4dda00f260f9f2d6136da92b86cd2c6))
* implement builds unattended, and apply goes ([#63](https://github.com/mascah/grove/issues/63)) ([9508c4b](https://github.com/mascah/grove/commit/9508c4be4b09037be4666fdca75afdcfed6c31e5))
* **instruct:** a built-in bugfix schema ([#65](https://github.com/mascah/grove/issues/65)) ([3e40d51](https://github.com/mascah/grove/commit/3e40d5104439fb2c032da868524000d25cc4f27e))
* **instruct:** fewer words in Grove's instructions ([#73](https://github.com/mascah/grove/issues/73)) ([8705c42](https://github.com/mascah/grove/commit/8705c422637e1744b76691ba46b113b51b059425))
* **instruct:** implement suggests a rules line (261023-mfgew) ([#84](https://github.com/mascah/grove/issues/84)) ([c1c058d](https://github.com/mascah/grove/commit/c1c058d4351b4a07b5a66d43675435f5f5b7d7ea))
* **instruct:** maintained specs list the changes that shaped them (261023-xkehc) ([#83](https://github.com/mascah/grove/issues/83)) ([bb85eb5](https://github.com/mascah/grove/commit/bb85eb598ba248968369dd2624060b35f9b8df0f))
* **instruct:** propose asks what matters ([#66](https://github.com/mascah/grove/issues/66)) ([1db9cae](https://github.com/mascah/grove/commit/1db9cae3a598e196281621538f293c7edb94ac07))
* **instruct:** verify before done ([#58](https://github.com/mascah/grove/issues/58)) ([aa245f1](https://github.com/mascah/grove/commit/aa245f140938be460365ca643081493968a749bf))
* specs without .grove.yaml, and no design.md ([#71](https://github.com/mascah/grove/issues/71)) ([d214b5e](https://github.com/mascah/grove/commit/d214b5ef8990457ac6ef2b4d67e281263655b539))


### Bug Fixes

* **cli:** a stale cycle is reported whether or not it is marked (261024-8ykrz) ([#81](https://github.com/mascah/grove/issues/81)) ([f110c04](https://github.com/mascah/grove/commit/f110c04bce9b3c07ca94a2d4415b78043c294190))

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
