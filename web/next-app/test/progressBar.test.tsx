import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProgressBar from '../components/ProgressBar'

describe('ProgressBar', () => {
  it('renders determinate state with aria values and clamped width', () => {
    const { container } = render(<ProgressBar answered={3} estimatedTotal={20} />)
    expect(screen.getByText('Progress bar')).toBeInTheDocument()
    expect(screen.getByText('3 / 20')).toBeInTheDocument()
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-valuenow')).toBe('3')
    expect(bar.getAttribute('aria-valuemax')).toBe('20')
    const fill = container.querySelector('.bg-mathua-blue') as HTMLElement
    expect(fill.style.width).toBe('15%')
  })

  it('clamps answered above the estimate at 100%', () => {
    const { container } = render(<ProgressBar answered={60} estimatedTotal={45} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-valuenow')).toBe('45')
    const fill = container.querySelector('.bg-mathua-blue') as HTMLElement
    expect(fill.style.width).toBe('100%')
  })

  it('renders indeterminate fallback when estimate is unknown', () => {
    render(<ProgressBar answered={1} estimatedTotal={0} />)
    const bar = screen.getByRole('progressbar', { name: 'Progress bar' })
    expect(bar.getAttribute('aria-busy')).toBe('true')
    expect(bar.hasAttribute('aria-valuenow')).toBe(false)
    expect(screen.getByText(/finding your frontier/i)).toBeInTheDocument()
  })
})
