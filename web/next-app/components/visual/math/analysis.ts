// Real analysis and metric spaces.
//
// A sequence readout reports the actual tail distance after the chosen index.
// The ball readout computes the metric it names, so the L1, L2 and L∞ pictures
// can differ without the labels pretending they are the same geometry.

export type SequenceId = 'reciprocal' | 'ratio' | 'alternating' | 'shifted'
export type MetricId = 'euclidean' | 'manhattan' | 'chebyshev'

export const SEQUENCES: Record<SequenceId, { label: string; formula: string; limit: number }> = {
  reciprocal: { label: '1/n', formula: '1/n', limit: 0 },
  ratio: { label: '(n+1)/n', formula: '(n+1)/n', limit: 1 },
  alternating: { label: '(−1)ⁿ/n', formula: '(−1)ⁿ/n', limit: 0 },
  shifted: { label: '2 + 1/n', formula: '2 + 1/n', limit: 2 },
}

export function sequenceTerm(id: SequenceId, n: number): number {
  if (!Number.isInteger(n) || n < 1) return Number.NaN
  if (id === 'reciprocal') return 1 / n
  if (id === 'ratio') return (n + 1) / n
  if (id === 'alternating') return (n % 2 === 0 ? 1 : -1) / n
  return 2 + 1 / n
}

export interface SequenceView {
  id: SequenceId
  terms: number[]
  limit: number
  epsilon: number
  threshold: number | null
  tailDistance: number | null
  within: boolean
}

export function sequenceView(id: SequenceId, epsilon: number, count = 30): SequenceView {
  const spec = SEQUENCES[id]
  const terms = Array.from({ length: Math.max(1, Math.min(60, Math.floor(count))) }, (_v, index) =>
    sequenceTerm(id, index + 1),
  )
  let threshold: number | null = null
  for (let index = 0; index < terms.length; index++) {
    if (terms.slice(index).every((term) => Math.abs(term - spec.limit) < epsilon)) {
      threshold = index + 1
      break
    }
  }
  const tailDistance = terms.length ? Math.abs(terms[terms.length - 1] - spec.limit) : null
  return {
    id,
    terms,
    limit: spec.limit,
    epsilon,
    threshold,
    tailDistance,
    within: tailDistance !== null && tailDistance < epsilon,
  }
}

export function partialSums(id: SequenceId, count: number): number[] {
  let total = 0
  return Array.from({ length: Math.max(1, Math.min(60, Math.floor(count))) }, (_v, index) => {
    total += sequenceTerm(id, index + 1)
    return total
  })
}

export interface MetricBall {
  metric: MetricId
  radius: number
  center: [number, number]
  pointsInside: number
  points: Array<{ x: number; y: number; distance: number; inside: boolean }>
}

export function metricDistance(metric: MetricId, x: number, y: number): number {
  if (metric === 'manhattan') return Math.abs(x) + Math.abs(y)
  if (metric === 'chebyshev') return Math.max(Math.abs(x), Math.abs(y))
  return Math.hypot(x, y)
}

export function metricBall(metric: MetricId, radius: number): MetricBall {
  const points: MetricBall['points'] = []
  for (let x = -2; x <= 2; x++) {
    for (let y = -2; y <= 2; y++) {
      const distance = metricDistance(metric, x, y)
      points.push({ x, y, distance, inside: distance <= radius + 1e-12 })
    }
  }
  return {
    metric,
    radius,
    center: [0, 0],
    pointsInside: points.filter((point) => point.inside).length,
    points,
  }
}
