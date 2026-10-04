// Diagnosis: when a learner taps Check Answer and nothing happens, which stage is
// stuck? This spec only measures. It changes no behaviour and asserts nothing about
// a fix, because the fix has not been chosen yet -- the brief was to establish the
// failure point first.
//
// The four candidates, and the mark that separates each pair:
//
//   1. field not interactive yet   -> no `mathlive-mounted` before the tap
//   2. value not propagating       -> `math-input` fires, `check-disabled-change`
//                                    to enabled never arrives
//   3. button state stale          -> value propagated, but the button never enabled
//   4. tap never reaches submit    -> `check-pointerdown` and `check-click` present
//                                    but no `api-request-start`
//
//
// MEASURED RESULT (2026-10-03, mobile-chromium / mobile-320 / chromium)
//
// The failure is (1), field readiness. (2), (3) and (4) are not implicated: once the
// field is interactive, `input -> React state` is 2-7ms and `first tap -> API
// response` is 12-14ms, on every platform and at every chunk latency.
//
// There is a window between the answer UI appearing and the editor being usable,
// during which the only thing on screen where the answer goes is
// `<div class="h-12" aria-hidden>`. It is not focusable, takes no input, and is the
// right height in the right place, so a tap on it is swallowed silently. Check Answer
// is correctly disabled throughout, because no answer exists. The learner's tap, then
// their typing, then their tap on Check all do nothing, and they do it again.
//
// Window length = MathLive chunk latency + ~310ms of parse and evaluate, measured by
// delaying the chunk and watching the window track it:
//
//   chunk delay      window (mobile-chromium)
//        0ms            322ms
//      500ms            611ms
//     1500ms           1611ms
//     3000ms           3116ms
//
// The same 3000ms delay on the waited path gives 3112ms, so waiting for the field
// costs the learner nothing once it is there.
//
// CAVEAT -- the 310ms floor is a floor for *this* harness, not a phone. It is measured
// on a desktop CPU with the chunk served from Chromium's memory cache, which is why
// CDP network throttling could not move it: the chunk reported transferSize 0 and
// completed in 39ms under a 4 Mbps profile while the main chunks correctly took
// 1100-1400ms. So the network contribution is unmeasured here, and the CPU cost of
// parsing a 794 kB decoded bundle on a mid-range phone is unknown. Both are larger
// than what is above. Treat these numbers as a floor.
//
// SECONDARY OBSERVATION, needs a device. At 320px with the field interactive and the
// answer present, Check Answer was enabled but the tap did not land within 1500ms in
// two scenarios, while the identical flow succeeded on mobile-chromium. Whether the
// keypad is covering the button cannot be settled without a real 320px device.
//
// Output is a printed timeline per scenario. A diagnosis you cannot read is not a
// diagnosis, and the numbers are the deliverable here.
//
// Runs under `chromium` (desktop) and both mobile profiles.

import { test, expect, type Page } from '@playwright/test'
import { answerField, focusAnswerField, setAnswer, typeIntoField } from './helpers/answer'
import {
  installInstrumentation,
  readTrace,
  throttleMobile4G,
  clearThrottle,
  delayMathLiveChunk,
  chunkDelayFired,
} from './helpers/instrument'

const CONCEPT = 'arith.add.single'

async function mockApi(page: Page) {
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem(
      'mathua_user',
      JSON.stringify({
        student_id: 's1',
        name: 'Tester',
        username: 'tester',
        concepts_mastered: 3,
        current_streak: 0,
        level: 'Novice',
        diagnostic_completed: true,
      }),
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
    if (url.includes('/kp'))
      return route.fulfill({ json: { concept_id: CONCEPT, kps: [{ label: 'Add', section: 'Add', subgoals: [], worked_example: '2+3=5' }] } })
    if (url.includes('/readiness'))
      return route.fulfill({ json: { concept_id: CONCEPT, ready: true, weak: [], missing: [] } })
    return route.fulfill({ json: {} })
  })
  await page.route('**/api/study/answer', route =>
    route.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 1, halted: false } }),
  )
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
}

