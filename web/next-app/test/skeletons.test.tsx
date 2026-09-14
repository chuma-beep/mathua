import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProfileSkeleton from '../components/skeletons/ProfileSkeleton'
import StudySkeleton from '../components/skeletons/StudySkeleton'
import { LeaderboardTableSkeleton, LeaguesSkeleton } from '../components/skeletons/LeaderboardSkeleton'

describe('skeletons', () => {
  it('ProfileSkeleton announces loading state', () => {
    render(<ProfileSkeleton />)
    expect(screen.getByRole('status', { name: 'Loading profile' })).toBeInTheDocument()
  })

  it('StudySkeleton announces loading state', () => {
    render(<StudySkeleton />)
    expect(screen.getByRole('status', { name: 'Loading lessons' })).toBeInTheDocument()
  })

  it('LeaderboardTableSkeleton renders placeholder rows', () => {
    const { container } = render(
      <table>
        <tbody>
          <LeaderboardTableSkeleton rows={8} />
        </tbody>
      </table>,
    )
    expect(container.querySelectorAll('tbody > tr')).toHaveLength(8)
  })

  it('LeaguesSkeleton announces loading state', () => {
    render(<LeaguesSkeleton />)
    expect(screen.getByRole('status', { name: 'Loading leagues' })).toBeInTheDocument()
  })
})
