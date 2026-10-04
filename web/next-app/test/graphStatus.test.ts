import { describe, it, expect } from 'vitest'
import { deriveStatuses } from '../lib/graphStatus'

const concepts = [
  { id: 'count', prerequisites: [] },
  { id: 'add', prerequisites: ['count'] },
  { id: 'sub', prerequisites: ['add'] },
  { id: 'frac', prerequisites: ['div', 'unknown-prereq'] },
  { id: 'div', prerequisites: ['mul'] },
  { id: 'mul', prerequisites: ['add'] },
]

describe('deriveStatuses', () => {
  it('marks roots as unseen when no progress exists', () => {
    const s = deriveStatuses(concepts)
    expect(s['count']).toBe('unseen')
    expect(s['add']).toBe('locked')
  })

  it('locks concepts whose prerequisites are not mastered', () => {
    const s = deriveStatuses(concepts, { count: { status: 'MASTERED' } })
    expect(s['add']).toBe('unseen')
    expect(s['sub']).toBe('locked')
    expect(s['mul']).toBe('locked')
  })

  it('maps known progress statuses', () => {
    const s = deriveStatuses(concepts, {
      count: { status: 'MASTERED' },
      add: { status: 'PRACTICING' },
      sub: { status: 'LEARNING' },
    })
    expect(s['count']).toBe('mastered')
    expect(s['add']).toBe('practicing')
    expect(s['sub']).toBe('learning')
  })

  it('treats started-but-unmastered prereqs as blocking', () => {
    const s = deriveStatuses(concepts, { add: { status: 'LEARNING' } })
    expect(s['mul']).toBe('locked')
  })

  it('ignores unknown progress entries and unknown prereq ids', () => {
    const s = deriveStatuses(concepts, { 'no-such-concept': { status: 'MASTERED' } })
    expect(s['frac']).toBe('locked')
    const withMulMastered = deriveStatuses(concepts, {
      count: { status: 'MASTERED' },
      add: { status: 'MASTERED' },
      mul: { status: 'MASTERED' },
      div: { status: 'MASTERED' },
    })
    expect(withMulMastered['frac']).toBe('unseen')
  })

  it('treats missing or unrecognized progress status as no progress', () => {
    const s = deriveStatuses(concepts, { count: {} , add: { status: 'WEIRD' } })
    expect(s['count']).toBe('unseen')
    expect(s['add']).toBe('locked')
  })
})

// Decay is computed at read time by the server, so a concept mastered a fortnight ago
// arrives as DECAYING rather than MASTERED. This file used to not know that state at
// all, which is the defect these tests pin.
describe('decaying', () => {
  const withPrereq = () => [
    { id: 'add', prerequisites: [] as string[] },
    { id: 'mul', prerequisites: ['add'] },
    { id: 'div', prerequisites: ['mul'] },
  ]

  it('is its own state, not unseen', () => {
    const s = deriveStatuses(withPrereq(), { add: { status: 'DECAYING' } })
    expect(s['add']).toBe('decaying')
  })

  it('does not relock anything downstream', () => {
    // The failure this replaces: DECAYING was missing from the recognised statuses, so
    // the concept was neither `started` nor `mastered`. It rendered as `unseen` — with
    // no progress bar, as if never attempted — and its absence from the `mastered` set
    // locked its successor. `mul` and `div` are visible again on the graph and in the
    // scheduler at the same time, and disagreeing here is what made that happen.
    const s = deriveStatuses(withPrereq(), {
      add: { status: 'DECAYING' },
      mul: { status: 'MASTERED' },
      div: { status: 'MASTERED' },
    })
    expect(s['mul']).toBe('mastered')
    expect(s['div']).toBe('mastered')
  })

  it('still unlocks an unstarted successor', () => {
    const s = deriveStatuses(withPrereq(), { add: { status: 'DECAYING' } })
    expect(s['mul']).toBe('unseen')
  })

  it('keeps a successor reachable that a merely-practising prerequisite would lock', () => {
    // `div` depends on `mul`, so `mul` is the prerequisite under test.
    const decayed = deriveStatuses(withPrereq(), {
      add: { status: 'MASTERED' },
      mul: { status: 'DECAYING' },
    })
    const practising = deriveStatuses(withPrereq(), {
      add: { status: 'MASTERED' },
      mul: { status: 'PRACTICING' },
    })
    expect(decayed['div']).toBe('unseen')
    expect(practising['div']).toBe('locked')
  })
})
