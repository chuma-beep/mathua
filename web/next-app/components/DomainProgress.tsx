'use client'

import { useMemo } from 'react'
import conceptsData from '../data/concepts.json'

const DOMAIN_ORDER = [
  'counting', 'arithmetic', 'fractions', 'prealgebra',
  'algebra', 'geometry', 'trigonometry', 'calculus',
  'statistics', 'linear_algebra', 'discrete_math',
  'complex_numbers', 'number_theory', 'differential_equations',
  'abstract_algebra', 'topology',
]

const DOMAIN_LABELS: Record<string, string> = {
  counting: 'Counting',
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
}

interface Props {
  progress: Record<string, { status: string }>
}

export default function DomainProgress({ progress }: Props) {
  const domains = useMemo(() => {
    // Group concepts by domain
    const byDomain = new Map<string, { total: number; mastered: number; learning: number }>()
    for (const c of conceptsData as { id: string; domain: string }[]) {
      const entry = byDomain.get(c.domain) || { total: 0, mastered: 0, learning: 0 }
      entry.total++
      const p = progress[c.id]
      if (p?.status === 'MASTERED') entry.mastered++
      else if (p?.status === 'LEARNING' || p?.status === 'PRACTICING') entry.learning++
      byDomain.set(c.domain, entry)
    }

    return DOMAIN_ORDER
      .filter((d) => byDomain.has(d))
      .map((d) => {
        const e = byDomain.get(d)!
        return { domain: d, label: DOMAIN_LABELS[d] ?? d, ...e }
      })
  }, [progress])

  if (domains.length === 0) return null

  return (
    <div>
      <h3
        style={{
          fontFamily: "'IBM Plex Serif', serif",
          fontSize: '1.05rem',
          fontWeight: 400,
          color: 'var(--text-primary)',
          marginBottom: 12,
        }}
      >
        Domain Progress
      </h3>
      <div
        style={{
          border: '0.5px solid var(--border)',
          fontFamily: "'IBM Plex Mono', monospace",
          fontSize: 11,
        }}
      >
        {domains.map((d, i) => {
          const pct = d.total > 0 ? Math.round((d.mastered / d.total) * 100) : 0
          return (
            <div
              key={d.domain}
              style={{
                display: 'flex',
                alignItems: 'center',
                padding: '6px 10px',
                borderBottom: i < domains.length - 1 ? '0.5px solid var(--border)' : 'none',
                gap: 10,
              }}
            >
              <div style={{ width: 130, color: 'var(--text-primary)' }}>{d.label}</div>
              <div style={{ flex: 1, height: 8, background: 'var(--border)', position: 'relative' }}>
                <div
                  style={{
                    position: 'absolute',
                    left: 0,
                    top: 0,
                    height: '100%',
                    width: `${pct}%`,
                    background: 'var(--accent-blue)',
                    transition: 'width 0.3s',
                  }}
                />
              </div>
              <div style={{ color: 'var(--text-muted)', width: 70, textAlign: 'right' }}>
                {d.mastered}/{d.total} · {pct}%
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
