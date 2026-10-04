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

describe('LessonDetail concept reveal', () => {  it('chip click delegates navigation to the page handler', () => {
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

describe('LessonDetail dedupe and single Start', () => {
  it('renders each concept chip once for repeated concept ids', () => {
    render(
      <LessonDetail
        lesson={{ ...lesson, concepts: ['a', 'a', 'b'] }}
        domain={null}
        onBack={() => {}}
      />,
    )
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    expect(links.filter((l) => l.getAttribute('href') === '/study?concept=a')).toHaveLength(1)
    expect(links.filter((l) => l.getAttribute('href') === '/study?concept=b')).toHaveLength(1)
  })

  it('collapses repeated worked-example labels instead of also X, X', async () => {
    vi.mocked(getLessonKPs).mockImplementation(async (cid: string) => ({
      concept_id: cid,
      kps: [
        { label: 'Same steps', section: 'S', subgoals: [], worked_example: 'shared $x$' },
        { label: 'Same steps', section: 'S', subgoals: [], worked_example: 'shared $x$' },
      ],
      diagram: '',
    }))
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    await screen.findByText(/Same steps/)
    expect(screen.queryByText(/· also/)).toBeNull()
  })

  it('shows one primary Start for the first concept still being learned', () => {
    render(
      <LessonDetail
        lesson={{ ...lesson, progress: { a: { status: 'MASTERED', streak: 3 } } }}
        domain={null}
        onBack={() => {}}
      />,
    )
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    const starts = links.filter((l) => (l.getAttribute('href') ?? '').startsWith('/learn?concept='))
    // Primary heads b, since a is already learned.
    expect(starts.filter((l) => l.getAttribute('href') === '/learn?concept=b&return=b')).toHaveLength(1)
    // Exactly one primary row; the rest live in the disclosure.
    expect(document.querySelectorAll('a.justify-between')).toHaveLength(1)

    // `a` is no longer offered as a Start target. It used to be — "a remains only under
    // the disclosure" — and every one of those links is a route into /learn, which opens
    // with the lesson's worked example unconditionally. Offering a concept the learner has
    // already demonstrated as something to start learning is the tutorial purgatory; the
    // reference itself is still one click away via the concept chips.
    expect(starts.filter((l) => l.getAttribute('href') === '/learn?concept=a&return=a')).toHaveLength(0)
    // And with nothing learnable left besides the head, the disclosure is gone rather than
    // present and empty — a count in "Or pick something else (N)" that does not match the
    // N links under it is its own kind of lie.
    expect(screen.queryByText(/Or pick something else/)).toBeNull()
  })

  it('does not offer Start learning at all when every concept is already learned', () => {
    render(
      <LessonDetail
        lesson={{ ...lesson, progress: { a: { status: 'MASTERED', streak: 3 }, b: { status: 'MASTERED', streak: 3 } } }}
        domain={null}
        onBack={() => {}}
      />,
    )
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    const starts = links.filter((l) => (l.getAttribute('href') ?? '').startsWith('/learn?concept='))

    // The head used to fall back to `learnTargets[0]` here — so a fully-mastered lesson
    // presented a concept the learner had already demonstrated as its primary call to
    // action, labelled "Start learning".
    expect(starts).toHaveLength(0)
    expect(screen.getByText(/Every concept here is learned/)).toBeTruthy()
  })

  it('treats a decayed concept as learned, not as a Start target', () => {
    // The server reports stale mastery as DECAYING rather than MASTERED, so a strict
    // `status === 'MASTERED'` test reads a fortnight's break as "not started" and offers
    // the tutorial again. Decay is a retrieval question, not a teaching one.
    render(
      <LessonDetail
        lesson={{ ...lesson, progress: { a: { status: 'DECAYING', streak: 3 }, b: { status: 'DECAYING', streak: 3 } } }}
        domain={null}
        onBack={() => {}}
      />,
    )
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    expect(links.filter((l) => (l.getAttribute('href') ?? '').startsWith('/learn?concept='))).toHaveLength(0)
    expect(screen.getByText(/Every concept here is learned/)).toBeTruthy()
  })

  it('heads the first concept when nothing is mastered', () => {
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    const starts = links.filter((l) => (l.getAttribute('href') ?? '').startsWith('/learn?concept='))
    expect(starts.filter((l) => l.getAttribute('href') === '/learn?concept=a&return=a')).toHaveLength(1)
    expect(starts.filter((l) => l.getAttribute('href') === '/learn?concept=b&return=b')).toHaveLength(1)
    expect(screen.getByText('Or pick something else (1)')).toBeTruthy()
  })
})

// The Study → Learn → Study round trip. Without a way back, a learner who
// clicked "Start learning" to check something had no way to return to the lesson
// they were reading, and `?return=` — which Learn already understood — was
// only ever set by Learn's own prerequisite links.
describe('LessonDetail Learn round trip', () => {
  it('offers a way back to practice when the learner arrived from Learn', () => {
    render(
      <LessonDetail
        lesson={lesson}
        domain={null}
        fromConcept="arith.add.single"
        onBack={() => {}}
      />,
    )
    // Both the primary row and each alternative offer the way back.
    expect(screen.getAllByText(/Back to practice/).length).toBeGreaterThan(1)
    expect(screen.queryByText('Start learning →')).toBeNull()
    // And the label says which concept, not just "practice".
    expect(screen.getByText(/Back to [A-Z]/)).toBeTruthy()
  })

  it('offers only Start when the learner arrived from navigation', () => {
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    expect(screen.getAllByText('Start learning →').length).toBeGreaterThan(0)
    expect(screen.queryByText(/Back to practice/)).toBeNull()
  })

  it('does not claim there is an interactive loop on the page', () => {
    // This section contains links and nothing else; describing it as the loop
    // invited a learner to look for a question that is not here.
    render(<LessonDetail lesson={lesson} domain={null} onBack={() => {}} />)
    expect(screen.queryByText(/interactive loop/)).toBeNull()
    expect(screen.queryByText(/answer questions on any of these/)).toBeTruthy()
  })
})
