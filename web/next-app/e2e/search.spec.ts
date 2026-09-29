import { test, expect } from '@playwright/test'

// SearchBar: exactly one clear control, and the results list scrolls.
test('study search has one clear button and scrollable results', async ({ page }) => {
  const lessons = Array.from({ length: 15 }, (_, i) => ({
    title: `Fraction Lesson ${i + 1}`,
    body: `# Fraction Lesson ${i + 1}\n\n${'Adding fractions with common denominators. '.repeat(8)}`,
    concepts: [`frac.add.${i + 1}`],
  }))
  await page.route('**/api/lessons', route => route.fulfill({ json: { lessons: { fractions: lessons } } }))
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 0, weekly_score: 0, speed_bonus: 0, concepts_mastered: 0, current_streak: 0, level: 'Novice', xp_total: 0, xp_today: 0, daily_xp_goal: 10, spaced_reps: {}, avg_learning_speed: 1.0 } }),
  )

  await page.goto('/study')
  const input = page.getByPlaceholder('Search lessons…')
  await expect(input).toBeVisible({ timeout: 30_000 })
  await input.fill('fraction')

  // Exactly one clear control: no native search-cancel plus custom X.
  await expect(page.getByRole('button', { name: 'Clear search' })).toHaveCount(1)

  const list = page.getByRole('listbox', { name: 'Search results' })
  await expect(list).toBeVisible()
  const scrollable = await list.evaluate(el => el.scrollHeight > el.clientHeight)
  expect(scrollable).toBe(true)
  const before = await list.evaluate(el => el.scrollTop)
  await list.hover()
  await page.mouse.wheel(0, 400)
  await expect.poll(() => list.evaluate(el => el.scrollTop), { timeout: 5_000 }).not.toBe(before)

  // Escape dismisses.
  await page.keyboard.press('Escape')
  await expect(list).toHaveCount(0)
})
