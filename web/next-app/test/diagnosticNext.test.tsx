import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { createRef } from 'react'
import { DiagnosticStep, WelcomeStep } from '../app/onboard/steps'
import OnboardPage from '../app/onboard/page'
import QuizHost from '../app/goals/QuizHost'
import {
  submitGoalAnswer,
  skipGoalQuestion,
  retryGoalQuestion,
  getGoalPlan,
  submitQuizAnswer,
  skipQuizQuestion,
} from '../lib/api'
import { formatForGradingType } from '../lib/answerFormat'
import { answerValue, findAnswerField, findMathField, isMathField, typeAnswer } from './helpers/answerInput'

// The onboard page routes to the results screen and back through the router.
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn(), back: vi.fn(), refresh: vi.fn() }),
  usePathname: () => '/onboard',
  useSearchParams: () => new URLSearchParams(),
}))

describe('formatForGradingType', () => {
  it('maps known types to hint + keyboard, unknowns to generic text', () => {
    expect(formatForGradingType('numeric')).toEqual({ hint: 'Answer with a number, fraction or mixed number (e.g. 4 1/10)', inputMode: 'numeric' })
    expect(formatForGradingType('expression').inputMode).toBe('text')
    expect(formatForGradingType('multiple_choice').hint).toBe('Select an option or type its letter')
    expect(formatForGradingType(undefined)).toEqual({ hint: 'Answer in the form shown', inputMode: 'text' })
    expect(formatForGradingType('bogus')).toEqual({ hint: 'Answer in the form shown', inputMode: 'text' })
  })
})

const PROG0 = {
  answered: 0,
  estimated_total: 25,
  min_total: 25,
  max_total: 45,
  probe_budget: 45,
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
    progress: { ...PROG0, answered: 1 },
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
    progress: { ...PROG0, answered: 1 },
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
    needsName: false,
    demoName: '',
    onDemoName: () => {},
    onStartDemo: () => {},
    onCancelDemo: () => {},
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
        submitError={{ message: 'boom', status: 503 }}
        onSubmit={onRetry}
      />,
    )
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText('(503) boom')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })

  it('retries a failed submission as the kind of answer it was', () => {
    // A failed "I don't know" retried as a graded answer posts an empty answer
    // and records a miss the learner never made, so the page passes a retry that
    // remembers which kind of submission failed.
    const onSubmit = vi.fn()
    const onRetrySubmit = vi.fn()
    render(
      <DiagnosticStep
        {...stepProps}
        lastResult={null}
        submitError={{ message: 'boom', status: 503 }}
        onSubmit={onSubmit}
        onRetrySubmit={onRetrySubmit}
      />,
    )
    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))
    expect(onRetrySubmit).toHaveBeenCalledTimes(1)
    expect(onSubmit).not.toHaveBeenCalled()
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

const PLAN_FIXTURE = {
  readiness: 0.62,
  total_tested: 12,
  correct_count: 9,
  weak_areas: { arithmetic: [] },
  strong_areas: {},
  frontier_label: 'Linear equations',
  frontier_idx: 3,
  frontier_conditional: false,
  conditionally_completed: [],
  placement_course_id: 'a1',
  completion_estimates: {},
}

