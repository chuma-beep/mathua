import graphMetaJson from '../../data/graph.meta.json'
import { DOMAIN_LABELS, DOMAIN_ORDER } from '../../lib/graphDomains'

// Tiny build-generated stats. The full concept corpus is fetched by the graph
// chunk at idle; the landing bundle only carries counts and labels.
interface GraphMeta {
  conceptCount: number
  connectionCount: number
  domainCount: number
  domainCounts: Record<string, number>
  layoutVersion: number
}

const graphMeta = graphMetaJson as GraphMeta

export const PIPELINE_STATES = [
  { label: 'UNSEEN', status: 'unseen' as const },
  { label: 'LEARNING', status: 'learning' as const },
  { label: 'PRACTICING', status: 'practicing' as const },
  { label: 'MASTERED', status: 'mastered' as const },
  { label: 'DECAYING', status: 'decaying' as const },
]

export const conceptCount = graphMeta.conceptCount
export const connectionCount = graphMeta.connectionCount
export const domainCount = graphMeta.domainCount
export const domainCounts = graphMeta.domainCounts

export const domainOrder = [...DOMAIN_ORDER]

export const domainLabels = DOMAIN_LABELS

export const levels = [
  { num: '01', name: 'Novice', range: '0–31' },
  { num: '02', name: 'Apprentice', range: '32–63' },
  { num: '03', name: 'Student', range: '64–95' },
  { num: '04', name: 'Scholar', range: '96–127' },
  { num: '05', name: 'Adept', range: '128–159' },
  { num: '06', name: 'Expert', range: '160–191' },
  { num: '07', name: 'Master', range: '192–223' },
  { num: '08', name: 'Grandmaster', range: '224–255' },
  { num: '09', name: 'Math Architect', range: '256–284', elite: true },
]