/**
 * Press something the way the project can.
 *
 * The desktop profile has no touch support, so `tap()` throws outright there. Using
 * the wrong one would have produced a desktop run full of failures that say nothing
 * about the bug, and a mobile run that quietly used a mouse.
 */
async function activate(
  page: Page,
  target: { tap: (o?: { timeout?: number }) => Promise<void>; click: (o?: { timeout?: number }) => Promise<void> },
  timeout = 15_000,
): Promise<string> {
  const touch = await page.evaluate(() => typeof window !== 'undefined' && 'ontouchstart' in window)
  try {
    if (touch) await target.tap({ timeout })
    else await target.click({ timeout })
    return 'landed'
  } catch (e) {
    // Worth recording rather than throwing: "the tap could not even be delivered"
    // is a distinct failure from "the tap was delivered and nothing happened", and
    // Playwright's actionability wait conflating them would hide the difference.
    return `NOT DELIVERED (${String(e instanceof Error ? e.message : e).split('\n')[0].slice(0, 80)})`
  }
}

/** Navigate to the point where the answer UI exists but may not be interactive. */
async function reachAnswerUi(page: Page) {
  await page.goto(`/learn?concept=${CONCEPT}`)
  await page.getByRole('button', { name: 'Next →' }).click()
}

function report(label: string, trace: Awaited<ReturnType<typeof readTrace>>) {
  const t = (n: string) => trace.at(n)
  const first = (n: string) => trace.all(n)[0]?.data
  console.log(
    [
      `\n===== ${label} =====`,
      `field mounted      ${fmt(t('mathlive-mounted'))}`,
      `field upgraded     ${fmt(t('mathlive-upgraded'))}`,
      `shell mounted      ${fmt(t('mathinput-shell'))}`,
      `first input        ${fmt(t('math-input'))}`,
      `first focus        ${fmt(t('math-focus'))}`,
      `button enabled     ${fmt(t('check-disabled-change'))}`,
      `pointerdown        ${fmt(t('check-pointerdown'))} disabled=${first('check-pointerdown')?.disabled}`,
      `click              ${fmt(t('check-click'))} disabled=${first('check-click')?.disabled}`,
      `api request        ${fmt(t('api-request-start'))}`,
      `api response       ${fmt(t('api-response'))}`,
      `--- raw ---`,
      trace.format(),
    ].join('\n'),
  )
  function fmt(v: number | undefined) {
    return v === undefined ? '   (never)' : `${v.toFixed(0).padStart(6)}ms`
  }
}

/**
 * Each scenario runs its interaction, then prints the timeline **whether or not it
 * worked**. The failing runs are the interesting ones, so a trace that only prints on
 * success is a trace that never shows you the bug.
 */
