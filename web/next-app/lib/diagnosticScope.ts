import { concepts, type ConceptRecord } from './conceptData'

// The real scope of a diagnostic, which is not the number of concepts a learner
// ticked.
//
// Selecting domains is a starting point for the assessment, the way Math Academy
// describes choosing a course at registration: "it is just a starting point for our
// adaptive diagnostic assessment". The server expands the selected concept ids to
// their transitive prerequisite closure before it builds anything
// (`internal/planning/planner.go`, `collectChain`), then compresses that to a
// covering set of at most 30 and asks roughly two probes per member, capped at 45
// (`internal/diagnostic/cat.go`).
//
// So the flat sum of concepts in the selected domains is wrong in both directions.
// It overstates: "Select everything" reads as 657 questions when the assessment
// asks about 45. It also understates: picking one advanced domain silently pulls
// in every prerequisite beneath it, so the number shown omits what is about to be
// assessed.
//
// The closure is the one figure worth putting in front of a learner, and it is a
// pure walk of the prerequisite lists — no constants to mirror, so nothing here
// can drift the way a copied cover size would. The server walks the same edges
// from the same corpus; test/diagnosticScope.test.ts pins the two against each
// other so a corpus edit cannot separate them quietly.

const byId = new Map<string, ConceptRecord>(
  concepts.map(c => [c.id, c] as const),
)

/**
 * Every concept the assessment could reach from `ids`: the ids themselves plus
 * their prerequisites, transitively.
 *
 * Mirrors `collectChain`: a breadth-first walk over `prerequisites`, upward only.
 * Unknown ids and unknown prerequisite references are skipped rather than thrown,
 * because this runs during render and a corpus gap must not blank the onboarding
 * screen — the server is still the authority and will reject what it cannot
 * resolve.
 */
export function prerequisiteClosure(ids: Iterable<string>): string[] {
  const inChain = new Set<string>()
  const queue: string[] = []

  for (const id of ids) {
    if (inChain.has(id)) continue
    inChain.add(id)
    queue.push(id)
  }

  for (let i = 0; i < queue.length; i++) {
    const record = byId.get(queue[i])
    if (!record) continue
    for (const prereq of record.prerequisites ?? []) {
      if (inChain.has(prereq)) continue
      inChain.add(prereq)
      queue.push(prereq)
    }
  }

  return [...inChain]
}

/** How many concepts the assessment covers for this selection. */
export function scopeSize(ids: Iterable<string>): number {
  return prerequisiteClosure(ids).length
}