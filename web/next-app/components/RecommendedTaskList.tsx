'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { hrefConceptId, type Shelf, type ShelfItem } from '../lib/nextUp'
import { getLessonEligibility, type EligibilityPrereq } from '../lib/api'
import PrerequisitePanel from './PrerequisitePanel'

/**
 * The adaptive task set, presented as a small group of equally actionable cards.
 *
 * The set comes from the engine (`/api/next` via `fetchShelf`); this component
 * never ranks or invents a task. Selecting one asks the server for the
 * eligibility verdict and then enters the one canonical learning loop at
 * `/learn?concept=`. A locked topic is explained, not started: the choice is
 * reachability only and cannot stand in for demonstrated knowledge (invariant 3).
 *
 * It lives on the profile and is the same component `/learn` can reuse, so a
 * recommended concept and a manually chosen one converge on one flow.
 */
function xpLabel(xp: number): string {
  return xp > 0 ? `+${xp} XP` : ''
}

export function taskItems(shelf: Shelf | null, cap = 5): ShelfItem[] {
  if (!shelf) return []
  const seen = new Set<string>()
  const out: ShelfItem[] = []
  for (const item of [shelf.next, ...shelf.alternatives]) {
    if (seen.has(item.href)) continue
    seen.add(item.href)
    out.push(item)
    if (out.length >= cap) break
  }
  return out
}

interface Blocked {
  conceptId: string
  title: string
  prerequisites: EligibilityPrereq[]
}

export default function RecommendedTaskList({
  shelf,
  loading = false,
}: {
  shelf: Shelf | null
  loading?: boolean
}) {
  const router = useRouter()
  const [busy, setBusy] = useState<string | null>(null)
  const [blocked, setBlocked] = useState<Blocked | null>(null)
  const [error, setError] = useState('')

  async function choose(item: ShelfItem) {
    const conceptId = hrefConceptId(item.href)
    if (!conceptId) {
      // Review / mastery check / orientation: the server already worded the
      // destination, so it is followed as-is.
      router.push(item.href)
      return
    }
    setBusy(conceptId)
    setError('')
    setBlocked(null)
    try {
      const eligibility = await getLessonEligibility(conceptId)
      if (eligibility.eligible === false) {
        setBlocked({
          conceptId,
          title: item.title,
          prerequisites: eligibility.prerequisites ?? [],
        })
        return
      }
      router.push(`/learn?concept=${encodeURIComponent(conceptId)}`)
    } catch {
      setError("Couldn't start this topic.")
    } finally {
      setBusy(null)
    }
  }

  const items = taskItems(shelf)

  return (
    <section aria-label="Your learning" className="mt-8 min-w-0">
      <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-1">
        Your learning
      </h2>
      <p className="font-mono text-[11px] text-mathua-muted mb-4">Recommended tasks</p>

      {loading && (
        <div
          role="status"
          aria-live="polite"
          className="border border-mathua-border bg-mathua-surface p-4 font-mono text-xs text-mathua-muted"
        >
          Loading recommended tasks…
        </div>
      )}

      {!loading && items.length === 0 && (
        <div className="border border-mathua-border bg-mathua-surface p-4 min-w-0">
          <p className="font-mono text-xs text-mathua-secondary">
            No learning tasks are available right now.
          </p>
          <Link
            href="/domains"
            className="mt-3 inline-flex border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-xs min-h-[36px] items-center"
          >
            Choose what to learn
          </Link>
        </div>
      )}

      {!loading && items.length > 0 && (
        <ul className="grid grid-cols-1 sm:grid-cols-3 gap-3 min-w-0">
          {items.map((item, i) => {
            const conceptId = hrefConceptId(item.href)
            const isBusy = conceptId !== null && busy === conceptId
            return (
              <li key={`${item.kind}:${item.href}:${i}`} className="min-w-0">
                <button
                  type="button"
                  onClick={() => void choose(item)}
                  disabled={isBusy}
                  className="w-full h-full text-left border border-mathua-border bg-mathua-surface p-4 hover:border-mathua-blue transition-colors min-w-0 disabled:opacity-60"
                >
                  <span className="flex items-center justify-between gap-2 min-w-0">
                    <span className="shrink-0 bg-mathua-blue-faint text-mathua-blue px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
                      {item.badge}
                    </span>
                    {xpLabel(item.xp) && (
                      <span className="shrink-0 font-mono text-[10px] text-yellow-400">{xpLabel(item.xp)}</span>
                    )}
                  </span>
                  <span className="mt-2 block font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere] leading-snug">
                    {item.title}
                  </span>
                  <span className="mt-1 block font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
                    {item.detail}
                  </span>
                  <span className="mt-2 block font-mono text-[11px] text-mathua-blue">
                    {isBusy ? 'Starting…' : item.cta}
                  </span>
                </button>
              </li>
            )
          })}
        </ul>
      )}

      {error && (
        <p role="alert" className="mt-3 font-mono text-[11px] text-mathua-secondary">
          {error}{' '}
          <button type="button" onClick={() => setError('')} className="text-mathua-blue hover:underline">
            Try again
          </button>
        </p>
      )}

      {blocked && (
        <PrerequisitePanel
          topicTitle={blocked.title}
          prerequisites={blocked.prerequisites}
          onDismiss={() => setBlocked(null)}
        />
      )}

      <Link
        href="/domains"
        className="mt-4 inline-flex font-mono text-xs text-mathua-blue hover:text-mathua-blue-hover"
      >
        Choose what to learn →
      </Link>
    </section>
  )
}
