'use client'

import type { WeaknessRes } from '../lib/api'

interface Props {
  weaknesses: WeaknessRes | null
}

const DOMAIN_LABELS: Record<string, string> = {
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

export default function StrugglesSection({ weaknesses }: Props) {
  const entries = weaknesses?.by_domain ?? {}
  const domains = Object.keys(entries)

  if (domains.length === 0) return null

  return (
    <div className="w-full max-w-full min-w-0 overflow-hidden">
      <h3 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
        Struggles
      </h3>
      <div className="border-[0.5px] border-mathua-border font-mono text-[11px] w-full max-w-full min-w-0 overflow-hidden">
        {domains.map((domain) => {
          const items = entries[domain]
          return (
            <div
              key={domain}
              className="px-3 py-2 border-b-[0.5px] border-mathua-border last:border-b-0 w-full max-w-full min-w-0 overflow-hidden"
            >
              <div className="text-mathua-muted text-[10px] mb-1 truncate">
                {DOMAIN_LABELS[domain] ?? domain}
              </div>
              {items.map((item) => (
                <div
                  key={item.id}
                  className="flex min-w-0 items-center gap-2 py-0.5"
                >
                  <div className="w-10 h-1 bg-mathua-border shrink-0">
                    <div
                      className="h-full bg-mathua-blue transition-all duration-300"
                      style={{ width: `${Math.round(item.weakness * 100)}%` }}
                    />
                  </div>
                  <span className="min-w-0 flex-1 truncate text-mathua-primary text-[11px]">
                    {item.label}
                  </span>
                  <span className="shrink-0 text-mathua-muted text-[10px]">
                    {Math.round(item.weakness * 100)}%
                  </span>
                </div>
              ))}
            </div>
          )
        })}
      </div>
    </div>
  )
}
