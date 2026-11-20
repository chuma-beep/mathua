import { describe, expect, it } from 'vitest'
import {
  classifyGeneralConic,
  clipLineToRect,
  conicParameterRange,
  conicPointAt,
  discriminant,
  focalMeasurements,
  generalConic,
  projectToConic,
  type ConicGeometry,
  type ConicKind,
  type Vec2,
} from '../components/visual/math/conic'

const near = (a: number, b: number, tol = 1e-9) => expect(Math.abs(a - b)).toBeLessThan(tol)
const close = (a: number | undefined, b: number, tol = 1e-9) => {
  expect(a).toBeTypeOf('number')
  near(a as number, b, tol)
}
const relative = (a: number | undefined, b: number, tol = 1e-6) => {
  expect(a).toBeTypeOf('number')
  expect(Math.abs((a as number) - b)).toBeLessThanOrEqual(tol * Math.max(1, Math.abs(a as number), Math.abs(b)))
}
const norm = (v: Vec2) => Math.hypot(v[0], v[1])
const parallel = (a: Vec2, b: Vec2, tol = 1e-9) => {
  near(norm(a), 1, tol)
  near(norm(b), 1, tol)
  near(a[0] * b[1] - a[1] * b[0], 0, tol)
}
const point = (a: Vec2 | undefined, b: Vec2, tol = 1e-9) => {
  expect(a).toBeDefined()
  near((a as Vec2)[0], b[0], tol)
  near((a as Vec2)[1], b[1], tol)
}
const kinds: ConicKind[] = [
  'circle',
  'ellipse',
  'hyperbola',
  'parabola',
  'point',
  'intersecting',
  'parallel',
  'line',
  'empty',
  'all-plane',
  'invalid',
]
const finiteFields = (g: ConicGeometry) => {
  for (const value of [g.a, g.b, g.p, g.theta, g.eccentricity, g.focalConstant]) {
    if (value !== undefined && value !== null) expect(Number.isFinite(value)).toBe(true)
  }
  for (const v of [g.center, g.vertex, g.major, g.minor, g.transverse, g.conjugate, g.squared, g.opening]) {
    if (v) {
      expect(Number.isFinite(v[0])).toBe(true)
      expect(Number.isFinite(v[1])).toBe(true)
    }
  }
  for (const focus of g.foci ?? []) {
    expect(Number.isFinite(focus[0])).toBe(true)
    expect(Number.isFinite(focus[1])).toBe(true)
  }
  if (g.directrix) {
    expect(Number.isFinite(g.directrix.point[0])).toBe(true)
    expect(Number.isFinite(g.directrix.point[1])).toBe(true)
    expect(Number.isFinite(g.directrix.direction[0])).toBe(true)
    expect(Number.isFinite(g.directrix.direction[1])).toBe(true)
  }
}

const residual = (g: ConicGeometry, p: Vec2): number => {
  if (g.kind === 'circle' && g.center && g.a !== undefined) {
    return Math.abs(Math.hypot(p[0] - g.center[0], p[1] - g.center[1]) - g.a)
  }
  if (g.kind === 'ellipse' && g.center && g.major && g.minor && g.a !== undefined && g.b !== undefined) {
    const dx = p[0] - g.center[0]
    const dy = p[1] - g.center[1]
    const x = (dx * g.major[0] + dy * g.major[1]) / g.a
    const y = (dx * g.minor[0] + dy * g.minor[1]) / g.b
    return Math.abs(x * x + y * y - 1)
  }
  if (g.kind === 'hyperbola' && g.center && g.transverse && g.conjugate && g.a !== undefined && g.b !== undefined) {
    const dx = p[0] - g.center[0]
    const dy = p[1] - g.center[1]
    const x = (dx * g.transverse[0] + dy * g.transverse[1]) / g.a
    const y = (dx * g.conjugate[0] + dy * g.conjugate[1]) / g.b
    return Math.abs(x * x - y * y - 1)
  }
  if (g.kind === 'parabola' && g.vertex && g.squared && g.opening && g.p !== undefined) {
    const dx = p[0] - g.vertex[0]
    const dy = p[1] - g.vertex[1]
    const along = dx * g.squared[0] + dy * g.squared[1]
    const axial = dx * g.opening[0] + dy * g.opening[1]
    return Math.abs(axial - (along * along) / (4 * g.p))
  }
  return Number.POSITIVE_INFINITY
}

const onConic = (g: ConicGeometry, p: Vec2 | null, tol = 1e-8) => {
  expect(p).not.toBeNull()
  near(residual(g, p as Vec2), 0, tol)
}

