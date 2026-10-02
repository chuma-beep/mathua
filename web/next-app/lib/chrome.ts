'use client'

// One source of truth for whether the app chrome -- the sticky header and the
// mobile tab bar -- is showing. Both read it, so one scroll or one tap moves
// them together instead of letting each nav invent its own rules.
//
// Module-level rather than a React provider: pages own their own chrome (see
// e2e/chrome.spec.ts), so a provider would have to be threaded through every
// page to buy nothing. useSyncExternalStore keeps this safe under the static
// export -- the server snapshot is always "visible", so the first client render
// matches the pre-rendered HTML and there is no hydration mismatch.

export type ChromePin = 'shown' | 'hidden' | null

interface ChromeState {
  scrollHidden: boolean
  pinned: ChromePin
}

// Hiding is cheap to trigger by accident -- a downward flick while reading --
// so it takes very little movement. Bringing the chrome back is deliberately
// harder. Scroll is animated and inertial, and a settling animation routinely
// emits a small upward correction (measured: a single -23px bounce right after
// a 900px jump). With a symmetric threshold that bounce undid the hide about
// 16ms later, so the chrome flickered instead of hiding.
const HIDE_THRESHOLD = 12
const SHOW_THRESHOLD = 48
const TOP_OFFSET = 40
const IDLE_MS = 300
const BOTTOM_OFFSET = 24

// At the very bottom the bar is revealed and then gets out of the way again.
// There is nothing left to scroll, so a bar left sitting over the last screen
// of content is just an obstruction -- and it makes the gesture feel like it
// did not take.
const BOTTOM_AUTO_HIDE_MS = 2000

// How far the page must travel before a pinned choice is handed back to the
// scroll logic. Measured as distance from where the pin was set, not per-event
// delta: the app scrolls smoothly, and a tap landing mid-animation would
// otherwise be undone by the animation's trailing scroll events.
const PIN_RELEASE_DISTANCE = 48

// Changing the chrome changes the height of the document above the viewport:
// revealing the header and its --chrome-top padding puts ~149px back. The
// browser scroll-anchors that shift, which arrives as a scroll event of its
// own -- measured at +23px, comfortably past the hide threshold. Left alone the
// two feed each other: show -> anchor -> looks like a downward scroll -> hide
// -> collapse -> anchor -> show, and the bar oscillates while the page drifts.
//
// So scroll is inert for a moment after we move the chrome ourselves. The
// window is shorter than IDLE_MS on purpose: it only needs to cover the anchor
// event, which lands within a frame or two, and it must not outlive the idle
// timer that brings the chrome back on its own. lastY still advances during
// the window, so real movement inside it is not lost -- it surfaces as a larger
// delta once the window closes.
const SETTLE_MS = 250

// Desktop keeps its chrome. The tab bar does not exist there (it is lg:hidden),
// so hiding a lone header on a pointer device buys nothing.
const DESKTOP_QUERY = '(min-width: 1024px)'

let state: ChromeState = { scrollHidden: false, pinned: null }
const listeners = new Set<() => void>()

let lastY = 0
let pinAnchorY = 0
let settleUntil = 0
let atBottomMode = false
let ticking = false
let idleTimer: ReturnType<typeof setTimeout> | null = null
let bottomTimer: ReturnType<typeof setTimeout> | null = null

function isVisible(s: ChromeState = state): boolean {
  if (s.pinned === 'shown') return true
  if (s.pinned === 'hidden') return false
  return !s.scrollHidden
}

function clearTimers(): void {
  if (idleTimer !== null) clearTimeout(idleTimer)
  if (bottomTimer !== null) clearTimeout(bottomTimer)
  idleTimer = null
  bottomTimer = null
}

function syncDataset(): void {
  if (typeof document === 'undefined') return
  document.documentElement.dataset.chrome = isVisible() ? 'visible' : 'hidden'
}

function emit(): void {
  syncDataset()
  for (const l of listeners) l()
}

function set(patch: Partial<ChromeState>): void {
  const before = isVisible()
  state = { ...state, ...patch }
  // Only a real change in what is on screen invalidates the scroll signal.
  if (isVisible() !== before) settleUntil = Date.now() + SETTLE_MS
  emit()
}

function onScroll(): void {
  if (ticking) return
  ticking = true
  requestAnimationFrame(() => {
    ticking = false
    if (typeof window === 'undefined') return
    if (window.matchMedia(DESKTOP_QUERY).matches) return

    const y = window.scrollY
    const delta = y - lastY
    lastY = y

    // Scroll the chrome's own layout change provoked: keep the position
    // tracking, but do not let it drive visibility.
    if (Date.now() < settleUntil) return

    const atBottom =
      window.innerHeight + y >= document.documentElement.scrollHeight - BOTTOM_OFFSET

    // Resting on the bottom edge, the timers own the chrome. Collapsing the
    // header shortens the document, so the browser clamps scrollY to stay at the
    // end and reports the clamp as a scroll event -- which then re-armed the
    // bottom timer, which hid the chrome, which clamped again. A two-second
    // cycle that never settled. Leaving the bottom is the only thing scroll
    // decides here.
    if (atBottomMode) {
      lastY = y
      if (delta < -SHOW_THRESHOLD) atBottomMode = false
      else return
    }

    clearTimers()
    // Stopping always brings the chrome back. This is what makes the bar feel
    // like it is answering the scroll rather than fighting it.
    idleTimer = setTimeout(() => set({ scrollHidden: false }), IDLE_MS)

    // A pin is a deliberate choice, not a lock. Scrolling meaningfully far
    // hands control back to the scroll logic -- but only after real distance,
    // so a tap is not immediately reverted by scroll-animation residue.
    if (state.pinned !== null && Math.abs(y - pinAnchorY) > PIN_RELEASE_DISTANCE) {
      state.pinned = null
    }

    if (atBottom) {
      atBottomMode = true
      set({ scrollHidden: false })
      bottomTimer = setTimeout(() => set({ scrollHidden: true }), BOTTOM_AUTO_HIDE_MS)
    } else if (y <= TOP_OFFSET) {
      set({ scrollHidden: false })
    } else if (delta > HIDE_THRESHOLD) {
      set({ scrollHidden: true })
    } else if (delta < -SHOW_THRESHOLD) {
      set({ scrollHidden: false })
    }
  })
}

function attach(): void {
  lastY = window.scrollY
  settleUntil = 0
  syncDataset()
  window.addEventListener('scroll', onScroll, { passive: true })
}

function detach(): void {
  window.removeEventListener('scroll', onScroll)
  clearTimers()
  ticking = false
  state = { scrollHidden: false, pinned: null }
  settleUntil = 0
  pinAnchorY = 0
  atBottomMode = false
  if (typeof document !== 'undefined') delete document.documentElement.dataset.chrome
}

/** Tap target: pin the chrome to the opposite of what it is doing now. */
export function toggleChrome(): void {
  clearTimers()
  pinAnchorY = window.scrollY
  atBottomMode = false
  state = { scrollHidden: false, pinned: isVisible() ? 'hidden' : 'shown' }
  emit()
}

export function subscribeChrome(cb: () => void): () => void {
  listeners.add(cb)
  if (listeners.size === 1) attach()
  return () => {
    listeners.delete(cb)
    if (listeners.size === 0) detach()
  }
}

export function getChromeVisible(): boolean {
  return isVisible()
}

// The pre-rendered HTML always has the chrome showing, so the first client
// render has to agree with it.
export function getChromeServerVisible(): boolean {
  return true
}
