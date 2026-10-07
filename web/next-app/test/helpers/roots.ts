import path from 'node:path'

/**
 * Where things are, for tests that read the repository rather than render it.
 *
 * This exists because three separate guards got the root wrong in one session, each of them
 * producing a test that passed while inspecting nothing: `ROOT` resolved one level short, the
 * scan visited no files, and "no findings" is exactly what a scan that looks at nothing reports.
 * `test/noEmoji.test.ts` failed for the same reason and needed a sentinel-file assertion to catch
 * it — which is the tell that a shared constant is worth more than four careful resolutions.
 *
 * This file is at `web/next-app/test/helpers/`, so from here the walk up is
 * `test` -> `next-app` -> `web` -> repo. Two `..` is the Next app root; four is the repo root.
 */
export const APP_ROOT = path.resolve(__dirname, '..', '..')
export const REPO_ROOT = path.resolve(__dirname, '..', '..', '..', '..')
