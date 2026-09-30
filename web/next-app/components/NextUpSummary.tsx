'use client'

import Link from 'next/link'
import type { Shelf } from '../lib/nextUp'

/**
 * NextUpSummary: Profile's compact view of the ranked shelf. It names the
 * Next head item and links to /learn, which owns the full NextUpCard.
 * No shelf, no alternatives here — Learn is the only place that proposes.
 */
export default function NextUpSummary({ shelf }: { shelf: Shelf }) {
  const { next } = shelf
  return (
    <section aria-label="Next up" className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-border bg-mathua-surface p-4">
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <div className="min-w-0">
          <div className="flex items-center gap-2 min-w-0">
            <span className="shrink-0 bg-mathua-blue/10 text-mathua-blue px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
              {next.badge}
            </span>
            <span className="font-mono text-[10px] text-mathua-muted uppercase tracking-wider">Next up</span>
          </div>
          <p className="mt-2 font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere] leading-snug">
            {next.title}
          </p>
          <p className="font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
            {next.detail}
          </p>
        </div>
        <Link
          href="/learn"
          className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue/10 px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
        >
          Open Learn →
        </Link>
      </div>
    </section>
  )
}
