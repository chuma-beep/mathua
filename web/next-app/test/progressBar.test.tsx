import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProgressBar from '../components/ProgressBar'

function fill(container: HTMLElement) {
  return container.querySelector('.bg-mathua-blue') as HTMLElement
}

describe('ProgressBar', () => {
  it('renders only the progress bar label and question number', () => {
    const { container } = render(<ProgressBar answered={2} budget={20} />)
    expect(screen.getByText('Progress bar')).toBeInTheDocument()
    expect(screen.getByText('Question 3')).toBeInTheDocument()
    expect(screen.queryByText(/30–45 min/i)).toBeNull()
    expect(screen.queryByText(/pause/i)).toBeNull()
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-valuenow')).toBe('2')
    expect(bar.getAttribute('aria-valuemax')).toBe('20')
    expect(fill(container).style.width).toBe('10%')
  })

  it('clamps answers past the budget at 100% (never overflows)', () => {
    const { container } = render(<ProgressBar answered={60} budget={45} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-valuenow')).toBe('45')
    expect(fill(container).style.width).toBe('100%')
  })

  it('renders indeterminate fallback with only the question number', () => {
    render(<ProgressBar answered={1} budget={0} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-busy')).toBe('true')
    expect(bar.hasAttribute('aria-valuenow')).toBe(false)
    expect(screen.getByText('Question 2')).toBeInTheDocument()
    expect(screen.queryByText(/finding your frontier/i)).toBeNull()
  })

  it('moves the fill on the first answer of a full 45-question budget', () => {
    const first = render(<ProgressBar answered={0} budget={45} />)
    expect(fill(first.container).style.width).toBe('0%')
    first.unmount()

    const second = render(<ProgressBar answered={1} budget={45} />)
    expect(fill(second.container).style.width).toBe('2.2222222222222223%')
  })
})
