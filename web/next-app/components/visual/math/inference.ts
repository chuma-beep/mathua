// Two-sided inference in standard-error units.
//
// The figure shows one null distribution, one alternative distribution, the
// observed statistic, the two tails that define alpha, and the interval the
// same critical value produces. Everything is computed here, so the test
// statistic, the p-value, the decision and the power cannot drift apart.

export type InferenceDistribution = 'z' | 't'

export interface InferenceInput {
  dist: InferenceDistribution
  df: number
  alpha: number
  effect: number
  observed: number
}

export interface InferenceModel {
  dist: InferenceDistribution
  df: number
  alpha: number
  effect: number
  observed: number
  critical: number
  statistic: number
  pValue: number
  reject: boolean
  power: number
  typeTwo: number
  /** Interval in standard-error units. */
  interval: [number, number]
  status: 'ok' | 'invalid'
  note?: string
}

const EPS = 1e-12
const MAX_Z = 100

export function normalPdf(x: number): number {
  return Math.exp(-0.5 * x * x) / Math.sqrt(2 * Math.PI)
}

function logGamma(value: number): number {
  const coefficients = [
    676.5203681218851,
    -1259.1392167224028,
    771.3234287776531,
    -176.6150291621406,
    12.507343278686905,
    -0.13857109526572012,
    9.984369578019572e-6,
    1.5056327351493116e-7,
  ]
  if (value < 0.5) {
    return Math.log(Math.PI / Math.sin(Math.PI * value)) - logGamma(1 - value)
  }
  let x = 0.9999999999998099
  const shifted = value - 1
  for (let index = 0; index < coefficients.length; index++) {
    x += coefficients[index] / (shifted + index + 1)
  }
  const t = shifted + coefficients.length - 0.5
  return 0.5 * Math.log(2 * Math.PI) + (shifted + 0.5) * Math.log(t) - t + Math.log(x)
}

/** Student t density with a positive integer or real degrees of freedom. */
export function studentPdf(x: number, df: number): number {
  if (!Number.isFinite(x) || !Number.isFinite(df) || df <= 0) return Number.NaN
  const logCoefficient =
    logGamma((df + 1) / 2) -
    logGamma(df / 2) -
    0.5 * Math.log(df * Math.PI)
  return Math.exp(logCoefficient - ((df + 1) / 2) * Math.log1p((x * x) / df))
}

function simpson(f: (x: number) => number, a: number, b: number, intervals: number): number {
  const n = intervals % 2 === 0 ? intervals : intervals + 1
  const h = (b - a) / n
  let total = f(a) + f(b)
  for (let index = 1; index < n; index++) {
    total += f(a + index * h) * (index % 2 === 0 ? 2 : 4)
  }
  return (h / 3) * total
}

export function normalCdf(x: number): number {
  const scaled = x / Math.SQRT2
  return 0.5 * (1 + erf(scaled))
}

/** Abramowitz and Stegun 7.1.26, accurate to about 1.5e-7. */
function erf(x: number): number {
  const sign = x < 0 ? -1 : 1
  const absolute = Math.abs(x)
  const t = 1 / (1 + 0.3275911 * absolute)
  const polynomial =
    (((((1.061405429 * t - 1.453152027) * t) + 1.421413741) * t - 0.284496736) * t +
      0.254829592) *
    t
  return sign * (1 - polynomial * Math.exp(-absolute * absolute))
}

export function distributionCdf(dist: InferenceDistribution, df: number, x: number): number {
  if (!Number.isFinite(x)) return x < 0 ? 0 : 1
  if (dist === 'z') return normalCdf(x)
  if (!Number.isFinite(df) || df <= 0) return Number.NaN
  if (x === 0) return 0.5
  if (x < 0) return 1 - distributionCdf(dist, df, -x)
  const mass = simpson((value) => studentPdf(value, df), 0, x, 400)
  return Math.min(1, Math.max(0, 0.5 + mass))
}

export function distributionPdf(dist: InferenceDistribution, df: number, x: number): number {
  return dist === 'z' ? normalPdf(x) : studentPdf(x, df)
}

export function twoSidedPValue(dist: InferenceDistribution, df: number, statistic: number): number {
  if (!Number.isFinite(statistic)) return Number.NaN
  const tail = 1 - distributionCdf(dist, df, Math.abs(statistic))
  return Math.min(1, Math.max(0, 2 * tail))
}

/** Positive critical value with total tail mass alpha. */
export function criticalValue(dist: InferenceDistribution, df: number, alpha: number): number {
  if (!Number.isFinite(alpha) || alpha <= 0 || alpha >= 1) return Number.NaN
  let low = 0
  let high = MAX_Z
  for (let iteration = 0; iteration < 80; iteration++) {
    const middle = (low + high) / 2
    const tail = 1 - distributionCdf(dist, df, middle)
    if (2 * tail > alpha) low = middle
    else high = middle
  }
  return (low + high) / 2
}

export function inferenceModel(input: InferenceInput): InferenceModel {
  const invalid = (note: string): InferenceModel => ({
    ...input,
    critical: Number.NaN,
    statistic: Number.NaN,
    pValue: Number.NaN,
    reject: false,
    power: Number.NaN,
    typeTwo: Number.NaN,
    interval: [Number.NaN, Number.NaN],
    status: 'invalid',
    note,
  })
  const values = [input.df, input.alpha, input.effect, input.observed]
  if (!values.every(Number.isFinite)) return invalid('enter finite numbers')
  if (input.alpha <= 0 || input.alpha >= 1) return invalid('alpha must lie between 0 and 1')
  if (input.dist === 't' && input.df < 1) return invalid('the t distribution needs at least one degree of freedom')

  const critical = criticalValue(input.dist, input.df, input.alpha)
  const statistic = input.observed
  const pValue = twoSidedPValue(input.dist, input.df, statistic)
  const reject = pValue <= input.alpha + EPS
  const upper = 1 - distributionCdf(input.dist, input.df, critical - input.effect)
  const lower = distributionCdf(input.dist, input.df, -critical - input.effect)
  const power = Math.min(1, Math.max(0, upper + lower))
  const typeTwo = 1 - power
  return {
    dist: input.dist,
    df: input.df,
    alpha: input.alpha,
    effect: input.effect,
    observed: input.observed,
    critical,
    statistic,
    pValue,
    reject,
    power,
    typeTwo,
    interval: [statistic - critical, statistic + critical],
    status: 'ok',
  }
}

/** Sample points for a density curve over [min, max]. */
export function densityCurve(
  dist: InferenceDistribution,
  df: number,
  min: number,
  max: number,
  count = 120,
): Array<[number, number]> {
  const safeCount = Math.max(2, Math.floor(count))
  const points: Array<[number, number]> = []
  for (let index = 0; index <= safeCount; index++) {
    const x = min + ((max - min) * index) / safeCount
    points.push([x, distributionPdf(dist, df, x)])
  }
  return points
}
