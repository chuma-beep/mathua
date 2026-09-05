import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import WeeklyChart from '../components/WeeklyChart'
import type { DailyActivity } from '../lib/api'

function day(offset: number, questions: number, correct: number): DailyActivity {
  const d = new Date()
  d.setDate(d.getDate() - offset)
  return { date: d.toISOString().slice(0, 10), questions, correct, concepts: [] }
}

describe('WeeklyChart', () => {
  it('aggregates last 7 days with accessible summary', () => {
    render(<WeeklyChart data={[day(6, 4, 3), day(0, 6, 6)]} />)
    expect(screen.getByRole('img')).toHaveAttribute(
      'aria-label',
      expect.stringContaining('10 questions in the last 7 days'),
    )
    expect(screen.getByText(/10 questions · 90% correct/)).toBeInTheDocument()
    expect(screen.getByText('Daily breakdown')).toBeInTheDocument()
  })

  it('shows empty state with Study CTA when no activity', () => {
    render(<WeeklyChart data={[]} />)
    expect(screen.getByText(/no activity this week yet/i)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /open study/i })).toHaveAttribute('href', '/study')
  })
})
