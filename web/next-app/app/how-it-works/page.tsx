'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import './globals.css'

const sections = [
  { id: 'concept-graph', label: 'Concept graph' },
  { id: 'student-model', label: 'Student model' },
  { id: 'spaced-repetition', label: 'Spaced repetition' },
  { id: 'diagnostic', label: 'Diagnostic' },
  { id: 'scheduler', label: 'Scheduler' },
  { id: 'scoring', label: 'Scoring' },
  { id: 'generators', label: 'Generators' },
  { id: 'symbolic-grading', label: 'Symbolic grading' },
]

function NavSidebar({ activeSection }: { activeSection: string }) {
  return (
    <nav className="nav-sidebar">
      {sections.map(section => (
        <a
          key={section.id}
          href={`#${section.id}`}
          className={`nav-link ${activeSection === section.id ? 'active' : ''}`}
        >
          {section.label}
        </a>
      ))}
    </nav>
  )
}

export default function HowItWorksPage() {
  const [activeSection, setActiveSection] = useState('concept-graph')
  const observerRef = useRef<IntersectionObserver | null>(null)
  const sectionRefs = useRef<Map<string, HTMLElement>>(new Map())

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
        sectionRefs.current.set(id, el)
        observerRef.current?.observe(el)
      }
    })

    return () => observerRef.current?.disconnect()
  }, [])

  return (
    <div className="how-page">
      <header className="how-header">
        <Link href="/" className="how-logo">Mathua</Link>
        <nav className="how-nav">
          <Link href="/" className="how-nav-link">Home</Link>
          <Link href="/" className="how-nav-link">Web App</Link>
        </nav>
      </header>

      <div className="how-layout">
        <NavSidebar activeSection={activeSection} />
        
        <main className="how-content">
          <section id="intro" className="how-intro">
            <h1>How Mathua Works</h1>
            <p className="how-subtitle">
              The engine behind the learning — a technical explanation of the concept graph, student model, diagnostic algorithm, task selection, and scoring system.
            </p>
          </section>

          <section id="concept-graph" className="how-section">
            <h2>The Concept Graph</h2>
            <p>
              Mathua represents all mathematical knowledge as a directed acyclic graph — a DAG. Each node in the graph is an atomic concept: the smallest unit of mathematical knowledge that can be practiced and mastered independently. Each directed edge is a prerequisite relationship. If concept B has an edge from concept A, then A must be mastered before B is ever shown to the student.
            </p>
            <p>
              The graph currently contains 81 concepts spanning five domains: Counting, Arithmetic, Fractions, Pre-Algebra, and Algebra (v1.1). Every concept has a unique identifier in dot-notation that encodes its domain and sub-domain.
            </p>
            <pre className="code-block">{`{
  "id": "frac.add.diff",
  "label": "Add fractions with different denominators",
  "domain": "fractions",
  "subdomain": "fractions.addition",
  "prerequisites": ["frac.add.same", "arith.factor.lcm"],
  "mastery_threshold": { "streak": 5, "avg_time_seconds": 18.0 }
}`}</pre>
            <p>
              The graph is stored as a flat JSON file — <code>data/concepts.json</code> — and is community-editable. A graph validator runs on every pull request and rejects the change if it introduces a cycle (a loop in the prerequisite chain), an orphaned concept (a node with no path to any root concept), or a dangling reference (a prerequisite ID that does not exist in the graph). This keeps the graph structurally sound regardless of how many contributors add to it.
            </p>
            <p>
              The graph is loaded at startup using Kahn's algorithm — a topological sort that processes concepts in dependency order. The sorted order is used by the scheduler to determine which concepts are currently available to a given student.
            </p>
            <div className="concept-chain">
              <span className="concept-pill">arith.factor.gcf</span>
              <span className="arrow">→</span>
              <span className="concept-pill">arith.factor.lcm</span>
              <span className="arrow">→</span>
              <span className="concept-pill">frac.add.diff</span>
              <span className="arrow">→</span>
              <span className="concept-pill">frac.mixed.add</span>
            </div>
          </section>

          <section id="student-model" className="how-section">
            <h2>The Student Model</h2>
            <p>
              Every concept in the graph has a state for each student. A concept can be in one of five states:
            </p>
            <div className="state-flow">
              <span className="state-pill unseen">UNSEEN</span>
              <span className="state-arrow">→</span>
              <span className="state-pill learning">LEARNING</span>
              <span className="state-arrow">→</span>
              <span className="state-pill practicing">PRACTICING</span>
              <span className="state-arrow">→</span>
              <span className="state-pill mastered">MASTERED</span>
              <span className="state-arrow">→</span>
              <span className="state-pill decaying">DECAYING</span>
            </div>
            <p>
              <strong>UNSEEN</strong> means the concept's prerequisites have not all been mastered yet. The concept is locked and will not appear in any session.
            </p>
            <p>
              <strong>LEARNING</strong> means the concept is unlocked — all prerequisites are mastered — and the student has begun attempting it but has not yet built a streak.
            </p>
            <p>
              <strong>PRACTICING</strong> means the student has answered correctly at least twice in a row and is building toward the mastery threshold.
            </p>
            <p>
              <strong>MASTERED</strong> means the student has met both the streak threshold and the time threshold simultaneously. The concept is considered known.
            </p>
            <p>
              <strong>DECAYING</strong> means the concept was mastered but has not been seen in more than 14 days. The spaced repetition algorithm has flagged it for review. Decaying concepts are treated as high priority by the scheduler.
            </p>
            <div className="threshold-table-wrap">
              <table className="threshold-table">
                <thead>
                  <tr>
                    <th>Concept</th>
                    <th>Streak required</th>
                    <th>Time limit</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>arith.add.single</td>
                    <td>5</td>
                    <td>8s</td>
                  </tr>
                  <tr>
                    <td>arith.mult.tables</td>
                    <td>7</td>
                    <td>6s</td>
                  </tr>
                  <tr>
                    <td>frac.add.diff</td>
                    <td>5</td>
                    <td>18s</td>
                  </tr>
                  <tr>
                    <td>prealg.eq.one_step_add</td>
                    <td>5</td>
                    <td>12s</td>
                  </tr>
                  <tr>
                    <td>arith.div.long</td>
                    <td>5</td>
                    <td>20s</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p>
              Both thresholds must be met simultaneously. A student who consistently gets correct answers but takes twice the expected time has not truly mastered the concept — they are still computing, not yet fluent.
            </p>
          </section>

          <section id="spaced-repetition" className="how-section">
            <h2>Spaced Repetition</h2>
            <p>
              When a concept reaches MASTERED, Mathua immediately schedules its next review using a simplified SM-2 algorithm. The review interval grows with each successful review:
            </p>
            <div className="formula-block">
              <span>first review:   1 day</span>
              <span>second review:  3 days</span>
              <span>third review:  interval × ease_factor</span>
              <span className="formula-label">ease_factor: starts at 2.5</span>
              <span className="formula-label">decreases by 0.2 on each failed review</span>
              <span className="formula-label">minimum value: 1.3</span>
            </div>
            <p>
              The ease factor is personal to each student on each concept. A concept that a student consistently answers quickly and correctly will have a high ease factor — reviews become infrequent because the knowledge is consolidating. A concept the student struggles with will have a lower ease factor — reviews remain frequent until the knowledge is solid.
            </p>
            <p>
              When a concept's next review date passes without the student seeing it, it transitions from MASTERED to DECAYING. The scheduler treats DECAYING concepts with a priority bonus of +5.0 in the priority formula — higher than any new concept — because allowing decay to proceed too long forces the student backward through the mastery process.
            </p>
            <p>
              Reviews are never presented as a separate "review mode." They are woven into every session by the scheduler, which manages the 70/30 balance between new material and review automatically.
            </p>
          </section>

          <section id="diagnostic" className="how-section">
            <h2>The Diagnostic Algorithm</h2>
            <p>
              When a student first opens Mathua, or chooses to retake the diagnostic at any time, they enter a Computerised Adaptive Testing (CAT) session. The goal of the diagnostic is to locate the student's knowledge frontier — the boundary between concepts they have mastered and concepts they are not yet ready for — using as few questions as possible.
            </p>
            <p>
              The algorithm works as follows:
            </p>
            <ol className="step-list">
              <li>The concept graph is sorted topologically. The diagnostic starts at the concept at the midpoint of the sorted order.</li>
              <li>If the student answers correctly within the time limit, the algorithm moves forward — it next tests a concept further along the prerequisite chain.</li>
              <li>If the student answers incorrectly or exceeds twice the expected time, the algorithm moves backward — it tests a concept earlier in the chain.</li>
              <li>This binary search continues until three consecutive correct answers are recorded in a region, or three consecutive failures. At that point the frontier is considered located.</li>
              <li>The diagnostic records a starting mastery estimate for every concept the student passed through. Concepts answered correctly count as LEARNING. Concepts answered correctly under time count as PRACTICING. Concepts answered quickly and accurately count as conditionally MASTERED and are skipped in early sessions.</li>
            </ol>
            <p>
              The diagnostic takes between 15 and 25 questions for most students. Without this algorithm, a naive assessment of 81 concepts would require up to 81 questions. The CAT approach reduces this by roughly 70%.
            </p>
            <p>
              Explain that the diagnostic can be retaken at any time from the settings menu. Retaking does not delete progress — it creates a new knowledge estimate that is merged with existing progress data, always preferring the more optimistic estimate so that students are never penalised for retaking.
            </p>
          </section>

          <section id="scheduler" className="how-section">
            <h2>The Task Selection Algorithm</h2>
            <p>
              The scheduler runs after every question is submitted. It selects the next concept to present by computing a priority score for every concept that is currently in LEARNING, PRACTICING, or DECAYING state, plus any newly unlocked concepts in UNSEEN state whose prerequisites just reached MASTERED.
            </p>
            <div className="formula-block">
              <span>priority = (0.7 × days_since_last_seen)</span>
              <span>+ (0.3 × (1.0 − mastery_score))</span>
              <span>+ 5.0   if status == DECAYING</span>
              <span>+ 2.0   if newly unlocked this session</span>
            </div>
            <p>
              <strong>Days since last seen</strong> weighted at 0.7 is the dominant term. Concepts the student has not seen recently float to the top. This implements the spacing effect automatically — the longer since a concept was practiced, the more urgent it becomes.
            </p>
            <p>
              <strong>Weakness score</strong> weighted at 0.3 adjusts for accuracy. A concept with a mastery score of 0.4 gets a weakness contribution of 0.18, while a mastered concept at 1.0 contributes 0. This means weak concepts are gently prioritised over strong ones when recency is equal.
            </p>
            <p>
              <strong>The DECAYING bonus</strong> of 5.0 is large enough to override everything else. A concept that is overdue for review will always be shown before new material. This prevents silent forgetting.
            </p>
            <p>
              <strong>The unlock bonus</strong> of 2.0 gives newly unlocked concepts a small head-start so they appear promptly after their prerequisite is mastered, rather than waiting for natural decay to elevate them.
            </p>
            <p>
              Three hard rules the scheduler enforces regardless of priority scores:
            </p>
            <div className="rule-block">
              A concept whose prerequisites are not all MASTERED is never surfaced, even if it somehow receives a high priority score. The prerequisite graph is the final authority.
            </div>
            <div className="rule-block">
              The same concept is never shown twice in a row. After any concept is presented, its priority score is set to zero for the next selection cycle.
            </div>
            <div className="rule-block">
              The scheduler targets a session composition of 70% new and practicing material, 30% review. If the review queue is empty, 100% new material is shown. If the student has many overdue reviews, the ratio can temporarily invert, but the scheduler always returns toward 70/30 as reviews are cleared.
            </div>
          </section>

          <section id="scoring" className="how-section">
            <h2>Scoring and Ranking</h2>
            <p>
              Mathua tracks two separate scoring systems that serve different purposes.
            </p>
            <p>
              <strong>Topic scores</strong> are personal and never reset. They reflect the depth and quality of a student's mastery within each domain and sub-domain. Topic scores are computed as:
            </p>
            <div className="formula-block">
              <span>topic_score = (mastered_concepts × 100)</span>
              <span>+ Σ speed_bonus per mastered concept</span>
              <span>+ domain_completion_bonus (200 pts if all concepts mastered)</span>
              <span className="formula-label">speed_bonus per concept = max(0, (time_limit − avg_time) / time_limit × 50)</span>
            </div>
            <p>
              A student who masters a concept in exactly the time limit earns zero speed bonus. A student who masters it in half the time earns 25 bonus points. Speed bonuses are capped at 50 points per concept and can never produce a negative score — a student who is slow but correct still earns their base 100 points.
            </p>
            <p>
              <strong>Weekly rank scores</strong> reset every Monday at 00:00 UTC. They reflect a student's activity and performance in the current week only. The formula:
            </p>
            <div className="formula-block">
              <span>weekly_score = (concepts_mastered_this_week × 100)</span>
              <span>+ speed_bonus_this_week</span>
              <span>+ (current_day_streak × 10)</span>
            </div>
            <p>
              The leaderboard shows global rank and level-tier rank simultaneously. Level is permanent and based on all-time mastered concept count. The nine levels from Newcomer to Math Architect reflect depth of knowledge and never decrease regardless of weekly performance.
            </p>
            <p>
              The purpose of separating these two systems is intentional: weekly scores create short-term motivation and competition without threatening the permanent record of what a student has learned. A bad week does not undo months of work.
            </p>
          </section>

          <section id="generators" className="how-section">
            <h2>Generator-Based Problems</h2>
            <p>
              Every problem in Mathua is generated on demand by a parameterised Go function — a generator. There is no static question bank. This has two consequences:
            </p>
            <p>
              The student cannot memorise answers. Each session produces different numbers, different coefficients, different framing. The only way to score well is to understand the underlying concept.
            </p>
            <p>
              The system is infinitely scalable. Adding a new concept requires writing one generator function and one fuzz test — not authoring hundreds of individual questions.
            </p>
            <pre className="code-block">{`// internal/generator/arithmetic/add.go
type AddSingleGen struct{}

func (g *AddSingleGen) Generate(difficulty float64) generator.Problem {
    max := int(5 + difficulty*4)
    a   := rand.Intn(max) + 2   // never 0 or 1
    b   := rand.Intn(max) + 2
    ans := a + b
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, ans),
    }
}`}</pre>
            <p>
              The <code>difficulty</code> parameter: it scales operand size from 0.0 (smallest valid inputs) to 1.0 (largest inputs the concept supports). The scheduler passes a difficulty value based on the student's current mastery score — students who are solidly practicing receive harder variations, newly unlocked concepts start at low difficulty.
            </p>
            <p>
              Every generator has a fuzz test that runs 10,000 samples and asserts: no panics, no empty question strings, no unparseable answer strings, no trivially obvious problems (such as <code>0 + n</code> or <code>1 × n</code>), and that the difficulty parameter visibly changes the output distribution. A generator that fails this test cannot be merged.
            </p>
          </section>

          <section id="symbolic-grading" className="how-section">
            <h2>Symbolic Grading — v1.1</h2>
            <p>
              Arithmetic concepts use numeric grading — the student's input is parsed as a float64 and compared to the expected answer within a tolerance of 1×10⁻⁶. This is fast, dependency-free, and covers all concepts through pre-algebra.
            </p>
            <p>
              From algebra onward, answers are expressions — <code>x = 4</code>, <code>2x + 1</code>, <code>(x+2)(x+3)</code> — and numeric comparison is no longer sufficient. Two expressions that look different can be algebraically identical. Mathua v1.1 introduces symbolic grading via a Python subprocess running SymPy.
            </p>
            <div className="flow-diagram">
              <span>Go grader</span>
              <span className="flow-arrow">→</span>
              <span>JSON line over stdin</span>
              <span className="flow-arrow">→</span>
              <span>sympy_service.py</span>
              <span className="flow-arrow">→</span>
              <span>JSON response over stdout</span>
              <span className="flow-arrow">→</span>
              <span>Go grader</span>
            </div>
            <p>
              The SymPy service is launched once at startup as a long-running subprocess. For each symbolic grading request, the Go grader writes a single JSON line to the subprocess's stdin and reads the response from stdout. The service uses <code>sympy.simplify(user - expected) == 0</code> to test algebraic equivalence.
            </p>
            <pre className="code-block">{`# grading/sympy_service.py
for line in sys.stdin:
    req = json.loads(line)
    try:
        user = sympify(req['user_answer'])
        exp  = sympify(req['expected'])
        ok   = simplify(user - exp) == 0
        print(json.dumps({'correct': bool(ok)}))
    except Exception as e:
        print(json.dumps({'correct': False, 'error': str(e)}))
    sys.stdout.flush()`}</pre>
            <p>
              Symbolic grading requires Python 3.10+ and SymPy (<code>pip install sympy</code>). The desktop TUI detects the absence of Python at startup and gracefully disables symbolic concepts, showing a one-line install prompt. All arithmetic and pre-algebra concepts remain available without Python.
            </p>
          </section>
        </main>
      </div>
    </div>
  )
}