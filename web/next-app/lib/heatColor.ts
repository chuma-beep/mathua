const HEAT_MAX = 4

function clampLevel(level: number): number {
  if (!Number.isFinite(level)) return 0
  return Math.min(HEAT_MAX, Math.max(0, Math.trunc(level)))
}

export function heatColor(level: number): string {
  return `var(--heat-${clampLevel(level)})`
}

export function heatTextColor(level: number): string {
  const lvl = clampLevel(level)
  return lvl === 0 ? 'var(--text-muted)' : `var(--heat-fg-${lvl})`
}
