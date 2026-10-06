import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, within, waitFor } from '@testing-library/react'
import LearnStepper from '../components/LearnStepper'
import {
  getLessonKPs,
  getLessonPractice,
  getLessonReadiness,
  submitStudyAnswer,
  getProgress,
  getNext,
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
    getProgress: vi.fn(),
    getNext: vi.fn(),
  }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn() }
})

Element.prototype.scrollIntoView = vi.fn() as unknown as typeof Element.prototype.scrollIntoView

const CID = 'arith.add.single'

const STEPS = ['Name the operation', 'Rewrite it as an equation']

const rec = (id: string) => ({
  id,
  conceptId: id,
  conceptTitle: id.toUpperCase(),
  kind: 'learn' as const,
  reason: 'new' as const,
  priority: 1,
  action: { type: 'learn' as const, href: `/learn?concept=${id}` },
  badge: 'New',
  detail: 'Ready to learn',
  cta: 'Start →',
})

const setup = (opts: { missing?: { id: string; label: string; status?: string }[] } = {}) => {
  vi.mocked(getLessonKPs).mockResolvedValue({
    concept_id: CID,
    kps: [{ label: 'Single-digit addition', section: 'S', subgoals: STEPS, worked_example: 'two plus three' }],
    diagram: '',
  })
  vi.mocked(getLessonPractice).mockResolvedValue({
    concept_id: CID,
    questions: [{ question: '2 + 3 = ?', answer: '5', explanation: 'Add them.' }],
  })
  vi.mocked(getLessonReadiness).mockResolvedValue({
    concept_id: CID,
    ready: opts.missing === undefined,
    missing: (opts.missing ?? []).map((m) => ({ ...m, status: m.status ?? 'UNSEEN', mastery_pct: 0 })),
    weak: [],
  })
  vi.mocked(submitStudyAnswer).mockResolvedValue({
    correct: true,
    expected: '5',
    xp: 3,
    band: 'building',
    feedback: 'right',
  } as never)
  vi.mocked(getProgress).mockResolvedValue({})
  vi.mocked(getNext).mockResolvedValue({ primary: rec('b'), alternatives: [], generatedAt: '' })
  vi.mocked(getUserInfo).mockResolvedValue(null)
}

const answerOnce = async () => {
  fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
  await waitFor(() => expect(answerField()).toBeTruthy())
  const input = answerField()
  typeAnswer(input, '5')
  fireEvent.click(within(answerCard(input)).getByRole('button', { name: 'Check' }))
}

describe('/learn chrome', () => {
  beforeEach(() => {
    setup()
  })

  // `level 0.40` was the generator's difficulty float, printed in the learner's header. It is
  // not a goal, a score, or anything they can act on, and it moved with every question, so it
  // read as a fluctuating verdict on them. The band beside it is server-decided and stays.
  it('does not print the difficulty float', async () => {
    render(<LearnStepper conceptId={CID} />)
    expect(await screen.findByText(/Single-digit addition/)).toBeTruthy()
    expect(screen.queryByText(/level \d/)).toBeNull()
    expect(document.body.textContent).not.toMatch(/level 0\.\d\d/)
  })

  // "KP" is what the codebase calls the unit. The learner's unit is a section of the lesson.
  it('names the position in learner words', async () => {
    render(<LearnStepper conceptId={CID} />)
    expect(await screen.findByText('Section 1 of 1')).toBeTruthy()
    expect(document.body.textContent).not.toMatch(/\bKP \d/)
  })

  // The narration instructed a method before the learner had met the problem, and asked them
  // to do something the app cannot see them do.
  it('does not tell the learner to work on paper', async () => {
    render(<LearnStepper conceptId={CID} />)
    expect(await screen.findByText(/Single-digit addition/)).toBeTruthy()
    expect(document.body.textContent).not.toMatch(/on paper/i)
  })

  // Steps listed before the question read as a syllabus, and the effect was to tick them off
  // and skip the work. They now appear once the section has actually been attempted.
  it('reveals the steps only after the section is attempted', async () => {
    render(<LearnStepper conceptId={CID} />)
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    expect(await screen.findByText('2 + 3 = ?')).toBeTruthy()
    expect(screen.queryByText('Steps (2)')).toBeNull()

    await answerOnce()
    expect(await screen.findByText('Steps (2)')).toBeTruthy()
  })

  // The split is right — a prerequisite never learned goes to /learn, one already demonstrated
  // goes to /review, which is problems-first. The verb was wrong: labelling an unseen concept
  // "Review" offered a tutorial under the name of a check.
  it('labels an unlearned prerequisite Learn, and a due one Check', async () => {
    setup({ missing: [{ id: 'arith.mul.single', label: 'Multiplication' }] })
    render(<LearnStepper conceptId={CID} />)
    const banner = await screen.findByRole('note', { name: 'Prerequisite suggestion' })
    expect(within_(banner, /Learn: Multiplication/)).toBeTruthy()
    expect(banner.textContent).not.toMatch(/Review: Multiplication/)
    expect(banner.textContent).not.toMatch(/Check: Multiplication/)
  })
})

// `within` needs the subtree; kept local so the assertion reads as the banner's own text.
function within_(root: HTMLElement, matcher: RegExp): HTMLElement | null {
  const found = Array.from(root.querySelectorAll('a')).find((a) => matcher.test(a.textContent ?? ''))
  return (found as HTMLElement) ?? null
}
