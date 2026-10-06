import { getNext, type Recommendation, type RecommendationRes } from './api'
import type { Shelf, ShelfItem } from './nextUp'

/**
 * The engine's recommendation, in the shape the surfaces already render.
 *
 * This is a presentation adapter and nothing else. The decision — which concept, in what order,
 * why, and with what wording — is made once, in `internal/scheduler`, and arrives over
 * `/api/next`. It used to be made here, in the browser, from five API calls plus a bundled copy
 * of the corpus, while the scheduler sat unreachable between `/learn` and `/profile`.
 *
 * The mapping exists so the components that render a recommendation do not each grow their own
 * interpretation of the payload. `Badge`, `Detail`, `CTA` and `Href` pass straight through:
 * they are the server's copy, and rewriting them here is how two surfaces came to describe one
 * head differently.
 */
export function recommendationToShelfItem(r: Recommendation): ShelfItem {
  return {
    kind: shelfKind(r),
    badge: r.badge,
    title: r.conceptTitle || r.id,
    detail: r.detail,
    href: r.action.href,
    cta: r.cta,
    xp: r.xp ?? 0,
  }
}

/**
 * The client's kind vocabulary is finer than the wire kind: it distinguishes resuming a
 * concept from starting one, which the surfaces render differently ("Practice again" versus
 * "Start"). `reason` carries that distinction, so it maps rather than guessing.
 */
function shelfKind(r: Recommendation): ShelfItem['kind'] {
  if (r.kind === 'review') return 'review'
  if (r.kind === 'practice') return 'weakness'
  if (r.kind === 'mastery_check') return 'diagnostic'
  if (r.id === 'diagnostic') return 'diagnostic'
  if (r.id === 'study') return 'browse'
  return r.reason === 'prerequisite' ? 'resume' : 'new'
}

export function recommendationResToShelf(res: RecommendationRes): Shelf {
  const primary = res.primary
    ? recommendationToShelfItem(res.primary)
    : {
        kind: 'browse' as const,
        badge: 'Library',
        title: 'Lesson library',
        detail: 'Nothing is recommended right now',
        href: '/study',
        cta: 'Browse lessons →',
        xp: 0,
      }
  return {
    next: primary,
    // Defaulted here as well as in getNext: this function is exported, and every caller
    // iterates the result without a null check. A proxy or an older build that omits the field
    // must not turn the home page into a render error.
    alternatives: (res.alternatives ?? []).map(recommendationToShelfItem),
  }
}

/**
 * Fetch the engine's answer and adapt it.
 *
 * `exclude` drops concepts the caller has just finished, so a concept is never offered back as
 * the next thing to do — which is why the server takes the parameter rather than the client
 * filtering afterwards.
 */
export async function fetchShelf(exclude?: string[]): Promise<Shelf> {
  return recommendationResToShelf(await getNext(exclude))
}

/** The raw response, for callers that need `reason` or `priority` rather than a shelf item. */
export async function fetchRecommendations(exclude?: string[]): Promise<RecommendationRes> {
  return getNext(exclude)
}
