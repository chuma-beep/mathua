'use client'

import { useState, useEffect, useMemo, useCallback, Suspense } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import Header from '../../components/Header'
import BottomTabs from '../../components/BottomTabs'
import SectionHeader from '../../components/SectionHeader'
import Footer from '../../components/Footer'
import AsciiDivider from '../../components/AsciiDivider'
import KatexContent from '../../components/KatexContent'
import SearchBar from '../../components/SearchBar'
import LessonQuiz from '../../components/LessonQuiz'
import MasteryBadge from '../../components/MasteryBadge'
import { getLessons, getLessonBody, getScores, type LessonInfo, type LessonsRes, type Scores } from '../../lib/api'
import { getUserInfo, getGuestId } from '../../lib/auth'
import conceptsData from '../../data/concepts.json'
import Loading from '../../components/Loading'

type ConceptEntry = { id: string; label: string; domain: string }
const conceptLabels = new Map<string, string>((conceptsData as ConceptEntry[]).map(c => [c.id, c.label]))

// ── Domain metadata ──────────────────────────────────────

const domainOrder = [
  'arith', 'fractions', 'prealgebra',
  'alg', 'trig', 'calc', 'calc.integral', 'calc.deriv', 'calc.limit',
  'linalg', 'linalg.vector', 'linalg.matrix', 'linalg.eigen', 'linalg.det',
  'linalg.transformations', 'linalg.span', 'linalg.basis',
  'complex', 'complex.concept',
  'alg.poly', 'alg.factor', 'alg.quad', 'alg.log', 'alg.systems',
  'arith.exp', 'arith.sqrt',
  'abstract', 'abstract.group', 'abstract.ring', 'abstract.rings', 'abstract.homomorphism', 'abstract.structures', 'abstract.subgroup',
  'discrete', 'discrete.sets', 'discrete.combinatorics', 'discrete.graphs', 'discrete.induction', 'discrete.logic', 'discrete.proof', 'discrete.recurrence',
  'fractions', 'frac.add', 'frac.benchmark', 'frac.compare', 'frac.concept', 'frac.div', 'frac.equivalent', 'frac.mixed', 'frac.mult', 'frac.on_number_line', 'frac.parts', 'frac.simplify', 'frac.sub', 'frac.to_decimal',
  'geo', 'geo.basic', 'geo.circle', 'geo.coord', 'geo.quad', 'geo.solid', 'geo.triangle',
  'nt', 'nt.adv', 'nt.basics', 'nt.congruence', 'nt.crypto', 'nt.diophantine', 'nt.divisibility', 'nt.euler_phi', 'nt.fermat_little', 'nt.gcd_euclidean', 'nt.modular',
  'ode', 'ode.basics', 'ode.concept', 'ode.exact', 'ode.homogeneous', 'ode.laplace', 'ode.linear_first', 'ode.nonhomogeneous', 'ode.separable', 'ode.systems',
  'pct', 'pct.concept', 'pct.discount', 'pct.find_rate', 'pct.from_dec', 'pct.increase', 'pct.of_number', 'pct.tax_tip', 'pct.to_dec',
  'ratio', 'ratio.concept', 'ratio.proportion', 'ratio.rate', 'ratio.scale', 'ratio.simplify',
  'dec', 'dec.add', 'dec.compare', 'dec.div', 'dec.from_frac', 'dec.mult', 'dec.round', 'dec.sub', 'dec.to_frac',
  'stat', 'stat.adv', 'stat.data', 'stat.dist', 'stat.infer', 'stat.prob',
  'topo', 'topo.continuous', 'topo.metric', 'topo.open_closed',
  'ml', 'ml.backpropagation',
  'alg.conic', 'alg.eq', 'alg.exp', 'alg.func', 'alg.ineq', 'alg.linear', 'alg.seq',
  'arith.abs_value', 'arith.add', 'arith.dec', 'arith.div', 'arith.factor', 'arith.mult', 'arith.neg', 'arith.order_ops', 'arith.place', 'arith.round', 'arith.sci_notation', 'arith.sub',
  'calc.seq', 'calc.series',
  'complex.add_sub', 'complex.adv', 'complex.basics', 'complex.conjugate', 'complex.de_moivre', 'complex.divide', 'complex.mult', 'complex.ops', 'complex.polar', 'complex.roots',
  'linalg.cramer', 'linalg.lintrans', 'linalg.sys', 'linalg.systems', 'linalg.vec',
  'prealg.eq', 'prealg.expr', 'prealg.ineq', 'prealg.real', 'prealg.types', 'prealg.var',
  'trig.adv', 'trig.basics', 'trig.eq', 'trig.graph', 'trig.graph_cos', 'trig.graph_sin', 'trig.hyperbolic', 'trig.ident', 'trig.ineq', 'trig.inverse', 'trig.law_cosines', 'trig.law_sines', 'trig.period', 'trig.pythagorean_id', 'trig.radians', 'trig.reciprocal', 'trig.reference_angle', 'trig.sin_cos_def', 'trig.special_angles', 'trig.tan_def', 'trig.unit_circle',
]

