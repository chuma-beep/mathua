import { describe, it, expect } from 'vitest'
import { selectNextUp, type NextUpInput } from '../lib/nextUp'

const base: NextUpInput = {
  dueReviews: 0,
  weaknesses: { by_domain: {} },
  progress: {},
  activity: [],
  diagnosticCompleted: false,
  conceptsMastered: 0,
}

describe('selectNextUp', () => {
  it('prefers due reviews', () => {
    const n = selectNextUp({ ...base, dueReviews: 3 })
    expect(n.kind).toBe('review')
    expect(n.href).toBe('/session')
  })

  it('recommends weakest non-mastered concept via supported ?concept= link', () => {
    const n = selectNextUp({
      ...base,
      conceptsMastered: 2,
      diagnosticCompleted: true,
      weaknesses: {
        by_domain: {
          fractions: [
            { id: 'frac.add', label: 'Add fractions', weakness: 0.8 },
            { id: 'frac.sub', label: 'Subtract fractions', weakness: 0.4 },
          ],
        },
      },
      progress: { 'frac.add': { status: 'learning', streak: 1 } },
    })
    expect(n.kind).toBe('weakness')
    expect(n.href).toBe('/study?concept=frac.add')
  })

  it('skips mastered weaknesses', () => {
    const n = selectNextUp({
      ...base,
      conceptsMastered: 5,
      diagnosticCompleted: true,
      weaknesses: { by_domain: { arithmetic: [{ id: 'arith.x', label: 'X', weakness: 0.9 }] } },
      progress: { 'arith.x': { status: 'MASTERED', streak: 2 } },
    })
    expect(n.kind).toBe('browse')
  })

  it('resumes most recent in-progress concept', () => {
    const n = selectNextUp({
      ...base,
      conceptsMastered: 3,
      diagnosticCompleted: true,
      progress: { 'alg.lin': { status: 'learning', streak: 1 } },
      activity: [
        { date: '2026-09-01', questions: 4, correct: 3, concepts: ['arith.old'] },
        { date: '2026-09-02', questions: 2, correct: 1, concepts: ['alg.lin'] },
      ],
    })
    expect(n.kind).toBe('resume')
    expect(n.href).toBe('/study?concept=alg.lin')
  })

  it('sends brand-new users to Diagnostic', () => {
    const n = selectNextUp(base)
    expect(n.kind).toBe('diagnostic')
    expect(n.href).toBe('/onboard')
  })

  it('falls back to Study library', () => {
    const n = selectNextUp({ ...base, conceptsMastered: 4, diagnosticCompleted: true })
    expect(n.kind).toBe('browse')
    expect(n.href).toBe('/study')
  })
})
