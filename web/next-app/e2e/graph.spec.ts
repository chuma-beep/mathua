import { test, expect } from '@playwright/test'
import conceptsData from '../data/concepts.json'

const sample = conceptsData.slice(0, 300) as Array<{ id: string; label: string; domain: string; prerequisites: string[] }>

test('graph renders the full fallback concept set', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  await page.goto('/graph')
  const nodes = page.locator('.react-flow__node')
  await expect(nodes.first()).toBeVisible({ timeout: 45_000 })
  await expect.poll(async () => nodes.count(), { timeout: 45_000 }).toBeGreaterThan(100)
})

test('search selects a concept and syncs the URL param', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  await page.goto('/graph')
  const input = page.locator('input[placeholder="Search concepts…"]')
  await expect.poll(async () => await input.isVisible(), { timeout: 30_000 }).toBe(true)
  const target = sample.find(c => c.prerequisites.length > 0)!
  await input.fill(target.label.slice(0, Math.max(4, target.label.length - 2)))
  await page.locator('button', { hasText: target.label }).first().click()
  await expect.poll(async () => page.url().includes(`concept=${encodeURIComponent(target.id)}`), { timeout: 30_000 }).toBe(true)
  await expect(page.getByText('Open concept →')).toBeVisible({ timeout: 10_000 })
})

test('deep-link restores selection via ?concept=', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  const target = sample.find(c => c.prerequisites.length > 0)!
  await page.goto(`/graph?concept=${encodeURIComponent(target.id)}`)
  await expect.poll(async () => await page.getByText(target.label).first().isVisible().catch(() => false), { timeout: 30_000 }).toBe(true)
  await expect(page.getByText('Open concept →')).toBeVisible({ timeout: 10_000 })
})

test('ambient flow toggle persists across reloads', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  await page.goto('/graph')
  await expect(page.locator('.react-flow__node').first()).toBeVisible({ timeout: 30_000 })
  await page.getByTitle('Toggle edge flow animation').click()
  await expect.poll(() => page.evaluate(() => localStorage.getItem('mathua_graph_flow'))).toBe('1')
  await page.reload()
  await expect(page.evaluate(() => localStorage.getItem('mathua_graph_flow'))).resolves.toBe('1')
})
