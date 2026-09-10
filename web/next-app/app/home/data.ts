import { concepts as conceptsData } from '../../lib/conceptData'

export const heroConcepts = conceptsData.map((c) => ({
  id: c.id,
  label: c.label,
  domain: c.domain,
  prerequisites: c.prerequisites,
}))

export const PIPELINE_STATES = [
  { label: 'UNSEEN', status: 'unseen' as const },
  { label: 'LEARNING', status: 'learning' as const },
  { label: 'PRACTICING', status: 'practicing' as const },
  { label: 'MASTERED', status: 'mastered' as const },
  { label: 'DECAYING', status: 'decaying' as const },
]

export const conceptCount = conceptsData.length
export const connectionCount = conceptsData.reduce(
  (sum, c) => sum + (c.prerequisites?.length ?? 0),
  0
)
export const domainCount = new Set(conceptsData.map((c) => c.domain)).size

export const domainOrder = [
  'arithmetic',
  'fractions',
  'prealgebra',
  'algebra',
  'geometry',
  'trigonometry',
  'calculus',
  'statistics',
  'linear_algebra',
  'discrete_math',
  'complex_numbers',
  'number_theory',
  'differential_equations',
  'abstract_algebra',
  'topology',
]

export const domainLabels = {
  arithmetic: 'Arithmetic',
  fractions: 'Fractions',
  prealgebra: 'Pre-Algebra',
  algebra: 'Algebra',
  geometry: 'Geometry',
  trigonometry: 'Trigonometry',
  calculus: 'Calculus',
  statistics: 'Statistics',
  linear_algebra: 'Linear Algebra',
  discrete_math: 'Discrete Math',
  complex_numbers: 'Complex Numbers',
  number_theory: 'Number Theory',
  differential_equations: 'Differential Equations',
  abstract_algebra: 'Abstract Algebra',
  topology: 'Topology',
} satisfies Record<string, string>

export const domainCounts = conceptsData.reduce(
  (acc: Record<string, number>, c) => {
    acc[c.domain] = (acc[c.domain] || 0) + 1
    return acc
  },
  {} as Record<string, number>
)

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
