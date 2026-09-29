'use client'

import Link from 'next/link'
import type { ShelfItem } from '../lib/nextUp'

/** Task shelf: the algorithm proposes eligible tasks, the learner disposes. */
export default function TaskShelf({ items }: { items: ShelfItem[] }) {
  if (items.length === 0) return null
  return (
    <section aria-label="What do you want to work on" className="mt-6 w-full max-w-full min-w-0 overflow-hidden">
      <h2 className="font-mono text-[11px] text-mathua-muted uppercase tracking-wider mb-3">
        What do you want to work on?
      </h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {items.map((it, i) => (
          <Link
            key={`${it.kind}:${it.href}:${i}`}
            href={it.href}
            className="border border-mathua-border bg-mathua-surface p-4 hover:border-mathua-blue transition-colors block min-w-0"
          >
            <div className="flex items-center justify-between gap-2 min-w-0">
              <span className="shrink-0 bg-mathua-blue/10 text-mathua-blue px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
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
    </section>
  )
}
