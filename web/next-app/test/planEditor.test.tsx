import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import PlanEditor from '../components/PlanEditor'
import { getDestinations, getEstimate, getScores, savePlan, getCurrentPlan } from '../lib/api'
import { getUserInfo, getGuestId } from '../lib/auth'
import { DEFAULT_DAILY_GOAL } from '../lib/plan'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    getDestinations: vi.fn(),
    getEstimate: vi.fn(),
    getScores: vi.fn(),
    savePlan: vi.fn(),
    getCurrentPlan: vi.fn(),
  }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn(), getGuestId: vi.fn() }
})

const dests = [{ id: 'alg', name: 'Algebra', description: '', total: 10, mastered: 2, pct: 0.2 }]

const estimate = {
  destination: 'alg',
  estimate: {
    total: 5,
    mastered: 1,
    remaining: ['a', 'b'],
    xp_remaining: 100,
    time_min_remaining: 125,
    reviews_due: 1,
    quizzes_ahead: 1,
    assessment_min: 10,
    diagnostic_min: 5,
    days: 10,
    finish_date: '2026-06-01',
    required_per_day: 10,
    feasible: true,
    lines: { learning: 80, assessment: 15, diagnostic: 5 },
  },
  pace: { rate: 10, source: 'measured', adherence: 1, active_days: 5, trailing_days: 7 },
  probes: ['a'],
  plan_delta_days: 2,
}

beforeEach(() => {
  vi.mocked(getUserInfo).mockReturnValue({ student_id: 's1' } as ReturnType<typeof getUserInfo>)
  vi.mocked(getGuestId).mockReturnValue('')
  vi.mocked(getDestinations).mockResolvedValue(dests)
  vi.mocked(getScores).mockResolvedValue({ daily_xp_goal: 10 } as Awaited<ReturnType<typeof getScores>>)
  vi.mocked(getEstimate).mockResolvedValue(estimate as Awaited<ReturnType<typeof getEstimate>>)
  vi.mocked(getCurrentPlan).mockResolvedValue({ plan: null })
  vi.mocked(savePlan).mockImplementation(async (body) => ({
    destination: body.destination,
    daily_goal: body.daily_goal,
    deadline_days: body.deadline_days,
    rest_days: body.rest_days,
    xp_remaining: 100,
    created: '2026-01-01',
  }))
})

describe('PlanEditor', () => {
  it('loads destinations, renders the estimate, and saves the plan', async () => {
    render(<PlanEditor />)
    expect(await screen.findByText(/Remaining workload/)).toBeTruthy()
    expect(screen.getByText(/Topics remaining/)).toBeTruthy()
    expect(screen.getByText(/What if/)).toBeTruthy()
    fireEvent.click(screen.getByRole('button', { name: /Plan my learning/ }))
    expect(await screen.findByRole('button', { name: /Update my plan/ })).toBeTruthy()
    expect(vi.mocked(savePlan)).toHaveBeenCalledWith({
      destination: 'alg',
      daily_goal: 10,
      deadline_days: 0,
      rest_days: 1,
    })
  })

  it('deadline presets set the deadline, not the daily goal', async () => {
    render(<PlanEditor />)
    await screen.findByText(/Remaining workload/)
    fireEvent.click(screen.getByRole('button', { name: 'Plan by deadline' }))
    fireEvent.click(screen.getByRole('button', { name: '3 mo' }))
    const last = vi.mocked(getEstimate).mock.calls.at(-1)![1]
    expect(last).toMatchObject({ deadline_days: 90 })
    expect(last.daily_goal).toBe(10)
  })

  it('compact hides detail, what-if, and delta but keeps save', async () => {
    render(<PlanEditor compact />)
    expect(await screen.findByText(/Remaining workload/)).toBeTruthy()
    expect(screen.queryByText(/Topics remaining/)).toBeNull()
    expect(screen.queryByLabelText(/What if/)).toBeNull()
    fireEvent.click(screen.getByRole('button', { name: /Plan my learning/ }))
    expect(await screen.findByRole('button', { name: /Update my plan/ })).toBeTruthy()
  })

  it('surfaces save failures', async () => {
    vi.mocked(savePlan).mockRejectedValue(new Error('offline'))
    render(<PlanEditor />)
    await screen.findByText(/Remaining workload/)
    fireEvent.click(screen.getByRole('button', { name: /Plan my learning/ }))
    expect(await screen.findByText('offline')).toBeTruthy()
  })

  // The regression that made the planner unusable. The save button lived inside the
  // `{e && ...}` estimate block, so any estimate failure removed the only way to save a
  // plan from the page — and the estimate fails whenever destinations fail, which used to
  // be swallowed into an empty array.
  it('can still save when the estimate fails', async () => {
    vi.mocked(getEstimate).mockRejectedValue(new Error('planner unavailable'))
    render(<PlanEditor />)
    const save = await screen.findByRole('button', { name: /Plan my learning/ })
    expect(save).toBeTruthy()
    expect(screen.queryByText(/Remaining workload/)).toBeNull()
    expect(await screen.findByText('planner unavailable')).toBeTruthy()
    fireEvent.click(save)
    await waitFor(() => expect(vi.mocked(savePlan)).toHaveBeenCalled())
  })

  it('says so when there are no destinations, instead of showing an empty dropdown', async () => {
    vi.mocked(getDestinations).mockRejectedValue(new Error('destinations unavailable'))
    render(<PlanEditor />)
    expect(await screen.findByText(/No learning paths are available/)).toBeTruthy()
    expect(screen.queryByRole('combobox')).toBeNull()
    expect(screen.getByText('destinations unavailable')).toBeTruthy()
  })

  it('falls back to the shared default goal, not a literal 30', async () => {
    // getScores failing used to be indistinguishable from "you have no stored goal", and the
    // fallback was written as a bare 30 in three places while the shipped default was 10.
    vi.mocked(getScores).mockRejectedValue(new Error('scores unavailable'))
    render(<PlanEditor />)
    await screen.findByText(/Remaining workload/)
    const last = vi.mocked(getEstimate).mock.calls.at(-1)![1]
    expect(last.daily_goal).toBe(DEFAULT_DAILY_GOAL)
  })
})
