// The general conic Ax² + Bxy + Cy² + Dx + Ey + F = 0.
//
// Classification is a separate, framework-free calculation. It diagonalises the
// quadratic form, translates to the centre where one exists, and reports the
// geometry a renderer may draw. The renderer never inspects A..F and never
// decides a kind from the discriminant: a positive discriminant can be a
// hyperbola or a pair of intersecting lines, and a zero discriminant can be a
// parabola, two parallel lines, one line or no points. Only the canonical form
// separates those cases.

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

export type ConicKind =
  | 'circle'
  | 'ellipse'
  | 'hyperbola'
  | 'parabola'
  | 'point'
  | 'intersecting'
  | 'parallel'
  | 'line'
  | 'empty'
  | 'all-plane'
  | 'invalid'

export type Vec2 = [number, number]

export interface ConicLine {
  point: Vec2
  direction: Vec2
}

export type InterceptStatus = 'two' | 'one' | 'double' | 'none' | 'coincident'

export interface Intercepts {
  status: InterceptStatus
  values: number[]
}

export type Certainty = 'exact' | 'near-degenerate' | 'invalid'

/**
 * The classified conic. Fields that do not belong to a kind are absent, so a
 * renderer that switches on `kind` cannot read a stale axis from another case.
 */
export interface ConicGeometry {
  kind: ConicKind
  center?: Vec2
  vertex?: Vec2
  /** Ellipse only. */
  major?: Vec2
  /** Ellipse only. */
  minor?: Vec2
  /** Hyperbola only. */
  transverse?: Vec2
  /** Hyperbola only. */
  conjugate?: Vec2
  /** Parabola only: the direction that is squared. */
  squared?: Vec2
  /** Parabola only: a unit vector from the vertex toward the focus. */
  opening?: Vec2
  /** Semi-axis along `major` or `transverse`. */
  a?: number
  /** Semi-axis along `minor` or `conjugate`. */
  b?: number
  /** Parabola: nonnegative vertex-to-focus distance. */
  p?: number
  /** Angle of the named principal axis, in radians. */
  theta?: number
  /**
   * Foci in the same coordinates as the curve. Circles have one entry because
   * both foci coincide at the centre; a parabola has one focus.
   */
  foci?: Vec2[]
  /** Parabola only: the directrix line, with a direction parallel to it. */
  directrix?: ConicLine
  /** Ellipse and circle: 2a, the constant focal sum. Hyperbola: 2a, the constant difference. */
  focalConstant?: number
  /** Null when the conic has no eccentricity. */
  eccentricity: number | null
  /** Explicit lines for the degenerate kinds. */
  lines: ConicLine[]
  intercepts: { x: Intercepts; y: Intercepts }
  certainty: Certainty
  note?: string
}

/** Exact rank and zero tests, relative to the local coefficient scale. */
export const EPS = 1e-10
/** A value inside this band is kept but reported as near-degenerate. */
export const NEAR = 1e-6

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

/**
 * A deterministic unit eigenvector for `lambda`, with the sign chosen so the
 * angle is in (−π/2, π/2]. The two candidate rows are used because one of them
 * can be the zero vector when the eigenvalue equals an entry.
 */
function eigenvector(A: number, B: number, C: number, lambda: number): Vec2 {
  const half = B / 2
  const first: Vec2 = [half, lambda - A]
  const second: Vec2 = [lambda - C, half]
  const n1 = Math.hypot(first[0], first[1])
  const n2 = Math.hypot(second[0], second[1])
  let vector = n1 >= n2 ? first : second
  const length = Math.hypot(vector[0], vector[1])
  if (length <= EPS) return [1, 0]
  vector = [vector[0] / length, vector[1] / length]
  if (vector[0] < -1e-15 || (Math.abs(vector[0]) <= 1e-15 && vector[1] < 0)) {
    vector = [-vector[0], -vector[1]]
  }
  return vector
}

const perpendicular = (v: Vec2): Vec2 => [-v[1], v[0]]
const dot = (a: Vec2, b: Vec2) => a[0] * b[0] + a[1] * b[1]
const scale = (v: Vec2, k: number): Vec2 => [v[0] * k, v[1] * k]
const add = (a: Vec2, b: Vec2): Vec2 => [a[0] + b[0], a[1] + b[1]]
const normalize = (v: Vec2): Vec2 => {
  const length = Math.hypot(v[0], v[1])
  return length <= EPS ? [1, 0] : [v[0] / length, v[1] / length]
}
const angleOf = (v: Vec2) => Math.atan2(v[1], v[0])
/** Keep −0 out of the readout: it prints as "-0" but means the same point. */
const clean = (value: number) => (value === 0 ? 0 : value)

