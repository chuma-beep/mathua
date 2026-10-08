import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { createRef } from 'react'
import { DiagnosticStep, WelcomeStep } from '../app/onboard/steps'
import DiagnosticHost from '../app/goals/DiagnosticHost'
import QuizHost from '../app/goals/QuizHost'
import {
  submitGoalAnswer,
  skipGoalQuestion,
  retryGoalQuestion,
  submitQuizAnswer,
  skipQuizQuestion,
} from '../lib/api'
import { formatForGradingType } from '../lib/answerFormat'
import { answerValue, findAnswerField, findMathField, isMathField, typeAnswer } from './helpers/answerInput'

describe('formatForGradingType', () => {
  it('maps known types to hint + keyboard, unknowns to generic text', () => {
    expect(formatForGradingType('numeric')).toEqual({ hint: 'Answer with a number, fraction, or mixed number (e.g. 4 1/10)', inputMode: 'numeric' })
    expect(formatForGradingType('expression').inputMode).toBe('text')
    expect(formatForGradingType('multiple_choice').hint).toBe('Type the exact answer')
    expect(formatForGradingType(undefined)).toEqual({ hint: 'Answer in the form shown', inputMode: 'text' })
    expect(formatForGradingType('bogus')).toEqual({ hint: 'Answer in the form shown', inputMode: 'text' })
  })
})

const PROG0 = {
  answered: 0,
  estimated_total: 25,
  min_total: 25,
  max_total: 45,
  cover_done: 0,
  cover_size: 5,
  done: false,
}

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    startGoalDiagnostic: vi.fn(async () => ({
    session_id: 's1',
    concept_id: 'c1',
    concept_name: 'Q1 concept',
    question: 'Q1 text',
    done: false,
    grading_type: 'numeric',
    progress: PROG0,
  })),
  submitGoalAnswer: vi.fn(async () => ({
    done: false,
    correct: true,
    feedback: 'Nice!',
    concept_id: 'c2',
    concept_name: 'Q2 concept',
    question: 'Q2 text',
    grading_type: 'numeric',
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
    grading_type: 'numeric',
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
    grading_type: 'numeric',
    time_limit_seconds: 30,
  })),
  skipGoalQuestion: vi.fn(async () => ({
    done: false,
    correct: false,
    feedback: 'Skipped — no evidence recorded.',
    concept_id: 'c3',
    concept_name: 'Q3 concept',
    question: 'Q3 text',
    grading_type: 'numeric',
    progress: { ...PROG0, answered: 1, cover_done: 1 },
  })),
  skipQuizQuestion: vi.fn(async () => ({
    done: false,
    correct: false,
    feedback: 'Skipped — no XP awarded.',
    xp: 0,
    concept_id: 'c3',
    concept_name: 'Q3 quiz',
    question: 'Quiz Q3 text',
    time_limit_seconds: 30,
  })),
  retryGoalQuestion: vi.fn(async () => ({
    done: false,
    concept_id: 'c1',
    concept_name: 'Q1 concept',
    question: 'Q1 text',
    grading_type: 'numeric',
  })),
  }
})

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
  onDontKnow: () => {},
  onNext: () => {},
  done: false,
  answerFormat: formatForGradingType('numeric'),
  submitError: null,
  planError: '',
  onSkip: () => {},
  skipsLeft: 3,
  onRestart: () => {},
  onRetryPlan: () => {},
  retryAvailable: false,
  onRetryQuestion: () => {},
}

