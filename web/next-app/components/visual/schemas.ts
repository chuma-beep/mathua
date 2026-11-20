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

export const vectorCombinationSchema = z.object({
  a: pair.default([1, 2]),
  b: pair.default([2, 0.5]),
  k: z.number().finite().min(-3).max(3).default(-1.5),
  showSum: z.boolean().default(true),
})

export const basisChangeSchema = z.object({
  u: pair.default([1.5, 0.5]),
  v: pair.default([0.5, 1.5]),
  point: pair.default([2, 1]),
})

export const spreadSchema = z.object({
  values: z.array(z.number().finite().min(0).max(10)).min(2).max(8).default([1, 2, 2, 4, 6]),
  squared: z.boolean().default(true),
  sample: z.boolean().default(true),
})

// The general six-term equation. F defaults to -1 so the old homogeneous
// params {A, B, C} keep drawing the same curve: Ax² + Bxy + Cy² - 1 = 0.
export const conicSchema = z.object({
  A: z.number().finite().min(-2).max(2).default(1),
  B: z.number().finite().min(-2).max(2).default(0),
  C: z.number().finite().min(-2).max(2).default(0.5),
  D: z.number().finite().min(-4).max(4).default(0),
  E: z.number().finite().min(-4).max(4).default(0),
  F: z.number().finite().min(-4).max(4).default(-1),
})

// The unit circle is drawn with its angle in degrees, because the lesson's
// common angles are written that way. The graph converts once to radians.
export const unitCircleSchema = z.object({
  angle: z.number().finite().min(-360).max(360).default(45),
  showSin: z.boolean().default(true),
  showCos: z.boolean().default(true),
})

export const contourSchema = z.object({
  field: z.enum(['bowl', 'saddle', 'product', 'ripple']).default('bowl'),
  levels: z.number().int().min(3).max(12).default(7),
  probe: pair.default([1, 1]),
})

export const complexMapSchema = z.object({
  map: z.enum(['square', 'exp', 'inverse', 'conjugate']).default('square'),
  z: pair.default([0.8, 0.4]),
  turn: z.number().finite().min(10).max(170).default(60),
})

export const complexFieldSchema = z.object({
  fn: z.enum(['square', 'cube', 'exp', 'inverse', 'sine']).default('square'),
  probe: pair.default([1, 0.6]),
  levels: z.number().int().min(3).max(8).default(5),
  showImag: z.boolean().default(true),
})

export const coordinateSchema = z.object({
  mode: z.enum(['measure', 'intersect']).default('measure'),
  a: pair.default([-3, 1]),
  b: pair.default([2, 4]),
  c: pair.default([1, -2]),
  slope: z.number().finite().min(-3).max(3).default(1),
})

export const familySchema = z.object({
  emphasis: z.enum(['normal', 'uniform', 'exponential']).default('uniform'),
  x: z.number().finite().min(-4.5).max(5.5).default(1),
})

export const solid3dSchema = z.object({
  cut: z.number().finite().min(0).max(4).default(2),
  showCut: z.boolean().default(true),
})

export const contourIntegralSchema = z.object({
  path: z.enum(['circle', 'ellipse', 'wavy', 'line']).default('circle'),
  fn: z.enum(['inv', 'identity', 'square', 'conj', 'exp']).default('inv'),
  reverse: z.boolean().default(false),
  radius: z.number().finite().min(0.3).max(2.5).default(1),
  centerX: z.number().finite().min(-2).max(2).default(0),
  t: z.number().finite().min(0).max(1).default(1),
})

const triple = z.tuple([z.number().finite(), z.number().finite(), z.number().finite()])

export const gaussianSchema = z.object({
  rows: z.tuple([triple, triple]).default([
    [2, 1, 5],
    [1, -1, 1],
  ]),
})

export const odeSlopeSchema = z.object({
  field: z.enum(['exponential', 'linear', 'logistic', 'oscillating', 'sum']).default('exponential'),
  x0: z.number().finite().min(-3).max(3).default(0),
  y0: z.number().finite().min(-3).max(3).default(1),
  step: z.number().finite().min(0.02).max(0.5).default(0.1),
  method: z.enum(['rk4', 'euler']).default('rk4'),
})

