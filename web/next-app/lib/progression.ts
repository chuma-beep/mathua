// Progression rules for lesson practice (MA pedagogy: worked example →
// ≤5 problems, 2-in-a-row correct to advance, improve.md:39).
export const REQUIRED_IN_A_ROW = 2

export interface StreakState {
  consecutive: number
  advanced: boolean
}

export function nextStreak(prev: number, correct: boolean): number {
  if (correct) return prev + 1
  return 0
}

export function isAdvanced(consecutive: number): boolean {
  return consecutive >= REQUIRED_IN_A_ROW
}

export function initialState(): StreakState {
  return { consecutive: 0, advanced: false }
}

export function applyResult(state: StreakState, correct: boolean): StreakState {
  const consecutive = nextStreak(state.consecutive, correct)
  return { consecutive, advanced: state.advanced || isAdvanced(consecutive) }
}

export interface Attempt {
  correct: boolean
  difficulty: number
  elapsed: number
  variation: string
}

export interface MasteryEstimate {
  score: number
  band: 'unseen' | 'developing' | 'strong' | 'well retained'
  decision: 'advance' | 'continue' | 'intervene'
}

// Mastery estimate from attempt history: accuracy (recent-weighted) ×
// difficulty × variation × time. Streak is one input, not the decision.
export function masteryEstimate(history: Attempt[]): MasteryEstimate {
  if (history.length === 0) return { score: 0, band: 'unseen', decision: 'continue' }
  const recent = history.slice(-6)
  let num = 0
  let den = 0
  recent.forEach((a, i) => {
    const w = 0.5 + (i / Math.max(1, recent.length - 1)) * 0.5 // recent counts more
    den += w
    if (a.correct) num += w * (0.6 + 0.4 * a.difficulty)
  })
  const accuracy = den > 0 ? num / den : 0
  const distinct = new Set(recent.map(a => a.variation)).size
  const variationBonus = Math.min(0.15, 0.05 * Math.max(0, distinct - 1))
  const times: number[] = []
  for (const a of recent) {
    if (a.correct) times.push(a.elapsed)
  }
  times.sort((x, y) => x - y)
  const median = times.length > 0 ? times[Math.floor(times.length / 2)] : 0
  const timeFactor = median <= 0 ? 1 : median > 60 ? 0.8 : 1
  const score = Math.min(1, (accuracy + variationBonus) * timeFactor)
  const band = score >= 0.85 ? 'well retained' : score >= 0.6 ? 'strong' : score >= 0.3 ? 'developing' : 'unseen'
  // Advance needs both quality and evidence: score + at least 3 attempts with
  // 2-in-a-row at the end. Otherwise continue, or intervene on repeated misses.
  const tail2 = recent.slice(-2)
  const streak2 = tail2.length === 2 && tail2.every(a => a.correct)
  const misses = recent.slice(-3).filter(a => !a.correct).length
  const decision = score >= 0.6 && recent.length >= 3 && streak2 ? 'advance' : misses >= 3 ? 'intervene' : 'continue'
  return { score, band, decision }
}
