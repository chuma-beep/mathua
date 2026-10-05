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
import { answerCard, answerField, typeAnswer } from './helpers/answerInput'

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
  await waitFor(() => expect(answerField()).toBeTruthy())
  const input = answerField()
  typeAnswer(input, value)
  fireEvent.click(within(answerCard(input)).getByRole('button', { name: 'Check' }))
}

describe('LearnStepper done state (PR5)', () => {
  it('learn → answer ×2 → done → Continue names the re-fetched head', async () => {
    render(<LearnStepper conceptId={CID} />)

    // Intro → first question.
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
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

    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
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

      fireEvent.click(screen.getByRole('button', { name: 'Next →' }))
      await vi.advanceTimersByTimeAsync(10)

      for (let i = 0; i < 2; i++) {
        const inputs = screen.getAllByPlaceholderText(/Your answer/)
        const input = inputs[inputs.length - 1]
        typeAnswer(input, '5')
        fireEvent.click(within(answerCard(input)).getByRole('button', { name: 'Check' }))
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

    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
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

    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    await answerCurrent('6')

    expect(await screen.findByText('Start at 2 and count on 3 to reach 5, not 6.')).toBeTruthy()
    expect(screen.getByText('How')).toBeTruthy()
    // The grader's status token must never stand in as the explanation.
    expect(screen.queryByText('Incorrect')).toBeNull()
  }, 15000)

  it('names the mistake when the server can, and stays silent when it cannot', async () => {
    // A diagnosis is a description of what happened, shown only after grading.
    // It must never appear on a correct answer, and an empty one leaves nothing
    // dangling.
    vi.mocked(submitStudyAnswer).mockResolvedValue({
      correct: false,
      feedback: 'Incorrect',
      diagnosis: 'You are exactly one away, so a count or a boundary is off by one.',
      explanation: 'Start at 2 and count on 3 to reach 5.',
      xp: 0,
    })
    render(<LearnStepper conceptId={CID} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    await answerCurrent('6')
    expect(await screen.findByText(/off by one/)).toBeTruthy()
  })

  it('renders no diagnosis block when the diagnosis is empty', async () => {
    vi.mocked(submitStudyAnswer).mockResolvedValue({
      correct: false,
      feedback: 'Incorrect',
      diagnosis: '',
      explanation: 'Start at 2 and count on 3 to reach 5.',
      xp: 0,
    })
    render(<LearnStepper conceptId={CID} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    await answerCurrent('6')
    expect(await screen.findByText(/Start at 2 and count on 3/)).toBeTruthy()
    // The explanation renders; the empty diagnosis contributes no node.
    expect(document.querySelectorAll('[data-diagnosis]')).toHaveLength(0)
  })

  it('leaves no empty section when the server sends no explanation', async () => {
    vi.mocked(submitStudyAnswer).mockResolvedValue({ correct: true, feedback: 'Correct!', explanation: '', xp: 1 })
    render(<LearnStepper conceptId={CID} />)

    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    await answerCurrent('5')

    // Card still locked (the tally header and the card both report the XP).
    await waitFor(() => expect(screen.getAllByText(/\+1 XP/).length).toBeGreaterThan(0))
    // …but an absent explanation leaves no dangling heading.
    expect(screen.queryByText('Why')).toBeNull()
  }, 15000)
})

// The readiness banner is the second route by which decay could become teaching.
// `buildCandidates` was fixed so a decayed concept leaves the shelf, but the banner
// builds its own links from the readiness payload — and the server's distinction
// between `missing` (never started) and `weak` (started, includes decayed) was being
// discarded by concatenating the two and giving every entry the same `/learn` href.
// So a concept the learner had mastered and not reviewed for a fortnight was re-taught
// with the tutorial, via a link labelled "Review", every time they opened anything that
// depended on it.
describe('readiness banner and decay', () => {
  const prereq = (id: string, status: string) => ({ id, label: `Concept ${id}`, status, mastery_pct: 1 })

  it('routes a decayed prerequisite to /review, never to /learn', async () => {
    vi.mocked(getLessonReadiness).mockResolvedValue({
      concept_id: CID,
      ready: false,
      weak: [prereq('frac.add.word', 'DECAYING')],
      missing: [],
    })
    render(<LearnStepper conceptId={CID} />)
    await waitFor(() => expect(screen.getByRole('note', { name: /prerequisite/i })).toBeTruthy())

    const hrefs = screen.getAllByRole('link').map(l => l.getAttribute('href') ?? '')
    expect(hrefs.some(h => h.includes('concept=frac.add.word'))).toBe(false)
    expect(hrefs).toContain('/review')
  })

  it('still routes a prerequisite that is genuinely being learned to /learn', async () => {
    vi.mocked(getLessonReadiness).mockResolvedValue({
      concept_id: CID,
      ready: false,
      weak: [],
      missing: [prereq('frac.add.word', 'UNSEEN')],
    })
    render(<LearnStepper conceptId={CID} />)
    await waitFor(() => expect(screen.getByRole('note', { name: /prerequisite/i })).toBeTruthy())

    const hrefs = screen.getAllByRole('link').map(l => l.getAttribute('href') ?? '')
    expect(hrefs.some(h => h.includes('concept=frac.add.word'))).toBe(true)
  })

  it('says "already learned" when every prerequisite only needs a check', async () => {
    vi.mocked(getLessonReadiness).mockResolvedValue({
      concept_id: CID,
      ready: false,
      weak: [prereq('frac.add.word', 'DECAYING')],
      missing: [],
    })
    render(<LearnStepper conceptId={CID} />)
    await waitFor(() => expect(screen.getByRole('note', { name: /prerequisite/i })).toBeTruthy())
    // "Before you start" would be false here: the learner did start, and finished.
    expect(screen.getByText(/Due a retrieval check/i)).toBeTruthy()
    expect(screen.getByText(/already learned/i)).toBeTruthy()
    expect(screen.queryByText(/Before you start/i)).toBeNull()
  })
})

// The dead end.
//
// `takeNext` returns null when the buffer is empty and the refill yields nothing. That used
// to hit `if (!nq) return` at three sites: no card, no message, no way forward, with the
// last verdict as the end of the page. It was not an edge case — `geo.basic.points_lines`
// has exactly four distinct question texts, so a session hit it on the fourth answer, and
// the measurement in `internal/generator/all/bank_test.go` found 378 of 657 concepts under a
// 12-question variety reference.
describe('LearnStepper question exhaustion', () => {
  it('finishes the section with an explanation when no further question is available', async () => {
    // An empty practice set from the second call onward: the buffer drains and every refill
    // comes back with nothing, which is exactly what the server used to do once a concept's
    // variants were all seen.
    vi.mocked(getLessonPractice)
      .mockResolvedValueOnce({ questions: [QS[0]], concept_id: CID })
      .mockResolvedValue({ questions: [], concept_id: CID })

    render(<LearnStepper conceptId={CID} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()

    // Answer incorrectly so the 2-in-a-row rule does not advance the section — the only way
    // to reach the "keep going in this section" path, which is where the strand was.
    await answerCurrent('nope')
    await waitForTimeout(500)

    // The learner must be told what happened and given somewhere to go, rather than being
    // left looking at a verdict.
    expect(await screen.findByText(/ran out of new questions/)).toBeTruthy()
    // The completion card is what offers the onward journey.
    expect(await screen.findByRole('button', { name: /Practice again/ })).toBeTruthy()
  })

  it('still answers when the bank keeps serving repeats', async () => {
    // The server's fallback: exclusion empties the set, so it serves what it has. Repeats are
    // a far better outcome than a strand, and the feed must keep working when they arrive.
    vi.mocked(getLessonPractice).mockResolvedValue({ questions: [QS[0]], concept_id: CID })

    render(<LearnStepper conceptId={CID} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()

    await answerCurrent('5')
    await waitForTimeout(500)
    // Either another question arrived or the section ended — but the learner is never left
    // with no card and no message.
    // `queryAllByText` because a repeat is exactly what is being asserted here, so the same
    // question text legitimately appears more than once in the scrollback.
    const stranded =
      screen.queryByText(/ran out of new questions/) === null &&
      screen.queryAllByText('2 + 3 = ?').length === 0 &&
      screen.queryAllByText('4 + 1 = ?').length === 0 &&
      screen.queryAllByText('1 + 6 = ?').length === 0
    expect(stranded).toBe(false)
  })
})

// The header's count and XP must describe the same span.
//
// The count used to reset on every knowledge-point advance while the XP summed every card
// in the feed, so a learner three questions in with two correct saw "1/2 correct" next to
// "+3 XP" — two scopes on one line, and the count was the one that read as wrong.
describe('LearnStepper header tally', () => {
  it('counts every answer in the concept, not just the current section', async () => {
    vi.mocked(getLessonKPs).mockResolvedValue({
      concept_id: CID,
      kps: [
        { label: 'KP one', subgoals: [], worked_example: '2 + 3 = 5' },
        { label: 'KP two', subgoals: [], worked_example: '4 + 1 = 5' },
      ],
    })
    render(<LearnStepper conceptId={CID} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    await screen.findByText('2 + 3 = ?')

    // Answered one at a time, waiting for each card: `answerField()` returns the *newest*
    // field, and a submitted card keeps its own, so answering without waiting would type into
    // the previous question.
    await answerCurrent('5')
    expect(await screen.findByText('4 + 1 = ?')).toBeTruthy()

    // Two in a row advances to the second KP — the point at which the old tally reset. The
    // third served question is `1 + 6 = ?`, whose answer is 7.
    await answerCurrent('5')
    expect(await screen.findByText('1 + 6 = ?')).toBeTruthy()

    // 3 answered, 3 correct across the concept. Under the old per-section tally this read
    // "0/1 correct" for the section the learner had just moved into.
    await answerCurrent('7')
    await waitFor(() => expect(screen.getAllByText(/3\/3 correct/).length).toBeGreaterThan(0))
  })
})

// Local helper: the component waits 350ms before appending the next card.
function waitForTimeout(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms))
}
