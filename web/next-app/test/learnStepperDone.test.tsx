import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, within, waitFor } from '@testing-library/react'
import LearnStepper, { headConceptId } from '../components/LearnStepper'
import {
  getLessonKPs,
  getLessonPractice,
  getLessonReadiness,
  submitStudyAnswer,
  getActivity,
  getProgress,
  getWeaknesses,
  getDueReviews,
  getScores,
} from '../lib/api'
import { getUserInfo } from '../lib/auth'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    getLessonKPs: vi.fn(),
    getLessonPractice: vi.fn(),
    getLessonReadiness: vi.fn(),
    submitStudyAnswer: vi.fn(),
    getActivity: vi.fn(),
    getProgress: vi.fn(),
    getWeaknesses: vi.fn(),
    getDueReviews: vi.fn(),
    getScores: vi.fn(),
  }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn() }
})

// jsdom has no scrollIntoView; the stepper calls it after each append.
Element.prototype.scrollIntoView = vi.fn() as unknown as typeof Element.prototype.scrollIntoView

const CID = 'arith.add.single'

const QS = [
  { question: '2 + 3 = ?', answer: '5', explanation: 'Add them.' },
  { question: '4 + 1 = ?', answer: '5', explanation: 'Add them.' },
  { question: '1 + 6 = ?', answer: '7', explanation: 'Add them.' },
]

beforeEach(() => {
  vi.mocked(getLessonKPs).mockResolvedValue({
    concept_id: CID,
    kps: [{ label: 'Add single digits', subgoals: [], worked_example: '2 + 3 = 5' }],
  })
  vi.mocked(getLessonPractice).mockResolvedValue({ questions: QS, concept_id: CID })
  vi.mocked(getLessonReadiness).mockResolvedValue({ concept_id: CID, ready: true, weak: [], missing: [] })
  // The server grades by looking the question up in its own anchor; the mock
  // mirrors that rather than trusting a caller-supplied expected answer.
  vi.mocked(submitStudyAnswer).mockImplementation(async (_cid, answer, _elapsed, question) => ({
    correct: answer.trim() === (QS.find(q => q.question === question)?.answer ?? '').trim(),
    feedback: 'Correct!',
    explanation: '',
    xp: 1,
  }))
  // Shelf head fetch on done: an experienced learner, nothing due.
  vi.mocked(getUserInfo).mockReturnValue({ student_id: 's1', diagnostic_completed: true } as ReturnType<typeof getUserInfo>)
  vi.mocked(getActivity).mockResolvedValue([])
  vi.mocked(getProgress).mockResolvedValue({})
  vi.mocked(getWeaknesses).mockResolvedValue({ by_domain: {} })
  vi.mocked(getDueReviews).mockResolvedValue({ count: 0 })
  vi.mocked(getScores).mockResolvedValue({ concepts_mastered: 3 } as Awaited<ReturnType<typeof getScores>>)
})

async function answerCurrent(value: string) {
  const inputs = await screen.findAllByPlaceholderText(/Your answer/)
  const input = inputs[inputs.length - 1]
  fireEvent.change(input, { target: { value } })
  const card = input.closest('div[class*="border"]') ?? document.body
  fireEvent.click(within(card as HTMLElement).getByRole('button', { name: 'Check' }))
}

