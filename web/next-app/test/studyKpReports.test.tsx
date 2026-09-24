import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import type { KpInfo } from '../lib/api'

interface KpStub {
  kps: KpInfo[]
  diagram: string
}

const { submitReportMock, getDiagram, kpOverride } = vi.hoisted(() => ({
  submitReportMock: vi.fn(async () => {}),
  getDiagram: { value: '/diagrams/algebrica/example.svg' as string | null },
  kpOverride: { value: null as null | Record<string, KpStub> },
}))

vi.mock('../lib/api', async importOriginal => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getLessonKPs: (cid: string) => {
    const over = kpOverride.value?.[cid]
    if (over) return Promise.resolve({ concept_id: cid, ...over })
    return Promise.resolve({
      concept_id: cid,
      kps: [
        { label: 'KP one', section: 'S1', subgoals: [], worked_example: 'first $x$' },
        { label: 'KP two', section: 'S2', subgoals: [], worked_example: 'second $y$' },
      ],
      diagram: getDiagram.value ?? '',
    })
  },
  submitReport: submitReportMock,
}))

vi.mock('../components/LessonQuiz', () => ({
  default: () => <div data-testid="quiz-stub" />,
}))

vi.mock('next/image', () => ({
  default: (props: Record<string, unknown>) => <img alt="" {...props} />,
}))

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))

import { LessonDetail } from '../app/study/components'

const lesson = {
  title: 'Test lesson',
  body: 'intro $x$',
  concepts: ['c1'],
  prerequisites: [],
}

describe('study KP report buttons', () => {
  beforeEach(() => {
    submitReportMock.mockClear()
    getDiagram.value = '/diagrams/algebrica/example.svg'
  })

  it('hoists the diagram report to one labeled button per concept', async () => {
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)

    // One diagram button total (not one per KP card), with a distinct label.
    const diagramButtons = await screen.findAllByText('Report diagram problem')
    expect(diagramButtons).toHaveLength(1)

    // Each KP card keeps exactly its worked-example button.
    expect(screen.getAllByText('Report a problem')).toHaveLength(2)

    // Diagram image renders once beside it.
    expect(screen.getByAltText('Worked diagram for c1')).toBeTruthy()
  })

  it('sends concept, kind, and asset path so admins can locate the problem', async () => {
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    fireEvent.click(await screen.findByText('Report diagram problem'))
    fireEvent.click(screen.getByText('Send report'))
    await waitFor(() => expect(submitReportMock).toHaveBeenCalledTimes(1))
    const payload = (submitReportMock.mock.calls[0] as unknown[])[0] as Record<string, unknown>
    expect(payload).toMatchObject({
      concept_id: 'c1',
      kind: 'diagram',
      question: '/diagrams/algebrica/example.svg',
    })
  })

  it('omits the diagram block entirely when a concept has no diagram', async () => {
    getDiagram.value = null
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP one/)
    expect(screen.queryByText('Report diagram problem')).toBeNull()
    // Worked-example buttons are unaffected.
    expect(screen.getAllByText('Report a problem')).toHaveLength(2)
  })
})

describe('study KP dedupe across concepts', () => {
  const twoConceptLesson = { ...lesson, concepts: ['c1', 'c2'] }

  beforeEach(() => {
    submitReportMock.mockClear()
    kpOverride.value = null
  })

  it('renders a shared worked-example body once and points repeats back', async () => {
    kpOverride.value = {
      c1: {
        kps: [{ label: 'KP one', section: 'S1', subgoals: [], worked_example: 'shared body' }],
        diagram: '',
      },
      c2: {
        kps: [
          { label: 'KP repeat', section: 'S1', subgoals: [], worked_example: 'shared body' },
          { label: 'KP fresh', section: 'S2', subgoals: [], worked_example: 'other body' },
        ],
        diagram: '',
      },
    }
    render(<LessonDetail lesson={twoConceptLesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP fresh/)
    // Full "Worked example" expanders: one for the shared body, one fresh.
    // (Per-concept header badges also read "Worked example"; count summaries.)
    const summaries = document.querySelectorAll('details > summary')
    expect(summaries).toHaveLength(2)
    // The repeat keeps its label and points back instead of re-rendering.
    expect(screen.getByText(/KP repeat/)).toBeTruthy()
    expect(screen.getByText(/Same worked example as/)).toBeTruthy()
  })

  it('renders distinct bodies fully with no pointers', async () => {
    kpOverride.value = {
      c1: {
        kps: [{ label: 'KP one', section: 'S1', subgoals: [], worked_example: 'first body' }],
        diagram: '',
      },
      c2: {
        kps: [{ label: 'KP two', section: 'S2', subgoals: [], worked_example: 'second body' }],
        diagram: '',
      },
    }
    render(<LessonDetail lesson={twoConceptLesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP two/)
    expect(document.querySelectorAll('details > summary')).toHaveLength(2)
    expect(screen.queryByText(/Same worked example as/)).toBeNull()
  })

  it('renders a shared diagram once and points repeats back', async () => {
    kpOverride.value = {
      c1: {
        kps: [{ label: 'KP one', section: 'S1', subgoals: [], worked_example: 'first body' }],
        diagram: '/diagrams/algebrica/shared.svg',
      },
      c2: {
        kps: [{ label: 'KP two', section: 'S2', subgoals: [], worked_example: 'second body' }],
        diagram: '/diagrams/algebrica/shared.svg',
      },
    }
    render(<LessonDetail lesson={twoConceptLesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP two/)
    expect(screen.getAllByText('Report diagram problem')).toHaveLength(1)
    expect(screen.getByText(/Same diagram as/)).toBeTruthy()
  })

  it('renders distinct diagrams fully with no pointers', async () => {
    kpOverride.value = {
      c1: {
        kps: [{ label: 'KP one', section: 'S1', subgoals: [], worked_example: 'first body' }],
        diagram: '/diagrams/algebrica/one.svg',
      },
      c2: {
        kps: [{ label: 'KP two', section: 'S2', subgoals: [], worked_example: 'second body' }],
        diagram: '/diagrams/algebrica/two.svg',
      },
    }
    render(<LessonDetail lesson={twoConceptLesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP two/)
    expect(screen.getAllByText('Report diagram problem')).toHaveLength(2)
    expect(screen.queryByText(/Same diagram as/)).toBeNull()
  })
})
