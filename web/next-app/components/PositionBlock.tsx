'use client'

import Link from 'next/link'
import { countOverall, type CatalogueConcept } from '../lib/progress'
import type { ConceptProgress } from '../lib/api'

interface Props {
  catalogue: CatalogueConcept[]
  progress: Record<string, ConceptProgress>
  /** Label of the concept the scheduler would serve next, if known. */
  frontierLabel?: string | null
  frontierHref?: string | null
}

// "Where am I" — the question the rest of the report hangs off, and the one
// nothing answered before. `concepts_mastered` was a bare number in a stat
// strip and the only other measure was per-domain percentages.
//
// Deliberately not a percentage of 657. The catalogue runs to topology and
// abstract algebra, so "7% of all of math" tells a learner nothing they can
// act on. Mastered against unlocked is the pair that means something: one is
// what you have, the other is what you can currently reach.
export default function PositionBlock({ catalogue, progress, frontierLabel, frontierHref }: Props) {
  const counts = countOverall(catalogue, progress)
  if (counts.total === 0) return null

  const pct = counts.total > 0 ? Math.round((counts.mastered / counts.total) * 100) : 0

  return (
    <section id="position" aria-label="Where you are" className="min-w-0">
      <h2 className="font-serif text-[1.05rem] font-normal text-mathua-primary mb-3">
        Where you are
      </h2>

      <div className="border border-mathua-border bg-mathua-surface p-4 min-w-0">
        <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1">
          <span className="font-mono text-3xl text-mathua-blue leading-none">{counts.mastered}</span>
          <span className="font-mono text-sm text-mathua-muted">
            of {counts.total} concepts mastered
          </span>
        </div>

        {/* The bar is against unlocked, not against the whole catalogue: it is
            the part of the graph the learner can actually walk right now. */}
        <div className="mt-3 h-1 w-full bg-mathua-border" role="presentation">
          <div
            className="h-full bg-mathua-blue transition-all duration-300"
            style={{ width: `${counts.unlocked > 0 ? Math.round((counts.mastered / counts.unlocked) * 100) : 0}%` }}
          />
        </div>

        <dl className="mt-3 grid grid-cols-2 sm:grid-cols-4 gap-3 font-mono text-[11px]">
          <div className="min-w-0">
            <dt className="text-mathua-muted uppercase">Mastered</dt>
            <dd className="text-mathua-primary text-base">{counts.mastered}</dd>
          </div>
          <div className="min-w-0">
            <dt className="text-mathua-muted uppercase">Unlocked</dt>
            <dd className="text-mathua-primary text-base">{counts.unlocked}</dd>
          </div>
          <div className="min-w-0">
            <dt className="text-mathua-muted uppercase">In progress</dt>
            <dd className="text-mathua-primary text-base">{counts.learning + counts.completed}</dd>
          </div>
          <div className="min-w-0">
            <dt className="text-mathua-muted uppercase">Still locked</dt>
            <dd className="text-mathua-secondary text-base">{counts.locked}</dd>
          </div>
        </dl>

        {counts.mastered > 0 && (
          <p className="mt-3 font-mono text-[11px] text-mathua-secondary">
            {pct}% of everything you have access to so far.
          </p>
        )}

        {counts.dueForReview > 0 && (
          // Named rather than folded into "Mastered", because it is a different fact with
          // a different consequence. These concepts are counted as mastered — the bar
          // above is honest about that — and this line is the review debt, so the number
          // the learner acts on is visible instead of having been subtracted from a
          // headline they would then compare against the graph's full bars.
          <p className="mt-2 font-mono text-[11px] text-mathua-secondary">
            {counts.dueForReview} due for review — still learned, worth a retrieval check.{' '}
            <Link href="/review" className="text-mathua-blue hover:underline">
              Review →
            </Link>
          </p>
        )}

        {frontierLabel && (
          <div className="mt-4 border-t border-mathua-border pt-3 flex flex-wrap items-center justify-between gap-2">
            <span className="min-w-0">
              <span className="font-mono text-[10px] uppercase text-mathua-muted block">
                {counts.mastered > 0 ? 'Currently working towards' : 'Start here'}
              </span>
              <span className="font-mono text-xs text-mathua-primary break-words">
                {frontierLabel}
              </span>
            </span>
            {frontierHref && (
              <Link
                href={frontierHref}
                className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-3 py-1.5 font-mono text-[11px] inline-flex items-center min-h-[36px]"
              >
                Answer it →
              </Link>
            )}
          </div>
        )}
      </div>
    </section>
  )
}