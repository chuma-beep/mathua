import { describe, expect, it } from 'vitest'
import {
  circulation,
  flux,
  loopArea,
  loopPath,
  streamline,
  vfCurl,
  vfDivergence,
  vfValue,
} from '../components/visual/math/vectorfield'
import { project } from '../components/visual/math/project3d'
import {
  cellStats,
  exactIntegral,
  regionCells,
  regionSum,
  reversedBounds,
} from '../components/visual/math/region'
import { innerRadius, outerRadius, sliceArea, sliceVolume, solidVolume } from '../components/visual/math/solidrev'
import { accumulationModel, areaCurve, integrand, primitive } from '../components/visual/math/accumulation'
import { directional, fieldCritical, fieldHessian, saddleTest, tangentPlane } from '../components/visual/math/field'
import { partialSum, ratioTerm, seriesLimit, sumStatus, sumTerm } from '../components/visual/math/series'
import { BIPARTITE_GRAPH, COLOUR_GRAPH, bipartiteCheck, greedyColour } from '../components/visual/math/graph'
import { conditionalArea } from '../components/visual/math/probability'

const close = (a: number, b: number, tolerance = 1e-6) =>
  expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('vector fields', () => {
  it('uses the field itself for the loop integrals', () => {
    close(circulation('vortex', 1.2), 2 * Math.PI * 1.44)
    close(circulation('shear', 1.2), -Math.PI * 1.44)
    close(circulation('source', 1.2), 0)
    close(flux('source', 1.2), 2 * Math.PI * 1.44)
    close(flux('vortex', 1.2), 0)
  })

  it('reports the exact divergence and curl of each preset', () => {
    expect(vfDivergence('source')).toBe(2)
    expect(vfDivergence('saddle')).toBe(0)
    expect(vfCurl('vortex')).toBe(2)
    expect(vfCurl('shear')).toBe(-1)
    const [vx, vy] = vfValue('vortex', 1, 0)
    close(vx, 0)
    close(vy, 1)
    close(vfCurl('vortex') * loopArea(1.2), circulation('vortex', 1.2))
    close(vfDivergence('source') * loopArea(1.2), flux('source', 1.2))
  })

  it('closes the loop and integrates a streamline', () => {
    const path = loopPath(1, 40)
    expect(path[0]).toEqual(path[path.length - 1])
    const line = streamline('vortex', [1, 0], 2, 80)
    expect(line.length).toBeGreaterThan(2)
    const radius = Math.hypot(line[line.length - 1][0], line[line.length - 1][1])
    close(radius, 1, 1e-3)
  })
})

describe('projection', () => {
  it('keeps depth monotone and rotates consistently', () => {
    const near = project({ x: 0, y: 0, z: 1 }, { theta: 0, phi: 0.5 })
    const far = project({ x: 0, y: 0, z: -1 }, { theta: 0, phi: 0.5 })
    expect(near.depth).toBeGreaterThan(far.depth)
    const turned = project({ x: 1, y: 0, z: 0 }, { theta: Math.PI / 2, phi: 0 })
    close(turned.x, 0, 1e-9)
  })
})

describe('region integrals', () => {
  it('matches the closed forms on each region', () => {
    expect(exactIntegral('rect', 'one')).toBe(2)
    expect(exactIntegral('rect', 'x')).toBe(2)
    expect(exactIntegral('rect', 'y')).toBe(1)
    close(exactIntegral('between', 'one'), 1 / 6)
    close(exactIntegral('between', 'x'), 1 / 12)
    close(exactIntegral('between', 'y'), 1 / 15)
    close(exactIntegral('polar', 'one'), Math.PI / 4)
    close(exactIntegral('polar', 'x'), 1 / 3)
  })

  it('approaches the exact value as the grid refines', () => {
    const coarse = Math.abs(regionSum('between', 'one', 8) - 1 / 6)
    const fine = Math.abs(regionSum('between', 'one', 40) - 1 / 6)
    expect(fine).toBeLessThan(coarse)
    expect(fine).toBeLessThan(0.02)
    const stats = cellStats('between', 10)
    expect(stats.count).toBe(regionCells('between', 10).filter((cell) => cell.inside).length)
  })

  it('names the reversed bounds', () => {
    expect(reversedBounds('between')).toContain('√y')
    expect(reversedBounds('polar')).toContain('θ')
  })
})