const hasFocus = (g: ConicGeometry, expected: Vec2, tol = 1e-9) => {
  const found = (g.foci ?? []).some(
    (focus) => Math.abs(focus[0] - expected[0]) <= tol && Math.abs(focus[1] - expected[1]) <= tol,
  )
  expect(found).toBe(true)
}

describe('classifyGeneralConic: the twenty canonical cases', () => {
  it('classifies the default ellipse', () => {
    const g = classifyGeneralConic(1, 0, 0.5, 0, 0, -1)
    expect(g.kind).toBe('ellipse')
    point(g.center, [0, 0])
    close(g.a, Math.SQRT2)
    close(g.b, 1)
    close(g.eccentricity as number, Math.sqrt(0.5))
    parallel(g.major as Vec2, [0, 1])
    parallel(g.minor as Vec2, [1, 0])
  })

  it('classifies a circle', () => {
    const g = classifyGeneralConic(1, 0, 1, 0, 0, -1)
    expect(g.kind).toBe('circle')
    point(g.center, [0, 0])
    close(g.a, 1)
    close(g.b, 1)
    expect(g.eccentricity).toBe(0)
  })

  it('classifies a translated circle', () => {
    const g = classifyGeneralConic(1, 0, 1, -2, 2, -2)
    expect(g.kind).toBe('circle')
    point(g.center, [1, -1])
    close(g.a, 2)
    close(g.b, 2)
    expect(g.eccentricity).toBe(0)
  })

  it('classifies a translated ellipse', () => {
    const g = classifyGeneralConic(1, 0, 4, -2, 16, 13)
    expect(g.kind).toBe('ellipse')
    point(g.center, [1, -2])
    close(g.a, 2)
    close(g.b, 1)
    close(g.eccentricity as number, Math.sqrt(3) / 2)
    parallel(g.major as Vec2, [1, 0])
    parallel(g.minor as Vec2, [0, 1])
  })

  it('classifies a rotated ellipse and names the minor axis', () => {
    const g = classifyGeneralConic(1, 1, 1, 0, 0, -1)
    expect(g.kind).toBe('ellipse')
    point(g.center, [0, 0])
    close(g.a, Math.SQRT2)
    close(g.b, Math.sqrt(2 / 3))
    close(g.eccentricity as number, Math.sqrt(2 / 3))
    parallel(g.minor as Vec2, [1 / Math.SQRT2, 1 / Math.SQRT2])
    parallel(g.major as Vec2, [1 / Math.SQRT2, -1 / Math.SQRT2])
    close(g.theta as number, Math.PI / 4)
  })

  it('classifies a hyperbola with the right transverse axis', () => {
    const g = classifyGeneralConic(1, 0, -1, 0, 0, -1)
    expect(g.kind).toBe('hyperbola')
    point(g.center, [0, 0])
    close(g.a, 1)
    close(g.b, 1)
    close(g.eccentricity as number, Math.SQRT2)
    parallel(g.transverse as Vec2, [1, 0])
    parallel(g.conjugate as Vec2, [0, 1])
  })

  it('classifies a rotated hyperbola and keeps the transverse direction', () => {
    const g = classifyGeneralConic(1, 4, 1, 0, 0, -1)
    expect(g.kind).toBe('hyperbola')
    point(g.center, [0, 0])
    close(g.a, 1 / Math.sqrt(3))
    close(g.b, 1)
    close(g.eccentricity as number, 2)
    parallel(g.transverse as Vec2, [1 / Math.SQRT2, 1 / Math.SQRT2])
    parallel(g.conjugate as Vec2, [1 / Math.SQRT2, -1 / Math.SQRT2])
  })

  it('classifies a vertical parabola with its vertex and opening', () => {
    const g = classifyGeneralConic(1, 0, 0, 0, -1, 0)
    expect(g.kind).toBe('parabola')
    point(g.vertex, [0, 0])
    close(g.p, 0.25)
    parallel(g.opening as Vec2, [0, 1])
    parallel(g.squared as Vec2, [1, 0])
    expect(g.eccentricity).toBe(1)
  })

  it('classifies a horizontal parabola', () => {
    const g = classifyGeneralConic(0, 0, 1, -1, 0, 0)
    expect(g.kind).toBe('parabola')
    point(g.vertex, [0, 0])
    close(g.p, 0.25)
    parallel(g.opening as Vec2, [1, 0])
    parallel(g.squared as Vec2, [0, 1])
  })

  it('classifies a translated parabola', () => {
    const g = classifyGeneralConic(1, 0, 0, -2, -1, 1)
    expect(g.kind).toBe('parabola')
    point(g.vertex, [1, 0])
    close(g.p, 0.25)
    parallel(g.opening as Vec2, [0, 1])
  })

  it('classifies a rotated parabola even with two written squared terms', () => {
    const g = classifyGeneralConic(1, 2, 1, 1, 0, 0)
    expect(g.kind).toBe('parabola')
    point(g.vertex, [-1 / 16, -3 / 16])
    close(g.p, Math.SQRT2 / 16, 1e-9)
    parallel(g.opening as Vec2, [-1 / Math.SQRT2, 1 / Math.SQRT2])
    parallel(g.squared as Vec2, [1 / Math.SQRT2, 1 / Math.SQRT2])
  })

  it('classifies intersecting lines, not a hyperbola', () => {
    const g = classifyGeneralConic(1, 0, -1, 0, 0, 0)
    expect(g.kind).toBe('intersecting')
    expect(g.kind).not.toBe('hyperbola')
    point(g.center, [0, 0])
    expect(g.lines).toHaveLength(2)
    parallel(g.lines[0].direction, [1 / Math.SQRT2, 1 / Math.SQRT2])
    parallel(g.lines[1].direction, [1 / Math.SQRT2, -1 / Math.SQRT2])
  })

  it('classifies parallel lines, not a parabola', () => {
    const g = classifyGeneralConic(1, 0, 0, 0, 0, -1)
    expect(g.kind).toBe('parallel')
    expect(g.kind).not.toBe('parabola')
    expect(g.lines).toHaveLength(2)
    for (const item of g.lines) parallel(item.direction, [0, 1])
    const xs = g.lines.map((item) => item.point[0]).sort((a, b) => a - b)
    near(xs[0], -1)
    near(xs[1], 1)
  })

  it('classifies a single line of rank one', () => {
    const g = classifyGeneralConic(1, 0, 0, 0, 0, 0)
    expect(g.kind).toBe('line')
    expect(g.lines).toHaveLength(1)
    parallel(g.lines[0].direction, [0, 1])
    near(g.lines[0].point[0], 0)
  })

  it('classifies an empty rank-one equation', () => {
    const g = classifyGeneralConic(1, 0, 0, 0, 0, 1)
    expect(g.kind).toBe('empty')
    expect(g.lines).toHaveLength(0)
  })

  it('classifies a point', () => {
    const g = classifyGeneralConic(1, 0, 1, 0, 0, 0)
    expect(g.kind).toBe('point')
    point(g.center, [0, 0])
  })

  it('classifies an empty rank-two equation', () => {
    const g = classifyGeneralConic(1, 0, 1, 0, 0, 1)
    expect(g.kind).toBe('empty')
  })

  it('classifies a rank-zero line', () => {
    const g = classifyGeneralConic(0, 0, 0, 2, 3, -6)
    expect(g.kind).toBe('line')
    expect(g.lines).toHaveLength(1)
    const { point: p, direction } = g.lines[0]
    near(2 * p[0] + 3 * p[1], 6)
    parallel(direction, [-3 / Math.sqrt(13), 2 / Math.sqrt(13)])
  })

  it('classifies an empty rank-zero equation', () => {
    const g = classifyGeneralConic(0, 0, 0, 0, 0, 1)
    expect(g.kind).toBe('empty')
  })

  it('classifies the whole plane', () => {
    const g = classifyGeneralConic(0, 0, 0, 0, 0, 0)
    expect(g.kind).toBe('all-plane')
    expect(g.lines).toHaveLength(0)
    expect(g.note).toBe('every point satisfies the equation')
  })
})

