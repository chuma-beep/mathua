// Domain cluster colors + curriculum order shared by the landing hero
// (legend, poster-adjacent UI, WebGL scene) and the build script.
// Zero-dependency so static sections can import it without pulling three.js
// into the page bundle.
const DOMAIN_COLORS = {
  arithmetic:           '#4db8a0',
  fractions:            '#a8a0f0',
  prealgebra:           '#7dd3fc',
  algebra:              '#e8a849',
  geometry:             '#86efac',
  trigonometry:         '#fda4af',
  complex_numbers:      '#e879f9',
  precalculus:          '#fda4af',
  calc1:                '#e879f9',
  calc2:                '#d946ef',
  calc3:                '#a21caf',
  linear_algebra:       '#86efac',
  statistics:           '#67e8f9',
  machine_learning:     '#818cf8',
  discrete_math:        '#f0abab',
  number_theory:        '#a8e6cf',
  differential_equations:'#fdba74',
  abstract_algebra:     '#c4b5fd',
  topology:             '#f9a8d4',
  graph_theory:         '#a3e635',
  probability_theory:   '#5eead4',
  numerical_analysis:   '#fbbf24',
  real_analysis:        '#a78bfa',
} satisfies Record<string, string>

const DOMAIN_COLORS_LIGHT = {
  arithmetic:           '#3a9a8a',
  fractions:            '#8a80d8',
  prealgebra:           '#5ab8dc',
  algebra:              '#c08a30',
  geometry:             '#68c88c',
  trigonometry:         '#e88a99',
  complex_numbers:      '#c868e8',
  precalculus:          '#e88a99',
  calc1:                '#c868e8',
  calc2:                '#a830c0',
  calc3:                '#802090',
  linear_algebra:       '#68c88c',
  statistics:           '#50c8d8',
  machine_learning:     '#5b64d6',
  discrete_math:        '#d08a8a',
  number_theory:        '#80c8a8',
  differential_equations:'#d09050',
  abstract_algebra:     '#a090d0',
  topology:             '#d080b0',
  graph_theory:         '#6a9a18',
  probability_theory:   '#199e8c',
  numerical_analysis:   '#b07d0a',
  real_analysis:        '#6d46c9',
} satisfies Record<string, string>

export const DOMAIN_ORDER = [
  'arithmetic',
  'fractions',
  'prealgebra',
  'algebra',
  'geometry',
  'trigonometry',
  'complex_numbers',
  'precalculus',
  'calc1',
  'calc2',
  'calc3',
  'linear_algebra',
  'statistics',
  'machine_learning',
  'discrete_math',
  'number_theory',
  'differential_equations',
  'abstract_algebra',
  'topology',
  'graph_theory',
  'probability_theory',
  'numerical_analysis',
  'real_analysis',
] as const

export const DOMAIN_LABELS = {
  arithmetic: 'Arithmetic',
  fractions: 'Fractions',
  prealgebra: 'Pre-Algebra',
  algebra: 'Algebra',
  geometry: 'Geometry',
  trigonometry: 'Trigonometry',
  complex_numbers: 'Complex Numbers',
  precalculus: 'Precalculus',
  calc1: 'Calculus I',
  calc2: 'Calculus II',
  calc3: 'Calculus III',
  linear_algebra: 'Linear Algebra',
  statistics: 'Statistics',
  machine_learning: 'Machine Learning',
  discrete_math: 'Discrete Math',
  number_theory: 'Number Theory',
  differential_equations: 'Differential Equations',
  abstract_algebra: 'Abstract Algebra',
  topology: 'Topology',
  graph_theory: 'Graph Theory',
  probability_theory: 'Probability Theory',
  numerical_analysis: 'Numerical Analysis',
  real_analysis: 'Real Analysis',
} satisfies Record<string, string>

const FALLBACK_COLOR = '#5a6577'
const FALLBACK_COLOR_LIGHT = '#888'

export function domainColor(domain: string, theme: 'dark' | 'light'): string {
  const map = theme === 'dark' ? DOMAIN_COLORS : DOMAIN_COLORS_LIGHT
  return map[domain] ?? (theme === 'dark' ? FALLBACK_COLOR : FALLBACK_COLOR_LIGHT)
}

export function domainLabel(domain: string): string {
  return DOMAIN_LABELS[domain] ?? domain.replace(/_/g, ' ')
}
