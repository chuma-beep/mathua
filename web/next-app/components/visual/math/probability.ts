// A two-stage probability tree for a binary prior and a binary test.
//
// The tree is the sample space, not an illustration of it: every leaf is the
// product of the branch probabilities above it and the leaves sum to one. The
// posterior is computed from those leaves, so a figure cannot show a posterior
// that disagrees with its own branches.

export interface ProbabilityTree {
  prior: number
  sensitivity: number
  specificity: number
  disease: number
  healthy: number
  positiveGivenDisease: number
  negativeGivenDisease: number
  positiveGivenHealthy: number
  negativeGivenHealthy: number
  diseasePositive: number
  diseaseNegative: number
  healthyPositive: number
  healthyNegative: number
  positive: number
  negative: number
  diseaseGivenPositive: number
  healthyGivenPositive: number
  diseaseGivenNegative: number
  healthyGivenNegative: number
}

const clamp01 = (value: number) => Math.min(1, Math.max(0, value))

export interface ConditionalArea {
  pB: number
  pAgivenB: number
  pAgivenBc: number
  pA: number
  pAandB: number
  pAandBc: number
  pBc: number
  /** true when the two conditional probabilities agree. */
  independent: boolean
}

/**
 * The area model behind conditioning: the unit square is split into a column
 * for B of width P(B) and a column for its complement, and the event A is a
 * horizontal band inside each column whose height is the conditional
 * probability. The band's area inside the B column is P(B)·P(A|B), which is
 * P(A ∩ B) by definition, so the picture cannot disagree with the formula.
 */
export function conditionalArea(pB: number, pAgivenB: number, pAgivenBc: number): ConditionalArea | null {
  if (![pB, pAgivenB, pAgivenBc].every(Number.isFinite)) return null
  const b = clamp01(pB)
  const givenB = clamp01(pAgivenB)
  const givenBc = clamp01(pAgivenBc)
  const bc = 1 - b
  const aAndB = b * givenB
  const aAndBc = bc * givenBc
  return {
    pB: b,
    pAgivenB: givenB,
    pAgivenBc: givenBc,
    pA: aAndB + aAndBc,
    pAandB: aAndB,
    pAandBc: aAndBc,
    pBc: bc,
    independent: Math.abs(givenB - givenBc) < 1e-9,
  }
}


/**
 * Build the tree from a prior, a sensitivity and a specificity. Returns null
 * when a value is not finite, because there is no tree to draw in that case.
 */
export function probabilityTree(
  prior: number,
  sensitivity: number,
  specificity: number,
): ProbabilityTree | null {
  if (![prior, sensitivity, specificity].every(Number.isFinite)) return null
  const p = clamp01(prior)
  const se = clamp01(sensitivity)
  const sp = clamp01(specificity)
  const disease = p
  const healthy = 1 - p
  const positiveGivenDisease = se
  const negativeGivenDisease = 1 - se
  const positiveGivenHealthy = 1 - sp
  const negativeGivenHealthy = sp
  const diseasePositive = disease * positiveGivenDisease
  const diseaseNegative = disease * negativeGivenDisease
  const healthyPositive = healthy * positiveGivenHealthy
  const healthyNegative = healthy * negativeGivenHealthy
  const positive = diseasePositive + healthyPositive
  const negative = diseaseNegative + healthyNegative
  return {
    prior: p,
    sensitivity: se,
    specificity: sp,
    disease,
    healthy,
    positiveGivenDisease,
    negativeGivenDisease,
    positiveGivenHealthy,
    negativeGivenHealthy,
    diseasePositive,
    diseaseNegative,
    healthyPositive,
    healthyNegative,
    positive,
    negative,
    diseaseGivenPositive: positive > 0 ? diseasePositive / positive : 0,
    healthyGivenPositive: positive > 0 ? healthyPositive / positive : 0,
    diseaseGivenNegative: negative > 0 ? diseaseNegative / negative : 0,
    healthyGivenNegative: negative > 0 ? healthyNegative / negative : 0,
  }
}
