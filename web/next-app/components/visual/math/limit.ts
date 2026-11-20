// The three limit shapes a visual offers. Each is a closed form rather than a
// parsed expression, so a one-sided value can be computed exactly or reported as
// unbounded, and the figure never guesses at what a graph "looks like".

export type LimitKind = 'removable' | 'jump' | 'vertical_asymptote'

/** f(x) at a point away from a, and NaN where the function is genuinely undefined at a. */
export function limitValue(kind: LimitKind, a: number, x: number): number {
  switch (kind) {
    case 'removable': {
      // (x^2 - a^2)/(x - a) is defined for every x except a, where it has a hole.
      if (Math.abs(x - a) < 1e-12) return Number.NaN
      return x + a
    }
    case 'jump':
      return x < a ? 1 : 3
    case 'vertical_asymptote':
      if (Math.abs(x - a) < 1e-9) return Number.NaN
      return 1 / (x - a)
  }
}

/**
 * The value the function has at a, which is not the same question as its limit.
 * A hole and a vertical asymptote have no value there; the jump is defined, by the
 * piece the point falls in, and its value is 3 rather than undefined. This delegates
 * to `limitValue` so the drawn curve and the printed value cannot disagree.
 */
export function valueAtPoint(kind: LimitKind, a: number): number {
  return limitValue(kind, a, a)
}

/**
 * The one-sided limits. A vertical asymptote is unbounded, so it reports
 * ±Infinity rather than pretending a number exists.
 */
export function oneSidedLimits(kind: LimitKind, a: number): { left: number; right: number } {
  switch (kind) {
    case 'removable':
      return { left: 2 * a, right: 2 * a }
    case 'jump':
      return { left: 1, right: 3 }
    case 'vertical_asymptote':
      return { left: Number.NEGATIVE_INFINITY, right: Number.POSITIVE_INFINITY }
  }
}

/** A two-sided limit exists when both sides are finite and agree. */
export function limitExists(kind: LimitKind, a: number): boolean {
  const { left, right } = oneSidedLimits(kind, a)
  return Number.isFinite(left) && Number.isFinite(right) && Math.abs(left - right) < 1e-12
}

/** Continuous at a means the limit exists and equals the value there. */
export function continuousAt(kind: LimitKind, a: number): boolean {
  if (!limitExists(kind, a)) return false
  const value = valueAtPoint(kind, a)
  const { left } = oneSidedLimits(kind, a)
  return Number.isFinite(value) && Math.abs(value - left) < 1e-12
}
