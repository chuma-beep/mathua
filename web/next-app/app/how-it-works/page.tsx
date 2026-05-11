'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import AsciiDivider from '../../components/AsciiDivider'

const sections = [
  { id: 'concept-graph', label: 'Concept graph' },
  { id: 'student-model', label: 'Student model' },
  { id: 'spaced-repetition', label: 'Spaced repetition' },
  { id: 'diagnostic', label: 'Diagnostic' },
  { id: 'scheduler', label: 'Scheduler' },
  { id: 'scoring', label: 'Scoring' },
  { id: 'generators', label: 'Generators' },
  { id: 'symbolic-grading', label: 'Polynomial grading' },
]

function NavSidebar({ activeSection }: { activeSection: string }) {
  return (
    <nav className="w-[220px] flex-shrink-0 sticky top-[100px] h-fit max-md:fixed max-md:top-[57px] max-md:left-0 max-md:right-0 max-md:w-full max-md:bg-[var(--bg)] max-md:border-b max-md:border-mathua-border max-md:flex max-md:overflow-x-auto max-md:p-[8px_16px] max-md:gap-2 max-md:z-[99]">
      {sections.map((section) => (
        <a
          key={section.id}
          href={`#${section.id}`}
          className={`block text-[13px] py-2 px-3 border-l-2 mb-1 transition-colors max-md:border-l-0 max-md:border-b-2 max-md:mb-0 max-md:whitespace-nowrap ${
            activeSection === section.id
              ? 'text-mathua-blue border-l-mathua-blue max-md:border-b-mathua-blue'
              : 'text-mathua-muted border-l-transparent max-md:border-b-transparent hover:text-mathua-secondary'
          }`}
        >
          {section.label}
        </a>
      ))}
    </nav>
  )
}

