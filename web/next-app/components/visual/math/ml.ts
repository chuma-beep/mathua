// A tiny deterministic regression example: four points, squared error and the
// exact gradient of that error. No external ML framework and no randomness.

export interface Point {
  x: number
  y: number
}

export const REGRESSION_DATA: Point[] = [
  { x: 1, y: 2 },
  { x: 2, y: 3 },
  { x: 3, y: 5 },
  { x: 4, y: 6 },
]

export function predict(slope: number, intercept: number, x: number): number {
  return slope * x + intercept
}

export function mse(slope: number, intercept: number): number {
  const total = REGRESSION_DATA.reduce((sum, point) => {
    const error = predict(slope, intercept, point.x) - point.y
    return sum + error * error
  }, 0)
  return total / REGRESSION_DATA.length
}

export function mseGradient(slope: number, intercept: number): [number, number] {
  let ds = 0
  let db = 0
  for (const point of REGRESSION_DATA) {
    const error = predict(slope, intercept, point.x) - point.y
    ds += (2 * error * point.x) / REGRESSION_DATA.length
    db += (2 * error) / REGRESSION_DATA.length
  }
  return [ds, db]
}

export interface FitStep {
  index: number
  slope: number
  intercept: number
  loss: number
  gradSlope: number
  gradIntercept: number
}

export interface FitRun {
  steps: FitStep[]
  status: 'converged' | 'running' | 'diverged'
}

export function trainRegression(rate: number, count: number): FitRun {
  if (!Number.isFinite(rate) || rate <= 0 || count < 1) {
    return { steps: [], status: 'diverged' }
  }
  const steps: FitStep[] = []
  let slope = 0
  let intercept = 0
  for (let index = 0; index <= Math.min(60, Math.floor(count)); index++) {
    const [gs, gb] = mseGradient(slope, intercept)
    steps.push({ index, slope, intercept, loss: mse(slope, intercept), gradSlope: gs, gradIntercept: gb })
    if (Math.hypot(gs, gb) < 1e-8) return { steps, status: 'converged' }
    slope -= rate * gs
    intercept -= rate * gb
    if (!Number.isFinite(slope) || !Number.isFinite(intercept) || Math.abs(slope) > 1e6) {
      return { steps, status: 'diverged' }
    }
  }
  return { steps, status: 'running' }
}

export function sigmoid(z: number): number {
  return 1 / (1 + Math.exp(-z))
}

export function logisticProbability(weight: number, bias: number, x: number): number {
  return sigmoid(weight * x + bias)
}