describe('onboard diagnostic flow', () => {
  beforeEach(() => {
    sessionStorage.clear()
  })

  // The diagnostic runs inside /onboard now; /goals redirects unless it is the
  // quiz. These tests used to render a component nothing mounted, so they proved
  // the flow on a path no learner could reach.
  async function beginDiagnostic() {
    render(<OnboardPage />)
    fireEvent.click(await screen.findByRole('button', { name: /Arithmetic/ }))
    fireEvent.click(screen.getByRole('button', { name: /Start diagnostic test/ }))
    fireEvent.click(await screen.findByRole('button', { name: /Begin diagnostic/ }))
    await screen.findByText('Q1 text')
  }

  it('holds feedback until Next reveals the staged question and focuses input', async () => {
    await beginDiagnostic()

    expect(screen.getByText('Answer with a number, fraction or mixed number (e.g. 4 1/10)')).toBeInTheDocument()
    // /onboard does not pass a grading type to the input, so the answer field is
    // the plain input rather than the math editor; the hint comes from the served type.
    expect(await findAnswerField()).toBeInTheDocument()

    typeAnswer(await findAnswerField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByText('Nice!')
    expect(screen.queryByText('Q2 text')).toBeNull()
    const next = screen.getByRole('button', { name: 'Next →' })
    fireEvent.click(next)
    await screen.findByText('Q2 text')
    expect(screen.queryByRole('button', { name: 'Next →' })).toBeNull()
    expect(await findAnswerField()).toBe(document.activeElement)
  })

  it('numbers the question on screen, and reports how long it took', async () => {
    const now = vi.spyOn(Date, 'now').mockReturnValue(1_000_000)
    await beginDiagnostic()
    expect(screen.getAllByText('Question 1')).toHaveLength(2)

    typeAnswer(await findAnswerField(), '5')
    now.mockReturnValue(1_012_000)
    vi.mocked(submitGoalAnswer).mockClear()
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))
    await screen.findByText('Nice!')
    expect(vi.mocked(submitGoalAnswer).mock.calls[0]?.[3]).toBe(12)
    expect(screen.getAllByText('Question 1')).toHaveLength(2)
    expect(screen.queryAllByText('Question 2')).toHaveLength(0)

    fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
    await screen.findByText('Q2 text')
    expect(screen.getAllByText('Question 2')).toHaveLength(2)
    now.mockRestore()
  })

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
    await beginDiagnostic()
    typeAnswer(await findAnswerField(), '4')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    expect(await screen.findByText(/Not quite/)).toBeInTheDocument()
    expect(screen.queryByText(/count on from 2/)).toBeNull()
    expect(screen.queryByText(/one away/)).toBeNull()
    expect(screen.getByRole('link', { name: /Learn this concept/ })).toHaveAttribute('href', '/learn?concept=c1')
  })

  it('shows the error block with Retry and Skip on submit failure', async () => {
    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 500'), {
        status: 503,
        serverMessage: 'failed to get next question',
      }),
    )
    await beginDiagnostic()

    typeAnswer(await findAnswerField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByRole('alert')
    expect(screen.getByText('(503) failed to get next question')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Retry' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Skip →' })).toBeInTheDocument()
    expect(answerValue(await findAnswerField())).toBe('5')

    fireEvent.click(screen.getByRole('button', { name: 'Retry' }))
    await screen.findByText('Nice!')
    fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
    await screen.findByText('Q2 text')
  })

  it('completes a run whose session ended after the last answer', async () => {
    // The final answer committed and the server deleted the session, so the
    // client's retry finds nothing. A plan exists, which means the run finished.
    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 404'), { status: 404, serverMessage: 'diagnostic session not found' }),
    )
    vi.mocked(getGoalPlan).mockResolvedValueOnce(PLAN_FIXTURE as never)
    await beginDiagnostic()
    typeAnswer(await findAnswerField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    expect(await screen.findByText(/Start here: Linear equations/)).toBeInTheDocument()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('shows only Restart on an expired session', async () => {
    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 404'), {
        status: 404,
        serverMessage: 'diagnostic session not found',
      }),
    )
    await beginDiagnostic()

    typeAnswer(await findAnswerField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByRole('alert')
    expect(screen.getByText(/Session expired/)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Retry' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Skip →' })).toBeNull()
    expect(vi.mocked(skipGoalQuestion)).not.toHaveBeenCalled()
  })

  it('Skip serves the next staged question without recording', async () => {
    vi.mocked(submitGoalAnswer).mockRejectedValueOnce(
      Object.assign(new Error('Goal answer failed: 500'), {
        status: 503,
        serverMessage: 'failed to get next question',
      }),
    )
    await beginDiagnostic()

    typeAnswer(await findAnswerField(), '5')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))
    await screen.findByRole('alert')

    fireEvent.click(screen.getByRole('button', { name: 'Skip →' }))
    await screen.findByText('Q3 text')
    expect(vi.mocked(skipGoalQuestion)).toHaveBeenCalledTimes(1)
    expect(screen.queryByRole('alert')).toBeNull()
    expect(await findAnswerField()).toBe(document.activeElement)
  })

  it("I don't know submits empty with the flag and is labelled No answer, not Not quite", async () => {
    vi.mocked(submitGoalAnswer).mockResolvedValueOnce({
      done: false,
      correct: false,
      concept_id: 'c2',
      concept_name: 'Q2 concept',
      question: 'Q2 text',
      grading_type: 'numeric',
      progress: { ...PROG0, answered: 1 },
    })
    await beginDiagnostic()

    expect(screen.getByRole('button', { name: 'Check Answer' })).toBeDisabled()
    fireEvent.click(screen.getByRole('button', { name: "I don't know" }))

    expect(await screen.findByText('Q2 text')).toBeInTheDocument()
    expect(screen.queryByText(/Not quite/)).toBeNull()
    expect(vi.mocked(submitGoalAnswer)).toHaveBeenCalledWith('s1', 'c1', '', expect.any(Number), true)
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
      progress: { ...PROG0, answered: 1 },
    })
    await beginDiagnostic()

    typeAnswer(await findAnswerField(), 'wrong')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByText('Nope')
    const retry = screen.getByRole('button', { name: 'I made a silly mistake — retry' })
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
        status: 503,
        serverMessage: 'failed to get next quiz question',
      }),
    )
    render(<QuizHost />)
    fireEvent.click(screen.getByRole('button', { name: 'Start quiz →' }))
    await screen.findByText('Quiz Q1 text')

    typeAnswer(await findAnswerField(), '8')
    fireEvent.click(screen.getByRole('button', { name: 'Check Answer' }))

    await screen.findByRole('alert')
    expect(screen.getByText('(503) failed to get next quiz question')).toBeInTheDocument()

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
