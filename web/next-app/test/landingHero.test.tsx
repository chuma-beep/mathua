import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import HomePage from '../app/page'
import { ensureGuestId, ensureGuestToken } from '../lib/auth'

const pushMock = vi.fn()
const replaceMock = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock, replace: replaceMock }),
  usePathname: () => '/',
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return {
    ...mod,
    ensureGuestId: vi.fn(),
    ensureGuestToken: vi.fn(async () => 'guest-token'),
  }
})

// The hero's "Get started" is the most prominent action on the marketing page.
// It used to send a signed-out visitor to /profile — an empty dashboard with a
// zeroed heatmap and no question anywhere on it, which is a poor first
// experience and the opposite of Learn being the primary destination.
describe('landing hero', () => {
  beforeEach(() => {
    localStorage.clear()
    sessionStorage.clear()
    pushMock.mockClear()
    replaceMock.mockClear()
    vi.mocked(ensureGuestId).mockClear()
  })

  it('sends "Get started" to Learn, not to an empty Profile', async () => {
    render(<HomePage />)
    const cta = await screen.findByRole('button', { name: 'Get started' })
    cta.click()
    await waitFor(() => expect(pushMock).toHaveBeenCalledWith('/learn'))
    expect(pushMock).not.toHaveBeenCalledWith('/profile')
  })

  it('establishes a guest identity before leaving, so the first answer is gradeable', async () => {
    render(<HomePage />)
    ;(await screen.findByRole('button', { name: 'Get started' })).click()
    await waitFor(() => expect(pushMock).toHaveBeenCalledWith('/learn'))
    // The anchor that grades the first answer is keyed by student, and the
    // practice endpoint only resolves one from a bearer token — so the guest
    // id and token have to exist before /learn fetches anything.
    expect(ensureGuestId).toHaveBeenCalled()
    expect(ensureGuestToken).toHaveBeenCalled()
  })

  it('still leaves a signed-in visitor on the landing page', async () => {
    localStorage.setItem('mathua_token', 't')
    render(<HomePage />)
    await waitFor(() => expect(replaceMock).toHaveBeenCalledWith('/profile'))
  })

  it('lists Learn in the landing header and does not advertise the closed library', async () => {
    render(<HomePage />)
    const header = await screen.findByRole('banner')
    const learn = header.querySelector('a[href="/learn"]')
    expect(learn).not.toBeNull()
    // The landing hardcodes its own header rather than inheriting lib/nav.ts, so it is the one
    // place the product's own navigation decision could quietly disappear — in either
    // direction. Learn leads, and the closed reference library is not offered at all.
    const hrefs = [...header.querySelectorAll('a')].map((a) => a.getAttribute('href'))
    // The brand mark is the first anchor, so the ordering claim is about the nav entries:
    // Learn leads them.
    expect(hrefs.indexOf('/learn')).toBeGreaterThan(-1)
    expect(hrefs.indexOf('/learn')).toBeLessThan(hrefs.indexOf('/how-it-works'))
    expect(hrefs).not.toContain('/study')
    expect([...header.querySelectorAll('a')].some((a) => /study/i.test(a.textContent ?? ''))).toBe(false)
  })
})
