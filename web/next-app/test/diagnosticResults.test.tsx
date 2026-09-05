import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import DiagnosticResults from '../components/DiagnosticResults'
import type { GoalPlanRes } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, getLessons: () => Promise.resolve({ lessons: {} }) }
})

vi.mock('../lib/auth', () => ({
  getUserInfo: () => null,
}))

const base: GoalPlanRes = {
  readiness: 0.7,
  total_tested: 20,
  correct_count: 14,
  weak_areas: {},
  strong_areas: {},
}

describe('DiagnosticResults', () => {
  it('renders placement banner with frontier, course, and estimates', () => {
    render(
      <DiagnosticResults
        plan={{
          ...base,
          frontier_label: 'Adding Fractions',
          placement_course_id: '5',
          completion_estimates: { '150': '≈5 days', '300': '≈10 days' },
        }}
      />,
    )
    expect(screen.getByText('Start here: Adding Fractions')).toBeInTheDocument()
    expect(screen.getByText('Suggested course: 5')).toBeInTheDocument()
    expect(screen.getByText(/150 XP → ≈5 days/)).toBeInTheDocument()
    expect(screen.getByText(/no timeout/i)).toBeInTheDocument()
  })

  it('renders provisional frontier and conditional topics', () => {
    render(
      <DiagnosticResults
        plan={{ ...base, frontier_label: 'X', frontier_conditional: true, conditionally_completed: ['Y'] }}
      />,
    )
    expect(screen.getAllByText(/provisional/i).length).toBeGreaterThan(0)
    expect(screen.getByText(/fall back on struggle/i)).toBeInTheDocument()
  })

  it('omits the banner when no placement data is present', () => {
    render(<DiagnosticResults plan={base} />)
    expect(screen.queryByText(/Start here:/)).toBeNull()
  })
})
