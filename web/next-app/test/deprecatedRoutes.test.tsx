import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render } from '@testing-library/react'
import ConceptPage from '../app/concept/page'
import SessionPage from '../app/session/page'

const replaceMock = vi.fn()
let params = new URLSearchParams()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: replaceMock, push: vi.fn() }),
  useSearchParams: () => params,
}))

// /concept and /session are client redirect stubs with no internal caller: no
// nav entry, not in the indexable routes, no link in the app, and — until this
// file — no test. The only thing keeping them is external links, which cannot
// be measured from here.
//
// So this test does not decide whether they stay. It pins what they currently
// do, so that deleting them is a deliberate, visible change rather than a
// quiet removal, and so that whoever makes the call knows exactly what breaks
// for an old link.
describe('deprecated redirect stubs', () => {
  beforeEach(() => {
    replaceMock.mockClear()
    params = new URLSearchParams()
  })

  it('/concept?id=X forwards to the Study page for that concept', () => {
    params = new URLSearchParams('id=frac.add.diff')
    render(<ConceptPage />)
    expect(replaceMock).toHaveBeenCalledWith('/study?concept=frac.add.diff')
  })

  it('/concept with no id forwards to the Study landing', () => {
    render(<ConceptPage />)
    expect(replaceMock).toHaveBeenCalledWith('/study')
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
    expect(replaceMock).toHaveBeenCalledWith('/study?concept=a%20b%2Fc')
  })
})