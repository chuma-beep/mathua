import type { Scores } from '../lib/api'
import { domainLabel } from '../lib/graphDomains'

interface ProgressSummaryProps {
  scores: Scores
  weakByDomain?: Record<string, { id: string; label: string }[]>
}



export default function ProgressSummary({ scores, weakByDomain }: ProgressSummaryProps) {
  const stats = [
    { label: 'Concepts mastered', value: scores.concepts_mastered, color: 'var(--accent-blue)' },
    { label: 'Day streak', value: scores.current_streak, color: 'var(--accent-teal)' },
    { label: 'Level', value: scores.level, color: 'var(--accent-blue)' },
    { label: 'Weekly score', value: scores.weekly_score, color: 'var(--accent-blue)' },
  ]

  const weakDomains = weakByDomain
    ? Object.entries(weakByDomain).filter(([, concepts]) => concepts.length > 0).sort((a, b) => b[1].length - a[1].length).slice(0, 5)
    : []

  return (
    <>
      <div className="mb-4 grid grid-cols-2 gap-2 justify-center sm:mb-6 sm:flex sm:flex-wrap sm:gap-4">
        {stats.map((s) => (
          <div
            key={s.label}
            className="bg-mathua-surface border border-mathua-border p-3 flex-1 min-w-0 max-w-none text-center sm:p-5 sm:min-w-[140px] sm:max-w-[200px]"
          >
            <div
              className="font-mono text-xl font-light mt-1 sm:text-2xl"
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
      {weakDomains.length > 0 && (
        <div className="mb-4 grid grid-cols-2 gap-2 justify-center sm:mb-6 sm:flex sm:flex-wrap sm:gap-3">
          {weakDomains.map(([domain, concepts]) => {
            const label = domainLabel(domain)
            return (
              <div key={domain} className="bg-mathua-surface border border-mathua-border px-3 py-2 text-center min-w-0 sm:min-w-[120px] sm:px-4">
                <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">{label}</div>
                <div className="font-mono text-[11px] text-mathua-red mt-0.5 sm:text-sm">{concepts.length} need review</div>
              </div>
            )
          })}
        </div>
      )}
    </>
  )
}
