// The four fields the contour figure offers.
//
// A closed list rather than an expression for the same reason the polynomial
// figures take coefficients: the gradient is then an exact formula, so the arrow
// drawn beside a level curve is the real gradient and not a finite difference
// that could point slightly wrong.

export type Field = 'bowl' | 'saddle' | 'product' | 'ripple'

export const FIELD_LABELS: Record<Field, string> = {
  bowl: 'x² + y²',
  saddle: 'x² − y²',
  product: 'x · y',
  ripple: 'sin x + cos y',
}

export function fieldValue(field: Field, x: number, y: number): number {
  switch (field) {
    case 'bowl':
      return x * x + y * y
    case 'saddle':
      return x * x - y * y
    case 'product':
      return x * y
    case 'ripple':
      return Math.sin(x) + Math.cos(y)
  }
}

/** The exact gradient, which is perpendicular to the level curve through the point. */
export function fieldGradient(field: Field, x: number, y: number): [number, number] {
  switch (field) {
    case 'bowl':
      return [2 * x, 2 * y]
    case 'saddle':
      return [2 * x, -2 * y]
    case 'product':
      return [y, x]
    case 'ripple':
      return [Math.cos(x), -Math.sin(y)]
  }
}

export function gradientLength(field: Field, x: number, y: number): number {
  const [gx, gy] = fieldGradient(field, x, y)
  return Math.hypot(gx, gy)
}

/**
 * The levels to draw, spread between the field's range over the window.
 *
 * The lowest level is deliberately one step above the minimum: for a bowl the
 * minimum is a single point, and a contour through one point draws nothing.
 */
export function fieldLevels(field: Field, extent: number, count: number): number[] {
  const samples = 24
  let lo = Number.POSITIVE_INFINITY
  let hi = Number.NEGATIVE_INFINITY
  for (let i = 0; i <= samples; i++) {
    for (let j = 0; j <= samples; j++) {
      const x = -extent + (2 * extent * i) / samples
      const y = -extent + (2 * extent * j) / samples
      const v = fieldValue(field, x, y)
      if (!Number.isFinite(v)) continue
      lo = Math.min(lo, v)
      hi = Math.max(hi, v)
    }
  }
  if (!Number.isFinite(lo) || !Number.isFinite(hi) || hi <= lo) return []
  const step = (hi - lo) / (count + 1)
  return Array.from({ length: count }, (_, i) => Number((lo + step * (i + 1)).toFixed(6)))
}
