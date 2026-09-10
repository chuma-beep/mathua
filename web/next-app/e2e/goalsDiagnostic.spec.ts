import { test, expect } from '@playwright/test'

// Goals diagnostic flow (select → answer → results). This is the safety net
// for the DiagnosticHost extraction: it must keep passing unchanged.
const PROGRESS = {
  answered: 0,
  estimated_total: 25,
  min_total: 25,
  max_total: 45,
  cover_done: 0,
  cover_size: 5,
  done: false,
}

const PLAN = {
  readiness: 0.62,
  total_tested: 12,
  correct_count: 9,
  weak_areas: { arithmetic: [{ id: 'arith.div.long', label: 'Long division' }] },
  strong_areas: { fractions: ['frac.add.same'] },
  frontier_label: 'Linear equations',
  frontier_idx: 3,
  frontier_conditional: false,
  conditionally_completed: [],
  placement_course_id: 'a1',
  completion_estimates: { 150: '≈5 days', 300: '≈10 days' },
}

test('goals diagnostic: select domain → answer → results', async ({ page }) => {
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 0, weekly_score: 0, speed_bonus: 0, concepts_mastered: 0, current_streak: 0, level: 'Novice', xp_total: 0, xp_today: 0, daily_xp_goal: 30 } }),
  )
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/goal/diagnostic', route =>
    route.fulfill({ json: { session_id: 's1', concept_id: 'arith.add.single', concept_name: 'Single-digit addition', question: '2 + 3 = ?', done: false, progress: PROGRESS } }),
  )
  await page.route('**/api/goal/diagnostic/answer', route =>
    route.fulfill({ json: { done: true, correct: true, feedback: 'Correct!', progress: { ...PROGRESS, answered: 1 } } }),
  )
  await page.route('**/api/goal/plan**', route => route.fulfill({ json: PLAN }))

  await page.goto('/goals')
  await expect(page.getByRole('button', { name: /Arithmetic/ }).first()).toBeVisible({ timeout: 20_000 })

  await page.getByRole('button', { name: /Arithmetic/ }).first().click()
  await page.getByRole('button', { name: /Start diagnostic test/ }).click()

  await expect(page.getByText('2 + 3 = ?')).toBeVisible({ timeout: 20_000 })
  const input = page.locator('input[placeholder*="Your answer"]').first()
  await input.fill('5')
  await page.getByRole('button', { name: 'Check Answer' }).first().click()

  await expect(page.getByText("Here's what we found")).toBeVisible({ timeout: 20_000 })
  await expect(page.getByText('Linear equations')).toBeVisible()
})
