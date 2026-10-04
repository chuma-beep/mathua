// Timeline instrumentation for the answer → submit path.
//
// Exists to answer one question with measurements instead of argument: when a
// learner taps Check Answer and nothing happens, which of these is stuck?
//
//   1. the field is not interactive yet (chunk/fonts still loading)
//   2. the field's value changed but never reached React state
//   3. React state is right but the button is stale/disabled
//   4. the button is enabled but the tap never reaches the submit handler
//
// Every mark below is `performance.now()`, which is already relative to
// `performance.timeOrigin` — i.e. to navigation start. That is what makes the
// numbers comparable between a cold throttled load and a warm one.
//
// Nothing here records an answer. The brief was explicit that complete answers do
// not belong in instrumentation, and it is the right call beyond privacy: a log
// full of `answer=42` is a log full of student data nobody needs. Presence,
// length, and whether the field is empty are enough to place a failure on this
// timeline.
//
// The recorder is installed with `addInitScript`, so it is in place before the
// first byte of application code runs. That matters: the Suspense fallback is
// mounted and replaced long before any test-side hook could attach.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import type { Page } from '@playwright/test'

export interface TraceMark {
  /** ms since navigation start. */
  t: number
  name: string
  data?: Record<string, unknown>
}

export interface AnswerTrace {
  marks: TraceMark[]
  /** Convenience: the first mark with this name. */
  at: (name: string) => number | undefined
  /** Every mark with this name. */
  all: (name: string) => TraceMark[]
  /** Ordered `name @ t` lines, for eyeballing a timeline. */
  format: () => string
}

