import { describe, expect, it } from 'vitest'
import {
  TAU,
  degreesToRadians,
  normalizeAngle,
  normalizeSigned,
  pythagorean,
  projectToUnitCircle,
  quadrantOf,
  radiansToDegrees,
  referenceAngle,
  unitPoint,
} from '../components/visual/math/trig'

const close = (actual: number, expected: number, tolerance = 1e-12) =>
  expect(Math.abs(actual - expected)).toBeLessThan(tolerance)

const pointClose = (actual: [number, number], expected: [number, number], tolerance = 1e-12) => {
  close(actual[0], expected[0], tolerance)
  close(actual[1], expected[1], tolerance)
}

describe('unit-circle angles', () => {
  it('reduces positive and negative angles to one turn', () => {
    close(normalizeAngle(0), 0)
    close(normalizeAngle(Math.PI / 4), Math.PI / 4)
    close(normalizeAngle(-Math.PI / 4), (7 * Math.PI) / 4)
    close(normalizeAngle(TAU + 0.2), 0.2, 1e-12)
    close(normalizeAngle(-TAU - 0.2), TAU - 0.2, 1e-12)
    close(normalizeAngle(3 * Math.PI), Math.PI)
  })

  it('keeps non-finite input non-finite', () => {
    expect(Number.isNaN(normalizeAngle(Number.NaN))).toBe(true)
    expect(Number.isNaN(normalizeAngle(Number.POSITIVE_INFINITY))).toBe(true)
  })

  it('reduces to the signed branch used by a graph marker', () => {
    close(normalizeSigned((5 * Math.PI) / 4), (-3 * Math.PI) / 4)
    close(normalizeSigned((3 * Math.PI) / 2), -Math.PI / 2)
    close(normalizeSigned(Math.PI), Math.PI)
    close(normalizeSigned(-Math.PI / 3), -Math.PI / 3)
  })

  it('converts degrees and radians in both directions', () => {
    close(degreesToRadians(180), Math.PI)
    close(degreesToRadians(30), Math.PI / 6)
    close(radiansToDegrees(Math.PI / 2), 90)
    close(radiansToDegrees((7 * Math.PI) / 4), 315)
  })

  it('reports the acute reference angle in every quadrant', () => {
    close(referenceAngle(degreesToRadians(30)), degreesToRadians(30))
    close(referenceAngle(degreesToRadians(150)), degreesToRadians(30))
    close(referenceAngle(degreesToRadians(210)), degreesToRadians(30))
    close(referenceAngle(degreesToRadians(330)), degreesToRadians(30))
    close(referenceAngle(degreesToRadians(120)), degreesToRadians(60))
    close(referenceAngle(degreesToRadians(270)), degreesToRadians(90))
    close(referenceAngle(degreesToRadians(-45)), degreesToRadians(45))
  })

  it('names quadrants and distinguishes the axes', () => {
    expect(quadrantOf(degreesToRadians(45))).toBe('I')
    expect(quadrantOf(degreesToRadians(135))).toBe('II')
    expect(quadrantOf(degreesToRadians(225))).toBe('III')
    expect(quadrantOf(degreesToRadians(315))).toBe('IV')
    expect(quadrantOf(0)).toBe('axis')
    expect(quadrantOf(degreesToRadians(90))).toBe('axis')
    expect(quadrantOf(degreesToRadians(180))).toBe('axis')
    expect(quadrantOf(degreesToRadians(270))).toBe('axis')
    expect(quadrantOf(degreesToRadians(360))).toBe('axis')
  })
})

describe('unit-circle point values', () => {
  it('reads the four cardinal points exactly', () => {
    pointClose(unitPoint(0).point, [1, 0])
    close(unitPoint(degreesToRadians(90)).point[0], 0, 1e-15)
    close(unitPoint(degreesToRadians(90)).point[1], 1, 1e-15)
    pointClose(unitPoint(degreesToRadians(180)).point, [-1, 0])
    close(unitPoint(degreesToRadians(270)).point[0], 0, 1e-15)
    close(unitPoint(degreesToRadians(270)).point[1], -1, 1e-15)
  })

  it('reads the common 30-60-90 point', () => {
    const point = unitPoint(degreesToRadians(30))
    close(point.cos, Math.sqrt(3) / 2)
    close(point.sin, 0.5)
    close(point.tan ?? Number.NaN, 1 / Math.sqrt(3))
  })

  it('returns a null tangent on the vertical axis', () => {
    expect(unitPoint(degreesToRadians(90)).tan).toBeNull()
    expect(unitPoint(degreesToRadians(270)).tan).toBeNull()
    close(unitPoint(degreesToRadians(45)).tan ?? Number.NaN, 1)
    close(unitPoint(degreesToRadians(135)).tan ?? Number.NaN, -1)
  })

  it('keeps the angle convention counterclockwise', () => {
    const quarter = unitPoint(degreesToRadians(90))
    const threeQuarters = unitPoint(degreesToRadians(270))
    expect(quarter.cos).toBeGreaterThan(-1e-15)
    expect(threeQuarters.sin).toBeLessThan(-0.99)
    expect(quadrantOf(degreesToRadians(-45))).toBe('IV')
    expect(quadrantOf(degreesToRadians(-135))).toBe('III')
  })

  it('has the same point after a full turn', () => {
    for (const degrees of [0, 30, 45, 120, 210, 330]) {
      const base = unitPoint(degreesToRadians(degrees))
      const turned = unitPoint(degreesToRadians(degrees + 360))
      close(turned.cos, base.cos)
      close(turned.sin, base.sin)
      close(turned.theta, base.theta, 1e-12)
    }
  })

  it('satisfies the Pythagorean identity at sampled angles', () => {
    for (const degrees of [0, 17, 45, 89, 120, 180, 233, 300, 359]) {
      close(pythagorean(degreesToRadians(degrees)), 1, 1e-12)
    }
  })

  it('is deterministic for repeated calls', () => {
    const first = unitPoint(degreesToRadians(137))
    const second = unitPoint(degreesToRadians(137))
    expect(second).toEqual(first)
  })
})

describe('dragging projects to the unit circle', () => {
  it('uses the pointer direction and ignores its length', () => {
    close(projectToUnitCircle([3, 4]), Math.atan2(4, 3))
    const theta = projectToUnitCircle([3, 4])
    pointClose([Math.cos(theta), Math.sin(theta)], [0.6, 0.8])
    close(projectToUnitCircle([0.3, 0.4]), Math.atan2(0.4, 0.3))
  })

  it('handles the origin, non-finite input and negative directions', () => {
    expect(projectToUnitCircle([0, 0])).toBe(0)
    expect(projectToUnitCircle([Number.NaN, 1])).toBe(0)
    close(projectToUnitCircle([-1, 0]), Math.PI)
    close(projectToUnitCircle([0, -1]), (3 * Math.PI) / 2)
    close(projectToUnitCircle([-1, -1]), (5 * Math.PI) / 4)
  })
})
