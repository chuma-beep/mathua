import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import GoalStepper from '../components/GoalStepper'
import { setDailyXPGoal } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, setDailyXPGoal: vi.fn() }
})

const mockSave = vi.mocked(setDailyXPGoal)

describe('GoalStepper', () => {
  beforeEach(() => {
    mockSave.mockReset()
    mockSave.mockResolvedValue(undefined)
  })

  it('renders current goal and presets with active state', () => {
    render(<GoalStepper goal={30} onGoalChange={() => {}} />)
    expect(screen.getByText(/XP\/day/)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Set daily goal to 30 XP' })).toHaveAttribute('aria-pressed', 'true')
    expect(screen.getByRole('button', { name: 'Set daily goal to 15 XP' })).toHaveAttribute('aria-pressed', 'false')
  })

  it('optimistically updates and persists on increment', async () => {
    const onGoalChange = vi.fn()
    render(<GoalStepper goal={30} onGoalChange={onGoalChange} />)
    fireEvent.click(screen.getByRole('button', { name: /increase daily xp goal/i }))
    expect(onGoalChange).toHaveBeenCalledWith(35)
    await waitFor(() => expect(mockSave).toHaveBeenCalledWith(35))
  })

  it('rolls back and shows alert on save failure', async () => {
    mockSave.mockRejectedValueOnce(new Error('offline'))
    const onGoalChange = vi.fn()
    render(<GoalStepper goal={30} onGoalChange={onGoalChange} />)
    fireEvent.click(screen.getByRole('button', { name: 'Set daily goal to 50 XP' }))
    expect(onGoalChange).toHaveBeenCalledWith(50)
    await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument())
    // rollback to previous goal
    expect(onGoalChange).toHaveBeenLastCalledWith(30)
  })

  it('clamps preset below minimum via step buttons', async () => {
    const onGoalChange = vi.fn()
    render(<GoalStepper goal={6} onGoalChange={onGoalChange} />)
    fireEvent.click(screen.getByRole('button', { name: /decrease daily xp goal/i }))
    expect(onGoalChange).toHaveBeenCalledWith(5)
  })
})