describe('degenerate counterexamples for each discriminant sign', () => {
  it('negative discriminant covers an ellipse, a point and an empty set', () => {
    for (const [A, B, C, F] of [
      [1, 0, 1, -1],
      [1, 0, 1, 0],
      [1, 0, 1, 1],
    ]) {
      expect(discriminant(A, B, C)).toBeLessThan(0)
    }
    expect(['ellipse', 'circle']).toContain(classifyGeneralConic(1, 0, 1, 0, 0, -1).kind)
    expect(classifyGeneralConic(1, 0, 1, 0, 0, 0).kind).toBe('point')
    expect(classifyGeneralConic(1, 0, 1, 0, 0, 1).kind).toBe('empty')
  })

  it('zero discriminant covers a parabola and three degenerate curves', () => {
    for (const [A, B, C] of [
      [1, 0, 0],
      [1, 2, 1],
      [0, 0, 1],
    ]) {
      near(discriminant(A, B, C), 0, 1e-12)
    }
    expect(classifyGeneralConic(1, 0, 0, 0, -1, 0).kind).toBe('parabola')
    expect(classifyGeneralConic(1, 0, 0, 0, 0, -1).kind).toBe('parallel')
    expect(classifyGeneralConic(1, 0, 0, 0, 0, 0).kind).toBe('line')
    expect(classifyGeneralConic(1, 0, 0, 0, 0, 1).kind).toBe('empty')
  })

  it('positive discriminant covers a hyperbola and intersecting lines', () => {
    expect(discriminant(1, 0, -1)).toBeGreaterThan(0)
    expect(classifyGeneralConic(1, 0, -1, 0, 0, -1).kind).toBe('hyperbola')
    expect(classifyGeneralConic(1, 0, -1, 0, 0, 0).kind).toBe('intersecting')
  })
})

