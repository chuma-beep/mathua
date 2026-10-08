import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render } from '@testing-library/react'
import ConceptPage from '../app/concept/page'
import SessionPage from '../app/session/page'
import StudyPage from '../app/study/page'

const replaceMock = vi.fn()
let params = new URLSearchParams()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: replaceMock, push: vi.fn() }),
  useSearchParams: () => params,
}))

// /concept, /session and /study are client redirect stubs with no internal caller: no nav
// entry, not in the indexable routes, no link in the app, and — until this file — no test.
// The only thing keeping them is external links and old bookmarks, which cannot be measured
// from here.
//
// So this test does not decide whether they stay. It pins what they currently do, so that
// deleting them is a deliberate, visible change rather than a quiet removal, and so that
// whoever makes the call knows exactly what breaks for an old link.
//
// `/study` is the one that matters most: it was a full destination — a browsable reference
// library, linked from every header, tab bar, sidebar and footer — until the learning path was
// closed. It forwards now, and it forwards to the loop rather than 404ing, because a bookmark
// to a concept's reference is a bookmark to that concept.
describe('deprecated redirect stubs', () => {
  beforeEach(() => {
    replaceMock.mockClear()
    params = new URLSearchParams()
  })

  it('/concept?id=X forwards to the Learn loop for that concept', () => {
    params = new URLSearchParams('id=frac.add.diff')
    render(<ConceptPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn?concept=frac.add.diff')
  })

  it('/concept with no id forwards to the Learn entry, which loads the Next head', () => {
    render(<ConceptPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn')
  })

  it('/study?concept=X forwards to the Learn loop for that concept', () => {
    params = new URLSearchParams('concept=frac.add.diff')
    render(<StudyPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn?concept=frac.add.diff')
  })

  it('/study with no concept forwards to the Learn entry', () => {
    render(<StudyPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn')
  })

  it('/study drops ?domain= rather than guessing where a domain index went', () => {
    // A domain list is orientation, and orientation is `/domains`. But the honest behaviour is
    // to forward to the loop rather than to pick a destination on the learner's behalf, and a
    // test that pinned `/domains` would be pinning a guess.
    params = new URLSearchParams('domain=fractions&lesson=Adding+up')
    render(<StudyPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn')
  })

  it('/study ignores ?from=, which only ever existed to come back from it', () => {
    params = new URLSearchParams('concept=arith.add.single&from=arith.add.single')
    render(<StudyPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn?concept=arith.add.single')
  })

  it('/session?concept=X forwards to the Learn loop for that concept', () => {
    params = new URLSearchParams('concept=frac.add.diff')
    render(<SessionPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn?concept=frac.add.diff')
  })

  it('/session with no concept forwards to the Learn entry, which loads the Next head', () => {
    render(<SessionPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn')
  })

  it('escapes an id rather than splicing it into the query string', () => {
    // An id with a space and a slash must not be able to reshape the query.
    params = new URLSearchParams('id=a b/c')
    render(<ConceptPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn?concept=a%20b%2Fc')
  })

  it('escapes a study concept the same way', () => {
    params = new URLSearchParams('concept=a b/c')
    render(<StudyPage />)
    expect(replaceMock).toHaveBeenCalledWith('/learn?concept=a%20b%2Fc')
  })

  it('never forwards to itself', () => {
    // A forwarder that can reach its own route is a loop, not a redirect. This is cheap to
    // assert and impossible to notice by reading the page.
    for (const q of ['', 'concept=a', 'domain=b', 'lesson=c', 'from=d', 'concept=a&from=a']) {
      replaceMock.mockClear()
      params = new URLSearchParams(q)
      render(<StudyPage />)
      for (const [href] of replaceMock.mock.calls) {
        expect(String(href).startsWith('/study')).toBe(false)
      }
    }
  })
})