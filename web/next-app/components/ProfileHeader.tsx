'use client'

import type { ReactNode } from 'react'
import Link from 'next/link'
import type { ConceptProgress, Scores } from '../lib/api'
import { countOverall, type CatalogueConcept } from '../lib/progress'
import { LEVEL_NAMES } from './ProfileStats'
import Avatar from './Avatar'

/**
 * The learner header as a few short sentences, not a wall of tiles.
 *
 * It replaced a bordered stat strip (avatar, LV tile, XP/TOTAL/STREAK/WEEKLY tiles, one progress
 * bar) stacked on a position card (big number, second bar, 3-column grid). Every fact below used
 * to have its own tile and label; now each line carries several, which is the whole reduction.
 * The weekly leaderboard link was cut as part of it, because it is navigation, not a position.
 * All-time XP stays on the reach line because it is the learner's cumulative work, not a
 * same-day vanity number.
 *
 * What the rewrite must not lose (all pinned by tests): the `Where you are` region, the literal
 * `of {n} you can reach`, the all-time XP value, the review-debt line with its `/review` link and
 * the frontier sentence. `ProfileStats` is untouched and still serves `/share`.
 */
export default function ProfileHeader({
  name,
  scores,
  avatarSeed,
  avatarUrl,
  avatarPreset,
  catalogue,
  progress,
  frontierLabel,
  footer,
}: {
  name: string
  scores: Scores
  avatarSeed?: string
  avatarUrl?: string
  avatarPreset?: number | null
  catalogue: CatalogueConcept[]
  progress: Record<string, ConceptProgress>
  frontierLabel?: string | null
  /** Small actions that belong on the header line rather than as their own blocks. */
  footer?: ReactNode
}) {
  const counts = countOverall(catalogue, progress)
  if (counts.total === 0) return null
  const levelNum = LEVEL_NAMES[scores.level] ?? '·'

  return (
    <div className="border-[0.5px] border-mathua-border p-4 sm:px-5 w-full max-w-full min-w-0 overflow-hidden">
      <div className="flex items-start gap-4 min-w-0">
        {/* The profile picture, square and boxed, at the left of the header. */}
        <div className="shrink-0 border-[0.5px] border-mathua-border p-1">
          <Avatar
            seed={avatarSeed ?? name}
            name={name}
            size={72}
            url={avatarUrl}
            preset={avatarPreset ?? undefined}
            className="!rounded-none"
          />
        </div>
        <section id="position" aria-label="Where you are" className="min-w-0 flex-1">
          <p className="font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere]">
            {name} · LV {levelNum} {scores.level}
          </p>
          <p className="mt-1.5 font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
            {counts.mastered} of {counts.unlocked} you can reach · {scores.xp_total.toLocaleString()} XP all-time · {scores.xp_today}/{scores.daily_xp_goal} XP today (net) · {scores.current_streak} day{scores.current_streak !== 1 ? 's' : ''} streak
          </p>
          <p className="mt-1.5 font-mono text-[11px] text-mathua-muted break-words [overflow-wrap:anywhere]">
            Mastered {counts.mastered} · Unlocked {counts.unlocked} · In progress {counts.learning + counts.completed}
            {frontierLabel ? ` · ${counts.mastered > 0 ? 'Currently working towards' : 'Start here'}: ${frontierLabel}` : ''}
          </p>
          {counts.dueForReview > 0 && (
            <p className="mt-1.5 font-mono text-[11px] text-mathua-secondary">
              {counts.dueForReview} due for review, still learned, worth a retrieval check.{' '}
              <Link href="/review" className="text-mathua-blue hover:underline">
                Review →
              </Link>
            </p>
          )}
        </section>
      </div>
      {footer && (
        <div className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-2 font-mono text-[11px]">
          {footer}
        </div>
      )}
    </div>
  )
}
