import { test, expect } from '@playwright/test'

const LESSON = {
  title: 'Addition Basics',
  body: '# Addition Basics\n\nLearn to add numbers.',
  concepts: ['arith.add.single'],
  progress: { 'arith.add.single': { status: 'UNSEEN', mastery: 0, streak: 0 } },
}

const PRACTICE = {
  questions: [
    { question: '2 + 2 = ?', answer: '4', explanation: '2+2=4', source: 'curated' },
    { question: '3 + 5 = ?', answer: '8', explanation: '3+5=8', source: 'curated' },
  ],
  concept_id: 'arith.add.single',
}

const KP = {
  concept_id: 'arith.add.single',
  kps: [
    { label: 'Use addition notation', section: 'Use Addition Notation', subgoals: ['Identify addends', 'Read 3+4', 'Translate'], worked_example: 'Add 3+4' },
  ],
  diagram: '/diagrams/algebrica/number-line-real.svg',
}

test('study → answer → XP persists (Study seam)', async ({ page }) => {
  await page.route('**/api/lessons**', route => {
    const url = route.request().url()
    if (url.includes('/practice')) return route.fulfill({ json: PRACTICE })
    if (url.includes('/kp')) return route.fulfill({ json: KP })
    if (url.includes('/body')) return route.fulfill({ json: { title: LESSON.title, body: LESSON.body } })
    return route.fulfill({ json: { lessons: { arithmetic: [LESSON] } } })
  })
  await page.route('**/api/study/answer', route =>
    route.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 10, expected_answer: '4', halted: false } }),
  )
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30, spaced_reps: {}, avg_learning_speed: 1.0 } }),
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: { 'arith.add.single': { status: 'LEARNING', streak: 1, mastery: 0.3 } } }))
  await page.route('**/api/activity**', route => route.fulfill({ json: [] }))
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 0 } }))
  await page.route('**/api/courses**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/transcript**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/efficacy**', route => route.fulfill({ json: { concepts_touched: 1, first_pass_rate: 1, second_pass_rate: 1, avg_attempts_per_concept: 1, total_attempts: 1 } }))

  await page.goto('/study?lesson=' + encodeURIComponent(LESSON.title))
  await expect(page.locator('body')).toContainText('Addition Basics', { timeout: 30_000 })
  await expect(page.getByText('Worked example').first()).toBeVisible()
  await expect(page.getByText('Use addition notation').first()).toBeVisible()

  const input = page.locator('input[placeholder*="Your answer"]').first()
  await expect(input).toBeVisible({ timeout: 20_000 })
  await input.fill('4')
  await page.getByRole('button', { name: 'Submit' }).first().click()
  await expect(page.getByText('+10 XP').first()).toBeVisible({ timeout: 20_000 })
})

test('quiz gate banner appears at 150 XP and links to quiz host', async ({ page }) => {
  const LESSONS = { arithmetic: [LESSON] }
  await page.addInitScript(() => localStorage.setItem('mathua_guest_id', 'guest_e2e_quiz'))
  await page.route('**/api/lessons**', route => {
    const url = route.request().url()
    if (url.includes('/practice') || url.includes('/kp')) return route.fulfill({ json: { questions: [], concept_id: 'arith.add.single', kps: [] } })
    if (url.includes('/body')) return route.fulfill({ json: { title: LESSON.title, body: LESSON.body } })
    return route.fulfill({ json: { lessons: LESSONS } })
  })
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 1500, weekly_score: 100, speed_bonus: 0, concepts_mastered: 15, current_streak: 5, level: 'Student', xp_total: 150, xp_today: 30, daily_xp_goal: 30, spaced_reps: {}, avg_learning_speed: 1.1 } }),
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
  await page.route('**/api/activity**', route => route.fulfill({ json: [] }))
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 0 } }))
  await page.route('**/api/courses**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/transcript**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/efficacy**', route => route.fulfill({ json: { concepts_touched: 10, first_pass_rate: 0.8, second_pass_rate: 1, avg_attempts_per_concept: 1.2, total_attempts: 12 } }))
  await page.route('**/api/leagues**', route => route.fulfill({ json: { week: '2026-W35', leagues: [] } }))

  await page.goto('/study')
  await expect(page.getByText('Quiz due').first()).toBeVisible({ timeout: 30_000 })
  await expect(page.getByRole('link', { name: 'Take Test' }).first()).toBeVisible()
})

test('quiz reuse host at /goals?quiz=1 starts actionable quiz (guest unlimited retake)', async ({ page }) => {
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 150, xp_today: 10, daily_xp_goal: 30 } }),
  )
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/quiz/session', route =>
    route.fulfill({ json: { session_id: 'q1', concept_id: 'arith.add.single', concept_name: 'Single-digit addition', question: '5 + 3 = ?', done: false } }),
  )
  await page.route('**/api/quiz/answer', route =>
    route.fulfill({ json: { done: false, correct: true, feedback: 'Correct!', xp: 20, concept_id: 'arith.sub.single', concept_name: 'Single-digit subtraction', question: '8 - 3 = ?' } }),
  )

  await page.goto('/goals?quiz=1')
  await expect(page.getByText('Quiz question 1').first()).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText('5 + 3 = ?').first()).toBeVisible()

  const input = page.locator('input[placeholder*="Your answer"]').first()
  await input.fill('8')
  await page.getByRole('button', { name: 'Check Answer' }).first().click()
  await expect(page.getByText('Correct').first()).toBeVisible({ timeout: 20_000 })
  await expect(page.getByText('+20 XP').first()).toBeVisible()
})

test('share link: settings enable → copyable URL → public share page', async ({ page }) => {
  await page.route('**/api/settings', route => {
    if (route.request().method() === 'GET') return route.fulfill({ json: {} })
    return route.continue()
  })
  await page.route('**/api/share', route => {
    if (route.request().method() === 'POST') return route.fulfill({ json: { enabled: true, token: 's_testtoken123' } })
    return route.continue()
  })
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 } }),
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
  await page.route('**/api/activity**', route => route.fulfill({ json: [] }))
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 0 } }))
  await page.route('**/api/courses**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/transcript**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/efficacy**', route => route.fulfill({ json: { concepts_touched: 0, first_pass_rate: 0, second_pass_rate: 0, avg_attempts_per_concept: 0, total_attempts: 0 } }))

  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 1, current_streak: 1, level: 'Novice', diagnostic_completed: true }))
  })
  await page.goto('/settings')
  await expect(page.getByText('Share with parent / teacher').first()).toBeVisible({ timeout: 20_000 })
  await page.getByRole('button', { name: 'Enable share link' }).click()
  const input = page.locator('input[value*="s_testtoken123"]').first()
  await expect(input).toBeVisible({ timeout: 20_000 })
  await expect(page.getByRole('button', { name: 'Copy link' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Disable share' })).toBeVisible()

  // Public share page
  await page.route('**/api/share/s_testtoken123', route =>
    route.fulfill({
      json: {
        student_id: 's1',
        name: 'Tester',
        scores: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 },
        activity: [],
        progress: {},
        weakness: {},
      },
    }),
  )
  await page.goto('/share?token=s_testtoken123')
  await expect(page.getByText("Tester's progress").first()).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText('Read-only report').first()).toBeVisible()
})
