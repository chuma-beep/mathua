'use client'

import Link from 'next/link'
import type { EligibilityPrereq } from '../lib/api'

/**
 * The prerequisite explanation shown when a learner selects a topic the server
 * has not cleared them for. It names what is missing and routes the learner into
 * the learning system for a prerequisite, never into the reference library:
 * eligibility is proved by answering, not by reading (ADR-021, ADR-045, ADR-047).
 *
 * Every link is a function of the server's own state for that prerequisite, so
 * a locked prerequisite is shown as locked rather than offered as a start.
 */
export function prereqHref(p: EligibilityPrereq): string | null {
  switch (p.state) {
    case 'mastered':
    case 'due_for_review':
      return '/review'
    case 'unlocked':
    case 'in_progress':
      return `/learn?concept=${encodeURIComponent(p.conceptId)}`
    default:
      return null // locked / upcoming: not startable yet
  }
}

export function prereqCta(p: EligibilityPrereq): string {
  switch (p.state) {
    case 'mastered':
    case 'due_for_review':
      return 'Review prerequisite →'
    case 'unlocked':
      return 'Learn prerequisite →'
    case 'in_progress':
      return 'Continue prerequisite →'
    case 'upcoming':
      return 'Not available yet'
    default:
      return 'Locked'
  }
}

export default function PrerequisitePanel({
  topicTitle,
  prerequisites,
  onDismiss,
}: {
  topicTitle: string
  prerequisites: EligibilityPrereq[]
  onDismiss?: () => void
}) {
  const unmet = prerequisites.filter((p) => !p.met)
  const shown = unmet.length > 0 ? unmet : prerequisites
  return (
    <section
      aria-label="Prerequisite explanation"
      role="note"
      className="mt-4 border border-mathua-border bg-mathua-surface p-4 min-w-0"
    >
      <h3 className="font-mono text-xs text-mathua-primary break-words [overflow-wrap:anywhere]">
        {topicTitle}
      </h3>
      <p className="mt-2 font-mono text-[11px] text-mathua-secondary break-words [overflow-wrap:anywhere]">
        You&apos;re not ready for this topic yet. Mathua recommends learning the
        {shown.length === 1 ? ' prerequisite' : ' prerequisites'} below first.
      </p>
      <ul className="mt-3 space-y-2">
        {shown.map((p) => {
          const href = prereqHref(p)
          return (
            <li key={p.conceptId} className="flex flex-col sm:flex-row sm:items-center gap-2 min-w-0">
              <span className="min-w-0 flex-1">
                <span className="block font-mono text-[11px] text-mathua-primary break-words [overflow-wrap:anywhere]">
                  {p.title}
                </span>
                <span className="block font-mono text-[10px] uppercase tracking-wider text-mathua-muted">
                  {p.met ? 'Done' : p.state.replace(/_/g, ' ')}
                </span>
              </span>
              {href ? (
                <Link
                  href={href}
                  className="w-full sm:w-auto shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap"
                >
                  {prereqCta(p)}
                </Link>
              ) : (
                <span className="w-full sm:w-auto shrink-0 border border-mathua-border text-mathua-muted px-4 py-2 font-mono text-[11px] min-h-[36px] inline-flex items-center justify-center text-center whitespace-nowrap">
                  {prereqCta(p)}
                </span>
              )}
            </li>
          )
        })}
      </ul>
      {onDismiss && (
        <button
          type="button"
          onClick={onDismiss}
          className="mt-3 font-mono text-[10px] uppercase tracking-wider text-mathua-muted hover:text-mathua-blue"
        >
          Close
        </button>
      )}
    </section>
  )
}
