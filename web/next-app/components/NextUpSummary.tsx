'use client'

import Link from 'next/link'
import { hrefConceptId, type Shelf, type ShelfItem } from '../lib/nextUp'

/**
 * NextUpSummary: Profile's dashboard view of the ranked shelf.
 * - Item 0 renders as the ONE primary button: "Continue: {label}" with the
 *   reason badge and estimated XP. It links to the head href as-is
 *   (/learn?concept= for learn items; /review, /onboard, /study otherwise).
 * - Alternatives render as quiet rows (reason + label + XP), visually
 *   secondary. Only the head is a button.
 * - Head href, head concept, and duplicate hrefs are filtered from the rows,
 *   so no two items share a destination.
 */
function queueRows(next: ShelfItem, alternatives: ShelfItem[]): ShelfItem[] {
  const headCid = hrefConceptId(next.href)
  const seen = new Set([next.href])
  const rows: ShelfItem[] = []
  for (const a of alternatives) {
    if (seen.has(a.href)) continue
    const cid = hrefConceptId(a.href)
    if (cid !== null && cid === headCid) continue
    seen.add(a.href)
    rows.push(a)
    if (rows.length >= 4) break
  }
  return rows
}

export function xpLabel(xp: number): string {
  return xp > 0 ? `+${xp} XP` : ''
}

export default function NextUpSummary({ shelf }: { shelf: Shelf }) {
  const { next } = shelf
  const rows = queueRows(next, shelf.alternatives)
  return (
    <section aria-label="Next up" className="mt-6 w-full max-w-full min-w-0 overflow-hidden border border-mathua-border bg-mathua-surface p-4">
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3">
        <div className="min-w-0">
          <div className="flex items-center gap-2 min-w-0">
            <span className="shrink-0 bg-mathua-blue-faint text-mathua-blue px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
              {next.badge}
            </span>
            <span className="font-mono text-[10px] text-mathua-muted uppercase tracking-wider">Next up</span>
          </div>
          <p className="mt-2 font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere] leading-snug">
            Continue: {next.title}
          </p>
          <p className="font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
            {next.detail}{xpLabel(next.xp) ? ` · ${xpLabel(next.xp)}` : ''}
          </p>
        </div>
        <Link
          href={next.href}
          className="w-full sm:w-auto sm:shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
        >
          {next.cta}
        </Link>
      </div>
      {rows.length > 0 && (
        <div className="mt-3 border-t border-mathua-border pt-3">
          <p className="font-mono text-[10px] text-mathua-muted uppercase tracking-wider mb-2">Today&apos;s queue</p>
          <ul className="space-y-1.5">
            {rows.map((r) => (
              <li key={`${r.kind}:${r.href}`}>
                <Link
                  href={r.href}
                  className="flex items-baseline gap-2 font-mono text-[11px] text-mathua-secondary hover:text-mathua-blue min-w-0"
                >
                  <span className="shrink-0 uppercase tracking-wider text-mathua-muted">{r.badge}</span>
                  <span className="truncate">{r.title}</span>
                  {xpLabel(r.xp) && <span className="ml-auto shrink-0 text-mathua-blue">{xpLabel(r.xp)}</span>}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      )}
    </section>
  )
}
