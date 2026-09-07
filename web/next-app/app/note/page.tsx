import type { Metadata } from 'next'
import Link from 'next/link'
import Header from '../../components/Header'
import Footer from '../../components/Footer'
import BottomTabs from '../../components/BottomTabs'

export const metadata: Metadata = {
  title: "Creator's Note — Mathua",
  description: 'Why Mathua exists, from the person building it.',
}

const monoFont = "'IBM Plex Mono', monospace"
const bodyFont = "'IBM Plex Serif', serif"

const textStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const h2Style: React.CSSProperties = {
  fontFamily: bodyFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  marginTop: '2.5rem',
  marginBottom: '1rem',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
}

export default function NotePage() {
  return (
    <>
      <Header />
      <div className="max-w-[680px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
        <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
          ← Back
        </Link>
        <div style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)', marginTop: '1.5rem', marginBottom: '0.5rem' }}>
          WHY MATHUA EXISTS
        </div>
        <h1 style={{
          fontFamily: "'IBM Plex Serif', serif",
          fontWeight: 400,
          fontSize: '1.9rem',
          color: 'var(--text-primary)',
          marginBottom: '2rem',
          borderBottom: '0.5px solid var(--border)',
          paddingBottom: '0.5rem',
        }}>
          Creator&apos;s Note
        </h1>

        <section>
          <h2 style={h2Style}>Motivation</h2>
          <p style={textStyle}>
            TODO(draft) — why this project started. Replace with the creator&apos;s own words.
          </p>
        </section>

        <section>
          <h2 style={h2Style}>Frustration</h2>
          <p style={textStyle}>
            TODO(draft) — what was broken about learning math that made a new engine necessary.
          </p>
        </section>

        <section>
          <h2 style={h2Style}>The bet</h2>
          <p style={textStyle}>
            TODO(draft) — the core bet: mastery gating, concept graph, spaced repetition. Why these three.
          </p>
        </section>

        <section>
          <h2 style={h2Style}>Invitation</h2>
          <p style={textStyle}>
            TODO(draft) — who this is for, and what to do next.
          </p>
          <div className="flex gap-3 justify-center mt-6 max-sm:flex-col max-sm:items-center">
            <Link href="/onboard" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--bg)', background: 'var(--accent-blue)', padding: '10px 22px', textDecoration: 'none' }}>
              Start diagnostic test →
            </Link>
            <Link href="/study" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-secondary)', border: '0.5px solid var(--border-strong)', padding: '10px 22px', textDecoration: 'none' }}>
              Open Study →
            </Link>
          </div>
        </section>
      </div>
      <Footer />
      <BottomTabs />
    </>
  )
}
