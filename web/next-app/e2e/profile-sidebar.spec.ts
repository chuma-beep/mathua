import { test, expect, type Page } from '@playwright/test'

test.beforeEach(async ({ context }) => {
  await context.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 2, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 } }),
  )
  await context.route('**/api/progress/**', route => route.fulfill({ json: {} }))
  await context.route('**/api/activity**', route => route.fulfill({ json: [] }))
  await context.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await context.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 3 } }))
  await context.route('**/api/efficacy**', route => route.fulfill({ json: { concepts_touched: 0, first_pass_rate: 0, second_pass_rate: 0, avg_attempts_per_concept: 0, total_attempts: 0 } }))
  await context.route('**/api/settings**', route => route.fulfill({ json: {} }))
})

async function loginAs(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 1, current_streak: 2, level: 'Novice', diagnostic_completed: true }))
  })
}

test('profile sidebar renders nav without duplicating main CTAs', async ({ page }) => {
  await loginAs(page)
  await page.goto('/profile')
  const sidebar = page.getByTestId('profile-sidebar')
  await expect(sidebar).toBeVisible({ timeout: 30_000 })
  for (const label of ['Study', 'Start', 'Graph', 'Leaderboard', 'Settings']) {
    await expect(sidebar.getByText(label, { exact: true }).first()).toBeVisible()
  }
  // No CTA duplication: Diagnostic / Quiz / Review live in main-column cards only
  await expect(sidebar.getByText('Diagnostic', { exact: true })).toHaveCount(0)
  await expect(sidebar.getByText('Take Test')).toHaveCount(0)
  // Removed groups never render in sidebar: Goals, On this page
  await expect(sidebar.getByText('Goals', { exact: true })).toHaveCount(0)
  await expect(sidebar.getByText('On this page', { exact: true })).toHaveCount(0)
  // Due-review count surfaces as a badge, not a second CTA
  await expect(sidebar.locator('[data-sidebar="menu-badge"]')).toHaveText('3')
  // Main column keeps its own cards (test user already completed Diagnostic → Retake)
  await expect(page.getByText(/diagnostic test →/i).first()).toBeVisible()
})

test('profile sidebar collapses to icons via trigger', async ({ page }) => {
  await loginAs(page)
  await page.goto('/profile')
  const sidebar = page.getByTestId('profile-sidebar')
  await expect(sidebar).toBeVisible({ timeout: 30_000 })
  const rail = sidebar.locator('xpath=ancestor::div[@data-state][1]')
  await expect(rail).toHaveAttribute('data-state', 'expanded')
  const expandedWidth = await sidebar.boundingBox().then((b) => b!.width)
  const collapseToggle = page.getByRole('button', { name: 'Collapse sidebar', exact: true })
  await collapseToggle.click()
  // Collapsed: state flips, rail narrows to icon width, group labels hide
  await expect(rail).toHaveAttribute('data-state', 'collapsed', { timeout: 10_000 })
  await expect.poll(async () => sidebar.boundingBox().then((b) => b!.width), { timeout: 10_000 }).toBeLessThan(expandedWidth)
  // Group labels fade out in icon mode (opacity transition, not display:none)
  await expect
    .poll(
      async () =>
        sidebar
          .getByText('Navigate', { exact: true })
          .evaluate((el) => getComputedStyle(el).opacity),
      { timeout: 10_000 }
    )
    .toBe('0')
  // Collapsed rail: nav links each show a centered Lucide icon (no text overflow)
  for (const href of ['/study', '/session', '/graph', '/leaderboard', '/settings']) {
    const link = sidebar.locator(`a[href="${href}"]`).first()
    await expect(link.locator('svg').first()).toBeVisible()
  }
  // Collapsed rail shows icons only — label spans are display:none in icon mode
  for (const label of ['Mathua', 'Study', 'Start', 'Graph', 'Leaderboard', 'Settings', 'Sign out']) {
    await expect(sidebar.locator(`span:text-is("${label}")`).first()).toBeHidden()
  }
  await page.getByRole('button', { name: 'Expand sidebar', exact: true }).click()
  await expect(rail).toHaveAttribute('data-state', 'expanded', { timeout: 10_000 })
  await expect(sidebar.getByText('Leaderboard', { exact: true }).first()).toBeVisible()
})

test('profile sidebar collapses to icons only at large viewport', async ({ page }) => {
  await page.setViewportSize({ width: 1920, height: 1080 })
  await loginAs(page)
  await page.goto('/profile')
  const sidebar = page.getByTestId('profile-sidebar')
  await expect(sidebar).toBeVisible({ timeout: 30_000 })
  const rail = sidebar.locator('xpath=ancestor::div[@data-state][1]')
  await expect(rail).toHaveAttribute('data-state', 'expanded')
  const expandedWidth = await sidebar.boundingBox().then((b) => b!.width)
  await page.getByRole('button', { name: 'Collapse sidebar', exact: true }).click()
  await expect(rail).toHaveAttribute('data-state', 'collapsed', { timeout: 10_000 })
  await expect.poll(async () => sidebar.boundingBox().then((b) => b!.width), { timeout: 10_000 }).toBeLessThan(expandedWidth)
  // Icons only: every nav link shows its svg, every label is hidden
  for (const href of ['/study', '/session', '/graph', '/leaderboard', '/settings']) {
    const link = sidebar.locator(`a[href="${href}"]`).first()
    await expect(link.locator('svg').first()).toBeVisible()
  }
  for (const label of ['Mathua', 'Study', 'Start', 'Graph', 'Leaderboard', 'Settings', 'Sign out']) {
    await expect(sidebar.locator(`span:text-is("${label}")`).first()).toBeHidden()
  }
})