function interceptsOf(a: number, b: number, c: number): Intercepts {
  const local = Math.max(1, Math.abs(a), Math.abs(b), Math.abs(c))
  const zero = EPS * local
  if (Math.abs(a) <= zero) {
    if (Math.abs(b) <= zero) {
      return Math.abs(c) <= zero
        ? { status: 'coincident', values: [] }
        : { status: 'none', values: [] }
    }
    return { status: 'one', values: [clean(-c / b)] }
  }
  const d = b * b - 4 * a * c
  const tau = EPS * Math.max(1, b * b, Math.abs(4 * a * c))
  if (d < -tau) return { status: 'none', values: [] }
  if (Math.abs(d) <= tau) return { status: 'double', values: [clean(-b / (2 * a))] }
  const root = Math.sqrt(d)
  const values = [(-b - root) / (2 * a), (-b + root) / (2 * a)].map(clean).sort((x, y) => x - y)
  return { status: 'two', values }
}

function emptyLines(): ConicLine[] {
  return []
}

function invalid(note: string): ConicGeometry {
  return {
    kind: 'invalid',
    eccentricity: null,
    lines: emptyLines(),
    intercepts: { x: { status: 'none', values: [] }, y: { status: 'none', values: [] } },
    certainty: 'invalid',
    note,
  }
}

/**
 * Clip an infinite line to an axis-aligned rectangle. The renderer uses this for
 * degenerate lines, so a line through the plane never becomes a coordinate of
 * 10^15 in the SVG.
 */
export function clipLineToRect(
  point: Vec2,
  direction: Vec2,
  bounds: { x0: number; x1: number; y0: number; y1: number },
): [Vec2, Vec2] | null {
  const dx = direction[0]
  const dy = direction[1]
  if (Math.abs(dx) <= EPS && Math.abs(dy) <= EPS) return null
  let tMin = Number.NEGATIVE_INFINITY
  let tMax = Number.POSITIVE_INFINITY
  const test = (origin: number, delta: number, low: number, high: number) => {
    if (Math.abs(delta) <= EPS) {
      return origin >= low && origin <= high
    }
    const t1 = (low - origin) / delta
    const t2 = (high - origin) / delta
    tMin = Math.max(tMin, Math.min(t1, t2))
    tMax = Math.min(tMax, Math.max(t1, t2))
    return true
  }
  if (!test(point[0], dx, bounds.x0, bounds.x1)) return null
  if (!test(point[1], dy, bounds.y0, bounds.y1)) return null
  if (tMin > tMax) return null
  return [
    [point[0] + dx * tMin, point[1] + dy * tMin],
    [point[0] + dx * tMax, point[1] + dy * tMax],
  ]
}

/**
 * Classify the general six-term conic and return the geometry to draw.
 *
 * The quadratic block, the affine vector and the constant term are scaled
 * separately. A single max-all normalisation would let a large constant term
 * erase a small but genuine quadratic rank, which is exactly the case a
 * translated conic produces.
 */
