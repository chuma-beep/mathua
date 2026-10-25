import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import ContactPage from '../app/contact/page'

vi.mock('next/navigation', () => ({
  usePathname: () => '/contact',
  useRouter: () => ({ push: vi.fn() }),
}))
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))

// The page's value is that it routes three different problems to three different
// places. If a later edit collapses them into "email me", the account-problem
// channel becomes a public issue tracker and the question channel loses the
// captured context that makes a bug reproducible at all.
describe('the contact page', () => {
  it('sends a wrong question to the report control, not to a mailbox', () => {
    render(<ContactPage />)
    expect(screen.getByRole('heading', { name: 'A question is wrong' })).toBeInTheDocument()
    expect(screen.getByText(/Report a problem/)).toBeInTheDocument()
    expect(screen.getByText(/Learn, a quiz, a review, a diagnostic/)).toBeInTheDocument()
    expect(screen.getByText(/In History, open the concept/)).toBeInTheDocument()
  })

  it('offers the box on the page itself, not only an external channel', () => {
    render(<ContactPage />)
    expect(screen.getByRole('heading', { name: 'Leave a message' })).toBeInTheDocument()
    expect(screen.getByLabelText('Your message')).toBeInTheDocument()
  })

  it('publishes the two channels the owner reads', () => {
    render(<ContactPage />)

    const x = screen.getByRole('link', { name: '@wisdomsglow' })
    expect(x).toHaveAttribute('href', 'https://x.com/wisdomsglow')
    expect(x).toHaveAttribute('target', '_blank')
    expect(x).toHaveAttribute('rel', 'noreferrer')

    const issues = screen.getByRole('link', { name: /github\.com\/chuma-beep\/mathua\/issues/ })
    expect(issues).toHaveAttribute('href', 'https://github.com/chuma-beep/mathua/issues')
  })

  it('never asks for a password and publishes no address to harvest', () => {
    const { container } = render(<ContactPage />)

    expect(container.querySelector('a[href^="mailto:"]')).toBeNull()
    expect(container.textContent).not.toMatch(/[\w.+-]+@[\w-]+\.[\w.]+/)
    expect(container.textContent).toMatch(/Never send your password/)
  })
})
