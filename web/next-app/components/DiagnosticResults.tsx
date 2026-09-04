'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import type { GoalPlanRes } from '../lib/api'
import { getLessons } from '../lib/api'
import { getUserInfo } from '../lib/auth'

const DOMAIN_LABELS: Record<string, string> = {
  arithmetic: 'Arithmetic', fractions: 'Fractions', prealgebra: 'Pre-Algebra',
  algebra: 'Algebra', geometry: 'Geometry', trigonometry: 'Trigonometry',
  complex_numbers: 'Complex Numbers', precalculus: 'Precalculus', calculus: 'Calculus',
  linear_algebra: 'Linear Algebra', statistics: 'Statistics', discrete_math: 'Discrete Math',
  number_theory: 'Number Theory', differential_equations: 'Differential Equations',
  abstract_algebra: 'Abstract Algebra', topology: 'Topology',
}

interface Props {
  plan: GoalPlanRes
  onStartPractice?: () => void
}

export default function DiagnosticResults({ plan, onStartPractice }: Props) {
  const [recommended, setRecommended] = useState<Record<string, { title: string; concepts: string[] }[]>>({})
  const [recLoading, setRecLoading] = useState(true)
  const allWeakIds = Object.values(plan.weak_areas).flat().map((c) => c.id)
  const [expanded, setExpanded] = useState<string | null>(null)

  useEffect(() => {
    if (allWeakIds.length === 0) { setRecLoading(false); return }
    const user = getUserInfo()
    getLessons(user?.student_id).then((res) => {
      const rec: Record<string, { title: string; concepts: string[] }[]> = {}
      for (const [domain, lessons] of Object.entries(res.lessons)) {
        const filtered = lessons.filter((l) => l.concepts.some((cid) => allWeakIds.includes(cid))).slice(0, 3)
        if (filtered.length) rec[domain] = filtered
      }
      setRecommended(rec)
    }).catch(() => {}).finally(() => setRecLoading(false))
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  const domains = new Set([...Object.keys(plan.weak_areas), ...Object.keys(plan.strong_areas)])
  const readinessPct = Math.round(plan.readiness * 100)

  return (
    <div className="max-w-3xl mx-auto">
      <div className="text-center mb-8">
        <div className="relative inline-flex items-center justify-center">
          <svg width="140" height="140" className="-rotate-90">
            <circle cx="70" cy="70" r="60" fill="none" stroke="var(--code-bg)" strokeWidth="10" />
            <circle cx="70" cy="70" r="60" fill="none" stroke="var(--accent-blue)" strokeWidth="10" strokeLinecap="round" strokeDasharray={`${2 * Math.PI * 60}`} strokeDashoffset={`${2 * Math.PI * 60 * (1 - plan.readiness)}`} className="transition-all duration-1000" />
          </svg>
          <span className="absolute text-3xl font-mono font-light text-mathua-primary">{readinessPct}%</span>
        </div>
        <p className="text-mathua-secondary text-sm mt-2">diagnostic test score {plan.total_tested > 0 && `(${plan.correct_count}/${plan.total_tested} correct)`}</p>
      </div>

      <div className="space-y-3 mb-8">
        {Array.from(domains).sort().map((domain) => {
          const label = DOMAIN_LABELS[domain] || domain
          const weak = plan.weak_areas[domain] || []
          const strong = plan.strong_areas[domain] || []
          const total = weak.length + strong.length
          const weakPct = total > 0 ? Math.round((weak.length / total) * 100) : 0
          const isExpanded = expanded === domain
          return (
            <div key={domain} className="border border-mathua-border bg-mathua-surface">
              <button onClick={() => setExpanded(isExpanded ? null : domain)} className="w-full flex items-center justify-between p-3 text-left">
                <span className="font-mono text-xs text-mathua-primary">{label}</span>
                <span className="font-mono text-[11px] text-mathua-muted">{weak.length} to review · {strong.length} strong</span>
              </button>
              <div className="px-3 pb-2">
                <div className="h-1.5 bg-mathua-code overflow-hidden">
                  <div className="h-full bg-mathua-blue transition-all" style={{ width: `${weakPct}%` }} />
                </div>
              </div>
              {isExpanded && (
                <div className="px-3 pb-3 border-t border-mathua-border mt-1 pt-2">
                  {weak.length > 0 && (
                    <div className="mb-2">
                      <p className="font-mono text-[10px] uppercase text-mathua-red mb-1">Needs review</p>
                      <div className="flex flex-wrap gap-1.5">
                        {weak.map((c) => (
                          <Link key={c.id} href={`/concept?id=${encodeURIComponent(c.id)}`} className="text-xs font-mono border border-mathua-border px-2 py-1 hover:border-mathua-blue hover:text-mathua-blue">
                            {c.label}
                          </Link>
                        ))}
                      </div>
                    </div>
                  )}
                  {strong.length > 0 && (
                    <div>
                      <p className="font-mono text-[10px] uppercase text-mathua-green mb-1">Strong</p>
                      <p className="font-mono text-xs text-mathua-secondary">{strong.join(', ')}</p>
                    </div>
                  )}
                </div>
              )}
            </div>
          )
        })}
      </div>

      {recLoading && allWeakIds.length > 0 && (
        <p className="font-mono text-xs text-mathua-muted text-center mb-8">Finding matching lessons…</p>
      )}
      {!recLoading && allWeakIds.length > 0 && Object.keys(recommended).length === 0 && (
        <p className="font-mono text-xs text-mathua-muted text-center mb-8">
          No matching lessons found — <Link href="/study" className="text-mathua-blue hover:text-mathua-blue-hover">browse Study →</Link>
        </p>
      )}
      {Object.keys(recommended).length > 0 && (
        <div className="mb-8">
          <h3 className="font-serif text-lg text-mathua-primary mb-3">Recommended lessons</h3>
          <div className="grid gap-2">
            {Object.entries(recommended).map(([domain, lessons]) => (
              <div key={domain} className="border border-mathua-border p-3 bg-mathua-surface-elevated">
                <p className="font-mono text-[11px] uppercase text-mathua-muted mb-2">{DOMAIN_LABELS[domain] || domain}</p>
                <div className="space-y-1">
                  {lessons.map((l) => (
                    <Link key={l.title} href={`/study?lesson=${encodeURIComponent(l.title)}`} className="block font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover">
                      → {l.title}
                    </Link>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {onStartPractice && (
        <div className="text-center">
          <button onClick={onStartPractice} className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-10 h-12 font-mono text-sm">
            Start practicing →
          </button>
        </div>
      )}
    </div>
  )
}
