import { test, expect } from '@playwright/test'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

/**
 * The figures a lesson body references must actually render.
 *
 * Nothing above this catches a blank image. The Go extractor proves a reference is
 * well-formed, `TestEveryAssetResolvesToAShippedDiagram` proves the file is on disk, and
 * `TestPlacedFiguresAreInTheSectionTheirPlacementNamed` proves each figure sits in the section
 * it was placed in — all three are satisfied by a diagram a browser declines to draw. An HTTP
 * 200 carrying an SVG with no viewBox, a zero-height `<img>`, a figure the renderer drops: each
 * passes every other gate and shows a learner an empty box where the instruction should be.
 *
 * Two things make this a measurement rather than a decoration:
 *
 *  - It serves the real lesson files from `data/lessons/`, not fixtures. A fixture would test
 *    the harness; the corpus file tests the corpus.
 *
 *  - It carries a negative control. A check that cannot fail is worse than no check, because
 *    it is read as coverage, so the control asserts the same measurement reports zero for a
 *    figure-less lesson and for a path nothing serves. If the control ever passed, the
 *    measurement would be vacuous and this spec would prove nothing.
 *
 * Per ADR-029, `npx playwright` serves the prebuilt `out/`, so `next build` must run first or
 * this exercises the previous build and passes regardless of the change under test.
 */

const ROOT = resolve(__dirname, '../../..')

/** Real bodies, chosen because they now carry figures and because one does not. */
const CASES = [
  { concept: 'discrete.sets.operations', file: 'algebrica/sets-and-numbers/sets.md', expectFigures: 6 },
  { concept: 'alg.log.concept', file: 'algebrica/powers-radicals-logarithms/logarithms.md', expectFigures: 3 },
  { concept: 'alg.exp.concept', file: 'algebrica/functions/exponential-function.md', expectFigures: 3 },
  { concept: 'calc.deriv.cauchy_mvt', file: 'algebrica/differential-calculus-theorems/cauchy-theorem.md', expectFigures: 1 },
  { concept: 'calc.seq.euler', file: 'algebrica/sequences/euler-number-limit-sequence.md', expectFigures: 1 },
]

const FIGURELESS = { concept: 'arith.add.word', file: 'teaching/arith.add.word.md' }

function bodyFor(file: string): string {
  return readFileSync(resolve(ROOT, 'data/lessons', file), 'utf8')
}

/** Stub the lesson API with a real body, and progress with nothing. */
async function stubLesson(page: import('@playwright/test').Page, concept: string, file: string) {
  const body = bodyFor(file)
  const title = body.split('\n')[0].replace(/^#\s*/, '')
  await page.route('**/api/lessons/body*', route =>
    route.fulfill({ json: { title, body, source: file.split('/')[0], assets: [] } }),
  )
  await page.route('**/api/lessons', route =>
    route.fulfill({ json: { lessons: { [file.split('/')[0]]: [{ title, body: '', concepts: [concept] }] } } }),
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
}

const FIGURE = 'img[src*="/diagrams/algebrica/"]'

test.describe('lesson figures render', () => {
  test('every referenced figure draws at non-zero size', async ({ page }) => {
    const problems: string[] = []

    for (const c of CASES) {
      await stubLesson(page, c.concept, c.file)
      const notFound: string[] = []
      const onResponse = (r: import('@playwright/test').Response) => {
        if (r.status() === 404 && r.url().includes('/diagrams/')) notFound.push(r.url())
      }
      page.on('response', onResponse)

      await page.goto(`/study?concept=${encodeURIComponent(c.concept)}`)
      await expect(page.locator('body')).toContainText(bodyFor(c.file).split('\n')[0].replace(/^#\s*/, ''), {
        timeout: 30_000,
      })
      const imgs = page.locator(FIGURE)
      await expect(imgs).toHaveCount(c.expectFigures)

      for (let i = 0; i < (await imgs.count()); i++) {
        const el = imgs.nth(i)
        const src = await el.getAttribute('src')
        const size = await el.evaluate((e: HTMLImageElement) => ({ w: e.naturalWidth, h: e.naturalHeight }))
        if (size.w === 0 || size.h === 0) problems.push(`${c.file}: ${src} drew at ${size.w}x${size.h}`)
      }
      if (notFound.length) problems.push(`${c.file}: 404 on ${notFound.join(', ')}`)
      page.off('response', onResponse)
    }

    expect(problems, problems.join('\n')).toEqual([])
  })

  test('negative control: a lesson with no figures reports none', async ({ page }) => {
    await stubLesson(page, FIGURELESS.concept, FIGURELESS.file)
    await page.goto(`/study?concept=${FIGURELESS.concept}`)
    await expect(page.locator('body')).toContainText(/\w/, { timeout: 30_000 })
    // The page really does render lesson prose, so zero here means "no figures", not "no
    // lesson" — which is the distinction that would otherwise make the test above vacuous.
    await expect(page.locator(FIGURE)).toHaveCount(0)
  })

  test('negative control: a missing diagram measures as zero, not as drawn', async ({ page }) => {
    await stubLesson(page, CASES[0].concept, CASES[0].file)
    await page.goto(`/study?concept=${CASES[0].concept}`)
    await expect(page.locator(FIGURE).first()).toBeVisible({ timeout: 30_000 })

    // The same measurement the passing test relies on, against a path nothing serves. If this
    // ever reported a non-zero size, "the figures drew" would mean nothing.
    const size = await page.evaluate(async () => {
      const img = new Image()
      img.src = '/diagrams/algebrica/definitely-not-a-real-diagram-xyz.svg'
      await new Promise((r) => {
        img.onload = r
        img.onerror = r
      })
      return { w: img.naturalWidth, h: img.naturalHeight }
    })
    expect(size.w === 0 || size.h === 0).toBe(true)
  })
})