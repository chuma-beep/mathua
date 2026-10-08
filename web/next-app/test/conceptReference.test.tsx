import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import type { KpInfo } from '../lib/api'
import ConceptReference from '../components/ConceptReference'

interface KpOverrides {
  [conceptId: string]: { kps: KpInfo[]; diagram?: string }
}

const { submitReportMock, getDiagram, kpOverride, getLessonsMock } = vi.hoisted(() => ({
  submitReportMock: vi.fn(async () => {}),
  getDiagram: { value: '/diagrams/algebrica/example.svg' as string | null },
  kpOverride: { value: null as null | KpOverrides },
  getLessonsMock: vi.fn(async () => ({ lessons: {} as Record<string, { title: string; concepts: string[] }[]> })),
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
  getLessons: getLessonsMock,
  getLessonBody: async () => 'the whole article',
  submitReport: submitReportMock,
}))

vi.mock('../lib/auth', async importOriginal => ({
  ...(await importOriginal<typeof import('../lib/auth')>()),
  getUserInfo: () => ({ student_id: 's1' }),
}))

vi.mock('next/image', () => ({
  default: (props: { src?: string; alt?: string }) => <img alt={props.alt ?? ''} {...props} />,
}))

beforeEach(() => {
  submitReportMock.mockClear()
  kpOverride.value = null
  getDiagram.value = '/diagrams/algebrica/example.svg'
  getLessonsMock.mockClear()
  getLessonsMock.mockResolvedValue({ lessons: {} })
})

// KatexContent injects a stylesheet into the DOM, and it is full of percentages. Text
// assertions have to look at what a learner reads, not at a stylesheet.
function visibleText(): string {
  const clone = document.body.cloneNode(true) as HTMLElement
  clone.querySelectorAll('style, script').forEach((n) => n.remove())
  return clone.textContent ?? ''
}

async function openFirstMenu() {
  // Wait for the full list: findAllByRole resolves on the first match, which can hand back a
  // node replaced by a later block's render.
  await screen.findByText(/KP two/)
  const menus = await screen.findAllByRole('button', { name: /Report options/ })
  fireEvent.click(menus[0])
}

describe('reference blocks', () => {
  it('shows one report menu per example block, with the diagram attached once', async () => {
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/KP two/)

    expect(screen.queryByText('Report a problem')).toBeNull()
    const menus = await screen.findAllByRole('button', { name: /Report options/ })
    expect(menus).toHaveLength(2)

    // The concept's own diagram rides one block, not all of them: the same figure three times
    // on one screen is a rendering bug, not emphasis.
    expect(screen.getAllByAltText(/Worked diagram/)).toHaveLength(1)
  })

  it('sends concept, kind and block id, so a report can be located', async () => {
    render(<ConceptReference conceptId="c1" />)
    await openFirstMenu()
    fireEvent.click(screen.getByText('Report a problem'))
    fireEvent.click(screen.getByText('Send report'))
    await waitFor(() => expect(submitReportMock).toHaveBeenCalledTimes(1))
    const payload = (submitReportMock.mock.calls[0] as unknown[])[0] as {
      concept_id?: string
      kind?: string
      question?: string
    }
    expect(payload).toMatchObject({ concept_id: 'c1', kind: 'diagram' })
    expect(payload.question).toMatch(/\[block c1\/0\].*example\.svg/)
  })

  it('omits the diagram from every block when a concept has none', async () => {
    getDiagram.value = null
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/KP one/)
    expect(screen.queryByAltText(/Worked diagram/)).toBeNull()
    expect(await screen.findAllByRole('button', { name: /Report options/ })).toHaveLength(2)
  })

  it('collapses a body that repeats under a different label, keeping both names', async () => {
    // A shard can name one lesson section twice under byte-different headings. Rendering
    // both reads the same prose twice, and dropping the second silently loses the label.
    kpOverride.value = {
      c1: {
        kps: [
          { label: 'Adding', section: 'S1', subgoals: [], worked_example: 'shared body' },
          { label: 'Sums', section: 'S1', subgoals: [], worked_example: 'shared body' },
          { label: 'Subtracting', section: 'S2', subgoals: [], worked_example: 'other body' },
        ],
      },
    }
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/Subtracting/)
    expect(document.querySelectorAll('[data-testid="reference-blocks"] details > summary')).toHaveLength(2)
    expect(screen.getAllByText('shared body')).toHaveLength(1)
    expect(screen.getByText(/also Sums/)).toBeTruthy()
  })

  it('does not repeat a label that only differs in formatting', async () => {
    // The dedupe key is (section, body), so two identical bodies collapse; the combined-label
    // guard is a second, separate check. Without it the list reads "also Sums, Sums".
    kpOverride.value = {
      c1: {
        kps: [
          { label: 'Same steps', section: 'S1', subgoals: [], worked_example: 'shared $x$' },
          { label: 'Same steps', section: 'S1', subgoals: [], worked_example: 'shared $x$' },
        ],
      },
    }
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/Same steps/)
    expect(screen.queryByText(/· also/)).toBeNull()
  })

  it('does not collapse two sections that happen to share a body', async () => {
    // Same prose, different sections, is still two things to read. Keying on the body alone
    // would hide one of them.
    kpOverride.value = {
      c1: {
        kps: [
          { label: 'First', section: 'S1', subgoals: [], worked_example: 'shared body' },
          { label: 'Second', section: 'S2', subgoals: [], worked_example: 'shared body' },
        ],
      },
    }
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/Second/)
    expect(document.querySelectorAll('[data-testid="reference-blocks"] details > summary')).toHaveLength(2)
  })
})