describe('LearnStepper done state (PR5)', () => {
  it('learn → answer ×2 → done → Continue names the re-fetched head', async () => {
    render(<LearnStepper conceptId={CID} />)

    // Intro → first question.
    fireEvent.click(await screen.findByRole('button', { name: 'Skip the example →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()

    // Two in a row advances the single KP → done (REQUIRED_IN_A_ROW = 2).
    await answerCurrent('5')
    expect(await screen.findByText('4 + 1 = ?')).toBeTruthy()
    await answerCurrent('5')

    // Closure card with mastery feedback…
    expect(await screen.findByText(/Complete — 2\/2 correct/)).toBeTruthy()

    // …and a Continue naming the head re-fetched AFTER the last answer,
    // labeled with the reason so the jump reads as intentional.
    const cont = await screen.findByRole('link', { name: /^New: / })
    const href = cont.getAttribute('href')!
    expect(href).toMatch(/^\/learn\?concept=/)
    // The finished concept is excluded from its own shelf: no self-link.
    expect(href).not.toContain(`concept=${encodeURIComponent(CID)}`)
    expect(submitStudyAnswer).toHaveBeenCalledTimes(2)

    // Practice again is always offered as the secondary action.
    expect(screen.getByRole('button', { name: 'Practice again' })).toBeTruthy()

    // Alternatives stay available in done — learner disposes, no bounce.
    expect(screen.getByText(/Or pick something else/)).toBeTruthy()

    // No auto-advance: only explicit navigation, Reference still offered.
    // It carries ?from= so the reference page can offer a way back here.
    expect(screen.getByRole('link', { name: 'Reference' }).getAttribute('href')).toBe(
      `/study?concept=${encodeURIComponent(CID)}&from=${encodeURIComponent(CID)}`,
    )
  }, 15000)

  it('still offers Continue when shelf APIs fail (fail-soft defaults)', async () => {
    vi.mocked(getActivity).mockRejectedValue(new Error('offline'))
    vi.mocked(getWeaknesses).mockRejectedValue(new Error('offline'))
    vi.mocked(getDueReviews).mockRejectedValue(new Error('offline'))
    vi.mocked(getScores).mockRejectedValue(new Error('offline'))
    vi.mocked(getProgress).mockRejectedValue(new Error('offline'))
    render(<LearnStepper conceptId={CID} />)

    fireEvent.click(await screen.findByRole('button', { name: 'Skip the example →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    await answerCurrent('5')
    expect(await screen.findByText('4 + 1 = ?')).toBeTruthy()
    await answerCurrent('5')

    // Per-call catch defaults keep the head satisfiable (browse fallback),
    // so done still offers Continue — never a dead end.
    expect(await screen.findByText(/Complete — 2\/2 correct/)).toBeTruthy()
    const cont = await screen.findByRole('link', { name: /^(New|Due now|Recommended|Continue|Study): / })
    expect(cont.getAttribute('href')!).toMatch(/^\/learn\?concept=|^\/study$|^\/onboard$/)
  }, 15000)

  it('guards same-concept heads by id even when query params differ', () => {
    expect(headConceptId(`/learn?concept=${CID}&seed=42&difficulty=0.7&exclude=a`)).toBe(CID)
    expect(headConceptId('/learn?concept=other.id')).toBe('other.id')
    expect(headConceptId('/review')).toBeNull()
    expect(headConceptId('/study')).toBeNull()
    expect(headConceptId('not a url at all')).toBeNull()
  })

  it('falls back to Practice again plus Back to Profile on shelf timeout', async () => {
    vi.useFakeTimers()
    try {
      const hang = () => new Promise<never>(() => {})
      vi.mocked(getActivity).mockImplementation(hang)
      vi.mocked(getProgress).mockImplementation(hang)
      vi.mocked(getWeaknesses).mockImplementation(hang)
      vi.mocked(getDueReviews).mockImplementation(hang)
      vi.mocked(getScores).mockImplementation(hang)
      render(<LearnStepper conceptId={CID} />)
      await vi.advanceTimersByTimeAsync(10)

      fireEvent.click(screen.getByRole('button', { name: 'Skip the example →' }))
      await vi.advanceTimersByTimeAsync(10)

      for (let i = 0; i < 2; i++) {
        const inputs = screen.getAllByPlaceholderText(/Your answer/)
        const input = inputs[inputs.length - 1]
        fireEvent.change(input, { target: { value: '5' } })
        const card = input.closest('div[class*="border"]') ?? document.body
        fireEvent.click(within(card as HTMLElement).getByRole('button', { name: 'Check' }))
        await vi.advanceTimersByTimeAsync(500)
      }

      expect(screen.getByText(/Complete — 2\/2 correct/)).toBeTruthy()
      // Shelf fetch hangs: past the timeout the failure branch renders.
      await vi.advanceTimersByTimeAsync(9000)
      expect(screen.getByRole('button', { name: 'Practice again' })).toBeTruthy()
      const profile = screen.getByRole('link', { name: 'Back to Profile' })
      expect(profile.getAttribute('href')).toBe('/profile')
      // No self-link anywhere on the done card.
      for (const link of screen.getAllByRole('link')) {
        const href = link.getAttribute('href') ?? ''
        expect(href.startsWith('/learn') && href.includes(`concept=${CID}`)).toBe(false)
      }
    } finally {
      vi.useRealTimers()
    }
  }, 30000)
})

// The loop teaches on both verdicts. A correct answer used to show only the
// verdict and XP, and a miss showed the grader's "Incorrect" in place of the
// worked solution the server had sent.
describe('LearnStepper explanation after submission', () => {
  it('shows the served explanation after a correct answer', async () => {
    vi.mocked(submitStudyAnswer).mockResolvedValue({
      correct: true,
      feedback: 'Correct!',
      explanation: '2 and 3 make 5: start at 2 and count on 3.',
      xp: 2,
    })
    render(<LearnStepper conceptId={CID} />)

    fireEvent.click(await screen.findByRole('button', { name: 'Skip the example →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    await answerCurrent('5')

    // The reasoning for *this* instance, on the correct branch too.
    expect(await screen.findByText('2 and 3 make 5: start at 2 and count on 3.')).toBeTruthy()
    expect(screen.getByText('Why')).toBeTruthy()
    // The card locked: the answer is restated, not just scored.
    expect(screen.getByText('Answer:')).toBeTruthy()
  }, 15000)

  it('shows the served explanation after a wrong answer, and never "Incorrect"', async () => {
    vi.mocked(submitStudyAnswer).mockResolvedValue({
      correct: false,
      feedback: 'Incorrect',
      explanation: 'Start at 2 and count on 3 to reach 5, not 6.',
      xp: 0,
    })
    render(<LearnStepper conceptId={CID} />)

    fireEvent.click(await screen.findByRole('button', { name: 'Skip the example →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    await answerCurrent('6')

    expect(await screen.findByText('Start at 2 and count on 3 to reach 5, not 6.')).toBeTruthy()
    expect(screen.getByText('How')).toBeTruthy()
    // The grader's status token must never stand in as the explanation.
    expect(screen.queryByText('Incorrect')).toBeNull()
  }, 15000)

  it('leaves no empty section when the server sends no explanation', async () => {
    vi.mocked(submitStudyAnswer).mockResolvedValue({ correct: true, feedback: 'Correct!', explanation: '', xp: 1 })
    render(<LearnStepper conceptId={CID} />)

    fireEvent.click(await screen.findByRole('button', { name: 'Skip the example →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    await answerCurrent('5')

    // Card still locked (the tally header and the card both report the XP).
    await waitFor(() => expect(screen.getAllByText(/\+1 XP/).length).toBeGreaterThan(0))
    // …but an absent explanation leaves no dangling heading.
    expect(screen.queryByText('Why')).toBeNull()
  }, 15000)
})