const domainLabels: Record<string, string> = {
  'arith': 'Arithmetic',
  'fractions': 'Fractions',
  'prealgebra': 'Pre-Algebra',
  'alg': 'Algebra',
  'trig': 'Trigonometry',
  'calc': 'Calculus',
  'calc.integral': 'Calculus · Integrals',
  'calc.deriv': 'Calculus · Derivatives',
  'calc.limit': 'Calculus · Limits',
  'linalg': 'Linear Algebra',
  'linalg.vector': 'Linear Algebra · Vectors',
  'linalg.matrix': 'Linear Algebra · Matrices',
  'linalg.eigen': 'Linear Algebra · Eigen',
  'linalg.det': 'Linear Algebra · Determinants',
  'linalg.transformations': 'Linear Algebra · Transformations',
  'linalg.span': 'Linear Algebra · Span',
  'linalg.basis': 'Linear Algebra · Basis',
  'complex': 'Complex Numbers',
  'complex.concept': 'Complex Numbers',
  'alg.poly': 'Algebra · Polynomials',
  'alg.factor': 'Algebra · Factoring',
  'alg.quad': 'Algebra · Quadratics',
  'alg.log': 'Algebra · Logarithms',
  'alg.systems': 'Algebra · Linear Systems',
  'arith.exp': 'Arithmetic · Exponents',
  'arith.sqrt': 'Arithmetic · Radicals',
  'abstract': 'Abstract Algebra',
  'abstract.group': 'Abstract Algebra · Groups',
  'abstract.ring': 'Abstract Algebra · Rings',
  'abstract.rings': 'Abstract Algebra · Rings',
  'abstract.homomorphism': 'Abstract Algebra · Homomorphisms',
  'abstract.structures': 'Abstract Algebra · Structures',
  'abstract.subgroup': 'Abstract Algebra · Subgroups',
  'discrete': 'Discrete Math',
  'discrete.sets': 'Discrete Math · Sets',
  'discrete.combinatorics': 'Discrete Math · Combinatorics',
  'discrete.graphs': 'Discrete Math · Graphs',
  'discrete.induction': 'Discrete Math · Induction',
  'discrete.logic': 'Discrete Math · Logic',
  'discrete.proof': 'Discrete Math · Proofs',
  'discrete.recurrence': 'Discrete Math · Recurrence',
  'frac.add': 'Fractions · Addition',
  'frac.benchmark': 'Fractions · Benchmark',
  'frac.compare': 'Fractions · Comparing',
  'frac.concept': 'Fractions · Concept',
  'frac.div': 'Fractions · Division',
  'frac.equivalent': 'Fractions · Equivalent',
  'frac.mixed': 'Fractions · Mixed Numbers',
  'frac.mult': 'Fractions · Multiplication',
  'frac.on_number_line': 'Fractions · Number Line',
  'frac.parts': 'Fractions · Parts of a Whole',
  'frac.simplify': 'Fractions · Simplifying',
  'frac.sub': 'Fractions · Subtraction',
  'frac.to_decimal': 'Fractions · to Decimal',
  'geo': 'Geometry',
  'geo.basic': 'Geometry · Basics',
  'geo.circle': 'Geometry · Circles',
  'geo.coord': 'Geometry · Coordinate',
  'geo.quad': 'Geometry · Quadrilaterals',
  'geo.solid': 'Geometry · Solids',
  'geo.triangle': 'Geometry · Triangles',
  'nt': 'Number Theory',
  'nt.adv': 'Number Theory · Advanced',
  'nt.basics': 'Number Theory · Basics',
  'nt.congruence': 'Number Theory · Congruence',
  'nt.crypto': 'Number Theory · Cryptography',
  'nt.diophantine': 'Number Theory · Diophantine Equations',
  'nt.divisibility': 'Number Theory · Divisibility',
  'nt.euler_phi': 'Number Theory · Euler\'s Totient',
  'nt.fermat_little': 'Number Theory · Fermat\'s Little Theorem',
  'nt.gcd_euclidean': 'Number Theory · GCD & Euclidean',
  'nt.modular': 'Number Theory · Modular Arithmetic',
  'ode': 'Differential Equations',
  'ode.basics': 'ODEs · Basics',
  'ode.concept': 'ODEs · Concept',
  'ode.exact': 'ODEs · Exact Equations',
  'ode.homogeneous': 'ODEs · Homogeneous',
  'ode.laplace': 'ODEs · Laplace Transform',
  'ode.linear_first': 'ODEs · First-Order Linear',
  'ode.nonhomogeneous': 'ODEs · Non-Homogeneous',
  'ode.separable': 'ODEs · Separable',
  'ode.systems': 'ODEs · Systems',
  'pct': 'Percentages',
  'pct.concept': 'Percentages · Concept',
  'pct.discount': 'Percentages · Discount',
  'pct.find_rate': 'Percentages · Find Rate',
  'pct.from_dec': 'Percentages · from Decimal',
  'pct.increase': 'Percentages · Increase',
  'pct.of_number': 'Percentages · of a Number',
  'pct.tax_tip': 'Percentages · Tax & Tip',
  'pct.to_dec': 'Percentages · to Decimal',
  'ratio': 'Ratios',
  'ratio.concept': 'Ratios · Concept',
  'ratio.proportion': 'Ratios · Proportion',
  'ratio.rate': 'Ratios · Rate',
  'ratio.scale': 'Ratios · Scale',
  'ratio.simplify': 'Ratios · Simplifying',
  'dec': 'Decimals',
  'dec.add': 'Decimals · Addition',
  'dec.compare': 'Decimals · Comparing',
  'dec.div': 'Decimals · Division',
  'dec.from_frac': 'Decimals · from Fractions',
  'dec.mult': 'Decimals · Multiplication',
  'dec.round': 'Decimals · Rounding',
  'dec.sub': 'Decimals · Subtraction',
  'dec.to_frac': 'Decimals · to Fractions',
  'stat': 'Statistics',
  'stat.adv': 'Statistics · Advanced',
  'stat.data': 'Statistics · Data',
  'stat.dist': 'Statistics · Distributions',
  'stat.infer': 'Statistics · Inference',
  'stat.prob': 'Statistics · Probability',
  'topo': 'Topology',
  'topo.continuous': 'Topology · Continuity',
  'topo.metric': 'Topology · Metric Spaces',
  'topo.open_closed': 'Topology · Open/Closed Sets',
  'ml': 'Machine Learning',
  'ml.backpropagation': 'ML · Backpropagation',
  'alg.conic': 'Algebra · Conic Sections',
  'alg.eq': 'Algebra · Equations',
  'alg.exp': 'Algebra · Exponents',
  'alg.func': 'Algebra · Functions',
  'alg.ineq': 'Algebra · Inequalities',
  'alg.linear': 'Algebra · Linear Equations',
  'alg.seq': 'Algebra · Sequences',
  'arith.abs_value': 'Arithmetic · Absolute Value',
  'arith.add': 'Arithmetic · Addition',
  'arith.dec': 'Arithmetic · Decimals',
  'arith.div': 'Arithmetic · Division',
  'arith.factor': 'Arithmetic · Factoring',
  'arith.mult': 'Arithmetic · Multiplication',
  'arith.neg': 'Arithmetic · Negatives',
  'arith.order_ops': 'Arithmetic · Order of Operations',
  'arith.place': 'Arithmetic · Place Value',
  'arith.round': 'Arithmetic · Rounding',
  'arith.sci_notation': 'Arithmetic · Scientific Notation',
  'arith.sub': 'Arithmetic · Subtraction',
  'calc.seq': 'Calculus · Sequences',
  'calc.series': 'Calculus · Series',
  'complex.add_sub': 'Complex · Add/Subtract',
  'complex.adv': 'Complex · Advanced',
  'complex.basics': 'Complex · Basics',
  'complex.conjugate': 'Complex · Conjugates',
  'complex.de_moivre': 'Complex · De Moivre',
  'complex.divide': 'Complex · Division',
  'complex.mult': 'Complex · Multiplication',
  'complex.ops': 'Complex · Operations',
  'complex.polar': 'Complex · Polar Form',
  'complex.roots': 'Complex · Roots',
  'linalg.cramer': 'Linear Algebra · Cramer\'s Rule',
  'linalg.vec': 'Linear Algebra · Vectors',
  'linalg.lintrans': 'Linear Algebra · Transformations',
  'linalg.sys': 'Linear Algebra · Systems',
  'linalg.systems': 'Linear Algebra · Systems',
  'prealg.eq': 'Pre-Algebra · Equations',
  'prealg.expr': 'Pre-Algebra · Expressions',
  'prealg.ineq': 'Pre-Algebra · Inequalities',
  'prealg.real': 'Pre-Algebra · Real Numbers',
  'prealg.types': 'Pre-Algebra · Number Types',
  'prealg.var': 'Pre-Algebra · Variables',
  'trig.adv': 'Trigonometry · Advanced',
  'trig.basics': 'Trigonometry · Basics',
  'trig.eq': 'Trigonometry · Equations',
  'trig.graph_cos': 'Trigonometry · Cosine Graphs',
  'trig.graph_sin': 'Trigonometry · Sine Graphs',
  'trig.hyperbolic': 'Trigonometry · Hyperbolic',
  'trig.ident': 'Trigonometry · Identities',
  'trig.ineq': 'Trigonometry · Inequalities',
  'trig.inverse': 'Trigonometry · Inverse',
  'trig.law_cosines': 'Trigonometry · Law of Cosines',
  'trig.law_sines': 'Trigonometry · Law of Sines',
  'trig.period': 'Trigonometry · Periodicity',
  'trig.pythagorean_id': 'Trigonometry · Pythagorean Identities',
  'trig.radians': 'Trigonometry · Radians',
  'trig.reciprocal': 'Trigonometry · Reciprocal',
  'trig.reference_angle': 'Trigonometry · Reference Angles',
  'trig.sin_cos_def': 'Trigonometry · Sine/Cosine Definitions',
  'trig.graph': 'Trigonometry · Graphs',
  'trig.special_angles': 'Trigonometry · Special Angles',
  'trig.tan_def': 'Trigonometry · Tangent Definition',
  'trig.unit_circle': 'Trigonometry · Unit Circle',
}

