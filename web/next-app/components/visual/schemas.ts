import { z } from 'zod'

// Parameter schemas, one per template. They are the single definition of what a
// shard may send: the registry validates with the same object the template parses,
// so the two cannot disagree about a default or a bound.

const pair = z.tuple([z.number().finite(), z.number().finite()])
const coeffs = z.array(z.number().finite()).min(1).max(8)

const quad = z.tuple([z.number().finite(), z.number().finite(), z.number().finite(), z.number().finite()])

export const matrix2dSchema = z.object({
  matrix: quad.default([1, 0.5, 0.5, 1]),
  // An optional second matrix, applied after the first, so a lesson can show
  // A, then B, against the composed BA. Absent means the figure is one matrix.
  then: quad.optional(),
  vector: pair.default([2, 1]),
  showGrid: z.boolean().default(true),
  showBasis: z.boolean().default(true),
})

export const tangentSchema = z.object({
  coeffs: coeffs.default([0, 0, 1]),
  x: z.number().finite().default(1),
  h: z.number().finite().min(0.05).max(3).default(0.75),
  window: z.tuple([z.number().finite(), z.number().finite()]).default([-3, 3]),
})

export const riemannSchema = z.object({
  coeffs: coeffs.default([0, 0, 1]),
  a: z.number().finite().default(0),
  b: z.number().finite().default(2),
  n: z.number().int().min(1).max(256).default(8),
  rule: z.enum(['left', 'right', 'mid']).default('left'),
})

export const complexMultiplySchema = z.object({
  z: pair.default([1.5, 0.8]),
  w: pair.default([0.8, 1.2]),
  showUnitCircle: z.boolean().default(true),
})

export const normalSchema = z.object({
  mu: z.number().finite().default(0),
  sigma: z.number().finite().min(0.1).max(4).default(1),
  x: z.number().finite().default(1),
  bands: z.boolean().default(true),
})

export const transform2dSchema = z.object({
  kind: z.enum(['rotation', 'reflection', 'translation', 'dilation', 'composition']).default('rotation'),
  angle: z.number().finite().default(30),
  scale: z.number().finite().min(0.25).max(2.5).default(1.25),
  dx: z.number().finite().default(1),
  dy: z.number().finite().default(0.5),
  shape: z.array(pair).min(3).max(8).default([
    [0, 0],
    [2, 0],
    [1, 2],
  ]),
  showOrbit: z.boolean().default(true),
})

export const eigenSchema = z.object({
  matrix: quad.default([2, 1, 1, 2]),
  vector: pair.default([1, 0]),
  showCircle: z.boolean().default(true),
})

export const taylorSchema = z.object({
  fn: z.enum(['sin', 'cos', 'exp', 'geometric']).default('sin'),
  center: z.number().finite().default(0),
  degree: z.number().int().min(0).max(14).default(3),
  window: z.tuple([z.number().finite(), z.number().finite()]).default([-3, 3]),
})

export const limitSchema = z.object({
  kind: z.enum(['removable', 'jump', 'vertical_asymptote']).default('removable'),
  a: z.number().finite().default(2),
  delta: z.number().finite().min(0.02).max(1.5).default(0.6),
})

export const samplingSchema = z.object({
  population: z.enum(['uniform', 'normal', 'skewed']).default('skewed'),
  n: z.number().int().min(1).max(64).default(9),
  draws: z.number().int().min(100).max(2000).default(600),
  seed: z.number().int().min(1).max(9999).default(7),
})

export const regressionSchema = z.object({
  points: z.array(pair).min(2).max(12).default([
    [1, 2],
    [2, 3],
    [3, 3.5],
    [4, 6],
    [5, 6.5],
  ]),
  showResiduals: z.boolean().default(true),
})
