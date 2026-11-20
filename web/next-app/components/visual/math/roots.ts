// Numerical root finding, shown as the actual sequence of brackets and tangent
// intercepts. The functions are deliberately low-degree so every iterate can be
// checked by hand; approximation and exact values stay separate in the readout.

export type RootFunctionId = 'quadratic' | 'cubic' | 'cosine' | 'exp'
export type RootMethod = 'bisection' | 'newton'

export interface RootFunction {
  id: RootFunctionId
  label: string
  expression: string
  f: (x: number) => number
  df: (x: number) => number
  /** Human-readable exact value. */
  exact: string
  /** Numeric exact value for error readouts. */
  exactValue: number
}

export const ROOT_FUNCTIONS: RootFunction[] = [
  {
    id: 'quadratic',
    label: 'x² − 2',
    expression: 'x² − 2',
    f: (x) => x * x - 2,
    df: (x) => 2 * x,
    exact: '√2 = 1.41421356…',
    exactValue: Math.SQRT2,
  },
  {
    id: 'cubic',
    label: 'x³ − x − 1',
    expression: 'x³ − x − 1',
    f: (x) => x * x * x - x - 1,
    df: (x) => 3 * x * x - 1,
    exact: '1.32471795…',
    exactValue: 1.324717957244746,
  },
  {
    id: 'cosine',
    label: 'cos x − x',
    expression: 'cos x − x',
    f: (x) => Math.cos(x) - x,
    df: (x) => -Math.sin(x) - 1,
    exact: '0.73908513…',
    exactValue: 0.7390851332151607,
  },
  {
    id: 'exp',
    label: 'eˣ − 3',
    expression: 'eˣ − 3',
    f: (x) => Math.exp(x) - 3,
    df: (x) => Math.exp(x),
    exact: 'ln 3 = 1.09861229…',
    exactValue: Math.log(3),
  },
]

export function rootFunction(id: RootFunctionId): RootFunction {
  return ROOT_FUNCTIONS.find((item) => item.id === id) ?? ROOT_FUNCTIONS[0]
}

export interface BisectionStep {
  index: number
  a: number
  b: number
  midpoint: number
  fMidpoint: number
  width: number
  keep: 'left' | 'right' | 'stop'
}

export interface BisectionResult {
  steps: BisectionStep[]
  root: number | null
  status: 'ok' | 'invalid' | 'no-sign-change'
  note?: string
}

export function bisection(
  f: (x: number) => number,
  a: number,
  b: number,
  iterations = 6,
  tolerance = 1e-9,
): BisectionResult {
  if (!Number.isFinite(a) || !Number.isFinite(b) || !Number.isFinite(tolerance) || tolerance <= 0) {
    return { steps: [], root: null, status: 'invalid', note: 'The bracket or tolerance is invalid.' }
  }
  if (a >= b) {
    return { steps: [], root: null, status: 'invalid', note: 'The left endpoint must be smaller than the right endpoint.' }
  }
  let left = a
  let right = b
  let fLeft = f(left)
  let fRight = f(right)
  if (!Number.isFinite(fLeft) || !Number.isFinite(fRight)) {
    return { steps: [], root: null, status: 'invalid', note: 'The function is not finite at an endpoint.' }
  }
  if (fLeft === 0) return { steps: [], root: left, status: 'ok', note: 'The left endpoint is already a root.' }
  if (fRight === 0) return { steps: [], root: right, status: 'ok', note: 'The right endpoint is already a root.' }
  if (fLeft * fRight > 0) {
    return { steps: [], root: null, status: 'no-sign-change', note: 'The endpoints have the same sign, so bisection cannot start.' }
  }

  const steps: BisectionStep[] = []
  const count = Math.max(1, Math.min(40, Math.floor(iterations)))
  for (let index = 0; index < count; index++) {
    const midpoint = (left + right) / 2
    const fMidpoint = f(midpoint)
    const width = right - left
    if (!Number.isFinite(fMidpoint)) {
      return { steps, root: null, status: 'invalid', note: 'The midpoint value is not finite.' }
    }
    const atRoot = Math.abs(fMidpoint) <= tolerance
    const keep = atRoot ? 'stop' : fLeft * fMidpoint < 0 ? 'left' : 'right'
    steps.push({ index, a: left, b: right, midpoint, fMidpoint, width, keep })
    if (atRoot || width / 2 <= tolerance) {
      return { steps, root: midpoint, status: 'ok' }
    }
    if (keep === 'left') {
      right = midpoint
      fRight = fMidpoint
    } else {
      left = midpoint
      fLeft = fMidpoint
    }
    if (index === count - 1) {
      return { steps, root: (left + right) / 2, status: 'ok' }
    }
  }
  return { steps, root: (left + right) / 2, status: 'ok' }
}

export interface NewtonStep {
  index: number
  x: number
  fx: number
  slope: number
  next: number
}

export interface NewtonResult {
  steps: NewtonStep[]
  root: number | null
  status: 'ok' | 'invalid' | 'stationary'
  note?: string
}

export function newton(
  f: (x: number) => number,
  df: (x: number) => number,
  start: number,
  iterations = 5,
): NewtonResult {
  if (!Number.isFinite(start)) {
    return { steps: [], root: null, status: 'invalid', note: 'The starting guess is not finite.' }
  }
  let x = start
  const steps: NewtonStep[] = []
  const count = Math.max(1, Math.min(20, Math.floor(iterations)))
  for (let index = 0; index < count; index++) {
    const fx = f(x)
    const slope = df(x)
    if (!Number.isFinite(fx) || !Number.isFinite(slope)) {
      return { steps, root: null, status: 'invalid', note: 'The function or derivative is not finite.' }
    }
    if (Math.abs(slope) <= 1e-14) {
      return { steps, root: null, status: 'stationary', note: 'The tangent is horizontal, so the next intercept is undefined.' }
    }
    const next = x - fx / slope
    steps.push({ index, x, fx, slope, next })
    if (!Number.isFinite(next)) {
      return { steps, root: null, status: 'invalid', note: 'The next iterate is not finite.' }
    }
    if (Math.abs(next - x) <= 1e-12) {
      return { steps, root: next, status: 'ok', note: 'The iterates stopped changing at display precision.' }
    }
    x = next
  }
  return { steps, root: x, status: 'ok' }
}

/** One Newton tangent line at (x, f(x)), for the SVG construction. */
export function newtonTangent(f: (x: number) => number, df: (x: number) => number, x: number) {
  const y = f(x)
  const slope = df(x)
  return { x, y, slope, intercept: Number.isFinite(slope) && slope !== 0 ? x - y / slope : null }
}