const domainIcons: Record<string, string> = {
  arith: '+',
  fractions: '½',
  frac: '½',
  prealgebra: '○',
  alg: 'x',
  trig: 'π',
  calc: '∫',
  'calc.integral': '∫',
  'calc.deriv': '∂',
  'calc.limit': '→',
  linalg: '⊕',
  complex: 'i',
  abstract: 'λ',
  discrete: '∈',
  nt: 'ℤ',
  geo: '△',
  ode: '∵',
  pct: '%',
  ratio: ':',
  dec: '.',
  stat: 'Σ',
  topo: '○',
  ml: '✦',
}

function domainIcon(domain: string): string {
  const parts = domain.split('.')
  for (let i = parts.length; i > 0; i--) {
    const key = parts.slice(0, i).join('.')
    if (domainIcons[key]) return domainIcons[key]
  }
  return '◇'
}

function lessonProgress(lesson: LessonInfo): { mastered: number; total: number } {
  if (!lesson.progress) return { mastered: 0, total: lesson.concepts.length }
  let mastered = 0
  for (const cid of lesson.concepts) {
    const p = lesson.progress[cid]
    if (p && p.status === 'MASTERED') mastered++
  }
  return { mastered, total: lesson.concepts.length }
}

function QuizGateBanner({ scores }: { scores: Scores | null }) {
  const xp = scores?.xp_total ?? 0
  const goal = 150 // CONTEXT.md Quiz 150 XP gate MA verbatim
  const done = xp >= goal
  const pct = Math.min((xp / goal) * 100, 100)
  return (
    <div className={`mt-6 border p-4 flex flex-col sm:flex-row items-center justify-between gap-3 ${done ? 'border-mathua-blue bg-mathua-surface' : 'border-mathua-border bg-mathua-surface'}`}>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className={`px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider ${done ? 'bg-mathua-blue text-white' : 'bg-mathua-border text-mathua-muted'}`}>
            {done ? 'Quiz due' : '150 XP gate'}
          </span>
          <span className="font-mono text-xs text-mathua-primary truncate">
            {done ? '150 XP reached — take your mastery check' : `${xp} / ${goal} XP toward next quiz`}
          </span>
        </div>
        <div className="mt-2 h-1 bg-mathua-code overflow-hidden">
          <div className={`h-full transition-all ${done ? 'bg-mathua-blue' : 'bg-mathua-blue/60'}`} style={{ width: `${pct}%` }} />
        </div>
      </div>
      {done ? (
        <Link href="/goals?quiz=1" className="shrink-0 border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-6 py-2 font-mono text-xs min-h-[36px] inline-flex items-center justify-center">
          Take Test →
        </Link>
      ) : null}
    </div>
  )
}

