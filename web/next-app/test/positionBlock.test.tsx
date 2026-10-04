import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import PositionBlock from '../components/PositionBlock'
import { countOverall, type CatalogueConcept } from '../lib/progress'
import type { ConceptProgress } from '../lib/api'

const catalogue: CatalogueConcept[] = [
  { id: 'a', domain: 'arithmetic', prerequisites: [] },
  { id: 'b', domain: 'arithmetic', prerequisites: ['a'] },
  { id: 'c', domain: 'arithmetic', prerequisites: ['a'] },
]

const p = (status: string): ConceptProgress => ({ status, streak: 3 })

const renderBlock = (progress: Record<string, ConceptProgress>) =>
  render(
    <PositionBlock
      catalogue={catalogue}
      progress={progress}
      frontierLabel=""
      frontierHref=""
    />,
  )
const counts = (progress: Record<string, ConceptProgress>) => countOverall(catalogue, progress)

describe('PositionBlock', () => {
  it('says nothing about review when nothing has decayed', () => {
    const { container } = renderBlock({ a: p('MASTERED') })
    expect(screen.queryByText(/due for review/i)).toBeNull()
    expect(container.textContent).toContain('1')
  })

  // The brief's §11 problem: "mastered" was quietly excluding decayed concepts, so this
  // headline disagreed with the graph's full bars. The fix was not to shrink the number
  // but to report the review debt in its own right — which only helps if it is visible.
  it('reports the review debt separately from the mastered count', () => {
    const c = counts({ a: p('DECAYING'), b: p('MASTERED') })
    // Both count as mastered, because both were demonstrated.
    expect(c.mastered).toBe(2)
    expect(c.dueForReview).toBe(1)

    renderBlock({ a: p('DECAYING'), b: p('MASTERED') })
    expect(screen.getByText(/1 due for review/i)).toBeTruthy()
    // And it says plainly that decay is not forgetting, so a learner reading a review
    // prompt next to a "2 of 3 mastered" headline is not left guessing which moved.
    expect(screen.getByText(/still learned/i)).toBeTruthy()
    expect(screen.getByRole('link', { name: /review/i }).getAttribute('href')).toBe('/review')
  })
})
