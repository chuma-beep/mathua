import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProgressBar from '../components/ProgressBar'

describe('ProgressBar', () => {
  it('renders only the progress bar label and question number', () => {
    const { container } = render(<ProgressBar answered={3} coverDone={2} coverSize={20} />)
    expect(screen.getByText('Progress bar')).toBeInTheDocument()
    expect(screen.getByText('Question 3')).toBeInTheDocument()
    expect(screen.queryByText(/30–45 min/i)).toBeNull()
    expect(screen.queryByText(/pause/i)).toBeNull()
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

  it('renders indeterminate fallback with only the question number', () => {
    render(<ProgressBar answered={1} coverDone={0} coverSize={0} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-busy')).toBe('true')
    expect(bar.hasAttribute('aria-valuenow')).toBe(false)
    expect(screen.getByText('Question 1')).toBeInTheDocument()
    expect(screen.queryByText(/finding your frontier/i)).toBeNull()
  })
})