// ── Domain Overview ─────────────────────────────────────

function DomainOverview({
  domains,
  lessonsByDomain,
  domainAgg,
  allLessons,
  onSelectDomain,
  onSelectLesson,
}: {
  domains: string[]
  lessonsByDomain: Record<string, LessonInfo[]>
  domainAgg: Record<string, { mastered: number; total: number; pct: number }>
  allLessons: { title: string; body: string; domain: string; concepts: string[]; conceptLabels: string[] }[]
  onSelectDomain: (d: string) => void
  onSelectLesson: (lesson: LessonInfo) => void
}) {
  return (
    <div className="max-w-4xl mx-auto mt-8">
      <SectionHeader label="Study" title="Browse Lessons" />

      <SearchBar
        items={allLessons}
        onSelect={(item) => {
          const lesson = lessonsByDomain[item.domain]?.find(l => l.title === item.title)
          if (lesson) onSelectLesson(lesson)
        }}
      />

      {domains.length === 0 && (
        <p className="text-mathua-muted text-sm text-center">No lessons available.</p>
      )}

      <div className="space-y-2">
        {domains.map((domain, idx) => {
          const lessons = lessonsByDomain[domain]
          const agg = domainAgg[domain]
          const label = domainLabels[domain] || domain
          const icon = domainIcon(domain)
          const allMastered = agg.total > 0 && agg.mastered === agg.total

          return (
            <button
              key={domain}
              onClick={() => onSelectDomain(domain)}
              className={`w-full text-left transition-all duration-200 group ${
                allMastered ? 'opacity-50 hover:opacity-70' : ''
              } animate-fadeIn`}
              style={{ animationDelay: `${idx * 30}ms` }}
            >
              <div className="flex items-stretch border border-mathua-border hover:shadow-card-hover transition-shadow bg-mathua-surface-elevated">
                <div
                  className={`w-[3px] shrink-0 transition-colors ${
                    allMastered
                      ? 'bg-mathua-green'
                      : agg.total > 0 && agg.mastered > 0
                      ? 'bg-mathua-blue'
                      : 'bg-mathua-border group-hover:bg-mathua-blue'
                  }`}
                />

                <div className="flex-1 min-w-0 p-3 sm:p-4 overflow-hidden">
                  <div className="flex items-center gap-2 sm:gap-3 min-w-0">
                    <span className="font-mono text-lg text-mathua-blue shrink-0 w-6 text-center select-none">
                      {icon}
                    </span>
                    <span className="font-mono text-sm text-mathua-primary group-hover:text-mathua-blue transition-colors min-w-0 flex-1 truncate">
                      {label}
                    </span>
                    <span className="font-mono text-[11px] text-mathua-muted shrink-0 ml-1">
                      {lessons.length} lesson{lessons.length !== 1 ? 's' : ''}
                    </span>
                  </div>

                  {agg.total > 0 && (
                    <div className="mt-2 pl-9">
                      <div className="h-[3px] bg-mathua-code">
                        <div
                          className={`h-full transition-all duration-500 ${
                            allMastered ? 'bg-mathua-green' : 'bg-mathua-blue'
                          }`}
                          style={{ width: `${agg.pct}%` }}
                        />
                      </div>
                      <div className="flex items-center gap-2 mt-1">
                        <span className="font-mono text-[10px] text-mathua-muted">
                          {agg.mastered} / {agg.total} concepts mastered
                        </span>
                        <span className={`font-mono text-[10px] ${
                          allMastered ? 'text-mathua-green' : 'text-mathua-blue'
                        }`}>
                          {agg.pct}%
                        </span>
                      </div>
                    </div>
                  )}
                </div>

                <div className="hidden sm:flex items-center pr-4 shrink-0">
                  <span className="font-mono text-[11px] text-mathua-blue opacity-0 group-hover:opacity-100 transition-opacity">
                    →
                  </span>
                </div>
              </div>
            </button>
          )
        })}
      </div>
      <Link href="/leaderboard" className="mt-4 flex items-center justify-between border border-mathua-border bg-mathua-surface p-3 hover:border-mathua-blue transition-colors">
        <span className="font-mono text-xs text-mathua-primary">Leaderboard</span>
        <span className="font-mono text-[11px] text-mathua-blue">See weekly ranking →</span>
      </Link>
    </div>
  )
}

