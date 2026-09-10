export type MasteryStatus = 'mastered' | 'practicing' | 'learning' | 'unseen' | 'locked'

export interface GraphConceptInput {
  id: string
  prerequisites: string[]
}

export interface ProgressInput {
  status?: string
}

const KNOWN_PROGRESS = {
  MASTERED: 'mastered',
  PRACTICING: 'practicing',
  LEARNING: 'learning',
} satisfies Record<string, Exclude<MasteryStatus, 'unseen' | 'locked'>>

export function deriveStatuses(
  concepts: GraphConceptInput[],
  progress: Record<string, ProgressInput> = {}
) {
  const knownIds = new Set(concepts.map(c => c.id))
  const mastered = new Set<string>()
  const started = new Map<string, MasteryStatus>()

  for (const [id, p] of Object.entries(progress)) {
    if (!knownIds.has(id)) continue
    const mapped = KNOWN_PROGRESS[p?.status ?? '']
    if (!mapped) continue
    started.set(id, mapped)
    if (mapped === 'mastered') mastered.add(id)
  }

  const result: Record<string, MasteryStatus> = {}

  for (const c of concepts) {
    const s = started.get(c.id)
    if (s) {
      result[c.id] = s
      continue
    }
    const blockingPrereq = c.prerequisites.some(p => knownIds.has(p) && !mastered.has(p))
    result[c.id] = blockingPrereq ? 'locked' : 'unseen'
  }

  return result
}
