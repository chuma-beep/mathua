import { test, expect, type Page } from '@playwright/test'
import conceptsData from '../data/concepts.json'

const sample = conceptsData as Array<{ id: string; label: string; domain: string; prerequisites: string[] }>

test.beforeEach(async ({ context }) => {
  await context.route('**/api/**', route => route.fulfill({ status: 404, body: '' }))
})

async function overflowPx(page: Page): Promise<number> {
  return page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
}

const ROUTES = [
  '/',
  '/profile',
  '/study',
  '/session',
  '/concept?id=arith.add.single',
  '/graph',
  '/leaderboard',
  '/login',
  '/onboard',
  '/settings',
  '/goals',
  '/diagnose',
  '/how-it-works',
  '/docs',
  '/docs/architecture',
  '/docs/contributing',
]

for (const route of ROUTES) {
  test(`no horizontal overflow on ${route}`, async ({ page }) => {
    await page.goto(route)
    await page.waitForTimeout(2500)
    const overflow = await overflowPx(page)
    expect(overflow).toBeLessThanOrEqual(1)
  })
}

test('bottom tabs present and scroll-aware on mobile', async ({ page }) => {
  await page.goto('/how-it-works')
  const tabs = page.locator('nav.lg\\:hidden')
  await expect(tabs).toBeVisible()

  await page.evaluate(() => { document.documentElement.style.scrollBehavior = 'auto' })

  // Scroll down mid-page (not the bottom — bottom force-shows) → hidden while scrolling
  await page.evaluate(() => window.scrollTo(0, Math.min(800, document.documentElement.scrollHeight - window.innerHeight - 300)))
  // poll for hidden — rAF + idleMs 300 needs >150ms
  await expect.poll(async () => tabs.evaluate(el => getComputedStyle(el).transform), { timeout: 5000 }).not.toContain('matrix(1, 0, 0, 1, 0, 0)')

  // Stopping scroll re-shows after the 300ms idle window
  await page.waitForTimeout(600)
  await expect(tabs).toBeInViewport()

  await page.evaluate(() => window.scrollTo(0, 0))
  await page.waitForTimeout(600)
  await expect(tabs).toBeInViewport()
})

test('mobile graph renders and is usable', async ({ page }) => {
  await page.goto('/graph')
  const nodes = page.locator('.react-flow__node')
  await expect(nodes.first()).toBeVisible({ timeout: 45_000 })
  await expect.poll(async () => nodes.count(), { timeout: 45_000 }).toBeGreaterThan(100)

  expect(await overflowPx(page)).toBeLessThanOrEqual(1)

  // search overlay stays inside viewport and remains usable
  const input = page.getByPlaceholder('Search concepts…')
  await input.fill(sample[5].label.slice(0, 5))
  const box = await input.boundingBox()
  expect(box).toBeTruthy()
  const vw = page.viewportSize()!.width
  expect(box!.x).toBeGreaterThanOrEqual(0)
  expect(box!.x + box!.width).toBeLessThanOrEqual(vw)
  expect(box!.height).toBeGreaterThanOrEqual(36)

  // ambient toggle must NOT be offered on mobile (forced off there)
  await expect(page.getByTitle('Toggle edge flow animation')).toHaveCount(0)
})

test('mobile pinch-zoom works on the concept map', async ({ page }) => {
  await page.goto('/graph')
  await expect(page.locator('.react-flow__node').first()).toBeVisible({ timeout: 45_000 })

  const pane = page.locator('.react-flow__viewport')
  const before = await pane.evaluate(el => el.style.transform)
  await page.locator('.react-flow').hover({ position: { x: 180, y: 160 } })
  await page.mouse.wheel(0, -240)
  await page.waitForTimeout(1000)
  await expect.poll(async () => pane.evaluate(el => el.style.transform), { timeout: 5000 }).not.toBe(before)
})

test('mobile info panel wraps without overflowing', async ({ page }) => {
  const target = sample.find(c => c.prerequisites.length > 2) ?? sample[10]
  await page.goto(`/graph?concept=${encodeURIComponent(target.id)}`)
  const openBtn = page.getByText('Open concept →')
  await expect(openBtn).toBeVisible({ timeout: 30_000 })
  const box = await openBtn.boundingBox()
  expect(box).toBeTruthy()
  const vw = page.viewportSize()!.width
  expect(box!.x).toBeGreaterThanOrEqual(0)
  expect(box!.x + box!.width).toBeLessThanOrEqual(vw)
})

test('graph list toggle and search buttons meet 36px hit area', async ({ page }) => {
  await page.goto('/graph')
  const listBtn = page.locator('button[aria-label="Open concept list"]')
  await expect(listBtn).toBeVisible({ timeout: 30_000 })
  const listBox = await listBtn.boundingBox()
  expect(listBox).toBeTruthy()
  expect(listBox!.height).toBeGreaterThanOrEqual(36)
})

test('graph zoom controls are 44px on mobile', async ({ page }) => {
  await page.goto('/graph')
  await expect(page.locator('.react-flow__node').first()).toBeVisible({ timeout: 30_000 })
  const ctrl = page.locator('.react-flow__controls-button').first()
  const box = await ctrl.boundingBox()
  expect(box).toBeTruthy()
  expect(box!.width).toBeGreaterThanOrEqual(44)
  expect(box!.height).toBeGreaterThanOrEqual(44)
})

