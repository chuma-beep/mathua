// The classify/find pair is the only thing standing between "this corpus has 517
// broken spans" and "this corpus has 517 false accusations", so it gets tested
// against the cases that actually decided the design.

import { describe, it, expect } from 'vitest'
import { ComputeEngine } from '@cortex-js/compute-engine'
import {
  classify,
  classifyResult,
  findError,
  normalizeReason,
  snippet,
} from '../scripts/normalize-latex.mjs'

const ce = new ComputeEngine()
const verdict = (s: string) => classifyResult(ce.parse(s))

describe('normalizeReason', () => {
  it('strips CE serializer tags and token quotes', () => {
    // Both layers are load-bearing: without the tag strip every anchored pattern
    // misses, and without the quote strip the reason still starts with a quote.
    expect(normalizeReason("'unexpected-command'")).toBe('unexpected-command')
    expect(normalizeReason("ErrorCode,'incompatible-type','number'")).toBe(
      'incompatible-type,number',
    )
    expect(normalizeReason("'no-product-between-points','cross-applies'")).toBe(
      'no-product-between-points,cross-applies',
    )
  })
})

describe('classify (pre-parse)', () => {
  it('skips environments, prose-in-math and diagram arrows', () => {
    expect(classify('\\begin{align} a &= b \\end{align}')).toEqual({
      outcome: 'unsupported',
      why: 'environment',
    })
    expect(classify('\\text{if } x > 0')).toEqual({
      outcome: 'unsupported',
      why: 'prose inside math',
    })
    expect(classify('C_0 \\xrightarrow{d} C_1')).toEqual({
      outcome: 'unsupported',
      why: 'diagram arrow',
    })
    expect(classify('\\mathbb{R}^2 \\setminus \\{0\\}')).toEqual({
      outcome: 'unsupported',
      why: 'set membership',
    })
  })

  it('lets real arithmetic through to the parser', () => {
    expect(classify('\\frac{1}{2}+\\frac{1}{3}')).toEqual({ outcome: 'parse' })
    expect(classify('\\int_0^1 x\\,dx')).toEqual({ outcome: 'parse' })
  })
})

describe('findError', () => {
  it('returns null for a clean expression', () => {
    expect(findError(ce.parse('\\frac{1}{2}'))).toBeNull()
  })

  it('finds a failure that is not at the top level', () => {
    // The reason a recursive walk exists: the span as a whole parses, and the
    // damage is buried two nodes down.
    const box = ce.parse('2+\\frac{3}{')
    expect(findError(box)).not.toBeNull()
  })

  it('reports something for a missing box rather than throwing', () => {
    expect(findError(null)).toEqual({ reason: 'parsed to nothing' })
  })
})

describe('classifyResult', () => {
  it('calls clean expressions clean', () => {
    expect(verdict('\\frac{1}{2}+\\frac{1}{3}')).toBeNull()
    expect(verdict('x^2+2x+1')).toBeNull()
    expect(verdict('\\int_0^1 x\\,dx')).toBeNull()
  })

  it('treats an unclosed delimiter as a real defect', () => {
    expect(verdict('\\frac{3}{')?.bucket).toBe('syntax')
    expect(verdict('(f, f(a)))')?.bucket).toBe('syntax')
    expect(verdict('2+')?.bucket).toBe('syntax')
  })

  it('treats a typographic quote in math source as a real defect', () => {
    // 162 of the corpus's 517 defects are this: a Unicode apostrophe from the
    // scrape where LaTeX wants ASCII. KaTeX renders it anyway, which is why no
    // existing gate caught it.
    expect(verdict('f’(x) = 18x - 4')?.bucket).toBe('syntax')
  })

  it('does not accuse the corpus of valid math that CE cannot type', () => {
    // CE parsed these perfectly and then declined to type them as numbers. The
    // whole point of the bucket: reporting these as defects would bury the 517.
    expect(verdict('\\zeta^k,')?.bucket).toBe('not-numeric')
  })

  it('does not accuse the corpus of commands CE simply lacks', () => {
    expect(verdict('(-3,\\;6,\\;-3)\\cdot(4,\\;5,\\;6)')?.bucket).toBe('ce-gap')
  })

  it('never returns a bucket without a reason', () => {
    for (const s of ['\\frac{3}{', '2+', '\\zeta^k,', '(f, f(a)))']) {
      const v = verdict(s)
      expect(v).not.toBeNull()
      expect(v!.reason.length).toBeGreaterThan(0)
      expect(v!.reason.length).toBeLessThanOrEqual(120)
    }
  })
})

describe('snippet', () => {
  it('keeps the head and the tail of a long span', () => {
    // The useful part of a failing span is usually where the parser stopped, which
    // is rarely the first 120 characters. The middle is elided, by design.
    const long = 'HEAD-' + 'x'.repeat(200) + 'TAIL'
    const s = snippet(long, 60)
    expect(s.length).toBeLessThanOrEqual(60)
    expect(s.startsWith('HEAD-')).toBe(true)
    expect(s.endsWith('TAIL')).toBe(true)
    expect(s).toContain('…')
  })

  it('collapses whitespace and leaves short spans alone', () => {
    expect(snippet('  a\n\n  b  ')).toBe('a b')
    expect(snippet('\\frac{1}{2}')).toBe('\\frac{1}{2}')
  })
})
