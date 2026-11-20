// The sampling distribution of the mean, from a deterministic pseudo-random
// stream.
//
// No Math.random: a figure that redraws a different histogram on every keystroke
// teaches nothing, because the learner cannot tell their change from the noise.
// A seeded generator makes the picture reproducible, so the standard error is the
// only thing that moves when n moves.

export type Population = 'uniform' | 'normal' | 'skewed'

/** A small deterministic generator (mulberry32). Returns values in [0, 1). */
export function rng(seed: number): () => number {
  let state = seed >>> 0
  return () => {
    state = (state + 0x6d2b79f5) >>> 0
    let t = state
    t = Math.imul(t ^ (t >>> 15), t | 1)
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61)
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

/** One draw from the named population. All three have mean 0 and sd 1. */
export function draw(population: Population, next: () => number): number {
  switch (population) {
    case 'uniform': {
      const u = next()
      return (u - 0.5) * Math.sqrt(12)
    }
    case 'normal': {
      // Box-Muller, guarded against log(0).
      const u = Math.max(next(), 1e-12)
      const v = next()
      return Math.sqrt(-2 * Math.log(u)) * Math.cos(2 * Math.PI * v)
    }
    case 'skewed': {
      // Exponential, centred and scaled to mean 0 sd 1.
      const u = Math.max(next(), 1e-12)
      return -Math.log(u) - 1
    }
  }
}

export function mean(values: number[]): number {
  if (values.length === 0) return 0
  return values.reduce((sum, v) => sum + v, 0) / values.length
}

/** The sample standard deviation with the n-1 denominator. */
export function sd(values: number[]): number {
  if (values.length < 2) return 0
  const m = mean(values)
  const variance = values.reduce((sum, v) => sum + (v - m) * (v - m), 0) / (values.length - 1)
  return Math.sqrt(variance)
}

/** `draws` sample means, each from `n` independent draws. */
export function sampleMeans(population: Population, n: number, draws: number, seed = 1): number[] {
  const next = rng(seed)
  const out: number[] = []
  for (let d = 0; d < draws; d++) {
    let sum = 0
    for (let i = 0; i < n; i++) sum += draw(population, next)
    out.push(sum / n)
  }
  return out
}

/** Counts into equal-width bins spanning [lo, hi]. */
export function histogram(values: number[], lo: number, hi: number, bins: number): number[] {
  const counts = new Array<number>(bins).fill(0)
  const width = (hi - lo) / bins
  if (width <= 0) return counts
  for (const v of values) {
    const index = Math.floor((v - lo) / width)
    if (index >= 0 && index < bins) counts[index]++
  }
  return counts
}

/**
 * The density of the population itself, used only as a shape reference beside the
 * histogram. All three have mean 0 and sd 1, so the comparison is fair.
 */
export function populationPdf(population: Population, x: number): number {
  switch (population) {
    case 'uniform': {
      const half = Math.sqrt(3)
      return Math.abs(x) <= half ? 1 / (2 * half) : 0
    }
    case 'normal':
      return Math.exp(-0.5 * x * x) / Math.sqrt(2 * Math.PI)
    case 'skewed':
      return x >= -1 ? Math.exp(-(x + 1)) : 0
  }
}
