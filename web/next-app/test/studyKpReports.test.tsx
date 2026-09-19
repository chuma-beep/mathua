import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'

const { submitReportMock, getDiagram } = vi.hoisted(() => ({
  submitReportMock: vi.fn(async () => {}),
  getDiagram: { value: '/diagrams/algebrica/example.svg' as string | null },
}))

vi.mock('../lib/api', async importOriginal => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getLessonKPs: (cid: string) =>
    Promise.resolve({
      concept_id: cid,
      kps: [
        { label: 'KP one', section: 'S1', subgoals: [], worked_example: 'first $x$' },
        { label: 'KP two', section: 'S2', subgoals: [], worked_example: 'second $y$' },
      ],
      diagram: getDiagram.value ?? '',
    }),
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
