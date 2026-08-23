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
        Struggles
      </h3>
      <div
        style={{
          border: '0.5px solid var(--border)',
          fontFamily: "'IBM Plex Mono', monospace",
          fontSize: 11,
        }}
      >
        {domains.map((domain) => {
          const items = entries[domain]
          return (
            <div
              key={domain}
              style={{
                padding: '8px 10px',
                borderBottom: '0.5px solid var(--border)',
              }}
            >
              <div
                style={{
                  color: 'var(--text-muted)',
                  fontSize: 10,
                  marginBottom: 4,
                }}
              >
                {DOMAIN_LABELS[domain] ?? domain}
              </div>
              {items.map((item) => (
                <div
                  key={item.id}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 8,
                    padding: '2px 0',
                  }}
                >
                  <div
                    style={{
                      width: 40,
                      height: 4,
                      background: 'var(--border)',
                      flexShrink: 0,
                    }}
                  >
                    <div
                      style={{
                        height: '100%',
                        width: `${Math.round(item.weakness * 100)}%`,
                        background: 'var(--accent-blue)',
                        transition: 'width 0.3s',
                      }}
                    />
                  </div>
                  <span style={{ color: 'var(--text-primary)', fontSize: 11 }}>
                    {item.label}
                  </span>
                  <span style={{ color: 'var(--text-muted)', fontSize: 10, marginLeft: 'auto' }}>
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
