// The conic Ax² + Bxy + Cy² = 1, classified from its quadratic part.
//
// The classification is the eigenvalues of [[A, B/2], [B/2, C]], not the sign of
// B² − 4AC alone: a negative definite form has a negative discriminant and no
// points at all, which the discriminant cannot distinguish from an ellipse. The
// readout prints both, because the lesson's own steps are "count the squared
// terms", "compare the coefficients" then "check the signs", and a learner should
// be able to see each of those in the figure.

export type ConicType = 'circle' | 'ellipse' | 'hyperbola' | 'parallel' | 'empty'

export interface ConicShape {
  type: ConicType
  /** How many of A and C are nonzero. Zero means the equation is linear or empty. */
  squaredTerms: number
  /** True when A and C are both nonzero and share a sign. */
  sameSign: boolean
  /** 0 for a circle, between 0 and 1 for an ellipse, above 1 for a hyperbola, NaN when there is no curve. */
  eccentricity: number
}

export function quadraticForm(A: number, B: number, C: number, x: number, y: number): number {
  return A * x * x + B * x * y + C * y * y
}

/** B² − 4AC, the test the lesson states. */
export function discriminant(A: number, B: number, C: number): number {
  return B * B - 4 * A * C
}

/** Eigenvalues of the symmetric matrix [[A, B/2], [B/2, C]], ascending. */
export function eigenvalues(A: number, B: number, C: number): [number, number] {
  const half = B / 2
  const mean = (A + C) / 2
  const spread = Math.sqrt(Math.max(0, ((A - C) / 2) * ((A - C) / 2) + half * half))
  return [mean - spread, mean + spread]
}

export function classifyQuadratic(A: number, B: number, C: number): ConicShape {
  const [low, high] = eigenvalues(A, B, C)
  const eps = 1e-9
  const squaredTerms = (Math.abs(A) > eps ? 1 : 0) + (Math.abs(C) > eps ? 1 : 0)
  const sameSign = Math.abs(A) > eps && Math.abs(C) > eps && A * C > 0

  // A zero discriminant here means the quadratic form is degenerate, and
  // Ax² + Bxy + Cy² = 1 with one zero eigenvalue is a pair of parallel lines, not
  // a parabola: a parabola needs a linear term, and this equation has none. The
  // earlier version reported the textbook rule "zero discriminant means parabola"
  // and drew two parallel lines underneath the word.
  let type: ConicType
  if (Math.abs(low) < eps && Math.abs(high) < eps) type = 'empty'
  else if (high <= eps) type = 'empty'
  else if (Math.abs(low) < eps) type = 'parallel'
  else if (low < 0) type = 'hyperbola'
  else type = Math.abs(low - high) < 1e-6 ? 'circle' : 'ellipse'

  // One formula covers every case that has a curve: the smaller eigenvalue
  // reaches zero exactly when the ellipse is about to open into parallel lines.
  // Only a curve with two directions has an eccentricity. A degenerate form is
  // two lines or no points, and reporting 1 for it would name a parabola again.
  const eccentricity =
    type === 'circle' || type === 'ellipse' || type === 'hyperbola'
      ? Math.sqrt(Math.max(0, 1 - low / high))
      : Number.NaN

  return { type, squaredTerms, sameSign, eccentricity }
}

/** Where the curve meets the x axis, if it does. */
export function xIntercepts(A: number, C: number): number[] {
  if (A <= 1e-9) return []
  const r = 1 / Math.sqrt(A)
  return [-r, r]
}

/** Where the curve meets the y axis, if it does. */
export function yIntercepts(A: number, C: number): number[] {
  if (C <= 1e-9) return []
  const r = 1 / Math.sqrt(C)
  return [-r, r]
}
