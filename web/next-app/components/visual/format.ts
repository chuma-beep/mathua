// Presentation formatting. Every number a learner reads goes through here so a
// template cannot show `0.30000000000000004` beside another that shows `0.3`.

export function fmt(value: number, digits = 2): string {
  if (!Number.isFinite(value)) return '—'
  const safe = Math.abs(value) < 1e-10 ? 0 : value
  const text = safe.toFixed(digits)
  if (!text.includes('.')) return text
  return text.replace(/0+$/, '').replace(/\.$/, '')
}

export function fmtSigned(value: number, digits = 2): string {
  const text = fmt(value, digits)
  return value > 0 ? `+${text}` : text
}

export function fmtVec(v: [number, number], digits = 2): string {
  return `(${fmt(v[0], digits)}, ${fmt(v[1], digits)})`
}

export function fmtRadians(value: number): string {
  return `${fmt(value, 2)} rad`
}