export function classifyGeneralConic(
  A: number,
  B: number,
  C: number,
  D: number,
  E: number,
  F: number,
): ConicGeometry {
  const coefficients = [A, B, C, D, E, F]
  if (coefficients.some((value) => !Number.isFinite(value))) {
    return invalid('coefficients must be finite')
  }

  const qScale = Math.max(Math.abs(A), Math.abs(B), Math.abs(C))
  const intercepts = {
    x: interceptsOf(A, D, F),
    y: interceptsOf(C, E, F),
  }

  if (qScale === 0) {
    // Rank 0: the equation is linear. Scale the linear test by its own terms.
    const linearScale = Math.max(Math.abs(D), Math.abs(E), Math.abs(F))
    const tau = EPS * linearScale
    const linearPart = Math.max(Math.abs(D), Math.abs(E))
    const length = Math.hypot(D, E)
    if (linearScale === 0) {
      return {
        kind: 'all-plane',
        eccentricity: null,
        lines: emptyLines(),
        intercepts,
        certainty: 'exact',
        note: 'every point satisfies the equation',
      }
    }
    if (linearPart <= tau) {
      return {
        kind: 'empty',
        eccentricity: null,
        lines: emptyLines(),
        intercepts,
        certainty: 'exact',
        note: 'no real points',
      }
    }
    const point: Vec2 = [(-F * D) / (length * length), (-F * E) / (length * length)]
    return {
      kind: 'line',
      eccentricity: null,
      lines: [{ point, direction: normalize([-E, D]) }],
      intercepts,
      certainty: 'exact',
    }
  }

  // Normalise the quadratic block before its eigenvalues, so a genuinely small
  // quadratic part keeps its rank. The rank test is relative to the block.
  const Ah = A / qScale
  const Bh = B / qScale
  const Ch = C / qScale
  const [lowHat, highHat] = eigenvalues(Ah, Bh, Ch)
  const rank2 = Math.abs(lowHat) > EPS && Math.abs(highHat) > EPS
  const rank1 = !rank2 && (Math.abs(lowHat) > EPS || Math.abs(highHat) > EPS)

  const affineScale = Math.hypot(D, E)
  const L0 = Math.max(1, affineScale / qScale, Math.sqrt(Math.abs(F) / qScale))
  const tauE = EPS * qScale * L0
  const tauK = EPS * qScale * L0 * L0

  if (rank2) {
    const determinant = Ah * Ch - (Bh / 2) * (Bh / 2)
    // Rank 2 guarantees this is away from zero, but guard it anyway.
    if (Math.abs(determinant) <= EPS) return invalid('quadratic form is numerically singular')

    const cx = (-Ch * D + (Bh / 2) * E) / (2 * qScale * determinant)
    const cy = ((Bh / 2) * D - Ah * E) / (2 * qScale * determinant)
    const center: Vec2 = [cx, cy]
    const K = qScale * (Ah * cx * cx + Bh * cx * cy + Ch * cy * cy) - F
    const nearRank = Math.min(Math.abs(lowHat), Math.abs(highHat)) <= NEAR
    const nearK = Math.abs(K) <= NEAR * qScale * L0 * L0
    const certainty: Certainty = nearRank || (nearK && K !== 0) ? 'near-degenerate' : 'exact'
    const note = certainty === 'near-degenerate' ? 'near a degenerate case; the classification is kept' : undefined

    const highVector = eigenvector(Ah, Bh, Ch, highHat)
    const lowVector = eigenvector(Ah, Bh, Ch, lowHat)
    const definite = lowHat * highHat > 0

    if (Math.abs(K) <= tauK) {
      if (definite) {
        return {
          kind: 'point',
          center,
          eccentricity: null,
          lines: emptyLines(),
          intercepts,
          certainty,
          note: note ?? 'the curve is a single point',
        }
      }
      const m = Math.sqrt(Math.max(0, -lowHat / highHat))
      const first = normalize(add(highVector, scale(lowVector, m)))
      const second = normalize(add(highVector, scale(lowVector, -m)))
      return {
        kind: 'intersecting',
        center,
        eccentricity: null,
        lines: [
          { point: center, direction: first },
          { point: center, direction: second },
        ],
        intercepts,
        certainty,
        note,
      }
    }

    const highRatio = K / qScale / highHat
    const lowRatio = K / qScale / lowHat

    if (highRatio > 0 && lowRatio > 0) {
      const highLength = Math.sqrt(highRatio)
      const lowLength = Math.sqrt(lowRatio)
      const major = highLength <= lowLength ? lowVector : highVector
      const minor = highLength <= lowLength ? highVector : lowVector
      const a = Math.max(highLength, lowLength)
      const b = Math.min(highLength, lowLength)
      const circle = Math.abs(a - b) <= NEAR * Math.max(1, a, b)
      const kind: ConicKind = circle ? 'circle' : 'ellipse'
      const theta = angleOf(circle ? [1, 0] : minor)
      const c = circle ? 0 : Math.sqrt(Math.max(0, a * a - b * b))
      return {
        kind,
        center,
        major: circle ? [1, 0] : major,
        minor: circle ? [0, 1] : minor,
        a,
        b,
        theta,
        foci: circle
          ? [center]
          : [add(center, scale(major, -c)), add(center, scale(major, c))],
        focalConstant: 2 * a,
        eccentricity: circle ? 0 : Math.sqrt(Math.max(0, 1 - (b * b) / (a * a))),
        lines: emptyLines(),
        intercepts,
        certainty,
        note,
      }
    }

    if (highRatio < 0 && lowRatio < 0) {
      return {
        kind: 'empty',
        center,
        eccentricity: null,
        lines: emptyLines(),
        intercepts,
        certainty,
        note: note ?? 'no real points',
      }
    }

    const transverse = highRatio > 0 ? highVector : lowVector
    const conjugate = highRatio > 0 ? lowVector : highVector
    const a = Math.sqrt(Math.abs(highRatio > 0 ? highRatio : lowRatio))
    const b = Math.sqrt(Math.abs(highRatio > 0 ? lowRatio : highRatio))
    const c = Math.sqrt(a * a + b * b)
    return {
      kind: 'hyperbola',
      center,
      transverse,
      conjugate,
      a,
      b,
      theta: angleOf(transverse),
      foci: [add(center, scale(transverse, -c)), add(center, scale(transverse, c))],
      focalConstant: 2 * a,
      eccentricity: Math.sqrt(1 + (b * b) / (a * a)),
      lines: emptyLines(),
      intercepts,
      certainty,
      note,
    }
  }

  if (rank1) {
    const lambdaHat = Math.abs(lowHat) > Math.abs(highHat) ? lowHat : highHat
    const lambda = lambdaHat * qScale
    if (lambda === 0) return invalid('the quadratic part is too small to resolve')
    const squared = eigenvector(Ah, Bh, Ch, lambdaHat)
    const nullDirection = perpendicular(squared)
    const Dp = dot([D, E], squared)
    const Ep = dot([D, E], nullDirection)
    const h = -Dp / (2 * lambda)
    const nearE = Math.abs(Ep) <= NEAR * qScale * L0

    if (Math.abs(Ep) > tauE) {
      const k = ((Dp * Dp) / (4 * lambda) - F) / Ep
      const vertex = add(scale(squared, h), scale(nullDirection, k))
      const pSigned = -Ep / (4 * lambda)
      const opening = scale(nullDirection, Math.sign(pSigned) || 1)
      const p = Math.abs(pSigned)
      return {
        kind: 'parabola',
        vertex,
        squared,
        opening,
        p,
        theta: angleOf(squared),
        foci: [add(vertex, scale(opening, p))],
        directrix: { point: add(vertex, scale(opening, -p)), direction: squared },
        eccentricity: 1,
        lines: emptyLines(),
        intercepts,
        certainty: nearE ? 'near-degenerate' : 'exact',
        note: nearE ? 'a very shallow parabola is near numerical resolution' : undefined,
      }
    }

    const R = ((Dp * Dp) / (4 * lambda) - F) / lambda
    const tauR = EPS * L0 * L0
    const nearR = Math.abs(R) <= NEAR * L0 * L0
    if (R > tauR) {
      const root = Math.sqrt(R)
      const pointA = scale(squared, h - root)
      const pointB = scale(squared, h + root)
      return {
        kind: 'parallel',
        lines: [
          { point: pointA, direction: nullDirection },
          { point: pointB, direction: nullDirection },
        ],
        eccentricity: null,
        intercepts,
        certainty: nearR ? 'near-degenerate' : 'exact',
        note: nearR ? 'close to a degenerate case; the classification is kept' : undefined,
      }
    }
    if (Math.abs(R) <= tauR) {
      return {
        kind: 'line',
        lines: [{ point: scale(squared, h), direction: nullDirection }],
        eccentricity: null,
        intercepts,
        certainty: 'exact',
      }
    }
    return {
      kind: 'empty',
      eccentricity: null,
      lines: emptyLines(),
      intercepts,
      certainty: nearR ? 'near-degenerate' : 'exact',
      note: nearR ? 'close to a degenerate case; the classification is kept' : undefined,
    }
  }

  return invalid('the equation has no quadratic or linear part')
}

