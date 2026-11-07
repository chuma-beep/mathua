import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProfileStats from '../components/ProfileStats'
import type { Scores } from '../lib/api'

const scores = {
  lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1,
  current_streak: 2, level: 'Novice', xp_total: -40, xp_today: 0, daily_xp_goal: 10,
} as Scores

describe('ProfileStats', () => {
  it('shows a negative lifetime XP total and keeps daily XP non-negative', () => {
    render(<ProfileStats name="Tester" scores={scores} />)
    expect(screen.getByText('TOTAL XP')).toBeDefined()
    expect(screen.getByText('-40')).toBeDefined()
    expect(screen.getByText('0')).toBeDefined()
  })
})
