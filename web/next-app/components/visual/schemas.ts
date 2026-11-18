import { z } from 'zod'

// Parameter schemas, one per template. They are the single definition of what a
// shard may send: the registry validates with the same object the template parses,
// so the two cannot disagree about a default or a bound.

const pair = z.tuple([z.number().finite(), z.number().finite()])
const coeffs = z.array(z.number().finite()).min(1).max(8)

export const matrix2dSchema = z.object({
  matrix: z
    .tuple([z.number().finite(), z.number().finite(), z.number().finite(), z.number().finite()])
    .default([1, 0.5, 0.5, 1]),
  vector: pair.default([2, 1]),
  showGrid: z.boolean().default(true),
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
