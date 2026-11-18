import type { LeaderboardView } from './api'

// Period keys are always YYYY-MM-DD in UTC. The calendar works in local days,
// so both directions are converted here, the same contract the admin date range
// picker uses. A Wednesday chosen east of UTC therefore still resolves to the
// Monday the server will use.

export function parsePeriodKey(key: string): Date | undefined {
  const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(key)
  if (!match) return undefined
  const value = new Date(Date.UTC(Number(match[1]), Number(match[2]) - 1, Number(match[3])))
  return Number.isNaN(value.getTime()) ? undefined : value
}

export function periodKey(value: Date): string {
  return value.toISOString().slice(0, 10)
}

// fromCalendarDay turns the local day the calendar draws into a UTC midnight.
export function fromCalendarDay(value: Date): Date {
  return new Date(Date.UTC(value.getFullYear(), value.getMonth(), value.getDate()))
}

// toCalendarDay turns a UTC midnight back into the local day the calendar
// should highlight.
export function toCalendarDay(value: Date): Date {
  return new Date(value.getUTCFullYear(), value.getUTCMonth(), value.getUTCDate())
}

export function weekStartKey(value: Date): string {
  const utc = fromCalendarDay(value)
  const day = utc.getUTCDay() || 7
  utc.setUTCDate(utc.getUTCDate() - (day - 1))
  return periodKey(utc)
}

export function monthStartKey(value: Date): string {
  return periodKey(new Date(Date.UTC(value.getFullYear(), value.getMonth(), 1)))
}

export function shiftWeekKey(key: string, delta: number): string {
  const value = parsePeriodKey(key)
  if (!value) return key
  value.setUTCDate(value.getUTCDate() + 7 * delta)
  return periodKey(value)
}

export function shiftMonthKey(key: string, delta: number): string {
  const value = parsePeriodKey(key)
  if (!value) return key
  return periodKey(new Date(Date.UTC(value.getUTCFullYear(), value.getUTCMonth() + delta, 1)))
}

export function weekLabel(key: string): string {
  const value = parsePeriodKey(key)
  if (!value) return 'This week'
  return `Week of ${value.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', timeZone: 'UTC' })}`
}

export function monthLabel(key: string): string {
  const value = parsePeriodKey(key)
  if (!value) return 'This month'
  return value.toLocaleDateString('en-GB', { month: 'long', year: 'numeric', timeZone: 'UTC' })
}

export function periodLabel(view: LeaderboardView, start?: string): string {
  if (view === 'all_time') return 'All time'
  if (!start) return view === 'monthly' ? 'This month' : 'This week'
  return view === 'monthly' ? monthLabel(start) : weekLabel(start)
}

export function scoreLabel(view: LeaderboardView, final = false): string {
  if (view === 'all_time') return 'Lifetime'
  if (view === 'monthly') return final ? 'That month' : 'This month'
  return final ? 'That week' : 'This week'
}

// monthsAvailable hides the Month tab until at least one past month exists.
// An empty second tab on a site launched last week reads as a broken feature.
export function monthsAvailable(earliest: string | undefined, now = new Date()): boolean {
  if (!earliest) return false
  const first = parsePeriodKey(earliest)
  if (!first) return false
  return monthStartKey(first) < monthStartKey(now)
}

export function todayKey(now = new Date()): string {
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${now.getFullYear()}-${month}-${day}`
}