/** Install the recorder. Must be called before `page.goto`. */
export async function installInstrumentation(page: Page): Promise<void> {
  await page.addInitScript(() => {
    const marks: { t: number; name: string; data?: Record<string, unknown> }[] = []
    let stopped = false
    const mark = (name: string, data?: Record<string, unknown>) => {
      if (stopped) return
      marks.push({ t: performance.now(), name, ...(data ? { data } : {}) })
    }
    const w = window as unknown as { __answerTrace?: unknown }
    w.__answerTrace = marks

    mark('script-installed')

    // --- the field ------------------------------------------------------
    //
    // Two distinct transitions, and the gap between them is the whole point:
    //   mathinput-shell   the wrapper div exists, but Suspense is showing its
    //                     fallback, so there is no editable control at all
    //   mathlive-mounted  the <math-field> custom element is in the DOM
    //   mathlive-upgraded the element is actually usable (property access works)
    let shellSeen = false
    let fieldSeen = false
    let upgraded = false

    const checkField = () => {
      const shell = document.querySelector('.mathua-math-input')
      if (shell && !shellSeen) {
        shellSeen = true
        mark('mathinput-shell', {
          hasMathField: !!shell.querySelector('math-field'),
          hasFallbackDiv: !!shell.querySelector('div[aria-hidden]'),
          html: (shell.innerHTML || '').slice(0, 120),
        })
      }
      const mf = document.querySelector('math-field') as (HTMLElement & { value?: string }) | null
      if (!mf || fieldSeen) return
      fieldSeen = true
      mark('mathlive-mounted')

      // Upgraded means the custom element's own properties are live. A
      // not-yet-upgraded element is in the DOM but has no `value`, which is the
      // difference between "there is a box" and "there is an editor".
      try {
        const probe = typeof mf.value
        upgraded = true
        mark('mathlive-upgraded', { valueType: probe, hasFocusFn: typeof (mf as unknown as { focus?: unknown }).focus })
      } catch (e) {
        mark('mathlive-upgrade-failed', { error: String(e) })
      }

      // Has the custom element definition landed? This is separate from mounting:
      // React can insert the element before the bundle registers it.
      mark('mathlive-registered', { registered: !!customElements.get('math-field') })

      mf.addEventListener('input', () => {
        const v = (mf as unknown as { value?: string }).value ?? ''
        mark('math-input', { len: v.length, empty: v.trim().length === 0 })
      })
      mf.addEventListener('change', () => mark('math-change'))
      mf.addEventListener('focus', () => mark('math-focus'))
      mf.addEventListener('blur', () => mark('math-blur'))
      mf.addEventListener('pointerdown', () => mark('math-pointerdown'))
    }

    // Observe `document`, not `document.documentElement`. Init scripts run before the
    // root element is parsed, so `documentElement` is null here and `observe(null)`
    // throws — which killed every mark after `script-installed` and left a trace that
    // looked like the page never rendered. `document` is always safe to observe.
    const obs = new MutationObserver(checkField)
    const WATCH = { childList: true, subtree: true } as const
    obs.observe(document, WATCH)
    checkField()

    // --- the Check button ------------------------------------------------
    //
    // `disabled` on this button is `checking || !answer.trim()`, so it is the
    // only externally observable proxy for "does the host have a non-empty
    // answer". There is no way to read React state from the page, and inventing
    // one would mean shipping a debug hook in production, so the button is the
    // instrument. Transition times are what matter here.
    let button: HTMLButtonElement | null = null
    const attachButton = () => {
      if (button && button.isConnected) return
      const b = document.querySelector('button[type="submit"]') as HTMLButtonElement | null
      if (!b || b === button) return
      button = b
      mark('check-button-found', { disabled: b.disabled })
      b.addEventListener('pointerdown', () =>
        mark('check-pointerdown', { disabled: (b as HTMLButtonElement).disabled }),
      )
      b.addEventListener('click', () =>
        mark('check-click', { disabled: (b as HTMLButtonElement).disabled }),
      )
    }
    const bobs = new MutationObserver(attachButton)
    bobs.observe(document, WATCH)
    attachButton()

    // Every disabled transition on whatever the current button is.
    let lastDisabled: boolean | null = null
    const watchDisabled = () => {
      const b = document.querySelector('button[type="submit"]') as HTMLButtonElement | null
      if (!b) return
      if (b.disabled !== lastDisabled) {
        lastDisabled = b.disabled
        mark('check-disabled-change', { disabled: b.disabled })
      }
    }
    const dobs = new MutationObserver(watchDisabled)
    dobs.observe(document, { ...WATCH, attributes: true, attributeFilter: ['disabled'] })

    // --- the network -----------------------------------------------------
    //
    // Patched in-page so request start and response land on the same clock as
    // everything else. Node-side `page.on('request')` would need its timestamp
    // converted across the process boundary.
    const ANSWER = /\/(api\/(study|goal|quiz)\/answer|api\/answer)/
    const originalFetch = window.fetch
    window.fetch = function patchedFetch(this: unknown, ...args: Parameters<typeof fetch>) {
      const url = String(args[0])
      const isAnswer = ANSWER.test(url)
      if (isAnswer) mark('api-request-start', { path: safePath(url) })
      return originalFetch.apply(this, args).then(
        r => {
          if (isAnswer) mark('api-response', { path: safePath(url), status: r.status })
          return r
        },
        e => {
          if (isAnswer) mark('api-rejected', { path: safePath(url), error: String(e) })
          throw e
        },
      )
    }

    const safePath = (u: string) => {
      try {
        return new URL(u, location.href).pathname
      } catch {
        return u
      }
    }

    // Let the app settle, then stop recording so a late observer callback cannot
    // add noise after the trace is read.
    // Resource timing, so "is the throttle actually applied?" is answerable rather
    // than assumed. If a 4 Mbps profile does not slow the MathLive chunk down, the
    // no-input window is CPU-bound (parse + evaluate), not network-bound, and
    // throttling the network is measuring the wrong thing entirely.
    ;(window as unknown as { __resources?: () => unknown }).__resources = () => {
      const rows = performance.getEntriesByType('resource') as PerformanceResourceTiming[]
      const nav = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming | undefined
      // The largest transfers, not the slowest: the question is whether the no-input
      // window is waiting on bytes or on the main thread.
      const big = rows
        .filter(r => r.transferSize > 30_000 || r.decodedBodySize > 200_000)
        .sort((a, b) => b.decodedBodySize - a.decodedBodySize)
        .slice(0, 6)
        .map(r => ({
          name: r.name.split('/').pop()?.slice(0, 44),
          kb: Math.round(r.transferSize / 1024),
          decodedKb: Math.round(r.decodedBodySize / 1024),
          ms: Math.round(r.duration),
          startMs: Math.round(r.startTime),
        }))
      return {
        navMs: nav ? Math.round(nav.loadEventEnd) : undefined,
        domContentLoaded: nav ? Math.round(nav.domContentLoadedEventEnd) : undefined,
        resources: big,
        totalKb: Math.round(rows.reduce((n, r) => n + r.transferSize, 0) / 1024),
        count: rows.length,
      }
    }

    ;(window as unknown as { __stopTrace?: () => void }).__stopTrace = () => {
      stopped = true
      obs.disconnect()
      bobs.disconnect()
      dobs.disconnect()
    }
  })
}

