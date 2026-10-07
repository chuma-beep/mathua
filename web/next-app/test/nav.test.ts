import { describe, it, expect } from 'vitest'
import { desktopLinks, overflowLinks, TAB_HREFS, LEARN_LINKS } from '../lib/nav'

// These assertions are about a product decision, not a rendering detail: Learn
// is the only surface that asks a question, grades it and explains it, and
// reference must not outrank doing in any navigation surface. A reorder that
// looks cosmetic here silently inverts the product's own hierarchy.
describe('navigation order', () => {
  const order = (links: { href: string }[]) => links.map(l => l.href)

  it('puts Learn before Study everywhere', () => {
    const links = desktopLinks(true)
    const learn = links.findIndex(l => l.href === '/learn')
    const study = links.findIndex(l => l.href === '/study')
    expect(learn).toBeGreaterThanOrEqual(0)
    expect(study).toBeGreaterThanOrEqual(0)
    expect(learn).toBeLessThan(study)
  })

  it('leads the header row with Learn', () => {
    expect(desktopLinks(true)[0].href).toBe('/learn')
    expect(desktopLinks(false)[0].href).toBe('/learn')
  })

  it('offers Learn to a guest, who can answer questions', () => {
    // Learn was previously logged-in-only, which meant a signed-out visitor
    // had no route to anything they could actually do.
    expect(desktopLinks(false).map(l => l.href)).toContain('/learn')
    expect(order(desktopLinks(false))).toEqual(['/learn', '/study', '/graph', '/domains', '/leaderboard', '/login'])
  })

  it('shows Review beside Learn for a signed-in learner', () => {
    // Review grades, awards XP and reschedules memory: it is a Learn sub-mode
    // and had no navigation entry at all, reachable only from Profile.
    const links = desktopLinks(true).map(l => l.href)
    const learn = links.indexOf('/learn')
    const review = links.indexOf('/review')
    expect(review).toBeGreaterThan(learn)
  })

  it('hides Review from a guest, whose review request would 401', () => {
    expect(desktopLinks(false).map(l => l.href)).not.toContain('/review')
  })

  it('keeps the bottom tab bar at five so Review reaches the compass', () => {
    expect(TAB_HREFS.size).toBe(5)
    expect(TAB_HREFS.has('/review')).toBe(false)
    // …and therefore does show up in the mobile overflow.
    expect(overflowLinks(true).map(l => l.href)).toContain('/review')
  })

  it('never repeats a destination in the compass overflow', () => {
    for (const loggedIn of [true, false]) {
      const shown = overflowLinks(loggedIn)
      expect(new Set(shown.map(l => l.href)).size).toBe(shown.length)
      for (const link of shown) {
        expect(TAB_HREFS.has(link.href)).toBe(false)
      }
    }
  })

  it('gives every link a label', () => {
    for (const loggedIn of [true, false]) {
      for (const link of desktopLinks(loggedIn)) {
        expect(link.label.length).toBeGreaterThan(0)
      }
    }
  })

  it('keeps Study in the tab bar: it is the discovery surface', () => {
    // Learn has no browse, no search and no domain list, so removing Study
    // from the tabs would leave no way to look a concept up.
    expect(TAB_HREFS.has('/study')).toBe(true)
    expect(TAB_HREFS.has('/learn')).toBe(true)
  })

  it('exposes Learn as its own group for callers that want just it', () => {
    expect(LEARN_LINKS.map(l => l.href)).toEqual(['/learn'])
  })
})
