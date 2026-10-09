import { describe, it, expect } from 'vitest'
import { desktopLinks, overflowLinks, TAB_HREFS, LEARN_LINKS, ADMIN_LINKS } from '../lib/nav'

// These assertions are about a product decision, not a rendering detail: Learn is the only
// surface that asks a question, grades it and explains it, and nothing may offer the learner a
// second way to approach the same mathematics.
//
// The one that used to fail here is the reference library. `/study` was a peer of Learn in
// every navigation surface, which made reading look like an alternative to being taught — and
// reading is not evidence, so a learner who found it had found a way around the only thing that
// counts. It is closed; its material lives inside `/learn` as a reference panel.
describe('navigation order', () => {
  const order = (links: { href: string }[]) => links.map(l => l.href)

  it('offers Learn to a guest, who can answer questions', () => {
    // Learn was previously logged-in-only, which meant a signed-out visitor
    // had no route to anything they could actually do.
    expect(desktopLinks(false).map(l => l.href)).toContain('/learn')
    expect(order(desktopLinks(false))).toEqual(['/learn', '/graph', '/domains', '/leaderboard', '/login'])
  })

  it('leads every surface with Learn', () => {
    expect(desktopLinks(true)[0].href).toBe('/learn')
    expect(desktopLinks(false)[0].href).toBe('/learn')
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

  it('keeps the orientation surfaces in the tab bar, since Learn has no browse', () => {
    // Domains took the slot Study held here. It is the honest replacement: it orients a learner
    // across the curriculum and links every row onward into /learn, where it is orientation
    // rather than a second curriculum.
    expect(TAB_HREFS.has('/domains')).toBe(true)
    expect(TAB_HREFS.has('/learn')).toBe(true)
  })

  it('exposes Learn as its own group for callers that want just it', () => {
    expect(LEARN_LINKS.map(l => l.href)).toEqual(['/learn'])
  })
})

// The Admin entry reached mobile only by typing a URL. The desktop rail on /profile has carried it
// since Admin V2, but on a phone /profile is one tap from the tab bar and no other surface offered
// it, so an administrator working from their phone had no route to administration at all. It belongs
// in the mobile overflow, and only for staff.
describe('the Admin entry reaches mobile', () => {
  it('appears in the compass overflow for staff', () => {
    const hrefs = overflowLinks(true, true).map(l => l.href)
    expect(hrefs).toContain('/admin')
  })

  it('is withheld from a learner and from a guest', () => {
    expect(overflowLinks(true, false).map(l => l.href)).not.toContain('/admin')
    expect(overflowLinks(false, false).map(l => l.href)).not.toContain('/admin')
  })

  it('does not become a bottom tab, so nothing appears twice on one screen', () => {
    expect(TAB_HREFS.has('/admin')).toBe(false)
    // …which is also why it must not collide with a tab href in the overflow.
    const hrefs = overflowLinks(true, true).map(l => l.href)
    expect(new Set(hrefs).size).toBe(hrefs.length)
  })

  it('leaves the learner overflow untouched', () => {
    // Adding staff must not perturb what a learner sees.
    expect(overflowLinks(true, false)).toEqual(overflowLinks(true))
  })

  it('is named the same way the desktop rail names it', () => {
    // The rail calls it "Admin"; a second wording for one destination is how two surfaces came to
    // disagree about it.
    expect(ADMIN_LINKS.map(l => l.label)).toEqual(['Admin'])
  })
})

// The decision, asserted on every function that builds a navigation surface. It is separated
// from the ordering tests above because a revert here is invisible in a snapshot and obvious in
// this list: the route still works, it just starts being linked again.
describe('the closed reference library stays closed', () => {
  it('is in no navigation surface at all', () => {
    for (const loggedIn of [true, false]) {
      expect(desktopLinks(loggedIn).map(l => l.href)).not.toContain('/study')
      expect(overflowLinks(loggedIn).map(l => l.href)).not.toContain('/study')
    }
    expect(TAB_HREFS.has('/study')).toBe(false)
    expect(LEARN_LINKS.map(l => l.href)).not.toContain('/study')
  })

  it('is labelled nowhere', () => {
    // Catches a link that was retargeted to /learn but kept its old wording, which is the
    // half-finished version of this change: the route is unused, the invitation is not.
    for (const loggedIn of [true, false]) {
      for (const link of desktopLinks(loggedIn)) {
        expect(link.label).not.toMatch(/study/i)
      }
    }
  })
})