describe('solids of revolution', () => {
  it('integrates each method from its own slice', () => {
    const disk = solidVolume('disk')
    const washer = solidVolume('washer')
    close(washer, disk - Math.PI * innerRadius(1) ** 2 * 2, 1e-6)
    expect(washer).toBeLessThan(disk)
    expect(solidVolume('shell')).toBeGreaterThan(washer)
    close(sliceVolume('disk', 1, 0.01), Math.PI * outerRadius(1) ** 2 * 0.01)
    close(sliceArea('washer', 1), Math.PI * (outerRadius(1) ** 2 - innerRadius(1) ** 2))
  })
})

describe('accumulation and the fundamental theorem', () => {
  it('differentiates the area back to the integrand', () => {
    for (const fn of ['square', 'sine', 'decay'] as const) {
      const model = accumulationModel(fn, 0.5, 2)
      expect(model.status).toBe('ok')
      close(model.area, primitive(fn, 2) - primitive(fn, 0.5))
      close(model.measured, integrand(fn, 2), 1e-6)
      expect(model.error).toBeLessThan(1e-6)
    }
  })

  it('produces a finite area curve', () => {
    const curve = areaCurve('sine', 0, 0, 3, 30)
    expect(curve).toHaveLength(31)
    expect(curve.every(([, value]) => Number.isFinite(value))).toBe(true)
    expect(accumulationModel('square', Number.NaN, 1).status).toBe('invalid')
  })
})

describe('tangent plane and saddle test', () => {
  it('matches value and slopes at the probe', () => {
    const plane = tangentPlane('bowl', 1, 0.6)
    close(plane.z0, 1.36)
    close(plane.at(1, 0.6), plane.z0)
    close(plane.at(1.3, 0.6), 1.96)
    close(plane.fx, 2)
    close(plane.fy, 1.2)
    expect(plane.equation).toContain('1.36')
  })

  it('classifies critical points from the Hessian', () => {
    expect(saddleTest('bowl', 0, 0).classification).toBe('minimum')
    expect(saddleTest('saddle', 0, 0).classification).toBe('saddle')
    expect(saddleTest('saddle', 0, 0).discriminant).toBe(-4)
    expect(fieldHessian('product', 0, 0).xy).toBe(1)
    expect(fieldCritical('ripple').length).toBeGreaterThan(0)
    close(directional('bowl', 1, 0, 0), 2)
    close(directional('bowl', 1, 0, Math.PI / 2), 0)
  })
})

describe('series sums', () => {
  it('matches the known limits and the ratio test', () => {
    close(sumTerm('geometric', 3), 0.125)
    close(partialSum('geometric', 40), 1, 1e-9)
    close(seriesLimit('p2'), (Math.PI * Math.PI) / 6)
    close(seriesLimit('alternating'), Math.LN2)
    expect(seriesLimit('harmonic')).toBe(Number.POSITIVE_INFINITY)
    expect(sumStatus('harmonic')).toBe('diverges')
    expect(sumStatus('p2')).toBe('converges')
    close(ratioTerm('geometric', 5), 0.5)
    expect(ratioTerm('harmonic', 4)).toBeLessThan(1)
  })
})

describe('graph bipartition and colouring', () => {
  it('splits the even cycle into two sides', () => {
    const result = bipartiteCheck(BIPARTITE_GRAPH)
    expect(result.status).toBe('ok')
    expect(result.sides?.left).toEqual(['A', 'C', 'E'])
    expect(result.sides?.right).toEqual(['B', 'D', 'F'])
  })

  it('needs three colours for the odd cycle with a chord', () => {
    const result = greedyColour(COLOUR_GRAPH)
    expect(result.colourCount).toBe(3)
    expect(result.colours?.A).toBe(0)
    expect(result.colours?.C).toBe(1)
    expect(result.colours?.B).toBe(2)
  })
})

describe('conditional area', () => {
  it('multiplies the column width by the band height', () => {
    const model = conditionalArea(0.4, 0.8, 0.3)
    expect(model).not.toBeNull()
    close(model!.pAandB, 0.32)
    close(model!.pAandBc, 0.18)
    close(model!.pA, 0.5)
    expect(model!.independent).toBe(false)
    expect(conditionalArea(0.4, 0.5, 0.5)!.independent).toBe(true)
    expect(conditionalArea(Number.NaN, 0.5, 0.5)).toBeNull()
  })
})
