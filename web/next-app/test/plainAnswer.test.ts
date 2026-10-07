// The MathLive → grader bridge, tested against what MathLive actually emits.
//
// The expectations in `plain` below are not guesses: the `raw` column is
// `MathfieldElement.getValue('plain-text')` recorded from mathlive 0.111.0 running
// in Chromium (see test/fixtures/mathlive-plain-text.json). If a future MathLive
// changes its plain-text serialization, this table fails and names the shape that
// moved — which is the point, because that shape is what the server grades.

import { describe, it, expect } from 'vitest'
import {
  fixDivision,
  fixMixedNumbers,
  fixFractionParens,
  fixNameSpacing,
  fixImplicitLetterSpacing,
  toPlainAnswer,
} from '../components/math/plainAnswer'
import corpus from './fixtures/mathlive-plain-text.json'

type Case = (typeof corpus.cases)[number]

describe('toPlainAnswer', () => {
  it('normalizes the shapes the numeric grader cannot parse', () => {
    const rows: [string, string][] = [
      // ASCIIMath division token; without the fix `7 -: 2` reads as subtraction.
      // The spaces MathLive puts around the token survive, which is harmless: the
      // numeric grader strips whitespace before parsing a fraction
      // (numeric.go `parseNumeric`), so `7 / 2` grades as 7/2.
      ['7 -: 2', '7 / 2'],
      // MathLive parenthesises both fraction operands; the numeric grader wants `1/2`.
      ['(1)/(2)', '1/2'],
      ['(30)/(-50)', '30/-50'],
      // A mixed number loses its separator without the fix: `41/10` is not `4 1/10`.
      ['4(1)/(10)', '4 1/10'],
      ['-1(1)/(2)', '-1 1/2'],
      // Space before the argument list hides the function from SymPy's
      // implicit-multiplication pass, which only matches `sin(`.
      ['sin (x)+cos (x)', 'sin(x)+cos(x)'],
      ['sin ^2(x)+cos ^2(x)', 'sin^2(x)+cos^2(x)'],
    ]
    for (const [raw, want] of rows) expect(toPlainAnswer(raw)).toBe(want)
  })

  it('leaves already-acceptable answers untouched', () => {
    for (const good of ['5', '-5', '0.75', '4 1/10', '3.14e-2', '2x+3', 'x+4=10', '6x^2', 'sqrt(9)', '3sqrt(2)', '2:3', '5 R 3', '2,3', '(2,3)', '<']) {
      expect(toPlainAnswer(good)).toBe(good)
    }
  })

  // The recorded corpus, end to end. `numeric` rows are the strict ones: the
  // numeric grader accepts only an integer, a decimal, `n/d` or `n d/d`, so
  // anything else in that group is a false miss.
  const byType = (t: string) => (corpus.cases as Case[]).filter(c => c.gradingType === t)
  const numericCases = byType('numeric')

  it.each(numericCases.map(c => [c.latex, c.plain, toPlainAnswer(c.plain)] as const))(
    'numeric: %s → %s',
    (_latex, _raw, answer) => {
      // Whatever the corpus says the generator expects, this must not silently
      // become an unparseable string. The Go side asserts the grading half.
      expect(answer).not.toContain('-:')
      expect(answer).not.toMatch(/\)\s*\/\s*\(/)
    },
  )

  it('produces the numeric grader’s mixed-number form from a whole-plus-fraction', () => {
    expect(toPlainAnswer('4(1)/(10)')).toMatch(/^-?\d+ \d+\/\d+$/)
  })

  it('does not reorder or reinterpret: every rule is a local substitution', () => {
    // Same input, applied rule by rule — the order is load-bearing (mixed numbers
    // before fraction parens), so assert the chain, not just the end state.
    const start = '4(1)/(10)'
    expect(fixDivision(start)).toBe(start)
    expect(fixMixedNumbers(start)).toBe('4 1/10')
    expect(fixFractionParens('4 1/10')).toBe('4 1/10')
    expect(fixNameSpacing('4 1/10')).toBe('4 1/10')
  })

  it('keeps a space that a learner needs, and only eats the ones MathLive adds', () => {
    expect(fixNameSpacing('4 1/10')).toBe('4 1/10')
    expect(fixNameSpacing('2 x + 1')).toBe('2 x + 1')
    expect(fixNameSpacing('2x + 1')).toBe('2x + 1')
    expect(fixNameSpacing('(x + 1)^2')).toBe('(x + 1)^2')
  })

  it('records the library version the fixture came from', () => {
    // A silent MathLive bump changes the corpus under us; this is the canary.
    expect(corpus.mathliveVersion).toBe('0.111.0')
  })

describe('fixImplicitLetterSpacing', () => {
  it('preserves the word boundary in line segment', () => {
    expect(toPlainAnswer('line segment')).toBe('line segment')
  })

  it('recovers explicit word spaces from MathLive for a bare alphabetic phrase', () => {
    for (const c of corpus.cases.filter(c => c.expected === 'line segment')) {
      expect(toPlainAnswer(c.plain, c.latex)).toBe('line segment')
    }
  })

  it('does not recover word spaces inside mathematical expressions', () => {
    expect(toPlainAnswer('2pi r', '2\\pi\\,r')).toBe('2pi r')
    expect(toPlainAnswer('4(1)/(10)', '4\\frac{1}{10}')).toBe('4 1/10')
    expect(toPlainAnswer('5 R 3', '5\\,R\\,3')).toBe('5 R 3')
    // Explicitly separated variables must remain separate symbols.
    expect(toPlainAnswer('x y', 'x\\,y')).toBe('x y')
    expect(toPlainAnswer('(x + 1)/(y)', '\\frac{x+1}{y}')).toBe('(x + 1)/(y)')
  })
  // MathLive reads bare adjacent letters as implicit multiplication and serialises them
  // with spaces. These are the answer shapes the corpus actually contains -- `yes` is the
  // single most common expected answer in it (673 occurrences), `no` second (171) -- and
  // the multiple-choice grader compares with strings.EqualFold, so `y e s` never matches
  // `yes`. Every value below is a recording from the real library via
  // scripts/record-mathlive-corpus.mjs, not a guess.
  const rows: [string, string][] = [
    ['y e s', 'yes'],
    ['n o', 'no'],
    ['o d d', 'odd'],
    ['a d d i t i o n', 'addition'],
    ['p i/2', 'pi/2'],
    ['x y z', 'xyz'],
  ]
  for (const [raw, want] of rows) {
    it(`rejoins ${JSON.stringify(raw)} as ${JSON.stringify(want)}`, () => {
      expect(fixImplicitLetterSpacing(raw)).toBe(want)
      expect(toPlainAnswer(raw)).toBe(want)
    })
  }

  // The guard is the whole difficulty. A letter-space-letter also appears inside real
  // products, and collapsing one turns `2pi r` into `2pir`, which is a different answer.
  const untouched: [string, string][] = [
    ['2pi r', '2pi r'],
    ['2 xy', '2 xy'],
    ['2 x + y', '2 x + y'],
    ['sin (x)', 'sin (x)'],
    ['x^2 + 1', 'x^2 + 1'],
  ]
  for (const [raw, want] of untouched) {
    it(`leaves the product ${JSON.stringify(raw)} alone`, () => {
      expect(fixImplicitLetterSpacing(raw)).toBe(want)
    })
  }

  it('is idempotent, because it runs after the other rules in a pipeline', () => {
    for (const [raw] of rows) {
      const once = fixImplicitLetterSpacing(raw)
      expect(fixImplicitLetterSpacing(once)).toBe(once)
    }
  })
})

describe('fixFractionParens with a symbolic numerator', () => {
  it('unwraps (pi)/(2) as well as (1)/(2)', () => {
    expect(toPlainAnswer('(1)/(2)')).toBe('1/2')
    expect(toPlainAnswer('(pi)/(2)')).toBe('pi/2')
  })

  it('keeps parentheses on a real expression numerator', () => {
    // Unwrapping this would read as `x + 1/y`, which is a different answer.
    expect(toPlainAnswer('(x + 1)/(y)')).toBe('(x + 1)/(y)')
  })
})
})