async function scenario(
  page: Page,
  label: string,
  run: () => Promise<{ note?: string } | void>,
): Promise<void> {
  let note: string | undefined
  let failure: string | undefined
  try {
    const r = await run()
    note = (r as { note?: string } | undefined)?.note
  } catch (e) {
    failure = String(e instanceof Error ? e.message : e).split('\n')[0]
  }
  const trace = await readTrace(page)
  const fired = await chunkDelayFired(page)
  const resources = await page
    .evaluate(() => (window as unknown as { __resources?: () => unknown }).__resources?.())
    .catch(() => undefined)
  const t = (n: string) => trace.at(n)
  const first = (n: string) => trace.all(n)[0]?.data
  const fmt = (v: number | undefined) => (v === undefined ? '  (never)' : `${v.toFixed(0).padStart(6)}ms`)
  const mounted = t('mathlive-mounted')
  const shell = t('mathinput-shell')
  console.log(
    [
      ``,
      `===== ${label} =====`,
      `  SHELL (no input yet)   ${fmt(shell)}`,
      `  FIELD interactive at  ${fmt(mounted)}`,
      `  no-input window        ${mounted !== undefined && shell !== undefined ? (mounted - shell).toFixed(0) + 'ms' : 'n/a'}`,
      `  first input            ${fmt(t('math-input'))}`,
      `  input -> React         ${t('math-input') !== undefined && trace.all('check-disabled-change').some(m => m.data?.disabled === false) ? (() => { const en = trace.all('check-disabled-change').find(m => m.data?.disabled === false); return en ? (en.t - t('math-input')!).toFixed(0) + 'ms' : 'never enabled' })() : 'never enabled'}`,
      `  tap -> api response    ${t('check-click') !== undefined && t('api-response') !== undefined ? (t('api-response')! - t('check-click')!).toFixed(0) + 'ms' : 'NO SUBMIT'}`,
      `  pointerdown disabled   ${String(first('check-pointerdown')?.disabled)}`,
      `  click disabled         ${String(first('check-click')?.disabled)}`,
      note ? `  note: ${note}` : '',
      note && note.includes('delayApplied')
        ? ''
        : `  chunk delay intercepted: ${fired} request(s)`,
      failure ? `  FAILURE: ${failure}` : '',
      `  --- resources ---`,
      JSON.stringify(resources, null, 1),
      `  --- raw ---`,
      trace.format(),
    ]
      .filter(Boolean)
      .join('\n'),
  )
}