export const rootsSchema = z.object({
  fn: z.enum(['quadratic', 'cubic', 'cosine', 'exp']).default('quadratic'),
  method: z.enum(['bisection', 'newton']).default('newton'),
  a: z.number().finite().min(-3).max(3).default(0),
  b: z.number().finite().min(-3).max(3).default(2),
  start: z.number().finite().min(-3).max(3).default(1),
  iterations: z.number().int().min(1).max(12).default(5),
})

export const graphSchema = z.object({
  algorithm: z
    .enum(['bfs', 'dfs', 'dijkstra', 'kruskal', 'topological', 'maxflow', 'bipartite', 'colouring'])
    .default('bfs'),
  start: z.enum(['A', 'B', 'C', 'D', 'E', 'F']).default('A'),
  showTree: z.boolean().default(true),
})

export const probabilityTreeSchema = z.object({
  prior: z.number().finite().min(0).max(1).default(0.01),
  sensitivity: z.number().finite().min(0).max(1).default(0.95),
  specificity: z.number().finite().min(0).max(1).default(0.9),
})

export const fractionLineSchema = z.object({
  numerator: z.number().int().min(0).max(8).default(3),
  denominator: z.number().int().min(1).max(8).default(4),
  otherNumerator: z.number().int().min(0).max(8).default(1),
  otherDenominator: z.number().int().min(1).max(8).default(2),
})

export const functionGraphSchema = z.object({
  family: z.enum(['quadratic', 'cubic', 'exponential']).default('quadratic'),
  a: z.number().finite().min(-2).max(2).default(1),
  h: z.number().finite().min(-2).max(2).default(0),
  k: z.number().finite().min(-2).max(2).default(0),
  showInverse: z.boolean().default(false),
})

export const numberLineSchema = z.object({
  mode: z.enum(['signed', 'decimal', 'percent', 'ratio', 'balance']).default('signed'),
  a: z.number().finite().min(-10).max(10).default(-4),
  b: z.number().finite().min(-10).max(10).default(2),
  tenths: z.number().int().min(0).max(10).default(7),
  percent: z.number().finite().min(-100).max(100).default(25),
  base: z.number().finite().min(0).max(200).default(80),
  ratioA: z.number().int().min(1).max(12).default(6),
  ratioB: z.number().int().min(1).max(12).default(4),
  guess: z.number().finite().min(-10).max(10).default(0),
  constant: z.number().finite().min(-10).max(10).default(3),
  total: z.number().finite().min(-10).max(10).default(8),
})

export const modularSchema = z.object({
  mode: z.enum(['residue', 'power', 'euclid', 'factor']).default('residue'),
  a: z.number().int().min(-50).max(100).default(17),
  n: z.number().int().min(2).max(12).default(5),
  base: z.number().int().min(1).max(12).default(2),
  exponent: z.number().int().min(1).max(30).default(10),
  b: z.number().int().min(1).max(100).default(30),
})

export const groupSchema = z.object({
  group: z.enum(['z4', 'z5', 's3']).default('z4'),
  element: z.number().int().min(0).max(5).default(1),
})

export const discreteSchema = z.object({
  mode: z.enum(['truth', 'sets', 'combinations']).default('truth'),
  op: z.enum(['and', 'or', 'implies', 'xor']).default('and'),
  leftMask: z.number().int().min(0).max(63).default(7),
  rightMask: z.number().int().min(0).max(63).default(28),
  n: z.number().int().min(1).max(8).default(5),
  k: z.number().int().min(0).max(8).default(2),
})

export const optimizationSchema = z.object({
  mode: z.enum(['descent', 'feasible', 'lagrange', 'newton']).default('descent'),
  startX: z.number().finite().min(-3).max(3).default(2.5),
  startY: z.number().finite().min(-3).max(3).default(2),
  rate: z.number().finite().min(0.01).max(0.5).default(0.15),
  steps: z.number().int().min(1).max(20).default(6),
  candidateX: z.number().finite().min(0).max(3).default(1),
  candidateY: z.number().finite().min(0).max(3).default(1),
  rhs: z.enum(['none', 'first', 'second']).default('none'),
  delta: z.number().finite().min(-1).max(1).default(0),
})

