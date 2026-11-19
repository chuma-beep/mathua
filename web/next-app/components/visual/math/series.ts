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