describe('the parabola convention', () => {
  it('reports p as the nonnegative vertex-to-focus distance', () => {
    const cases: Array<[number, number, number, number, number, number, number, Vec2, Vec2]> = [
      [1, 0, 0, 0, -1, 0, 0.25, [0, 0], [0, 1]],
      [1, 0, 0, 0, -4, 0, 1, [0, 0], [0, 1]],
      [1, 0, 0, 0, 4, 0, 1, [0, 0], [0, -1]],
      [0, 0, 1, -4, 0, 0, 1, [0, 0], [1, 0]],
    ]
    for (const [A, B, C, D, E, F, p, vertex, opening] of cases) {
      const g = classifyGeneralConic(A, B, C, D, E, F)
      expect(g.kind).toBe('parabola')
      point(g.vertex, vertex)
      close(g.p, p)
      parallel(g.opening as Vec2, opening)
    }
  })

  it('keeps the focus direction, not just the line of the axis', () => {
    const up = classifyGeneralConic(1, 0, 0, 0, -1, 0)
    const down = classifyGeneralConic(1, 0, 0, 0, 1, 0)
    near((up.opening as Vec2)[1], 1)
    near((down.opening as Vec2)[1], -1)
  })
})

describe('focal geometry comes from the classifier', () => {
  it('places both ellipse foci on the major axis, including translated and rotated cases', () => {
    const base = classifyGeneralConic(1, 0, 0.5, 0, 0, -1)
    expect(base.foci).toHaveLength(2)
    hasFocus(base, [0, 1])
    hasFocus(base, [0, -1])
    close(base.focalConstant, 2 * Math.SQRT2)

    const translated = classifyGeneralConic(1, 0, 4, -2, 16, 13)
    expect(translated.foci).toHaveLength(2)
    hasFocus(translated, [1 + Math.sqrt(3), -2])
    hasFocus(translated, [1 - Math.sqrt(3), -2])
    close(translated.focalConstant, 4)

    const rotated = classifyGeneralConic(1, 1, 1, 0, 0, -1)
    expect(rotated.foci).toHaveLength(2)
    hasFocus(rotated, [Math.sqrt(2 / 3), -Math.sqrt(2 / 3)])
    hasFocus(rotated, [-Math.sqrt(2 / 3), Math.sqrt(2 / 3)])
    close(rotated.focalConstant, 2 * Math.SQRT2)
  })

  it('places hyperbola foci on the transverse axis for both orientations', () => {
    const right = classifyGeneralConic(1, 0, -1, 0, 0, -1)
    expect(right.foci).toHaveLength(2)
    hasFocus(right, [Math.SQRT2, 0])
    hasFocus(right, [-Math.SQRT2, 0])
    close(right.focalConstant, 2)

    const vertical = classifyGeneralConic(1, 0, -1, 0, 0, 1)
    expect(vertical.foci).toHaveLength(2)
    hasFocus(vertical, [0, Math.SQRT2])
    hasFocus(vertical, [0, -Math.SQRT2])

    const rotated = classifyGeneralConic(1, 4, 1, 0, 0, -1)
    expect(rotated.foci).toHaveLength(2)
    hasFocus(rotated, [Math.sqrt(2 / 3), Math.sqrt(2 / 3)])
    hasFocus(rotated, [-Math.sqrt(2 / 3), -Math.sqrt(2 / 3)])
    close(rotated.focalConstant, 2 / Math.sqrt(3))
  })

  it('puts a circle focus at its centre and keeps the focal constant at the diameter', () => {
    const g = classifyGeneralConic(1, 0, 1, -2, 2, -2)
    expect(g.kind).toBe('circle')
    expect(g.foci).toHaveLength(1)
    point(g.foci?.[0], [1, -1])
    close(g.focalConstant, 4)
    expect(g.eccentricity).toBe(0)
  })

  it('builds the focus and directrix of a parabola in every opening direction', () => {
    const cases: Array<[number, number, number, number, number, number, Vec2, Vec2, Vec2, Vec2]> = [
      // A, B, C, D, E, F, focus, directrix point, directrix direction, opening
      [1, 0, 0, 0, -1, 0, [0, 0.25], [0, -0.25], [1, 0], [0, 1]],
      [1, 0, 0, 0, 1, 0, [0, -0.25], [0, 0.25], [1, 0], [0, -1]],
      [0, 0, 1, -1, 0, 0, [0.25, 0], [-0.25, 0], [0, 1], [1, 0]],
      [0, 0, 1, 1, 0, 0, [-0.25, 0], [0.25, 0], [0, 1], [-1, 0]],
    ]
    for (const [A, B, C, D, E, F, focus, directrixPoint, direction, opening] of cases) {
      const g = classifyGeneralConic(A, B, C, D, E, F)
      expect(g.kind).toBe('parabola')
      expect(g.foci).toHaveLength(1)
      point(g.foci?.[0], focus)
      point(g.directrix?.point, directrixPoint)
      parallel(g.directrix?.direction as Vec2, direction)
      parallel(g.opening as Vec2, opening)
    }
  })

  it('rotates the parabola focal geometry with the curve', () => {
    const g = classifyGeneralConic(1, 2, 1, 1, 0, 0)
    expect(g.kind).toBe('parabola')
    point(g.vertex, [-1 / 16, -3 / 16])
    point(g.foci?.[0], [-1 / 8, -1 / 8])
    point(g.directrix?.point, [0, -0.25])
    parallel(g.directrix?.direction as Vec2, [1 / Math.SQRT2, 1 / Math.SQRT2])
    close(g.p, Math.SQRT2 / 16)
  })

  it('omits focal geometry for degenerate and invalid conics', () => {
    const cases: Array<[number, number, number, number, number, number]> = [
      [1, 0, 1, 0, 0, 0],
      [1, 0, -1, 0, 0, 0],
      [1, 0, 0, 0, 0, -1],
      [1, 0, 0, 0, 0, 0],
      [1, 0, 1, 0, 0, 1],
      [0, 0, 0, 2, 3, -6],
      [0, 0, 0, 0, 0, 1],
      [0, 0, 0, 0, 0, 0],
      [Number.NaN, 0, 1, 0, 0, -1],
    ]
    for (const [A, B, C, D, E, F] of cases) {
      const g = classifyGeneralConic(A, B, C, D, E, F)
      expect(g.foci ?? []).toHaveLength(0)
      expect(g.directrix).toBeUndefined()
      expect(g.focalConstant).toBeUndefined()
      finiteFields(g)
    }
  })
})

