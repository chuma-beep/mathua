// The complex line integral along a parameterised path.
//
// This is not an area. A real definite integral can be drawn as the region under
// a curve because its value is that area; a complex line integral accumulates a
// complex number whose value depends on the path, so the figure shows the path in
// the z-plane and the running value in the w-plane. Everything here is plain
// arithmetic with no React and no rendering, so the quadrature can be tested
// against closed-form results without a DOM.

export type Complex = [number, number]

export type Path =
  | { kind: 'line'; a: Complex; b: Complex }
  | { kind: 'circle'; center: Complex; radius: number }
  | { kind: 'ellipse'; center: Complex; rx: number; ry: number }
  | { kind: 'wavy'; center: Complex; radius: number; amplitude: number; lobes: number }

export type ContourFn = 'inv' | 'identity' | 'square' | 'conj' | 'exp'

export const FnLabels: Record<ContourFn, string> = {
  inv: '1 / z',
  identity: 'z',
  square: 'z²',
  conj: 'conj z',
  exp: 'e^z',
}

export function mul(a: Complex, b: Complex): Complex {
  return [a[0] * b[0] - a[1] * b[1], a[0] * b[1] + a[1] * b[0]]
}

export function add(a: Complex, b: Complex): Complex {
  return [a[0] + b[0], a[1] + b[1]]
}

/** f(z), or null at a pole where the function has no value. */
export function fnValue(fn: ContourFn, z: Complex): Complex | null {
  const [x, y] = z
  switch (fn) {
    case 'inv': {
      const d = x * x + y * y
      if (d < 1e-12) return null
      return [x / d, -y / d]
    }
    case 'identity':
      return [x, y]
    case 'square':
      return [x * x - y * y, 2 * x * y]
    case 'conj':
      return [x, -y]
    case 'exp': {
      const e = Math.exp(x)
      return [e * Math.cos(y), e * Math.sin(y)]
    }
  }
}

/**
 * The point and the derivative at parameter t. Direction flips the parameter, so
 * traversing backwards negates the integral by negating ds rather than by a
 * second code path.
 */
export function pathAt(path: Path, direction: 1 | -1, t: number): { z: Complex; dz: Complex } {
  const s = direction === 1 ? t : 1 - t
  const ds = direction === 1 ? 1 : -1
  const twoPi = 2 * Math.PI
  switch (path.kind) {
    case 'line': {
      const [ax, ay] = path.a
      const [bx, by] = path.b
      return { z: [ax + (bx - ax) * s, ay + (by - ay) * s], dz: [(bx - ax) * ds, (by - ay) * ds] }
    }
    case 'circle': {
      const angle = twoPi * s
      const [cx, cy] = path.center
      return {
        z: [cx + path.radius * Math.cos(angle), cy + path.radius * Math.sin(angle)],
        dz: [-path.radius * twoPi * Math.sin(angle) * ds, path.radius * twoPi * Math.cos(angle) * ds],
      }
    }
    case 'ellipse': {
      const angle = twoPi * s
      const [cx, cy] = path.center
      return {
        z: [cx + path.rx * Math.cos(angle), cy + path.ry * Math.sin(angle)],
        dz: [-path.rx * twoPi * Math.sin(angle) * ds, path.ry * twoPi * Math.cos(angle) * ds],
      }
    }
    case 'wavy': {
      const angle = twoPi * s
      const r = path.radius + path.amplitude * Math.cos(path.lobes * angle)
      const dr = -path.amplitude * path.lobes * Math.sin(path.lobes * angle)
      const [cx, cy] = path.center
      return {
        z: [cx + r * Math.cos(angle), cy + r * Math.sin(angle)],
        dz: [
          (dr * Math.cos(angle) - r * twoPi * Math.sin(angle)) * ds,
          (dr * Math.sin(angle) + r * twoPi * Math.cos(angle)) * ds,
        ],
      }
    }
  }
}

function pointSegmentDistance(p: Complex, a: Complex, b: Complex): number {
  const vx = b[0] - a[0]
  const vy = b[1] - a[1]
  const wx = p[0] - a[0]
  const wy = p[1] - a[1]
  const lengthSquared = vx * vx + vy * vy
  if (lengthSquared === 0) return Math.hypot(wx, wy)
  const t = Math.max(0, Math.min(1, (wx * vx + wy * vy) / lengthSquared))
  return Math.hypot(wx - t * vx, wy - t * vy)
}

export type PoleCertainty = 'exact' | 'converged' | 'uncertain'

export type PoleInfo = {
  crosses: boolean
  closest: number
  /** How the number was obtained, so the figure never overstates its evidence. */
  certainty: PoleCertainty
}

/**
 * A path within this distance of the pole is treated as passing through it.
 *
 * Deliberately tiny. A numerical epsilon is not a geometric tolerance: a curve
 * that misses the pole by 10^-4 is a near miss, and widening this to something
 * like 10^-3 would report a divergent integral for a path on which every value
 * exists.
 */
const CROSS_EPS = 1e-9

const TWO_PI = 2 * Math.PI

function distanceToOrigin(z: Complex): number {
  return Math.hypot(z[0], z[1])
}

/**
 * The minimum distance from a closed curve to the origin.
 *
 * Samples only seed the search: each sample whose distance is no larger than
 * both neighbours brackets a candidate minimum, and golden-section search inside
 * that bracket produces the value. Minimising the distance and minimising its
 * square share the same minimiser, and the distance needs no square root in the
 * branch condition. This is robust where bracketing the derivative is not: a
 * curve that touches the pole is tangent there, so its derivative changes sign
 * twice at one point and a sign-change test can miss it entirely.
 *
 * `certainty` is `converged` only when at least one bracket actually narrowed to
 * below 1e-9; otherwise the number exists but the search failed, and the figure
 * must not present it as settled.
 */
