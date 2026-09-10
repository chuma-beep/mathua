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
            Mathua started from a simple observation: most math software lets you
            move on whether or not you actually learned the thing. A student can
            score 70 percent on a chapter test, get pushed into the next chapter,
            and carry the missing 30 percent as a quiet debt that compounds for
            years. Math Academy showed there is another way: never advance until
            the prerequisite is truly mastered. Mathua is an independent, open
            implementation of that idea, built so anyone can run it, inspect it,
            and improve it.
          </p>
        </section>

        <section>
          <h2 style={h2Style}>Frustration</h2>
          <p style={textStyle}>
            Three things were broken enough to build around. First, advancement
            without mastery: curricula move by calendar, not by understanding.
            Second, accuracy without fluency: getting the right answer in 45
            seconds is not the same as knowing it, and fragile knowledge collapses
            under the next topic. Third, forgetting without review: even solid
            knowledge decays, and almost nothing students use schedules reviews
            before that happens. Add memorizable problem banks on top, and you get
            grades that certify exposure instead of ability.
          </p>
        </section>

        <section>
          <h2 style={h2Style}>The bet</h2>
          <p style={textStyle}>
            The bet is that three mechanisms, combined, fix most of that. A
            concept graph that refuses to surface anything whose prerequisites
            you have not mastered. Gating on speed and accuracy together, so
            advancement means fluency, not luck. And spaced repetition woven into
            normal practice, so old concepts resurface before they decay. Each
            one is old science (Bloom, mastery learning, SM-2). The combination,
            enforced without exceptions, is the whole product.
          </p>
        </section>

        <section>
          <h2 style={h2Style}>Invitation</h2>
          <p style={textStyle}>
            Mathua is for self-learners working through math on their own, and
            for parents and teachers who want an honest picture of where a
            student stands. Take the diagnostic to find your frontier, or browse
            the Study library directly. If something is wrong, report it from
            any question: every report goes to a real triage queue.
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
