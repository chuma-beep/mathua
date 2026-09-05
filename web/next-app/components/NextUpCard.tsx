'use client'

import Link from 'next/link'
import type { NextUp } from '../lib/nextUp'

/** Full-width hero card for #next; compact=true renders the sidebar slot variant. */
export default function NextUpCard({ next, compact = false }: { next: NextUp; compact?: boolean }) {
  if (compact) {
    return (
      <div className="border border-mathua-blue bg-mathua-surface p-3 min-w-0">
        <span className="inline-block bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
          {next.badge}
        </span>
        <div className="font-mono text-xs text-mathua-primary mt-2 break-words [overflow-wrap:anywhere] leading-snug">
          {next.title}
        </div>
        <Link
          href={next.href}
          className="mt-2 w-full border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 py-2 font-mono text-xs min-h-[44px] inline-flex items-center justify-center text-center whitespace-nowrap transition-colors"
        >
          {next.cta}
        </Link>
      </div>
    )
  }

  return (
    <section id="next" aria-label="Next up" className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-blue bg-mathua-surface p-4 scroll-mt-28">
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <div className="min-w-0 flex-1">
          <div className="flex flex-col items-start gap-1.5 sm:flex-row sm:items-center sm:gap-2 min-w-0">
            <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
              {next.badge}
            </span>
            <span className="min-w-0 break-words [overflow-wrap:anywhere] leading-snug font-mono text-[11px] sm:text-xs text-mathua-primary">
              {next.title}
            </span>
          </div>
          <p className="font-mono text-xs text-mathua-secondary mt-1 break-words [overflow-wrap:anywhere]">
            {next.detail}
          </p>
        </div>
        <Link
          href={next.href}
          className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[44px] inline-flex items-center justify-center text-center whitespace-nowrap transition-colors"
        >
          {next.cta}
        </Link>
      </div>
    </section>
  )
}
