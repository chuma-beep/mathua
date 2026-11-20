// Phase lines and planar phase portraits.
//
// A one-dimensional autonomous equation is drawn as a line with arrows that
// follow the sign of f(x). A two-dimensional linear system is drawn from its
// matrix: the vector field, the nullclines and the trajectories all use the
// same entries, and the classification comes from the eigenvalues.

export type PhaseLineModel = 'decay' | 'logistic' | 'bistable' | 'threshold'

export interface PhaseEquilibrium {
  x: number
  stable: boolean
}

export interface PhaseLine {
  model: PhaseLineModel
  label: string
  equation: string
  equilibria: PhaseEquilibrium[]
}

export type PhaseSystemId = 'saddle' | 'stable-node' | 'spiral' | 'center'

export interface PhaseSystem {
  id: PhaseSystemId
  label: string
  matrix: [number, number, number, number]
  trace: number
  determinant: number
  eigenvalues: [number, number] | { real: number; imaginary: number }
  classification: string
  stability: 'stable' | 'unstable' | 'neutral'
}

export interface PhaseTrajectory {
  points: Array<[number, number]>
  status: 'ok' | 'left-window' | 'invalid'
}

const LINE_MODELS: Record<PhaseLineModel, Omit<PhaseLine, 'equilibria'>> = {
  decay: { model: 'decay', label: 'Exponential decay', equation: "x' = −x" },
  logistic: { model: 'logistic', label: 'Logistic growth', equation: "x' = x(1 − x)" },
  bistable: { model: 'bistable', label: 'Bistable switch', equation: "x' = x(1 − x)(x − 0.5)" },
  threshold: { model: 'threshold', label: 'Threshold', equation: "x' = x² − 1" },
}

export function phaseLineSlope(model: PhaseLineModel, x: number): number {
  switch (model) {
    case 'decay':
      return -x
    case 'logistic':
      return x * (1 - x)
    case 'bistable':
      return x * (1 - x) * (x - 0.5)
    case 'threshold':
      return x * x - 1
  }
}

export function phaseLine(model: PhaseLineModel): PhaseLine {
  const equilibria: PhaseEquilibrium[] = []
  const add = (x: number, stable: boolean) => equilibria.push({ x, stable })
  switch (model) {
    case 'decay':
      add(0, true)
      break
    case 'logistic':
      add(0, false)
      add(1, true)
      break
    case 'bistable':
      add(0, true)
      add(0.5, false)
      add(1, true)
      break
    case 'threshold':
      add(-1, true)
      add(1, false)
      break
  }
  return { ...LINE_MODELS[model], equilibria }
}

const SYSTEMS: Record<PhaseSystemId, PhaseSystem> = {
  saddle: {
    id: 'saddle',
    label: 'Saddle',
    matrix: [2, 0, 0, -3],
    trace: -1,
    determinant: -6,
    eigenvalues: [2, -3],
    classification: 'saddle',
    stability: 'unstable',
  },
  'stable-node': {
    id: 'stable-node',
    label: 'Stable node',
    matrix: [-1, 0, 0, -2],
    trace: -3,
    determinant: 2,
    eigenvalues: [-1, -2],
    classification: 'stable node',
    stability: 'stable',
  },
  spiral: {
    id: 'spiral',
    label: 'Stable spiral',
    matrix: [-0.5, 1, -1, -0.5],
    trace: -1,
    determinant: 1.25,
    eigenvalues: { real: -0.5, imaginary: 1 },
    classification: 'stable spiral',
    stability: 'stable',
  },
  center: {
    id: 'center',
    label: 'Centre',
    matrix: [0, 1, -1, 0],
    trace: 0,
    determinant: 1,
    eigenvalues: { real: 0, imaginary: 1 },
    classification: 'centre',
    stability: 'neutral',
  },
}

export function phaseSystem(id: PhaseSystemId): PhaseSystem {
  return SYSTEMS[id]
}

export function phaseVector(system: PhaseSystem, x: number, y: number): [number, number] {
  const [a, b, c, d] = system.matrix
  return [a * x + b * y, c * x + d * y]
}

/** Nullclines as point/direction pairs: x' = 0 and y' = 0. */
export function phaseNullclines(system: PhaseSystem): Array<{
  label: string
  point: [number, number]
  direction: [number, number]
}> {
  const [a, b, c, d] = system.matrix
  const lines: Array<{ label: string; point: [number, number]; direction: [number, number] }> = []
  if (Math.abs(b) > 1e-12) lines.push({ label: "x' = 0", point: [0, 0], direction: [b, -a] })
  else if (Math.abs(a) > 1e-12) lines.push({ label: "x' = 0", point: [0, 0], direction: [0, 1] })
  if (Math.abs(d) > 1e-12) lines.push({ label: "y' = 0", point: [0, 0], direction: [d, -c] })
  else if (Math.abs(c) > 1e-12) lines.push({ label: "y' = 0", point: [0, 0], direction: [1, 0] })
  return lines
}

/** RK4 integration of the linear system from one seed, in one time direction. */
export function phaseTrajectory(
  system: PhaseSystem,
  start: [number, number],
  duration = 4,
  steps = 150,
): PhaseTrajectory {
  if (!start.every(Number.isFinite) || !Number.isFinite(duration) || steps < 1) {
    return { points: [], status: 'invalid' }
  }
  const dt = duration / Math.floor(steps)
  let [x, y] = start
  const points: Array<[number, number]> = [[x, y]]
  for (let index = 0; index < steps; index++) {
    const k1 = phaseVector(system, x, y)
    const k2 = phaseVector(system, x + (dt * k1[0]) / 2, y + (dt * k1[1]) / 2)
    const k3 = phaseVector(system, x + (dt * k2[0]) / 2, y + (dt * k2[1]) / 2)
    const k4 = phaseVector(system, x + dt * k3[0], y + dt * k3[1])
    x += (dt / 6) * (k1[0] + 2 * k2[0] + 2 * k3[0] + k4[0])
    y += (dt / 6) * (k1[1] + 2 * k2[1] + 2 * k3[1] + k4[1])
    if (!Number.isFinite(x) || !Number.isFinite(y)) return { points, status: 'invalid' }
    points.push([x, y])
    if (Math.abs(x) > 8 || Math.abs(y) > 8) return { points, status: 'left-window' }
  }
  return { points, status: 'ok' }
}
