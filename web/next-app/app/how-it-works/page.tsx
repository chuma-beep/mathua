'use client'

import { useState, useEffect, useRef } from 'react'
import dynamic from 'next/dynamic'
import Header from '../../components/Header'
import AsciiDivider from '../../components/AsciiDivider'

const StudentModelFlow = dynamic(() => import('../../components/StudentModelFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 320, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px' }}>Loading&hellip;</div>,
})

const DiagnosticFlow = dynamic(() => import('../../components/DiagnosticFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 420, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px' }}>Loading&hellip;</div>,
})

const SchedulerFlow = dynamic(() => import('../../components/SchedulerFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 320, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px' }}>Loading&hellip;</div>,
})

const GradingFlow = dynamic(() => import('../../components/GradingFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 300, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px' }}>Loading&hellip;</div>,
})

const sections = [
  { id: 'concept-graph', label: 'Concept graph' },
  { id: 'student-model', label: 'Student model' },
  { id: 'spaced-repetition', label: 'Spaced repetition' },
  { id: 'diagnostic', label: 'Diagnostic' },
  { id: 'scheduler', label: 'Scheduler' },
  { id: 'scoring', label: 'Scoring' },
  { id: 'generators', label: 'Generators' },
  { id: 'symbolic-grading', label: 'Expression grading' },
]

function NavSidebar({ activeSection }: { activeSection: string }) {
  return (
    <nav className="w-[220px] flex-shrink-0 sticky top-[100px] h-fit max-md:fixed max-md:top-[57px] max-md:left-0 max-md:right-0 max-md:w-full max-md:z-[99] max-md:flex max-md:overflow-x-auto max-md:p-[8px_16px] max-md:gap-2"
      style={{ background: 'var(--bg)' }}
    >
      {sections.map((section) => (
        <a
          key={section.id}
          href={`#${section.id}`}
          className="block text-[12px] py-2 px-3 mb-1 transition-colors max-md:mb-0 max-md:whitespace-nowrap"
          style={{
            fontFamily: "'IBM Plex Mono', monospace",
            color: activeSection === section.id ? 'var(--accent-blue)' : 'var(--text-muted)',
            borderLeft: activeSection === section.id ? '2px solid var(--accent-blue)' : '2px solid transparent',
            textDecoration: 'none',
            borderRadius: 0,
            background: 'transparent',
          }}
        >
          {section.label}
        </a>
      ))}
    </nav>
  )
}

const headingFont = "'IBM Plex Serif', serif"
const bodyFont = "'IBM Plex Serif', serif"
const monoFont = "'IBM Plex Mono', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: 'clamp(1.4rem, 4vw, 1.7rem)',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1.5rem',
  letterSpacing: '-0.01em',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '1rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const codeBlockStyle: React.CSSProperties = {
  background: 'transparent',
  border: 'none',
  borderLeft: '2px solid var(--accent-blue)',
  borderRadius: 0,
  padding: '0.5rem 0 0.5rem 1.5rem',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '1rem 0',
}

const tableHeaderStyle: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '0.9rem',
  textTransform: 'uppercase' as const,
  letterSpacing: '0.1em',
  color: 'var(--accent-blue)',
  padding: '10px 14px 10px 0',
  borderBottom: '1px solid var(--accent-blue)',
  textAlign: 'left' as const,
  background: 'transparent',
}

const tableCellStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '1rem',
  color: 'var(--text-secondary)',
  padding: '10px 14px 10px 0',
  borderBottom: '0.5px solid var(--border)',
  background: 'transparent',
}

const h1Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.9rem',
  color: 'var(--text-primary)',
  marginBottom: '1rem',
  letterSpacing: '-0.01em',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
}

const stepLiStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.7,
  marginBottom: '0.75rem',
  paddingLeft: '2.25rem',
  position: 'relative',
}

