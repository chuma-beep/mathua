// Spread of a one-dimensional sample.
//
// The figure's whole argument is that variance is an average of squared gaps, so
// every number it prints comes from these functions and the bars it draws are
// those same gaps. The sample denominator is the default because that is what a
// learner meets in a statistics course, and the population form is offered beside
// it so the difference is visible rather than asserted.

export function mean(values: number[]): number {
  if (values.length === 0) return 0
  return values.reduce((sum, v) => sum + v, 0) / values.length
}

/** Signed gaps from the mean. They always sum to zero, which is the point of squaring. */
export function deviations(values: number[]): number[] {
  const m = mean(values)
  return values.map((v) => v - m)
}

export function squaredDeviations(values: number[]): number[] {
  return deviations(values).map((d) => d * d)
}

export function sumSquares(values: number[]): number {
  return squaredDeviations(values).reduce((sum, d) => sum + d, 0)
}

/**
 * The average squared gap. `sample` divides by n-1, which is the unbiased
 * estimator and the one a course teaches; `false` divides by n.
 */
export function variance(values: number[], sample = true): number {
  if (values.length === 0) return 0
  const denominator = sample ? values.length - 1 : values.length
  if (denominator <= 0) return 0
  return sumSquares(values) / denominator
}

export function stdDev(values: number[], sample = true): number {
  return Math.sqrt(variance(values, sample))
}
