import { countsAsMastered, masteryState, type MasteryState } from './progress'

/**
 * A concept's state on the graph, plus the one thing that is not a mastery state.
 *
 * `locked` is not stored and not derived from a status: it means *a prerequisite is
 * not yet satisfied*, which is a position in the DAG rather than a fact about this
 * learner. It is kept in the same union because every caller has to handle both.
 */
export type GraphStatus = MasteryState | 'locked'

export interface GraphConceptInput {
  id: string
  prerequisites: string[]
}

export { type MasteryState }

/**
 * Derives every node's graph status from the catalogue and the learner's progress.
 *
 * This used to carry its own copies of both mastery predicates, and had already drifted
 * from them: it mapped MASTERED, PRACTICING and LEARNING but not DECAYING, so a concept
 * the learner had mastered and not reviewed for a fortnight was neither recorded as
 * started nor as mastered. It therefore rendered as `unseen` — with no progress bar, as
 * if never attempted — and, because it was missing from the `mastered` set, it locked
 * every concept downstream of it.
 *
 * That is precisely the failure `lib/progress.ts` was extracted to prevent, and its own
 * comment predicted it: "they would drift the moment one was edited, and a drift here is
 * invisible: Profile would count a concept as mastered while the scheduler scheduled it
 * for review." Here it was worse — the graph relocked the frontier. Both questions are
 * now asked of `lib/progress.ts`, which is also what
 * `scheduler.prereqsMet` mirrors on the server.
 */
export function deriveStatuses(
  concepts: GraphConceptInput[],
  progress: Record<string, { status?: string }> = {},
) {
  const knownIds = new Set(concepts.map(c => c.id))
  const states = new Map<string, MasteryState>()

  for (const [id, p] of Object.entries(progress)) {
    if (!knownIds.has(id)) continue
    const state = masteryState(p)
    // `unseen` here means "no usable progress row", which is the same as having no row
    // at all — so it is not recorded and the concept is judged on its prerequisites.
    if (state === 'unseen') continue
    states.set(id, state)
  }

  const result: Record<string, GraphStatus> = {}

  for (const c of concepts) {
    const state = states.get(c.id)
    if (state) {
      result[c.id] = state
      continue
    }
    const blockingPrereq = (c.prerequisites ?? []).some(
      pid => knownIds.has(pid) && !countsAsMastered(progress[pid]),
    )
    result[c.id] = blockingPrereq ? 'locked' : 'unseen'
  }

  return result
}