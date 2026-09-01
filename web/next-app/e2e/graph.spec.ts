import { test, expect } from '@playwright/test'
import conceptsData from '../data/concepts.json'

const sample = conceptsData.slice(0, 300) as Array<{ id: string; label: string; domain: string; prerequisites: string[] }>

test('graph renders the full fallback concept set', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  await page.goto('/graph')
  const nodes = page.locator('.react-flow__node')
  await expect(nodes.first()).toBeVisible({ timeout: 30_000 })
  const nodeCount = await nodes.count()
  expect(nodeCount).toBeGreaterThan(100)
})

test('search selects a concept and syncs the URL param', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  await page.goto('/graph')
  const input = page.getByPlaceholder('Search concepts…')
  await expect(input).toBeVisible({ timeout: 30_000 })
  const target = sample.find(c => c.prerequisites.length > 0)!
  await input.fill(target.label.slice(0, Math.max(4, target.label.length - 2)))
  await page.locator('button', { hasText: target.label }).first().click()
  await expect(page).toHaveURL(new RegExp(`concept=${encodeURIComponent(target.id)}`))
  await expect(page.getByText('Open concept →')).toBeVisible()
})

test('deep-link restores selection via ?concept=', async ({ page }) => {
  await page.route('**/api/**', route => route.fulfill({ status: 404, body: 'not found' }))
  const target = sample.find(c => c.prerequisites.length > 0)!
  await page.goto(`/graph?concept=${encodeURIComponent(target.id)}`)
  await expect(page.getByText(target.label).first()).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText('Open concept →')).toBeVisible()
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
