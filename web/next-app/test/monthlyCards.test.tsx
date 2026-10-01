import { describe, it, expect } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import MonthlyCards from '../components/MonthlyCards'

function iso(year: number, month: number, day: number): string {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
}

function entry(year: number, month: number, day: number, questions = 5) {
  return { date: iso(year, month, day), questions, correct: 4, concepts: ['a'] }
}

describe('MonthlyCards', () => {
  it('renders only active months plus the current anchor, newest first', () => {
    const now = new Date()
    const old = new Date(now.getFullYear(), now.getMonth() - 4, 1)
    render(
      <MonthlyCards
        data={[
          entry(old.getFullYear(), old.getMonth(), 3),
          entry(old.getFullYear(), old.getMonth(), 4, 2),
        ]}
      />
    )
    const buttons = screen.getAllByRole('button')
    expect(buttons).toHaveLength(2)
    // Newest (current anchor) first.
    expect(buttons[0].getAttribute('aria-label')).toContain(String(now.getFullYear()))
  })

  it('renders a single current-month card for empty data', () => {
    render(<MonthlyCards data={[]} />)
    expect(screen.getAllByRole('button')).toHaveLength(1)
  })

  it('expands same-month-different-year independently', () => {
    const now = new Date()
    const month = now.getMonth()
    render(
      <MonthlyCards
        data={[
          entry(now.getFullYear() - 1, month, 3),
          entry(now.getFullYear(), month, 5),
        ]}
      />
    )
    const buttons = screen.getAllByRole('button')
    expect(buttons).toHaveLength(2)
    // Newest expanded by default; exactly one open.
    expect(buttons.filter((b) => b.getAttribute('aria-expanded') === 'true')).toHaveLength(1)
    // Opening the older year closes the newer one — never both.
    fireEvent.click(buttons[1])
    expect(buttons.filter((b) => b.getAttribute('aria-expanded') === 'true')).toHaveLength(1)
    expect(buttons[1].getAttribute('aria-expanded')).toBe('true')
  })

  it('shows heat cells without day numbers, tooltips intact', () => {
    const now = new Date()
    const { container } = render(
      <MonthlyCards data={[entry(now.getFullYear(), now.getMonth(), 4, 5)]} />
    )
    const cells = Array.from(container.querySelectorAll('div[title]')).filter((d) =>
      /q, \d+%/.test(d.getAttribute('title') ?? '')
    )
    expect(cells.length).toBeGreaterThan(0)
    for (const c of cells) expect(c.textContent).toBe('')
    expect(container.textContent).toMatch(/5 questions/)
  })

  it('trims day headers to Mo/We/Fr', () => {
    const { container } = render(<MonthlyCards data={[]} />)
    const text = container.textContent ?? ''
    for (const dh of ['Mo', 'We', 'Fr']) expect(text).toContain(dh)
    for (const dh of ['Su', 'Tu', 'Th', 'Sa']) expect(text).not.toContain(dh)
  })
})