function minimumDistance(point: (t: number) => Complex, samples: number): { closest: number; certainty: PoleCertainty } {
  const at = (t: number) => distanceToOrigin(point(t))
  const step = TWO_PI / samples
  let best = Number.POSITIVE_INFINITY
  let converged = 0
  for (let i = 0; i < samples; i++) {
    const lo0 = (i - 1) * step
    const mid = i * step
    const hi0 = (i + 1) * step
    if (at(mid) > at(lo0) || at(mid) > at(hi0)) continue
    let lo = lo0
    let hi = hi0
    const phi = (Math.sqrt(5) - 1) / 2
    let x1 = hi - phi * (hi - lo)
    let x2 = lo + phi * (hi - lo)
    let f1 = at(x1)
    let f2 = at(x2)
    for (let k = 0; k < 120; k++) {
      if (f1 < f2) {
        hi = x2
        x2 = x1
        f2 = f1
        x1 = hi - phi * (hi - lo)
        f1 = at(x1)
      } else {
        lo = x1
        x1 = x2
        f1 = f2
        x2 = lo + phi * (hi - lo)
        f2 = at(x2)
      }
    }
    best = Math.min(best, at((lo + hi) / 2))
    if (hi - lo < 1e-9) converged++
  }
  if (!Number.isFinite(best)) return { closest: best, certainty: 'uncertain' }
  return { closest: best, certainty: converged > 0 ? 'converged' : 'uncertain' }
}

function ellipseClosest(path: Extract<Path, { kind: 'ellipse' }>): { closest: number; certainty: PoleCertainty } {
  const [cx, cy] = path.center
  const { rx, ry } = path
  return minimumDistance((t) => [cx + rx * Math.cos(t), cy + ry * Math.sin(t)], 512)
}

function wavyClosest(path: Extract<Path, { kind: 'wavy' }>): { closest: number; certainty: PoleCertainty } {
  const [cx, cy] = path.center
  const { radius, amplitude, lobes } = path
  return minimumDistance((t) => {
    const r = radius + amplitude * Math.cos(lobes * t)
    return [cx + r * Math.cos(t), cy + r * Math.sin(t)]
  }, 1024)
}

/**
 * How close the path comes to the origin, where 1/z has its pole.
 *
 * Each family is solved on its own terms rather than by sampling the curve. A
 * segment is exact point-to-segment distance and a circle is exact too, because
 * its closest point lies on the line from the centre to the origin. An ellipse
 * and the wavy family have no closed form here, so they are minimised by
 * bracketing the stationary points of the squared distance and bisecting each
 * bracket; `certainty` reports whether that convergence succeeded. The figure
 * must not present a `converged` or `uncertain` result as a proof.
 */
export function poleInfo(path: Path): PoleInfo {
  if (path.kind === 'line') {
    const closest = pointSegmentDistance([0, 0], path.a, path.b)
    return { crosses: closest < CROSS_EPS, closest, certainty: 'exact' }
  }
  if (path.kind === 'circle') {
    const centreDistance = Math.hypot(path.center[0], path.center[1])
    const closest = Math.abs(centreDistance - path.radius)
    return { crosses: closest < CROSS_EPS, closest, certainty: 'exact' }
  }
  const solved = path.kind === 'ellipse' ? ellipseClosest(path) : wavyClosest(path)
  return { crosses: solved.closest < CROSS_EPS, closest: solved.closest, certainty: solved.certainty }
}

export type Integration = { partial: Complex[]; total: Complex }

/**
 * Midpoint rule over the parameter. `partial[k]` is the integral from 0 up to
 * (k+1)/samples, which is what the figure draws as the running value.
 */
export function integrate(path: Path, fn: ContourFn, direction: 1 | -1, samples = 2000): Integration {
  if (poleInfo(path).crosses) return { partial: [], total: [Number.NaN, Number.NaN] }
  const partial: Complex[] = []
  let sum: Complex = [0, 0]
  const step = 1 / samples
  for (let i = 0; i < samples; i++) {
    const mid = (i + 0.5) * step
    const { z, dz } = pathAt(path, direction, mid)
    const value = fnValue(fn, z)
    if (!value) return { partial: [], total: [Number.NaN, Number.NaN] }
    sum = add(sum, mul(value, dz))
    partial.push([sum[0] * step, sum[1] * step])
  }
  return { partial, total: [sum[0] * step, sum[1] * step] }
}

/**
 * The closed-form value where one exists, used by the tests and never by the
 * drawing. `conj` is not analytic, so its circle integral is 2πi r² rather than
 * the zero the Cauchy theorem would give an analytic function.
 */
export function analytic(path: Path, fn: ContourFn, direction: 1 | -1): Complex | null {
  if (path.kind === 'line') return null
  const sign = direction === 1 ? 1 : -1
  const twoPi = 2 * Math.PI
  if (fn === 'inv') {
    const inside =
      path.kind === 'circle'
        ? Math.hypot(path.center[0], path.center[1]) < path.radius
        : path.kind === 'ellipse'
          ? (path.center[0] * path.center[0]) / (path.rx * path.rx) + (path.center[1] * path.center[1]) / (path.ry * path.ry) < 1
          : Math.hypot(path.center[0], path.center[1]) < path.radius - path.amplitude
    return inside ? [0, sign * twoPi] : [0, 0]
  }
  if (fn === 'identity' || fn === 'square' || fn === 'exp') return [0, 0]
  if (path.kind === 'circle') return [0, sign * twoPi * path.radius * path.radius]
  return null
}
