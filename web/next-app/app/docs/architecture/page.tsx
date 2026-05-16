'use client'

import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'
import MermaidDiagram from '../../../components/MermaidDiagram'

const headingFont = "'IBM Plex Serif', serif"
const bodyFont = "'IBM Plex Serif', serif"
const monoFont = "'IBM Plex Mono', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1rem',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const calloutStyle: React.CSSProperties = {
  borderLeft: '2px solid var(--accent-blue)',
  paddingLeft: '1.5rem',
  fontFamily: bodyFont,
  fontStyle: 'italic',
  fontSize: '0.9rem',
  color: 'var(--text-muted)',
  lineHeight: 1.7,
  marginTop: '1.5rem',
}

export default function ArchitecturePage() {
  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <SectionHeader label="Documentation" title="Architecture" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          Mathua is a single Go binary with two delivery modes. The engine core is identical:
          only the presentation layer differs.
        </p>
        <div style={{ textAlign: 'center', marginBottom: '2rem' }}>
          <a
            href="https://github.com/chuma-beep/mathua/blob/main/docs/architecture.md"
            className="link-underline"
            style={{
              fontFamily: monoFont,
              fontSize: '12px',
              color: 'var(--text-muted)',
            }}
          >
            View source on GitHub
          </a>
        </div>
      </section>

      {/* High-level diagram */}
      <MermaidDiagram code={`graph TD
    subgraph UI[UI Layer]
        Web[Web Browser\\nReact + KaTeX]
        TUI[Desktop\\nBubble Tea]
    end
    subgraph API[API Layer]
        Session[Session Service]
        Graph[Graph Service]
        Leaderboard[Leaderboard Service]
        Diag[Diagnostic Service]
    end
    subgraph Engine[Core Engine]
        Dag[DAG Loader]
        Sched[Scheduler SM-2]
        Gen[Generators]
        Grade[Graders]
        Score[Scoring]
    end
    subgraph Grading[Grading Layer]
        Numeric[Numeric Go]
        Poly[Polynomial]
    end
    subgraph Storage[Storage]
        SQLite[SQLite desktop]
        PG[PostgreSQL web]
    end
    Web --> Session
    Web --> Graph
    Web --> Leaderboard
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
    Gen --> Numeric
    Gen --> Poly
    Score --> SQLite
    Score --> PG`} />
      <style>{`
        :root { --subgraph-bg: #f4f6f8; }
        .dark { --subgraph-bg: #1a1b2e; }
        .subgraph { --_group-fill: var(--subgraph-bg); }
        .subgraph[data-id="UI"] { --_group-hdr: color-mix(in srgb, var(--accent-blue) 10%, var(--subgraph-bg)); }
        .subgraph[data-id="API"] { --_group-hdr: color-mix(in srgb, var(--accent-teal) 10%, var(--subgraph-bg)); }
        .subgraph[data-id="Engine"] { --_group-hdr: color-mix(in srgb, var(--accent-green) 10%, var(--subgraph-bg)); }
        .subgraph[data-id="Grading"] { --_group-hdr: color-mix(in srgb, var(--accent-blue) 10%, var(--subgraph-bg)); }
        .subgraph[data-id="Storage"] { --_group-hdr: color-mix(in srgb, var(--accent-teal) 10%, var(--subgraph-bg)); }
      `}</style>

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
            <div style={{ fontFamily: monoFont, fontSize: '2rem', fontWeight: 400, color: 'var(--accent-blue)', lineHeight: 1 }}>
              {num}
            </div>
            <div style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '0.05em', marginTop: '4px' }}>
              {label}
            </div>
          </div>
        ))}
      </div>

      <AsciiDivider pattern="wave" />

      {/* Layer 1: UI */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>I. UI: Two presentations, one engine</h2>
        <h2 style={h2Style}>II. API: Four services, one router</h2>
        <h2 style={h2Style}>III. Core Engine: Five modules</h2>
        <h2 style={h2Style}>IV. Grading: Numeric and symbolic</h2>
        <h2 style={h2Style}>V. Storage: SQLite and Postgres</h2>
        <h2 style={h2Style}>Diagnostic: Computerised Adaptive Testing</h2>
        <p style={bodyStyle}>
          The diagnostic engine locates a student's knowledge frontier using binary search on
          the topologically sorted concept graph. This reduces the assessment from 284 questions
          (one per concept) to approximately 20–35.
        </p>

        <MermaidDiagram code={`graph TD
    Start((Start)) --> Mid[Test midpoint concept]
    Mid --> Correct{Correct within time limit?}
    Correct -->|Yes| Forward[Move forward: harder concepts]
    Correct -->|No| Backward[Move backward: foundational concepts]
    Forward --> Check{3 correct in a row?}
    Backward --> Check
    Check -->|Yes| Lock([Frontier locked])
    Check -->|No| Mid`} />

        <div style={{ marginTop: '2rem' }}>
          {[
            'The concept graph is sorted topologically. The diagnostic starts at the midpoint.',
            'Correct answers within the time limit move the probe forward toward harder concepts.',
            'Incorrect or slow answers move backward toward foundational material.',
            'After 3 consecutive correct answers in a region, the frontier is considered located.',
            'The diagnostic records a mastery estimate for every concept passed through.',
          ].map((step, i) => (
            <div key={step} style={{ ...bodyStyle, marginBottom: '0.5rem' }}>
              <span style={{ color: 'var(--accent-blue)', fontFamily: monoFont, fontSize: '13px', fontWeight: 500 }}>
                {i + 1}.
              </span>
              {' '}{step}
            </div>
          ))}
        </div>
        <div style={calloutStyle}>
          The diagnostic can be retaken at any time. Retaking does not delete progress: it creates
          a new estimate that is merged with existing data, always preferring the more optimistic
          estimate.
        </div>
      </section>
    </div>
  )
}
