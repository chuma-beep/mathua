// The equivalence primitive, against the real Compute Engine.
//
// These are not unit tests of a mock: Compute Engine is installed and imported
// directly, because every assertion here is about its actual behaviour, including
// the parts that are inconvenient. The four edge cases this file exists for were each
// measured first and the measurement is what got asserted — see the note on each.

import { describe, it, expect, beforeAll } from 'vitest'
import { ComputeEngine } from '@cortex-js/compute-engine'
import { loadEngine, equivalence, simplificationHint, type Engine, resetEngineForTests } from '../components/math/equivalence'

let engine: Engine

beforeAll(async () => {
  resetEngineForTests()
  const e = await loadEngine()
  if (!e) throw new Error('Compute Engine failed to load; these tests are about its real behaviour')
  engine = e
})

describe('equivalence — the four edge cases', () => {
  it('implicit multiplication: 2x+3 is 2*x+3', () => {
    expect(equivalence(engine, '2x+3', '2*x+3')).toBe('equivalent')
    expect(equivalence(engine, '3(4+5)', '3*9')).toBe('equivalent')
    expect(equivalence(engine, '2x+3', '2x+4')).not.toBe('equivalent')
  })

  it('differently factored: (x+1)² is x²+2x+1', () => {
    // `isEqual` alone returns `undefined` here — it compares values and cannot
    // establish a symbolic identity. That is why the primitive uses
    // `isIdenticallyEqual` and falls back to canonical forms.
    expect(equivalence(engine, '(x+1)^2', 'x^2+2x+1')).toBe('equivalent')
    expect(equivalence(engine, '2x+2', '2(x+1)')).toBe('equivalent')
  })

  it('trig identities: sin²x + cos²x is 1', () => {
    // Also `undefined` from `isIdenticallyEqual` on its own; the simplify fallback
    // is what closes this one.
    expect(equivalence(engine, '\\sin(x)^2+\\cos(x)^2', '1')).toBe('equivalent')
    expect(equivalence(engine, '\\sin(2x)^2+\\cos(2x)^2', '1')).toBe('equivalent')
  })

  it('domain restrictions: sqrt(x²) is NOT x', () => {
    // The important one. `sqrt(x^2)` equals `|x|`, not `x`, so the primitive must
    // refuse rather than confirm. It returns 'unknown', never 'equivalent'.
    expect(equivalence(engine, '\\sqrt{x^2}', 'x')).not.toBe('equivalent')
    // …and it knows the right answer when asked.
    expect(equivalence(engine, '\\sqrt{x^2}', '|x|')).toBe('equivalent')
  })

  it('reports a proof, not a guess', () => {
    // Everything it can prove equal…
    expect(equivalence(engine, '\\frac{1}{2}+\\frac{1}{3}', '\\frac{5}{6}')).toBe('equivalent')
    expect(equivalence(engine, '0.9', '\\frac{9}{10}')).toBe('equivalent')
    expect(equivalence(engine, '\\sqrt{9}', '3')).toBe('equivalent')
    expect(equivalence(engine, '\\frac{2}{4}', '\\frac{1}{2}')).toBe('equivalent')
    expect(equivalence(engine, '2^10', '1024')).toBe('different')
  })

  it('returns unknown rather than guessing when it cannot decide', () => {
    // CE declines to call these different. Round them to "different" and a learner
    // would be told their expression is wrong on the tool's say-so alone.
    expect(equivalence(engine, 'x^2+1', 'x^2+2')).toBe('unknown')
  })

  it('never throws on input it cannot parse', () => {
    for (const bad of ['', '   ', '(((', 'not math at all', '\\frac{1}{']) {
      expect(equivalence(engine, bad, 'x')).toBe('unknown')
      expect(simplificationHint(engine, bad)).toBeNull()
    }
  })
})

describe('simplificationHint', () => {
  it('offers the exact value of a constant expression', () => {
    expect(simplificationHint(engine, '\\frac{1}{2}+\\frac{1}{3}')?.latex).toBe('5/6')
    expect(simplificationHint(engine, '\\sqrt{9}')?.latex).toBe('3')
    expect(simplificationHint(engine, '3(4+5)')?.latex).toBe('27')
    expect(simplificationHint(engine, 'x^2+x^2')?.latex).toBe('2x^2')
  })

  it('does not offer a fraction the learner already wrote in lowest-looking form', () => {
    // Rewriting \frac{1}{2} as "1/2" is a different rendering, not a simplification.
    expect(simplificationHint(engine, '\\frac{1}{2}')).toBeNull()
  })

  it('gives |x| for sqrt(x²), not x', () => {
    // The reason this uses a CAS at all. A regex would happily return `x`, which is
    // wrong for every negative x — and this is the exact mistake a learner makes.
    expect(simplificationHint(engine, '\\sqrt{x^2}')?.latex).toBe('|x|')
  })

  it('stays quiet when there is nothing to gain', () => {
    // Same expression wearing different spaces is not a simplification, and a hint
    // that fires on every keystroke is a hint nobody reads.
    expect(simplificationHint(engine, '2x+3')).toBeNull()
    expect(simplificationHint(engine, '2x + 3')).toBeNull()
    // A lone variable has nothing to simplify either.
    expect(simplificationHint(engine, 'x')).toBeNull()
  })

  it('says something the learner can act on', () => {
    expect(simplificationHint(engine, '2(x+1)')?.text).toBe('That simplifies to 2x + 2.')
  })
})

describe('loadEngine', () => {
  it('is loaded once and cached', async () => {
    resetEngineForTests()
    const a = await loadEngine()
    const b = await loadEngine()
    expect(a).not.toBeNull()
    // Same instance: the 1.1 MB bundle is not fetched twice.
    expect(a).toBe(b)
  })

  it('registers the global MathLive looks for', async () => {
    // MathLive resolves Compute Engine through
    // globalThis[Symbol.for("io.cortexjs.compute-engine")] and logs an error when it
    // is missing, so importing the library also switches on its MathJSON paste path.
    // Nothing else has to be wired; see docs/math-input.md.
    expect((globalThis as Record<symbol, unknown>)[Symbol.for('io.cortexjs.compute-engine')]).toBeTruthy()
  })

  it('parses what MathLive emits', async () => {
    // The two libraries have to agree on the same string, or the hint would be
    // computed from something the editor never produces.
    for (const latex of ['\\frac{1}{2}', '2^{3}', '\\sqrt{9}', '4\\frac{1}{10}', '\\div']) {
      expect(() => engine.parse(latex)).not.toThrow()
    }
  })
})