export default function HowItWorksPage() {
  const [activeSection, setActiveSection] = useState('concept-graph')
  const [theme, setTheme] = useState<'dark' | 'light'>('dark')
  const [mounted, setMounted] = useState(false)
  const observerRef = useRef<IntersectionObserver | null>(null)

  useEffect(() => {
    setMounted(true)
    const saved = localStorage.getItem('mathua-theme')
    if (saved === 'light' || saved === 'dark') {
      setTheme(saved)
    }
  }, [])

  useEffect(() => {
    if (!mounted) return
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [theme, mounted])

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setTheme(next)
    localStorage.setItem('mathua-theme', next)
  }

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
    <div className="min-h-screen bg-mathua-bg">
      <button
        onClick={toggleTheme}
        className="fixed top-[50px] right-5 z-[1000] bg-mathua-surface border border-mathua-border-strong text-mathua-primary px-3.5 py-2 rounded-md font-mono text-xs cursor-pointer transition-all duration-200 hover:border-mathua-blue hover:text-mathua-blue"
        aria-label="Toggle theme"
      >
        {theme === 'dark' ? '\u2600' : '\u263E'}
      </button>

      <header className="flex justify-between items-center p-[16px_24px] border-b border-mathua-border sticky top-0 bg-mathua-bg z-[100]">
        <Link href="/" className="text-lg font-semibold text-mathua-blue">
          Mathua
        </Link>
        <nav className="flex gap-4">
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            Home
          </Link>
          <Link href="/" className="text-mathua-secondary text-sm hover:text-mathua-primary">
            Web App
          </Link>
        </nav>
      </header>

      <div className="flex max-w-[960px] mx-auto p-[32px_24px] gap-10 max-md:flex-col max-md:p-4">
        <NavSidebar activeSection={activeSection} />

        <main className="max-w-[720px] flex-1 max-md:mt-20">
          <section id="intro" className="mb-12 pb-8 border-b border-mathua-border">
            <h1 className="font-serif text-[1.6rem] font-semibold text-mathua-blue mb-4">
              How Mathua Works
            </h1>
            <p className="text-mathua-secondary text-base leading-relaxed">
              The engine behind the learning — a technical explanation of the concept graph, student
              model, diagnostic algorithm, task selection, and scoring system.
            </p>
          </section>

          {/* Concept Graph */}
          <section id="concept-graph" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              The Concept Graph
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              Mathua represents all mathematical knowledge as a directed acyclic graph — a DAG. Each
              node in the graph is an atomic concept: the smallest unit of mathematical knowledge
              that can be practiced and mastered independently. Each directed edge is a prerequisite
              relationship. If concept B has an edge from concept A, then A must be mastered before
              B is ever shown to the student.
            </p>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The graph currently contains 284 concepts spanning 16 domains — from early Counting
              through Calculus, Linear Algebra, and Topology.
            </p>
            <pre className="bg-mathua-code border border-mathua-border rounded-md p-5 font-mono text-[13px] text-mathua-primary whitespace-pre overflow-x-auto my-4">
{`{
  "id": "frac.add.diff",
  "label": "Add fractions with different denominators",
  "domain": "fractions",
  "subdomain": "fractions.addition",
  "prerequisites": ["frac.add.same", "arith.factor.lcm"],
  "mastery_threshold": { "streak": 5, "avg_time_seconds": 18.0 }
}`}
            </pre>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The graph is stored as a flat JSON file — <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">data/concepts.json</code> — and
              is community-editable. A graph validator runs on every pull request and rejects the
              change if it introduces a cycle.
            </p>
            <div className="flex items-center gap-3 flex-wrap my-6">
              <span className="bg-mathua-surface border border-mathua-border text-mathua-blue font-mono text-xs px-3 py-1 rounded-full">arith.factor.gcf</span>
              <span className="text-mathua-border-strong">→</span>
              <span className="bg-mathua-surface border border-mathua-border text-mathua-blue font-mono text-xs px-3 py-1 rounded-full">arith.factor.lcm</span>
              <span className="text-mathua-border-strong">→</span>
              <span className="bg-mathua-surface border border-mathua-border text-mathua-blue font-mono text-xs px-3 py-1 rounded-full">frac.add.diff</span>
              <span className="text-mathua-border-strong">→</span>
              <span className="bg-mathua-surface border border-mathua-border text-mathua-blue font-mono text-xs px-3 py-1 rounded-full">frac.mixed.add</span>
            </div>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Student Model */}
          <section id="student-model" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              The Student Model
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              Every concept in the graph has a state for each student:
            </p>
            <div className="flex items-center gap-2 flex-wrap my-4">
              <span className="bg-mathua-surface border border-mathua-border text-mathua-muted font-mono text-[10px] px-2.5 py-1 rounded">UNSEEN</span>
              <span className="text-mathua-muted text-xs">→</span>
              <span className="bg-mathua-surface border border-mathua-border text-mathua-muted font-mono text-[10px] px-2.5 py-1 rounded">LEARNING</span>
              <span className="text-mathua-muted text-xs">→</span>
              <span className="bg-mathua-surface border border-mathua-border text-mathua-muted font-mono text-[10px] px-2.5 py-1 rounded">PRACTICING</span>
              <span className="text-mathua-muted text-xs">→</span>
              <span className="bg-mathua-green text-white border border-mathua-green font-mono text-[10px] px-2.5 py-1 rounded">MASTERED</span>
              <span className="text-mathua-muted text-xs">→</span>
              <span className="bg-mathua-gold text-[var(--bg)] border border-mathua-gold font-mono text-[10px] px-2.5 py-1 rounded">DECAYING</span>
            </div>
            <div className="overflow-x-auto my-4">
              <table className="bg-mathua-surface border border-mathua-border border-collapse w-full text-[13px]">
                <thead>
                  <tr>
                    <th className="bg-mathua-surface-elevated text-mathua-blue font-mono text-[11px] uppercase p-[10px_14px] text-left border-b border-mathua-border">
                      Concept
                    </th>
                    <th className="bg-mathua-surface-elevated text-mathua-blue font-mono text-[11px] uppercase p-[10px_14px] text-left border-b border-mathua-border">
                      Streak required
                    </th>
                    <th className="bg-mathua-surface-elevated text-mathua-blue font-mono text-[11px] uppercase p-[10px_14px] text-left border-b border-mathua-border">
                      Time limit
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {[
                    ['arith.add.single', '5', '8s'],
                    ['arith.mult.tables', '7', '6s'],
                    ['frac.add.diff', '5', '18s'],
                    ['prealg.eq.one_step_add', '5', '12s'],
                    ['arith.div.long', '5', '20s'],
                  ].map(([concept, streak, time], i) => (
                    <tr key={i} className={i % 2 === 0 ? 'bg-mathua-surface-elevated' : ''}>
                      <td className="text-mathua-secondary p-[10px_14px] border-b border-mathua-border">
                        {concept}
                      </td>
                      <td className="text-mathua-secondary p-[10px_14px] border-b border-mathua-border">
                        {streak}
                      </td>
                      <td className="text-mathua-secondary p-[10px_14px] border-b border-mathua-border">
                        {time}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Spaced Repetition */}
          <section id="spaced-repetition" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              Spaced Repetition
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              When a concept reaches MASTERED, Mathua immediately schedules its next review using a
              simplified SM-2 algorithm.
            </p>
            <pre className="bg-mathua-code border border-mathua-border border-l-[2px] border-l-mathua-blue rounded-r-md p-5 font-mono text-[13px] text-mathua-secondary whitespace-pre leading-relaxed my-4">
{`first review:   1 day
second review:  3 days
third review:  interval × ease_factor
ease_factor: starts at 2.5
decreases by 0.2 on each failed review
minimum value: 1.3`}
            </pre>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              Reviews are never presented as a separate "review mode." They are woven into every
              session by the scheduler, which manages the 70/30 balance between new material and
              review automatically.
            </p>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Diagnostic */}
          <section id="diagnostic" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              The Diagnostic Algorithm
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              When a student first opens Mathua, they enter a Computerised Adaptive Testing (CAT)
              session. The goal is to locate the student's knowledge frontier using as few questions
              as possible.
            </p>
            <ol className="list-none my-4">
              {[
                'The concept graph is sorted topologically. The diagnostic starts at the concept at the midpoint of the sorted order.',
                'If the student answers correctly within the time limit, the algorithm moves forward — it next tests a concept further along the prerequisite chain.',
                'If the student answers incorrectly or exceeds twice the expected time, the algorithm moves backward — it tests a concept earlier in the chain.',
                'This binary search continues until three consecutive correct answers are recorded in a region, or three consecutive failures.',
                'The diagnostic records a starting mastery estimate for every concept the student passed through. Concepts answered correctly count as LEARNING. Concepts answered quickly and accurately count as conditionally MASTERED and are skipped in early sessions.',
              ].map((step, i) => (
                <li
                  key={i}
                  className="text-mathua-secondary text-[0.95rem] leading-[1.7] mb-3 pl-9 relative before:content-[counter(step)] before:absolute before:left-0 before:text-mathua-blue before:font-mono before:text-[13px]"
                  style={{ counterIncrement: 'step-counter 1' }}
                >
                  {step}
                </li>
              ))}
            </ol>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The diagnostic takes 20–35 questions for most students. Without this algorithm, a
              naive assessment of 284 concepts would require up to 284 questions. The CAT approach,
              combining binary search with the topological ordering, reduces this by roughly 90%.
            </p>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The diagnostic can be retaken at any time from the settings menu. Retaking does not
              delete progress — it creates a new knowledge estimate that is merged with existing
              data, always preferring the more optimistic estimate so students are never penalised
              for reassessing.
            </p>
            <style>{`ol { counter-reset: step-counter; }`}</style>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Scheduler */}
          <section id="scheduler" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              The Task Selection Algorithm
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The scheduler runs after every question is submitted. It selects the next concept by
              computing a priority score for every eligible concept.
            </p>
            <pre className="bg-mathua-code border border-mathua-border border-l-[2px] border-l-mathua-blue rounded-r-md p-5 font-mono text-[13px] text-mathua-secondary whitespace-pre leading-relaxed my-4">
{`priority = (0.7 × days_since_last_seen)
+ (0.3 × (1.0 − mastery_score))
+ 5.0   if status == DECAYING
+ 2.0   if newly unlocked this session`}
            </pre>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              Three hard rules the scheduler enforces regardless of priority scores:
            </p>
            {[
              'A concept whose prerequisites are not all MASTERED is never surfaced.',
              'The same concept is never shown twice in a row.',
              'The scheduler targets a session composition of 70% new and practicing material, 30% review.',
            ].map((rule, i) => (
              <div
                key={i}
                className="border-l-2 border-mathua-border-strong pl-4 my-2 text-[0.95rem] text-mathua-secondary"
              >
                {rule}
              </div>
            ))}
          </section>

          <AsciiDivider pattern="dash" />

          {/* Scoring */}
          <section id="scoring" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              Scoring and Ranking
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              Mathua tracks two separate scoring systems:
            </p>
            <pre className="bg-mathua-code border border-mathua-border border-l-[2px] border-l-mathua-blue rounded-r-md p-5 font-mono text-[13px] text-mathua-secondary whitespace-pre leading-relaxed my-4">
{`topic_score = (mastered_concepts × 100)
+ Σ speed_bonus per mastered concept
+ domain_completion_bonus (200 pts if all concepts mastered)
speed_bonus = max(0, (time_limit − avg_time) / time_limit × 50)`}
            </pre>
            <pre className="bg-mathua-code border border-mathua-border border-l-[2px] border-l-mathua-blue rounded-r-md p-5 font-mono text-[13px] text-mathua-secondary whitespace-pre leading-relaxed my-4">
{`weekly_score = (concepts_mastered_this_week × 100)
+ speed_bonus_this_week
+ (current_day_streak × 10)`}
            </pre>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Generators */}
          <section id="generators" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              Generator-Based Problems
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              Every problem in Mathua is generated on demand by a parameterised Go function — a
              generator. There is no static question bank.
            </p>
            <pre className="bg-mathua-code border border-mathua-border rounded-md p-5 font-mono text-[13px] text-mathua-primary whitespace-pre overflow-x-auto my-4">
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
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">difficulty</code> parameter scales operand size from 0.0 to 1.0.
              The scheduler passes a difficulty value based on the student's current mastery score.
            </p>
          </section>

          <AsciiDivider pattern="dash" />

          {/* Polynomial & Expression Grading */}
          <section id="symbolic-grading" className="mt-12 pt-6">
            <h2 className="font-serif text-[1.6rem] font-medium text-mathua-blue border-b border-mathua-border pb-2 mb-6">
              Polynomial &amp; Expression Grading
            </h2>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              From algebra onward, answers are expressions —{' '}
              <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">x = 4</code>,{' '}
              <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">(x+2)(x+3)</code>,{' '}
              <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">2x&#178; + 3x - 5</code>{' '}
              — and numeric comparison is no longer sufficient. Mathua uses a pure-Go polynomial
              grader with no external dependencies: no Python, no SymPy, no subprocess.
            </p>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The grading router selects the grader by concept type. Arithmetic concepts go to the
              numeric grader. Algebra concepts go to the polynomial grader, which parses the
              student's input into a coefficient representation, normalises it, and compares against
              the expected answer.
            </p>
            <div className="flex items-center gap-2 flex-wrap my-4 text-xs">
              <span className="text-mathua-secondary">Input</span>
              <span className="text-mathua-muted">→</span>
              <span className="text-mathua-secondary">Tokenizer</span>
              <span className="text-mathua-muted">→</span>
              <span className="text-mathua-secondary">Parser (AST)</span>
              <span className="text-mathua-muted">→</span>
              <span className="text-mathua-secondary">Expander / Normaliser</span>
              <span className="text-mathua-muted">→</span>
              <span className="bg-mathua-green text-white px-2 py-0.5 rounded font-mono text-[10px]">Match</span>
            </div>
            <pre className="bg-mathua-code border border-mathua-border rounded-md p-5 font-mono text-[13px] text-mathua-primary whitespace-pre overflow-x-auto my-4">
{`// internal/grader/polynomial.go
type PolyGrader struct{}

func (g *PolyGrader) Grade(input, expected string) (bool, error) {
    userPoly   := parsePolynomial(input)    // e.g. (x+2)(x+3) → coeffs
    expectPoly := parsePolynomial(expected) // e.g. x^2+5x+6 → coeffs

    expand(&userPoly)   // multiply out factored form
    normalise(&userPoly)
    normalise(&expectPoly)

    return userPoly.Equals(expectPoly), nil
}`}
            </pre>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              The polynomial grader handles factoring (parse and expand), simplification
              (normalise), and equation-solving (isolate variable, compare). The same pure-Go
              pipeline runs identically in the web server and the desktop TUI — no Python,
              no environment dependencies, no disabled features.
            </p>
            <p className="text-mathua-secondary text-base leading-[1.85] mb-4">
              For display, the TUI renders exponents using <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">^</code> notation
              (<code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">x^2 + 5x + 6</code>)
              while the web frontend uses KaTeX with{' '}
              <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">{'{'}^{'}'}{'{'}^{'}'}</code>{' '}
              for proper superscripts. Every generator receives a render mode flag
              (<code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">RenderTUI</code> or{' '}
              <code className="bg-mathua-surface-elevated px-1.5 py-0.5 rounded font-mono text-[0.9em] text-mathua-blue">RenderWeb</code>) so
              output is always correct for the target display.
            </p>
          </section>
        </main>
      </div>
    </div>
  )
}
