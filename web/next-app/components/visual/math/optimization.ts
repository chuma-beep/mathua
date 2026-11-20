// Optimization: a real gradient-descent iteration and a small linear program.
//
// The descent path is the sequence the algorithm actually produces, not an
// interpolated curve. The linear program reports every constraint at the
// candidate point and the exact vertex optimum of the fixed feasible polygon.

export interface DescentPoint {
  index: number
  x: number
  y: number
  value: number
  gradX: number
  gradY: number
}

export interface DescentRun {
  points: DescentPoint[]
  status: 'converged' | 'running' | 'diverged'
  note?: string
}

export function objective(x: number, y: number): number {
  return x * x + 2 * y * y
}

export function gradient(x: number, y: number): [number, number] {
  return [2 * x, 4 * y]
}

export function gradientDescent(
  startX: number,
  startY: number,
  rate: number,
  steps: number,
): DescentRun {
  if (![startX, startY, rate, steps].every(Number.isFinite) || rate <= 0 || steps < 1) {
    return { points: [], status: 'diverged', note: 'enter a positive rate and at least one step' }
  }
  const points: DescentPoint[] = []
  let x = startX
  let y = startY
  let status: DescentRun['status'] = 'running'
  for (let index = 0; index <= Math.min(30, Math.floor(steps)); index++) {
    if (!Number.isFinite(x) || !Number.isFinite(y) || Math.abs(x) > 1e6 || Math.abs(y) > 1e6) {
      status = 'diverged'
      break
    }
    const [gx, gy] = gradient(x, y)
    points.push({ index, x, y, value: objective(x, y), gradX: gx, gradY: gy })
    if (Math.hypot(gx, gy) < 1e-7) {
      status = 'converged'
      break
    }
    x -= rate * gx
    y -= rate * gy
  }
  return {
    points,
    status,
    note: status === 'diverged' ? 'the iteration left the finite range' : undefined,
  }
}

export interface Constraint {
  label: string
  evaluate: (x: number, y: number) => number
}

export interface ConstraintCheck {
  label: string
  value: number
  satisfied: boolean
}

export interface LinearProgram {
  constraints: Constraint[]
  vertices: Array<{ x: number; y: number; value: number }>
  optimum: { x: number; y: number; value: number }
  candidate: ConstraintCheck[]
  feasible: boolean
  objectiveValue: number
}

const LP_CONSTRAINTS: Constraint[] = [
  { label: 'x ≥ 0', evaluate: (x) => x },
  { label: 'y ≥ 0', evaluate: (_x, y) => y },
  { label: 'x + 2y ≤ 4', evaluate: (x, y) => 4 - x - 2 * y },
  { label: '2x + y ≤ 5', evaluate: (x, y) => 5 - 2 * x - y },
]

const VERTICES = [
  { x: 0, y: 0 },
  { x: 2.5, y: 0 },
  { x: 2, y: 1 },
  { x: 0, y: 2 },
]

export function linearProgram(x: number, y: number): LinearProgram {
  const candidate = LP_CONSTRAINTS.map((constraint) => {
    const value = constraint.evaluate(x, y)
    return { label: constraint.label, value, satisfied: value >= -1e-12 }
  })
  const vertices = VERTICES.map((vertex) => ({ ...vertex, value: 3 * vertex.x + 2 * vertex.y }))
  const optimum = vertices.reduce((best, vertex) => (vertex.value > best.value ? vertex : best))
  return {
    constraints: LP_CONSTRAINTS,
    vertices,
    optimum,
    candidate,
    feasible: candidate.every((item) => item.satisfied),
    objectiveValue: 3 * x + 2 * y,
  }
}
