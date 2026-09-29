import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import type { KpInfo } from '../lib/api'

interface KpStub {
  kps: KpInfo[]
  diagram: string
}

interface KpOverrides {
  [conceptId: string]: KpStub
}

const { submitReportMock, getDiagram, kpOverride } = vi.hoisted(() => ({
  submitReportMock: vi.fn(async () => {}),
  getDiagram: { value: '/diagrams/algebrica/example.svg' as string | null },
  kpOverride: { value: null as null | KpOverrides },
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

vi.mock('next/image', () => ({
  default: (props: { src?: string; alt?: string }) => <img alt={props.alt ?? ''} {...props} />,
}))

import { LessonDetail } from '../app/study/components'

const lesson = {
  title: 'Test lesson',
  body: 'intro $x$',
  concepts: ['c1'],
  prerequisites: [],
}

async function openFirstMenu() {
  // Wait for the full list: findAllByRole resolves on the first match,
  // which can hand back a node replaced by a later block's render.
  await screen.findByText(/KP two/)
  const menuButtons = await screen.findAllByRole('button', { name: /Report options/ })
  fireEvent.click(menuButtons[0])
}

describe('study KP report menus', () => {
  beforeEach(() => {
    submitReportMock.mockClear()
    getDiagram.value = '/diagrams/algebrica/example.svg'
  })

  it('shows one icon menu per example block, diagram attached once', async () => {
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP two/)

    // No visible report text until a ⋯ menu opens.
    expect(screen.queryByText('Report a problem')).toBeNull()
    const menus = await screen.findAllByRole('button', { name: /Report options/ })
    expect(menus).toHaveLength(2)

    // Diagram image renders once beside the first block.
    expect(screen.getByAltText('Worked diagram for c1')).toBeTruthy()
  })

  it('sends concept, kind, block id, and asset path so admins can locate the problem', async () => {
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    await openFirstMenu()
    fireEvent.click(screen.getByText('Report a problem'))
    fireEvent.click(screen.getByText('Send report'))
    await waitFor(() => expect(submitReportMock).toHaveBeenCalledTimes(1))
    const payload = (submitReportMock.mock.calls[0] as unknown[])[0] as { concept_id?: string; kind?: string; question?: string }
    expect(payload).toMatchObject({
      concept_id: 'c1',
      kind: 'diagram',
    })
    expect(payload.question).toMatch(/\[block c1\/0\].*example\.svg/)
    expect(String(payload.question)).toMatch(/\[block c1\/0\].*example\.svg/)
  })

  it('omits the diagram from blocks when a concept has none', async () => {
    getDiagram.value = null
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/KP one/)
    expect(screen.queryByAltText(/Worked diagram/)).toBeNull()
    expect(await screen.findAllByRole('button', { name: /Report options/ })).toHaveLength(2)
  })
})

describe('study KP collapse across concepts', () => {
  const twoConceptLesson = { ...lesson, concepts: ['c1', 'c2'] }

  beforeEach(() => {
    submitReportMock.mockClear()
    kpOverride.value = null
  })

  it('renders a shared worked-example body once with combined labels', async () => {
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
    // Two unique bodies → two collapsible blocks, no pointer links.
    expect(document.querySelectorAll('details > summary')).toHaveLength(2)
    expect(screen.queryByText(/Same worked example as/)).toBeNull()
    expect(screen.getAllByText('shared body')).toHaveLength(1)
    // The repeat's label survives as a combined "also" label.
    expect(screen.getByText(/KP repeat/)).toBeTruthy()
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

  it('renders a shared diagram once with no pointers', async () => {
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
    expect(screen.getAllByAltText(/Worked diagram/)).toHaveLength(1)
    expect(screen.queryByText(/Same diagram as/)).toBeNull()
  })

  it('renders distinct diagrams on their own blocks', async () => {
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
    expect(screen.getAllByAltText(/Worked diagram/)).toHaveLength(2)
    expect(screen.queryByText(/Same diagram as/)).toBeNull()
  })
})
