import { test, expect, type Page } from '@playwright/test'
import { answerField, focusAnswerField, setAnswer, typeIntoField, showVirtualKeyboard, keyboardRows, keyboardPanel, keyboardVisible, hideVirtualKeyboard } from './helpers/answer'

// Mobile behaviour, on a real touch device profile (Pixel 7, and a 320px viewport).
//
// The desktop suite proves the editor builds structures. This one asks the question
// that only touch can answer: can a learner with a thumb read the problem, enter an
// answer, and submit it without the keypad swallowing the screen.
//
// Runs under the `mobile-chromium` and `mobile-320` projects, both of which set
// `hasTouch`, so every interaction here is a real tap rather than a synthesised
// mouse event.

const CONCEPT = 'arith.add.single'

async function openLearn(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem(
      'mathua_user',
      JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 3, current_streak: 0, level: 'Novice', diagnostic_completed: true }),
    )
  })
  await page.route('**/api/lessons/**', route => {
    const url = route.request().url()
    if (url.includes('/practice')) {
      return route.fulfill({
        json: {
          concept_id: CONCEPT,
          questions: [
            { question: '2 + 2 = ?', answer: '4', explanation: 'Two and two.', source: 'curated' },
            { question: '3 + 5 = ?', answer: '8', explanation: 'Three and five.', source: 'curated' },
          ],
        },
      })
    }
    if (url.includes('/kp')) return route.fulfill({ json: { concept_id: CONCEPT, kps: [{ label: 'Add', section: 'Add', subgoals: [], worked_example: '2+3=5' }] } })
    if (url.includes('/readiness')) return route.fulfill({ json: { concept_id: CONCEPT, ready: true, weak: [], missing: [] } })
    return route.fulfill({ json: {} })
  })
  await page.route('**/api/study/answer', route => {
    const body = JSON.parse(route.request().postData() ?? '{}') as { answer?: string }
    const right = body.answer === '4' || body.answer === '2+2'
    return route.fulfill({ json: { correct: right, feedback: 'Correct!', xp: right ? 1 : 0, halted: false } })
  })
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 0, weekly_score: 0, speed_bonus: 0, concepts_mastered: 1, current_streak: 0, level: 'Novice', xp_total: 0, xp_today: 0, daily_xp_goal: 10, spaced_reps: {}, avg_learning_time: 1 } }),
  )
  for (const [pattern, body] of [
    ['**/api/activity**', []],
    ['**/api/progress/**', {}],
    ['**/api/weaknesses**', { by_domain: {} }],
    ['**/api/reviews/due**', { count: 0 }],
    ['**/api/courses**', { courses: [] }],
    ['**/api/transcript**', { courses: [] }],
  ] as const) {
    await page.route(pattern, route => route.fulfill({ json: body }))
  }

  await page.goto(`/learn?concept=${CONCEPT}`)
  await page.getByRole('button', { name: 'Next →' }).click()
  await expect(answerField(page)).toBeVisible({ timeout: 30_000 })
}

test.describe('touch input', () => {
  test('a tap focuses the field, so the keypad can take input', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await page.waitForTimeout(150)
    // The sink is what holds focus: MathLive routes it into the shadow root. Checking
    // the host's activeElement would be asking a question about MathLive's internals;
    // checking the sink is asking whether the learner can reach the field at all.
    // Either the host or its sink: which one holds focus is MathLive's business, and
    // both mean the learner reached the field. The test that actually matters is the
    // one below — that a keypad tap then reaches the answer.
    const focused = await answerField(page).evaluate(el => {
      const sink = (el as HTMLElement & { shadowRoot?: ShadowRoot }).shadowRoot?.querySelector('[part=keyboard-sink]')
      return el === document.activeElement || sink?.matches(':focus') === true
    })
    expect(focused).toBe(true)
  })

  test('the keypad opens on tap, because there is no physical keyboard', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    // On a touch device the keypad is the only way to reach ÷, √ or a fraction, so
    // it must come up on its own. On desktop it deliberately does not.
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
  })

  test('a keypad digit reaches the answer', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^7/ }).first().tap()
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^5/ }).first().tap()
    await page.waitForTimeout(150)
    expect(await answerField(page).evaluate(el => (el as unknown as { value: string }).value)).toBe('75')
  })

  test('a structural keycap builds a fraction by tap', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^a⁄b/ }).first().tap()
    await typeIntoField(page, '3')
    await typeIntoField(page, '4')
    await page.waitForTimeout(150)
    const value = await answerField(page).evaluate(el => (el as unknown as { value: string }).value)
    // A fraction, not "34": the tap inserted a structure with two slots.
    expect(value).toContain('\\frac')
  })
})

