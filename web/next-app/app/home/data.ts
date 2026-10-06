import graphMetaJson from '../../data/graph.meta.json'
import { DOMAIN_LABELS, DOMAIN_ORDER } from '../../lib/graphDomains'
// A type-only import, so it costs the landing bundle nothing: graphPayload.ts is otherwise
// pulled in by the idle-loaded graph chunk and the build script, and zod with it.
// This file redeclared the interface locally, which meant the landing page's view of the
// artifact was a second schema — adding quizGateXP to the real one type-checked here as a
// missing property rather than as a change to the page, which is a quieter way to drift.
import type { GraphMeta } from '../../lib/graphPayload'

const graphMeta = graphMetaJson as GraphMeta

export const PIPELINE_STATES = [
  { label: 'UNSEEN', status: 'unseen' as const },
  { label: 'LEARNING', status: 'learning' as const },
  { label: 'PRACTICING', status: 'practicing' as const },
  { label: 'MASTERED', status: 'mastered' as const },
  { label: 'DECAYING', status: 'decaying' as const },
]

export const conceptCount = graphMeta.conceptCount

// The landing page used to hardcode 150 while the engine gates the quiz at 50 (xp.QuizGateXP,
// rescaled by ADR-020). One number now, guarded from Go by TestLandingQuizGateMatchesEngine.
export const quizGateXP = graphMeta.quizGateXP
export const connectionCount = graphMeta.connectionCount
export const domainCount = graphMeta.domainCount
export const domainCounts = graphMeta.domainCounts

export const domainOrder = [...DOMAIN_ORDER]

export const domainLabels = DOMAIN_LABELS

// Ranges are derived from the corpus size, because the hand-written ladder was a snapshot.
// It stopped at "256–284" for a catalogue of 657, so the top four levels a learner could
// actually reach were not on the page. ADR-013's rule for corpus counts applies verbatim:
// a number that describes the corpus must be computed from the corpus.
const LEVEL_NAMES = [
  'Novice',
  'Apprentice',
  'Student',
  'Scholar',
  'Adept',
  'Expert',
  'Master',
  'Grandmaster',
  'Math Architect',
] as const

export const levels = LEVEL_NAMES.map((name, i) => {
  const lo = Math.floor((conceptCount * i) / LEVEL_NAMES.length)
  const hi = Math.ceil((conceptCount * (i + 1)) / LEVEL_NAMES.length) - 1
  const elite = i === LEVEL_NAMES.length - 1
  return {
    num: String(i + 1).padStart(2, '0'),
    name,
    range: i === 0 ? `0–${hi}` : `${lo}–${hi}`,
    ...(elite ? { elite: true } : {}),
  }
})