export const analysisSchema = z.object({
  mode: z.enum(['sequence', 'series', 'metric']).default('sequence'),
  sequence: z.enum(['reciprocal', 'ratio', 'alternating', 'shifted']).default('reciprocal'),
  epsilon: z.number().finite().min(0.01).max(1).default(0.2),
  metric: z.enum(['euclidean', 'manhattan', 'chebyshev']).default('euclidean'),
  radius: z.number().finite().min(0.5).max(3).default(2),
})

export const modelFitSchema = z.object({
  mode: z.enum(['regression', 'logistic']).default('regression'),
  rate: z.number().finite().min(0.01).max(0.5).default(0.1),
  steps: z.number().int().min(1).max(20).default(8),
  threshold: z.number().finite().min(-2).max(2).default(0),
})

// y = A f(B(x − C)) + D. B is positive in the figure, so the period is 2π/B
// and a negative A is the reflection the amplitude lesson talks about.
export const sinusoidSchema = z.object({
  form: z.enum(['sin', 'cos']).default('sin'),
  amplitude: z.number().finite().min(-3).max(3).default(1),
  frequency: z.number().finite().min(0.25).max(4).default(1),
  phase: z.number().finite().min(-3.14).max(3.14).default(0),
  vertical: z.number().finite().min(-8).max(8).default(0),
  showReference: z.boolean().default(true),
})

export const inferenceSchema = z.object({
  dist: z.enum(['z', 't']).default('z'),
  df: z.number().int().min(1).max(60).default(10),
  alpha: z.number().finite().min(0.01).max(0.2).default(0.05),
  effect: z.number().finite().min(-1).max(3).default(0.8),
  observed: z.number().finite().min(-4).max(6).default(1.96),
})

export const phaseSchema = z.object({
  mode: z.enum(['line', 'plane']).default('line'),
  lineModel: z.enum(['decay', 'logistic', 'bistable', 'threshold']).default('logistic'),
  system: z.enum(['saddle', 'stable-node', 'spiral', 'center']).default('saddle'),
  showTrajectories: z.boolean().default(true),
})

export const vectorFieldSchema = z.object({
  preset: z.enum(['source', 'vortex', 'saddle', 'shear']).default('vortex'),
  radius: z.number().finite().min(0.4).max(2).default(1.2),
  showStreamlines: z.boolean().default(true),
})

export const surfaceSchema = z.object({
  field: z.enum(['bowl', 'saddle', 'product', 'ripple']).default('bowl'),
  x0: z.number().finite().min(-2).max(2).default(1),
  y0: z.number().finite().min(-2).max(2).default(0.6),
})

export const regionSchema = z.object({
  region: z.enum(['rect', 'between', 'polar']).default('between'),
  density: z.enum(['one', 'x', 'y']).default('one'),
  n: z.number().int().min(4).max(40).default(10),
})

export const seriesSumSchema = z.object({
  kind: z.enum(['geometric', 'harmonic', 'p2', 'alternating']).default('geometric'),
  terms: z.number().int().min(2).max(40).default(12),
})

export const accumulationSchema = z.object({
  fn: z.enum(['square', 'sine', 'decay']).default('square'),
  a: z.number().finite().min(0).max(2).default(0.5),
  x: z.number().finite().min(0).max(3).default(2),
})

export const solidRevSchema = z.object({
  mode: z.enum(['disk', 'washer', 'shell']).default('washer'),
  slice: z.number().finite().min(0.1).max(2).default(1.2),
})

export const conditionalAreaSchema = z.object({
  pB: z.number().finite().min(0.1).max(0.9).default(0.4),
  pAgivenB: z.number().finite().min(0).max(1).default(0.8),
  pAgivenBc: z.number().finite().min(0).max(1).default(0.3),
})
