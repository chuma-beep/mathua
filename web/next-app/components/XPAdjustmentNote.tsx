'use client'

type Props = {
  adjustment?: number
  note?: string
  source?: 'performance' | 'effort' | 'legacy' | string
  className?: string
}

function sourceLabel(source?: string): string {
  if (source === 'performance') return 'Performance'
  if (source === 'effort') return 'Effort'
  if (source === 'legacy') return 'Fast-answer rule'
  return 'XP'
}

/**
 * One neutral line explaining an accountability adjustment.
 *
 * The server sends the amount that actually reached the ledger, so this
 * component never computes a deduction of its own. When the amount is zero but
 * a note exists, the note still renders: a learner should see why a task
 * produced no XP as clearly as why it lost some.
 */
export default function XPAdjustmentNote({ adjustment, note, source, className = '' }: Props) {
  if (!adjustment && !note) return null
  const value = adjustment ?? 0
  const tone = value < 0 ? 'text-mathua-red' : value > 0 ? 'text-mathua-green' : 'text-mathua-muted'
  return (
    <p className={`font-mono text-[11px] leading-relaxed ${tone} ${className}`} role="status">
      <span className="uppercase tracking-[0.14em] text-mathua-muted">{sourceLabel(source)}</span>
      {' · '}
      {note ?? `${value > 0 ? '+' : ''}${value} XP`}
    </p>
  )
}
