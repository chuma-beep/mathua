import { test, expect, type Page } from './fixtures'

// A knowledge point that names a visual, without the production API. The
// dispatcher, template registry and lazy chunk all run exactly as they do in a
// lesson, which is what this spec is for: the other visual tests render the
// templates directly and never exercise the loading path or the lesson mount.

const TANGENT = 'calc.deriv.tangent_line'
const MATRIX = 'linalg.matrix.mult'

const tangentVisual = {
  type: 'calc.tangent',
  title: 'Tangent and secant at x',
  caption: 'Move x and h: the secant slope approaches the tangent slope as h shrinks.',
  params: { coeffs: [0, 0, 1], x: 1, h: 0.8, window: [-2, 3] },
}

// The same params the lesson shard carries, so this covers the figure a learner
// actually meets rather than a convenient matrix for the test.
const matrixVisual = {
  type: 'linalg.matrix2d',
  title: 'Matrix multiplying a vector',
  caption: 'The columns are the arrows, and the shaded square has area |det A|.',
  params: { matrix: [1.2, 0.6, 0.4, 1.4], vector: [2, 1], showGrid: true },
  predict: {
    question: 'det A is 1.44. What does that say about the image of the unit square?',
    options: ['Its area is unchanged', 'Its area is 1.44 times the original'],
    answer: 1,
    explain: "The determinant is the area scale factor, so the unit square's image has area 1.44.",
  },
}

async function stubVisualStep(
  page: Page,
  concept: string,
  title: string,
  visual: Record<string, unknown>,
) {
  await page.route('**/api/lessons/body*', (route) =>
    route.fulfill({ json: { title, body: '', source: 'teaching', assets: [] } }),
  )
  await page.route('**/api/lessons/*/kp*', (route) =>
    route.fulfill({
      json: {
        concept_id: concept,
        diagram: '',
        kps: [
          {
            label: title,
            section: title,
            subgoals: [],
            worked_example: '',
            visual,
          },
        ],
      },
    }),
  )
  await page.route('**/api/lessons/*/readiness*', (route) =>
    route.fulfill({
      json: { concept_id: concept, eligible: true, ready: true, weak: [], missing: [], prerequisites: [] },
    }),
  )
  await page.route('**/api/lessons/*/practice*', (route) =>
    route.fulfill({ json: { concept_id: concept, questions: [] } }),
  )
  await page.route('**/api/progress/**', (route) => route.fulfill({ json: {} }))
  await page.route('**/api/lessons*', (route) =>
    route.fulfill({
      json: { lessons: { teaching: [{ title, body: '', concepts: [concept] }] } },
    }),
  )
}

/** The frame must not overflow the viewport; a horizontal scrollbar on the
 *  lesson page is the visible failure of a fixed-width SVG. */
async function expectNoOverflow(page: Page) {
  const overflow = await page.evaluate(
    () => document.documentElement.scrollWidth - document.documentElement.clientWidth,
  )
  expect(overflow).toBeLessThanOrEqual(1)
}

test('mounts the tangent visual and reacts to its native slider', async ({ page }, testInfo) => {
  await stubVisualStep(page, TANGENT, 'Tangent lines', tangentVisual)
  await page.goto('/learn?concept=' + encodeURIComponent(TANGENT))

  const readout = page.getByTestId('visual-readout')
  await expect(readout).toBeVisible({ timeout: 20_000 })
  await expect(readout).toContainText('secant')

  const before = await readout.textContent()
  const slider = page.getByLabel('x', { exact: true })
  await slider.focus()
  for (let i = 0; i < 10; i++) await slider.press('ArrowRight')
  await expect(readout).not.toHaveText(before ?? '')

  await expectNoOverflow(page)

  await page.screenshot({ path: `/tmp/opencode/visual-tangent-${testInfo.project.name}.png` })
})

test('mounts the matrix figure and plays it back to A', async ({ page }, testInfo) => {
  await stubVisualStep(page, MATRIX, 'Matrix multiplying a vector', matrixVisual)
  await page.goto('/learn?concept=' + encodeURIComponent(MATRIX))

  const readout = page.getByTestId('visual-readout')
  await expect(readout).toBeVisible({ timeout: 20_000 })
  // The inverse is the readout the compact revision had dropped; it is back, and
  // this is the only test that sees it through the lazy chunk and the dispatcher.
  await expect(readout).toContainText('A⁻¹')
  await expect(readout).toContainText('full parallelogram')
  await expect(page.locator('svg[aria-label^="A equals the matrix"]')).toBeVisible()
  // Two options, which is the approved prediction block rather than the old three.
  await expect(page.getByRole('radio')).toHaveCount(2)

  // The animation is opt-in: nothing is scheduled until Play is pressed.
  const t = page.getByLabel('t', { exact: true })
  await expect(t).toHaveValue('1')
  await page.getByRole('button', { name: 'Play' }).click()
  await expect(page.getByRole('button', { name: 'Pause' })).toBeVisible()
  await expect(page.locator('span', { hasText: /^M\(t\), t =/ })).toBeVisible()

  await page.getByRole('button', { name: 'Pause' }).click()
  await page.getByRole('button', { name: 'Show A' }).click()
  await expect(t).toHaveValue('1')
  await expect(page.getByRole('button', { name: 'Play' })).toBeVisible()
  // A's own readout is untouched by the intermediate state.
  await expect(readout).toContainText('det A')
  await expect(readout).toContainText('1.44')

  await expectNoOverflow(page)

  await page.screenshot({ path: `/tmp/opencode/visual-matrix-${testInfo.project.name}.png` })
})