describe('probe projection and focal invariants', () => {
  it('parameterises every non-degenerate kind and rejects degenerate kinds', () => {
    expect(conicParameterRange(classifyGeneralConic(1, 0, 0.5, 0, 0, -1))).toEqual({
      min: 0,
      max: 2 * Math.PI,
      step: Math.PI / 90,
    })
    expect(conicParameterRange(classifyGeneralConic(1, 0, -1, 0, 0, -1))).toEqual({
      min: -4,
      max: 4,
      step: 0.04,
    })
    const parabola = conicParameterRange(classifyGeneralConic(1, 0, 0, 0, -1, 0))
    expect(parabola).not.toBeNull()
    expect((parabola as { min: number }).min).toBeLessThan(0)
    expect((parabola as { max: number }).max).toBeGreaterThan(0)
    expect(conicParameterRange(classifyGeneralConic(1, 0, -1, 0, 0, 0))).toBeNull()
  })

  it('projects a circle pointer to the circle and returns a genuine curve point', () => {
    const g = classifyGeneralConic(1, 0, 1, 0, 0, -1)
    const projection = projectToConic(g, [3, 4])
    expect(projection).not.toBeNull()
    point(projection?.point, [0.6, 0.8])
    close(projection?.distance, 4)
    onConic(g, projection?.point ?? null)
  })

  it('projects points onto translated and rotated ellipses', () => {
    for (const g of [classifyGeneralConic(1, 0, 4, -2, 16, 13), classifyGeneralConic(1, 1, 1, 0, 0, -1)]) {
      expect(g.kind).toBe('ellipse')
      const projection = projectToConic(g, [3, 2])
      expect(projection).not.toBeNull()
      onConic(g, projection?.point ?? null)
      const exact = conicPointAt(g, projection?.parameter as number)
      point(projection?.point, exact as Vec2)
    }
  })

  it('projects onto either hyperbola branch and keeps the chosen branch', () => {
    const g = classifyGeneralConic(1, 4, 1, 0, 0, -1)
    const first = projectToConic(g, [2, 2], 1)
    const second = projectToConic(g, [-2, -2], -1)
    expect(first).not.toBeNull()
    expect(second).not.toBeNull()
    onConic(g, first?.point ?? null)
    onConic(g, second?.point ?? null)
    expect((first?.point as Vec2)[0] + (first?.point as Vec2)[1]).toBeGreaterThan(0)
    expect((second?.point as Vec2)[0] + (second?.point as Vec2)[1]).toBeLessThan(0)
    expect(first?.branch).toBe(1)
    expect(second?.branch).toBe(-1)
  })

  it('projects onto a rotated parabola', () => {
    const g = classifyGeneralConic(1, 2, 1, 1, 0, 0)
    const projection = projectToConic(g, [0.4, 0.6])
    expect(projection).not.toBeNull()
    onConic(g, projection?.point ?? null)
  })

  it('keeps the ellipse focal sum equal to 2a at many probe positions', () => {
    for (const g of [classifyGeneralConic(1, 0, 4, -2, 16, 13), classifyGeneralConic(1, 1, 1, 0, 0, -1)]) {
      const expected = 2 * (g.a as number)
      for (const parameter of [0, 0.6, 1.7, 3.3, 5.1]) {
        const p = conicPointAt(g, parameter)
        onConic(g, p)
        const m = focalMeasurements(g, p as Vec2)
        close(m.focalSum ?? undefined, expected, 1e-8)
      }
    }
  })

  it('keeps the hyperbola focal difference equal to 2a on both branches', () => {
    for (const g of [classifyGeneralConic(1, 0, -1, 0, 0, -1), classifyGeneralConic(1, 4, 1, 0, 0, -1)]) {
      const expected = 2 * (g.a as number)
      for (const branch of [1, -1] as const) {
        for (const parameter of [-2, -0.4, 0, 0.6, 2]) {
          const p = conicPointAt(g, parameter, branch)
          onConic(g, p)
          const m = focalMeasurements(g, p as Vec2)
          close(m.focalDifference ?? undefined, expected, 1e-8)
        }
      }
    }
  })

  it('keeps focus distance equal to directrix distance for every parabola direction', () => {
    const parabolas = [
      classifyGeneralConic(1, 0, 0, 0, -1, 0),
      classifyGeneralConic(1, 0, 0, 0, 1, 0),
      classifyGeneralConic(0, 0, 1, -1, 0, 0),
      classifyGeneralConic(0, 0, 1, 1, 0, 0),
      classifyGeneralConic(1, 2, 1, 1, 0, 0),
    ]
    for (const g of parabolas) {
      expect(g.kind).toBe('parabola')
      for (const parameter of [-2, -0.5, 0, 0.7, 1.8]) {
        const p = conicPointAt(g, parameter)
        onConic(g, p)
        const m = focalMeasurements(g, p as Vec2)
        close(m.focusDistance ?? undefined, m.directrixDistance as number, 1e-8)
      }
    }
  })

  it('keeps eccentricity independent of the probe position', () => {
    for (const g of [
      classifyGeneralConic(1, 0, 0.5, 0, 0, -1),
      classifyGeneralConic(1, 0, -1, 0, 0, -1),
      classifyGeneralConic(1, 0, 0, 0, -1, 0),
    ]) {
      const eccentricity = g.eccentricity
      for (const parameter of [0, 0.8, 1.9, 3.1]) {
        const p = conicPointAt(g, parameter)
        expect(p).not.toBeNull()
        expect(g.eccentricity).toBe(eccentricity)
      }
    }
  })
})

