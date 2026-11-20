// Fractions as exact integer pairs, with the arithmetic the figures need.
// Values are never compared as floating-point decimals: 1/3 and 2/6 are equal
// because the cross products are equal, not because both rounded to 0.333.

export interface Fraction {
  numerator: number
  denominator: number
}

export function gcd(a: number, b: number): number {
  let x = Math.abs(Math.trunc(a))
  let y = Math.abs(Math.trunc(b))
  while (y !== 0) {
    const next = x % y
    x = y
    y = next
  }
  return x || 1
}

export function lcm(a: number, b: number): number {
  const x = Math.abs(Math.trunc(a))
  const y = Math.abs(Math.trunc(b))
  if (x === 0 || y === 0) return 0
  return (x / gcd(x, y)) * y
}

export function simplify(fraction: Fraction): Fraction {
  const { numerator, denominator } = fraction
  if (denominator === 0) return { numerator, denominator }
  const sign = denominator < 0 ? -1 : 1
  const divisor = gcd(numerator, denominator)
  return {
    numerator: (sign * numerator) / divisor,
    denominator: Math.abs(denominator) / divisor,
  }
}

export function fractionValue(fraction: Fraction): number {
  return fraction.numerator / fraction.denominator
}

export function compareFractions(a: Fraction, b: Fraction): -1 | 0 | 1 {
  const left = a.numerator * b.denominator
  const right = b.numerator * a.denominator
  if (left < right) return -1
  if (left > right) return 1
  return 0
}

export function addFractions(a: Fraction, b: Fraction): Fraction {
  const denominator = lcm(a.denominator, b.denominator)
  const left = a.numerator * (denominator / a.denominator)
  const right = b.numerator * (denominator / b.denominator)
  return simplify({ numerator: left + right, denominator })
}

export function equivalentFractions(fraction: Fraction, multiplier: number): Fraction {
  return {
    numerator: fraction.numerator * multiplier,
    denominator: fraction.denominator * multiplier,
  }
}
