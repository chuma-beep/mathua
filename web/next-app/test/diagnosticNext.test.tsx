import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { createRef } from 'react'
import { DiagnosticStep } from '../app/onboard/steps'
import DiagnosticHost from '../app/goals/DiagnosticHost'
import QuizHost from '../app/goals/QuizHost'

const PROG0 = {
  answered: 0,
  estimated_total: 25,
  min_total: 25,
  max_total: 45,
  cover_done: 0,
  cover_size: 5,
  done: false,
}

vi.mock('../lib/api', () => ({
  startGoalDiagnostic: vi.fn(async () => ({
    session_id: 's1',
    concept_id: 'c1',
    concept_name: 'Q1 concept',
    question: 'Q1 text',
    done: false,
    progress: PROG0,
  })),
  submitGoalAnswer: vi.fn(async () => ({
    done: false,
    correct: true,
    feedback: 'Nice!',
    concept_id: 'c2',
    concept_name: 'Q2 concept',
    question: 'Q2 text',
    progress: { ...PROG0, answered: 1, cover_done: 1 },
  })),
  resumeGoalDiagnostic: vi.fn(),
  getGoalPlan: vi.fn(),
  startQuizSession: vi.fn(async () => ({
    session_id: 'q1',
    concept_id: 'c1',
    concept_name: 'Q1 quiz',
    question: 'Quiz Q1 text',
    done: false,
    time_limit_seconds: 30,
  })),
  submitQuizAnswer: vi.fn(async () => ({
    done: false,
    correct: true,
    feedback: 'Good!',
    xp: 20,
    concept_id: 'c2',
    concept_name: 'Q2 quiz',
    question: 'Quiz Q2 text',
    time_limit_seconds: 30,
  })),
}))

const stepProps = {
  question: 'Q1 text',
  conceptName: 'Q1 concept',
  progress: null,
  questionCount: 1,
  answerInput: '',
  accuracy: { correct: 0, total: 0 },
  loading: false,
  inputRef: createRef<HTMLInputElement>(),
  conceptId: 'c1',
  sessionId: 's1',
  onInputChange: () => {},
  onSubmit: () => {},
  onNext: () => {},
}

describe('DiagnosticStep Next button', () => {
  it('shows no Next button before answering', () => {
    render(<DiagnosticStep {...stepProps} lastResult={null} />)
    expect(screen.getByRole('button', { name: 'Check Answer' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
  })

  it('shows Next after feedback and calls onNext on click', () => {
    const onNext = vi.fn()
    render(
      <DiagnosticStep
        {...stepProps}
        lastResult={{ correct: true, feedback: 'Nice!' }}
        onNext={onNext}
      />,
    )
    const next = screen.getByRole('button', { name: 'Next →' })
    expect(next).toBeInTheDocument()
    fireEvent.click(next)
    expect(onNext).toHaveBeenCalledTimes(1)
  })
})

describe('DiagnosticHost manual advance', () => {
  beforeEach(() => {
    sessionStorage.clear()
  })

  it('holds feedback until Next reveals the staged question and focuses input', async () => {
    render(
      <DiagnosticHost
        startIds={[]}
        resumeId={null}
        onComplete={() => {}}
        onResumeExpired={() => {}}
      />,
    )
    await screen.findByText('Q1 text')

    fireEvent.change(screen.getByPlaceholderText(/your answer/i), { target: { value: '5' } })
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    // Feedback shows, staged Q2 stays hidden — no auto-advance.
    await screen.findByText('Nice!')
    expect(screen.queryByText('Q2 text')).toBeNull()
    const next = screen.getByRole('button', { name: 'Next →' })
    expect(next).toBeInTheDocument()

    fireEvent.click(next)
    await screen.findByText('Q2 text')
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
    expect(screen.getByPlaceholderText(/your answer/i)).toBe(document.activeElement)
  })
})

describe('QuizHost manual advance', () => {
  it('holds feedback until Next reveals the staged question and focuses input', async () => {
    render(<QuizHost />)
    await screen.findByText('Quiz Q1 text')

    fireEvent.change(screen.getByPlaceholderText(/your answer/i), { target: { value: '8' } })
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByText('Good!')
    expect(screen.queryByText('Quiz Q2 text')).toBeNull()

    fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
    await screen.findByText('Quiz Q2 text')
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
    expect(screen.getByPlaceholderText(/your answer/i)).toBe(document.activeElement)
  })
})
