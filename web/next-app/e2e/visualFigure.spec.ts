import { test, expect, type Page } from './fixtures'

// A knowledge point that names a visual, without the production API. The
// dispatcher, template registry and lazy chunk all run exactly as they do in a
// lesson, which is what this spec is for: the other visual tests render the
// templates directly and never exercise the loading path or the lesson mount.
const CONCEPT = 'calc.deriv.tangent_line'

async function stubVisualStep(page: Page) {
  await page.route('**/api/lessons/body*', (route) =>
    route.fulfill({ json: { title: 'Tangent lines', body: '', source: 'teaching', assets: [] } }),
  )
  await page.route('**/api/lessons/*/kp*', (route) =>
    route.fulfill({
      json: {
        concept_id: CONCEPT,
        diagram: '',
        kps: [
          {
            label: 'Tangent line equation',
            section: 'Equation of the tangent line',
            subgoals: [],
            worked_example: 'The tangent at $x$ has slope $f\'(x)$.',
            visual: {
              type: 'calc.tangent',
              title: 'Tangent and secant at x',
              caption: 'Move x and h: the secant slope approaches the tangent slope as h shrinks.',
              params: { coeffs: [0, 0, 1], x: 1, h: 0.8, window: [-2, 3] },
            },
          },
        ],
      },
    }),
  )
  await page.route('**/api/lessons/*/readiness*', (route) =>
    route.fulfill({
      json: { concept_id: CONCEPT, eligible: true, ready: true, weak: [], missing: [], prerequisites: [] },
    }),
  )
  await page.route('**/api/lessons/*/practice*', (route) =>
    route.fulfill({ json: { concept_id: CONCEPT, questions: [] } }),
  )
  await page.route('**/api/progress/**', (route) => route.fulfill({ json: {} }))
  await page.route('**/api/lessons*', (route) =>
    route.fulfill({
      json: { lessons: { teaching: [{ title: 'Tangent lines', body: '', concepts: [CONCEPT] }] } },
    }),
  )
}

test('mounts the tangent visual and reacts to its native slider', async ({ page }, testInfo) => {
  await stubVisualStep(page)
  await page.goto('/learn?concept=' + encodeURIComponent(CONCEPT))

  const readout = page.getByTestId('visual-readout')
  await expect(readout).toBeVisible({ timeout: 20_000 })
  await expect(readout).toContainText('secant')

  const before = await readout.textContent()
  const slider = page.getByLabel('x', { exact: true })
  await slider.focus()
  for (let i = 0; i < 10; i++) await slider.press('ArrowRight')
  await expect(readout).not.toHaveText(before ?? '')

  // The frame must not overflow the mobile viewport; a horizontal scrollbar on
  // the lesson page is the visible failure of a fixed-width SVG.
  const overflow = await page.evaluate(
    () => document.documentElement.scrollWidth - document.documentElement.clientWidth,
  )
  expect(overflow).toBeLessThanOrEqual(1)

  await page.screenshot({ path: `/tmp/opencode/visual-tangent-${testInfo.project.name}.png` })
})