test('stacked table cards are centered with equal gutters', async ({ page }) => {
  for (const path of ['/docs/contributing', '/how-it-works']) {
    await page.goto(path)
    const stack = page.locator('div.sm\\:hidden.space-y-2').first()
    await expect(stack).toBeVisible()
    const firstCard = stack.locator('> div').first()
    const box = await firstCard.boundingBox()
    expect(box).toBeTruthy()
    const vw = page.viewportSize()!.width
    expect(Math.abs(box!.x - (vw - (box!.x + box!.width)))).toBeLessThanOrEqual(1)
  }
})

test('scroll containers are edge-to-edge and horizontally scrollable', async ({ page }) => {
  // Leaderboard table
  await page.goto('/leaderboard')
  const lb = page.locator('div.overflow-x-auto').first()
  await expect(lb).toBeVisible()
  let lbBox = await lb.boundingBox()
  expect(lbBox).toBeTruthy()
  expect(lbBox!.x).toBeLessThanOrEqual(1)
  expect(Math.abs(lbBox!.x + lbBox!.width - page.viewportSize()!.width)).toBeLessThanOrEqual(1)
  // Only scrollable when content overflows (320px table fits inside 390px viewport)
  const lbScrollable = await lb.evaluate(el => el.scrollWidth > el.clientWidth)
  if (lbScrollable) {
    await lb.evaluate(el => { el.scrollLeft = 50 })
    expect(await lb.evaluate(el => el.scrollLeft)).toBeGreaterThan(0)
  }

  // Yearly heatmap (switch view first)
  await page.goto('/profile')
  await page.getByRole('button', { name: 'Yearly' }).click()
  const yearly = page.locator('div.overflow-x-auto').first()
  await expect(yearly).toBeVisible()
  const yBox = await yearly.boundingBox()
  expect(yBox).toBeTruthy()
  expect(yBox!.x).toBeLessThanOrEqual(1)
  expect(Math.abs(yBox!.x + yBox!.width - page.viewportSize()!.width)).toBeLessThanOrEqual(1)
  const yScrollable = await yearly.evaluate(el => el.scrollWidth > el.clientWidth)
  if (yScrollable) {
    await yearly.evaluate(el => { el.scrollLeft = 50 })
    expect(await yearly.evaluate(el => el.scrollLeft)).toBeGreaterThan(0)
  }
})

test('double-clicking Check Answer fires exactly one POST', async ({ page }) => {
  let answerPosts = 0
  page.on('request', req => {
    if (req.method() === 'POST' && (req.url().includes('/api/study/answer') || req.url().includes('/api/quiz/answer'))) answerPosts++
  })

  // Mock Study seam: LessonQuiz → SubmitAnswer (CONTEXT.md Seam) — same as flows.spec.ts:26
  const LESSON = { title: 'Addition Basics', body: '# Addition Basics\nLearn', concepts: ['arith.add.single'] }
  const PRACTICE = { questions: [{ question: '5+4=?', answer: '9', explanation: '', source: 'curated' }], concept_id: 'arith.add.single' }
  const KP = { concept_id: 'arith.add.single', kps: [{ label: 'Use addition notation', section: 'Use Addition Notation', subgoals: [], worked_example: 'Add' }], diagram: '' }
  await page.route('**/api/lessons**', route => {
    const url = route.request().url()
    if (url.includes('/practice')) return route.fulfill({ json: PRACTICE })
    if (url.includes('/kp')) return route.fulfill({ json: KP })
    if (url.includes('/body')) return route.fulfill({ json: { title: LESSON.title, body: LESSON.body } })
    return route.fulfill({ json: { lessons: { arithmetic: [LESSON] } } })
  })
  await page.route('**/api/study/answer', r => r.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 10, expected_answer: '9' } }))
  await page.route('**/api/scores/**', r => r.fulfill({ json: { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 } }))
  await page.route('**/api/progress/**', r => r.fulfill({ json: {} }))
  await page.route('**/api/activity**', r => r.fulfill({ json: [] }))
  await page.route('**/api/weaknesses**', r => r.fulfill({ json: { by_domain: {} } }))
  await page.route('**/api/reviews/due**', r => r.fulfill({ json: { count: 0 } }))
  await page.route('**/api/courses**', r => r.fulfill({ json: { courses: [] } }))
  await page.route('**/api/transcript**', r => r.fulfill({ json: { courses: [] } }))
  await page.route('**/api/efficacy**', r => r.fulfill({ json: { concepts_touched: 1, first_pass_rate: 1, second_pass_rate: 1, avg_attempts_per_concept: 1, total_attempts: 1 } }))

  await page.goto('/study?lesson=' + encodeURIComponent(LESSON.title))
  await expect(page.getByText('5+4=?')).toBeVisible({ timeout: 30_000 })

  const input = page.locator('input[placeholder*="Your answer"]').first()
  await expect(input).toBeVisible({ timeout: 20_000 })
  await input.fill('9')
  await page.getByRole('button', { name: /Submit|Check Answer/ }).first().dblclick()
  await page.waitForTimeout(500)

  expect(answerPosts).toBe(1)
  await expect(page.getByText('Failed to submit answer')).toHaveCount(0)
})
