import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import HomePage from '../app/page'

const { push, ensureGuestToken, ensureGuestId } = vi.hoisted(() => ({
  push: vi.fn(),
  ensureGuestToken: vi.fn(),
  ensureGuestId: vi.fn(),
}))

vi.mock('next/navigation', () => ({ useRouter: () => ({ push }) }))
vi.mock('../lib/auth', () => ({ ensureGuestId, ensureGuestToken }))
vi.mock('../hooks/useTheme', () => ({ useTheme: () => ({ mounted: true }) }))
vi.mock('../hooks/useHomeRedirect', () => ({ useHomeRedirect: () => true }))
vi.mock('../hooks/useScrollReveal', () => ({ useScrollReveal: () => {} }))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
// Only the hero matters here; the rest of the landing is copy.
vi.mock('../app/home/sections', () => ({
  HeroSection: ({ onGetStarted }: { onGetStarted: () => void }) => (
    <button type="button" onClick={onGetStarted}>Get started</button>
  ),
  TrustSection: () => null,
  FeaturesSection: () => null,
  StatesSection: () => null,
  CoverageSection: () => null,
  RankingSection: () => null,
  FaqSection: () => null,
}))

// Reported as "Get started page won't load on desktop". The page was fine: the
// click waited for the guest token before navigating, so a stalled token request
// (a filtered or blocked API origin) kept the button dead for the request's full
// ten-second timeout. Navigation must not depend on it.
describe('Get started', () => {
  it('opens the learn page without waiting for the guest token', () => {
    ensureGuestToken.mockReturnValue(new Promise(() => {})) // never settles
    render(<HomePage />)

    fireEvent.click(screen.getByRole('button', { name: 'Get started' }))

    expect(push).toHaveBeenCalledWith('/learn')
    expect(ensureGuestId).toHaveBeenCalled()
    expect(ensureGuestToken).toHaveBeenCalled()
  })
})
