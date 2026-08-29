import { test, expect } from '@playwright/test'

test('graph controls are desktop-sized and minimap is readable at 1280px', async ({ page }) => {
  await page.goto('/graph')
  await expect(page.locator('.react-flow__node').first()).toBeVisible({ timeout: 20_000 })

  const ctrl = page.locator('.react-flow__controls-button').first()
  const ctrlBox = await ctrl.boundingBox()
  expect(ctrlBox).toBeTruthy()
  expect(ctrlBox!.width).toBeLessThanOrEqual(30)
  expect(ctrlBox!.height).toBeLessThanOrEqual(30)

  const mm = page.locator('.react-flow__minimap')
  const mmBox = await mm.boundingBox()
  expect(mmBox).toBeTruthy()
  expect(mmBox!.width).toBeGreaterThanOrEqual(140)
  expect(mmBox!.height).toBeGreaterThanOrEqual(90)

  // desktop header nav visible, bottom tabs hidden
  await expect(page.locator('header nav')).toBeVisible()
  await expect(page.locator('nav.lg\\:hidden')).not.toBeVisible()

  // no horizontal overflow
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
  expect(overflow).toBeLessThanOrEqual(1)
})

test('stacked table cards hidden on desktop, real table shown', async ({ page }) => {
  for (const path of ['/docs/contributing', '/how-it-works']) {
    await page.goto(path)
    await expect(page.locator('div.sm\\:hidden.space-y-2')).not.toBeVisible()
    await expect(page.locator('table')).toBeVisible()
  }
})
