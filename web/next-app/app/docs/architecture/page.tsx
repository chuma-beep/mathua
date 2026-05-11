'use client'

import { useTheme } from '../../../hooks/useTheme'
import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'
import FormulaBlock from '../../../components/FormulaBlock'
import MermaidDiagram from '../../../components/MermaidDiagram'

const headingFont = "'Cormorant Garamond', serif"
const bodyFont = "'Inter', -apple-system, sans-serif"
const monoFont = "'JetBrains Mono', 'Fira Code', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1rem',
}

const h3Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1rem',
  color: 'var(--text-primary)',
  marginTop: '1.5rem',
  marginBottom: '0.5rem',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const codeBlockStyle: React.CSSProperties = {
  background: 'var(--surface)',
  border: 'none',
  padding: '20px 24px',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '1rem 0',
}

export default function ArchitecturePage() {
  const { theme } = useTheme()

  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <SectionHeader label="Documentation" title="Architecture" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          Mathua is a single Go binary with two delivery modes. The engine core is identical —
          only the presentation layer differs.
        </p>
      </section>

      {/* High-level diagram */}
      <MermaidDiagram
        theme={theme}
        chart={`flowchart TB
    subgraph UI["UI Layer"]
        Web["Web Browser\\nReact + Next.js"]
        TUI["Desktop TUI\\nBubble Tea"]
    end
    subgraph API["API Layer"]
        Session["Session Service"]
        Graph["Graph Service"]
        Leader["Leaderboard Service"]
        Diag["Diagnostic Service"]
    end
    subgraph Engine["Core Engine"]
        Dag["DAG Loader"]
        Sched["Scheduler SM-2"]
        Gen["Generators"]
        Grade["Graders"]
        Score["Scoring"]
    end
    subgraph Storage["Storage"]
        SQLite[("SQLite\\nDesktop")]
        PG[("PostgreSQL\\nWeb")]
    end
    Web --> Session
    Web --> Graph
    Web --> Leader
    Web --> Diag
    TUI --> Dag
    TUI --> Sched
    TUI --> Gen
    Session --> Dag
    Session --> Sched
    Graph --> Dag
    Diag --> Sched
    Dag --> Sched
    Sched --> Gen
    Gen --> Grade
    Grade --> Score
    Score --> SQLite
    Score --> PG`}
      />

      {/* Metric cards */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
        gap: '16px',
        marginBottom: '48px',
        marginTop: '24px',
      }}>
        {[
          { num: '1', label: 'Binary' },
          { num: '2', label: 'UI Modes' },
          { num: '5', label: 'Engine Modules' },
          { num: '2', label: 'Graders' },
        ].map(({ num, label }) => (
          <div
            key={label}
            style={{
              background: 'var(--surface)',
              border: '0.5px solid var(--border)',
              padding: '20px',
              textAlign: 'center',
            }}
          >
            <div style={{ fontFamily: monoFont, fontSize: '2rem', fontWeight: 400, color: 'var(--accent-gold)', lineHeight: 1 }}>
              {num}
            </div>
            <div style={{ fontFamily: monoFont, fontSize: '11px', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '0.05em', marginTop: '4px' }}>
              {label}
            </div>
          </div>
        ))}
      </div>

      <AsciiDivider pattern="wave" />

      {/* Layer 1: UI */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>I. UI — Two presentations, one engine</h2>
        <p style={bodyStyle}>
          Mathua ships with two user interfaces that share the same engine core. The desktop
          interface requires no account, no network, and stores everything in SQLite. The web
          interface uses React with Next.js and connects to a shared Postgres database.
        </p>

        <div style={{ overflowX: 'auto', margin: '1.5rem 0' }}>
          <table style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}>
            <thead>
              <tr style={{ background: 'transparent' }}>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>
                  Property
                </th>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>
                  Desktop TUI
                </th>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>
                  Web App
                </th>
              </tr>
            </thead>
            <tbody>
              {[
                ['Framework', 'Bubble Tea (Go)', 'React + Next.js'],
                ['Account', 'Not required', 'Required'],
                ['Network', 'Fully offline', 'Requires server'],
                ['Storage', 'SQLite', 'PostgreSQL'],
                ['Leaderboard', 'Local only', 'Global weekly'],
                ['Graph View', 'ASCII tree', 'Interactive 3D'],
              ].map(([prop, tui, web], i) => (
                <tr key={i} style={{ background: 'transparent' }}>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>
                    {prop}
                  </td>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>
                    {tui}
                  </td>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>
                    {web}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Layer 2: API */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>II. API — Four services, one router</h2>
        <p style={bodyStyle}>
          The web server exposes a REST API through Go's standard <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>net/http</code> package.
          No external HTTP framework is used. The API is organized into four services:
        </p>
        <div style={{ ...bodyStyle, marginBottom: '0.3rem' }}>
          <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
          <strong style={{ color: 'var(--text-primary)', fontWeight: 400, fontFamily: monoFont }}>Session Service</strong> — Starts and manages practice sessions. Returns the next problem based on scheduler output.
        </div>
        <div style={{ ...bodyStyle, marginBottom: '0.3rem' }}>
          <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
          <strong style={{ color: 'var(--text-primary)', fontWeight: 400, fontFamily: monoFont }}>Graph Service</strong> — Exposes the concept DAG for visualization. Returns node metadata and prerequisite edges.
        </div>
        <div style={{ ...bodyStyle, marginBottom: '0.3rem' }}>
          <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
          <strong style={{ color: 'var(--text-primary)', fontWeight: 400, fontFamily: monoFont }}>Leaderboard Service</strong> — Computes and returns the weekly leaderboard. Scores reset every Monday at 00:00 UTC.
        </div>
        <div style={{ ...bodyStyle }}>
          <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
          <strong style={{ color: 'var(--text-primary)', fontWeight: 400, fontFamily: monoFont }}>Diagnostic Service</strong> — Runs the Computerised Adaptive Testing engine to locate a student's knowledge frontier.
        </div>
        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
          marginTop: '1.5rem',
        }}>
          The TUI bypasses the API layer entirely — it calls the engine modules directly through Go
          function interfaces. The API layer exists solely for the web client.
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Layer 3: Core Engine */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>III. Core Engine — Five modules</h2>

        <h3 style={h3Style}>DAG Loader</h3>
        <p style={bodyStyle}>
          Loads <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>concepts.json</code> and builds an in-memory directed acyclic graph.
          Validates that the graph contains no cycles and no orphaned prerequisites. Runs on every
          server start and on every PR through CI. The in-memory representation supports efficient
          topological sorting and prerequisite lookups.
        </p>

        <h3 style={h3Style}>Scheduler (SM-2)</h3>
        <p style={bodyStyle}>
          The scheduler selects the next concept after every submitted answer. It computes a priority
          score for every eligible concept — those whose prerequisites are all mastered. Three hard
          rules apply: prerequisites must be mastered, no repeats, and a 70/30 balance between new
          material and review.
        </p>
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <FormulaBlock
            code={`priority = (0.7 × days_since_last_seen)
+ (0.3 × (1.0 − mastery_score))
+ 5.0  if status == DECAYING
+ 2.0  if newly unlocked`}
          />
        </div>

        <h3 style={h3Style}>Mastery & Spaced Repetition</h3>
        <p style={bodyStyle}>
          Each concept tracks a mastery score based on streak and response time. When mastered,
          a concept enters a spaced repetition schedule using a simplified SM-2 algorithm. Reviews
          are woven into normal sessions — there is no separate review mode.
        </p>
        <pre style={codeBlockStyle}>
{`first review:   1 day
second review:  3 days
third review:   interval × 2.5
interval minimum: 1
ease factor:     2.5, min 1.3`}
        </pre>

        <h3 style={h3Style}>Generators</h3>
        <p style={bodyStyle}>
          Every problem is generated on demand by a parameterised Go function — a generator.
          There is no static question bank. Each generator accepts a difficulty parameter (0.0–1.0)
          and returns a problem struct with a question, answer, and explanation.
        </p>
        <pre style={codeBlockStyle}>
{`type Generator interface {
    Generate(difficulty float64) Problem
}

type Problem struct {
    Question    string
    Answer      string
    Explanation string
    Difficulty  float64
}`}
        </pre>

        <h3 style={h3Style}>Scoring & Leaderboard</h3>
        <p style={bodyStyle}>
          Mathua tracks two scores. The lifetime topic score reflects permanent mastery — it never
          decreases. The weekly score resets every Monday and is used for the global leaderboard.
          Both incorporate concept counts, speed bonuses, and streak multipliers.
        </p>
        <div style={{ display: 'flex', justifyContent: 'center' }}>
          <FormulaBlock
            code={`topic_score  = (mastered_concepts × 100)
              + Σ speed_bonus
weekly_score = (mastered_this_week × 100)
              + speed_bonus
              + (current_streak × 10)`}
          />
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Layer 4: Grading */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>IV. Grading — Numeric and symbolic</h2>
        <p style={bodyStyle}>
          Mathua uses two grading strategies, selected by concept type. Arithmetic concepts
          use the numeric grader — simple string or integer comparison. Algebra concepts use
          a pure-Go polynomial grader with no external dependencies.
        </p>
        <div style={{ overflowX: 'auto', margin: '1.5rem 0' }}>
          <table style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}>
            <thead>
              <tr style={{ background: 'transparent' }}>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>Grader</th>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>Used for</th>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>Strategy</th>
              </tr>
            </thead>
            <tbody>
              {[
                ['Numeric', 'Counting, arithmetic, fractions', 'String / integer equality'],
                ['Polynomial', 'Algebra (future)', 'Parse → expand → normalise → compare coefficients'],
              ].map(([grader, used, strategy], i) => (
                <tr key={i} style={{ background: 'transparent' }}>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>{grader}</td>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>{used}</td>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>{strategy}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Layer 5: Storage */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>V. Storage — SQLite and Postgres</h2>
        <p style={bodyStyle}>
          The desktop TUI uses SQLite for fully offline operation — no installation, no
          configuration, no network. The web server uses PostgreSQL for concurrent access
          and global state. The data schema is identical across both databases, and the
          engine core uses a repository interface that abstracts the storage backend.
        </p>
        <div style={{ overflowX: 'auto', margin: '1.5rem 0' }}>
          <table style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}>
            <thead>
              <tr style={{ background: 'transparent' }}>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>Entity</th>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>Stored in</th>
                <th style={{
                  fontFamily: monoFont,
                  fontWeight: 400,
                  fontSize: '0.85rem',
                  textTransform: 'uppercase',
                  letterSpacing: '0.08em',
                  color: 'var(--accent-gold)',
                  padding: '10px 14px 10px 0',
                  borderBottom: '1px solid var(--accent-gold)',
                  textAlign: 'left',
                }}>Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                ['Concept graph', 'concepts.json', 'Community-editable, validated on PR'],
                ['Student progress', 'SQLite / Postgres', 'Per-concept state, mastery, intervals'],
                ['Session history', 'SQLite / Postgres', 'Every problem, answer, and timing'],
                ['Weekly scores', 'Postgres only', 'Global leaderboard, reset every Monday'],
              ].map(([entity, stored, desc], i) => (
                <tr key={i} style={{ background: 'transparent' }}>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>{entity}</td>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>{stored}</td>
                  <td style={{
                    fontFamily: bodyFont,
                    fontSize: '0.9rem',
                    color: 'var(--text-secondary)',
                    padding: '10px 14px 10px 0',
                    borderBottom: '0.5px solid var(--border)',
                  }}>{desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* CAT Engine */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Diagnostic — Computerised Adaptive Testing</h2>
        <p style={bodyStyle}>
          The diagnostic engine locates a student's knowledge frontier using binary search on
          the topologically sorted concept graph. This reduces the assessment from 284 questions
          (one per concept) to approximately 20–35.
        </p>

        <MermaidDiagram
          theme={theme}
          chart={`flowchart LR
    Start(["Start: Topologically Sorted Graph"]) --> Midpoint{"Test midpoint concept"}
    Midpoint -->|"Correct + Fast"| Forward["Move forward\\nharder concepts"]
    Midpoint -->|"Incorrect / Slow"| Backward["Move backward\\nfoundational concepts"]
    Forward --> Check{"3 consecutive\\ncorrect?"}
    Backward --> Check
    Check -->|"Yes"| Locked(["FRONTIER LOCKED\\nRecord mastery estimates"])
    Check -->|"No"| Midpoint`}
        />

        <div style={{ marginTop: '2rem' }}>
          {[
            'The concept graph is sorted topologically. The diagnostic starts at the midpoint.',
            'Correct answers within the time limit move the probe forward toward harder concepts.',
            'Incorrect or slow answers move backward toward foundational material.',
            'After 3 consecutive correct answers in a region, the frontier is considered located.',
            'The diagnostic records a mastery estimate for every concept passed through.',
          ].map((step, i) => (
            <div key={i} style={{ ...bodyStyle, marginBottom: '0.5rem' }}>
              <span style={{ color: 'var(--accent-gold)', fontFamily: monoFont, fontSize: '13px', fontWeight: 500 }}>
                {i + 1}.
              </span>
              {' '}{step}
            </div>
          ))}
        </div>
        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
          marginTop: '1.5rem',
        }}>
          The diagnostic can be retaken at any time. Retaking does not delete progress — it creates
          a new estimate that is merged with existing data, always preferring the more optimistic
          estimate.
        </div>
      </section>
    </div>
  )
}
