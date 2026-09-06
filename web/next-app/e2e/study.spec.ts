import { test, expect } from '@playwright/test'

const LESSON = {
  title: 'Addition Basics',
  body: '# Addition Basics\n\nLearn to add numbers.',
  concepts: ['addition_whole_numbers'],
}

test('?concept= resolves to the lesson containing that concept', async ({ page }) => {
  await page.route('**/api/lessons**', route =>
    route.fulfill({ json: { lessons: { arithmetic: [LESSON] } } })
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
  await page.goto(`/study?concept=${encodeURIComponent(LESSON.concepts[0])}`)
  await expect(page.locator('body')).toContainText(LESSON.title, { timeout: 30_000 })
})

test('study page lists domains normally without a concept param', async ({ page }) => {
  await page.route('**/api/lessons**', route =>
    route.fulfill({ json: { lessons: { arithmetic: [LESSON] } } })
  )
  await page.goto('/study')
  await expect(page.getByText('Arithmetic').first()).toBeVisible({ timeout: 30_000 })
})

test('domain list paginates at 20 per page', async ({ page }) => {
  const lessons: Record<string, typeof LESSON[]> = {}
  for (let i = 0; i < 25; i++) {
    const key = `zzdomain${String(i).padStart(2, '0')}`
    lessons[key] = [{ ...LESSON, title: `Lesson ${key}` }]
  }
  await page.route('**/api/lessons**', route => route.fulfill({ json: { lessons } }))
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
  await page.goto('/study')
  await expect(page.getByText('Page 1 of 2', { exact: true })).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText('zzdomain00').first()).toBeVisible()
  await expect(page.getByText('zzdomain24')).toHaveCount(0)
  await page.getByRole('button', { name: 'Next page' }).click()
  await expect(page.getByText('Page 2 of 2', { exact: true })).toBeVisible()
  await expect(page.getByText('zzdomain24').first()).toBeVisible()
  await expect(page.getByText('zzdomain00')).toHaveCount(0)
  await page.getByRole('button', { name: 'Previous page' }).click()
  await expect(page.getByText('Page 1 of 2', { exact: true })).toBeVisible()
})
