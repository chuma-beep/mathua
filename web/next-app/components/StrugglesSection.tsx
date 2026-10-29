'use client'

import Link from 'next/link'
import type { WeaknessRes } from '../lib/api'

interface Props {
  weaknesses: WeaknessRes | null
}

// A one-line "needs attention" summary: how many concepts are weak, and the way to practise them.
// It was a boxed panel with a per-concept bar list and a per-domain expander, the heaviest block
// on the page and every number in it was decided by the same count this now states outright.
// This is a weakness count, not the spaced-review due count, so the action must open practice.
export default function StrugglesSection({ weaknesses }: Props) {
  const entries = weaknesses?.by_domain ?? {}
  const domains = Object.keys(entries)
  if (domains.length === 0) return null

  let struggling = 0
  let total = 0
  for (const items of Object.values(entries)) {
    for (const it of items) {
      total++
      if (it.weakness >= 0.6) struggling++
    }
  }
  const needsPractice = total - struggling

  return (
    <div className="min-w-0">
      <h3 className="font-mono text-[11px] uppercase tracking-wider text-mathua-muted mb-2">
        Needs attention
      </h3>
      <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-2 border-[0.5px] border-mathua-border px-3 py-2 min-w-0">
        <span className="font-mono text-[11px] min-w-0">
          {struggling > 0 && <span className="text-mathua-red">{struggling} struggling</span>}
          {struggling > 0 && needsPractice > 0 && <span className="text-mathua-muted"> · </span>}
          {needsPractice > 0 && <span className="text-mathua-muted">{needsPractice} need practice</span>}
          <span className="text-mathua-muted">
            {total > 0 ? ' · ' : ''}{total} concept{total !== 1 ? 's' : ''} · {domains.length} domain{domains.length !== 1 ? 's' : ''}
          </span>
        </span>
        <Link
          href="/learn"
          className="shrink-0 font-mono text-[11px] text-mathua-blue hover:text-mathua-blue-hover whitespace-nowrap"
        >
          Practice now →
        </Link>
      </div>
    </div>
  )
}
