import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import {
  subscribeChrome,
  getChromeVisible,
  toggleChrome,
} from '../lib/chrome'

// Mirrors the constants in lib/chrome.ts. Kept literal rather than exported so a
// change to the real values shows up here as a failing expectation instead of
// silently moving the goalposts.
const IDLE_MS = 300
const BOTTOM_AUTO_HIDE_MS = 2000
const SETTLE_MS = 250

// lib/chrome.ts is a module singleton (that is the point of it), so each test
// subscribes and unsubscribes to get a clean slate: detaching resets the state.
describe('chrome visibility', () => {
  let unsubscribe: (() => void) | null = null

  // jsdom reports every layout as zero, so the scroll geometry is faked.
  const setViewport = (scrollHeight: number, innerHeight: number) => {
    Object.defineProperty(document.documentElement, 'scrollHeight', {
      value: scrollHeight,
      configurable: true,
    })
    Object.defineProperty(window, 'innerHeight', {
      value: innerHeight,
      configurable: true,
    })
  }

  const scrollTo = (y: number) => {
    Object.defineProperty(window, 'scrollY', { value: y, configurable: true })
    window.dispatchEvent(new Event('scroll'))
  }

  // The handler defers to requestAnimationFrame before touching state.
  const flushFrame = () => new Promise<void>(r => requestAnimationFrame(() => r()))

  const scrollToAndSettle = async (y: number) => {
    scrollTo(y)
    await flushFrame()
  }

  // Chrome is inert to scroll for a moment after it moves itself, so a
  // deliberate gesture has to start outside that window.
  const settle = () => vi.advanceTimersByTime(SETTLE_MS + 10)

  beforeEach(() => {
    vi.useFakeTimers()
    vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
      cb(0)
      return 0
    })
    setViewport(4000, 800)
    // Seed the scroll position before subscribing: attach() reads it as the
    // baseline for direction, so a leftover scrollY from a previous test would
    // make the first scroll a no-op (delta 0).
    Object.defineProperty(window, 'scrollY', { value: 0, configurable: true })
    unsubscribe = subscribeChrome(() => {})
  })

  afterEach(() => {
    unsubscribe?.()
    unsubscribe = null
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  it('starts visible', () => {
    expect(getChromeVisible()).toBe(true)
  })

  it('hides on a downward scroll and shows again when scrolling up', async () => {
    await scrollToAndSettle(600)
    expect(getChromeVisible()).toBe(false)

    settle()
    await scrollToAndSettle(400)
    expect(getChromeVisible()).toBe(true)
  })

  it('brings the chrome back when scrolling stops', async () => {
    await scrollToAndSettle(600)
    expect(getChromeVisible()).toBe(false)

    vi.advanceTimersByTime(IDLE_MS + 100)
    expect(getChromeVisible()).toBe(true)
  })

  it('does not let a settling bounce undo the hide', async () => {
    await scrollToAndSettle(900)
    expect(getChromeVisible()).toBe(false)

    settle()
    // Real scroll animations settle with a small upward correction. Under a
    // symmetric threshold this single event showed the chrome again ~16ms after
    // hiding, so the bar flickered rather than hiding.
    await scrollToAndSettle(877)
    expect(getChromeVisible()).toBe(false)

    // A deliberate scroll back up still brings it.
    await scrollToAndSettle(800)
    expect(getChromeVisible()).toBe(true)
  })

  it('ignores the scroll its own layout change provokes', async () => {
    await scrollToAndSettle(600)
    expect(getChromeVisible()).toBe(false)

    // Stopping brings it back, which re-adds the header and its padding; the
    // browser scroll-anchors that and reports a downward scroll. Without the
    // settle window this event re-hid the chrome, and hide/show ping-ponged.
    vi.advanceTimersByTime(IDLE_MS + 50)
    expect(getChromeVisible()).toBe(true)

    await scrollToAndSettle(640)
    expect(getChromeVisible()).toBe(true)

    // Once the window closes, scrolling hides it again as normal.
    settle()
    await scrollToAndSettle(900)
    expect(getChromeVisible()).toBe(false)
  })

  it('shows at the top regardless of direction', async () => {
    await scrollToAndSettle(600)
    expect(getChromeVisible()).toBe(false)

    settle()
    await scrollToAndSettle(0)
    expect(getChromeVisible()).toBe(true)
  })

  it('toggles on tap and holds through the idle window', async () => {
    await scrollToAndSettle(600)
    vi.advanceTimersByTime(IDLE_MS + 100)
    expect(getChromeVisible()).toBe(true)

    toggleChrome()
    expect(getChromeVisible()).toBe(false)

    // The idle timer must not undo the choice the user just made.
    vi.advanceTimersByTime(5000)
    expect(getChromeVisible()).toBe(false)

    toggleChrome()
    expect(getChromeVisible()).toBe(true)
    vi.advanceTimersByTime(5000)
    expect(getChromeVisible()).toBe(true)
  })

  it('hands control back to the scroll after a deliberate scroll', async () => {
    toggleChrome()
    expect(getChromeVisible()).toBe(false)

    // A pin is not a lock: scrolling on takes over again.
    settle()
    await scrollToAndSettle(600)
    expect(getChromeVisible()).toBe(false)
    vi.advanceTimersByTime(IDLE_MS + 100)
    expect(getChromeVisible()).toBe(true)
  })

  it('reveals at the bottom, then gets out of the way', async () => {
    await scrollToAndSettle(600)
    expect(getChromeVisible()).toBe(false)

    // scrollHeight 4000, innerHeight 800 -> bottom is y >= 3176.
    settle()
    await scrollToAndSettle(3200)
    expect(getChromeVisible()).toBe(true)

    vi.advanceTimersByTime(BOTTOM_AUTO_HIDE_MS + 500)
    expect(getChromeVisible()).toBe(false)
  })

  it('stays hidden at the bottom instead of cycling', async () => {
    settle()
    await scrollToAndSettle(3200)
    vi.advanceTimersByTime(BOTTOM_AUTO_HIDE_MS + 500)
    expect(getChromeVisible()).toBe(false)

    // Hiding shortens the document, so the browser clamps scrollY to stay at
    // the end and reports the clamp as a scroll event. Treating that as intent
    // re-showed the chrome and re-armed the timer: a two-second loop that never
    // settled. Resting at the bottom, only leaving it is scroll's decision.
    await scrollToAndSettle(3190)
    expect(getChromeVisible()).toBe(false)

    vi.advanceTimersByTime(BOTTOM_AUTO_HIDE_MS * 3)
    expect(getChromeVisible()).toBe(false)

    // Scrolling well clear of the end hands control back.
    await scrollToAndSettle(2000)
    expect(getChromeVisible()).toBe(true)
  })

  it('marks the document so CSS can reclaim the header spacing', () => {
    expect(document.documentElement.dataset.chrome).toBe('visible')

    toggleChrome()
    expect(document.documentElement.dataset.chrome).toBe('hidden')

    toggleChrome()
    expect(document.documentElement.dataset.chrome).toBe('visible')
  })

  it('clears the marker when the last consumer unmounts', () => {
    toggleChrome()
    expect(document.documentElement.dataset.chrome).toBe('hidden')

    unsubscribe?.()
    unsubscribe = null
    expect(document.documentElement.dataset.chrome).toBeUndefined()
    expect(getChromeVisible()).toBe(true)
  })
})
