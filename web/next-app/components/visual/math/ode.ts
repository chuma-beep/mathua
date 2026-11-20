// First-order slope fields and numerical solution curves.
//
// A field is a real right-hand side f(x, y); the drawn segment at a point has
// slope f(x, y), never a decorative angle. Solution curves are produced by a
// fixed-step integrator and carry their method and step size into the readout,
// so an Euler approximation is never presented as the exact solution.

export type OdeFieldId = 'exponential' | 'linear' | 'logistic' | 'oscillating' | 'sum'

export type OdeMethod = 'rk4' | 'euler'

export interface OdeField {
  id: OdeFieldId
  label: string
  equation: string
}

export const ODE_FIELDS: OdeField[] = [
  { id: 'exponential', label: 'Growth', equation: "y' = y" },
  { id: 'linear', label: 'Relaxation', equation: "y' = x − y" },
  { id: 'logistic', label: 'Logistic', equation: "y' = y(1 − y)" },
  { id: 'oscillating', label: 'Driven decay', equation: "y' = sin(x) − y/2" },
  { id: 'sum', label: 'Coupled growth', equation: "y' = x + y" },
]

export function odeField(id: OdeFieldId): OdeField {
  return ODE_FIELDS.find((field) => field.id === id) ?? ODE_FIELDS[0]
}

/** The actual right-hand side, with no clipping of the mathematical value. */
export function slopeAt(id: OdeFieldId, x: number, y: number): number {
  switch (id) {
    case 'exponential':
      return y
    case 'linear':
      return x - y
    case 'logistic':
      return y * (1 - y)
    case 'oscillating':
      return Math.sin(x) - y / 2
    case 'sum':
      return x + y
  }
}

export interface OdePoint {
  x: number
  y: number
}

export type OdeStatus = 'ok' | 'boundary' | 'diverged' | 'invalid'

export interface OdeTrajectory {
  points: OdePoint[]
  status: OdeStatus
  steps: number
  note?: string
}

export interface OdeIntegrationOptions {
  xMin: number
  xMax: number
  yMin: number
  yMax: number
  step: number
  method?: OdeMethod
  maxSteps?: number
}

const finitePoint = (point: OdePoint) => Number.isFinite(point.x) && Number.isFinite(point.y)

function rk4Step(id: OdeFieldId, x: number, y: number, h: number): OdePoint {
  const k1 = slopeAt(id, x, y)
  const k2 = slopeAt(id, x + h / 2, y + (h * k1) / 2)
  const k3 = slopeAt(id, x + h / 2, y + (h * k2) / 2)
  const k4 = slopeAt(id, x + h, y + h * k3)
  return { x: x + h, y: y + (h / 6) * (k1 + 2 * k2 + 2 * k3 + k4) }
}

function eulerStep(id: OdeFieldId, x: number, y: number, h: number): OdePoint {
  return { x: x + h, y: y + h * slopeAt(id, x, y) }
}

function march(
  id: OdeFieldId,
  start: OdePoint,
  sign: 1 | -1,
  options: OdeIntegrationOptions,
): OdeTrajectory {
  const h = sign * options.step
  const method = options.method ?? 'rk4'
  const maxSteps = options.maxSteps ?? 1200
  const points: OdePoint[] = [start]
  let current = start
  let status: OdeStatus = 'ok'
  let steps = 0

  while (steps < maxSteps) {
    const next = method === 'euler' ? eulerStep(id, current.x, current.y, h) : rk4Step(id, current.x, current.y, h)
    if (!finitePoint(next)) {
      status = 'diverged'
      break
    }
    if (next.y < options.yMin || next.y > options.yMax) {
      status = 'boundary'
      break
    }
    if (next.x < options.xMin || next.x > options.xMax) {
      break
    }
    points.push(next)
    current = next
    steps++
  }

  return {
    points,
    status,
    steps,
    note:
      status === 'boundary'
        ? 'The curve leaves the plotted window.'
        : status === 'diverged'
          ? 'The numerical step stopped at a non-finite value.'
          : undefined,
  }
}

/**
 * Integrate through the initial condition in both directions. The returned
 * points are ordered left to right, so a renderer can draw one polyline.
 */
export function integrateSlopeField(
  id: OdeFieldId,
  x0: number,
  y0: number,
  options: OdeIntegrationOptions,
): OdeTrajectory {
  if (!Number.isFinite(x0) || !Number.isFinite(y0) || !Number.isFinite(options.step) || options.step <= 0) {
    return { points: [], status: 'invalid', steps: 0, note: 'The initial condition or step size is invalid.' }
  }
  if (x0 < options.xMin || x0 > options.xMax || y0 < options.yMin || y0 > options.yMax) {
    return { points: [], status: 'invalid', steps: 0, note: 'The initial condition is outside the plotted window.' }
  }

  const forward = march(id, { x: x0, y: y0 }, 1, options)
  const backward = march(id, { x: x0, y: y0 }, -1, options)
  const backwardPoints = backward.points.slice(1).reverse()
  const points = [...backwardPoints, ...forward.points]
  const status = forward.status === 'invalid' || backward.status === 'invalid'
    ? 'invalid'
    : forward.status === 'diverged' || backward.status === 'diverged'
      ? 'diverged'
      : forward.status === 'boundary' || backward.status === 'boundary'
        ? 'boundary'
        : 'ok'
  return {
    points,
    status,
    steps: forward.steps + backward.steps,
    note: forward.note ?? backward.note,
  }
}
