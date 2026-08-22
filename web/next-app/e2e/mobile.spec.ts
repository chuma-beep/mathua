import { test, expect } from '@playwright/test'
import conceptsData from '../data/concepts.json'

const sample = conceptsData as Array<{ id: string; label: string; domain: string; prerequisites: string[] }>

test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

test('mobile graph renders and is usable at Pixel 7 viewport', async ({ page }) => {
  await page.goto('/graph')
  const nodes = page.locator('.react-flow__node')
  await expect(nodes.first()).toBeVisible({ timeout: 20_000 })
  expect(await nodes.count()).toBeGreaterThan(100)

  // no horizontal overflow
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
  expect(overflow).toBeLessThanOrEqual(1)

  // search overlay stays inside viewport and remains usable
  const input = page.getByPlaceholder('Search concepts…')
  await input.fill(sample[5].label.slice(0, 5))
  const box = await input.boundingBox()
  expect(box).toBeTruthy()
  expect(box!.x).toBeGreaterThanOrEqual(0)
  expect(box!.x + box!.width).toBeLessThanOrEqual(390)

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
  expect(box!.x).toBeGreaterThanOrEqual(0)
  expect(box!.x + box!.width).toBeLessThanOrEqual(391)
})

test('mobile header collapses to hamburger menu', async ({ page }) => {
  await page.goto('/graph')
  const header = page.locator('header')
  const burger = header.locator('button[aria-label="Toggle navigation menu"]')
  await expect(burger).toBeVisible()
  await burger.click()
  await expect(header.locator('nav.flex-col').getByText('Login')).toBeVisible()
})
