import { test, expect } from '@playwright/test'
import conceptsData from '../data/concepts.json'

test('landing page renders hero and lazy-mounts the 3D graph', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1 })).toBeVisible()
  const hero = page.locator('section').first()
  await hero.scrollIntoViewIfNeeded()
  await page.mouse.wheel(0, 400)
  await expect(page.locator('canvas').first()).toBeVisible({ timeout: 30_000 })
  // Domain cluster labels are painted by the imperative label overlay.
  await expect(page.locator('.graph-label--domain').first()).toBeVisible({ timeout: 30_000 })
})

test('landing page does not ship the concept corpus in First Load JS', async ({ page }) => {
  // Inspect the scripts in the exported HTML itself: prefetches for other
  // routes (e.g. /graph) may still pull the corpus at runtime, but it must
  // not be part of the landing page's first load.
  const html = await (await page.request.get('/')).text()
  const srcs = [...new Set([...html.matchAll(/src="([^"]+\.js)"/g)].map(m => m[1]))]
  expect(srcs.length).toBeGreaterThan(0)
  for (const src of srcs) {
    const body = await (await page.request.get(src)).text()
    // Corpus-only label; the build-generated graph payload is fetched, not bundled.
    expect(body).not.toContain('Subgroups and cosets')
  }
  // The fetched graph payload stays available for the scene.
  const payload = await page.request.get('/data/concepts.graph.json')
  expect(payload.ok()).toBe(true)
  const json = await payload.json()
  expect(json.nodes.length).toBeGreaterThan(500)
  expect(json.edges.length).toBeGreaterThan(500)
})

test('concept count is shown from the dataset', async ({ page }) => {
  await page.goto('/')
  const count = conceptsData.length
  await expect(page.getByText(String(count)).first()).toBeVisible()
})
