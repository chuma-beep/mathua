import type { Scores } from '../lib/api'

interface ProgressSummaryProps {
  scores: Scores
  weakByDomain?: Record<string, { id: string; label: string }[]>
}

const domainLabels = {
  arithmetic: 'Arithmetic',
  fractions: 'Fractions',
  prealgebra: 'Pre-Algebra',
  algebra: 'Algebra',
  geometry: 'Geometry',
  trigonometry: 'Trigonometry',
  complex_numbers: 'Complex Numbers',
  precalculus: 'Precalculus',
  calculus: 'Calculus',
  linear_algebra: 'Linear Algebra',
  statistics: 'Statistics',
  discrete_math: 'Discrete Math',
  number_theory: 'Number Theory',
  differential_equations: 'Differential Equations',
  abstract_algebra: 'Abstract Algebra',
  topology: 'Topology',
} satisfies Record<string, string>

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
      <div className="flex flex-wrap gap-4 justify-center mb-6">
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
      {weakDomains.length > 0 && (
        <div className="flex flex-wrap gap-3 justify-center mb-6">
          {weakDomains.map(([domain, concepts]) => {
            const label = domainLabels[domain] || domain
            return (
              <div key={domain} className="bg-mathua-surface border border-mathua-border rounded-lg px-4 py-2 text-center min-w-[120px]">
                <div className="font-mono text-[10px] uppercase tracking-[0.1em] text-mathua-muted">{label}</div>
                <div className="font-mono text-sm text-mathua-red mt-0.5">{concepts.length} need review</div>
              </div>
            )
          })}
        </div>
      )}
    </>
  )
}
