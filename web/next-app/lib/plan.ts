// Small pure helpers for goal + planner math (client-side what-if).
// Server remains the authority for estimates; these only recompute
// displayed days from an already-fetched xp_remaining.

export const GOAL_MIN = 1
export const GOAL_MAX = 10000
export const GOAL_PRESETS = [5, 10, 20, 30]

export function clampGoal(v: number): number {
  if (!Number.isFinite(v)) return 30
  return Math.min(GOAL_MAX, Math.max(GOAL_MIN, Math.round(v)))
}

// daysFor mirrors the server's effective-rate math: rate scaled by
// active days per week. Returns -1 when unknowable.
export function daysFor(xpRemaining: number, dailyRate: number, restDaysPerWeek: number): number {
  const rest = Math.min(6, Math.max(0, restDaysPerWeek))
  const effective = dailyRate * (7 - rest) / 7
  if (!(effective > 0) || !(xpRemaining > 0)) return -1
  return xpRemaining / effective
}

// monthLabel renders a YYYY-MM-DD date as "June 2026" — estimates are
// ranges, never promises, so day precision is deliberately dropped.
export function monthLabel(isoDate: string): string {
  const d = new Date(`${isoDate}T12:00:00Z`)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleString('en-US', { month: 'long', year: 'numeric', timeZone: 'UTC' })
}
