import { test, expect, type Page } from '@playwright/test'
import { answerField, answerLatex, focusAnswerField, setAnswer, typeIntoField, showVirtualKeyboard, keyboardRows, keyboardPanel, keyboardVisible, hideVirtualKeyboard, keycaps, glyphKeycap, loadingAnswerField, captureSubmittedAnswer, pressKeycap, useAlphabeticLayout, mathLiveSpacebar } from './helpers/answer'
import { installInstrumentation, delayMathLiveChunk } from './helpers/instrument'

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

/**
 * `waitForEditor: false` returns as soon as the answer area exists, without waiting
 * for MathLive. The loading-state tests need that: `openLearn`'s normal wait is for the
 * editor, and the editor arriving *is* the end of the window they are trying to act
 * inside. Using the normal wait made them pass or fail depending on how the two
 * deadlines interleaved -- which is exactly what they did.
 */
/**
 * Press MathLive's own layer switcher in the keypad toolbar.
 *
 * This is the only layer-navigation mechanism that works in 0.111.0: a keycap with a
 * `layer` property gains a `layer-switch` *class* and does nothing when pressed, while
 * the toolbar emits `data-layer` on its entries and a click assigns
 * `virtualKeyboard.currentLayer`. Layers were unreachable before this for exactly that
 * reason, so the switcher is asserted rather than assumed.
 */
async function switchLayer(page: Page, label: string) {
  // Forced: the switcher entries are small toolbar targets that Playwright's
  // stability check rejects, and MathLive's own handler reads `data-layer` from the
  // event target regardless of how the event was produced.
  await page
    // `:visible` matters: MathLive renders one toolbar per layout and only the active
    // one is shown, so the first match in DOM order is usually a hidden copy.
    .locator('.MLK__toolbar .layer-switch:visible')
    .filter({ hasText: new RegExp(`^${label}$`) })
    .first()
    .click({ force: true })
}

async function openLearn(
  page: Page,
  opts: { waitForEditor?: boolean; concept?: string } = {},
) {
  const concept = opts.concept ?? CONCEPT
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
          concept_id: concept,
          questions: [
            { question: '2 + 2 = ?', answer: '4', explanation: 'Two and two.', source: 'curated' },
            { question: '3 + 5 = ?', answer: '8', explanation: 'Three and five.', source: 'curated' },
          ],
        },
      })
    }
    if (url.includes('/kp')) return route.fulfill({ json: { concept_id: concept, kps: [{ label: 'Add', section: 'Add', subgoals: [], worked_example: '2+3=5' }] } })
    if (url.includes('/readiness')) return route.fulfill({ json: { concept_id: concept, ready: true, weak: [], missing: [] } })
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

  await page.goto(`/learn?concept=${concept}`)
  await page.getByRole('button', { name: 'Next →' }).click()
  if (opts.waitForEditor === false) {
    await expect(page.locator('.mathua-math-input')).toBeAttached({ timeout: 30_000 })
    return
  }
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

  // Both of these keys were broken in all five layouts and neither test existed.
  // The layouts named their shortcuts in `latex`, which MathLive does not consult:
  // it resolves a keycap as command -> insert -> key -> latex -> typedText(label),
  // and merges its own KEYCAP_SHORTCUTS table only for a bare-string keycap or one
  // carrying a `label` or `key`. So both keys missed the table *and* fell through to
  // the insert branch, typing their own source text into the field. A screenshot
  // looked fine and every existing assertion still passed, because the digit tests
  // only ever tapped 1-9 and nothing ever pressed delete.
  test('the zero key inserts 0, not the text [0]', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await keycaps(page).filter({ hasText: /^0/ }).first().tap()
    await page.waitForTimeout(150)
    // Was '[0]' before the fix.
    expect(await answerLatex(page)).toBe('0')
  })

  test('backspace deletes rather than inserting its own name', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await typeIntoField(page, '4')
    await typeIntoField(page, '2')
    expect(await answerLatex(page)).toBe('42')

    // The merged shortcut renders MathLive's own delete glyph, and `renderKeycap`
    // gives it a class of `action ...` *without* `MLK__keycap` -- so both the glyph
    // and the missing class are why this key was invisible to the old selectors.
    await glyphKeycap(page, 'delete-backward').first().tap()
    await page.waitForTimeout(150)
    // Was '4[backspace]' before the fix.
    expect(await answerLatex(page)).toBe('4')
  })
})

