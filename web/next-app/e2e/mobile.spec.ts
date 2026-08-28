import { test, expect, type Page } from '@playwright/test'
import conceptsData from '../data/concepts.json'

const sample = conceptsData as Array<{ id: string; label: string; domain: string; prerequisites: string[] }>

test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

async function overflowPx(page: Page): Promise<number> {
  return page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
}

const ROUTES = [
  '/',
  '/profile',
  '/study',
  '/session',
  '/concept?id=arith.add.single',
  '/graph',
  '/leaderboard',
  '/login',
  '/onboard',
  '/settings',
  '/goals',
  '/diagnose',
  '/how-it-works',
  '/docs',
  '/docs/architecture',
  '/docs/contributing',
]

for (const route of ROUTES) {
  test(`no horizontal overflow on ${route}`, async ({ page }) => {
    await page.goto(route)
    await page.waitForTimeout(2500)
    const overflow = await overflowPx(page)
    expect(overflow).toBeLessThanOrEqual(1)
  })
}

test('bottom tabs present and scroll-aware on mobile', async ({ page }) => {
  await page.goto('/how-it-works')
  const tabs = page.locator('nav.lg\\:hidden')
  await expect(tabs).toBeVisible()

  // Disable smooth scrolling so the scroll settles synchronously for the hook.
  await page.evaluate(() => { document.documentElement.style.scrollBehavior = 'auto' })

  // Scroll down mid-page (not the bottom — bottom force-shows) → hidden while scrolling
  await page.evaluate(() => window.scrollTo(0, Math.min(800, document.documentElement.scrollHeight - window.innerHeight - 300)))
  await page.waitForTimeout(150)
  await expect(tabs).not.toBeInViewport()

  // Stopping scroll re-shows after the 300ms idle window
  await page.waitForTimeout(600)
  await expect(tabs).toBeInViewport()

  // Scroll back to top → still visible
  await page.evaluate(() => window.scrollTo(0, 0))
  await page.waitForTimeout(600)
  await expect(tabs).toBeInViewport()
})

test('mobile graph renders and is usable', async ({ page }) => {
  await page.goto('/graph')
  const nodes = page.locator('.react-flow__node')
  await expect(nodes.first()).toBeVisible({ timeout: 20_000 })
  expect(await nodes.count()).toBeGreaterThan(100)

  expect(await overflowPx(page)).toBeLessThanOrEqual(1)

  // search overlay stays inside viewport and remains usable
  const input = page.getByPlaceholder('Search concepts…')
  await input.fill(sample[5].label.slice(0, 5))
  const box = await input.boundingBox()
  expect(box).toBeTruthy()
  const vw = page.viewportSize()!.width
  expect(box!.x).toBeGreaterThanOrEqual(0)
  expect(box!.x + box!.width).toBeLessThanOrEqual(vw)
  expect(box!.height).toBeGreaterThanOrEqual(36)

  // ambient toggle must NOT be offered on mobile (forced off there)
  await expect(page.getByTitle('Toggle edge flow animation')).toHaveCount(0)
})

test('mobile pinch-zoom works on the concept map', async ({ page }) => {
  await page.goto('/graph')
  await expect(page.locator('.react-flow__node').first()).toBeVisible({ timeout: 20_000 })

  const pane = page.locator('.react-flow__viewport')
  const before = await pane.evaluate(el => el.style.transform)
  await page.locator('.react-flow').hover({ position: { x: 180, y: 160 } })
  // ReactFlow maps ctrl+wheel to zoom on touch devices too; simulate pinch via wheel
  await page.mouse.wheel(0, -240)
  await page.waitForTimeout(600)
  const after = await pane.evaluate(el => el.style.transform)
  expect(after).not.toBe(before)
})

test('mobile info panel wraps without overflowing', async ({ page }) => {
  const target = sample.find(c => c.prerequisites.length > 2) ?? sample[10]
  await page.goto(`/graph?concept=${encodeURIComponent(target.id)}`)
  const openBtn = page.getByText('Open concept →')
  await expect(openBtn).toBeVisible({ timeout: 20_000 })
  const box = await openBtn.boundingBox()
  expect(box).toBeTruthy()
  const vw = page.viewportSize()!.width
  expect(box!.x).toBeGreaterThanOrEqual(0)
  expect(box!.x + box!.width).toBeLessThanOrEqual(vw)
})

test('graph list toggle and search buttons meet 36px hit area', async ({ page }) => {
  await page.goto('/graph')
  const listBtn = page.locator('button[aria-label="Open concept list"]')
  await expect(listBtn).toBeVisible({ timeout: 20_000 })
  const listBox = await listBtn.boundingBox()
  expect(listBox).toBeTruthy()
  expect(listBox!.height).toBeGreaterThanOrEqual(36)
})
