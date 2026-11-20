// Optimization: a real gradient-descent iteration and a small linear program.
//
// The descent path is the sequence the algorithm actually produces, not an
// interpolated curve. The linear program reports every constraint at the
// candidate point, the exact vertex optimum, and the KKT multipliers of the
// active set. The multipliers are solved from the same normals and objective
// gradient the figure draws, so the dual readout cannot drift from the picture.

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

export type LpShiftWhich = 'none' | 'first' | 'second'

export interface LpShift {
  which: LpShiftWhich
  delta: number
}

export interface Constraint {
  label: string
  /** The left-hand normal of a·(x, y) ≤ b. */
  normal: [number, number]
  bound: number
  /** The slack b − a·(x, y). Nonnegative exactly when the point is feasible. */
  evaluate: (x: number, y: number) => number
}

export interface ConstraintCheck {
  label: string
  value: number
  satisfied: boolean
  active: boolean
}

export interface KktMultiplier {
  label: string
  value: number
  normal: [number, number]
}

export interface KktCheck {
  active: string[]
  multipliers: KktMultiplier[]
  residual: number
  stationary: boolean
  dualFeasible: boolean
  kkt: boolean
  note?: string
}

export interface ShadowPrice {
  label: string
  value: number
}

export interface AllowableRange {
  label: string
  min: number
  max: number
}

export interface LpVertex {
  x: number
  y: number
  value: number
}

export interface LinearProgram {
  constraints: Constraint[]
  vertices: LpVertex[]
  optimum: LpVertex
  candidate: ConstraintCheck[]
  feasible: boolean
  objectiveValue: number
  status: 'optimal' | 'infeasible'
  active: string[]
  kkt: KktCheck
  optimumKkt: KktCheck
  dualValue: number
  gap: number
  shadow: ShadowPrice | null
  allowable: AllowableRange | null
}

const TOLERANCE = 1e-9

function constraint(label: string, normal: [number, number], bound: number): Constraint {
  return {
    label,
    normal,
    bound,
    evaluate: (x, y) => bound - normal[0] * x - normal[1] * y,
  }
}

function constraintsFor(shift: LpShift): Constraint[] {
  const delta = Number.isFinite(shift.delta) ? shift.delta : 0
  const firstBound = 4 + (shift.which === 'first' ? delta : 0)
  const secondBound = 5 + (shift.which === 'second' ? delta : 0)
  return [
    constraint('x ≥ 0', [-1, 0], 0),
    constraint('y ≥ 0', [0, -1], 0),
    constraint(`x + 2y ≤ ${firstBound}`, [1, 2], firstBound),
    constraint(`2x + y ≤ ${secondBound}`, [2, 1], secondBound),
  ]
}

function intersection(first: Constraint, second: Constraint): [number, number] | null {
  const determinant =
    first.normal[0] * second.normal[1] - first.normal[1] * second.normal[0]
  if (Math.abs(determinant) <= 1e-12) return null
  const x =
    (first.bound * second.normal[1] - first.normal[1] * second.bound) / determinant
  const y =
    (first.normal[0] * second.bound - first.bound * second.normal[0]) / determinant
  return [x, y]
}

function solveMultipliers(
  active: Constraint[],
  target: [number, number],
): { multipliers: KktMultiplier[]; residual: number } {
  if (active.length === 0) {
    return { multipliers: [], residual: Math.hypot(target[0], target[1]) }
  }
  if (active.length === 1) {
    const item = active[0]
    const lengthSquared =
      item.normal[0] * item.normal[0] + item.normal[1] * item.normal[1]
    if (lengthSquared <= 1e-15) {
      return { multipliers: [], residual: Math.hypot(target[0], target[1]) }
    }
    const value =
      (target[0] * item.normal[0] + target[1] * item.normal[1]) / lengthSquared
    const residual = Math.hypot(
      target[0] - value * item.normal[0],
      target[1] - value * item.normal[1],
    )
    return { multipliers: [{ label: item.label, value, normal: item.normal }], residual }
  }

  let best: { multipliers: KktMultiplier[]; residual: number } | null = null
  for (let i = 0; i < active.length; i++) {
    for (let j = i + 1; j < active.length; j++) {
      const first = active[i]
      const second = active[j]
      const determinant =
        first.normal[0] * second.normal[1] - first.normal[1] * second.normal[0]
      if (Math.abs(determinant) <= 1e-12) continue
      const firstValue =
        (target[0] * second.normal[1] - target[1] * second.normal[0]) / determinant
      const secondValue =
        (first.normal[0] * target[1] - first.normal[1] * target[0]) / determinant
      const residual = Math.hypot(
        target[0] -
          firstValue * first.normal[0] -
          secondValue * second.normal[0],
        target[1] -
          firstValue * first.normal[1] -
          secondValue * second.normal[1],
      )
      if (!best || residual < best.residual) {
        best = {
          multipliers: [
            { label: first.label, value: firstValue, normal: first.normal },
            { label: second.label, value: secondValue, normal: second.normal },
          ],
          residual,
        }
      }
    }
  }
  return best ?? { multipliers: [], residual: Math.hypot(target[0], target[1]) }
}

