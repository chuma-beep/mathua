// Polynomials in ascending coefficient order: coeffs[i] multiplies x^i.
// A visual names coefficients, never an arbitrary expression, so every derivative
// and integral here is exact and has no parser to disagree with the renderer.

export type Coeffs = number[]

export function evalPoly(coeffs: Coeffs, x: number): number {
  let value = 0
  for (let i = coeffs.length - 1; i >= 0; i--) value = value * x + coeffs[i]
  return value
}

export function derivative(coeffs: Coeffs): Coeffs {
  return coeffs.slice(1).map((c, i) => c * (i + 1))
}

/** The antiderivative with constant 0. */
export function antiderivative(coeffs: Coeffs): Coeffs {
  return [0, ...coeffs.map((c, i) => c / (i + 1))]
}

export function definiteIntegral(coeffs: Coeffs, a: number, b: number): number {
  const f = antiderivative(coeffs)
  return evalPoly(f, b) - evalPoly(f, a)
}