export async function readTrace(page: Page): Promise<AnswerTrace> {
  const marks = (await page.evaluate(
    () => (window as unknown as { __answerTrace?: { t: number; name: string; data?: Record<string, unknown> }[] }).__answerTrace ?? [],
  )) as TraceMark[]

  const all = (name: string) => marks.filter(m => m.name === name)
  return {
    marks,
    all,
    at: name => marks.find(m => m.name === name)?.t,
    format: () =>
      marks
        .map(m => {
          const data = m.data && Object.keys(m.data).length ? ' ' + JSON.stringify(m.data) : ''
          return `${m.t.toFixed(0).padStart(6)}ms  ${m.name}${data}`
        })
        .join('\n'),
  }
}

/**
 * Apply a realistic mobile network profile via CDP.
 *
 * Not decorative. The hypothesis under test is that the field is not interactive
 * for a meaningful time on a real connection, and an unthrottled localhost run
 * cannot distinguish "slow" from "instant" -- it will always look instant.
 */
export async function throttleMobile4G(page: Page): Promise<void> {
  const client = await page.context().newCDPSession(page)
  await client.send('Network.enable')
  // Cache off. Without this the measurement is meaningless: the MathLive chunk came
  // back with `transferSize: 0`, i.e. served from cache, so throttling could not slow
  // it and the no-input window sat at ~310ms whether the network was 4 Mbps or local.
  // A cold phone is the case being investigated, so the cache has to go.
  await client.send('Network.setCacheDisabled', { cacheDisabled: true })
  await client.send('Network.emulateNetworkConditions', {
    offline: false,
    // Regular 4G: ~9ms RTT, 4 Mbps down, 2 Mbps up. Chosen over "slow 3G"
    // because that is what most phones are actually on.
    latency: 90,
    downloadThroughput: (4 * 1024 * 1024) / 8,
    uploadThroughput: (2 * 1024 * 1024) / 8,
  })
}

/** Turn throttling off again, so a warm measurement is possible on the same page. */
export async function clearThrottle(page: Page): Promise<void> {
  const client = await page.context().newCDPSession(page)
  await client.send('Network.enable')
  await client.send('Network.setCacheDisabled', { cacheDisabled: false })
  await client.send('Network.emulateNetworkConditions', {
    offline: false,
    latency: 0,
    downloadThroughput: -1,
    uploadThroughput: -1,
  })
}

/**
 * Delay the MathLive chunk by a known amount, so the no-input window can be measured
 * as a function of chunk latency.
 *
 * Two things this had to get right, both learned from numbers that were wrong first.
 *
 * The chunk is identified by scanning the build output for MathLive's own markers
 * rather than by a hardcoded hash. A hash silently stops matching the moment anything
 * is rebuilt, and a delay that does not apply produces a confident, plausible,
 * completely fictional number.
 *
 * And the delay reports whether it fired. `window.__chunkDelayFired` is read after the
 * scenario: a run where the route never matched must not be allowed to contribute a
 * measurement, which is exactly the failure this whole spec exists to catch.
 */
export async function delayMathLiveChunk(page: Page, ms: number): Promise<void> {
  if (ms <= 0) return
  // process.cwd(), not import.meta: Playwright loads specs through CommonJS here, so
  // import.meta is a syntax error. Playwright runs from the package root.
  const root = join(process.cwd(), 'out', '_next', 'static', 'chunks')
  const names: string[] = []
  for (const file of readdirSync(root)) {
    if (!file.endsWith('.js')) continue
    const body = readFileSync(join(root, file), 'utf8')
    if (body.includes('svg-delete-backward') && body.includes('mathVirtualKeyboard')) names.push(file)
  }
  if (names.length === 0) {
    throw new Error(
      'delayMathLiveChunk: no MathLive chunk found in out/_next/static/chunks. Run `next build` first.',
    )
  }
  // Context-level, not page-level. The chunk is served from Chromium's memory cache
  // in this harness, so no network request is made for it and a page-level route
  // never fires -- verified by the interception counter, which reported 0 for every
  // delay. Intercepting at the context is what forces the request to happen.
  await page.context().route(
    url => url.pathname.endsWith('.js') && names.some(n => url.pathname.endsWith(n)),
    async route => {
      await page.evaluate(() => {
        const w = window as unknown as { __chunkDelayFired?: number }
        w.__chunkDelayFired = (w.__chunkDelayFired ?? 0) + 1
      }).catch(() => {})
      await new Promise(r => setTimeout(r, ms))
      await route.continue()
    },
  )
}

/** How many times the injected chunk delay actually intercepted a request. */
export async function chunkDelayFired(page: Page): Promise<number> {
  return page
    .evaluate(() => (window as unknown as { __chunkDelayFired?: number }).__chunkDelayFired ?? 0)
    .catch(() => 0)
}
