// The normal density and distribution function.
//
// The CDF is the standard Abramowitz-Stegun 7.1.26 rational approximation, whose
// absolute error is below 7.5e-8. That is far more precision than a figure needs,
// and it keeps the module free of a statistics dependency. `normalCdf` takes a
// z-score; `normalProbability` takes the raw value and parameters.

export function normalPdf(x: number, mu: number, sigma: number): number {
  const z = (x - mu) / sigma
  return Math.exp(-0.5 * z * z) / (sigma * Math.sqrt(2 * Math.PI))
}

export function normalCdf(z: number): number {
  const t = 1 / (1 + 0.2316419 * Math.abs(z))
  const d = 0.3989422804014327 * Math.exp(-0.5 * z * z)
  const p =
    d * t * (0.319381530 + t * (-0.356563782 + t * (1.781477937 + t * (-1.821255978 + t * 1.330274429))))
  return z >= 0 ? 1 - p : p
}

export function normalProbability(x: number, mu: number, sigma: number): number {
  return normalCdf((x - mu) / sigma)
}
