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
// It used to send a signed-out visitor to /profile, an empty dashboard with a
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

  it('sends "Get started" to Login', async () => {
    render(<HomePage />)
    const cta = await screen.findByRole('button', { name: 'Get started' })
    cta.click()
    await waitFor(() => expect(pushMock).toHaveBeenCalledWith('/login'))
    expect(pushMock).not.toHaveBeenCalledWith('/profile')
  })

  it('does not mint a guest identity when opening Login', async () => {
    render(<HomePage />)
    ;(await screen.findByRole('button', { name: 'Get started' })).click()
    await waitFor(() => expect(pushMock).toHaveBeenCalledWith('/login'))
    // The anchor that grades the first answer is keyed by student, and the
    // practice endpoint only resolves one from a bearer token, so the guest
    // id and token have to exist before /learn fetches anything.
    expect(ensureGuestId).not.toHaveBeenCalled()
    expect(ensureGuestToken).not.toHaveBeenCalled()
  })

  it('still leaves a signed-in visitor on the landing page', async () => {
    localStorage.setItem('mathua_token', 't')
    render(<HomePage />)
    await waitFor(() => expect(replaceMock).toHaveBeenCalledWith('/profile'))
  })

  it('keeps the two hero actions on separate, non-conflicting backgrounds', async () => {
    render(<HomePage />)
    const primary = await screen.findByRole('button', { name: 'Get started' })
    const ghost = await screen.findByRole('link', { name: 'View on GitHub' })
    // The ghost button was once the primary button plus overrides, so both
    // `bg-mathua-blue` and `bg-transparent` (and both hover backgrounds) sat on one
    // element. Same property, same specificity: Tailwind's emitted order decided the
    // winner, and the GitHub action hovered to blue over the artwork. The ghost now
    // carries its own transparent pair below lg and themed `lg:` overrides for the
    // desktop panel, compared as exact class tokens, because the legitimate
    // `lg:hover:bg-mathua-blue-faint` contains `bg-mathua-blue` as a substring.
    const primaryTokens = primary.className.split(/\s+/)
    const ghostTokens = ghost.className.split(/\s+/)
    expect(primaryTokens).toContain('bg-mathua-blue')
    expect(primaryTokens).toContain('hover:bg-mathua-blue-hover')
    expect(ghostTokens).toContain('bg-transparent')
    expect(ghostTokens).not.toContain('bg-mathua-blue')
    expect(ghostTokens).not.toContain('hover:bg-mathua-blue-hover')
  })

  it('lists Learn in the landing header and does not advertise the closed library', async () => {
    render(<HomePage />)
    const header = await screen.findByRole('banner')
    const learn = header.querySelector('a[href="/learn"]')
    expect(learn).not.toBeNull()
    // The landing hardcodes its own header rather than inheriting lib/nav.ts, so it is the one
    // place the product's own navigation decision could quietly disappear, in either
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
