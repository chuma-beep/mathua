import { test, expect } from '@playwright/test'

const USER = {
  student_id: 's1',
  name: 'Tester',
  username: 'tester',
  concepts_mastered: 3,
  current_streak: 2,
  level: 'beginner',
  diagnostic_completed: true,
}

test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

test('header shows Login for anonymous visitors', async ({ page }) => {
  await page.goto('/graph')
  await expect(page.locator('header').getByText('Login').first()).toBeVisible()
})

test('header shows authenticated links when a token exists', async ({ page }) => {
  await page.addInitScript(([user]) => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify(user))
  }, [USER] as unknown as string[])
  await page.goto('/graph')
  const header = page.locator('header')
  await expect(header.getByText('Profile').first()).toBeVisible()
  await expect(header.getByText('Start').first()).toBeVisible()
  await expect(header.getByText('Login')).toHaveCount(0)
})
