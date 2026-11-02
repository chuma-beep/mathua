import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import NotePage from '../app/note/page'

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))

describe("the Creator's Note", () => {
  it('offers support without implying a paywall', () => {
    render(<NotePage />)

    const support = screen.getByRole('link', { name: 'tip the project on Myhappr' })
    expect(support).toHaveAttribute('href', 'https://myhappr.com/wisddom')
    expect(support).toHaveAttribute('target', '_blank')
    expect(support).toHaveAttribute('rel', 'noreferrer')
    expect(screen.getByText(/keep Mathua free/)).toBeInTheDocument()
    expect(screen.getByText(/future app work such as a mobile app/)).toBeInTheDocument()
    expect(screen.getByText(/There is no paywall/)).toBeInTheDocument()
  })
})
