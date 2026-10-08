import { test, expect } from '@playwright/test'
import { answerField, typeAnswer } from './helpers/answer'

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

// The whole loop, on the one surface that now has it.
//
// This used to be a two-page test: open `/study?lesson=…`, read the worked examples, then click
// "Start learning →" into `/learn` to be asked a question. That pairing was the bug, not the
// test — a learner could read a concept instead of being taught it, and the app made the
// reading feel like progress. Both halves now live on one page, so the test does too.
async function stubLearn(page: import('@playwright/test').Page) {
  await page.route('**/api/lessons**', route => {
    const url = route.request().url()
    if (url.includes('/practice')) return route.fulfill({ json: PRACTICE })
    if (url.includes('/kp')) return route.fulfill({ json: KP })
    if (url.includes('/body')) return route.fulfill({ json: { title: LESSON.title, body: LESSON.body } })
    if (url.includes('/readiness')) return route.fulfill({ json: { concept_id: 'arith.add.single', ready: true, weak: [], missing: [] } })
    return route.fulfill({ json: { lessons: { arithmetic: [LESSON] } } })
  })
  await page.route('**/api/study/answer', route =>
    route.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 1, expected_answer: '4', halted: false } }),
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
}

test('learn teaches, answers and awards XP on one page', async ({ page }) => {
  await stubLearn(page)

  await page.goto('/learn?concept=arith.add.single')
  // Scoped to the intro's own KP label. The reference panel repeats the same knowledge-point
  // labels inside a collapsed `<details>`, and Playwright treats a `<summary>` in a closed
  // `<details>` as hidden — so a page-wide `.first()` resolves to the panel's copy and fails for a
  // reason unrelated to what this test is about.
  const kpLabel = page.locator('p.font-mono.uppercase', { hasText: 'Use addition notation' })
  await expect(kpLabel.first()).toBeVisible({ timeout: 30_000 })
  await page.getByRole('button', { name: 'Next →' }).click()

  // The answer control is MathLive for arithmetic concepts and the plain input for
  // choice-based ones; `answerField` finds whichever rendered.
  const input = answerField(page)
  await expect(input).toBeVisible({ timeout: 20_000 })
  // Answer input is at least the standard 48px control at all widths.
  const inputBox = await input.boundingBox()
  expect(inputBox?.height).toBeGreaterThanOrEqual(48)
  await typeAnswer(page, '4')
  await page.getByRole('button', { name: 'Check', exact: true }).first().click()
  await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
})

test('the reference material is on the learn page, and holds no way out of it', async ({ page }) => {
  await stubLearn(page)
  await page.goto('/learn?concept=arith.add.single')

  // The worked example the library used to render as a whole article is here, behind a
  // disclosure, without leaving the page that asks the question.
  const panel = page.getByText('Reference — the full lesson')
  await expect(panel).toBeVisible({ timeout: 30_000 })
  await panel.click()
  await expect(page.getByText('Add 3+4').first()).toBeVisible()

  // The behavioural property, asserted in a real browser rather than by reading source: nothing
  // inside the panel navigates. A link out of here — to another concept, to a library, anywhere
  // — is the second learning path this closure exists to prevent.
  const linksInPanel = await page.locator('details:has-text("Reference — the full lesson") a').all()
  expect(linksInPanel).toHaveLength(0)

  // And it says so itself: reading is not progress.
  await expect(page.getByText(/only the questions below move this concept/)).toBeVisible()
})

test('/study forwards a concept deep link into the loop', async ({ page }) => {
  // Old bookmarks and inbound links exist, and a 404 would lose the one thing the old URL
  // still carried: which concept. Static export cannot use next.config redirects, so this is a
  // client forwarder and the URL changes in the browser rather than over the wire.
  await stubLearn(page)
  await page.goto('/study?concept=arith.add.single')
  await expect(page).toHaveURL(/\/learn\?concept=arith\.add\.single/, { timeout: 30_000 })
  // Not a page-wide `.first()`. The reference panel repeats the same knowledge-point labels
  // inside a collapsed `<details>`, and Playwright treats a `<summary>` in a closed `<details>`
  // as hidden — so a page-wide `.first()` can resolve to the panel's copy and fail for a reason
  // that has nothing to do with the forwarder. Scoped to the intro's own label, which is what
  // "the concept reached the teaching loop" actually looks like.
  const kpLabel = page.locator('p.font-mono.uppercase', { hasText: 'Use addition notation' })
  await expect(kpLabel.first()).toBeVisible()
})