export function generalConic(
  A: number,
  B: number,
  C: number,
  D: number,
  E: number,
  F: number,
  x: number,
  y: number,
): number {
  return A * x * x + B * x * y + C * y * y + D * x + E * y + F
}

export type ConicBranch = 1 | -1

export interface ConicParameterRange {
  min: number
  max: number
  step: number
}

export interface ConicProjection {
  point: Vec2
  parameter: number
  distance: number
  branch: ConicBranch
}

export interface FocalMeasurements {
  toFoci: number[]
  /** Ellipse and circle only. */
  focalSum: number | null
  /** Hyperbola only. */
  focalDifference: number | null
  /** Parabola only. */
  focusDistance: number | null
  /** Parabola only. */
  directrixDistance: number | null
}

/** The parameter interval a probe slider should expose for this geometry. */
export function conicParameterRange(geometry: ConicGeometry): ConicParameterRange | null {
  if (geometry.kind === 'circle' || geometry.kind === 'ellipse') {
    return { min: 0, max: 2 * Math.PI, step: Math.PI / 90 }
  }
  if (geometry.kind === 'hyperbola') {
    return { min: -4, max: 4, step: 0.04 }
  }
  if (geometry.kind === 'parabola' && geometry.p !== undefined && geometry.p > 0) {
    // u is the signed distance along the squared direction. This range reaches
    // a little past the plot on both sides without making the slider useless
    // for a very large or very small p.
    const span = Math.max(1, Math.sqrt(32 * geometry.p))
    return { min: -span, max: span, step: span / 100 }
  }
  return null
}

