import type { Vec2 } from './matrix'

export function distance(a: Vec2, b: Vec2): number {
  return Math.hypot(a[0] - b[0], a[1] - b[1])
}

export function sideLengths(points: Vec2[]): number[] {
  const out: number[] = []
  for (let i = 0; i < points.length; i++) out.push(distance(points[i], points[(i + 1) % points.length]))
  return out
}

/** The largest change in any side length, which is 0 for an isometry. */
export function maxSideDelta(before: Vec2[], after: Vec2[]): number {
  const a = sideLengths(before)
  const b = sideLengths(after)
  if (a.length !== b.length) return Number.POSITIVE_INFINITY
  return a.reduce((max, value, i) => Math.max(max, Math.abs(value - b[i])), 0)
}

export function polygonPath(
  points: Vec2[],
  X: (value: number) => number = (value) => value,
  Y: (value: number) => number = (value) => value,
): string {
  if (points.length === 0) return ''
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${X(p[0])} ${Y(p[1])}`).join(' ') + ' Z'
}
