import type { ConceptProgress, DailyActivity, WeaknessRes } from './api'

export type NextUpKind = 'review' | 'weakness' | 'resume' | 'diagnostic' | 'browse'

export interface NextUp {
  kind: NextUpKind
  badge: string
  title: string
  detail: string
  href: string
  cta: string
}

export interface NextUpInput {
  dueReviews: number
  weaknesses: WeaknessRes | null
  progress: Record<string, ConceptProgress>
  activity: DailyActivity[]
  diagnosticCompleted: boolean
  conceptsMastered: number
}

function isMastered(p: ConceptProgress | undefined): boolean {
  return (p?.status ?? '').toLowerCase() === 'mastered'
}

function topWeakness(input: NextUpInput): { id: string; label: string } | null {
  const entries = input.weaknesses?.by_domain ?? {}
  let best: { id: string; label: string; weakness: number } | null = null
  for (const items of Object.values(entries)) {
    for (const it of items) {
      if (isMastered(input.progress[it.id])) continue
      if (!best || it.weakness > best.weakness) best = it
    }
  }
  return best
}

function mostRecentInProgress(input: NextUpInput): string | null {
  for (let i = input.activity.length - 1; i >= 0; i--) {
    const day = input.activity[i]
    for (let j = day.concepts.length - 1; j >= 0; j--) {
      const cid = day.concepts[j]
      const p = input.progress[cid]
      // In progress = has a progress record but not yet mastered
      if (p && !isMastered(p)) return cid
    }
  }
  return null
}

/**
 * Deterministic Next-up selector. Priority:
 * 1) due reviews → review session host
 * 2) weakest non-mastered concept → Study deep link (supported ?concept= param)
 * 3) most recently touched in-progress concept → Study deep link (resume)
 * 4) brand-new user without Diagnostic → Diagnostic entry
 * 5) fallback → Study library
 */
export function selectNextUp(input: NextUpInput): NextUp {
  if (input.dueReviews > 0) {
    return {
      kind: 'review',
      badge: 'Due now',
      title: `${input.dueReviews} concept${input.dueReviews !== 1 ? 's' : ''} due for review`,
      detail: 'Spaced repetition — review before decay',
      href: '/session',
      cta: 'Review now →',
    }
  }

  const weak = topWeakness(input)
  if (weak) {
    return {
      kind: 'weakness',
      badge: 'Recommended',
      title: `Study ${weak.label}`,
      detail: 'Weakest concept — targeted Study session',
      href: `/study?concept=${encodeURIComponent(weak.id)}`,
      cta: 'Open in Study →',
    }
  }

  const resume = mostRecentInProgress(input)
  if (resume) {
    return {
      kind: 'resume',
      badge: 'Continue',
      title: `Resume ${resume}`,
      detail: 'Pick up where you left off',
      href: `/study?concept=${encodeURIComponent(resume)}`,
      cta: 'Continue →',
    }
  }

  if (input.conceptsMastered === 0 && !input.diagnosticCompleted) {
    return {
      kind: 'diagnostic',
      badge: 'Recommended',
      title: 'Take a Diagnostic to find your frontier',
      detail: 'Adaptive — finds where to start',
      href: '/onboard',
      cta: 'Start Diagnostic →',
    }
  }

  return {
    kind: 'browse',
    badge: 'Study',
    title: 'Browse the Study library',
    detail: 'Pick a Lesson — worked example first, then answer',
    href: '/study',
    cta: 'Browse Study →',
  }
}