/**
 * A point on the classified curve. The parameter is the natural one for the
 * kind: an angle for an ellipse or circle, a hyperbolic parameter for a
 * hyperbola, and the signed offset along the squared direction for a parabola.
 */
export function conicPointAt(
  geometry: ConicGeometry,
  parameter: number,
  branch: ConicBranch = 1,
): Vec2 | null {
  if (!Number.isFinite(parameter)) return null
  if (geometry.kind === 'circle' || geometry.kind === 'ellipse') {
    if (
      !geometry.center ||
      !geometry.major ||
      !geometry.minor ||
      geometry.a === undefined ||
      geometry.b === undefined ||
      geometry.a <= 0 ||
      geometry.b <= 0
    ) {
      return null
    }
    return add(
      geometry.center,
      add(scale(geometry.major, geometry.a * Math.cos(parameter)), scale(geometry.minor, geometry.b * Math.sin(parameter))),
    )
  }
  if (geometry.kind === 'hyperbola') {
    if (
      !geometry.center ||
      !geometry.transverse ||
      !geometry.conjugate ||
      geometry.a === undefined ||
      geometry.b === undefined ||
      geometry.a <= 0 ||
      geometry.b <= 0
    ) {
      return null
    }
    const t = Math.max(-20, Math.min(20, parameter))
    const along = geometry.a * Math.cosh(t)
    const across = geometry.b * Math.sinh(t)
    return add(
      geometry.center,
      add(
        scale(geometry.transverse, branch * along),
        scale(geometry.conjugate, branch * across),
      ),
    )
  }
  if (
    geometry.kind === 'parabola' &&
    geometry.vertex &&
    geometry.squared &&
    geometry.opening &&
    geometry.p !== undefined &&
    geometry.p > 0
  ) {
    const u = parameter
    return add(
      geometry.vertex,
      add(scale(geometry.opening, (u * u) / (4 * geometry.p)), scale(geometry.squared, u)),
    )
  }
  return null
}

/**
 * Project a pointer onto the conic. The returned point is a genuine point of
 * the classified curve, so every readout is computed from a point on the curve
 * rather than from the pointer's approximate position.
 */
export function projectToConic(
  geometry: ConicGeometry,
  target: Vec2,
  branch: ConicBranch = 1,
): ConicProjection | null {
  const range = conicParameterRange(geometry)
  if (!range) return null
  if (!Number.isFinite(target[0]) || !Number.isFinite(target[1])) return null

  if (geometry.kind === 'circle' && geometry.center && geometry.a !== undefined) {
    const dx = target[0] - geometry.center[0]
    const dy = target[1] - geometry.center[1]
    const length = Math.hypot(dx, dy)
    const unit: Vec2 = length <= EPS ? [1, 0] : [dx / length, dy / length]
    const point: Vec2 = add(geometry.center, scale(unit, geometry.a))
    const parameter = ((Math.atan2(unit[1], unit[0]) % (2 * Math.PI)) + 2 * Math.PI) % (2 * Math.PI)
    return { point, parameter, distance: Math.abs(length - geometry.a), branch }
  }

  const evaluate = (parameter: number): number => {
    const candidate = conicPointAt(geometry, parameter, branch)
    return candidate ? Math.hypot(candidate[0] - target[0], candidate[1] - target[1]) : Number.POSITIVE_INFINITY
  }

  const samples = 96
  let best = range.min
  let bestValue = Number.POSITIVE_INFINITY
  for (let index = 0; index <= samples; index++) {
    const parameter = range.min + ((range.max - range.min) * index) / samples
    const value = evaluate(parameter)
    if (value < bestValue) {
      bestValue = value
      best = parameter
    }
  }

  // Golden-section search around the best sample. Both curve families here are
  // smooth and the sampled bracket contains the local minimum; the fallback is
  // the sample itself when the bracket is flat.
  const span = (range.max - range.min) / samples
  let low = Math.max(range.min, best - span)
  let high = Math.min(range.max, best + span)
  const ratio = (Math.sqrt(5) - 1) / 2
  let left = high - ratio * (high - low)
  let right = low + ratio * (high - low)
  let leftValue = evaluate(left)
  let rightValue = evaluate(right)
  for (let index = 0; index < 60 && high - low > 1e-10; index++) {
    if (leftValue <= rightValue) {
      high = right
      right = left
      rightValue = leftValue
      left = high - ratio * (high - low)
      leftValue = evaluate(left)
    } else {
      low = left
      left = right
      leftValue = rightValue
      right = low + ratio * (high - low)
      rightValue = evaluate(right)
    }
  }
  const parameter = (low + high) / 2
  const point = conicPointAt(geometry, parameter, branch)
  if (!point) return null
  return {
    point,
    parameter,
    distance: Math.hypot(point[0] - target[0], point[1] - target[1]),
    branch,
  }
}

