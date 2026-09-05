'use client'

import { useState, useMemo } from 'react'
import Link from 'next/link'
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
  machine_learning: 'Machine Learning',
  machinelearning: 'Machine Learning',
  precalculus: 'Precalculus',
}

function tierFor(w: number): { label: string; color: string } {
  if (w >= 0.6) return { label: 'Struggling', color: '#ef4444' }
  return { label: 'Needs practice', color: '#f59e0b' }
}

export default function StrugglesSection({ weaknesses }: Props) {
  const entries = useMemo(() => weaknesses?.by_domain ?? {}, [weaknesses])
  const domains = Object.keys(entries)

  const [expanded, setExpanded] = useState(false)
  const [openDomains, setOpenDomains] = useState<Set<string>>(new Set())

  const flat = useMemo(() => {
    const out: { id: string; label: string; domain: string; weakness: number }[] = []
    for (const [domain, items] of Object.entries(entries)) {
      for (const it of items as { id: string; label: string; weakness: number }[]) {
        out.push({ domain, ...it })
      }
    }
    out.sort((a, b) => b.weakness - a.weakness)
    return out
  }, [entries])

  if (domains.length === 0) return null

  const struggling = flat.filter(f => f.weakness >= 0.6).length
  const needsPractice = flat.length - struggling
  const top = flat.slice(0, 3)

  const toggleDomain = (d: string) => {
    setOpenDomains(prev => {
      const next = new Set(prev)
      if (next.has(d)) next.delete(d)
      else next.add(d)
      return next
    })
  }

  return (
    <div className="w-full max-w-full min-w-0 overflow-hidden">
      <h3 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
        Needs attention
      </h3>
      <div className="border-[0.5px] border-mathua-border bg-mathua-surface w-full max-w-full min-w-0 overflow-hidden">
        {/* Summary row */}
        <div className="px-3 py-3 flex items-center gap-3 min-w-0">
          <div className="flex-1 min-w-0">
            <div className="font-mono text-[11px] text-mathua-primary">
              {struggling > 0 && <span className="text-mathua-red">{struggling} Struggling</span>}
              {struggling > 0 && needsPractice > 0 && <span className="text-mathua-muted"> · </span>}
              {needsPractice > 0 && <span className="text-mathua-muted">{needsPractice} Need practice</span>}
              {flat.length === 0 && <span className="text-mathua-muted">No struggles yet</span>}
            </div>
            <div className="font-mono text-[10px] text-mathua-muted truncate">{flat.length} concept{flat.length !== 1 ? 's' : ''} · {domains.length} domain{domains.length !== 1 ? 's' : ''}</div>
          </div>
          <Link
            href="/session"
            className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none px-3 h-8 font-mono text-[11px] inline-flex items-center justify-center min-h-[32px]"
          >
            Review now →
          </Link>
        </div>

        {/* Top 3 preview */}
        {top.length > 0 && (
          <div className="px-3 pb-2 space-y-1.5">
            {top.map(item => {
              const tier = tierFor(item.weakness)
              return (
                <div key={item.id} className="flex min-w-0 items-center gap-2 py-0.5">
                  <div className="w-10 h-1 bg-mathua-border shrink-0">
                    <div className="h-full transition-all duration-300" style={{ width: `${Math.round(item.weakness * 100)}%`, background: tier.color }} />
                  </div>
                  <span className="min-w-0 flex-1 truncate text-mathua-primary text-[11px] font-mono">{item.label}</span>
                  <span className="shrink-0 font-mono text-[10px]" style={{ color: tier.color }}>{tier.label}</span>
                </div>
              )
            })}
          </div>
        )}

        {/* View all toggle */}
        <div className="px-3 pb-3">
          <button
            onClick={() => setExpanded(v => !v)}
            aria-expanded={expanded}
            className="font-mono text-[11px] text-mathua-blue hover:text-mathua-blue-hover underline underline-offset-4"
          >
            {expanded ? 'Hide details' : `View all — ${flat.length} items`}
          </button>
        </div>

        {/* Expanded per-domain collapsible list */}
        {expanded && (
          <div className="border-t-[0.5px] border-mathua-border font-mono text-[11px] w-full max-w-full min-w-0 overflow-hidden">
            {Object.keys(entries)
              .sort((a, b) => {
                const aw = Math.max(...(entries[a] as { weakness: number }[]).map(x => x.weakness))
                const bw = Math.max(...(entries[b] as { weakness: number }[]).map(x => x.weakness))
                return bw - aw
              })
              .map(domain => {
                const items = (entries[domain] as { id: string; label: string; weakness: number }[]).slice().sort((a, b) => b.weakness - a.weakness)
                const isOpen = openDomains.has(domain)
                const maxW = Math.max(...items.map(i => i.weakness))
                const tier = tierFor(maxW)
                return (
                  <div key={domain} className="border-b-[0.5px] border-mathua-border last:border-b-0">
                    <button
                      onClick={() => toggleDomain(domain)}
                      aria-expanded={isOpen}
                      className="w-full flex items-center gap-2 px-3 py-2 min-h-[44px] text-left hover:bg-mathua-surface-elevated transition-colors"
                    >
                      <span className="w-10 h-1 bg-mathua-border shrink-0">
                        <span className="block h-full transition-all" style={{ width: `${Math.round(maxW * 100)}%`, background: tier.color }} />
                      </span>
                      <span className="flex-1 min-w-0 truncate text-mathua-primary text-[11px]">{DOMAIN_LABELS[domain] ?? domain}</span>
                      <span className="shrink-0 text-mathua-muted text-[10px]">{items.length} weak</span>
                      <span className="shrink-0 text-mathua-muted text-[11px]">{isOpen ? '▾' : '▸'}</span>
                    </button>
                    {isOpen && (
                      <div className="px-3 pb-2 space-y-0.5 bg-mathua-surface-elevated/50">
                        {items.slice(0, 8).map(item => {
                          const t = tierFor(item.weakness)
                          return (
                            <div key={item.id} className="flex min-w-0 items-center gap-2 py-1">
                              <div className="w-8 h-1 bg-mathua-border shrink-0">
                                <div className="h-full" style={{ width: `${Math.round(item.weakness * 100)}%`, background: t.color }} />
                              </div>
                              <span className="min-w-0 flex-1 truncate text-mathua-secondary text-[11px]">{item.label}</span>
                              <span className="shrink-0 font-mono text-[10px]" style={{ color: t.color }}>{Math.round(item.weakness * 100)}%</span>
                            </div>
                          )
                        })}
                        {items.length > 8 && (
                          <div className="pt-1">
                            <Link href={`/study?domain=${encodeURIComponent(domain)}`} className="font-mono text-[10px] text-mathua-blue hover:underline">
                              +{items.length - 8} more — practice in Study →
                            </Link>
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                )
              })}
          </div>
        )}
      </div>
    </div>
  )
}