function kktCheck(
  x: number,
  y: number,
  feasible: boolean,
  constraints: Constraint[],
): KktCheck {
  const activeConstraints = constraints.filter(
    (item) => Math.abs(item.evaluate(x, y)) <= TOLERANCE,
  )
  const { multipliers, residual } = solveMultipliers(activeConstraints, [3, 2])
  const stationary = residual <= 1e-8
  const dualFeasible = multipliers.every((item) => item.value >= -1e-8)
  const kkt = feasible && stationary && dualFeasible
  let note: string | undefined
  if (!feasible) note = 'candidate is not primal feasible'
  else if (!stationary) note = 'gradient is not a combination of the active normals'
  else if (!dualFeasible) note = 'an active multiplier is negative'
  return {
    active: activeConstraints.map((item) => item.label),
    multipliers,
    residual,
    stationary,
    dualFeasible,
    kkt,
    note,
  }
}

const EMPTY_KKT: KktCheck = {
  active: [],
  multipliers: [],
  residual: Number.NaN,
  stationary: false,
  dualFeasible: false,
  kkt: false,
}

export function linearProgram(
  x: number,
  y: number,
  shift: LpShift = { which: 'none', delta: 0 },
): LinearProgram {
  const constraints = constraintsFor(shift)
  const candidate = constraints.map((item) => {
    const value = item.evaluate(x, y)
    return {
      label: item.label,
      value,
      satisfied: value >= -TOLERANCE,
      active: Math.abs(value) <= TOLERANCE,
    }
  })
  const feasible = candidate.every((item) => item.satisfied)

  const vertices: LpVertex[] = []
  for (let i = 0; i < constraints.length; i++) {
    for (let j = i + 1; j < constraints.length; j++) {
      const point = intersection(constraints[i], constraints[j])
      if (!point) continue
      if (constraints.every((item) => item.evaluate(point[0], point[1]) >= -TOLERANCE)) {
        const duplicate = vertices.some(
          (vertex) =>
            Math.abs(vertex.x - point[0]) <= TOLERANCE &&
            Math.abs(vertex.y - point[1]) <= TOLERANCE,
        )
        if (!duplicate) {
          vertices.push({ x: point[0], y: point[1], value: 3 * point[0] + 2 * point[1] })
        }
      }
    }
  }

  const status: LinearProgram['status'] = vertices.length > 0 ? 'optimal' : 'infeasible'
  const optimum = vertices.length
    ? vertices.reduce((best, vertex) => (vertex.value > best.value ? vertex : best))
    : { x: 0, y: 0, value: 0 }

  const kkt = kktCheck(x, y, feasible, constraints)
  const optimumKkt = kktCheck(optimum.x, optimum.y, status === 'optimal', constraints)
  const dualValue = optimumKkt.multipliers.reduce((total, multiplier) => {
    const item = constraints.find((constraintItem) => constraintItem.label === multiplier.label)
    return total + (item ? multiplier.value * item.bound : 0)
  }, 0)

  let shadow: ShadowPrice | null = null
  let allowable: AllowableRange | null = null
  if (shift.which !== 'none' && status === 'optimal') {
    const firstBound = constraints[2].bound
    const secondBound = constraints[3].bound
    if (shift.which === 'first') {
      shadow = {
        label: 'x + 2y',
        value: optimumKkt.multipliers.find((item) => item.label === constraints[2].label)?.value ?? 0,
      }
      allowable = { label: 'x + 2y', min: secondBound / 2, max: 2 * secondBound }
    } else {
      shadow = {
        label: '2x + y',
        value: optimumKkt.multipliers.find((item) => item.label === constraints[3].label)?.value ?? 0,
      }
      allowable = { label: '2x + y', min: firstBound / 2, max: 2 * firstBound }
    }
  }

  return {
    constraints,
    vertices,
    optimum,
    candidate,
    feasible,
    objectiveValue: 3 * x + 2 * y,
    status,
    active: kkt.active,
    kkt,
    optimumKkt,
    dualValue: status === 'optimal' ? dualValue : 0,
    gap: status === 'optimal' ? optimum.value - dualValue : Number.NaN,
    shadow,
    allowable,
  }
}