describe('WelcomeStep briefing', () => {
  const welcomeProps = {
    domains: [],
    loading: false,
    hasPaused: false,
    selectedCount: 0,
    confirming: false,
    onToggle: () => {},
    onSelectAll: () => {},
    onStart: () => {},
    onResume: () => {},
    onBegin: () => {},
    onCancel: () => {},
  }

  it('states duration and expectations in the review step', () => {
    render(<WelcomeStep {...welcomeProps} confirming />)
    expect(screen.getByText('Before you begin')).toBeInTheDocument()
    expect(screen.getByText('What to expect')).toBeInTheDocument()
    expect(screen.getByText('25–45 questions · about 30–45 minutes')).toBeInTheDocument()
    expect(screen.getByText('Answer from what you know — no searching')).toBeInTheDocument()
    expect(screen.getByText('Speed counts')).toBeInTheDocument()
    expect(screen.getByText('Pause anytime')).toBeInTheDocument()
    expect(screen.queryByText(/frontier/i)).toBeNull()
  })

  it('Start opens review, Begin starts, Back cancels', () => {
    const onStart = vi.fn()
    const onBegin = vi.fn()
    const onCancel = vi.fn()
    const { rerender } = render(
      <WelcomeStep {...welcomeProps} selectedCount={3} onStart={onStart} onBegin={onBegin} onCancel={onCancel} />,
    )
    expect(screen.queryByText('Before you begin')).toBeNull()
    fireEvent.click(screen.getByRole('button', { name: /Start diagnostic test/ }))
    expect(onStart).toHaveBeenCalledTimes(1)

    rerender(
      <WelcomeStep
        {...welcomeProps}
        confirming
        onStart={onStart}
        onBegin={onBegin}
        onCancel={onCancel}
      />,
    )
    fireEvent.click(screen.getByRole('button', { name: 'Begin diagnostic →' }))
    expect(onBegin).toHaveBeenCalledTimes(1)
    fireEvent.click(screen.getByRole('button', { name: '← Back' }))
    expect(onCancel).toHaveBeenCalledTimes(1)
  })
})

