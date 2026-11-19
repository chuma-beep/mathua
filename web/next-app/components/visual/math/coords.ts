// Coordinate geometry on the plane.
//
// Every number the figure prints comes from here, so the working shown beside
// the drawing is the same arithmetic as the answer. Two shapes of line are
// kept distinct: a vertical line has no slope, and pretending it has a very
// large one would make the intersection formula disagree with the picture.

export type Point = [number, number]

export type Line =
  | { kind: 'vertical'; x: number }
  | { kind: 'slope'; m: number; b: number }

export function distance(a: Point, b: Point): number {
  return Math.hypot(b[0] - a[0], b[1] - a[1])
}

export function midpoint(a: Point, b: Point): Point {
  return [(a[0] + b[0]) / 2, (a[1] + b[1]) / 2]
}

/** The line through two points, or a vertical line when the x values match. */
export function lineThrough(a: Point, b: Point): Line {
  const dx = b[0] - a[0]
  if (Math.abs(dx) < 1e-9) return { kind: 'vertical', x: a[0] }
  const m = (b[1] - a[1]) / dx
  return { kind: 'slope', m, b: a[1] - m * a[0] }
}

/** The line through a point with a given slope. */
export function lineWithSlope(through: Point, m: number): Line {
  return { kind: 'slope', m, b: through[1] - m * through[0] }
}

export type Intersection =
  | { kind: 'point'; x: number; y: number }
  | { kind: 'parallel' }
  | { kind: 'coincident' }

/**
 * Where two lines meet. Parallel and coincident are separate answers on purpose:
 * a pair of identical lines has infinitely many shared points, and reporting a
 * coordinate for them would be a number that is true of one point and false of
 * the rest.
 */
export function intersect(l1: Line, l2: Line): Intersection {
  if (l1.kind === 'vertical' && l2.kind === 'vertical') {
    return Math.abs(l1.x - l2.x) < 1e-9 ? { kind: 'coincident' } : { kind: 'parallel' }
  }
  if (l1.kind === 'vertical') {
    const other = l2 as Extract<Line, { kind: 'slope' }>
    return { kind: 'point', x: l1.x, y: other.m * l1.x + other.b }
  }
  if (l2.kind === 'vertical') {
    const first = l1 as Extract<Line, { kind: 'slope' }>
    return { kind: 'point', x: l2.x, y: first.m * l2.x + first.b }
  }
  if (Math.abs(l1.m - l2.m) < 1e-9) {
    return Math.abs(l1.b - l2.b) < 1e-9 ? { kind: 'coincident' } : { kind: 'parallel' }
  }
  const x = (l2.b - l1.b) / (l1.m - l2.m)
  return { kind: 'point', x, y: l1.m * x + l1.b }
}

/** A point on the line at a given x, for drawing it across the plot. */
export function yAt(line: Line, x: number): number | null {
  return line.kind === 'vertical' ? null : line.m * x + line.b
}

/** Whether a point lies on the line, within a tolerance. */
export function onLine(point: Point, line: Line, tolerance = 1e-9): boolean {
  if (line.kind === 'vertical') return Math.abs(point[0] - line.x) < tolerance
  return Math.abs(point[1] - (line.m * point[0] + line.b)) < tolerance
}
