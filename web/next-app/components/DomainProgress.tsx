'use client'

import { useMemo } from 'react'
import conceptsData from '../data/concepts.json'

const DOMAIN_ORDER = [
'arithmetic', 'fractions', 'prealgebra',
  'algebra', 'geometry', 'trigonometry', 'calculus',
  'statistics', 'linear_algebra', 'discrete_math',
  'complex_numbers', 'number_theory', 'differential_equations',
  'abstract_algebra', 'topology',
]

const DOMAIN_LABELS = {
  arithmetic: 'Arithmetic',
  fractions: 'Fractions',
  prealgebra: 'Pre-Algebra',
  algebra: 'Algebra',
  geometry: 'Geometry',
  trigonometry: 'Trigonometry',
  calculus: 'Calculus',
  statistics: 'Statistics',
  linear_algebra: 'Linear Algebra',
  discrete_math: 'Discrete Math',
  complex_numbers: 'Complex Numbers',
  number_theory: 'Number Theory',
  differential_equations: 'Differential Equations',
  abstract_algebra: 'Abstract Algebra',
  topology: 'Topology',
} satisfies Record<string, string>

interface Props {
  progress: Record<string, { status: string }>
}

export default function DomainProgress({ progress }: Props) {
  const domains = useMemo(() => {
    const byDomain = new Map<string, { total: number; mastered: number; learning: number }>()
    for (const c of conceptsData as { id: string; domain: string }[]) {
      const entry = byDomain.get(c.domain) || { total: 0, mastered: 0, learning: 0 }
      entry.total++
      const p = progress[c.id]
      if (p?.status === 'MASTERED') entry.mastered++
      else if (p?.status === 'LEARNING' || p?.status === 'PRACTICING') entry.learning++
      byDomain.set(c.domain, entry)
    }

    return DOMAIN_ORDER.flatMap((d) => {
      const e = byDomain.get(d)
      if (!e) return []
      return [{ domain: d, label: DOMAIN_LABELS[d] ?? d, ...e }]
    })
  }, [progress])

  if (domains.length === 0) return null

  return (
    <div className="w-full max-w-full min-w-0 overflow-hidden">
      <h3 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
        Domain Progress
      </h3>
      <div className="border-[0.5px] border-mathua-border font-mono text-[11px] w-full max-w-full min-w-0 overflow-hidden">
        {domains.map((d, i) => {
          const pct = d.total > 0 ? Math.round((d.mastered / d.total) * 100) : 0
          return (
            <div
              key={d.domain}
              className={`flex flex-wrap sm:flex-nowrap items-center gap-2 sm:gap-[10px] px-3 py-2 sm:px-[10px] sm:py-1.5 min-w-0 ${i < domains.length - 1 ? 'border-b-[0.5px] border-mathua-border' : ''}`}
            >
              <div className="min-w-0 w-full sm:w-[110px] lg:w-[130px] shrink-0 truncate text-mathua-primary text-[11px] sm:shrink-0">
                {d.label}
              </div>
              <div className="flex-1 min-w-[60px] h-2 bg-mathua-border relative shrink">
                <div
                  className="absolute left-0 top-0 h-full bg-mathua-blue transition-all duration-300"
                  style={{ width: `${pct}%` }}
                />
              </div>
              <div className="hidden sm:block w-[70px] text-right shrink-0 text-mathua-muted text-[11px]">
                {d.mastered}/{d.total} · {pct}%
              </div>
              <div className="sm:hidden w-full text-right text-mathua-muted text-[10px] leading-none">
                {d.mastered}/{d.total} · {pct}%
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
