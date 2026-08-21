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