describe('reference is not a progression surface', () => {
  it('offers no way to another concept', async () => {
    // The single most important property of the new reference layer. `/study` failed it by
    // design: it had concept chips, a domain index, a search box and a lesson list, so a
    // learner could pick whatever they liked next. Here the only destination is `/learn`,
    // and it carries no concept at all.
    kpOverride.value = { c1: { kps: [{ label: 'KP one', section: 'S1', subgoals: [], worked_example: 'b' }] } }
    render(<ConceptReference conceptId="c1" prerequisites={[{ id: 'p1', label: 'Addition', status: 'UNSEEN' }]} />)
    await screen.findByText(/KP one/)
    // Not "no link off to another concept" but "no link at all": the panel is reached from
    // inside the learning loop, so a way back to it is redundant, and a way onward to
    // something else is the thing being closed.
    expect(screen.queryAllByRole('link')).toHaveLength(0)
  })

  it('lists prerequisites as words, and says they are not progress', async () => {
    // §8 of the closure: prerequisites stay visible, but they are not a link farm. The
    // actionable subset already sits in `/learn`'s own banner with links, so this list is
    // the complete set without duplicating it — and a prerequisite that needs work is sent
    // into the learning flow rather than to a library.
    render(
      <ConceptReference
        conceptId="c1"
        prerequisites={[
          { id: 'p1', label: 'Addition', status: 'UNSEEN' },
          { id: 'p2', label: 'Counting', status: 'MASTERED' },
          { id: 'p3', label: 'Fractions', status: 'DECAYING' },
        ]}
      />,
    )
    await screen.findByText(/KP two/)
    expect(screen.getByText(/Addition/)).toBeTruthy()
    expect(screen.getByText(/not started/)).toBeTruthy()
    expect(screen.getByText(/learned, due a check/)).toBeTruthy()
    expect(screen.getByText(/offered above, in the learning flow/)).toBeTruthy()
  })

  it('states that reading is not evidence, and names no completion control', async () => {
    // §10: no "I studied this", no "mark as learned", no skip. Asserting the absence is the
    // point — the panel is prose next to the questions, which is where such a control would
    // be expected if one ever arrived.
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/KP two/)
    expect(screen.getByText(/only the questions below move this/)).toBeTruthy()
    for (const phrase of [/mark as (read|learned)/i, /i (studied|read) this/i, /skip practice/i, /understand/i]) {
      expect(document.body.textContent).not.toMatch(phrase)
    }
    expect(screen.queryByRole('button', { name: /mark|complete|skip/i })).toBeNull()
  })

  it('renders no mastery percentage or progress quantity', async () => {
    // The rule ADR-021 held for `/study`, and it still holds here: the panel explains, so it
    // must not invite an action against a number. Status words are the part that helps.
    render(<ConceptReference conceptId="c1" prerequisites={[{ id: 'p1', label: 'Addition', status: 'MASTERED' }]} />)
    await screen.findByText(/KP two/)
    expect(visibleText()).not.toMatch(/\d+\s*%/)
    expect(screen.getByText(/learned/)).toBeTruthy()
  })
})

describe('lesson text fallback', () => {
  it('serves the article only for a concept with no knowledge points at all', async () => {
    // Every concept in the corpus has a shard, so this is the rare path — and it is the last
    // place the full article is rendered anywhere in the product.
    kpOverride.value = { c1: { kps: [] } }
    getLessonsMock.mockResolvedValue({ lessons: { arith: [{ title: 'Adding up', concepts: ['c1'] }] } })
    render(<ConceptReference conceptId="c1" />)
    await waitFor(() => expect(screen.getByText(/the whole article/)).toBeTruthy())
  })

  it('does not fetch the lesson index on the common path', async () => {
    render(<ConceptReference conceptId="c1" />)
    await screen.findByText(/KP two/)
    expect(getLessonsMock).not.toHaveBeenCalled()
  })
})