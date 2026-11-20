import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import {
  accumulationSchema,
  conditionalAreaSchema,
  graphSchema,
  regionSchema,
  seriesSumSchema,
  solidRevSchema,
  surfaceSchema,
  vectorFieldSchema,
} from '../components/visual/schemas'
import { circulation, vfCurl } from '../components/visual/math/vectorfield'
import { exactIntegral, regionSum } from '../components/visual/math/region'
import { solidVolume } from '../components/visual/math/solidrev'
import { accumulationModel } from '../components/visual/math/accumulation'
import { saddleTest } from '../components/visual/math/field'
import { partialSum, seriesLimit } from '../components/visual/math/series'
import { BIPARTITE_GRAPH, COLOUR_GRAPH, bipartiteCheck, greedyColour } from '../components/visual/math/graph'
import { conditionalArea } from '../components/visual/math/probability'

const REPO_ROOT = join(__dirname, '../../..')

interface Step {
  label?: string
  visual?: {
    type: string
    title?: string
    caption?: string
    params?: Record<string, unknown>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const shard = (id: string): Step[] =>
  JSON.parse(readFileSync(join(REPO_ROOT, 'data', 'lessons', 'kp', `${id}.json`), 'utf8'))

const visualOf = (id: string) => {
  const steps = shard(id)
  expect(steps[0].visual, `${id} has no visual on step 0`).toBeTruthy()
  expect(steps.filter((step) => step.visual)).toHaveLength(1)
  return steps[0].visual!
}

const CASES: Array<{ id: string; type: string; answer: string }> = [
  { id: 'calc.vec.vector_field', type: 'calc.vectorfield', answer: 'The vector assigned to that point' },
  { id: 'calc.vec.divergence', type: 'calc.vectorfield', answer: 'Outflow per unit area' },
  { id: 'calc.vec.curl', type: 'calc.vectorfield', answer: '−1' },
  { id: 'calc.vec.curl_3d', type: 'calc.vectorfield', answer: 'Each component is a planar curl in the plane normal to it' },
  { id: 'calc.vec.line_integral', type: 'calc.vectorfield', answer: 'It changes its sign' },
  { id: 'calc.vec.greens_theorem', type: 'calc.vectorfield', answer: 'The double integral of the curl over the interior' },
  { id: 'calc.vec.stokes', type: 'calc.vectorfield', answer: 'The surface is flat and lies in a plane' },
  { id: 'calc.vec.divergence_theorem', type: 'calc.vectorfield', answer: 'The triple integral of the divergence inside' },
  { id: 'calc.vec.vector_potential', type: 'calc.vectorfield', answer: 'F = curl A' },
  { id: 'calc.deriv.partial', type: 'calc.surface', answer: '2' },
  { id: 'calc.multi.tangent_plane', type: 'calc.surface', answer: 'The two partial derivatives' },
  { id: 'calc.multi.taylor2', type: 'calc.surface', answer: 'The Hessian' },
  { id: 'calc.multi.saddle_test', type: 'calc.surface', answer: '−4, a saddle' },
  { id: 'calc.multi.implicit_function', type: 'calc.surface', answer: 'F_y ≠ 0' },
  { id: 'calc.multi.double_integral_rect', type: 'calc.region', answer: 'x' },
  { id: 'calc.multi.double_integral_region', type: 'calc.region', answer: '1/6' },
  { id: 'calc.multi.change_order', type: 'calc.region', answer: 'y ≤ x ≤ √y' },
  { id: 'calc.multi.triple_integral_box', type: 'calc.region', answer: 'Innermost variable first' },
  { id: 'calc.multi.coordinates_jacobian', type: 'calc.region', answer: 'r dr dθ' },
  { id: 'calc.integral.area_between', type: 'calc.region', answer: '1/6' },
  { id: 'calc.series.concept', type: 'calc.seriessum', answer: 'The limit of its partial sums' },
  { id: 'calc.series.geometric', type: 'calc.seriessum', answer: '1' },
  { id: 'calc.series.harmonic', type: 'calc.seriessum', answer: 'Diverges' },
  { id: 'calc.series.p_series', type: 'calc.seriessum', answer: 'p > 1' },
  { id: 'calc.series.ratio_test', type: 'calc.seriessum', answer: 'L < 1' },
  { id: 'calc.series.estimate', type: 'calc.seriessum', answer: 'The next term' },
  { id: 'calc.integral.accumulation', type: 'calc.accumulation', answer: 'f(x), the integrand' },
  { id: 'calc.integral.ftc', type: 'calc.accumulation', answer: 'Area and antiderivatives' },
  { id: 'calc.integral.net_change', type: 'calc.accumulation', answer: 'f(b) − f(a)' },
  { id: 'calc.integral.variable_limits', type: 'calc.accumulation', answer: 'f(x)' },
  { id: 'calc.integral.disc', type: 'calc.solidrev', answer: 'πR² dx' },
  { id: 'calc.integral.washer', type: 'calc.solidrev', answer: 'π(R² − r²)' },
  { id: 'calc.integral.shell', type: 'calc.solidrev', answer: '2π · radius · height · thickness' },
  { id: 'calc.integral.volume', type: 'calc.solidrev', answer: 'The distance from the axis of rotation' },
  { id: 'calc.integral.surface_area', type: 'calc.solidrev', answer: '2π f(x) √(1 + f′(x)²) dx' },
  { id: 'graph.grundy', type: 'graph.network', answer: '3' },
  { id: 'graph.konig', type: 'graph.network', answer: 'It has no odd cycle' },
  { id: 'prob.cond.multiplication', type: 'prob.area', answer: 'P(B) · P(A|B)' },
  { id: 'prob.independence.pairwise', type: 'prob.area', answer: 'A and B are independent' },
  { id: 'stat.infer.sampling', type: 'stat.sampling', answer: 'Narrows around the population mean' },
]

describe('batch 6 shard attachments', () => {
  for (const item of CASES) {
    it(`${item.id} carries ${item.type} and an agreeing prediction`, () => {
      const visual = visualOf(item.id)
      expect(visual.type).toBe(item.type)
      expect(visual.predict?.options[visual.predict.answer]).toBe(item.answer)
      expect(visual.predict?.explain.length).toBeGreaterThan(10)
      expect(visual.title?.length ?? 0).toBeLessThanOrEqual(160)
      expect(visual.caption?.length ?? 0).toBeLessThanOrEqual(600)
    })
  }
})

describe('batch 6 numbers match the engine', () => {
  it('vector field loop integrals agree with the drawn field', () => {
    const parsed = vectorFieldSchema.parse(visualOf('calc.vec.greens_theorem').params)
    expect(parsed.preset).toBe('vortex')
    expect(circulation(parsed.preset, parsed.radius)).toBeCloseTo(
      vfCurl(parsed.preset) * Math.PI * parsed.radius ** 2,
      9,
    )
  })

  it('region estimate approaches the exact shard value', () => {
    const parsed = regionSchema.parse(visualOf('calc.multi.double_integral_region').params)
    expect(parsed.region).toBe('between')
    expect(regionSum(parsed.region, parsed.density, 200)).toBeCloseTo(exactIntegral(parsed.region, parsed.density), 2)
  })

  it('the area-between lesson uses the same region', () => {
    const parsed = regionSchema.parse(visualOf('calc.integral.area_between').params)
    expect(exactIntegral(parsed.region, parsed.density)).toBeCloseTo(1 / 6, 9)
  })

  it('the washer volume is the integral of its slice', () => {
    const parsed = solidRevSchema.parse(visualOf('calc.integral.washer').params)
    expect(parsed.mode).toBe('washer')
    expect(solidVolume(parsed.mode)).toBeCloseTo(8.1262, 3)
  })

  it('the accumulation readout differentiates to f(x)', () => {
    const parsed = accumulationSchema.parse(visualOf('calc.integral.ftc').params)
    const model = accumulationModel(parsed.fn, parsed.a, parsed.x)
    expect(model.error).toBeLessThan(1e-6)
  })

  it('the surface lesson lands on a true saddle', () => {
    const parsed = surfaceSchema.parse(visualOf('calc.multi.saddle_test').params)
    expect(saddleTest(parsed.field, parsed.x0, parsed.y0).classification).toBe('saddle')
  })

  it('series limits agree with the shard claims', () => {
    const geometric = seriesSumSchema.parse(visualOf('calc.series.geometric').params)
    expect(seriesLimit(geometric.kind)).toBe(1)
    expect(partialSum('geometric', 30)).toBeCloseTo(1, 6)
    const harmonic = seriesSumSchema.parse(visualOf('calc.series.harmonic').params)
    expect(seriesLimit(harmonic.kind)).toBe(Number.POSITIVE_INFINITY)
  })

  it('the colouring count agrees with the shard answer', () => {
    const parsed = graphSchema.parse(visualOf('graph.grundy').params)
    expect(parsed.algorithm).toBe('colouring')
    expect(greedyColour(COLOUR_GRAPH).colourCount).toBe(3)
    const bipartite = graphSchema.parse(visualOf('graph.konig').params)
    expect(bipartite.algorithm).toBe('bipartite')
    expect(bipartiteCheck(BIPARTITE_GRAPH).status).toBe('ok')
  })

  it('the conditional area multiplies to the stated overlap', () => {
    const parsed = conditionalAreaSchema.parse(visualOf('prob.cond.multiplication').params)
    const model = conditionalArea(parsed.pB, parsed.pAgivenB, parsed.pAgivenBc)
    expect(model?.pAandB).toBeCloseTo(parsed.pB * parsed.pAgivenB, 9)
  })
})
