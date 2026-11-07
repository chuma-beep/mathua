import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProfileHeader from '../components/ProfileHeader'
import type { CatalogueConcept } from '../lib/progress'
import type { ConceptProgress, Scores } from '../lib/api'

// The learner header is a few short sentences in one card: identity and level, the reach
// headline with today's numbers, and the mastery breakdown. It replaced a bordered stat strip
// (avatar, LV tile, XP/TOTAL/STREAK/WEEKLY tiles, a progress bar) stacked on a position card
// (big number, second bar, 3-column grid). The guarded properties: the reach sentence, mastery
// named once, and the region the rest of the page orders against.

const catalogue: CatalogueConcept[] = [
  { id: 'a', domain: 'arithmetic', prerequisites: [] },
  { id: 'b', domain: 'arithmetic', prerequisites: ['a'] },
  { id: 'c', domain: 'arithmetic', prerequisites: ['a'] },
]

const scores = {
  lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1,
  current_streak: 2, level: 'Novice', xp_total: 120, xp_today: 6, daily_xp_goal: 10,
} as Scores

const progress = { a: { status: 'MASTERED', streak: 3 } } satisfies Record<string, ConceptProgress>

describe('ProfileHeader', () => {
  it('shows identity, reach and mastery in a few short lines', () => {
    render(
      <ProfileHeader name="Tester" scores={scores} catalogue={catalogue} progress={progress} />,
    )
    expect(screen.getByText(/Tester · LV 01 Novice/)).toBeTruthy()
    const position = screen.getByRole('region', { name: 'Where you are' })
    expect(position.textContent).toMatch(/1 of 3 you can reach/)
    expect(position.textContent).toMatch(/120 XP all-time/)
    expect(position.textContent).toMatch(/Mastered 1 · Unlocked 3 · In progress 0/)
  })

  it('shows a negative lifetime XP balance without clamping it', () => {
    render(
      <ProfileHeader
        name="Tester"
        scores={{ ...scores, xp_total: -40 }}
        catalogue={catalogue}
        progress={progress}
      />,
    )
    expect(screen.getByText(/-40 XP all-time/)).toBeTruthy()
  })

  it('names mastery once, not as a separate tile', () => {
    render(
      <ProfileHeader name="Tester" scores={scores} catalogue={catalogue} progress={progress} />,
    )
    // No stat strip anymore: the only "Mastered" is the breakdown sentence.
    expect(screen.queryByText('MASTERED')).toBeNull()
    // The old tile label is gone, but the cumulative XP value stays on the reach line.
    expect(screen.queryByText('TOTAL XP')).toBeNull()
    expect(screen.getByText(/120 XP all-time/)).toBeTruthy()
    expect(screen.getByText(/Mastered 1 · Unlocked/)).toBeTruthy()
  })

  it('reports the review debt separately, with a way to act on it', () => {
    render(
      <ProfileHeader
        name="Tester"
        scores={scores}
        catalogue={catalogue}
        progress={{ a: { status: 'DECAYING', streak: 3 }, b: { status: 'MASTERED', streak: 3 } }}
      />,
    )
    expect(screen.getByText(/1 due for review/i)).toBeTruthy()
    expect(screen.getByText(/still learned/i)).toBeTruthy()
    expect(screen.getByRole('link', { name: /review/i }).getAttribute('href')).toBe('/review')
  })

  it('renders folded actions when provided', () => {
    render(
      <ProfileHeader
        name="Tester"
        scores={scores}
        catalogue={catalogue}
        progress={progress}
        footer={<a href="/plan">Plan your learning →</a>}
      />,
    )
    expect(screen.getByRole('link', { name: /Plan your learning/ }).getAttribute('href')).toBe('/plan')
  })

  it('shows the profile picture square and boxed on the left', () => {
    render(
      <ProfileHeader
        name="Tester"
        scores={scores}
        catalogue={catalogue}
        progress={progress}
        avatarSeed="s1"
      />,
    )
    const pic = screen.getByLabelText('Tester avatar')
    // Square, not the round default: the override flattens the border radius.
    expect(pic.className).toMatch(/rounded-none/)
    // And it sits in its own bordered box, not floating bare.
    const box = pic.parentElement!
    expect(box.className).toMatch(/border-mathua-border/)
    // Before the text, not after it.
    const region = screen.getByRole('region', { name: 'Where you are' })
    expect(box.compareDocumentPosition(region) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
  })
})
