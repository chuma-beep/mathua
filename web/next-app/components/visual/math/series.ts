// Taylor polynomials of the four functions a visual offers.
//
// A visual names a function from a closed list rather than accepting an
// expression, for the same reason the Riemann sum takes coefficients: the kth
// derivative at the centre is then a closed form, so the polynomial is exact and
// no parser can disagree with the renderer.

export type SeriesFn = 'sin' | 'cos' | 'exp' | 'geometric'

export type SeriesSpec = {
  fn: SeriesFn
  center: number
}

export function fnValue(fn: SeriesFn, x: number): number {
  switch (fn) {
    case 'sin':
      return Math.sin(x)
    case 'cos':
      return Math.cos(x)
    case 'exp':
      return Math.exp(x)
    case 'geometric':
      return Math.abs(1 - x) < 1e-9 ? Number.POSITIVE_INFINITY : 1 / (1 - x)
  }
}

/**
 * The kth derivative at a, as a number. sin and cos shift phase by k quarter
 * turns, exp repeats itself, and 1/(1-x) has derivative k!/(1-a)^(k+1).
 */
export function nthDerivative(fn: SeriesFn, a: number, k: number): number {
  const turn = (k * Math.PI) / 2
  switch (fn) {
    case 'sin':
      return Math.sin(a + turn)
    case 'cos':
      return Math.cos(a + turn)
    case 'exp':
      return Math.exp(a)
    case 'geometric': {
      const base = 1 - a
      if (Math.abs(base) < 1e-9) return Number.POSITIVE_INFINITY
      return factorial(k) / Math.pow(base, k + 1)
    }
  }
}

function factorial(n: number): number {
  let out = 1
  for (let i = 2; i <= n; i++) out *= i
  return out
}

/** Coefficients of the Taylor polynomial about `center`, ascending in (x - center). */
export function taylorCoeffs(fn: SeriesFn, center: number, degree: number): number[] {
  const out: number[] = []
  for (let k = 0; k <= degree; k++) out.push(nthDerivative(fn, center, k) / factorial(k))
  return out
}

export function taylorValue(coeffs: number[], center: number, x: number): number {
  let value = 0
  for (let i = coeffs.length - 1; i >= 0; i--) value = value * (x - center) + coeffs[i]
  return value
}

/** The largest |f(x) - P(x)| over a sampled window, which is the visual's claim. */
export type SumKind = 'geometric' | 'harmonic' | 'p2' | 'alternating'

export const SUM_LABELS: Record<SumKind, string> = {
  geometric: 'Σ (1/2)^k',
  harmonic: 'Σ 1/k',
  p2: 'Σ 1/k²',
  alternating: 'Σ (−1)^(k+1)/k',
}

export const SUM_NAMES: Record<SumKind, string> = {
  geometric: 'Geometric',
  harmonic: 'Harmonic',
  p2: 'p-series, p = 2',
  alternating: 'Alternating',
}

/** The kth term, with k counted from 1. */
export function sumTerm(kind: SumKind, k: number): number {
  if (!Number.isFinite(k) || k < 1) return Number.NaN
  switch (kind) {
    case 'geometric':
      return Math.pow(0.5, k)
    case 'harmonic':
      return 1 / k
    case 'p2':
      return 1 / (k * k)
    case 'alternating':
      return (k % 2 === 1 ? 1 : -1) / k
  }
}

/** The sum of the first n terms, accumulated in order so it matches the bars. */
export function partialSum(kind: SumKind, n: number): number {
  const count = Math.max(0, Math.floor(n))
  let total = 0
  for (let k = 1; k <= count; k++) total += sumTerm(kind, k)
  return total
}

/** The value the partial sums approach, or Infinity when they do not. */
export function seriesLimit(kind: SumKind): number {
  switch (kind) {
    case 'geometric':
      return 1
    case 'harmonic':
      return Number.POSITIVE_INFINITY
    case 'p2':
      return (Math.PI * Math.PI) / 6
    case 'alternating':
      return Math.LN2
  }
}

/** |a_{k+1}/a_k|, which is what the ratio test reads. */
export function ratioTerm(kind: SumKind, k: number): number {
  const next = Math.abs(sumTerm(kind, k + 1))
  const current = Math.abs(sumTerm(kind, k))
  if (current === 0) return Number.NaN
  return next / current
}

export function sumStatus(kind: SumKind): 'converges' | 'diverges' {
  return kind === 'harmonic' ? 'diverges' : 'converges'
}

export function maxError(fn: SeriesFn, center: number, degree: number, a: number, b: number, samples = 200): number {
  // A sample grid can step over the pole of 1/(1-x) and report a finite maximum for
  // a window where the function has no value. The pole is known in closed form, so
  // it is checked directly rather than hoped for.
  if (fn === 'geometric' && a <= 1 && b >= 1) return Number.NaN
  const coeffs = taylorCoeffs(fn, center, degree)
  let worst = 0
  for (let i = 0; i <= samples; i++) {
    const x = a + ((b - a) * i) / samples
    const d = Math.abs(fnValue(fn, x) - taylorValue(coeffs, center, x))
    // Returning 0 for an undefined sample is worse than returning nothing: the
    // window of 1/(1-x) contains its pole, so at least one sample has no value,
    // and a figure that printed "max error 0" there would be claiming a match it
    // never computed.
    if (!Number.isFinite(d)) return Number.NaN
    worst = Math.max(worst, d)
  }
  return worst
}
