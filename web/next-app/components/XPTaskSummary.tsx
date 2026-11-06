'use client'

import XPAdjustmentNote from './XPAdjustmentNote'

type Props = {
  answerXP: number
  adjustment: number
  note?: string
  source?: 'performance' | 'effort' | 'legacy' | string
  className?: string
}

function adjustmentLabel(source?: string): string {
  if (source === 'performance') return 'Performance adjustment'
  if (source === 'effort') return 'Effort adjustment'
  if (source === 'legacy') return 'Fast-answer rule'
  return 'Task adjustment'
}

function signed(value: number): string {
  return `${value > 0 ? '+' : ''}${value} XP`
}

function tone(value: number): string {
  return value < 0 ? 'text-mathua-red' : value > 0 ? 'text-mathua-green' : 'text-mathua-muted'
}

/**
 * The task-level accounting view: what the answers earned, what the actual
 * accountability adjustment was, and the resulting net. The adjustment is the
 * server's applied amount, never a client-side estimate.
 */
export default function XPTaskSummary({ answerXP, adjustment, note, source, className = '' }: Props) {
  const net = answerXP + adjustment
  return (
    <div className={`border border-mathua-border bg-mathua-surface p-4 ${className}`} role="group" aria-label="Task XP summary">
      <p className="font-mono text-[11px] uppercase tracking-[0.16em] text-mathua-muted">Task XP summary</p>
      <div className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div>
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">XP from answers</p>
          <p className={`mt-1 font-mono text-sm ${tone(answerXP)}`}>{signed(answerXP)}</p>
        </div>
        <div>
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">{adjustmentLabel(source)}</p>
          <p className={`mt-1 font-mono text-sm ${tone(adjustment)}`}>{signed(adjustment)}</p>
        </div>
        <div>
          <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted">Net XP</p>
          <p className={`mt-1 font-mono text-sm ${tone(net)}`}>{signed(net)}</p>
        </div>
      </div>
      {note ? (
        <XPAdjustmentNote adjustment={adjustment} note={note} source={source} className="mt-3" />
      ) : adjustment === 0 ? (
        <p className="mt-3 font-mono text-[11px] text-mathua-muted">No accountability adjustment on this task.</p>
      ) : null}
    </div>
  )
}