export default function HowItWorksPage() {
  const [activeSection, setActiveSection] = useState('concept-graph')
  const observerRef = useRef<IntersectionObserver | null>(null)

  useEffect(() => {
    observerRef.current = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setActiveSection(entry.target.id)
          }
        })
      },
      { rootMargin: '-20% 0px -70% 0px' }
    )

    sections.forEach(({ id }) => {
      const el = document.getElementById(id)
      if (el) {
        observerRef.current?.observe(el)
      }
    })

    return () => observerRef.current?.disconnect()
  }, [])

  return (
    <>
      <Header links={[{ label: 'Docs', href: '/docs' }]} />
      <div className="min-h-screen">
      <div className="flex max-w-[960px] mx-auto p-[32px_24px] gap-10 max-md:flex-col max-md:p-4">
        <NavSidebar activeSection={activeSection} />

        <main className="max-w-[720px] flex-1 max-md:mt-20">
          <section id="intro" className="mb-12 pb-8" style={{ borderBottom: '0.5px solid var(--border)' }}>
            <h1 style={h1Style}>
              How Mathua Works
            </h1>
            <p style={bodyStyle}>
              The engine behind the learning: a technical explanation of the concept graph, student
              model, diagnostic algorithm, task selection, and scoring system.
            </p>
          </section>

          {/* Concept Graph */}
          <section id="concept-graph" className="mt-12 pt-6">
            <h2 style={h2Style}>The Concept Graph</h2>
            <p style={bodyStyle}>
              Mathua represents all mathematical knowledge as a directed acyclic graph (a DAG). Each
              node in the graph is an atomic concept: the smallest unit of mathematical knowledge
              that can be practiced and mastered independently. Each directed edge is a prerequisite
              relationship. If concept B has an edge from concept A, then A must be mastered before
              B is ever shown to the student.
            </p>
            <p style={bodyStyle}>
              The graph currently contains 284 concepts spanning 16 domains: from early Counting
              through Calculus, Linear Algebra, and Topology.
            </p>
            <pre style={codeBlockStyle}>
{`{
  "id": "frac.add.diff",
  "label": "Add fractions with different denominators",
  "domain": "fractions",
  "subdomain": "fractions.addition",
  "prerequisites": ["frac.add.same", "arith.factor.lcm"],
  "mastery_threshold": { "streak": 5, "avg_time_seconds": 18.0 }
}`}
            </pre>
            <p style={bodyStyle}>
              The graph is stored as a flat JSON file:{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-blue)' }}>data/concepts.json</code>, and{' '}
              is community-editable. A graph validator runs on every pull request and rejects the
              change if it introduces a cycle.
            </p>
            <div className="flex items-center gap-2 flex-wrap my-6">
              <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>arith.factor.gcf</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>arith.factor.lcm</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>frac.add.diff</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--accent-blue)' }}>frac.mixed.add</span>
            </div>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Student Model */}
          <section id="student-model" className="mt-12 pt-6">
            <h2 style={h2Style}>The Student Model</h2>
            <p style={bodyStyle}>
              Every concept in the graph has a state for each student:
            </p>
            <div className="flex items-center gap-1.5 flex-wrap my-4">
              <span style={{ fontFamily: monoFont, fontSize: '13px', letterSpacing: '0.05em', color: 'var(--text-secondary)' }}>UNSEEN</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont, fontSize: '13px' }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', letterSpacing: '0.05em', color: 'var(--text-secondary)' }}>LEARNING</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont, fontSize: '13px' }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', letterSpacing: '0.05em', color: 'var(--text-secondary)' }}>PRACTICING</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont, fontSize: '13px' }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', letterSpacing: '0.05em', color: 'var(--accent-blue)' }}>MASTERED</span>
              <span style={{ color: 'var(--border-strong)', fontFamily: monoFont, fontSize: '13px' }}>→</span>
              <span style={{ fontFamily: monoFont, fontSize: '13px', letterSpacing: '0.05em', color: 'var(--accent-teal)' }}>DECAYING</span>
            </div>
            <StudentModelFlow />
            <div className="overflow-x-auto my-4">
              <table
                style={{
                  width: '100%',
                  borderCollapse: 'collapse',
                  border: 'none',
                  background: 'transparent',
                }}
              >
                <thead>
                  <tr style={{ background: 'transparent' }}>
                    <th style={tableHeaderStyle}>Concept</th>
                    <th style={tableHeaderStyle}>Streak required</th>
                    <th style={tableHeaderStyle}>Time limit</th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['arith.add.single', '5', '8s'],
                    ['arith.mult.tables', '7', '6s'],
                    ['frac.add.diff', '5', '18s'],
                    ['prealg.eq.one_step_add', '5', '12s'],
                    ['arith.div.long', '5', '20s'],
                  ].map(([concept, streak, time]) => (
                    <tr key={`concept-${concept}`} style={{ background: 'transparent' }}>
                      <td style={tableCellStyle}>{concept}</td>
                      <td style={tableCellStyle}>{streak}</td>
                      <td style={tableCellStyle}>{time}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Spaced Repetition */}
          <section id="spaced-repetition" className="mt-12 pt-6">
            <h2 style={h2Style}>Spaced Repetition</h2>
            <p style={bodyStyle}>
              When a concept reaches MASTERED, Mathua immediately schedules its next review using a
              simplified SM-2 algorithm.
            </p>
            <pre style={codeBlockStyle}>
{`first review:   1 day
second review:  3 days
third review:  interval × ease_factor
ease_factor: starts at 2.5
decreases by 0.2 on each failed review
minimum value: 1.3`}
            </pre>
            <p style={bodyStyle}>
              Reviews are never presented as a separate "review mode." They are woven into every
              session by the scheduler, which manages the 70/30 balance between new material and
              review automatically.
            </p>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Diagnostic */}
          <section id="diagnostic" className="mt-12 pt-6">
            <h2 style={h2Style}>The Diagnostic Algorithm</h2>
            <p style={bodyStyle}>
              When a student first opens Mathua, they enter a Computerised Adaptive Testing (CAT)
              session. The goal is to locate the student's knowledge frontier using as few questions
              as possible.
            </p>
            <ol className="list-none my-4">
              {[
                'The concept graph is sorted topologically. The diagnostic starts at the concept at the midpoint of the sorted order.',
                'If the student answers correctly within the time limit, the algorithm moves forward: it next tests a concept further along the prerequisite chain.',
                'If the student answers incorrectly or exceeds twice the expected time, the algorithm moves backward: it tests a concept earlier in the chain.',
                'This binary search continues with multiple probes per concept until the student\'s knowledge frontier is clearly established.',
                'The diagnostic records a starting mastery estimate for every concept the student passed through. Concepts answered correctly count as LEARNING. Concepts answered quickly and accurately count as conditionally MASTERED and are skipped in early sessions.',
              ].map((step, i) => (
                <li
                  key={step}
                  style={stepLiStyle}
                >
                  <span
                    style={{
                      position: 'absolute',
                      left: 0,
                      fontFamily: monoFont,
                      fontSize: '13px',
                      color: 'var(--accent-blue)',
                    }}
                  >
                    {roman(i + 1)}.
                  </span>
                  {step}
                </li>
              ))}
            </ol>
            <DiagnosticFlow />
            <p style={bodyStyle}>
              The diagnostic takes 20–35 questions for most students. Without this algorithm, a
              naive assessment of 284 concepts would require up to 284 questions. The CAT approach,
              combining binary search with the topological ordering, reduces this by roughly 90%.
            </p>
            <p style={bodyStyle}>
              The diagnostic can be retaken at any time from the settings menu. Retaking does not
              delete progress: it creates a new knowledge estimate that is merged with existing
              data, always preferring the more optimistic estimate so students are never penalised
              for reassessing.
            </p>
            <style>{`ol { counter-reset: step-counter; }`}</style>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Scheduler */}
          <section id="scheduler" className="mt-12 pt-6">
            <h2 style={h2Style}>The Task Selection Algorithm</h2>
            <p style={bodyStyle}>
              The scheduler runs after every question is submitted. It selects the next concept by
              computing a priority score for every eligible concept.
            </p>
            <pre style={codeBlockStyle}>
{`priority = (0.7 × days_since_last_seen)
+ (0.3 × (1.0 − mastery_score))
+ 5.0   if status == DECAYING
+ 2.0   if newly unlocked this session`}
            </pre>
            <p style={bodyStyle}>
              Three hard rules the scheduler enforces regardless of priority scores:
            </p>
            {[
              'A concept whose prerequisites are not all MASTERED is never surfaced.',
              'The same concept is never shown twice in a row.',
              'The scheduler targets a session composition of 70% new and practicing material, 30% review.',
            ].map((rule) => (
              <div
                key={rule}
                style={{
                  fontFamily: bodyFont,
                  fontSize: '0.95rem',
                  color: 'var(--text-secondary)',
                  lineHeight: 1.7,
                  borderLeft: '2px solid var(--accent-blue)',
                  paddingLeft: '1rem',
                  margin: '0.5rem 0',
                }}
              >
              {rule}
            </div>
            ))}
            <SchedulerFlow />
          </section>

          <AsciiDivider pattern="dash" />

          {/* Scoring */}
          <section id="scoring" className="mt-12 pt-6">
            <h2 style={h2Style}>Scoring and Ranking</h2>
            <p style={bodyStyle}>
              Mathua tracks two separate scoring systems:
            </p>
            <pre style={codeBlockStyle}>
{`topic_score = (mastered_concepts × 100)
+ Σ speed_bonus per mastered concept
+ domain_completion_bonus (200 pts if all concepts mastered)
speed_bonus = max(0, (time_limit − avg_time) / time_limit × 50)`}
            </pre>
            <pre style={codeBlockStyle}>
{`weekly_score = (concepts_mastered_this_week × 100)
+ speed_bonus_this_week
+ (current_day_streak × 10)`}
            </pre>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Generators */}
          <section id="generators" className="mt-12 pt-6">
            <h2 style={h2Style}>Generator-Based Problems</h2>
            <p style={bodyStyle}>
              Every problem in Mathua is generated on demand by a parameterised Go function: a
              generator. There is no static question bank.
            </p>
            <pre style={{
              ...codeBlockStyle,
              whiteSpace: 'pre',
              overflowX: 'auto',
            }}>
{`// internal/generator/arithmetic/add.go
type AddSingleGen struct{}

func (g *AddSingleGen) Generate(difficulty float64) generator.Problem {
    max := int(5 + difficulty*4)
    a   := rand.Intn(max) + 2
    b   := rand.Intn(max) + 2
    ans := a + b
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, ans),
    }
}`}
            </pre>
            <p style={bodyStyle}>
              The <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-blue)' }}>difficulty</code> parameter scales operand size from 0.0 to 1.0.
              The scheduler passes a difficulty value based on the student's current mastery score.
            </p>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Polynomial & Expression Grading */}
          <section id="symbolic-grading" className="mt-12 pt-6">
            <h2 style={h2Style}>Expression Grading with SymPy</h2>
            <p style={bodyStyle}>
              From algebra onward, answers are expressions:{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--text-secondary)' }}>x = 4</code>,{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--text-secondary)' }}>(x+2)(x+3)</code>,{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--text-secondary)' }}>2x&#178; + 3x - 5</code>{' '}
              where numeric comparison is no longer sufficient. Mathua uses a mixed Go/Python
              grading system: a Go router dispatches to six grader types, and mathematical
              equivalence for algebra, calculus, differential equations, and trigonometry is
              handled by a Python subprocess running{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-blue)' }}>sympy</code>.
            </p>
            <p style={bodyStyle}>
              The <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-blue)' }}>Router</code> selects the
              grader by <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--text-secondary)' }}>grading_type</code>.
              Concepts with type <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--text-secondary)' }}>polynomial</code> or{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--text-secondary)' }}>expression</code> route to{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-blue)' }}>sympyGrade()</code>, which
              spawns a long-lived Python 3 subprocess. If Python or SymPy are not installed, it
              falls back to a pure-Go string normaliser (symbolic grader).
            </p>
            <GradingFlow />
            <pre style={{
              ...codeBlockStyle,
              whiteSpace: 'pre',
              overflowX: 'auto',
            }}>
{`// internal/grader/router.go
func (r *Router) Grade(t GradingType, expected, answer string) Result {
    switch t {
    case GradingNumeric:
        return r.numeric.grade(expected, answer)
    case GradingPolynomial, GradingExpression:
        return r.sympyGrade(expected, answer)
    case GradingMultipleChoice:
        return r.choice.grade(expected, answer)
    case GradingComparison:
        return r.comparison.grade(expected, answer)
    case GradingOrdering:
        return r.ordering.grade(expected, answer)
    }
}

func (r *Router) sympyGrade(expected, answer string) Result {
    if r.sympy == nil {
        var err error
        r.sympy, err = newSympyGrader()
        if err != nil {
            return r.symbolic.grade(expected, answer) // fallback
        }
    }
    return r.sympy.grade(expected, answer)
}`}
            </pre>
            <p style={bodyStyle}>
              The Go client sends JSON requests to the Python subprocess via stdin. SymPy parses
              both expressions into trees and checks equivalence with{' '}
              <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-blue)' }}>simplify(expected - answer) == 0</code>.
              This catches identities that string or coefficient comparison never could:
            </p>
            <pre style={codeBlockStyle}>
{`# grading/sympy_service.py
import sys, json
from sympy import simplify, parse_expr, Symbol

for line in sys.stdin:
    req = json.loads(line)
    expected = parse_expr(req["expected"])
    answer   = parse_expr(req["answer"])
    correct  = simplify(expected - answer) == 0
    print(json.dumps({"id": req["id"], "correct": correct}))`}
            </pre>
            <p style={bodyStyle}>
              Examples of equivalences SymPy can detect:
              <span style={{ display: 'block', fontFamily: monoFont, fontSize: '13px', marginTop: '0.5rem', color: 'var(--text-muted)' }}>
                x² + 2x + 1 == (x+1)²{' · '}
                sin²(x) + cos²(x) == 1{' · '}
                xe^x − e^x + C == e^x(x−1) + C{' · '}
                2e^(2x) == 2exp(2x)
              </span>
            </p>

          </section>
        </main>
      </div>
    </div>
    </>
  )
}

function roman(n: number): string {
  const r = ['I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X']
  return r[n - 1] ?? String(n)
}
