'use client'

import dynamic from 'next/dynamic'
import Link from 'next/link'
import { useTheme } from '../hooks/useTheme'
import Header from '../components/Header'
import AsciiDivider from '../components/AsciiDivider'
import SectionHeader from '../components/SectionHeader'
import Pipeline from '../components/Pipeline'
import FormulaBlock from '../components/FormulaBlock'
import DomainTable from '../components/DomainTable'
import ProgressionLevels from '../components/ProgressionLevels'
import D2Diagram from '../components/D2Diagram'
import Footer from '../components/Footer'
import conceptsData from '../data/concepts.json'

const loadingGraphStyle: React.CSSProperties = {
  height: 'clamp(320px, 50vh, 520px)',
  background: 'transparent',
  borderRadius: 0,
  border: '0.5px solid var(--border)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'var(--text-muted)',
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '13px',
}

const statCountStyle: React.CSSProperties = {
  color: 'var(--bg)',
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '30px',
  fontWeight: 400,
  lineHeight: 1,
  marginBottom: '4px',
}

const MathConceptGraph3D = dynamic(() => import('../components/MathConceptGraph3D'), {
  ssr: false,
  loading: () => (
    <div style={loadingGraphStyle}>
      Loading graph…
    </div>
  ),
})

const PIPELINE_STATES = [
  { label: 'UNSEEN', status: 'unseen' as const },
  { label: 'LEARNING', status: 'learning' as const },
  { label: 'PRACTICING', status: 'practicing' as const },
  { label: 'MASTERED', status: 'mastered' as const },
  { label: 'DECAYING', status: 'decaying' as const },
]