describe('DiagnosticStep Next button', () => {  it('shows no Next button before answering', () => {
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

  it('hides Next once finished and shows the error block on submit failure', () => {
    const onRetry = vi.fn()
    const { rerender } = render(
      <DiagnosticStep {...stepProps} lastResult={{ correct: true, feedback: 'Nice!' }} done />,
    )
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
    expect(screen.getByText('Preparing results…')).toBeInTheDocument()

    rerender(
      <DiagnosticStep
        {...stepProps}
        lastResult={null}
        submitError={{ message: 'boom', status: 500 }}
        onSubmit={onRetry}
      />,
    )
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText('(500) boom')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it('shows only Restart on a 404 expired session', () => {    render(
      <DiagnosticStep
        {...stepProps}
        lastResult={null}
        submitError={{ message: 'diagnostic session not found', status: 404 }}
      />,
    )
    expect(screen.getByText(/Session expired/)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Retry' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Skip →' })).toBeNull()
    expect(screen.getByRole('button', { name: 'Restart diagnostic' })).toBeInTheDocument()
  })

  it("shows I don't know before answering, hides it after, and sends the flag", () => {
    const onDontKnow = vi.fn()
    const { rerender } = render(
      <DiagnosticStep {...stepProps} lastResult={null} onDontKnow={onDontKnow} />,
    )
    const btn = screen.getByRole('button', { name: "I don't know" })
    expect(btn).toBeInTheDocument()
    fireEvent.click(btn)
    expect(onDontKnow).toHaveBeenCalledTimes(1)

    rerender(
      <DiagnosticStep
        {...stepProps}
        lastResult={{ correct: false, feedback: 'Nope' }}
        onDontKnow={onDontKnow}
      />,
    )
    expect(screen.queryByRole('button', { name: "I don't know" })).toBeNull()
  })

  it('offers retry only on eligible incorrect misses', () => {
    const onRetryQuestion = vi.fn()
    const { rerender } = render(
      <DiagnosticStep
        {...stepProps}
        lastResult={{ correct: false, feedback: 'Nope' }}
        retryAvailable
        onRetryQuestion={onRetryQuestion}
      />,
    )
    const retry = screen.getByRole('button', { name: 'I made a silly mistake — retry' })
    expect(retry).toBeInTheDocument()
    fireEvent.click(retry)
    expect(onRetryQuestion).toHaveBeenCalledTimes(1)

    // Correct answers never offer retry.
    rerender(
      <DiagnosticStep
        {...stepProps}
        lastResult={{ correct: true, feedback: 'Nice!' }}
        retryAvailable
        onRetryQuestion={onRetryQuestion}
      />,
    )
    expect(screen.queryByRole('button', { name: 'I made a silly mistake — retry' })).toBeNull()

    // Ineligible misses never offer retry.
    rerender(
      <DiagnosticStep
        {...stepProps}
        lastResult={{ correct: false, feedback: 'Nope' }}
        retryAvailable={false}
        onRetryQuestion={onRetryQuestion}
      />,
    )
    expect(screen.queryByRole('button', { name: 'I made a silly mistake — retry' })).toBeNull()
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

    // Format hint from the served grading_type, with numeric keyboard.
    expect(screen.getByText('Answer with a number, fraction, or mixed number (e.g. 4 1/10)')).toBeInTheDocument()
    // `numeric` selects the math editor, which brings its own virtual keyboard
    // rather than the OS numeric pad — so the control itself is the assertion.
    // `findMathField`, not `findAnswerField`: the loading state is a working plain
    // input, so the generic helper would hand back that and pass a test meant for MathLive.
    expect(isMathField(await findMathField())).toBe(true)

    typeAnswer(await findMathField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    // Feedback shows, staged Q2 stays hidden — no auto-advance.
    await screen.findByText('Nice!')
    expect(screen.queryByText('Q2 text')).toBeNull()
    const next = screen.getByRole('button', { name: 'Next →' })
    expect(next).toBeInTheDocument()

    fireEvent.click(next)
    await screen.findByText('Q2 text')
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
    expect(await findAnswerField()).toBe(document.activeElement)
  })

  // A diagnostic measures prior knowledge. A wrong answer must not hand back the solution, and a
  // regressed server that sends one anyway must not get it onto the screen — the client holds its
  // own copy of the contract. The cast is past the response type on purpose: the type no longer
  // has these fields, and the guard must not depend on the type staying honest.
  it('shows no solution on a wrong diagnostic answer, and routes the miss to Learn', async () => {
    vi.mocked(submitGoalAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      feedback: 'Not quite',
      concept_id: 'c2',
      concept_name: 'Q2 concept',
      question: 'Q2 text',
      grading_type: 'numeric',
      explanation: '2 + 3 = 5: count on from 2.',
      diagnosis: 'You are exactly one away, so a count or a boundary is off by one.',
    } as never)
    render(
      <DiagnosticHost startIds={[]} resumeId={null} onComplete={() => {}} onResumeExpired={() => {}} />,
    )
    await screen.findByText('Q1 text')
    typeAnswer(await findMathField(), '4')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    expect(await screen.findByText('✗ Not quite')).toBeInTheDocument()
    expect(screen.queryByText(/count on from 2/)).toBeNull()
    expect(screen.queryByText(/one away/)).toBeNull()
    // The concept just answered is the one the hand-off names, not the staged next one.
    expect(screen.getByRole('link', { name: /Learn this concept/ })).toHaveAttribute('href', '/learn?concept=c1')
  })

  it('shows the error block with Retry and Skip on submit failure', async () => {
    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 500'), {
        status: 500,
        serverMessage: 'failed to get next question',
      }),
    )
    render(
      <DiagnosticHost
        startIds={[]}
        resumeId={null}
        onComplete={() => {}}
        onResumeExpired={() => {}}
      />,
    )
    await screen.findByText('Q1 text')

    typeAnswer(await findMathField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    // No alert: persistent inline error with the server's own words.
    await screen.findByRole('alert')
    expect(screen.getByText('(500) failed to get next question')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Retry' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Skip →' })).toBeInTheDocument()
    // Answer preserved for one-tap retry.
    expect(answerValue(await findAnswerField())).toBe('5')

    // Retry succeeds (mock default) and advances on Next.
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))
    await screen.findByText('Nice!')
    fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
    await screen.findByText('Q2 text')
  })

    it('shows only Restart on an expired session', async () => {    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 404'), {
        status: 404,
        serverMessage: 'diagnostic session not found',
      }),
    )
    render(
      <DiagnosticHost
        startIds={[]}
        resumeId={null}
        onComplete={() => {}}
        onResumeExpired={() => {}}
      />,
    )
    await screen.findByText('Q1 text')

    typeAnswer(await findMathField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    // 404: session gone — Retry/Skip pointless, only Restart.
    await screen.findByRole('alert')
    expect(screen.getByText(/Session expired/)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Retry' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Skip →' })).toBeNull()
    expect(vi.mocked(skipGoalQuestion)).not.toHaveBeenCalled()
  })

  it('Skip serves the next staged question without recording', async () => {
    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 500'), {
        status: 500,
        serverMessage: 'failed to get next question',
      }),
    )
    render(
      <DiagnosticHost
        startIds={[]}
        resumeId={null}
        onComplete={() => {}}
        onResumeExpired={() => {}}
      />,
    )
    await screen.findByText('Q1 text')

    typeAnswer(await findMathField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))
    await screen.findByRole('alert')

    fireEvent.click(screen.getByRole('button', { name: 'Skip →' }))
    await screen.findByText('Q3 text')
    expect(vi.mocked(skipGoalQuestion)).toHaveBeenCalledTimes(1)
    expect(screen.queryByRole('alert')).toBeNull()
    expect(await findAnswerField()).toBe(document.activeElement)
  })

  it("I don't know submits empty with the flag and is labelled No answer, not Not quite", async () => {
    // The server sends no feedback for an admitted unknown — the boundary strips it — so the
    // label is the verdict, and it must not borrow the wrong answer's words.
    vi.mocked(submitGoalAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      concept_id: 'c2',
      concept_name: 'Q2 concept',
      question: 'Q2 text',
      grading_type: 'numeric',
      progress: { ...PROG0, answered: 1, cover_done: 1 },
    })
    render(
      <DiagnosticHost
        startIds={[]}
        resumeId={null}
        onComplete={() => {}}
        onResumeExpired={() => {}}
      />,
    )
    await screen.findByText('Q1 text')

    // Empty input: Check Answer stays disabled, I don't know works.
    expect(screen.getByRole('button', { name: 'Check Answer' })).toBeDisabled()
    fireEvent.click(screen.getByRole('button', { name: "I don't know" }))

    expect(await screen.findByText('No answer')).toBeInTheDocument()
    expect(screen.queryByText(/Not quite/)).toBeNull()
    expect(vi.mocked(submitGoalAnswer)).toHaveBeenCalledWith(
      's1',
      'c1',
      '',
      expect.any(Number),
      true,
    )
  })

  it('retry voids the slip and re-serves the same question', async () => {
    vi.mocked(submitGoalAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      feedback: 'Nope',
      concept_id: 'c2',
      concept_name: 'Q2 concept',
      question: 'Q2 text',
      grading_type: 'numeric',
      retry_available: true,
      progress: { ...PROG0, answered: 1, cover_done: 1 },
    })
    render(
      <DiagnosticHost
        startIds={[]}
        resumeId={null}
        onComplete={() => {}}
        onResumeExpired={() => {}}
      />,
    )
    await screen.findByText('Q1 text')

    typeAnswer(await findAnswerField(), 'wrong')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByText('Nope')
    const retry = screen.getByRole('button', { name: 'I made a silly mistake — retry' })
    expect(retry).toBeInTheDocument()

    fireEvent.click(retry)
    await screen.findByText('Q1 text')
    expect(vi.mocked(retryGoalQuestion)).toHaveBeenCalledWith('s1', 'c1')
    expect(screen.queryByRole('button', { name: 'I made a silly mistake — retry' })).toBeNull()
    expect(await findAnswerField()).toBe(document.activeElement)
  })
})

