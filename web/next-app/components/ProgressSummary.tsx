import type { Scores } from '../lib/api'

interface ProgressSummaryProps {
  scores: Scores
}

export default function ProgressSummary({ scores }: ProgressSummaryProps) {
  const stats = [
    { label: 'Concepts mastered', value: scores.concepts_mastered, color: 'var(--accent-gold)' },
    { label: 'Day streak', value: scores.current_streak, color: 'var(--accent-teal)' },
    { label: 'Level', value: scores.level, color: 'var(--accent-blue)' },
    { label: 'Weekly score', value: scores.weekly_score, color: 'var(--accent-gold)' },
  ]

  return (
    <div className="flex flex-wrap gap-4 justify-center mb-8">
      {stats.map((s) => (
        <div
          key={s.label}
          className="bg-mathua-surface border border-mathua-border rounded-lg p-5 flex-1 min-w-[140px] max-w-[200px] text-center"
        >
          <div
            className="font-mono text-2xl font-light mt-1"
            style={{ color: s.color }}
          >
            {s.value}
          </div>
          <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted mt-1">
            {s.label}
          </div>
        </div>
      ))}
    </div>
  )
}
