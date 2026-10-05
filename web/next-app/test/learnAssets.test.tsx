import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
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
import type { LessonAsset } from '../lib/api'

vi.mock('next/image', () => ({
  default: ({ alt, src }: { alt?: string; src?: string }) => <img alt={alt ?? ''} data-src={src} />,
}))

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

Element.prototype.scrollIntoView = vi.fn() as unknown as typeof Element.prototype.scrollIntoView

const CID = 'trig.unit-circle'

const fetchMock = vi.fn<(url: string) => Promise<unknown>>()
vi.stubGlobal('fetch', fetchMock)

const svgDoc =
  '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><circle cx="50" cy="50" r="40"/></svg>'

// An asset exactly as `internal/lessons` derives it: `src` is the raw markdown reference,
// not a served URL. That is the whole reason the client has to resolve it.
const unitCircle: LessonAsset = {
  id: 'unit-circle-labeled.svg',
  kind: 'image',
  src: 'svg/unit-circle-labeled.svg',
  alt: 'Unit circle with notable angles',
  line: 4,
  section: 'unit-circle',
  concept_ids: [CID],
}

const tangent: LessonAsset = {
  id: 'unit-circle-tangent.svg',
  kind: 'image',
  src: '/diagrams/algebrica/unit-circle-tangent.svg',
  alt: 'Tangent line to the unit circle',
  line: 9,
  section: 'tangent-line',
}

function kpsFor(assets?: LessonAsset[]) {
  return {
    concept_id: CID,
    diagram: '/diagrams/algebrica/concept-diagram.svg',
    kps: [
      { label: 'The unit circle', subgoals: [], worked_example: 'A point on the circle.', assets },
      { label: 'The tangent line', subgoals: [], worked_example: 'The tangent is perpendicular.', assets: [tangent] },
    ],
  }
}

async function renderStepper() {
  const LearnStepper = (await import('../components/LearnStepper')).default
  return render(<LearnStepper conceptId={CID} />)
}

beforeEach(() => {
  vi.mocked(getLessonKPs).mockResolvedValue(kpsFor([unitCircle]) as never)
  vi.mocked(getLessonPractice).mockResolvedValue({ questions: [], concept_id: CID })
  vi.mocked(getLessonReadiness).mockResolvedValue({ concept_id: CID, ready: true, weak: [], missing: [] })
  vi.mocked(getActivity).mockResolvedValue({ days: [], total: 0, streak: 0 } as never)
  vi.mocked(getProgress).mockResolvedValue({} as never)
  vi.mocked(getWeaknesses).mockResolvedValue({ weaknesses: [] } as never)
  vi.mocked(getDueReviews).mockResolvedValue([] as never)
  vi.mocked(getScores).mockResolvedValue({} as never)
  vi.mocked(submitStudyAnswer).mockResolvedValue({ correct: true, xp: 1 } as never)
  vi.mocked(getUserInfo).mockResolvedValue(null as never)
  fetchMock.mockReset()
  fetchMock.mockImplementation(async (url: string) => {
    if (url === '/diagrams/meta.json') {
      return { ok: true, status: 200, text: async () => '{}', json: async () => ({}) }
    }
    return { ok: true, status: 200, text: async () => svgDoc }
  })
})

describe('/learn figures', () => {
  it('shows the concept diagram the endpoint already returned', async () => {
    // The KP response has always carried `diagram` and the client has always discarded it,
    // so the one figure guaranteed relevant to the concept being studied was the one figure
    // /learn never showed.
    const { container } = await renderStepper()
    await waitFor(() => expect(getLessonKPs).toHaveBeenCalled())
    await waitFor(() =>
      expect(container.querySelector('[data-src="/diagrams/algebrica/concept-diagram.svg"]')).toBeTruthy(),
    )
  })

  it('shows the figures belonging to the current knowledge point', async () => {
    const { container } = await renderStepper()
    await waitFor(() => expect(container.textContent).toContain('The unit circle'))
    await waitFor(() =>
      // A relative markdown reference is resolved to the served URL, exactly as the prose
      // around it is. One rule, or the same figure renders in one place and 404s in another.
      expect(
        container.querySelector('[data-src="/diagrams/algebrica/unit-circle-labeled.svg"]'),
      ).toBeTruthy(),
    )
  })

  it('does not show another knowledge point\'s figures on the first card', async () => {
    // Selection is the point: a shorter surface should carry less irrelevant text, not less
    // instruction. The tangent figure belongs to KP 2 and must not leak into KP 1.
    const { container } = await renderStepper()
    await waitFor(() => expect(container.textContent).toContain('The unit circle'))
    expect(container.querySelector('[data-src="/diagrams/algebrica/unit-circle-tangent.svg"]')).toBeNull()
  })

  it('uses the asset alt text, and falls back to the id when there is none', async () => {
    const bare: LessonAsset = { id: 'floor-and-ceiling-functions-1.svg', kind: 'image', src: 'svg/x.svg', line: 2 }
    vi.mocked(getLessonKPs).mockResolvedValue(kpsFor([bare]) as never)
    const { container } = await renderStepper()
    await waitFor(() => expect(container.textContent).toContain('The unit circle'))
    // 133 of the corpus's figures carry neither alt text nor a caption. Falling back to the
    // identifier keeps the image addressable rather than rendering an unlabelled graphic.
    await waitFor(() =>
      expect(container.querySelector('img[alt="floor-and-ceiling-functions-1.svg"]')).toBeTruthy(),
    )
  })

  it('renders nothing for a lesson with no figures', async () => {
    vi.mocked(getLessonKPs).mockResolvedValue(kpsFor([]) as never)
    const { container } = await renderStepper()
    await waitFor(() => expect(container.textContent).toContain('The unit circle'))
    expect(container.querySelectorAll('figure')).toHaveLength(1) // the concept diagram only
  })

  it('survives a malformed asset without losing the lesson', async () => {
    // `src` is optional in the type, and `table` is a kind with no renderer yet. Neither may
    // take the learning loop down with it.
    vi.mocked(getLessonKPs).mockResolvedValue(
      kpsFor([{ id: 'no-src', kind: 'image', line: 1 }, { id: 'a-table', kind: 'table', line: 2 }]) as never,
    )
    vi.mocked(getLessonPractice).mockResolvedValue({
      concept_id: CID,
      questions: [{ question: '2 + 3 = ?', answer: '5', explanation: 'Add them.' }],
    })
    const { container } = await renderStepper()
    // The loop must still be live: a lesson can answer every question it is asked.
    fireEvent.click(await screen.findByRole('button', { name: 'Next →' }))
    await waitFor(() => expect(container.textContent).toContain('2 + 3 = ?'))
    expect(screen.getByRole('button', { name: 'Check' })).toBeTruthy()
  })
})