describe('QuizHost manual advance', () => {
  it('holds feedback until Next reveals the staged question and focuses input', async () => {
    render(<QuizHost />)
    // Intro briefing gates the start — timing begins at Start, not mount.
    expect(screen.getByText('Before you start')).toBeInTheDocument()
    expect(screen.getByText('Up to 5 questions · timed · closed book')).toBeInTheDocument()
    expect(screen.queryByText('Quiz Q1 text')).toBeNull()
    fireEvent.click(screen.getByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')

    typeAnswer(await findAnswerField(), '8')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByText('Good!')
    expect(screen.queryByText('Quiz Q2 text')).toBeNull()

    fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
    await screen.findByText('Quiz Q2 text')
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
    expect(await findAnswerField()).toBe(document.activeElement)
  })

  // A test verifies knowledge and does not teach. These two used to assert the opposite — that
  // the served solution was shown on both verdicts — which meant a learner could read the answer
  // to a question they had just been measured on, and the questions after it measured whether
  // they had read rather than what they knew.
  it('shows no solution on a wrong quiz answer, and routes the miss to Learn', async () => {
    vi.mocked(submitQuizAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      // A grader token is the only feedback the server sends now. Modelled here so the assertion
      // covers what the client does with it.
      feedback: 'Incorrect',
      xp: 0,
      concept_id: 'c2',
      concept_name: 'Q2 quiz',
      question: 'Quiz Q2 text',
      grading_type: 'numeric',
      time_limit_seconds: 30,
    })
    render(<QuizHost />)
    fireEvent.click(await screen.findByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')
    typeAnswer(await findAnswerField(), '9')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    expect(await screen.findByText('✗ Not quite')).toBeInTheDocument()
    // The verdict, and no worked solution anywhere on the page.
    expect(screen.queryByText(/start at 3 and count on 5/)).toBeNull()
    // The one route to instruction leaves the assessment for the concept just answered.
    const learn = screen.getByRole('link', { name: /Learn this concept/ })
    expect(learn).toHaveAttribute('href', '/learn?concept=c1')
  })

  it('shows no solution on a correct quiz answer either', async () => {
    vi.mocked(submitQuizAnswer).mockResolvedValueOnce({
      done: false,
      correct: true,
      feedback: '',
      xp: 2,
      concept_id: 'c2',
      concept_name: 'Q2 quiz',
      question: 'Quiz Q2 text',
      grading_type: 'numeric',
      time_limit_seconds: 30,
    })
    render(<QuizHost />)
    fireEvent.click(await screen.findByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')
    typeAnswer(await findAnswerField(), '8')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    expect(await screen.findByText('✓ Correct!')).toBeInTheDocument()
    // Correct offers no "Learn this concept" hand-off: there is nothing to remediate, and the
    // hand-off's absence on a hit is what keeps it a signal of a miss rather than decoration.
    expect(screen.queryByText(/checked by counting on from 3/)).toBeNull()
    expect(screen.queryByRole('link', { name: /Learn this concept/ })).toBeNull()
  })

  // Belt as well as braces. The server is the boundary — it is what strips the solution, and
  // the Go tests measure that — but the client holds its own copy of the contract, so a server
  // that regressed to sending `explanation` and `diagnosis` would still not put the answer on
  // screen. The response is cast past the type on purpose: the type no longer has these fields,
  // and the point is that the rendering does not depend on the type being honest.
  it('does not render a solution even if one arrives, and routes the miss to Learn', async () => {
    vi.mocked(submitQuizAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      feedback: 'Incorrect',
      xp: 0,
      concept_id: 'c2',
      concept_name: 'Q2 quiz',
      question: 'Quiz Q2 text',
      grading_type: 'numeric',
      time_limit_seconds: 30,
      // A regressed server. Both fields are solution-bearing: the worked answer outright, and
      // "you are one away" which reconstructs it from the learner's own wrong answer.
      explanation: '3 and 5 make 8: start at 3 and count on 5.',
      diagnosis: 'You are exactly one away, so a count or a boundary is off by one.',
    } as never)
    render(<QuizHost />)
    fireEvent.click(await screen.findByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')
    typeAnswer(await findAnswerField(), '9')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    expect(await screen.findByText('✗ Not quite')).toBeInTheDocument()
    expect(screen.queryByText(/start at 3 and count on 5/)).toBeNull()
    expect(screen.queryByText(/one away/)).toBeNull()
    expect(screen.getByRole('link', { name: /Learn this concept/ })).toHaveAttribute('href', '/learn?concept=c1')
  })

  it('QuizHost shows the error block and Skip advances without XP', async () => {
    vi.mocked(submitQuizAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Quiz answer failed: 500'), {
        status: 500,
        serverMessage: 'failed to get next quiz question',
      }),
    )
    render(<QuizHost />)
    fireEvent.click(screen.getByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')

    typeAnswer(await findAnswerField(), '8')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByRole('alert')
    expect(screen.getByText('(500) failed to get next quiz question')).toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: 'Skip →' }))
    await screen.findByText('Quiz Q3 text')
    expect(vi.mocked(skipQuizQuestion)).toHaveBeenCalledTimes(1)
  })

  it("QuizHost I don't know submits empty with the flag and is labelled No answer", async () => {
    vi.mocked(submitQuizAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      xp: 0,
      concept_id: 'c2',
      concept_name: 'Q2 quiz',
      question: 'Quiz Q2 text',
      grading_type: 'numeric',
      time_limit_seconds: 30,
    })
    render(<QuizHost />)
    fireEvent.click(screen.getByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')

    expect(screen.getByRole('button', { name: 'Check Answer' })).toBeDisabled()
    fireEvent.click(screen.getByRole('button', { name: "I don't know" }))

    expect(await screen.findByText('No answer')).toBeInTheDocument()
    expect(screen.queryByText(/Not quite/)).toBeNull()
    expect(vi.mocked(submitQuizAnswer)).toHaveBeenCalledWith(
      'q1',
      'c1',
      '',
      expect.any(Number),
      true,
    )
  })
})