describe('eccentricity is a constant of the conic', () => {
  it('matches the canonical formula and never depends on a point', () => {
    const ellipse = classifyGeneralConic(1, 0, 0.5, 0, 0, -1)
    const hyperbola = classifyGeneralConic(1, 0, -1, 0, 0, -1)
    const parabola = classifyGeneralConic(1, 0, 0, 0, -1, 0)
    close(ellipse.eccentricity as number, Math.sqrt(1 - (ellipse.b as number) ** 2 / (ellipse.a as number) ** 2))
    close(
      hyperbola.eccentricity as number,
      Math.sqrt(1 + (hyperbola.b as number) ** 2 / (hyperbola.a as number) ** 2),
    )
    expect(parabola.eccentricity).toBe(1)
    expect(classifyGeneralConic(1, 0, 1, 0, 0, -1).eccentricity).toBe(0)
    expect(classifyGeneralConic(1, 0, -1, 0, 0, 0).eccentricity).toBeNull()
  })
})

describe('scale invariance', () => {
  const cases: Array<[number, number, number, number, number, number]> = [
    [1, 0, 0.5, 0, 0, -1],
    [1, 0, 1, -2, 2, -2],
    [1, 0, 4, -2, 16, 13],
    [1, 1, 1, 0, 0, -1],
    [1, 0, -1, 0, 0, -1],
    [1, 4, 1, 0, 0, -1],
    [1, 0, 0, 0, -1, 0],
    [0, 0, 1, -1, 0, 0],
    [1, 2, 1, 1, 0, 0],
    [1, 0, -1, 0, 0, 0],
    [1, 0, 0, 0, 0, -1],
    [1, 0, 0, 0, 0, 0],
    [1, 0, 0, 0, 0, 1],
    [1, 0, 1, 0, 0, 0],
    [1, 0, 1, 0, 0, 1],
    [0, 0, 0, 2, 3, -6],
    [0, 0, 0, 0, 0, 1],
    [0, 0, 0, 0, 0, 0],
  ]

  for (const scale of [-1, 1e-6, 1e6]) {
    it(`keeps every kind and length under a factor of ${scale}`, () => {
      for (const [A, B, C, D, E, F] of cases) {
        const base = classifyGeneralConic(A, B, C, D, E, F)
        const scaled = classifyGeneralConic(A * scale, B * scale, C * scale, D * scale, E * scale, F * scale)
        expect(scaled.kind).toBe(base.kind)
        expect(scaled.lines.length).toBe(base.lines.length)
        for (const key of ['a', 'b', 'p', 'eccentricity'] as const) {
          const left = base[key]
          const right = scaled[key]
          if (left === null || left === undefined) expect(right ?? null).toBeNull()
          else close(right as number, left, 1e-6)
        }
        if (base.center) point(scaled.center, base.center, 1e-6)
        if (base.vertex) point(scaled.vertex, base.vertex, 1e-6)
        for (const key of ['major', 'minor', 'transverse', 'conjugate', 'squared', 'opening'] as const) {
          const left = base[key]
          const right = scaled[key]
          if (left) parallel(right as Vec2, left, 1e-6)
        }
        finiteFields(scaled)
      }
    })
  }
})

