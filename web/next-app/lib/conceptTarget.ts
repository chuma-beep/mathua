// Concept-chip target resolution for Study LessonDetail.
//
// One behavior, not two: the URL is the state (`/study?concept=X`), and an
// effect reveals the block. Chips use `replace` for an in-lesson jump (no
// back-button history spam) and `push` when the resolved lesson changes.
// Unknown concepts resolve cross-lesson so the URL sync can no-op gracefully.
export interface ConceptTarget {
  href: string
  sameLesson: boolean
}

export function resolveConceptTarget(
  cid: string,
  lessonByConcept: Map<string, { title: string }>,
  currentTitle: string | null,
): ConceptTarget {
  const href = `/study?concept=${encodeURIComponent(cid)}`
  const target = lessonByConcept.get(cid)
  return { href, sameLesson: !!target && !!currentTitle && target.title === currentTitle }
}

// One navigation plan per chip click. Re-clicking the concept already in
// the URL changes nothing, so no effect fires — the handler reveals
// directly instead. Everything else goes through the URL-sync effect.
export type ConceptNavAction = 'push' | 'replace' | 'reveal'

export interface ConceptNav {
  action: ConceptNavAction
  href: string
}

export function planConceptNavigation(
  cid: string,
  lessonByConcept: Map<string, { title: string }>,
  currentTitle: string | null,
  conceptParam: string | null,
): ConceptNav {
  const { href, sameLesson } = resolveConceptTarget(cid, lessonByConcept, currentTitle)
  if (sameLesson && conceptParam === cid) return { action: 'reveal', href }
  return { action: sameLesson ? 'replace' : 'push', href }
}