test.describe('layout and reachability', () => {
  test('every visible keycap is at least a 44px touch target', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })

    // `keycaps()`, not `.MLK__keycap`: MathLive omits that class from any keycap whose
    // class contains `separator`, `action`, `shift`, `fnbutton` or `bigfnbutton`, so a
    // `.MLK__keycap` audit measured only the keys that already worked and never looked
    // at backspace, the caret arrows, the dismiss key or a separator. Separators are
    // still excluded below, by class rather than by the accident of the old selector.
    const undersized = await keycaps(page).evaluateAll(els =>
      els
        .map(el => {
          const r = el.getBoundingClientRect()
          return {
            text: (el.textContent ?? '').trim().slice(0, 12),
            glyph: el.querySelector('use')?.getAttribute('href') ?? '',
            cls: el.className,
            w: Math.round(r.width),
            h: Math.round(r.height),
          }
        })
        // A separator is not a tap target. Everything else is.
        .filter(k => !/(^|\s)separator(\s|$)/.test(k.cls))
        .filter(k => k.w < 44 || k.h < 44),
    )
    expect(undersized, `undersized keycaps: ${JSON.stringify(undersized)}`).toEqual([])

    // And the audit must actually be reaching the glyph keys, or "no undersized
    // keycaps" would again mean "no keycaps measured". A `.MLK__keycap` selector
    // reported a clean pass over a keypad whose backspace and caret arrows were
    // invisible to it.
    // MathLive writes the glyph reference as `xlink:href`, so a bare `href` read
    // returns null and the guard would pass on an empty set.
    const measured = await keycaps(page).evaluateAll(els =>
      els.map(el => {
        const use = el.querySelector('use')
        return use?.getAttribute('href') ?? use?.getAttribute('xlink:href') ?? ''
      }),
    )
    expect(measured.filter(h => h.includes('delete-backward')).length,
      'the touch-target audit did not reach the backspace key').toBeGreaterThan(0)
    expect(measured.filter(h => h.includes('arrow-')).length,
      'the touch-target audit did not reach the caret arrows').toBeGreaterThan(0)
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

test.describe('the spacebar', () => {
  test('the space key inserts a space', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })

    await pressKeycap(page, '5')
    await pressKeycap(page, '\u2423')
    await pressKeycap(page, '3')

    // `\,` and not a bare space. TeX discards literal spaces in math mode, so a spacebar
    // that inserted `" "` would render as nothing and serialise as nothing; `\,` is the
    // payload the recorded corpus proves survives the round trip.
    expect(await answerLatex(page)).toBe('5\\,3')
  })

  test('MathLive\'s own spacebar works, which is the key that was reported broken', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await useAlphabeticLayout(page)

    const spacebar = mathLiveSpacebar(page)
    await expect(spacebar, 'no spacebar found in the alphabetic layout').toBeVisible()
    await spacebar.click()
    await page.waitForTimeout(120)

    expect(await answerLatex(page)).toContain('\\,')

    // And it is a space in the answer that gets submitted, which is the only reason to
    // want one: `Q R`, `a sqrt(b)` and `P=[[...]] D=[[...]]` are generator answer
    // formats a learner cannot type without it.
    await pressKeycap(page, 'r')
    expect(await answerLatex(page)).toBe('\\,r')
  })

  test('a space still rejoins a word typed with one, so letters are not fragmented', async ({ page }) => {
    // The regression risk in making spaces real. `fixImplicitLetterSpacing` collapses
    // `y e s` to `yes` on the *plain-text* answer, and `\,` has already become a plain
    // space by the time that runs -- but only because the collapse happens downstream of
    // the serialisation. If it were ever applied to the LaTeX instead, every letter a
    // learner typed would be welded to its neighbour. This is the test that notices.
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await useAlphabeticLayout(page)

    for (const k of ['y', 'e', 's']) {
      await pressKeycap(page, k)
      await mathLiveSpacebar(page).click()
      await page.waitForTimeout(60)
    }
    expect(await answerLatex(page)).toBe('y\\,e\\,s\\,')

    const submitted = captureSubmittedAnswer(page, '/api/study/answer')
    await page.getByRole('button', { name: 'Check', exact: true }).first().tap()
    await expect.poll(submitted).toBe('yes')
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

  test('no sound asset is ever requested', async ({ page }) => {
    // MathLive 0.111.0 defaults `MathfieldElement._soundsDirectory` to './sounds' and
    // `keypressSound` to 'keypress-standard.wav', loaded with a fetch on every
    // keystroke. Mathua ships `public/fonts` but never shipped `public/sounds`, so
    // every key was requesting a file that does not exist. The fix sets the static to
    // `null`, which is the documented way to load no sounds at all -- verified against
    // the 0.111.0 types, because the *instance* accessor throws rather than working.
    const sounds: string[] = []
    page.on('request', req => {
      if (/\/sounds\/|\.wav($|\?)/i.test(req.url())) sounds.push(req.url())
    })
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await keycaps(page).filter({ hasText: /^7/ }).first().tap()
    await keycaps(page).filter({ hasText: /^5/ }).first().tap()
    await page.waitForTimeout(600)
    expect(sounds, `MathLive asked for sound assets: ${sounds.join(', ')}`).toEqual([])
  })

  test('the keypad height is published so page layout can clear it', async ({ page }) => {
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    // `geometrychange` carries the new bounding rectangle; Mathua turns it into one CSS
    // custom property and a state attribute, so stylesheets can adapt without a
    // hardcoded pixel height that is wrong on every device.
    await expect
      .poll(() => page.evaluate(() => document.documentElement.dataset.mathuaKeyboard ?? null), {
        timeout: 10_000,
      })
      .toBe('open')
    const height = await page.evaluate(() =>
      document.documentElement.style.getPropertyValue('--mathua-keyboard-height'),
    )
    expect(Number.parseInt(height, 10)).toBeGreaterThan(0)
  })

  // Deliberately no assertion that the keypad can be dismissed.
  //
  // `mathVirtualKeyboard.hide()` does not work in MathLive 0.111.0 under Chromium —
  // called directly from the page it left `visible === true` and the panel in the DOM,
  // under both the `auto` and `manual` policies. The wrapper calls it on an outside
  // touch, but the behaviour is unproven, so it is not claimed. Until it is fixed,
  // the usable exits on a phone are submitting (which blurs the field) or the
  // element's own hide-keyboard keycap.
  test('an outside tap does not hide the keypad while the field keeps focus', async ({ page }) => {
    // The regression, asserted at the mechanism rather than the symptom.
    //
    // Mathua used to hide the keypad from a capture-phase `pointerdown` listener on
    // `document`, on any tap outside the field. MathLive raises the keyboard from
    // exactly one place -- a `focusin` listener -- so hiding it while focus stayed in
    // the field left no way back: no focus transition meant no re-show. Tapping a
    // non-focusable area did precisely that, and the learner had to keep tapping.
    //
    // Playwright cannot reproduce the strand itself, because `hide()` is a no-op under
    // Chromium (documented in docs/math-input.md). So this asserts the *call* never
    // happens, which is the part that is ours and which does fail on the old code.
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    expect(await keyboardVisible(page)).toBe(true)

    await page.evaluate(() => {
      const w = window as unknown as { mathVirtualKeyboard?: { hide: () => void }; __hideCalls?: number }
      const kb = w.mathVirtualKeyboard
      if (!kb) throw new Error('no virtual keyboard')
      w.__hideCalls = 0
      const original = kb.hide.bind(kb)
      kb.hide = () => {
        w.__hideCalls = (w.__hideCalls ?? 0) + 1
        original()
      }
    })

    // A heading: not focusable, so tapping it cannot move focus out of the field.
    //
    // `:visible` matters now. `/learn` carries a collapsed reference panel whose own headings sit
    // earlier in the document than anything else, and a heading inside a closed `<details>` is
    // display:none — so the first match was a heading no one can see or tap. Filtering by
    // visibility keeps the locator pointing at something a finger could actually reach, which is
    // what the rest of the test assumes.
    const heading = page.locator('h1:visible, h2:visible, h3:visible').first()
    await expect(heading).toBeVisible()
    await heading.tap()

    expect(await page.evaluate(() => (window as unknown as { __hideCalls?: number }).__hideCalls ?? 0)).toBe(0)

    // Whether focus survives an outside tap is platform behaviour, not ours:
    // Chromium blurs to the body, iOS Safari frequently does not. So this test does
    // not assert it. What matters, and is platform-neutral, is that another tap on
    // the field brings the keypad back — that is the recovery path, and it is what the
    // learner experiences as "I have to tap it several times".
    await page.waitForTimeout(400) // MathLive's own focusout -> hide timer is 300ms
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    expect(await keyboardVisible(page)).toBe(true)
  })

  // The chain that works, asserted in order, so a regression in any link is
  // attributable rather than showing up as one vague "submit did not happen".
  //
  // No sleeps anywhere: every step is an assertion that polls. The measured timings
  // behind this (e2e/answer-diagnosis.spec.ts) put input -> React at 2-7ms and
  // tap -> API response at 12-14ms, so there is nothing here worth waiting for.
  test('field interactive -> type -> React receives -> Check enables -> first tap submits once', async ({ page }) => {
    let posts = 0
    await page.on('request', req => {
      if (req.method() === 'POST' && req.url().includes('/api/study/answer')) posts++
    })
    await openLearn(page)

    // 1. The field exists and is usable. Polling, not a fixed wait: the no-input
    //    window is real and its length depends on the chunk (300ms cached here,
    //    3.1s with a 3s chunk delay on the network), so a sleep would be wrong in
    //    both directions.
    await expect(answerField(page)).toBeVisible({ timeout: 20_000 })
    await setAnswer(page, '')

    // 2. Enter an answer by tapping the keypad, and let the answer reach the host.
    //    Keycaps rather than `keyboard.type`, because at 320px the open virtual
    //    keyboard swallows synthetic keystrokes -- the sibling loop test taps for the
    //    same reason, and it is closer to what a phone user does anyway. The button's
    //    `disabled` is `checking || !answer.trim()`, so it is the only observable for
    //    React having received the answer.
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await keycaps(page).filter({ hasText: /^4/ }).first().tap()
    const check = page.getByRole('button', { name: 'Check', exact: true }).first()
    await expect(check).toBeEnabled({ timeout: 10_000 })

    // 3. The first tap reaches the submit handler.
    await check.tap()

    // 4. Exactly one POST, and the verdict comes back. The explanation text is the
    //    verdict here; asserting on the feedback body rather than a "correct" string
    //    because the panel shows the explanation, not the word correct.
    await expect(page.getByText('Two and two.').first()).toBeVisible({ timeout: 20_000 })
    expect(posts).toBe(1)
  })

  // The reported "I have to type Check Answer multiple times", as two assertions.
  //
  // The loading state used to be `<div className="h-12" aria-hidden />` — the right
  // size in the right place, so a tap looked like it landed on an answer field and then
  // did nothing, and every keystroke went nowhere. Measured in
  // `e2e/answer-diagnosis.spec.ts`: the window is 322ms with the chunk cached and 3116ms
  // with a 3s chunk delay, so on a real phone connection it is seconds. Value
  // propagation was never at fault (input to React state in 2-7ms, first tap to submit
  // in 12-14ms); the control simply was not there yet.
  //
  // The chunk is delayed so the window is wide enough to act inside deterministically.
  // Nothing here waits for a fixed duration to pass: each step waits for the state it
  // needs, and the delay is the condition under test rather than a sleep.
  // The editor is universal now, so the keypad has to answer what every grading type
  // asks for. Three of these were missing before and each blocked a real question:
  // no comma meant no tuple or ordering answer could be entered at all, and no
  // comparison key meant `dec.basics.compare` and `frac.ops.compare` could not be
  // answered on a phone without digging for a symbol row.
  test('the keypad can enter a tuple, so tuple and ordering concepts are answerable', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    await keycaps(page).filter({ hasText: /^3/ }).first().tap()
    await keycaps(page).filter({ hasText: /^,/ }).first().tap()
    await keycaps(page).filter({ hasText: /^7/ }).first().tap()
    await page.waitForTimeout(150)
    // A tuple, not "37": the comma is a separator, and without it these 14 concepts have
    // no way to enter their answer at all.
    expect(await answerLatex(page)).toBe('3,7')
  })

  test('the keypad carries the comparison signs the two comparison concepts need', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
// Reached through MathLive's own toolbar switcher, because that is the only
    // mechanism that works: a keycap with a `layer` property only gains a CSS class and
    // does nothing when pressed, while the toolbar emits `data-layer` and assigns
    // `currentLayer`. Layers were unreachable before this for the same reason.
    //
    // The two switches are inside the retry on purpose. MathLive assigns `currentLayer` and
    // re-renders the keypad asynchronously, and under CI's two workers that re-render can land
    // after both switches were sent — so switching once and then reading the keypad asserts
    // that the first switch took and measures the second one's re-render, which fails as
    // "the ≥ key is missing". Retrying the whole sequence states the actual property: eventually
    // both switches land and the key is there.
    const ge = keycaps(page).filter({ hasText: /≥/ })
    await expect(async () => {
      await switchLayer(page, '±÷')
      await switchLayer(page, 'abc')
      await expect(ge).toBeVisible({ timeout: 1_000 })
    }).toPass({ timeout: 20_000 })
    await ge.first().tap()
    await page.waitForTimeout(150)
    expect(await answerLatex(page)).toBe('\\ge')
  })

  test('offers an alphabet layout, so a yes/no answer can be typed at all', async ({ page }) => {
    // `yes` is the most common expected answer in the corpus and `no` the second. They
    // stay tappable through ChoiceOptions, but the typed path has to work too -- and
    // MathLive reads bare letters as implicit multiplication, so it needs `toPlainAnswer`
    // to rejoin them (test/plainAnswer.test.ts).
    await openLearn(page)
    await answerField(page).tap()
    await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
    // MathLive renders the switcher itself once `kb.layouts` has more than one entry.
    await expect(page.locator('.MLK__toolbar').first()).toContainText('abc')
  })

  test('the answer is typeable while the editor is still loading', async ({ page }) => {
    await installInstrumentation(page)
    await delayMathLiveChunk(page, 3000)
    await openLearn(page, { waitForEditor: false })

    // The window is genuinely open: no editor yet.
    await expect(loadingAnswerField(page)).toBeVisible({ timeout: 15_000 })
    expect(await page.locator('math-field').count()).toBe(0)

    // Type into what is actually on screen.
    await loadingAnswerField(page).fill('4')

    // React receives it, Check enables, and the first tap submits.
    const check = page.getByRole('button', { name: 'Check', exact: true }).first()
    await expect(check).toBeEnabled({ timeout: 15_000 })
    await check.tap()
    await expect(page.getByText('Two and two.').first()).toBeVisible({ timeout: 20_000 })
  })

  test('an answer typed while loading survives the editor arriving', async ({ page }) => {
    await installInstrumentation(page)
    await delayMathLiveChunk(page, 1500)
    await openLearn(page, { waitForEditor: false })

    await loadingAnswerField(page).fill('4')

    // The editor mounts and must adopt the answer rather than replacing it with an
    // empty field -- which is what a naive swap would do.
    await expect(answerField(page)).toBeVisible({ timeout: 30_000 })
    // `toHaveValue` does not work here: `math-field` is a custom element, not an
    // <input>, and Playwright rejects it. The LaTeX is read off `.value` directly.
    await expect.poll(() => answerLatex(page), { timeout: 15_000 }).toBe('4')
  })

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

