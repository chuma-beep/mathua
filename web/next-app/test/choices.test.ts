import { describe, it, expect } from 'vitest'
import { parseChoices } from '../lib/choices'

describe('parseChoices', () => {
  it('parses trailing "A) …" option blocks', () => {
    const q = 'Which shape?\n\nA) circle\n\nB) square\n\nC) triangle'
    expect(parseChoices(q)).toEqual([
      { letter: 'A', text: 'circle' },
      { letter: 'B', text: 'square' },
      { letter: 'C', text: 'triangle' },
    ])
  })

  it('returns null when there is no option list', () => {
    expect(parseChoices('Simplify: \\((5^2)^2\\)')).toBeNull()
    expect(parseChoices('2 + 2 = ?')).toBeNull()
  })

  it('returns null when letters are not sequential from A', () => {
    expect(parseChoices('Q\n\nA) x\n\nC) y')).toBeNull()
  })

  it('preserves math in option text', () => {
    const q = 'Pick one\n\nA) $x^2$\n\nB) $x^3$'
    expect(parseChoices(q)).toEqual([
      { letter: 'A', text: '$x^2$' },
      { letter: 'B', text: '$x^3$' },
    ])
  })
})
