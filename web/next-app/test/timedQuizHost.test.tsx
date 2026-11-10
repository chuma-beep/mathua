import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import TimedQuizHost from '../app/timed-quiz/TimedQuizHost'
import { getSettings, getTimedQuiz, startTimedQuizSession, submitTimedQuizAnswer } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    getSettings: vi.fn(),
    getTimedQuiz: vi.fn(),
    startTimedQuizSession: vi.fn(),
    submitTimedQuizAnswer: vi.fn(),
    skipTimedQuizItem: vi.fn(),
    finalizeTimedQuiz: vi.fn(),
    getTimedQuizResult: vi.fn(),
  }
})

vi.mock('next/link', () => ({
  default: ({ href, children, ...rest }: { href: string; children: React.ReactNode }) => (
    <a href={href} {...rest}>{children}</a>
  ),
}))

vi.mock('../components/KatexContent', () => ({
  default: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}))

const summary = {
  eligible_xp: 150,
  gate_xp: 150,
  due: true,
  server_now: Date.now(),
} as Awaited<ReturnType<typeof getTimedQuiz>>

const session = {
  run_id: 'run-1',
  status: 'active',
  started_at: Date.now(),
  deadline_at: Date.now() + 10 * 60 * 1000,
  server_now: Date.now(),
  duration_seconds: 600,
  questions_total: 1,
  closed_book: true,
  items: [
    {
      id: 'run-1:0',
      item_id: 'run-1:0',
      concept_id: 'a',
      concept_name: 'Concept A',
      position: 0,
      question: '2 + 2',
      grading_type: 'numeric',
      answer_hint: 'A number',
      answer_input_mode: 'numeric' as const,
      served_at: Date.now(),
      status: 'served',
    },
  ],
} as Awaited<ReturnType<typeof startTimedQuizSession>>

const result = {
  run_id: 'run-1',
  status: 'completed',
  started_at: Date.now() - 60_000,
  deadline_at: Date.now() + 540_000,
  completed_at: Date.now(),
  question_count: 1,
  correct_count: 1,
  unanswered_count: 0,
  xp_earned: 5,
  score_percent: 100,
  items: [
    {
      item_id: 'run-1:0',
      concept_id: 'a',
      concept_name: 'Concept A',
      position: 0,
      question: '2 + 2',
      expected_answer: '4',
      answer: '4',
      correct: true,
      status: 'answered',
      explanation: 'Two plus two is four.',
      xp_earned: 5,
      elapsed_seconds: 4,
    },
  ],
  remediation: [],
}

describe('TimedQuizHost', () => {
  beforeEach(() => {
    vi.mocked(getSettings).mockResolvedValue({})
    vi.mocked(getTimedQuiz).mockResolvedValue(summary)
    vi.mocked(startTimedQuizSession).mockResolvedValue(session)
    vi.mocked(submitTimedQuizAnswer).mockResolvedValue({
      run_id: 'run-1',
      item_id: 'run-1:0',
      correct: true,
      feedback: 'Correct',
      explanation: 'Two plus two is four.',
      xp: 5,
      completed: true,
      result,
    })
  })

  it('starts a due timed quiz and shows the persisted result', async () => {
    render(<TimedQuizHost />)
    const start = await screen.findByRole('button', { name: 'Start timed quiz' })
    fireEvent.click(start)

    const input = await screen.findByLabelText('Your answer')
    fireEvent.change(input, { target: { value: '4' } })
    fireEvent.click(screen.getByRole('button', { name: 'Check answer' }))

    await waitFor(() => expect(screen.getByText('100% correct')).toBeDefined())
    expect(submitTimedQuizAnswer).toHaveBeenCalledWith(
      'run-1',
      'run-1:0',
      '4',
      expect.any(Number),
      expect.any(String),
      false,
    )
    expect(screen.getByText('Two plus two is four.')).toBeDefined()
  })

  it('opens the full device keyboard by default, not the numeric pad', async () => {
    render(<TimedQuizHost />)
    fireEvent.click(await screen.findByRole('button', { name: 'Start timed quiz' }))
    const input = await screen.findByLabelText('Your answer')
    expect(input).toHaveAttribute('inputmode', 'text')
  })

  it('uses the numeric pad only when the learner opts in', async () => {
    vi.mocked(getSettings).mockResolvedValue({ timed_quiz_number_pad: true })
    render(<TimedQuizHost />)
    fireEvent.click(await screen.findByRole('button', { name: 'Start timed quiz' }))
    const input = await screen.findByLabelText('Your answer')
    await waitFor(() => expect(input).toHaveAttribute('inputmode', 'decimal'))
  })

  it('does not offer a start button when the gate is not due', async () => {
    vi.mocked(getTimedQuiz).mockResolvedValue({ ...summary, eligible_xp: 20, due: false })
    render(<TimedQuizHost />)
    await screen.findByText(/Earn 130 more eligible XP/)
    expect(screen.queryByRole('button', { name: 'Start timed quiz' })).toBeNull()
  })
})
