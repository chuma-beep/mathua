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
  await expect(page.locator('body')).toContainText(LESSON.title, { timeout: 15_000 })
})

test('study page lists domains normally without a concept param', async ({ page }) => {
  await page.route('**/api/lessons**', route =>
    route.fulfill({ json: { lessons: { arithmetic: [LESSON] } } })
  )
  await page.goto('/study')
  await expect(page.getByText('Arithmetic').first()).toBeVisible({ timeout: 15_000 })
})