/** The contour 3x + 2y = level, as a point and a direction. */
export function objectiveContour(level: number): { point: [number, number]; direction: [number, number] } {
  if (!Number.isFinite(level)) return { point: [0, 0], direction: [-2, 3] }
  return { point: [level / 3, 0], direction: [-2, 3] }
}

export { EMPTY_KKT }

export interface LagrangeModel {
  x: number
  y: number
  lambda: number
  objectiveValue: number
  constraintValue: number
  gradient: [number, number]
  /** Derivative of the objective along x + y = 2. Zero at the constrained minimum. */
  along: number
  optimal: boolean
  optimum: { x: number; y: number; lambda: number; value: number }
}

/** Minimise x² + y² subject to x + y = 2. The candidate always lies on the line. */
export function lagrangeModel(candidateX: number): LagrangeModel {
  const x = Number.isFinite(candidateX) ? candidateX : 1
  const y = 2 - x
  const gradient: [number, number] = [2 * x, 2 * y]
  const lambda = x + y
  const along = (gradient[0] - gradient[1]) / Math.SQRT2
  return {
    x,
    y,
    lambda,
    objectiveValue: x * x + y * y,
    constraintValue: x + y,
    gradient,
    along,
    optimal: Math.abs(along) <= 1e-9,
    optimum: { x: 1, y: 1, lambda: 2, value: 2 },
  }
}

export interface NewtonPoint {
  index: number
  x: number
  f: number
  df: number
  d2f: number
  next: number
}

export interface NewtonRun {
  points: NewtonPoint[]
  status: 'converged' | 'running' | 'stationary' | 'invalid'
  note?: string
}

export function newtonObjective(x: number): number {
  return x ** 4 - 3 * x * x + x
}

export function newtonFirst(x: number): number {
  return 4 * x ** 3 - 6 * x + 1
}

export function newtonSecond(x: number): number {
  return 12 * x * x - 6
}

/** Newton's method for a stationary point of the one-dimensional objective. */
export function newtonMinimise(start: number, steps: number): NewtonRun {
  if (!Number.isFinite(start) || !Number.isFinite(steps) || steps < 1) {
    return { points: [], status: 'invalid', note: 'enter a finite start and at least one step' }
  }
  const points: NewtonPoint[] = []
  let x = start
  let status: NewtonRun['status'] = 'running'
  for (let index = 0; index < Math.min(20, Math.floor(steps)); index++) {
    const f = newtonObjective(x)
    const df = newtonFirst(x)
    const d2f = newtonSecond(x)
    if (![f, df, d2f].every(Number.isFinite)) {
      status = 'invalid'
      break
    }
    if (Math.abs(df) <= 1e-9) {
      points.push({ index, x, f, df, d2f, next: x })
      status = 'converged'
      break
    }
    if (d2f <= 1e-10) {
      points.push({ index, x, f, df, d2f, next: Number.NaN })
      status = 'stationary'
      break
    }
    const next = x - df / d2f
    points.push({ index, x, f, df, d2f, next })
    if (!Number.isFinite(next) || Math.abs(next) > 1e6) {
      status = 'invalid'
      break
    }
    x = next
    if (Math.abs(df) <= 1e-9) status = 'converged'
  }
  return {
    points,
    status,
    note:
      status === 'stationary'
        ? 'the curvature is not positive, so the Newton step is not a descent step'
        : status === 'invalid'
          ? 'the iteration left the finite range'
          : undefined,
  }
}
