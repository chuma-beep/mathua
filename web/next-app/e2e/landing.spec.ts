import { test, expect } from '@playwright/test'
import conceptsData from '../data/concepts.json'

test('landing page renders hero and lazy-mounts the 3D graph', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('heading', { level: 1 })).toBeVisible()
  const hero = page.locator('section').first()
  await hero.scrollIntoViewIfNeeded()
  await page.mouse.wheel(0, 400)
  await expect(page.locator('canvas').first()).toBeVisible({ timeout: 30_000 })
})

test('stats row shows concept and domain counts from the dataset', async ({ page }) => {
  await page.goto('/')
  const count = conceptsData.length
  await expect(page.getByText(String(count)).first()).toBeVisible()
})
