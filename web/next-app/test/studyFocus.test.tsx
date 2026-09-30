import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { LessonDetail } from '../app/study/components'
import { getLessonKPs } from '../lib/api'
import { resolveConceptTarget, planConceptNavigation } from '../lib/conceptTarget'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, getLessonKPs: vi.fn() }
})

vi.mock('next/image', () => ({
  default: (props: { src?: string; alt?: string }) => <img alt={props.alt ?? ''} {...props} />,
}))

const scrollIntoViewMock = vi.fn()
Object.defineProperty(Element.prototype, 'scrollIntoView', {
  value: scrollIntoViewMock,
  writable: true,
  configurable: true,
})

const lesson = {
  title: 'Lesson L',
  body: '',
  concepts: ['a', 'b'],
  prerequisites: [],
  dependents: [],
}

beforeEach(() => {
  vi.mocked(getLessonKPs).mockImplementation(async (cid: string) => ({
    concept_id: cid,
    kps: [{ label: `${cid} KP`, section: 'S', subgoals: [], worked_example: `body ${cid} $x$` }],
    diagram: '',
  }))
  vi.mocked(scrollIntoViewMock).mockClear()
})

describe('resolveConceptTarget / planConceptNavigation', () => {
  const byConcept = new Map([
    ['a', { title: 'Lesson L' }],
    ['z', { title: 'Other lesson' }],
  ])

  it('same lesson resolves self', () => {
    expect(resolveConceptTarget('a', byConcept, 'Lesson L')).toEqual({
      href: '/study?concept=a',
      sameLesson: true,
    })
  })

  it('other lesson resolves cross-lesson', () => {
    expect(resolveConceptTarget('z', byConcept, 'Lesson L').sameLesson).toBe(false)
  })

  it('unknown concept resolves cross-lesson (URL sync no-ops)', () => {
    expect(resolveConceptTarget('zzz', byConcept, 'Lesson L').sameLesson).toBe(false)
  })

  it('plans replace for a fresh in-lesson jump', () => {
    expect(planConceptNavigation('a', byConcept, 'Lesson L', null)).toEqual({
      action: 'replace',
      href: '/study?concept=a',
    })
  })

  it('plans push across lessons', () => {
    expect(planConceptNavigation('z', byConcept, 'Lesson L', 'a').action).toBe('push')
  })

  it('plans direct reveal on re-click (URL unchanged)', () => {
    expect(planConceptNavigation('a', byConcept, 'Lesson L', 'a')).toEqual({
      action: 'reveal',
      href: '/study?concept=a',
    })
  })
})

describe('LessonDetail concept reveal', () => {
  it('chip click delegates navigation to the page handler', () => {
    const onConceptSelect = vi.fn()
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} onConceptSelect={onConceptSelect} />)
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    const chipB = links.find((l) => l.getAttribute('href') === '/study?concept=b')
    expect(chipB).toBeTruthy()
    fireEvent.click(chipB!)
    expect(onConceptSelect).toHaveBeenCalledTimes(1)
    expect(onConceptSelect.mock.calls[0][0]).toBe('b')
  })

  it('reveal request opens the concept block, highlights and scrolls', async () => {
    const { rerender } = render(
      <LessonDetail lesson={lesson} domain={null} onBack={() => {}} revealReq={null} />,
    )
    await screen.findByText(/a KP/)
    rerender(
      <LessonDetail
        lesson={lesson}
        domain={null}
        onBack={() => {}}
        revealReq={{ cid: 'b', n: 1 }}
      />,
    )
    // b's <details> opens (a's was seeded open already).
    await waitFor(() => {
      const details = Array.from(document.querySelectorAll('details')) as HTMLDetailsElement[]
      expect(details.some((d) => d.open && d.textContent?.includes('b KP'))).toBe(true)
    })
    // Highlight + scroll ride requestAnimationFrame; wait a frame.
    await waitFor(() => {
      expect(document.querySelector('.ring-mathua-blue')).not.toBeNull()
    })
    expect(scrollIntoViewMock).toHaveBeenCalled()
  })

  it('no-block concept falls back to the lesson header', async () => {
    render(
      <LessonDetail
        lesson={lesson}
        domain={null}
        onBack={() => {}}
        revealReq={{ cid: 'zzz', n: 1 }}
      />,
    )
    await screen.findByText(/a KP/)
    // Header highlight, no crash, scroll still happens.
    await waitFor(() => {
      expect(document.querySelector('.ring-mathua-blue')).not.toBeNull()
    })
    expect(scrollIntoViewMock).toHaveBeenCalled()
  })
})