test.describe('Check Answer diagnosis', () => {
  // A cold throttled load fetches the app, the MathLive chunk (~215 kB gzip) and
  // ~300 kB of webfonts over a 4 Mbps link, and these scenarios wait for states
  // rather than sleeping for fixed periods.
  test.setTimeout(240_000)

  // The window as a function of how long the MathLive chunk takes to arrive. The
  // chunk is ~215 kB gzip / 794 kB decoded, so these delays stand in for roughly
  // "4G", "poor 4G", "3G" and "edge" on a real phone.
  for (const delayMs of [0, 500, 1500, 3000]) {
    test(`G · no-input window with a ${delayMs}ms MathLive chunk delay`, async ({ page }) => {
      await installInstrumentation(page)
      await mockApi(page)
      await delayMathLiveChunk(page, delayMs)
      await scenario(page, `G · chunk delay ${delayMs}ms · typed immediately`, async () => {
        await reachAnswerUi(page)
        const shell = page.locator('.mathua-math-input')
        await shell.waitFor({ state: 'attached', timeout: 60_000 })
        const fieldTap = await activate(page, shell, 1500)
        const typed = await page.keyboard.type('4').then(() => 'typed').catch(() => 'TYPE LOST')
        const check = page.getByRole('button', { name: 'Check', exact: true }).first()
        await check.waitFor({ state: 'attached', timeout: 60_000 })
        const enabledAtTap = await check.isEnabled()
        const checkTap = await activate(page, check, 1500)
        await page.waitForTimeout(2000)
        return {
          note: `field tap: ${fieldTap} · ${typed} · Check enabled at tap: ${enabledAtTap} · Check tap: ${checkTap}${
            delayMs > 0 ? ` · delayApplied=${await chunkDelayFired(page)}` : ''
          }`,
        }
      })
    })
  }

  test('H · no-input window with a 3000ms chunk delay, waited for the field', async ({ page }) => {
    await installInstrumentation(page)
    await mockApi(page)
    await delayMathLiveChunk(page, 3000)
    await scenario(page, 'H · chunk delay 3000ms · waited for field', async () => {
      await reachAnswerUi(page)
      await expect(answerField(page)).toBeVisible({ timeout: 120_000 })
      await setAnswer(page, '')
      await focusAnswerField(page)
      await typeIntoField(page, '4')
      const check = page.getByRole('button', { name: 'Check', exact: true }).first()
      await expect(check).toBeEnabled({ timeout: 60_000 })
      await activate(page, check)
      await page.waitForTimeout(1500)
    })
  })

  test('C · warm, unthrottled, typed the instant the UI appears', async ({ page }) => {
    await installInstrumentation(page)
    await mockApi(page)
    await scenario(page, 'C · WARM / unthrottled / typed immediately', async () => {
      await reachAnswerUi(page)
      const shell = page.locator('.mathua-math-input')
      await shell.waitFor({ state: 'attached', timeout: 60_000 })
      // 1500ms is Playwright's default actionability budget, not an arbitrary sleep:
      // it is how long a tap can take to be *delivered* on a loaded page. Anything
      // slower than that and a human would have tapped again.
      const fieldTap = await activate(page, shell, 1500)
      const typed = await page.keyboard.type('4').then(() => 'typed').catch(() => 'TYPE LOST')
      const check = page.getByRole('button', { name: 'Check', exact: true }).first()
      await check.waitFor({ state: 'attached', timeout: 60_000 })
      const enabledAtTap = await check.isEnabled()
      const checkTap = await activate(page, check, 1500)
      await page.waitForTimeout(2000)
      return {
        note: `immediate field tap: ${fieldTap} · ${typed} · Check enabled at tap: ${enabledAtTap} · immediate Check tap: ${checkTap}`,
      }
    })
  })

  test('D · warm, waited, single tap', async ({ page }) => {
    await installInstrumentation(page)
    await mockApi(page)
    await scenario(page, 'D · WARM / waited / single tap', async () => {
      await reachAnswerUi(page)
      await expect(answerField(page)).toBeVisible({ timeout: 60_000 })
      await setAnswer(page, '')
      await focusAnswerField(page)
      await typeIntoField(page, '4')
      const check = page.getByRole('button', { name: 'Check', exact: true }).first()
      await expect(check).toBeEnabled({ timeout: 30_000 })
      await activate(page, check)
      await page.waitForTimeout(1500)
    })
  })

  test('E · replacing an existing answer, then clearing it', async ({ page }) => {
    await installInstrumentation(page)
    await mockApi(page)
    await scenario(page, 'E · WARM / replace then clear', async () => {
      await reachAnswerUi(page)
      await expect(answerField(page)).toBeVisible({ timeout: 60_000 })
      await setAnswer(page, '')
      await focusAnswerField(page)
      await typeIntoField(page, '4')
      const check = page.getByRole('button', { name: 'Check', exact: true }).first()
      await expect(check).toBeEnabled({ timeout: 30_000 })

      // Replace: select all, then retype.
      await page.keyboard.press('ControlOrMeta+a')
      await typeIntoField(page, '9')
      const enabledAfterReplace = await check.isEnabled()
      await activate(page, check)
      await page.waitForTimeout(1200)

      // Clear entirely, then confirm the button disables again.
      await page.keyboard.press('ControlOrMeta+a')
      await page.keyboard.press('Backspace')
      await page.waitForTimeout(400)
      const enabledAfterClear = await check.isEnabled()
      return {
        note: `enabled after replace=${enabledAfterReplace} · enabled after clearing everything=${enabledAfterClear}`,
      }
    })
  })

  test('F · desktop, no touch, throttled', async ({ page }) => {
    await installInstrumentation(page)
    await mockApi(page)
    await throttleMobile4G(page)
    await scenario(page, 'F · DESKTOP / throttled / waited', async () => {
      await reachAnswerUi(page)
      await expect(answerField(page)).toBeVisible({ timeout: 120_000 })
      await setAnswer(page, '')
      await focusAnswerField(page)
      await typeIntoField(page, '4')
      const check = page.getByRole('button', { name: 'Check', exact: true }).first()
      await expect(check).toBeEnabled({ timeout: 60_000 })
      await activate(page, check)
      await page.waitForTimeout(1500)
    })
  })
})