export default function HomePage() {
  const { theme, mounted } = useTheme()

  const conceptCount = conceptsData.length
  const connectionCount = conceptsData.reduce(
    (sum: number, c: any) => sum + ((c.prerequisites as any[])?.length || 0),
    0
  )
  const domainCount = new Set(conceptsData.map((c: any) => c.domain)).size

  const domainOrder = [
    'counting',
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
  const domainLabels: Record<string, string> = {
    counting: 'Counting',
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
  }
  const domainCounts = conceptsData.reduce(
    (acc: Record<string, number>, c: any) => {
      acc[c.domain] = (acc[c.domain] || 0) + 1
      return acc
    },
    {} as Record<string, number>
  )

  const levels = [
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

  if (!mounted) {
    return <div style={{ background: 'var(--bg)', minHeight: '100vh' }} />
  }

  const headingFont = "'IBM Plex Serif', serif"
  const bodyFont = "'IBM Plex Serif', serif"
  const monoFont = "'IBM Plex Mono', monospace"

  const ctaPrimaryStyle: React.CSSProperties = {
    background: 'var(--accent-blue)',
    color: 'var(--bg)',
    fontFamily: monoFont,
    fontSize: '13px',
    padding: '10px 22px',
    border: 'none',
    borderRadius: '2px',
    letterSpacing: '0.04em',
    textTransform: 'none',
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    textDecoration: 'none',
    transition: 'background 0.2s',
  }

  const ctaSecondaryStyle: React.CSSProperties = {
    background: 'transparent',
    color: 'var(--text-secondary)',
    fontFamily: monoFont,
    fontSize: '13px',
    padding: '10px 22px',
    border: '0.5px solid var(--border-strong)',
    borderRadius: '2px',
    letterSpacing: '0.04em',
    textTransform: 'none',
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    textDecoration: 'none',
    transition: 'background 0.2s, color 0.2s',
  }

  const statsRowStyle: React.CSSProperties = {
    fontFamily: monoFont,
    fontSize: 'clamp(11px, 3vw, 12px)',
    letterSpacing: '0.04em',
    color: 'var(--text-muted)',
    marginBottom: '2.5rem',
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'center',
    gap: '2px 8px',
  }

  const codeQuoteStyle: React.CSSProperties = {
    background: 'transparent',
    border: 'none',
    borderLeft: '2px solid var(--accent-blue)',
    borderRadius: 0,
    padding: '0.5rem 0 0.5rem 1rem',
    fontFamily: monoFont,
    fontSize: 'clamp(11px, 2.5vw, 13px)',
    color: 'var(--text-secondary)',
    whiteSpace: 'pre',
    overflowX: 'auto',
    textAlign: 'left',
    lineHeight: 1.6,
    display: 'inline-block',
  }

  return (
    <>
      <Header links={[{ label: 'How it works', href: '/how-it-works' }, { label: 'Docs', href: '/docs' }, { label: 'Leaderboard', href: '/leaderboard' }, { label: 'Login', href: '/login' }]} />
      <div className="max-w-container mx-auto px-6 max-sm:px-4">
      {/* ── Hero ── */}
      <section
        className="py-20 max-sm:py-12 text-center"
        style={{ background: 'transparent' }}
      >
        <h1
          style={{
            fontFamily: headingFont,
            fontWeight: 400,
            fontSize: 'clamp(2.2rem, 5vw, 3.8rem)',
            lineHeight: 1.2,
            color: 'var(--text-primary)',
            marginBottom: '1rem',
            letterSpacing: '-0.01em',
          }}
        >
          Master the foundation.
          <br />
          Earn the abstraction.
        </h1>
        <p
          style={{
            fontFamily: bodyFont,
            fontSize: 'clamp(0.95rem, 2.5vw, 1.1rem)',
            color: 'var(--text-secondary)',
            maxWidth: '600px',
            margin: '0 auto 2rem',
            lineHeight: 1.85,
          }}
        >
          Mathua is an open-source adaptive math learning engine. It never lets you advance until
          you have truly mastered the prerequisite, both speed and accuracy must be proven.
        </p>

        <div className="flex gap-3 justify-center items-center mb-10 max-sm:flex-col max-sm:[&_a]:w-full max-sm:[&_a]:max-w-[280px]">
          <Link href="/login" style={ctaPrimaryStyle}>
            Open the web app
          </Link>
          <a href="https://github.com/chuma-beep/mathua" style={ctaSecondaryStyle}>
            View on GitHub
          </a>
        </div>

        <div style={statsRowStyle}>
          <span>{conceptCount} topics</span>
          <span style={{ color: 'var(--border-strong)' }}>·</span>
          <span>{connectionCount} connections</span>
          <span style={{ color: 'var(--border-strong)' }}>·</span>
          <span>{domainCount} domains</span>
          <span style={{ color: 'var(--border-strong)' }}>·</span>
          <span>web + desktop</span>
          <span style={{ color: 'var(--border-strong)' }}>·</span>
          <span>open source</span>
        </div>

        {/* 3D Concept Graph */}
          <div
          className="max-sm:p-4"
          style={{
            border: '0.5px solid var(--border)',
            borderRadius: 0,
            padding: '1.5rem',
            background: 'var(--graph-surface)',
          }}
        >
          <MathConceptGraph3D
            theme={theme}
            concepts={conceptsData.map((c: any) => ({
              id: c.id,
              label: c.label,
              domain: c.domain,
              prerequisites: c.prerequisites,
            }))}
          />
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* ── How It Works ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader title="A different kind of math system" />

        <div style={{ marginTop: '2rem' }}>
          <div style={{ borderTop: '0.5px solid var(--border)', padding: '1.2rem 0' }}>
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.3rem' }}>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--accent-blue)' }}>I.</span>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--text-primary)' }}>Mastery gating</span>
            </div>
            <p style={{ fontFamily: bodyFont, fontSize: '1rem', color: 'var(--text-secondary)', lineHeight: 1.85, marginLeft: 'clamp(0.5rem, 3vw, 1.5rem)' }}>
              You cannot advance until your streak and response time both meet the threshold. Knowing
              the answer is not enough; you must know it fast.
            </p>
          </div>

          <div style={{ borderTop: '0.5px solid var(--border)', padding: '1.2rem 0' }}>
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.3rem' }}>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--accent-blue)' }}>II.</span>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--text-primary)' }}>Concept graph</span>
            </div>
            <p style={{ fontFamily: bodyFont, fontSize: '1rem', color: 'var(--text-secondary)', lineHeight: 1.85, marginLeft: 'clamp(0.5rem, 3vw, 1.5rem)' }}>
              Every concept is a node with explicit prerequisites. The scheduler reads the graph and
              your progress to decide what you see next.
            </p>
          </div>

          <div style={{ borderTop: '0.5px solid var(--border)', padding: '1.2rem 0' }}>
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.3rem' }}>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--accent-blue)' }}>III.</span>
              <span style={{ fontFamily: headingFont, fontSize: '1.1rem', color: 'var(--text-primary)' }}>Spaced repetition</span>
            </div>
            <p style={{ fontFamily: bodyFont, fontSize: '1rem', color: 'var(--text-secondary)', lineHeight: 1.85, marginLeft: 'clamp(0.5rem, 3vw, 1.5rem)' }}>
              Concepts you master resurface automatically before they decay. Reviews are woven into
              your session; there is no separate review mode.
            </p>
          </div>
        </div>
      </section>

      {/* ── Mastery Pipeline ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader title="Five states. One direction." />
        <Pipeline states={PIPELINE_STATES} className="my-7" />
        <div className="flex justify-center">
          <FormulaBlock
            code={`priority = (0.7 × days_since_last_seen)
         + (0.3 × (1 − mastery))
         + 5.0  if DECAYING
         + 2.0  if newly_unlocked`}
          />
        </div>
        <p
          style={{
            fontFamily: bodyFont,
            fontSize: '0.95rem',
            color: 'var(--text-secondary)',
            textAlign: 'center',
            marginTop: '1.25rem',
            lineHeight: 1.85,
          }}
        >
          The scheduler enforces three hard rules: prerequisites must be mastered before a concept
          unlocks, the same concept never appears twice in a row, and roughly 70% of each session is
          new material.
        </p>
      </section>

      <AsciiDivider pattern="dash" />

      {/* ── Curriculum ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader title="What Mathua covers" />
        <DomainTable
          className="mt-5"
          rows={domainOrder.map((d) => ({
            domain: domainLabels[d],
            count: domainCounts[d],
          }))}
        />
        <p
          style={{
            fontFamily: bodyFont,
            fontStyle: 'italic',
            fontSize: '13px',
            color: 'var(--text-muted)',
            marginTop: '0.75rem',
          }}
        >
          Problems are generated on demand, never stored. There is nothing to memorise.
        </p>
      </section>

      <AsciiDivider pattern="wave" />

      {/* ── Progression ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader title="Rank by mastery. Level by depth." />
        <div className="flex flex-wrap gap-8 items-start mt-7 max-sm:flex-col max-sm:gap-4">
          <div className="flex-1 min-w-[280px] max-sm:min-w-0">
            <p
              style={{
                fontFamily: bodyFont,
                fontSize: '0.95rem',
                color: 'var(--text-secondary)',
                lineHeight: 1.85,
                textAlign: 'center',
              }}
            >
              The leaderboard resets every Monday at 00:00 UTC. Your score is calculated from three
              components:
            </p>
            <div className="flex justify-center mt-4">
              <FormulaBlock
                code={`score = (mastered_count × 100)
    + speed_bonus
     + (current_streak × 10)`}
              />
            </div>
          </div>
          <div className="flex-1 min-w-[280px] max-sm:min-w-0">
            <ProgressionLevels levels={levels} />
          </div>
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* ── Platforms ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader title="One engine. Two ways to run it." />

        <div className="flex flex-wrap gap-0 items-start mt-7 max-sm:flex-col" style={{ fontFamily: bodyFont }}>
          {/* Web column */}
          <div className="flex-1 min-w-[280px] max-sm:pr-0" style={{ paddingRight: '2rem' }}>
            <h3 style={{ fontFamily: headingFont, fontSize: '1.1rem', fontWeight: 400, color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
              Web (browser)
            </h3>
            <div style={{ borderBottom: '0.5px solid var(--border)', width: '100%', marginBottom: '0.75rem' }} />
            {[
              'React + KaTeX',
              'Account required',
              'Global leaderboard',
              'Graph view',
            ].map((feature) => (
              <div key={feature} style={{ fontSize: '0.95rem', color: 'var(--text-secondary)', lineHeight: 1.8 }}>
                <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem' }}>·</span>
                {feature}
              </div>
            ))}
          </div>

          {/* Divider */}
          <div className="max-sm:hidden" style={{ borderLeft: '0.5px solid var(--border)', alignSelf: 'stretch', minHeight: '160px' }} />

          {/* Desktop column */}
          <div className="flex-1 min-w-[280px] max-sm:pl-0" style={{ paddingLeft: '2rem' }}>
            <h3 style={{ fontFamily: headingFont, fontSize: '1.1rem', fontWeight: 400, color: 'var(--text-primary)', marginBottom: '0.75rem' }}>
              Desktop (terminal)
            </h3>
            <div style={{ borderBottom: '0.5px solid var(--border)', width: '100%', marginBottom: '0.75rem' }} />
            {[
              'Bubble Tea TUI',
              'No account needed',
              'Fully offline',
              'SQLite storage',
            ].map((feature) => (
              <div key={feature} style={{ fontSize: '0.95rem', color: 'var(--text-secondary)', lineHeight: 1.8 }}>
                <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem' }}>·</span>
                {feature}
              </div>
            ))}
          </div>
        </div>

        <div className="flex justify-center mt-8">
          <D2Diagram name="platforms" theme={theme} />
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* ── Contributing ── */}
      <section className="py-20 max-sm:py-12">
        <SectionHeader title="Built to be extended." />
        <p
          style={{
            fontFamily: bodyFont,
            fontSize: '0.95rem',
            color: 'var(--text-secondary)',
            lineHeight: 1.85,
            textAlign: 'center',
            maxWidth: '600px',
            margin: '0 auto',
          }}
        >
          Every concept is a JSON node. Every problem is a Go generator function. Every contribution
          goes through a graph validator that rejects cycles and orphaned nodes automatically.
        </p>
        <div className="flex justify-center mt-4">
          <pre style={codeQuoteStyle}>
{`{
   "id":                "arith.add.multi",
   "label":             "Multi-digit addition",
   "domain":            "arithmetic",
   "prerequisites":     ["arith.add.single", "arith.add.carry"],
   "mastery_threshold": {
     "streak":          5,
     "avg_time_seconds": 8
   }
 }`}
          </pre>
        </div>
        <div className="flex justify-center mt-6">
          <a
            href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md"
            className="link-underline"
            style={{
              fontFamily: headingFont,
              fontStyle: 'italic',
              fontSize: '1rem',
              color: 'var(--accent-blue)',
            }}
          >
            Read CONTRIBUTING.md →
          </a>
        </div>
      </section>

      <Footer />
    </div>
    </>
  )
}
