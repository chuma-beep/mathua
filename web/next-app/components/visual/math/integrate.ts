import { evalPoly, type Coeffs } from './poly'

export type RiemannRule = 'left' | 'right' | 'mid'

/** One Riemann sum of a polynomial over [a, b]. */
export function riemannSum(coeffs: Coeffs, a: number, b: number, n: number, rule: RiemannRule): number {
  const dx = (b - a) / n
  let sum = 0
  for (let i = 0; i < n; i++) {
    const left = a + i * dx
    const sample = rule === 'left' ? left : rule === 'right' ? left + dx : left + dx / 2
    sum += evalPoly(coeffs, sample) * dx
  }
  return sum
}
