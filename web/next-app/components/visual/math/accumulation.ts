// An accumulation function and the theorem that differentiates it.
//
// A(x) is the area under f from a fixed left end to x. The figure draws A and
// f on the same x axis, and the readout compares the slope of A with f(x). The
// area is computed from a closed-form antiderivative, not from a sampled sum,
// so the equality is a fact about the functions rather than an artefact of a
// fine grid.

export type AccumulationFn = 'square' | 'sine' | 'decay'

export const ACC_LABELS: Record<AccumulationFn, string> = {
  square: 'f(x) = x²',
  sine: 'f(x) = sin x',
  decay: 'f(x) = e^−x',
}

/** The integrand. */
export function integrand(fn: AccumulationFn, x: number): number {
  switch (fn) {
    case 'square':
      return x * x
    case 'sine':
      return Math.sin(x)
    case 'decay':
      return Math.exp(-x)
  }
}

/** A primitive F with F′ = f, so the area is F(x) − F(a). */
export function primitive(fn: AccumulationFn, x: number): number {
  switch (fn) {
    case 'square':
      return (x * x * x) / 3
    case 'sine':
      return -Math.cos(x)
    case 'decay':
      return -Math.exp(-x)
  }
}

export interface AccumulationModel {
  fn: AccumulationFn
  a: number
  x: number
  /** ∫ₐˣ f, signed. */
  area: number
  /** f(x), which the theorem says is A′(x). */
  value: number
  /** The same derivative measured from A at x ± h, to show agreement. */
  measured: number
  /** |measured − value|. */
  error: number
  status: 'ok' | 'invalid'
}

export function accumulationModel(fn: AccumulationFn, a: number, x: number, h = 1e-4): AccumulationModel {
  if (![a, x].every(Number.isFinite) || !Number.isFinite(h) || h <= 0) {
    return { fn, a, x, area: Number.NaN, value: Number.NaN, measured: Number.NaN, error: Number.NaN, status: 'invalid' }
  }
  const area = primitive(fn, x) - primitive(fn, a)
  const value = integrand(fn, x)
  const measured = (primitive(fn, x + h) - primitive(fn, x - h)) / (2 * h)
  return { fn, a, x, area, value, measured, error: Math.abs(measured - value), status: 'ok' }
}

/** The area function A(x) on a grid, for drawing its own curve. */
export function areaCurve(fn: AccumulationFn, a: number, lo: number, hi: number, count = 120): Array<[number, number]> {
  const out: Array<[number, number]> = []
  const steps = Math.max(2, Math.floor(count))
  for (let index = 0; index <= steps; index++) {
    const x = lo + ((hi - lo) * index) / steps
    out.push([x, primitive(fn, x) - primitive(fn, a)])
  }
  return out
}
