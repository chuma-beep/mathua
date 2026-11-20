// Three continuous families compared on one axis.
//
// The point of the figure is that a mean and a standard deviation do not name a
// distribution. All three here have standard deviation 1, so the shapes can be
// compared fairly, and they still differ in support, in where the mass sits and
// in the probability they give the same interval. The exponential's mean is 1
// rather than 0 because it starts at zero; that is stated in the readout instead
// of being hidden by shifting it.

import { normalPdf, normalProbability } from './normal'

export type Family = 'normal' | 'uniform' | 'exponential'

export const FAMILY_LABELS: Record<Family, string> = {
  normal: 'Normal (0, 1)',
  uniform: 'Uniform, width 2√3',
  exponential: 'Exponential, rate 1',
}

/** The half-width of the uniform family. √3 gives it standard deviation 1. */
const UNIFORM_HALF = Math.sqrt(3)

export function familyPdf(family: Family, x: number): number {
  switch (family) {
    case 'normal':
      return normalPdf(x, 0, 1)
    case 'uniform':
      return Math.abs(x) <= UNIFORM_HALF ? 1 / (2 * UNIFORM_HALF) : 0
    case 'exponential':
      return x >= 0 ? Math.exp(-x) : 0
  }
}

/** P(X ≤ x). The uniform and exponential cases are closed forms, not samples. */
export function familyCdf(family: Family, x: number): number {
  switch (family) {
    case 'normal':
      return normalProbability(x, 0, 1)
    case 'uniform': {
      const t = (x + UNIFORM_HALF) / (2 * UNIFORM_HALF)
      return Math.min(1, Math.max(0, t))
    }
    case 'exponential':
      return x <= 0 ? 0 : 1 - Math.exp(-x)
  }
}

/** P(X > x), which is the right tail and is what the caption calls the tail. */
export function familyTail(family: Family, x: number): number {
  return 1 - familyCdf(family, x)
}

export function familyMean(family: Family): number {
  switch (family) {
    case 'normal':
    case 'uniform':
      return 0
    case 'exponential':
      return 1
  }
}

/** All three are scaled to standard deviation 1, which is what makes the comparison fair. */
export function familySd(family: Family): number {
  switch (family) {
    case 'normal':
      return 1
    case 'uniform':
      return UNIFORM_HALF / Math.sqrt(3)
    case 'exponential':
      return 1
  }
}

export const FAMILIES: Family[] = ['normal', 'uniform', 'exponential']
