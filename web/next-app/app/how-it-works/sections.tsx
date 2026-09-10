'use client'

import dynamic from 'next/dynamic'
import {
  bodyStyle,
  codeBlockStyle,
  h2Style,
  inlineCodeStyle,
  monoFont,
  mutedCodeStyle,
  roman,
  stepLiStyle,
  tableCellStyle,
  tableHeaderStyle,
  bodyFont,
} from './styles'

const StudentModelFlow = dynamic(() => import('../../components/StudentModelFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 320, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: monoFont, fontSize: '13px' }}>Loading&hellip;</div>,
})

const DiagnosticFlow = dynamic(() => import('../../components/DiagnosticFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 420, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: monoFont, fontSize: '13px' }}>Loading&hellip;</div>,
})

const SchedulerFlow = dynamic(() => import('../../components/SchedulerFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 320, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: monoFont, fontSize: '13px' }}>Loading&hellip;</div>,
})

const GradingFlow = dynamic(() => import('../../components/GradingFlow'), {
  ssr: false,
  loading: () => <div style={{ height: 300, border: '0.5px solid var(--border)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontFamily: monoFont, fontSize: '13px' }}>Loading&hellip;</div>,
})

export function ConceptGraphSection() {
  return (
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
        The graph currently contains 630 concepts spanning 17 domains: from early Counting
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
        <code style={inlineCodeStyle}>data/concepts.json</code>, and{' '}
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
  )
}

export function StudentModelSection() {
  const rows: [string, string, string][] = [
    ['Single-digit addition', '5', '8 seconds'],
    ['Multiplication tables', '7', '6 seconds'],
    ['Add fractions with different denominators', '5', '18 seconds'],
    ['Solve one-step addition equations', '5', '12 seconds'],
    ['Long division', '5', '20 seconds'],
  ]
  return (
    <section id="student-model" className="mt-12 pt-6">
      <h2 style={h2Style}>The Student Model</h2>
      <p style={bodyStyle}>
        Every concept in the graph has a state for each student:
      </p>

      <StudentModelFlow />

      {/* Mobile: stacked labeled cards */}
      <div className="sm:hidden space-y-2 w-full max-w-[720px] mx-auto min-w-0 my-6">
        {rows.map(([concept, streak, time]) => (
          <div key={`card-${concept}`} className="border-[0.5px] border-mathua-border p-3 min-w-0 overflow-hidden text-center">
            <div className="font-mono text-[12px] text-mathua-primary mb-1 break-words">{concept}</div>
            <div className="flex flex-wrap justify-center gap-x-4 gap-y-0.5 font-mono text-[11px] text-mathua-muted">
              <span>{streak} correct in a row</span>
              <span>{time}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Desktop: real table */}
      <div className="hidden sm:block overflow-x-auto overscroll-x-contain w-[calc(100%+2rem)] -mx-4 px-4 sm:w-full sm:mx-0 sm:px-0 my-6">
        <table className="w-full"
          style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}
        >
          <thead>
            <tr style={{ background: 'transparent' }}>
              <th style={tableHeaderStyle}>Skill</th>
              <th style={tableHeaderStyle}>Correct answers in a row needed</th>
              <th style={tableHeaderStyle}>Time per question</th>
            </tr>
          </thead>
          <tbody>
            {rows.map(([concept, streak, time]) => (
              <tr key={`row-${concept}`} style={{ background: 'transparent' }}>
                <td style={tableCellStyle}>{concept}</td>
                <td style={tableCellStyle}>{streak}</td>
                <td style={tableCellStyle}>{time}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}

export function SpacedRepetitionSection() {
  return (
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
        Reviews are never presented as a separate &quot;review mode.&quot; They are woven into every
        session by the scheduler, which manages the 70/30 balance between new material and
        review automatically.
      </p>
    </section>
  )
}

export function DiagnosticSection() {
  const steps = [
    'The concept graph is sorted topologically. The diagnostic starts at the concept at the midpoint of the sorted order.',
    'If the student answers correctly within the time limit, the algorithm moves forward: it next tests a concept further along the prerequisite chain.',
    'If the student answers incorrectly or exceeds twice the expected time, the algorithm moves backward: it tests a concept earlier in the chain.',
    'This binary search continues with multiple probes per concept until the student\'s knowledge frontier is clearly established.',
    'The diagnostic records a starting mastery estimate for every concept the student passed through. Concepts answered correctly count as LEARNING. Concepts answered quickly and accurately count as conditionally MASTERED and are skipped in early sessions.',
  ]
  return (
    <section id="diagnostic" className="mt-12 pt-6">
      <h2 style={h2Style}>The Diagnostic Algorithm</h2>
      <p style={bodyStyle}>
        When a student first opens Mathua, they enter a Computerised Adaptive Testing (CAT)
        session. The goal is to locate the student&apos;s knowledge frontier using as few questions
        as possible.
      </p>
      <ol className="list-none my-4">
        {steps.map((step, i) => (
          <li key={step} style={stepLiStyle}>
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
        The diagnostic test asks as few questions as possible. Without this algorithm, a
        naive assessment of 630 concepts would require up to 630 questions. The CAT approach,
        combining binary search with the topological ordering, reduces this by roughly 90%.
      </p>
      <p style={bodyStyle}>
        The diagnostic test can be retaken at any time from your profile. Retaking does not
        delete progress: it creates a new knowledge estimate that is merged with existing
        data, always preferring the more optimistic estimate so students are never penalised
        for reassessing.
      </p>
      <style>{`ol { counter-reset: step-counter; }`}</style>
    </section>
  )
}

export function SchedulerSection() {
  const rules = [
    'A concept whose prerequisites are not all MASTERED is never surfaced.',
    'The same concept is never shown twice in a row.',
    'The scheduler targets a session composition of 70% new and practicing material, 30% review.',
  ]
  return (
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
      {rules.map((rule) => (
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
  )
}

export function ScoringSection() {
  return (
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
  )
}

export function GeneratorsSection() {
  return (
    <section id="generators" className="mt-12 pt-6">
      <h2 style={h2Style}>Generator-Based Problems</h2>
      <p style={bodyStyle}>
        Every problem in Mathua is generated on demand by a parameterised Go function: a
        generator. There is no static question bank.
      </p>
      <pre style={{ ...codeBlockStyle, whiteSpace: 'pre', overflowX: 'auto' }}>
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
        The <code style={inlineCodeStyle}>difficulty</code> parameter scales operand size from 0.0 to 1.0.
        The scheduler passes a difficulty value based on the student&apos;s current mastery score.
      </p>
    </section>
  )
}

export function SymbolicGradingSection() {
  return (
    <section id="symbolic-grading" className="mt-12 pt-6">
      <h2 style={h2Style}>Expression Grading with SymPy</h2>
      <p style={bodyStyle}>
        From algebra onward, answers are expressions:{' '}
        <code style={mutedCodeStyle}>x = 4</code>,{' '}
        <code style={mutedCodeStyle}>(x+2)(x+3)</code>,{' '}
        <code style={mutedCodeStyle}>2x&#178; + 3x - 5</code>{' '}
        where numeric comparison is no longer sufficient. Mathua uses a mixed Go/Python
        grading system: a Go router dispatches to six grader types, and mathematical
        equivalence for algebra, calculus, differential equations, and trigonometry is
        handled by a Python subprocess running{' '}
        <code style={inlineCodeStyle}>sympy</code>.
      </p>
      <p style={bodyStyle}>
        The <code style={inlineCodeStyle}>Router</code> selects the
        grader by <code style={mutedCodeStyle}>grading_type</code>.
        Concepts with type <code style={mutedCodeStyle}>polynomial</code> or{' '}
        <code style={mutedCodeStyle}>expression</code> route to{' '}
        <code style={inlineCodeStyle}>sympyGrade()</code>, which
        spawns a long-lived Python 3 subprocess. If Python or SymPy are not installed, it
        falls back to a pure-Go string normaliser (symbolic grader).
      </p>
      <GradingFlow />
      <pre style={{ ...codeBlockStyle, whiteSpace: 'pre', overflowX: 'auto' }}>
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
        <code style={inlineCodeStyle}>simplify(expected - answer) == 0</code>.
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
  )
}
