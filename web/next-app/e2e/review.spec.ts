import { test, expect } from '@playwright/test'

// Review runner flow (intro → answer → done). Safety net for the dead-end
// fix: both Review-now links land here instead of the session chooser.
test('review: start → answer → done', async ({ page }) => {
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 } }),
  )
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 2 } }))
  await page.route('**/api/reviews/session', route =>
    route.fulfill({ json: { student_id: 's1', session_id: 'r1', question: { concept_id: 'arith.add.single', concept_name: 'Single-digit addition', question: '2 + 3 = ?', is_review: true, attempt_id: 'a1' } } }),
  )
  await page.route('**/api/reviews/answer', route =>
    route.fulfill({ json: { result: { correct: true, feedback: 'Correct!', new_status: 'PRACTICING', explanation: '', streak: 1, required_streak: 2, xp: 5, expected_answer: '5' }, next_question: null, done: true } }),
  )

  await page.goto('/review')
  await expect(page.getByRole('button', { name: 'Start review →' })).toBeVisible({ timeout: 30_000 })
  await page.getByRole('button', { name: 'Start review →' }).click()

  await expect(page.getByText('2 + 3 = ?').first()).toBeVisible({ timeout: 20_000 })
  const input = page.locator('input[placeholder*="Your answer"]').first()
  await input.fill('5')
  await page.getByRole('button', { name: 'Check Answer' }).first().click()

  await expect(page.getByText('Review complete').first()).toBeVisible({ timeout: 20_000 })
  await expect(page.getByText('1/1 correct').first()).toBeVisible()
})

test('review: nothing due shows the caught-up state', async ({ page }) => {
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 } }),
  )
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 0 } }))
  await page.route('**/api/reviews/session', route =>
    route.fulfill({ json: { student_id: 's1', session_id: 'r1', question: null } }),
  )

  await page.goto('/review')
  await expect(page.getByRole('button', { name: 'Start review →' })).toBeVisible({ timeout: 30_000 })
  await page.getByRole('button', { name: 'Start review →' }).click()
  await expect(page.getByText('All caught up').first()).toBeVisible({ timeout: 20_000 })
})
