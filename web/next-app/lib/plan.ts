// Small pure helpers for goal + planner math (client-side what-if).
// Server remains the authority for estimates; these only recompute
// displayed days from an already-fetched xp_remaining.

// Daily XP goal policy. One place, because the default used to be written four times —
// `PlanEditor`, `clampGoal`, the preset list, and `server.go`'s estimate route — and they
// had drifted: the shipped default is DEFAULT_DAILY_GOAL, but the planner's fallbacks
// still said 30 while `engine/plan.go` and `POST /api/plans/current` said 10. A learner with
// no stored goal was therefore planning against 3x the XP/day they were actually banking,
// which is what made the plan's day and month figures wrong.
//
// The server range is 10..100 and is enforced in internal/server. Keep these three numbers
// and the server's check in agreement; `test/plan.test.ts` and the Go goal test both pin
// the same bounds from each side.
export const GOAL_MIN = 10
export const GOAL_MAX = 100
export const DEFAULT_DAILY_GOAL = 30
export const GOAL_PRESETS = [10, 20, 30, 50, 100]

export function clampGoal(v: number): number {
  // A NaN used to become 30 silently, so clearing the custom field and blurring it
  // silently rewrote the learner's goal. Fall back to the documented default instead.
  if (!Number.isFinite(v)) return DEFAULT_DAILY_GOAL
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
