# Vendored: dmmulroy/anti-slop

Opinionated Oxlint rules that reject low-evidence TypeScript/JavaScript patterns.

- Upstream: https://github.com/dmmulroy/anti-slop
- Commit: `95a56e5d24fb3d849673c2d51eb0908b8bd2d33b`
- License: MIT (see upstream `LICENSE`)
- Vendored on: 2026-09-10
- Pinned toolchain: `oxlint@1.78.0`, `@oxlint/plugins@1.78.0`

This directory is a copy of upstream `src/` plus a scoped `package.json`
(`{"type":"module"}`) so Node parses the `.ts` entry as ESM without setting
`"type": "module"` on the app itself. It is excluded from linting via
`ignorePatterns` in `.oxlintrc.json`.

## Local policy

The rules are configured in `web/next-app/.oxlintrc.json`, all as `"warn"`
(report-only; never fails the build):

- `anti-slop/require-safety-comment-for-type-assertion`: **off**. It flagged
  115 assertions across the app and is a stylistic discipline, not a defect
  signal; revisit as an opt-in later.
- `anti-slop/no-module-mocking`: **off in test files** (`test/**`,
  `**/*.test.*`, `e2e/**`) — `vi.mock` is our deliberate test seam.
- `anti-slop/no-runtime-typeof`: **off**. Browser platform guards
  (`typeof window === 'undefined'`) and React `ReactNode` narrowing are
  idiomatic and not "representation without a contract"; the rule assumes a
  Node I/O boundary we do not have in the web app.
- `anti-slop/no-shape-in-symbol-names`: **off**. False positives on the SVG
  DOM property `shapeRendering` and Zod's `.shape` accessor.

## Updating

Vendored by hand, not a dependency. To update, re-copy upstream `src/`, diff,
and port reviewed changes. Do not `npm install` it.