describe('numerical edges', () => {
  it('keeps a near-degenerate hyperbola a hyperbola', () => {
    const g = classifyGeneralConic(1, 0, -1, 0, 0, -1e-8)
    expect(g.kind).toBe('hyperbola')
    expect(g.certainty).toBe('near-degenerate')
    expect(g.note).toMatch(/near a degenerate case/)
  })

  it('reports the constant below resolution without changing the primitive', () => {
    const g = classifyGeneralConic(1, 0, -1, 0, 0, -1e-12)
    expect(g.kind).toBe('intersecting')
    expect(g.certainty).toBe('near-degenerate')
    expect(g.lines).toHaveLength(2)
  })

  it('marks a nearly rank-one ellipse as near-degenerate', () => {
    const g = classifyGeneralConic(1, 2 - 1e-7, 1, 0, 0, -1)
    expect(g.kind).toBe('ellipse')
    expect(g.certainty).toBe('near-degenerate')
    expect(g.a as number).toBeGreaterThan(1000)
  })

  it('places B at the balance point on the parallel-line boundary', () => {
    const at = classifyGeneralConic(1, Math.SQRT2, 0.5, 0, 0, -1)
    const below = classifyGeneralConic(1, 0, 0.5, 0, 0, -1)
    const above = classifyGeneralConic(1, 2, 0.5, 0, 0, -1)
    expect(at.kind).toBe('parallel')
    expect(below.kind).toBe('ellipse')
    expect(above.kind).toBe('hyperbola')
  })

  it('classifies a very small quadratic part as a parabola, not as a line', () => {
    const g = classifyGeneralConic(1e-12, 0, 0, 0, 1e12, 0)
    expect(g.kind).toBe('parabola')
    close(g.p, 2.5e23, 1e4)
    point(g.vertex, [0, 0])
    parallel(g.opening as Vec2, [0, -1])
    finiteFields(g)
  })

  it('returns finite geometry for a small but genuine circle', () => {
    const g = classifyGeneralConic(1e-8, 0, 1e-8, 0, 0, -1e-8)
    expect(g.kind).toBe('circle')
    close(g.a, 1)
    close(g.b, 1)
    finiteFields(g)
  })

  it('returns invalid for non-finite coefficients and never NaN geometry', () => {
    for (const bad of [Number.NaN, Number.POSITIVE_INFINITY, Number.NEGATIVE_INFINITY]) {
      const g = classifyGeneralConic(bad, 0, 1, 0, 0, -1)
      expect(g.kind).toBe('invalid')
      expect(g.certainty).toBe('invalid')
      expect(g.eccentricity).toBeNull()
      expect(g.note).toBe('coefficients must be finite')
      finiteFields(g)
    }
  })

  it('is deterministic', () => {
    const first = classifyGeneralConic(1, 2, 3, -4, 5, -6)
    const second = classifyGeneralConic(1, 2, 3, -4, 5, -6)
    expect(JSON.stringify(second)).toBe(JSON.stringify(first))
  })

  it('holds seeded random properties', () => {
    let state = 0x9e3779b9
    const random = () => {
      state = (Math.imul(state, 1664525) + 1013904223) >>> 0
      return state / 4294967296
    }
    const signed = () => (random() < 0.5 ? -1 : 1)
    const magnitude = () => 10 ** (-6 + 12 * random())

    for (let index = 0; index < 200; index++) {
      const coefficients: [number, number, number, number, number, number] = [
        signed() * magnitude(),
        random() < 0.4 ? 0 : signed() * magnitude(),
        signed() * magnitude(),
        random() < 0.4 ? 0 : signed() * magnitude(),
        random() < 0.4 ? 0 : signed() * magnitude(),
        signed() * magnitude(),
      ]
      const g = classifyGeneralConic(...coefficients)
      expect(kinds).toContain(g.kind)
      expect(g.kind).not.toBe('invalid')
      finiteFields(g)
      for (const item of g.lines) finiteFields({ ...g, lines: [], center: item.point, major: item.direction })
      const scaled = classifyGeneralConic(...(coefficients.map((value) => value * -3) as typeof coefficients))
      expect(scaled.kind).toBe(g.kind)
      if (g.kind === 'ellipse' || g.kind === 'circle' || g.kind === 'hyperbola') {
        relative(scaled.a, g.a as number)
        relative(scaled.b ?? 0, g.b ?? 0)
      }
      if (g.kind === 'parabola') {
        relative(scaled.p, g.p as number)
        if (g.vertex) point(scaled.vertex, g.vertex, 1e-6)
      }
    }
  })
})

