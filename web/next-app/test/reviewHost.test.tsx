import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ReviewHost from '../app/review/ReviewHost'
import { startReviewSession, submitReviewAnswer } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, startReviewSession: vi.fn(), submitReviewAnswer: vi.fn() }
})

const Q1 = {
  concept_id: 'arith.add.single',
  concept_name: 'Single-digit addition',
  question: '2 + 3 = ?',
  is_review: true,
  attempt_id: 'a1',
}

const Q2 = {
  concept_id: 'arith.sub.single',
  concept_name: 'Single-digit subtraction',
  question: '8 - 3 = ?',
  is_review: true,
  attempt_id: 'a2',
}

function result(correct: boolean, xp: number) {
  return {
    correct,
    feedback: correct ? 'Correct!' : 'Not quite.',
    new_status: 'PRACTICING',
    explanation: '',
    streak: 1,
    required_streak: 2,
    xp,
    expected_answer: '5',
  }
}

beforeEach(() => {
  vi.mocked(startReviewSession).mockResolvedValue({ student_id: 's1', session_id: 'r1', question: Q1 })
  vi.mocked(submitReviewAnswer).mockReset()
})

describe('ReviewHost', () => {
  it('starts a review run and shows the first question', async () => {
    render(<ReviewHost />)
    fireEvent.click(screen.getByRole('button', { name: /Start review/ }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    expect(screen.getByText('Review question 1')).toBeTruthy()
  })

  it('answers, advances manually, then completes with a score', async () => {
    vi.mocked(submitReviewAnswer)
      .mockResolvedValueOnce({ result: result(true, 5), next_question: Q2, done: false })
      .mockResolvedValueOnce({ result: result(true, 5), next_question: null, done: true })
    render(<ReviewHost />)
    fireEvent.click(screen.getByRole('button', { name: /Start review/ }))
    await screen.findByText('2 + 3 = ?')

    fireEvent.change(screen.getByPlaceholderText(/Your answer/), { target: { value: '5' } })
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))
    expect(await screen.findByText('✓ Correct!')).toBeTruthy()

    // Manual advance: staged next question reveals on Next.
    fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
    expect(await screen.findByText('8 - 3 = ?')).toBeTruthy()
    expect(screen.getByText('Review question 2')).toBeTruthy()

    fireEvent.change(screen.getByPlaceholderText(/Your answer/), { target: { value: '5' } })
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))
    expect(await screen.findByText('Review complete', {}, { timeout: 3000 })).toBeTruthy()
    expect(screen.getByText('2/2 correct')).toBeTruthy()
    expect(submitReviewAnswer).toHaveBeenCalledTimes(2)
    // attempt_id rides along for grading.
    expect(vi.mocked(submitReviewAnswer).mock.calls[0][3]).toBe('a1')
    expect(vi.mocked(submitReviewAnswer).mock.calls[1][3]).toBe('a2')
  })

  it('shows the caught-up empty state when nothing is due', async () => {
    vi.mocked(startReviewSession).mockResolvedValue({ student_id: 's1', session_id: 'r1', question: null })
    render(<ReviewHost />)
    fireEvent.click(screen.getByRole('button', { name: /Start review/ }))
    expect(await screen.findByText('All caught up')).toBeTruthy()
  })

  it('surfaces submit failures with retry', async () => {
    vi.mocked(submitReviewAnswer).mockRejectedValueOnce(Object.assign(new Error('boom'), { status: 500 }))
    render(<ReviewHost />)
    fireEvent.click(screen.getByRole('button', { name: /Start review/ }))
    await screen.findByText('2 + 3 = ?')
    fireEvent.change(screen.getByPlaceholderText(/Your answer/), { target: { value: '5' } })
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))
    expect(await screen.findByText(/Couldn't submit your answer/)).toBeTruthy()
  })
})