test('/study with no concept forwards to the learn entry', async ({ page }) => {
  await stubLearn(page)
  await page.route('**/api/next', route =>
    route.fulfill({
      json: {
        primary: {
          id: 'arith.add.single', conceptId: 'arith.add.single', conceptTitle: 'Single-digit addition',
          kind: 'learn', reason: 'new', priority: 1,
          action: { type: 'learn', href: '/learn?concept=arith.add.single' },
          badge: 'New', detail: 'ready to learn', cta: 'Start →',
        },
        alternatives: [],
        generatedAt: '',
      },
    }),
  )
  await page.goto('/study')
  await expect(page).toHaveURL(/\/learn(\?|$)/, { timeout: 30_000 })
})

test('no learner navigation offers a way into the closed library', async ({ page }) => {
  // Enumerated over every surface the library used to appear in, rather than one of them. Each
  // was a separate place to reintroduce it, and a single-page assertion would pass while four
  // of the other six still linked there.
  await page.route('**/api/**', route => route.fulfill({ status: 404, json: {} }))
  const surfaces: [string, string][] = [
    ['/', 'landing'],
    ['/learn', 'learn entry'],
    ['/profile', 'profile'],
    ['/graph', 'graph'],
    ['/domains', 'domains'],
    ['/review', 'review'],
    ['/leaderboard', 'leaderboard'],
    ['/history', 'history'],
    ['/settings', 'settings'],
    ['/how-it-works', 'how it works'],
    ['/docs', 'docs'],
    ['/note', 'note'],
  ]
  for (const [path, name] of surfaces) {
    await page.goto(path)
    await page.waitForLoadState('domcontentloaded')
    const chrome = page.locator('header, nav, footer')
    const hrefs = await chrome.locator('a').evaluateAll(els => els.map((e) => e.getAttribute('href') ?? ''))
    expect(hrefs.filter((h) => h.startsWith('/study')), `${name} links to /study`).toEqual([])
    const words = await chrome.locator('a, button').evaluateAll(els => els.map((e) => (e.textContent ?? '').trim()))
    expect(words.filter((w) => /^(open |browse )?study$/i.test(w)), `${name} offers a Study control`).toEqual([])
  }
})

test('quiz gate banner appears at 50 XP on profile and links to quiz host', async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 15, current_streak: 5, level: 'Student', diagnostic_completed: true }))
  })
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 1500, weekly_score: 100, speed_bonus: 0, concepts_mastered: 15, current_streak: 5, level: 'Student', xp_total: 50, xp_today: 10, daily_xp_goal: 10, quiz_due: true, xp_since_quiz: 50, quiz_gate_xp: 50, spaced_reps: {}, avg_learning_speed: 1.1 } }),
  )
  await page.route('**/api/progress/**', route => route.fulfill({ json: {} }))
  await page.route('**/api/activity**', route => route.fulfill({ json: [] }))
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', route => route.fulfill({ json: { count: 0 } }))
  await page.route('**/api/courses**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/transcript**', route => route.fulfill({ json: { courses: [] } }))
  await page.route('**/api/efficacy**', route => route.fulfill({ json: { concepts_touched: 10, first_pass_rate: 0.8, second_pass_rate: 1, avg_attempts_per_concept: 1.2, total_attempts: 12 } }))
  await page.route('**/api/efficacy/trend**', route => route.fulfill({ json: { weeks: [], total_students: 0, returning_students: 0, retention_rate: 0, first_pass_trend: 0 } }))
  await page.route('**/api/leagues**', route => route.fulfill({ json: { week: '2026-W35', leagues: [] } }))
  await page.route('**/api/settings**', route => route.fulfill({ json: {} }))

  await page.goto('/profile')
  await expect(page.getByText('Quiz due').first()).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText('50 XP reached: mastery check recommended').first()).toBeVisible()
  await expect(page.getByRole('link', { name: 'Start the quiz' }).first()).toBeVisible()
})