describe('intercepts', () => {
  it('reports two, one, double, none and coincident cases', () => {
    const ellipse = classifyGeneralConic(1, 0, 0.5, 0, 0, -1)
    expect(ellipse.intercepts.x).toEqual({ status: 'two', values: [-1, 1] })
    expect(ellipse.intercepts.y.status).toBe('two')
    near(ellipse.intercepts.y.values[1], Math.SQRT2)

    const parabola = classifyGeneralConic(1, 0, 0, 0, -1, 0)
    expect(parabola.intercepts.x.status).toBe('double')
    expect(parabola.intercepts.x.values).toHaveLength(1)
    near(parabola.intercepts.x.values[0], 0)

    const horizontal = classifyGeneralConic(0, 0, 0, 0, 1, 0)
    expect(horizontal.intercepts.y).toEqual({ status: 'one', values: [0] })
    expect(horizontal.intercepts.x.status).toBe('coincident')

    const empty = classifyGeneralConic(1, 0, 1, 0, 0, 1)
    expect(empty.intercepts.x.status).toBe('none')
    expect(empty.intercepts.y.status).toBe('none')
  })
})

describe('the equation and its clipping', () => {
  it('evaluates the general equation', () => {
    near(generalConic(1, 2, 3, 4, 5, 6, 1, 2), 1 + 4 + 12 + 4 + 10 + 6)
    near(generalConic(1, 0, 0.5, 0, 0, -1, Math.SQRT2, 0), 2 - 1)
  })

  it('clips a line to a rectangle or reports that it misses', () => {
    const bounds = { x0: -2, x1: 2, y0: -2, y1: 2 }
    expect(clipLineToRect([0, 0], [1, 1], bounds)).toEqual([
      [-2, -2],
      [2, 2],
    ])
    expect(clipLineToRect([0, 0], [0, 1], bounds)).toEqual([
      [0, -2],
      [0, 2],
    ])
    expect(clipLineToRect([3, 0], [0, 1], bounds)).toBeNull()
    expect(clipLineToRect([0, 5], [1, 1], bounds)).toBeNull()
    expect(clipLineToRect([0, 0], [0, 0], bounds)).toBeNull()
  })
})