// ── Domain Drill-Down ───────────────────────────────────

function DomainDrillDown({
  domain,
  lessons,
  agg,
  onBack,
  onSelectLesson,
}: {
  domain: string
  lessons: LessonInfo[]
  agg?: { mastered: number; total: number; pct: number }
  onBack: () => void
  onSelectLesson: (lesson: LessonInfo) => void
}) {
  const label = domainLabels[domain] || domain
  const icon = domainIcon(domain)

  return (
    <div className="max-w-7xl mx-auto mt-8 mb-16">
      <button
        onClick={onBack}
        className="text-mathua-secondary text-xs font-mono hover:text-mathua-blue mb-6"
      >
        ← All domains
      </button>

      <div className="flex items-center gap-3 mb-4">
        <span className="font-mono text-2xl text-mathua-blue select-none">{icon}</span>
        <SectionHeader label="Domain" title={label} className="flex-1" />
      </div>

      {agg && agg.total > 0 && (
        <div className="border border-mathua-border p-4 mb-6 bg-mathua-surface">
          <div className="flex items-center gap-3 mb-2">
            <span className="font-mono text-[11px] text-mathua-muted">Domain progress</span>
            <span className="font-mono text-[11px] text-mathua-blue">
              {agg.mastered} / {agg.total} concepts
            </span>
            <span className={`font-mono text-[11px] ${
              agg.mastered === agg.total ? 'text-mathua-green' : 'text-mathua-blue'
            }`}>
              {agg.pct}%
            </span>
          </div>
          <div className="h-[3px] bg-mathua-code">
            <div
              className={`h-full transition-all duration-500 ${
                agg.mastered === agg.total ? 'bg-mathua-green' : 'bg-mathua-blue'
              }`}
              style={{ width: `${agg.pct}%` }}
            />
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        {lessons.map((lesson, i) => {
          const { mastered, total } = lessonProgress(lesson)
          const pct = total > 0 ? Math.round((mastered / total) * 100) : 0
          const allDone = total > 0 && mastered === total

          return (
            <button
              key={i}
              onClick={() => onSelectLesson(lesson)}
              className={`text-left transition-all duration-200 group animate-fadeIn ${
                allDone ? 'opacity-50 hover:opacity-70' : ''
              }`}
              style={{ animationDelay: `${i * 30}ms` }}
            >
              <div className="border border-mathua-border hover:shadow-card-hover transition-shadow bg-mathua-surface-elevated p-3 sm:p-4 min-w-0 overflow-hidden">
                <div className="flex items-center justify-between gap-2 min-w-0">
                  <div className="font-mono text-sm text-mathua-primary group-hover:text-mathua-blue transition-colors min-w-0 flex-1 line-clamp-2 break-words">
                    {lesson.title}
                  </div>
                  <span className="hidden sm:inline font-mono text-[11px] text-mathua-blue opacity-0 group-hover:opacity-100 transition-opacity shrink-0 ml-2">
                    View →
                  </span>
                </div>

                <div className="flex items-center gap-2 mt-2">
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {total} concept{total !== 1 ? 's' : ''}
                  </span>
                  {lesson.progress && (
                    <>
                      <div className="flex-1 h-[3px] bg-mathua-code max-w-[120px]">
                        <div
                          className={`h-full transition-all ${
                            allDone ? 'bg-mathua-green' : 'bg-mathua-blue'
                          }`}
                          style={{ width: `${pct}%` }}
                        />
                      </div>
                      <span className="font-mono text-[10px] text-mathua-muted">
                        {mastered}/{total}
                      </span>
                    </>
                  )}
                </div>

                {lesson.progress && (
                  <div className="flex items-center gap-1 mt-2">
                    {lesson.concepts.map((cid) => (
                      <MasteryBadge key={cid} status={lesson.progress?.[cid]?.status} size="sm" />
                    ))}
                  </div>
                )}
              </div>
            </button>
          )
        })}
      </div>
    </div>
  )
}

// ── Lesson Detail ───────────────────────────────────────

function LessonDetail({
  lesson,
  domain,
  onBack,
}: {
  lesson: LessonInfo
  domain: string | null
  onBack: () => void
}) {
  return (
    <div className="max-w-7xl mx-auto mt-8 mb-16">
      <button
        onClick={onBack}
        className="text-mathua-secondary text-xs font-mono hover:text-mathua-blue mb-6"
      >
        ← {domain ? domainLabels[domain] || domain : 'All domains'}
      </button>
      <SectionHeader label="Lesson" title={lesson.title} />

      <div className="bg-mathua-surface border border-mathua-border p-4 sm:p-6 mt-6 mb-6 min-w-0 overflow-hidden">
        <div className="flex items-center gap-2 sm:gap-3 flex-wrap min-w-0">
          <span className="text-mathua-muted text-xs font-mono shrink-0">Concepts:</span>
          {lesson.concepts.map((cid) => {
            const p = lesson.progress?.[cid]
            return (
              <Link
                key={cid}
                href={`/concept?id=${encodeURIComponent(cid)}`}
                className="inline-flex items-center gap-1.5 border border-mathua-border px-2.5 py-1.5 min-h-[36px] text-xs font-mono text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue transition-colors max-w-full"
              >
                <MasteryBadge status={p?.status} size="sm" />
                <span className="truncate">{cid}</span>
              </Link>
            )
          })}
        </div>
      </div>

      {lesson.prerequisites && lesson.prerequisites.length > 0 && (
        <div className="mb-8">
          <h3 className="font-mono text-[11px] text-mathua-muted mb-3 border-b border-mathua-border pb-2 uppercase tracking-wider">
            Before you start ({lesson.prerequisites.length})
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {lesson.prerequisites.map((p) => (
              <Link
                key={p.id}
                href={`/concept?id=${encodeURIComponent(p.id)}`}
                className="bg-mathua-surface border border-mathua-border p-3 hover:border-mathua-blue transition-colors block"
              >
                <div className="flex items-center gap-2">
                  <MasteryBadge status={p.status} size="sm" />
                  <div className="font-mono text-xs text-mathua-primary">{p.label}</div>
                </div>
                <div className="flex items-center gap-2 mt-1.5">
                  <div className="flex-1 h-1 bg-mathua-bg overflow-hidden">
                    <div
                      className="h-full bg-mathua-blue transition-all"
                      style={{ width: `${Math.round(p.mastery_pct * 100)}%` }}
                    />
                  </div>
                  <span className="font-mono text-[10px] text-mathua-muted">
                    {Math.round(p.mastery_pct * 100)}%
                  </span>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      <div className="bg-mathua-surface border border-mathua-border p-4 sm:p-6 md:p-8 lg:p-10 w-full max-w-full min-w-0 overflow-hidden">
        <div className="flex items-center gap-2 mb-4 border-b border-mathua-border pb-3">
          <span className="bg-mathua-blue text-white px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
            Worked example
          </span>
          <span className="font-mono text-[10px] text-mathua-muted">
            study this first, then practice below
          </span>
        </div>
        <div className="w-full max-w-full min-w-0 overflow-hidden">
          <KatexContent>{lesson.body}</KatexContent>
        </div>
      </div>

      <div className="mt-2 flex items-center gap-2">
        <span className="bg-mathua-border text-mathua-primary px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider">
          Practice
        </span>
        <span className="font-mono text-[10px] text-mathua-muted">
          2 in a row to advance
        </span>
      </div>
      {lesson.concepts.slice(0, 3).map(cid => (
        <LessonQuiz key={cid} conceptId={cid} limit={4} />
      ))}

      {lesson.prerequisites && lesson.prerequisites.length > 0 && (
        <div className="mt-8">
          <p className="font-mono text-xs text-mathua-muted border-b border-mathua-border pb-2 mb-3">
            Requires: {lesson.prerequisites.map((p, i) => (
              <span key={p.id}>
                {i > 0 && <span className="mx-1 text-mathua-border">·</span>}
                <Link href={`/concept?id=${encodeURIComponent(p.id)}`} className="text-mathua-blue hover:text-mathua-blue-hover transition-colors">
                  {p.label}
                </Link>
              </span>
            ))}
          </p>
        </div>
      )}

      {lesson.dependents && lesson.dependents.length > 0 && (
        <div className="mt-6 mb-8">
          <h3 className="font-mono text-[11px] text-mathua-muted mb-3 border-b border-mathua-border pb-2 uppercase tracking-wider">
            What to study next ({lesson.dependents.length})
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {lesson.dependents.map((d) => (
              <Link
                key={d.id}
                href={`/concept?id=${encodeURIComponent(d.id)}`}
                className="bg-mathua-surface border border-mathua-border p-3 hover:border-mathua-blue transition-colors block"
              >
                <div className="flex items-center gap-2">
                  <MasteryBadge status={d.status} size="sm" />
                  <div className="font-mono text-xs text-mathua-primary">{d.label}</div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}

      <div className="mt-8 text-center">
        <Link
          href={`/session`}
          className="inline-block border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white h-12 px-8 font-medium text-sm leading-[48px]"
        >
          Practice these concepts
        </Link>
      </div>
    </div>
  )
}

let lessonsCache: { key: string; res: LessonsRes } | null = null
function getLessonsCached(studentId?: string): Promise<LessonsRes> {
  const key = studentId ?? ''
  if (lessonsCache && lessonsCache.key === key) return Promise.resolve(lessonsCache.res)
  return getLessons(studentId).then(res => {
    lessonsCache = { key, res }
    return res
  })
}

function StudyContent() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const lessonParam = searchParams.get('lesson')
  const domainParam = searchParams.get('domain')
  const conceptParam = searchParams.get('concept')

  const [lessonsByDomain, setLessonsByDomain] = useState<Record<string, LessonInfo[]>>({})
  const [selectedBody, setSelectedBody] = useState<string | null>(null)
  const [selectedLesson, setSelectedLesson] = useState<LessonInfo | null>(null)
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [scores, setScores] = useState<Scores | null>(null)

  useEffect(() => {
    const user = getUserInfo()
    const studentId = user?.student_id
    getLessonsCached(studentId).then(res => {
      setLessonsByDomain(res.lessons)
      setLoading(false)
    }).catch((e) => { console.error('getLessons failed:', e); setLoading(false) })
    const sid = studentId || getGuestId() || ''
    if (sid) {
      getScores(sid).then(setScores).catch(() => {})
    }
  }, [])

  // URL → state sync
  useEffect(() => {
    if (lessonParam) {
      for (const lessons of Object.values(lessonsByDomain)) {
        const found = lessons.find(l => l.title === lessonParam)
        if (found) { setSelectedLesson(found); return }
      }
      setSelectedLesson(null)
      return
    }
    if (conceptParam && Object.keys(lessonsByDomain).length > 0) {
      for (const lessons of Object.values(lessonsByDomain)) {
        const found = lessons.find(l => l.concepts.includes(conceptParam))
        if (found) { setSelectedLesson(found); return }
      }
    }

    if (domainParam && lessonsByDomain[domainParam]) {
      setSelectedDomain(domainParam)
    } else if (!domainParam) {
      setSelectedDomain(null)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lessonParam, domainParam, conceptParam, lessonsByDomain])

  // Lazily fetch the selected lesson's markdown body.
  useEffect(() => {
    if (!selectedLesson || selectedLesson.body) {
      setSelectedBody(selectedLesson?.body ?? null)
      return
    }
    let cancelled = false
    setSelectedBody(null)
    getLessonBody(selectedLesson.title)
      .then(body => { if (!cancelled) setSelectedBody(body) })
      .catch((e) => { console.error('lesson body failed:', e) })
    return () => { cancelled = true }
  }, [selectedLesson])

  const hydratedLesson = useMemo(
    () => (selectedLesson ? { ...selectedLesson, body: selectedBody ?? '' } : null),
    [selectedLesson, selectedBody]
  )

  // Popstate: browser back/forward  // Popstate: browser back/forward
  useEffect(() => {
    const onPop = () => {
      const params = window.location.search
      if (!params.includes('lesson=') && !params.includes('domain=') && !params.includes('concept=')) {
        setSelectedLesson(null)
        setSelectedDomain(null)
      }
    }
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  const sortedDomains = useMemo(() => {
    return Object.keys(lessonsByDomain).sort((a, b) => {
      const ai = domainOrder.indexOf(a)
      const bi = domainOrder.indexOf(b)
      return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
    })
  }, [lessonsByDomain])

  // Aggregate progress per domain
  const domainAgg = useMemo(() => {
    const agg: Record<string, { mastered: number; total: number; pct: number }> = {}
    for (const [domain, lessons] of Object.entries(lessonsByDomain)) {
      let mastered = 0, total = 0
      for (const lesson of lessons) {
        const p = lessonProgress(lesson)
        mastered += p.mastered
        total += p.total
      }
      agg[domain] = { mastered, total, pct: total > 0 ? Math.round((mastered / total) * 100) : 0 }
    }
    return agg
  }, [lessonsByDomain])

  const allLessons = useMemo(() => {
    const items: { title: string; body: string; domain: string; concepts: string[]; conceptLabels: string[] }[] = []
    for (const [domain, lessons] of Object.entries(lessonsByDomain)) {
      for (const l of lessons) {
        items.push({
          title: l.title,
          body: l.body ?? '',
          domain,
          concepts: l.concepts,
          conceptLabels: l.concepts.map(cid => conceptLabels.get(cid) || cid),
        })
      }
    }
    return items
  }, [lessonsByDomain])

  if (loading) {
    return (
      <>
        <Header />
        <div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0">
          <Loading label="LOADING LESSONS" />
        </div>
        <BottomTabs />
        <Footer />
      </>
    )
  }

  const backHref = selectedLesson
    ? (selectedDomain ? `/study?domain=${encodeURIComponent(selectedDomain)}` : '/study')
    : selectedDomain
    ? '/study'
    : '/'

  return (
    <>
      <Header />
      <div className="max-w-container mx-auto px-4 sm:px-6 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 overflow-x-hidden min-w-0">
        <section className="pt-8 min-w-0 overflow-hidden">
          {!selectedLesson && (
            <span className="flex justify-between items-center mb-4">
              <Link href={backHref} className="text-mathua-secondary text-sm hover:text-mathua-primary">
                ← Back
              </Link>
            </span>
          )}

          <QuizGateBanner scores={scores} />

          {selectedLesson ? (
            // ── Lesson Detail ──
            <LessonDetail
              lesson={hydratedLesson ?? selectedLesson}
              domain={selectedDomain}
              onBack={() => {
                setSelectedLesson(null)
                const url = selectedDomain
                  ? '/study?domain=' + encodeURIComponent(selectedDomain)
                  : '/study'
                window.history.replaceState(null, '', url)
              }}
            />
          ) : selectedDomain ? (
            // ── Domain Drill-Down ──
            <DomainDrillDown
              domain={selectedDomain}
              lessons={lessonsByDomain[selectedDomain] || []}
              agg={domainAgg[selectedDomain]}
              onBack={() => {
                setSelectedDomain(null)
                window.history.replaceState(null, '', '/study')
              }}
              onSelectLesson={(lesson) => {
                setSelectedLesson(lesson)
                router.push('/study?domain=' + encodeURIComponent(selectedDomain) + '&lesson=' + encodeURIComponent(lesson.title))
              }}
            />
          ) : (
            // ── Domain Overview ──
            <DomainOverview
              domains={sortedDomains}
              lessonsByDomain={lessonsByDomain}
              domainAgg={domainAgg}
              allLessons={allLessons}
              onSelectDomain={(d) => {
                setSelectedDomain(d)
                router.push('/study?domain=' + encodeURIComponent(d))
              }}
              onSelectLesson={(lesson) => {
                setSelectedLesson(lesson)
                router.push('/study?lesson=' + encodeURIComponent(lesson.title))
              }}
            />
          )}
        </section>

        <AsciiDivider pattern="wave" />
        <Footer />
      </div>
      <BottomTabs />
    </>
  )
}

export default function StudyPage() {
  return (
    <Suspense fallback={
      <><Header /><div className="max-w-container mx-auto px-4 sm:px-6 pt-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-0 text-center overflow-x-hidden min-w-0"><Loading label="LOADING LESSONS" /></div><Footer /></>
    }>
      <StudyContent />
    </Suspense>
  )
}
