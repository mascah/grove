# Changelog

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
