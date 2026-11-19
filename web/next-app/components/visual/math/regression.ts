// Ordinary least squares for a scatter a learner can drag.
//
// Every number the figure shows comes from here, so the line drawn is the line
// whose coefficients are printed. A degenerate cloud (all x equal) has no unique
// least-squares line, and the slope is reported as 0 rather than as an infinity
// the drawing could not place.

export type Point = [number, number]

export function means(points: Point[]): Point {
  if (points.length === 0) return [0, 0]
  let sx = 0
  let sy = 0
  for (const [x, y] of points) {
    sx += x
    sy += y
  }
  return [sx / points.length, sy / points.length]
}

export function slope(points: Point[]): number {
  if (points.length < 2) return 0
  const [mx, my] = means(points)
  let sxy = 0
  let sxx = 0
  for (const [x, y] of points) {
    sxy += (x - mx) * (y - my)
    sxx += (x - mx) * (x - mx)
  }
  return sxx === 0 ? 0 : sxy / sxx
}

/**
 * Whether x varies at all. When every x is the same, the normal equations are
 * singular: no line is better than any other, so `slope` returns 0 as a stand-in
 * and the figure must say the line is undefined rather than draw that stand-in.
 */
export function hasXSpread(points: Point[]): boolean {
  if (points.length < 2) return false
  const [mx] = means(points)
  return points.reduce((sum, [x]) => sum + (x - mx) * (x - mx), 0) > 1e-12
}

export function intercept(points: Point[]): number {
  const [mx, my] = means(points)
  return my - slope(points) * mx
}

export function predict(points: Point[], x: number): number {
  return intercept(points) + slope(points) * x
}

/** Pearson's r, 0 when either variable has no spread. */
export function correlation(points: Point[]): number {
  if (points.length < 2) return 0
  const [mx, my] = means(points)
  let sxy = 0
  let sxx = 0
  let syy = 0
  for (const [x, y] of points) {
    sxy += (x - mx) * (y - my)
    sxx += (x - mx) * (x - mx)
    syy += (y - my) * (y - my)
  }
  if (sxx === 0 || syy === 0) return 0
  return sxy / Math.sqrt(sxx * syy)
}

/** The sum of squared residuals, which the least-squares line minimizes. */
export function sse(points: Point[]): number {
  if (points.length === 0) return 0
  const b = intercept(points)
  const m = slope(points)
  return points.reduce((sum, [x, y]) => {
    const d = y - (b + m * x)
    return sum + d * d
  }, 0)
}
