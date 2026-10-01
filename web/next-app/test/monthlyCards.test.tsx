import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import MonthlyCards from '../components/MonthlyCards'

function iso(year: number, month: number, day: number): string {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
}

function entry(year: number, month: number, day: number, questions = 5) {
  return { date: iso(year, month, day), questions, correct: 4, concepts: ['a'] }
}

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

describe('MonthlyCards strip', () => {
  it('renders active months plus the current anchor, oldest first', () => {
    const now = new Date()
    const old = new Date(now.getFullYear(), now.getMonth() - 4, 1)
    const { container } = render(
      <MonthlyCards
        data={[
          entry(old.getFullYear(), old.getMonth(), 3),
          entry(old.getFullYear(), old.getMonth(), 4, 2),
        ]}
      />
    )
    const strip = container.querySelector('[data-testid="month-strip"]')
    expect(strip).not.toBeNull()
    expect(strip?.className).toMatch(/overflow-x-auto/)
    const headings = Array.from(strip?.querySelectorAll('p') ?? []).map((p) => p.textContent ?? '')
    expect(headings).toHaveLength(2)
    // Oldest left, current anchor right.
    expect(headings[0]).toContain(MONTHS[old.getMonth()])
    expect(headings[1]).toContain(MONTHS[now.getMonth()])
  })

  it('renders a single current-month card for empty data', () => {
    const { container } = render(<MonthlyCards data={[]} />)
    const strip = container.querySelector('[data-testid="month-strip"]')
    expect(strip?.querySelectorAll('p').length).toBe(1)
  })

  it('shows every grid without toggles', () => {
    const now = new Date()
    render(
      <MonthlyCards
        data={[
          entry(now.getFullYear() - 1, now.getMonth(), 3),
          entry(now.getFullYear(), now.getMonth(), 5),
        ]}
      />
    )
    expect(screen.queryByRole('button')).toBeNull()
    // Both months' grids visible: tooltips from each year present.
    expect(document.querySelectorAll('div[title*="q,"]').length).toBeGreaterThan(0)
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
