import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'

const apiMocks = vi.hoisted(() => ({
  getGoal: vi.fn(),
  getGraph: vi.fn(),
  getProgress: vi.fn(),
}))

const authMocks = vi.hoisted(() => ({
  getUserInfo: vi.fn(),
}))

vi.mock('../lib/api', () => apiMocks)
vi.mock('../lib/auth', () => authMocks)

import GoalContinueCard from '../components/GoalContinueCard'

describe('GoalContinueCard', () => {
  beforeEach(() => {
    apiMocks.getGoal.mockReset()
    apiMocks.getGraph.mockReset()
    apiMocks.getProgress.mockReset()
    authMocks.getUserInfo.mockReset()
  })

  it('shows goal progress and the next step', async () => {
    authMocks.getUserInfo.mockReturnValue({ student_id: 's1' })
    apiMocks.getGoal.mockResolvedValue('b')
    apiMocks.getGraph.mockResolvedValue({
      count: 2,
      nodes: [
        { id: 'a', label: 'Alpha', domain: 'arithmetic', prerequisites: [], grading_type: 'numeric' },
        { id: 'b', label: 'Beta', domain: 'algebra', prerequisites: ['a'], grading_type: 'numeric' },
      ],
    })
    apiMocks.getProgress.mockResolvedValue({
      a: { status: 'MASTERED', streak: 5 },
    })

    render(<GoalContinueCard />)

    expect(await screen.findByRole('heading', { name: 'Beta' })).toBeInTheDocument()
    expect(screen.getByText('1 of 2 concepts learned')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Continue toward Beta/ })).toHaveAttribute(
      'href',
      '/learn?concept=b',
    )
    expect(screen.getByRole('link', { name: 'View path' })).toHaveAttribute(
      'href',
      '/graph?goal=b',
    )
  })

  it('renders nothing for a signed-out learner', () => {
    authMocks.getUserInfo.mockReturnValue(null)
    render(<GoalContinueCard />)
    expect(screen.queryByLabelText('Your goal')).not.toBeInTheDocument()
  })
})
