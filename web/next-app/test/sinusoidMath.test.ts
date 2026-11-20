import { describe, expect, it } from 'vitest'
import {
  formatPi,
  sinusoid,
  sinusoidExtremes,
  sinusoidFormula,
  sinusoidValue,
} from '../components/visual/math/sinusoid'
import { TAU } from '../components/visual/math/trig'

const near = (a: number, b: number, tolerance = 1e-12) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

const wave = (overrides: Partial<Parameters<typeof sinusoid>[0]> = {}) =>
  sinusoid({ form: 'sin', amplitude: 1, frequency: 1, phase: 0, vertical: 0, ...overrides })

describe('sinusoid parameters', () => {
  it('derives amplitude, period, frequency, midline and extrema from A, B, C and D', () => {
    const w = wave({ amplitude: 2, frequency: 3, phase: Math.PI / 2, vertical: 1 })
    near(w.peak, 2)
    near(w.period, TAU / 3)
    near(w.frequency, 3)
    near(w.midline, 1)
    near(w.max, 3)
    near(w.min, -1)
  })

  it('treats a negative amplitude as a reflection but reports a positive amplitude', () => {
    const w = wave({ amplitude: -2 })
    near(w.peak, 2)
    near(sinusoidValue(w, Math.PI / 2), -2)
    expect(sinusoidFormula(w)).toContain('−2sin')
  })

  it('uses the absolute value of B for the period', () => {
    near(wave({ frequency: 2 }).period, Math.PI)
    near(wave({ frequency: 0.5 }).period, 4 * Math.PI)
    near(wave({ frequency: 4 }).period, Math.PI / 2)
  })

  it('keeps non-finite input as NaN rather than throwing', () => {
    const w = wave({ amplitude: Number.NaN })
    expect(Number.isNaN(w.peak)).toBe(true)
    expect(Number.isNaN(w.period)).toBe(true)
  })

  it('is deterministic', () => {
    const a = wave({ amplitude: 1.7, frequency: 2.3, phase: -0.4, vertical: 0.8 })
    const b = wave({ amplitude: 1.7, frequency: 2.3, phase: -0.4, vertical: 0.8 })
    expect(sinusoidFormula(a)).toBe(sinusoidFormula(b))
    near(sinusoidValue(a, 1.25), sinusoidValue(b, 1.25))
  })
})

describe('sinusoid values', () => {
  it('computes sine and cosine at the quadrant points', () => {
    const sine = wave()
    near(sinusoidValue(sine, 0), 0)
    near(sinusoidValue(sine, Math.PI / 2), 1)
    near(sinusoidValue(sine, Math.PI), 0)
    near(sinusoidValue(sine, (3 * Math.PI) / 2), -1)

    const cosine = wave({ form: 'cos' })
    near(sinusoidValue(cosine, 0), 1)
    near(sinusoidValue(cosine, Math.PI / 2), 0)
    near(sinusoidValue(cosine, Math.PI), -1)
  })

  it('applies the phase shift in the direction of the sign', () => {
    const right = wave({ phase: Math.PI / 2 })
    near(sinusoidValue(right, Math.PI / 2), 0)
    near(sinusoidValue(right, Math.PI), 1)

    const left = wave({ phase: -Math.PI / 2 })
    near(sinusoidValue(left, -Math.PI / 2), 0)
    near(sinusoidValue(left, 0), 1)
  })

  it('applies the vertical shift to every value', () => {
    const w = wave({ amplitude: 2, vertical: 3 })
    near(sinusoidValue(w, Math.PI / 2), 5)
    near(sinusoidValue(w, (3 * Math.PI) / 2), 1)
  })

  it('repeats with the derived period', () => {
    const w = wave({ amplitude: 1.4, frequency: 2.5, phase: 0.3, vertical: -0.7 })
    for (const x of [-2, -0.4, 0, 1.1, 2.7]) {
      near(sinusoidValue(w, x), sinusoidValue(w, x + w.period), 1e-10)
    }
  })
})

describe('sinusoid presentation', () => {
  it('formats multiples of pi as fractions', () => {
    expect(formatPi(0)).toBe('0')
    expect(formatPi(Math.PI)).toBe('π')
    expect(formatPi(-Math.PI / 2)).toBe('−π/2')
    expect(formatPi((3 * Math.PI) / 2)).toBe('3π/2')
    expect(formatPi(TAU)).toBe('2π')
    expect(formatPi(TAU / 3)).toBe('2π/3')
  })

  it('falls back to a decimal multiple for a slider value that is not a small fraction', () => {
    expect(formatPi(0.37)).toBe('0.12π')
  })

  it('writes the equation with the sign, shift and lift', () => {
    const w = wave({ amplitude: 2, frequency: 3, phase: Math.PI / 2, vertical: 1 })
    expect(sinusoidFormula(w)).toBe('y = 2sin(3(x − π/2)) + 1')
    expect(sinusoidFormula(wave({ form: 'cos', phase: -Math.PI / 2, vertical: -1 }))).toBe(
      'y = cos(x + π/2) − 1',
    )
  })

  it('lists the first peak and trough of one period', () => {
    const sine = sinusoidExtremes(wave())
    near(sine.peaks[0], Math.PI / 2)
    near(sine.troughs[0], (3 * Math.PI) / 2)

    const cosine = sinusoidExtremes(wave({ form: 'cos', frequency: 2 }))
    near(cosine.peaks[0], 0)
    near(cosine.peaks[1], Math.PI)
    near(cosine.troughs[0], Math.PI / 2)
  })
})