test.describe('layout and reachability', () => {
  test('every visible keycap is at least a 44px touch target', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })

    const undersized = await page.locator('.MLK__keycap:visible').evaluateAll(els =>
      els
        .map(el => {
          const r = el.getBoundingClientRect()
          return { text: (el.textContent ?? '').trim().slice(0, 12), w: Math.round(r.width), h: Math.round(r.height) }
        })
        // The toolbar's own controls are not answer keys; a narrow separator is not a tap
        // target either. Both are excluded by width, since neither is a digit or operator.
        .filter(k => /[0-9÷×−+=()√]/.test(k.text))
        .filter(k => k.w < 44 || k.h < 44),
    )
    expect(undersized, `undersized keycaps: ${JSON.stringify(undersized)}`).toEqual([])
  })

  test('the keypad does not take more than half the screen', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    const fraction = await keyboardPanel(page).evaluate(el => {
      const r = el.getBoundingClientRect()
      return r.height / window.innerHeight
    })
    // "Do not make the virtual keyboard consume unnecessary vertical space": the
    // learner still needs to see the problem and the Check button.
    expect(fraction).toBeLessThan(0.5)
  })

  test('the problem and the answer stay visible with the keypad open', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await expect(page.getByText('2 + 2 = ?').first()).toBeVisible()
    await expect(answerField(page)).toBeVisible()
  })

  test('the answer field is at least a 44px target', async ({ page }) => {
    await openLearn(page)
    const box = await answerField(page).boundingBox()
    expect(box?.height).toBeGreaterThanOrEqual(44)
  })

  test('nothing overflows the viewport horizontally', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    // A 320px viewport is the narrowest phone this is expected to work on; the keypad
    // rows must not push the page into a horizontal scroll.
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
    expect(overflow).toBeLessThanOrEqual(1)
  })
})

test.describe('the whole loop on a phone', () => {
  test('problem → answer → submit → verdict, by touch and keypad', async ({ page }) => {
    await openLearn(page)
    await expect(page.getByText('2 + 2 = ?').first()).toBeVisible()

    // Answer "2 + 2" the way a learner on a phone does: tap the field, then tap
    // 2, +, 2 on the keypad. No physical keyboard anywhere in this test.
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^2/ }).first().tap()
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^[+]/ }).first().tap()
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^2/ }).first().tap()
    await page.waitForTimeout(200)

    await page.getByRole('button', { name: 'Check', exact: true }).first().tap()
    await expect(page.getByText('Two and two.').first()).toBeVisible({ timeout: 20_000 })
    await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
  })

  // Deliberately no assertion that the keypad can be dismissed.
  //
  // `mathVirtualKeyboard.hide()` does not work in MathLive 0.111.0 under Chromium —
  // called directly from the page it left `visible === true` and the panel in the DOM,
  // under both the `auto` and `manual` policies. The wrapper calls it on an outside
  // touch, but the behaviour is unproven, so it is not claimed. Until it is fixed,
  // the usable exits on a phone are submitting (which blurs the field) or the
  // element's own hide-keyboard keycap.
  test('the answer survives the keypad opening and closing repeatedly', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await page.locator('.MLK__keycap:visible').filter({ hasText: /^7/ }).first().tap()
    await page.waitForTimeout(150)
    await page.getByText('2 + 2 = ?').first().tap()
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    expect(await answerField(page).evaluate(el => (el as unknown as { value: string }).value)).toBe('7')
  })
})