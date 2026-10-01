'use client'

import Link from 'next/link'
import type { Shelf } from '../lib/nextUp'

/**
 * NextUpCard: one primary "Next up" head card, with the remaining shelf
 * items collapsed under "Or pick something else". The scheduler stays the
 * source of truth; the UI just commits to item 0 as the head. Learner
 * disposes (ADR-017) via the alternatives or by ignoring the head.
 */
export default function NextUpCard({ shelf }: { shelf: Shelf }) {
  const { next, alternatives } = shelf
  return (
    <section aria-label="Next up" className="mt-6 w-full max-w-full min-w-0 overflow-hidden">
      <h2 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">
        Next up
      </h2>
      <Link
        href={next.href}
        className="border-2 border-mathua-blue bg-mathua-surface p-4 hover:bg-mathua-blue-faint transition-colors block min-w-0"
      >
        <div className="flex items-center justify-between gap-2 min-w-0">
          <span className="shrink-0 bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
            {next.badge}
          </span>
          {next.xp > 0 && (
            <span className="shrink-0 font-mono text-[10px] text-yellow-400">{next.xp} XP</span>
          )}
        </div>
        <div className="mt-2 font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere] leading-snug">
          {next.title}
        </div>
        <div className="mt-1 font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
          {next.detail}
        </div>
        <div className="mt-2 font-mono text-[11px] text-mathua-blue">{next.cta}</div>
      </Link>
      {alternatives.length > 0 && (
        <details className="mt-3 border border-mathua-border bg-mathua-surface">
          <summary className="font-mono text-[11px] text-mathua-secondary cursor-pointer px-4 py-3">
            Or pick something else ({alternatives.length})
          </summary>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 px-4 pb-4">
            {alternatives.map((it, i) => (
              <Link
                key={`${it.kind}:${it.href}:${i}`}
                href={it.href}
                className="border border-mathua-border bg-mathua-surface p-4 hover:border-mathua-blue transition-colors block min-w-0"
              >
                <div className="flex items-center justify-between gap-2 min-w-0">
                  <span className="shrink-0 bg-mathua-blue-faint text-mathua-blue px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                    {it.badge}
                  </span>
                  {it.xp > 0 && (
                    <span className="shrink-0 font-mono text-[10px] text-yellow-400">{it.xp} XP</span>
                  )}
                </div>
                <div className="mt-2 font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere] leading-snug">
                  {it.title}
                </div>
                <div className="mt-1 font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
                  {it.detail}
                </div>
                <div className="mt-2 font-mono text-[11px] text-mathua-blue">{it.cta}</div>
              </Link>
            ))}
          </div>
        </details>
      )}
    </section>
  )
}