/** Distances from a point to the classified focal objects. */
export function focalMeasurements(geometry: ConicGeometry, point: Vec2): FocalMeasurements {
  const foci = geometry.foci ?? []
  const toFoci = foci.map((focus) => Math.hypot(point[0] - focus[0], point[1] - focus[1]))
  let focalSum: number | null = null
  let focalDifference: number | null = null
  let focusDistance: number | null = null
  let directrixDistance: number | null = null

  if (geometry.kind === 'circle' && toFoci.length === 1) {
    focalSum = 2 * toFoci[0]
  } else if (geometry.kind === 'ellipse' && toFoci.length === 2) {
    focalSum = toFoci[0] + toFoci[1]
  } else if (geometry.kind === 'hyperbola' && toFoci.length === 2) {
    focalDifference = Math.abs(toFoci[0] - toFoci[1])
  } else if (geometry.kind === 'parabola' && toFoci.length === 1) {
    focusDistance = toFoci[0]
  }

  if (geometry.directrix) {
    const normal = geometry.opening ?? perpendicular(geometry.directrix.direction)
    directrixDistance = Math.abs(
      (point[0] - geometry.directrix.point[0]) * normal[0] +
        (point[1] - geometry.directrix.point[1]) * normal[1],
    )
  }

  return { toFoci, focalSum, focalDifference, focusDistance, directrixDistance }
}

/** The old homogeneous form, kept for callers and tests that still use it. */
export function classifyQuadratic(A: number, B: number, C: number): ConicShape {
  const [low, high] = eigenvalues(A, B, C)
  const eps = 1e-9
  const squaredTerms = (Math.abs(A) > eps ? 1 : 0) + (Math.abs(C) > eps ? 1 : 0)
  const sameSign = Math.abs(A) > eps && Math.abs(C) > eps && A * C > 0

  // A zero discriminant here means the quadratic form is degenerate, and
  // Ax² + Bxy + Cy² = 1 with one zero eigenvalue is a pair of parallel lines, not
  // a parabola: a parabola needs a linear term, and this equation has none.
  let type: ConicType
  if (Math.abs(low) < eps && Math.abs(high) < eps) type = 'empty'
  else if (high <= eps) type = 'empty'
  else if (Math.abs(low) < eps) type = 'parallel'
  else if (low < 0) type = 'hyperbola'
  else type = Math.abs(low - high) < 1e-6 ? 'circle' : 'ellipse'

  const eccentricity =
    type === 'circle' || type === 'ellipse' || type === 'hyperbola'
      ? Math.sqrt(Math.max(0, 1 - low / high))
      : Number.NaN

  return { type, squaredTerms, sameSign, eccentricity }
}

/** Where the old homogeneous curve meets the x axis, if it does. */
export function xIntercepts(A: number, C: number): number[] {
  if (A <= 1e-9) return []
  const r = 1 / Math.sqrt(A)
  return [-r, r]
}

/** Where the old homogeneous curve meets the y axis, if it does. */
export function yIntercepts(A: number, C: number): number[] {
  if (C <= 1e-9) return []
  const r = 1 / Math.sqrt(C)
  return [-r, r]
}
