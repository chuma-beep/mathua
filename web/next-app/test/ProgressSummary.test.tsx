import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProgressSummary from '../components/ProgressSummary'
import type { Scores } from '../lib/api'

const SCORES: Scores = {
  lifetime_points: 100,
  weekly_score: 0,
  speed_bonus: 0,
  concepts_mastered: 0,
  current_streak: 0,
  level: 'Novice',
  xp_total: 0,
  xp_today: 0,
  daily_xp_goal: 10,
}

describe('ProgressSummary', () => {
  it('shows the four headline numbers', () => {
    render(<ProgressSummary scores={SCORES} />)
    // Rendered by /graph and /goals, so these are the labels a learner reads on two
    // different pages.
    for (const label of ['Concepts mastered', 'Day streak', 'Level', 'Weekly score']) {
      expect(screen.getByText(label)).toBeInTheDocument()
    }
    expect(screen.getByText('Novice')).toBeInTheDocument()
  })

  it('draws no bordered box rounded', () => {
    // These tiles carried `rounded-lg` while the rest of the app — including
    // `components/ui/button.tsx`, whose comment reads "No rounded-*" — did not. They are
    // the most visible boxes on /graph, so the inconsistency read as the graph page
    // being built to a different spec.
    //
    // Same rule as the graph guard: a bordered box may not be rounded, which leaves
    // genuinely circular things alone without enumerating them.
    const { container } = render(
      <ProgressSummary
        scores={SCORES}
        weakByDomain={{ arithmetic: [{ id: 'frac.add.word', label: 'Add fractions' }] }}
      />,
    )
    const offenders = Array.from(container.querySelectorAll<HTMLElement>('*'))
      .map(el => {
        const cls = el.getAttribute('class') ?? ''
        const hasBorder = /\bborder/.test(cls)
        const rounded = /(?:^|\s)rounded-(?!none\b)/.test(cls) || parseFloat(getComputedStyle(el).borderRadius) > 0
        return hasBorder && rounded ? `${el.tagName.toLowerCase()}.${cls}` : null
      })
      .filter((v): v is string => v !== null)

    expect(offenders, `rounded bordered boxes: ${JSON.stringify(offenders)}`).toEqual([])
  })
})