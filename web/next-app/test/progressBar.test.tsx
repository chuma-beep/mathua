import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProgressBar from '../components/ProgressBar'

describe('ProgressBar', () => {
  it('renders MA-style copy without promising a total', () => {
    const { container } = render(<ProgressBar answered={3} coverDone={2} coverSize={20} />)
    expect(screen.getByText('Progress bar')).toBeInTheDocument()
    expect(screen.getByText(/Question 3 · ~30–45 min · can pause/i)).toBeInTheDocument()
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-valuenow')).toBe('2')
    expect(bar.getAttribute('aria-valuemax')).toBe('20')
    const fill = container.querySelector('.bg-mathua-blue') as HTMLElement
    expect(fill.style.width).toBe('10%')
  })

  it('clamps cover done at 100% (never regresses, never overflows)', () => {
    const { container } = render(<ProgressBar answered={60} coverDone={99} coverSize={45} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-valuenow')).toBe('45')
    const fill = container.querySelector('.bg-mathua-blue') as HTMLElement
    expect(fill.style.width).toBe('100%')
  })

  it('renders indeterminate fallback when cover is unknown', () => {
    render(<ProgressBar answered={1} coverDone={0} coverSize={0} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-busy')).toBe('true')
    expect(bar.hasAttribute('aria-valuenow')).toBe(false)
    expect(screen.getByText(/finding your frontier/i)).toBeInTheDocument()
  })
})
