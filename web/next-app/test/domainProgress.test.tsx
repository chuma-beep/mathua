import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import DomainProgress from '../components/DomainProgress'

describe('DomainProgress counts', () => {
  it('shows completed and locked words with empty progress', () => {
    render(<DomainProgress progress={{}} />)
    expect(screen.getAllByText(/0 done/).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/locked/).length).toBeGreaterThan(0)
  })

  it('counts a completed concept separately from mastered', () => {
    const { container } = render(
      <DomainProgress progress={{ 'arith.add.single': { status: 'learning', completed: true, streak: 1 } }} />
    )
    expect(container.textContent).toMatch(/1 done/)
    expect(container.textContent).not.toMatch(/1\/\d+ · [1-9]/)
  })
})
