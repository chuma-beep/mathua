'use client'

import { useTheme } from '../../../hooks/useTheme'
import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'
import D2Diagram from '../../../components/D2Diagram'

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

export default function ArchitecturePage() {
  const { theme } = useTheme()

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
      <D2Diagram name="architecture" theme={theme} />

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

        <D2Diagram name="cat-diagnostic" theme={theme} />

        <div style={{ marginTop: '2rem' }}>
          {[
            'The concept graph is sorted topologically. The diagnostic starts at the midpoint.',
            'Correct answers within the time limit move the probe forward toward harder concepts.',
            'Incorrect or slow answers move backward toward foundational material.',
            'After 3 consecutive correct answers in a region, the frontier is considered located.',
            'The diagnostic records a mastery estimate for every concept passed through.',
          ].map((step, i) => (
            <div key={`cat-step-${i}`} style={{ ...bodyStyle, marginBottom: '0.5rem' }}>
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
          The diagnostic can be retaken at any time. Retaking does not delete progress: it creates
          a new estimate that is merged with existing data, always preferring the more optimistic
          estimate.
        </div>
      </section>
    </div>
  )
}
