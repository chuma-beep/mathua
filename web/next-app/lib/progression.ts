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