test('quiz reuse host at /goals?quiz=1 starts actionable quiz (guest unlimited retake)', async ({ page }) => {
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 150, xp_today: 10, daily_xp_goal: 30 } }),
  )
  await page.route('**/api/weaknesses**', route => route.fulfill({ json: { by_domain: {} } }))
  // time_limit_seconds is deliberately present. The quiz used to render a live countdown from
  // it — "30s", ticking down, then "Time up (counts as slow)" — and the field is still sent by
  // the server because TimeLimitFor/accommodatedThreshold scale it for grading. Supplying it here
  // is what makes this test able to catch the countdown coming back: with the field absent the
  // old UI rendered nothing either, so the assertion below would pass either way.
  await page.route('**/api/quiz/session', route =>
    route.fulfill({ json: { session_id: 'q1', concept_id: 'arith.add.single', concept_name: 'Single-digit addition', question: '5 + 3 = ?', grading_type: 'numeric', time_limit_seconds: 30, done: false } }),
  )
  let sawDontKnow = false
  await page.route('**/api/quiz/answer', route => {
    if (route.request().postDataJSON()?.dont_know === true) sawDontKnow = true
    return route.fulfill({ json: { done: false, correct: true, feedback: 'Correct!', xp: 3, concept_id: 'arith.sub.single', concept_name: 'Single-digit subtraction', question: '8 - 3 = ?', grading_type: 'numeric' } })
  })

  await page.goto('/goals?quiz=1')
  await expect(page.getByText('Before you start').first()).toBeVisible({ timeout: 30_000 })
  await page.getByRole('button', { name: 'Start quiz →' }).click()
  await expect(page.getByText('Quiz question 1').first()).toBeVisible({ timeout: 30_000 })
  await expect(page.getByText('5 + 3 = ?').first()).toBeVisible()
  await expect(page.getByText('Answer with a number').first()).toBeVisible()

  // Time is measured, not shown. A countdown and a "Time up" warning made the clock the task:
  // it rewarded speed over reasoning and put a learner who thinks for 30 seconds into a visibly
  // failing state. elapsed_seconds still goes to the server, so slow classification, XP and
  // analytics are unaffected — the learner is simply not told about any of it while answering.
  const quizChrome = await page.locator('body').innerText()
  expect(quizChrome).not.toMatch(/Time up/)
  expect(quizChrome).not.toMatch(/\b\d+s\b/)

  await typeAnswer(page, '8')
  await page.getByRole('button', { name: 'Check Answer' }).first().click()
  await expect(page.getByText('Correct').first()).toBeVisible({ timeout: 20_000 })
  await expect(page.getByText('+3 XP').first()).toBeVisible()
  // Manual advance: feedback holds until Next reveals the staged question.
  await page.getByRole('button', { name: 'Next →' }).click()
  await expect(page.getByText('8 - 3 = ?').first()).toBeVisible({ timeout: 20_000 })
  // I don't know: empty input still submits, flagged for the engine.
  await page.getByRole('button', { name: "I don't know" }).click()
  await expect(page.getByText('Correct').first()).toBeVisible({ timeout: 20_000 })
  await expect.poll(async () => sawDontKnow, { timeout: 20_000 }).toBe(true)
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

test('history: rows with corrections, missed filter, source filter', async ({ page }) => {
  const rows = [
    { session_id: 's1', student_id: 's1', concept_id: 'arith.add.single', concept_name: 'Single-digit addition', answer: '5', expected: '4', correct: false, elapsed_seconds: 3, timestamp: '2026-09-10T10:00:00Z', question: '2 + 3 = ?', source: 'diagnostic', explanation: '2 + 3 = 5' },
    { session_id: 's1', student_id: 's1', concept_id: 'arith.add.single', concept_name: 'Single-digit addition', answer: '5', expected: '5', correct: true, elapsed_seconds: 2, timestamp: '2026-09-11T10:00:00Z', question: '2 + 3 = ?', source: 'quiz', explanation: '' },
  ]
  await page.route('**/api/attempts**', route => {
    const url = new URL(route.request().url())
    let out = rows
    if (url.searchParams.get('incorrect_only') === '1') out = out.filter(r => !r.correct)
    const source = url.searchParams.get('source')
    if (source) out = out.filter(r => r.source === source)
    return route.fulfill({ json: { attempts: out, total: out.length } })
  })
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 1, current_streak: 1, level: 'Novice', diagnostic_completed: true }))
  })
  await page.goto('/history')
  await expect(page.getByText("Every question you've answered").first()).toBeVisible({ timeout: 30_000 })
  // Default: missed only.
  await expect(page.getByText('Showing 1 of 1').first()).toBeVisible()
  await expect(page.getByText('2 + 3 = ?').first()).toBeVisible()
  await expect(page.getByText(/Correct: 4/).first()).toBeVisible()
  // Switch to all + filter by quiz source.
  await page.getByRole('button', { name: 'All', exact: true }).click()
  await expect(page.getByText('Showing 2 of 2').first()).toBeVisible()
  await page.getByRole('button', { name: 'quiz' }).click()
  await expect(page.getByText('Showing 1 of 1').first()).toBeVisible()
})
